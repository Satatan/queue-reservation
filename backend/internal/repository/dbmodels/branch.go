package dbmodels

type Branch struct {
	ID           uint64
	RestaurantId uint64
	Name         string
	IsInit       bool
}

var (
	BranchDB = []Branch{
		{
			ID:           1,
			RestaurantId: 1,
			Name:         "Bangkok",
			IsInit:       false,
		},
	}
)
