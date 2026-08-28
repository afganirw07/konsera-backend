package promo

import (
	"database/sql"
	"errors"
	dto "konsera-backend/internal/DTO/promo"
	"konsera-backend/internal/helpers"
	svc "konsera-backend/internal/services/promo"

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
func bind(c *gin.Context, v any) bool {
	if e := c.ShouldBindJSON(v); e != nil {
		helpers.Error(c, 400, "Invalid request", e.Error())
		return false
	}
	return true
}

// @Summary List promo codes
// @Tags Promo Codes
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.PromoCode
// @Router /promo-codes [get]
func (h *Handler) List(c *gin.Context) {
	x, e := h.s.List(c)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Promo codes retrieved successfully", x)
}

// @Summary Get promo code
// @Tags Promo Codes
// @Produce json
// @Security BearerAuth
// @Param promo_id path string true "Promo ID"
// @Success 200 {object} models.PromoCode
// @Router /promo-codes/{promo_id} [get]
func (h *Handler) Get(c *gin.Context) {
	x, e := h.s.Get(c, c.Param("promo_id"))
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Promo code retrieved successfully", x)
}

// @Summary Create promo code
// @Tags Promo Codes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param promo body promo.CreatePromoCodeRequest true "Promo code"
// @Success 201 {object} models.PromoCode
// @Router /promo-codes [post]
func (h *Handler) Create(c *gin.Context) {
	q := dto.CreatePromoCodeRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.Create(c, &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 201, "Promo code created successfully", x)
}

// @Summary Update promo code
// @Tags Promo Codes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param promo_id path string true "Promo ID"
// @Param promo body promo.UpdatePromoCodeRequest true "Promo changes"
// @Success 200 {object} models.PromoCode
// @Router /promo-codes/{promo_id} [put]
func (h *Handler) Update(c *gin.Context) {
	q := dto.UpdatePromoCodeRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.Update(c, c.Param("promo_id"), &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Promo code updated successfully", x)
}

// @Summary Delete promo code
// @Tags Promo Codes
// @Produce json
// @Security BearerAuth
// @Param promo_id path string true "Promo ID"
// @Success 204 {object} helpers.APIResponse
// @Router /promo-codes/{promo_id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	if e := h.s.Delete(c, c.Param("promo_id")); e != nil {
		fail(c, e)
		return
	}
	helpers.Message(c, 204, "Promo code deleted successfully")
}

// @Summary Apply promo code
// @Tags Promo Codes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param booking_id path string true "Booking ID"
// @Param promo body promo.ApplyPromoRequest true "Promo code"
// @Success 200 {object} models.PromoCodeUsage
// @Router /checkout/bookings/{booking_id}/promo [post]
func (h *Handler) Apply(c *gin.Context) {
	q := dto.ApplyPromoRequest{}
	if !bind(c, &q) {
		return
	}
	user, _ := c.Get("user_id")
	x, e := h.s.Apply(c, user.(string), c.Param("booking_id"), q.Code)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Promo code applied successfully", x)
}
