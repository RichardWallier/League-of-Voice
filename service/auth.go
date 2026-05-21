package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"lov/dto"
	"lov/entity"
)

var ErrWrongPasswordOrUserNotFound = fmt.Errorf("Wrong password or user not found")

type AuthService struct {
	e *entity.AuthEntity
	UserService	*UserService
}

func NewAuthService(e *entity.AuthEntity) *AuthService {
	return &AuthService{e: e}
}

func (s *AuthService) Login(ctx context.Context, user dto.LoginRequest) (string, error) {
	dbUser, dbSalt, err := s.e.GetUserPasswordAndSalt(ctx, user.Email)
	if err != nil {
		return "", fmt.Errorf("failed to get user password and salt: %w", err)
	}

	salt, err := base64.StdEncoding.DecodeString(dbSalt)
	if err != nil {
		return "", fmt.Errorf("failed to decode salt: %w", err)
	}

	userPassword, _, err := s.UserService.CreateUserPassword(user.Password, salt)
	if err != nil {
		return "", fmt.Errorf("failed to create user password: %w", err)

	}
	if userPassword == "" {
		return ""	, ErrWrongPasswordOrUserNotFound
	}
	if userPassword != dbUser {
		return "", ErrWrongPasswordOrUserNotFound
	}

	jwtToken, err := GenerateJWTToken(ctx, user.Email)
	if err != nil {
		fmt.Println("Failed to generate JWT token:", err)
	} else {
		fmt.Println("Generated JWT token:", jwtToken)
	}
	return jwtToken, nil
}
