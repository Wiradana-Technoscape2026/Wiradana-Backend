package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/wiradana/backend/internal/entity"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/repository"
)

type ModuleUsecase interface {
	List(ctx context.Context, coopID string) ([]model.ModuleResponse, error)
	Update(ctx context.Context, coopID, key string, req *model.UpdateModuleRequest) (*model.ModuleResponse, error)
}

type moduleUsecase struct {
	moduleRepo repository.ModuleRepository
}

func NewModuleUsecase(moduleRepo repository.ModuleRepository) ModuleUsecase {
	return &moduleUsecase{moduleRepo: moduleRepo}
}

var defaultModules = []struct {
	Key     string
	Name    string
	Enabled bool
}{
	{Key: "simpan_pinjam", Name: "Simpan Pinjam", Enabled: true},
	{Key: "inventory", Name: "Inventory", Enabled: false},
}

func (u *moduleUsecase) List(ctx context.Context, coopID string) ([]model.ModuleResponse, error) {
	modules, err := u.moduleRepo.FindAll(ctx, coopID)
	if err != nil {
		return nil, err
	}

	// Merge with defaults
	existingMap := make(map[string]bool)
	for _, m := range modules {
		existingMap[m.ModuleKey] = m.Enabled
	}

	result := make([]model.ModuleResponse, 0, len(defaultModules))
	for _, dm := range defaultModules {
		enabled := dm.Enabled
		if _, ok := existingMap[dm.Key]; ok {
			enabled = existingMap[dm.Key]
		}
		result = append(result, model.ModuleResponse{
			Key:     dm.Key,
			Name:    dm.Name,
			Enabled: enabled,
		})
	}
	return result, nil
}

func (u *moduleUsecase) Update(ctx context.Context, coopID, key string, req *model.UpdateModuleRequest) (*model.ModuleResponse, error) {
	coopUUID, err := uuid.Parse(coopID)
	if err != nil {
		return nil, errors.New("cooperative_id tidak valid")
	}

	// Find module name from defaults
	name := key
	for _, dm := range defaultModules {
		if dm.Key == key {
			name = dm.Name
			break
		}
	}

	m := &entity.CoopModule{
		CooperativeID: coopUUID,
		ModuleKey:     key,
		Enabled:       req.Enabled,
	}

	if err := u.moduleRepo.Upsert(ctx, m); err != nil {
		return nil, err
	}

	return &model.ModuleResponse{
		Key:     key,
		Name:    name,
		Enabled: req.Enabled,
	}, nil
}
