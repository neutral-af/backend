package distance

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAirportGeo(t *testing.T) {
	airport := getAirportGeo("CYYZ")
	assert.NotZero(t, airport.Latitude)
	assert.NotZero(t, airport.Longitude)
}

func TestTwoAirports(t *testing.T) {
	dist, err := TwoAirports("CYYZ", "EDDB")
	assert.NoError(t, err)
	assert.Greater(t, dist, 6000.0)
}
