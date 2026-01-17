package loop

import (
	"miners_game/internal/game/domain"
	"sync"
)

type Service struct {
	games map[string]*domain.GameState
	mu    sync.RWMutex
}

func NewService() *Service {
	return &Service{
		games: make(map[string]*domain.GameState),
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
}

func (s *Service) Unregister(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.games, id)
}
