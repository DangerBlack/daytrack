package client

import (
	"fmt"
	"net/http"

	"512b.it/daytrack/src/models"
	"512b.it/daytrack/tests/utils"
)

func CreateApiKey(token string, name string) (*models.ApiKey, error) {
	var err error
	var response models.ApiKey
	url := BASE_URL + "/v1/api_keys?name=" + name

	if err = utils.DoRequest(
		url,
		utils.WithRequestMethod(http.MethodPost),
		utils.WithExpectedStatusCode(http.StatusCreated),
		utils.WithAccessToken(token),
		utils.ExtractGenericModel(&response),
	); err != nil {
		return nil, fmt.Errorf("failed unable to create get account request: %w", err)
	}

	return &response, nil
}

func ListApiKeys(token string) (*models.List[models.ApiKey], error) {
	var err error
	var response models.List[models.ApiKey]
	url := BASE_URL + "/v1/api_keys"

	if err = utils.DoRequest(
		url,
		utils.WithRequestMethod(http.MethodGet),
		utils.WithExpectedStatusCode(http.StatusOK),
		utils.WithAccessToken(token),
		utils.ExtractGenericModel(&response),
	); err != nil {
		return nil, fmt.Errorf("failed unable to create get account request: %w", err)
	}

	return &response, nil
}
