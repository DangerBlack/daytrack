package models

type TrackVisibility string
type TrackStatus string

const (
	TrackVisibilityPrivate     TrackVisibility = "private"
	TrackVisibilityPublicRead  TrackVisibility = "public_r"
	TrackVisibilityPublicWrite TrackVisibility = "public_rw"

	TrackStatusEnabled  TrackStatus = "enabled"
	TrackStatusDisabled TrackStatus = "disabled"
)

type Track struct {
	ID          int64           `json:"id"`
	UserID      int64           `json:"user_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Visibility  TrackVisibility `json:"visibility"`
	Status      TrackStatus     `json:"status"`
	CreatedAt   string          `json:"created_at"`
	DeleteAt    string          `json:"delete_at"`
}
