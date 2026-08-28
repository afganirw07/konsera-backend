package models

import (
	"time"

	"github.com/google/uuid"
)

type EventCategory struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Slug    string    `json:"slug"`
	IconURL *string   `json:"icon_url,omitempty"`
}
type Event struct {
	ID                 uuid.UUID   `json:"id"`
	OrganizerID        uuid.UUID   `json:"organizer_id"`
	VenueID            uuid.UUID   `json:"venue_id"`
	Title              string      `json:"title"`
	Slug               string      `json:"slug"`
	Description        *string     `json:"description,omitempty"`
	PosterURL          *string     `json:"poster_url,omitempty"`
	BannerURL          *string     `json:"banner_url,omitempty"`
	Status             string      `json:"status"`
	IsFeatured         bool        `json:"is_featured"`
	MinPrice           *float64    `json:"min_price,omitempty"`
	MaxPrice           *float64    `json:"max_price,omitempty"`
	TermsAndConditions *string     `json:"terms_and_conditions,omitempty"`
	RefundPolicy       *string     `json:"refund_policy,omitempty"`
	ApprovedBy         *uuid.UUID  `json:"approved_by,omitempty"`
	ApprovedAt         *time.Time  `json:"approved_at,omitempty"`
	RejectionReason    *string     `json:"rejection_reason,omitempty"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
	DeletedAt          *time.Time  `json:"deleted_at,omitempty"`
	CategoryIDs        []uuid.UUID `json:"category_ids,omitempty"`
}
type EventSession struct {
	ID          uuid.UUID `json:"id"`
	EventID     uuid.UUID `json:"event_id"`
	SessionName *string   `json:"session_name,omitempty"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
	GateOpenAt  time.Time `json:"gate_open_at"`
	GateCloseAt time.Time `json:"gate_close_at"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
type Artist struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Bio       *string   `json:"bio,omitempty"`
	PhotoURL  *string   `json:"photo_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
type EventArtist struct {
	ID               uuid.UUID  `json:"id"`
	EventSessionID   uuid.UUID  `json:"event_session_id"`
	ArtistID         uuid.UUID  `json:"artist_id"`
	Role             string     `json:"role"`
	PerformanceOrder *int       `json:"performance_order,omitempty"`
	StageTime        *time.Time `json:"stage_time,omitempty"`
}
