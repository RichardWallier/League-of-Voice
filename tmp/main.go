package tmp

import "context"


type DomainUser struct {}

type (
	UserRepository interface{
		CreateUser(ctx context.Context, newUser DomainUser)
		// ...
	}
)

type (
	CreateUserCommand struct {}

	CreateUserUseCase  interface {
		Execute(ctx context.Context, cmd CreateUserCommand)
	}
)

type (


)
