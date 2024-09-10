package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/alghurabi0/whatsapp-webhook-server-golang/internal/models"
)

func (app *application) unmarshal(r *http.Request, payload *models.Payload) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("couldn't read body, error: %v", err)
	}
	app.infoLog.Println(string(body))
	err = json.Unmarshal(body, &payload)
	if err != nil {
		return fmt.Errorf("failed to unmarshal json, error: %v", err)
	}
	return nil
}

func (app *application) validatePayload(payload *models.Payload) (string, error) {
	if payload.Entry == nil {
		return "", errors.New("entry doesn't exist")
	}
	if len(payload.Entry) < 1 {
		return "", errors.New("entry is an empty slice")
	}
	entry := payload.Entry[0]
	if entry.Changes == nil {
		return "", errors.New("changes doesn't exist")
	}
	if len(entry.Changes) < 1 {
		return "", errors.New("changes is an empty slice")
	}
	if payload.HasMessages() && !payload.HasStatuses() {
		return "msg", nil
	} else if payload.HasStatuses() && !payload.HasMessages() {
		return "status", nil
	} else {
		return "other", errors.New("payload doens't contain messages or statuses")
	}
}

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
		Text:             nil,
		Template: &models.Template{
			Name: name,
			Language: models.Language{
				Code: "ar",
			},
		},
	}
	jsonData, err := json.Marshal(msg)
	app.infoLog.Println(string(jsonData))
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
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body, error: %v", err)
	}
	app.infoLog.Println(string(body))

	return nil
}
