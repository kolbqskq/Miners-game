package game

import "miners_game/internal/game/domain"

type IGameRepository interface {
	Load(userID, gameID string) (*domain.GameState, error)
	Save(gameState *domain.GameState) error
}

type ILoopService interface {
	Tick(now int64)
	Register(id string, game *domain.GameState)
	Unregister(id string)
}

type ISessionService interface {
	MarkActive(id string)
	IsActive(id string) bool
	GetExpired() []string
}
