package models

import "time"

type ApiKey struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	Key       string     `json:"key"`
	CreatedAt time.Time  `json:"created_at"`
	DeleteAt  *time.Time `json:"delete_at"`
}
