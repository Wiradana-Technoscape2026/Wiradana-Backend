package middleware_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/wiradana/backend/internal/delivery/http/middleware"
)

func newApp(role string, handler fiber.Handler, guards ...fiber.Handler) *fiber.App {
	app := fiber.New()
	// Simulate auth middleware by setting Locals manually
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("role", role)
		c.Locals("cooperative_id", "00000000-0000-0000-0000-000000000001")
		c.Locals("member_id", "00000000-0000-0000-0000-000000000002")
		return c.Next()
	})
	for _, g := range guards {
		app.Use(g)
	}
	app.Get("/test", handler)
	return app
}

func doGet(app *fiber.App) *http.Response {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, _ := app.Test(req)
	return resp
}

func TestRequireRole_AnggotaAccessesPengurusEndpoint_Returns403(t *testing.T) {
	app := newApp("anggota", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	}, middleware.RequireRole("pengurus"))

	resp := doGet(app)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("want 403, got %d", resp.StatusCode)
	}
}

func TestRequireRole_PengurusAccessesPengurusEndpoint_Returns200(t *testing.T) {
	app := newApp("pengurus", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	}, middleware.RequireRole("pengurus"))

	resp := doGet(app)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("want 200, got %d", resp.StatusCode)
	}
}

func TestRequireRole_AnggotaAccessesPortalEndpoint_Returns200(t *testing.T) {
	app := newApp("anggota", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	}, middleware.RequireRole("anggota"))

	resp := doGet(app)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("want 200, got %d", resp.StatusCode)
	}
}

func TestRequireRole_PengurusAccessesPortalEndpoint_Returns403(t *testing.T) {
	app := newApp("pengurus", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	}, middleware.RequireRole("anggota"))

	resp := doGet(app)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("want 403, got %d", resp.StatusCode)
	}
}

func TestRequireRole_NoRole_Returns403(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.RequireRole("pengurus"))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, _ := app.Test(req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("want 403, got %d", resp.StatusCode)
	}
}

func TestAuth_MissingToken_Returns401(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.Auth("secret"))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, _ := app.Test(req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("want 401, got %d", resp.StatusCode)
	}
	_ = io.Discard
}

func TestAuth_InvalidToken_Returns401(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.Auth("secret"))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	resp, _ := app.Test(req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("want 401, got %d", resp.StatusCode)
	}
}
