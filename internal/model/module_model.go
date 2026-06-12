package model

type ModuleResponse struct {
	Key     string `json:"key"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

type UpdateModuleRequest struct {
	Enabled bool `json:"enabled"`
}

// ── Inventory Models ──

type InventoryFieldDefResponse struct {
	ID        string `json:"id"`
	FieldKey  string `json:"field_key"`
	Label     string `json:"label"`
	DataType  string `json:"data_type"`
	Options   any    `json:"options"`
	Required  bool   `json:"required"`
	SortOrder int    `json:"sort_order"`
}

type CreateInventoryFieldDefRequest struct {
	FieldKey  string `json:"field_key" validate:"required"`
	Label     string `json:"label" validate:"required"`
	DataType  string `json:"data_type" validate:"required,oneof=text number date boolean select"`
	Options   any    `json:"options"`
	Required  bool   `json:"required"`
	SortOrder int    `json:"sort_order"`
}

type InventoryProductResponse struct {
	ID               string `json:"id"`
	SKU              string `json:"sku"`
	Name             string `json:"name"`
	Unit             string `json:"unit"`
	CustomAttributes any    `json:"custom_attributes"`
	Stock            int64  `json:"stock"`
	CreatedAt        string `json:"created_at"`
}

type CreateInventoryProductRequest struct {
	SKU              string `json:"sku" validate:"required"`
	Name             string `json:"name" validate:"required"`
	Unit             string `json:"unit" validate:"required"`
	CustomAttributes any    `json:"custom_attributes"`
}

type UpdateInventoryProductRequest struct {
	SKU              string `json:"sku"`
	Name             string `json:"name"`
	Unit             string `json:"unit"`
	CustomAttributes any    `json:"custom_attributes"`
}

type InventoryMovementResponse struct {
	ID               string `json:"id"`
	ProductID        string `json:"product_id"`
	Direction        string `json:"direction"`
	Quantity         float64 `json:"quantity"`
	Note             string `json:"note"`
	CustomAttributes any    `json:"custom_attributes"`
	CreatedAt        string `json:"created_at"`
}

type CreateInventoryMovementRequest struct {
	Direction        string  `json:"direction" validate:"required,oneof=masuk keluar"`
	Quantity         float64 `json:"quantity" validate:"required,gt=0"`
	Note             string  `json:"note"`
	CustomAttributes any     `json:"custom_attributes"`
}
