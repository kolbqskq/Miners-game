package game

import (
	"miners_game/internal/equipments"
	"miners_game/internal/miners"
	"sync"
	"time"
)

type GameState struct {
	UserID       string
	SaveID       string
	Balance      int64
	LastUpdateAt int64
	LaseSeen     time.Time
	Miners       map[string]*miners.Miner
	Equipments   []equipments.Equipment
	mu           sync.RWMutex
}

func NewGameState(userID, saveID string) *GameState {
	equipments := []equipments.Equipment{
		{
			Name: "1",
			Own:  false,
		},
		{
			Name: "2",
			Own:  false,
		},
		{
			Name: "3",
			Own:  false,
		},
	}
	return &GameState{
		UserID:       userID,
		SaveID:       saveID,
		Balance:      0,
		LastUpdateAt: time.Now().Unix(),
		LaseSeen:     time.Now(),
		Miners:       make(map[string]*miners.Miner),
		Equipments:   equipments,
	}
}
