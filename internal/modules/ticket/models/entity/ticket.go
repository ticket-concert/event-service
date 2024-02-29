package entity

import "time"

type Country struct {
	Name  string `json:"name" bson:"name"`
	Code  string `json:"code" bson:"code"`
	City  string `json:"city" bson:"city"`
	Place string `json:"place" bson:"place"`
}

type Ticket struct {
	TicketId       string    `json:"ticketId" bson:"ticketId"`
	EventId        string    `json:"eventId" bson:"eventId"`
	TicketType     string    `json:"ticketType" bson:"ticketType"`
	TicketPrice    int       `json:"ticketPrice" bson:"ticketPrice"`
	TotalQuota     int       `json:"totalQuota" bson:"totalQuota"`
	TotalRemaining int       `json:"totalRemaining" bson:"totalRemaining"`
	ContinentName  string    `json:"continentName" bson:"continentName"`
	ContinentCode  string    `json:"continentCode" bson:"continentCode"`
	Country        Country   `json:"country" bson:"country"`
	Tag            string    `json:"tag" bson:"tag"`
	CreatedAt      time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt" bson:"updatedAt"`
}

type OnlineTicketConfig struct {
	Tag         string        `json:"tag" bson:"tag"`
	TotalQuota  int           `json:"totalQuota" bson:"totalQuota"`
	CountryList []CountryList `json:"countryList" bson:"countryList"`
	CreatedAt   time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt" bson:"updatedAt"`
	CreatedBy   string        `json:"createdBy" bson:"createdBy"`
	UpdatedBy   string        `json:"updatedBy" bson:"updatedBy"`
}

type CountryList struct {
	CountryNumber int    `json:"countryNumber" bson:"countryNumber"`
	Percentage    int    `json:"percentage" bson:"percentage"`
	CountryCode   string `json:"countryCode" bson:"countryCode"`
}

type AggregateTotalTicket struct {
	Id                   string `json:"_id" bson:"_id"`
	CountryName          string `json:"countryName" bson:"countryName"`
	TotalAvailableTicket int    `json:"totalAvailableTicket" bson:"totalAvailableTicket"`
}
