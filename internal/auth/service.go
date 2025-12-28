package auth

import (
	"fmt"
	"miners_game/internal/user"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepo user.IUserRepository
	authRepo IAuthRepository
}

type ServiceDeps struct {
	UserRepository user.IUserRepository
	AuthRepository IAuthRepository
}

func NewService(deps ServiceDeps) *Service {
	return &Service{
		userRepo: deps.UserRepository,
		authRepo: deps.AuthRepository,
	}
}

func (s *Service) Register(email, password, userName string) (string, error) {
	existedUser, _ := s.userRepo.FindByEmail(email)
	if existedUser != nil {
		return "", fmt.Errorf("Пользователь уже существует") //errs
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user := user.NewUser(email, string(hashPassword), userName)

	if err := s.userRepo.SaveUser(user); err != nil {
		return "", err
	}

	return user.ID, nil
}

func (s *Service) Login(email, password string) (string, error) {
	user, _ := s.userRepo.FindByEmail(email)
	if user == nil {
		return "", fmt.Errorf("Непральный email или пароль") //errs
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("Непральный email или пароль") //errs
	}
	return user.ID, nil
}
