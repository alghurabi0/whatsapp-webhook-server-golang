package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/alghurabi0/whatsapp-webhook-server-golang/internal/models"
)

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

func (app *application) processPayload(w http.ResponseWriter, r *http.Request) {
	var payload models.Payload
	err := app.unmarshal(r, &payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// determine message or status
	info, err := app.validatePayload(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	contact := payload.Entry[0].Changes[0].Value.Contacts[0]
	_, err = app.contact.Get(ctx, contact.WaId)
	if err != nil {
		c := &models.Contact{
			WaId: contact.WaId,
			Name: contact.Profile.Name,
		}
		_, err = app.contact.Create(ctx, c)
		if err != nil {
			http.Error(w, fmt.Sprintf("couldn't create new contact, err: %v\n", err), http.StatusInternalServerError)
			app.errorLog.Println("couldn't create new contact")
			app.errorLog.Println(err)
			return
		}
	}

	if info == "msg" {
		msgType := payload.Entry[0].Changes[0].Value.Messages[0].Type
		// determine type or location, is there a referral?
		switch msgType {
		case "":
			// determine if there is a location
			break
		case "text":
			_, err = app.message.Create(ctx, &payload)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				app.errorLog.Println("couldn't create new messages")
				app.errorLog.Println(err)
				return
			}

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
	} else if info == "status" {
		// determine what to do with status
		app.infoLog.Println("Staaaaatus")
	} else {
		http.Error(w, "unexpected error", http.StatusBadRequest)
		app.errorLog.Println("unexpected error")
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
	msg := &models.SendMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               phone,
		Type:             "text",
		Template:         nil,
		Text: &models.Text{
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
