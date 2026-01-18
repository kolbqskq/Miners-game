package equipments

type Equipment struct {
	Name string
	Own  bool
}

type EquipmentConfig struct {
	Price int64
	Value int64
}

var EquipmentPresets = map[string]EquipmentConfig{
	"1": {
		Price: 450,
		Value: 110,
	},
	"2": {
		Price: 450,
		Value: 130,
	},
	"3": {
		Price: 450,
		Value: 200,
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
