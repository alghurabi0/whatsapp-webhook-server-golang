package main

import (
	"context"
	"net/http"

	"github.com/alghurabi0/whatsapp-webhook-server-golang/internal/models"
)

func (app *application) processPayload(w http.ResponseWriter, r *http.Request) {
	var payload models.Payload
	err := app.unmarshalPayload(r, &payload)
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
	switch info {
	case "msg":
		err := app.getOrCreateContact(ctx, &payload)
		if err != nil {
			app.serverError(w, err)
			return
		}
		err = app.saveMessage(ctx, &payload)
		if err != nil {
			app.serverError(w, err)
			return
		}
	case "status":
		// determine what to do with status
		err := app.saveStatus(ctx, &payload)
		if err != nil {
			app.serverError(w, err)
			return
		}
		app.infoLog.Println("Staaaaatus")
	default:
		http.Error(w, "unexpected error", http.StatusBadRequest)
		app.errorLog.Println("unexpected error")
		app.errorLog.Println(payload)
		return
	}
}
