package converter

import (
	"github.com/wiradana/backend/internal/entity"
	"github.com/wiradana/backend/internal/model"
	"gorm.io/datatypes"
)

func ToFieldDefResponse(f *entity.InventoryFieldDef) model.FieldDefResponse {
	opts := f.Options
	if opts == nil {
		opts = datatypes.JSON("[]")
	}
	return model.FieldDefResponse{
		ID:        f.ID.String(),
		FieldKey:  f.FieldKey,
		Label:     f.Label,
		DataType:  f.DataType,
		Options:   opts,
		Required:  f.Required,
		SortOrder: f.SortOrder,
	}
}

func ToProductResponse(p *entity.InventoryProduct, stock int64) model.ProductResponse {
	attrs := p.CustomAttributes
	if attrs == nil {
		attrs = datatypes.JSON("{}")
	}
	return model.ProductResponse{
		ID:               p.ID.String(),
		SKU:              p.SKU,
		Name:             p.Name,
		Unit:             p.Unit,
		Stock:            stock,
		CustomAttributes: attrs,
		CreatedAt:        p.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func ToMovementResponse(m *entity.InventoryMovement) model.MovementResponse {
	attrs := m.CustomAttributes
	if attrs == nil {
		attrs = datatypes.JSON("{}")
	}
	return model.MovementResponse{
		ID:               m.ID.String(),
		ProductID:        m.ProductID.String(),
		Direction:        m.Direction,
		Quantity:         m.Quantity,
		Note:             m.Note,
		CustomAttributes: attrs,
		CreatedAt:        m.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
