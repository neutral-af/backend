package distance

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTwoAirports(t *testing.T) {
	dist, err := TwoAirports("CYYZ", "EDDB")
	assert.NoError(t, err)
	assert.Greater(t, dist, 6000.0)
}

func TestAirportNotFound(t *testing.T) {
	dist, err := TwoAirports("invalid", "EDDB")
	assert.Error(t, err)
	assert.Equal(t, dist, 0.0)
}
