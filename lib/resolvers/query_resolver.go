package resolvers

import (
	"context"

	"github.com/honeycombio/beeline-go"
	models "github.com/neutral-af/backend/lib/graphql-models"
)

type queryResolver struct{ *Resolver }

func (r *queryResolver) Health(ctx context.Context) (bool, error) {
	ctx, span := beeline.StartSpan(ctx, "health")
	defer span.Send()

	return true, nil
}

func (r *queryResolver) Estimate(ctx context.Context) (*models.GetEstimate, error) {
	return &models.GetEstimate{}, nil
}
