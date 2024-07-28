// Code generated by MockGen. DO NOT EDIT.
// Source: contract.go

// Package mock_handlers is a generated GoMock package.
package handlers

import (
	context "context"
	model "gophermart/internal/model"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockgmService is a mock of gmService interface.
type MockgmService struct {
	ctrl     *gomock.Controller
	recorder *MockgmServiceMockRecorder
}

// MockgmServiceMockRecorder is the mock recorder for MockgmService.
type MockgmServiceMockRecorder struct {
	mock *MockgmService
}

// NewMockgmService creates a new mock instance.
func NewMockgmService(ctrl *gomock.Controller) *MockgmService {
	mock := &MockgmService{ctrl: ctrl}
	mock.recorder = &MockgmServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockgmService) EXPECT() *MockgmServiceMockRecorder {
	return m.recorder
}

// AddAuthInfo mocks base method.
func (m *MockgmService) AddAuthInfo(ctx context.Context, login, pass string, passKey []byte) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddAuthInfo", ctx, login, pass, passKey)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddAuthInfo indicates an expected call of AddAuthInfo.
func (mr *MockgmServiceMockRecorder) AddAuthInfo(ctx, login, pass, passKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddAuthInfo", reflect.TypeOf((*MockgmService)(nil).AddAuthInfo), ctx, login, pass, passKey)
}

// AddOrder mocks base method.
func (m *MockgmService) AddOrder(ctx context.Context, orderID string, userID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddOrder", ctx, orderID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddOrder indicates an expected call of AddOrder.
func (mr *MockgmServiceMockRecorder) AddOrder(ctx, orderID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOrder", reflect.TypeOf((*MockgmService)(nil).AddOrder), ctx, orderID, userID)
}

// GetAuthInfo mocks base method.
func (m *MockgmService) GetAuthInfo(ctx context.Context, login, pass string, passKey []byte) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAuthInfo", ctx, login, pass, passKey)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAuthInfo indicates an expected call of GetAuthInfo.
func (mr *MockgmServiceMockRecorder) GetAuthInfo(ctx, login, pass, passKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAuthInfo", reflect.TypeOf((*MockgmService)(nil).GetAuthInfo), ctx, login, pass, passKey)
}

// GetBalance mocks base method.
func (m *MockgmService) GetBalance(ctx context.Context, userID int64) (model.Balance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBalance", ctx, userID)
	ret0, _ := ret[0].(model.Balance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBalance indicates an expected call of GetBalance.
func (mr *MockgmServiceMockRecorder) GetBalance(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalance", reflect.TypeOf((*MockgmService)(nil).GetBalance), ctx, userID)
}

// GetOrders mocks base method.
func (m *MockgmService) GetOrders(ctx context.Context, userID int64) ([]model.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrders", ctx, userID)
	ret0, _ := ret[0].([]model.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrders indicates an expected call of GetOrders.
func (mr *MockgmServiceMockRecorder) GetOrders(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrders", reflect.TypeOf((*MockgmService)(nil).GetOrders), ctx, userID)
}

// GetWithdrawals mocks base method.
func (m *MockgmService) GetWithdrawals(ctx context.Context, userID int64) ([]model.Withdraw, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithdrawals", ctx, userID)
	ret0, _ := ret[0].([]model.Withdraw)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWithdrawals indicates an expected call of GetWithdrawals.
func (mr *MockgmServiceMockRecorder) GetWithdrawals(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithdrawals", reflect.TypeOf((*MockgmService)(nil).GetWithdrawals), ctx, userID)
}

// Withdraw mocks base method.
func (m *MockgmService) Withdraw(ctx context.Context, withdraw model.Withdraw) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Withdraw", ctx, withdraw)
	ret0, _ := ret[0].(error)
	return ret0
}

// Withdraw indicates an expected call of Withdraw.
func (mr *MockgmServiceMockRecorder) Withdraw(ctx, withdraw interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Withdraw", reflect.TypeOf((*MockgmService)(nil).Withdraw), ctx, withdraw)
}