package client

import (
	"fmt"
	"net/http"

	"512b.it/daytrack/tests/utils"
)

func TrackEvent(key string, username, trackName string, quantity int) error {
	var err error
	url := BASE_URL + fmt.Sprintf("/v1/events/%s/%s?quantity=%d&key=%s", username, trackName, quantity, key)

	if err = utils.DoRequest(
		url,
		utils.WithRequestMethod(http.MethodPost),
		utils.WithExpectedStatusCode(http.StatusCreated),
	); err != nil {
		return fmt.Errorf("failed unable to register an event: %w", err)
	}

	return nil
}
