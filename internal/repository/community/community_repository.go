package community

import (
	"context"
	"database/sql"
	dto "konsera-backend/internal/DTO/community"
	"konsera-backend/internal/models"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }
func scanReview(s interface{ Scan(...any) error }) (*models.Review, error) {
	x := &models.Review{}
	e := s.Scan(&x.ID, &x.EventID, &x.UserID, &x.BookingID, &x.Rating, &x.Comment, pq.Array(&x.PhotoURLs), &x.Status, &x.CreatedAt)
	return x, e
}

const reviewCols = `id,event_id,user_id,booking_id,rating,comment,photo_urls,status,created_at`

func (r *Repository) Reviews(ctx context.Context, event *uuid.UUID) ([]*models.Review, error) {
	rows, e := r.db.QueryContext(ctx, `SELECT `+reviewCols+` FROM reviews WHERE ($1::uuid IS NULL OR event_id=$1) ORDER BY created_at DESC`, event)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []*models.Review{}
	for rows.Next() {
		x, e := scanReview(rows)
		if e != nil {
			return nil, e
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *Repository) Review(ctx context.Context, user, id uuid.UUID) (*models.Review, error) {
	return scanReview(r.db.QueryRowContext(ctx, `SELECT `+reviewCols+` FROM reviews WHERE user_id=$1 AND id=$2`, user, id))
}
func (r *Repository) CreateReview(ctx context.Context, x *models.Review) error {
	y, e := scanReview(r.db.QueryRowContext(ctx, `INSERT INTO reviews(event_id,user_id,booking_id,rating,comment,photo_urls) VALUES($1,$2,$3,$4,$5,$6) RETURNING `+reviewCols, x.EventID, x.UserID, x.BookingID, x.Rating, x.Comment, pq.Array(x.PhotoURLs)))
	if e != nil {
		return e
	}
	*x = *y
	return nil
}
func (r *Repository) UpdateReview(ctx context.Context, user, id uuid.UUID, q *dto.UpdateReviewRequest) (*models.Review, error) {
	return scanReview(r.db.QueryRowContext(ctx, `UPDATE reviews SET rating=COALESCE($3,rating),comment=COALESCE($4,comment),photo_urls=COALESCE($5,photo_urls),status=COALESCE($6,status) WHERE user_id=$1 AND id=$2 RETURNING `+reviewCols, user, id, q.Rating, q.Comment, pq.Array(q.PhotoURLs), q.Status))
}
func (r *Repository) DeleteReview(ctx context.Context, user, id uuid.UUID) error {
	res, e := r.db.ExecContext(ctx, `DELETE FROM reviews WHERE user_id=$1 AND id=$2`, user, id)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
func scanFavorite(s interface{ Scan(...any) error }) (*models.Favorite, error) {
	x := &models.Favorite{}
	e := s.Scan(&x.ID, &x.UserID, &x.EventID, &x.CreatedAt)
	return x, e
}
func (r *Repository) Favorites(ctx context.Context, user uuid.UUID) ([]*models.Favorite, error) {
	rows, e := r.db.QueryContext(ctx, `SELECT id,user_id,event_id,created_at FROM favorites WHERE user_id=$1 ORDER BY created_at DESC`, user)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []*models.Favorite{}
	for rows.Next() {
		x, e := scanFavorite(rows)
		if e != nil {
			return nil, e
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *Repository) CreateFavorite(ctx context.Context, x *models.Favorite) error {
	return r.db.QueryRowContext(ctx, `INSERT INTO favorites(user_id,event_id) VALUES($1,$2) RETURNING id,user_id,event_id,created_at`, x.UserID, x.EventID).Scan(&x.ID, &x.UserID, &x.EventID, &x.CreatedAt)
}
func (r *Repository) DeleteFavorite(ctx context.Context, user, event uuid.UUID) error {
	res, e := r.db.ExecContext(ctx, `DELETE FROM favorites WHERE user_id=$1 AND event_id=$2`, user, event)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
