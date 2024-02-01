// Code generated by MockGen. DO NOT EDIT.
// Source: services/loader.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockLoader is a mock of Loader interface.
type MockLoader struct {
	ctrl     *gomock.Controller
	recorder *MockLoaderMockRecorder
}

// MockLoaderMockRecorder is the mock recorder for MockLoader.
type MockLoaderMockRecorder struct {
	mock *MockLoader
}

// NewMockLoader creates a new mock instance.
func NewMockLoader(ctrl *gomock.Controller) *MockLoader {
	mock := &MockLoader{ctrl: ctrl}
	mock.recorder = &MockLoaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLoader) EXPECT() *MockLoaderMockRecorder {
	return m.recorder
}

// GetCpuUsage mocks base method.
func (m *MockLoader) GetCpuUsage() (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCpuUsage")
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCpuUsage indicates an expected call of GetCpuUsage.
func (mr *MockLoaderMockRecorder) GetCpuUsage() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCpuUsage", reflect.TypeOf((*MockLoader)(nil).GetCpuUsage))
}
