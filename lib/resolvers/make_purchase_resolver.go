package resolvers

import (
	"context"
	"errors"

	"github.com/honeycombio/beeline-go"
	models "github.com/neutral-af/backend/lib/graphql-models"
	providers "github.com/neutral-af/backend/lib/offset-providers"
)

type makePurchaseResolver struct{ *Resolver }

func (r *makePurchaseResolver) FromEstimate(ctx context.Context, mp *models.MakePurchase, estimateID *string, selectedProvider *models.Provider) (*models.Purchase, error) {
	ctx, span := beeline.StartSpan(ctx, "fromEstimate")
	defer span.Send()

	beeline.AddField(ctx, "provider", *selectedProvider)

	var provider providers.Provider
	if *selectedProvider == models.ProviderCloverly {
		provider = &cloverlyAPI
	} else {
		return nil, errors.New("Provider unknown or not set")
	}

	return provider.Purchase(*estimateID)
}
