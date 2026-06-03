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
	return TrackEventWithStatus(key, username, trackName, quantity, createAt, http.StatusCreated)
}

func TrackEventWithStatus(key string, username, trackName string, quantity int, createAt *time.Time, expectedStatus int) error {
	var err error
	url := BASE_URL + fmt.Sprintf("/v1/events/%s/%s?quantity=%d", username, trackName, quantity)

	if key != "" {
		url += fmt.Sprintf("&key=%s", key)
	}

	if createAt != nil {
		url += fmt.Sprintf("&created_at=%s", url_escape.QueryEscape(createAt.Format(time.RFC3339)))
	}

	if err = utils.DoRequest(
		url,
		utils.WithRequestMethod(http.MethodPost),
		utils.WithExpectedStatusCode(expectedStatus),
	); err != nil {
		return fmt.Errorf("failed unable to register an event: %w", err)
	}

	return nil
}

func TrackEventAnonymous(username, trackName string, quantity int, createAt *time.Time) error {
	return TrackEventWithStatus("", username, trackName, quantity, createAt, http.StatusCreated)
}

func ListEvents(key string, username, trackName string, listBy *models.ListBy) (*models.List[models.Day], error) {
	return ListEventsWithStatus(key, username, trackName, listBy, http.StatusOK)
}

func ListEventsWithStatus(key string, username, trackName string, listBy *models.ListBy, expectedStatus int) (*models.List[models.Day], error) {
	var err error
	var response models.List[models.Day]

	url := BASE_URL + fmt.Sprintf("/v1/events/%s/%s", username, trackName)

	if key != "" {
		url += fmt.Sprintf("?key=%s", key)
	}

	if listBy != nil {
		if key != "" {
			url += fmt.Sprintf("&list_by=%s", *listBy)
		} else {
			url += fmt.Sprintf("?list_by=%s", *listBy)
		}
	}

	if err = utils.DoRequest(
		url,
		utils.WithRequestMethod(http.MethodGet),
		utils.WithExpectedStatusCode(expectedStatus),
		utils.ExtractGenericModel(&response),
	); err != nil {
		return nil, fmt.Errorf("failed unable to create api key request: %w", err)
	}

	return &response, nil
}

func ListEventsAnonymous(username, trackName string, listBy *models.ListBy) (*models.List[models.Day], error) {
	return ListEventsWithStatus("", username, trackName, listBy, http.StatusOK)
}
