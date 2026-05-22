package command

type (
	CreateUserInput struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	CreateUserOutput struct {
		Token []byte `json:"token"`
	}
)
