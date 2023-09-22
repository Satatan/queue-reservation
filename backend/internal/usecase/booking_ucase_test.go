package usecase

import (
	"errors"
	domainMocks "queue_reservation/internal/domain/mocks"
	"queue_reservation/internal/models"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewBookingUsecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		bookingRepo         *domainMocks.MockBookingRepositoryInterface
		bookingRepoBehavior func(*domainMocks.MockBookingRepositoryInterface)
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				bookingRepo:         domainMocks.NewMockBookingRepositoryInterface(ctrl),
				bookingRepoBehavior: func(meri *domainMocks.MockBookingRepositoryInterface) {},
			},
		},
	}
	for _, tt := range tests {

		tt.args.bookingRepoBehavior(tt.args.bookingRepo)

		t.Run(tt.name, func(t *testing.T) {
			got := NewBookingUsecase(tt.args.bookingRepo)
			assert.NotNil(t, got)
		})
	}
}

func Test_bookingUsecase_CreateBooking(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		id     uint64 = 1
		fooErr        = errors.New("foo")
	)

	type fields struct {
		bookingRepo         *domainMocks.MockBookingRepositoryInterface
		bookingRepoBehavior func(*domainMocks.MockBookingRepositoryInterface)
	}
	type args struct {
		branch  models.Branch
		booking models.Booking
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "error with creating",
			fields: fields{
				bookingRepo: domainMocks.NewMockBookingRepositoryInterface(ctrl),
				bookingRepoBehavior: func(meri *domainMocks.MockBookingRepositoryInterface) {
					meri.EXPECT().CreateBooking(gomock.Any(), gomock.Any()).Return(
						nil, nil, fooErr,
					)
				},
			},
			args: args{
				branch: models.Branch{
					ID:           1,
					RestaurantId: 1,
				},
				booking: models.Booking{
					Number: 4,
				},
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				bookingRepo: domainMocks.NewMockBookingRepositoryInterface(ctrl),
				bookingRepoBehavior: func(meri *domainMocks.MockBookingRepositoryInterface) {
					meri.EXPECT().CreateBooking(gomock.Any(), gomock.Any()).Return(
						&id,
						&models.TableCount{
							NumberOfNewReservedTables: 1,
							NumberOfRemainingTables:   1,
						},
						nil,
					)
				},
			},
			args: args{
				branch: models.Branch{
					ID:           1,
					RestaurantId: 1,
				},
				booking: models.Booking{
					Number: 4,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt.fields.bookingRepoBehavior(tt.fields.bookingRepo)

		t.Run(tt.name, func(t *testing.T) {
			u := bookingUsecase{
				bookingRepo: tt.fields.bookingRepo,
			}

			got, got1, err := u.CreateBooking(tt.args.branch, tt.args.booking)
			if (err != nil) != tt.wantErr {
				t.Errorf("bookingUsecase.CreateBooking() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
				assert.NotNil(t, got1)
			} else {
				assert.Nil(t, got)
				assert.Nil(t, got1)
			}
		})
	}
}

func Test_bookingUsecase_CancelBooking(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		fooErr = errors.New("foo")
	)

	type fields struct {
		bookingRepo         *domainMocks.MockBookingRepositoryInterface
		bookingRepoBehavior func(*domainMocks.MockBookingRepositoryInterface)
	}
	type args struct {
		branch  models.Branch
		booking models.Booking
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "error with canceling",
			fields: fields{
				bookingRepo: domainMocks.NewMockBookingRepositoryInterface(ctrl),
				bookingRepoBehavior: func(meri *domainMocks.MockBookingRepositoryInterface) {
					meri.EXPECT().CancelBooking(gomock.Any(), gomock.Any()).Return(
						nil, fooErr,
					)
				},
			},
			args: args{
				branch: models.Branch{
					ID:           1,
					RestaurantId: 1,
				},
				booking: models.Booking{
					ID: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				bookingRepo: domainMocks.NewMockBookingRepositoryInterface(ctrl),
				bookingRepoBehavior: func(meri *domainMocks.MockBookingRepositoryInterface) {
					meri.EXPECT().CancelBooking(gomock.Any(), gomock.Any()).Return(
						&models.TableCount{
							NumberOfRemainingTables: 1,
						},
						nil,
					)
				},
			},
			args: args{
				branch: models.Branch{
					ID:           1,
					RestaurantId: 1,
				},
				booking: models.Booking{
					ID: 1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt.fields.bookingRepoBehavior(tt.fields.bookingRepo)

		t.Run(tt.name, func(t *testing.T) {
			u := bookingUsecase{
				bookingRepo: tt.fields.bookingRepo,
			}

			got, err := u.CancelBooking(tt.args.branch, tt.args.booking)
			if (err != nil) != tt.wantErr {
				t.Errorf("bookingUsecase.CancelBooking() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			} else {
				assert.Nil(t, got)
			}
		})
	}
}
