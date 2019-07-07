package schema

import (
	"context"

	"github.com/jasongwartz/carbon-offset-backend/lib/currency"
	"github.com/jasongwartz/carbon-offset-backend/lib/schema/generated"
)

type estimateResolver struct{ *Resolver }

func (r *estimateResolver) Price(ctx context.Context, e *generated.Estimate, inputCurrency *generated.Currency) (*generated.Price, error) {
	localCents, err := currency.ConvertFromUSD(*e.Price.Cents, string(*inputCurrency))
	if err != nil {
		return &generated.Price{}, err
	}

	e.Price.Cents = &localCents
	e.Price.Currency = inputCurrency

	return e.Price, nil
}
