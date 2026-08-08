package user

type UserRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Email         string  `json:"email"`
	Phone         *string `json:"phone,omitempty"`
	Auth_Provider *string `json:"auth_provider,omitempty"`
	Provider_UID  *string `json:"provider_uid,omitempty"`
	Status        string  `json:"status"`
}
