package game_test

import (
	"errors"
	"miners_game/internal/game"
	"miners_game/internal/game/domain"
	"miners_game/pkg/errs"
	"testing"
)

// Repository:
type MockGameRepository struct {
	LoadCalled bool
	SaveCalled bool
	MockLoad   func(userID, gameID string) (*domain.GameState, error)
	MockSave   func(gameState *domain.GameState) error
}

func (m *MockGameRepository) Load(userID, gameID string) (*domain.GameState, error) {
	m.LoadCalled = true
	return m.MockLoad(userID, gameID)
}

func (m *MockGameRepository) Save(gameState *domain.GameState) error {
	m.SaveCalled = true
	return m.MockSave(gameState)
}

// Session:
type MockSessionService struct {
	isActive         bool
	MarkActiveCalled bool
	MockGetExpired   func() []string
}

func (m *MockSessionService) MarkActive(id string) {
	m.MarkActiveCalled = true
	m.isActive = true
}
func (m *MockSessionService) IsActive(id string) bool {
	return m.isActive
}
func (m *MockSessionService) GetExpired() []string {
	return m.MockGetExpired()
}

// Loop:
type MockLoopService struct {
	RegisterCalled   bool
	UnregisterCalled bool
}

func (m *MockLoopService) Tick(now int64) {

}
func (m *MockLoopService) Register(id string, game *domain.GameState) {
	m.RegisterCalled = true
}
func (m *MockLoopService) Unregister(id string) {
	m.UnregisterCalled = true

}

func TestEnterGameSuccess(t *testing.T) {
	repo := MockGameRepository{
		MockLoad: func(userID, gameID string) (*domain.GameState, error) {
			return &domain.GameState{}, nil
		},
		MockSave: func(gameState *domain.GameState) error {
			return nil
		},
	}
	sessions := MockSessionService{
		isActive: false,
	}
	loop := MockLoopService{}

	userID := "testUserID"
	gameID := "testGameID"

	gameService := game.NewService(game.ServiceDeps{
		Repo:     &repo,
		Loop:     &loop,
		Sessions: &sessions,
	})
	if _, err := gameService.EnterGame(userID, gameID); err != nil {
		t.Fatalf("expected success, got err %v:", err)
	}
	if !loop.RegisterCalled {
		t.Fatalf("expected game to be registered in loop")
	}

}

func TestEnterGameFromMemorySuccess(t *testing.T) {
	repo := MockGameRepository{
		MockLoad: func(userID, gameID string) (*domain.GameState, error) {
			return &domain.GameState{}, nil
		},
		MockSave: func(gameState *domain.GameState) error {
			return nil
		},
	}
	session := MockSessionService{
		isActive: true,
	}
	loop := MockLoopService{}

	userID := "testUserID"
	gameID := "testGameID"

	gameService := game.NewService(game.ServiceDeps{
		Repo:     &repo,
		Loop:     &loop,
		Sessions: &session,
	})
	gameState := domain.NewGameState(userID, gameID)
	game.PutGameToMemory(gameService, userID, gameID, gameState)
	if _, err := gameService.EnterGame(userID, gameID); err != nil {
		t.Fatalf("expected success, got err %v:", err)
	}
	if repo.LoadCalled {
		t.Fatalf("load sound not be called")
	}
	if session.isActive == false {
		t.Fatalf("expected session should be active")
	}

}

func TestEnterCreateNewSuccess(t *testing.T) {
	repo := MockGameRepository{
		MockLoad: func(userID, gameID string) (*domain.GameState, error) {
			return &domain.GameState{}, errs.ErrGameNotFound
		},
		MockSave: func(gameState *domain.GameState) error {
			return nil
		},
	}
	sessions := MockSessionService{
		isActive: false,
	}
	loop := MockLoopService{}

	userID := "testUserID"
	gameID := "testGameID"

	gameService := game.NewService(game.ServiceDeps{
		Repo:     &repo,
		Loop:     &loop,
		Sessions: &sessions,
	})
	if _, err := gameService.EnterGame(userID, gameID); err != nil {
		t.Fatalf("expected success, got err %v:", err)
	}
	if !loop.RegisterCalled {
		t.Fatalf("expected game to be registered in loop")
	}
	if !repo.SaveCalled {
		t.Fatalf("expected game to be save in repo")
	}

}

func TestEnterRepositoryError(t *testing.T) {
	repo := MockGameRepository{
		MockLoad: func(userID, gameID string) (*domain.GameState, error) {
			return &domain.GameState{}, errors.New("db down")
		},
		MockSave: func(gameState *domain.GameState) error {
			return nil
		},
	}
	sessions := MockSessionService{
		isActive: false,
	}
	loop := MockLoopService{}

	userID := "testUserID"
	gameID := "testGameID"

	gameService := game.NewService(game.ServiceDeps{
		Repo:     &repo,
		Loop:     &loop,
		Sessions: &sessions,
	})
	if _, err := gameService.EnterGame(userID, gameID); err == nil {
		t.Fatalf("expected error")
	}
	if loop.RegisterCalled {
		t.Fatalf("expected game to not be registered in loop")
	}
	if repo.SaveCalled {
		t.Fatalf("expected game to not be save in repo")
	}

}

func TestGetGameStateSessionNotActive(t *testing.T) {
	userID := "testUserID"
	gameID := "testGameID"

	sessions := MockSessionService{
		isActive: false,
	}
	gameService := game.NewService(game.ServiceDeps{
		Sessions: &sessions,
		Repo:     nil,
		Loop:     nil,
	})
	_, err := gameService.GetGameState(userID, gameID)
	if err == nil {
		t.Fatalf("expected error")
	}

	if !errors.Is(err, errs.ErrSessionIsNotActive) {
		t.Fatalf("expected ErrSessionIsNotActive, got %v:", err)
	}
}

func TestGetGameStateGameNotFound(t *testing.T) {
	userID := "testUserID"
	gameID := "testGameID"

	sessions := MockSessionService{
		isActive: true,
	}
	gameService := game.NewService(game.ServiceDeps{
		Sessions: &sessions,
		Repo:     nil,
		Loop:     nil,
	})
	_, err := gameService.GetGameState(userID, gameID)
	if err == nil {
		t.Fatalf("expected error")
	}

	if !errors.Is(err, errs.ErrGameNotFound) {
		t.Fatalf("expected ErrGameNotFound, got %v:", err)
	}
}

func TestBuyMinerSuccess(t *testing.T) {
	userID := "testUserID"
	gameID := "testGameID"

	sessions := MockSessionService{
		isActive: true,
	}
	gameService := game.NewService(game.ServiceDeps{
		Sessions: &sessions,
		Repo:     nil,
		Loop:     nil,
	})
	gameState := domain.NewGameState(userID, gameID)
	gameState.Balance = 1000000
	game.PutGameToMemory(gameService, userID, gameID, gameState)
	_, err := gameService.BuyMiner(userID, gameID, "small", "miner")
	if err != nil {
		t.Fatalf("expected success, got %v:", err)
	}
}

func TestBuyMinerNotEnoughBalance(t *testing.T) {
	userID := "testUserID"
	gameID := "testGameID"

	sessions := MockSessionService{
		isActive: true,
	}
	gameService := game.NewService(game.ServiceDeps{
		Sessions: &sessions,
		Repo:     nil,
		Loop:     nil,
	})
	gameState := domain.NewGameState(userID, gameID)

	game.PutGameToMemory(gameService, userID, gameID, gameState)
	_, err := gameService.BuyMiner(userID, gameID, "small", "miner")
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, errs.ErrNotEnoughBalance) {
		t.Fatalf("expected ErrNotEnoughBalance, got %v:", err)
	}
}

func TestBuyEquipmentAlreadyOwn(t *testing.T) {
	userID := "testUserID"
	gameID := "testGameID"

	sessions := MockSessionService{
		isActive: true,
	}
	gameService := game.NewService(game.ServiceDeps{
		Sessions: &sessions,
		Repo:     nil,
		Loop:     nil,
	})
	gameState := domain.NewGameState(userID, gameID)
	gameState.AddEquipment("1")

	game.PutGameToMemory(gameService, userID, gameID, gameState)
	_, err := gameService.BuyEquipment(userID, gameID, "1", "equipment")
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, errs.ErrAlreadyOwn) {
		t.Fatalf("expected ErrAlreadyOwn, got %v:", err)
	}
}

func TestDeleteExpiredSessionsSuccess(t *testing.T) {
	userID := "testUserID"
	gameID := "testGameID"
	repo := MockGameRepository{
		MockSave: func(gameState *domain.GameState) error {
			return nil
		},
	}
	loop := MockLoopService{}
	sessions := MockSessionService{
		isActive: true,
		MockGetExpired: func() []string {
			return []string{userID + "/" + gameID}
		},
	}

	gameService := game.NewService(game.ServiceDeps{
		Sessions: &sessions,
		Repo:     &repo,
		Loop:     &loop,
	})
	gameState := domain.NewGameState(userID, gameID)
	game.PutGameToMemory(gameService, userID, gameID, gameState)
	gameService.DeleteExpiredSessions()
	if !loop.UnregisterCalled {
		t.Fatalf("expected game to be unregister in loop")
	}
	if !repo.SaveCalled {
		t.Fatalf("expected game to be save in repo")
	}

}

func TestGetHudSuccess(t *testing.T) {
	userID := "testUserID"
	gameID := "testGameID"
	repo := MockGameRepository{
		MockLoad: func(userID, gameID string) (*domain.GameState, error) {
			return &domain.GameState{}, nil
		},
		MockSave: func(gameState *domain.GameState) error {
			return nil
		},
	}
	loop := MockLoopService{}
	sessions := MockSessionService{
		isActive: true,
	}

	gameService := game.NewService(game.ServiceDeps{
		Sessions: &sessions,
		Repo:     &repo,
		Loop:     &loop,
	})

	balance, income, err := gameService.GetHud(userID, gameID)
	if err != nil {
		t.Fatalf("expected success, got %v:", err)
	}
	if balance == "" || income == "" {
		t.Fatalf("expected game to be return hud")
	}
	if !sessions.MarkActiveCalled {
		t.Fatalf("expected game to be mark active")
	}
}
