package repository

import (
	"errors"
	"queue_reservation/internal/domain"
	"queue_reservation/internal/models"
	"queue_reservation/internal/models/enum"
	"queue_reservation/internal/repository/dbmodels"
	"queue_reservation/pkg/logx"
	"time"
)

type bookingRepository struct {
}

func NewBookingRepository() domain.BookingRepositoryInterface {
	return &bookingRepository{}
}

func (r *bookingRepository) CreateBooking(branch models.Branch, booking models.Booking) (*uint64, *models.TableCount, error) {
	tableDB := dbmodels.TableDB
	branchDB := dbmodels.BranchDB
	remainingTables := 0
	newTables := 0
	numberRemainingCustomers := booking.Number
	numberRemainingSeats := 0
	remainingTableIdx := []int{}

	for bidx := range branchDB {
		// check correct restaurant and branch
		if branchDB[bidx].ID == branch.ID && branchDB[bidx].RestaurantId == branch.RestaurantId {
			if branchDB[bidx].IsInit {
				for tidx := range tableDB {
					if tableDB[tidx].BranchID == branch.ID {
						// table is not reserved
						if tableDB[tidx].BookingId == nil {
							// count table empty
							remainingTables += 1
							// store table that empty
							remainingTableIdx = append(remainingTableIdx, tidx)
							// store remaining seat
							numberRemainingSeats += tableDB[tidx].SeatMax
						}
					}
				}
			} else {
				// Return error if this API is called before initialization.
				err := errors.New(enum.ErrorTableNotInitialize)
				logx.GetLog().Errorf("CreateBooking Repo Error: %s", err.Error())
				return nil, nil, err
			}
		}
	}

	// there are not enough tables for all customers in the reservation
	if numberRemainingSeats < booking.Number {
		err := errors.New(enum.ErrorTableNotEnouugh)
		logx.GetLog().Errorf("CreateBooking Repo Error: %s", err.Error())
		return nil, nil, err
	}

	// create id for new booking
	newBookingId := enum.NextBookingId
	// create new booking
	newBooking := dbmodels.Booking{
		ID:        newBookingId,
		BranchID:  branch.ID,
		Number:    booking.Number,
		DeletedAt: nil,
	}
	dbmodels.BookingDB = append(dbmodels.BookingDB, newBooking)

	// set id for next booking
	enum.NextBookingId += 1

	// reserve in remaining table
	for _, tableID := range remainingTableIdx {
		if numberRemainingCustomers > 0 {
			tableDB[tableID].BookingId = &newBookingId
			numberRemainingCustomers -= tableDB[tableID].SeatMax
			remainingTables -= 1
			newTables += 1
		} else {
			break
		}
	}

	tableCount := models.TableCount{
		NumberOfNewReservedTables: newTables,
		NumberOfRemainingTables:   remainingTables,
	}

	return &newBookingId, &tableCount, nil
}

func (r *bookingRepository) CancelBooking(branch models.Branch, booking models.Booking) (*models.TableCount, error) {
	tableDB := dbmodels.TableDB
	branchDB := dbmodels.BranchDB
	bookingDB := dbmodels.BookingDB
	remainingTables := 0
	hasBooking := false

	for bidx := range branchDB {
		// check correct restaurant and branch
		if branchDB[bidx].ID == branch.ID && branchDB[bidx].RestaurantId == branch.RestaurantId {
			if branchDB[bidx].IsInit {
				for boidx := range bookingDB {
					if bookingDB[boidx].BranchID == branch.ID {
						// table is not reserved
						if bookingDB[boidx].ID == booking.ID && bookingDB[boidx].DeletedAt == nil {
							// soft delete this booking id
							now := time.Now().UTC()
							bookingDB[boidx].DeletedAt = &now
							hasBooking = true
						}
					}
				}

				if !hasBooking {
					err := errors.New(enum.ErrorRecordNotFound)
					logx.GetLog().Errorf("CancelBooking Repo Error: %s", err.Error())
					return nil, err
				}

				for tidx := range tableDB {
					if tableDB[tidx].BranchID == branch.ID {
						if tableDB[tidx].BookingId == nil {
							// count table empty
							remainingTables += 1
						} else if *tableDB[tidx].BookingId == booking.ID { // table is reserved
							// cancel the reservation
							tableDB[tidx].BookingId = nil
							remainingTables += 1

						}
					}
				}

			} else {
				// Return error if this API is called again after initialization.
				err := errors.New(enum.ErrorTableNotInitialize)
				logx.GetLog().Errorf("CancelBooking Repo Error: %s", err.Error())
				return nil, err
			}
		}
	}

	tableCount := models.TableCount{
		NumberOfRemainingTables: remainingTables,
	}

	return &tableCount, nil
}
