package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	os.Setenv("CLOVERLY_API_KEY", "testkey")
	c := New()

	assert.Equal(t, c.CloverlyAPIKey, "testkey")
}
