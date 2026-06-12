package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/usecase"
)

type ModuleController struct {
	moduleUC usecase.ModuleUsecase
	validate *validator.Validate
}

func NewModuleController(moduleUC usecase.ModuleUsecase, validate *validator.Validate) *ModuleController {
	return &ModuleController{moduleUC: moduleUC, validate: validate}
}

func (ctrl *ModuleController) List(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)

	modules, err := ctrl.moduleUC.List(c.Context(), coopID)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "gagal mengambil modul")
	}
	return OKList(c, modules, int64(len(modules)))
}

func (ctrl *ModuleController) Update(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	key := c.Params("key")

	var req model.UpdateModuleRequest
	if err := c.BodyParser(&req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "request body tidak valid")
	}

	resp, err := ctrl.moduleUC.Update(c.Context(), coopID, key, &req)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "gagal update modul")
	}
	return OK(c, resp)
}
