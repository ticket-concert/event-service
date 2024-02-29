package commands

import (
	"context"
	"event-service/internal/modules/event"
	"event-service/internal/modules/event/models/entity"
	"event-service/internal/pkg/databases/mongodb"
	"event-service/internal/pkg/log"

	wrapper "event-service/internal/pkg/helpers"
)

type commandMongodbRepository struct {
	mongoDb mongodb.Collections
	logger  log.Logger
}

func NewCommandMongodbRepository(mongodb mongodb.Collections, log log.Logger) event.MongodbRepositoryCommand {
	return &commandMongodbRepository{
		mongoDb: mongodb,
		logger:  log,
	}
}

func (c commandMongodbRepository) InsertOneEventCollection(ctx context.Context, event entity.Event) <-chan wrapper.Result {
	output := make(chan wrapper.Result)

	go func() {
		resp := <-c.mongoDb.InsertOne(mongodb.InsertOne{
			CollectionName: "event",
			Document:       event,
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}
