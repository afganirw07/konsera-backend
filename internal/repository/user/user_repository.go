package user

import (
	"context"
	"database/sql"
	dto "konsera-backend/internal/DTO/user"
	user "konsera-backend/internal/models"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(
	ctx context.Context,
	user *user.User,
	profile *user.UserProfile,
) (*dto.UserResponse, error) {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// for users
	queryUser := `
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

	err = tx.QueryRowContext(
		ctx,
		queryUser,
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

	if err != nil {
		return nil, err
	}

	// for user_profiles
	queryProfile := `
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

	err = tx.QueryRowContext(
		ctx,
		queryProfile,
		user.ID,
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

	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:            user.ID.String(),
		Name:          profile.FullName,
		Email:         user.Email,
		Phone:         user.Phone,
		Auth_Provider: user.AuthProvider,
		Provider_UID:  user.ProviderUID,
		Status:        user.Status,
	}, nil
}

func (r *UserRepository) CreateUserPreference(ctx context.Context, req *user.UserPreference) (*user.UserPreference, error) {
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

func (r *UserRepository) CreateOTP(ctx context.Context, profileID string, code string) error {
	query := `
		INSERT INTO otp_codes (
			profile_id, code,  expires_at, created_at, 
		) VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(ctx, query, profileID, code, time.Now().Add(5*time.Minute), time.Now())
	return err
}

func (r *UserRepository) VerifyOTP(ctx context.Context, profileID string, code string) (bool, error) {
	query := `
		SELECT COUNT(*) FROM otp_codes
		WHERE profile_id = $1 AND code = $2 AND expires_at > NOW()
	`
	var count int
	err := r.db.QueryRowContext(ctx, query, profileID, code).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
