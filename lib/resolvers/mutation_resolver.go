package resolvers

import (
	"context"

	models "github.com/jasongwartz/carbon-offset-backend/lib/graphql-models"
)

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) Purchase(ctx context.Context) (*models.MakePurchase, error) {
	return &models.MakePurchase{}, nil
}
