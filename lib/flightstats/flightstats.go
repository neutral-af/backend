package flightstats

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/levigross/grequests"
	"github.com/neutral-af/backend/lib/config"
)

const baseURL = "https://api.flightstats.com"

var flightNumberRegex = regexp.MustCompile("([A-Z]+)([0-9]+)")

func splitFlightNumber(flightNumber string) (airline string, flight string, err error) {
	// The first element will be the whole string, successive elements are capture groups
	reMatches := flightNumberRegex.FindStringSubmatch(flightNumber)

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

func (f *FlightStats) GetAirportsForFlight(flightNumber string, date time.Time) (Details, error) {
	var path string
	airlineCode, flightCode, err := splitFlightNumber(flightNumber)
	if date.Before(time.Now()) {
		// Use historical API
		path = fmt.Sprintf(
			"/flex/flightstatus/historical/rest/v3/json/flight/status/%s/%s/dep/%d/%d/%d",
			airlineCode, flightCode, date.Year(), date.Month(), date.Day(),
		)
	} else {
		// Use schedules API
		path = fmt.Sprintf(
			"/flex/schedules/rest/v1/json/flight/%s/%s/departing/%d/%d/%d",
			airlineCode, flightCode, date.Year(), date.Month(), date.Day(),
		)
	}

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
