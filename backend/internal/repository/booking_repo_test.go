package repository

import (
	"queue_reservation/internal/models"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewBookingRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			got := NewBookingRepository()
			assert.NotNil(t, got)
		})
	}
}

func Test_bookingRepository_CreateBooking(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		id uint64 = 1
	)

	type fields struct {
	}
	type args struct {
		branch  models.Branch
		booking models.Booking
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		want          *uint64
		want1         *models.TableCount
		wantErr       bool
		setInitialize bool
	}{
		{
			name:   "error called before initialization",
			fields: fields{},
			args: args{
				branch: models.Branch{
					ID:           1,
					RestaurantId: 1,
				},
				booking: models.Booking{
					Number: 4,
				},
			},
			want:          nil,
			want1:         nil,
			wantErr:       true,
			setInitialize: false,
		},
		{
			name:   "error not enough tables for all customers",
			fields: fields{},
			args: args{
				branch: models.Branch{
					ID:           1,
					RestaurantId: 1,
				},
				booking: models.Booking{
					Number: 400,
				},
			},
			want:          nil,
			want1:         nil,
			wantErr:       true,
			setInitialize: true,
		},
		{
			name:   "success",
			fields: fields{},
			args: args{
				branch: models.Branch{
					ID:           1,
					RestaurantId: 1,
				},
				booking: models.Booking{
					ID:     1,
					Number: 4,
				},
			},
			want: &id,
			want1: &models.TableCount{
				NumberOfNewReservedTables: 1,
				NumberOfRemainingTables:   9,
			},
			wantErr:       false,
			setInitialize: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bookingRepo := &bookingRepository{}
			tableRepo := &tableRepository{}

			if tt.setInitialize {
				tableRepo.InitializeTables(tt.args.branch)
			}
			got, got1, err := bookingRepo.CreateBooking(tt.args.branch, tt.args.booking)
			if (err != nil) != tt.wantErr {
				t.Errorf("bookingRepository.CreateBooking() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else {
				bookingRepo.CancelBooking(tt.args.branch, tt.args.booking)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bookingRepository.CreateBooking() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("bookingRepository.CreateBooking() got1 = %v, want %v", got1, tt.want1)
			}

		})
	}
}

func Test_bookingRepository_CancelBooking(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
	}
	type args struct {
		branch  models.Branch
		booking models.Booking
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		want          *models.TableCount
		wantErr       bool
		setInitialize bool
		hasRecord     bool
	}{
		{
			name:   "error called before initialization",
			fields: fields{},
			args: args{
				branch: models.Branch{
					ID:           1,
					RestaurantId: 1,
				},
				booking: models.Booking{
					ID: 10,
				},
			},
			want:          nil,
			wantErr:       true,
			setInitialize: false,
		},
		{
			name:   "error with record not found",
			fields: fields{},
			args: args{
				branch: models.Branch{
					ID:           1,
					RestaurantId: 1,
				},
				booking: models.Booking{
					ID: 10,
				},
			},
			want:          nil,
			wantErr:       true,
			setInitialize: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bookingRepo := &bookingRepository{}
			tableRepo := &tableRepository{}

			if tt.setInitialize {
				tableRepo.InitializeTables(tt.args.branch)
			}

			got, err := bookingRepo.CancelBooking(tt.args.branch, tt.args.booking)
			if (err != nil) != tt.wantErr {
				t.Errorf("bookingRepository.CancelBooking() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bookingRepository.CancelBooking() = %v, want %v", got, tt.want)
			}

		})
	}
}
