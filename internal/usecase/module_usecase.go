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

var ErrModuleNotFound = errors.New("modul tidak ditemukan")

type ModuleUsecase interface {
	List(ctx context.Context, coopID string) ([]model.ModuleResponse, error)
	Update(ctx context.Context, coopID, key string, req *model.UpdateModuleRequest) (*model.ModuleResponse, error)
}

type moduleUsecase struct {
	repo repository.ModuleRepository
}

func NewModuleUsecase(repo repository.ModuleRepository) ModuleUsecase {
	return &moduleUsecase{repo: repo}
}

func (u *moduleUsecase) List(ctx context.Context, coopID string) ([]model.ModuleResponse, error) {
	mods, err := u.repo.FindAll(ctx, coopID)
	if err != nil {
		return nil, err
	}
	// If no modules, seed default
	if len(mods) == 0 {
		coopUUID, _ := uuid.Parse(coopID)
		defaults := []entity.CoopModule{
			{CooperativeID: coopUUID, ModuleKey: "simpan_pinjam", Enabled: true},
			{CooperativeID: coopUUID, ModuleKey: "inventory", Enabled: false},
		}
		for _, d := range defaults {
			_ = u.repo.Upsert(ctx, &d)
		}
		mods, _ = u.repo.FindAll(ctx, coopID)
	}
	result := make([]model.ModuleResponse, len(mods))
	for i, m := range mods {
		result[i] = converter.ToModuleResponse(&m)
	}
	return result, nil
}

func (u *moduleUsecase) Update(ctx context.Context, coopID, key string, req *model.UpdateModuleRequest) (*model.ModuleResponse, error) {
	m, err := u.repo.FindByKey(ctx, coopID, key)
	if err != nil {
		// Auto-create if not exists
		coopUUID, _ := uuid.Parse(coopID)
		m = &entity.CoopModule{CooperativeID: coopUUID, ModuleKey: key}
	}
	if req.Enabled != nil {
		m.Enabled = *req.Enabled
	}
	if err := u.repo.Upsert(ctx, m); err != nil {
		return nil, err
	}
	resp := converter.ToModuleResponse(m)
	return &resp, nil
}
