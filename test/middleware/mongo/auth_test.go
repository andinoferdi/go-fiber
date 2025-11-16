package middleware_test

import (
	"net/http/httptest"
	"os"
	"testing"

	mw "go-fiber/middleware/mongo"
	utilsmongo "go-fiber/utils/mongo"

	"github.com/gofiber/fiber/v2"
)

func TestAuthRequiredSuccess(t *testing.T) {
	os.Setenv("JWT_SECRET", "secret-for-test-32-chars-minimum-123")
	token, err := utilsmongo.GenerateToken(utilsmongo.AlumniToken{ID: "507f1f77bcf86cd799439011", Email: "a@a.com", Role: "admin"})
	if err != nil {
		t.Fatalf("GenerateToken error %v", err)
	}
	app := fiber.New()
	app.Get("/", mw.AuthRequired(), func(c *fiber.Ctx) error { return c.SendStatus(200) })
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ := app.Test(req)
	if resp.StatusCode != 200 {
		t.Fatalf("status got %d want %d", resp.StatusCode, 200)
	}
}

func TestAuthRequiredMissingHeader(t *testing.T) {
	app := fiber.New()
	app.Get("/", mw.AuthRequired(), func(c *fiber.Ctx) error { return c.SendStatus(200) })
	req := httptest.NewRequest("GET", "/", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != 401 {
		t.Fatalf("status got %d want %d", resp.StatusCode, 401)
	}
}


