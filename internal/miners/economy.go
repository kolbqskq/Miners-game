package miners

func (m *Miner) CalcIncome(from, to int64) int64 {
	cfg := GetMinerConfig(m.Class)
	if from < m.StartAt {
		from = m.StartAt
	}
	if to > m.EndAt {
		to = m.EndAt
	}
	if to <= from {
		return 0
	}
	seconds := to - from

	return seconds * cfg.Power

}
