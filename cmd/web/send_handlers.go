package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/alghurabi0/whatsapp-webhook-server-golang/internal/models"
)

func (app *application) sendMessage(w http.ResponseWriter, r *http.Request) {
	token := os.Getenv("ACCESS_TOKEN")
	if token == "" {
		app.errorLog.Fatal("empty ACCESS_TOKEN")
		return
	}
	phone_number_id := os.Getenv("PHONE_NUMBER_ID")
	if phone_number_id == "" {
		app.errorLog.Fatal("empty phone number id")
	}
	url := fmt.Sprintf("https://graph.facebook.com/v20.0/%s/messages", phone_number_id)
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}

	wa_id := r.FormValue("wa_id")
	if wa_id == "" {
		http.Error(w, "empty wa_id", http.StatusBadRequest)
		return
	}
	text := r.FormValue("message_content")
	if wa_id == "" {
		http.Error(w, "empty message_content", http.StatusBadRequest)
		return
	}
	msg := &models.Message{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               wa_id,
		Type:             "text",
		Text: &models.Text{
			PreviewUrl: false,
			Body:       text,
		},
		Context:  nil,
		Referral: nil,
		Reaction: nil,
		Image:    nil,
		Sticker:  nil,
		Location: nil,
		Button:   nil,
	}
	jsonData, err := json.Marshal(msg)
	if err != nil {
		app.errorLog.Printf("couldn't marshal to json, error: %v\n", err)
		return
	}
	app.infoLog.Println(string(jsonData))

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
	body, err := io.ReadAll(res.Body)
	if err != nil {
		app.errorLog.Fatal(err)
		return
	}
	app.infoLog.Println(string(body))
	result := &models.Response{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		app.errorLog.Printf("failed to unmarshal response: %v\n", err)
		return
	}
	app.infoLog.Printf("Response struct: %v\n", result)

	msg.Id = result.Messages[0].Id

	ctx := context.Background()
	_, err = app.message.Create(ctx, wa_id, msg)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (app *application) chat(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "empty name", http.StatusBadRequest)
		return
	}
	wa_id := r.FormValue("wa_id")
	if wa_id == "" {
		http.Error(w, "empty wa_id", http.StatusBadRequest)
		return
	}
	contact := &models.Contact{
		Name: name,
		WaId: wa_id,
	}

	data := app.newTemplateData(r)
	data.Contact = contact
	app.renderPart(w, http.StatusOK, "chat.tmpl.html", "chat", data)
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

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.infoLog.Println("not home")
		return
	}
	ctx := context.Background()
	contacts, err := app.contact.GetAll(ctx)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Contacts = contacts
	app.render(w, http.StatusOK, "home.tmpl.html", data)
}
