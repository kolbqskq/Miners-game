package domain

func (g *GameState) IsOwnEquipment(name string) bool {
	for _, v := range g.Equipments {
		if v.Name == name {
			return v.Own
		}
	}
	return false
}

func (g *GameState) IsOwnUpgrade(name string) bool {
	for _, v := range g.Upgrades {
		if v.Name == name {
			return v.Own
		}
	}
	return false
}