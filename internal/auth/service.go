package auth

import (
	"fmt"
	"miners_game/config"
	"miners_game/internal/user"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepo    user.IUserRepository
	gmailConfig *config.GmailConfig
	logger      zerolog.Logger
}

type ServiceDeps struct {
	UserRepository user.IUserRepository
	GmailConfig    *config.GmailConfig
	Logger         zerolog.Logger
}

func NewService(deps ServiceDeps) *Service {
	return &Service{
		userRepo:    deps.UserRepository,
		gmailConfig: deps.GmailConfig,
		logger:      deps.Logger,
	}
}

func (s *Service) Register(email, username, hashedPassword string) (string, error) {
	existedUser, _ := s.userRepo.FindByEmail(email)
	if existedUser != nil {
		return "", fmt.Errorf("Пользователь уже существует") //errs
	}
	user := user.NewUser(email, hashedPassword, username)

	if err := s.userRepo.SaveUser(user); err != nil {
		return "", err
	}

	return user.ID, nil
}

func (s *Service) Login(email, password string) (string, string, error) {
	user, _ := s.userRepo.FindByEmail(email)
	if user == nil {
		return "", "", fmt.Errorf("Непральный email или пароль") //errs
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", fmt.Errorf("Непральный email или пароль") //errs
	}
	return user.ID, user.UserName, nil
}

func (s *Service) SendEmail(to, code string) error {
	client := resty.New()

	payload := map[string]interface{}{
		"from": map[string]string{
			"email": "test@test-nrw7gymqexog2k8e.mlsender.net",
			"name":  "Verification code",
		},
		"to": []map[string]string{
			{"email": to},
		},
		"subject": "Verification code",
		"text":    code,
	}

	_, err := client.R().
		SetHeader("Authorization", "Bearer "+s.gmailConfig.AppPassword).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(payload).
		Post("https://api.mailersend.com/v1/email")

	if err != nil {
		return err
	}

	return nil
}
