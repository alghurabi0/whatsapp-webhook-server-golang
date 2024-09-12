package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/alghurabi0/whatsapp-webhook-server-golang/internal/models"
)

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

	if info == "msg" {
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
		msgType := payload.Entry[0].Changes[0].Value.Messages[0].Type
		// determine type or location, is there a referral?
		switch msgType {
		case "":
			// determine if there is a location
			break
		case "text":
			msg := payload.Entry[0].Changes[0].Value.Messages[0]
			wa_id := payload.Entry[0].Changes[0].Value.Contacts[0].WaId
			_, err = app.message.Create(ctx, wa_id, &msg)
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
