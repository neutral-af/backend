package resolvers

import (
	"context"

	models "github.com/jasongwartz/carbon-offset-backend/lib/graphql-models"
	"github.com/jasongwartz/carbon-offset-backend/lib/payments"
)

type paymentActionsResolver struct{ *Resolver }

func (r *paymentActionsResolver) Checkout(ctx context.Context, pa *models.PaymentActions, paymentMethod string, amount int, currency models.Currency, opts *models.PaymentOptions) (*models.PaymentResponse, error) {
	return payments.Checkout(paymentMethod, amount, currency, opts)
}

func (r *paymentActionsResolver) Confirm(ctx context.Context, pa *models.PaymentActions, paymentIntent string, opts *models.PaymentOptions) (*models.PaymentResponse, error) {
	return payments.Confirm(paymentIntent, opts)
}
