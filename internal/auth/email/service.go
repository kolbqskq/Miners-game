package email

import (
	"github.com/rs/zerolog"
)

type Service struct {
	logger zerolog.Logger
}

type ServiceDeps struct {
	Logger zerolog.Logger
}

func NewService(deps ServiceDeps) *Service {
	return &Service{
		logger: deps.Logger,
	}
}

func (s *Service) Send(to, code string) error {
	s.logger.Info().Str("to", to).Str("code", code).Msg("send verification email")
	return nil
}
