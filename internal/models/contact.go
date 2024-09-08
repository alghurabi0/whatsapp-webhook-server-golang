package models

type Contact struct {
	WaId    string  `json:"wa_id"`
	Profile Profile `json:"profile"`
}

type Profile struct {
	Name string `json:"name"`
}
