package game

type IGameRepository interface {
	GetGameState(userID, saveID string) (*GameState, error)
	SaveGameState(gameState *GameState) error
}
