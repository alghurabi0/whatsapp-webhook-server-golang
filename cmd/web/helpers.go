package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/alghurabi0/whatsapp-webhook-server-golang/internal/models"
)

func (app *application) unmarshalPayload(r *http.Request, payload *models.Payload) error {
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

func (app *application) getOrCreateContact(ctx context.Context, payload *models.Payload) error {
	contact := payload.Entry[0].Changes[0].Value.Contacts[0]
	_, err := app.contact.Get(ctx, contact.WaId)
	if err != nil {
		c := &models.Contact{
			WaId: contact.WaId,
			Name: contact.Profile.Name,
		}
		_, err = app.contact.Create(ctx, c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (app *application) saveStatus(ctx context.Context, payload *models.Payload) error {
	status := payload.Entry[0].Changes[0].Value.Statuses[0]
	_, err := app.status.Create(ctx, &status)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) saveMessage(ctx context.Context, payload *models.Payload) error {
	msg := payload.Entry[0].Changes[0].Value.Messages[0]
	msgType := payload.Entry[0].Changes[0].Value.Messages[0].Type
	wa_id := payload.Entry[0].Changes[0].Value.Contacts[0].WaId

	// determine type or location, is there a referral?
	switch msgType {
	case "":
		// determine if there is a location
		break
	case "text":
		_, err := app.message.Create(ctx, wa_id, &msg)
		if err != nil {
			return err
		}
	case "reaction":
		break
	case "image":
		imgId := msg.Image.Id
		url, err := app.getImgUrl(imgId)
		if err != nil {
			return err
		}
		strgUrl, err := app.services.DownloadAndUploadImg(url, imgId)
		if err != nil {
			return err
		}
		msg.Image.Link = strgUrl
		_, err = app.message.Create(ctx, wa_id, &msg)
		if err != nil {
			return err
		}
	case "sticker":
		break
	case "button":
		break
	default:
		break
	}
	return nil
}

func (app *application) saveTextMessage(ctx context.Context, msg *models.Message, wa_id string) error {
	_, err := app.message.Create(ctx, wa_id, msg)
	if err != nil {
		return err
	}
	return nil
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

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{}
}

func (app *application) getImgUrl(id string) (string, error) {
	token := os.Getenv("ACCESS_TOKEN")
	url := fmt.Sprintf("https://graph.facebook.com/v20.0/%s/", id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}
	url, ok := data["url"].(string)
	if !ok {
		return "", fmt.Errorf("error: 'url' key not found or not a string")
	}

	return url, nil
}
