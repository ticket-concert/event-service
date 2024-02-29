package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.elastic.co/apm"

	"event-service/internal/modules/address"
	"event-service/internal/modules/event"
	"event-service/internal/modules/event/models/entity"
	"event-service/internal/modules/event/models/request"
	"event-service/internal/modules/ticket"
	"event-service/internal/pkg/errors"
	"event-service/internal/pkg/log"

	addressEntity "event-service/internal/modules/address/models/entity"
	ticketEntity "event-service/internal/modules/ticket/models/entity"
	kafkaConfluent "event-service/internal/pkg/kafka/confluent"
)

type commandUsecase struct {
	eventRepositoryQuery    event.MongodbRepositoryQuery
	eventRepositoryCommand  event.MongodbRepositoryCommand
	ticketRepositoryCommand ticket.MongodbRepositoryCommand
	addressRepositoryQuery  address.MongodbRepositoryQuery
	kafkaProducer           kafkaConfluent.Producer
	logger                  log.Logger
}

func NewCommandUsecase(erq event.MongodbRepositoryQuery, erc event.MongodbRepositoryCommand,
	trc ticket.MongodbRepositoryCommand, arq address.MongodbRepositoryQuery, kp kafkaConfluent.Producer, log log.Logger) event.UsecaseCommand {
	return commandUsecase{
		eventRepositoryQuery:    erq,
		eventRepositoryCommand:  erc,
		ticketRepositoryCommand: trc,
		addressRepositoryQuery:  arq,
		kafkaProducer:           kp,
		logger:                  log,
	}
}

func (c commandUsecase) CreateEvent(origCtx context.Context, payload request.EventReq) (*string, error) {
	domain := "eventUsecase-CreateEvent"
	span, ctx := apm.StartSpanOptions(origCtx, domain, "function", apm.SpanOptions{
		Start:  time.Now(),
		Parent: apm.TraceContext{},
	})
	defer span.End()

	currentEvent := <-c.eventRepositoryQuery.FindEventByName(ctx, payload.Name)
	if currentEvent.Error != nil {
		return nil, currentEvent.Error
	}

	if currentEvent.Data != nil {
		return nil, errors.BadRequest("event already exist")
	}

	eventTagData := <-c.eventRepositoryQuery.FindEventByTag(ctx, payload.Tag)
	if eventTagData.Error != nil {
		return nil, eventTagData.Error
	}

	if eventTagData.Data != nil {
		eventTag, ok := eventTagData.Data.(*entity.Event)
		if !ok {
			return nil, errors.InternalServerError("failed marshal event")
		}

		if payload.UserId != eventTag.CreatedBy {
			return nil, errors.BadRequest("tag already exist, please create event with the same user to use this tag")
		}
	}

	dateTime, err := time.Parse("2006-01-02 15:04", payload.DateTime)
	if err != nil {
		return nil, errors.BadRequest("Format dateTime must be 'YYYY-MM-DD HH:MM'")
	}

	continentData := <-c.addressRepositoryQuery.FindOneContinentByCode(ctx, payload.ContinentCode)
	if continentData.Error != nil {
		return nil, currentEvent.Error
	}

	if continentData.Data == nil {
		return nil, errors.BadRequest("continent not found")
	}

	continent, ok := continentData.Data.(*addressEntity.Continent)
	if !ok {
		return nil, errors.InternalServerError("failed marshal continentData")
	}

	countryData := <-c.addressRepositoryQuery.FindOneCountry(ctx, payload.Country.Id)
	if countryData.Error != nil {
		return nil, currentEvent.Error
	}

	if countryData.Data == nil {
		return nil, errors.BadRequest("continent not found")
	}

	country, ok := countryData.Data.(*addressEntity.Country)
	if !ok {
		return nil, errors.InternalServerError("failed marshal countryData")
	}

	eventId := uuid.New().String()
	tiketIds := make([]string, 0)
	tickets := make([]ticketEntity.Ticket, 0)
	for _, v := range payload.Tickets {
		ticketId := uuid.New().String()
		ticket := ticketEntity.Ticket{
			TicketId:       ticketId,
			EventId:        eventId,
			TicketType:     v.TicketType,
			TicketPrice:    v.TicketPrice,
			TotalQuota:     v.TotalQuota,
			TotalRemaining: v.TotalQuota,
			Tag:            payload.Tag,
			Country: ticketEntity.Country{
				Name:  country.Name,
				Code:  country.Code,
				City:  payload.Country.City,
				Place: payload.Country.Place,
			},
			ContinentName: payload.ContinentName,
			ContinentCode: payload.ContinentCode,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		tickets = append(tickets, ticket)
		tiketIds = append(tiketIds, ticketId)
	}
	respTicket := <-c.ticketRepositoryCommand.InsertManyTicketCollection(ctx, tickets)
	if respTicket.Error != nil {
		return nil, respTicket.Error
	}

	event := entity.Event{
		EventId:  eventId,
		Name:     payload.Name,
		DateTime: dateTime,
		// Location:      payload.Location,
		ContinentName: continent.Name,
		ContinentCode: continent.Code,
		Country: entity.Country{
			Name:  country.Name,
			Code:  country.Code,
			City:  payload.Country.City,
			Place: payload.Country.Place,
		},
		Description: payload.Description,
		Tag:         payload.Tag,
		TicketIds:   tiketIds,
		CreatedBy:   payload.UserId,
		UpdatedBy:   payload.UserId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	respEvent := <-c.eventRepositoryCommand.InsertOneEventCollection(ctx, event)
	if respEvent.Error != nil {
		return nil, errors.InternalServerError("failed save event")
	}

	for _, ticketId := range tiketIds {
		createTicketReq := request.CreateTicketReq{
			TicketId: ticketId,
			EventId:  eventId,
		}
		marshaledKafkaData, _ := json.Marshal(createTicketReq)
		topic := "concert-create-bank-ticket"
		c.kafkaProducer.Publish(topic, marshaledKafkaData, nil)
		c.logger.Info(ctx, fmt.Sprintf("Send kafka create bank ticket, ticketId : %s", ticketId), fmt.Sprintf("%+v", createTicketReq))
	}

	rs := "Success create event"
	return &rs, nil
}

func (c commandUsecase) CreateOnlineTicketConfig(origCtx context.Context, payload request.OnlineTicketReq) (*string, error) {
	domain := "eventUsecase-CreateOnlineTicketConfig"
	span, ctx := apm.StartSpanOptions(origCtx, domain, "function", apm.SpanOptions{
		Start:  time.Now(),
		Parent: apm.TraceContext{},
	})
	defer span.End()

	eventTagData := <-c.eventRepositoryQuery.FindEventByTag(ctx, payload.Tag)
	if eventTagData.Error != nil {
		return nil, eventTagData.Error
	}

	if eventTagData.Data != nil {
		eventTag, ok := eventTagData.Data.(*entity.Event)
		if !ok {
			return nil, errors.InternalServerError("failed marshal event")
		}

		fmt.Println(eventTag)

		if payload.UserId != eventTag.CreatedBy {
			return nil, errors.BadRequest("tag already exist, please create event with the same user to use this tag")
		}
	}

	countryList := make([]ticketEntity.CountryList, 0)
	totalPercentage := 0
	for _, v := range payload.CountryList {
		totalPercentage = totalPercentage + v.Percentage
		countryList = append(countryList, ticketEntity.CountryList{
			CountryNumber: v.CountryNumber,
			Percentage:    v.Percentage,
		})
	}
	if totalPercentage != 100 {
		return nil, errors.BadRequest("total percentage must be 100")
	}
	otConfig := ticketEntity.OnlineTicketConfig{
		Tag:         payload.Tag,
		TotalQuota:  payload.TotalQuota,
		CountryList: countryList,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		CreatedBy:   payload.UserId,
		UpdatedBy:   payload.UserId,
	}
	resp := <-c.ticketRepositoryCommand.UpsertOneOnlineTicketConfig(ctx, otConfig)
	if resp.Error != nil {
		return nil, eventTagData.Error
	}
	result := "Success create online ticket config"
	return &result, nil
}
