package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"media-agent/dtos"
	"net/http"
	"net/url"
	"time"
)

type S3HttpClient struct {
	http.Client
}

func NewS3HttpClient() *S3HttpClient {
	return &S3HttpClient{
		Client: http.Client{
			Timeout: 5 * time.Minute,
		},
	}
}

func (s *S3HttpClient) CreatePresignedURL(ctx context.Context, lambdaURL, bucketName, key string) (string, error) {
	url, err := url.Parse(lambdaURL)
	if err != nil {
		return "", nil
	}

	query := url.Query()

	query.Add("bucket-name", bucketName)
	query.Add("key", key)

	url.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), http.NoBody)
	if err != nil {
		return "", err
	}

	resp, err := s.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("expected status code 200, got %d", resp.StatusCode)
	}

	var presignedURLCreatedResponse dtos.PresignedURLCreatedResponse

	if err := json.NewDecoder(resp.Body).Decode(&presignedURLCreatedResponse); err != nil {
		return "", err
	}

	return presignedURLCreatedResponse.Url, nil
}

func (s *S3HttpClient) UploadMedia(ctx context.Context, lambdaURL string, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, lambdaURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	resp, err := s.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status code 200, got %d", resp.StatusCode)
	}

	return nil
}
