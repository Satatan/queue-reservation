package delivery

import (
	"errors"
	"net/http"
	"net/http/httptest"
	domainMocks "queue_reservation/internal/domain/mocks"
	"queue_reservation/internal/models"
	"queue_reservation/internal/models/enum"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewTableHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		e                    *echo.Echo
		tableUsecase         *domainMocks.MockTableUsecaseInterface
		tableUsecaseBehavior func(*domainMocks.MockTableUsecaseInterface)
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				e:                    echo.New(),
				tableUsecase:         domainMocks.NewMockTableUsecaseInterface(ctrl),
				tableUsecaseBehavior: func(meui *domainMocks.MockTableUsecaseInterface) {},
			},
		},
	}
	for _, tt := range tests {
		tt.args.tableUsecaseBehavior(tt.args.tableUsecase)

		t.Run(tt.name, func(t *testing.T) {
			NewTableHandler(tt.args.e, tt.args.tableUsecase)
		})
	}
}

func Test_tableHandler_InitializeTables(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		fooErr                = errors.New("foo")
		errorTableInitialized = errors.New(enum.ErrorTableInitialized)
	)

	type fields struct {
		tableUsecase         *domainMocks.MockTableUsecaseInterface
		tableUsecaseBehavior func(*domainMocks.MockTableUsecaseInterface)
	}
	type httpBuilder struct {
		restaurantId string
		branchId     string
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
				tableUsecase:         domainMocks.NewMockTableUsecaseInterface(ctrl),
				tableUsecaseBehavior: func(meui *domainMocks.MockTableUsecaseInterface) {},
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
				tableUsecase:         domainMocks.NewMockTableUsecaseInterface(ctrl),
				tableUsecaseBehavior: func(meui *domainMocks.MockTableUsecaseInterface) {},
			},
			httpBuilder: httpBuilder{
				restaurantId: "1",
				branchId:     "foo",
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "error already initialized",
			fields: fields{
				tableUsecase: domainMocks.NewMockTableUsecaseInterface(ctrl),
				tableUsecaseBehavior: func(meui *domainMocks.MockTableUsecaseInterface) {
					meui.EXPECT().InitializeTables(gomock.Any()).Return(nil, errorTableInitialized)
				},
			},
			httpBuilder: httpBuilder{
				restaurantId: "1",
				branchId:     "1",
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "error with initialize",
			fields: fields{
				tableUsecase: domainMocks.NewMockTableUsecaseInterface(ctrl),
				tableUsecaseBehavior: func(meui *domainMocks.MockTableUsecaseInterface) {
					meui.EXPECT().InitializeTables(gomock.Any()).Return(nil, fooErr)
				},
			},
			httpBuilder: httpBuilder{
				restaurantId: "1",
				branchId:     "1",
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "success",
			fields: fields{
				tableUsecase: domainMocks.NewMockTableUsecaseInterface(ctrl),
				tableUsecaseBehavior: func(meui *domainMocks.MockTableUsecaseInterface) {
					meui.EXPECT().InitializeTables(gomock.Any()).Return(&models.TableCount{
						NumberOfRemainingTables: 10,
					}, nil)
				},
			},
			httpBuilder: httpBuilder{
				restaurantId: "1",
				branchId:     "1",
			},
			statusCode: http.StatusOK,
		},
	}
	for _, tt := range tests {

		tt.fields.tableUsecaseBehavior(tt.fields.tableUsecase)

		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/restaurants/:restaurant_id/branchs:branch_id/tables", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		tt.args.c = e.NewContext(req, rec)

		tt.args.c.SetParamNames("restaurant_id", "branch_id")
		tt.args.c.SetParamValues(tt.httpBuilder.restaurantId, tt.httpBuilder.branchId)

		t.Run(tt.name, func(t *testing.T) {
			h := tableHandler{
				tableUsecase: tt.fields.tableUsecase,
			}

			err := h.InitializeTables(tt.args.c)
			if err != nil {
				httpError := err.(*echo.HTTPError)
				assert.Equal(t, tt.statusCode, httpError.Code)
			} else {
				assert.Equal(t, tt.statusCode, rec.Code)
			}
		})
	}
}
