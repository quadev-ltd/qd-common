// Code generated by MockGen. DO NOT EDIT.
// Source: token_signer.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	jwt "github.com/quadev-ltd/qd-common/pkg/jwt"
)

// MockTokenSignerer is a mock of TokenSignerer interface.
type MockTokenSignerer struct {
	ctrl     *gomock.Controller
	recorder *MockTokenSignererMockRecorder
}

// MockTokenSignererMockRecorder is the mock recorder for MockTokenSignerer.
type MockTokenSignererMockRecorder struct {
	mock *MockTokenSignerer
}

// NewMockTokenSignerer creates a new mock instance.
func NewMockTokenSignerer(ctrl *gomock.Controller) *MockTokenSignerer {
	mock := &MockTokenSignerer{ctrl: ctrl}
	mock.recorder = &MockTokenSignererMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenSignerer) EXPECT() *MockTokenSignererMockRecorder {
	return m.recorder
}

// SignToken mocks base method.
func (m *MockTokenSignerer) SignToken(email string, expiry time.Time, tokenType jwt.TokenType) (*string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignToken", email, expiry, tokenType)
	ret0, _ := ret[0].(*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignToken indicates an expected call of SignToken.
func (mr *MockTokenSignererMockRecorder) SignToken(email, expiry, tokenType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignToken", reflect.TypeOf((*MockTokenSignerer)(nil).SignToken), email, expiry, tokenType)
}