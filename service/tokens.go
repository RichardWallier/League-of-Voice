package service

import (
	"context"
	"lov/entity"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type TokenService struct {
	e *entity.TokenEntity
}

func NewTokenService(e *entity.TokenEntity) *TokenService {
	return &TokenService{
		e,
	}
}

func GenerateJWTToken(ctx context.Context, email string) (string, error) {
	unsignedJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Minute * 15).Unix(),
	})
	token, err := unsignedJwt.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}

func ValidateJWTToken(ctx context.Context, tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email := claims["email"].(string)
		return email, nil
	} else {
		return "", jwt.ErrSignatureInvalid
	}
}
