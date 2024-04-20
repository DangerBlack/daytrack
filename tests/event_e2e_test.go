package tests

import (
	"testing"

	"512b.it/daytrack/src/models"
	"512b.it/daytrack/tests/client"
)

func TestCreateEvent(t *testing.T) {
	t.Parallel()

	var err error
	var user *client.SignedUser
	var apiKey *models.ApiKey
	trackName := "water-plant"

	if user, err = client.CreateUser(); err != nil {
		t.Fatalf("unable to sign up %v", err)
		return
	}

	if apiKey, err = client.CreateApiKey(user.Token, "test"); err != nil {
		t.Fatalf("unable to create api key %v", err)
		return
	}

	if err = client.CreateTrack(user.Token, trackName); err != nil {
		t.Fatalf("unable to create api key %v", err)
		return
	}

	if err = client.TrackEvent(apiKey.Key, user.Username, trackName, 1); err != nil {
		t.Fatalf("unable to track event %v", err)
		return
	}
}
