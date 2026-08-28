package models

import (
	"time"

	"github.com/google/uuid"
)

type Ticket struct {
	ID                uuid.UUID `json:"id"`
	BookingItemID     uuid.UUID `json:"booking_item_id"`
	OwnerUserID       uuid.UUID `json:"owner_user_id"`
	TicketCode        string    `json:"ticket_code"`
	QRToken           string    `json:"qr_token,omitempty"`
	QRTokenUpdatedAt  time.Time `json:"qr_token_updated_at"`
	Status            string    `json:"status"`
	WatermarkUserHash *string   `json:"watermark_user_hash,omitempty"`
	IssuedAt          time.Time `json:"issued_at"`
}
type TicketTransfer struct {
	ID          uuid.UUID  `json:"id"`
	TicketID    uuid.UUID  `json:"ticket_id"`
	FromUserID  uuid.UUID  `json:"from_user_id"`
	ToUserID    uuid.UUID  `json:"to_user_id"`
	Status      string     `json:"status"`
	RequestedAt time.Time  `json:"requested_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
}
type CheckIn struct {
	ID             uuid.UUID `json:"id"`
	TicketID       uuid.UUID `json:"ticket_id"`
	EventSessionID uuid.UUID `json:"event_session_id"`
	ScannedBy      uuid.UUID `json:"scanned_by"`
	Result         string    `json:"result"`
	DeviceID       *string   `json:"device_id,omitempty"`
	ScannedAt      time.Time `json:"scanned_at"`
	IsOfflineSync  bool      `json:"is_offline_sync"`
}
