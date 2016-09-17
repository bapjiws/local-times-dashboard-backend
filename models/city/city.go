package city

import "gopkg.in/olivere/elastic.v3"

// City implements the Document interface.
// This struct partially maps to the ones that we get from our city database:
// https://www.maxmind.com/en/free-world-cities-database
type City struct {
	Id          string                `json:"id"`
	Name        string                `json:"name"`        // aka "City"
	AccentName  string                `json:"accentName"`  // aka "AccentCity"
	CountryCode string                `json:"countryCode"` // aka "Country"
	Latitude    float64               `json:"latitude"`
	Longitude   float64               `json:"longitude"`
	Suggest     *elastic.SuggestField `json:"suggest"`
}

// City implements the Stringer interface (see https://golang.org/pkg/fmt/), so it can be printed by, say, AddDocument.
func (c City) String() string {
	return c.Name
}
