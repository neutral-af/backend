package resolvers

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jasongwartz/carbon-offset-backend/lib/cloverly"
	"github.com/jasongwartz/carbon-offset-backend/lib/config"
	"github.com/jasongwartz/carbon-offset-backend/lib/distance"
	"github.com/jasongwartz/carbon-offset-backend/lib/emissions"
	models "github.com/jasongwartz/carbon-offset-backend/lib/graphql-models"
)

var cloverlyAPI cloverly.Cloverly

func init() {
	cloverlyAPI = cloverly.New(config.C.CloverlyAPIKey)
}

type getEstimateResolver struct{ *Resolver }

func (r *getEstimateResolver) FromFlights(ctx context.Context, get *models.GetEstimate, flights []*models.Flight, options *models.EstimateOptions) (*models.Estimate, error) {
	totalCarbon := 0.0

	for _, f := range flights {
		if f.Departure != nil && *f.Departure != "" && f.Arrival != nil && *f.Arrival != "" {
			distance := distance.TwoAirports(*f.Departure, *f.Arrival)
			emissions := emissions.FlightCarbon(distance)
			totalCarbon += emissions
		} else if f.FlightNumber != nil && *f.FlightNumber != "" && f.Date != nil && *f.Date != "" {
			return nil, errors.New("Calculating from flight number not yet implemented")
		} else {
			return nil, errors.New("Invalid flight input: either (departure,arrival) or (flightNumber,date) must be provided")
		}
	}

	if *options.Provider == models.ProviderCloverly {
		estimate, err := cloverlyAPI.Estimate(totalCarbon)
		if err != nil {
			return nil, err
		}

		detailsBytes, err := json.Marshal(estimate)
		details := string(detailsBytes)
		if err != nil {
			return nil, err
		}

		return &models.Estimate{
			ID: estimate.Slug,
			Price: &models.Price{
				Cents:    estimate.TotalCostInUSDCents, // This value should get rewritten into a local currency below
				Currency: models.CurrencyUsd,           // This should also get rewritten (based on user-selected currency)
				Breakdown: []*models.PriceElement{
					&models.PriceElement{
						Name:     "Cloverly processing fee",
						Cents:    estimate.TransactionCostInUSDCents,
						Currency: models.CurrencyUsd,
					},
				},
			},
			Carbon:   &estimate.EquivalentCarbonInKG,
			Provider: options.Provider,
			Details:  &details,
		}, nil
	}

	return nil, errors.New("Provider not set")

}
