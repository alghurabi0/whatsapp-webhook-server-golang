package models

type Payload struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}
