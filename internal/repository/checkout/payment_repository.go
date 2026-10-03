package checkout

import (
	"context"
	"database/sql"
	"konsera-backend/internal/DTO/checkout"
	"konsera-backend/internal/models"

	"github.com/google/uuid"
)

func scanMethod(s interface{ Scan(...any) error }) (*models.PaymentMethod, error) {
	x := &models.PaymentMethod{}
	e := s.Scan(&x.ID, &x.Name, &x.Type, &x.Provider, &x.IsActive, &x.IconURL)
	return x, e
}
func (r *Repository) Methods(ctx context.Context) ([]*models.PaymentMethod, error) {
	rows, e := r.db.QueryContext(ctx, `SELECT id,name,type,provider,is_active,icon_url FROM payment_methods ORDER BY name`)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []*models.PaymentMethod{}
	for rows.Next() {
		x, e := scanMethod(rows)
		if e != nil {
			return nil, e
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *Repository) CreateMethod(ctx context.Context, x *models.PaymentMethod) error {
	return r.db.QueryRowContext(ctx, `INSERT INTO payment_methods(name,type,provider,icon_url) VALUES($1,$2,$3,$4) RETURNING id,name,type,provider,is_active,icon_url`, x.Name, x.Type, x.Provider, x.IconURL).Scan(&x.ID, &x.Name, &x.Type, &x.Provider, &x.IsActive, &x.IconURL)
}
func (r *Repository) UpdateMethod(ctx context.Context, id uuid.UUID, q *checkout.UpdatePaymentMethodRequest) (*models.PaymentMethod, error) {
	return scanMethod(r.db.QueryRowContext(ctx, `UPDATE payment_methods SET name=COALESCE($2,name),type=COALESCE($3,type),provider=COALESCE($4,provider),is_active=COALESCE($5,is_active),icon_url=COALESCE($6,icon_url) WHERE id=$1 RETURNING id,name,type,provider,is_active,icon_url`, id, q.Name, q.Type, q.Provider, q.IsActive, q.IconURL))
}
func (r *Repository) DeleteMethod(ctx context.Context, id uuid.UUID) error {
	res, e := r.db.ExecContext(ctx, `DELETE FROM payment_methods WHERE id=$1`, id)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
func scanPayment(s interface{ Scan(...any) error }) (*models.Payment, error) {
	x := &models.Payment{}
	e := s.Scan(&x.ID, &x.BookingID, &x.PaymentMethodID, &x.ProviderTransactionID, &x.Amount, &x.Status, &x.ExpiresAt, &x.PaidAt, &x.CreatedAt, &x.UpdatedAt)
	return x, e
}
func (r *Repository) Payments(ctx context.Context, user uuid.UUID, id *uuid.UUID) ([]*models.Payment, error) {
	rows, e := r.db.QueryContext(ctx, `SELECT p.id,p.booking_id,p.payment_method_id,p.provider_transaction_id,p.amount,p.status,p.expires_at,p.paid_at,p.created_at,p.updated_at FROM payments p JOIN bookings b ON b.id=p.booking_id WHERE b.user_id=$1 AND ($2::uuid IS NULL OR p.booking_id=$2) ORDER BY p.created_at DESC`, user, id)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []*models.Payment{}
	for rows.Next() {
		x, e := scanPayment(rows)
		if e != nil {
			return nil, e
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *Repository) CreatePayment(ctx context.Context, user uuid.UUID, x *models.Payment) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO payments(booking_id,payment_method_id,amount,expires_at)
		SELECT b.id, $2, b.total_amount, $3
		FROM bookings b
		JOIN payment_methods pm ON pm.id=$2 AND pm.is_active=TRUE
		WHERE b.id=$1 AND b.user_id=$4 AND b.status IN ('pending','awaiting_payment')
		RETURNING id,booking_id,payment_method_id,provider_transaction_id,amount,status,expires_at,paid_at,created_at,updated_at
	`, x.BookingID, x.PaymentMethodID, x.ExpiresAt, user).Scan(&x.ID, &x.BookingID, &x.PaymentMethodID, &x.ProviderTransactionID, &x.Amount, &x.Status, &x.ExpiresAt, &x.PaidAt, &x.CreatedAt, &x.UpdatedAt)
}
func (r *Repository) UpdatePayment(ctx context.Context, user, id uuid.UUID, q *checkout.UpdatePaymentRequest) (*models.Payment, error) {
	return scanPayment(r.db.QueryRowContext(ctx, `UPDATE payments p SET status=COALESCE($3,status),provider_transaction_id=COALESCE($4,provider_transaction_id),paid_at=COALESCE($5,paid_at),updated_at=NOW() FROM bookings b WHERE p.booking_id=b.id AND b.user_id=$1 AND p.id=$2 RETURNING p.id,p.booking_id,p.payment_method_id,p.provider_transaction_id,p.amount,p.status,p.expires_at,p.paid_at,p.created_at,p.updated_at`, user, id, q.Status, q.ProviderTransactionID, q.PaidAt))
}
