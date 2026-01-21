package game

import (
	"errors"
	"miners_game/internal/game/domain"
	"miners_game/internal/game/equipments"
	"miners_game/internal/game/loop"
	"miners_game/internal/game/sessions"
	"miners_game/internal/game/shop"
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

func (a *Service) enterGame(userID, gameID string) (*domain.GameState, error) {
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

func (a *Service) buyMiner(userID, gameID, class, kind string) (shop.ShopCard, error) {
	game, err := a.buy(userID, gameID)
	if err != nil {
		return shop.ShopCard{}, err
	}
	price := miners.GetMinerConfig(class).Price
	if err := game.SpendBalance(price); err != nil {
		return getErrShopCard(class, kind, err.Error()), err
	}
	game.AddMiner(class)

	return shop.ShopCard{}, nil
}

func (a *Service) buyEquipment(userID, gameID, name, kind string) (shop.ShopCard, error) {
	game, err := a.buy(userID, gameID)
	if err != nil {
		return shop.ShopCard{}, err
	}
	if game.IsOwnEquipment(name) {
		return getErrShopCard(name, kind, errs.ErrAlreadyOwn.Error()), errs.ErrAlreadyOwn
	}

	price := equipments.GetEquipmentConfig(name).Price
	if err := game.SpendBalance(price); err != nil {
		return getErrShopCard(name, kind, err.Error()), err
	}
	game.AddEquipment(name)
	return shop.ShopCard{}, nil
}

func (a *Service) buyUpgrade(userID, gameID, name, kind string) (shop.ShopCard, error) {
	game, err := a.buy(userID, gameID)
	if err != nil {
		return shop.ShopCard{}, err
	}
	if game.IsOwnUpgrade(name) {
		return getErrShopCard(name, kind, errs.ErrAlreadyOwn.Error()), errs.ErrAlreadyOwn
	}

	price := upgrades.GetUpgradesConfig(name).Price
	if err := game.SpendBalance(price); err != nil {
		return getErrShopCard(name, kind, err.Error()), err
	}
	game.AddUpgrade(name)

	return shop.ShopCard{}, nil
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

func (a *Service) getHud(userID, gameID string) (string, string, error) {
	id := userID + "/" + gameID

	game, err := a.enterGame(userID, gameID)
	if err != nil {
		return "", "", err
	}

	game.Mu.Lock()
	balance := strconv.Itoa(int(game.Balance))
	income := strconv.Itoa(int(game.IncomePerSec))
	game.Mu.Unlock()

	a.sessions.MarkActive(id)

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
		return nil, errs.ErrSessionIsNotActive
	}

	a.mu.RLock()
	game, ok := a.games[id]
	a.mu.RUnlock()

	if !ok {
		return nil, errs.ErrGameNotFound
	}
	return game, nil
}

func (a *Service) getShopState(kind string) []shop.ShopCard {
	switch kind {
	case "miner":
		return miners.MinerShopCards()
	case "equipment":
		return equipments.EquipmentShopCards()
	case "upgrade":
		return upgrades.UpgradeShopCards()
	}
	return nil
}

func getErrShopCard(name, kind, reason string) shop.ShopCard {
	card := getShopCardByName(name, kind)
	card.Disabled = true
	card.Reason = reason
	return card
}

func getShopCardByName(name, kind string) shop.ShopCard {
	cases := map[string]func() []shop.ShopCard{
		"miner":     miners.MinerShopCards,
		"equipment": equipments.EquipmentShopCards,
		"upgrade":   upgrades.UpgradeShopCards,
	}
	cards := cases[kind]()
	for _, v := range cards {
		if v.Name == name {
			return v
		}
	}
	return shop.ShopCard{}
}
