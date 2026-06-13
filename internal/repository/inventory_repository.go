package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/wiradana/backend/internal/entity"
	"gorm.io/gorm"
)

var (
	ErrProductNotFound   = errors.New("produk tidak ditemukan")
	ErrFieldDefNotFound  = errors.New("field def tidak ditemukan")
	ErrInsufficientStock = errors.New("stok tidak mencukupi")
)

type ProductWithStock struct {
	entity.InventoryProduct
	Stock int64 `gorm:"column:stock"`
}

type InventoryRepository interface {
	// Field definitions
	CreateFieldDef(ctx context.Context, fd *entity.InventoryFieldDef) error
	FindAllFieldDefs(ctx context.Context, coopID string) ([]entity.InventoryFieldDef, error)
	FindFieldDefByID(ctx context.Context, coopID, id string) (*entity.InventoryFieldDef, error)
	DeleteFieldDef(ctx context.Context, id string) error

	// Products
	CreateProduct(ctx context.Context, p *entity.InventoryProduct) error
	FindAllProducts(ctx context.Context, coopID string) ([]ProductWithStock, error)
	FindProductByID(ctx context.Context, coopID, id string) (*ProductWithStock, error)
	UpdateProduct(ctx context.Context, p *entity.InventoryProduct) error

	// Movements (INSERT-ONLY)
	CreateMovement(ctx context.Context, m *entity.InventoryMovement) error
	GetStock(ctx context.Context, productID string) (float64, error)
}

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

// ── Field Definitions ──────────────────────────────────────────────────────────

func (r *inventoryRepository) CreateFieldDef(ctx context.Context, fd *entity.InventoryFieldDef) error {
	return r.db.WithContext(ctx).Create(fd).Error
}

func (r *inventoryRepository) FindAllFieldDefs(ctx context.Context, coopID string) ([]entity.InventoryFieldDef, error) {
	var defs []entity.InventoryFieldDef
	err := r.db.WithContext(ctx).
		Where("cooperative_id = ?", coopID).
		Order("sort_order ASC").
		Find(&defs).Error
	return defs, err
}

func (r *inventoryRepository) FindFieldDefByID(ctx context.Context, coopID, id string) (*entity.InventoryFieldDef, error) {
	var fd entity.InventoryFieldDef
	err := r.db.WithContext(ctx).
		Where("id = ? AND cooperative_id = ?", id, coopID).
		First(&fd).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFieldDefNotFound
		}
		return nil, err
	}
	return &fd, nil
}

func (r *inventoryRepository) DeleteFieldDef(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.InventoryFieldDef{}).Error
}

// ── Products ───────────────────────────────────────────────────────────────────

func (r *inventoryRepository) CreateProduct(ctx context.Context, p *entity.InventoryProduct) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *inventoryRepository) FindAllProducts(ctx context.Context, coopID string) ([]ProductWithStock, error) {
	var products []ProductWithStock
	err := r.db.WithContext(ctx).Raw(`
		SELECT p.*,
			COALESCE(
				(
					SELECT
						COALESCE(SUM(CASE WHEN m.direction = 'masuk' THEN m.quantity ELSE 0 END), 0)
					  - COALESCE(SUM(CASE WHEN m.direction = 'keluar' THEN m.quantity ELSE 0 END), 0)
					FROM inventory_movement m
					WHERE m.product_id = p.id
				),
				0
			)::bigint AS stock
		FROM inventory_product p
		WHERE p.cooperative_id = ?
		ORDER BY p.name ASC
	`, coopID).Scan(&products).Error
	return products, err
}

func (r *inventoryRepository) FindProductByID(ctx context.Context, coopID, id string) (*ProductWithStock, error) {
	var p ProductWithStock
	err := r.db.WithContext(ctx).Raw(`
		SELECT p.*,
			COALESCE(
				(
					SELECT
						COALESCE(SUM(CASE WHEN m.direction = 'masuk' THEN m.quantity ELSE 0 END), 0)
					  - COALESCE(SUM(CASE WHEN m.direction = 'keluar' THEN m.quantity ELSE 0 END), 0)
					FROM inventory_movement m
					WHERE m.product_id = p.id
				),
				0
			)::bigint AS stock
		FROM inventory_product p
		WHERE p.id = ? AND p.cooperative_id = ?
	`, id, coopID).Scan(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	if p.ID == uuid.Nil {
		return nil, ErrProductNotFound
	}
	return &p, nil
}

func (r *inventoryRepository) UpdateProduct(ctx context.Context, p *entity.InventoryProduct) error {
	return r.db.WithContext(ctx).Save(p).Error
}

// ── Movements (INSERT-ONLY) ────────────────────────────────────────────────────

func (r *inventoryRepository) CreateMovement(ctx context.Context, m *entity.InventoryMovement) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *inventoryRepository) GetStock(ctx context.Context, productID string) (float64, error) {
	var stock float64
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			COALESCE(SUM(quantity) FILTER (WHERE direction = 'masuk'), 0)
		  - COALESCE(SUM(quantity) FILTER (WHERE direction = 'keluar'), 0)
		FROM inventory_movement
		WHERE product_id = ?
	`, productID).Scan(&stock).Error
	return stock, err
}
