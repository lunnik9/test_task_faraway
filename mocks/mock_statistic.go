// Code generated by MockGen. DO NOT EDIT.
// Source: repository/statistic.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	repository "test_task_faraway/repository"

	gomock "github.com/golang/mock/gomock"
)

// MockStatModule is a mock of StatModule interface.
type MockStatModule struct {
	ctrl     *gomock.Controller
	recorder *MockStatModuleMockRecorder
}

// MockStatModuleMockRecorder is the mock recorder for MockStatModule.
type MockStatModuleMockRecorder struct {
	mock *MockStatModule
}

// NewMockStatModule creates a new mock instance.
func NewMockStatModule(ctrl *gomock.Controller) *MockStatModule {
	mock := &MockStatModule{ctrl: ctrl}
	mock.recorder = &MockStatModuleMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStatModule) EXPECT() *MockStatModuleMockRecorder {
	return m.recorder
}

// SaveUserStat mocks base method.
func (m *MockStatModule) SaveUserStat(ctx context.Context, stat repository.UserStat) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveUserStat", ctx, stat)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveUserStat indicates an expected call of SaveUserStat.
func (mr *MockStatModuleMockRecorder) SaveUserStat(ctx, stat interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveUserStat", reflect.TypeOf((*MockStatModule)(nil).SaveUserStat), ctx, stat)
}