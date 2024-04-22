package tests

import (
	"testing"

	"512b.it/daytrack/src/models"
	"512b.it/daytrack/tests/client"
)

func TestCreateTrack(t *testing.T) {
	t.Parallel()

	var err error
	var user *client.SignedUser

	if user, err = client.CreateUser(); err != nil {
		t.Fatalf("unable to sign up %v", err)
		return
	}

	if err = client.CreateTrack(user.Token, "test"); err != nil {
		t.Fatalf("unable to create api key %v", err)
		return
	}
}

func TestCreateShortTrack(t *testing.T) {
	t.Parallel()

	var err error
	var user *client.SignedUser

	if user, err = client.CreateUser(); err != nil {
		t.Fatalf("unable to sign up %v", err)
		return
	}

	for _, name := range []string{"a", "aa", "a a"} {
		if err = client.CreateTrack(user.Token, name); err == nil {
			t.Fatalf("track name '%s' should have problem", name)
			return
		}
	}
}

func TestListTracks(t *testing.T) {
	t.Parallel()

	var err error
	var user *client.SignedUser
	var tracks *models.List[models.Track]

	if user, err = client.CreateUser(); err != nil {
		t.Fatalf("unable to sign up %v", err)
		return
	}

	if err = client.CreateTrack(user.Token, "test"); err != nil {
		t.Fatalf("unable to create api key %v", err)
		return
	}

	if err = client.CreateTrack(user.Token, "test-2"); err != nil {
		t.Fatalf("unable to create api key %v", err)
		return
	}

	if tracks, err = client.ListTracks(user.Token); err != nil {
		t.Fatalf("unable to list tracks %v", err)
		return
	}

	if len(tracks.Items) != 2 {
		t.Fatalf("expected 2 api keys, got %d", len(tracks.Items))
		return
	}

	for _, track := range tracks.Items {
		if track.Name == "" {
			t.Fatal("Failed unable to extract a valid track name")
			return
		}

		if track.Name != "test" && track.Name != "test-2" {
			t.Fatalf("expected track name to be test or test2, got %s", track.Name)
			return
		}
	}
}

func TestUpdateTrack(t *testing.T) {
	t.Parallel()

	var err error
	var user *client.SignedUser
	trackName := "water-plant"

	if user, err = client.CreateUser(); err != nil {
		t.Fatalf("unable to sign up %v", err)
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
}
