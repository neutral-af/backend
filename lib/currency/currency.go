package currency

import (
	"time"

	"github.com/levigross/grequests"
	models "github.com/neutral-af/backend/lib/graphql-models"
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

func (c *currencyExchanger) exchange(cents int, fromCurrency models.Currency, toCurrency models.Currency) (int, error) {
	from := fromCurrency.String()
	to := toCurrency.String()

	if c.rates[from].expiry.Before(time.Now()) {
		updatedRates, err := refreshRates(from)
		if err != nil {
			return 0, errors.Wrap(err, "Unable to refresh currency rates cache")
		}
		c.rates[from] = rateCache{
			values: updatedRates,
			expiry: time.Now().Add(ratesTTL),
		}
	}

	result := float64(cents) * c.rates[from].values[to]
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

func ConvertFromUSD(cents int, currency models.Currency) (int, error) {
	// Skip conversion if it's already in USD
	if currency == models.CurrencyUsd {
		return cents, nil
	}

	return exchanger.exchange(cents, "USD", currency)
}

func Convert(cents int, from models.Currency, to models.Currency) (int, error) {
	if from == to {
		return cents, nil
	}

	return exchanger.exchange(cents, from, to)
}
