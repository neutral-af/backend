package currency

import (
	"errors"
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestRefreshRates(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.exchangeratesapi.io").
		Get("/latest").
		Reply(200).
		JSON(struct {
			Rates map[string]float64
		}{
			map[string]float64{
				"EUR": 1.3,
			},
		})

	results, err := refreshRates("USD")

	assert.NoError(t, err)
	assert.Contains(t, results, "EUR")

	assert.True(t, gock.IsDone())
}

func TestRefreshRatesError(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.exchangeratesapi.io").
		Get("/latest").
		Reply(500).
		SetError(errors.New("simulated error"))

	results, err := refreshRates("USD")

	assert.Error(t, err)
	assert.NotContains(t, results, "EUR")

	assert.True(t, gock.IsDone())
}

func TestRefreshRatesBadJSON(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.exchangeratesapi.io").
		Get("/latest").
		Reply(200).
		BodyString(`malformed json`)

	results, err := refreshRates("USD")

	assert.Error(t, err)
	assert.NotContains(t, results, "EUR")

	assert.True(t, gock.IsDone())
}

func TestExchangeCache(t *testing.T) {
	defer gock.Off()

	gock.New("*").Reply(500).SetError(errors.New("shouldn't have made this request!"))

	testExchanger := currencyExchanger{
		rates: map[string]rateCache{
			"USD": rateCache{
				values: map[string]float64{
					"EUR": 1.3,
				},
				expiry: time.Now().Add(time.Hour),
			},
		},
	}

	cents, err := testExchanger.exchange(200, "USD", "EUR")
	assert.NoError(t, err)
	assert.NotZero(t, cents)

	// Make sure the http request wasn't made (and intercepted)
	assert.True(t, gock.IsPending())
}

func TestExchangeUpdate(t *testing.T) {
	defer gock.Off()

	testExchanger := currencyExchanger{
		rates: map[string]rateCache{},
	}

	exchangeAmount := 200
	exchangeRate := 1.3

	gock.New("https://api.exchangeratesapi.io").
		Get("/latest").
		Reply(200).
		JSON(struct {
			Rates map[string]float64
		}{
			map[string]float64{
				"EUR": exchangeRate,
			},
		})

	cents, err := testExchanger.exchange(exchangeAmount, "USD", "EUR")
	assert.NoError(t, err)
	assert.Equal(t, int(float64(exchangeAmount)*exchangeRate), cents)

	assert.True(t, gock.IsDone())
}

func TestExchangeUpdateError(t *testing.T) {
	defer gock.Off()

	testExchanger := currencyExchanger{
		rates: map[string]rateCache{},
	}

	gock.New("https://api.exchangeratesapi.io").
		Get("/latest").
		Reply(200).
		SetError(errors.New("unexpected error"))

	cents, err := testExchanger.exchange(100, "USD", "EUR")
	assert.Error(t, err)
	assert.Zero(t, cents)

	assert.True(t, gock.IsDone())
}

func TestCovertNoop(t *testing.T) {
	defer gock.Off()
	gock.New("*").Reply(500).SetError(errors.New("shouldn't have made this request!"))

	cents, err := Convert(100, "USD", "USD")
	assert.NoError(t, err)
	assert.NotZero(t, cents)

	assert.True(t, gock.IsPending())
}

func TestConvert(t *testing.T) {
	defer gock.Off()
	gock.New("*").Reply(500).SetError(errors.New("shouldn't have made this request!"))

	exchanger = currencyExchanger{
		rates: map[string]rateCache{
			"USD": rateCache{
				values: map[string]float64{
					"EUR": 1.3,
				},
				expiry: time.Now().Add(time.Hour),
			},
		},
	}

	cents, err := Convert(200, "USD", "EUR")
	assert.NoError(t, err)
	assert.NotZero(t, cents)
	assert.Equal(t, 260, cents)

	assert.True(t, gock.IsPending())
}
