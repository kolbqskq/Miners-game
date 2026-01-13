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
	Price     int64
	Power     int64
	Energy    int64
	BreakTime int64
	Progress  int64
}

var MinerPresets = map[string]MinerConfig{
	"small": {
		Price:     5,
		Power:     1,
		Energy:    30,
		BreakTime: 3,
		Progress:  0,
	},
	"normal": {
		Price:     50,
		Energy:    45,
		Power:     3,
		BreakTime: 2,
		Progress:  0,
	},
	"strong": {
		Price:     450,
		Energy:    60,
		Power:     10,
		BreakTime: 1,
		Progress:  3,
	},
}

func GetMinerConfig(class string) MinerConfig {
	return MinerPresets[class]
}

func NewMiner(class string) *Miner {

	cfg := GetMinerConfig(class)
	lifetime := cfg.Energy * cfg.BreakTime

	now := time.Now().Unix()
	miner := &Miner{
		ID:      uuid.NewString(),
		Class:   class,
		StartAt: now,
		EndAt:   now + int64(lifetime),
	}
	return miner
}

func (m *Miner) Gif() string {
	switch m.Class {
	case "small":
		return "/public/gif/miner_small.mp4"
	case "normal":
		return "/public/gif/miner_normal.gif"
	case "strong":
		return "/public/gif/miner_strong.gif"
	default:
		return "public/gif/miner_small.gif"
	}
}