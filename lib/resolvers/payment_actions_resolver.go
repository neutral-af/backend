package resolvers

import (
	"context"

	"github.com/honeycombio/beeline-go"
	models "github.com/neutral-af/backend/lib/graphql-models"
	providers "github.com/neutral-af/backend/lib/offset-providers"
	"github.com/neutral-af/backend/lib/payments"
)

type paymentActionsResolver struct{ *Resolver }

func (r *paymentActionsResolver) Checkout(ctx context.Context, pa *models.PaymentActions, estimateIn models.EstimateIn, paymentMethod string, currency models.Currency, opts *models.PaymentOptions) (*models.PaymentResponse, error) {
	ctx, span := beeline.StartSpan(ctx, "checkout")
	defer span.Send()

	provider, err := providers.GetProviderAPI(*estimateIn.Options.Provider)
	if err != nil {
		return nil, err
	}

	estimate, err := provider.RetrieveEstimate(*estimateIn.ID)
	if err != nil {
		return nil, err
	}

	price, err := processTotalPrice(*estimate.Price, currency)
	if err != nil {
		return nil, err
	}

	estimate.Price = &price

	response, err := payments.Checkout(paymentMethod, estimate.Price.Cents, currency, opts)
	if err != nil {
		return nil, err
	}

	return purchaseIfReady(ctx, response, estimateIn)
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

	provider, err := providers.GetProviderAPI(*estimate.Options.Provider)
	if err != nil {
		return nil, err
	}

	purchase, err := provider.Purchase(estimate)
	if err != nil {
		// Risky error here! Money already taken from stripe, offset not yet confirmed
		return nil, err
	}

	response.Purchase = purchase
	return response, nil
}
