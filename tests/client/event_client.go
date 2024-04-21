package client

import (
	"fmt"
	"net/http"
	url_escape "net/url"
	"time"

	"512b.it/daytrack/src/models"
	"512b.it/daytrack/tests/utils"
)

func TrackEvent(key string, username, trackName string, quantity int, createAt *time.Time) error {
	var err error
	url := BASE_URL + fmt.Sprintf("/v1/events/%s/%s?quantity=%d&key=%s", username, trackName, quantity, key)

	if createAt != nil {
		url += fmt.Sprintf("&created_at=%s", url_escape.QueryEscape(createAt.Format(time.RFC3339)))
	}

	if err = utils.DoRequest(
		url,
		utils.WithRequestMethod(http.MethodPost),
		utils.WithExpectedStatusCode(http.StatusCreated),
	); err != nil {
		return fmt.Errorf("failed unable to register an event: %w", err)
	}

	return nil
}

func ListEvents(key string, username, trackName string, listBy *models.ListBy) (*models.List[models.Day], error) {
	var err error
	var response models.List[models.Day]

	url := BASE_URL + fmt.Sprintf("/v1/events/%s/%s?key=%s", username, trackName, key)

	if listBy != nil {
		url += fmt.Sprintf("&list_by=%s", *listBy)
	}

	if err = utils.DoRequest(
		url,
		utils.WithRequestMethod(http.MethodGet),
		utils.WithExpectedStatusCode(http.StatusOK),
		utils.ExtractGenericModel(&response),
	); err != nil {
		return nil, fmt.Errorf("failed unable to create api key request: %w", err)
	}

	return &response, nil
}
