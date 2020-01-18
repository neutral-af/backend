package utils

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapToJSONEstimate(t *testing.T) {
	r, err := MapToJSON(map[string]interface{}{
		"weight": map[string]interface{}{
			"value": 1000,
			"units": "kg",
		},
	})
	assert.NoError(t, err)

	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	assert.Equal(t, `{"weight":{"units":"kg","value":1000}}`, buf.String())
}

func TestMapToJSONPurchase(t *testing.T) {
	r, err := MapToJSON(map[string]interface{}{
		"estimate_slug": "testtest",
	})
	assert.NoError(t, err)

	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	assert.Equal(t, buf.String(), `{"estimate_slug":"testtest"}`)
}

func TestMapToJSONError(t *testing.T) {
	r, err := MapToJSON(map[string]interface{}{
		"example": make(chan int),
	})
	assert.Error(t, err)

	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	assert.Empty(t, buf.String())
}
