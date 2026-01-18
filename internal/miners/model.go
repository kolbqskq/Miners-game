package miners

import (
	"time"

	"github.com/google/uuid"
)

const (
	MaxMiners int = 20
)

type Miner struct {
	ID      string
	Class   string
	StartAt int64
	EndAt   int64
}

type MinerConfig struct {
	Price    int64
	Power    int64
	Energy   int64
}

var MinerPresets = map[string]MinerConfig{
	"small": {
		Price:    5,
		Power:    1,
		Energy:   30,
	},
	"normal": {
		Price:    50,
		Energy:   45,
		Power:    3,
	},
	"strong": {
		Price:    450,
		Energy:   60,
		Power:    10,
	},
}

func GetMinerConfig(class string) MinerConfig {
	return MinerPresets[class]
}

func NewMiner(class string) *Miner {

	cfg := GetMinerConfig(class)

	now := time.Now().Unix()
	miner := &Miner{
		ID:      uuid.NewString(),
		Class:   class,
		StartAt: now,
		EndAt:   now + cfg.Energy,
	}
	return miner
}
