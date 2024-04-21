package client

import (
	"fmt"
	"net/http"

	"512b.it/daytrack/src/models"
	"512b.it/daytrack/tests/utils"
)

func CreateTrack(token string, name string) error {
	var err error
	url := BASE_URL + "/v1/tracks"

	if err = utils.DoRequest(
		url,
		utils.WithRequestMethod(http.MethodPost),
		utils.WithExpectedStatusCode(http.StatusCreated),
		utils.WithAccessToken(token),
		utils.WithRequestBody(map[string]string{
			"name": name,
		}),
	); err != nil {
		return fmt.Errorf("failed unable to create track request: %w", err)
	}

	return nil
}

func UpdateTrack(token string, name string, visibility models.TrackVisibility) error {
	var err error
	url := BASE_URL + "/v1/tracks/" + name

	if err = utils.DoRequest(
		url,
		utils.WithRequestMethod(http.MethodPatch),
		utils.WithExpectedStatusCode(http.StatusOK),
		utils.WithAccessToken(token),
		utils.WithRequestBody(map[string]string{
			"visibility": string(visibility),
		}),
	); err != nil {
		return fmt.Errorf("failed unable to create track request: %w", err)
	}

	return nil
}

func ListTracks(token string) (*models.List[models.Track], error) {
	var err error
	var response models.List[models.Track]
	url := BASE_URL + "/v1/tracks"

	if err = utils.DoRequest(
		url,
		utils.WithRequestMethod(http.MethodGet),
		utils.WithExpectedStatusCode(http.StatusOK),
		utils.WithAccessToken(token),
		utils.ExtractGenericModel(&response),
	); err != nil {
		return nil, fmt.Errorf("failed unable to create get tracks request: %w", err)
	}

	return &response, nil
}
