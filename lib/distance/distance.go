package distance

import (
	"errors"

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
		airportMap[i.ICAO] = airportGeo{
			Latitude:  i.Latitude,
			Longitude: i.Longitude,
		}
	}
}

func getAirportGeo(icao string) airportGeo {
	return airportMap[icao]
}

// TwoAirports returns the Great Circle Distance between the airports given by the
// two ICAO codes
func TwoAirports(departureCode string, arrivalCode string) (float64, error) {
	departureData := getAirportGeo(departureCode)
	arrivalData := getAirportGeo(arrivalCode)
	if departureData == (airportGeo{}) || arrivalData == (airportGeo{}) {
		return 0, errors.New("Airport not found")
	}

	departureGeo := geo.NewPoint(departureData.Latitude, departureData.Longitude)
	arrivalGeo := geo.NewPoint(arrivalData.Latitude, arrivalData.Longitude)

	distance := departureGeo.GreatCircleDistance(arrivalGeo)
	return distance, nil
}
