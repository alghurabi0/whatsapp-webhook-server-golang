package models

type Message struct {
	Id        string `json:"id"`
	From      string `json:"from"`
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
	Text      Text   `json:"text"`
}

type Text struct {
	Body string `json:"body"`
}
