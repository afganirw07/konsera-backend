package organizer

import (
	"context"
	"database/sql"
	"fmt"
	dto "konsera-backend/internal/DTO/organizer"
	"konsera-backend/internal/models"
	repository "konsera-backend/internal/repository/organizer"
	"strings"

	"github.com/google/uuid"
)

type OrganizerService struct {
	repo *repository.OrganizerRepository
}

func NewOrganizerService(repo *repository.OrganizerRepository) *OrganizerService {
	return &OrganizerService{repo: repo}
}

func parseID(value, field string) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid %s", field)
	}
	return id, nil
}

func (s *OrganizerService) Create(ctx context.Context, req *dto.CreateOrganizerRequest) (*models.Organizer, error) {
	userID, err := parseID(req.UserID, "user_id")
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(req.CompanyName) == "" {
		return nil, fmt.Errorf("company_name is required")
	}
	item := &models.Organizer{UserID: userID, CompanyName: req.CompanyName,
		LegalEntityType: req.LegalEntityType, TaxID: req.TaxID, BankAccountName: req.BankAccountName,
		BankAccountNumber: req.BankAccountNumber, BankName: req.BankName,
		VerificationStatus: "pending", CommissionRateOverride: req.CommissionRateOverride}
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *OrganizerService) List(ctx context.Context) ([]*models.Organizer, error) {
	return s.repo.List(ctx)
}

func (s *OrganizerService) Get(ctx context.Context, id string) (*models.Organizer, error) {
	parsed, err := parseID(id, "organizer_id")
	if err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, parsed)
}

func (s *OrganizerService) Update(ctx context.Context, id string, req *dto.UpdateOrganizerRequest) (*models.Organizer, error) {
	parsed, err := parseID(id, "organizer_id")
	if err != nil {
		return nil, err
	}
	if req.CompanyName != nil && strings.TrimSpace(*req.CompanyName) == "" {
		return nil, fmt.Errorf("company_name cannot be empty")
	}
	return s.repo.Update(ctx, parsed, &models.Organizer{CompanyName: stringValue(req.CompanyName),
		LegalEntityType: req.LegalEntityType, TaxID: req.TaxID, BankAccountName: req.BankAccountName,
		BankAccountNumber: req.BankAccountNumber, BankName: req.BankName,
		VerificationStatus: stringValue(req.VerificationStatus), CommissionRateOverride: req.CommissionRateOverride})
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func (s *OrganizerService) Delete(ctx context.Context, id string) error {
	parsed, err := parseID(id, "organizer_id")
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, parsed)
}

func (s *OrganizerService) CreateVerification(ctx context.Context, organizerID string, req *dto.CreateOrganizerVerificationRequest) (*models.OrganizerVerification, error) {
	parsed, err := parseID(organizerID, "organizer_id")
	if err != nil {
		return nil, err
	}
	if err := s.ensureOrganizer(ctx, parsed); err != nil {
		return nil, err
	}
	item := &models.OrganizerVerification{OrganizerID: parsed, DocumentType: req.DocumentType, DocumentURL: req.DocumentURL}
	if err := s.repo.CreateVerification(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *OrganizerService) ListVerifications(ctx context.Context, organizerID string) ([]*models.OrganizerVerification, error) {
	parsed, err := parseID(organizerID, "organizer_id")
	if err != nil {
		return nil, err
	}
	if err := s.ensureOrganizer(ctx, parsed); err != nil {
		return nil, err
	}
	return s.repo.ListVerifications(ctx, parsed)
}

func (s *OrganizerService) GetVerification(ctx context.Context, organizerID, id string) (*models.OrganizerVerification, error) {
	organizerUUID, err := parseID(organizerID, "organizer_id")
	if err != nil {
		return nil, err
	}
	verificationUUID, err := parseID(id, "verification_id")
	if err != nil {
		return nil, err
	}
	return s.repo.GetVerification(ctx, organizerUUID, verificationUUID)
}

func (s *OrganizerService) UpdateVerification(ctx context.Context, organizerID, id string, req *dto.UpdateOrganizerVerificationRequest) (*models.OrganizerVerification, error) {
	organizerUUID, err := parseID(organizerID, "organizer_id")
	if err != nil {
		return nil, err
	}
	verificationUUID, err := parseID(id, "verification_id")
	if err != nil {
		return nil, err
	}
	return s.repo.UpdateVerification(ctx, organizerUUID, verificationUUID, &models.OrganizerVerification{
		DocumentType: stringValue(req.DocumentType), DocumentURL: stringValue(req.DocumentURL), Status: stringValue(req.Status),
		ReviewedBy: parseOptionalUUID(req.ReviewedBy), ReviewedAt: req.ReviewedAt, RejectionReason: req.RejectionReason})
}

func parseOptionalUUID(value *string) *uuid.UUID {
	if value == nil {
		return nil
	}
	id, err := uuid.Parse(*value)
	if err != nil {
		return nil
	}
	return &id
}

func (s *OrganizerService) DeleteVerification(ctx context.Context, organizerID, id string) error {
	organizerUUID, err := parseID(organizerID, "organizer_id")
	if err != nil {
		return err
	}
	verificationUUID, err := parseID(id, "verification_id")
	if err != nil {
		return err
	}
	return s.repo.DeleteVerification(ctx, organizerUUID, verificationUUID)
}

func (s *OrganizerService) ensureOrganizer(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetByID(ctx, id)
	if err == sql.ErrNoRows {
		return fmt.Errorf("organizer not found")
	}
	return err
}
