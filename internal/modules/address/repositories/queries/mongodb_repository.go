package queries

import (
	"context"
	"event-service/internal/modules/address"
	"event-service/internal/modules/address/models/entity"
	"event-service/internal/pkg/databases/mongodb"
	wrapper "event-service/internal/pkg/helpers"
	"event-service/internal/pkg/log"

	"go.mongodb.org/mongo-driver/bson"
)

type queryMongodbRepository struct {
	mongoDb mongodb.Collections
	logger  log.Logger
}

func NewQueryMongodbRepository(mongodb mongodb.Collections, log log.Logger) address.MongodbRepositoryQuery {
	return &queryMongodbRepository{
		mongoDb: mongodb,
		logger:  log,
	}
}

func (q queryMongodbRepository) FindOneCountry(ctx context.Context, id int) <-chan wrapper.Result {
	var country entity.Country
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.FindOne(mongodb.FindOne{
			Result:         &country,
			CollectionName: "country",
			Filter: bson.M{
				"id": id,
			},
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}

func (q queryMongodbRepository) FindOneContinentByCode(ctx context.Context, code string) <-chan wrapper.Result {
	var continent entity.Continent
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.FindOne(mongodb.FindOne{
			Result:         &continent,
			CollectionName: "continent",
			Filter: bson.M{
				"code": code,
			},
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}
