package schema

import (
	"context"
	"encoding/json"

	"github.com/jasongwartz/carbon-offset-backend/lib/schema/generated"
)

type makePurchaseResolver struct{ *Resolver }

func (r *makePurchaseResolver) FromEstimate(ctx context.Context, mp *generated.MakePurchase, estimateID *string, provider *generated.Provider) (*generated.Purchase, error) {
	var resp *generated.Purchase

	if *provider == generated.ProviderCloverly {
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
