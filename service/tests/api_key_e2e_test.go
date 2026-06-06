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

func TestCreateWebClientApiKey(t *testing.T) {
	t.Parallel()

	var err error
	var user *client.SignedUser
	var apiKey *models.ApiKey

	if user, err = client.CreateUser(); err != nil {
		t.Fatalf("unable to sign up %v", err)
		return
	}

	// Create a web-client key — triggers DeleteOldWebClientKeys cleanup
	if apiKey, err = client.CreateApiKey(user.Token, "web-client"); err != nil {
		t.Fatalf("unable to create web-client api key %v", err)
		return
	}

	if apiKey.Key == "" {
		t.Fatal("Failed unable to extract a valid api key")
		return
	}

	// Create another web-client key — first one should NOT be deleted
	// since it was created seconds ago (cleanup threshold is JWT_ACCESS_TOKEN_DURATION, default 720h)
	var apiKey2 *models.ApiKey
	if apiKey2, err = client.CreateApiKey(user.Token, "web-client"); err != nil {
		t.Fatalf("unable to create second web-client api key %v", err)
		return
	}

	if apiKey2.Key == "" {
		t.Fatal("Failed unable to extract a valid second api key")
		return
	}

	// Both keys should still be active (neither is old enough to be cleaned up)
	var apiKeys *models.List[models.ApiKey]
	if apiKeys, err = client.ListApiKeys(user.Token); err != nil {
		t.Fatalf("unable to list api keys %v", err)
		return
	}

	if len(apiKeys.Items) != 2 {
		t.Fatalf("expected 2 web-client api keys, got %d", len(apiKeys.Items))
		return
	}
}

func TestListApiKey(t *testing.T) {
	t.Parallel()

	var err error
	var user *client.SignedUser
	var apiKey, apiKey2 *models.ApiKey
	var apiKeys *models.List[models.ApiKey]

	if user, err = client.CreateUser(); err != nil {
		t.Fatalf("unable to sign up %v", err)
		return
	}

	if apiKey, err = client.CreateApiKey(user.Token, "test"); err != nil {
		t.Fatalf("unable to create api key %v", err)
		return
	}

	if apiKey2, err = client.CreateApiKey(user.Token, "test2"); err != nil {
		t.Fatalf("unable to create api key %v", err)
		return
	}

	if apiKeys, err = client.ListApiKeys(user.Token); err != nil {
		t.Fatalf("unable to list api keys %v", err)
		return
	}

	if len(apiKeys.Items) != 2 {
		t.Fatalf("expected 2 api keys, got %d", len(apiKeys.Items))
		return
	}

	for _, key := range apiKeys.Items {
		if key.Key == "" {
			t.Fatal("Failed unable to extract a valid api key")
			return
		}
		if key.Name == "" {
			t.Fatal("Failed unable to extract a valid api key name")
			return
		}
		if key.Name != apiKey.Name && key.Name != apiKey2.Name {
			t.Fatalf("expected key name to be test or test2, got %s", key.Name)
			return
		}
	}
}

func TestDeleteApiKey(t *testing.T) {
	t.Parallel()

	var err error
	var user *client.SignedUser
	var apiKey *models.ApiKey
	var apiKeys *models.List[models.ApiKey]

	if user, err = client.CreateUser(); err != nil {
		t.Fatalf("unable to sign up %v", err)
		return
	}

	if apiKey, err = client.CreateApiKey(user.Token, "test"); err != nil {
		t.Fatalf("unable to create api key %v", err)
		return
	}

	if err = client.DeleteApiKey(user.Token, apiKey.ID); err != nil {
		t.Fatalf("unable to delete api key %v", err)
		return
	}

	if apiKeys, err = client.ListApiKeys(user.Token); err != nil {
		t.Fatalf("unable to list api keys %v", err)
		return
	}

	if len(apiKeys.Items) != 0 {
		t.Fatalf("expected 0 api keys, got %d", len(apiKeys.Items))
		return
	}
}
