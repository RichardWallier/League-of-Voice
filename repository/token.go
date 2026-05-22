package repository

import (
	"context"
	"lov/db"
)

type TokenEntity struct {
	query *db.Queries
}

func NewTokenEntity(query *db.Queries) *TokenEntity {
	return &TokenEntity{query: query}
}

func (t *TokenEntity) MissingAnyPermissionsByEmail(ctx context.Context, userId int32, permissions []db.Permission) (bool, error) {
	hasPermission, err := t.query.ListPermissionsByUser(ctx, userId)
	if err != nil {
		return false, err
	}

	hasPermissionMap := make(map[string]bool)
	for _, p := range hasPermission {
		hasPermissionMap[p.Name] = true
	}

	for _, p := range permissions {
		if !hasPermissionMap[p.Name] {
			return true, nil
		}
	}

	return false, nil
}
