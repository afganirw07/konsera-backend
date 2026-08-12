package user

import (
	"context"
	"database/sql"
	"fmt"
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

func (r *UserRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (*user.User, error) {
	query := `
		SELECT id, email, phone, password, auth_provider,
		       provider_uid, status, email_verified_at,
		       phone_verified_at, last_login_at, created_at,
		       updated_at, deleted_at
		FROM users
		WHERE email = $1
		  AND deleted_at IS NULL
	`

	u := &user.User{}

	err := r.db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.Phone,
		&u.Password,
		&u.AuthProvider,
		&u.ProviderUID,
		&u.Status,
		&u.EmailVerifiedAt,
		&u.PhoneVerifiedAt,
		&u.LastLoginAt,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) GetUserRoles(
	ctx context.Context,
	userID uuid.UUID,
) ([]string, error) {
	query := `
		SELECT r.name
		FROM roles r
		JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = $1
		ORDER BY r.name
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []string

	for rows.Next() {
		var role string

		if err := rows.Scan(&role); err != nil {
			return nil, err
		}

		roles = append(roles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *UserRepository) LoginUser(
	ctx context.Context,
	email string,
) (*user.User, error) {
	query := `
		SELECT id, email, phone, password, auth_provider,
		       provider_uid, status, email_verified_at,
		       phone_verified_at, last_login_at, created_at,
		       updated_at, deleted_at
		FROM users
		WHERE email = $1
		  AND deleted_at IS NULL
	`

	user := &user.User{}

	err := r.db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Phone,
		&user.Password,
		&user.AuthProvider,
		&user.ProviderUID,
		&user.Status,
		&user.EmailVerifiedAt,
		&user.PhoneVerifiedAt,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserProfileByUserID(
	ctx context.Context,
	userID uuid.UUID,
) (*user.UserProfile, error) {
	query := `
		SELECT id, user_id, full_name, avatar_url, date_of_birth,
		       gender, address_line, city, province, postal_code,
		       country, created_at, updated_at
		FROM user_profiles
		WHERE user_id = $1
	`

	profile := &user.UserProfile{}

	err := r.db.QueryRowContext(
		ctx,
		query,
		userID,
	).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.FullName,
		&profile.AvatarURL,
		&profile.DateOfBirth,
		&profile.Gender,
		&profile.AddressLine,
		&profile.City,
		&profile.Province,
		&profile.PostalCode,
		&profile.Country,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return profile, nil
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

func (r *UserRepository) GetRoleByNameTx(
	ctx context.Context,
	tx *sql.Tx,
	name string,
) (*user.Role, error) {
	query := `
		SELECT id, name, description, created_at
		FROM roles
		WHERE name = $1
	`

	role := &user.Role{}
	err := tx.QueryRowContext(ctx, query, name).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return role, nil
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

func (r *UserRepository) ResendOTP(
	ctx context.Context,
	profileID string,
	code int,
) error {
	query := `
		UPDATE otp_codes
		SET code = $1, expires_at = $2, created_at = $3
		WHERE profile_id = $4
	`

	now := time.Now()

	_, err := r.db.ExecContext(
		ctx,
		query,
		code,
		now.Add(5*time.Minute),
		now,
		profileID,
	)

	return err
}

func (r *UserRepository) VerifyOTP(
	ctx context.Context,
	profileID string,
	code int,
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
		return false, fmt.Errorf("[OTP]failed to verify OTP: %w", err)
	}

	return count > 0, nil
}

func (r *UserRepository) UpdateUserStatus(
	ctx context.Context,
	userID string,
	status string,
) error {
	query := `
		UPDATE users
		SET status = $1, email_verified_at = $2, updated_at = $3
		WHERE id = $4
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		status,
		time.Now(),
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

func (r *UserRepository) GetUserByProfileID(
	ctx context.Context,
	profileID string,
) (*user.User, error) {
	query := `
		SELECT u.id, u.email, u.phone, u.password, u.auth_provider,
		       u.provider_uid, u.status, u.email_verified_at,
		       u.phone_verified_at, u.last_login_at, u.created_at,
		       u.updated_at, u.deleted_at
		FROM users u
		JOIN user_profiles up ON u.id = up.user_id
		WHERE up.id = $1
	`

	user := &user.User{}

	err := r.db.QueryRowContext(
		ctx,
		query,
		profileID,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Phone,
		&user.Password,
		&user.AuthProvider,
		&user.ProviderUID,
		&user.Status,
		&user.EmailVerifiedAt,
		&user.PhoneVerifiedAt,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) CheckUserActive(
	ctx context.Context,
	userID string,
) (bool, error) {
	query := `
		SELECT status
		FROM users
		WHERE id = $1
	`

	var status string

	err := r.db.QueryRowContext(
		ctx,
		query,
		userID,
	).Scan(&status)

	if err != nil {
		return false, err
	}

	return status == "active", nil
}

func (r *UserRepository) UpdateLastLoginAt(
	ctx context.Context,
	userID uuid.UUID,
) error {
	query := `
		UPDATE users
		SET last_login_at = $1,
		    updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		time.Now(),
		time.Now(),
		userID,
	)

	return err
}
