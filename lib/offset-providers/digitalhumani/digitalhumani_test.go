package digitalhumani

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestGetAllProjects(t *testing.T) {
	defer gock.Off()

	mockURL := "https://mock.digitalhumani.com"

	gock.New(mockURL).
		Get("/.*project.*").
		Reply(200).
		JSON([]map[string]string{{"reforestationCompanyName_en": "test"}})

	c := New()
	c.baseURL = mockURL

	resp, err := c.GetAllProjects()
	assert.NoError(t, err)
	assert.Equal(t, len(resp), 1)
	assert.Equal(t, resp[0].CompanyName, "test")

	assert.True(t, gock.IsDone())
}

// func TestCreateBodyFromMapEstimate(t *testing.T) {
// 	r, err := createBodyFromMap(map[string]interface{}{
// 		"weight": map[string]interface{}{
// 			"value": 1000,
// 			"units": "kg",
// 		},
// 	})
// 	assert.NoError(t, err)

// 	buf := new(bytes.Buffer)
// 	buf.ReadFrom(r)
// 	assert.Equal(t, `{"weight":{"units":"kg","value":1000}}`, buf.String())
// }

// func TestCreateBodyFromMapPurchase(t *testing.T) {
// 	r, err := createBodyFromMap(map[string]interface{}{
// 		"estimate_slug": "testtest",
// 	})
// 	assert.NoError(t, err)

// 	buf := new(bytes.Buffer)
// 	buf.ReadFrom(r)
// 	assert.Equal(t, buf.String(), `{"estimate_slug":"testtest"}`)
// }

// func TestSimple(t *testing.T) {
// 	defer gock.Off()

// 	mockUrl := "https://mock.cloverly.com"

// 	gock.New(mockUrl).
// 		Post("/estimates.+").
// 		Reply(200).
// 		JSON(map[string]string{"slug": "test"})

// 	c := New("mock_key")
// 	c.baseURL = mockUrl

// 	resp, err := c.CreateCarbonEstimate(6000)
// 	assert.NoError(t, err)
// 	assert.Equal(t, resp.Slug, "test")

// 	assert.True(t, gock.IsDone())
// }

// func TestRetrieve(t *testing.T) {
// 	defer gock.Off()

// 	mockUrl := "https://mock.cloverly.com"

// 	gock.New(mockUrl).
// 		Get("/estimates.+").
// 		Reply(200).
// 		JSON(map[string]string{"slug": "test"})

// 	c := New("mock_key")
// 	c.baseURL = mockUrl

// 	resp, err := c.RetrieveEstimate("test")
// 	assert.NoError(t, err)
// 	assert.Equal(t, resp.Slug, "test")

// 	assert.True(t, gock.IsDone())
// }

// func TestRetrieveError(t *testing.T) {
// 	defer gock.Off()

// 	mockUrl := "https://mock.cloverly.com"

// 	gock.New(mockUrl).
// 		Get("/estimates.+").
// 		Reply(200).
// 		JSON(map[string]string{"error": "test"})

// 	c := New("mock_key")
// 	c.baseURL = mockUrl

// 	_, err := c.RetrieveEstimate("test")
// 	assert.Error(t, err)

// 	assert.True(t, gock.IsDone())
// }
