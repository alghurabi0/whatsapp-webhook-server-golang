package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/alghurabi0/whatsapp-webhook-server-golang/internal/models"
)

func (app *application) sendMessage(w http.ResponseWriter, r *http.Request) {
	msg, err := app.prepareMessage(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		app.errorLog.Printf("couldn't marshal to json, error: %v\n", err)
		app.serverError(w, err)
		return
	}
	app.infoLog.Println("sending this json")
	app.infoLog.Println(string(jsonData))

	resp, code, err := app.sendMsgReq(jsonData)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if code != 200 {
		app.errorLog.Println("not a 200 code")
		app.errorLog.Println(code)
		app.serverError(w, err)
	}
	if len(resp.Messages[0].Errors) != 0 {
		app.errorLog.Println("got message errors")
		app.clientError(w, http.StatusBadRequest)
	}

	ctx := context.Background()
	wa_id := msg.To
	_, err = app.message.Create(ctx, wa_id, msg)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
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
