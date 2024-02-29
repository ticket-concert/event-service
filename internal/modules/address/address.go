package address

import (
	"context"
	wrapper "event-service/internal/pkg/helpers"
)

type MongodbRepositoryQuery interface {
	FindOneCountry(ctx context.Context, id int) <-chan wrapper.Result
	FindOneContinentByCode(ctx context.Context, code string) <-chan wrapper.Result
}
