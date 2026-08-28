package models

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ID        uuid.UUID `json:"id"`
	EventID   uuid.UUID `json:"event_id"`
	UserID    uuid.UUID `json:"user_id"`
	BookingID uuid.UUID `json:"booking_id"`
	Rating    int16     `json:"rating"`
	Comment   *string   `json:"comment,omitempty"`
	PhotoURLs []string  `json:"photo_urls,omitempty"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
type Favorite struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	EventID   uuid.UUID `json:"event_id"`
	CreatedAt time.Time `json:"created_at"`
}
type CommissionSetting struct {
	ID             uuid.UUID  `json:"id"`
	Scope          string     `json:"scope"`
	ScopeRefID     *uuid.UUID `json:"scope_ref_id,omitempty"`
	RatePercentage float64    `json:"rate_percentage"`
	EffectiveFrom  time.Time  `json:"effective_from"`
	CreatedAt      time.Time  `json:"created_at"`
}
type Payout struct {
	ID               uuid.UUID  `json:"id"`
	OrganizerID      uuid.UUID  `json:"organizer_id"`
	EventID          uuid.UUID  `json:"event_id"`
	GrossAmount      float64    `json:"gross_amount"`
	CommissionAmount float64    `json:"commission_amount"`
	NetAmount        float64    `json:"net_amount"`
	Status           string     `json:"status"`
	ScheduledAt      time.Time  `json:"scheduled_at"`
	PaidAt           *time.Time `json:"paid_at,omitempty"`
	BankReference    *string    `json:"bank_reference,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}
type AuditLog struct {
	ID          uuid.UUID  `json:"id"`
	ActorUserID *uuid.UUID `json:"actor_user_id,omitempty"`
	Action      string     `json:"action"`
	EntityType  string     `json:"entity_type"`
	EntityID    *uuid.UUID `json:"entity_id,omitempty"`
	BeforeState []byte     `json:"before_state,omitempty"`
	AfterState  []byte     `json:"after_state,omitempty"`
	IPAddress   *string    `json:"ip_address,omitempty"`
	UserAgent   *string    `json:"user_agent,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}
