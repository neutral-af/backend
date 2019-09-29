package resolvers

import (
	"context"

	"github.com/honeycombio/beeline-go"
	"github.com/neutral-af/backend/lib/currency"
	models "github.com/neutral-af/backend/lib/graphql-models"
)

type estimateResolver struct{ *Resolver }

func (r *estimateResolver) Price(ctx context.Context, e *models.Estimate, inputCurrency *models.Currency) (*models.Price, error) {
	ctx, span := beeline.StartSpan(ctx, "calculatePrice")
	defer span.Send()

	userCurrency := *inputCurrency
	beeline.AddField(ctx, "currency", userCurrency)

	priceElements := e.Price.Breakdown
	priceElements = append(priceElements, &models.PriceElement{
		Name:     "Payment processing fee",
		Cents:    30,
		Currency: models.CurrencyUsd,
	})

	totalCents := 0
	for _, f := range priceElements {
		cents, err := currency.Convert(f.Cents, f.Currency, userCurrency)
		if err != nil {
			return nil, err
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

	e.Price.Breakdown = priceElements
	e.Price.Cents = totalCents
	e.Price.Currency = userCurrency

	return e.Price, nil
}
