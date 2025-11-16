package service_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	service "go-fiber/app/service/mongo"
	model "go-fiber/app/model/mongo"
	utilsmongo "go-fiber/utils/mongo"

	"github.com/gofiber/fiber/v2"
)

type mockAlumniRepoAuth struct {
	alumni *model.Alumni
	err    error
}

func (m *mockAlumniRepoAuth) CreateAlumni(ctx context.Context, alumni *model.Alumni) (*model.Alumni, error) { return nil, nil }
func (m *mockAlumniRepoAuth) FindAlumniByID(ctx context.Context, id string) (*model.Alumni, error)          { return nil, nil }
func (m *mockAlumniRepoAuth) FindAlumniByEmail(ctx context.Context, email string) (*model.Alumni, error)    { return m.alumni, m.err }
func (m *mockAlumniRepoAuth) FindAlumniByNIM(ctx context.Context, nim string) (*model.Alumni, error)        { return nil, nil }
func (m *mockAlumniRepoAuth) FindAllAlumni(ctx context.Context) ([]model.Alumni, error)                     { return nil, nil }
func (m *mockAlumniRepoAuth) UpdateAlumni(ctx context.Context, id string, alumni *model.Alumni) (*model.Alumni, error) {
	return nil, nil
}
func (m *mockAlumniRepoAuth) DeleteAlumni(ctx context.Context, id string) error { return nil }

func TestLoginServiceSuccess(t *testing.T) {
	os.Setenv("JWT_SECRET", "secret-for-test-32-chars-minimum-123")
	hash, _ := utilsmongo.HashPassword("secret123")
	repo := &mockAlumniRepoAuth{
		alumni: &model.Alumni{
			NIM:          "2021001",
			Nama:         "A",
			Jurusan:      "TI",
			Angkatan:     2021,
			TahunLulus:   2025,
			Email:        "a@a.com",
			PasswordHash: hash,
			Role:         "admin",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}
	svc := service.NewAuthService(repo)
	app := fiber.New()
	app.Post("/login", func(c *fiber.Ctx) error { return svc.LoginService(c) })
	body, _ := json.Marshal(map[string]string{"email": "a@a.com", "password": "secret123"})
	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != 200 {
		t.Fatalf("status got %d want %d", resp.StatusCode, 200)
	}
}

func TestLoginServiceWrongPassword(t *testing.T) {
	hash, _ := utilsmongo.HashPassword("secret123")
	repo := &mockAlumniRepoAuth{
		alumni: &model.Alumni{Email: "a@a.com", PasswordHash: hash, Role: "user", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	svc := service.NewAuthService(repo)
	app := fiber.New()
	app.Post("/login", func(c *fiber.Ctx) error { return svc.LoginService(c) })
	body, _ := json.Marshal(map[string]string{"email": "a@a.com", "password": "wrong"})
	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != 401 {
		t.Fatalf("status got %d want %d", resp.StatusCode, 401)
	}
}


