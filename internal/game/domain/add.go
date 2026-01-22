package domain

import (
	"miners_game/internal/miners"
	"sort"
	"strconv"

	"github.com/google/uuid"
)

func (g *GameState) AddMiner(class string) {
	g.Mu.Lock()
	defer g.Mu.Unlock()
	g.Miners[uuid.New().String()] = miners.NewMiner(class)
}

func (g *GameState) AddEquipment(name string) {
	g.Mu.Lock()
	defer g.Mu.Unlock()
	for k := range g.Equipments {
		if g.Equipments[k].Name == name {
			g.Equipments[k].Own = true
			return
		}
	}
}

func (g *GameState) AddUpgrade(name string) {
	g.Mu.Lock()
	defer g.Mu.Unlock()
	for k := range g.Upgrades {
		if g.Upgrades[k].Name == name {
			g.Upgrades[k].Own = true
			return
		}
	}
}

func (g *GameState) GetMaxUpgrade() string {
	allOwnedUpgrades := make([]string, 0, 3)
	for _, v := range g.Upgrades {
		if v.Own {
			allOwnedUpgrades = append(allOwnedUpgrades, v.Name)
		}
	}
	if len(allOwnedUpgrades) == 0 {
		return "0"
	}
	sort.Slice(allOwnedUpgrades, func(i, j int) bool {
		pi, _ := strconv.Atoi(allOwnedUpgrades[i])
		pj, _ := strconv.Atoi(allOwnedUpgrades[j])
		return pi > pj
	})
	return allOwnedUpgrades[0]
}
