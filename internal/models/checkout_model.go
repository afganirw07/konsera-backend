package models

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID             uuid.UUID  `json:"id"`
	UserID         uuid.UUID  `json:"user_id"`
	TicketTierID   uuid.UUID  `json:"ticket_tier_id"`
	EventSessionID uuid.UUID  `json:"event_session_id"`
	SeatID         *uuid.UUID `json:"seat_id,omitempty"`
	Quantity       int        `json:"quantity"`
	HeldUntil      time.Time  `json:"held_until"`
	CreatedAt      time.Time  `json:"created_at"`
}
type Booking struct {
	ID                uuid.UUID  `json:"id"`
	BookingCode       string     `json:"booking_code"`
	UserID            uuid.UUID  `json:"user_id"`
	EventID           uuid.UUID  `json:"event_id"`
	Status            string     `json:"status"`
	SubtotalAmount    float64    `json:"subtotal_amount"`
	DiscountAmount    float64    `json:"discount_amount"`
	PlatformFeeAmount float64    `json:"platform_fee_amount"`
	TotalAmount       float64    `json:"total_amount"`
	PromoCodeID       *uuid.UUID `json:"promo_code_id,omitempty"`
	HoldExpiresAt     *time.Time `json:"hold_expires_at,omitempty"`
	PaidAt            *time.Time `json:"paid_at,omitempty"`
	CancelledAt       *time.Time `json:"cancelled_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}
type BookingItem struct {
	ID             uuid.UUID `json:"id"`
	BookingID      uuid.UUID `json:"booking_id"`
	TicketTierID   uuid.UUID `json:"ticket_tier_id"`
	EventSessionID uuid.UUID `json:"event_session_id"`
	UnitPrice      float64   `json:"unit_price"`
	Quantity       int       `json:"quantity"`
	LineTotal      float64   `json:"line_total"`
	CreatedAt      time.Time `json:"created_at"`
}
type PaymentMethod struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	Provider string    `json:"provider"`
	IsActive bool      `json:"is_active"`
	IconURL  *string   `json:"icon_url,omitempty"`
}
type Payment struct {
	ID                    uuid.UUID  `json:"id"`
	BookingID             uuid.UUID  `json:"booking_id"`
	PaymentMethodID       uuid.UUID  `json:"payment_method_id"`
	ProviderTransactionID *string    `json:"provider_transaction_id,omitempty"`
	Amount                float64    `json:"amount"`
	Status                string     `json:"status"`
	ExpiresAt             *time.Time `json:"expires_at,omitempty"`
	PaidAt                *time.Time `json:"paid_at,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}
type Refund struct {
	ID          uuid.UUID  `json:"id"`
	BookingID   uuid.UUID  `json:"booking_id"`
	PaymentID   uuid.UUID  `json:"payment_id"`
	Amount      float64    `json:"amount"`
	Reason      *string    `json:"reason,omitempty"`
	Status      string     `json:"status"`
	RequestedBy uuid.UUID  `json:"requested_by"`
	ProcessedBy *uuid.UUID `json:"processed_by,omitempty"`
	ProcessedAt *time.Time `json:"processed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}
