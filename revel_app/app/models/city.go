package models

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
	CountryCode string  `json:"countryCode"` // aka "Country"
	Name        string  `json:"name"`        // aka "AccentCity"
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

func (c City) String() string {
	return c.Name
}