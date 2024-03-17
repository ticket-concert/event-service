package handlers

import (
	"event-service/internal/modules/event"
	"event-service/internal/modules/event/models/request"
	"event-service/internal/pkg/errors"
	"event-service/internal/pkg/helpers"
	"event-service/internal/pkg/log"
	"event-service/internal/pkg/redis"
	"fmt"

	middlewares "event-service/configs/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type EventHttpHandler struct {
	EventUsecaseCommand event.UsecaseCommand
	EventUsecaseQuery   event.UsecaseQuery
	Logger              log.Logger
	Validator           *validator.Validate
}

func InitEventHttpHandler(app *fiber.App, euc event.UsecaseCommand, euq event.UsecaseQuery, log log.Logger, redisClient redis.Collections) {
	handler := &EventHttpHandler{
		EventUsecaseCommand: euc,
		EventUsecaseQuery:   euq,
		Logger:              log,
		Validator:           validator.New(),
	}
	middlewares := middlewares.NewMiddlewares(redisClient)
	route := app.Group("/api/event")

	route.Post("/v1/create-event", middlewares.VerifyBearer(), handler.CreateEvent)
	route.Post("/v1/create-online-ticket-config", middlewares.VerifyBearer(), handler.CreateOnlineTicketConfig)
	route.Get("/v1/list", middlewares.VerifyBearer(), handler.GetEvents)
}

func (e EventHttpHandler) CreateEvent(c *fiber.Ctx) error {
	req := new(request.EventReq)
	if err := c.BodyParser(req); err != nil {
		return helpers.RespError(c, e.Logger, errors.BadRequest("bad request"))
	}

	userId := c.Locals("userId").(string)
	req.UserId = userId

	if err := e.Validator.Struct(req); err != nil {
		return helpers.RespError(c, e.Logger, errors.BadRequest(err.Error()))
	}
	resp, err := e.EventUsecaseCommand.CreateEvent(c.Context(), *req)
	if err != nil {
		return helpers.RespCustomError(c, e.Logger, err)
	}
	return helpers.RespSuccess(c, e.Logger, resp, "Create event success")
}

func (e EventHttpHandler) GetEvents(c *fiber.Ctx) error {
	req := new(request.AllEventReq)
	if err := c.QueryParser(req); err != nil {
		return helpers.RespError(c, e.Logger, errors.BadRequest("bad request"))
	}

	if err := e.Validator.Struct(req); err != nil {
		return helpers.RespError(c, e.Logger, errors.BadRequest(err.Error()))
	}
	resp, err := e.EventUsecaseQuery.FindEvents(c.Context(), *req)
	if err != nil {
		return helpers.RespCustomError(c, e.Logger, err)
	}
	return helpers.RespPagination(c, e.Logger, resp.CollectionData, resp.MetaData, "Get country success")
}

func (e EventHttpHandler) CreateOnlineTicketConfig(c *fiber.Ctx) error {
	req := new(request.OnlineTicketReq)
	if err := c.BodyParser(req); err != nil {
		return helpers.RespError(c, e.Logger, errors.BadRequest("bad request"))
	}

	userId := c.Locals("userId").(string)
	req.UserId = userId

	if err := e.Validator.Struct(req); err != nil {
		fmt.Println(err)
		return helpers.RespError(c, e.Logger, errors.BadRequest(err.Error()))
	}
	resp, err := e.EventUsecaseCommand.CreateOnlineTicketConfig(c.Context(), *req)
	if err != nil {
		return helpers.RespCustomError(c, e.Logger, err)
	}
	return helpers.RespSuccess(c, e.Logger, resp, "Create online ticket config success")
}
