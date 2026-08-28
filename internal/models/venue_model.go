package models

import (
	"time"

	"github.com/google/uuid"
)

type Venue struct {
	ID            uuid.UUID  `json:"id"`
	Name          string     `json:"name"`
	Type          string     `json:"type"`
	Address       string     `json:"address"`
	City          string     `json:"city"`
	Province      *string    `json:"province,omitempty"`
	Country       string     `json:"country"`
	Latitude      *float64   `json:"latitude,omitempty"`
	Longitude     *float64   `json:"longitude,omitempty"`
	TotalCapacity int        `json:"total_capacity"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}

type VenueSection struct {
	ID        uuid.UUID `json:"id"`
	VenueID   uuid.UUID `json:"venue_id"`
	Name      string    `json:"name"`
	Capacity  int       `json:"capacity"`
	IsSeated  bool      `json:"is_seated"`
	CreatedAt time.Time `json:"created_at"`
}

type Seat struct {
	ID             uuid.UUID `json:"id"`
	VenueSectionID uuid.UUID `json:"venue_section_id"`
	RowLabel       string    `json:"row_label"`
	SeatNumber     string    `json:"seat_number"`
	CoordX         *float64  `json:"coord_x,omitempty"`
	CoordY         *float64  `json:"coord_y,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}
