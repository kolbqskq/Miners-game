package upgrades

import (
	"miners_game/internal/game/shop"
	"sort"
	"strconv"
)

type Upgrade struct {
	Name string
	Own  bool
}

type UpgradesConfig struct {
	Price int64
	Value int64
}

var UpgradesPresets = map[string]UpgradesConfig{
	"1": {
		Price: 450,
		Value: 200,
	},
	"2": {
		Price: 450,
		Value: 200,
	},
	"3": {
		Price: 450,
		Value: 200,
	},
}

func GetUpgradesConfig(name string) UpgradesConfig {
	return UpgradesPresets[name]
}

func NewUpgrades() []Upgrade {
	return []Upgrade{
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
}

func UpgradeShopCards() []shop.ShopCard {
	cards := make([]shop.ShopCard, 0, len(UpgradesPresets))
	for k, v := range UpgradesPresets {
		card := shop.ShopCard{
			ID:       "upgrade-" + k,
			Title:    k,
			Income:   strconv.Itoa(int(v.Value)),
			Duration: "",
			Price:    strconv.Itoa(int(v.Price)),
			Name:     k,
			Kind:     "upgrade",
			Icon:     "/public/icons/shop/upgrade-" + k + ".png",
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
