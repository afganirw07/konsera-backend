package event

import (
	"context"
	"database/sql"
	dto "konsera-backend/internal/DTO/event"
	"konsera-backend/internal/models"

	"github.com/google/uuid"
)

type EventRepository struct{ db *sql.DB }

func NewEventRepository(db *sql.DB) *EventRepository                 { return &EventRepository{db: db} }
func scan(scanner interface{ Scan(...any) error }, dst ...any) error { return scanner.Scan(dst...) }

func (r *EventRepository) Category(ctx context.Context, id uuid.UUID) (*models.EventCategory, error) {
	x := &models.EventCategory{}
	e := r.db.QueryRowContext(ctx, "SELECT id,name,slug,icon_url FROM event_categories WHERE id=$1", id)
	err := scan(e, &x.ID, &x.Name, &x.Slug, &x.IconURL)
	return x, err
}
func (r *EventRepository) Categories(ctx context.Context) ([]*models.EventCategory, error) {
	rows, e := r.db.QueryContext(ctx, "SELECT id,name,slug,icon_url FROM event_categories ORDER BY name")
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []*models.EventCategory{}
	for rows.Next() {
		x := &models.EventCategory{}
		if e = scan(rows, &x.ID, &x.Name, &x.Slug, &x.IconURL); e != nil {
			return nil, e
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *EventRepository) CreateCategory(ctx context.Context, x *models.EventCategory) error {
	return scan(r.db.QueryRowContext(ctx, "INSERT INTO event_categories(name,slug,icon_url) VALUES($1,$2,$3) RETURNING id,name,slug,icon_url", x.Name, x.Slug, x.IconURL), &x.ID, &x.Name, &x.Slug, &x.IconURL)
}
func (r *EventRepository) UpdateCategory(ctx context.Context, id uuid.UUID, q *dto.UpdateCategoryRequest) (*models.EventCategory, error) {
	x := &models.EventCategory{}
	e := r.db.QueryRowContext(ctx, "UPDATE event_categories SET name=COALESCE($2,name),slug=COALESCE($3,slug),icon_url=COALESCE($4,icon_url) WHERE id=$1 RETURNING id,name,slug,icon_url", id, q.Name, q.Slug, q.IconURL)
	err := scan(e, &x.ID, &x.Name, &x.Slug, &x.IconURL)
	return x, err
}
func (r *EventRepository) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	res, e := r.db.ExecContext(ctx, "DELETE FROM event_categories WHERE id=$1", id)
	if e != nil {
		return e
	}
	n, e := res.RowsAffected()
	if e != nil {
		return e
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

const eventCols = `id,organizer_id,venue_id,title,slug,description,poster_url,banner_url,status,is_featured,min_price,max_price,terms_and_conditions,refund_policy,approved_by,approved_at,rejection_reason,created_at,updated_at,deleted_at`

func scanEvent(s interface{ Scan(...any) error }) (*models.Event, error) {
	x := &models.Event{}
	e := scan(s, &x.ID, &x.OrganizerID, &x.VenueID, &x.Title, &x.Slug, &x.Description, &x.PosterURL, &x.BannerURL, &x.Status, &x.IsFeatured, &x.MinPrice, &x.MaxPrice, &x.TermsAndConditions, &x.RefundPolicy, &x.ApprovedBy, &x.ApprovedAt, &x.RejectionReason, &x.CreatedAt, &x.UpdatedAt, &x.DeletedAt)
	return x, e
}
func (r *EventRepository) Event(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	return scanEvent(r.db.QueryRowContext(ctx, "SELECT "+eventCols+" FROM events WHERE id=$1 AND deleted_at IS NULL", id))
}
func (r *EventRepository) Events(ctx context.Context) ([]*models.Event, error) {
	rows, e := r.db.QueryContext(ctx, "SELECT "+eventCols+" FROM events WHERE deleted_at IS NULL ORDER BY created_at DESC")
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []*models.Event{}
	for rows.Next() {
		x, e := scanEvent(rows)
		if e != nil {
			return nil, e
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *EventRepository) CreateEvent(ctx context.Context, x *models.Event) error {
	y, e := scanEvent(r.db.QueryRowContext(ctx, "INSERT INTO events(organizer_id,venue_id,title,slug,description,poster_url,banner_url,terms_and_conditions,refund_policy) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "+eventCols, x.OrganizerID, x.VenueID, x.Title, x.Slug, x.Description, x.PosterURL, x.BannerURL, x.TermsAndConditions, x.RefundPolicy))
	if e != nil {
		return e
	}
	*x = *y
	return nil
}
func (r *EventRepository) UpdateEvent(ctx context.Context, id uuid.UUID, q *dto.UpdateEventRequest) (*models.Event, error) {
	return scanEvent(r.db.QueryRowContext(ctx, "UPDATE events SET organizer_id=COALESCE($2,organizer_id),venue_id=COALESCE($3,venue_id),title=COALESCE($4,title),slug=COALESCE($5,slug),description=COALESCE($6,description),poster_url=COALESCE($7,poster_url),banner_url=COALESCE($8,banner_url),status=COALESCE($9,status),is_featured=COALESCE($10,is_featured),terms_and_conditions=COALESCE($11,terms_and_conditions),refund_policy=COALESCE($12,refund_policy),rejection_reason=COALESCE($13,rejection_reason) WHERE id=$1 AND deleted_at IS NULL RETURNING "+eventCols, id, q.OrganizerID, q.VenueID, q.Title, q.Slug, q.Description, q.PosterURL, q.BannerURL, q.Status, q.IsFeatured, q.TermsAndConditions, q.RefundPolicy, q.RejectionReason))
}
func (r *EventRepository) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	res, e := r.db.ExecContext(ctx, "UPDATE events SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL", id)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
func (r *EventRepository) SetCategories(ctx context.Context, eventID uuid.UUID, ids []uuid.UUID) error {
	tx, e := r.db.BeginTx(ctx, nil)
	if e != nil {
		return e
	}
	defer tx.Rollback()
	if _, e = tx.ExecContext(ctx, "DELETE FROM event_category_pivot WHERE event_id=$1", eventID); e != nil {
		return e
	}
	for _, id := range ids {
		if _, e = tx.ExecContext(ctx, "INSERT INTO event_category_pivot(event_id,category_id) VALUES($1,$2)", eventID, id); e != nil {
			return e
		}
	}
	return tx.Commit()
}

func (r *EventRepository) Sessions(ctx context.Context, eventID uuid.UUID) ([]*models.EventSession, error) {
	rows, e := r.db.QueryContext(ctx, "SELECT id,event_id,session_name,start_at,end_at,gate_open_at,gate_close_at,status,created_at FROM event_sessions WHERE event_id=$1 ORDER BY start_at", eventID)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []*models.EventSession{}
	for rows.Next() {
		x := &models.EventSession{}
		if e = scan(rows, &x.ID, &x.EventID, &x.SessionName, &x.StartAt, &x.EndAt, &x.GateOpenAt, &x.GateCloseAt, &x.Status, &x.CreatedAt); e != nil {
			return nil, e
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *EventRepository) CreateSession(ctx context.Context, x *models.EventSession) error {
	return scan(r.db.QueryRowContext(ctx, "INSERT INTO event_sessions(event_id,session_name,start_at,end_at,gate_open_at,gate_close_at,status) VALUES($1,$2,$3,$4,$5,$6,COALESCE(NULLIF($7,''),'scheduled')) RETURNING id,event_id,session_name,start_at,end_at,gate_open_at,gate_close_at,status,created_at", x.EventID, x.SessionName, x.StartAt, x.EndAt, x.GateOpenAt, x.GateCloseAt, x.Status), &x.ID, &x.EventID, &x.SessionName, &x.StartAt, &x.EndAt, &x.GateOpenAt, &x.GateCloseAt, &x.Status, &x.CreatedAt)
}
func (r *EventRepository) Session(ctx context.Context, eventID, id uuid.UUID) (*models.EventSession, error) {
	x := &models.EventSession{}
	e := r.db.QueryRowContext(ctx, "SELECT id,event_id,session_name,start_at,end_at,gate_open_at,gate_close_at,status,created_at FROM event_sessions WHERE event_id=$1 AND id=$2", eventID, id)
	err := scan(e, &x.ID, &x.EventID, &x.SessionName, &x.StartAt, &x.EndAt, &x.GateOpenAt, &x.GateCloseAt, &x.Status, &x.CreatedAt)
	return x, err
}
func (r *EventRepository) UpdateSession(ctx context.Context, eventID, id uuid.UUID, q *dto.UpdateSessionRequest) (*models.EventSession, error) {
	x := &models.EventSession{}
	e := r.db.QueryRowContext(ctx, "UPDATE event_sessions SET session_name=COALESCE($3,session_name),start_at=COALESCE($4,start_at),end_at=COALESCE($5,end_at),gate_open_at=COALESCE($6,gate_open_at),gate_close_at=COALESCE($7,gate_close_at),status=COALESCE($8,status) WHERE event_id=$1 AND id=$2 RETURNING id,event_id,session_name,start_at,end_at,gate_open_at,gate_close_at,status,created_at", eventID, id, q.SessionName, q.StartAt, q.EndAt, q.GateOpenAt, q.GateCloseAt, q.Status)
	err := scan(e, &x.ID, &x.EventID, &x.SessionName, &x.StartAt, &x.EndAt, &x.GateOpenAt, &x.GateCloseAt, &x.Status, &x.CreatedAt)
	return x, err
}
func (r *EventRepository) DeleteSession(ctx context.Context, eventID, id uuid.UUID) error {
	res, e := r.db.ExecContext(ctx, "DELETE FROM event_sessions WHERE event_id=$1 AND id=$2", eventID, id)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *EventRepository) Artists(ctx context.Context) ([]*models.Artist, error) {
	rows, e := r.db.QueryContext(ctx, "SELECT id,name,bio,photo_url,created_at FROM artists ORDER BY name")
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []*models.Artist{}
	for rows.Next() {
		x := &models.Artist{}
		if e = scan(rows, &x.ID, &x.Name, &x.Bio, &x.PhotoURL, &x.CreatedAt); e != nil {
			return nil, e
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *EventRepository) Artist(ctx context.Context, id uuid.UUID) (*models.Artist, error) {
	x := &models.Artist{}
	e := r.db.QueryRowContext(ctx, "SELECT id,name,bio,photo_url,created_at FROM artists WHERE id=$1", id)
	err := scan(e, &x.ID, &x.Name, &x.Bio, &x.PhotoURL, &x.CreatedAt)
	return x, err
}
func (r *EventRepository) CreateArtist(ctx context.Context, x *models.Artist) error {
	return scan(r.db.QueryRowContext(ctx, "INSERT INTO artists(name,bio,photo_url) VALUES($1,$2,$3) RETURNING id,name,bio,photo_url,created_at", x.Name, x.Bio, x.PhotoURL), &x.ID, &x.Name, &x.Bio, &x.PhotoURL, &x.CreatedAt)
}
func (r *EventRepository) UpdateArtist(ctx context.Context, id uuid.UUID, q *dto.UpdateArtistRequest) (*models.Artist, error) {
	x := &models.Artist{}
	e := r.db.QueryRowContext(ctx, "UPDATE artists SET name=COALESCE($2,name),bio=COALESCE($3,bio),photo_url=COALESCE($4,photo_url) WHERE id=$1 RETURNING id,name,bio,photo_url,created_at", id, q.Name, q.Bio, q.PhotoURL)
	err := scan(e, &x.ID, &x.Name, &x.Bio, &x.PhotoURL, &x.CreatedAt)
	return x, err
}
func (r *EventRepository) DeleteArtist(ctx context.Context, id uuid.UUID) error {
	res, e := r.db.ExecContext(ctx, "DELETE FROM artists WHERE id=$1", id)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
func (r *EventRepository) EventArtists(ctx context.Context, sessionID uuid.UUID) ([]*models.EventArtist, error) {
	rows, e := r.db.QueryContext(ctx, "SELECT id,event_session_id,artist_id,role,performance_order,stage_time FROM event_artists WHERE event_session_id=$1 ORDER BY performance_order NULLS LAST,id", sessionID)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []*models.EventArtist{}
	for rows.Next() {
		x := &models.EventArtist{}
		if e = scan(rows, &x.ID, &x.EventSessionID, &x.ArtistID, &x.Role, &x.PerformanceOrder, &x.StageTime); e != nil {
			return nil, e
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *EventRepository) EventArtist(ctx context.Context, sessionID, id uuid.UUID) (*models.EventArtist, error) {
	x := &models.EventArtist{}
	e := r.db.QueryRowContext(ctx, "SELECT id,event_session_id,artist_id,role,performance_order,stage_time FROM event_artists WHERE event_session_id=$1 AND id=$2", sessionID, id)
	err := scan(e, &x.ID, &x.EventSessionID, &x.ArtistID, &x.Role, &x.PerformanceOrder, &x.StageTime)
	return x, err
}
func (r *EventRepository) CreateEventArtist(ctx context.Context, x *models.EventArtist) error {
	return scan(r.db.QueryRowContext(ctx, "INSERT INTO event_artists(event_session_id,artist_id,role,performance_order,stage_time) VALUES($1,$2,COALESCE(NULLIF($3,''),'lineup'),$4,$5) RETURNING id,event_session_id,artist_id,role,performance_order,stage_time", x.EventSessionID, x.ArtistID, x.Role, x.PerformanceOrder, x.StageTime), &x.ID, &x.EventSessionID, &x.ArtistID, &x.Role, &x.PerformanceOrder, &x.StageTime)
}
func (r *EventRepository) UpdateEventArtist(ctx context.Context, sessionID, id uuid.UUID, q *dto.UpdateEventArtistRequest) (*models.EventArtist, error) {
	x := &models.EventArtist{}
	e := r.db.QueryRowContext(ctx, "UPDATE event_artists SET artist_id=COALESCE($3,artist_id),role=COALESCE($4,role),performance_order=COALESCE($5,performance_order),stage_time=COALESCE($6,stage_time) WHERE event_session_id=$1 AND id=$2 RETURNING id,event_session_id,artist_id,role,performance_order,stage_time", sessionID, id, q.ArtistID, q.Role, q.PerformanceOrder, q.StageTime)
	err := scan(e, &x.ID, &x.EventSessionID, &x.ArtistID, &x.Role, &x.PerformanceOrder, &x.StageTime)
	return x, err
}
func (r *EventRepository) DeleteEventArtist(ctx context.Context, sessionID, id uuid.UUID) error {
	res, e := r.db.ExecContext(ctx, "DELETE FROM event_artists WHERE event_session_id=$1 AND id=$2", sessionID, id)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
