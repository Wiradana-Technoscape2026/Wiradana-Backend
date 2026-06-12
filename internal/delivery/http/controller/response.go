package controller

import "github.com/gofiber/fiber/v2"

func OK(c *fiber.Ctx, data any) error {
	return c.JSON(fiber.Map{"success": true, "data": data})
}

func OKList(c *fiber.Ctx, data any, total int64) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
		"meta":    fiber.Map{"total": total},
	})
}

func Fail(c *fiber.Ctx, status int, code, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"success": false,
		"error":   fiber.Map{"code": code, "message": message},
	})
}
