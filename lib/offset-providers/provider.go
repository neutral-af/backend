package providers

import (
	models "github.com/neutral-af/backend/lib/graphql-models"
)

// Provider is the interface that describes the behaviour a carbon
// offset provider is expected to implement.
type Provider interface {
	CreateCarbonEstimate(carbon int) (*models.Estimate, error)
	RetrieveEstimate(slug string) (*models.Estimate, error)
	Purchase(estimateID string) (*models.Purchase, error)
}
