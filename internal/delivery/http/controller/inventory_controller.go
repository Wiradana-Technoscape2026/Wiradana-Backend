package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/usecase"
)

type InventoryController struct {
	inventoryUC usecase.InventoryUsecase
	validate    *validator.Validate
}

func NewInventoryController(inventoryUC usecase.InventoryUsecase, validate *validator.Validate) *InventoryController {
	return &InventoryController{inventoryUC: inventoryUC, validate: validate}
}

// ── Field Definitions ──────────────────────────────────────────────────────────

func (ctrl *InventoryController) ListFieldDefs(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	defs, err := ctrl.inventoryUC.ListFieldDefs(c.Context(), coopID)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "gagal mengambil field defs")
	}
	return OKList(c, defs, int64(len(defs)))
}

func (ctrl *InventoryController) CreateFieldDef(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	var req model.CreateInventoryFieldDefRequest
	if err := c.BodyParser(&req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "request body tidak valid")
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}
	resp, err := ctrl.inventoryUC.CreateFieldDef(c.Context(), coopID, &req)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "gagal membuat field def")
	}
	return OK(c, resp)
}

func (ctrl *InventoryController) DeleteFieldDef(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	id := c.Params("id")
	if err := ctrl.inventoryUC.DeleteFieldDef(c.Context(), coopID, id); err != nil {
		return Fail(c, fiber.StatusNotFound, "NOT_FOUND", "field def tidak ditemukan")
	}
	return OK(c, fiber.Map{"deleted": true})
}

// ── Products ───────────────────────────────────────────────────────────────────

func (ctrl *InventoryController) ListProducts(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	products, err := ctrl.inventoryUC.ListProducts(c.Context(), coopID)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "gagal mengambil produk")
	}
	return OKList(c, products, int64(len(products)))
}

func (ctrl *InventoryController) CreateProduct(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	var req model.CreateInventoryProductRequest
	if err := c.BodyParser(&req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "request body tidak valid")
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}
	resp, err := ctrl.inventoryUC.CreateProduct(c.Context(), coopID, &req)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "gagal membuat produk")
	}
	return OK(c, resp)
}

func (ctrl *InventoryController) UpdateProduct(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	id := c.Params("id")
	var req model.UpdateInventoryProductRequest
	if err := c.BodyParser(&req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "request body tidak valid")
	}
	resp, err := ctrl.inventoryUC.UpdateProduct(c.Context(), coopID, id, &req)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "gagal update produk")
	}
	return OK(c, resp)
}

// ── Movements ──────────────────────────────────────────────────────────────────

func (ctrl *InventoryController) RecordMovement(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	productID := c.Params("id")
	var req model.CreateInventoryMovementRequest
	if err := c.BodyParser(&req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "request body tidak valid")
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}
	resp, err := ctrl.inventoryUC.RecordMovement(c.Context(), coopID, productID, &req)
	if err != nil {
		if err == usecase.ErrStockInsufficient {
			return Fail(c, fiber.StatusConflict, "CONFLICT", "stok tidak mencukupi untuk pengeluaran")
		}
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "gagal mencatat pergerakan")
	}
	return OK(c, resp)
}
