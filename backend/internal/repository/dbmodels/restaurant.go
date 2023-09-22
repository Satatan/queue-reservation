package dbmodels

type Restaurant struct {
	ID   uint64
	Name string
}

var (
	RestaurantDB = []Restaurant{
		{
			ID:   1,
			Name: "OneSiam",
		},
	}
)
