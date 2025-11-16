package service_test

import (
	"context"
	"net/http/httptest"
	"strings"
	testing "testing"

	service "go-fiber/app/service/mongo"
	model "go-fiber/app/model/mongo"

	"github.com/gofiber/fiber/v2"
)

type mockPekerjaanRepo struct {
	byID *model.PekerjaanAlumni
	all  []model.PekerjaanAlumni
	err  error
}

func (m *mockPekerjaanRepo) CreatePekerjaanAlumni(ctx context.Context, p *model.PekerjaanAlumni) (*model.PekerjaanAlumni, error) {
	return p, m.err
}
func (m *mockPekerjaanRepo) FindPekerjaanAlumniByID(ctx context.Context, id string) (*model.PekerjaanAlumni, error) {
	return m.byID, m.err
}
func (m *mockPekerjaanRepo) FindAllPekerjaanAlumni(ctx context.Context) ([]model.PekerjaanAlumni, error) {
	return m.all, m.err
}
func (m *mockPekerjaanRepo) FindPekerjaanAlumniByAlumniID(ctx context.Context, id string) ([]model.PekerjaanAlumni, error) {
	return m.all, m.err
}
func (m *mockPekerjaanRepo) UpdatePekerjaanAlumni(ctx context.Context, id string, p *model.PekerjaanAlumni) (*model.PekerjaanAlumni, error) {
	return p, m.err
}
func (m *mockPekerjaanRepo) DeletePekerjaanAlumni(ctx context.Context, id string) error { return m.err }

func TestCreatePekerjaanAlumniServiceBadTanggal(t *testing.T) {
	repo := &mockPekerjaanRepo{}
	svc := service.NewPekerjaanAlumniService(repo)
	app := fiber.New()
	app.Post("/", func(c *fiber.Ctx) error { return svc.CreatePekerjaanAlumniService(c) })
	body := `{"alumni_info":{"alumni_id":"507f1f77bcf86cd799439011","nim":"1","nama":"A","email":"a@a.com"},"nama_perusahaan":"X","posisi_jabatan":"Y","bidang_industri":"Z","lokasi_kerja":"K","status_pekerjaan":"aktif","tanggal_mulai_kerja":"2025-13-01"}`
	req := httptest.NewRequest("POST", "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != 400 {
		t.Fatalf("status got %d want %d", resp.StatusCode, 400)
	}
}

func TestGetPekerjaanAlumniByIDNotFound(t *testing.T) {
	repo := &mockPekerjaanRepo{byID: nil}
	svc := service.NewPekerjaanAlumniService(repo)
	app := fiber.New()
	app.Get("/:id", func(c *fiber.Ctx) error { return svc.GetPekerjaanAlumniByIDService(c) })
	req := httptest.NewRequest("GET", "/507f1f77bcf86cd799439011", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != 404 {
		t.Fatalf("status got %d want %d", resp.StatusCode, 404)
	}
}


