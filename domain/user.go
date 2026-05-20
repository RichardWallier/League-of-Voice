package domain

import (
	"fmt"
	"lov/db"
	"time"
)

type User struct {
	ID        int32  `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"-"`
	Salt      string `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUser(user db.User) User {
	return User{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		Password:  user.Password,
		CreatedAt: user.CreatedAt.Time,
	}
}

func NewUserList(user []db.User) []User {
	var users []User
	for _, u := range user {
		users = append(users, User{
			ID:        u.ID,
			Email:     u.Email,
			Username:  u.Username,
			Password:  u.Password,
			CreatedAt: u.CreatedAt.Time,
		})
	}
	fmt.Println(users)
	return users
}
