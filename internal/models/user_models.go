package models

import (
	"time"
)

type User struct {
	ID                int        `json:"id"`
	Email             string     `json:"email"`
	Phone             string     `json:"phone"`
	Password          string     `json:"password"`
	Auth_Provider     string     `json:"auth_provider"`
	Provider_UID      string     `json:"provider_uid"`
	Status            string     `json:"status" default:"pending_verification"`
	Email_Verified_At *time.Time `json:"email_verified_at"`
	Phone_Verified_At *time.Time `json:"phone_verified_at"`
	Last_Login_At     *time.Time `json:"last_login_at"`
	Created_At        *time.Time `json:"created_at"`
	Updated_At        *time.Time `json:"updated_at"`
	Deleted_At        *time.Time `json:"deleted_at"`
}
