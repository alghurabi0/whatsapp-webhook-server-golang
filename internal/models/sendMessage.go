package models

type SendMessage struct {
	MessagingProduct string    `json:"messaging_product"`
	RecipientType    string    `json:"recipient_type"`
	To               string    `json:"to"`
	Type             string    `json:"type"`
	Text             *SendText `json:"text,omitempty"`
	Template         *Template `json:"template,omitempty"`
}

type SendText struct {
	PreviewUrl bool   `json:"preview_url"`
	Body       string `json:"body"`
}

type Template struct {
	Name     string   `json:"name"`
	Language Language `json:"language"`
}

type Language struct {
	Code string `json:"code"`
}
