package ticket

import (
	"context"
	"database/sql"

	dto "konsera-backend/internal/DTO/ticket"
	"konsera-backend/internal/models"

	"github.com/google/uuid"
)

type TicketRepository struct{ db *sql.DB }

func NewTicketRepository(db *sql.DB) *TicketRepository { return &TicketRepository{db: db} }

const tierColumns = `id,event_id,venue_section_id,name,tier_type,price,max_per_transaction,
sale_start_at,sale_end_at,is_seated,created_at,updated_at`

func scanTier(s interface{ Scan(...any) error }) (*models.TicketTier, error) {
	x := &models.TicketTier{}
	err := s.Scan(&x.ID, &x.EventID, &x.VenueSectionID, &x.Name, &x.TierType, &x.Price,
		&x.MaxPerTransaction, &x.SaleStartAt, &x.SaleEndAt, &x.IsSeated, &x.CreatedAt, &x.UpdatedAt)
	return x, err
}

func (r *TicketRepository) CreateTier(ctx context.Context, x *models.TicketTier) error {
	created, err := scanTier(r.db.QueryRowContext(ctx, `INSERT INTO ticket_tiers
		(event_id,venue_section_id,name,tier_type,price,max_per_transaction,sale_start_at,sale_end_at,is_seated)
		VALUES($1,$2,$3,$4,$5,COALESCE(NULLIF($6,0),6),$7,$8,$9)
		RETURNING `+tierColumns, x.EventID, x.VenueSectionID, x.Name, x.TierType, x.Price,
		x.MaxPerTransaction, x.SaleStartAt, x.SaleEndAt, x.IsSeated))
	if err != nil {
		return err
	}
	*x = *created
	return nil
}

func (r *TicketRepository) ListTiers(ctx context.Context, eventID uuid.UUID) ([]*models.TicketTier, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+tierColumns+` FROM ticket_tiers WHERE event_id=$1 ORDER BY price,id`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]*models.TicketTier, 0)
	for rows.Next() {
		item, err := scanTier(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *TicketRepository) GetTier(ctx context.Context, eventID, id uuid.UUID) (*models.TicketTier, error) {
	return scanTier(r.db.QueryRowContext(ctx, `SELECT `+tierColumns+` FROM ticket_tiers WHERE event_id=$1 AND id=$2`, eventID, id))
}

func (r *TicketRepository) UpdateTier(ctx context.Context, eventID, id uuid.UUID, q *dto.UpdateTicketTierRequest) (*models.TicketTier, error) {
	return scanTier(r.db.QueryRowContext(ctx, `UPDATE ticket_tiers SET
		venue_section_id=COALESCE($3,venue_section_id),name=COALESCE($4,name),tier_type=COALESCE($5,tier_type),
		price=COALESCE($6,price),max_per_transaction=COALESCE($7,max_per_transaction),sale_start_at=COALESCE($8,sale_start_at),
		sale_end_at=COALESCE($9,sale_end_at),is_seated=COALESCE($10,is_seated)
		WHERE event_id=$1 AND id=$2 RETURNING `+tierColumns, eventID, id, q.VenueSectionID, q.Name, q.TierType,
		q.Price, q.MaxPerTransaction, q.SaleStartAt, q.SaleEndAt, q.IsSeated))
}

func (r *TicketRepository) DeleteTier(ctx context.Context, eventID, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM ticket_tiers WHERE event_id=$1 AND id=$2`, eventID, id)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

const inventoryColumns = `id,ticket_tier_id,event_session_id,total_quota,sold_count,held_count,version,updated_at`

func scanInventory(s interface{ Scan(...any) error }) (*models.TicketInventory, error) {
	x := &models.TicketInventory{}
	err := s.Scan(&x.ID, &x.TicketTierID, &x.EventSessionID, &x.TotalQuota, &x.SoldCount, &x.HeldCount, &x.Version, &x.UpdatedAt)
	return x, err
}
func (r *TicketRepository) CreateInventory(ctx context.Context, x *models.TicketInventory) error {
	created, err := scanInventory(r.db.QueryRowContext(ctx, `INSERT INTO ticket_inventories(ticket_tier_id,event_session_id,total_quota) VALUES($1,$2,$3) RETURNING `+inventoryColumns, x.TicketTierID, x.EventSessionID, x.TotalQuota))
	if err != nil {
		return err
	}
	*x = *created
	return nil
}
func (r *TicketRepository) ListInventory(ctx context.Context, tierID uuid.UUID) ([]*models.TicketInventory, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+inventoryColumns+` FROM ticket_inventories WHERE ticket_tier_id=$1 ORDER BY event_session_id`, tierID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*models.TicketInventory{}
	for rows.Next() {
		x, e := scanInventory(rows)
		if e != nil {
			return nil, e
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *TicketRepository) GetInventory(ctx context.Context, tierID, sessionID uuid.UUID) (*models.TicketInventory, error) {
	return scanInventory(r.db.QueryRowContext(ctx, `SELECT `+inventoryColumns+` FROM ticket_inventories WHERE ticket_tier_id=$1 AND event_session_id=$2`, tierID, sessionID))
}
func (r *TicketRepository) UpdateInventory(ctx context.Context, tierID, sessionID uuid.UUID, q *dto.UpdateInventoryRequest) (*models.TicketInventory, error) {
	return scanInventory(r.db.QueryRowContext(ctx, `UPDATE ticket_inventories SET total_quota=COALESCE($3,total_quota),sold_count=COALESCE($4,sold_count),held_count=COALESCE($5,held_count),version=version+1,updated_at=NOW() WHERE ticket_tier_id=$1 AND event_session_id=$2 AND COALESCE($3,total_quota)>=COALESCE($4,sold_count)+COALESCE($5,held_count) RETURNING `+inventoryColumns, tierID, sessionID, q.TotalQuota, q.SoldCount, q.HeldCount))
}
func (r *TicketRepository) DeleteInventory(ctx context.Context, tierID, sessionID uuid.UUID) error {
	res, e := r.db.ExecContext(ctx, `DELETE FROM ticket_inventories WHERE ticket_tier_id=$1 AND event_session_id=$2`, tierID, sessionID)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
