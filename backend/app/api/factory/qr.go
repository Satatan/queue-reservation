package factory

import (
	"github.com/labstack/echo/v4"

	qrDelivery "queue_reservation/internal/delivery"
	qrRepository "queue_reservation/internal/repository"
	qrUsercase "queue_reservation/internal/usecase"
)

func QRDependencyResolve(e *echo.Echo) {

	qrTableRepo := qrRepository.NewTableRepository()
	qrBookingRepo := qrRepository.NewBookingRepository()

	qrTableUsecase := qrUsercase.NewTableUsecase(qrTableRepo)
	qrBookingUsecase := qrUsercase.NewBookingUsecase(qrBookingRepo)

	qrDelivery.NewTableHandler(e, qrTableUsecase)
	qrDelivery.NewBookingHandler(e, qrBookingUsecase)
}
