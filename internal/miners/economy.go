package miners

func (m *Miner) CalcIncome(from, to int64) int64 {
	cfg := GetMinerConfig(m.Class)
	if from < m.StartAt {
		from = m.StartAt
	}
	if to > m.EndAt {
		to = m.EndAt
	}
	times := (to - from) / cfg.BreakTime
	if times <= 0 {
		return 0
	}
	income := times*cfg.Power + (times*(times-1)/2)*cfg.Progress

	return income

}
