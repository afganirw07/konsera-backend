package event

import (
	"context"
	"fmt"
	dto "konsera-backend/internal/DTO/event"
	"konsera-backend/internal/models"
	repo "konsera-backend/internal/repository/event"

	"github.com/google/uuid"
)

type EventService struct{ repo *repo.EventRepository }

func NewEventService(r *repo.EventRepository) *EventService { return &EventService{repo: r} }
func id(v, n string) (uuid.UUID, error) {
	x, e := uuid.Parse(v)
	if e != nil {
		return uuid.Nil, fmt.Errorf("invalid %s", n)
	}
	return x, nil
}
func ids(values []string) ([]uuid.UUID, error) {
	out := []uuid.UUID{}
	for _, v := range values {
		x, e := id(v, "category_id")
		if e != nil {
			return nil, e
		}
		out = append(out, x)
	}
	return out, nil
}
func (s *EventService) Categories(c context.Context) ([]*models.EventCategory, error) {
	return s.repo.Categories(c)
}
func (s *EventService) Category(c context.Context, v string) (*models.EventCategory, error) {
	x, e := id(v, "category_id")
	if e != nil {
		return nil, e
	}
	return s.repo.Category(c, x)
}
func (s *EventService) CreateCategory(c context.Context, q *dto.CreateCategoryRequest) (*models.EventCategory, error) {
	x := &models.EventCategory{Name: q.Name, Slug: q.Slug, IconURL: q.IconURL}
	e := s.repo.CreateCategory(c, x)
	return x, e
}
func (s *EventService) UpdateCategory(c context.Context, v string, q *dto.UpdateCategoryRequest) (*models.EventCategory, error) {
	x, e := id(v, "category_id")
	if e != nil {
		return nil, e
	}
	return s.repo.UpdateCategory(c, x, q)
}
func (s *EventService) DeleteCategory(c context.Context, v string) error {
	x, e := id(v, "category_id")
	if e != nil {
		return e
	}
	return s.repo.DeleteCategory(c, x)
}
func (s *EventService) Create(c context.Context, q *dto.CreateEventRequest) (*models.Event, error) {
	o, e := id(q.OrganizerID, "organizer_id")
	if e != nil {
		return nil, e
	}
	v, e := id(q.VenueID, "venue_id")
	if e != nil {
		return nil, e
	}
	cs, e := ids(q.CategoryIDs)
	if e != nil {
		return nil, e
	}
	x := &models.Event{OrganizerID: o, VenueID: v, Title: q.Title, Slug: q.Slug, Description: q.Description, PosterURL: q.PosterURL, BannerURL: q.BannerURL, TermsAndConditions: q.TermsAndConditions, RefundPolicy: q.RefundPolicy, Status: "draft", CategoryIDs: cs}
	e = s.repo.CreateEvent(c, x)
	if e == nil && len(cs) > 0 {
		e = s.repo.SetCategories(c, x.ID, cs)
	}
	return x, e
}
func (s *EventService) Events(c context.Context) ([]*models.Event, error) { return s.repo.Events(c) }
func (s *EventService) Event(c context.Context, v string) (*models.Event, error) {
	x, e := id(v, "event_id")
	if e != nil {
		return nil, e
	}
	return s.repo.Event(c, x)
}
func (s *EventService) Update(c context.Context, v string, q *dto.UpdateEventRequest) (*models.Event, error) {
	x, e := id(v, "event_id")
	if e != nil {
		return nil, e
	}
	if q.OrganizerID != nil {
		if _, e = id(*q.OrganizerID, "organizer_id"); e != nil {
			return nil, e
		}
	}
	if q.VenueID != nil {
		if _, e = id(*q.VenueID, "venue_id"); e != nil {
			return nil, e
		}
	}
	out, e := s.repo.UpdateEvent(c, x, q)
	if e == nil && q.CategoryIDs != nil {
		out.CategoryIDs, e = ids(q.CategoryIDs)
		if e == nil {
			e = s.repo.SetCategories(c, x, out.CategoryIDs)
		}
	}
	return out, e
}
func (s *EventService) Delete(c context.Context, v string) error {
	x, e := id(v, "event_id")
	if e != nil {
		return e
	}
	return s.repo.DeleteEvent(c, x)
}
func (s *EventService) CreateSession(c context.Context, ev string, q *dto.CreateSessionRequest) (*models.EventSession, error) {
	eID, e := id(ev, "event_id")
	if e != nil {
		return nil, e
	}
	if !q.EndAt.After(q.StartAt) {
		return nil, fmt.Errorf("end_at must be after start_at")
	}
	st := q.Status
	if st == "" {
		st = "scheduled"
	}
	x := &models.EventSession{EventID: eID, SessionName: q.SessionName, StartAt: q.StartAt, EndAt: q.EndAt, GateOpenAt: q.GateOpenAt, GateCloseAt: q.GateCloseAt, Status: st}
	e = s.repo.CreateSession(c, x)
	return x, e
}
func (s *EventService) Sessions(c context.Context, v string) ([]*models.EventSession, error) {
	x, e := id(v, "event_id")
	if e != nil {
		return nil, e
	}
	return s.repo.Sessions(c, x)
}
func (s *EventService) Session(c context.Context, eid, sid string) (*models.EventSession, error) {
	a, e := id(eid, "event_id")
	if e != nil {
		return nil, e
	}
	b, e := id(sid, "session_id")
	if e != nil {
		return nil, e
	}
	return s.repo.Session(c, a, b)
}
func (s *EventService) UpdateSession(c context.Context, eid, sid string, q *dto.UpdateSessionRequest) (*models.EventSession, error) {
	a, e := id(eid, "event_id")
	if e != nil {
		return nil, e
	}
	b, e := id(sid, "session_id")
	if e != nil {
		return nil, e
	}
	if q.StartAt != nil && q.EndAt != nil && !q.EndAt.After(*q.StartAt) {
		return nil, fmt.Errorf("end_at must be after start_at")
	}
	return s.repo.UpdateSession(c, a, b, q)
}
func (s *EventService) DeleteSession(c context.Context, eid, sid string) error {
	a, e := id(eid, "event_id")
	if e != nil {
		return e
	}
	b, e := id(sid, "session_id")
	if e != nil {
		return e
	}
	return s.repo.DeleteSession(c, a, b)
}
func (s *EventService) Artists(c context.Context) ([]*models.Artist, error) { return s.repo.Artists(c) }
func (s *EventService) Artist(c context.Context, v string) (*models.Artist, error) {
	x, e := id(v, "artist_id")
	if e != nil {
		return nil, e
	}
	return s.repo.Artist(c, x)
}
func (s *EventService) CreateArtist(c context.Context, q *dto.CreateArtistRequest) (*models.Artist, error) {
	x := &models.Artist{Name: q.Name, Bio: q.Bio, PhotoURL: q.PhotoURL}
	e := s.repo.CreateArtist(c, x)
	return x, e
}
func (s *EventService) UpdateArtist(c context.Context, v string, q *dto.UpdateArtistRequest) (*models.Artist, error) {
	x, e := id(v, "artist_id")
	if e != nil {
		return nil, e
	}
	return s.repo.UpdateArtist(c, x, q)
}
func (s *EventService) DeleteArtist(c context.Context, v string) error {
	x, e := id(v, "artist_id")
	if e != nil {
		return e
	}
	return s.repo.DeleteArtist(c, x)
}
func (s *EventService) Lineup(c context.Context, v string) ([]*models.EventArtist, error) {
	x, e := id(v, "session_id")
	if e != nil {
		return nil, e
	}
	return s.repo.EventArtists(c, x)
}
func (s *EventService) AddArtist(c context.Context, v string, q *dto.CreateEventArtistRequest) (*models.EventArtist, error) {
	sid, e := id(v, "session_id")
	if e != nil {
		return nil, e
	}
	aid, e := id(q.ArtistID, "artist_id")
	if e != nil {
		return nil, e
	}
	role := q.Role
	if role == "" {
		role = "lineup"
	}
	x := &models.EventArtist{EventSessionID: sid, ArtistID: aid, Role: role, PerformanceOrder: q.PerformanceOrder, StageTime: q.StageTime}
	e = s.repo.CreateEventArtist(c, x)
	return x, e
}
func (s *EventService) UpdateArtistLineup(c context.Context, sid, lid string, q *dto.UpdateEventArtistRequest) (*models.EventArtist, error) {
	a, e := id(sid, "session_id")
	if e != nil {
		return nil, e
	}
	b, e := id(lid, "event_artist_id")
	if e != nil {
		return nil, e
	}
	if q.ArtistID != nil {
		if _, e = id(*q.ArtistID, "artist_id"); e != nil {
			return nil, e
		}
	}
	return s.repo.UpdateEventArtist(c, a, b, q)
}
func (s *EventService) DeleteArtistLineup(c context.Context, sid, lid string) error {
	a, e := id(sid, "session_id")
	if e != nil {
		return e
	}
	b, e := id(lid, "event_artist_id")
	if e != nil {
		return e
	}
	return s.repo.DeleteEventArtist(c, a, b)
}
