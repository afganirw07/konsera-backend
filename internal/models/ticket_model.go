package models

import (
	"time"

	"github.com/google/uuid"
)

type TicketTier struct {
	ID                uuid.UUID  `json:"id"`
	EventID           uuid.UUID  `json:"event_id"`
	VenueSectionID    *uuid.UUID `json:"venue_section_id,omitempty"`
	Name              string     `json:"name"`
	TierType          string     `json:"tier_type"`
	Price             float64    `json:"price"`
	MaxPerTransaction int        `json:"max_per_transaction"`
	SaleStartAt       *time.Time `json:"sale_start_at,omitempty"`
	SaleEndAt         *time.Time `json:"sale_end_at,omitempty"`
	IsSeated          bool       `json:"is_seated"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type TicketInventory struct {
	ID             uuid.UUID `json:"id"`
	TicketTierID   uuid.UUID `json:"ticket_tier_id"`
	EventSessionID uuid.UUID `json:"event_session_id"`
	TotalQuota     int       `json:"total_quota"`
	SoldCount      int       `json:"sold_count"`
	HeldCount      int       `json:"held_count"`
	Version        int       `json:"version"`
	UpdatedAt      time.Time `json:"updated_at"`
}
