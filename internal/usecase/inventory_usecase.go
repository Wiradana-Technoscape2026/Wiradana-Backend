package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/wiradana/backend/internal/entity"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/repository"
	"gorm.io/datatypes"
)

var (
	ErrStockInsufficient = errors.New("stok tidak mencukupi untuk pengeluaran")
)

type InventoryUsecase interface {
	// Field definitions
	CreateFieldDef(ctx context.Context, coopID string, req *model.CreateInventoryFieldDefRequest) (*model.InventoryFieldDefResponse, error)
	ListFieldDefs(ctx context.Context, coopID string) ([]model.InventoryFieldDefResponse, error)
	DeleteFieldDef(ctx context.Context, coopID, id string) error

	// Products
	CreateProduct(ctx context.Context, coopID string, req *model.CreateInventoryProductRequest) (*model.InventoryProductResponse, error)
	ListProducts(ctx context.Context, coopID string) ([]model.InventoryProductResponse, error)
	UpdateProduct(ctx context.Context, coopID, id string, req *model.UpdateInventoryProductRequest) (*model.InventoryProductResponse, error)

	// Movements
	RecordMovement(ctx context.Context, coopID, productID string, req *model.CreateInventoryMovementRequest) (*model.InventoryMovementResponse, error)
}

type inventoryUsecase struct {
	repo repository.InventoryRepository
}

func NewInventoryUsecase(repo repository.InventoryRepository) InventoryUsecase {
	return &inventoryUsecase{repo: repo}
}

// jsonToJSONB converts an `any` value to datatypes.JSON.
// If val is nil, returns a JSON representation of defaultStr (e.g. "[]" or "{}").
func jsonToJSONB(val any, defaultStr string) datatypes.JSON {
	if val == nil {
		return datatypes.JSON(defaultStr)
	}
	b, err := json.Marshal(val)
	if err != nil {
		return datatypes.JSON(defaultStr)
	}
	return datatypes.JSON(b)
}

// ── Field Definitions ──────────────────────────────────────────────────────────

func (u *inventoryUsecase) CreateFieldDef(ctx context.Context, coopID string, req *model.CreateInventoryFieldDefRequest) (*model.InventoryFieldDefResponse, error) {
	coopUUID, err := uuid.Parse(coopID)
	if err != nil {
		return nil, errors.New("cooperative_id tidak valid")
	}

	opts := jsonToJSONB(req.Options, "[]")

	fd := &entity.InventoryFieldDef{
		CooperativeID: coopUUID,
		FieldKey:      req.FieldKey,
		Label:         req.Label,
		DataType:      req.DataType,
		Options:       opts,
		Required:      req.Required,
		SortOrder:     req.SortOrder,
	}

	if err := u.repo.CreateFieldDef(ctx, fd); err != nil {
		return nil, err
	}

	return &model.InventoryFieldDefResponse{
		ID:        fd.ID.String(),
		FieldKey:  fd.FieldKey,
		Label:     fd.Label,
		DataType:  fd.DataType,
		Options:   fd.Options,
		Required:  fd.Required,
		SortOrder: fd.SortOrder,
	}, nil
}

func (u *inventoryUsecase) ListFieldDefs(ctx context.Context, coopID string) ([]model.InventoryFieldDefResponse, error) {
	defs, err := u.repo.FindAllFieldDefs(ctx, coopID)
	if err != nil {
		return nil, err
	}
	result := make([]model.InventoryFieldDefResponse, 0, len(defs))
	for _, d := range defs {
		result = append(result, model.InventoryFieldDefResponse{
			ID:        d.ID.String(),
			FieldKey:  d.FieldKey,
			Label:     d.Label,
			DataType:  d.DataType,
			Options:   d.Options,
			Required:  d.Required,
			SortOrder: d.SortOrder,
		})
	}
	return result, nil
}

func (u *inventoryUsecase) DeleteFieldDef(ctx context.Context, coopID, id string) error {
	_, err := u.repo.FindFieldDefByID(ctx, coopID, id)
	if err != nil {
		return errors.New("field def tidak ditemukan")
	}
	return u.repo.DeleteFieldDef(ctx, id)
}

// ── Products ───────────────────────────────────────────────────────────────────

func (u *inventoryUsecase) CreateProduct(ctx context.Context, coopID string, req *model.CreateInventoryProductRequest) (*model.InventoryProductResponse, error) {
	coopUUID, err := uuid.Parse(coopID)
	if err != nil {
		return nil, errors.New("cooperative_id tidak valid")
	}

	attrs := jsonToJSONB(req.CustomAttributes, "{}")

	product := &entity.InventoryProduct{
		CooperativeID:    coopUUID,
		SKU:              req.SKU,
		Name:             req.Name,
		Unit:             req.Unit,
		CustomAttributes: attrs,
	}

	if err := u.repo.CreateProduct(ctx, product); err != nil {
		return nil, err
	}

	return &model.InventoryProductResponse{
		ID:               product.ID.String(),
		SKU:              product.SKU,
		Name:             product.Name,
		Unit:             product.Unit,
		CustomAttributes: product.CustomAttributes,
		Stock:            0,
		CreatedAt:        product.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (u *inventoryUsecase) ListProducts(ctx context.Context, coopID string) ([]model.InventoryProductResponse, error) {
	products, err := u.repo.FindAllProducts(ctx, coopID)
	if err != nil {
		return nil, err
	}
	result := make([]model.InventoryProductResponse, 0, len(products))
	for _, p := range products {
		result = append(result, model.InventoryProductResponse{
			ID:               p.ID.String(),
			SKU:              p.SKU,
			Name:             p.Name,
			Unit:             p.Unit,
			CustomAttributes: p.CustomAttributes,
			Stock:            p.Stock,
			CreatedAt:        p.CreatedAt.Format(time.RFC3339),
		})
	}
	return result, nil
}

func (u *inventoryUsecase) UpdateProduct(ctx context.Context, coopID, id string, req *model.UpdateInventoryProductRequest) (*model.InventoryProductResponse, error) {
	existing, err := u.repo.FindProductByID(ctx, coopID, id)
	if err != nil {
		return nil, err
	}

	if req.SKU != "" {
		existing.SKU = req.SKU
	}
	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Unit != "" {
		existing.Unit = req.Unit
	}
	if req.CustomAttributes != nil {
		existing.CustomAttributes = jsonToJSONB(req.CustomAttributes, "{}")
	}

	if err := u.repo.UpdateProduct(ctx, &existing.InventoryProduct); err != nil {
		return nil, err
	}

	return &model.InventoryProductResponse{
		ID:               existing.ID.String(),
		SKU:              existing.SKU,
		Name:             existing.Name,
		Unit:             existing.Unit,
		CustomAttributes: existing.CustomAttributes,
		Stock:            existing.Stock,
		CreatedAt:        existing.CreatedAt.Format(time.RFC3339),
	}, nil
}

// ── Movements (INSERT-ONLY) ────────────────────────────────────────────────────

func (u *inventoryUsecase) RecordMovement(ctx context.Context, coopID, productID string, req *model.CreateInventoryMovementRequest) (*model.InventoryMovementResponse, error) {
	coopUUID, err := uuid.Parse(coopID)
	if err != nil {
		return nil, errors.New("cooperative_id tidak valid")
	}
	prodUUID, err := uuid.Parse(productID)
	if err != nil {
		return nil, errors.New("product_id tidak valid")
	}

	// Check product exists
	_, err = u.repo.FindProductByID(ctx, coopID, productID)
	if err != nil {
		return nil, err
	}

	// Guard: stock check for "keluar"
	if req.Direction == "keluar" {
		currentStock, err := u.repo.GetStock(ctx, productID)
		if err != nil {
			return nil, err
		}
		if req.Quantity > currentStock {
			return nil, ErrStockInsufficient
		}
	}

	attrs := jsonToJSONB(req.CustomAttributes, "{}")

	var note *string
	if req.Note != "" {
		note = &req.Note
	}

	movement := &entity.InventoryMovement{
		CooperativeID:    coopUUID,
		ProductID:        prodUUID,
		Direction:        req.Direction,
		Quantity:         req.Quantity,
		Note:             note,
		CustomAttributes: attrs,
	}

	if err := u.repo.CreateMovement(ctx, movement); err != nil {
		return nil, err
	}

	noteStr := ""
	if movement.Note != nil {
		noteStr = *movement.Note
	}

	return &model.InventoryMovementResponse{
		ID:               movement.ID.String(),
		ProductID:        movement.ProductID.String(),
		Direction:        movement.Direction,
		Quantity:         movement.Quantity,
		Note:             noteStr,
		CustomAttributes: movement.CustomAttributes,
		CreatedAt:        movement.CreatedAt.Format(time.RFC3339),
	}, nil
}
