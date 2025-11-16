package service_test

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	service "go-fiber/app/service/mongo"
	model "go-fiber/app/model/mongo"

	"github.com/gofiber/fiber/v2"
)

type mockAlumniRepo struct {
	byID *model.Alumni
	all  []model.Alumni
	err  error
}

func (m *mockAlumniRepo) CreateAlumni(ctx context.Context, alumni *model.Alumni) (*model.Alumni, error) {
	alumni.CreatedAt = time.Now()
	alumni.UpdatedAt = time.Now()
	return alumni, nil
}
func (m *mockAlumniRepo) FindAlumniByID(ctx context.Context, id string) (*model.Alumni, error) { return m.byID, m.err }
func (m *mockAlumniRepo) FindAlumniByEmail(ctx context.Context, email string) (*model.Alumni, error) {
	return nil, nil
}
func (m *mockAlumniRepo) FindAlumniByNIM(ctx context.Context, nim string) (*model.Alumni, error) { return nil, nil }
func (m *mockAlumniRepo) FindAllAlumni(ctx context.Context) ([]model.Alumni, error)               { return m.all, m.err }
func (m *mockAlumniRepo) UpdateAlumni(ctx context.Context, id string, alumni *model.Alumni) (*model.Alumni, error) {
	alumni.UpdatedAt = time.Now()
	return alumni, m.err
}
func (m *mockAlumniRepo) DeleteAlumni(ctx context.Context, id string) error { return m.err }

func TestGetAllAlumniService(t *testing.T) {
	repo := &mockAlumniRepo{
		all: []model.Alumni{
			{NIM: "1", Nama: "A", Jurusan: "TI", Angkatan: 2021, TahunLulus: 2025, Email: "a@a.com", Role: "user", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{NIM: "2", Nama: "B", Jurusan: "SI", Angkatan: 2020, TahunLulus: 2024, Email: "b@b.com", Role: "admin", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
	}
	svc := service.NewAlumniService(repo)
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error { return svc.GetAllAlumniService(c) })
	req := httptest.NewRequest("GET", "/", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != 200 {
		t.Fatalf("status got %d want %d", resp.StatusCode, 200)
	}
}

func TestUpdateAlumniServiceBadRole(t *testing.T) {
	repo := &mockAlumniRepo{byID: &model.Alumni{}}
	svc := service.NewAlumniService(repo)
	app := fiber.New()
	app.Put("/:id", func(c *fiber.Ctx) error { return svc.UpdateAlumniService(c) })
	body := `{"nama":"A","jurusan":"TI","angkatan":2021,"tahun_lulus":2025,"email":"a@a.com","role":"guest"}`
	req := httptest.NewRequest("PUT", "/507f1f77bcf86cd799439011", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != 400 {
		t.Fatalf("status got %d want %d", resp.StatusCode, 400)
	}
}


