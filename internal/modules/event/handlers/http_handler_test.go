// user_http_handler_test.go

package handlers_test

import (
	"bytes"
	"encoding/json"
	"event-service/internal/modules/event/handlers"
	"event-service/internal/modules/event/models/request"
	"event-service/internal/modules/event/models/response"
	"event-service/internal/pkg/constants"
	"event-service/internal/pkg/errors"
	mockcert "event-service/mocks/modules/event"
	mocklog "event-service/mocks/pkg/log"
	mockredis "event-service/mocks/pkg/redis"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
)

type EventHttpHandlerTestSuite struct {
	suite.Suite

	cUC       *mockcert.UsecaseCommand
	cUQ       *mockcert.UsecaseQuery
	cLog      *mocklog.Logger
	validator *validator.Validate
	handler   *handlers.EventHttpHandler
	cRedis    *mockredis.Collections
	app       *fiber.App
}

func (suite *EventHttpHandlerTestSuite) SetupTest() {
	suite.cUC = new(mockcert.UsecaseCommand)
	suite.cUQ = new(mockcert.UsecaseQuery)
	suite.cLog = new(mocklog.Logger)
	suite.validator = validator.New()
	suite.cRedis = new(mockredis.Collections)
	suite.handler = &handlers.EventHttpHandler{
		EventUsecaseCommand: suite.cUC,
		EventUsecaseQuery:   suite.cUQ,
		Logger:              suite.cLog,
		Validator:           suite.validator,
	}
	suite.app = fiber.New()
	handlers.InitEventHttpHandler(suite.app, suite.cUC, suite.cUQ, suite.cLog, suite.cRedis)
}

func TestEventHttpHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(EventHttpHandlerTestSuite))
}

func (suite *EventHttpHandlerTestSuite) TestCreateEvent() {
	var res string
	suite.cUC.On("CreateEvent", mock.Anything, mock.Anything).Return(&res, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqM := request.EventReq{
		Name:          "name",
		DateTime:      "2024-03-09",
		ContinentName: "name",
		ContinentCode: "code",
		Country: request.Country{
			Name:  "name",
			Id:    1,
			City:  "city",
			Place: "place",
		},
		Description: "desc",
		Tag:         "tag",
		Tickets: []request.Ticket{
			{
				TicketType: "Gold",
			},
		},
	}
	requestBody, _ := json.Marshal(reqM)
	req := httptest.NewRequest(fiber.MethodPost, "/v1/create-event", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	ctx.Request().SetRequestURI("/v1/create-event")
	ctx.Request().Header.SetMethod(fiber.MethodPost)
	ctx.Request().Header.SetContentType("application/json")
	ctx.Request().SetBody(requestBody)

	err := suite.handler.CreateEvent(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *EventHttpHandlerTestSuite) TestCreateEventErrParser() {
	var res string
	suite.cUC.On("CreateEvent", mock.Anything, mock.Anything).Return(&res, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqM := request.EventReq{
		Name:          "name",
		DateTime:      "2024-03-09",
		ContinentName: "name",
		ContinentCode: "code",
		Country: request.Country{
			Name:  "name",
			Id:    1,
			City:  "city",
			Place: "place",
		},
		Description: "desc",
		Tag:         "tag",
		Tickets: []request.Ticket{
			{
				TicketType: "Gold",
			},
		},
	}
	requestBody, _ := json.Marshal(reqM)
	req := httptest.NewRequest(fiber.MethodPost, "/v1/create-event", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	ctx.Request().SetRequestURI("/v1/create-event")
	ctx.Request().Header.SetMethod(fiber.MethodPost)
	ctx.Request().Header.SetContentType("application/json")
	// ctx.Request().SetBody(requestBody)

	err := suite.handler.CreateEvent(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *EventHttpHandlerTestSuite) TestCreateEventErrValidation() {
	var res string
	suite.cUC.On("CreateEvent", mock.Anything, mock.Anything).Return(&res, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqM := request.EventReq{}
	requestBody, _ := json.Marshal(reqM)
	req := httptest.NewRequest(fiber.MethodPost, "/v1/create-event", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	ctx.Request().SetRequestURI("/v1/create-event")
	ctx.Request().Header.SetMethod(fiber.MethodPost)
	ctx.Request().Header.SetContentType("application/json")
	ctx.Request().SetBody(requestBody)

	err := suite.handler.CreateEvent(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *EventHttpHandlerTestSuite) TestCreateEventErr() {
	suite.cUC.On("CreateEvent", mock.Anything, mock.Anything).Return(nil, errors.BadRequest("error"))
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqM := request.EventReq{
		Name:          "name",
		DateTime:      "2024-03-09",
		ContinentName: "name",
		ContinentCode: "code",
		Country: request.Country{
			Name:  "name",
			Id:    1,
			City:  "city",
			Place: "place",
		},
		Description: "desc",
		Tag:         "tag",
		Tickets: []request.Ticket{
			{
				TicketType: "Gold",
			},
		},
	}
	requestBody, _ := json.Marshal(reqM)
	req := httptest.NewRequest(fiber.MethodPost, "/v1/create-event", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	ctx.Request().SetRequestURI("/v1/create-event")
	ctx.Request().Header.SetMethod(fiber.MethodPost)
	ctx.Request().Header.SetContentType("application/json")
	ctx.Request().SetBody(requestBody)

	err := suite.handler.CreateEvent(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *EventHttpHandlerTestSuite) TestGetEvents() {

	response := &response.EventResp{
		CollectionData: []response.Event{
			{
				EventId: "id",
			},
		},
		MetaData: constants.MetaData{},
	}
	suite.cUQ.On("FindEvents", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	req := httptest.NewRequest(fiber.MethodGet, "/v1/list?page=1&size=1", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/list?page=1&size=1")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetEvents(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *EventHttpHandlerTestSuite) TestGetEventsErrParse() {

	response := &response.EventResp{
		CollectionData: []response.Event{
			{
				EventId: "id",
			},
		},
		MetaData: constants.MetaData{},
	}
	suite.cUQ.On("FindEvents", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	req := httptest.NewRequest(fiber.MethodGet, "/v1/list?page=aa&size=aa", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/list?page=aa&size=aa")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetEvents(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *EventHttpHandlerTestSuite) TestGetEventsErrValidate() {

	response := &response.EventResp{
		CollectionData: []response.Event{
			{
				EventId: "id",
			},
		},
		MetaData: constants.MetaData{},
	}
	suite.cUQ.On("FindEvents", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	req := httptest.NewRequest(fiber.MethodGet, "/v1/list?page=&size=", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/list?page=&size=")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetEvents(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *EventHttpHandlerTestSuite) TestGetEventsErr() {
	suite.cUQ.On("FindEvents", mock.Anything, mock.Anything).Return(nil, errors.BadRequest("error"))
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	req := httptest.NewRequest(fiber.MethodGet, "/v1/list?page=1&size=1", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/list?page=1&size=1")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetEvents(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *EventHttpHandlerTestSuite) TestCreateOnlineTicketConfig() {
	var res string
	suite.cUC.On("CreateOnlineTicketConfig", mock.Anything, mock.Anything).Return(&res, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqM := request.OnlineTicketReq{
		UserId:     "id",
		Tag:        "tag",
		TotalQuota: 10,
		CountryList: []request.CountryList{
			{
				CountryNumber: 1,
				Percentage:    30,
			},
		},
	}
	requestBody, _ := json.Marshal(reqM)
	req := httptest.NewRequest(fiber.MethodPost, "/v1/create-online-ticket-config", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	ctx.Request().SetRequestURI("/v1/create-online-ticket-config")
	ctx.Request().Header.SetMethod(fiber.MethodPost)
	ctx.Request().Header.SetContentType("application/json")
	ctx.Request().SetBody(requestBody)

	err := suite.handler.CreateOnlineTicketConfig(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *EventHttpHandlerTestSuite) TestCreateOnlineTicketConfigErrParser() {
	var res string
	suite.cUC.On("CreateOnlineTicketConfig", mock.Anything, mock.Anything).Return(&res, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqM := request.OnlineTicketReq{
		UserId:     "id",
		Tag:        "tag",
		TotalQuota: 10,
		CountryList: []request.CountryList{
			{
				CountryNumber: 1,
				Percentage:    30,
			},
		},
	}
	requestBody, _ := json.Marshal(reqM)
	req := httptest.NewRequest(fiber.MethodPost, "/v1/create-online-ticket-config", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	ctx.Request().SetRequestURI("/v1/create-online-ticket-config")
	ctx.Request().Header.SetMethod(fiber.MethodPost)
	ctx.Request().Header.SetContentType("application/json")
	// ctx.Request().SetBody(requestBody)

	err := suite.handler.CreateOnlineTicketConfig(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *EventHttpHandlerTestSuite) TestCreateOnlineTicketConfigErrValidate() {
	var res string
	suite.cUC.On("CreateOnlineTicketConfig", mock.Anything, mock.Anything).Return(&res, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqM := request.OnlineTicketReq{}
	requestBody, _ := json.Marshal(reqM)
	req := httptest.NewRequest(fiber.MethodPost, "/v1/create-online-ticket-config", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	ctx.Request().SetRequestURI("/v1/create-online-ticket-config")
	ctx.Request().Header.SetMethod(fiber.MethodPost)
	ctx.Request().Header.SetContentType("application/json")
	ctx.Request().SetBody(requestBody)

	err := suite.handler.CreateOnlineTicketConfig(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *EventHttpHandlerTestSuite) TestCreateOnlineTicketConfigErr() {
	suite.cUC.On("CreateOnlineTicketConfig", mock.Anything, mock.Anything).Return(nil, errors.BadRequest("error"))
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqM := request.OnlineTicketReq{
		UserId:     "id",
		Tag:        "tag",
		TotalQuota: 10,
		CountryList: []request.CountryList{
			{
				CountryNumber: 1,
				Percentage:    30,
			},
		},
	}
	requestBody, _ := json.Marshal(reqM)
	req := httptest.NewRequest(fiber.MethodPost, "/v1/create-online-ticket-config", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	ctx.Request().SetRequestURI("/v1/create-online-ticket-config")
	ctx.Request().Header.SetMethod(fiber.MethodPost)
	ctx.Request().Header.SetContentType("application/json")
	ctx.Request().SetBody(requestBody)

	err := suite.handler.CreateOnlineTicketConfig(ctx)
	assert.Nil(suite.T(), err)
}
