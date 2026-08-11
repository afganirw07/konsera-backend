package user

import (
	"context"
	"database/sql"
	"fmt"
	emailDTO "konsera-backend/internal/DTO/email"
	dto "konsera-backend/internal/DTO/user"
	"konsera-backend/internal/models"
	userRepo "konsera-backend/internal/repository/user"
	emailService "konsera-backend/internal/services/email"
	"konsera-backend/internal/utils"
	"os"
	"time"
)

type UserService struct {
	repo         *userRepo.UserRepository
	emailService *emailService.Service
	db           *sql.DB
}

func NewUserService(repo *userRepo.UserRepository, emailService *emailService.Service, db *sql.DB) *UserService {
	return &UserService{repo: repo, emailService: emailService, db: db}
}

func (s *UserService) CreateUser(
	ctx context.Context,
	user *dto.CreateUserRequest,
	provider string,
	emailMeta *emailDTO.RegisterOTPEmailMeta,
) (*dto.UserResponse, error) {

	if user.FullName == "" ||
		user.Email == "" ||
		user.Phone == "" ||
		user.Password == "" {
		return nil, fmt.Errorf("[ERROR] Missing required fields")
	}

	emailExists, err := s.repo.CheckEmailExists(ctx, user.Email)
	if err != nil {
		return nil, fmt.Errorf(
			"[ERROR] Failed to check email existence: %w",
			err,
		)
	}

	if emailExists {
		return nil, fmt.Errorf("[ERROR] Email already exists")
	}

	phoneExists, err := s.repo.CheckPhoneExists(ctx, user.Phone)
	if err != nil {
		return nil, fmt.Errorf(
			"[ERROR] Failed to check phone existence: %w",
			err,
		)
	}

	if phoneExists {
		return nil, fmt.Errorf("[ERROR] Phone number already exists")
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf(
			"[ERROR] Failed to hash password: %w",
			err,
		)
	}

	now := time.Now()

	phone := user.Phone
	providerUID := utils.GenerateRandomString(6)

	newUser := &models.User{
		Email:        user.Email,
		Phone:        &phone,
		Password:     &hashedPassword,
		Status:       "pending_verification",
		AuthProvider: &provider,
		ProviderUID:  &providerUID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	newProfile := &models.UserProfile{
		FullName:  user.FullName,
		CreatedAt: now,
		UpdatedAt: now,
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf(
			"[ERROR] Failed to begin transaction: %w",
			err,
		)
	}

	committed := false

	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	if err := s.repo.CreateUserTx(
		ctx,
		tx,
		newUser,
	); err != nil {
		return nil, fmt.Errorf(
			"[ERROR] Failed to create user: %w",
			err,
		)
	}

	newProfile.UserID = newUser.ID

	if err := s.repo.CreateUserProfileTx(
		ctx,
		tx,
		newProfile,
	); err != nil {
		return nil, fmt.Errorf(
			"[ERROR] Failed to create user profile: %w",
			err,
		)
	}

	customerRole, _ := s.repo.GetRoleByNameTx(
		ctx,
		tx,
		"customer",
	)

	newUserRole := &models.UserRole{
		UserID:     newUser.ID,
		RoleID:     customerRole.ID,
		AssignedAt: now,
	}

	if err := s.repo.CreateUserRoleTx(
		ctx,
		tx,
		newUserRole,
	); err != nil {
		return nil, fmt.Errorf(
			"[ERROR] Failed to create user role: %w",
			err,
		)
	}

	otpCode, err := utils.GenerateOTP(6)
	if err != nil {
		return nil, fmt.Errorf(
			"[ERROR] Failed to generate OTP: %w",
			err,
		)
	}

	if err := s.repo.CreateOTPTx(
		ctx,
		tx,
		newUser.ID,
		otpCode,
	); err != nil {
		return nil, fmt.Errorf(
			"[ERROR] Failed to save OTP: %w",
			err,
		)
	}

	konseraFEURL := os.Getenv("KONSERA_FE_URL")

	emailData := emailDTO.RegisterOTPEmailData{
		Device:           emailMeta.Device,
		Location:         emailMeta.Location,
		IPAddress:        emailMeta.IPAddress,
		Time:             emailMeta.Time,
		VerifyURL:        fmt.Sprintf("%sauth/verify-email", konseraFEURL),
		OTPCode:          otpCode,
		OTPExpiry:        5,
		SecureAccountURL: fmt.Sprintf("%ssecurity", konseraFEURL),
		UnsubscribeURL:   fmt.Sprintf("%sunsubscribe", konseraFEURL),
		PrivacyURL:       fmt.Sprintf("%sprivacy", konseraFEURL),
		HelpCenterURL:    fmt.Sprintf("%shelp", konseraFEURL),
	}

	if err := s.emailService.SendRegisterOTP(
		newUser.Email,
		emailData,
	); err != nil {
		return nil, fmt.Errorf(
			"failed to send verification email: %w",
			err,
		)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf(
			"[ERROR] Failed to commit transaction: %w",
			err,
		)
	}

	committed = true

	return &dto.UserResponse{
		ID:            newUser.ID.String(),
		ProfileID:     newProfile.ID.String(),
		Name:          newProfile.FullName,
		Email:         newUser.Email,
		Phone:         newUser.Phone,
		Auth_Provider: newUser.AuthProvider,
		Provider_UID:  newUser.ProviderUID,
		Status:        newUser.Status,
	}, nil
}


func (s *UserService) VerifyOTP(
	ctx context.Context,
	profileID string,
	code int,
) (bool, error) {
	if profileID == "" {
		return false, fmt.Errorf("[ERROR] Missing profile ID")
	}

	if code == 0 {
		return false, fmt.Errorf("[ERROR] Missing OTP code")
	}

	isValid, err := s.repo.VerifyOTP(ctx, profileID, code)
	if err != nil {
		return false, fmt.Errorf(
			"[ERROR] Failed to verify OTP: %w",
			err,
		)
	}

	if !isValid {
		return false, fmt.Errorf("[ERROR] Invalid or expired OTP code")
	}

	if err := s.repo.UpdateUserStatus(
		ctx,
		profileID,
		"active",
	); err != nil {
		return false, fmt.Errorf(
			"[ERROR] Failed to update user status: %w",
			err,
		)
	}

	return true, nil
}

func (s *UserService) ResendOTP(
	ctx context.Context,
	profileID string,
	emailMeta *emailDTO.RegisterOTPEmailMeta,
) error {
	if profileID == "" {
		return fmt.Errorf("[ERROR] Missing profile ID")
	}

	otpCode, err := utils.GenerateOTP(6)
	if err != nil {
		return fmt.Errorf(
			"[ERROR] Failed to generate OTP: %w",
			err,
		)
	}

	if err := s.repo.ResendOTP(ctx, profileID, otpCode); err != nil {
		return fmt.Errorf(
			"[ERROR] Failed to resend OTP: %w",
			err,
		)
	}

	user, err := s.repo.GetUserByProfileID(ctx, profileID)
	if err != nil {
		return fmt.Errorf(
			"[ERROR] Failed to get user by profile ID: %w",
			err,
		)
	}

	konseraFEURL := os.Getenv("KONSERA_FE_URL")

	emailData := emailDTO.RegisterOTPEmailData{
		Device:           emailMeta.Device,
		Location:         emailMeta.Location,
		IPAddress:        emailMeta.IPAddress,
		Time:             emailMeta.Time,
		VerifyURL:        fmt.Sprintf("%sauth/verify-email", konseraFEURL),
		OTPCode:          otpCode,
		OTPExpiry:        5,
		SecureAccountURL: fmt.Sprintf("%ssecurity", konseraFEURL),
		UnsubscribeURL:   fmt.Sprintf("%sunsubscribe", konseraFEURL),
		PrivacyURL:       fmt.Sprintf("%sprivacy", konseraFEURL),
		HelpCenterURL:    fmt.Sprintf("%shelp", konseraFEURL),
	}

	if err := s.emailService.SendRegisterOTP(
		user.Email,
		emailData,
	); err != nil {
		return fmt.Errorf(
			"failed to send verification email: %w",
			err,
		)
	}

	return nil
}