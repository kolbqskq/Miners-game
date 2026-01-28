package auth

import (
	"miners_game/config"
	"miners_game/internal/user"
	"miners_game/pkg/errs"

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

func (s *Service) register(email, username, hashedPassword string) (string, error) {
	existedUser, _ := s.userRepo.FindByEmail(email)
	if existedUser != nil {
		s.logger.Warn().Str("email", email).Msg("failed register account already exist")
		return "", errs.ErrEmailAlreadyExist
	}
	user := user.NewUser(email, hashedPassword, username)

	if err := s.userRepo.SaveUser(user); err != nil {
		s.logger.Error().Err(err).Str("email", email).Msg("failed to save user")
		return "", errs.ErrServer
	}

	return user.ID, nil
}

func (s *Service) login(email, password string) (string, string, error) {
	user, _ := s.userRepo.FindByEmail(email)
	if user == nil {
		s.logger.Warn().Msg("failed to login incorrect email")
		return "", "", errs.ErrIncorrectLogin
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		s.logger.Warn().Msg("failed to login incorrect password")
		return "", "", errs.ErrIncorrectLogin
	}
	return user.ID, user.UserName, nil
}

func (s *Service) sendEmail(to, code string) error {
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
		s.logger.Error().Err(err).Msg("failed to send Email")
		return errs.ErrServer
	}
	s.logger.Debug().Str("email", to).Msg("email sended")
	return nil
}

func (s *Service) emailExist(email string) bool {
	user, _ := s.userRepo.FindByEmail(email)
	if user != nil {
		return true
	}
	return false
}
