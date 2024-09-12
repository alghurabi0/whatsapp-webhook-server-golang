package models

import (
	"context"

	"cloud.google.com/go/firestore"
)

type MessageModel struct {
	DB *firestore.Client
}

func (m *MessageModel) Create(ctx context.Context, wa_id string, msg *Message) (string, error) {
	_, err := m.DB.Collection("contacts").Doc(wa_id).Collection("messages").Doc(msg.Id).Set(ctx, msg)
	if err != nil {
		return "", err
	}

	return msg.Id, nil
}
