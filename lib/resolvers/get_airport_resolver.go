package resolvers

import (
	"context"

	"github.com/sahilm/fuzzy"

	"github.com/neutral-af/backend/lib/airports"
	models "github.com/neutral-af/backend/lib/graphql-models"
)

type getAirportResolver struct{ *Resolver }

func (r *getAirportResolver) FuzzySearch(ctx context.Context, get *models.GetAirport, query string) ([]*models.Airport, error) {
	airports := airports.GetAll()
	results := fuzzy.FindFrom(query, airports)

	matches := []*models.Airport{}
	for _, r := range results {
		if len(matches) > 10 {
			break
		}
		a := airports[r.Index].ToModel()
		matches = append(matches, &a)
	}

	return matches, nil
}

func (r *getAirportResolver) FromIcao(ctx context.Context, get *models.GetAirport, code string) (*models.Airport, error) {
	a, err := airports.GetFromICAO(code)
	if err != nil {
		return nil, err
	}

	m := a.ToModel()
	return &m, nil

}
