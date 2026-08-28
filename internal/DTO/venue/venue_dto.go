package venue

type CreateVenueRequest struct {
	Name          string   `json:"name" binding:"required,max=200"`
	Type          string   `json:"type" binding:"required,oneof=indoor outdoor hybrid"`
	Address       string   `json:"address" binding:"required"`
	City          string   `json:"city" binding:"required,max=100"`
	Province      *string  `json:"province,omitempty"`
	Country       string   `json:"country,omitempty" binding:"omitempty,len=2"`
	Latitude      *float64 `json:"latitude,omitempty" binding:"omitempty,gte=-90,lte=90"`
	Longitude     *float64 `json:"longitude,omitempty" binding:"omitempty,gte=-180,lte=180"`
	TotalCapacity int      `json:"total_capacity" binding:"required,gt=0"`
}

type UpdateVenueRequest struct {
	Name          *string  `json:"name,omitempty" binding:"omitempty,max=200"`
	Type          *string  `json:"type,omitempty" binding:"omitempty,oneof=indoor outdoor hybrid"`
	Address       *string  `json:"address,omitempty"`
	City          *string  `json:"city,omitempty" binding:"omitempty,max=100"`
	Province      *string  `json:"province,omitempty"`
	Country       *string  `json:"country,omitempty" binding:"omitempty,len=2"`
	Latitude      *float64 `json:"latitude,omitempty" binding:"omitempty,gte=-90,lte=90"`
	Longitude     *float64 `json:"longitude,omitempty" binding:"omitempty,gte=-180,lte=180"`
	TotalCapacity *int     `json:"total_capacity,omitempty" binding:"omitempty,gt=0"`
}

type CreateVenueSectionRequest struct {
	Name     string `json:"name" binding:"required,max=100"`
	Capacity int    `json:"capacity" binding:"required,gt=0"`
	IsSeated *bool  `json:"is_seated,omitempty"`
}

type UpdateVenueSectionRequest struct {
	Name     *string `json:"name,omitempty" binding:"omitempty,max=100"`
	Capacity *int    `json:"capacity,omitempty" binding:"omitempty,gt=0"`
	IsSeated *bool   `json:"is_seated,omitempty"`
}

type CreateSeatRequest struct {
	RowLabel   string   `json:"row_label" binding:"required,max=10"`
	SeatNumber string   `json:"seat_number" binding:"required,max=10"`
	CoordX     *float64 `json:"coord_x,omitempty"`
	CoordY     *float64 `json:"coord_y,omitempty"`
}

type UpdateSeatRequest struct {
	RowLabel   *string  `json:"row_label,omitempty" binding:"omitempty,max=10"`
	SeatNumber *string  `json:"seat_number,omitempty" binding:"omitempty,max=10"`
	CoordX     *float64 `json:"coord_x,omitempty"`
	CoordY     *float64 `json:"coord_y,omitempty"`
}
