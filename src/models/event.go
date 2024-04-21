package models

import "time"

type Day struct {
	Date     time.Time `json:"date"`
	Quantity int       `json:"quantity"`
}
