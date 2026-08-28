package checkout

import (
	"context"
	"database/sql"
	"fmt"
	"konsera-backend/internal/models"

	"github.com/google/uuid"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }
func scanCart(s interface{ Scan(...any) error }) (*models.Cart, error) {
	x := &models.Cart{}
	e := s.Scan(&x.ID, &x.UserID, &x.TicketTierID, &x.EventSessionID, &x.SeatID, &x.Quantity, &x.HeldUntil, &x.CreatedAt)
	return x, e
}
func (r *Repository) Cart(ctx context.Context, user, id uuid.UUID) (*models.Cart, error) {
	return scanCart(r.db.QueryRowContext(ctx, `SELECT id,user_id,ticket_tier_id,event_session_id,seat_id,quantity,held_until,created_at FROM carts WHERE user_id=$1 AND id=$2`, user, id))
}
func (r *Repository) Carts(ctx context.Context, user uuid.UUID) ([]*models.Cart, error) {
	rows, e := r.db.QueryContext(ctx, `SELECT id,user_id,ticket_tier_id,event_session_id,seat_id,quantity,held_until,created_at FROM carts WHERE user_id=$1 AND held_until>NOW() ORDER BY created_at DESC`, user)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []*models.Cart{}
	for rows.Next() {
		x, e := scanCart(rows)
		if e != nil {
			return nil, e
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *Repository) CreateCart(ctx context.Context, x *models.Cart) error {
	created, e := scanCart(r.db.QueryRowContext(ctx, `INSERT INTO carts(user_id,ticket_tier_id,event_session_id,seat_id,quantity,held_until) VALUES($1,$2,$3,$4,$5,$6) RETURNING id,user_id,ticket_tier_id,event_session_id,seat_id,quantity,held_until,created_at`, x.UserID, x.TicketTierID, x.EventSessionID, x.SeatID, x.Quantity, x.HeldUntil))
	if e != nil {
		return e
	}
	*x = *created
	return nil
}
func (r *Repository) UpdateCart(ctx context.Context, user, id uuid.UUID, qty int) (*models.Cart, error) {
	return scanCart(r.db.QueryRowContext(ctx, `UPDATE carts SET quantity=$3 WHERE user_id=$1 AND id=$2 AND held_until>NOW() RETURNING id,user_id,ticket_tier_id,event_session_id,seat_id,quantity,held_until,created_at`, user, id, qty))
}
func (r *Repository) DeleteCart(ctx context.Context, user, id uuid.UUID) error {
	res, e := r.db.ExecContext(ctx, `DELETE FROM carts WHERE user_id=$1 AND id=$2`, user, id)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
func (r *Repository) Checkout(ctx context.Context, user, event uuid.UUID, expiry *string) (*models.Booking, error) {
	tx, e := r.db.BeginTx(ctx, nil)
	if e != nil {
		return nil, e
	}
	defer tx.Rollback()
	var hold interface{}
	if expiry != nil {
		hold = *expiry
	}
	row := tx.QueryRowContext(ctx, `INSERT INTO bookings(booking_code,user_id,event_id,status,hold_expires_at) VALUES(upper(substr(md5(random()::text),1,12)),$1,$2,'awaiting_payment',NULLIF($3,'')::timestamptz) RETURNING id,booking_code,user_id,event_id,status,subtotal_amount,discount_amount,platform_fee_amount,total_amount,promo_code_id,hold_expires_at,paid_at,cancelled_at,created_at,updated_at`, user, event, hold)
	x, e := scanBooking(row)
	if e != nil {
		return nil, e
	}
	rows, e := tx.QueryContext(ctx, `SELECT c.ticket_tier_id,c.event_session_id,c.quantity,tt.price FROM carts c JOIN ticket_tiers tt ON tt.id=c.ticket_tier_id WHERE c.user_id=$1 AND c.held_until>NOW() AND EXISTS(SELECT 1 FROM events ev WHERE ev.id=$2 AND ev.id=tt.event_id) FOR UPDATE`, user, event)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var subtotal float64
	count := 0
	for rows.Next() {
		var tier, session uuid.UUID
		var qty int
		var price float64
		if e = rows.Scan(&tier, &session, &qty, &price); e != nil {
			return nil, e
		}
		line := price * float64(qty)
		_, e = tx.ExecContext(ctx, `INSERT INTO booking_items(booking_id,ticket_tier_id,event_session_id,unit_price,quantity,line_total) VALUES($1,$2,$3,$4,$5,$6)`, x.ID, tier, session, price, qty, line)
		if e != nil {
			return nil, e
		}
		subtotal += line
		count++
	}
	if e = rows.Err(); e != nil {
		return nil, e
	}
	if count == 0 {
		return nil, fmt.Errorf("cart is empty")
	}
	_, e = tx.ExecContext(ctx, `UPDATE bookings SET subtotal_amount=$2,total_amount=$2 WHERE id=$1`, x.ID, subtotal)
	if e != nil {
		return nil, e
	}
	if _, e = tx.ExecContext(ctx, `DELETE FROM carts WHERE user_id=$1 AND held_until>NOW() AND EXISTS(SELECT 1 FROM ticket_tiers tt WHERE tt.id=carts.ticket_tier_id AND tt.event_id=$2)`, user, event); e != nil {
		return nil, e
	}
	x.SubtotalAmount = subtotal
	x.TotalAmount = subtotal
	if e = tx.Commit(); e != nil {
		return nil, e
	}
	return x, nil
}
func scanBooking(s interface{ Scan(...any) error }) (*models.Booking, error) {
	x := &models.Booking{}
	e := s.Scan(&x.ID, &x.BookingCode, &x.UserID, &x.EventID, &x.Status, &x.SubtotalAmount, &x.DiscountAmount, &x.PlatformFeeAmount, &x.TotalAmount, &x.PromoCodeID, &x.HoldExpiresAt, &x.PaidAt, &x.CancelledAt, &x.CreatedAt, &x.UpdatedAt)
	return x, e
}
func (r *Repository) Bookings(ctx context.Context, user uuid.UUID) ([]*models.Booking, error) {
	rows, e := r.db.QueryContext(ctx, `SELECT id,booking_code,user_id,event_id,status,subtotal_amount,discount_amount,platform_fee_amount,total_amount,promo_code_id,hold_expires_at,paid_at,cancelled_at,created_at,updated_at FROM bookings WHERE user_id=$1 ORDER BY created_at DESC`, user)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []*models.Booking{}
	for rows.Next() {
		x, e := scanBooking(rows)
		if e != nil {
			return nil, e
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *Repository) Booking(ctx context.Context, user, id uuid.UUID) (*models.Booking, error) {
	return scanBooking(r.db.QueryRowContext(ctx, `SELECT id,booking_code,user_id,event_id,status,subtotal_amount,discount_amount,platform_fee_amount,total_amount,promo_code_id,hold_expires_at,paid_at,cancelled_at,created_at,updated_at FROM bookings WHERE user_id=$1 AND id=$2`, user, id))
}
func (r *Repository) UpdateBooking(ctx context.Context, user, id uuid.UUID, status string) (*models.Booking, error) {
	return scanBooking(r.db.QueryRowContext(ctx, `UPDATE bookings SET status=$3 WHERE user_id=$1 AND id=$2 RETURNING id,booking_code,user_id,event_id,status,subtotal_amount,discount_amount,platform_fee_amount,total_amount,promo_code_id,hold_expires_at,paid_at,cancelled_at,created_at,updated_at`, user, id, status))
}
func (r *Repository) DeleteBooking(ctx context.Context, user, id uuid.UUID) error {
	res, e := r.db.ExecContext(ctx, `UPDATE bookings SET status='cancelled',cancelled_at=NOW() WHERE user_id=$1 AND id=$2`, user, id)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
