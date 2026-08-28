package community

import "encoding/json"

type CreateReviewRequest struct {
	EventID   string   `json:"event_id" binding:"required"`
	BookingID string   `json:"booking_id" binding:"required"`
	Rating    int16    `json:"rating" binding:"required,gte=1,lte=5"`
	Comment   *string  `json:"comment,omitempty"`
	PhotoURLs []string `json:"photo_urls,omitempty"`
}
type UpdateReviewRequest struct {
	Rating    *int16   `json:"rating,omitempty" binding:"omitempty,gte=1,lte=5"`
	Comment   *string  `json:"comment,omitempty"`
	PhotoURLs []string `json:"photo_urls,omitempty"`
	Status    *string  `json:"status,omitempty" binding:"omitempty,oneof=visible hidden reported"`
}
type CreateFavoriteRequest struct {
	EventID string `json:"event_id" binding:"required"`
}
type CreateCommissionRequest struct {
	Scope          string  `json:"scope" binding:"required,oneof=global organizer event"`
	ScopeRefID     *string `json:"scope_ref_id,omitempty"`
	RatePercentage float64 `json:"rate_percentage" binding:"gte=0,lte=100"`
	EffectiveFrom  *string `json:"effective_from,omitempty"`
}
type UpdateCommissionRequest struct {
	RatePercentage *float64 `json:"rate_percentage,omitempty" binding:"omitempty,gte=0,lte=100"`
	EffectiveFrom  *string  `json:"effective_from,omitempty"`
}
type CreatePayoutRequest struct {
	OrganizerID      string  `json:"organizer_id" binding:"required"`
	EventID          string  `json:"event_id" binding:"required"`
	GrossAmount      float64 `json:"gross_amount" binding:"gte=0"`
	CommissionAmount float64 `json:"commission_amount" binding:"gte=0"`
	NetAmount        float64 `json:"net_amount" binding:"gte=0"`
	ScheduledAt      string  `json:"scheduled_at" binding:"required"`
	BankReference    *string `json:"bank_reference,omitempty"`
}
type UpdatePayoutRequest struct {
	Status        *string `json:"status,omitempty" binding:"omitempty,oneof=scheduled processing paid failed on_hold"`
	PaidAt        *string `json:"paid_at,omitempty"`
	BankReference *string `json:"bank_reference,omitempty"`
}
type CreateAuditRequest struct {
	Action      string          `json:"action" binding:"required,oneof=create update delete approve reject login payout refund"`
	EntityType  string          `json:"entity_type" binding:"required,max=50"`
	EntityID    *string         `json:"entity_id,omitempty"`
	BeforeState json.RawMessage `json:"before_state,omitempty"`
	AfterState  json.RawMessage `json:"after_state,omitempty"`
	IPAddress   *string         `json:"ip_address,omitempty"`
	UserAgent   *string         `json:"user_agent,omitempty"`
}
