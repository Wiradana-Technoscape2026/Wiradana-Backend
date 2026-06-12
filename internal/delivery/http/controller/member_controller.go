package controller

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/usecase"
)

type MemberController struct {
	memberUC usecase.MemberUsecase
	validate *validator.Validate
}

func NewMemberController(memberUC usecase.MemberUsecase, validate *validator.Validate) *MemberController {
	return &MemberController{memberUC: memberUC, validate: validate}
}

func (ctrl *MemberController) Create(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)

	var req model.CreateMemberRequest
	if err := c.BodyParser(&req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "request body tidak valid")
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}

	member, err := ctrl.memberUC.Create(c.Context(), coopID, &req)
	if err != nil {
		if errors.Is(err, usecase.ErrDuplicateNIK) {
			return Fail(c, fiber.StatusConflict, "CONFLICT", "NIK sudah terdaftar di koperasi ini")
		}
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "terjadi kesalahan server")
	}

	return OK(c, member)
}

func (ctrl *MemberController) List(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	search := c.Query("search")
	status := c.Query("status")

	members, err := ctrl.memberUC.FindAll(c.Context(), coopID, search, status)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "terjadi kesalahan server")
	}

	return OKList(c, members, int64(len(members)))
}

func (ctrl *MemberController) GetByID(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	memberID := c.Params("id")

	member, err := ctrl.memberUC.FindByID(c.Context(), coopID, memberID)
	if err != nil {
		if errors.Is(err, usecase.ErrMemberNotFound) {
			return Fail(c, fiber.StatusNotFound, "NOT_FOUND", "anggota tidak ditemukan")
		}
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "terjadi kesalahan server")
	}

	return OK(c, member)
}

func (ctrl *MemberController) Update(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	memberID := c.Params("id")

	var req model.UpdateMemberRequest
	if err := c.BodyParser(&req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "request body tidak valid")
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}

	member, err := ctrl.memberUC.Update(c.Context(), coopID, memberID, &req)
	if err != nil {
		if errors.Is(err, usecase.ErrMemberNotFound) {
			return Fail(c, fiber.StatusNotFound, "NOT_FOUND", "anggota tidak ditemukan")
		}
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "terjadi kesalahan server")
	}

	return OK(c, member)
}
