package game

import (
	"errors"
	"time"
)

const (
	passiveIncome int64 = 1
)

func (g *GameState) RecalculateBalance() int64 {
	now := time.Now().Unix()

	g.mu.Lock()
	defer g.mu.Unlock()
	
	if now <= g.LastUpdateAt {
		return 0
	}
	income := g.CalcIncome(g.LastUpdateAt, now)

	g.Balance += income
	g.LastUpdateAt = now
	g.DeleteExpiredMiners(now)
	return income
}

func (g *GameState) CalcIncome(from, to int64) int64 {
	total := passiveIncome * (to - from)
	for _, v := range g.Miners {
		total += v.CalcIncome(from, to)
	}
	return total
}

func (g *GameState) DeleteExpiredMiners(now int64) {
	for k, v := range g.Miners {
		if v.EndAt <= now {
			delete(g.Miners, k)
		}
	}
}

func (g *GameState) SpendBalance(price int64) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Balance < price {
		return errors.New("") //error
	}
	g.Balance -= price
	return nil
}
