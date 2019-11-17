package resolvers

import (
	"context"
	"errors"

	"github.com/honeycombio/beeline-go"
	models "github.com/neutral-af/backend/lib/graphql-models"
	providers "github.com/neutral-af/backend/lib/offset-providers"
	"github.com/neutral-af/backend/lib/payments"
)

type paymentActionsResolver struct{ *Resolver }

func (r *paymentActionsResolver) Checkout(ctx context.Context, pa *models.PaymentActions, estimate models.EstimateIn, paymentMethod string, amount int, currency models.Currency, opts *models.PaymentOptions) (*models.PaymentResponse, error) {
	ctx, span := beeline.StartSpan(ctx, "checkout")
	defer span.Send()

	response, err := payments.Checkout(paymentMethod, amount, currency, opts)
	if err != nil {
		return nil, err
	}

	return purchaseIfReady(ctx, response, estimate)
}

func (r *paymentActionsResolver) Confirm(ctx context.Context, pa *models.PaymentActions, estimate models.EstimateIn, paymentIntent string, opts *models.PaymentOptions) (*models.PaymentResponse, error) {
	ctx, span := beeline.StartSpan(ctx, "confirm")
	defer span.Send()

	response, err := payments.Confirm(paymentIntent, opts)
	if err != nil {
		return nil, err
	}

	return purchaseIfReady(ctx, response, estimate)
}

func purchaseIfReady(ctx context.Context, response *models.PaymentResponse, estimate models.EstimateIn) (*models.PaymentResponse, error) {
	ctx, span := beeline.StartSpan(ctx, "fromEstimate")
	beeline.AddField(ctx, "provider", estimate.Options.Provider)
	defer span.Send()

	if !*response.Success {
		return response, nil
	}

	var provider providers.Provider
	if *estimate.Options.Provider == models.ProviderCloverly {
		provider = &cloverlyAPI
	} else {
		return nil, errors.New("Provider unknown or not set")
	}

	purchase, err := provider.Purchase(*estimate.ID)
	if err != nil {
		// Risky error here! Money already taken from stripe, offset not yet confirmed
		return nil, err
	}

	response.Purchase = purchase
	return response, nil
}
