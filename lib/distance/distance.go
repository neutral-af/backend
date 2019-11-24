package distance

import (
	geo "github.com/kellydunn/golang-geo"

	"github.com/neutral-af/backend/lib/airports"
)

type airportGeo struct {
	Latitude  float64 `json:"latitude,string"`
	Longitude float64 `json:"longitude,string"`
}

// TwoAirports returns the Great Circle Distance between the airports given by the
// two ICAO codes
func TwoAirports(departureCode string, arrivalCode string) (float64, error) {
	departureData, err := airports.GetFromICAO(departureCode)
	if err != nil {
		return 0, err
	}

	arrivalData, err := airports.GetFromICAO(arrivalCode)
	if err != nil {
		return 0, err
	}

	departureGeo := geo.NewPoint(departureData.Latitude, departureData.Longitude)
	arrivalGeo := geo.NewPoint(arrivalData.Latitude, arrivalData.Longitude)

	distance := departureGeo.GreatCircleDistance(arrivalGeo)
	return distance, nil
}
