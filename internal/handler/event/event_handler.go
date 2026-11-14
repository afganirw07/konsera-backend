package event

import (
	"database/sql"
	"errors"
	dto "konsera-backend/internal/DTO/event"
	"konsera-backend/internal/helpers"
	service "konsera-backend/internal/services/event"
	"konsera-backend/internal/services/imageupload"
	"net/http"

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
func uploadFail(c *gin.Context, e error) {
	st := 500
	switch {
	case errors.Is(e, sql.ErrNoRows):
		st = 404
	case errors.Is(e, imageupload.ErrInvalidImage):
		st = 400
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

// @Summary Upload category icon
// @Tags Event Categories
// @Accept mpfd
// @Produce json
// @Security BearerAuth
// @Param category_id path string true "Category ID"
// @Param file formData file true "Category icon"
// @Success 200 {object} models.EventCategory
// @Router /event-categories/{category_id}/icon [post]
func (h *EventHandler) UploadCategoryIcon(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		helpers.Error(c, 400, "file is required", err.Error())
		return
	}
	x, e := h.s.UploadCategoryIcon(c, c.Param("category_id"), file)
	if e != nil {
		uploadFail(c, e)
		return
	}
	helpers.Success(c, 200, "Category icon uploaded successfully", x)
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

// @Summary Upload event images
// @Tags Events
// @Accept mpfd
// @Produce json
// @Security BearerAuth
// @Param event_id path string true "Event ID"
// @Param poster formData file false "Event poster"
// @Param banner formData file false "Event banner"
// @Success 200 {object} models.Event
// @Router /events/{event_id}/images [post]
func (h *EventHandler) UploadEventImages(c *gin.Context) {
	poster, _ := c.FormFile("poster")
	banner, _ := c.FormFile("banner")
	if poster == nil && banner == nil {
		helpers.Error(c, 400, "poster or banner file is required", nil)
		return
	}
	x, e := h.s.UploadEventImages(c, c.Param("event_id"), poster, banner)
	if e != nil {
		uploadFail(c, e)
		return
	}
	helpers.Success(c, 200, "Event images uploaded successfully", x)
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
// @Accept json,mpfd
// @Produce json
// @Security BearerAuth
// @Param artist body event.CreateArtistRequest true "Artist"
// @Param photo formData file false "Artist photo"
// @Success 201 {object} models.Artist
// @Router /artists [post]
func (h *EventHandler) CreateArtist(c *gin.Context) {
	q := dto.CreateArtistRequest{}
	if e := c.ShouldBind(&q); e != nil {
		helpers.Error(c, 400, "Invalid request", e.Error())
		return
	}
	file, err := c.FormFile("photo")
	if err != nil && !errors.Is(err, http.ErrMissingFile) && !errors.Is(err, http.ErrNotMultipart) {
		helpers.Error(c, 400, "Invalid photo", err.Error())
		return
	}
	x, e := h.s.CreateArtist(c, &q, file)
	if e != nil {
		if errors.Is(e, imageupload.ErrInvalidImage) {
			uploadFail(c, e)
			return
		}
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

// @Summary Upload artist photo
// @Tags Artists
// @Accept mpfd
// @Produce json
// @Security BearerAuth
// @Param artist_id path string true "Artist ID"
// @Param file formData file true "Artist photo"
// @Success 200 {object} models.Artist
// @Router /artists/{artist_id}/photo [post]
func (h *EventHandler) UploadArtistPhoto(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		helpers.Error(c, 400, "file is required", err.Error())
		return
	}
	x, e := h.s.UploadArtistPhoto(c, c.Param("artist_id"), file)
	if e != nil {
		uploadFail(c, e)
		return
	}
	helpers.Success(c, 200, "Artist photo uploaded successfully", x)
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
