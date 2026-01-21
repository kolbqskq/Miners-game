package domain

import (
	"miners_game/pkg/errs"
)

func (g *GameState) SpendBalance(price int64) error {
	g.Mu.Lock()
	defer g.Mu.Unlock()
	if g.Balance < price {
		return errs.ErrNotEnoughBalance
	}
	g.Balance -= price
	return nil
}
