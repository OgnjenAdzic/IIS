package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey []byte
	issuer    string
}

func NewJWTService(secret string, issuer string) *JWTService {
	return &JWTService{
		secretKey: []byte(secret),
		issuer:    issuer,
	}
}

func (j *JWTService) GenerateToken(userId, role, username string) (string, error) {
	claims := jwt.MapClaims{
		"sub":      userId,   // Standard Subject claim
		"role":     role,     // Custom claim
		"username": username, // Custom claim
		"iss":      j.issuer,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *JWTService) ValidateToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})
}
