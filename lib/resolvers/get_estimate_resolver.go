package resolvers

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/honeycombio/beeline-go"
	"github.com/neutral-af/backend/lib/distance"
	"github.com/neutral-af/backend/lib/emissions"
	"github.com/neutral-af/backend/lib/flightstats"
	models "github.com/neutral-af/backend/lib/graphql-models"
	providers "github.com/neutral-af/backend/lib/offset-providers"
)

var flightStatsAPI flightstats.FlightStats

func init() {
	flightStatsAPI = flightstats.New()
}

type getEstimateResolver struct{ *Resolver }

func (r *getEstimateResolver) FromFlights(ctx context.Context, get *models.GetEstimate, flights []*models.Flight, options *models.EstimateOptions) (*models.Estimate, error) {
	ctx, span := beeline.StartSpan(ctx, "get estimate from flights")
	defer span.Send()

	accumDistance := 0.0
	accumCarbon := 0.0

	for _, f := range flights {
		var departure, arrival string
		if f.Departure != nil && *f.Departure != "" && f.Arrival != nil && *f.Arrival != "" {
			if *f.Departure == *f.Arrival {
				return nil, errors.New("Departure and Arrival cannot be the same")
			}
			departure = *f.Departure
			arrival = *f.Arrival
		} else if f.FlightNumber != nil && *f.FlightNumber != "" && f.Date != nil && *f.Date != "" {
			date, err := time.Parse(time.RFC3339, *f.Date)
			if err != nil {
				return nil, err
			}

			details, err := flightStatsAPI.GetAirportsForFlight(*f.FlightNumber, date)
			if err != nil {
				return nil, err
			}
			departure = details.Departure
			arrival = details.Arrival
		} else {
			return nil, errors.New("Invalid flight input: either (departure,arrival) or (flightNumber,date) must be provided")
		}

		distance, err := distance.TwoAirports(departure, arrival)
		if err != nil {
			return nil, err
		}
		accumDistance += distance

		emissions := emissions.FlightCarbon(distance)
		accumCarbon += emissions
	}

	totalDistance := int(math.Round(accumDistance))
	totalCarbon := int(math.Round(accumCarbon))

	beeline.AddField(ctx, "provider", *options.Provider)
	beeline.AddField(ctx, "carbon", totalCarbon)

	provider, err := providers.GetProviderAPI(*options.Provider)
	if err != nil {
		return nil, err
	}

	ctx, providerSpan := beeline.StartSpan(ctx, "call estimate provider")
	estimate, err := provider.CreateCarbonEstimate(totalCarbon)
	if err != nil {
		return nil, err
	}
	providerSpan.Send()

	beeline.AddField(ctx, "estimateID", estimate.ID)

	estimate.Km = &totalDistance

	return estimate, nil

}

func (r *getEstimateResolver) FromID(ctx context.Context, get *models.GetEstimate, id *string, provider *models.Provider) (*models.Estimate, error) {
	ctx, span := beeline.StartSpan(ctx, "retrieve estimate from ID")
	defer span.Send()

	p, err := providers.GetProviderAPI(*provider)
	if err != nil {
		return nil, err
	}

	estimate, err := p.RetrieveEstimate(*id)
	if err != nil {
		return nil, err
	}

	return estimate, nil
}
