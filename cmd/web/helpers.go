package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/alghurabi0/whatsapp-webhook-server-golang/internal/models"
)

func (app *application) sendTemplate(name string) error {
	token := os.Getenv("ACCESS_TOKEN")
	phone_number_id := os.Getenv("PHONE_NUMBER_ID")
	url := fmt.Sprintf("https://graph.facebook.com/v20.0/%s/messages", phone_number_id)
	if token == "" {
		return errors.New("empty access token")
	}
	phone := "9647802089950"
	msg := &models.SendMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               phone,
		Type:             "template",
		Template: models.Template{
			Name: name,
			Language: models.Language{
				Code: "ar",
			},
		},
	}
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("couldn't marshal to json, error: %v", err)

	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("couldn't get new req: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send a req: %v", err)
	}

	defer res.Body.Close()
	app.infoLog.Printf("message status: %s\n", res.Status)
	return nil
}
