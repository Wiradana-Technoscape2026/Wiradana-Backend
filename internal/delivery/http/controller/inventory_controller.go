package controller

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/usecase"
)

type InventoryController struct {
	uc       usecase.InventoryUsecase
	validate *validator.Validate
}

func NewInventoryController(uc usecase.InventoryUsecase, validate *validator.Validate) *InventoryController {
	return &InventoryController{uc: uc, validate: validate}
}

// ---- Field Defs ----

func (ctrl *InventoryController) ListFieldDefs(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	fds, err := ctrl.uc.ListFieldDefs(c.Context(), coopID)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
	}
	return OKList(c, fds, int64(len(fds)))
}

func (ctrl *InventoryController) CreateFieldDef(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	var req model.CreateFieldDefRequest
	if err := c.BodyParser(&req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "request body tidak valid")
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}
	resp, err := ctrl.uc.CreateFieldDef(c.Context(), coopID, &req)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
	}
	return OK(c, resp)
}

func (ctrl *InventoryController) DeleteFieldDef(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	id := c.Params("id")
	if err := ctrl.uc.DeleteFieldDef(c.Context(), coopID, id); err != nil {
		if errors.Is(err, usecase.ErrFieldDefNotFound) {
			return Fail(c, fiber.StatusNotFound, "NOT_FOUND", err.Error())
		}
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
	}
	return OK(c, fiber.Map{"deleted": true})
}

// ---- Products ----

func (ctrl *InventoryController) ListProducts(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	products, err := ctrl.uc.ListProducts(c.Context(), coopID)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
	}
	return OKList(c, products, int64(len(products)))
}

func (ctrl *InventoryController) CreateProduct(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	var req model.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "request body tidak valid")
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}
	resp, err := ctrl.uc.CreateProduct(c.Context(), coopID, &req)
	if err != nil {
		if errors.Is(err, usecase.ErrDuplicateSKU) {
			return Fail(c, fiber.StatusConflict, "CONFLICT", err.Error())
		}
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
	}
	return OK(c, resp)
}

func (ctrl *InventoryController) UpdateProduct(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	productID := c.Params("id")
	var req model.UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "request body tidak valid")
	}
	resp, err := ctrl.uc.UpdateProduct(c.Context(), coopID, productID, &req)
	if err != nil {
		if errors.Is(err, usecase.ErrProductNotFound) {
			return Fail(c, fiber.StatusNotFound, "NOT_FOUND", err.Error())
		}
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
	}
	return OK(c, resp)
}

// ---- Movements ----

func (ctrl *InventoryController) RecordMovement(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	productID := c.Params("id")
	var req model.RecordMovementRequest
	if err := c.BodyParser(&req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "request body tidak valid")
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}
	resp, err := ctrl.uc.RecordMovement(c.Context(), coopID, productID, &req)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrProductNotFound):
			return Fail(c, fiber.StatusNotFound, "NOT_FOUND", err.Error())
		case errors.Is(err, usecase.ErrInsufficientStock):
			return Fail(c, fiber.StatusConflict, "CONFLICT", err.Error())
		default:
			return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
		}
	}
	return OK(c, resp)
}
