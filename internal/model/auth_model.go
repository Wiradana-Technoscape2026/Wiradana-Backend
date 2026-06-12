package model

// LoginRequest — POST /auth/login
// identifier = email (pengurus) atau NIK (anggota)
type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password"   validate:"required"`
}

// RegisterPengurusRequest — POST /auth/register/pengurus
type RegisterPengurusRequest struct {
	CooperativeID string `json:"cooperative_id" validate:"required,uuid"`
	Email         string `json:"email" validate:"required,email"`
	Password      string `json:"password" validate:"required,min=8"`
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
