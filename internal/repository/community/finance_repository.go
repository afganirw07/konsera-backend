package community

import (
	"context"
	"database/sql"
	dto "konsera-backend/internal/DTO/community"
	"konsera-backend/internal/models"

	"github.com/google/uuid"
)

func scanCommission(s interface{ Scan(...any) error }) (*models.CommissionSetting, error) {
	x := &models.CommissionSetting{}
	e := s.Scan(&x.ID, &x.Scope, &x.ScopeRefID, &x.RatePercentage, &x.EffectiveFrom, &x.CreatedAt)
	return x, e
}

const commissionCols = `id,scope,scope_ref_id,rate_percentage,effective_from,created_at`

func (r *Repository) Commissions(ctx context.Context) ([]*models.CommissionSetting, error) {
	rows, e := r.db.QueryContext(ctx, `SELECT `+commissionCols+` FROM commission_settings ORDER BY effective_from DESC`)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []*models.CommissionSetting{}
	for rows.Next() {
		x, e := scanCommission(rows)
		if e != nil {
			return nil, e
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *Repository) CreateCommission(ctx context.Context, x *models.CommissionSetting) error {
	y, e := scanCommission(r.db.QueryRowContext(ctx, `INSERT INTO commission_settings(scope,scope_ref_id,rate_percentage,effective_from) VALUES($1,$2,$3,$4) RETURNING `+commissionCols, x.Scope, x.ScopeRefID, x.RatePercentage, x.EffectiveFrom))
	if e != nil {
		return e
	}
	*x = *y
	return nil
}
func (r *Repository) UpdateCommission(ctx context.Context, id uuid.UUID, q *dto.UpdateCommissionRequest) (*models.CommissionSetting, error) {
	return scanCommission(r.db.QueryRowContext(ctx, `UPDATE commission_settings SET rate_percentage=COALESCE($2,rate_percentage),effective_from=COALESCE($3,effective_from) WHERE id=$1 RETURNING `+commissionCols, id, q.RatePercentage, q.EffectiveFrom))
}
func (r *Repository) DeleteCommission(ctx context.Context, id uuid.UUID) error {
	res, e := r.db.ExecContext(ctx, `DELETE FROM commission_settings WHERE id=$1`, id)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
func scanPayout(s interface{ Scan(...any) error }) (*models.Payout, error) {
	x := &models.Payout{}
	e := s.Scan(&x.ID, &x.OrganizerID, &x.EventID, &x.GrossAmount, &x.CommissionAmount, &x.NetAmount, &x.Status, &x.ScheduledAt, &x.PaidAt, &x.BankReference, &x.CreatedAt)
	return x, e
}

const payoutCols = `id,organizer_id,event_id,gross_amount,commission_amount,net_amount,status,scheduled_at,paid_at,bank_reference,created_at`

func (r *Repository) Payouts(ctx context.Context) ([]*models.Payout, error) {
	rows, e := r.db.QueryContext(ctx, `SELECT `+payoutCols+` FROM payouts ORDER BY scheduled_at DESC`)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []*models.Payout{}
	for rows.Next() {
		x, e := scanPayout(rows)
		if e != nil {
			return nil, e
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *Repository) CreatePayout(ctx context.Context, x *models.Payout) error {
	y, e := scanPayout(r.db.QueryRowContext(ctx, `INSERT INTO payouts(organizer_id,event_id,gross_amount,commission_amount,net_amount,scheduled_at,bank_reference) VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING `+payoutCols, x.OrganizerID, x.EventID, x.GrossAmount, x.CommissionAmount, x.NetAmount, x.ScheduledAt, x.BankReference))
	if e != nil {
		return e
	}
	*x = *y
	return nil
}
func (r *Repository) UpdatePayout(ctx context.Context, id uuid.UUID, q *dto.UpdatePayoutRequest) (*models.Payout, error) {
	return scanPayout(r.db.QueryRowContext(ctx, `UPDATE payouts SET status=COALESCE($2,status),paid_at=COALESCE($3,paid_at),bank_reference=COALESCE($4,bank_reference) WHERE id=$1 RETURNING `+payoutCols, id, q.Status, q.PaidAt, q.BankReference))
}
func (r *Repository) DeletePayout(ctx context.Context, id uuid.UUID) error {
	res, e := r.db.ExecContext(ctx, `DELETE FROM payouts WHERE id=$1`, id)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
func scanAudit(s interface{ Scan(...any) error }) (*models.AuditLog, error) {
	x := &models.AuditLog{}
	e := s.Scan(&x.ID, &x.ActorUserID, &x.Action, &x.EntityType, &x.EntityID, &x.BeforeState, &x.AfterState, &x.IPAddress, &x.UserAgent, &x.CreatedAt)
	return x, e
}

const auditCols = `id,actor_user_id,action,entity_type,entity_id,before_state,after_state,ip_address,user_agent,created_at`

func (r *Repository) Audits(ctx context.Context) ([]*models.AuditLog, error) {
	rows, e := r.db.QueryContext(ctx, `SELECT `+auditCols+` FROM audit_logs ORDER BY created_at DESC`)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []*models.AuditLog{}
	for rows.Next() {
		x, e := scanAudit(rows)
		if e != nil {
			return nil, e
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *Repository) CreateAudit(ctx context.Context, x *models.AuditLog) error {
	y, e := scanAudit(r.db.QueryRowContext(ctx, `INSERT INTO audit_logs(actor_user_id,action,entity_type,entity_id,before_state,after_state,ip_address,user_agent) VALUES($1,$2,$3,$4,$5,$6,$7,$8) RETURNING `+auditCols, x.ActorUserID, x.Action, x.EntityType, x.EntityID, x.BeforeState, x.AfterState, x.IPAddress, x.UserAgent))
	if e != nil {
		return e
	}
	*x = *y
	return nil
}
func (r *Repository) DeleteAudit(ctx context.Context, id uuid.UUID) error {
	res, e := r.db.ExecContext(ctx, `DELETE FROM audit_logs WHERE id=$1`, id)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
