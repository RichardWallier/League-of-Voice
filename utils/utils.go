package utils

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/argon2"
)

func NewSecret(length uint32) ([]byte, error){
	secret := make([]byte, length)
	_, err := rand.Read(secret)
	if err != nil {
		return nil, err
	}
	return secret, nil
}

func CreateHashPassword(password string, secret []byte) (string, string, error) {
	newSecret := secret
	if secret == nil {
		var err error
		newSecret, err = NewSecret(32)
		if err != nil {
			return "", "", fmt.Errorf("failed to generate secret: %w", err)
		}
	}
	hashPassword := argon2.Key([]byte(password), newSecret, 3, 32*1024, 4, 32)
	encodeHashPassword := base64.StdEncoding.EncodeToString(hashPassword)
	encodeSalt := base64.StdEncoding.EncodeToString(newSecret)

	return encodeHashPassword, encodeSalt, nil
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
	trimmedToken := strings.Trim(strings.TrimPrefix(tokenString, "Bearer"), " ")

	token, err := jwt.Parse(trimmedToken, func(token *jwt.Token) (any, error) {
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
