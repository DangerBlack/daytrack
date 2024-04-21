package tests

import (
	"testing"
	"time"

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

	if err = client.TrackEvent(apiKey.Key, user.Username, trackName, 1, nil); err != nil {
		t.Fatalf("unable to track event %v", err)
		return
	}
}

func TestCreateEventAtTime(t *testing.T) {
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

	now := time.Now()

	if _, err := time.Parse(time.RFC3339, now.Format(time.RFC3339)); err != nil {
		t.Fatalf("unable to parse time %v", err)
		return
	}

	if err = client.TrackEvent(apiKey.Key, user.Username, trackName, 1, &now); err != nil {
		t.Fatalf("unable to track event %v", err)
		return
	}
}

func TestListEventsByDay(t *testing.T) {
	t.Parallel()

	var err error
	var user *client.SignedUser
	var apiKey *models.ApiKey
	var events *models.List[models.Day]
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

	t1, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	t2, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	t3, _ := time.Parse(time.RFC3339, "2021-01-02T00:00:00Z")

	for _, c := range []time.Time{t1, t2, t3} {
		if err = client.TrackEvent(apiKey.Key, user.Username, trackName, 1, &c); err != nil {
			t.Fatalf("unable to track event %v", err)
			return
		}
	}

	if events, err = client.ListEvents(apiKey.Key, user.Username, trackName, nil); err != nil {
		t.Fatalf("unable to list events %v", err)
		return
	}

	if len(events.Items) != 2 {
		t.Fatalf("expected 1 event, got %d", len(events.Items))
		return
	}

	listBy := models.ListByRaw
	if events, err = client.ListEvents(apiKey.Key, user.Username, trackName, &listBy); err != nil {
		t.Fatalf("unable to list events %v", err)
		return
	}

	if len(events.Items) != 3 {
		t.Fatalf("expected 1 event, got %d", len(events.Items))
		return
	}
}

func TestListPublicEventsByDay(t *testing.T) {
	t.Parallel()

	var err error
	var user, user2 *client.SignedUser
	var apiKey, apiKey2 *models.ApiKey
	var events *models.List[models.Day]
	trackName := "water-plant"

	if user, err = client.CreateUser(); err != nil {
		t.Fatalf("unable to sign up %v", err)
		return
	}

	if apiKey, err = client.CreateApiKey(user.Token, "test"); err != nil {
		t.Fatalf("unable to create api key %v", err)
		return
	}

	if user2, err = client.CreateUser(); err != nil {
		t.Fatalf("unable to sign up %v", err)
		return
	}

	if apiKey2, err = client.CreateApiKey(user2.Token, "test"); err != nil {
		t.Fatalf("unable to create api key %v", err)
		return
	}

	if err = client.CreateTrack(user.Token, trackName); err != nil {
		t.Fatalf("unable to create api key %v", err)
		return
	}

	if err = client.UpdateTrack(user.Token, trackName, models.TrackVisibilityPublicRead); err != nil {
		t.Fatalf("unable to update track %v", err)
		return
	}

	var tracks *models.List[models.Track]
	if tracks, err = client.ListTracks(user.Token); err != nil {
		t.Fatalf("unable to get track %v", err)
		return
	}

	if len(tracks.Items) != 1 {
		t.Fatalf("expected 1 track, got %d", len(tracks.Items))
		return
	}

	if tracks.Items[0].Visibility != models.TrackVisibilityPublicRead {
		t.Fatalf("expected track to be public read, got %s", tracks.Items[0].Visibility)
		return
	}

	t1, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	t2, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	t3, _ := time.Parse(time.RFC3339, "2021-01-02T00:00:00Z")

	for _, c := range []time.Time{t1, t2, t3} {
		if err = client.TrackEvent(apiKey.Key, user.Username, trackName, 1, &c); err != nil {
			t.Fatalf("unable to track event %v", err)
			return
		}
	}

	if events, err = client.ListEvents(apiKey2.Key, user.Username, trackName, nil); err != nil {
		t.Fatalf("unable to list events %v", err)
		return
	}

	if len(events.Items) != 2 {
		t.Fatalf("expected 1 event, got %d", len(events.Items))
		return
	}

	listBy := models.ListByRaw
	if events, err = client.ListEvents(apiKey2.Key, user.Username, trackName, &listBy); err != nil {
		t.Fatalf("unable to list events %v", err)
		return
	}

	if len(events.Items) != 3 {
		t.Fatalf("expected 1 event, got %d", len(events.Items))
		return
	}
}
