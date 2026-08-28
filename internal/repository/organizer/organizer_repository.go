package organizer

import (
	"context"
	"database/sql"
	"konsera-backend/internal/models"

	"github.com/google/uuid"
)

type OrganizerRepository struct { db *sql.DB }

func NewOrganizerRepository(db *sql.DB) *OrganizerRepository { return &OrganizerRepository{db: db} }

const organizerColumns = `id, user_id, company_name, legal_entity_type, tax_id,
bank_account_name, bank_account_number, bank_name, verification_status,
commission_rate_override, created_at, updated_at, deleted_at`

func scanOrganizer(scanner interface{ Scan(...any) error }) (*models.Organizer, error) {
	item := &models.Organizer{}
	err := scanner.Scan(&item.ID, &item.UserID, &item.CompanyName, &item.LegalEntityType,
		&item.TaxID, &item.BankAccountName, &item.BankAccountNumber, &item.BankName,
		&item.VerificationStatus, &item.CommissionRateOverride, &item.CreatedAt,
		&item.UpdatedAt, &item.DeletedAt)
	if err != nil { return nil, err }
	return item, nil
}

func (r *OrganizerRepository) Create(ctx context.Context, item *models.Organizer) error {
	created, err := scanOrganizer(r.db.QueryRowContext(ctx, `INSERT INTO organizers
		(user_id, company_name, legal_entity_type, tax_id, bank_account_name,
		 bank_account_number, bank_name, commission_rate_override)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING `+organizerColumns,
		item.UserID, item.CompanyName, item.LegalEntityType, item.TaxID,
		item.BankAccountName, item.BankAccountNumber, item.BankName, item.CommissionRateOverride))
	if err != nil { return err }
	*item = *created
	return nil
}

func (r *OrganizerRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Organizer, error) {
	return scanOrganizer(r.db.QueryRowContext(ctx,
		`SELECT `+organizerColumns+` FROM organizers WHERE id = $1 AND deleted_at IS NULL`, id))
}

func (r *OrganizerRepository) List(ctx context.Context) ([]*models.Organizer, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT `+organizerColumns+` FROM organizers WHERE deleted_at IS NULL ORDER BY created_at DESC`)
	if err != nil { return nil, err }
	defer rows.Close()
	items := make([]*models.Organizer, 0)
	for rows.Next() {
		item, err := scanOrganizer(rows)
		if err != nil { return nil, err }
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *OrganizerRepository) Update(ctx context.Context, id uuid.UUID, item *models.Organizer) (*models.Organizer, error) {
	return scanOrganizer(r.db.QueryRowContext(ctx, `UPDATE organizers SET
		company_name = COALESCE($2, company_name), legal_entity_type = COALESCE($3, legal_entity_type),
		tax_id = COALESCE($4, tax_id), bank_account_name = COALESCE($5, bank_account_name),
		bank_account_number = COALESCE($6, bank_account_number), bank_name = COALESCE($7, bank_name),
		verification_status = COALESCE($8, verification_status),
		commission_rate_override = COALESCE($9, commission_rate_override)
		WHERE id = $1 AND deleted_at IS NULL RETURNING `+organizerColumns,
		id, nullableString(item.CompanyName), item.LegalEntityType, item.TaxID,
		item.BankAccountName, item.BankAccountNumber, item.BankName, nullableString(item.VerificationStatus),
		item.CommissionRateOverride))
}

func nullableString(value string) *string { if value == "" { return nil }; return &value }

func (r *OrganizerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE organizers SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil { return err }
	count, err := result.RowsAffected()
	if err != nil { return err }
	if count == 0 { return sql.ErrNoRows }
	return nil
}

const verificationColumns = `id, organizer_id, document_type, document_url, status,
reviewed_by, reviewed_at, rejection_reason, created_at`

func scanVerification(scanner interface{ Scan(...any) error }) (*models.OrganizerVerification, error) {
	item := &models.OrganizerVerification{}
	err := scanner.Scan(&item.ID, &item.OrganizerID, &item.DocumentType, &item.DocumentURL,
		&item.Status, &item.ReviewedBy, &item.ReviewedAt, &item.RejectionReason, &item.CreatedAt)
	if err != nil { return nil, err }
	return item, nil
}

func (r *OrganizerRepository) CreateVerification(ctx context.Context, item *models.OrganizerVerification) error {
	created, err := scanVerification(r.db.QueryRowContext(ctx, `INSERT INTO organizer_verifications
		(organizer_id, document_type, document_url) VALUES ($1, $2, $3)
		RETURNING `+verificationColumns, item.OrganizerID, item.DocumentType, item.DocumentURL))
	if err != nil { return err }
	*item = *created
	return nil
}

func (r *OrganizerRepository) ListVerifications(ctx context.Context, organizerID uuid.UUID) ([]*models.OrganizerVerification, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+verificationColumns+
		` FROM organizer_verifications WHERE organizer_id = $1 ORDER BY created_at DESC`, organizerID)
	if err != nil { return nil, err }
	defer rows.Close()
	items := make([]*models.OrganizerVerification, 0)
	for rows.Next() {
		item, err := scanVerification(rows)
		if err != nil { return nil, err }
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *OrganizerRepository) GetVerification(ctx context.Context, organizerID, id uuid.UUID) (*models.OrganizerVerification, error) {
	return scanVerification(r.db.QueryRowContext(ctx, `SELECT `+verificationColumns+
		` FROM organizer_verifications WHERE organizer_id = $1 AND id = $2`, organizerID, id))
}

func (r *OrganizerRepository) UpdateVerification(ctx context.Context, organizerID, id uuid.UUID, item *models.OrganizerVerification) (*models.OrganizerVerification, error) {
	return scanVerification(r.db.QueryRowContext(ctx, `UPDATE organizer_verifications SET
		document_type = COALESCE($3, document_type), document_url = COALESCE($4, document_url),
		status = COALESCE($5, status), reviewed_by = COALESCE($6, reviewed_by),
		reviewed_at = COALESCE($7, reviewed_at), rejection_reason = COALESCE($8, rejection_reason)
		WHERE organizer_id = $1 AND id = $2 RETURNING `+verificationColumns,
		organizerID, id, nullableString(item.DocumentType), nullableString(item.DocumentURL),
		nullableString(item.Status), item.ReviewedBy, item.ReviewedAt, item.RejectionReason))
}

func (r *OrganizerRepository) DeleteVerification(ctx context.Context, organizerID, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx,
		`DELETE FROM organizer_verifications WHERE organizer_id = $1 AND id = $2`, organizerID, id)
	if err != nil { return err }
	count, err := result.RowsAffected()
	if err != nil { return err }
	if count == 0 { return sql.ErrNoRows }
	return nil
}