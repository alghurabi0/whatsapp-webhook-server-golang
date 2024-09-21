package models

import (
	"context"

	"cloud.google.com/go/firestore"
)

type StatusModel struct {
	DB *firestore.Client
}

func (s *StatusModel) Create(ctx context.Context, status *Status) (string, error) {
	_, err := s.DB.Collection("contacts").Doc(status.RecipientId).Collection("messages").Doc(status.Id).Collection("status").Doc(status.Id).Set(ctx, status)
	if err != nil {
		return "", err
	}

	return status.Id, nil
}
