package user

import (
	"context"
	wrapper "event-service/internal/pkg/helpers"
)

type MongodbRepositoryQuery interface {
	FindOneUserId(ctx context.Context, userId string) <-chan wrapper.Result
}
