package converter

import (
	"github.com/wiradana/backend/internal/entity"
	"github.com/wiradana/backend/internal/model"
)

<<<<<<< HEAD
func ToSavingsTransactionResponse(tx *entity.SavingsTransaction) model.SavingsTransactionResponse {
	return model.SavingsTransactionResponse{
		ID:          tx.ID.String(),
		MemberID:    tx.MemberID.String(),
		SavingsType: tx.SavingsType,
		Direction:   tx.Direction,
		Amount:      tx.Amount,
		CreatedAt:   tx.CreatedAt,
=======
func ToSavingsResponse(e *entity.SavingsTransaction) model.SavingsTransactionResponse {
	return model.SavingsTransactionResponse{
		ID:          e.ID.String(),
		MemberID:    e.MemberID.String(),
		SavingsType: e.SavingsType,
		Direction:   e.Direction,
		Amount:      e.Amount,
		CreatedAt:   e.CreatedAt,
>>>>>>> e6c7f422c936b4876b95b9366e0dc7eebfff82ed
	}
}
