package models

import "time"

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
	ID          int64           `json:"id" example:"1"`
	UserID      int64           `json:"user_id" example:"1"`
	Name        string          `json:"name" example:"default"`
	Description string          `json:"description" example:"default description"`
	Visibility  TrackVisibility `json:"visibility" example:"private"`
	Status      TrackStatus     `json:"status" example:"enabled"`
	CreatedAt   string          `json:"created_at" example:"2021-01-01T00:00:00Z"`
	DeleteAt    *time.Time      `json:"delete_at" example:"2021-01-01T00:00:00Z"`
}
