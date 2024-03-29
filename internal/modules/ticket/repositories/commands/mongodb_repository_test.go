package commands_test

import (
	"context"
	"event-service/internal/modules/ticket"
	"event-service/internal/modules/ticket/models/entity"
	mongoRC "event-service/internal/modules/ticket/repositories/commands"
	"event-service/internal/pkg/helpers"
	mocks "event-service/mocks/pkg/databases/mongodb"
	mocklog "event-service/mocks/pkg/log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CommandTestSuite struct {
	suite.Suite
	mockMongodb *mocks.Collections
	mockLogger  *mocklog.Logger
	repository  ticket.MongodbRepositoryCommand
	ctx         context.Context
}

func (suite *CommandTestSuite) SetupTest() {
	suite.mockMongodb = new(mocks.Collections)
	suite.mockLogger = &mocklog.Logger{}
	suite.repository = mongoRC.NewCommandMongodbRepository(
		suite.mockMongodb,
		suite.mockLogger,
	)
	suite.ctx = context.Background()
}

func TestCommandTestSuite(t *testing.T) {
	suite.Run(t, new(CommandTestSuite))
}

func (suite *CommandTestSuite) TestInsertManyTicketCollection() {
	payload := []entity.Ticket{
		{
			TicketId: "id",
			EventId:  "id",
		},
	}

	// Mock UpsertOne
	expectedResult := make(chan helpers.Result)
	suite.mockMongodb.On("InsertMany", mock.Anything, mock.Anything).Return((<-chan helpers.Result)(expectedResult))

	// Act
	result := suite.repository.InsertManyTicketCollection(suite.ctx, payload)
	// Asset
	assert.NotNil(suite.T(), result, "Expected a result")

	// Simulate receiving a result from the channel
	go func() {
		expectedResult <- helpers.Result{Data: "result not nil", Error: nil}
		close(expectedResult)
	}()

	// Wait for the goroutine to complete
	<-result

	// Assert UpsertOne
	suite.mockMongodb.AssertCalled(suite.T(), "InsertMany", mock.Anything, mock.Anything)
}

func (suite *CommandTestSuite) TestUpsertOneOnlineTicketConfig() {
	payload := entity.OnlineTicketConfig{
		Tag:        "tag",
		TotalQuota: 10,
	}

	// Mock UpsertOne
	expectedResult := make(chan helpers.Result)
	suite.mockMongodb.On("UpsertOne", mock.Anything, mock.Anything).Return((<-chan helpers.Result)(expectedResult))

	// Act
	result := suite.repository.UpsertOneOnlineTicketConfig(suite.ctx, payload)
	// Asset
	assert.NotNil(suite.T(), result, "Expected a result")

	// Simulate receiving a result from the channel
	go func() {
		expectedResult <- helpers.Result{Data: "result not nil", Error: nil}
		close(expectedResult)
	}()

	// Wait for the goroutine to complete
	<-result

	// Assert UpsertOne
	suite.mockMongodb.AssertCalled(suite.T(), "UpsertOne", mock.Anything, mock.Anything)
}
