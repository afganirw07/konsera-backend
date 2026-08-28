package promo

import "time"

type CreatePromoCodeRequest struct {
	Code              string    `json:"code" binding:"required,max=50"`
	EventID           *string   `json:"event_id,omitempty"`
	DiscountType      string    `json:"discount_type" binding:"required,oneof=percentage fixed_amount"`
	DiscountValue     float64   `json:"discount_value" binding:"gte=0"`
	MaxDiscountAmount *float64  `json:"max_discount_amount,omitempty" binding:"omitempty,gte=0"`
	UsageLimit        *int      `json:"usage_limit,omitempty" binding:"omitempty,gt=0"`
	PerUserLimit      int       `json:"per_user_limit,omitempty" binding:"omitempty,gt=0"`
	ValidFrom         time.Time `json:"valid_from" binding:"required"`
	ValidUntil        time.Time `json:"valid_until" binding:"required"`
}

type UpdatePromoCodeRequest struct {
	Code              *string    `json:"code,omitempty" binding:"omitempty,max=50"`
	EventID           *string    `json:"event_id,omitempty"`
	DiscountType      *string    `json:"discount_type,omitempty" binding:"omitempty,oneof=percentage fixed_amount"`
	DiscountValue     *float64   `json:"discount_value,omitempty" binding:"omitempty,gte=0"`
	MaxDiscountAmount *float64   `json:"max_discount_amount,omitempty" binding:"omitempty,gte=0"`
	UsageLimit        *int       `json:"usage_limit,omitempty" binding:"omitempty,gt=0"`
	PerUserLimit      *int       `json:"per_user_limit,omitempty" binding:"omitempty,gt=0"`
	ValidFrom         *time.Time `json:"valid_from,omitempty"`
	ValidUntil        *time.Time `json:"valid_until,omitempty"`
	IsActive          *bool      `json:"is_active,omitempty"`
}

type ApplyPromoRequest struct {
	Code string `json:"code" binding:"required,max=50"`
}
