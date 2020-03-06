package providers

import (
	"errors"

	models "github.com/neutral-af/backend/lib/graphql-models"
	"github.com/neutral-af/backend/lib/offset-providers/cloverly"
	"github.com/neutral-af/backend/lib/offset-providers/digitalhumani"
)

// Provider is the interface that describes the behaviour a carbon
// offset provider is expected to implement.
type Provider interface {
	CreateCarbonEstimate(carbon int) (*models.Estimate, error)
	RetrieveEstimate(slug string) (*models.Estimate, error)
	Purchase(estimate models.EstimateIn) (*models.Purchase, error)
}

var cloverlyAPI cloverly.Cloverly
var digitalHumaniAPI digitalhumani.DigitalHumani

func init() {
	cloverlyAPI = cloverly.New()
	digitalHumaniAPI = digitalhumani.New()
}

// GetProviderAPI is used to get the active API object for
// the given provider model name
func GetProviderAPI(p models.Provider) (Provider, error) {
	var provider Provider
	switch p {
	case models.ProviderCloverly:
		provider = &cloverlyAPI
	case models.ProviderDigitalHumani:
		provider = &digitalHumaniAPI
	default:
		return nil, errors.New("Provider unknown or not set")
	}

	return provider, nil
}
