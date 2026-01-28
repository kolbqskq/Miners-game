package loop

import (
	"miners_game/internal/game/domain"
	"sync"

	"github.com/rs/zerolog"
)

type Service struct {
	games  map[string]*domain.GameState
	mu     sync.RWMutex
	logger zerolog.Logger
}

type ServiceDeps struct {
	Logger zerolog.Logger
}

func NewService(deps ServiceDeps) *Service {
	return &Service{
		games:  make(map[string]*domain.GameState),
		logger: deps.Logger,
	}
}

func (s *Service) Tick(now int64) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, game := range s.games {
		game.Tick(now)
	}
}

func (s *Service) Register(id string, game *domain.GameState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.games[id] = game

	s.logger.Debug().Str("game_id/save_id", id).Msg("game registered in loop")
}

func (s *Service) Unregister(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.games, id)

	s.logger.Debug().Str("game_id/save_id", id).Msg("game unregistered from loop")
}
