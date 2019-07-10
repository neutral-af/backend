package resolvers

import (
	"context"

	models "github.com/jasongwartz/carbon-offset-backend/lib/graphql-models"
)

type queryResolver struct{ *Resolver }

func (r *queryResolver) Health(ctx context.Context) (bool, error) {
	return true, nil
}

func (r *queryResolver) Estimate(ctx context.Context) (*models.GetEstimate, error) {
	return &models.GetEstimate{}, nil
}
