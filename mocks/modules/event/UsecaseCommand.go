// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	request "event-service/internal/modules/event/models/request"
)

// UsecaseCommand is an autogenerated mock type for the UsecaseCommand type
type UsecaseCommand struct {
	mock.Mock
}

// CreateEvent provides a mock function with given fields: origCtx, payload
func (_m *UsecaseCommand) CreateEvent(origCtx context.Context, payload request.EventReq) (*string, error) {
	ret := _m.Called(origCtx, payload)

	if len(ret) == 0 {
		panic("no return value specified for CreateEvent")
	}

	var r0 *string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, request.EventReq) (*string, error)); ok {
		return rf(origCtx, payload)
	}
	if rf, ok := ret.Get(0).(func(context.Context, request.EventReq) *string); ok {
		r0 = rf(origCtx, payload)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, request.EventReq) error); ok {
		r1 = rf(origCtx, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateOnlineTicketConfig provides a mock function with given fields: origCtx, payload
func (_m *UsecaseCommand) CreateOnlineTicketConfig(origCtx context.Context, payload request.OnlineTicketReq) (*string, error) {
	ret := _m.Called(origCtx, payload)

	if len(ret) == 0 {
		panic("no return value specified for CreateOnlineTicketConfig")
	}

	var r0 *string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, request.OnlineTicketReq) (*string, error)); ok {
		return rf(origCtx, payload)
	}
	if rf, ok := ret.Get(0).(func(context.Context, request.OnlineTicketReq) *string); ok {
		r0 = rf(origCtx, payload)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, request.OnlineTicketReq) error); ok {
		r1 = rf(origCtx, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUsecaseCommand creates a new instance of UsecaseCommand. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUsecaseCommand(t interface {
	mock.TestingT
	Cleanup(func())
}) *UsecaseCommand {
	mock := &UsecaseCommand{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
