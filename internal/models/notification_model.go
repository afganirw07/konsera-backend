package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificationTemplate struct {
	ID           uuid.UUID `json:"id"`
	Code         string    `json:"code"`
	Channel      string    `json:"channel"`
	Subject      *string   `json:"subject,omitempty"`
	BodyTemplate string    `json:"body_template"`
	CreatedAt    time.Time `json:"created_at"`
}

type Notification struct {
	ID         uuid.UUID              `json:"id"`
	UserID     uuid.UUID              `json:"user_id"`
	TemplateID *uuid.UUID             `json:"template_id,omitempty"`
	Channel    string                 `json:"channel"`
	Title      *string                `json:"title,omitempty"`
	Body       *string                `json:"body,omitempty"`
	Status     *string                `json:"status,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	ReadAt     *time.Time             `json:"read_at,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
}
