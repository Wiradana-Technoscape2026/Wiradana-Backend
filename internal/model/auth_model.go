package model

// LoginRequest — POST /auth/login
// identifier = email (pengurus) atau NIK (anggota)
type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password"   validate:"required"`
}

// SelectCooperativeRequest — POST /auth/select-cooperative
// Digunakan ketika anggota terdaftar di lebih dari satu koperasi.
type SelectCooperativeRequest struct {
	Identifier    string `json:"identifier"     validate:"required"`
	Password      string `json:"password"       validate:"required"`
	CooperativeID string `json:"cooperative_id" validate:"required,uuid"`
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

// CooperativeOption — satu entry koperasi dalam daftar pilihan saat anggota punya lebih dari satu.
type CooperativeOption struct {
	CooperativeID   string `json:"cooperative_id"`
	CooperativeName string `json:"cooperative_name"`
	MemberID        string `json:"member_id"` // kosong untuk pengurus
}

// LoginResponse — api_planning §3.1
// Dua bentuk:
//   - Login langsung (1 koperasi): Token + User terisi, RequiresCooperativeSelection false.
//   - Pilih koperasi (>1 koperasi): RequiresCooperativeSelection true + Cooperatives terisi.
type LoginResponse struct {
	Token                        string              `json:"token,omitempty"`
	User                         *AppUserResponse    `json:"user,omitempty"`
	RequiresCooperativeSelection bool                `json:"requires_cooperative_selection,omitempty"`
	Cooperatives                 []CooperativeOption `json:"cooperatives,omitempty"`
}
