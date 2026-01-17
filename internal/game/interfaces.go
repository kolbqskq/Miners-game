package game

import "miners_game/internal/game/domain"

type IGameRepository interface {
	Load(userID, gameID string) (*domain.GameState, error)
	Save(gameState *domain.GameState) error
}
