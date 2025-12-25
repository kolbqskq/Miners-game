package game

import (
	"miners_game/internal/miners"
	"sync"
)

type Service struct {
	memoryGameState map[string]*GameState
	repo            IGameRepository
	mu              sync.RWMutex
}

type ServiceDeps struct {
	GameRepository IGameRepository
}

func NewService(deps ServiceDeps) *Service {
	return &Service{
		repo: deps.GameRepository,
	}
}

func (s *Service) BuyMiner(class, userID, saveID string) (*GameState, error) {
	miner := miners.NewMiner(class)
	price := miners.GetMinerConfig(class).Price

	gameState, err := s.repo.GetGameState(userID, saveID)
	if err != nil {
		return gameState, err
	}

	gameState.RecalculateBalance()
	k := userID + "/" + saveID
	s.memoryGameState[k] = gameState
	if err := gameState.ValidateBalance(price); err != nil {
		return gameState, err
	}
	gameState.Miners[miner.ID] = *miner

	if err := s.repo.SaveGameState(gameState); err != nil {
		return gameState, err
	}
	return gameState, nil

}

func (s *Service) GetGameStateFromMemory(userID, saveID string) *GameState {
	k := userID + "/" + saveID
	s.mu.RLock()
	gs, _ := s.memoryGameState[k]
	s.mu.RUnlock()
	return gs
}
