package converter

import (
	"github.com/wiradana/backend/internal/entity"
	"github.com/wiradana/backend/internal/model"
)

func ToSavingsResponse(e *entity.SavingsTransaction) model.SavingsTransactionResponse {
	return model.SavingsTransactionResponse{
		ID:          e.ID.String(),
		MemberID:    e.MemberID.String(),
		SavingsType: e.SavingsType,
		Direction:   e.Direction,
		Amount:      e.Amount,
		CreatedAt:   e.CreatedAt,
	}
}
