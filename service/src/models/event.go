package models

import "time"

type ListBy string

const (
	ListByDay   ListBy = "day"
	ListByMonth ListBy = "month"
	ListByRaw   ListBy = "raw"
)

type Format string

const (
	FormatJSON Format = "json"
	FormatHTML Format = "html"
)

type Day struct {
	Date     time.Time `json:"date" example:"2021-01-01T00:00:00Z"`
	Quantity int       `json:"quantity" example:"1"`
}
