package game

import (
	"errors"
	"time"
)

const (
	passiveIncome int64 = 1
)

func (g *GameState) RecalculateBalance() {
	now := time.Now().Unix()

	if now <= g.LastUpdateAt {
		return
	}
	income := g.CalcIncome(g.LastUpdateAt, now)

	g.mu.Lock()
	defer g.mu.Unlock()
	g.Balance += income
	g.LastUpdateAt = now
	g.DeleteMiners(now)
}

func (g *GameState) CalcIncome(from, to int64) int64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	total := passiveIncome * (to - from)

	for _, v := range g.Miners {
		total += v.CalcIncome(from, to)
	}

	return total
}

func (g *GameState) DeleteMiners(now int64) {
	for k, v := range g.Miners {
		if v.EndAt <= now {
			delete(g.Miners, k)
		}
	}
}

func (g *GameState) ValidateBalance(price int64) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Balance < price {
		return errors.New("") //error
	}
	return nil
}
