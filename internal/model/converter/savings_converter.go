package converter

import (
	"github.com/wiradana/backend/internal/entity"
	"github.com/wiradana/backend/internal/model"
)

func ToSavingsTransactionResponse(tx *entity.SavingsTransaction) model.SavingsTransactionResponse {
	return model.SavingsTransactionResponse{
		ID:          tx.ID.String(),
		MemberID:    tx.MemberID.String(),
		SavingsType: tx.SavingsType,
		Direction:   tx.Direction,
		Amount:      tx.Amount,
		CreatedAt:   tx.CreatedAt,
	}
}
