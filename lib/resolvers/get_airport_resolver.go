package resolvers

import (
	"context"

	"github.com/neutral-af/backend/lib/airports"
	models "github.com/neutral-af/backend/lib/graphql-models"
)

type getAirportResolver struct{ *Resolver }

func (r *getAirportResolver) Search(ctx context.Context, get *models.GetAirport, query string) ([]*models.Airport, error) {
	matches := airports.Search(query)
	var matchesModels []*models.Airport

	for _, m := range matches {
		model := m.ToModel()
		matchesModels = append(matchesModels, &model)
	}

	return matchesModels, nil
}

func (r *getAirportResolver) FromIcao(ctx context.Context, get *models.GetAirport, code string) (*models.Airport, error) {
	a, err := airports.GetFromICAO(code)
	if err != nil {
		return nil, err
	}

	m := a.ToModel()
	return &m, nil

}
