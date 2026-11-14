package checkout

import (
	"context"
	"encoding/json"
	"fmt"
	dto "konsera-backend/internal/DTO/checkout"
	"konsera-backend/internal/models"
	repo "konsera-backend/internal/repository/checkout"
	"konsera-backend/internal/services/imageupload"
	"konsera-backend/internal/services/midtrans"
	"mime/multipart"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo     *repo.Repository
	uploader imageupload.Uploader
	midtrans *midtrans.Service
}

func NewService(r *repo.Repository, uploader imageupload.Uploader, gateway *midtrans.Service) *Service {
	return &Service{repo: r, uploader: uploader, midtrans: gateway}
}
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
func (s *Service) Checkout(c context.Context, u string, q *dto.CreateBookingRequest) (*dto.CheckoutResponse, error) {
	user, e := uid(u, "user_id")
	if e != nil {
		return nil, e
	}
	event, e := uid(q.EventID, "event_id")
	if e != nil {
		return nil, e
	}
	if !s.midtrans.IsConfigured() {
		return nil, midtrans.ErrNotConfigured
	}

	booking, e := s.repo.Checkout(c, user, event)
	if e != nil {
		return nil, e
	}

	snap, e := s.midtrans.CreateSnapTransaction(c, midtrans.SnapRequest{
		OrderID:     booking.BookingCode,
		GrossAmount: booking.TotalAmount,
	})
	if e != nil {
		return nil, e
	}

	return &dto.CheckoutResponse{
		Booking:     booking,
		SnapToken:   snap.Token,
		RedirectURL: snap.RedirectURL,
	}, nil
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
func (s *Service) UploadPaymentMethodIcon(c context.Context, v string, file *multipart.FileHeader) (*models.PaymentMethod, error) {
	id, e := uid(v, "payment_method_id")
	if e != nil {
		return nil, e
	}
	url, e := s.uploader.Upload(c, file)
	if e != nil {
		return nil, e
	}
	return s.repo.UpdatePaymentMethodIcon(c, id, url)
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
func (s *Service) CreatePayment(c context.Context, u string, q *dto.CreatePaymentRequest) (*models.Payment, error) {
	user, e := uid(u, "user_id")
	if e != nil {
		return nil, e
	}
	b, e := uid(q.BookingID, "booking_id")
	if e != nil {
		return nil, e
	}
	m, e := uid(q.PaymentMethodID, "payment_method_id")
	if e != nil {
		return nil, e
	}
	x := &models.Payment{BookingID: b, PaymentMethodID: m, ExpiresAt: q.ExpiresAt}
	e = s.repo.CreatePayment(c, user, x)
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
func (s *Service) UpdateRefund(c context.Context, actor, v string, q *dto.UpdateRefundRequest) (*models.Refund, error) {
	id, e := uid(v, "refund_id")
	if e != nil {
		return nil, e
	}
	actorID, e := uid(actor, "user_id")
	if e != nil {
		return nil, e
	}
	if q.Status == nil {
		return nil, fmt.Errorf("status is required")
	}
	if *q.Status != "approved" {
		return s.repo.UpdateRefund(c, id, actorID, q)
	}
	if !s.midtrans.IsConfigured() {
		return nil, midtrans.ErrNotConfigured
	}

	refund, e := s.repo.Refund(c, id)
	if e != nil {
		return nil, e
	}
	if refund.Status != "requested" && refund.Status != "approved" {
		return nil, fmt.Errorf("refund cannot be approved from status %s", refund.Status)
	}

	payment, e := s.repo.Payment(c, refund.PaymentID)
	if e != nil {
		return nil, e
	}
	if payment.Status != "success" && payment.Status != "partially_refunded" {
		return nil, fmt.Errorf("payment is not refundable")
	}
	if refund.Amount > payment.Amount {
		return nil, fmt.Errorf("refund amount exceeds payment amount")
	}
	refunded, e := s.repo.RefundedAmount(c, payment.ID)
	if e != nil {
		return nil, e
	}
	if refunded+refund.Amount > payment.Amount {
		return nil, fmt.Errorf("refund amount exceeds remaining payment amount")
	}

	booking, e := s.repo.BookingByID(c, refund.BookingID)
	if e != nil {
		return nil, e
	}

	reason := ""
	if refund.Reason != nil {
		reason = *refund.Reason
	}
	result, e := s.midtrans.Refund(c, midtrans.RefundRequest{
		OrderID:   booking.BookingCode,
		RefundKey: refund.ID.String(),
		Amount:    refund.Amount,
		Reason:    reason,
	})
	if e != nil {
		return nil, e
	}

	if result.Pending {
		approved := "approved"
		return s.repo.UpdateRefund(c, id, actorID, &dto.UpdateRefundRequest{Status: &approved})
	}
	if !result.Approved {
		return nil, fmt.Errorf("midtrans refund was not approved: %s", result.Message)
	}

	paymentStatus, bookingStatus := "partially_refunded", "partially_refunded"
	if refund.Amount >= payment.Amount {
		paymentStatus, bookingStatus = "refunded", "refunded"
	}

	return s.repo.CompleteRefund(c, refund.ID, payment.ID, refund.BookingID, actorID, paymentStatus, bookingStatus)
}

func (s *Service) HandleMidtransNotification(c context.Context, n *dto.MidtransNotificationRequest) error {
	if n == nil {
		return fmt.Errorf("notification is required")
	}
	if !s.midtrans.VerifySignature(n.OrderID, n.StatusCode, n.GrossAmount, n.SignatureKey) {
		return midtrans.ErrInvalidSignature
	}

	amount, e := strconv.ParseFloat(n.GrossAmount, 64)
	if e != nil || amount < 0 {
		return fmt.Errorf("invalid gross_amount")
	}

	raw, e := json.Marshal(n)
	if e != nil {
		return e
	}

	_, e = s.repo.ApplyPaymentNotification(c, repo.PaymentNotification{
		OrderID:       n.OrderID,
		TransactionID: n.TransactionID,
		Status:        midtrans.MapTransactionStatus(n.TransactionStatus, n.FraudStatus),
		Amount:        amount,
		PaidAt:        midtrans.ParseTransactionTime(n.TransactionTime),
		RawPayload:    string(raw),
	})
	return e
}
func (s *Service) DeleteRefund(c context.Context, v string) error {
	id, e := uid(v, "refund_id")
	if e != nil {
		return e
	}
	return s.repo.DeleteRefund(c, id)
}
