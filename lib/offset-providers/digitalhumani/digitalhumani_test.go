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
