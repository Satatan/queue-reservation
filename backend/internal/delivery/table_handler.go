package delivery

import (
	"net/http"
	"queue_reservation/internal/delivery/dto"
	"queue_reservation/internal/domain"
	"queue_reservation/internal/models"
	"queue_reservation/internal/models/enum"
	"queue_reservation/pkg/logx"
	"strconv"

	"github.com/labstack/echo/v4"
)

type tableHandler struct {
	tableUsecase domain.TableUsecaseInterface
}

func NewTableHandler(
	e *echo.Echo,
	tableUsecase domain.TableUsecaseInterface,
) {
	handler := &tableHandler{
		tableUsecase: tableUsecase,
	}

	api := e.Group("/api")
	v1 := api.Group("/v1")

	tableGroup := v1.Group("/restaurants/:restaurant_id/branchs/:branch_id/tables")
	tableGroup.POST("/initialize", handler.InitializeTables)
}

// @Summary Get a list of chairs
// @Description Get a list of chairs in the restaurant
// @Tags Chairs
// @Produce json
// @Success 200 {array} Chair
// @Router /chairs [get]
func (h *tableHandler) InitializeTables(c echo.Context) error {
	restaurantIdParam := c.Param("restaurant_id")
	restaurantId, err := strconv.ParseUint(restaurantIdParam, 10, 64)
	if err != nil {
		logx.GetLog().Errorf("InitializeTables Handler Error: %s", err.Error())
		return c.JSON(http.StatusBadRequest, dto.CoreResponse{
			Message: enum.MessageError,
			Result:  err.Error(),
		})
	}

	branchIdParam := c.Param("branch_id")
	branchId, err := strconv.ParseUint(branchIdParam, 10, 64)
	if err != nil {
		logx.GetLog().Errorf("InitializeTables Handler Error: %s", err.Error())
		return c.JSON(http.StatusBadRequest, dto.CoreResponse{
			Message: enum.MessageError,
			Result:  err.Error(),
		})
	}

	tableCount, err := h.tableUsecase.InitializeTables(models.Branch{
		ID:           branchId,
		RestaurantId: restaurantId,
	})
	if err != nil {
		if err.Error() == enum.ErrorTableInitialized {
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

	result := dto.TableInitializeResponse{}
	result.ToTableInitializeResponse(&tableCount.NumberOfRemainingTables)

	res := dto.CoreResponse{
		Message: enum.MessageSuccess,
		Result:  result,
	}

	return c.JSON(http.StatusOK, res)
}
