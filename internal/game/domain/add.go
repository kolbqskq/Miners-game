package domain

import (
	"miners_game/internal/miners"

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
