package domain

import (
	"miners_game/internal/game/equipments"
	"miners_game/internal/game/upgrades"
	"miners_game/internal/miners"
	"sync"
	"time"
)

type GameState struct {
	UserID string
	GameID string

	Balance      int64
	IncomePerSec int64

	LastUpdateAt int64

	Miners     map[string]*miners.Miner
	Equipments []equipments.Equipment
	Upgrades   []upgrades.Upgrade

	Mu sync.RWMutex
}

func NewGameState(userID, gameID string) *GameState {
	equipments := equipments.NewEquipments()
	upgrades := upgrades.NewUpgrades()
	return &GameState{
		UserID:       userID,
		GameID:       gameID,
		Balance:      0,
		IncomePerSec: 1,
		LastUpdateAt: time.Now().Unix(),
		Miners:       make(map[string]*miners.Miner),
		Equipments:   equipments,
		Upgrades:     upgrades,
	}
}
