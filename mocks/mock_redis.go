// Code generated by MockGen. DO NOT EDIT.
// Source: repository/redis.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCache is a mock of Cache interface.
type MockCache struct {
	ctrl     *gomock.Controller
	recorder *MockCacheMockRecorder
}

// MockCacheMockRecorder is the mock recorder for MockCache.
type MockCacheMockRecorder struct {
	mock *MockCache
}

// NewMockCache creates a new mock instance.
func NewMockCache(ctrl *gomock.Controller) *MockCache {
	mock := &MockCache{ctrl: ctrl}
	mock.recorder = &MockCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCache) EXPECT() *MockCacheMockRecorder {
	return m.recorder
}

// GetChallenge mocks base method.
func (m *MockCache) GetChallenge(ctx context.Context, userID int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChallenge", ctx, userID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChallenge indicates an expected call of GetChallenge.
func (mr *MockCacheMockRecorder) GetChallenge(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChallenge", reflect.TypeOf((*MockCache)(nil).GetChallenge), ctx, userID)
}

// GetUserID mocks base method.
func (m *MockCache) GetUserID(ctx context.Context) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserID", ctx)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserID indicates an expected call of GetUserID.
func (mr *MockCacheMockRecorder) GetUserID(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserID", reflect.TypeOf((*MockCache)(nil).GetUserID), ctx)
}

// InsertChallenge mocks base method.
func (m *MockCache) InsertChallenge(ctx context.Context, challenge string, userID int64, ttl int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertChallenge", ctx, challenge, userID, ttl)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertChallenge indicates an expected call of InsertChallenge.
func (mr *MockCacheMockRecorder) InsertChallenge(ctx, challenge, userID, ttl interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertChallenge", reflect.TypeOf((*MockCache)(nil).InsertChallenge), ctx, challenge, userID, ttl)
}

// InsertSusUser mocks base method.
func (m *MockCache) InsertSusUser(ctx context.Context, userID int64, ttl int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertSusUser", ctx, userID, ttl)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertSusUser indicates an expected call of InsertSusUser.
func (mr *MockCacheMockRecorder) InsertSusUser(ctx, userID, ttl interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertSusUser", reflect.TypeOf((*MockCache)(nil).InsertSusUser), ctx, userID, ttl)
}

// PresentSusUser mocks base method.
func (m *MockCache) PresentSusUser(ctx context.Context, userID int64, ttl int) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PresentSusUser", ctx, userID, ttl)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PresentSusUser indicates an expected call of PresentSusUser.
func (mr *MockCacheMockRecorder) PresentSusUser(ctx, userID, ttl interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PresentSusUser", reflect.TypeOf((*MockCache)(nil).PresentSusUser), ctx, userID, ttl)
}
