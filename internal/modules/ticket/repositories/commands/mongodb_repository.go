package commands

import (
	"context"
	"event-service/internal/modules/ticket"
	"event-service/internal/modules/ticket/models/entity"
	"event-service/internal/pkg/databases/mongodb"
	wrapper "event-service/internal/pkg/helpers"
	"event-service/internal/pkg/log"

	"go.mongodb.org/mongo-driver/bson"
)

type commandMongodbRepository struct {
	mongoDb mongodb.Collections
	logger  log.Logger
}

func NewCommandMongodbRepository(mongodb mongodb.Collections, log log.Logger) ticket.MongodbRepositoryCommand {
	return &commandMongodbRepository{
		mongoDb: mongodb,
		logger:  log,
	}
}

func (c commandMongodbRepository) InsertManyTicketCollection(ctx context.Context, ticket []entity.Ticket) <-chan wrapper.Result {
	output := make(chan wrapper.Result)

	go func() {
		var documentsInsert []interface{}
		for _, v := range ticket {
			documentsInsert = append(documentsInsert, v)
		}
		resp := <-c.mongoDb.InsertMany(mongodb.InsertMany{
			CollectionName: "ticket-detail",
			Documents:      documentsInsert,
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}

func (c commandMongodbRepository) UpsertOneOnlineTicketConfig(ctx context.Context, payload entity.OnlineTicketConfig) <-chan wrapper.Result {
	output := make(chan wrapper.Result)

	go func() {
		resp := <-c.mongoDb.UpsertOne(mongodb.UpdateOne{
			CollectionName: "online-ticket-config",
			Filter: bson.M{
				"tag": payload.Tag,
			},
			Document: payload,
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}
