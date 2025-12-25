package equipments

type Equipment struct {
	Name string
	Own  bool
}

type EquipmentConfig struct {
	Price int
}

var EquipmentPresets = map[string]EquipmentConfig{
	"1": {
		Price: 450,
	},
	"2": {
		Price: 450,
	},
	"3": {
		Price: 450,
	},
}

func GetEquipmentConfig(name string) EquipmentConfig {
	return EquipmentPresets[name]
}
