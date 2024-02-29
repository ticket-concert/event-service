package response

import (
	"event-service/internal/pkg/constants"
	"time"
)

type Country struct {
	Name  string `json:"name" bson:"name"`
	Code  string `json:"code" bson:"code"`
	City  string `json:"city" bson:"city"`
	Place string `json:"place" bson:"place"`
}

type Event struct {
	EventId       string    `json:"eventId" bson:"eventId"`
	Name          string    `json:"name" bson:"name"`
	DateTime      time.Time `json:"dateTime" bson:"dateTime"`
	ContinentName string    `json:"continentName" bson:"continentName"`
	ContinentCode string    `json:"continentCode" bson:"continentCode"`
	Country       Country   `json:"country" bson:"country"`
	Description   string    `json:"description" bson:"description"`
	Tag           string    `json:"tag" bson:"tag"`
	TicketIds     []string  `json:"ticketIds" bson:"ticketIds"`
}

type EventResp struct {
	CollectionData []Event
	MetaData       constants.MetaData
}
