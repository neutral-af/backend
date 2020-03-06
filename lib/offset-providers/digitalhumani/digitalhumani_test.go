package digitalhumani

import (
	"errors"
	"testing"

	"github.com/neutral-af/backend/lib/config"
	models "github.com/neutral-af/backend/lib/graphql-models"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

const mockURL = "https://mock.digitalhumani.com"

func TestURLEnvironment(t *testing.T) {
	testCases := []struct {
		env      config.Environment
		contains string
	}{
		{config.EnvironmentDev, "api-dev."},
		{config.EnvironmentStaging, "api-dev."},
		{config.EnvironmentProd, "api."},
	}

	for _, tc := range testCases {
		config.C.Environment = tc.env
		c := New()
		assert.Contains(t, c.baseURL, tc.contains)
	}
}

func TestGetAllProjects(t *testing.T) {
	defer gock.Off()

	gock.New(mockURL).
		Get("/.*project.*").
		Reply(200).
		JSON([]map[string]string{{"reforestationCompanyName_en": "test"}})

	c := New()
	c.baseURL = mockURL

	resp, err := c.GetAllProjects()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resp))
	assert.Equal(t, "test", resp[0].CompanyName)

	assert.True(t, gock.IsDone())
}

func TestCreateCarbonEstimate(t *testing.T) {
	defer gock.Off()

	gock.New(mockURL)

	c := New()
	c.baseURL = mockURL

	carbon := 2000
	resp, err := c.CreateCarbonEstimate(carbon)
	assert.NoError(t, err)
	assert.Equal(t, carbon/10, resp.Price.Breakdown[0].Cents)
	assert.Equal(t, carbon, *resp.Carbon)

	// Make sure no HTTP request was made/attempted
	assert.False(t, gock.HasUnmatchedRequest())
	assert.False(t, gock.IsDone())
}

func TestCreateCarbonEstimateString(t *testing.T) {
	defer gock.Off()

	gock.New(mockURL)

	c := New()
	c.baseURL = mockURL

	testCases := []struct {
		carbon   int
		contains string
	}{
		{1000, "tree)"},
		{2000, "trees)"},
	}

	for _, tc := range testCases {
		resp, err := c.CreateCarbonEstimate(tc.carbon)
		assert.NoError(t, err)

		assert.Contains(t, resp.Price.Breakdown[0].Name, tc.contains)
	}
}

func TestRetrieveEstimate(t *testing.T) {
	defer gock.Off()

	gock.New(mockURL)

	c := New()
	c.baseURL = mockURL

	resp, err := c.RetrieveEstimate("")
	assert.Error(t, err)
	assert.Nil(t, resp)

	// Make sure no HTTP request was made/attempted
	assert.False(t, gock.HasUnmatchedRequest())
	assert.False(t, gock.IsDone())
}

func TestPurchaseSuccess(t *testing.T) {
	defer gock.Off()

	c := New()
	c.baseURL = mockURL

	gock.New(mockURL).
		Post("/tree").
		Reply(200).
		JSON(map[string]interface{}{
			"uuid":      "1234",
			"projectId": "test_project",
			"user":      c.user,
			"treeCount": 2,
		})

	carbon := 2000
	estimate := models.EstimateIn{
		Carbon: &carbon,
	}

	resp, err := c.Purchase(estimate)
	assert.NoError(t, err)
	assert.Equal(t, carbon, resp.Carbon)
	assert.Equal(t, *resp.ID, "1234")
}

func TestPurchaseError(t *testing.T) {
	defer gock.Off()

	c := New()
	c.baseURL = mockURL

	gock.New(mockURL).
		Post("/tree").
		ReplyError(errors.New("failed request"))

	carbon := 2000
	estimate := models.EstimateIn{
		Carbon: &carbon,
	}

	resp, err := c.Purchase(estimate)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestPurchaseErrorStatus(t *testing.T) {
	defer gock.Off()

	c := New()
	c.baseURL = mockURL

	gock.New(mockURL).
		Post("/tree").
		Reply(500).
		BodyString(`{}`)

	carbon := 2000
	estimate := models.EstimateIn{
		Carbon: &carbon,
	}

	resp, err := c.Purchase(estimate)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestPurchaseInvalid(t *testing.T) {
	defer gock.Off()

	c := New()
	c.baseURL = mockURL

	gock.New(mockURL).
		Post("/tree").
		Reply(201).
		BodyString(`{ "invalid json": "" `)

	carbon := 2000
	estimate := models.EstimateIn{
		Carbon: &carbon,
	}

	resp, err := c.Purchase(estimate)
	assert.Error(t, err)
	assert.Nil(t, resp)
}
