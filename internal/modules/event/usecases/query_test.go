package usecases_test

import (
	"context"
	"event-service/internal/modules/event"
	"testing"

	"event-service/internal/modules/event/models/entity"
	"event-service/internal/modules/event/models/request"
	uc "event-service/internal/modules/event/usecases"
	"event-service/internal/pkg/errors"
	"event-service/internal/pkg/helpers"
	mockcert "event-service/mocks/modules/event"
	mocklog "event-service/mocks/pkg/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type QueryUsecaseTestSuite struct {
	suite.Suite
	mockOrderRepositoryQuery *mockcert.MongodbRepositoryQuery
	mockLogger               *mocklog.Logger
	usecase                  event.UsecaseQuery
	ctx                      context.Context
}

func (suite *QueryUsecaseTestSuite) SetupTest() {
	suite.mockOrderRepositoryQuery = &mockcert.MongodbRepositoryQuery{}
	suite.mockLogger = &mocklog.Logger{}
	suite.ctx = context.Background()
	suite.usecase = uc.NewQueryUsecase(
		suite.mockOrderRepositoryQuery,
		suite.mockLogger,
	)
}
func TestQueryUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(QueryUsecaseTestSuite))
}

func (suite *QueryUsecaseTestSuite) TestFindEvents() {
	payload := request.AllEventReq{
		Page: 1,
		Size: 1,
	}
	mockAllEvent := helpers.Result{
		Data: &[]entity.Event{
			{
				EventId: "id",
				Name:    "name",
			},
		},
		Error: nil,
	}

	suite.mockOrderRepositoryQuery.On("FindAllEvent", mock.Anything, payload).Return(mockChannel(mockAllEvent))

	_, err := suite.usecase.FindEvents(suite.ctx, payload)
	assert.NoError(suite.T(), err)
}

func (suite *QueryUsecaseTestSuite) TestFindEventsErr() {
	payload := request.AllEventReq{
		Page: 1,
		Size: 1,
	}
	mockAllEvent := helpers.Result{
		Data: &[]entity.Event{
			{
				EventId: "id",
				Name:    "name",
			},
		},
		Error: errors.BadRequest("error"),
	}

	suite.mockOrderRepositoryQuery.On("FindAllEvent", mock.Anything, payload).Return(mockChannel(mockAllEvent))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	_, err := suite.usecase.FindEvents(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *QueryUsecaseTestSuite) TestFindEventsErrNil() {
	payload := request.AllEventReq{
		Page: 1,
		Size: 1,
	}
	mockAllEvent := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockOrderRepositoryQuery.On("FindAllEvent", mock.Anything, payload).Return(mockChannel(mockAllEvent))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	_, err := suite.usecase.FindEvents(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *QueryUsecaseTestSuite) TestFindEventsErrParse() {
	payload := request.AllEventReq{
		Page: 1,
		Size: 1,
	}
	mockAllEvent := helpers.Result{
		Data: &entity.Country{
			Name: "name",
		},
		Error: nil,
	}

	suite.mockOrderRepositoryQuery.On("FindAllEvent", mock.Anything, payload).Return(mockChannel(mockAllEvent))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	_, err := suite.usecase.FindEvents(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

// Helper function to create a channel
func mockChannel(result helpers.Result) <-chan helpers.Result {
	responseChan := make(chan helpers.Result)

	go func() {
		responseChan <- result
		close(responseChan)
	}()

	return responseChan
}
