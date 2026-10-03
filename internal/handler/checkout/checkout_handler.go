package checkout

import (
	"database/sql"
	"errors"
	dto "konsera-backend/internal/DTO/checkout"
	"konsera-backend/internal/helpers"
	svc "konsera-backend/internal/services/checkout"

	"github.com/gin-gonic/gin"
)

type Handler struct{ s *svc.Service }

func NewHandler(s *svc.Service) *Handler { return &Handler{s: s} }
func fail(c *gin.Context, e error) {
	st := 400
	if errors.Is(e, sql.ErrNoRows) {
		st = 404
	}
	helpers.Error(c, st, e.Error(), nil)
}
func userID(c *gin.Context) string { v, _ := c.Get("user_id"); return v.(string) }
func bind(c *gin.Context, v any) bool {
	if e := c.ShouldBindJSON(v); e != nil {
		helpers.Error(c, 400, "Invalid request", e.Error())
		return false
	}
	return true
}
func (h *Handler) Carts(c *gin.Context) {
	x, e := h.s.Carts(c, userID(c))
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Carts retrieved successfully", x)
}
func (h *Handler) Cart(c *gin.Context) {
	x, e := h.s.Cart(c, userID(c), c.Param("cart_id"))
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Cart retrieved successfully", x)
}
func (h *Handler) CreateCart(c *gin.Context) {
	q := dto.CreateCartRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.CreateCart(c, userID(c), &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 201, "Cart created successfully", x)
}
func (h *Handler) UpdateCart(c *gin.Context) {
	q := dto.UpdateCartRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.UpdateCart(c, userID(c), c.Param("cart_id"), &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Cart updated successfully", x)
}
func (h *Handler) DeleteCart(c *gin.Context) {
	if e := h.s.DeleteCart(c, userID(c), c.Param("cart_id")); e != nil {
		fail(c, e)
		return
	}
	helpers.Message(c, 204, "Cart deleted successfully")
}
func (h *Handler) CreateBooking(c *gin.Context) {
	q := dto.CreateBookingRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.Checkout(c, userID(c), &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 201, "Booking created successfully", x)
}
func (h *Handler) Bookings(c *gin.Context) {
	x, e := h.s.Bookings(c, userID(c))
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Bookings retrieved successfully", x)
}
func (h *Handler) Booking(c *gin.Context) {
	x, e := h.s.Booking(c, userID(c), c.Param("booking_id"))
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Booking retrieved successfully", x)
}
func (h *Handler) UpdateBooking(c *gin.Context) {
	q := dto.UpdateBookingRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.UpdateBooking(c, userID(c), c.Param("booking_id"), &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Booking updated successfully", x)
}
func (h *Handler) DeleteBooking(c *gin.Context) {
	if e := h.s.DeleteBooking(c, userID(c), c.Param("booking_id")); e != nil {
		fail(c, e)
		return
	}
	helpers.Message(c, 204, "Booking cancelled successfully")
}
func (h *Handler) Methods(c *gin.Context) {
	x, e := h.s.Methods(c)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Payment methods retrieved successfully", x)
}
func (h *Handler) CreateMethod(c *gin.Context) {
	q := dto.CreatePaymentMethodRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.CreateMethod(c, &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 201, "Payment method created successfully", x)
}
func (h *Handler) UpdateMethod(c *gin.Context) {
	q := dto.UpdatePaymentMethodRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.UpdateMethod(c, c.Param("payment_method_id"), &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Payment method updated successfully", x)
}
func (h *Handler) DeleteMethod(c *gin.Context) {
	if e := h.s.DeleteMethod(c, c.Param("payment_method_id")); e != nil {
		fail(c, e)
		return
	}
	helpers.Message(c, 204, "Payment method deleted successfully")
}
func (h *Handler) Payments(c *gin.Context) {
	x, e := h.s.Payments(c, userID(c), c.Query("booking_id"))
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Payments retrieved successfully", x)
}
func (h *Handler) CreatePayment(c *gin.Context) {
	q := dto.CreatePaymentRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.CreatePayment(c, userID(c), &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 201, "Payment created successfully", x)
}
func (h *Handler) UpdatePayment(c *gin.Context) {
	q := dto.UpdatePaymentRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.UpdatePayment(c, userID(c), c.Param("payment_id"), &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Payment updated successfully", x)
}
func (h *Handler) Refunds(c *gin.Context) {
	x, e := h.s.Refunds(c, userID(c))
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Refunds retrieved successfully", x)
}
func (h *Handler) CreateRefund(c *gin.Context) {
	q := dto.CreateRefundRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.CreateRefund(c, userID(c), &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 201, "Refund requested successfully", x)
}
func (h *Handler) UpdateRefund(c *gin.Context) {
	q := dto.UpdateRefundRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.UpdateRefund(c, c.Param("refund_id"), &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Refund updated successfully", x)
}
func (h *Handler) DeleteRefund(c *gin.Context) {
	if e := h.s.DeleteRefund(c, c.Param("refund_id")); e != nil {
		fail(c, e)
		return
	}
	helpers.Message(c, 204, "Refund deleted successfully")
}
