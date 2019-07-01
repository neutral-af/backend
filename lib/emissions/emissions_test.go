package emissions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlightCarbon(t *testing.T) {
	carbon := FlightCarbon(5000.0)

	assert.NotZero(t, carbon)
	assert.Greater(t, carbon, 1000.0)
}
