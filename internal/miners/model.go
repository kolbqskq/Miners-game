package miners

import (
	"miners_game/internal/game/shop"
	"sort"
	"strconv"
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
	Price  int64
	Power  int64
	Energy int64
	Title  string
}

var MinerPresets = map[string]MinerConfig{
	"small": {
		Price:  5,
		Power:  1,
		Energy: 30,
		Title:  "Шахтер",
	},
	"normal": {
		Price:  50,
		Energy: 45,
		Power:  3,
		Title:  "Шахтер+",
	},
	"strong": {
		Price:  450,
		Energy: 60,
		Power:  10,
		Title:  "Гигабайт",
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

func MinerShopCards() []shop.ShopCard {
	cards := make([]shop.ShopCard, 0, len(MinerPresets))
	for k, v := range MinerPresets {
		card := shop.ShopCard{
			ID:       "miner-" + k,
			Title:    v.Title,
			Income:   "+" + strconv.Itoa(int(v.Power)) + ".0/сек",
			Duration: "⏱" + strconv.Itoa(int(v.Energy)) + " сек",
			Price:    strconv.Itoa(int(v.Price)),
			Name:     k,
			Kind:     "miner",
			Icon:     "/public/icons/shop/miner-" + k + ".png",
			Disabled: false,
			Reason:   "",
		}
		cards = append(cards, card)
	}
	sort.Slice(cards, func(i, j int) bool {
		pi, _ := strconv.Atoi(cards[i].Price)
		pj, _ := strconv.Atoi(cards[j].Price)
		return pi < pj
	})
	return cards
}
