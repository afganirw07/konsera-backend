package user

import (
	"context"
	"database/sql"
	dto "konsera-backend/internal/DTO/user"
	user "konsera-backend/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *user.User) (res *dto.UserResponse, err error) {
	query := `
		     INSERT INTO users (
			 	email, phone, password, auth_provider, provider_uid, status, email_verified_at, phone_verified_at, last_login_at, created_at, updated_at, deleted_at
			 )
				VALUES (
					$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
				)
				RETURNING id
	`

	err = r.db.QueryRowContext(ctx, query, user.Email, user.Phone, user.Password, user.Auth_Provider, user.Provider_UID, user.Status, user.Email_Verified_At, user.Phone_Verified_At, user.Last_Login_At, user.Created_At, user.Updated_At, user.Deleted_At).Scan(&user.ID)

	if err != nil {
		return nil, err
	}

	res = &dto.UserResponse{
		ID:            user.ID,
		Email:         user.Email,
		Phone:         user.Phone,
		Auth_Provider: user.Auth_Provider,
		Provider_UID:  user.Provider_UID,
		Status:        user.Status,
	}

	return res, nil
}
