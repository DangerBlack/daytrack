package tests

import (
	"net/http"
	"testing"
	"time"

	"512b.it/daytrack/src/models"
	"512b.it/daytrack/tests/client"
)

// --- Anonymous access ---

func TestAnonymousGetPublicReadWriteTrack(t *testing.T) {
	t.Parallel()

	user, err := client.CreateUser()
	if err != nil {
		t.Fatalf("unable to create user: %v", err)
	}
	trackName := "pub-rw-" + t.Name()

	if err = client.CreateTrack(user.Token, trackName); err != nil {
		t.Fatalf("unable to create track: %v", err)
	}
	if err = client.UpdateTrack(user.Token, trackName, models.TrackVisibilityPublicWrite); err != nil {
		t.Fatalf("unable to update track: %v", err)
	}

	events, err := client.ListEventsAnonymous(user.Username, trackName, nil)
	if err != nil {
		t.Fatalf("anonymous GET on public_rw track should succeed: %v", err)
	}
	if events == nil {
		t.Fatal("expected non-nil response")
	}
}

func TestAnonymousGetPublicReadOnlyTrack(t *testing.T) {
	t.Parallel()

	user, err := client.CreateUser()
	if err != nil {
		t.Fatalf("unable to create user: %v", err)
	}
	trackName := "pub-ro-" + t.Name()

	if err = client.CreateTrack(user.Token, trackName); err != nil {
		t.Fatalf("unable to create track: %v", err)
	}
	if err = client.UpdateTrack(user.Token, trackName, models.TrackVisibilityPublicRead); err != nil {
		t.Fatalf("unable to update track: %v", err)
	}

	events, err := client.ListEventsAnonymous(user.Username, trackName, nil)
	if err != nil {
		t.Fatalf("anonymous GET on public_ro track should succeed: %v", err)
	}
	if events == nil {
		t.Fatal("expected non-nil response")
	}
}

func TestAnonymousGetPrivateTrack(t *testing.T) {
	t.Parallel()

	user, err := client.CreateUser()
	if err != nil {
		t.Fatalf("unable to create user: %v", err)
	}
	trackName := "priv-" + t.Name()

	if err = client.CreateTrack(user.Token, trackName); err != nil {
		t.Fatalf("unable to create track: %v", err)
	}

	_, err = client.ListEventsWithStatus("", user.Username, trackName, nil, http.StatusForbidden)
	if err != nil {
		t.Fatalf("anonymous GET on private track should be 403: %v", err)
	}
}

func TestAnonymousPostPublicReadWriteTrack(t *testing.T) {
	t.Parallel()

	user, err := client.CreateUser()
	if err != nil {
		t.Fatalf("unable to create user: %v", err)
	}
	trackName := "pubrw-post-" + t.Name()

	if err = client.CreateTrack(user.Token, trackName); err != nil {
		t.Fatalf("unable to create track: %v", err)
	}
	if err = client.UpdateTrack(user.Token, trackName, models.TrackVisibilityPublicWrite); err != nil {
		t.Fatalf("unable to update track: %v", err)
	}

	if err = client.TrackEventWithStatus("", user.Username, trackName, 1, nil, http.StatusCreated); err != nil {
		t.Fatalf("anonymous POST on public_rw track should succeed: %v", err)
	}
}

func TestAnonymousPostPublicReadOnlyTrack(t *testing.T) {
	t.Parallel()

	user, err := client.CreateUser()
	if err != nil {
		t.Fatalf("unable to create user: %v", err)
	}
	trackName := "pubro-post-" + t.Name()

	if err = client.CreateTrack(user.Token, trackName); err != nil {
		t.Fatalf("unable to create track: %v", err)
	}
	if err = client.UpdateTrack(user.Token, trackName, models.TrackVisibilityPublicRead); err != nil {
		t.Fatalf("unable to update track: %v", err)
	}

	if err = client.TrackEventWithStatus("", user.Username, trackName, 1, nil, http.StatusForbidden); err != nil {
		t.Fatalf("anonymous POST on public_ro track should be 403: %v", err)
	}
}

func TestAnonymousPostPrivateTrack(t *testing.T) {
	t.Parallel()

	user, err := client.CreateUser()
	if err != nil {
		t.Fatalf("unable to create user: %v", err)
	}
	trackName := "priv-post-" + t.Name()

	if err = client.CreateTrack(user.Token, trackName); err != nil {
		t.Fatalf("unable to create track: %v", err)
	}

	if err = client.TrackEventWithStatus("", user.Username, trackName, 1, nil, http.StatusForbidden); err != nil {
		t.Fatalf("anonymous POST on private track should be 403: %v", err)
	}
}

// --- API key validation ---

func TestPostWithInvalidApiKey(t *testing.T) {
	t.Parallel()

	user, err := client.CreateUser()
	if err != nil {
		t.Fatalf("unable to create user: %v", err)
	}
	trackName := "invkey-" + t.Name()

	if err = client.CreateTrack(user.Token, trackName); err != nil {
		t.Fatalf("unable to create track: %v", err)
	}

	if err = client.TrackEventWithStatus("this-is-not-a-valid-key", user.Username, trackName, 1, nil, http.StatusUnauthorized); err != nil {
		t.Fatalf("POST with invalid key should be 401: %v", err)
	}
}

func TestListWithInvalidApiKey(t *testing.T) {
	t.Parallel()

	user, err := client.CreateUser()
	if err != nil {
		t.Fatalf("unable to create user: %v", err)
	}
	trackName := "invkey-list-" + t.Name()

	if err = client.CreateTrack(user.Token, trackName); err != nil {
		t.Fatalf("unable to create track: %v", err)
	}

	_, err = client.ListEventsWithStatus("this-is-not-a-valid-key", user.Username, trackName, nil, http.StatusUnauthorized)
	if err != nil {
		t.Fatalf("GET with invalid key should be 401: %v", err)
	}
}

func TestPostWithDeletedApiKey(t *testing.T) {
	t.Parallel()

	user, err := client.CreateUser()
	if err != nil {
		t.Fatalf("unable to create user: %v", err)
	}
	trackName := "delkey-" + t.Name()

	if err = client.CreateTrack(user.Token, trackName); err != nil {
		t.Fatalf("unable to create track: %v", err)
	}

	apiKey, err := client.CreateApiKey(user.Token, "delete-me")
	if err != nil {
		t.Fatalf("unable to create API key: %v", err)
	}

	if err = client.DeleteApiKey(user.Token, apiKey.ID); err != nil {
		t.Fatalf("unable to delete API key: %v", err)
	}

	if err = client.TrackEventWithStatus(apiKey.Key, user.Username, trackName, 1, nil, http.StatusUnauthorized); err != nil {
		t.Fatalf("POST with deleted key should be 401: %v", err)
	}
}

// --- Cross-access tests ---

func TestCrossAccessPrivateTrackPost(t *testing.T) {
	t.Parallel()

	owner, err := client.CreateUser()
	if err != nil {
		t.Fatalf("unable to create owner: %v", err)
	}
	intruder, err := client.CreateUser()
	if err != nil {
		t.Fatalf("unable to create intruder: %v", err)
	}
	trackName := "cross-priv-" + t.Name()

	if err = client.CreateTrack(owner.Token, trackName); err != nil {
		t.Fatalf("unable to create track: %v", err)
	}

	intruderKey, err := client.CreateApiKey(intruder.Token, "intruder")
	if err != nil {
		t.Fatalf("unable to create intruder API key: %v", err)
	}

	if err = client.TrackEventWithStatus(intruderKey.Key, owner.Username, trackName, 1, nil, http.StatusForbidden); err != nil {
		t.Fatalf("intruder POST on private track should be 403: %v", err)
	}
}

func TestCrossAccessPrivateTrackGet(t *testing.T) {
	t.Parallel()

	owner, err := client.CreateUser()
	if err != nil {
		t.Fatalf("unable to create owner: %v", err)
	}
	intruder, err := client.CreateUser()
	if err != nil {
		t.Fatalf("unable to create intruder: %v", err)
	}
	trackName := "cross-priv-get-" + t.Name()

	if err = client.CreateTrack(owner.Token, trackName); err != nil {
		t.Fatalf("unable to create track: %v", err)
	}

	intruderKey, err := client.CreateApiKey(intruder.Token, "intruder")
	if err != nil {
		t.Fatalf("unable to create intruder API key: %v", err)
	}

	_, err = client.ListEventsWithStatus(intruderKey.Key, owner.Username, trackName, nil, http.StatusForbidden)
	if err != nil {
		t.Fatalf("intruder GET on private track should be 403: %v", err)
	}
}

func TestCrossAccessPublicReadOnlyTrackGet(t *testing.T) {
	t.Parallel()

	owner, err := client.CreateUser()
	if err != nil {
		t.Fatalf("unable to create owner: %v", err)
	}
	reader, err := client.CreateUser()
	if err != nil {
		t.Fatalf("unable to create reader: %v", err)
	}
	trackName := "cross-ro-" + t.Name()

	if err = client.CreateTrack(owner.Token, trackName); err != nil {
		t.Fatalf("unable to create track: %v", err)
	}
	if err = client.UpdateTrack(owner.Token, trackName, models.TrackVisibilityPublicRead); err != nil {
		t.Fatalf("unable to update track: %v", err)
	}

	ownerKey, err := client.CreateApiKey(owner.Token, "owner")
	if err != nil {
		t.Fatalf("unable to create owner API key: %v", err)
	}

	now := time.Now()
	if err = client.TrackEvent(ownerKey.Key, owner.Username, trackName, 1, &now); err != nil {
		t.Fatalf("unable to seed event: %v", err)
	}

	readerKey, err := client.CreateApiKey(reader.Token, "reader")
	if err != nil {
		t.Fatalf("unable to create reader API key: %v", err)
	}

	events, err := client.ListEvents(readerKey.Key, owner.Username, trackName, nil)
	if err != nil {
		t.Fatalf("reader GET on public_ro track should succeed: %v", err)
	}
	if len(events.Items) < 1 {
		t.Fatalf("expected at least 1 event, got %d", len(events.Items))
	}
}

func TestCrossAccessPublicReadOnlyTrackPost(t *testing.T) {
	t.Parallel()

	owner, err := client.CreateUser()
	if err != nil {
		t.Fatalf("unable to create owner: %v", err)
	}
	writer, err := client.CreateUser()
	if err != nil {
		t.Fatalf("unable to create writer: %v", err)
	}
	trackName := "cross-ro-post-" + t.Name()

	if err = client.CreateTrack(owner.Token, trackName); err != nil {
		t.Fatalf("unable to create track: %v", err)
	}
	if err = client.UpdateTrack(owner.Token, trackName, models.TrackVisibilityPublicRead); err != nil {
		t.Fatalf("unable to update track: %v", err)
	}

	writerKey, err := client.CreateApiKey(writer.Token, "writer")
	if err != nil {
		t.Fatalf("unable to create writer API key: %v", err)
	}

	if err = client.TrackEventWithStatus(writerKey.Key, owner.Username, trackName, 1, nil, http.StatusForbidden); err != nil {
		t.Fatalf("writer POST on public_ro track should be 403: %v", err)
	}
}

func TestCrossAccessPublicReadWriteTrackGet(t *testing.T) {
	t.Parallel()

	owner, err := client.CreateUser()
	if err != nil {
		t.Fatalf("unable to create owner: %v", err)
	}
	reader, err := client.CreateUser()
	if err != nil {
		t.Fatalf("unable to create reader: %v", err)
	}
	trackName := "cross-rw-get-" + t.Name()

	if err = client.CreateTrack(owner.Token, trackName); err != nil {
		t.Fatalf("unable to create track: %v", err)
	}
	if err = client.UpdateTrack(owner.Token, trackName, models.TrackVisibilityPublicWrite); err != nil {
		t.Fatalf("unable to update track: %v", err)
	}

	ownerKey, err := client.CreateApiKey(owner.Token, "owner")
	if err != nil {
		t.Fatalf("unable to create owner API key: %v", err)
	}

	now := time.Now()
	if err = client.TrackEvent(ownerKey.Key, owner.Username, trackName, 1, &now); err != nil {
		t.Fatalf("unable to seed event: %v", err)
	}

	readerKey, err := client.CreateApiKey(reader.Token, "reader")
	if err != nil {
		t.Fatalf("unable to create reader API key: %v", err)
	}

	events, err := client.ListEvents(readerKey.Key, owner.Username, trackName, nil)
	if err != nil {
		t.Fatalf("reader GET on public_rw track should succeed: %v", err)
	}
	if len(events.Items) < 1 {
		t.Fatalf("expected at least 1 event, got %d", len(events.Items))
	}
}

func TestCrossAccessPublicReadWriteTrackPost(t *testing.T) {
	t.Parallel()

	owner, err := client.CreateUser()
	if err != nil {
		t.Fatalf("unable to create owner: %v", err)
	}
	writer, err := client.CreateUser()
	if err != nil {
		t.Fatalf("unable to create writer: %v", err)
	}
	trackName := "cross-rw-post-" + t.Name()

	if err = client.CreateTrack(owner.Token, trackName); err != nil {
		t.Fatalf("unable to create track: %v", err)
	}
	if err = client.UpdateTrack(owner.Token, trackName, models.TrackVisibilityPublicWrite); err != nil {
		t.Fatalf("unable to update track: %v", err)
	}

	writerKey, err := client.CreateApiKey(writer.Token, "writer")
	if err != nil {
		t.Fatalf("unable to create writer API key: %v", err)
	}

	if err = client.TrackEventWithStatus(writerKey.Key, owner.Username, trackName, 1, nil, http.StatusCreated); err != nil {
		t.Fatalf("writer POST on public_rw track should succeed: %v", err)
	}
}

// --- Owner access tests ---

func TestOwnerPostOnOwnPrivateTrack(t *testing.T) {
	t.Parallel()

	user, err := client.CreateUser()
	if err != nil {
		t.Fatalf("unable to create user: %v", err)
	}
	trackName := "own-priv-" + t.Name()

	if err = client.CreateTrack(user.Token, trackName); err != nil {
		t.Fatalf("unable to create track: %v", err)
	}

	apiKey, err := client.CreateApiKey(user.Token, "owner")
	if err != nil {
		t.Fatalf("unable to create API key: %v", err)
	}

	if err = client.TrackEventWithStatus(apiKey.Key, user.Username, trackName, 1, nil, http.StatusCreated); err != nil {
		t.Fatalf("owner POST on own private track should succeed: %v", err)
	}
}

// --- Anonymous event visibility ---

func TestAnonymousCreateAndListPublicReadWriteTrack(t *testing.T) {
	t.Parallel()

	user, err := client.CreateUser()
	if err != nil {
		t.Fatalf("unable to create user: %v", err)
	}
	trackName := "anon-rw-" + t.Name()

	if err = client.CreateTrack(user.Token, trackName); err != nil {
		t.Fatalf("unable to create track: %v", err)
	}
	if err = client.UpdateTrack(user.Token, trackName, models.TrackVisibilityPublicWrite); err != nil {
		t.Fatalf("unable to update track: %v", err)
	}

	// Anonymous creates an event
	if err = client.TrackEventWithStatus("", user.Username, trackName, 3, nil, http.StatusCreated); err != nil {
		t.Fatalf("anonymous POST on public_rw should succeed: %v", err)
	}

	// Anonymous lists events — should see what was created
	events, err := client.ListEventsAnonymous(user.Username, trackName, nil)
	if err != nil {
		t.Fatalf("anonymous GET on public_rw should succeed: %v", err)
	}
	if len(events.Items) < 1 {
		t.Fatalf("expected at least 1 event from anonymous list, got %d", len(events.Items))
	}
}

// --- Invalid track scenario ---

func TestEventOnNonexistentTrack(t *testing.T) {
	t.Parallel()

	user, err := client.CreateUser()
	if err != nil {
		t.Fatalf("unable to create user: %v", err)
	}

	apiKey, err := client.CreateApiKey(user.Token, "test")
	if err != nil {
		t.Fatalf("unable to create API key: %v", err)
	}

	if err = client.TrackEventWithStatus(apiKey.Key, user.Username, "this-track-does-not-exist", 1, nil, http.StatusForbidden); err != nil {
		t.Fatalf("POST on nonexistent track should be 403: %v", err)
	}
}
