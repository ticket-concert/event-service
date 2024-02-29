package ticket

import (
	"context"
	"event-service/internal/modules/ticket/models/entity"
	wrapper "event-service/internal/pkg/helpers"
)

type MongodbRepositoryCommand interface {
	InsertManyTicketCollection(ctx context.Context, ticket []entity.Ticket) <-chan wrapper.Result
	UpsertOneOnlineTicketConfig(ctx context.Context, payload entity.OnlineTicketConfig) <-chan wrapper.Result
}
