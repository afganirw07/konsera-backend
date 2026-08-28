package models

import (
	"time"

	"github.com/google/uuid"
)

type PromoCode struct {
	ID                uuid.UUID  `json:"id"`
	Code              string     `json:"code"`
	EventID           *uuid.UUID `json:"event_id,omitempty"`
	DiscountType      string     `json:"discount_type"`
	DiscountValue     float64    `json:"discount_value"`
	MaxDiscountAmount *float64   `json:"max_discount_amount,omitempty"`
	UsageLimit        *int       `json:"usage_limit,omitempty"`
	UsageCount        int        `json:"usage_count"`
	PerUserLimit      int        `json:"per_user_limit"`
	ValidFrom         time.Time  `json:"valid_from"`
	ValidUntil        time.Time  `json:"valid_until"`
	IsActive          bool       `json:"is_active"`
	CreatedAt         time.Time  `json:"created_at"`
}

type PromoCodeUsage struct {
	ID              uuid.UUID `json:"id"`
	PromoCodeID     uuid.UUID `json:"promo_code_id"`
	UserID          uuid.UUID `json:"user_id"`
	BookingID       uuid.UUID `json:"booking_id"`
	DiscountApplied float64   `json:"discount_applied"`
	UsedAt          time.Time `json:"used_at"`
}
