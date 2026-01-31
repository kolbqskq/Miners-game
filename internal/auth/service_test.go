package auth_test

import (
	"errors"
	"miners_game/internal/auth"
	"miners_game/internal/user"
	"miners_game/pkg/errs"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	MockFindByEmail func(email string) (*user.User, error)
}

func (m *MockUserRepository) SaveUser(user *user.User) error {
	return nil
}

func (m *MockUserRepository) FindByEmail(email string) (*user.User, error) {
	return m.MockFindByEmail(email)
}

type MockEmailService struct {
}

func (m *MockEmailService) Send(to, code string) error {
	return nil
}

func TestStartRegisterSuccess(t *testing.T) {
	repo := &MockUserRepository{
		MockFindByEmail: func(email string) (*user.User, error) {
			return nil, nil
		},
	}
	emailService := &MockEmailService{}
	form := auth.RegisterForm{
		Email:           "testReg@gmail.com",
		UserName:        "testUsername",
		Password:        "testPass",
		PasswordConfirm: "testPass",
	}
	authService := auth.NewService(auth.ServiceDeps{
		UserRepository: repo,
		EmailService:   emailService,
	})
	if _, err := authService.StartRegistration(form); err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestCompleteRegistrationSuccess(t *testing.T) {
	repo := &MockUserRepository{}
	code := "admin"
	sess := auth.RegisterSession{
		Email:          "testReg@gmail.com",
		Code:           "admin",
		Username:       "testUsername",
		HashedPassword: "testPass",
		ExpiresAt:      time.Now().Unix() + 100,
	}
	authService := auth.NewService(auth.ServiceDeps{
		UserRepository: repo,
	})
	if _, err := authService.CompleteRegistration(sess, code); err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestLoginSuccess(t *testing.T) {
	form := auth.LoginForm{
		Email:    "testLog@gmail.com",
		Password: "testPass",
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	repo := &MockUserRepository{
		MockFindByEmail: func(email string) (*user.User, error) {
			return &user.User{
				Email:    form.Email,
				Password: string(hashed),
			}, nil
		},
	}

	authService := auth.NewService(auth.ServiceDeps{
		UserRepository: repo,
	})
	if _, _, err := authService.Login(form); err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	wrongPass := "testWrongPass"
	form := auth.LoginForm{
		Email:    "testLog@gmail.com",
		Password: "testPass",
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(wrongPass), bcrypt.DefaultCost)
	repo := &MockUserRepository{
		MockFindByEmail: func(email string) (*user.User, error) {
			return &user.User{
				Email:    "testLog@gmail.com",
				Password: string(hashed),
			}, nil
		},
	}
	authService := auth.NewService(auth.ServiceDeps{
		UserRepository: repo,
	})
	_, _, err := authService.Login(form)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, errs.ErrIncorrectLogin) {
		t.Fatalf("expected ErrIncorrectLogin, got %v:", err)
	}
}
