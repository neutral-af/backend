package cloverly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/levigross/grequests"
)

// Response matches the schema of an estimate or purchase response from Cloverly
type Response struct {
	Slug                      string  `json:"slug"`
	Environment               string  `json:"environment"`
	State                     string  `json:"state"`
	MicroRecCount             int     `json:"micro_rec_count"`
	MicroUnits                int     `json:"micro_units"`
	TotalCostInUSDCents       int     `json:"total_cost_in_usd_cents"`
	EstimatedAt               string  `json:"estimated_at"`
	EquivalentCarbonInKG      float64 `json:"equivalent_carbon_in_kg"`
	ElectricityInKWH          float64 `json:"electricity_in_kwh"`
	RecCostInUSDCents         int     `json:"rec_cost_in_usd_cents"`
	TransactionCostInUSDCents int     `json:"transaction_cost_in_usd_cents"`
	PrettyURL                 string  `json:"pretty_url"`
	Offset                    struct {
		Slug          string
		Name          string
		City          string
		Province      string
		Country       string
		OffsetType    string `json:"offset_type"`
		TotalCapacity string `json:"total_capacity"`
		LatLong       struct {
			X float64
			Y float64
		} `json:"latlng"`
		TechnicalDetails    string  `json:"technical_details"`
		AvailableCarbonInKG float64 `json:"available_carbon_in_kg"`
		PrettyURL           string  `json:"pretty_url"`
	} `json:"offset"`
	RenewableEnergyCertificate struct {
		Slug          string
		Name          string
		City          string
		Province      string
		Country       string
		RenewableType string `json:"renewable_type"`
		TotalCapacity string `json:"total_capacity"`
		LatLong       struct {
			X float64
			Y float64
		} `json:"latlng"`
		TechnicalDetails string `json:"technical_details"`
		Deprecated       map[string]string
	} `json:"renewable_energy_certificate"`
}

func createBody(carbonKilograms float64) (io.Reader, error) {
	var body struct {
		Weight struct {
			Value float64 `json:"value"`
			Units string  `json:"units"`
		} `json:"weight"`
	}

	body.Weight.Units = "kg"
	body.Weight.Value = carbonKilograms

	b, err := json.Marshal(body)
	if err != nil {
		return bytes.NewReader([]byte{}), err
	}

	return bytes.NewReader(b), nil
}

// Estimate creates a Cloverly estimate for the given volume of carbon
func Estimate(carbon float64) (Response, error) {
	url := "https://api.cloverly.app/2019-03-beta/estimates/carbon"

	data, err := createBody(carbon)
	if err != nil {
		return Response{}, err
	}
	key := os.Getenv("CLOVERLY_API_KEY")

	resp, err := grequests.Post(url, &grequests.RequestOptions{
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer private_key:%s", key),
		},
		RequestBody: data,
	})
	if err != nil {
		return Response{}, err
	}

	var responseData Response
	err = resp.JSON(&responseData)
	if err != nil {
		return Response{}, err
	}

	return responseData, nil
}
