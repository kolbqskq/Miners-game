package game

import (
	"miners_game/internal/miners"
	"sync"
)

type GameState struct {
	UserID       string
	SaveID       string
	Balance      int64
	LastUpdateAt int64
	Miners       []miners.Miner
	mu           sync.RWMutex
}
