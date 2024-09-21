package models

type SendMessage struct {
	MessagingProduct string        `json:"messaging_product" firestore:"messaging_product"`
	RecipientType    string        `json:"recipient_type" firestore:"recipient_type"`
	To               string        `json:"to" firestore:"to"`
	Type             string        `json:"type" firestore:"type"`
	Context          *ReplyContext `json:"context,omitempty" firestore:"context"`
	Text             *Text         `json:"text,omitempty" firestore:"text"`
	Template         *Template     `json:"template,omitempty" firestore:"template"`
	Reaction         *Reaction     `json:"reaction,omitempty" firestore:"reaction"`
	Image            *Image        `json:"image,omitempty" firestore:"image"`
	Audio            *Audio        `json:"audio,omitempty" firestore:"audio"`
	Document         *Document     `json:"document,omitempty" firestore:"document"`
	Video            *Video        `json:"video,omitempty" firestore:"video"`
	Location         *Location     `json:"location,omitempty" firestore:"location"`
}

type ReplyContext struct {
	MessageId string `json:"message_id,omitempty" firestore:"message_id"`
}

type Template struct {
	Name       string      `json:"name,omitempty" firestore:"name"`
	Language   Language    `json:"language,omitempty" firestore:"language"`
	Components []Component `json:"components,omitempty" firestore:"components"`
}

type Language struct {
	Code string `json:"code,omitempty" firestore:"code"`
}

type Component struct {
	Type       string      `json:"type,omitempty" firestore:"type"`
	SubType    string      `json:"sub_type,omitempty" firestore:"sub_type"`
	Index      string      `json:"index,omitempty" firestore:"index"`
	Parameters []Parameter `json:"parameters,omitempty" firestore:"parameters"`
}

type Parameter struct {
	Type     string    `json:"type,omitempty" firestore:"type"`
	Image    *Image    `json:"image,omitempty" firestore:"image"`
	Text     string    `json:"text,omitempty" firestore:"text"`
	Currency *Currency `json:"currency,omitempty" firestore:"currency"`
	DateTime *DateTime `json:"date_time,omitempty" firestore:"date_time"`
	Payload  string    `json:"payload,omitempty" firestore:"payload"`
}

type Currency struct {
	FallbackValue string `json:"fallback_value,omitempty" firestore:"fallback_value"`
	Code          string `json:"code,omitempty" firestore:"code"`
	Amount1000    int    `json:"amount_1000,omitempty" firestore:"amount_1000"`
}

type DateTime struct {
	FallbackValue string `json:"fallback_value,omitempty" firestore:"fallback_value"`
	DayOfWeek     int    `json:"day_of_week,omitempty" firestore:"day_of_week"`
	Year          int    `json:"year,omitempty" firestore:"year"`
	Month         int    `json:"month,omitempty" firestore:"month"`
	DayOfMonth    int    `json:"day_of_month,omitempty" firestore:"day_of_month"`
	Hour          int    `json:"hour,omitempty" firestore:"hour"`
	Minute        int    `json:"minute,omitempty" firestore:"minute"`
	Calendar      string `json:"calendar,omitempty" firestore:"calendar"`
}

type Response struct {
	MessagingProduct string    `json:"messaging_product,omitempty"`
	Contacts         []ResCont `json:"contacts,omitempty"`
	Messages         []ResMsg  `json:"messages,omitempty"`
	Error            *Error    `json:"error,omitempty"`
}

type ResCont struct {
	Input string `json:"input,omitempty"`
	WaId  string `json:"wa_id,omitempty"`
}

type ResMsg struct {
	Id     string  `json:"id,omitempty"`
	Errors []Error `json:"errors"`
}

type Error struct {
	Title     string     `json:"title,omitempty"`
	Message   string     `json:"message,omitempty"`
	Type      string     `json:"type,omitempty"`
	Code      int        `json:"code,omitempty"`
	FbTraceId string     `json:"fbtrace_id,omitempty"`
	ErrorData *ErrorData `json:"error_data,omitempty"`
	Href      string     `json:"href,omitempty"`
}

type ErrorData struct {
	MessagingProduct string `json:"messaging_product,omitempty"`
	Details          string `json:"details,omitempty"`
}
