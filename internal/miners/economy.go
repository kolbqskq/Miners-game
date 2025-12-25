package miners

func (g *Miner) CalcIncome(from, to int64) int64 {
	cfg := GetMinerConfig(g.Class)

	var times int64
	if (to - from) >= (g.EndAt - g.StartAt) {
		times = (g.EndAt - g.StartAt) / cfg.BreakTime
	} else {
		times = (to - from) / cfg.BreakTime
	}
	if times <= 0 {
		return 0
	}
	income := times*cfg.Power + (times*(times-1)/2)*cfg.Progress

	return income

}