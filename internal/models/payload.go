package models

type Payload struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

type Entry struct {
	Id      string   `json:"id"`
	Changes []Change `json:"changes"`
}

type Change struct {
	Field string `json:"field"`
	Value Value  `json:"value"`
}

type Value struct {
	MessagingProduct string    `json:"messaging_product"`
	Metadata         Metadata  `json:"metadata"`
	Contacts         []Contact `json:"contacts,omitempty"`
	Messages         []Message `json:"messages,omitempty"`
	Statuses         []Status  `json:"statuses,omitempty"`
}

type Metadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberId      string `json:"phone_number_id"`
}

type Contact struct {
	WaId    string  `json:"wa_id" firestore:"-"`
	Profile Profile `json:"profile" firestore:"-"`
	Name    string  `json:"-" firestore:"name"`
}

type Profile struct {
	Name string `json:"name"`
}

type Status struct {
	Id          string `json:"id"`
	Status      string `json:"status"`
	Timestamp   string `json:"timestamp"`
	RecipientId string `json:"recipient_id"`
}

type Message struct {
	Id               string    `json:"id,omitempty" firestore:"-"`
	From             string    `json:"from,omitempty" firestore:"from"`
	To               string    `json:"to,omitempty" firestore:"to"`
	Timestamp        string    `json:"timestamp,omitempty" firestore:"timestamp"`
	Type             string    `json:"type" firestore:"type"`
	MessagingProduct string    `json:"messaging_product,omitempty" firestore:"messaging_product"`
	RecipientType    string    `json:"recipient_type,omitempty" firestore:"recipient_type"`
	Context          *Context  `json:"context,omitempty" firestore:"context,omitempty"`
	Referral         *Referral `json:"referral,omitempty" firestore:"referral,omitempty"`
	Text             *Text     `json:"text,omitempty" firestore:"text,omitempty"`
	Reaction         *Reaction `json:"reaction,omitempty" firestore:"reaction,omitempty"`
	Image            *Image    `json:"image,omitempty" firestore:"image,omitempty"`
	Sticker          *Sticker  `json:"sticker,omitempty" firestore:"sticker,omitempty"`
	Location         *Location `json:"location,omitempty" firestore:"location,omitempty"` // doesn't have a type field
	Button           *Button   `json:"button,omitempty" firestore:"button,omitempty"`
}

type Context struct {
	From string `json:"from"`
	Id   string `json:"id"`
}

type Referral struct {
	SourceUrl    string `json:"source_url"`
	SourceId     string `json:"source_id"`
	SourceType   string `json:"source_type"`
	Headline     string `json:"headline"`
	Body         string `json:"body"`
	MediaType    string `json:"media_type"`
	ImageUrl     string `json:"image_url"`
	VideoUrl     string `json:"video_url"`
	ThumbnailUrl string `json:"thumbnail_url"`
}

type Text struct {
	Body       string `json:"body" firestore:"body"`
	PreviewUrl bool   `json:"preview_url,omitempty" firestore:"preview_url"`
}

type Reaction struct {
	Emoji     string `json:"emoji" firestore:"emoji"`
	MessageId string `json:"message_id" firestore:"message_id"`
}

type Image struct {
	Caption  string `json:"caption,omitempty" firestore:"caption"`
	MimeType string `json:"mime_type,omitempty" firestore:"mime_type"`
	Sha256   string `json:"sha256,omitempty" firestore:"sha256"`
	Id       string `json:"id,omitempty" firestore:"id"`
	Link     string `json:"link,omitempty" firestore:"link"`
}

type Sticker struct {
	Id       string `json:"id,omitempty" firestore:"id"`
	Link     string `json:"link,omitempty" firestore:"link"`
	Animated bool   `json:"animated,omitempty" firestore:"animated"`
	MimeType string `json:"mime_type,omitempty" firestore:"mime_type"`
	Sha256   string `json:"sha256,omitempty" firestore:"sha256"`
}

type Location struct {
	Latitude  string `json:"latitude,omitempty" firestore:"latitude"`
	Longitude string `json:"longitude,omitempty" firestore:"longitude"`
	Name      string `json:"name,omitempty" firestore:"name"`
	Address   string `json:"address,omitempty" firestore:"address"`
}

type Button struct {
	Text    string `json:"text"`
	Payload string `json:"payload"`
}

type Audio struct {
	Id   string `json:"id,omitempty" firestore:"id"`
	Link string `json:"link,omitempty" firestore:"link"`
}

type Document struct {
	Caption  string `json:"caption,omitempty" firestore:"caption"`
	Link     string `json:"link,omitempty" firestore:"link"`
	Id       string `json:"id,omitempty" firestore:"id"`
	Filename string `json:"filename,omitempty" firestore:"filename"`
}

type Video struct {
	Caption string `json:"caption,omitempty" firestore:"caption"`
	Link    string `json:"link,omitempty" firestore:"link"`
	Id      string `json:"id,omitempty" firestore:"id"`
}

func (p *Payload) HasMessages() bool {
	return len(p.Entry[0].Changes[0].Value.Messages) > 0
}

func (p *Payload) HasStatuses() bool {
	return len(p.Entry[0].Changes[0].Value.Statuses) > 0
}
