package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/alghurabi0/whatsapp-webhook-server-golang/internal/models"
)

func (app *application) prepareMessage(r *http.Request) (*models.Message, error) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return nil, err
	}
	wa_id := r.FormValue("wa_id")
	if wa_id == "" {
		return nil, errors.New("empty wa_id")
	}
	text := r.FormValue("message_content")
	if text != "" {
		return app.sendText(text, wa_id), nil
	}
	img, handler, err := r.FormFile("img")
	if err != nil {
		if err == http.ErrMissingFile {
			fmt.Println("no file in this message")
		}
		return nil, fmt.Errorf("error with getting file from form: %v", err)
	}
	if img != nil {
		defer img.Close()
		link, err := app.services.UploadImg(img, handler.Filename)
		if err != nil {
			return nil, err
		}

		return app.sendImage(link, wa_id), nil
	}

	return nil, nil
}

func (app *application) sendText(text, wa_id string) *models.Message {
	msg := &models.Message{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               wa_id,
		Type:             "text",
		Text: &models.Text{
			PreviewUrl: false, Body: text,
		},
		Context:  nil,
		Referral: nil,
		Reaction: nil,
		Image:    nil,
		Sticker:  nil,
		Location: nil,
		Button:   nil,
	}
	return msg
}

func (app *application) sendImage(link, wa_id string) *models.Message {
	msg := &models.Message{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               wa_id,
		Type:             "image",
		Text:             nil,
		Context:          nil,
		Referral:         nil,
		Reaction:         nil,
		Image: &models.Image{
			Link: link,
		},
		Sticker:  nil,
		Location: nil,
		Button:   nil,
	}
	return msg
}

func (app *application) sendMsgReq(jsonData []byte) (*models.Response, int, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v20.0/%s/messages", app.phone_number_id)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", app.token))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, 0, err
	}

	app.infoLog.Println("recevied this json after sending a msg with facebook api")
	app.infoLog.Println(string(body))

	result := &models.Response{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, 0, err
	}

	return result, res.StatusCode, nil
}
