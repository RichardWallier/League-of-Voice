package entity

import (
	"lov/db"
)

type TokenEntity struct {
	query *db.Queries
}

func NewTokenEntity(query *db.Queries) *TokenEntity {
	return &TokenEntity{query: query}
}
