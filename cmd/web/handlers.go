package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/alghurabi0/whatsapp-webhook-server-golang/internal/models"
)

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

func (app *application) processPayload(w http.ResponseWriter, r *http.Request) {
	var payload models.Payload
	body, err := io.ReadAll(r.Body)
	if err != nil {
		app.errorLog.Printf("couldn't read body, error: %v\n", err)
		return
	}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		app.errorLog.Printf("failed to unmarshal json, error: %v\n", err)
		return
	}
	app.infoLog.Println(payload.Entry[0].Changes[0].Value.Messages[0].Text.Body)
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
	msg := &models.SendMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               phone,
		Type:             "text",
		Text: models.SendText{
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
