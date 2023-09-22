package delivery

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"queue_reservation/internal/delivery/dto"
	domainMocks "queue_reservation/internal/domain/mocks"
	"queue_reservation/internal/models"
	"queue_reservation/internal/models/enum"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewBookingHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		e                      *echo.Echo
		bookingUsecase         *domainMocks.MockBookingUsecaseInterface
		bookingUsecaseBehavior func(*domainMocks.MockBookingUsecaseInterface)
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				e:                      echo.New(),
				bookingUsecase:         domainMocks.NewMockBookingUsecaseInterface(ctrl),
				bookingUsecaseBehavior: func(meui *domainMocks.MockBookingUsecaseInterface) {},
			},
		},
	}

	for _, tt := range tests {
		tt.args.bookingUsecaseBehavior(tt.args.bookingUsecase)

		t.Run(tt.name, func(t *testing.T) {
			NewBookingHandler(tt.args.e, tt.args.bookingUsecase)
		})
	}
}

func Test_bookingHandler_CreateBooking(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		id                      uint64 = 1
		errorTableNotInitialize        = errors.New(enum.ErrorTableNotInitialize)
		fooErr                         = errors.New("foo")
	)

	type fields struct {
		bookingUsecase         *domainMocks.MockBookingUsecaseInterface
		bookingUsecaseBehavior func(*domainMocks.MockBookingUsecaseInterface)
	}
	type httpBuilder struct {
		restaurantId string
		branchId     string
		body         interface{}
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name        string
		fields      fields
		httpBuilder httpBuilder
		args        args
		statusCode  int
	}{
		{
			name: "error with binding",
			fields: fields{
				bookingUsecase:         domainMocks.NewMockBookingUsecaseInterface(ctrl),
				bookingUsecaseBehavior: func(meui *domainMocks.MockBookingUsecaseInterface) {},
			},
			httpBuilder: httpBuilder{
				restaurantId: "1",
				branchId:     "1",
				body:         "foo",
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "error with restaurant_id is not number",
			fields: fields{
				bookingUsecase:         domainMocks.NewMockBookingUsecaseInterface(ctrl),
				bookingUsecaseBehavior: func(meui *domainMocks.MockBookingUsecaseInterface) {},
			},
			httpBuilder: httpBuilder{
				restaurantId: "foo",
				branchId:     "1",
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "error with branch_id is not number",
			fields: fields{
				bookingUsecase:         domainMocks.NewMockBookingUsecaseInterface(ctrl),
				bookingUsecaseBehavior: func(meui *domainMocks.MockBookingUsecaseInterface) {},
			},
			httpBuilder: httpBuilder{
				restaurantId: "1",
				branchId:     "foo",
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "error validate",
			fields: fields{
				bookingUsecase:         domainMocks.NewMockBookingUsecaseInterface(ctrl),
				bookingUsecaseBehavior: func(meui *domainMocks.MockBookingUsecaseInterface) {},
			},
			httpBuilder: httpBuilder{
				restaurantId: "1",
				branchId:     "1",
				body:         dto.BookingRequest{},
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "error with not initialize",
			fields: fields{
				bookingUsecase: domainMocks.NewMockBookingUsecaseInterface(ctrl),
				bookingUsecaseBehavior: func(meui *domainMocks.MockBookingUsecaseInterface) {
					meui.EXPECT().CreateBooking(gomock.Any(), gomock.Any()).Return(
						nil, nil, errorTableNotInitialize,
					)
				},
			},
			httpBuilder: httpBuilder{
				restaurantId: "1",
				branchId:     "1",
				body: dto.BookingRequest{
					NumberOfCustomers: 4,
				},
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "error with initialization",
			fields: fields{
				bookingUsecase: domainMocks.NewMockBookingUsecaseInterface(ctrl),
				bookingUsecaseBehavior: func(meui *domainMocks.MockBookingUsecaseInterface) {
					meui.EXPECT().CreateBooking(gomock.Any(), gomock.Any()).Return(
						nil, nil, fooErr,
					)
				},
			},
			httpBuilder: httpBuilder{
				restaurantId: "1",
				branchId:     "1",
				body: dto.BookingRequest{
					NumberOfCustomers: 4,
				},
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "success",
			fields: fields{
				bookingUsecase: domainMocks.NewMockBookingUsecaseInterface(ctrl),
				bookingUsecaseBehavior: func(meui *domainMocks.MockBookingUsecaseInterface) {
					meui.EXPECT().CreateBooking(gomock.Any(), gomock.Any()).Return(
						&id,
						&models.TableCount{
							NumberOfNewReservedTables: 1,
							NumberOfRemainingTables:   1,
						},
						nil,
					)
				},
			},
			httpBuilder: httpBuilder{
				restaurantId: "1",
				branchId:     "1",
				body: dto.BookingRequest{
					NumberOfCustomers: 4,
				},
			},
			statusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {

		tt.fields.bookingUsecaseBehavior(tt.fields.bookingUsecase)

		e := echo.New()
		b, err := json.Marshal(tt.httpBuilder.body)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/restaurants/:restaurant_id/branchs/:branch_id/bookings", strings.NewReader(string(b)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		tt.args.c = e.NewContext(req, rec)

		tt.args.c.SetParamNames("restaurant_id", "branch_id")
		tt.args.c.SetParamValues(tt.httpBuilder.restaurantId, tt.httpBuilder.branchId)

		t.Run(tt.name, func(t *testing.T) {
			h := bookingHandler{
				bookingUsecase: tt.fields.bookingUsecase,
			}

			err := h.CreateBooking(tt.args.c)
			if err != nil {
				httpError := err.(*echo.HTTPError)
				assert.Equal(t, tt.statusCode, httpError.Code)
			} else {
				assert.Equal(t, tt.statusCode, rec.Code)
			}
		})
	}
}

func Test_bookingHandler_CancelBooking(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		fooErr              = errors.New("foo")
		errorRecordNotFound = errors.New(enum.ErrorRecordNotFound)
	)

	type fields struct {
		bookingUsecase         *domainMocks.MockBookingUsecaseInterface
		bookingUsecaseBehavior func(*domainMocks.MockBookingUsecaseInterface)
	}
	type httpBuilder struct {
		restaurantId string
		branchId     string
		bookingId    string
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name        string
		fields      fields
		httpBuilder httpBuilder
		args        args
		statusCode  int
	}{
		{
			name: "error with restaurant_id is not number",
			fields: fields{
				bookingUsecase:         domainMocks.NewMockBookingUsecaseInterface(ctrl),
				bookingUsecaseBehavior: func(meui *domainMocks.MockBookingUsecaseInterface) {},
			},
			httpBuilder: httpBuilder{
				restaurantId: "foo",
				branchId:     "1",
				bookingId:    "1",
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "error with branch_id is not number",
			fields: fields{
				bookingUsecase:         domainMocks.NewMockBookingUsecaseInterface(ctrl),
				bookingUsecaseBehavior: func(meui *domainMocks.MockBookingUsecaseInterface) {},
			},
			httpBuilder: httpBuilder{
				restaurantId: "1",
				branchId:     "foo",
				bookingId:    "1",
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "error with booking_id is not number",
			fields: fields{
				bookingUsecase:         domainMocks.NewMockBookingUsecaseInterface(ctrl),
				bookingUsecaseBehavior: func(meui *domainMocks.MockBookingUsecaseInterface) {},
			},
			httpBuilder: httpBuilder{
				restaurantId: "1",
				branchId:     "1",
				bookingId:    "foo",
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "error with record not found",
			fields: fields{
				bookingUsecase: domainMocks.NewMockBookingUsecaseInterface(ctrl),
				bookingUsecaseBehavior: func(meui *domainMocks.MockBookingUsecaseInterface) {
					meui.EXPECT().CancelBooking(gomock.Any(), gomock.Any()).Return(
						nil, errorRecordNotFound,
					)
				},
			},
			httpBuilder: httpBuilder{
				restaurantId: "1",
				branchId:     "1",
				bookingId:    "1",
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "error with canceling",
			fields: fields{
				bookingUsecase: domainMocks.NewMockBookingUsecaseInterface(ctrl),
				bookingUsecaseBehavior: func(meui *domainMocks.MockBookingUsecaseInterface) {
					meui.EXPECT().CancelBooking(gomock.Any(), gomock.Any()).Return(
						nil, fooErr,
					)
				},
			},
			httpBuilder: httpBuilder{
				restaurantId: "1",
				branchId:     "1",
				bookingId:    "1",
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "success",
			fields: fields{
				bookingUsecase: domainMocks.NewMockBookingUsecaseInterface(ctrl),
				bookingUsecaseBehavior: func(meui *domainMocks.MockBookingUsecaseInterface) {
					meui.EXPECT().CancelBooking(gomock.Any(), gomock.Any()).Return(
						&models.TableCount{
							NumberOfRemainingTables: 1,
						},
						nil,
					)
				},
			},
			httpBuilder: httpBuilder{
				restaurantId: "1",
				branchId:     "1",
				bookingId:    "1",
			},
			statusCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		tt.fields.bookingUsecaseBehavior(tt.fields.bookingUsecase)

		e := echo.New()

		req := httptest.NewRequest(http.MethodDelete, "/restaurants/:restaurant_id/branchs/:branch_id/bookings/booking_id", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		tt.args.c = e.NewContext(req, rec)

		tt.args.c.SetParamNames("restaurant_id", "branch_id", "booking_id")
		tt.args.c.SetParamValues(tt.httpBuilder.restaurantId, tt.httpBuilder.branchId, tt.httpBuilder.bookingId)

		t.Run(tt.name, func(t *testing.T) {
			h := bookingHandler{
				bookingUsecase: tt.fields.bookingUsecase,
			}

			err := h.CancelBooking(tt.args.c)
			if err != nil {
				httpError := err.(*echo.HTTPError)
				assert.Equal(t, tt.statusCode, httpError.Code)
			} else {
				assert.Equal(t, tt.statusCode, rec.Code)
			}
		})
	}
}
