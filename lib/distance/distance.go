package distance

import (
	geo "github.com/kellydunn/golang-geo"
	"github.com/mmcloughlin/openflights"
)

type airportGeo struct {
	Latitude  float64 `json:"latitude,string"`
	Longitude float64 `json:"longitude,string"`
}

var airportMap map[string]airportGeo

func init() {
	airportMap = make(map[string]airportGeo)
	for _, i := range openflights.Airports {
		airportMap[i.IATA] = airportGeo{
			Latitude:  i.Latitude,
			Longitude: i.Longitude,
		}
	}
}

func getAirportGeo(iata string) airportGeo {
	return airportMap[iata]
}

// TwoAirports returns the Great Circle Distance between the airports given by the
// two IATA codes
func TwoAirports(departureCode string, arrivalCode string) float64 {
	departureData := getAirportGeo(departureCode)
	arrivalData := getAirportGeo(arrivalCode)

	departureGeo := geo.NewPoint(departureData.Latitude, departureData.Longitude)
	arrivalGeo := geo.NewPoint(arrivalData.Latitude, arrivalData.Longitude)

	distance := departureGeo.GreatCircleDistance(arrivalGeo)
	return distance
}
