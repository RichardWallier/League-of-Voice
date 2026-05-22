package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"lov/dto"
	"lov/repository"
	"lov/utils"
)

var ErrWrongPasswordOrUserNotFound = fmt.Errorf("Wrong password or user not found")

type AuthService struct {
	e *repository.AuthEntity
	userEntity *repository.UserEntity
}

func NewAuthService(e *repository.AuthEntity, userEntity *repository.UserEntity) *AuthService {
	return &AuthService{e: e, userEntity: userEntity}}

func (s *AuthService) Login(ctx context.Context, user dto.LoginRequest) (string, error) {
	dbUser, dbSalt, err := s.e.GetUserPasswordAndSalt(ctx, user.Email)
	if err != nil {
		return "", fmt.Errorf("failed to get user password and salt: %w", err)
	}

	salt, err := base64.StdEncoding.DecodeString(dbSalt)
	if err != nil {
		return "", fmt.Errorf("failed to decode salt: %w", err)
	}

	userPassword, _, err := utils.CreateHashPassword(user.Password, salt)
	if err != nil {
		return "", fmt.Errorf("failed to create user password: %w", err)

	}
	if userPassword == "" || userPassword != dbUser {
		return ""	, ErrWrongPasswordOrUserNotFound
	}

	jwtToken, err := utils.GenerateJWTToken(ctx, user.Email)
	if err != nil {
		fmt.Println("Failed to generate JWT token:", err)
	} else {
		fmt.Println("Generated JWT token:", jwtToken)
	}

	return jwtToken, nil
}
