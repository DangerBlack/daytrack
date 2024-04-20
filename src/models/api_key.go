package models

import "time"

type ApiKey struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	Name      string     `json:"name"`
	Key       string     `json:"key"`
	CreatedAt time.Time  `json:"created_at"`
	DeleteAt  *time.Time `json:"delete_at"`
}
