package notification

import (
	"context"
	"database/sql"
	"encoding/json"

	dto "konsera-backend/internal/DTO/notification"
	"konsera-backend/internal/models"

	"github.com/google/uuid"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

const templateColumns = `id,code,channel,subject,body_template,created_at`

func scanTemplate(s interface{ Scan(...any) error }) (*models.NotificationTemplate, error) {
	x := &models.NotificationTemplate{}
	err := s.Scan(&x.ID, &x.Code, &x.Channel, &x.Subject, &x.BodyTemplate, &x.CreatedAt)
	return x, err
}

func (r *Repository) ListTemplates(ctx context.Context) ([]*models.NotificationTemplate, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+templateColumns+` FROM notification_templates ORDER BY code, channel`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]*models.NotificationTemplate, 0)
	for rows.Next() {
		item, err := scanTemplate(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) GetTemplate(ctx context.Context, id uuid.UUID) (*models.NotificationTemplate, error) {
	return scanTemplate(r.db.QueryRowContext(ctx, `SELECT `+templateColumns+` FROM notification_templates WHERE id=$1`, id))
}

func (r *Repository) CreateTemplate(ctx context.Context, item *models.NotificationTemplate) error {
	created, err := scanTemplate(r.db.QueryRowContext(ctx, `INSERT INTO notification_templates(code,channel,subject,body_template) VALUES($1,$2,$3,$4) RETURNING `+templateColumns, item.Code, item.Channel, item.Subject, item.BodyTemplate))
	if err != nil {
		return err
	}
	*item = *created
	return nil
}

func (r *Repository) UpdateTemplate(ctx context.Context, id uuid.UUID, q *dto.UpdateNotificationTemplateRequest) (*models.NotificationTemplate, error) {
	return scanTemplate(r.db.QueryRowContext(ctx, `UPDATE notification_templates SET code=COALESCE($2,code),channel=COALESCE($3,channel),subject=COALESCE($4,subject),body_template=COALESCE($5,body_template) WHERE id=$1 RETURNING `+templateColumns, id, q.Code, q.Channel, q.Subject, q.BodyTemplate))
}

func (r *Repository) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM notification_templates WHERE id=$1`, id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}

const notificationColumns = `id,user_id,template_id,channel,title,body,status,metadata,read_at,created_at`

func scanNotification(s interface{ Scan(...any) error }) (*models.Notification, error) {
	x := &models.Notification{}
	var rawMetadata []byte
	err := s.Scan(&x.ID, &x.UserID, &x.TemplateID, &x.Channel, &x.Title, &x.Body, &x.Status, &rawMetadata, &x.ReadAt, &x.CreatedAt)
	if err == nil && len(rawMetadata) > 0 {
		err = json.Unmarshal(rawMetadata, &x.Metadata)
	}
	return x, err
}

func (r *Repository) List(ctx context.Context, userID uuid.UUID) ([]*models.Notification, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+notificationColumns+` FROM notifications WHERE user_id=$1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]*models.Notification, 0)
	for rows.Next() {
		item, err := scanNotification(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) Get(ctx context.Context, userID, id uuid.UUID) (*models.Notification, error) {
	return scanNotification(r.db.QueryRowContext(ctx, `SELECT `+notificationColumns+` FROM notifications WHERE user_id=$1 AND id=$2`, userID, id))
}

func (r *Repository) Create(ctx context.Context, item *models.Notification) error {
	metadata, err := json.Marshal(item.Metadata)
	if err != nil {
		return err
	}
	created, err := scanNotification(r.db.QueryRowContext(ctx, `INSERT INTO notifications(user_id,template_id,channel,title,body,status,metadata) VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING `+notificationColumns, item.UserID, item.TemplateID, item.Channel, item.Title, item.Body, item.Status, metadata))
	if err != nil {
		return err
	}
	*item = *created
	return nil
}

func (r *Repository) Update(ctx context.Context, userID, id uuid.UUID, q *dto.UpdateNotificationRequest) (*models.Notification, error) {
	return scanNotification(r.db.QueryRowContext(ctx, `UPDATE notifications SET title=COALESCE($3,title),body=COALESCE($4,body),status=COALESCE($5,status),metadata=COALESCE($6,metadata),read_at=CASE WHEN $7 IS TRUE THEN COALESCE(read_at,NOW()) WHEN $7 IS FALSE THEN NULL ELSE read_at END WHERE user_id=$1 AND id=$2 RETURNING `+notificationColumns, userID, id, q.Title, q.Body, q.Status, q.Metadata, q.Read))
}

func (r *Repository) Delete(ctx context.Context, userID, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM notifications WHERE user_id=$1 AND id=$2`, userID, id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}
