package user

import "github.com/google/uuid"

type User struct {
	ID       string
	Email    string
	Password string
	UserName string
}

func NewUser(email, password, userName string) *User {
	return &User{
		ID:       uuid.NewString(),
		Email:    email,
		Password: password,
		UserName: userName,
	}
}
