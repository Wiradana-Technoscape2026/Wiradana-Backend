package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/wiradana/backend/internal/entity"
	"gorm.io/gorm"
)

type Claims struct {
	UserID        string `json:"user_id"`
	CooperativeID string `json:"cooperative_id"`
	Role          string `json:"role"`
	MemberID      string `json:"member_id"`
	jwt.RegisteredClaims
}

func failMiddleware(c *fiber.Ctx, status int, code, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"success": false,
		"error":   fiber.Map{"code": code, "message": message},
	})
}

func Auth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return failMiddleware(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "token tidak ditemukan")
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return failMiddleware(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "token tidak valid atau sudah kedaluwarsa")
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("cooperative_id", claims.CooperativeID)
		c.Locals("role", claims.Role)
		c.Locals("member_id", claims.MemberID)

		return c.Next()
	}
}

func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("role").(string)
		for _, r := range roles {
			if role == r {
				return c.Next()
			}
		}
		return failMiddleware(c, fiber.StatusForbidden, "FORBIDDEN", "akses ditolak")
	}
}

func RequireModule(db *gorm.DB, key string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		coopID, _ := c.Locals("cooperative_id").(string)
		var mod entity.CoopModule
		err := db.Where("cooperative_id = ? AND module_key = ? AND enabled = true", coopID, key).First(&mod).Error
		if err != nil {
			return failMiddleware(c, fiber.StatusForbidden, "FORBIDDEN", "modul tidak aktif")
		}
		return c.Next()
	}
}
