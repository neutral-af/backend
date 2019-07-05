package schema

import (
	"context"
	"encoding/json"

	"github.com/jasongwartz/carbon-offset-backend/lib/cloverly"
	"github.com/jasongwartz/carbon-offset-backend/lib/currency"
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

func (r *Resolver) EstimateResponse() generated.EstimateResponseResolver {
	return &estimateResponseResolver{r}
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

	detailsBytes, err := json.Marshal(estimate)
	details := string(detailsBytes)

	if err != nil {
		return nil, err
	}

	return &generated.EstimateResponse{
		Price: &generated.Price{
			Cents: &estimate.TotalCostInUSDCents, // This value should get rewritten into a local currency below
		},
		Carbon:  &estimate.EquivalentCarbonInKG,
		Details: &details,
	}, nil
}

type estimateResponseResolver struct{ *Resolver }

func (r *estimateResponseResolver) Price(ctx context.Context, e *generated.EstimateResponse, inputCurrency *generated.Currency) (*generated.Price, error) {
	localCents, err := currency.ConvertFromUSD(*e.Price.Cents, string(*inputCurrency))
	if err != nil {
		return &generated.Price{}, err
	}

	centsPtr := &localCents
	priceOut := &generated.Price{
		Cents:    centsPtr,
		Currency: inputCurrency,
	}

	return priceOut, nil
}
