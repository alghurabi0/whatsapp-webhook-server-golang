package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	gcloud "cloud.google.com/go/storage"
	"firebase.google.com/go/storage"
)

type StorageModel struct {
	ST *storage.Client
}

func (s *StorageModel) DownloadAndUploadImg(url, id string) (string, error) {
	token := os.Getenv("ACCESS_TOKEN")
	bkt, err := s.ST.DefaultBucket()
	if err != nil {
		return "", err
	}

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
		return "", fmt.Errorf("unexpected status code while getting image: %d", resp.StatusCode)
	}

	object := bkt.Object("images/" + id)
	ctx := context.Background()
	writer := object.NewWriter(ctx)

	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return "", fmt.Errorf("error while streaming img to firebase storage: %v", err)
	}
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("error closing writer: %v", err)
	}
	expiration := time.Now().Add(time.Hour * 8640)
	opts := &gcloud.SignedURLOptions{
		Expires: expiration,
		Method:  http.MethodGet,
	}
	link, err := bkt.SignedURL("images/"+id, opts)
	if err != nil {
		object.Delete(ctx)
		return "", fmt.Errorf("couldn't get signed url: %v", err)
	}
	if link == "" {
		object.Delete(ctx)
		return "", errors.New("empty photo signed file url")
	}

	return link, nil
}

func (s *StorageModel) UploadImg(img multipart.File, id string) (string, error) {
	bkt, err := s.ST.DefaultBucket()
	if err != nil {
		return "", err
	}

	object := bkt.Object("images/" + id)
	ctx := context.Background()
	writer := object.NewWriter(ctx)

	_, err = io.Copy(writer, img)
	if err != nil {
		return "", fmt.Errorf("error while streaming img to firebase storage: %v", err)
	}
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("error closing writer: %v", err)
	}
	expiration := time.Now().Add(time.Hour * 8640)
	opts := &gcloud.SignedURLOptions{
		Expires: expiration,
		Method:  http.MethodGet,
	}
	link, err := bkt.SignedURL("images/"+id, opts)
	if err != nil {
		object.Delete(ctx)
		return "", fmt.Errorf("couldn't get signed url: %v", err)
	}
	if link == "" {
		object.Delete(ctx)
		return "", errors.New("empty photo signed file url")
	}

	return link, nil
}
