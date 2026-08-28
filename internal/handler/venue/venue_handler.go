package venue

import (
	"database/sql"
	"errors"
	"net/http"

	dto "konsera-backend/internal/DTO/venue"
	"konsera-backend/internal/helpers"
	service "konsera-backend/internal/services/venue"

	"github.com/gin-gonic/gin"
)

type VenueHandler struct{ service *service.VenueService }

func NewVenueHandler(service *service.VenueService) *VenueHandler { return &VenueHandler{service: service} }

func respondError(c *gin.Context, err error) {
	status := http.StatusBadRequest
	if errors.Is(err, sql.ErrNoRows) || err.Error() == "venue not found" || err.Error() == "venue section not found" || err.Error() == "seat not found" {
		status = http.StatusNotFound
	}
	helpers.Error(c, status, err.Error(), nil)
}

// @Summary Create venue
// @Tags Venues
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param venue body venue.CreateVenueRequest true "Venue data"
// @Success 201 {object} models.Venue
// @Failure 400 {object} helpers.APIResponse
// @Router /venues [post]
func (h *VenueHandler) Create(c *gin.Context) {
	var req dto.CreateVenueRequest
	if err := c.ShouldBindJSON(&req); err != nil { helpers.Error(c, 400, "Invalid request", err.Error()); return }
	item, err := h.service.Create(c, &req)
	if err != nil { respondError(c, err); return }
	helpers.Success(c, http.StatusCreated, "Venue created successfully", item)
}

// @Summary List venues
// @Tags Venues
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Venue
// @Router /venues [get]
func (h *VenueHandler) List(c *gin.Context) {
	items, err := h.service.List(c)
	if err != nil { respondError(c, err); return }
	helpers.Success(c, 200, "Venues retrieved successfully", items)
}

// @Summary Get venue
// @Tags Venues
// @Produce json
// @Security BearerAuth
// @Param venue_id path string true "Venue ID"
// @Success 200 {object} models.Venue
// @Router /venues/{venue_id} [get]
func (h *VenueHandler) Get(c *gin.Context) {
	item, err := h.service.Get(c, c.Param("venue_id"))
	if err != nil { respondError(c, err); return }
	helpers.Success(c, 200, "Venue retrieved successfully", item)
}

// @Summary Update venue
// @Tags Venues
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param venue_id path string true "Venue ID"
// @Param venue body venue.UpdateVenueRequest true "Venue changes"
// @Success 200 {object} models.Venue
// @Router /venues/{venue_id} [put]
func (h *VenueHandler) Update(c *gin.Context) {
	var req dto.UpdateVenueRequest
	if err := c.ShouldBindJSON(&req); err != nil { helpers.Error(c, 400, "Invalid request", err.Error()); return }
	item, err := h.service.Update(c, c.Param("venue_id"), &req)
	if err != nil { respondError(c, err); return }
	helpers.Success(c, 200, "Venue updated successfully", item)
}

// @Summary Delete venue
// @Tags Venues
// @Produce json
// @Security BearerAuth
// @Param venue_id path string true "Venue ID"
// @Success 204 {object} helpers.APIResponse
// @Router /venues/{venue_id} [delete]
func (h *VenueHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c, c.Param("venue_id")); err != nil { respondError(c, err); return }
	helpers.Message(c, 204, "Venue deleted successfully")
}

// @Summary Create venue section
// @Tags Venue Sections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param venue_id path string true "Venue ID"
// @Param section body venue.CreateVenueSectionRequest true "Section data"
// @Success 201 {object} models.VenueSection
// @Router /venues/{venue_id}/sections [post]
func (h *VenueHandler) CreateSection(c *gin.Context) {
	var req dto.CreateVenueSectionRequest
	if err := c.ShouldBindJSON(&req); err != nil { helpers.Error(c, 400, "Invalid request", err.Error()); return }
	item, err := h.service.CreateSection(c, c.Param("venue_id"), &req)
	if err != nil { respondError(c, err); return }
	helpers.Success(c, 201, "Venue section created successfully", item)
}

// @Summary List venue sections
// @Tags Venue Sections
// @Produce json
// @Security BearerAuth
// @Param venue_id path string true "Venue ID"
// @Success 200 {array} models.VenueSection
// @Router /venues/{venue_id}/sections [get]
func (h *VenueHandler) ListSections(c *gin.Context) {
	items, err := h.service.ListSections(c, c.Param("venue_id"))
	if err != nil { respondError(c, err); return }
	helpers.Success(c, 200, "Venue sections retrieved successfully", items)
}

// @Summary Get venue section
// @Tags Venue Sections
// @Produce json
// @Security BearerAuth
// @Param venue_id path string true "Venue ID"
// @Param section_id path string true "Section ID"
// @Success 200 {object} models.VenueSection
// @Router /venues/{venue_id}/sections/{section_id} [get]
func (h *VenueHandler) GetSection(c *gin.Context) {
	item, err := h.service.GetSection(c, c.Param("venue_id"), c.Param("section_id"))
	if err != nil { respondError(c, err); return }
	helpers.Success(c, 200, "Venue section retrieved successfully", item)
}

// @Summary Update venue section
// @Tags Venue Sections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param venue_id path string true "Venue ID"
// @Param section_id path string true "Section ID"
// @Param section body venue.UpdateVenueSectionRequest true "Section changes"
// @Success 200 {object} models.VenueSection
// @Router /venues/{venue_id}/sections/{section_id} [put]
func (h *VenueHandler) UpdateSection(c *gin.Context) {
	var req dto.UpdateVenueSectionRequest
	if err := c.ShouldBindJSON(&req); err != nil { helpers.Error(c, 400, "Invalid request", err.Error()); return }
	item, err := h.service.UpdateSection(c, c.Param("venue_id"), c.Param("section_id"), &req)
	if err != nil { respondError(c, err); return }
	helpers.Success(c, 200, "Venue section updated successfully", item)
}

// @Summary Delete venue section
// @Tags Venue Sections
// @Produce json
// @Security BearerAuth
// @Param venue_id path string true "Venue ID"
// @Param section_id path string true "Section ID"
// @Success 204 {object} helpers.APIResponse
// @Router /venues/{venue_id}/sections/{section_id} [delete]
func (h *VenueHandler) DeleteSection(c *gin.Context) {
	if err := h.service.DeleteSection(c, c.Param("venue_id"), c.Param("section_id")); err != nil { respondError(c, err); return }
	helpers.Message(c, 204, "Venue section deleted successfully")
}

// @Summary Create seat
// @Tags Seats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param venue_id path string true "Venue ID"
// @Param section_id path string true "Section ID"
// @Param seat body venue.CreateSeatRequest true "Seat data"
// @Success 201 {object} models.Seat
// @Router /venues/{venue_id}/sections/{section_id}/seats [post]
func (h *VenueHandler) CreateSeat(c *gin.Context) {
	var req dto.CreateSeatRequest
	if err := c.ShouldBindJSON(&req); err != nil { helpers.Error(c, 400, "Invalid request", err.Error()); return }
	item, err := h.service.CreateSeat(c, c.Param("venue_id"), c.Param("section_id"), &req)
	if err != nil { respondError(c, err); return }
	helpers.Success(c, 201, "Seat created successfully", item)
}

// @Summary List seats
// @Tags Seats
// @Produce json
// @Security BearerAuth
// @Param venue_id path string true "Venue ID"
// @Param section_id path string true "Section ID"
// @Success 200 {array} models.Seat
// @Router /venues/{venue_id}/sections/{section_id}/seats [get]
func (h *VenueHandler) ListSeats(c *gin.Context) {
	items, err := h.service.ListSeats(c, c.Param("venue_id"), c.Param("section_id"))
	if err != nil { respondError(c, err); return }
	helpers.Success(c, 200, "Seats retrieved successfully", items)
}

// @Summary Get seat
// @Tags Seats
// @Produce json
// @Security BearerAuth
// @Param venue_id path string true "Venue ID"
// @Param section_id path string true "Section ID"
// @Param seat_id path string true "Seat ID"
// @Success 200 {object} models.Seat
// @Router /venues/{venue_id}/sections/{section_id}/seats/{seat_id} [get]
func (h *VenueHandler) GetSeat(c *gin.Context) {
	item, err := h.service.GetSeat(c, c.Param("venue_id"), c.Param("section_id"), c.Param("seat_id"))
	if err != nil { respondError(c, err); return }
	helpers.Success(c, 200, "Seat retrieved successfully", item)
}

// @Summary Update seat
// @Tags Seats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param venue_id path string true "Venue ID"
// @Param section_id path string true "Section ID"
// @Param seat_id path string true "Seat ID"
// @Param seat body venue.UpdateSeatRequest true "Seat changes"
// @Success 200 {object} models.Seat
// @Router /venues/{venue_id}/sections/{section_id}/seats/{seat_id} [put]
func (h *VenueHandler) UpdateSeat(c *gin.Context) {
	var req dto.UpdateSeatRequest
	if err := c.ShouldBindJSON(&req); err != nil { helpers.Error(c, 400, "Invalid request", err.Error()); return }
	item, err := h.service.UpdateSeat(c, c.Param("venue_id"), c.Param("section_id"), c.Param("seat_id"), &req)
	if err != nil { respondError(c, err); return }
	helpers.Success(c, 200, "Seat updated successfully", item)
}

// @Summary Delete seat
// @Tags Seats
// @Produce json
// @Security BearerAuth
// @Param venue_id path string true "Venue ID"
// @Param section_id path string true "Section ID"
// @Param seat_id path string true "Seat ID"
// @Success 204 {object} helpers.APIResponse
// @Router /venues/{venue_id}/sections/{section_id}/seats/{seat_id} [delete]
func (h *VenueHandler) DeleteSeat(c *gin.Context) {
	if err := h.service.DeleteSeat(c, c.Param("venue_id"), c.Param("section_id"), c.Param("seat_id")); err != nil { respondError(c, err); return }
	helpers.Message(c, 204, "Seat deleted successfully")
}
