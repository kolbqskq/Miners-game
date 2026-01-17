package upgrades

type Upgrade struct {
	Name string
	Own  bool
}

type UpgradesConfig struct {
	Price int64
	Value float32
}

var UpgradesPresets = map[string]UpgradesConfig{
	"1": {
		Price: 450,
		Value: 2,
	},
	"2": {
		Price: 450,
		Value: 2,
	},
	"3": {
		Price: 450,
		Value: 2,
	},
}

func GetUpgradesConfig(name string) UpgradesConfig {
	return UpgradesPresets[name]
}

func NewUpgrades() []Upgrade {
	return []Upgrade{
		{
			Name: "1",
			Own:  false,
		},

		{
			Name: "2",
			Own:  false,
		},

		{
			Name: "3",
			Own:  false,
		},
	}
}
