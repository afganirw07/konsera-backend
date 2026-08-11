package user

import (
	"context"
	"database/sql"
	user "konsera-backend/internal/models"
	"time"

	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUserTx(
	ctx context.Context,
	tx *sql.Tx,
	user *user.User,
) error {
	query := `
		INSERT INTO users (
			email,
			phone,
			password,
			auth_provider,
			provider_uid,
			status,
			email_verified_at,
			phone_verified_at,
			last_login_at,
			created_at,
			updated_at,
			deleted_at
		)
		VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11, $12
		)
		RETURNING id
	`

	return tx.QueryRowContext(
		ctx,
		query,
		user.Email,
		user.Phone,
		user.Password,
		user.AuthProvider,
		user.ProviderUID,
		user.Status,
		user.EmailVerifiedAt,
		user.PhoneVerifiedAt,
		user.LastLoginAt,
		user.CreatedAt,
		user.UpdatedAt,
		user.DeletedAt,
	).Scan(&user.ID)
}

func (r *UserRepository) CreateUserProfileTx(
	ctx context.Context,
	tx *sql.Tx,
	profile *user.UserProfile,
) error {
	query := `
		INSERT INTO user_profiles (
			user_id,
			full_name,
			avatar_url,
			date_of_birth,
			gender,
			address_line,
			city,
			province,
			postal_code,
			country
		)
		VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10
		)
		RETURNING id
	`

	return tx.QueryRowContext(
		ctx,
		query,
		profile.UserID,
		profile.FullName,
		profile.AvatarURL,
		profile.DateOfBirth,
		profile.Gender,
		profile.AddressLine,
		profile.City,
		profile.Province,
		profile.PostalCode,
		profile.Country,
	).Scan(&profile.ID)
}

func (r *UserRepository) CreateRoleTx(
	ctx context.Context,
	tx *sql.Tx,
	role *user.Role,
) error {
	query := `
		INSERT INTO roles (
			name,
			description
		)
		VALUES ($1, $2)
		RETURNING id
	`

	return tx.QueryRowContext(
		ctx,
		query,
		role.Name,
		role.Description,
	).Scan(&role.ID)
}

func (r *UserRepository) CreateUserRoleTx(
	ctx context.Context,
	tx *sql.Tx,
	userRole *user.UserRole,
) error {
	query := `
		INSERT INTO user_roles (
			user_id,
			role_id,
			assigned_at
		)
		VALUES ($1, $2, $3)
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		userRole.UserID,
		userRole.RoleID,
		userRole.AssignedAt,
	)

	return err
}

func (r *UserRepository) CreateUserPreference(
	ctx context.Context,
	req *user.UserPreference,
) (*user.UserPreference, error) {
	query := `
		INSERT INTO user_preferences (
			user_id,
			favorite_genres,
			notify_push,
			notify_email,
			notify_sms,
			marketing_opt_in,
			created_at,
			updated_at
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		req.UserID,
		req.FavoriteGenres,
		req.NotifyPush,
		req.NotifyEmail,
		req.NotifySMS,
		req.MarketingOptIn,
		req.CreatedAt,
		req.UpdatedAt,
	).Scan(&req.ID)

	if err != nil {
		return nil, err
	}

	return req, nil
}

func (r *UserRepository) CreateOTPTx(
	ctx context.Context,
	tx *sql.Tx,
	profileID uuid.UUID,
	code int,
) error {
	query := `
		INSERT INTO otp_codes (
			profile_id,
			code,
			expires_at,
			created_at
		)
		VALUES ($1, $2, $3, $4)
	`

	now := time.Now()

	_, err := tx.ExecContext(
		ctx,
		query,
		profileID,
		code,
		now.Add(5*time.Minute),
		now,
	)

	return err
}

func (r *UserRepository) VerifyOTP(
	ctx context.Context,
	profileID string,
	code string,
) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM otp_codes
		WHERE profile_id = $1
		  AND code = $2
		  AND expires_at > NOW()
	`

	var count int

	err := r.db.QueryRowContext(
		ctx,
		query,
		profileID,
		code,
	).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *UserRepository) UpdateUserStatusTx(
	ctx context.Context,
	tx *sql.Tx,
	userID uuid.UUID,
	status string,
) error {
	query := `
		UPDATE users
		SET email_verified_at = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		status,
		time.Now(),
		userID,
	)

	return err
}

func (r *UserRepository) CheckEmailExists(
	ctx context.Context,
	email string,
) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE email = $1
		)
	`

	var exists bool

	err := r.db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *UserRepository) CheckPhoneExists(
	ctx context.Context,
	phone string,
) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE phone = $1
		)
	`

	var exists bool

	err := r.db.QueryRowContext(
		ctx,
		query,
		phone,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}