package model

import (
	"time"

	"gorm.io/datatypes"
)

type CreateMemberRequest struct {
	NIK              string         `json:"nik" validate:"required,len=16,numeric"`
	FullName         string         `json:"full_name" validate:"required"`
	Address          string         `json:"address"`
	BirthDate        string         `json:"birth_date"` // "YYYY-MM-DD", opsional
	CustomAttributes datatypes.JSON `json:"custom_attributes"`
	// Opsional: buat akun login anggota sekaligus. Password min 6 karakter.
	// Login menggunakan NIK + password (bukan email).
	Password *string `json:"password" validate:"omitempty,min=6"`
}

type UpdateMemberRequest struct {
	FullName         *string        `json:"full_name"`
	Address          *string        `json:"address"`
	Status           *string        `json:"status" validate:"omitempty,oneof=aktif nonaktif keluar"`
	CustomAttributes datatypes.JSON `json:"custom_attributes"`
}

type SavingsSummary struct {
	Pokok    int64 `json:"pokok"`
	Wajib    int64 `json:"wajib"`
	Sukarela int64 `json:"sukarela"`
	Total    int64 `json:"total"`
}

type MemberResponse struct {
	ID               string         `json:"id"`
	NIK              string         `json:"nik"`
	FullName         string         `json:"full_name"`
	Address          string         `json:"address"`
	BirthDate        string         `json:"birth_date"`
	Status           string         `json:"status"`
	CustomAttributes datatypes.JSON `json:"custom_attributes"`
	JoinedAt         time.Time      `json:"joined_at"`
	SavingsSummary   SavingsSummary `json:"savings_summary"`
}
