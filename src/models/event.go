package models

import "time"

type ListBy string

const (
	ListByDay   ListBy = "day"
	ListByMonth ListBy = "month"
	ListByRaw   ListBy = "raw"
)

type Day struct {
	Date     time.Time `json:"date"`
	Quantity int       `json:"quantity"`
}
