package notification

import "encoding/json"

type CreateNotificationTemplateRequest struct {
	Code         string  `json:"code" binding:"required,max=80"`
	Channel      string  `json:"channel" binding:"required"`
	Subject      *string `json:"subject,omitempty" binding:"omitempty,max=200"`
	BodyTemplate string  `json:"body_template" binding:"required"`
}

type UpdateNotificationTemplateRequest struct {
	Code         *string `json:"code,omitempty" binding:"omitempty,max=80"`
	Channel      *string `json:"channel,omitempty"`
	Subject      *string `json:"subject,omitempty" binding:"omitempty,max=200"`
	BodyTemplate *string `json:"body_template,omitempty"`
}

type CreateNotificationRequest struct {
	UserID     string          `json:"user_id" binding:"required"`
	TemplateID *string         `json:"template_id,omitempty"`
	Channel    string          `json:"channel" binding:"required"`
	Title      *string         `json:"title,omitempty" binding:"omitempty,max=200"`
	Body       *string         `json:"body,omitempty"`
	Status     *string         `json:"status,omitempty"`
	Metadata   json.RawMessage `json:"metadata,omitempty"`
}

type UpdateNotificationRequest struct {
	Title    *string         `json:"title,omitempty" binding:"omitempty,max=200"`
	Body     *string         `json:"body,omitempty"`
	Status   *string         `json:"status,omitempty"`
	Metadata json.RawMessage `json:"metadata,omitempty"`
	Read     *bool           `json:"read,omitempty"`
}
