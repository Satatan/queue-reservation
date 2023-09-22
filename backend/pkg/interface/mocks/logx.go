// Code generated by MockGen. DO NOT EDIT.
// Source: ./logx.go

// Package mock_util is a generated GoMock package.
package mock_util

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	logrus "github.com/sirupsen/logrus"
)

// MockLogxInterface is a mock of LogxInterface interface.
type MockLogxInterface struct {
	ctrl     *gomock.Controller
	recorder *MockLogxInterfaceMockRecorder
}

// MockLogxInterfaceMockRecorder is the mock recorder for MockLogxInterface.
type MockLogxInterfaceMockRecorder struct {
	mock *MockLogxInterface
}

// NewMockLogxInterface creates a new mock instance.
func NewMockLogxInterface(ctrl *gomock.Controller) *MockLogxInterface {
	mock := &MockLogxInterface{ctrl: ctrl}
	mock.recorder = &MockLogxInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogxInterface) EXPECT() *MockLogxInterfaceMockRecorder {
	return m.recorder
}

// GetLog mocks base method.
func (m *MockLogxInterface) GetLog() *logrus.Entry {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLog")
	ret0, _ := ret[0].(*logrus.Entry)
	return ret0
}

// GetLog indicates an expected call of GetLog.
func (mr *MockLogxInterfaceMockRecorder) GetLog() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLog", reflect.TypeOf((*MockLogxInterface)(nil).GetLog))
}
