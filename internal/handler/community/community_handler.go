package community

import (
	"database/sql"
	"errors"
	dto "konsera-backend/internal/DTO/community"
	"konsera-backend/internal/helpers"
	svc "konsera-backend/internal/services/community"

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
func uid(c *gin.Context) string { v, _ := c.Get("user_id"); return v.(string) }
func bind(c *gin.Context, v any) bool {
	if e := c.ShouldBindJSON(v); e != nil {
		helpers.Error(c, 400, "Invalid request", e.Error())
		return false
	}
	return true
}
func (h *Handler) Reviews(c *gin.Context) {
	x, e := h.s.Reviews(c, c.Query("event_id"))
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Reviews retrieved successfully", x)
}
func (h *Handler) GetReview(c *gin.Context) {
	x, e := h.s.Review(c, uid(c), c.Param("review_id"))
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Review retrieved successfully", x)
}
func (h *Handler) CreateReview(c *gin.Context) {
	q := dto.CreateReviewRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.CreateReview(c, uid(c), &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 201, "Review created successfully", x)
}
func (h *Handler) UpdateReview(c *gin.Context) {
	q := dto.UpdateReviewRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.UpdateReview(c, uid(c), c.Param("review_id"), &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Review updated successfully", x)
}
func (h *Handler) DeleteReview(c *gin.Context) {
	if e := h.s.DeleteReview(c, uid(c), c.Param("review_id")); e != nil {
		fail(c, e)
		return
	}
	helpers.Message(c, 204, "Review deleted successfully")
}
func (h *Handler) Favorites(c *gin.Context) {
	x, e := h.s.Favorites(c, uid(c))
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Favorites retrieved successfully", x)
}
func (h *Handler) CreateFavorite(c *gin.Context) {
	q := dto.CreateFavoriteRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.CreateFavorite(c, uid(c), &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 201, "Favorite created successfully", x)
}
func (h *Handler) DeleteFavorite(c *gin.Context) {
	if e := h.s.DeleteFavorite(c, uid(c), c.Param("event_id")); e != nil {
		fail(c, e)
		return
	}
	helpers.Message(c, 204, "Favorite deleted successfully")
}
func (h *Handler) Commissions(c *gin.Context) {
	x, e := h.s.Commissions(c)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Commission settings retrieved successfully", x)
}
func (h *Handler) CreateCommission(c *gin.Context) {
	q := dto.CreateCommissionRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.CreateCommission(c, &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 201, "Commission setting created successfully", x)
}
func (h *Handler) UpdateCommission(c *gin.Context) {
	q := dto.UpdateCommissionRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.UpdateCommission(c, c.Param("commission_id"), &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Commission setting updated successfully", x)
}
func (h *Handler) DeleteCommission(c *gin.Context) {
	if e := h.s.DeleteCommission(c, c.Param("commission_id")); e != nil {
		fail(c, e)
		return
	}
	helpers.Message(c, 204, "Commission setting deleted successfully")
}
func (h *Handler) Payouts(c *gin.Context) {
	x, e := h.s.Payouts(c)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Payouts retrieved successfully", x)
}
func (h *Handler) CreatePayout(c *gin.Context) {
	q := dto.CreatePayoutRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.CreatePayout(c, &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 201, "Payout created successfully", x)
}
func (h *Handler) UpdatePayout(c *gin.Context) {
	q := dto.UpdatePayoutRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.UpdatePayout(c, c.Param("payout_id"), &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Payout updated successfully", x)
}
func (h *Handler) DeletePayout(c *gin.Context) {
	if e := h.s.DeletePayout(c, c.Param("payout_id")); e != nil {
		fail(c, e)
		return
	}
	helpers.Message(c, 204, "Payout deleted successfully")
}
func (h *Handler) Audits(c *gin.Context) {
	x, e := h.s.Audits(c)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Audit logs retrieved successfully", x)
}
func (h *Handler) CreateAudit(c *gin.Context) {
	q := dto.CreateAuditRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.CreateAudit(c, uid(c), &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 201, "Audit log created successfully", x)
}
func (h *Handler) DeleteAudit(c *gin.Context) {
	if e := h.s.DeleteAudit(c, c.Param("audit_id")); e != nil {
		fail(c, e)
		return
	}
	helpers.Message(c, 204, "Audit log deleted successfully")
}
