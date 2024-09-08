package models

type SendMessage struct {
	MessagingProduct string   `json:"messaging_product"`
	RecipientType    string   `json:"recipient_type"`
	To               string   `json:"to"`
	Type             string   `json:"type"`
	Text             SendText `json:"text"`
}

type SendText struct {
	PreviewUrl bool   `json:"preview_url"`
	Body       string `json:"body"`
}
