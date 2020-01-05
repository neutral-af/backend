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

func TestSplitFlightNumberWhitespace(t *testing.T) {
	airline, flight, err := splitFlightNumber("AC 123")
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

func TestBuildAPIPathTooOld(t *testing.T) {
	path, err := buildAPIPath("AC123", time.Now().Add(-time.Hour*24*8))
	assert.Empty(t, path)
	assert.Error(t, err)
}

func TestBuildAPIPathFarFuture(t *testing.T) {
	path, err := buildAPIPath("AC123", time.Now().Add(time.Hour*24*90))
	assert.NoError(t, err)
	assert.Contains(t, path, "schedules")
}

func TestBuildAPIPathToday(t *testing.T) {
	path, err := buildAPIPath("AC123", time.Now())
	assert.NoError(t, err)
	assert.Contains(t, path, "flightstatus")
}

func TestBuildAPIPathThreeDaysFuture(t *testing.T) {
	path, err := buildAPIPath("AC123", time.Now().Add(time.Hour*24*3))
	assert.NoError(t, err)
	assert.Contains(t, path, "flightstatus")
}

func TestBuildAPIPathFourDaysFuture(t *testing.T) {
	path, err := buildAPIPath("AC123", time.Now().Add(time.Hour*24*4))
	assert.NoError(t, err)
	assert.Contains(t, path, "schedules")
}

func TestGetAirportsForFlightTooOld(t *testing.T) {
	f := New()
	details, err := f.GetAirportsForFlight("AC123", time.Now().Add(-time.Hour*24*8))
	assert.Empty(t, details)
	assert.Error(t, err)
}
