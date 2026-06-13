package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/wiradana/backend/internal/entity"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/model/converter"
	"github.com/wiradana/backend/internal/repository"
)

var (
	ErrPokokAlreadyRecorded   = errors.New("simpanan pokok sudah pernah disetor")
	ErrCannotWithdrawMandatory = errors.New("simpanan pokok dan wajib tidak dapat ditarik")
	ErrInsufficientBalance    = errors.New("saldo sukarela tidak mencukupi")
)

type SavingsUsecase interface {
	Record(ctx context.Context, coopID, memberID, recordedByUserID string, req *model.CreateSavingsRequest) (*model.SavingsTransactionResponse, error)
	FindByMember(ctx context.Context, coopID, memberID string) ([]model.SavingsTransactionResponse, error)
	FindByMemberWithRecorder(ctx context.Context, coopID, memberID string) ([]model.SavingsTransactionResponse, error)
	GetCoopSummary(ctx context.Context, coopID string) (*model.SavingsSummaryResponse, error)
	FindAllRecent(ctx context.Context, coopID, savingsType, direction string, limit, offset int) ([]model.SavingsTransactionWithMemberResponse, int64, error)
}

type savingsUsecase struct {
	savingsRepo repository.SavingsRepository
	memberRepo  repository.MemberRepository
}

func NewSavingsUsecase(savingsRepo repository.SavingsRepository, memberRepo repository.MemberRepository) SavingsUsecase {
	return &savingsUsecase{savingsRepo: savingsRepo, memberRepo: memberRepo}
}

func (u *savingsUsecase) Record(ctx context.Context, coopID, memberID, recordedByUserID string, req *model.CreateSavingsRequest) (*model.SavingsTransactionResponse, error) {
	if _, err := u.memberRepo.FindByID(ctx, coopID, memberID); err != nil {
		if errors.Is(err, repository.ErrMemberNotFound) {
			return nil, ErrMemberNotFound
		}
		return nil, err
	}

	switch {
	case req.SavingsType == "pokok" && req.Direction == "setor":
		count, err := u.savingsRepo.CountPokok(ctx, memberID, coopID)
		if err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, ErrPokokAlreadyRecorded
		}

	case (req.SavingsType == "pokok" || req.SavingsType == "wajib") && req.Direction == "tarik":
		return nil, ErrCannotWithdrawMandatory

	case req.SavingsType == "sukarela" && req.Direction == "tarik":
		balance, err := u.savingsRepo.GetSukarelaBalance(ctx, memberID, coopID)
		if err != nil {
			return nil, err
		}
		if balance < req.Amount {
			return nil, ErrInsufficientBalance
		}
	}

	coopUUID, err := uuid.Parse(coopID)
	if err != nil {
		return nil, errors.New("cooperative_id tidak valid")
	}
	memberUUID, err := uuid.Parse(memberID)
	if err != nil {
		return nil, errors.New("member_id tidak valid")
	}

	tx := &entity.SavingsTransaction{
		CooperativeID: coopUUID,
		MemberID:      memberUUID,
		SavingsType:   req.SavingsType,
		Direction:     req.Direction,
		Amount:        req.Amount,
	}
	// Simpan siapa yang mencatat jika user_id tersedia
	if recordedByUserID != "" {
		if uid, err := uuid.Parse(recordedByUserID); err == nil {
			tx.RecordedBy = &uid
		}
	}

	if err := u.savingsRepo.Create(ctx, tx); err != nil {
		return nil, err
	}

	resp := converter.ToSavingsResponse(tx)
	return &resp, nil
}

func (u *savingsUsecase) FindByMember(ctx context.Context, coopID, memberID string) ([]model.SavingsTransactionResponse, error) {
	if _, err := u.memberRepo.FindByID(ctx, coopID, memberID); err != nil {
		if errors.Is(err, repository.ErrMemberNotFound) {
			return nil, ErrMemberNotFound
		}
		return nil, err
	}

	txs, err := u.savingsRepo.FindByMember(ctx, coopID, memberID)
	if err != nil {
		return nil, err
	}

	responses := make([]model.SavingsTransactionResponse, 0, len(txs))
	for _, tx := range txs {
		responses = append(responses, converter.ToSavingsResponse(tx))
	}
	return responses, nil
}

// FindByMemberWithRecorder — dipakai oleh portal anggota untuk tampilkan nama pencatat.
func (u *savingsUsecase) FindByMemberWithRecorder(ctx context.Context, coopID, memberID string) ([]model.SavingsTransactionResponse, error) {
	rows, err := u.savingsRepo.FindByMemberWithRecorder(ctx, coopID, memberID)
	if err != nil {
		return nil, err
	}
	responses := make([]model.SavingsTransactionResponse, 0, len(rows))
	for _, r := range rows {
		responses = append(responses, model.SavingsTransactionResponse{
			ID:             r.ID,
			MemberID:       r.MemberID,
			SavingsType:    r.SavingsType,
			Direction:      r.Direction,
			Amount:         r.Amount,
			RecordedByName: r.RecordedByName,
			CreatedAt:      r.CreatedAt,
		})
	}
	return responses, nil
}

func (u *savingsUsecase) GetCoopSummary(ctx context.Context, coopID string) (*model.SavingsSummaryResponse, error) {
	row, err := u.savingsRepo.GetCoopSummary(ctx, coopID)
	if err != nil {
		return nil, err
	}
	return &model.SavingsSummaryResponse{
		Pokok:    row.Pokok,
		Wajib:    row.Wajib,
		Sukarela: row.Sukarela,
		Total:    row.Pokok + row.Wajib + row.Sukarela,
	}, nil
}

func (u *savingsUsecase) FindAllRecent(ctx context.Context, coopID, savingsType, direction string, limit, offset int) ([]model.SavingsTransactionWithMemberResponse, int64, error) {
	rows, total, err := u.savingsRepo.FindAllRecent(ctx, coopID, savingsType, direction, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]model.SavingsTransactionWithMemberResponse, 0, len(rows))
	for _, r := range rows {
		responses = append(responses, model.SavingsTransactionWithMemberResponse{
			ID:             r.ID,
			MemberID:       r.MemberID,
			MemberName:     r.MemberName,
			SavingsType:    r.SavingsType,
			Direction:      r.Direction,
			Amount:         r.Amount,
			RecordedByName: r.RecordedByName,
			CreatedAt:      r.CreatedAt,
		})
	}
	return responses, total, nil
}
