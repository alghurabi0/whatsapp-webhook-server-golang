package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"firebase.google.com/go/storage"
)

type StorageModel struct {
	ST *storage.Client
}

func (s *StorageModel) GetImageUrl(id string) (string, error) {
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
