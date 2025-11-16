package utils_test

import (
	"os"
	"testing"

	utils "go-fiber/utils/mongo"
)

func TestGenerateAndValidateToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "secret-for-test-32-chars-minimum-123")
	token, err := utils.GenerateToken(utils.AlumniToken{ID: "507f1f77bcf86cd799439011", Email: "a@a.com", Role: "admin"})
	if err != nil {
		t.Fatalf("GenerateToken error %v", err)
	}
	claims, err := utils.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken error %v", err)
	}
	if claims.AlumniID != "507f1f77bcf86cd799439011" {
		t.Fatalf("AlumniID got %s want %s", claims.AlumniID, "507f1f77bcf86cd799439011")
	}
	if claims.Email != "a@a.com" {
		t.Fatalf("Email got %s want %s", claims.Email, "a@a.com")
	}
	if claims.Role != "admin" {
		t.Fatalf("Role got %s want %s", claims.Role, "admin")
	}
}

func TestExtractTokenFromHeader(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
	}{
		{"BearerWithSpace", "Bearer abc", "abc"},
		{"BearerLower", "bearer abc", "abc"},
		{"BearerNoSpace", "Bearerabc", "abc"},
		{"RawToken", "abc", "abc"},
		{"Empty", "", ""},
	}
	for _, tt := range tests {
		got := utils.ExtractTokenFromHeader(tt.in)
		if got != tt.out {
			t.Fatalf("%s got %s want %s", tt.name, got, tt.out)
		}
	}
}


