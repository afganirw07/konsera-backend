package event

import "time"

type CreateCategoryRequest struct {
	Name    string  `json:"name" binding:"required,max=100"`
	Slug    string  `json:"slug" binding:"required,max=120"`
	IconURL *string `json:"icon_url,omitempty"`
}
type UpdateCategoryRequest struct {
	Name    *string `json:"name,omitempty" binding:"omitempty,max=100"`
	Slug    *string `json:"slug,omitempty" binding:"omitempty,max=120"`
	IconURL *string `json:"icon_url,omitempty"`
}
type CreateEventRequest struct {
	OrganizerID        string   `json:"organizer_id" binding:"required"`
	VenueID            string   `json:"venue_id" binding:"required"`
	Title              string   `json:"title" binding:"required,max=250"`
	Slug               string   `json:"slug" binding:"required,max=280"`
	Description        *string  `json:"description,omitempty"`
	PosterURL          *string  `json:"poster_url,omitempty"`
	BannerURL          *string  `json:"banner_url,omitempty"`
	TermsAndConditions *string  `json:"terms_and_conditions,omitempty"`
	RefundPolicy       *string  `json:"refund_policy,omitempty"`
	CategoryIDs        []string `json:"category_ids,omitempty"`
}
type UpdateEventRequest struct {
	OrganizerID        *string  `json:"organizer_id,omitempty"`
	VenueID            *string  `json:"venue_id,omitempty"`
	Title              *string  `json:"title,omitempty" binding:"omitempty,max=250"`
	Slug               *string  `json:"slug,omitempty" binding:"omitempty,max=280"`
	Description        *string  `json:"description,omitempty"`
	PosterURL          *string  `json:"poster_url,omitempty"`
	BannerURL          *string  `json:"banner_url,omitempty"`
	Status             *string  `json:"status,omitempty"`
	IsFeatured         *bool    `json:"is_featured,omitempty"`
	TermsAndConditions *string  `json:"terms_and_conditions,omitempty"`
	RefundPolicy       *string  `json:"refund_policy,omitempty"`
	RejectionReason    *string  `json:"rejection_reason,omitempty"`
	CategoryIDs        []string `json:"category_ids,omitempty"`
}
type CreateSessionRequest struct {
	SessionName *string   `json:"session_name,omitempty"`
	StartAt     time.Time `json:"start_at" binding:"required"`
	EndAt       time.Time `json:"end_at" binding:"required"`
	GateOpenAt  time.Time `json:"gate_open_at" binding:"required"`
	GateCloseAt time.Time `json:"gate_close_at" binding:"required"`
	Status      string    `json:"status,omitempty"`
}
type UpdateSessionRequest struct {
	SessionName *string    `json:"session_name,omitempty"`
	StartAt     *time.Time `json:"start_at,omitempty"`
	EndAt       *time.Time `json:"end_at,omitempty"`
	GateOpenAt  *time.Time `json:"gate_open_at,omitempty"`
	GateCloseAt *time.Time `json:"gate_close_at,omitempty"`
	Status      *string    `json:"status,omitempty"`
}
type CreateArtistRequest struct {
	Name     string  `json:"name" binding:"required,max=150"`
	Bio      *string `json:"bio,omitempty"`
	PhotoURL *string `json:"photo_url,omitempty"`
}
type UpdateArtistRequest struct {
	Name     *string `json:"name,omitempty" binding:"omitempty,max=150"`
	Bio      *string `json:"bio,omitempty"`
	PhotoURL *string `json:"photo_url,omitempty"`
}
type CreateEventArtistRequest struct {
	ArtistID         string     `json:"artist_id" binding:"required"`
	Role             string     `json:"role,omitempty"`
	PerformanceOrder *int       `json:"performance_order,omitempty"`
	StageTime        *time.Time `json:"stage_time,omitempty"`
}
type UpdateEventArtistRequest struct {
	ArtistID         *string    `json:"artist_id,omitempty"`
	Role             *string    `json:"role,omitempty"`
	PerformanceOrder *int       `json:"performance_order,omitempty"`
	StageTime        *time.Time `json:"stage_time,omitempty"`
}
