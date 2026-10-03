package checkout

import (
	"context"
	"fmt"
	dto "konsera-backend/internal/DTO/checkout"
	"konsera-backend/internal/models"
	repo "konsera-backend/internal/repository/checkout"
	"time"

	"github.com/google/uuid"
)

type Service struct{ repo *repo.Repository }

func NewService(r *repo.Repository) *Service { return &Service{repo: r} }
func uid(v, n string) (uuid.UUID, error) {
	x, e := uuid.Parse(v)
	if e != nil {
		return uuid.Nil, fmt.Errorf("invalid %s", n)
	}
	return x, nil
}
func (s *Service) Carts(c context.Context, u string) ([]*models.Cart, error) {
	id, e := uid(u, "user_id")
	if e != nil {
		return nil, e
	}
	return s.repo.Carts(c, id)
}
func (s *Service) Cart(c context.Context, u, v string) (*models.Cart, error) {
	a, e := uid(u, "user_id")
	if e != nil {
		return nil, e
	}
	b, e := uid(v, "cart_id")
	if e != nil {
		return nil, e
	}
	return s.repo.Cart(c, a, b)
}
func (s *Service) CreateCart(c context.Context, u string, q *dto.CreateCartRequest) (*models.Cart, error) {
	user, e := uid(u, "user_id")
	if e != nil {
		return nil, e
	}
	tier, e := uid(q.TicketTierID, "ticket_tier_id")
	if e != nil {
		return nil, e
	}
	session, e := uid(q.EventSessionID, "event_session_id")
	if e != nil {
		return nil, e
	}
	if q.Quantity > 6 {
		return nil, fmt.Errorf("quantity cannot exceed 6 tickets")
	}
	seat := (*uuid.UUID)(nil)
	if q.SeatID != nil {
		x, e := uid(*q.SeatID, "seat_id")
		if e != nil {
			return nil, e
		}
		seat = &x
	}
	x := &models.Cart{UserID: user, TicketTierID: tier, EventSessionID: session, SeatID: seat, Quantity: q.Quantity, HeldUntil: time.Now().Add(15 * time.Minute)}
	e = s.repo.CreateCart(c, x)
	return x, e
}
func (s *Service) UpdateCart(c context.Context, u, v string, q *dto.UpdateCartRequest) (*models.Cart, error) {
	a, e := uid(u, "user_id")
	if e != nil {
		return nil, e
	}
	b, e := uid(v, "cart_id")
	if e != nil {
		return nil, e
	}
	if q.Quantity == nil {
		return nil, fmt.Errorf("quantity is required")
	}
	return s.repo.UpdateCart(c, a, b, *q.Quantity)
}
func (s *Service) DeleteCart(c context.Context, u, v string) error {
	a, e := uid(u, "user_id")
	if e != nil {
		return e
	}
	b, e := uid(v, "cart_id")
	if e != nil {
		return e
	}
	return s.repo.DeleteCart(c, a, b)
}
func (s *Service) Checkout(c context.Context, u string, q *dto.CreateBookingRequest) (*models.Booking, error) {
	user, e := uid(u, "user_id")
	if e != nil {
		return nil, e
	}
	event, e := uid(q.EventID, "event_id")
	if e != nil {
		return nil, e
	}
	return s.repo.Checkout(c, user, event)
}
func (s *Service) Bookings(c context.Context, u string) ([]*models.Booking, error) {
	id, e := uid(u, "user_id")
	if e != nil {
		return nil, e
	}
	return s.repo.Bookings(c, id)
}
func (s *Service) Booking(c context.Context, u, v string) (*models.Booking, error) {
	a, e := uid(u, "user_id")
	if e != nil {
		return nil, e
	}
	b, e := uid(v, "booking_id")
	if e != nil {
		return nil, e
	}
	return s.repo.Booking(c, a, b)
}
func (s *Service) UpdateBooking(c context.Context, u, v string, q *dto.UpdateBookingRequest) (*models.Booking, error) {
	a, e := uid(u, "user_id")
	if e != nil {
		return nil, e
	}
	b, e := uid(v, "booking_id")
	if e != nil {
		return nil, e
	}
	if q.Status == nil {
		return nil, fmt.Errorf("status is required")
	}
	return s.repo.UpdateBooking(c, a, b, *q.Status)
}
func (s *Service) DeleteBooking(c context.Context, u, v string) error {
	a, e := uid(u, "user_id")
	if e != nil {
		return e
	}
	b, e := uid(v, "booking_id")
	if e != nil {
		return e
	}
	return s.repo.DeleteBooking(c, a, b)
}
func (s *Service) Methods(c context.Context) ([]*models.PaymentMethod, error) {
	return s.repo.Methods(c)
}
func (s *Service) CreateMethod(c context.Context, q *dto.CreatePaymentMethodRequest) (*models.PaymentMethod, error) {
	x := &models.PaymentMethod{Name: q.Name, Type: q.Type, Provider: q.Provider, IconURL: q.IconURL}
	e := s.repo.CreateMethod(c, x)
	return x, e
}
func (s *Service) UpdateMethod(c context.Context, v string, q *dto.UpdatePaymentMethodRequest) (*models.PaymentMethod, error) {
	id, e := uid(v, "payment_method_id")
	if e != nil {
		return nil, e
	}
	return s.repo.UpdateMethod(c, id, q)
}
func (s *Service) DeleteMethod(c context.Context, v string) error {
	id, e := uid(v, "payment_method_id")
	if e != nil {
		return e
	}
	return s.repo.DeleteMethod(c, id)
}
func (s *Service) Payments(c context.Context, u, v string) ([]*models.Payment, error) {
	a, e := uid(u, "user_id")
	if e != nil {
		return nil, e
	}
	var b *uuid.UUID
	if v != "" {
		id, e := uid(v, "booking_id")
		if e != nil {
			return nil, e
		}
		b = &id
	}
	return s.repo.Payments(c, a, b)
}
func (s *Service) CreatePayment(c context.Context, q *dto.CreatePaymentRequest) (*models.Payment, error) {
	b, e := uid(q.BookingID, "booking_id")
	if e != nil {
		return nil, e
	}
	m, e := uid(q.PaymentMethodID, "payment_method_id")
	if e != nil {
		return nil, e
	}
	x := &models.Payment{BookingID: b, PaymentMethodID: m, Amount: q.Amount, ExpiresAt: q.ExpiresAt}
	e = s.repo.CreatePayment(c, x)
	return x, e
}
func (s *Service) UpdatePayment(c context.Context, u, v string, q *dto.UpdatePaymentRequest) (*models.Payment, error) {
	a, e := uid(u, "user_id")
	if e != nil {
		return nil, e
	}
	b, e := uid(v, "payment_id")
	if e != nil {
		return nil, e
	}
	return s.repo.UpdatePayment(c, a, b, q)
}
func (s *Service) Refunds(c context.Context, u string) ([]*models.Refund, error) {
	id, e := uid(u, "user_id")
	if e != nil {
		return nil, e
	}
	return s.repo.Refunds(c, id)
}
func (s *Service) CreateRefund(c context.Context, u string, q *dto.CreateRefundRequest) (*models.Refund, error) {
	user, e := uid(u, "user_id")
	if e != nil {
		return nil, e
	}
	b, e := uid(q.BookingID, "booking_id")
	if e != nil {
		return nil, e
	}
	p, e := uid(q.PaymentID, "payment_id")
	if e != nil {
		return nil, e
	}
	x := &models.Refund{BookingID: b, PaymentID: p, Amount: q.Amount, Reason: q.Reason, RequestedBy: user}
	e = s.repo.CreateRefund(c, x)
	return x, e
}
func (s *Service) UpdateRefund(c context.Context, v string, q *dto.UpdateRefundRequest) (*models.Refund, error) {
	id, e := uid(v, "refund_id")
	if e != nil {
		return nil, e
	}
	return s.repo.UpdateRefund(c, id, q)
}
func (s *Service) DeleteRefund(c context.Context, v string) error {
	id, e := uid(v, "refund_id")
	if e != nil {
		return e
	}
	return s.repo.DeleteRefund(c, id)
}
