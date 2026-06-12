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

func Auth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "missing or invalid authorization header",
			})
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
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "invalid or expired token",
			})
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
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error":   "insufficient permissions",
		})
	}
}

func RequireModule(db *gorm.DB, key string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		coopID, _ := c.Locals("cooperative_id").(string)
		var mod entity.CoopModule
		err := db.Where("cooperative_id = ? AND module_key = ? AND enabled = true", coopID, key).First(&mod).Error
		if err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error":   "module not enabled",
			})
		}
		return c.Next()
	}
}
