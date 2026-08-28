package organizer

import "time"

type CreateOrganizerRequest struct {
	UserID                 string  `json:"user_id" binding:"required"`
	CompanyName            string  `json:"company_name" binding:"required,max=200"`
	LegalEntityType        *string `json:"legal_entity_type,omitempty"`
	TaxID                  *string `json:"tax_id,omitempty"`
	BankAccountName        *string `json:"bank_account_name,omitempty"`
	BankAccountNumber      *string `json:"bank_account_number,omitempty"`
	BankName               *string `json:"bank_name,omitempty"`
	CommissionRateOverride *float64 `json:"commission_rate_override,omitempty" binding:"omitempty,gte=0,lte=100"`
}

type UpdateOrganizerRequest struct {
	CompanyName            *string  `json:"company_name,omitempty" binding:"omitempty,max=200"`
	LegalEntityType        *string  `json:"legal_entity_type,omitempty"`
	TaxID                  *string  `json:"tax_id,omitempty"`
	BankAccountName        *string  `json:"bank_account_name,omitempty"`
	BankAccountNumber      *string  `json:"bank_account_number,omitempty"`
	BankName               *string  `json:"bank_name,omitempty"`
	VerificationStatus     *string  `json:"verification_status,omitempty" binding:"omitempty,oneof=pending in_review approved rejected"`
	CommissionRateOverride *float64 `json:"commission_rate_override,omitempty" binding:"omitempty,gte=0,lte=100"`
}

type CreateOrganizerVerificationRequest struct {
	DocumentType string `json:"document_type" binding:"required,max=50"`
	DocumentURL  string `json:"document_url" binding:"required,url"`
}

type UpdateOrganizerVerificationRequest struct {
	DocumentType    *string `json:"document_type,omitempty" binding:"omitempty,max=50"`
	DocumentURL     *string `json:"document_url,omitempty" binding:"omitempty,url"`
	Status          *string `json:"status,omitempty" binding:"omitempty,oneof=pending in_review approved rejected"`
	ReviewedBy      *string `json:"reviewed_by,omitempty"`
	ReviewedAt      *time.Time `json:"reviewed_at,omitempty"`
	RejectionReason *string `json:"rejection_reason,omitempty"`
}