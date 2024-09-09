package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/alghurabi0/whatsapp-webhook-server-golang/internal/models/WA"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.infoLog.Println("not home")
		return
	}
	app.infoLog.Println("welcome home")
	w.Write([]byte("welcome home"))
}

func (app *application) processPayload(w http.ResponseWriter, r *http.Request) {
	var payload WA.Payload
	err := app.unmarshal(r, &payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// determine message or status
	if payload.HasMessages() && !payload.HasStatuses() {
		msgType := payload.Entry[0].Changes[0].Value.Messages[0].Type
		// determine type or location, is there a referral?
		switch msgType {
		case "":
			break
		case "text":
			break
		case "reaction":
			break
		case "image":
			break
		case "sticker":
			break
		case "button":
			break
		default:
			break
		}
	} else if !payload.HasMessages() && payload.HasStatuses() {
		// determine what to do with status
	} else {
		http.Error(w, "payload doesn't contain messages or statuses", http.StatusBadRequest)
		app.errorLog.Println("payload doesn't contain messages of statuses")
		app.errorLog.Println(payload)
		return
	}
}

func (app *application) sendMessage(w http.ResponseWriter, r *http.Request) {
	token := os.Getenv("ACCESS_TOKEN")
	phone_number_id := os.Getenv("PHONE_NUMBER_ID")
	url := fmt.Sprintf("https://graph.facebook.com/v20.0/%s/messages", phone_number_id)
	if token == "" {
		app.errorLog.Fatal("empty ACCESS_TOKEN")
		return
	}
	phone := "9647802089950"
	text := "mew mew mew weew"
	msg := &WA.SendMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               phone,
		Type:             "text",
		Template:         nil,
		Text: &WA.SendText{
			PreviewUrl: false,
			Body:       text,
		},
	}
	jsonData, err := json.Marshal(msg)
	if err != nil {
		app.errorLog.Printf("couldn't marshal to json, error: %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorLog.Printf("couldn't get new req: %v\n", err)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		app.errorLog.Printf("failed to send a req: %v\n", err)
		return
	}

	defer res.Body.Close()
	app.infoLog.Printf("message status: %s\n", res.Status)
}

func (app *application) verifyHook(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("hub.mode")
	if mode == "" {
		app.errorLog.Println("empty mode")
		return
	}
	token := r.URL.Query().Get("hub.verify_token")
	if token == "" {
		app.errorLog.Println("empty token")
		return
	}
	challenge := r.URL.Query().Get("hub.challenge")
	if challenge == "" {
		app.errorLog.Println("empty challenge")
		return
	}
	my_token := os.Getenv("WEBHOOK_VERIFY_TOKEN")
	if mode == "subscribe" && token == my_token {
		w.Write([]byte(challenge))
		app.infoLog.Println("webhook verified successfully")
	}
}
