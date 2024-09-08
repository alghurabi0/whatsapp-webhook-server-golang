package models

type Entry struct {
	Id      string   `json:"id"`
	Changes []Change `json:"changes"`
}
