package auth

import (
	"miners_game/config"
	"miners_game/internal/user"
	"miners_game/pkg/code"
	"miners_game/pkg/errs"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gookit/validate"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepo    user.IUserRepository
	gmailConfig *config.GmailConfig
	metrics     *Metrics
	logger      zerolog.Logger
}

type ServiceDeps struct {
	UserRepository user.IUserRepository
	GmailConfig    *config.GmailConfig
	Metrics        *Metrics
	Logger         zerolog.Logger
}

func NewService(deps ServiceDeps) *Service {
	return &Service{
		userRepo:    deps.UserRepository,
		gmailConfig: deps.GmailConfig,
		metrics:     deps.Metrics,
		logger:      deps.Logger,
	}
}

func (s *Service) login(form LoginForm) (userID, userName string, err error) {
	defer func() {
		if err != nil {
			s.metrics.LoginFailedTotal.Inc()
		}
	}()
	v := validate.Struct(&form)
	if !v.Validate() {
		s.logger.Warn().Err(v.Errors).Msg("failed to validate login form")
		return "", "", v.Errors.OneError()
	}
	user, err := s.userRepo.FindByEmail(form.Email)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to find user by email")
		return "", "", err
	}
	if user == nil {
		s.logger.Warn().Msg("failed to login incorrect email")
		return "", "", errs.ErrIncorrectLogin
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
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

func (s *Service) userExist(email string) (bool, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to find user by email")
		return false, err
	}
	return user != nil, nil
}

func (s *Service) startRegistration(form RegisterForm) (formSess RegisterSession, err error) {
	defer func() {
		s.metrics.RegisterAttemptsTotal.Inc()
		if err != nil {
			s.metrics.RegisterFailedTotal.Inc()
		}
	}()
	v := validate.Struct(&form)
	if !v.Validate() {
		s.logger.Warn().Err(v.Errors).Msg("failed to validate register form")
		return RegisterSession{}, v.Errors.OneError()
	}
	userExist, err := s.userExist(form.Email)
	if err != nil {
		return RegisterSession{}, errs.ErrServer
	}
	if userExist {
		s.logger.Warn().Msg("failed user already exist")
		return RegisterSession{}, errs.ErrEmailAlreadyExist
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to hash password")
		return RegisterSession{}, errs.ErrServer
	}
	code := code.Generate()
	if err := s.sendEmail(form.Email, code); err != nil {
		s.logger.Error().Err(err).Msg("failed to send email")
		return RegisterSession{}, errs.ErrServer
	}
	sess := RegisterSession{
		Email:          form.Email,
		Code:           code,
		Username:       form.UserName,
		HashedPassword: string(hashedPassword),
		ExpiresAt:      time.Now().Add(10 * time.Minute).Unix(),
	}
	return sess, nil
}

func (s *Service) completeRegistration(sess RegisterSession, enteredCode string) (userID string, err error) {
	defer func() {
		if err != nil {
			s.metrics.RegisterFailedTotal.Inc()
		} else {
			s.metrics.RegisterSuccessTotal.Inc()
		}
	}()
	if time.Now().Unix() > sess.ExpiresAt {
		s.logger.Warn().Msg("failed expire session")
		return "", errs.ErrExpireSession
	}
	if enteredCode == "" {
		return "", errs.ErrEmptyRegisterCode
	}
	if enteredCode != sess.Code {
		return "", errs.ErrRegisterCode
	}

	user := user.NewUser(sess.Email, sess.HashedPassword, sess.Username)

	if err := s.userRepo.SaveUser(user); err != nil {
		s.logger.Error().Err(err).Str("email", sess.Email).Msg("failed to save user")
		return "", errs.ErrServer
	}

	return user.ID, nil
}
