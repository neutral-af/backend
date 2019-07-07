package schema

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jasongwartz/carbon-offset-backend/lib/cloverly"
	"github.com/jasongwartz/carbon-offset-backend/lib/config"
	"github.com/jasongwartz/carbon-offset-backend/lib/distance"
	"github.com/jasongwartz/carbon-offset-backend/lib/emissions"
	"github.com/jasongwartz/carbon-offset-backend/lib/schema/generated"
)

type contextKey string

const contextEstimateKey contextKey = "estimate"

var cloverlyAPI cloverly.Cloverly

func init() {
	cloverlyAPI = cloverly.New(config.C.CloverlyAPIKey)
}

type flightEstimateResolver struct{ *Resolver }

func (r *flightEstimateResolver) FromAirports(ctx context.Context, f *generated.FlightEstimate, departure string, arrival string, options *generated.EstimateOptions) (*generated.Estimate, error) {
	distance := distance.TwoAirports(departure, arrival)
	emissions := emissions.FlightCarbon(distance)

	if *options.Provider == generated.ProviderCloverly {
		estimate, err := cloverlyAPI.Estimate(emissions)
		if err != nil {
			return nil, err
		}

		detailsBytes, err := json.Marshal(estimate)
		details := string(detailsBytes)
		if err != nil {
			return nil, err
		}

		return &generated.Estimate{
			Price: &generated.Price{
				Cents: &estimate.TotalCostInUSDCents, // This value should get rewritten into a local currency below
				Breakdown: []*generated.PriceElement{
					&generated.PriceElement{
						Name:  "Cloverly processing fee",
						Cents: &estimate.TransactionCostInUSDCents,
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
