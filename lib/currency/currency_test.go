package currency

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCovertFromUSD(t *testing.T) {
	fetcher = func() (map[string]float64, error) {
		return map[string]float64{
			"USD": 1.0,
			"CAD": 1.5,
			"EUR": 0.8,
		}, nil
	}

	usd := 150
	result, err := ConvertFromUSD(usd, "USD")
	assert.NoError(t, err)
	assert.Equal(t, usd, result)

	// Test that euros are worth more than dollars
	result, err = ConvertFromUSD(usd, "EUR")
	assert.NoError(t, err)
	assert.Greater(t, usd, result)

	// Test that Canadian dollars are less than USD
	result, err = ConvertFromUSD(usd, "CAD")
	assert.NoError(t, err)
	assert.Less(t, usd, result)
}
