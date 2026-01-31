package game

import (
	"errors"
	"miners_game/internal/game/domain"
	"miners_game/internal/game/equipments"
	"miners_game/internal/game/shop"
	"miners_game/internal/game/upgrades"
	"miners_game/internal/miners"
	"miners_game/pkg/errs"
	"strconv"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type Service struct {
	repo     IGameRepository
	loop     ILoopService
	sessions ISessionService

	games   map[string]*domain.GameState
	logger  zerolog.Logger
	metrics *Metrics
	mu      sync.RWMutex
}

type ServiceDeps struct {
	Repo     IGameRepository
	Loop     ILoopService
	Sessions ISessionService
	Metrics  *Metrics
	Logger   zerolog.Logger
}

func NewService(deps ServiceDeps) *Service {
	return &Service{
		repo:     deps.Repo,
		loop:     deps.Loop,
		sessions: deps.Sessions,
		logger:   deps.Logger,
		games:    make(map[string]*domain.GameState),
		metrics:  deps.Metrics,
	}
}

func (s *Service) EnterGame(userID, gameID string) (*domain.GameState, error) {
	id := userID + "/" + gameID

	s.mu.RLock()
	if game, ok := s.games[id]; ok {
		s.mu.RUnlock()
		s.sessions.MarkActive(id)
		return game, nil
	}
	s.mu.RUnlock()

	game, err := s.repo.Load(userID, gameID)
	if err != nil {
		if !errors.Is(err, errs.ErrGameNotFound) {
			return nil, err
		}
		game = domain.NewGameState(userID, gameID)
		s.repo.Save(game)
	}

	now := time.Now().Unix()

	if now-game.LastUpdateAt > 5 {
		game.LastUpdateAt = now
	}

	s.mu.Lock()
	s.games[id] = game
	s.mu.Unlock()

	s.loop.Register(id, game)
	s.sessions.MarkActive(id)

	s.logger.Info().Str("user_id", userID).Str("game_id", gameID).Msg("game entered")
	return game, nil
}

func (s *Service) BuyMiner(userID, gameID, class, kind string) (shop.ShopCard, error) {
	var err error
	defer s.buyMetrics(&err)()
	var game *domain.GameState
	game, err = s.GetGameState(userID, gameID)
	if err != nil {
		return shop.ShopCard{}, err
	}
	price := miners.GetMinerConfig(class).Price
	if err = game.SpendBalance(price); err != nil {
		return getErrShopCard(class, kind, err.Error()), err
	}
	game.AddMiner(class)

	return shop.ShopCard{}, nil
}

func (s *Service) BuyEquipment(userID, gameID, name, kind string) (shop.ShopCard, error) {
	var err error
	defer s.buyMetrics(&err)()
	var game *domain.GameState
	game, err = s.GetGameState(userID, gameID)
	if err != nil {
		return shop.ShopCard{}, err
	}
	if game.IsOwnEquipment(name) {
		return getErrShopCard(name, kind, errs.ErrAlreadyOwn.Error()), errs.ErrAlreadyOwn
	}

	price := equipments.GetEquipmentConfig(name).Price
	if err = game.SpendBalance(price); err != nil {
		return getErrShopCard(name, kind, err.Error()), err
	}
	game.AddEquipment(name)
	return shop.ShopCard{}, nil
}

func (s *Service) BuyUpgrade(userID, gameID, name, kind string) (shop.ShopCard, error) {
	var err error
	defer s.buyMetrics(&err)()
	var game *domain.GameState
	game, err = s.GetGameState(userID, gameID)
	if err != nil {
		return shop.ShopCard{}, err
	}
	if game.IsOwnUpgrade(name) {
		return getErrShopCard(name, kind, errs.ErrAlreadyOwn.Error()), errs.ErrAlreadyOwn
	}

	price := upgrades.GetUpgradesConfig(name).Price
	if err = game.SpendBalance(price); err != nil {
		return getErrShopCard(name, kind, err.Error()), err
	}
	game.AddUpgrade(name)

	return shop.ShopCard{}, nil
}

func (s *Service) getCurrUpgrade(userID, gameID string) (string, error) {
	game, err := s.GetGameState(userID, gameID)
	if err != nil {
		return "", err
	}
	game.Mu.RLock()
	curr := game.GetMaxUpgrade()
	game.Mu.RUnlock()
	return curr, nil
}

func (s *Service) DeleteExpiredSessions() {
	expired := s.sessions.GetExpired()
	for _, id := range expired {
		s.loop.Unregister(id)

		s.mu.Lock()
		game := s.games[id]
		delete(s.games, id)
		s.mu.Unlock()

		if game != nil {
			if err := s.repo.Save(game); err != nil {
				s.logger.Error().Err(err).Msg("failed to delete expired sessions")
				return
			}
		}
	}
	if len(expired) > 0 {
		s.logger.Info().Int("count", len(expired)).Msg("expired sessions deleted")
	}
}

func (s *Service) GetHud(userID, gameID string) (string, string, error) {
	id := userID + "/" + gameID

	game, err := s.EnterGame(userID, gameID)
	if err != nil {
		return "", "", err
	}

	game.Mu.Lock()
	balance := strconv.Itoa(int(game.Balance))
	income := strconv.Itoa(int(game.IncomePerSec))
	game.Mu.Unlock()

	s.sessions.MarkActive(id)

	return balance, income, nil
}

func (s *Service) SaveAll() {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, game := range s.games {
		game.Mu.RLock()
		s.repo.Save(game)
		game.Mu.RUnlock()
	}
	if len(s.games) > 0 {
		s.logger.Info().Int("count", len(s.games)).Msg("saved all games")
	}
}

func (s *Service) GetGameState(userID, gameID string) (*domain.GameState, error) {
	id := userID + "/" + gameID
	if !s.sessions.IsActive(id) {
		return nil, errs.ErrSessionIsNotActive
	}

	s.mu.RLock()
	game, ok := s.games[id]
	s.mu.RUnlock()

	if !ok {
		return nil, errs.ErrGameNotFound
	}
	return game, nil
}

func (s *Service) getShopState(kind string) []shop.ShopCard {
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
	card := GetShopCardByName(name, kind)
	card.Disabled = true
	card.Reason = reason
	return card
}

func GetShopCardByName(name, kind string) shop.ShopCard {
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

func (s *Service) buyMetrics(err *error) func() {
	return func() {
		if s.metrics == nil {
			return
		}
		s.metrics.BuyAttemptsTotal.Inc()
		if *err != nil {
			s.metrics.BuyFailedTotal.Inc()
		} else {
			s.metrics.BuySuccessTotal.Inc()
		}
	}
}
