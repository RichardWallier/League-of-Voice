package service

import (
	"context"
	"fmt"
	"lov/constants"
	"lov/domain"
	"lov/repository"
	"lov/utils"
)

var ErrUserNotAuthenticated = fmt.Errorf("User not authenticated")

type UserService struct {
	entity *repository.UserEntity
}

func NewUserService(e *repository.UserEntity) *UserService {
	return &UserService{
		entity: e,
	}
}

func (s *UserService) GetAllUsers(ctx context.Context) []domain.User {
	return s.entity.GetAllUsers(ctx)
}


func (s *UserService) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	encodePassword, encodeSalt, err := utils.CreateHashPassword(user.Password, nil)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to create user password: %w", err)
	}

	newUser := domain.User{
		Email:    user.Email,
		Username: user.Username,
		Password: encodePassword,
		Salt: encodeSalt,
	}
	fmt.Printf("1. Creating user with email: %s, username: %s\n", newUser.Email, newUser.Username)
	newUser, err = s.entity.CreateUser(ctx, newUser)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to create user: %w", err)
	}
	if err := s.entity.AttachRoleToUser(ctx, newUser.ID, constants.RoleUser); err != nil {
		return domain.User{}, fmt.Errorf("failed to attach role to user: %w", err)
	}

	return newUser, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	return s.entity.GetUserByEmail(ctx, email)
}
