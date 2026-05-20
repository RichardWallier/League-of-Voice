package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"lov/domain"
	"lov/dto"
	"lov/entity"
	"lov/utils"

	"golang.org/x/crypto/argon2"
)

type UserService struct {
	entity *entity.UserEntity
}

func NewUserService(e *entity.UserEntity) *UserService {
	return &UserService{
		entity: e,
	}
}

func (s *UserService) GetAllUsers(ctx context.Context) []domain.User {
		return s.entity.GetAllUsers(ctx)
}

func (s *UserService) CreateUserPassword(password string, secret []byte) (string, string, error) {
	newSecret := secret
	if secret == nil {
		var err error
		newSecret, err = utils.NewSecret(32)
		if err != nil {
			return "", "", fmt.Errorf("failed to generate secret: %w", err)
		}
	}
	hashPassword := argon2.Key([]byte(password), newSecret, 3, 32*1024, 4, 32)
	encodeHashPassword := base64.StdEncoding.EncodeToString(hashPassword)
	encodeSalt := base64.StdEncoding.EncodeToString(newSecret)
	return encodeHashPassword, encodeSalt, nil
}

func (s *UserService) CreateUser(ctx context.Context, user dto.CreateUserRequest) (domain.User, error) {
	encodePassword, encodeSalt, err := s.CreateUserPassword(user.Password, nil)
	newUser := domain.User{
		Email:    user.Email,
		Username: user.Username,
		Password: encodePassword,
		Salt: encodeSalt,
	}
	newUser, err = s.entity.CreateUser(ctx, newUser)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to create user: %w", err)
	}
	return newUser, nil
}
