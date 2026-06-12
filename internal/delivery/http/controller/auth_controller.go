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

	return OK(c, model.LoginResponse{
		Token: result.Token,
		User:  converter.ToAppUserResponse(result.User),
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

	return OK(c, model.LoginResponse{
		Token: result.Token,
		User:  converter.ToAppUserResponse(result.User),
	})
}
