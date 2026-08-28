package promo

import (
	"context"
	"database/sql"
	"konsera-backend/internal/DTO/promo"
	"konsera-backend/internal/models"
	"strings"

	"github.com/google/uuid"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

const cols = `id,code,event_id,discount_type,discount_value,max_discount_amount,usage_limit,usage_count,per_user_limit,valid_from,valid_until,is_active,created_at`

func scan(s interface{ Scan(...any) error }) (*models.PromoCode, error) {
	x := &models.PromoCode{}
	e := s.Scan(&x.ID, &x.Code, &x.EventID, &x.DiscountType, &x.DiscountValue, &x.MaxDiscountAmount, &x.UsageLimit, &x.UsageCount, &x.PerUserLimit, &x.ValidFrom, &x.ValidUntil, &x.IsActive, &x.CreatedAt)
	return x, e
}
func (r *Repository) List(ctx context.Context) ([]*models.PromoCode, error) {
	rows, e := r.db.QueryContext(ctx, `SELECT `+cols+` FROM promo_codes ORDER BY created_at DESC`)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []*models.PromoCode{}
	for rows.Next() {
		x, e := scan(rows)
		if e != nil {
			return nil, e
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *Repository) Get(ctx context.Context, id uuid.UUID) (*models.PromoCode, error) {
	return scan(r.db.QueryRowContext(ctx, `SELECT `+cols+` FROM promo_codes WHERE id=$1`, id))
}
func (r *Repository) Create(ctx context.Context, x *models.PromoCode) error {
	y, e := scan(r.db.QueryRowContext(ctx, `INSERT INTO promo_codes(code,event_id,discount_type,discount_value,max_discount_amount,usage_limit,per_user_limit,valid_from,valid_until) VALUES($1,$2,$3,$4,$5,$6,COALESCE(NULLIF($7,0),1),$8,$9) RETURNING `+cols, x.Code, x.EventID, x.DiscountType, x.DiscountValue, x.MaxDiscountAmount, x.UsageLimit, x.PerUserLimit, x.ValidFrom, x.ValidUntil))
	if e != nil {
		return e
	}
	*x = *y
	return nil
}
func (r *Repository) Update(ctx context.Context, id uuid.UUID, q *promo.UpdatePromoCodeRequest) (*models.PromoCode, error) {
	event := q.EventID
	var eventID *uuid.UUID
	if event != nil && strings.TrimSpace(*event) != "" {
		v, e := uuid.Parse(*event)
		if e != nil {
			return nil, e
		}
		eventID = &v
	}
	return scan(r.db.QueryRowContext(ctx, `UPDATE promo_codes SET code=COALESCE($2,code),event_id=COALESCE($3,event_id),discount_type=COALESCE($4,discount_type),discount_value=COALESCE($5,discount_value),max_discount_amount=COALESCE($6,max_discount_amount),usage_limit=COALESCE($7,usage_limit),per_user_limit=COALESCE($8,per_user_limit),valid_from=COALESCE($9,valid_from),valid_until=COALESCE($10,valid_until),is_active=COALESCE($11,is_active) WHERE id=$1 RETURNING `+cols, id, q.Code, eventID, q.DiscountType, q.DiscountValue, q.MaxDiscountAmount, q.UsageLimit, q.PerUserLimit, q.ValidFrom, q.ValidUntil, q.IsActive))
}
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	res, e := r.db.ExecContext(ctx, `DELETE FROM promo_codes WHERE id=$1`, id)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
func (r *Repository) Apply(ctx context.Context, user, booking uuid.UUID, code string) (*models.PromoCodeUsage, error) {
	tx, e := r.db.BeginTx(ctx, nil)
	if e != nil {
		return nil, e
	}
	defer tx.Rollback()
	p, e := scan(tx.QueryRowContext(ctx, `SELECT `+cols+` FROM promo_codes WHERE UPPER(code)=UPPER($1) AND is_active=TRUE AND NOW() BETWEEN valid_from AND valid_until FOR UPDATE`, code))
	if e != nil {
		return nil, e
	}
	var eventID uuid.UUID
	var subtotal float64
	if e = tx.QueryRowContext(ctx, `SELECT event_id,subtotal_amount FROM bookings WHERE id=$1 AND user_id=$2`, booking, user).Scan(&eventID, &subtotal); e != nil {
		return nil, e
	}
	if p.EventID != nil && *p.EventID != eventID {
		return nil, sql.ErrNoRows
	}
	var used int
	if e = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM promo_code_usages WHERE promo_code_id=$1 AND user_id=$2`, p.ID, user).Scan(&used); e != nil {
		return nil, e
	}
	if used >= p.PerUserLimit {
		return nil, sql.ErrNoRows
	}
	if p.UsageLimit != nil && p.UsageCount >= *p.UsageLimit {
		return nil, sql.ErrNoRows
	}
	discount := subtotal * p.DiscountValue / 100
	if p.DiscountType == "fixed_amount" {
		discount = p.DiscountValue
	}
	if p.MaxDiscountAmount != nil && discount > *p.MaxDiscountAmount {
		discount = *p.MaxDiscountAmount
	}
	if discount > subtotal {
		discount = subtotal
	}
	var out models.PromoCodeUsage
	e = tx.QueryRowContext(ctx, `INSERT INTO promo_code_usages(promo_code_id,user_id,booking_id,discount_applied) VALUES($1,$2,$3,$4) RETURNING id,promo_code_id,user_id,booking_id,discount_applied,used_at`, p.ID, user, booking, discount).Scan(&out.ID, &out.PromoCodeID, &out.UserID, &out.BookingID, &out.DiscountApplied, &out.UsedAt)
	if e != nil {
		return nil, e
	}
	_, e = tx.ExecContext(ctx, `UPDATE promo_codes SET usage_count=usage_count+1 WHERE id=$1`, p.ID)
	if e != nil {
		return nil, e
	}
	_, e = tx.ExecContext(ctx, `UPDATE bookings SET promo_code_id=$2,discount_amount=$3,total_amount=GREATEST(subtotal_amount+platform_fee_amount-$3,0) WHERE id=$1`, booking, p.ID, discount)
	if e != nil {
		return nil, e
	}
	if e = tx.Commit(); e != nil {
		return nil, e
	}
	return &out, nil
}
