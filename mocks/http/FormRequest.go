// Code generated by mockery v2.34.2. DO NOT EDIT.

package mocks

import (
	http "github.com/goravel/framework/contracts/http"
	mock "github.com/stretchr/testify/mock"

	validation "github.com/goravel/framework/contracts/validation"
)

// FormRequest is an autogenerated mock type for the FormRequest type
type FormRequest struct {
	mock.Mock
}

// Attributes provides a mock function with given fields: ctx
func (_m *FormRequest) Attributes(ctx http.Context) map[string]string {
	ret := _m.Called(ctx)

	var r0 map[string]string
	if rf, ok := ret.Get(0).(func(http.Context) map[string]string); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	return r0
}

// Authorize provides a mock function with given fields: ctx
func (_m *FormRequest) Authorize(ctx http.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(http.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Messages provides a mock function with given fields: ctx
func (_m *FormRequest) Messages(ctx http.Context) map[string]string {
	ret := _m.Called(ctx)

	var r0 map[string]string
	if rf, ok := ret.Get(0).(func(http.Context) map[string]string); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	return r0
}

// PrepareForValidation provides a mock function with given fields: ctx, data
func (_m *FormRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	ret := _m.Called(ctx, data)

	var r0 error
	if rf, ok := ret.Get(0).(func(http.Context, validation.Data) error); ok {
		r0 = rf(ctx, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Rules provides a mock function with given fields: ctx
func (_m *FormRequest) Rules(ctx http.Context) map[string]string {
	ret := _m.Called(ctx)

	var r0 map[string]string
	if rf, ok := ret.Get(0).(func(http.Context) map[string]string); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	return r0
}

// NewFormRequest creates a new instance of FormRequest. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFormRequest(t interface {
	mock.TestingT
	Cleanup(func())
}) *FormRequest {
	mock := &FormRequest{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
