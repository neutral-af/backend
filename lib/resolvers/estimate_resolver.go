package resolvers

import (
	"context"

	"github.com/jasongwartz/carbon-offset-backend/lib/currency"
	models "github.com/jasongwartz/carbon-offset-backend/lib/graphql-models"
)

type estimateResolver struct{ *Resolver }

func (r *estimateResolver) Price(ctx context.Context, e *models.Estimate, inputCurrency *models.Currency) (*models.Price, error) {
	userCurrency := *inputCurrency

	fees := []*models.PriceElement{
		&models.PriceElement{
			Name:     "Stripe processing fee (30 cents USD)",
			Cents:    30,
			Currency: models.CurrencyUsd,
		},
		&models.PriceElement{
			Name:     "Our fee (10%)",
			Cents:    e.Price.Cents / 10,
			Currency: e.Price.Currency,
		},
	}

	localCents, err := currency.Convert(e.Price.Cents, e.Price.Currency, userCurrency)
	if err != nil {
		return nil, err
	}

	totalCents := localCents

	for _, f := range fees {
		cents, err := currency.Convert(f.Cents, f.Currency, userCurrency)
		if err != nil {
			return nil, err
		}
		f.Cents = cents
		f.Currency = userCurrency
		totalCents = totalCents + cents
	}

	e.Price.Breakdown = append(e.Price.Breakdown, fees...)
	e.Price.Cents = totalCents
	e.Price.Currency = userCurrency

	return e.Price, nil
}
