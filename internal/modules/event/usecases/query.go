package usecases

import (
	"context"
	"event-service/internal/modules/event"
	"event-service/internal/modules/event/models/entity"
	"event-service/internal/modules/event/models/request"
	"event-service/internal/modules/event/models/response"
	"event-service/internal/pkg/errors"
	"event-service/internal/pkg/helpers"
	"event-service/internal/pkg/log"
	"fmt"
	"time"

	"go.elastic.co/apm"
)

type queryUsecase struct {
	eventRepositoryQuery event.MongodbRepositoryQuery
	logger               log.Logger
}

func NewQueryUsecase(emq event.MongodbRepositoryQuery, log log.Logger) event.UsecaseQuery {
	return queryUsecase{
		eventRepositoryQuery: emq,
		logger:               log,
	}
}

func (q queryUsecase) FindEvents(origCtx context.Context, payload request.AllEventReq) (*response.EventResp, error) {
	domain := "addressUsecase-FindCountries"
	span, ctx := apm.StartSpanOptions(origCtx, domain, "function", apm.SpanOptions{
		Start:  time.Now(),
		Parent: apm.TraceContext{},
	})
	defer span.End()

	resp := <-q.eventRepositoryQuery.FindAllEvent(ctx, payload)
	if resp.Error != nil {
		msg := "Error query event"
		q.logger.Error(ctx, msg, fmt.Sprintf("%+v", resp.Error))
		return nil, resp.Error
	}

	if resp.Data == nil {
		msg := "Event Not Found"
		q.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		return nil, errors.NotFound("event not found")
	}

	event, ok := resp.Data.(*[]entity.Event)
	if !ok {
		return nil, errors.InternalServerError("cannot parsing data")
	}

	var collectionData = make([]response.Event, 0)
	for _, value := range *event {
		collectionData = append(collectionData, response.Event{
			EventId:  value.EventId,
			Name:     value.Name,
			DateTime: value.DateTime,
			// Location:      value.Location,
			ContinentCode: value.ContinentCode,
			ContinentName: value.ContinentName,
			Country:       response.Country(value.Country),
			Description:   value.Description,
			Tag:           value.Tag,
			TicketIds:     value.TicketIds,
		})
	}

	return &response.EventResp{
		CollectionData: collectionData,
		MetaData:       helpers.GenerateMetaData(resp.Count, int64(len(*event)), payload.Page, payload.Size),
	}, nil

}
