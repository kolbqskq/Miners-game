package equipments

import (
	"miners_game/internal/game/shop"
	"sort"
	"strconv"
)

type Equipment struct {
	Name string
	Own  bool
}

type EquipmentConfig struct {
	Price int64
	Value int64
}

var EquipmentPresets = map[string]EquipmentConfig{
	"1": {
		Price: 450,
		Value: 110,
	},
	"2": {
		Price: 450,
		Value: 130,
	},
	"3": {
		Price: 450,
		Value: 200,
	},
}

func GetEquipmentConfig(name string) EquipmentConfig {
	return EquipmentPresets[name]
}

func NewEquipments() []Equipment {
	return []Equipment{
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

func EquipmentShopCards() []shop.ShopCard {
	cards := make([]shop.ShopCard, 0, len(EquipmentPresets))
	for k, v := range EquipmentPresets {
		card := shop.ShopCard{
			ID:       "equipment-" + k,
			Title:    k,
			Income:   strconv.Itoa(int(v.Value)),
			Duration: "",
			Price:    strconv.Itoa(int(v.Price)),
			Name:     k,
			Kind:     "equipment",
			Icon:     "/public/icons/shop/equipment-" + k + ".png",
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
