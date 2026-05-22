package service

import (
	"context"
	"lov/db"
	"lov/repository"
	"lov/utils"
)

type TokenService struct {
	tokenEntity *repository.TokenEntity
	userEntity *repository.UserEntity
}

func NewTokenService(tokenEntity *repository.TokenEntity, userEntity *repository.UserEntity) *TokenService {
	return &TokenService{
		tokenEntity: tokenEntity,
		userEntity: userEntity,
	}
}

func (t *TokenService) ValidatePermissions(ctx context.Context, tokenString string, permissions []db.Permission) (bool, error) {
	email, err := utils.ValidateJWTToken(ctx, tokenString)
	if err != nil {
		return false, err
	}
	user, err := t.userEntity.GetUserByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	missingPermissions, err := t.tokenEntity.MissingAnyPermissionsByEmail(ctx, user.ID, permissions)
	if err != nil {
		return false, err
	}

	return missingPermissions, nil
}
