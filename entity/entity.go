package entity

import "lov/db"

type Entities struct {
	UserEntity *UserEntity
	AuthEntity *AuthEntity
}

func NewEntities(db *db.PostgresDB) *Entities {
	userEntity := NewUserEntity(db.Queries)
	authEntity := NewAuthEntity(db.Queries)
	return &Entities{UserEntity: userEntity, AuthEntity: authEntity}
}
