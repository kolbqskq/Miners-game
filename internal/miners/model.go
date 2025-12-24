package miners

var MinerPresets = map[string]MinerConfig{
	"small": {
		Class:     "small",
		Price:     5,
		Power:     1,
		Energy:    30,
		BreakTime: 3,
		Progress:  0,
	},
	"normal": {
		Class:     "normal",
		Price:     50,
		Energy:    45,
		Power:     3,
		BreakTime: 2,
		Progress:  0,
	},
	"strong": {
		Class:     "strong",
		Price:     450,
		Energy:    60,
		Power:     10,
		BreakTime: 1,
		Progress:  3,
	},
}

type MinerConfig struct {
	Class     string
	Price     int64
	Power     int
	Energy    int
	BreakTime int
	Progress  int
}

type Miner struct {
	ID      string
	StartAt int64
	MinerConfig
}
