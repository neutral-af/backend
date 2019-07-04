package schema

import (
	"context"
	"encoding/json"

	"github.com/jasongwartz/carbon-offset-backend/lib/cloverly"
	"github.com/jasongwartz/carbon-offset-backend/lib/distance"
	"github.com/jasongwartz/carbon-offset-backend/lib/emissions"
	"github.com/jasongwartz/carbon-offset-backend/lib/schema/generated"
)

type Resolver struct{}

func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) FlightEstimate() generated.FlightEstimateResolver {
	return &flightEstimateResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Health(ctx context.Context) (bool, error) {
	return true, nil
}

func (r *queryResolver) FlightEstimate(ctx context.Context) (*generated.FlightEstimate, error) {
	return &generated.FlightEstimate{}, nil
}

type flightEstimateResolver struct{ *Resolver }

func (r *flightEstimateResolver) FromAirports(ctx context.Context, f *generated.FlightEstimate, departure string, arrival string) (*generated.EstimateResponse, error) {
	distance := distance.TwoAirports(departure, arrival)
	emissions := emissions.FlightCarbon(distance)
	estimate, err := cloverly.Estimate(emissions)
	if err != nil {
		return nil, err
	}

	// currency := "USD"
	detailsBytes, err := json.Marshal(estimate)
	details := string(detailsBytes)

	if err != nil {
		return nil, err
	}

	return &generated.EstimateResponse{
		Price:   &estimate.TotalCostInUSDCents,
		Carbon:  &estimate.EquivalentCarbonInKG,
		Details: &details,
	}, nil
}
