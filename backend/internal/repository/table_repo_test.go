package repository

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewTableRepository(t *testing.T) {
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
			got := NewTableRepository()
			assert.NotNil(t, got)
		})
	}
}
