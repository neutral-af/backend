package digitalhumani

import (
	"fmt"

	"github.com/levigross/grequests"

	"github.com/neutral-af/backend/lib/config"
	models "github.com/neutral-af/backend/lib/graphql-models"
	tree_carbon "github.com/neutral-af/backend/lib/tree-carbon"
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
	Message      string
}

func New() DigitalHumani {
	var baseURL string
	if config.C.Environment == config.EnvironmentProd {
		baseURL = "https://api.digitalhumani.com"
	} else {
		baseURL = "https://api-dev.digitalhumani.com"
	}

	return DigitalHumani{
		enterpriseID: config.C.DigitalHumaniEnterpriseID,
		baseURL:      baseURL,
		user:         "invoices@neutral.af",
	}
}

func (d *DigitalHumani) GetAllProjects() ([]Project, error) {
	resp, err := grequests.Get(d.baseURL+"/project", &grequests.RequestOptions{})
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
	var treeStr string
	if trees == 1 {
		treeStr = "tree"
	} else {
		treeStr = "trees"
	}

	return &models.Estimate{
		Carbon: &carbon,
		Price: &models.Price{
			Breakdown: []*models.PriceElement{
				&models.PriceElement{
					Name:     fmt.Sprintf("Your carbon offsets (%d %s)", trees, treeStr),
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

func (d *DigitalHumani) RetrieveEstimate(estimateID string) (*models.Estimate, error) {
	return nil, fmt.Errorf("Provider DigitalHumani does not implement retrieval")
}

func (d *DigitalHumani) Purchase(estimate models.EstimateIn) (*models.Purchase, error) {
	trees := tree_carbon.TreesForCarbonKG(*estimate.Carbon)

	resp, err := grequests.Post(d.baseURL+"/tree", &grequests.RequestOptions{
		Params: map[string]string{
			"enterpriseId": d.enterpriseID,
			"projectId":    "93333333",
			"user":         d.user,
			"treeCount":    string(trees),
		},
	})
	if err != nil {
		return nil, err
	}

	var responseData Tree
	err = resp.JSON(&responseData)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error in DigitalHumani call (status %d): %s", resp.StatusCode, responseData.Message)
	}

	purchase := &models.Purchase{}
	purchase.Carbon = tree_carbon.CarbonKGForTrees(responseData.TreeCount)
	purchase.ID = &responseData.UUID

	return purchase, nil
}
