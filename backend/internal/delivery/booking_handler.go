package delivery

import (
	"net/http"
	"queue_reservation/internal/delivery/dto"
	"queue_reservation/internal/domain"
	"queue_reservation/internal/helper"
	"queue_reservation/internal/models"
	"queue_reservation/internal/models/enum"
	"queue_reservation/pkg/logx"
	"strconv"

	"github.com/labstack/echo/v4"
)

type bookingHandler struct {
	bookingUsecase domain.BookingUsecaseInterface
}

func NewBookingHandler(
	e *echo.Echo,
	bookingUsecase domain.BookingUsecaseInterface,
) {
	handler := &bookingHandler{
		bookingUsecase: bookingUsecase,
	}

	api := e.Group("/api")
	v1 := api.Group("/v1")

	bookingGroup := v1.Group("/restaurants/:restaurant_id/branchs/:branch_id/bookings")
	bookingGroup.POST("", handler.CreateBooking)
	bookingGroup.DELETE("/:booking_id", handler.CancelBooking)
}

// BookingSpec godoc
// @Summary CreateBooking
// @Description CreateBooking
// @Tags Booking
// @Accept json
// @Produce json
// @Success 200 {object} dto.CoreResponse{Result=dto.BookingWithTableResponse}}
// @Error 400 {object} dto.CoreResponse{Result=string}}
// @Security ApiKeyAuth
// @Router /api/v1/restaurant/{restaurant_id}/branch/{branch_id}/booking [post]
// @ID qr_booking_create
func (h *bookingHandler) CreateBooking(c echo.Context) error {

	req := dto.BookingRequest{}

	err := c.Bind(&req)
	if err != nil {
		logx.GetLog().Errorf("CreateBooking Handler Bind Error: %s", err.Error())
		return c.JSON(http.StatusBadRequest, dto.CoreResponse{
			Message: enum.MessageError,
			Result:  err.Error(),
		})
	}

	restaurantIdParam := c.Param("restaurant_id")
	restaurantId, err := strconv.ParseUint(restaurantIdParam, 10, 64)
	if err != nil {
		logx.GetLog().Errorf("CreateBooking Handler Error: %s", err.Error())
		return c.JSON(http.StatusBadRequest, dto.CoreResponse{
			Message: enum.MessageError,
			Result:  err.Error(),
		})
	}

	branchIdParam := c.Param("branch_id")
	branchId, err := strconv.ParseUint(branchIdParam, 10, 64)
	if err != nil {
		logx.GetLog().Errorf("CreateBooking Handler Error: %s", err.Error())
		return c.JSON(http.StatusBadRequest, dto.CoreResponse{
			Message: enum.MessageError,
			Result:  err.Error(),
		})
	}

	if err := helper.NewValidator().Validate(req); err != nil {
		logx.GetLog().Errorf("CreateBooking Handler Validate Error: %s", err[0].Message)
		return c.JSON(http.StatusBadRequest, dto.CoreResponse{
			Message: enum.MessageError,
			Result:  err[0].Message,
		})

	}

	id, tableCount, err := h.bookingUsecase.CreateBooking(
		models.Branch{
			ID:           branchId,
			RestaurantId: restaurantId,
		},
		models.Booking{
			Number: req.NumberOfCustomers,
		},
	)

	if err != nil {
		if err.Error() == enum.ErrorTableNotInitialize || err.Error() == enum.ErrorTableNotEnouugh {
			return c.JSON(http.StatusBadRequest, dto.CoreResponse{
				Message: enum.MessageError,
				Result:  err.Error(),
			})
		} else {
			return c.JSON(http.StatusInternalServerError, dto.CoreResponse{
				Message: enum.MessageError,
				Result:  err.Error(),
			})
		}
	}

	result := dto.BookingWithTableResponse{}
	result.ToBookingWithTableResponse(id, tableCount)

	res := dto.CoreResponse{
		Message: enum.MessageSuccess,
		Result:  result,
	}

	return c.JSON(http.StatusOK, res)
}

func (h *bookingHandler) CancelBooking(c echo.Context) error {

	restaurantIdParam := c.Param("restaurant_id")
	restaurantId, err := strconv.ParseUint(restaurantIdParam, 10, 64)
	if err != nil {
		logx.GetLog().Errorf("CancelBooking Handler Error: %s", err.Error())
		return c.JSON(http.StatusBadRequest, dto.CoreResponse{
			Message: enum.MessageError,
			Result:  err.Error(),
		})
	}

	branchIdParam := c.Param("branch_id")
	branchId, err := strconv.ParseUint(branchIdParam, 10, 64)
	if err != nil {
		logx.GetLog().Errorf("CancelBooking Handler Error: %s", err.Error())
		return c.JSON(http.StatusBadRequest, dto.CoreResponse{
			Message: enum.MessageError,
			Result:  err.Error(),
		})
	}

	bookingIdParam := c.Param("booking_id")
	bookingId, err := strconv.ParseUint(bookingIdParam, 10, 64)
	if err != nil {
		logx.GetLog().Errorf("CancelBooking Handler Error: %s", err.Error())
		return c.JSON(http.StatusBadRequest, dto.CoreResponse{
			Message: enum.MessageError,
			Result:  err.Error(),
		})
	}

	tableCount, err := h.bookingUsecase.CancelBooking(
		models.Branch{
			ID:           branchId,
			RestaurantId: restaurantId,
		},
		models.Booking{
			ID: bookingId,
		},
	)

	if err != nil {
		if err.Error() == enum.ErrorTableNotInitialize || err.Error() == enum.ErrorRecordNotFound {
			return c.JSON(http.StatusBadRequest, dto.CoreResponse{
				Message: enum.MessageError,
				Result:  err.Error(),
			})
		} else {
			return c.JSON(http.StatusInternalServerError, dto.CoreResponse{
				Message: enum.MessageError,
				Result:  err.Error(),
			})
		}
	}

	result := dto.BookingWithTableResponse{}
	result.ToCancelBookingResponse(tableCount)

	res := dto.CoreResponse{
		Message: enum.MessageSuccess,
		Result:  result,
	}

	return c.JSON(http.StatusOK, res)
}
