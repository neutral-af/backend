package resolvers

import (
	"context"

	models "github.com/neutral-af/backend/lib/graphql-models"
)

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) Payment(ctx context.Context) (*models.PaymentActions, error) {
	return &models.PaymentActions{}, nil
}
