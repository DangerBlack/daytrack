package tests

import (
	"testing"

	"512b.it/daytrack/src/models"
	"512b.it/daytrack/tests/client"
)

func TestCreateApiKey(t *testing.T) {
	t.Parallel()

	var err error
	var user *client.SignedUser
	var apiKey *models.ApiKey

	if user, err = client.CreateUser(); err != nil {
		t.Fatalf("unable to sign up %v", err)
		return
	}

	if apiKey, err = client.CreateApiKey(user.Token, "test"); err != nil {
		t.Fatalf("unable to create api key %v", err)
		return
	}

	if apiKey.Key == "" {
		t.Fatal("Failed unable to extract a valid api key")
		return
	}
}
