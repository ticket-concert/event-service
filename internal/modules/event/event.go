package event

import (
	"context"
	"event-service/internal/modules/event/models/entity"
	"event-service/internal/modules/event/models/request"
	"event-service/internal/modules/event/models/response"
	wrapper "event-service/internal/pkg/helpers"
)

type UsecaseQuery interface {
	FindEvents(origCtx context.Context, payload request.AllEventReq) (*response.EventResp, error)
}

type UsecaseCommand interface {
	CreateEvent(origCtx context.Context, payload request.EventReq) (*string, error)
	CreateOnlineTicketConfig(origCtx context.Context, payload request.OnlineTicketReq) (*string, error)
}

type MongodbRepositoryQuery interface {
	FindEventByName(ctx context.Context, name string) <-chan wrapper.Result
	FindEventByTag(ctx context.Context, tag string) <-chan wrapper.Result
	FindAllEvent(ctx context.Context, payload request.AllEventReq) <-chan wrapper.Result
}

type MongodbRepositoryCommand interface {
	InsertOneEventCollection(ctx context.Context, event entity.Event) <-chan wrapper.Result
}
