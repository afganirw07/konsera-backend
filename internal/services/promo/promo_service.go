package promo

import (
	"context"
	"fmt"
	dto "konsera-backend/internal/DTO/promo"
	"konsera-backend/internal/models"
	repo "konsera-backend/internal/repository/promo"
	"strings"

	"github.com/google/uuid"
)

type Service struct{ repo *repo.Repository }

func NewService(r *repo.Repository) *Service { return &Service{repo: r} }
func pid(v, n string) (uuid.UUID, error) {
	x, e := uuid.Parse(v)
	if e != nil {
		return uuid.Nil, fmt.Errorf("invalid %s", n)
	}
	return x, nil
}
func (s *Service) List(c context.Context) ([]*models.PromoCode, error) { return s.repo.List(c) }
func (s *Service) Get(c context.Context, v string) (*models.PromoCode, error) {
	x, e := pid(v, "promo_id")
	if e != nil {
		return nil, e
	}
	return s.repo.Get(c, x)
}
func (s *Service) Create(c context.Context, q *dto.CreatePromoCodeRequest) (*models.PromoCode, error) {
	if q.ValidUntil.Before(q.ValidFrom) || q.ValidUntil.Equal(q.ValidFrom) {
		return nil, fmt.Errorf("valid_until must be after valid_from")
	}
	if q.DiscountType == "percentage" && q.DiscountValue > 100 {
		return nil, fmt.Errorf("percentage discount cannot exceed 100")
	}
	var event *uuid.UUID
	if q.EventID != nil && strings.TrimSpace(*q.EventID) != "" {
		x, e := pid(*q.EventID, "event_id")
		if e != nil {
			return nil, e
		}
		event = &x
	}
	x := &models.PromoCode{Code: strings.ToUpper(strings.TrimSpace(q.Code)), EventID: event, DiscountType: q.DiscountType, DiscountValue: q.DiscountValue, MaxDiscountAmount: q.MaxDiscountAmount, UsageLimit: q.UsageLimit, PerUserLimit: q.PerUserLimit, ValidFrom: q.ValidFrom, ValidUntil: q.ValidUntil, IsActive: true}
	e := s.repo.Create(c, x)
	return x, e
}
func (s *Service) Update(c context.Context, v string, q *dto.UpdatePromoCodeRequest) (*models.PromoCode, error) {
	x, e := pid(v, "promo_id")
	if e != nil {
		return nil, e
	}
	if q.DiscountType != nil && *q.DiscountType == "percentage" && q.DiscountValue != nil && *q.DiscountValue > 100 {
		return nil, fmt.Errorf("percentage discount cannot exceed 100")
	}
	if q.ValidFrom != nil && q.ValidUntil != nil && !q.ValidUntil.After(*q.ValidFrom) {
		return nil, fmt.Errorf("valid_until must be after valid_from")
	}
	return s.repo.Update(c, x, q)
}
func (s *Service) Delete(c context.Context, v string) error {
	x, e := pid(v, "promo_id")
	if e != nil {
		return e
	}
	return s.repo.Delete(c, x)
}
func (s *Service) Apply(c context.Context, user, booking, code string) (*models.PromoCodeUsage, error) {
	u, e := pid(user, "user_id")
	if e != nil {
		return nil, e
	}
	b, e := pid(booking, "booking_id")
	if e != nil {
		return nil, e
	}
	return s.repo.Apply(c, u, b, strings.TrimSpace(code))
}
