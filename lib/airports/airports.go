package airports

import (
	"fmt"

	"github.com/mmcloughlin/openflights"
	"github.com/sahilm/fuzzy"

	models "github.com/neutral-af/backend/lib/graphql-models"
)

type Airport openflights.Airport

type Airports []Airport

// Implement the searchable interface for the fuzzy-searcher
func (a Airports) String(i int) string {
	return a[i].SearchableString()
}

func (a Airports) Len() int {
	return len(a)
}

func (a Airport) SearchableString() string {
	return fmt.Sprintf("%s %s %s %s %s", a.Name, a.IATA, a.City, a.Country, a.ICAO)
}

func (a Airport) ToModel() models.Airport {
	return models.Airport{
		Name:    a.Name,
		Icao:    a.ICAO,
		Iata:    a.IATA,
		City:    a.City,
		Country: a.Country,
	}
}

var allAirports []Airport
var allAirportsSearchableString []string
var airportsByICAO = make(map[string]Airport)
var airportsByIATA = make(map[string]Airport)
var airportsBySearchableString = make(map[string]Airport)

func init() {

	for _, i := range openflights.Airports {
		a := Airport(i)
		allAirports = append(allAirports, a)
		allAirportsSearchableString = append(allAirportsSearchableString, a.SearchableString())
		airportsByICAO[a.ICAO] = a
		airportsByIATA[a.IATA] = a
		airportsBySearchableString[a.SearchableString()] = a
	}

}

func Search(query string) []Airport {
	results := fuzzy.FindFrom(query, Airports(allAirports))

	matches := []Airport{}
	for _, r := range results {
		if len(matches) > 10 {
			break
		}
		a := allAirports[r.Index]
		matches = append(matches, a)
	}

	return matches
}

func GetAll() []Airport {
	return allAirports
}

func GetFromICAO(code string) (Airport, error) {
	a, ok := airportsByICAO[code]

	if !ok {
		return Airport{}, fmt.Errorf("Airport for code not found: %s", code)
	}

	return a, nil
}

func GetFromIATA(code string) (Airport, error) {
	a, ok := airportsByIATA[code]

	if !ok {
		return Airport{}, fmt.Errorf("Airport for code not found: %s", code)
	}

	return a, nil
}
