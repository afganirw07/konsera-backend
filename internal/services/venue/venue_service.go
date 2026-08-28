package venue

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	dto "konsera-backend/internal/DTO/venue"
	"konsera-backend/internal/models"
	repository "konsera-backend/internal/repository/venue"

	"github.com/google/uuid"
)

type VenueService struct{ repo *repository.VenueRepository }

func NewVenueService(repo *repository.VenueRepository) *VenueService {
	return &VenueService{repo: repo}
}

func parseID(value, field string) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid %s", field)
	}
	return id, nil
}

func (s *VenueService) Create(ctx context.Context, req *dto.CreateVenueRequest) (*models.Venue, error) {
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Address) == "" || strings.TrimSpace(req.City) == "" {
		return nil, fmt.Errorf("name, address, and city are required")
	}
	country := req.Country
	if country == "" {
		country = "ID"
	}
	item := &models.Venue{Name: req.Name, Type: req.Type, Address: req.Address, City: req.City,
		Province: req.Province, Country: country, Latitude: req.Latitude, Longitude: req.Longitude,
		TotalCapacity: req.TotalCapacity}
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *VenueService) List(ctx context.Context) ([]*models.Venue, error) { return s.repo.List(ctx) }

func (s *VenueService) Get(ctx context.Context, value string) (*models.Venue, error) {
	id, err := parseID(value, "venue_id")
	if err != nil {
		return nil, err
	}
	return s.repo.Get(ctx, id)
}

func (s *VenueService) Update(ctx context.Context, value string, req *dto.UpdateVenueRequest) (*models.Venue, error) {
	id, err := parseID(value, "venue_id")
	if err != nil {
		return nil, err
	}
	return s.repo.Update(ctx, id, req)
}

func (s *VenueService) Delete(ctx context.Context, value string) error {
	id, err := parseID(value, "venue_id")
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func (s *VenueService) CreateSection(ctx context.Context, venueValue string, req *dto.CreateVenueSectionRequest) (*models.VenueSection, error) {
	venueID, err := s.requireVenue(ctx, venueValue)
	if err != nil {
		return nil, err
	}
	isSeated := true
	if req.IsSeated != nil {
		isSeated = *req.IsSeated
	}
	item := &models.VenueSection{VenueID: venueID, Name: req.Name, Capacity: req.Capacity, IsSeated: isSeated}
	if err := s.repo.CreateSection(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *VenueService) ListSections(ctx context.Context, venueValue string) ([]*models.VenueSection, error) {
	id, err := s.requireVenue(ctx, venueValue)
	if err != nil {
		return nil, err
	}
	return s.repo.ListSections(ctx, id)
}

func (s *VenueService) GetSection(ctx context.Context, venueValue, sectionValue string) (*models.VenueSection, error) {
	venueID, sectionID, err := s.parseSectionIDs(ctx, venueValue, sectionValue)
	if err != nil {
		return nil, err
	}
	return s.repo.GetSection(ctx, venueID, sectionID)
}

func (s *VenueService) UpdateSection(ctx context.Context, venueValue, sectionValue string, req *dto.UpdateVenueSectionRequest) (*models.VenueSection, error) {
	venueID, sectionID, err := s.parseSectionIDs(ctx, venueValue, sectionValue)
	if err != nil {
		return nil, err
	}
	return s.repo.UpdateSection(ctx, venueID, sectionID, req)
}

func (s *VenueService) DeleteSection(ctx context.Context, venueValue, sectionValue string) error {
	venueID, sectionID, err := s.parseSectionIDs(ctx, venueValue, sectionValue)
	if err != nil {
		return err
	}
	return s.repo.DeleteSection(ctx, venueID, sectionID)
}

func (s *VenueService) CreateSeat(ctx context.Context, venueValue, sectionValue string, req *dto.CreateSeatRequest) (*models.Seat, error) {
	_, sectionID, err := s.parseSectionIDs(ctx, venueValue, sectionValue)
	if err != nil {
		return nil, err
	}
	item := &models.Seat{VenueSectionID: sectionID, RowLabel: req.RowLabel, SeatNumber: req.SeatNumber, CoordX: req.CoordX, CoordY: req.CoordY}
	if err := s.repo.CreateSeat(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *VenueService) ListSeats(ctx context.Context, venueValue, sectionValue string) ([]*models.Seat, error) {
	_, sectionID, err := s.parseSectionIDs(ctx, venueValue, sectionValue)
	if err != nil {
		return nil, err
	}
	return s.repo.ListSeats(ctx, sectionID)
}

func (s *VenueService) GetSeat(ctx context.Context, venueValue, sectionValue, seatValue string) (*models.Seat, error) {
	sectionID, seatID, err := s.parseSeatIDs(ctx, venueValue, sectionValue, seatValue)
	if err != nil {
		return nil, err
	}
	return s.repo.GetSeat(ctx, sectionID, seatID)
}

func (s *VenueService) UpdateSeat(ctx context.Context, venueValue, sectionValue, seatValue string, req *dto.UpdateSeatRequest) (*models.Seat, error) {
	sectionID, seatID, err := s.parseSeatIDs(ctx, venueValue, sectionValue, seatValue)
	if err != nil {
		return nil, err
	}
	return s.repo.UpdateSeat(ctx, sectionID, seatID, req)
}

func (s *VenueService) DeleteSeat(ctx context.Context, venueValue, sectionValue, seatValue string) error {
	sectionID, seatID, err := s.parseSeatIDs(ctx, venueValue, sectionValue, seatValue)
	if err != nil {
		return err
	}
	return s.repo.DeleteSeat(ctx, sectionID, seatID)
}

func (s *VenueService) requireVenue(ctx context.Context, value string) (uuid.UUID, error) {
	id, err := parseID(value, "venue_id")
	if err != nil {
		return uuid.Nil, err
	}
	_, err = s.repo.Get(ctx, id)
	if err == sql.ErrNoRows {
		return uuid.Nil, fmt.Errorf("venue not found")
	}
	return id, err
}

func (s *VenueService) parseSectionIDs(ctx context.Context, venueValue, sectionValue string) (uuid.UUID, uuid.UUID, error) {
	venueID, err := s.requireVenue(ctx, venueValue)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	sectionID, err := parseID(sectionValue, "section_id")
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	if _, err = s.repo.GetSection(ctx, venueID, sectionID); err == sql.ErrNoRows {
		return uuid.Nil, uuid.Nil, fmt.Errorf("venue section not found")
	}
	return venueID, sectionID, err
}

func (s *VenueService) parseSeatIDs(ctx context.Context, venueValue, sectionValue, seatValue string) (uuid.UUID, uuid.UUID, error) {
	_, sectionID, err := s.parseSectionIDs(ctx, venueValue, sectionValue)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	seatID, err := parseID(seatValue, "seat_id")
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	if _, err = s.repo.GetSeat(ctx, sectionID, seatID); err == sql.ErrNoRows {
		return uuid.Nil, uuid.Nil, fmt.Errorf("seat not found")
	}
	return sectionID, seatID, err
}
