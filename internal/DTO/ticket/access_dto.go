package ticket

import "time"

type CreateTransferRequest struct {
	ToUserID string `json:"to_user_id" binding:"required"`
}
type UpdateTransferRequest struct {
	Status string `json:"status" binding:"required,oneof=accepted rejected cancelled"`
}
type CreateCheckInRequest struct {
	TicketID       string     `json:"ticket_id" binding:"required"`
	EventSessionID string     `json:"event_session_id" binding:"required"`
	Result         string     `json:"result" binding:"required,oneof=success duplicate invalid expired not_yet_open closed"`
	DeviceID       *string    `json:"device_id,omitempty"`
	IsOfflineSync  bool       `json:"is_offline_sync,omitempty"`
	ScannedAt      *time.Time `json:"scanned_at,omitempty"`
}
type UpdateTicketRequest struct {
	Status *string `json:"status,omitempty" binding:"omitempty,oneof=issued transferred checked_in expired cancelled void"`
}
