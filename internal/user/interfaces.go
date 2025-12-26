package user

type IUserRepository interface {
	SaveUser(user *User) error
	FindByEmail(email string) (*User, error)
}
