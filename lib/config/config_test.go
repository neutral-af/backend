package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDevFallback(t *testing.T) {
	os.Setenv("DEV_CLOVERLY_API_KEY", "testkey")
	c := New()

	assert.Equal(t, c.CloverlyAPIKey, "testkey")
}

func TestStaging(t *testing.T) {
	os.Setenv(branchEnvVar, "test")
	os.Setenv("CLOVERLY_API_KEY", "test0")
	os.Setenv("DEV_CLOVERLY_API_KEY", "test1")
	os.Setenv("STAGING_CLOVERLY_API_KEY", "test2")

	c := New()

	assert.Equal(t, c.CloverlyAPIKey, "test2")
}
