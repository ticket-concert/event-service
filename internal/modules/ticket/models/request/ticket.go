package request

type TicketReq struct {
	CountryCode string `json:"countryCode" validate:"required"`
}
