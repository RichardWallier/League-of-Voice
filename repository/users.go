package repository

import (
	"context"
	"errors"
	"fmt"
	"lov/db"
	"lov/domain"

	"github.com/jackc/pgx/v5/pgconn"
)

var ErrUserAlreadyExists = fmt.Errorf("user already exists")

type UserEntity struct {
	query *db.Queries
}

func NewUserEntity(query *db.Queries) *UserEntity {
	return &UserEntity{
		query: query,
	}
}

func (u *UserEntity) GetAllUsers(ctx context.Context) []domain.User {
	users, err := u.query.ListUsers(ctx)
	if err != nil {
		panic("failed to list users: " + err.Error())
	}

	return domain.NewUserList(users)
}

func (u *UserEntity) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	fmt.Printf("DEBUG: Creating user with email: %s, username: %s\n", user.Email, user.Username)
	newUser, err := u.query.CreateUser(ctx, db.CreateUserParams{
		Email: user.Email,
		Username: user.Username,
		Password: user.Password,
		Salt: user.Salt,
	})
	fmt.Printf("Created user: %+v\n", newUser)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.User{}, ErrUserAlreadyExists
		}
		return domain.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return domain.NewUser(newUser), nil
}

func (u *UserEntity) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := u.query.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get user by email: %w", err)
	}

	return domain.NewUser(user), nil
}

func (u *UserEntity) AttachRoleToUser(ctx context.Context, userID int32, roleID int32) error {
	err := u.query.AssignRoleToUser(ctx, db.AssignRoleToUserParams{
		UserID: userID,
		RoleID: roleID,
	})
	if err != nil {
		return fmt.Errorf("failed to attach permission to user: %w", err)
	}

	return nil
}
