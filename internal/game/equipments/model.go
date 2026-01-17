package equipments

type Equipment struct {
	Name string
	Own  bool
}

type EquipmentConfig struct {
	Price int64
	Value float32
}

var EquipmentPresets = map[string]EquipmentConfig{
	"1": {
		Price: 450,
		Value: 1.1,
	},
	"2": {
		Price: 450,
		Value: 1.3,
	},
	"3": {
		Price: 450,
		Value: 2,
	},
}

func GetEquipmentConfig(name string) EquipmentConfig {
	return EquipmentPresets[name]
}

func NewEquipments() []Equipment {
	return []Equipment{
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
