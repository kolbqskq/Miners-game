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
		memoryGameState: make(map[string]*GameState),
		repo:            deps.GameRepository,
	}
}

func (s *Service) StartNewGame(userID, saveID string) (*GameState, error) {
	gameState := NewGameState(userID, saveID)
	if err := s.repo.SaveGameState(gameState); err != nil {
		return gameState, err
	}
	k := userID + "/" + saveID
	s.mu.Lock()
	s.memoryGameState[k] = gameState
	s.mu.Unlock()
	return gameState, nil
}

func (s *Service) BuyMiner(class, userID, saveID string) (*GameState, int64, error) {
	miner := miners.NewMiner(class)
	price := miners.GetMinerConfig(class).Price
	k := userID + "/" + saveID

	gameState, err := s.repo.GetGameState(userID, saveID)
	if err != nil {
		s.mu.Lock()
		gameState = s.memoryGameState[k]
		s.mu.Unlock()
		return gameState, 0, err
	}

	income := gameState.RecalculateBalance()

	if err := gameState.SpendBalance(price); err != nil {
		return gameState, 0, err
	}

	gameState.mu.Lock()
	gameState.Miners[miner.ID] = miner
	gameState.mu.Unlock()

	s.mu.Lock()
	s.memoryGameState[k] = gameState
	s.mu.Unlock()

	if err := s.repo.SaveGameState(gameState); err != nil {
		return gameState, 0, err
	}
	return gameState, income, nil

}

func (s *Service) GetGameStateFromMemory(userID, saveID string) *GameState {
	k := userID + "/" + saveID
	s.mu.RLock()
	gs, _ := s.memoryGameState[k]
	s.mu.RUnlock()
	return gs
}
