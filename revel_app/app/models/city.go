package models

import "gopkg.in/olivere/elastic.v2"

// TODO: remove when done
/*
Country,City,AccentCity,Region,Population,Latitude,Longitude
us,new york,New York,NY,8107916,40.7141667,-74.0063889
ru,moscow,Moscow,48,10381288,55.752222,37.615556
*/

// City implements the Document interface.
// This struct partially maps to the ones that we get from our city database:
// https://www.maxmind.com/en/free-world-cities-database
type City struct {
	CityName    `json:"cityName"`
	AccentName  string  `json:"accentName"`  // aka "AccentCity"
	CountryCode string  `json:"countryCode"` // aka "Country"
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

// TODO: rename into something appropriate + change the mapping
type CityName struct {
	Name    string                `json:"name"` // aka "City"
	Suggest *elastic.SuggestField `json:"suggest"`
}

func (c City) String() string {
	return c.Name
}
