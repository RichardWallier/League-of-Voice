package entity

import "lov/db"

type Entities struct {
	UserEntity *UserEntity
	AuthEntity *AuthEntity
	TokenEntity *TokenEntity
}

func NewEntities(db *db.PostgresDB) *Entities {
	return &Entities{NewUserEntity(db.Queries), NewAuthEntity(db.Queries), NewTokenEntity(db.Queries)}
}
