package digitalhumani

import (
	"fmt"

	"github.com/levigross/grequests"

	"github.com/neutral-af/backend/lib/config"
	models "github.com/neutral-af/backend/lib/graphql-models"
	tree_carbon "github.com/neutral-af/backend/lib/tree-carbon"
	"github.com/neutral-af/backend/lib/utils"
)

const centsPerTree = 100

type DigitalHumani struct {
	baseURL      string
	enterpriseID string
	user         string
}

type DigitalHumaniOpts struct {
}

// Project maps the response shema of /project
type Project struct {
	ID             string
	CompanyName    string `json:"reforestationCompanyName_en"`
	CompanyWebsite string `json:"reforestationCompanyWebsite_en"`
	Country        string `json:"reforestationProjectCountry_en"`
	Description    string `json:"reforestationProjectDescription_en"`
	ImageURL       string `json:"reforestationProjectImageURL_en"`
	State          string `json:"reforestationProjectState_en"`
	Website        string `json:"reforestationProjectWebsite_en"`
}

// Tree maps the response schema of /tree
type Tree struct {
	UUID         string
	TreeCount    int
	EnterpriseID string
	ProjectID    string
	User         string
}

func New() DigitalHumani {
	return DigitalHumani{
		enterpriseID: config.C.DigitalHumaniEnterpriseID,
		baseURL:      "https://3ib0d53ao8.execute-api.ca-central-1.amazonaws.com",
		user:         "invoices@neutral.af",
	}
}

func (d *DigitalHumani) GetAllProjects() ([]Project, error) {
	resp, err := grequests.Get(d.baseURL+"/dev/project", &grequests.RequestOptions{})
	if err != nil {
		return nil, err
	}

	var responseData []Project
	err = resp.JSON(&responseData)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}

// Estimate creates a DigitalHumani estimate for the given volume of carbon
func (d *DigitalHumani) CreateCarbonEstimate(carbon int) (*models.Estimate, error) {
	trees := tree_carbon.TreesForCarbonKG(carbon)

	return &models.Estimate{
		Carbon: &carbon,
		Price: &models.Price{
			Breakdown: []*models.PriceElement{
				&models.PriceElement{
					Name:     fmt.Sprintf("Your carbon offsets: %d trees!", trees),
					Cents:    trees * centsPerTree,
					Currency: models.CurrencyUsd,
				},
				&models.PriceElement{
					Name:     "DigitalHumani processing fee",
					Cents:    0,
					Currency: models.CurrencyUsd,
				},
			},
		},
	}, nil
}

func (d *DigitalHumani) Purchase(estimateID string) (*models.Purchase, error) {
	body, err := utils.CreateBodyFromMap(map[string]interface{}{
		// POST /tree {"enterpriseId":"11111111","projectId":"93333333","user":"email@test.com,","treeCount":2}
		"enterpriseId": d.enterpriseID,
		"projectId":    "93333333",
		"user":         d.user,
		"treeCount":    1,
	})
	resp, err := grequests.Post(d.baseURL+"/dev/tree", &grequests.RequestOptions{
		RequestBody: body,
	})
	if err != nil {
		return nil, err
	}

	var responseData Tree
	err = resp.JSON(&responseData)
	if err != nil {
		return nil, err
	}

	purchase := &models.Purchase{}
	purchase.Carbon = tree_carbon.CarbonKGForTrees(responseData.TreeCount)
	purchase.ID = responseData.UUID

	return purchase, nil
}

// func (c *DigitalHumani) Purchase(estimateID string) (Response, error) {
// 	path := "/purchases"

// 	data := map[string]interface{}{
// 		"estimate_slug": estimateID,
// 	}

// 	return c.postWithBody(path, data)
// }
