package request

type TicketReq struct {
	// ContinentCode string `json:"continentCode"`
	CountryCode string `json:"countryCode" validate:"required"`
}

type Country struct {
	Name  string `json:"name" validate:"required"`
	Id    int    `json:"id" validate:"required"`
	City  string `json:"city" validate:"required"`
	Place string `json:"place" validate:"required"`
}

type EventReq struct {
	Name          string   `json:"name" validate:"required"`
	DateTime      string   `json:"dateTime" validate:"required"`
	ContinentName string   `json:"continentName" validate:"required"`
	ContinentCode string   `json:"continentCode" validate:"required"`
	Country       Country  `json:"country" validate:"required"`
	Description   string   `json:"description" validate:"required"`
	Tag           string   `json:"tag" validate:"required"`
	Tickets       []Ticket `json:"tickets" validate:"required"`
	UserId        string   `json:"userId" validate:"required"`
}

type OnlineTicketReq struct {
	UserId      string        `json:"userId" validate:"required"`
	Tag         string        `json:"tag" validate:"required"`
	TotalQuota  int           `json:"totalQuota" validate:"required"`
	CountryList []CountryList `json:"countryList" validate:"required"`
}

type CountryList struct {
	CountryNumber int `json:"countryNumber" validate:"required"`
	Percentage    int `json:"percentage" validate:"required"`
}

type Ticket struct {
	TicketType  string `json:"ticketType" validate:"required"`
	TicketPrice int    `json:"ticketPrice" validate:"required"`
	TotalQuota  int    `json:"totalQuota" validate:"required"`
	Tag         string `json:"tag" validate:"required"`
}

type AllEventReq struct {
	Page   int64  `query:"page" validate:"required"`
	Size   int64  `query:"size" validate:"required"`
	Search string `query:"search"`
}

type CreateTicketReq struct {
	TicketId string `json:"ticketId" bson:"ticketId"`
	EventId  string `json:"eventId" bson:"eventId"`
}
