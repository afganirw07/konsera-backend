package venue

import (
	"context"
	"database/sql"

	dto "konsera-backend/internal/DTO/venue"
	"konsera-backend/internal/models"

	"github.com/google/uuid"
)

type VenueRepository struct{ db *sql.DB }

func NewVenueRepository(db *sql.DB) *VenueRepository { return &VenueRepository{db: db} }

const venueColumns = `id, name, type, address, city, province, country, latitude, longitude,
total_capacity, created_at, updated_at, deleted_at`

func scanVenue(scanner interface{ Scan(...any) error }) (*models.Venue, error) {
	item := &models.Venue{}
	err := scanner.Scan(&item.ID, &item.Name, &item.Type, &item.Address, &item.City,
		&item.Province, &item.Country, &item.Latitude, &item.Longitude, &item.TotalCapacity,
		&item.CreatedAt, &item.UpdatedAt, &item.DeletedAt)
	return item, err
}

func (r *VenueRepository) Create(ctx context.Context, item *models.Venue) error {
	created, err := scanVenue(r.db.QueryRowContext(ctx, `INSERT INTO venues
		(name, type, address, city, province, country, latitude, longitude, total_capacity)
		VALUES ($1,$2,$3,$4,$5,COALESCE(NULLIF($6,''),'ID'),$7,$8,$9)
		RETURNING `+venueColumns, item.Name, item.Type, item.Address, item.City, item.Province,
		item.Country, item.Latitude, item.Longitude, item.TotalCapacity))
	if err != nil { return err }
	*item = *created
	return nil
}

func (r *VenueRepository) List(ctx context.Context) ([]*models.Venue, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+venueColumns+` FROM venues
		WHERE deleted_at IS NULL ORDER BY created_at DESC`)
	if err != nil { return nil, err }
	defer rows.Close()
	items := make([]*models.Venue, 0)
	for rows.Next() {
		item, err := scanVenue(rows)
		if err != nil { return nil, err }
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *VenueRepository) Get(ctx context.Context, id uuid.UUID) (*models.Venue, error) {
	return scanVenue(r.db.QueryRowContext(ctx, `SELECT `+venueColumns+`
		FROM venues WHERE id=$1 AND deleted_at IS NULL`, id))
}

func (r *VenueRepository) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateVenueRequest) (*models.Venue, error) {
	return scanVenue(r.db.QueryRowContext(ctx, `UPDATE venues SET
		name=COALESCE($2,name), type=COALESCE($3,type), address=COALESCE($4,address),
		city=COALESCE($5,city), province=COALESCE($6,province), country=COALESCE($7,country),
		latitude=COALESCE($8,latitude), longitude=COALESCE($9,longitude),
		total_capacity=COALESCE($10,total_capacity)
		WHERE id=$1 AND deleted_at IS NULL RETURNING `+venueColumns, id, req.Name, req.Type,
		req.Address, req.City, req.Province, req.Country, req.Latitude, req.Longitude, req.TotalCapacity))
}

func (r *VenueRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, `UPDATE venues SET deleted_at=NOW()
		WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil { return err }
	count, err := result.RowsAffected()
	if err != nil { return err }
	if count == 0 { return sql.ErrNoRows }
	return nil
}

const sectionColumns = `id, venue_id, name, capacity, is_seated, created_at`

func scanSection(scanner interface{ Scan(...any) error }) (*models.VenueSection, error) {
	item := &models.VenueSection{}
	err := scanner.Scan(&item.ID, &item.VenueID, &item.Name, &item.Capacity, &item.IsSeated, &item.CreatedAt)
	return item, err
}

func (r *VenueRepository) CreateSection(ctx context.Context, item *models.VenueSection) error {
	created, err := scanSection(r.db.QueryRowContext(ctx, `INSERT INTO venue_sections
		(venue_id,name,capacity,is_seated) VALUES ($1,$2,$3,COALESCE($4,TRUE)) RETURNING `+sectionColumns,
		item.VenueID, item.Name, item.Capacity, item.IsSeated))
	if err != nil { return err }
	*item = *created
	return nil
}

func (r *VenueRepository) ListSections(ctx context.Context, venueID uuid.UUID) ([]*models.VenueSection, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+sectionColumns+` FROM venue_sections
		WHERE venue_id=$1 ORDER BY created_at`, venueID)
	if err != nil { return nil, err }
	defer rows.Close()
	items := make([]*models.VenueSection, 0)
	for rows.Next() {
		item, err := scanSection(rows)
		if err != nil { return nil, err }
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *VenueRepository) GetSection(ctx context.Context, venueID, id uuid.UUID) (*models.VenueSection, error) {
	return scanSection(r.db.QueryRowContext(ctx, `SELECT `+sectionColumns+` FROM venue_sections WHERE venue_id=$1 AND id=$2`, venueID, id))
}

func (r *VenueRepository) UpdateSection(ctx context.Context, venueID, id uuid.UUID, req *dto.UpdateVenueSectionRequest) (*models.VenueSection, error) {
	return scanSection(r.db.QueryRowContext(ctx, `UPDATE venue_sections SET
		name=COALESCE($3,name), capacity=COALESCE($4,capacity), is_seated=COALESCE($5,is_seated)
		WHERE venue_id=$1 AND id=$2 RETURNING `+sectionColumns, venueID, id, req.Name, req.Capacity, req.IsSeated))
}

func (r *VenueRepository) DeleteSection(ctx context.Context, venueID, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM venue_sections WHERE venue_id=$1 AND id=$2`, venueID, id)
	if err != nil { return err }
	count, err := result.RowsAffected()
	if err != nil { return err }
	if count == 0 { return sql.ErrNoRows }
	return nil
}

const seatColumns = `id, venue_section_id, row_label, seat_number, coord_x, coord_y, created_at`

func scanSeat(scanner interface{ Scan(...any) error }) (*models.Seat, error) {
	item := &models.Seat{}
	err := scanner.Scan(&item.ID, &item.VenueSectionID, &item.RowLabel, &item.SeatNumber, &item.CoordX, &item.CoordY, &item.CreatedAt)
	return item, err
}

func (r *VenueRepository) CreateSeat(ctx context.Context, item *models.Seat) error {
	created, err := scanSeat(r.db.QueryRowContext(ctx, `INSERT INTO seats
		(venue_section_id,row_label,seat_number,coord_x,coord_y) VALUES ($1,$2,$3,$4,$5) RETURNING `+seatColumns,
		item.VenueSectionID, item.RowLabel, item.SeatNumber, item.CoordX, item.CoordY))
	if err != nil { return err }
	*item = *created
	return nil
}

func (r *VenueRepository) ListSeats(ctx context.Context, sectionID uuid.UUID) ([]*models.Seat, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+seatColumns+` FROM seats WHERE venue_section_id=$1 ORDER BY row_label, seat_number`, sectionID)
	if err != nil { return nil, err }
	defer rows.Close()
	items := make([]*models.Seat, 0)
	for rows.Next() {
		item, err := scanSeat(rows)
		if err != nil { return nil, err }
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *VenueRepository) GetSeat(ctx context.Context, sectionID, id uuid.UUID) (*models.Seat, error) {
	return scanSeat(r.db.QueryRowContext(ctx, `SELECT `+seatColumns+` FROM seats WHERE venue_section_id=$1 AND id=$2`, sectionID, id))
}

func (r *VenueRepository) UpdateSeat(ctx context.Context, sectionID, id uuid.UUID, req *dto.UpdateSeatRequest) (*models.Seat, error) {
	return scanSeat(r.db.QueryRowContext(ctx, `UPDATE seats SET
		row_label=COALESCE($3,row_label), seat_number=COALESCE($4,seat_number),
		coord_x=COALESCE($5,coord_x), coord_y=COALESCE($6,coord_y)
		WHERE venue_section_id=$1 AND id=$2 RETURNING `+seatColumns, sectionID, id, req.RowLabel, req.SeatNumber, req.CoordX, req.CoordY))
}

func (r *VenueRepository) DeleteSeat(ctx context.Context, sectionID, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM seats WHERE venue_section_id=$1 AND id=$2`, sectionID, id)
	if err != nil { return err }
	count, err := result.RowsAffected()
	if err != nil { return err }
	if count == 0 { return sql.ErrNoRows }
	return nil
}
