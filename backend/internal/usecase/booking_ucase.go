package usecase

import (
	"queue_reservation/internal/domain"
	"queue_reservation/internal/models"
)

type bookingUsecase struct {
	bookingRepo domain.BookingRepositoryInterface
}

func NewBookingUsecase(bookingRepository domain.BookingRepositoryInterface) domain.BookingUsecaseInterface {
	return &bookingUsecase{
		bookingRepo: bookingRepository,
	}
}

func (u *bookingUsecase) CreateBooking(branch models.Branch, booking models.Booking) (*uint64, *models.TableCount, error) {
	id, tableCount, err := u.bookingRepo.CreateBooking(branch, booking)
	if err != nil {
		return nil, nil, err
	}
	return id, tableCount, nil
}

func (u *bookingUsecase) CancelBooking(branch models.Branch, booking models.Booking) (*models.TableCount, error) {
	tableCount, err := u.bookingRepo.CancelBooking(branch, booking)
	if err != nil {
		return nil, err
	}
	return tableCount, nil
}
