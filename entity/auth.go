package entity

import (
	"context"
	"lov/db"
)

type AuthEntity struct {
	query *db.Queries
}

func NewAuthEntity(query *db.Queries) *AuthEntity {
	return &AuthEntity{query: query}
}

func (a *AuthEntity) GetUserPasswordAndSalt(ctx context.Context, email string) (string, string, error) {
	user, err := a.query.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}
	return user.Password, user.Salt, nil
}
