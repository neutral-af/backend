package flightstats

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/levigross/grequests"
	"github.com/neutral-af/backend/lib/airports"
	"github.com/neutral-af/backend/lib/config"
)

const baseURL = "https://api.flightstats.com/flex/schedules/rest/v1/json/flight"

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
	if date.Before(time.Now()) {
		return Details{}, fmt.Errorf("Cannot process flight for date because it is in the past: %s", date.Format("2006-01-02 15:04:05"))
	}

	airlineCode, flightCode, err := splitFlightNumber(flightNumber)
	if err != nil {
		return Details{}, err
	}
	path := fmt.Sprintf("/%s/%s/departing/%d/%d/%d", airlineCode, flightCode, date.Year(), date.Month(), date.Day())
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

	departureAirport, err := airports.GetFromIATA(responseData.ScheduledFlights[0].DepartureAirportFsCode)
	if err != nil {
		return Details{}, err
	}
	ArrivalAirport, err := airports.GetFromIATA(responseData.ScheduledFlights[0].ArrivalAirportFsCode)
	if err != nil {
		return Details{}, err
	}

	return Details{departureAirport.ICAO, ArrivalAirport.ICAO}, nil
}

type Response struct {
	ScheduledFlights []struct {
		DepartureAirportFsCode string `json:"departureAirportFsCode"`
		ArrivalAirportFsCode   string `json:"arrivalAirportFsCode"`
	} `json:"scheduledFlights"`
	Error struct {
		HTTPStatusCode int
		ErrorMessage   string
	}
}
