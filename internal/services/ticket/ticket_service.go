package ticket

import (
	"context"
	"fmt"

	dto "konsera-backend/internal/DTO/ticket"
	"konsera-backend/internal/models"
	repo "konsera-backend/internal/repository/ticket"

	"github.com/google/uuid"
)

type TicketService struct{ repo *repo.TicketRepository }

func NewTicketService(r *repo.TicketRepository) *TicketService { return &TicketService{repo: r} }
func parseID(v, n string) (uuid.UUID, error) {
	id, err := uuid.Parse(v)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid %s", n)
	}
	return id, nil
}
func (s *TicketService) CreateTier(c context.Context, eventValue string, q *dto.CreateTicketTierRequest) (*models.TicketTier, error) {
	eventID, e := parseID(eventValue, "event_id")
	if e != nil {
		return nil, e
	}
	if q.SaleStartAt != nil && q.SaleEndAt != nil && !q.SaleEndAt.After(*q.SaleStartAt) {
		return nil, fmt.Errorf("sale_end_at must be after sale_start_at")
	}
	section := (*uuid.UUID)(nil)
	if q.VenueSectionID != nil {
		x, err := parseID(*q.VenueSectionID, "venue_section_id")
		if err != nil {
			return nil, err
		}
		section = &x
	}
	max := q.MaxPerTransaction
	if max == 0 {
		max = 6
	}
	x := &models.TicketTier{EventID: eventID, VenueSectionID: section, Name: q.Name, TierType: q.TierType, Price: q.Price, MaxPerTransaction: max, SaleStartAt: q.SaleStartAt, SaleEndAt: q.SaleEndAt, IsSeated: q.IsSeated}
	e = s.repo.CreateTier(c, x)
	return x, e
}
func (s *TicketService) Tiers(c context.Context, v string) ([]*models.TicketTier, error) {
	id, e := parseID(v, "event_id")
	if e != nil {
		return nil, e
	}
	return s.repo.ListTiers(c, id)
}
func (s *TicketService) GetTier(c context.Context, ev, v string) (*models.TicketTier, error) {
	a, e := parseID(ev, "event_id")
	if e != nil {
		return nil, e
	}
	b, e := parseID(v, "tier_id")
	if e != nil {
		return nil, e
	}
	return s.repo.GetTier(c, a, b)
}
func (s *TicketService) UpdateTier(c context.Context, ev, v string, q *dto.UpdateTicketTierRequest) (*models.TicketTier, error) {
	a, e := parseID(ev, "event_id")
	if e != nil {
		return nil, e
	}
	b, e := parseID(v, "tier_id")
	if e != nil {
		return nil, e
	}
	if q.VenueSectionID != nil {
		if _, e = parseID(*q.VenueSectionID, "venue_section_id"); e != nil {
			return nil, e
		}
	}
	if q.SaleStartAt != nil && q.SaleEndAt != nil && !q.SaleEndAt.After(*q.SaleStartAt) {
		return nil, fmt.Errorf("sale_end_at must be after sale_start_at")
	}
	return s.repo.UpdateTier(c, a, b, q)
}
func (s *TicketService) DeleteTier(c context.Context, ev, v string) error {
	a, e := parseID(ev, "event_id")
	if e != nil {
		return e
	}
	b, e := parseID(v, "tier_id")
	if e != nil {
		return e
	}
	return s.repo.DeleteTier(c, a, b)
}
func (s *TicketService) CreateInventory(c context.Context, tierValue string, q *dto.CreateInventoryRequest) (*models.TicketInventory, error) {
	tier, e := parseID(tierValue, "tier_id")
	if e != nil {
		return nil, e
	}
	session, e := parseID(q.EventSessionID, "event_session_id")
	if e != nil {
		return nil, e
	}
	x := &models.TicketInventory{TicketTierID: tier, EventSessionID: session, TotalQuota: q.TotalQuota}
	e = s.repo.CreateInventory(c, x)
	return x, e
}
func (s *TicketService) Inventories(c context.Context, v string) ([]*models.TicketInventory, error) {
	id, e := parseID(v, "tier_id")
	if e != nil {
		return nil, e
	}
	return s.repo.ListInventory(c, id)
}
func (s *TicketService) GetInventory(c context.Context, tier, session string) (*models.TicketInventory, error) {
	a, e := parseID(tier, "tier_id")
	if e != nil {
		return nil, e
	}
	b, e := parseID(session, "event_session_id")
	if e != nil {
		return nil, e
	}
	return s.repo.GetInventory(c, a, b)
}
func (s *TicketService) UpdateInventory(c context.Context, tier, session string, q *dto.UpdateInventoryRequest) (*models.TicketInventory, error) {
	a, e := parseID(tier, "tier_id")
	if e != nil {
		return nil, e
	}
	b, e := parseID(session, "event_session_id")
	if e != nil {
		return nil, e
	}
	return s.repo.UpdateInventory(c, a, b, q)
}
func (s *TicketService) DeleteInventory(c context.Context, tier, session string) error {
	a, e := parseID(tier, "tier_id")
	if e != nil {
		return e
	}
	b, e := parseID(session, "event_session_id")
	if e != nil {
		return e
	}
	return s.repo.DeleteInventory(c, a, b)
}
