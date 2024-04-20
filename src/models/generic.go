package models

type List[T any] struct {
	Items []T `json:"items"`
}
