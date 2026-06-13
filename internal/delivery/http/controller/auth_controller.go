package controller

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/model/converter"
	"github.com/wiradana/backend/internal/usecase"
)

type AuthController struct {
	authUC   usecase.AuthUsecase
	validate *validator.Validate
}

func NewAuthController(authUC usecase.AuthUsecase, validate *validator.Validate) *AuthController {
	return &AuthController{authUC: authUC, validate: validate}
}

func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "request body tidak valid")
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}

	result, err := ctrl.authUC.Login(c.Context(), req.Identifier, req.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			return Fail(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "identifier atau password salah")
		}
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "terjadi kesalahan server")
	}

	// Anggota dengan lebih dari satu koperasi perlu memilih koperasi terlebih dahulu.
	if len(result.Memberships) > 0 {
		opts := make([]model.CooperativeOption, len(result.Memberships))
		for i, m := range result.Memberships {
			opts[i] = model.CooperativeOption{
				CooperativeID:   m.CooperativeID,
				CooperativeName: m.CooperativeName,
				MemberID:        m.MemberID,
			}
		}
		user := converter.ToAppUserResponse(result.User)
		return OK(c, model.LoginResponse{
			RequiresCooperativeSelection: true,
			User:                         &user,
			Cooperatives:                 opts,
		})
	}

	user := converter.ToAppUserResponse(result.User)
	return OK(c, model.LoginResponse{
		Token: result.Token,
		User:  &user,
	})
}

func (ctrl *AuthController) SelectCooperative(c *fiber.Ctx) error {
	var req model.SelectCooperativeRequest
	if err := c.BodyParser(&req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "request body tidak valid")
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}

	result, err := ctrl.authUC.SelectCooperative(c.Context(), req.Identifier, req.Password, req.CooperativeID)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidCredentials):
			return Fail(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "identifier atau password salah")
		case errors.Is(err, usecase.ErrNotMemberOfCooperative):
			return Fail(c, fiber.StatusForbidden, "FORBIDDEN", "user bukan anggota koperasi ini")
		default:
			return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "terjadi kesalahan server")
		}
	}

	user := converter.ToAppUserResponse(result.User)
	return OK(c, model.LoginResponse{
		Token: result.Token,
		User:  &user,
	})
}

func (ctrl *AuthController) RegisterPengurus(c *fiber.Ctx) error {
	var req model.RegisterPengurusRequest
	if err := c.BodyParser(&req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "request body tidak valid")
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}

	result, err := ctrl.authUC.RegisterPengurus(c.Context(), &req)
	if err != nil {
		if errors.Is(err, usecase.ErrEmailAlreadyUsed) {
			return Fail(c, fiber.StatusConflict, "CONFLICT", "email sudah terdaftar")
		}
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "terjadi kesalahan server")
	}

	user := converter.ToAppUserResponse(result.User)
	return OK(c, model.LoginResponse{
		Token: result.Token,
		User:  &user,
	})
}
