package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/wiradana/backend/internal/entity"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/model/converter"
	"github.com/wiradana/backend/internal/repository"
	"gorm.io/datatypes"
)

var (
	ErrProductNotFound   = errors.New("produk tidak ditemukan")
	ErrInsufficientStock = errors.New("stok tidak mencukupi untuk pengeluaran")
	ErrDuplicateSKU      = errors.New("SKU sudah digunakan")
	ErrFieldDefNotFound  = errors.New("definisi field tidak ditemukan")
)

type InventoryUsecase interface {
	// Field Defs
	ListFieldDefs(ctx context.Context, coopID string) ([]model.FieldDefResponse, error)
	CreateFieldDef(ctx context.Context, coopID string, req *model.CreateFieldDefRequest) (*model.FieldDefResponse, error)
	DeleteFieldDef(ctx context.Context, coopID, fieldDefID string) error

	// Products
	ListProducts(ctx context.Context, coopID string) ([]model.ProductResponse, error)
	CreateProduct(ctx context.Context, coopID string, req *model.CreateProductRequest) (*model.ProductResponse, error)
	UpdateProduct(ctx context.Context, coopID, productID string, req *model.UpdateProductRequest) (*model.ProductResponse, error)

	// Movements
	RecordMovement(ctx context.Context, coopID, productID string, req *model.RecordMovementRequest) (*model.MovementResponse, error)
}

type inventoryUsecase struct {
	repo repository.InventoryRepository
}

func NewInventoryUsecase(repo repository.InventoryRepository) InventoryUsecase {
	return &inventoryUsecase{repo: repo}
}

// ---- Field Defs ----

func (u *inventoryUsecase) ListFieldDefs(ctx context.Context, coopID string) ([]model.FieldDefResponse, error) {
	fds, err := u.repo.FindFieldDefs(ctx, coopID)
	if err != nil {
		return nil, err
	}
	result := make([]model.FieldDefResponse, len(fds))
	for i, f := range fds {
		result[i] = converter.ToFieldDefResponse(&f)
	}
	return result, nil
}

func (u *inventoryUsecase) CreateFieldDef(ctx context.Context, coopID string, req *model.CreateFieldDefRequest) (*model.FieldDefResponse, error) {
	coopUUID, _ := uuid.Parse(coopID)
	opts := req.Options
	if opts == nil {
		opts = datatypes.JSON("[]")
	}
	required := false
	if req.Required != nil {
		required = *req.Required
	}
	fd := &entity.InventoryFieldDef{
		CooperativeID: coopUUID,
		FieldKey:      req.FieldKey,
		Label:         req.Label,
		DataType:      req.DataType,
		Options:       opts,
		Required:      required,
		SortOrder:     req.SortOrder,
	}
	if err := u.repo.CreateFieldDef(ctx, fd); err != nil {
		return nil, err
	}
	resp := converter.ToFieldDefResponse(fd)
	return &resp, nil
}

func (u *inventoryUsecase) DeleteFieldDef(ctx context.Context, coopID, fieldDefID string) error {
	err := u.repo.DeleteFieldDef(ctx, coopID, fieldDefID)
	if errors.Is(err, repository.ErrFieldDefNotFound) {
		return ErrFieldDefNotFound
	}
	return err
}

// ---- Products ----

func (u *inventoryUsecase) ListProducts(ctx context.Context, coopID string) ([]model.ProductResponse, error) {
	products, err := u.repo.FindProducts(ctx, coopID)
	if err != nil {
		return nil, err
	}
	result := make([]model.ProductResponse, len(products))
	for i, p := range products {
		result[i] = converter.ToProductResponse(&p.InventoryProduct, p.Stock)
	}
	return result, nil
}

func (u *inventoryUsecase) CreateProduct(ctx context.Context, coopID string, req *model.CreateProductRequest) (*model.ProductResponse, error) {
	coopUUID, _ := uuid.Parse(coopID)
	attrs := req.CustomAttributes
	if attrs == nil {
		attrs = datatypes.JSON("{}")
	}
	unit := req.Unit
	if unit == "" {
		unit = "pcs"
	}
	p := &entity.InventoryProduct{
		CooperativeID:    coopUUID,
		SKU:              req.SKU,
		Name:             req.Name,
		Unit:             unit,
		CustomAttributes: attrs,
	}
	if err := u.repo.CreateProduct(ctx, p); err != nil {
		if errors.Is(err, repository.ErrDuplicateSKU) {
			return nil, ErrDuplicateSKU
		}
		return nil, err
	}
	resp := converter.ToProductResponse(p, 0)
	return &resp, nil
}

func (u *inventoryUsecase) UpdateProduct(ctx context.Context, coopID, productID string, req *model.UpdateProductRequest) (*model.ProductResponse, error) {
	pws, err := u.repo.FindProductByID(ctx, coopID, productID)
	if err != nil {
		return nil, ErrProductNotFound
	}
	p := &pws.InventoryProduct
	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.Unit != nil {
		p.Unit = *req.Unit
	}
	if req.CustomAttributes != nil {
		p.CustomAttributes = req.CustomAttributes
	}
	if err := u.repo.UpdateProduct(ctx, p); err != nil {
		return nil, err
	}
	resp := converter.ToProductResponse(p, pws.Stock)
	return &resp, nil
}

// ---- Movements ----

func (u *inventoryUsecase) RecordMovement(ctx context.Context, coopID, productID string, req *model.RecordMovementRequest) (*model.MovementResponse, error) {
	// Validate product exists
	_, err := u.repo.FindProductByID(ctx, coopID, productID)
	if err != nil {
		return nil, ErrProductNotFound
	}

	// If keluar, check stock
	if req.Direction == "keluar" {
		stock, _ := u.repo.GetProductStock(ctx, productID)
		if req.Quantity > stock {
			return nil, ErrInsufficientStock
		}
	}

	coopUUID, _ := uuid.Parse(coopID)
	productUUID, _ := uuid.Parse(productID)
	attrs := req.CustomAttributes
	if attrs == nil {
		attrs = datatypes.JSON("{}")
	}

	m := &entity.InventoryMovement{
		CooperativeID:    coopUUID,
		ProductID:        productUUID,
		Direction:        req.Direction,
		Quantity:         req.Quantity,
		Note:             req.Note,
		CustomAttributes: attrs,
	}

	if err := u.repo.CreateMovement(ctx, m); err != nil {
		return nil, err
	}

	resp := converter.ToMovementResponse(m)
	return &resp, nil
}
