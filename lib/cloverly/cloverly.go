package cloverly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"

	"github.com/levigross/grequests"
	"github.com/neutral-af/backend/lib/config"
	models "github.com/neutral-af/backend/lib/graphql-models"
	"github.com/pkg/errors"
)

type Cloverly struct {
	baseURL string
	apiKey  string
}

type CloverlyOpts struct {
}

func New() Cloverly {
	return Cloverly{
		apiKey:  config.C.CloverlyAPIKey,
		baseURL: "https://api.cloverly.app/2019-03-beta",
	}
}

func (c *Cloverly) get(path string) (Response, error) {
	resp, err := grequests.Get(c.baseURL+path, &grequests.RequestOptions{
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer private_key:%s", c.apiKey),
		},
	})
	if err != nil {
		return Response{}, err
	}

	var responseData Response
	err = resp.JSON(&responseData)
	if err != nil {
		return Response{}, err
	}

	if responseData.Error != "" {
		return Response{}, errors.Wrap(errors.New(responseData.Error), "Error in Cloverly response")
	}

	return responseData, nil
}

func (c *Cloverly) postWithBody(path string, body map[string]interface{}) (Response, error) {
	data, err := createBodyFromMap(body)
	if err != nil {
		return Response{}, err
	}

	resp, err := grequests.Post(c.baseURL+path, &grequests.RequestOptions{
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer private_key:%s", c.apiKey),
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

	if responseData.Error != "" {
		return Response{}, errors.Wrap(errors.New(responseData.Error), "Error in Cloverly response")
	}

	return responseData, nil
}

// Estimate creates a Cloverly estimate for the given volume of carbon
func (c *Cloverly) CreateCarbonEstimate(carbon int) (*models.Estimate, error) {
	path := "/estimates/carbon"

	data := map[string]interface{}{
		"weight": map[string]interface{}{
			"value": carbon,
			"units": "kg",
		},
	}

	response, err := c.postWithBody(path, data)
	if err != nil {
		return nil, err
	}

	return responseToEstimate(response)
}

func (c *Cloverly) RetrieveEstimate(slug string) (*models.Estimate, error) {
	path := fmt.Sprintf("/estimates/%s", slug)

	response, err := c.get((path))
	if err != nil {
		return nil, err
	}

	return responseToEstimate(response)
}

func (c *Cloverly) Purchase(estimateID string) (*models.Purchase, error) {
	path := "/purchases"

	data := map[string]interface{}{
		"estimate_slug": estimateID,
	}

	response, err := c.postWithBody(path, data)
	if err != nil {
		return nil, err
	}

	return responseToPurchase(response)
}

// Response matches the schema of an estimate or purchase response from Cloverly
type Response struct {
	Error                     string  `json:"error"`
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
		Deprecated       string
	} `json:"renewable_energy_certificate"`
}

func createBodyFromMap(data map[string]interface{}) (io.Reader, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return bytes.NewReader([]byte{}), err
	}

	return bytes.NewReader(b), nil
}

func responseToEstimate(response Response) (*models.Estimate, error) {
	provider := models.ProviderCloverly

	detailsBytes, err := json.Marshal(response)
	details := string(detailsBytes)
	if err != nil {
		return nil, err
	}

	totalCarbon := int(math.Round(response.EquivalentCarbonInKG))

	return &models.Estimate{
		ID: response.Slug,
		Price: &models.Price{
			Breakdown: []*models.PriceElement{
				&models.PriceElement{
					Name:     "Your carbon offsets contribution",
					Cents:    response.RecCostInUSDCents,
					Currency: models.CurrencyUsd,
				},
				&models.PriceElement{
					Name:     "Cloverly processing fee",
					Cents:    response.TransactionCostInUSDCents,
					Currency: models.CurrencyUsd,
				},
			},
		},
		Carbon:   &totalCarbon,
		Provider: &provider,
		Details:  &details,
	}, nil
}

func responseToPurchase(response Response) (*models.Purchase, error) {
	purchase := &models.Purchase{}
	purchase.Carbon = int(math.Round(response.EquivalentCarbonInKG))
	purchase.ID = response.Slug

	detailsBytes, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	details := string(detailsBytes)
	purchase.Details = &details
	return purchase, nil
}
