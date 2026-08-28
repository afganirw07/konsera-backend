package community

import (
	"context"
	"fmt"
	dto "konsera-backend/internal/DTO/community"
	"konsera-backend/internal/models"
	repo "konsera-backend/internal/repository/community"
	"time"

	"github.com/google/uuid"
)

type Service struct{ repo *repo.Repository }

func NewService(r *repo.Repository) *Service { return &Service{repo: r} }
func id(v, n string) (uuid.UUID, error) {
	x, e := uuid.Parse(v)
	if e != nil {
		return uuid.Nil, fmt.Errorf("invalid %s", n)
	}
	return x, nil
}
func (s *Service) Reviews(c context.Context, event string) ([]*models.Review, error) {
	var p *uuid.UUID
	if event != "" {
		x, e := id(event, "event_id")
		if e != nil {
			return nil, e
		}
		p = &x
	}
	return s.repo.Reviews(c, p)
}
func (s *Service) Review(c context.Context, user, v string) (*models.Review, error) {
	u, e := id(user, "user_id")
	if e != nil {
		return nil, e
	}
	x, e := id(v, "review_id")
	if e != nil {
		return nil, e
	}
	return s.repo.Review(c, u, x)
}
func (s *Service) CreateReview(c context.Context, user string, q *dto.CreateReviewRequest) (*models.Review, error) {
	u, e := id(user, "user_id")
	if e != nil {
		return nil, e
	}
	ev, e := id(q.EventID, "event_id")
	if e != nil {
		return nil, e
	}
	b, e := id(q.BookingID, "booking_id")
	if e != nil {
		return nil, e
	}
	x := &models.Review{EventID: ev, UserID: u, BookingID: b, Rating: q.Rating, Comment: q.Comment, PhotoURLs: q.PhotoURLs, Status: "visible"}
	e = s.repo.CreateReview(c, x)
	return x, e
}
func (s *Service) UpdateReview(c context.Context, user, v string, q *dto.UpdateReviewRequest) (*models.Review, error) {
	u, e := id(user, "user_id")
	if e != nil {
		return nil, e
	}
	x, e := id(v, "review_id")
	if e != nil {
		return nil, e
	}
	return s.repo.UpdateReview(c, u, x, q)
}
func (s *Service) DeleteReview(c context.Context, user, v string) error {
	u, e := id(user, "user_id")
	if e != nil {
		return e
	}
	x, e := id(v, "review_id")
	if e != nil {
		return e
	}
	return s.repo.DeleteReview(c, u, x)
}
func (s *Service) Favorites(c context.Context, user string) ([]*models.Favorite, error) {
	u, e := id(user, "user_id")
	if e != nil {
		return nil, e
	}
	return s.repo.Favorites(c, u)
}
func (s *Service) CreateFavorite(c context.Context, user string, q *dto.CreateFavoriteRequest) (*models.Favorite, error) {
	u, e := id(user, "user_id")
	if e != nil {
		return nil, e
	}
	ev, e := id(q.EventID, "event_id")
	if e != nil {
		return nil, e
	}
	x := &models.Favorite{UserID: u, EventID: ev}
	e = s.repo.CreateFavorite(c, x)
	return x, e
}
func (s *Service) DeleteFavorite(c context.Context, user, event string) error {
	u, e := id(user, "user_id")
	if e != nil {
		return e
	}
	ev, e := id(event, "event_id")
	if e != nil {
		return e
	}
	return s.repo.DeleteFavorite(c, u, ev)
}
func (s *Service) Commissions(c context.Context) ([]*models.CommissionSetting, error) {
	return s.repo.Commissions(c)
}
func (s *Service) CreateCommission(c context.Context, q *dto.CreateCommissionRequest) (*models.CommissionSetting, error) {
	var ref *uuid.UUID
	var e error
	if q.ScopeRefID != nil {
		v, e2 := id(*q.ScopeRefID, "scope_ref_id")
		if e2 != nil {
			return nil, e2
		}
		ref = &v
	}
	t := time.Now()
	if q.EffectiveFrom != nil {
		t, e = time.Parse(time.RFC3339, *q.EffectiveFrom)
		if e != nil {
			return nil, fmt.Errorf("invalid effective_from")
		}
	}
	x := &models.CommissionSetting{Scope: q.Scope, ScopeRefID: ref, RatePercentage: q.RatePercentage, EffectiveFrom: t}
	e = s.repo.CreateCommission(c, x)
	return x, e
}
func (s *Service) UpdateCommission(c context.Context, v string, q *dto.UpdateCommissionRequest) (*models.CommissionSetting, error) {
	x, e := id(v, "commission_id")
	if e != nil {
		return nil, e
	}
	return s.repo.UpdateCommission(c, x, q)
}
func (s *Service) DeleteCommission(c context.Context, v string) error {
	x, e := id(v, "commission_id")
	if e != nil {
		return e
	}
	return s.repo.DeleteCommission(c, x)
}
func (s *Service) Payouts(c context.Context) ([]*models.Payout, error) { return s.repo.Payouts(c) }
func (s *Service) CreatePayout(c context.Context, q *dto.CreatePayoutRequest) (*models.Payout, error) {
	o, e := id(q.OrganizerID, "organizer_id")
	if e != nil {
		return nil, e
	}
	ev, e := id(q.EventID, "event_id")
	if e != nil {
		return nil, e
	}
	t, e := time.Parse(time.RFC3339, q.ScheduledAt)
	if e != nil {
		return nil, fmt.Errorf("invalid scheduled_at")
	}
	x := &models.Payout{OrganizerID: o, EventID: ev, GrossAmount: q.GrossAmount, CommissionAmount: q.CommissionAmount, NetAmount: q.NetAmount, ScheduledAt: t, BankReference: q.BankReference}
	e = s.repo.CreatePayout(c, x)
	return x, e
}
func (s *Service) UpdatePayout(c context.Context, v string, q *dto.UpdatePayoutRequest) (*models.Payout, error) {
	x, e := id(v, "payout_id")
	if e != nil {
		return nil, e
	}
	return s.repo.UpdatePayout(c, x, q)
}
func (s *Service) DeletePayout(c context.Context, v string) error {
	x, e := id(v, "payout_id")
	if e != nil {
		return e
	}
	return s.repo.DeletePayout(c, x)
}
func (s *Service) Audits(c context.Context) ([]*models.AuditLog, error) { return s.repo.Audits(c) }
func (s *Service) CreateAudit(c context.Context, actor string, q *dto.CreateAuditRequest) (*models.AuditLog, error) {
	a, e := id(actor, "actor_user_id")
	if e != nil {
		return nil, e
	}
	var entity *uuid.UUID
	if q.EntityID != nil {
		x, e2 := id(*q.EntityID, "entity_id")
		if e2 != nil {
			return nil, e2
		}
		entity = &x
	}
	x := &models.AuditLog{ActorUserID: &a, Action: q.Action, EntityType: q.EntityType, EntityID: entity, BeforeState: q.BeforeState, AfterState: q.AfterState, IPAddress: q.IPAddress, UserAgent: q.UserAgent}
	e = s.repo.CreateAudit(c, x)
	return x, e
}
func (s *Service) DeleteAudit(c context.Context, v string) error {
	x, e := id(v, "audit_id")
	if e != nil {
		return e
	}
	return s.repo.DeleteAudit(c, x)
}
