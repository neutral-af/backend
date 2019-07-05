package currency

import (
	"github.com/levigross/grequests"
)

// Define an interface for making the upstream call
// so that it can be mocked in tests
var fetcher func() (map[string]float64, error)

func init() {
	fetcher = func() (map[string]float64, error) {
		url := "https://api.exchangeratesapi.io/latest?base=USD"
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
}

func ConvertFromUSD(cents int, currency string) (int, error) {
	latest, err := fetcher()
	if err != nil {
		return 0, err
	}
	result := float64(cents) * latest[currency]
	rounded := int(result) // round off the fractions of a cent
	return rounded, nil
}
