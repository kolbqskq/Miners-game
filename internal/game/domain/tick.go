package domain

import (
	"miners_game/internal/game/equipments"
	"miners_game/internal/game/upgrades"
)

const (
	passiveIncome int64 = 1
)

func (g *GameState) Tick(now int64) {

	g.Mu.Lock()
	defer g.Mu.Unlock()

	if now <= g.LastUpdateAt {
		return
	}
	income := g.CalcIncome(g.LastUpdateAt, now)
	g.IncomePerSec = g.CalcIncome(now-1, now)
	g.Balance += income
	g.LastUpdateAt = now
	g.deleteExpiredMiners(now)

}

func (g *GameState) CalcIncome(from, to int64) int64 {
	total := passiveIncome * (to - from)
	for _, v := range g.Miners {
		total += v.CalcIncome(from, to)
	}
	var rise float32 = 1
	for _, v := range g.Equipments {
		if v.Own {
			rise += equipments.GetEquipmentConfig(v.Name).Value
		}
	}
	for _, v := range g.Upgrades {
		if v.Own {
			rise += upgrades.GetUpgradesConfig(v.Name).Value
		}
	}
	return int64(float32(total) * rise)
}

func (g *GameState) deleteExpiredMiners(now int64) {
	for k, v := range g.Miners {
		if v.EndAt <= now {
			delete(g.Miners, k)
		}
	}
}
