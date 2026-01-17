package game

import (
	"errors"
	"miners_game/internal/game/domain"
	"miners_game/internal/game/equipments"
	"miners_game/internal/game/loop"
	"miners_game/internal/game/sessions"
	"miners_game/internal/game/upgrades"
	"miners_game/internal/miners"
	"miners_game/pkg/errs"
	"strconv"
	"sync"
	"time"
)

type Service struct {
	repo     IGameRepository
	loop     *loop.Service
	sessions *sessions.Service

	games map[string]*domain.GameState
	mu    sync.RWMutex
}

type ServiceDeps struct {
	Repo     IGameRepository
	Loop     *loop.Service
	Sessions *sessions.Service
}

func NewService(deps ServiceDeps) *Service {
	return &Service{
		repo:     deps.Repo,
		loop:     deps.Loop,
		sessions: deps.Sessions,
		games:    make(map[string]*domain.GameState),
	}
}

func (a *Service) EnterGame(userID, gameID string) (*domain.GameState, error) {
	id := userID + "/" + gameID

	a.mu.RLock()
	if game, ok := a.games[id]; ok {
		a.mu.RUnlock()
		a.sessions.MarkActive(id)
		return game, nil
	}
	a.mu.RUnlock()

	game, err := a.repo.Load(userID, gameID)
	if err != nil {
		if !errors.Is(err, errs.ErrGameNotFound) {
			return nil, err
		}
		game = domain.NewGameState(userID, gameID)
		a.repo.Save(game)
	}

	now := time.Now().Unix()

	if now-game.LastUpdateAt > 5 {
		game.LastUpdateAt = now
	}

	a.mu.Lock()
	a.games[id] = game
	a.mu.Unlock()

	a.loop.Register(id, game)
	a.sessions.MarkActive(id)

	return game, nil
}

func (a *Service) Heartbeat(userID, gameID string) {
	id := userID + "/" + gameID
	a.sessions.MarkActive(id)
}

func (a *Service) BuyMiner(userID, gameID, class string) error {
	game, err := a.buy(userID, gameID)
	if err != nil {
		return err
	}
	price := miners.GetMinerConfig(class).Price
	if err := game.SpendBalance(price); err != nil {
		return err
	}
	game.AddMiner(class)

	return nil
}

func (a *Service) BuyEquipment(userID, gameID, name string) error {
	game, err := a.buy(userID, gameID)
	if err != nil {
		return err
	}

	price := equipments.GetEquipmentConfig(name).Price
	if err := game.SpendBalance(price); err != nil {
		return err
	}
	game.AddEquipment(name)
	return nil
}

func (a *Service) BuyUpgrade(userID, gameID, name string) error {
	game, err := a.buy(userID, gameID)
	if err != nil {
		return err
	}

	price := upgrades.GetUpgradesConfig(name).Price
	if err := game.SpendBalance(price); err != nil {
		return err
	}
	game.AddUpgrade(name)

	return nil
}

func (a *Service) HandleExpiredSessions() {
	expired := a.sessions.CheckExpired()
	for _, id := range expired {
		a.loop.Unregister(id)

		a.mu.Lock()
		game := a.games[id]
		delete(a.games, id)
		a.mu.Unlock()

		if game != nil {
			a.repo.Save(game)
		}
	}
}

func (a *Service) GetHud(userID, gameID string) (string, string, error) {
	game, err := a.EnterGame(userID, gameID)
	if err != nil {
		return "", "", err
	}
	
	game.Mu.Lock()
	balance := strconv.Itoa(int(game.Balance))
	income := strconv.Itoa(int(game.IncomePerSec))
	game.Mu.Unlock()
	return balance, income, nil
}

func (a *Service) SaveAll() {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for _, game := range a.games {
		game.Mu.RLock()
		a.repo.Save(game)
		game.Mu.RUnlock()
	}
}

func (a *Service) buy(userID, gameID string) (*domain.GameState, error) {
	id := userID + "/" + gameID
	if !a.sessions.IsActive(id) {
		return nil, errors.New("NotActive") // error
	}

	a.mu.RLock()
	game, ok := a.games[id]
	a.mu.RUnlock()

	if !ok {
		return nil, errors.New("game not loaded")
	}

	return game, nil
}
