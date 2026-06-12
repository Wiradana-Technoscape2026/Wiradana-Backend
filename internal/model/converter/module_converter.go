package converter

import (
	"github.com/wiradana/backend/internal/entity"
	"github.com/wiradana/backend/internal/model"
)

var moduleNameMap = map[string]string{
	"simpan_pinjam": "Simpan Pinjam",
	"inventory":     "Inventory",
}

func ToModuleResponse(m *entity.CoopModule) model.ModuleResponse {
	name := moduleNameMap[m.ModuleKey]
	if name == "" {
		name = m.ModuleKey
	}
	return model.ModuleResponse{
		Key:     m.ModuleKey,
		Name:    name,
		Enabled: m.Enabled,
	}
}
