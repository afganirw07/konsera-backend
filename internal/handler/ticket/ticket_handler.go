package ticket

import (
	"database/sql"
	"errors"
	dto "konsera-backend/internal/DTO/ticket"
	"konsera-backend/internal/helpers"
	service "konsera-backend/internal/services/ticket"

	"github.com/gin-gonic/gin"
)

type TicketHandler struct{ s *service.TicketService }

func NewTicketHandler(s *service.TicketService) *TicketHandler { return &TicketHandler{s: s} }
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

// @Summary Create ticket tier
// @Tags Ticket Tiers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param event_id path string true "Event ID"
// @Param tier body ticket.CreateTicketTierRequest true "Ticket tier"
// @Success 201 {object} models.TicketTier
// @Router /events/{event_id}/ticket-tiers [post]
func (h *TicketHandler) CreateTier(c *gin.Context) {
	q := dto.CreateTicketTierRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.CreateTier(c, c.Param("event_id"), &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 201, "Ticket tier created successfully", x)
}

// @Summary List ticket tiers
// @Tags Ticket Tiers
// @Produce json
// @Security BearerAuth
// @Param event_id path string true "Event ID"
// @Success 200 {array} models.TicketTier
// @Router /events/{event_id}/ticket-tiers [get]
func (h *TicketHandler) ListTiers(c *gin.Context) {
	x, e := h.s.Tiers(c, c.Param("event_id"))
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Ticket tiers retrieved successfully", x)
}

// @Summary Get ticket tier
// @Tags Ticket Tiers
// @Produce json
// @Security BearerAuth
// @Param event_id path string true "Event ID"
// @Param tier_id path string true "Ticket tier ID"
// @Success 200 {object} models.TicketTier
// @Router /events/{event_id}/ticket-tiers/{tier_id} [get]
func (h *TicketHandler) GetTier(c *gin.Context) {
	x, e := h.s.GetTier(c, c.Param("event_id"), c.Param("tier_id"))
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Ticket tier retrieved successfully", x)
}

// @Summary Update ticket tier
// @Tags Ticket Tiers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param event_id path string true "Event ID"
// @Param tier_id path string true "Ticket tier ID"
// @Param tier body ticket.UpdateTicketTierRequest true "Ticket tier changes"
// @Success 200 {object} models.TicketTier
// @Router /events/{event_id}/ticket-tiers/{tier_id} [put]
func (h *TicketHandler) UpdateTier(c *gin.Context) {
	q := dto.UpdateTicketTierRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.UpdateTier(c, c.Param("event_id"), c.Param("tier_id"), &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Ticket tier updated successfully", x)
}

// @Summary Delete ticket tier
// @Tags Ticket Tiers
// @Produce json
// @Security BearerAuth
// @Param event_id path string true "Event ID"
// @Param tier_id path string true "Ticket tier ID"
// @Success 204 {object} helpers.APIResponse
// @Router /events/{event_id}/ticket-tiers/{tier_id} [delete]
func (h *TicketHandler) DeleteTier(c *gin.Context) {
	if e := h.s.DeleteTier(c, c.Param("event_id"), c.Param("tier_id")); e != nil {
		fail(c, e)
		return
	}
	helpers.Message(c, 204, "Ticket tier deleted successfully")
}

// @Summary Create inventory
// @Tags Ticket Inventories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tier_id path string true "Ticket tier ID"
// @Param inventory body ticket.CreateInventoryRequest true "Inventory"
// @Success 201 {object} models.TicketInventory
// @Router /ticket-tiers/{tier_id}/inventory [post]
func (h *TicketHandler) CreateInventory(c *gin.Context) {
	q := dto.CreateInventoryRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.CreateInventory(c, c.Param("tier_id"), &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 201, "Inventory created successfully", x)
}
func (h *TicketHandler) ListInventory(c *gin.Context) {
	x, e := h.s.Inventories(c, c.Param("tier_id"))
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Inventory retrieved successfully", x)
}
func (h *TicketHandler) GetInventory(c *gin.Context) {
	x, e := h.s.GetInventory(c, c.Param("tier_id"), c.Param("event_session_id"))
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Inventory retrieved successfully", x)
}
func (h *TicketHandler) UpdateInventory(c *gin.Context) {
	q := dto.UpdateInventoryRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.UpdateInventory(c, c.Param("tier_id"), c.Param("event_session_id"), &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Inventory updated successfully", x)
}
func (h *TicketHandler) DeleteInventory(c *gin.Context) {
	if e := h.s.DeleteInventory(c, c.Param("tier_id"), c.Param("event_session_id")); e != nil {
		fail(c, e)
		return
	}
	helpers.Message(c, 204, "Inventory deleted successfully")
}
