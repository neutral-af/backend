package resolvers

import (
	"context"
	"encoding/json"

	models "github.com/jasongwartz/carbon-offset-backend/lib/graphql-models"
)

type makePurchaseResolver struct{ *Resolver }

func (r *makePurchaseResolver) FromEstimate(ctx context.Context, mp *models.MakePurchase, estimateID *string, provider *models.Provider) (*models.Purchase, error) {
	var resp *models.Purchase

	if *provider == models.ProviderCloverly {
		result, err := cloverlyAPI.Purchase(*estimateID)
		if err != nil {
			return nil, err
		}
		resp.Carbon = result.EquivalentCarbonInKG
		resp.ID = result.Slug

		detailsBytes, err := json.Marshal(result)
		details := string(detailsBytes)
		resp.Details = &details
	}

	return resp, nil
}
