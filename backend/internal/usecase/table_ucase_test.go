package usecase

import (
	"errors"
	domainMocks "queue_reservation/internal/domain/mocks"
	"queue_reservation/internal/models"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewTableUsecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		tableRepo         *domainMocks.MockTableRepositoryInterface
		tableRepoBehavior func(*domainMocks.MockTableRepositoryInterface)
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				tableRepo:         domainMocks.NewMockTableRepositoryInterface(ctrl),
				tableRepoBehavior: func(meri *domainMocks.MockTableRepositoryInterface) {},
			},
		},
	}
	for _, tt := range tests {

		tt.args.tableRepoBehavior(tt.args.tableRepo)

		t.Run(tt.name, func(t *testing.T) {
			got := NewTableUsecase(tt.args.tableRepo)
			assert.NotNil(t, got)
		})
	}
}

func Test_tableUsecase_InitializeTables(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		fooErr = errors.New("foo")
	)

	type fields struct {
		tableRepo         *domainMocks.MockTableRepositoryInterface
		tableRepoBehavior func(*domainMocks.MockTableRepositoryInterface)
	}
	type args struct {
		branch models.Branch
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "error with initialize",
			fields: fields{
				tableRepo: domainMocks.NewMockTableRepositoryInterface(ctrl),
				tableRepoBehavior: func(meri *domainMocks.MockTableRepositoryInterface) {
					meri.EXPECT().InitializeTables(gomock.Any()).Return(
						nil, fooErr,
					)
				},
			},
			args: args{
				branch: models.Branch{
					ID:           1,
					RestaurantId: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				tableRepo: domainMocks.NewMockTableRepositoryInterface(ctrl),
				tableRepoBehavior: func(meri *domainMocks.MockTableRepositoryInterface) {
					meri.EXPECT().InitializeTables(gomock.Any()).Return(
						&models.TableCount{
							NumberOfRemainingTables: 10,
						}, nil,
					)
				},
			},
			args: args{
				branch: models.Branch{
					ID:           1,
					RestaurantId: 1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt.fields.tableRepoBehavior(tt.fields.tableRepo)

		t.Run(tt.name, func(t *testing.T) {
			u := tableUsecase{
				tableRepo: tt.fields.tableRepo,
			}

			got, err := u.InitializeTables(tt.args.branch)
			if (err != nil) != tt.wantErr {
				t.Errorf("tableUsecase.InitializeTables() error = %v, wantErr %v", err, tt.wantErr)
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
