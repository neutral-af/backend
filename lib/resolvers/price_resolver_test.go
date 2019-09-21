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
			Breakdown: []*models.PriceElement{
				&models.PriceElement{
					Name:     "Main fee",
					Currency: models.CurrencyEur,
					Cents:    cents,
				},
			},
		},
	}

	r := estimateResolver{&Resolver{}}
	price, err := r.Price(context.TODO(), estimate, &currency)

	assert.NoError(t, err)

	assert.Equal(t, price.Currency, currency)
	assert.Equal(t, price.Breakdown[0].Currency, currency)
	assert.Equal(t, price.Breakdown[0].Cents, cents)
	assert.GreaterOrEqual(t, price.Cents, cents)
}

func TestPriceResolverConvert(t *testing.T) {
	cents := 100
	baseCurrency := models.CurrencyEur
	userCurrency := models.CurrencyCad

	estimate := &models.Estimate{
		Price: &models.Price{
			Breakdown: []*models.PriceElement{
				&models.PriceElement{
					Name:     "Main fee",
					Currency: baseCurrency,
					Cents:    cents,
				},
			},
		},
	}

	r := estimateResolver{&Resolver{}}
	price, err := r.Price(context.TODO(), estimate, &userCurrency)

	assert.NoError(t, err)

	assert.NotEqual(t, price.Cents, cents)
	assert.Equal(t, price.Currency, userCurrency)
	assert.NotEqual(t, price.Breakdown[0].Cents, cents)
	assert.NotEqual(t, price.Breakdown[0].Currency, baseCurrency)
	assert.Equal(t, price.Breakdown[0].Currency, userCurrency)

	// NOTE! The following assertion very dangerously assumes that
	// the EUR is worth more than CAD. If this changes,
	// the test will break.
	assert.Greater(t, price.Breakdown[0].Cents, cents)
}
