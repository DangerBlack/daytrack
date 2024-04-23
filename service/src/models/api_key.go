package models

import "time"

type ApiKey struct {
	ID        int64      `json:"id" example:"1"`
	UserID    int64      `json:"user_id" example:"1"`
	Name      string     `json:"name" example:"default"`
	Key       string     `json:"key" example:"a1b2c3d4e5f6g7h8i9j0"`
	CreatedAt time.Time  `json:"created_at" example:"2021-01-01T00:00:00Z"`
	DeleteAt  *time.Time `json:"delete_at" example:"2021-01-01T00:00:00Z"`
}
