package repository

import (
	"context"
	"errors"

	"github.com/wiradana/backend/internal/entity"
	"gorm.io/gorm"
)

var (
	ErrProductNotFound    = errors.New("produk tidak ditemukan")
	ErrFieldDefNotFound   = errors.New("definisi field tidak ditemukan")
	ErrInsufficientStock  = errors.New("stok tidak mencukupi untuk pengeluaran")
	ErrDuplicateSKU       = errors.New("SKU sudah digunakan")
)

type InventoryRepository interface {
	// Field Defs
	CreateFieldDef(ctx context.Context, fd *entity.InventoryFieldDef) error
	FindFieldDefs(ctx context.Context, coopID string) ([]entity.InventoryFieldDef, error)
	DeleteFieldDef(ctx context.Context, coopID, fieldDefID string) error

	// Products
	CreateProduct(ctx context.Context, p *entity.InventoryProduct) error
	FindProducts(ctx context.Context, coopID string) ([]ProductWithStock, error)
	FindProductByID(ctx context.Context, coopID, productID string) (*ProductWithStock, error)
	UpdateProduct(ctx context.Context, p *entity.InventoryProduct) error
	GetProductStock(ctx context.Context, productID string) (int64, error)

	// Movements
	CreateMovement(ctx context.Context, m *entity.InventoryMovement) error
	FindMovementsByProduct(ctx context.Context, productID string) ([]entity.InventoryMovement, error)
}

type ProductWithStock struct {
	entity.InventoryProduct
	Stock int64 `gorm:"column:stock"`
}

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

// ---- Field Defs ----

func (r *inventoryRepository) CreateFieldDef(ctx context.Context, fd *entity.InventoryFieldDef) error {
	return r.db.WithContext(ctx).Create(fd).Error
}

func (r *inventoryRepository) FindFieldDefs(ctx context.Context, coopID string) ([]entity.InventoryFieldDef, error) {
	var fds []entity.InventoryFieldDef
	err := r.db.WithContext(ctx).
		Where("cooperative_id = ?", coopID).
		Order("sort_order ASC").
		Find(&fds).Error
	return fds, err
}

func (r *inventoryRepository) DeleteFieldDef(ctx context.Context, coopID, fieldDefID string) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND cooperative_id = ?", fieldDefID, coopID).
		Delete(&entity.InventoryFieldDef{})
	if result.RowsAffected == 0 {
		return ErrFieldDefNotFound
	}
	return result.Error
}

// ---- Products ----

func (r *inventoryRepository) CreateProduct(ctx context.Context, p *entity.InventoryProduct) error {
	err := r.db.WithContext(ctx).Create(p).Error
	if err != nil && (isDuplicateKey(err)) {
		return ErrDuplicateSKU
	}
	return err
}

func (r *inventoryRepository) FindProducts(ctx context.Context, coopID string) ([]ProductWithStock, error) {
	var rows []ProductWithStock
	err := r.db.WithContext(ctx).Raw(`
		SELECT p.*,
			COALESCE(
				(SELECT SUM(CASE WHEN m.direction='masuk' THEN m.quantity ELSE -m.quantity END)
				 FROM inventory_movement m WHERE m.product_id = p.id), 0) AS stock
		FROM inventory_product p
		WHERE p.cooperative_id = ?
		ORDER BY p.name ASC`, coopID).Scan(&rows).Error
	return rows, err
}

func (r *inventoryRepository) FindProductByID(ctx context.Context, coopID, productID string) (*ProductWithStock, error) {
	var row ProductWithStock
	err := r.db.WithContext(ctx).Raw(`
		SELECT p.*,
			COALESCE(
				(SELECT SUM(CASE WHEN m.direction='masuk' THEN m.quantity ELSE -m.quantity END)
				 FROM inventory_movement m WHERE m.product_id = p.id), 0) AS stock
		FROM inventory_product p
		WHERE p.id = ? AND p.cooperative_id = ?`, productID, coopID).Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.ID.String() == "00000000-0000-0000-0000-000000000000" {
		return nil, ErrProductNotFound
	}
	return &row, nil
}

func (r *inventoryRepository) UpdateProduct(ctx context.Context, p *entity.InventoryProduct) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *inventoryRepository) GetProductStock(ctx context.Context, productID string) (int64, error) {
	var stock int64
	err := r.db.WithContext(ctx).Raw(`
		SELECT COALESCE(
			SUM(CASE WHEN direction='masuk' THEN quantity ELSE -quantity END), 0)
		FROM inventory_movement WHERE product_id = ?`, productID).Scan(&stock).Error
	return stock, err
}

// ---- Movements ----

func (r *inventoryRepository) CreateMovement(ctx context.Context, m *entity.InventoryMovement) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *inventoryRepository) FindMovementsByProduct(ctx context.Context, productID string) ([]entity.InventoryMovement, error) {
	var movements []entity.InventoryMovement
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("created_at DESC").
		Find(&movements).Error
	return movements, err
}

// Helper
func isDuplicateKey(err error) bool {
	return err != nil && (contains(err.Error(), "23505") || contains(err.Error(), "duplicate key"))
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || findSubstr(s, sub))
}

func findSubstr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
