// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	time "time"

	mock "github.com/stretchr/testify/mock"

	tokenhandler "github.com/victor-nach/time-tracker/lib/tokenhandler"
)

// TokenHandler is an autogenerated mock type for the TokenHandler type
type TokenHandler struct {
	mock.Mock
}

// NewToken provides a mock function with given fields: userId, expirationTime
func (_m *TokenHandler) NewToken(userId string, expirationTime time.Time) (string, error) {
	ret := _m.Called(userId, expirationTime)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, time.Time) string); ok {
		r0 = rf(userId, expirationTime)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, time.Time) error); ok {
		r1 = rf(userId, expirationTime)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ValidateToken provides a mock function with given fields: token
func (_m *TokenHandler) ValidateToken(token string) (*tokenhandler.Claims, error) {
	ret := _m.Called(token)

	var r0 *tokenhandler.Claims
	if rf, ok := ret.Get(0).(func(string) *tokenhandler.Claims); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tokenhandler.Claims)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}