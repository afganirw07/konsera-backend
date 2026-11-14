package checkout

import (
	"context"
	"database/sql"
	"errors"
	"konsera-backend/internal/models"
	"time"

	"github.com/google/uuid"
)

const MidtransPaymentMethodID = "f47ac10b-58cc-4372-a567-0e02b2c3d479"

type PaymentNotification struct {
	OrderID       string
	TransactionID string
	Status        string
	Amount        float64
	PaidAt        *time.Time
	RawPayload    string
}

func (r *Repository) ApplyPaymentNotification(ctx context.Context, notification PaymentNotification) (*models.Payment, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var bookingID uuid.UUID
	if err := tx.QueryRowContext(ctx, `SELECT id FROM bookings WHERE booking_code=$1 FOR UPDATE`, notification.OrderID).Scan(&bookingID); err != nil {
		return nil, err
	}

	if _, err := tx.ExecContext(ctx, `INSERT INTO payment_methods(id,name,type,provider) VALUES($1,'Midtrans','midtrans','midtrans') ON CONFLICT (id) DO NOTHING`, MidtransPaymentMethodID); err != nil {
		return nil, err
	}

	var paymentID uuid.UUID
	var currentStatus string
	err = tx.QueryRowContext(ctx, `
		SELECT id,status FROM payments
		WHERE booking_id=$1 AND payment_method_id=$2
		ORDER BY created_at DESC
		LIMIT 1
		FOR UPDATE
	`, bookingID, MidtransPaymentMethodID).Scan(&paymentID, &currentStatus)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		err = tx.QueryRowContext(ctx, `
			INSERT INTO payments(booking_id,payment_method_id,provider_transaction_id,amount,status,raw_callback_payload)
			VALUES($1,$2,$3,$4,'pending',$5::jsonb)
			RETURNING id,status
		`, bookingID, MidtransPaymentMethodID, notification.TransactionID, notification.Amount, notification.RawPayload).Scan(&paymentID, &currentStatus)
		if err != nil {
			return nil, err
		}
	case err != nil:
		return nil, err
	}

	if currentStatus == notification.Status {
		_, err = tx.ExecContext(ctx, `
			UPDATE payments
			SET provider_transaction_id=COALESCE(NULLIF($2,''),provider_transaction_id),
			    amount=$3,
			    paid_at=COALESCE($4,paid_at),
			    raw_callback_payload=$5::jsonb,
			    updated_at=NOW()
			WHERE id=$1
		`, paymentID, notification.TransactionID, notification.Amount, notification.PaidAt, notification.RawPayload)
	} else {
		_, err = tx.ExecContext(ctx, `
			UPDATE payments
			SET status=$2,
			    provider_transaction_id=COALESCE(NULLIF($3,''),provider_transaction_id),
			    amount=$4,
			    paid_at=COALESCE($5,paid_at),
			    raw_callback_payload=$6::jsonb,
			    updated_at=NOW()
			WHERE id=$1
		`, paymentID, notification.Status, notification.TransactionID, notification.Amount, notification.PaidAt, notification.RawPayload)
	}
	if err != nil {
		return nil, err
	}

	if currentStatus != notification.Status && (notification.Status == "expired" || notification.Status == "failed") {
		if _, err = tx.ExecContext(ctx, `
			UPDATE ticket_inventories ti
			SET held_count=GREATEST(ti.held_count-bi.quantity,0), version=ti.version+1, updated_at=NOW()
			FROM (
				SELECT ticket_tier_id, event_session_id, SUM(quantity) AS quantity
				FROM booking_items
				WHERE booking_id=$1
				GROUP BY ticket_tier_id, event_session_id
			) bi
			WHERE ti.ticket_tier_id=bi.ticket_tier_id
			  AND ti.event_session_id=bi.event_session_id
		`, bookingID); err != nil {
			return nil, err
		}
	}

	switch notification.Status {
	case "expired":
		_, err = tx.ExecContext(ctx, `UPDATE bookings SET status='expired' WHERE id=$1 AND status IN ('pending','awaiting_payment')`, bookingID)
	case "failed":
		_, err = tx.ExecContext(ctx, `UPDATE bookings SET status='cancelled',cancelled_at=NOW() WHERE id=$1 AND status IN ('pending','awaiting_payment')`, bookingID)
	case "refunded":
		_, err = tx.ExecContext(ctx, `UPDATE bookings SET status='refunded' WHERE id=$1`, bookingID)
	case "partially_refunded":
		_, err = tx.ExecContext(ctx, `UPDATE bookings SET status='partially_refunded' WHERE id=$1`, bookingID)
	}
	if err != nil {
		return nil, err
	}

	payment, err := scanPayment(tx.QueryRowContext(ctx, `
		SELECT id,booking_id,payment_method_id,provider_transaction_id,amount,status,expires_at,paid_at,created_at,updated_at
		FROM payments WHERE id=$1
	`, paymentID))
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return payment, nil
}

func (r *Repository) Payment(ctx context.Context, id uuid.UUID) (*models.Payment, error) {
	return scanPayment(r.db.QueryRowContext(ctx, `
		SELECT id,booking_id,payment_method_id,provider_transaction_id,amount,status,expires_at,paid_at,created_at,updated_at
		FROM payments WHERE id=$1
	`, id))
}

func (r *Repository) BookingByID(ctx context.Context, id uuid.UUID) (*models.Booking, error) {
	return scanBooking(r.db.QueryRowContext(ctx, `
		SELECT id,booking_code,user_id,event_id,status,subtotal_amount,discount_amount,platform_fee_amount,total_amount,promo_code_id,hold_expires_at,paid_at,cancelled_at,created_at,updated_at
		FROM bookings WHERE id=$1
	`, id))
}

func (r *Repository) Refund(ctx context.Context, id uuid.UUID) (*models.Refund, error) {
	return scanRefund(r.db.QueryRowContext(ctx, `SELECT `+refundCols+` FROM refunds WHERE id=$1 AND deleted_at IS NULL`, id))
}

func (r *Repository) RefundedAmount(ctx context.Context, paymentID uuid.UUID) (float64, error) {
	var amount float64
	err := r.db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(amount),0)
		FROM refunds
		WHERE payment_id=$1 AND status='processed' AND deleted_at IS NULL
	`, paymentID).Scan(&amount)
	return amount, err
}

func (r *Repository) CompleteRefund(ctx context.Context, refundID, paymentID, bookingID, processedBy uuid.UUID, paymentStatus, bookingStatus string) (*models.Refund, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	refund, err := scanRefund(tx.QueryRowContext(ctx, `
		UPDATE refunds
		SET status='processed', processed_by=$2, processed_at=NOW()
		WHERE id=$1 AND deleted_at IS NULL
		RETURNING `+refundCols, refundID, processedBy))
	if err != nil {
		return nil, err
	}

	if _, err := tx.ExecContext(ctx, `UPDATE payments SET status=$2, updated_at=NOW() WHERE id=$1`, paymentID, paymentStatus); err != nil {
		return nil, err
	}

	if _, err := tx.ExecContext(ctx, `UPDATE bookings SET status=$2 WHERE id=$1`, bookingID, bookingStatus); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return refund, nil
}
