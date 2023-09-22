package domain

import "queue_reservation/internal/models"

type BookingUsecaseInterface interface {
	CreateBooking(branch models.Branch, booking models.Booking) (*uint64, *models.TableCount, error)
	CancelBooking(branch models.Branch, booking models.Booking) (*models.TableCount, error)
}

type BookingRepositoryInterface interface {
	CreateBooking(branch models.Branch, booking models.Booking) (*uint64, *models.TableCount, error)
	CancelBooking(branch models.Branch, booking models.Booking) (*models.TableCount, error)
}
