package resolvers

import (
	"context"

	models "github.com/neutral-af/backend/lib/graphql-models"
)

type queryResolver struct{ *Resolver }

func (r *queryResolver) Estimate(ctx context.Context) (*models.GetEstimate, error) {
	return &models.GetEstimate{}, nil
}

func (r *queryResolver) Airport(ctx context.Context) (*models.GetAirport, error) {
	return &models.GetAirport{}, nil
}
