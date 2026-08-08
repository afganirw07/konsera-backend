package user


type UserRequest struct {
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	Password          string `json:"password"`
}

type UserResponse struct {
	ID                int    `json:"id"`
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	Auth_Provider     string `json:"auth_provider"`
	Provider_UID      string `json:"provider_uid"`
	Status            string `json:"status"`
}