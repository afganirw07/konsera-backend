package models

import (
	"time"

	"github.com/google/uuid"
)

type Organizer struct {
	ID                     uuid.UUID  `json:"id"`
	UserID                 uuid.UUID  `json:"user_id"`
	CompanyName            string     `json:"company_name"`
	LegalEntityType        *string    `json:"legal_entity_type,omitempty"`
	TaxID                  *string    `json:"tax_id,omitempty"`
	BankAccountName        *string    `json:"bank_account_name,omitempty"`
	BankAccountNumber      *string    `json:"bank_account_number,omitempty"`
	BankName               *string    `json:"bank_name,omitempty"`
	VerificationStatus     string     `json:"verification_status"`
	CommissionRateOverride *float64   `json:"commission_rate_override,omitempty"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
	DeletedAt              *time.Time `json:"deleted_at,omitempty"`
}

type OrganizerVerification struct {
	ID              uuid.UUID  `json:"id"`
	OrganizerID     uuid.UUID  `json:"organizer_id"`
	DocumentType    string     `json:"document_type"`
	DocumentURL     string     `json:"document_url"`
	Status          string     `json:"status"`
	ReviewedBy      *uuid.UUID `json:"reviewed_by,omitempty"`
	ReviewedAt      *time.Time `json:"reviewed_at,omitempty"`
	RejectionReason *string    `json:"rejection_reason,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}
