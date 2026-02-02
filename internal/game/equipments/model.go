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
	Title string
	Price int64
	Value int64
}

var EquipmentPresets = map[string]EquipmentConfig{
	"1": {
		Title: "Кирка",
		Price: 200,
		Value: 10,
	},
	"2": {
		Title: "Бур",
		Price: 600,
		Value: 25,
	},
	"3": {
		Title: "Динамит",
		Price: 1800,
		Value: 40,
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
			Title:    v.Title,
			Income:   "+" + strconv.Itoa(int(v.Value)) + "%",
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
