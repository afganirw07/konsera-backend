package checkout

import "time"

type CreateCartRequest struct {
	TicketTierID   string  `json:"ticket_tier_id" binding:"required"`
	EventSessionID string  `json:"event_session_id" binding:"required"`
	SeatID         *string `json:"seat_id,omitempty"`
	Quantity       int     `json:"quantity" binding:"required,gt=0"`
}
type UpdateCartRequest struct {
	Quantity *int `json:"quantity,omitempty" binding:"omitempty,gt=0"`
}
type CreateBookingRequest struct {
	EventID       string     `json:"event_id" binding:"required"`
	HoldExpiresAt *time.Time `json:"hold_expires_at,omitempty"`
}
type UpdateBookingRequest struct {
	Status *string `json:"status,omitempty" binding:"omitempty,oneof=pending awaiting_payment paid expired cancelled refunded partially_refunded"`
}
type CreatePaymentMethodRequest struct {
	Name     string  `json:"name" binding:"required,max=100"`
	Type     string  `json:"type" binding:"required"`
	Provider string  `json:"provider" binding:"required,max=50"`
	IconURL  *string `json:"icon_url,omitempty"`
}
type UpdatePaymentMethodRequest struct {
	Name     *string `json:"name,omitempty"`
	Type     *string `json:"type,omitempty"`
	Provider *string `json:"provider,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
	IconURL  *string `json:"icon_url,omitempty"`
}
type CreatePaymentRequest struct {
	BookingID       string     `json:"booking_id" binding:"required"`
	PaymentMethodID string     `json:"payment_method_id" binding:"required"`
	Amount          float64    `json:"amount" binding:"gte=0"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`
}
type UpdatePaymentRequest struct {
	Status                *string    `json:"status,omitempty" binding:"omitempty,oneof=pending processing success failed expired refunded partially_refunded"`
	ProviderTransactionID *string    `json:"provider_transaction_id,omitempty"`
	PaidAt                *time.Time `json:"paid_at,omitempty"`
}
type CreateRefundRequest struct {
	BookingID string  `json:"booking_id" binding:"required"`
	PaymentID string  `json:"payment_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"gt=0"`
	Reason    *string `json:"reason,omitempty"`
}
type UpdateRefundRequest struct {
	Status *string `json:"status,omitempty" binding:"omitempty,oneof=requested approved rejected processed"`
}
type CreateBookingItemRequest struct {
	TicketTierID   string  `json:"ticket_tier_id" binding:"required"`
	EventSessionID string  `json:"event_session_id" binding:"required"`
	UnitPrice      float64 `json:"unit_price" binding:"gte=0"`
	Quantity       int     `json:"quantity" binding:"required,gt=0"`
}
