package schema

import (
	"context"

	"github.com/jasongwartz/carbon-offset-backend/lib/schema/generated"
)

type queryResolver struct{ *Resolver }

func (r *queryResolver) Health(ctx context.Context) (bool, error) {
	return true, nil
}

func (r *queryResolver) Estimate(ctx context.Context) (*generated.GetEstimate, error) {
	return &generated.GetEstimate{}, nil
}
