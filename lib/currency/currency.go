package currency

import (
	"time"

	"github.com/jasongwartz/carbon-offset-backend/lib/schema/generated"
	"github.com/levigross/grequests"
	"github.com/pkg/errors"
)

const ratesTTL = time.Hour
const baseURL = "https://api.exchangeratesapi.io"

var exchanger = currencyExchanger{
	rates: map[string]rateCache{},
}

type rateCache struct {
	values map[string]float64
	expiry time.Time
}

type currencyExchanger struct {
	rates map[string]rateCache
}

func (c *currencyExchanger) exchange(cents int, fromCurrency string, toCurrency string) (int, error) {
	if c.rates[fromCurrency].expiry.Before(time.Now()) {
		updatedRates, err := refreshRates(fromCurrency)
		if err != nil {
			return 0, errors.Wrap(err, "Unable to refresh currency rates cache")
		}
		c.rates[fromCurrency] = rateCache{
			values: updatedRates,
			expiry: time.Now().Add(ratesTTL),
		}
	}

	result := float64(cents) * c.rates[fromCurrency].values[toCurrency]
	rounded := int(result) // round off the fractions of a cent
	return rounded, nil
}

func refreshRates(baseCurrency string) (map[string]float64, error) {
	url := baseURL + "/latest?base=" + baseCurrency
	resp, err := grequests.Get(url, &grequests.RequestOptions{})

	if err != nil {
		return nil, err
	}

	type Response struct {
		Rates map[string]float64
	}

	var rates Response
	if err := resp.JSON(&rates); err != nil {
		return nil, err
	}
	return rates.Rates, nil
}

func ConvertFromUSD(cents int, currency string) (int, error) {
	// Skip conversion if it's already in USD
	if currency == generated.CurrencyUsd.String() {
		return cents, nil
	}

	return exchanger.exchange(cents, "USD", currency)
}
