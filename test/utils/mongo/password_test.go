package utils_test

import (
	"testing"

	utils "go-fiber/utils/mongo"
)

func TestHashAndCheckPassword(t *testing.T) {
	hash, err := utils.HashPassword("secret123")
	if err != nil {
		t.Fatalf("HashPassword error %v", err)
	}
	if !utils.CheckPassword("secret123", hash) {
		t.Fatalf("CheckPassword should be true")
	}
	if utils.CheckPassword("wrong", hash) {
		t.Fatalf("CheckPassword should be false")
	}
}


