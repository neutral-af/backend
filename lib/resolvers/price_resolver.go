package resolvers

import (
	"context"

	"github.com/honeycombio/beeline-go"
	"github.com/neutral-af/backend/lib/currency"
	models "github.com/neutral-af/backend/lib/graphql-models"
)

type estimateResolver struct{ *Resolver }

func processTotalPrice(p models.Price, userCurrency models.Currency) (models.Price, error) {
	priceElements := p.Breakdown
	priceElements = append(priceElements, &models.PriceElement{
		Name:     "Payment processing fee",
		Cents:    30,
		Currency: models.CurrencyUsd,
	})

	totalCents := 0
	for _, f := range priceElements {
		cents, err := currency.Convert(f.Cents, f.Currency, userCurrency)
		if err != nil {
			return models.Price{}, err
		}
		f.Cents = cents
		f.Currency = userCurrency
		totalCents = totalCents + cents
	}

	commission := totalCents / 10
	totalCents += commission
	priceElements = append(priceElements, &models.PriceElement{
		Name:     "Our fee (10%)",
		Cents:    commission,
		Currency: userCurrency,
	})

	p.Breakdown = priceElements
	p.Cents = totalCents
	p.Currency = userCurrency

	return p, nil
}

func (r *estimateResolver) Price(ctx context.Context, e *models.Estimate, inputCurrency *models.Currency) (*models.Price, error) {
	ctx, span := beeline.StartSpan(ctx, "calculate price")
	defer span.Send()

	price, err := processTotalPrice(*e.Price, *inputCurrency)
	if err != nil {
		return nil, err
	}

	return &price, nil
}
