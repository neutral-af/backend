package flightstats

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/levigross/grequests"
	"github.com/neutral-af/backend/lib/config"
)

const baseURL = "https://api.flightstats.com"

var flightNumberRegex = regexp.MustCompile("([A-Z]+)([0-9]+)")

func splitFlightNumber(flightNumber string) (airline string, flight string, err error) {
	// The first element will be the whole string, successive elements are capture groups
	reMatches := flightNumberRegex.FindStringSubmatch(strings.ReplaceAll(flightNumber, " ", ""))

	if len(reMatches) < 3 || reMatches[1] == "" || reMatches[2] == "" {
		return "", "", fmt.Errorf("Unable to parse airline and number from flight number: %s", flightNumber)
	}

	return reMatches[1], reMatches[2], nil
}

type FlightStats struct {
	baseURL    string
	authParams map[string]string
}

type Details struct {
	Departure string
	Arrival   string
}

func New() FlightStats {
	return FlightStats{
		authParams: map[string]string{
			"appId":  config.C.FlightStatsAppID,
			"appKey": config.C.FlightStatsAppKey,
		},
	}
}

func buildAPIPath(flightNumber string, date time.Time) (string, error) {
	airlineCode, flightCode, err := splitFlightNumber(flightNumber)
	if err != nil {
		return "", err
	}

	var path string
	if date.Before(time.Now().Add(-time.Hour * 24 * 7)) {
		return "", fmt.Errorf(
			`Flight number lookups can only be for flights in the future or in the past 7 days (got date: %s).
			For flights more than 7 days old, use the departure/arrival airports search instead`,
			date.Format("2006-01-02"),
		)
	} else if date.Before(time.Now().Add(time.Hour * 24 * 3)) {
		// "Flight Status API" contains previous 7 days and next 3 days
		path = fmt.Sprintf(
			"/flex/flightstatus/rest/v2/json/flight/status/%s/%s/dep/%d/%d/%d",
			airlineCode, flightCode, date.Year(), date.Month(), date.Day(),
		)
	} else {
		// Use "Schedules API" for flights more than 3 days in the future
		path = fmt.Sprintf(
			"/flex/schedules/rest/v1/json/flight/%s/%s/departing/%d/%d/%d",
			airlineCode, flightCode, date.Year(), date.Month(), date.Day(),
		)
	}
	return path, nil
}

func (f *FlightStats) GetAirportsForFlight(flightNumber string, date time.Time) (Details, error) {
	path, err := buildAPIPath(flightNumber, date)
	if err != nil {
		return Details{}, err
	}
	resp, err := grequests.Get(baseURL+path, &grequests.RequestOptions{
		Params: f.authParams,
	})
	if err != nil {
		return Details{}, err
	}

	var responseData Response
	err = resp.JSON(&responseData)
	if err != nil {
		return Details{}, nil
	}
	if responseData.Error.ErrorMessage != "" {
		return Details{}, fmt.Errorf("Error in FlightStats call: %w", errors.New(responseData.Error.ErrorMessage))
	}
	airports := responseData.Appendix.Airports
	if len(airports) < 2 {
		return Details{}, fmt.Errorf("Could not find flight number %s on %s", flightNumber, date.Format("2006-01-02"))
	}

	return Details{airports[0].ICAO, airports[1].ICAO}, nil
}

type Response struct {
	Appendix struct {
		Airports []struct {
			ICAO string `json:"icao"`
		}
	}
	Error struct {
		HTTPStatusCode int
		ErrorMessage   string
	}
}
