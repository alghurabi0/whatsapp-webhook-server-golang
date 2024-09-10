package models

import (
	"context"

	"cloud.google.com/go/firestore"
)

type ContactModel struct {
	DB *firestore.Client
}

func (c *ContactModel) Get(ctx context.Context, wa_id string) (*Contact, error) {
	doc, err := c.DB.Collection("contacts").Doc(wa_id).Get(ctx)
	if err != nil {
		return &Contact{}, err
	}
	var contact Contact
	err = doc.DataTo(&contact)
	if err != nil {
		return &Contact{}, err
	}
	contact.WaId = doc.Ref.ID
	return &contact, nil
}

func (c *ContactModel) Create(ctx context.Context, contact *Contact) (string, error) {
	_, err := c.DB.Collection("contacts").Doc(contact.WaId).Set(ctx, contact)
	if err != nil {
		return "", err
	}

	return contact.WaId, nil
}
