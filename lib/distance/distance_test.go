package distance

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAirportGeo(t *testing.T) {
	airport := getAirportGeo("YYZ")
	assert.NotZero(t, airport.Latitude)
	assert.NotZero(t, airport.Longitude)
}

func TestTwoAirports(t *testing.T) {
	dist := TwoAirports("YYZ", "SXF")
	assert.Greater(t, dist, 6000.0)
}
