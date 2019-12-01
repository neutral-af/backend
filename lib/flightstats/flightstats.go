package flightstats

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/levigross/grequests"
	"github.com/neutral-af/backend/lib/config"
)

const baseURL = "https://api.flightstats.com/flex/schedules/rest/v1/json/flight"

var flightNumberRegex = regexp.MustCompile("([A-Z]+)([0-9]+)")

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
	reMatches := flightNumberRegex.FindStringSubmatch(flightNumber)
	path := fmt.Sprintf("/%s/%s/departing/%d/%d/%d", reMatches[0], reMatches[1], date.Year(), date.Month(), date.Day())
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

	// todo: parse real values
	return Details{"EGLL", "EGDY"}, nil
}

type Response struct {
	Error struct {
		HTTPStatusCode int
		ErrorMessage   string
	}
}
