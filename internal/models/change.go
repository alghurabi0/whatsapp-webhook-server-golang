package models

type Change struct {
	Field string `json:"field"`
	Value Value  `json:"value"`
}
