// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package stores is a generated GoMock package.
package stores

import (
	reflect "reflect"

	models "developer.zopsmart.com/go/gofr/examples/data-layer-with-ycql/models"
	gofr "developer.zopsmart.com/go/gofr/pkg/gofr"
	gomock "github.com/golang/mock/gomock"
)

// MockShop is a mock of Shop interface.
type MockShop struct {
	ctrl     *gomock.Controller
	recorder *MockShopMockRecorder
}

// MockShopMockRecorder is the mock recorder for MockShop.
type MockShopMockRecorder struct {
	mock *MockShop
}

// NewMockShop creates a new mock instance.
func NewMockShop(ctrl *gomock.Controller) *MockShop {
	mock := &MockShop{ctrl: ctrl}
	mock.recorder = &MockShopMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockShop) EXPECT() *MockShopMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockShop) Create(ctx *gofr.Context, data models.Shop) ([]models.Shop, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, data)
	ret0, _ := ret[0].([]models.Shop)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockShopMockRecorder) Create(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockShop)(nil).Create), ctx, data)
}

// Delete mocks base method.
func (m *MockShop) Delete(ctx *gofr.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockShopMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockShop)(nil).Delete), ctx, id)
}

// Get mocks base method.
func (m *MockShop) Get(ctx *gofr.Context, filter models.Shop) []models.Shop {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, filter)
	ret0, _ := ret[0].([]models.Shop)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockShopMockRecorder) Get(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockShop)(nil).Get), ctx, filter)
}

// Update mocks base method.
func (m *MockShop) Update(ctx *gofr.Context, data models.Shop) ([]models.Shop, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, data)
	ret0, _ := ret[0].([]models.Shop)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockShopMockRecorder) Update(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockShop)(nil).Update), ctx, data)
}
