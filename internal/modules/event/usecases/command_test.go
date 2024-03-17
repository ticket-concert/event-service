package usecases_test

import (
	"context"
	"event-service/internal/modules/event"
	"event-service/internal/pkg/errors"
	"event-service/internal/pkg/helpers"
	"testing"

	addressEntity "event-service/internal/modules/address/models/entity"
	eventEntity "event-service/internal/modules/event/models/entity"
	"event-service/internal/modules/event/models/request"
	uc "event-service/internal/modules/event/usecases"
	mockcertAddress "event-service/mocks/modules/address"
	mockcert "event-service/mocks/modules/event"
	mockcertTicket "event-service/mocks/modules/ticket"
	mockkafka "event-service/mocks/pkg/kafka"
	mocklog "event-service/mocks/pkg/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CommandUsecaseTestSuite struct {
	suite.Suite
	mockEventRepositoryQuery    *mockcert.MongodbRepositoryQuery
	mockEventRepositoryCommand  *mockcert.MongodbRepositoryCommand
	mockTicketRepositoryCommand *mockcertTicket.MongodbRepositoryCommand
	mockAddressRepositoryQuery  *mockcertAddress.MongodbRepositoryQuery
	mockLogger                  *mocklog.Logger
	mockKafkaProducer           *mockkafka.Producer
	usecase                     event.UsecaseCommand
	ctx                         context.Context
}

func (suite *CommandUsecaseTestSuite) SetupTest() {
	suite.mockEventRepositoryQuery = &mockcert.MongodbRepositoryQuery{}
	suite.mockEventRepositoryCommand = &mockcert.MongodbRepositoryCommand{}
	suite.mockTicketRepositoryCommand = &mockcertTicket.MongodbRepositoryCommand{}
	suite.mockAddressRepositoryQuery = &mockcertAddress.MongodbRepositoryQuery{}
	suite.mockLogger = &mocklog.Logger{}
	suite.mockKafkaProducer = &mockkafka.Producer{}
	suite.ctx = context.Background()
	suite.usecase = uc.NewCommandUsecase(
		suite.mockEventRepositoryQuery,
		suite.mockEventRepositoryCommand,
		suite.mockTicketRepositoryCommand,
		suite.mockAddressRepositoryQuery,
		suite.mockKafkaProducer,
		suite.mockLogger,
	)
}

func TestCommandUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(CommandUsecaseTestSuite))
}

func (suite *CommandUsecaseTestSuite) TestCreateEvent() {
	payload := request.EventReq{
		UserId:        "userId",
		DateTime:      "2024-09-02 15:04",
		ContinentCode: "code",
		Tickets: []request.Ticket{
			{
				TicketType:  "Gold",
				TicketPrice: 50,
				TotalQuota:  10,
				Tag:         "tag",
			},
		},
	}

	mockEventByName := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockEventByTag := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag:       "tag",
			CreatedBy: "userId",
			UpdatedBy: "userId",
		},
		Error: nil,
	}

	mockContinentByCode := helpers.Result{
		Data: &addressEntity.Continent{
			Code: "code",
			Name: "name",
		},
		Error: nil,
	}

	mockCountry := helpers.Result{
		Data: &addressEntity.Country{
			Id:   1,
			Code: "code",
			Name: "name",
		},
		Error: nil,
	}

	mockInsertManyTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockEventRepositoryQuery.On("FindEventByName", mock.Anything, mock.Anything).Return(mockChannel(mockEventByName))
	suite.mockEventRepositoryQuery.On("FindEventByTag", mock.Anything, mock.Anything).Return(mockChannel(mockEventByTag))
	suite.mockAddressRepositoryQuery.On("FindOneContinentByCode", mock.Anything, mock.Anything).Return(mockChannel(mockContinentByCode))
	suite.mockAddressRepositoryQuery.On("FindOneCountry", mock.Anything, mock.Anything).Return(mockChannel(mockCountry))
	suite.mockTicketRepositoryCommand.On("InsertManyTicketCollection", mock.Anything, mock.Anything).Return(mockChannel(mockInsertManyTicket))
	suite.mockEventRepositoryCommand.On("InsertOneEventCollection", mock.Anything, mock.Anything).Return(mockChannel(mockInsertManyTicket))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)

	_, err := suite.usecase.CreateEvent(suite.ctx, payload)
	assert.NoError(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateEventErrEventName() {
	payload := request.EventReq{
		UserId:        "userId",
		DateTime:      "2024-09-02 15:04",
		ContinentCode: "code",
		Tickets: []request.Ticket{
			{
				TicketType:  "Gold",
				TicketPrice: 50,
				TotalQuota:  10,
				Tag:         "tag",
			},
		},
	}

	mockEventByName := helpers.Result{
		Data:  nil,
		Error: errors.BadRequest("error"),
	}

	suite.mockEventRepositoryQuery.On("FindEventByName", mock.Anything, mock.Anything).Return(mockChannel(mockEventByName))

	_, err := suite.usecase.CreateEvent(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateEventExistEventName() {
	payload := request.EventReq{
		UserId:        "userId",
		DateTime:      "2024-09-02 15:04",
		ContinentCode: "code",
		Tickets: []request.Ticket{
			{
				TicketType:  "Gold",
				TicketPrice: 50,
				TotalQuota:  10,
				Tag:         "tag",
			},
		},
	}

	mockEventByName := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
		},
		Error: nil,
	}

	suite.mockEventRepositoryQuery.On("FindEventByName", mock.Anything, mock.Anything).Return(mockChannel(mockEventByName))

	_, err := suite.usecase.CreateEvent(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateEventErrEventTag() {
	payload := request.EventReq{
		UserId:        "userId",
		DateTime:      "2024-09-02 15:04",
		ContinentCode: "code",
		Tickets: []request.Ticket{
			{
				TicketType:  "Gold",
				TicketPrice: 50,
				TotalQuota:  10,
				Tag:         "tag",
			},
		},
	}

	mockEventByName := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockEventByTag := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag:       "tag",
			CreatedBy: "userId",
			UpdatedBy: "userId",
		},
		Error: errors.BadRequest("error"),
	}

	suite.mockEventRepositoryQuery.On("FindEventByName", mock.Anything, mock.Anything).Return(mockChannel(mockEventByName))
	suite.mockEventRepositoryQuery.On("FindEventByTag", mock.Anything, mock.Anything).Return(mockChannel(mockEventByTag))

	_, err := suite.usecase.CreateEvent(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateEventByTagErr() {
	payload := request.EventReq{
		UserId:        "userId",
		DateTime:      "2024-09-02 15:04",
		ContinentCode: "code",
		Tickets: []request.Ticket{
			{
				TicketType:  "Gold",
				TicketPrice: 50,
				TotalQuota:  10,
				Tag:         "tag",
			},
		},
	}

	mockEventByName := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockEventByTag := helpers.Result{
		Data: &eventEntity.Country{
			Name: "name",
		},
		Error: nil,
	}

	suite.mockEventRepositoryQuery.On("FindEventByName", mock.Anything, mock.Anything).Return(mockChannel(mockEventByName))
	suite.mockEventRepositoryQuery.On("FindEventByTag", mock.Anything, mock.Anything).Return(mockChannel(mockEventByTag))

	_, err := suite.usecase.CreateEvent(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateEventErrUserId() {
	payload := request.EventReq{
		UserId:        "userId",
		DateTime:      "2024-09-02 15:04",
		ContinentCode: "code",
		Tickets: []request.Ticket{
			{
				TicketType:  "Gold",
				TicketPrice: 50,
				TotalQuota:  10,
				Tag:         "tag",
			},
		},
	}

	mockEventByName := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockEventByTag := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag:       "tag",
			CreatedBy: "tesId",
			UpdatedBy: "tesId",
		},
		Error: nil,
	}

	suite.mockEventRepositoryQuery.On("FindEventByName", mock.Anything, mock.Anything).Return(mockChannel(mockEventByName))
	suite.mockEventRepositoryQuery.On("FindEventByTag", mock.Anything, mock.Anything).Return(mockChannel(mockEventByTag))

	_, err := suite.usecase.CreateEvent(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateEventErrDate() {
	payload := request.EventReq{
		UserId:        "userId",
		DateTime:      "09-02-2024T15:04",
		ContinentCode: "code",
		Tickets: []request.Ticket{
			{
				TicketType:  "Gold",
				TicketPrice: 50,
				TotalQuota:  10,
				Tag:         "tag",
			},
		},
	}

	mockEventByName := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockEventByTag := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag:       "tag",
			CreatedBy: "userId",
			UpdatedBy: "userId",
		},
		Error: nil,
	}

	suite.mockEventRepositoryQuery.On("FindEventByName", mock.Anything, mock.Anything).Return(mockChannel(mockEventByName))
	suite.mockEventRepositoryQuery.On("FindEventByTag", mock.Anything, mock.Anything).Return(mockChannel(mockEventByTag))

	_, err := suite.usecase.CreateEvent(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateEventErrContinent() {
	payload := request.EventReq{
		UserId:        "userId",
		DateTime:      "2024-09-02 15:04",
		ContinentCode: "code",
		Tickets: []request.Ticket{
			{
				TicketType:  "Gold",
				TicketPrice: 50,
				TotalQuota:  10,
				Tag:         "tag",
			},
		},
	}

	mockEventByName := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockEventByTag := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag:       "tag",
			CreatedBy: "userId",
			UpdatedBy: "userId",
		},
		Error: nil,
	}

	mockContinentByCode := helpers.Result{
		Data:  nil,
		Error: errors.BadRequest("error"),
	}
	suite.mockEventRepositoryQuery.On("FindEventByName", mock.Anything, mock.Anything).Return(mockChannel(mockEventByName))
	suite.mockEventRepositoryQuery.On("FindEventByTag", mock.Anything, mock.Anything).Return(mockChannel(mockEventByTag))
	suite.mockAddressRepositoryQuery.On("FindOneContinentByCode", mock.Anything, mock.Anything).Return(mockChannel(mockContinentByCode))

	_, err := suite.usecase.CreateEvent(suite.ctx, payload)
	assert.NoError(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateEventErrContinentNil() {
	payload := request.EventReq{
		UserId:        "userId",
		DateTime:      "2024-09-02 15:04",
		ContinentCode: "code",
		Tickets: []request.Ticket{
			{
				TicketType:  "Gold",
				TicketPrice: 50,
				TotalQuota:  10,
				Tag:         "tag",
			},
		},
	}

	mockEventByName := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockEventByTag := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag:       "tag",
			CreatedBy: "userId",
			UpdatedBy: "userId",
		},
		Error: nil,
	}

	mockContinentByCode := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockEventRepositoryQuery.On("FindEventByName", mock.Anything, mock.Anything).Return(mockChannel(mockEventByName))
	suite.mockEventRepositoryQuery.On("FindEventByTag", mock.Anything, mock.Anything).Return(mockChannel(mockEventByTag))
	suite.mockAddressRepositoryQuery.On("FindOneContinentByCode", mock.Anything, mock.Anything).Return(mockChannel(mockContinentByCode))

	_, err := suite.usecase.CreateEvent(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateEventErrContinentParse() {
	payload := request.EventReq{
		UserId:        "userId",
		DateTime:      "2024-09-02 15:04",
		ContinentCode: "code",
		Tickets: []request.Ticket{
			{
				TicketType:  "Gold",
				TicketPrice: 50,
				TotalQuota:  10,
				Tag:         "tag",
			},
		},
	}

	mockEventByName := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockEventByTag := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag:       "tag",
			CreatedBy: "userId",
			UpdatedBy: "userId",
		},
		Error: nil,
	}

	mockContinentByCode := helpers.Result{
		Data: &addressEntity.Country{
			Code: "code",
			Name: "name",
		},
		Error: nil,
	}
	suite.mockEventRepositoryQuery.On("FindEventByName", mock.Anything, mock.Anything).Return(mockChannel(mockEventByName))
	suite.mockEventRepositoryQuery.On("FindEventByTag", mock.Anything, mock.Anything).Return(mockChannel(mockEventByTag))
	suite.mockAddressRepositoryQuery.On("FindOneContinentByCode", mock.Anything, mock.Anything).Return(mockChannel(mockContinentByCode))

	_, err := suite.usecase.CreateEvent(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateEventErrCountry() {
	payload := request.EventReq{
		UserId:        "userId",
		DateTime:      "2024-09-02 15:04",
		ContinentCode: "code",
		Tickets: []request.Ticket{
			{
				TicketType:  "Gold",
				TicketPrice: 50,
				TotalQuota:  10,
				Tag:         "tag",
			},
		},
	}

	mockEventByName := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockEventByTag := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag:       "tag",
			CreatedBy: "userId",
			UpdatedBy: "userId",
		},
		Error: nil,
	}

	mockContinentByCode := helpers.Result{
		Data: &addressEntity.Continent{
			Code: "code",
			Name: "name",
		},
		Error: nil,
	}

	mockCountry := helpers.Result{
		Data:  nil,
		Error: errors.BadRequest("error"),
	}

	suite.mockEventRepositoryQuery.On("FindEventByName", mock.Anything, mock.Anything).Return(mockChannel(mockEventByName))
	suite.mockEventRepositoryQuery.On("FindEventByTag", mock.Anything, mock.Anything).Return(mockChannel(mockEventByTag))
	suite.mockAddressRepositoryQuery.On("FindOneContinentByCode", mock.Anything, mock.Anything).Return(mockChannel(mockContinentByCode))
	suite.mockAddressRepositoryQuery.On("FindOneCountry", mock.Anything, mock.Anything).Return(mockChannel(mockCountry))

	_, err := suite.usecase.CreateEvent(suite.ctx, payload)
	assert.NoError(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateEventErrCountryNil() {
	payload := request.EventReq{
		UserId:        "userId",
		DateTime:      "2024-09-02 15:04",
		ContinentCode: "code",
		Tickets: []request.Ticket{
			{
				TicketType:  "Gold",
				TicketPrice: 50,
				TotalQuota:  10,
				Tag:         "tag",
			},
		},
	}

	mockEventByName := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockEventByTag := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag:       "tag",
			CreatedBy: "userId",
			UpdatedBy: "userId",
		},
		Error: nil,
	}

	mockContinentByCode := helpers.Result{
		Data: &addressEntity.Continent{
			Code: "code",
			Name: "name",
		},
		Error: nil,
	}

	mockCountry := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockEventRepositoryQuery.On("FindEventByName", mock.Anything, mock.Anything).Return(mockChannel(mockEventByName))
	suite.mockEventRepositoryQuery.On("FindEventByTag", mock.Anything, mock.Anything).Return(mockChannel(mockEventByTag))
	suite.mockAddressRepositoryQuery.On("FindOneContinentByCode", mock.Anything, mock.Anything).Return(mockChannel(mockContinentByCode))
	suite.mockAddressRepositoryQuery.On("FindOneCountry", mock.Anything, mock.Anything).Return(mockChannel(mockCountry))

	_, err := suite.usecase.CreateEvent(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateEventErrCountryParse() {
	payload := request.EventReq{
		UserId:        "userId",
		DateTime:      "2024-09-02 15:04",
		ContinentCode: "code",
		Tickets: []request.Ticket{
			{
				TicketType:  "Gold",
				TicketPrice: 50,
				TotalQuota:  10,
				Tag:         "tag",
			},
		},
	}

	mockEventByName := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockEventByTag := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag:       "tag",
			CreatedBy: "userId",
			UpdatedBy: "userId",
		},
		Error: nil,
	}

	mockContinentByCode := helpers.Result{
		Data: &addressEntity.Continent{
			Code: "code",
			Name: "name",
		},
		Error: nil,
	}

	mockCountry := helpers.Result{
		Data: &addressEntity.Continent{
			Code: "code",
		},
		Error: nil,
	}

	suite.mockEventRepositoryQuery.On("FindEventByName", mock.Anything, mock.Anything).Return(mockChannel(mockEventByName))
	suite.mockEventRepositoryQuery.On("FindEventByTag", mock.Anything, mock.Anything).Return(mockChannel(mockEventByTag))
	suite.mockAddressRepositoryQuery.On("FindOneContinentByCode", mock.Anything, mock.Anything).Return(mockChannel(mockContinentByCode))
	suite.mockAddressRepositoryQuery.On("FindOneCountry", mock.Anything, mock.Anything).Return(mockChannel(mockCountry))

	_, err := suite.usecase.CreateEvent(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateEventErrInsertTicket() {
	payload := request.EventReq{
		UserId:        "userId",
		DateTime:      "2024-09-02 15:04",
		ContinentCode: "code",
		Tickets: []request.Ticket{
			{
				TicketType:  "Gold",
				TicketPrice: 50,
				TotalQuota:  10,
				Tag:         "tag",
			},
		},
	}

	mockEventByName := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockEventByTag := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag:       "tag",
			CreatedBy: "userId",
			UpdatedBy: "userId",
		},
		Error: nil,
	}

	mockContinentByCode := helpers.Result{
		Data: &addressEntity.Continent{
			Code: "code",
			Name: "name",
		},
		Error: nil,
	}

	mockCountry := helpers.Result{
		Data: &addressEntity.Country{
			Id:   1,
			Code: "code",
			Name: "name",
		},
		Error: nil,
	}

	mockInsertManyTicket := helpers.Result{
		Data:  nil,
		Error: errors.BadRequest("error"),
	}
	suite.mockEventRepositoryQuery.On("FindEventByName", mock.Anything, mock.Anything).Return(mockChannel(mockEventByName))
	suite.mockEventRepositoryQuery.On("FindEventByTag", mock.Anything, mock.Anything).Return(mockChannel(mockEventByTag))
	suite.mockAddressRepositoryQuery.On("FindOneContinentByCode", mock.Anything, mock.Anything).Return(mockChannel(mockContinentByCode))
	suite.mockAddressRepositoryQuery.On("FindOneCountry", mock.Anything, mock.Anything).Return(mockChannel(mockCountry))
	suite.mockTicketRepositoryCommand.On("InsertManyTicketCollection", mock.Anything, mock.Anything).Return(mockChannel(mockInsertManyTicket))
	suite.mockEventRepositoryCommand.On("InsertOneEventCollection", mock.Anything, mock.Anything).Return(mockChannel(mockInsertManyTicket))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)

	_, err := suite.usecase.CreateEvent(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateEventErrInsertEvent() {
	payload := request.EventReq{
		UserId:        "userId",
		DateTime:      "2024-09-02 15:04",
		ContinentCode: "code",
		Tickets: []request.Ticket{
			{
				TicketType:  "Gold",
				TicketPrice: 50,
				TotalQuota:  10,
				Tag:         "tag",
			},
		},
	}

	mockEventByName := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockEventByTag := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag:       "tag",
			CreatedBy: "userId",
			UpdatedBy: "userId",
		},
		Error: nil,
	}

	mockContinentByCode := helpers.Result{
		Data: &addressEntity.Continent{
			Code: "code",
			Name: "name",
		},
		Error: nil,
	}

	mockCountry := helpers.Result{
		Data: &addressEntity.Country{
			Id:   1,
			Code: "code",
			Name: "name",
		},
		Error: nil,
	}

	mockInsertManyTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockInsertEvent := helpers.Result{
		Data:  nil,
		Error: errors.BadRequest("error"),
	}
	suite.mockEventRepositoryQuery.On("FindEventByName", mock.Anything, mock.Anything).Return(mockChannel(mockEventByName))
	suite.mockEventRepositoryQuery.On("FindEventByTag", mock.Anything, mock.Anything).Return(mockChannel(mockEventByTag))
	suite.mockAddressRepositoryQuery.On("FindOneContinentByCode", mock.Anything, mock.Anything).Return(mockChannel(mockContinentByCode))
	suite.mockAddressRepositoryQuery.On("FindOneCountry", mock.Anything, mock.Anything).Return(mockChannel(mockCountry))
	suite.mockTicketRepositoryCommand.On("InsertManyTicketCollection", mock.Anything, mock.Anything).Return(mockChannel(mockInsertManyTicket))
	suite.mockEventRepositoryCommand.On("InsertOneEventCollection", mock.Anything, mock.Anything).Return(mockChannel(mockInsertEvent))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)

	_, err := suite.usecase.CreateEvent(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateOnlineTicketConfig() {
	payload := request.OnlineTicketReq{
		UserId:     "userId",
		Tag:        "tag",
		TotalQuota: 10,
		CountryList: []request.CountryList{
			{
				CountryNumber: 1,
				Percentage:    100,
			},
		},
	}

	mockEventByTag := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag:       "tag",
			CreatedBy: "userId",
			UpdatedBy: "userId",
		},
		Error: nil,
	}

	mockUpsertOnlineTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockEventRepositoryQuery.On("FindEventByTag", mock.Anything, mock.Anything).Return(mockChannel(mockEventByTag))
	suite.mockTicketRepositoryCommand.On("UpsertOneOnlineTicketConfig", mock.Anything, mock.Anything).Return(mockChannel(mockUpsertOnlineTicket))

	_, err := suite.usecase.CreateOnlineTicketConfig(suite.ctx, payload)
	assert.NoError(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateOnlineTicketConfigErr() {
	payload := request.OnlineTicketReq{
		UserId:     "userId",
		Tag:        "tag",
		TotalQuota: 10,
		CountryList: []request.CountryList{
			{
				CountryNumber: 1,
				Percentage:    100,
			},
		},
	}

	mockEventByTag := helpers.Result{
		Data:  nil,
		Error: errors.BadRequest("error"),
	}

	mockUpsertOnlineTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockEventRepositoryQuery.On("FindEventByTag", mock.Anything, mock.Anything).Return(mockChannel(mockEventByTag))
	suite.mockTicketRepositoryCommand.On("UpsertOneOnlineTicketConfig", mock.Anything, mock.Anything).Return(mockChannel(mockUpsertOnlineTicket))

	_, err := suite.usecase.CreateOnlineTicketConfig(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateOnlineTicketConfigErrParse() {
	payload := request.OnlineTicketReq{
		UserId:     "userId",
		Tag:        "tag",
		TotalQuota: 10,
		CountryList: []request.CountryList{
			{
				CountryNumber: 1,
				Percentage:    100,
			},
		},
	}

	mockEventByTag := helpers.Result{
		Data: &eventEntity.Country{
			Name: "name",
		},
		Error: nil,
	}

	mockUpsertOnlineTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockEventRepositoryQuery.On("FindEventByTag", mock.Anything, mock.Anything).Return(mockChannel(mockEventByTag))
	suite.mockTicketRepositoryCommand.On("UpsertOneOnlineTicketConfig", mock.Anything, mock.Anything).Return(mockChannel(mockUpsertOnlineTicket))

	_, err := suite.usecase.CreateOnlineTicketConfig(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateOnlineTicketConfigErrUser() {
	payload := request.OnlineTicketReq{
		UserId:     "tesId",
		Tag:        "tag",
		TotalQuota: 10,
		CountryList: []request.CountryList{
			{
				CountryNumber: 1,
				Percentage:    100,
			},
		},
	}

	mockEventByTag := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag:       "tag",
			CreatedBy: "userId",
			UpdatedBy: "userId",
		},
		Error: nil,
	}

	mockUpsertOnlineTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockEventRepositoryQuery.On("FindEventByTag", mock.Anything, mock.Anything).Return(mockChannel(mockEventByTag))
	suite.mockTicketRepositoryCommand.On("UpsertOneOnlineTicketConfig", mock.Anything, mock.Anything).Return(mockChannel(mockUpsertOnlineTicket))

	_, err := suite.usecase.CreateOnlineTicketConfig(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateOnlineTicketConfigErrPercentage() {
	payload := request.OnlineTicketReq{
		UserId:     "userId",
		Tag:        "tag",
		TotalQuota: 10,
		CountryList: []request.CountryList{
			{
				CountryNumber: 1,
				Percentage:    10,
			},
		},
	}

	mockEventByTag := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag:       "tag",
			CreatedBy: "userId",
			UpdatedBy: "userId",
		},
		Error: nil,
	}

	mockUpsertOnlineTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockEventRepositoryQuery.On("FindEventByTag", mock.Anything, mock.Anything).Return(mockChannel(mockEventByTag))
	suite.mockTicketRepositoryCommand.On("UpsertOneOnlineTicketConfig", mock.Anything, mock.Anything).Return(mockChannel(mockUpsertOnlineTicket))

	_, err := suite.usecase.CreateOnlineTicketConfig(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateOnlineTicketConfigErrUpsert() {
	payload := request.OnlineTicketReq{
		UserId:     "userId",
		Tag:        "tag",
		TotalQuota: 10,
		CountryList: []request.CountryList{
			{
				CountryNumber: 1,
				Percentage:    100,
			},
		},
	}

	mockEventByTag := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag:       "tag",
			CreatedBy: "userId",
			UpdatedBy: "userId",
		},
		Error: nil,
	}

	mockUpsertOnlineTicket := helpers.Result{
		Data:  nil,
		Error: errors.BadRequest("error"),
	}
	suite.mockEventRepositoryQuery.On("FindEventByTag", mock.Anything, mock.Anything).Return(mockChannel(mockEventByTag))
	suite.mockTicketRepositoryCommand.On("UpsertOneOnlineTicketConfig", mock.Anything, mock.Anything).Return(mockChannel(mockUpsertOnlineTicket))

	_, err := suite.usecase.CreateOnlineTicketConfig(suite.ctx, payload)
	assert.NoError(suite.T(), err)
}
