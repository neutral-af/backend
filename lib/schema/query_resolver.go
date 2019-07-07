package schema

import (
	"context"

	"github.com/jasongwartz/carbon-offset-backend/lib/schema/generated"
)

type queryResolver struct{ *Resolver }

func (r *queryResolver) Health(ctx context.Context) (bool, error) {
	return true, nil
}

func (r *queryResolver) FlightEstimate(ctx context.Context) (*generated.FlightEstimate, error) {
	return &generated.FlightEstimate{}, nil
}
