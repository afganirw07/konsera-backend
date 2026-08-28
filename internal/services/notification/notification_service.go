package notification

import (
	"context"
	"fmt"
	"strings"

	dto "konsera-backend/internal/DTO/notification"
	"konsera-backend/internal/models"
	repo "konsera-backend/internal/repository/notification"

	"github.com/google/uuid"
)

type Service struct{ repo *repo.Repository }

func NewService(repository *repo.Repository) *Service { return &Service{repo: repository} }

func parseID(value, field string) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid %s", field)
	}
	return id, nil
}

func (s *Service) ListTemplates(ctx context.Context) ([]*models.NotificationTemplate, error) {
	return s.repo.ListTemplates(ctx)
}
func (s *Service) GetTemplate(ctx context.Context, value string) (*models.NotificationTemplate, error) {
	id, err := parseID(value, "template_id")
	if err != nil {
		return nil, err
	}
	return s.repo.GetTemplate(ctx, id)
}
func (s *Service) CreateTemplate(ctx context.Context, q *dto.CreateNotificationTemplateRequest) (*models.NotificationTemplate, error) {
	item := &models.NotificationTemplate{Code: strings.ToUpper(strings.TrimSpace(q.Code)), Channel: q.Channel, Subject: q.Subject, BodyTemplate: q.BodyTemplate}
	err := s.repo.CreateTemplate(ctx, item)
	return item, err
}
func (s *Service) UpdateTemplate(ctx context.Context, value string, q *dto.UpdateNotificationTemplateRequest) (*models.NotificationTemplate, error) {
	id, err := parseID(value, "template_id")
	if err != nil {
		return nil, err
	}
	return s.repo.UpdateTemplate(ctx, id, q)
}
func (s *Service) DeleteTemplate(ctx context.Context, value string) error {
	id, err := parseID(value, "template_id")
	if err != nil {
		return err
	}
	return s.repo.DeleteTemplate(ctx, id)
}

func (s *Service) List(ctx context.Context, user string) ([]*models.Notification, error) {
	id, err := parseID(user, "user_id")
	if err != nil {
		return nil, err
	}
	return s.repo.List(ctx, id)
}
func (s *Service) Get(ctx context.Context, user, value string) (*models.Notification, error) {
	userID, err := parseID(user, "user_id")
	if err != nil {
		return nil, err
	}
	id, err := parseID(value, "notification_id")
	if err != nil {
		return nil, err
	}
	return s.repo.Get(ctx, userID, id)
}
func (s *Service) Create(ctx context.Context, q *dto.CreateNotificationRequest) (*models.Notification, error) {
	userID, err := parseID(q.UserID, "user_id")
	if err != nil {
		return nil, err
	}
	var templateID *uuid.UUID
	if q.TemplateID != nil {
		id, err := parseID(*q.TemplateID, "template_id")
		if err != nil {
			return nil, err
		}
		templateID = &id
	}
	item := &models.Notification{UserID: userID, TemplateID: templateID, Channel: q.Channel, Title: q.Title, Body: q.Body, Status: q.Status, Metadata: q.Metadata}
	err = s.repo.Create(ctx, item)
	return item, err
}
func (s *Service) Update(ctx context.Context, user, value string, q *dto.UpdateNotificationRequest) (*models.Notification, error) {
	userID, err := parseID(user, "user_id")
	if err != nil {
		return nil, err
	}
	id, err := parseID(value, "notification_id")
	if err != nil {
		return nil, err
	}
	return s.repo.Update(ctx, userID, id, q)
}
func (s *Service) Delete(ctx context.Context, user, value string) error {
	userID, err := parseID(user, "user_id")
	if err != nil {
		return err
	}
	id, err := parseID(value, "notification_id")
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, userID, id)
}
