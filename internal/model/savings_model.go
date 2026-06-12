package model

import "time"

type CreateSavingsRequest struct {
	SavingsType string `json:"savings_type" validate:"required,oneof=pokok wajib sukarela"`
	Direction   string `json:"direction"    validate:"required,oneof=setor tarik"`
	Amount      int64  `json:"amount"       validate:"required,gt=0"`
}

type SavingsTransactionResponse struct {
	ID          string    `json:"id"`
	MemberID    string    `json:"member_id"`
	SavingsType string    `json:"savings_type"`
	Direction   string    `json:"direction"`
	Amount      int64     `json:"amount"`
	CreatedAt   time.Time `json:"created_at"`
}
