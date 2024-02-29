package queries

import (
	"context"
	"event-service/internal/modules/event"
	"event-service/internal/modules/event/models/entity"
	"event-service/internal/modules/event/models/request"
	"event-service/internal/pkg/databases/mongodb"
	wrapper "event-service/internal/pkg/helpers"
	"event-service/internal/pkg/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type queryMongodbRepository struct {
	mongoDb mongodb.Collections
	logger  log.Logger
}

func NewQueryMongodbRepository(mongodb mongodb.Collections, log log.Logger) event.MongodbRepositoryQuery {
	return &queryMongodbRepository{
		mongoDb: mongodb,
		logger:  log,
	}
}

func (q queryMongodbRepository) FindEventByName(ctx context.Context, name string) <-chan wrapper.Result {
	var event entity.Event
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.FindOne(mongodb.FindOne{
			Result:         &event,
			CollectionName: "event",
			Filter: bson.M{
				"name": name,
			},
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}

func (q queryMongodbRepository) FindEventByTag(ctx context.Context, tag string) <-chan wrapper.Result {
	var event entity.Event
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.FindOne(mongodb.FindOne{
			Result:         &event,
			CollectionName: "event",
			Filter: bson.M{
				"tag": tag,
			},
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}

func (q queryMongodbRepository) FindAllEvent(ctx context.Context, payload request.AllEventReq) <-chan wrapper.Result {
	var event []entity.Event
	var countData int64
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.FindAllData(mongodb.FindAllData{
			Result:         &event,
			CountData:      &countData,
			CollectionName: "event",
			Filter:         bson.M{"name": primitive.Regex{Pattern: ".*" + payload.Search + ".*", Options: "i"}},
			Sort: &mongodb.Sort{
				FieldName: "name",
				By:        mongodb.SortAscending,
			},
			Page: payload.Page,
			Size: payload.Size,
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}
