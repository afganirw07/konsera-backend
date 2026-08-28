package ticket

import "time"

type CreateTicketTierRequest struct {
	Name              string     `json:"name" binding:"required,max=100"`
	VenueSectionID    *string    `json:"venue_section_id,omitempty"`
	TierType          string     `json:"tier_type" binding:"required,oneof=early_bird regular vip vvip group_package"`
	Price             float64    `json:"price" binding:"gte=0"`
	MaxPerTransaction int        `json:"max_per_transaction,omitempty" binding:"omitempty,gt=0"`
	SaleStartAt       *time.Time `json:"sale_start_at,omitempty"`
	SaleEndAt         *time.Time `json:"sale_end_at,omitempty"`
	IsSeated          bool       `json:"is_seated,omitempty"`
}

type UpdateTicketTierRequest struct {
	Name              *string    `json:"name,omitempty" binding:"omitempty,max=100"`
	VenueSectionID    *string    `json:"venue_section_id,omitempty"`
	TierType          *string    `json:"tier_type,omitempty" binding:"omitempty,oneof=early_bird regular vip vvip group_package"`
	Price             *float64   `json:"price,omitempty" binding:"omitempty,gte=0"`
	MaxPerTransaction *int       `json:"max_per_transaction,omitempty" binding:"omitempty,gt=0"`
	SaleStartAt       *time.Time `json:"sale_start_at,omitempty"`
	SaleEndAt         *time.Time `json:"sale_end_at,omitempty"`
	IsSeated          *bool      `json:"is_seated,omitempty"`
}

type CreateInventoryRequest struct {
	EventSessionID string `json:"event_session_id" binding:"required"`
	TotalQuota     int    `json:"total_quota" binding:"gte=0"`
}

type UpdateInventoryRequest struct {
	TotalQuota *int `json:"total_quota,omitempty" binding:"omitempty,gte=0"`
	SoldCount  *int `json:"sold_count,omitempty" binding:"omitempty,gte=0"`
	HeldCount  *int `json:"held_count,omitempty" binding:"omitempty,gte=0"`
}
