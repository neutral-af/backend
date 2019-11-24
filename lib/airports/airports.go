package airports

import (
	"errors"
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

func init() {
	airportsByICAO = make(map[string]Airport)
	for _, i := range openflights.Airports {
		allAirports = append(allAirports, Airport(i))
		airportsByICAO[i.ICAO] = Airport(i)
	}
}

func GetAll() Airports {
	return allAirports
}

func GetFromICAO(code string) (Airport, error) {
	a, ok := airportsByICAO[code]

	if !ok {
		return Airport{}, errors.New("Airport for code not found")
	}

	return a, nil
}
