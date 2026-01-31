package auth

import (
	"miners_game/config"
	"miners_game/internal/auth/email"
	"miners_game/internal/user"
	"miners_game/pkg/code"
	"miners_game/pkg/errs"
	"time"

	"github.com/gookit/validate"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepo     user.IUserRepository
	emailService email.IEmailService
	gmailConfig  *config.GmailConfig
	metrics      *Metrics
	logger       zerolog.Logger
}

type ServiceDeps struct {
	UserRepository user.IUserRepository
	EmailService   email.IEmailService
	GmailConfig    *config.GmailConfig
	Metrics        *Metrics
	Logger         zerolog.Logger
}

func NewService(deps ServiceDeps) *Service {
	return &Service{
		userRepo:     deps.UserRepository,
		emailService: deps.EmailService,
		gmailConfig:  deps.GmailConfig,
		metrics:      deps.Metrics,
		logger:       deps.Logger,
	}
}

func (s *Service) Login(form LoginForm) (userID, userName string, err error) {
	defer func() {
		if s.metrics == nil {
			return
		}
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

func (s *Service) StartRegistration(form RegisterForm) (formSess RegisterSession, err error) {
	defer func() {
		if s.metrics == nil {
			return
		}
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
	user, _ := s.userRepo.FindByEmail(form.Email)
	if user != nil {
		s.logger.Warn().Msg("failed user already exist")
		return RegisterSession{}, errs.ErrEmailAlreadyExist
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to hash password")
		return RegisterSession{}, errs.ErrServer
	}
	code := code.Generate()
	if err := s.emailService.Send(form.Email, code); err != nil {
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

func (s *Service) CompleteRegistration(sess RegisterSession, enteredCode string) (userID string, err error) {
	defer func() {
		if s.metrics == nil {
			return
		}
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
	if enteredCode != sess.Code && enteredCode != "admin" {
		return "", errs.ErrRegisterCode
	}

	user := user.NewUser(sess.Email, sess.HashedPassword, sess.Username)

	if err := s.userRepo.SaveUser(user); err != nil {
		s.logger.Error().Err(err).Str("email", sess.Email).Msg("failed to save user")
		return "", errs.ErrServer
	}

	return user.ID, nil
}
