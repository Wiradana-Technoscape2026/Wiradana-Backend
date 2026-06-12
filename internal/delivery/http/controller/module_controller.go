package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/usecase"
)

type ModuleController struct {
	uc       usecase.ModuleUsecase
	validate *validator.Validate
}

func NewModuleController(uc usecase.ModuleUsecase, validate *validator.Validate) *ModuleController {
	return &ModuleController{uc: uc, validate: validate}
}

func (ctrl *ModuleController) List(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	mods, err := ctrl.uc.List(c.Context(), coopID)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
	}
	return OKList(c, mods, int64(len(mods)))
}

func (ctrl *ModuleController) Update(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	key := c.Params("key")

	var req model.UpdateModuleRequest
	if err := c.BodyParser(&req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "request body tidak valid")
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}

	resp, err := ctrl.uc.Update(c.Context(), coopID, key, &req)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
	}
	return OK(c, resp)
}
