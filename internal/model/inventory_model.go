package model

import "gorm.io/datatypes"

// ---- Field Defs ----

type CreateFieldDefRequest struct {
	FieldKey  string         `json:"field_key"  validate:"required"`
	Label     string         `json:"label"      validate:"required"`
	DataType  string         `json:"data_type"  validate:"required,oneof=text number select date"`
	Options   datatypes.JSON `json:"options"`
	Required  *bool          `json:"required"`
	SortOrder int            `json:"sort_order"`
}

type FieldDefResponse struct {
	ID        string         `json:"id"`
	FieldKey  string         `json:"field_key"`
	Label     string         `json:"label"`
	DataType  string         `json:"data_type"`
	Options   datatypes.JSON `json:"options"`
	Required  bool           `json:"required"`
	SortOrder int            `json:"sort_order"`
}

// ---- Products ----

type CreateProductRequest struct {
	SKU              string         `json:"sku"               validate:"required"`
	Name             string         `json:"name"              validate:"required"`
	Unit             string         `json:"unit"`
	CustomAttributes datatypes.JSON `json:"custom_attributes"`
}

type UpdateProductRequest struct {
	Name             *string        `json:"name"`
	Unit             *string        `json:"unit"`
	CustomAttributes datatypes.JSON `json:"custom_attributes"`
}

type ProductResponse struct {
	ID               string         `json:"id"`
	SKU              string         `json:"sku"`
	Name             string         `json:"name"`
	Unit             string         `json:"unit"`
	Stock            int64          `json:"stock"`
	CustomAttributes datatypes.JSON `json:"custom_attributes"`
	CreatedAt        string         `json:"created_at"`
}

// ---- Movements ----

type RecordMovementRequest struct {
	Direction        string         `json:"direction"         validate:"required,oneof=masuk keluar"`
	Quantity         int64          `json:"quantity"          validate:"required,gt=0"`
	Note             string         `json:"note"`
	CustomAttributes datatypes.JSON `json:"custom_attributes"`
}

type MovementResponse struct {
	ID               string         `json:"id"`
	ProductID        string         `json:"product_id"`
	Direction        string         `json:"direction"`
	Quantity         int64          `json:"quantity"`
	Note             string         `json:"note"`
	CustomAttributes datatypes.JSON `json:"custom_attributes"`
	CreatedAt        string         `json:"created_at"`
}
