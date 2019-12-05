package flightstats

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSplitFlightNumberSuccess(t *testing.T) {
	airline, flight, err := splitFlightNumber("AC123")
	assert.NoError(t, err)
	assert.Equal(t, "AC", airline)
	assert.Equal(t, "123", flight)
}

func TestSplitFlightNumberFail(t *testing.T) {
	airline, flight, err := splitFlightNumber("ACAC")
	assert.Error(t, err)
	assert.Empty(t, airline)
	assert.Empty(t, flight)
}

func TestSplitFlightNumberEmpty(t *testing.T) {
	airline, flight, err := splitFlightNumber("")
	assert.Error(t, err)
	assert.Empty(t, airline)
	assert.Empty(t, flight)
}

func TestGetAirportsForFlightFailPast(t *testing.T) {
	f := New()
	details, err := f.GetAirportsForFlight("", time.Now().Add(-time.Hour))
	assert.Empty(t, details)
	assert.Error(t, err)
}
