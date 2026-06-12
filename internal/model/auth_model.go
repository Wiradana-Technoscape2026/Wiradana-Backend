package model

// LoginRequest — POST /auth/login (api_planning §3.1)
type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AppUserResponse — api_planning §2.1
type AppUserResponse struct {
	ID       string  `json:"id"`
	Email    string  `json:"email"`
	Role     string  `json:"role"`
	MemberID *string `json:"member_id"`
}

// LoginResponse — api_planning §3.1
type LoginResponse struct {
	Token string          `json:"token"`
	User  AppUserResponse `json:"user"`
}
