package models

import (
	"cloud.google.com/go/firestore"
	"firebase.google.com/go/storage"
)

type ImageModel struct {
	DB *firestore.Client
	ST *storage.Client
}

func (i *ImageModel) SaveImage() {}
