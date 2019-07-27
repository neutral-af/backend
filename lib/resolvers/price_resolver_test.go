package resolvers

import (
	"context"
	"testing"

	models "github.com/jasongwartz/carbon-offset-backend/lib/graphql-models"
	"github.com/stretchr/testify/assert"
)

func TestPriceResolverNoop(t *testing.T) {
	cents := 100
	currency := models.CurrencyEur

	estimate := &models.Estimate{
		Price: &models.Price{
			Cents:    cents,
			Currency: currency,
		},
	}

	r := estimateResolver{&Resolver{}}
	price, err := r.Price(context.TODO(), estimate, &currency)

	assert.NoError(t, err)

	var fees int
	for _, f := range price.Breakdown {
		fees += f.Cents
	}
	assert.Equal(t, cents, price.Cents-fees)
}

func TestPriceResolverConvert(t *testing.T) {
	// NOTE! This test very dangerously assumes that
	// the EUR is worth more than CAD. If this changes,
	// the test will break.
	cents := 100
	baseCurrency := models.CurrencyEur
	userCurrency := models.CurrencyCad

	estimate := &models.Estimate{
		Price: &models.Price{
			Cents:    cents,
			Currency: baseCurrency,
		},
	}

	r := estimateResolver{&Resolver{}}
	price, err := r.Price(context.TODO(), estimate, &userCurrency)

	assert.NoError(t, err)

	var fees int
	for _, f := range price.Breakdown {
		fees += f.Cents
	}
	assert.Greater(t, price.Cents-fees, cents)
}
