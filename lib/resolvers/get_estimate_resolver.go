package resolvers

import (
	"context"
	"encoding/json"
	"errors"
	"math"

	"github.com/honeycombio/beeline-go"
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
	ctx, span := beeline.StartSpan(ctx, "fromFlights")
	defer span.Send()

	accumDistance := 0.0
	accumCarbon := 0.0

	for _, f := range flights {
		if f.Departure != nil && *f.Departure != "" && f.Arrival != nil && *f.Arrival != "" {
			distance, err := distance.TwoAirports(*f.Departure, *f.Arrival)
			if err != nil {
				return nil, err
			}
			accumDistance += distance

			emissions := emissions.FlightCarbon(distance)
			accumCarbon += emissions
		} else if f.FlightNumber != nil && *f.FlightNumber != "" && f.Date != nil && *f.Date != "" {
			return nil, errors.New("Calculating from flight number not yet implemented")
		} else {
			return nil, errors.New("Invalid flight input: either (departure,arrival) or (flightNumber,date) must be provided")
		}
	}

	totalDistance := int(math.Round(accumDistance))
	totalCarbon := int(math.Round(accumCarbon))

	beeline.AddField(ctx, "provider", *options.Provider)
	beeline.AddField(ctx, "carbon", totalCarbon)

	if *options.Provider == models.ProviderCloverly {
		cloverlyEstimate, err := cloverlyAPI.CreateCarbonEstimate(totalCarbon)
		if err != nil {
			return nil, err
		}

		beeline.AddField(ctx, "estimateID", cloverlyEstimate.Slug)

		estimate, err := cloverlyToEstimate(cloverlyEstimate)
		if err != nil {
			return nil, err
		}
		estimate.Km = &totalDistance

		return estimate, nil
	}

	return nil, errors.New("Provider unknown or not set")
}

func (r *getEstimateResolver) FromID(ctx context.Context, get *models.GetEstimate, id *string, provider *models.Provider) (*models.Estimate, error) {
	if *provider == models.ProviderCloverly {
		estimate, err := cloverlyAPI.RetrieveEstimate(*id)
		if err != nil {
			return nil, err
		}

		return cloverlyToEstimate((estimate))
	}

	return nil, errors.New("Cannot retrieve estimate for given provider")
}

func cloverlyToEstimate(response cloverly.Response) (*models.Estimate, error) {
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
