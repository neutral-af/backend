package resolvers

import (
	"context"

	models "github.com/neutral-af/backend/lib/graphql-models"
)

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) Purchase(ctx context.Context) (*models.MakePurchase, error) {
	return &models.MakePurchase{}, nil
}

func (r *mutationResolver) Payment(ctx context.Context) (*models.PaymentActions, error) {
	return &models.PaymentActions{}, nil
}
