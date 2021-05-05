// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package mock_store is a generated GoMock package.
package store

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/zopsmart/gofr/examples/using-ycql/entity"
	gofr "github.com/zopsmart/gofr/pkg/gofr"
)

// MockShop is a mock of Shop interface
type MockShop struct {
	ctrl     *gomock.Controller
	recorder *MockShopMockRecorder
}

// MockShopMockRecorder is the mock recorder for MockShop
type MockShopMockRecorder struct {
	mock *MockShop
}

// NewMockShop creates a new mock instance
func NewMockShop(ctrl *gomock.Controller) *MockShop {
	mock := &MockShop{ctrl: ctrl}
	mock.recorder = &MockShopMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockShop) EXPECT() *MockShopMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockShop) Get(ctx *gofr.Context, filter entity.Shop) []entity.Shop {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, filter)
	ret0, _ := ret[0].([]entity.Shop)
	return ret0
}

// Get indicates an expected call of Get
func (mr *MockShopMockRecorder) Get(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockShop)(nil).Get), ctx, filter)
}

// Create mocks base method
func (m *MockShop) Create(ctx *gofr.Context, data entity.Shop) ([]entity.Shop, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, data)
	ret0, _ := ret[0].([]entity.Shop)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockShopMockRecorder) Create(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockShop)(nil).Create), ctx, data)
}

// Delete mocks base method
func (m *MockShop) Delete(ctx *gofr.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockShopMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockShop)(nil).Delete), ctx, id)
}

// Update mocks base method
func (m *MockShop) Update(ctx *gofr.Context, data entity.Shop) ([]entity.Shop, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, data)
	ret0, _ := ret[0].([]entity.Shop)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update
func (mr *MockShopMockRecorder) Update(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockShop)(nil).Update), ctx, data)
}
