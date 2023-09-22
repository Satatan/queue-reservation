package dbmodels

import "time"

type Booking struct {
	ID        uint64
	BranchID  uint64
	Number    int
	DeletedAt *time.Time
}

var (
	BookingDB = []Booking{}
)
