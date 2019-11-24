package resolvers

import (
	"context"
	"errors"
	"math"

	"github.com/honeycombio/beeline-go"
	"github.com/neutral-af/backend/lib/distance"
	"github.com/neutral-af/backend/lib/emissions"
	models "github.com/neutral-af/backend/lib/graphql-models"
	providers "github.com/neutral-af/backend/lib/offset-providers"
	"github.com/neutral-af/backend/lib/offset-providers/cloverly"
)

var cloverlyAPI cloverly.Cloverly

func init() {
	cloverlyAPI = cloverly.New()
}

type getEstimateResolver struct{ *Resolver }

func (r *getEstimateResolver) FromFlights(ctx context.Context, get *models.GetEstimate, flights []*models.Flight, options *models.EstimateOptions) (*models.Estimate, error) {
	ctx, span := beeline.StartSpan(ctx, "fromFlights")
	defer span.Send()

	accumDistance := 0.0
	accumCarbon := 0.0

	for _, f := range flights {
		if f.Departure != nil && *f.Departure != "" && f.Arrival != nil && *f.Arrival != "" {
			if *f.Departure == *f.Arrival {
				return nil, errors.New("Departure and Arrival cannot be the same")
			}
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

	var provider providers.Provider
	if *options.Provider == models.ProviderCloverly {
		provider = &cloverlyAPI
	} else {
		return nil, errors.New("Provider unknown or not set")

	}
	estimate, err := provider.CreateCarbonEstimate(totalCarbon)
	if err != nil {
		return nil, err
	}

	beeline.AddField(ctx, "estimateID", estimate.ID)

	estimate.Km = &totalDistance

	return estimate, nil

}

func (r *getEstimateResolver) FromID(ctx context.Context, get *models.GetEstimate, id *string, provider *models.Provider) (*models.Estimate, error) {
	if *provider == models.ProviderCloverly {
		estimate, err := cloverlyAPI.RetrieveEstimate(*id)
		if err != nil {
			return nil, err
		}

		return estimate, nil
	}

	return nil, errors.New("Cannot retrieve estimate for given provider")
}
