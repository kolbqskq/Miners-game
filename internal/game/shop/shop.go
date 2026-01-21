package shop

type ShopCard struct {
	ID       string
	Title    string
	Income   string
	Duration string
	Price    string
	Name     string
	Kind     string
	Icon     string

	Disabled bool
	Reason   string
}
