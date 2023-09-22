package dbmodels

type Table struct {
	ID        uint64
	BranchID  uint64
	BookingId *uint64
	SeatMax   int
}

var (
	TableDB = []Table{
		{ID: 1, BranchID: 1, BookingId: nil, SeatMax: 4},
		{ID: 2, BranchID: 1, BookingId: nil, SeatMax: 4},
		{ID: 3, BranchID: 1, BookingId: nil, SeatMax: 4},
		{ID: 4, BranchID: 1, BookingId: nil, SeatMax: 4},
		{ID: 5, BranchID: 1, BookingId: nil, SeatMax: 4},
		{ID: 6, BranchID: 1, BookingId: nil, SeatMax: 4},
		{ID: 7, BranchID: 1, BookingId: nil, SeatMax: 4},
		{ID: 8, BranchID: 1, BookingId: nil, SeatMax: 4},
		{ID: 9, BranchID: 1, BookingId: nil, SeatMax: 4},
		{ID: 10, BranchID: 1, BookingId: nil, SeatMax: 4},
	}
)
