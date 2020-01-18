package cloverly

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestSimple(t *testing.T) {
	defer gock.Off()

	mockUrl := "https://mock.cloverly.com"

	gock.New(mockUrl).
		Post("/estimates.+").
		Reply(200).
		JSON(map[string]string{"slug": "test"})

	c := New()
	c.baseURL = mockUrl

	resp, err := c.CreateCarbonEstimate(6000)
	assert.NoError(t, err)
	assert.Equal(t, resp.ID, "test")

	assert.True(t, gock.IsDone())
}

func TestRetrieve(t *testing.T) {
	defer gock.Off()

	mockUrl := "https://mock.cloverly.com"

	gock.New(mockUrl).
		Get("/estimates.+").
		Reply(200).
		JSON(map[string]string{"slug": "test"})

	c := New()
	c.baseURL = mockUrl

	resp, err := c.RetrieveEstimate("test")
	assert.NoError(t, err)
	assert.Equal(t, resp.ID, "test")

	assert.True(t, gock.IsDone())
}

func TestRetrieveError(t *testing.T) {
	defer gock.Off()

	mockUrl := "https://mock.cloverly.com"

	gock.New(mockUrl).
		Get("/estimates.+").
		Reply(200).
		JSON(map[string]string{"error": "test"})

	c := New()
	c.baseURL = mockUrl

	_, err := c.RetrieveEstimate("test")
	assert.Error(t, err)

	assert.True(t, gock.IsDone())
}
