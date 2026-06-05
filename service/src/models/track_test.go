package models

import "testing"

func TestIsValidTrackNameValid(t *testing.T) {
	valid := []string{"abc", "my-track", "my_track", "MyTrack", "test123"}
	for _, name := range valid {
		if !IsValidTrackName(name) {
			t.Errorf("expected '%s' to be valid", name)
		}
	}
}

func TestIsValidTrackNameInvalid(t *testing.T) {
	invalid := []string{"", "a b", "track name", "my@track", "track.name", "hello world", "user/name"}
	for _, name := range invalid {
		if IsValidTrackName(name) {
			t.Errorf("expected '%s' to be invalid", name)
		}
	}
}
