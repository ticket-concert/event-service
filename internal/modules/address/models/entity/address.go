package entity

type Country struct {
	Id            int    `json:"id" bson:"id"`
	Code          string `json:"code" bson:"code"`
	Name          string `json:"name" bson:"name"`
	Iso3          string `json:"iso3" bson:"iso3"`
	Number        int    `json:"number" bson:"number"`
	ContinentCode string `json:"continentCode" bson:"continentCode"`
	ContinentName string `json:"continentName" bson:"continentName"`
	DisplayOrder  int    `json:"displayOrder" bson:"displayOrder"`
	FullName      string `json:"fullName" bson:"fullName"`
}

type Continent struct {
	Code string `json:"code" bson:"code"`
	Name string `json:"name" bson:"name"`
}
