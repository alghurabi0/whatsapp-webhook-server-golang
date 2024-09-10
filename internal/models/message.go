package models

import (
	"context"

	"cloud.google.com/go/firestore"
)

type MessageModel struct {
	DB *firestore.Client
}

func (m *MessageModel) Create(ctx context.Context, payload *Payload) (string, error) {
	wa_id := payload.Entry[0].Changes[0].Value.Contacts[0].WaId
	msg := payload.Entry[0].Changes[0].Value.Messages[0]
	_, err := m.DB.Collection("contacts").Doc(wa_id).Collection("messages").Doc(msg.Id).Set(ctx, msg)
	if err != nil {
		return "", err
	}

	return msg.Id, nil
}
