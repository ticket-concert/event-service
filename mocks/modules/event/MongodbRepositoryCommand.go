// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	context "context"
	entity "event-service/internal/modules/event/models/entity"

	helpers "event-service/internal/pkg/helpers"

	mock "github.com/stretchr/testify/mock"
)

// MongodbRepositoryCommand is an autogenerated mock type for the MongodbRepositoryCommand type
type MongodbRepositoryCommand struct {
	mock.Mock
}

// InsertOneEventCollection provides a mock function with given fields: ctx, _a1
func (_m *MongodbRepositoryCommand) InsertOneEventCollection(ctx context.Context, _a1 entity.Event) <-chan helpers.Result {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for InsertOneEventCollection")
	}

	var r0 <-chan helpers.Result
	if rf, ok := ret.Get(0).(func(context.Context, entity.Event) <-chan helpers.Result); ok {
		r0 = rf(ctx, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan helpers.Result)
		}
	}

	return r0
}

// NewMongodbRepositoryCommand creates a new instance of MongodbRepositoryCommand. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMongodbRepositoryCommand(t interface {
	mock.TestingT
	Cleanup(func())
}) *MongodbRepositoryCommand {
	mock := &MongodbRepositoryCommand{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
