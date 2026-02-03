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
	Title string
	Price int64
	Value int64
}

var UpgradesPresets = map[string]UpgradesConfig{
	"1": {
		Title: "Индастриал",
		Price: 1500,
		Value: 50,
	},
	"2": {
		Title: "Энергосети",
		Price: 5000,
		Value: 150,
	},
	"3": {
		Title: "Глобальная логистика",
		Price: 15000,
		Value: 300,
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
			Title:    v.Title,
			Income:   "+" + strconv.Itoa(int(v.Value)) + "%",
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
