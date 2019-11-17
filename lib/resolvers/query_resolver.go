package resolvers

import (
	"context"
	"time"

	"github.com/honeycombio/beeline-go"
	"github.com/neutral-af/backend/lib/config"
	models "github.com/neutral-af/backend/lib/graphql-models"
)

type queryResolver struct{ *Resolver }

var startTime = time.Now()

func (r *queryResolver) Health(ctx context.Context) (*models.Health, error) {
	ctx, span := beeline.StartSpan(ctx, "health")
	defer span.Send()

	return &models.Health{
		AliveSince:  int(startTime.Unix()),
		Environment: config.C.Environment.ToString(),
	}, nil
}

func (r *queryResolver) Estimate(ctx context.Context) (*models.GetEstimate, error) {
	return &models.GetEstimate{}, nil
}

func (r *queryResolver) Airport(ctx context.Context) (*models.GetAirport, error) {
	return &models.GetAirport{}, nil
}
