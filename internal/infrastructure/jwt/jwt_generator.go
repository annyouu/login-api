package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"login/internal/usecase"
)

type JWTGeneratorImpl struct {
	secretKey string
	expireMinutes int
}

func NewJWTGenerator(secretKey string, expireMinutes int) usecase.JWTGeneratorInterface {
	return &JWTGeneratorImpl{
		secretKey: secretKey,
		expireMinutes: expireMinutes,
	}
}

func (j *JWTGeneratorImpl) Generate(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp": time.Now().Add(time.Duration(j.expireMinutes) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

