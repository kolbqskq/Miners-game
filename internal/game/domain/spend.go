package domain

import "errors"

func (g *GameState) SpendBalance(price int64) error {
	g.Mu.Lock()
	defer g.Mu.Unlock()
	if g.Balance < price {
		return errors.New("") //error
	}
	g.Balance -= price
	return nil
}
