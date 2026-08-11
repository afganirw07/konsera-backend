package user


type CreateUserRequest struct {
	FullName string `json:"full_name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}
type UserResponse struct {
	ID            string  `json:"id"`
	ProfileID     string  `json:"profile_id"`
	Name          string  `json:"name"`
	Email         string  `json:"email"`
	Phone         *string `json:"phone,omitempty"`
	Auth_Provider *string `json:"auth_provider,omitempty"`
	Provider_UID  *string `json:"provider_uid,omitempty"`
	Status        string  `json:"status"`
}