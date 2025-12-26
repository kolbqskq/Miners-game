package auth

type LoginForm struct {
	Email    string `json:"email" validate:"required|email"`
	Password string `json:"password" validate:"required"`
}

type RegisterForm struct {
	Email           string `json:"email" validate:"required|email"`
	UserName        string `json:"userName" validate:"required|minLen:8"`
	Password        string `json:"password" validate:"required|minLen:8"`
	PasswordConfirm string `json:"passwordConfirm" validate:"required|eqField:password"`
}
