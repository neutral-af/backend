package resolvers

import (
	"context"
	"time"

	"github.com/honeycombio/beeline-go"
	"github.com/neutral-af/backend/lib/config"
	models "github.com/neutral-af/backend/lib/graphql-models"
)

var startTime = time.Now()

func (r *queryResolver) Health(ctx context.Context) (*models.Health, error) {
	ctx, span := beeline.StartSpan(ctx, "health")
	defer span.Send()

	return &models.Health{
		AliveSince:  int(startTime.Unix()),
		Environment: config.C.Environment.ToString(),
	}, nil
}
