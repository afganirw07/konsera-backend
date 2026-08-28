package checkout

import (
	"context"
	"database/sql"
	dto "konsera-backend/internal/DTO/checkout"
	"konsera-backend/internal/models"

	"github.com/google/uuid"
)

func scanRefund(s interface{ Scan(...any) error }) (*models.Refund, error) {
	x := &models.Refund{}
	e := s.Scan(&x.ID, &x.BookingID, &x.PaymentID, &x.Amount, &x.Reason, &x.Status, &x.RequestedBy, &x.ProcessedBy, &x.ProcessedAt, &x.CreatedAt)
	return x, e
}

const refundCols = `id,booking_id,payment_id,amount,reason,status,requested_by,processed_by,processed_at,created_at`

func (r *Repository) CreateRefund(ctx context.Context, x *models.Refund) error {
	y, e := scanRefund(r.db.QueryRowContext(ctx, `INSERT INTO refunds(booking_id,payment_id,amount,reason,requested_by) VALUES($1,$2,$3,$4,$5) RETURNING `+refundCols, x.BookingID, x.PaymentID, x.Amount, x.Reason, x.RequestedBy))
	if e != nil {
		return e
	}
	*x = *y
	return nil
}
func (r *Repository) Refunds(ctx context.Context, user uuid.UUID) ([]*models.Refund, error) {
	rows, e := r.db.QueryContext(ctx, `SELECT `+refundCols+` FROM refunds r JOIN bookings b ON b.id=r.booking_id WHERE b.user_id=$1 ORDER BY r.created_at DESC`, user)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []*models.Refund{}
	for rows.Next() {
		x, e := scanRefund(rows)
		if e != nil {
			return nil, e
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *Repository) UpdateRefund(ctx context.Context, id uuid.UUID, q *dto.UpdateRefundRequest) (*models.Refund, error) {
	return scanRefund(r.db.QueryRowContext(ctx, `UPDATE refunds SET status=$2,processed_at=CASE WHEN $2='processed' THEN NOW() ELSE processed_at END WHERE id=$1 RETURNING `+refundCols, id, q.Status))
}
func (r *Repository) DeleteRefund(ctx context.Context, id uuid.UUID) error {
	res, e := r.db.ExecContext(ctx, `DELETE FROM refunds WHERE id=$1`, id)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
