package helper

import (
	"go-fiber/app/model"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("SECURITY ERROR: JWT_SECRET environment variable is required but not set! Please add JWT_SECRET to your .env file")
	}
	
	if len(secret) < 32 {
		log.Fatal("SECURITY ERROR: JWT_SECRET must be at least 32 characters long for security")
	}
	
	return []byte(secret)
}

func GenerateToken(alumni model.Alumni) (string, error) {
	claims := model.JWTClaims{
		AlumniID: alumni.ID,
		Email:    alumni.Email,
		RoleID:   alumni.RoleID,
		RoleName: alumni.Role.Nama,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "alumni-api",
			Subject:   alumni.Email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

func ValidateToken(tokenString string) (*model.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return getJWTSecret(), nil
		})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*model.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}