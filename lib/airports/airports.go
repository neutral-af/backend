package airports

import (
	"fmt"

	"github.com/mmcloughlin/openflights"

	models "github.com/neutral-af/backend/lib/graphql-models"
)

type Airport openflights.Airport

type Airports []Airport

// Implement the searchable interface for the fuzzy-searcher
func (a Airports) String(i int) string {
	return fmt.Sprintf("%s %s %s %s %s", a[i].City, a[i].IATA, a[i].Name, a[i].Country, a[i].ICAO)
}

func (a Airports) Len() int {
	return len(a)
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

var allAirports Airports
var airportsByICAO map[string]Airport
var airportsByIATA map[string]Airport

func init() {
	airportsByICAO = make(map[string]Airport)
	airportsByIATA = make(map[string]Airport)
	for _, i := range openflights.Airports {
		a := Airport(i)
		allAirports = append(allAirports, a)
		airportsByICAO[i.ICAO] = a
		airportsByIATA[i.IATA] = a
	}
}

func GetAll() Airports {
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
