package event

import (
	"database/sql"
	"errors"
	dto "konsera-backend/internal/DTO/event"
	"konsera-backend/internal/helpers"
	service "konsera-backend/internal/services/event"

	"github.com/gin-gonic/gin"
)

type EventHandler struct{ s *service.EventService }

func NewEventHandler(s *service.EventService) *EventHandler { return &EventHandler{s: s} }
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

// @Summary Create category
// @Tags Event Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category body event.CreateCategoryRequest true "Category"
// @Success 201 {object} models.EventCategory
// @Router /event-categories [post]
func (h *EventHandler) CreateCategory(c *gin.Context) {
	q := dto.CreateCategoryRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.CreateCategory(c, &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 201, "Category created successfully", x)
}

// @Summary List categories
// @Tags Event Categories
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.EventCategory
// @Router /event-categories [get]
func (h *EventHandler) Categories(c *gin.Context) {
	x, e := h.s.Categories(c)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Categories retrieved successfully", x)
}

// @Summary Get category
// @Tags Event Categories
// @Produce json
// @Security BearerAuth
// @Param category_id path string true "Category ID"
// @Success 200 {object} models.EventCategory
// @Router /event-categories/{category_id} [get]
func (h *EventHandler) GetCategory(c *gin.Context) {
	x, e := h.s.Category(c, c.Param("category_id"))
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Category retrieved successfully", x)
}

// @Summary Update category
// @Tags Event Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category_id path string true "Category ID"
// @Param category body event.UpdateCategoryRequest true "Category changes"
// @Success 200 {object} models.EventCategory
// @Router /event-categories/{category_id} [put]
func (h *EventHandler) UpdateCategory(c *gin.Context) {
	q := dto.UpdateCategoryRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.UpdateCategory(c, c.Param("category_id"), &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Category updated successfully", x)
}

// @Summary Delete category
// @Tags Event Categories
// @Produce json
// @Security BearerAuth
// @Param category_id path string true "Category ID"
// @Success 204 {object} helpers.APIResponse
// @Router /event-categories/{category_id} [delete]
func (h *EventHandler) DeleteCategory(c *gin.Context) {
	if e := h.s.DeleteCategory(c, c.Param("category_id")); e != nil {
		fail(c, e)
		return
	}
	helpers.Message(c, 204, "Category deleted successfully")
}

// @Summary Create event
// @Tags Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param event body event.CreateEventRequest true "Event"
// @Success 201 {object} models.Event
// @Router /events [post]
func (h *EventHandler) Create(c *gin.Context) {
	q := dto.CreateEventRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.Create(c, &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 201, "Event created successfully", x)
}

// @Summary List events
// @Tags Events
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Event
// @Router /events [get]
func (h *EventHandler) Events(c *gin.Context) {
	x, e := h.s.Events(c)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Events retrieved successfully", x)
}

// @Summary Get event
// @Tags Events
// @Produce json
// @Security BearerAuth
// @Param event_id path string true "Event ID"
// @Success 200 {object} models.Event
// @Router /events/{event_id} [get]
func (h *EventHandler) Get(c *gin.Context) {
	x, e := h.s.Event(c, c.Param("event_id"))
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Event retrieved successfully", x)
}

// @Summary Update event
// @Tags Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param event_id path string true "Event ID"
// @Param event body event.UpdateEventRequest true "Event changes"
// @Success 200 {object} models.Event
// @Router /events/{event_id} [put]
func (h *EventHandler) Update(c *gin.Context) {
	q := dto.UpdateEventRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.Update(c, c.Param("event_id"), &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Event updated successfully", x)
}

// @Summary Delete event
// @Tags Events
// @Produce json
// @Security BearerAuth
// @Param event_id path string true "Event ID"
// @Success 204 {object} helpers.APIResponse
// @Router /events/{event_id} [delete]
func (h *EventHandler) Delete(c *gin.Context) {
	if e := h.s.Delete(c, c.Param("event_id")); e != nil {
		fail(c, e)
		return
	}
	helpers.Message(c, 204, "Event deleted successfully")
}

// @Summary Create event session
// @Tags Event Sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param event_id path string true "Event ID"
// @Param session body event.CreateSessionRequest true "Session"
// @Success 201 {object} models.EventSession
// @Router /events/{event_id}/sessions [post]
func (h *EventHandler) CreateSession(c *gin.Context) {
	q := dto.CreateSessionRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.CreateSession(c, c.Param("event_id"), &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 201, "Session created successfully", x)
}
func (h *EventHandler) Sessions(c *gin.Context) {
	x, e := h.s.Sessions(c, c.Param("event_id"))
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Sessions retrieved successfully", x)
}
func (h *EventHandler) GetSession(c *gin.Context) {
	x, e := h.s.Session(c, c.Param("event_id"), c.Param("session_id"))
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Session retrieved successfully", x)
}
func (h *EventHandler) UpdateSession(c *gin.Context) {
	q := dto.UpdateSessionRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.UpdateSession(c, c.Param("event_id"), c.Param("session_id"), &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Session updated successfully", x)
}
func (h *EventHandler) DeleteSession(c *gin.Context) {
	if e := h.s.DeleteSession(c, c.Param("event_id"), c.Param("session_id")); e != nil {
		fail(c, e)
		return
	}
	helpers.Message(c, 204, "Session deleted successfully")
}

// @Summary Create artist
// @Tags Artists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param artist body event.CreateArtistRequest true "Artist"
// @Success 201 {object} models.Artist
// @Router /artists [post]
func (h *EventHandler) CreateArtist(c *gin.Context) {
	q := dto.CreateArtistRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.CreateArtist(c, &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 201, "Artist created successfully", x)
}
func (h *EventHandler) Artists(c *gin.Context) {
	x, e := h.s.Artists(c)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Artists retrieved successfully", x)
}
func (h *EventHandler) GetArtist(c *gin.Context) {
	x, e := h.s.Artist(c, c.Param("artist_id"))
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Artist retrieved successfully", x)
}
func (h *EventHandler) UpdateArtist(c *gin.Context) {
	q := dto.UpdateArtistRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.UpdateArtist(c, c.Param("artist_id"), &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Artist updated successfully", x)
}
func (h *EventHandler) DeleteArtist(c *gin.Context) {
	if e := h.s.DeleteArtist(c, c.Param("artist_id")); e != nil {
		fail(c, e)
		return
	}
	helpers.Message(c, 204, "Artist deleted successfully")
}
func (h *EventHandler) Lineup(c *gin.Context) {
	x, e := h.s.Lineup(c, c.Param("session_id"))
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Lineup retrieved successfully", x)
}
func (h *EventHandler) AddArtist(c *gin.Context) {
	q := dto.CreateEventArtistRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.AddArtist(c, c.Param("session_id"), &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 201, "Artist added to lineup successfully", x)
}
func (h *EventHandler) UpdateLineup(c *gin.Context) {
	q := dto.UpdateEventArtistRequest{}
	if !bind(c, &q) {
		return
	}
	x, e := h.s.UpdateArtistLineup(c, c.Param("session_id"), c.Param("event_artist_id"), &q)
	if e != nil {
		fail(c, e)
		return
	}
	helpers.Success(c, 200, "Lineup updated successfully", x)
}
func (h *EventHandler) DeleteLineup(c *gin.Context) {
	if e := h.s.DeleteArtistLineup(c, c.Param("session_id"), c.Param("event_artist_id")); e != nil {
		fail(c, e)
		return
	}
	helpers.Message(c, 204, "Artist removed from lineup successfully")
}
