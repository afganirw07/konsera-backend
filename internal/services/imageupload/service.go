package imageupload

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

const MaxImageSize int64 = 5 * 1024 * 1024

var (
	ErrMissingConfig = errors.New("image upload configuration is missing")
	ErrInvalidImage  = errors.New("invalid image file")
)

var allowedImageTypes = map[string]struct{}{
	"image/jpeg": {},
	"image/png":  {},
	"image/webp": {},
	"image/gif":  {},
}

type Config struct {
	Path   string
	APIKey string
}

type Uploader interface {
	Upload(ctx context.Context, file *multipart.FileHeader) (string, error)
}

type Service struct {
	config Config
	client *http.Client
}

func NewService(config Config) *Service {
	return &Service{
		config: config,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *Service) Upload(ctx context.Context, file *multipart.FileHeader) (string, error) {
	if s == nil || strings.TrimSpace(s.config.Path) == "" || strings.TrimSpace(s.config.APIKey) == "" {
		return "", ErrMissingConfig
	}

	data, contentType, err := readImage(file)
	if err != nil {
		return "", err
	}

	return s.send(ctx, data, sanitizeFilename(file.Filename), contentType)
}

func readImage(file *multipart.FileHeader) ([]byte, string, error) {
	if file == nil {
		return nil, "", fmt.Errorf("%w: file is required", ErrInvalidImage)
	}
	if file.Size <= 0 {
		return nil, "", fmt.Errorf("%w: file is empty", ErrInvalidImage)
	}
	if file.Size > MaxImageSize {
		return nil, "", fmt.Errorf("%w: file exceeds %d MB limit", ErrInvalidImage, MaxImageSize/(1024*1024))
	}

	opened, err := file.Open()
	if err != nil {
		return nil, "", fmt.Errorf("%w: cannot open file", ErrInvalidImage)
	}
	defer opened.Close()

	data, err := io.ReadAll(io.LimitReader(opened, MaxImageSize+1))
	if err != nil {
		return nil, "", fmt.Errorf("%w: cannot read file", ErrInvalidImage)
	}
	if len(data) == 0 {
		return nil, "", fmt.Errorf("%w: file is empty", ErrInvalidImage)
	}
	if int64(len(data)) > MaxImageSize {
		return nil, "", fmt.Errorf("%w: file exceeds %d MB limit", ErrInvalidImage, MaxImageSize/(1024*1024))
	}

	contentType := http.DetectContentType(data)
	if _, ok := allowedImageTypes[contentType]; !ok {
		return nil, "", fmt.Errorf("%w: unsupported content type %s", ErrInvalidImage, contentType)
	}

	return data, contentType, nil
}

func (s *Service) send(ctx context.Context, data []byte, filename, contentType string) (string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if err := writer.WriteField("key", s.config.APIKey); err != nil {
		return "", fmt.Errorf("build image upload request: %w", err)
	}

	part, err := writer.CreateFormFile("image", filename)
	if err != nil {
		return "", fmt.Errorf("build image upload request: %w", err)
	}

	if _, err := part.Write(data); err != nil {
		return "", fmt.Errorf("build image upload request: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("build image upload request: %w", err)
	}

	endpoint, err := url.Parse(s.config.Path)
	if err != nil {
		return "", fmt.Errorf("invalid image upload endpoint")
	}
	query := endpoint.Query()
	query.Set("key", s.config.APIKey)
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), &body)
	if err != nil {
		return "", fmt.Errorf("build image upload request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("send image upload request: %w", err)
	}
	defer resp.Body.Close()

	payload, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", fmt.Errorf("read image upload response: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("image upload failed with status %d", resp.StatusCode)
	}

	var result struct {
		Success bool `json:"success"`
		Data    struct {
			URL string `json:"url"`
		} `json:"data"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(payload, &result); err != nil {
		return "", fmt.Errorf("decode image upload response: %w", err)
	}
	if !result.Success {
		if result.Error != nil && strings.TrimSpace(result.Error.Message) != "" {
			return "", fmt.Errorf("image upload rejected: %s", result.Error.Message)
		}
		return "", errors.New("image upload rejected")
	}

	imageURL := strings.TrimSpace(result.Data.URL)
	if imageURL == "" {
		return "", errors.New("image upload response did not contain a URL")
	}

	parsed, err := url.Parse(imageURL)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return "", errors.New("image upload response contained an invalid URL")
	}

	return imageURL, nil
}

func sanitizeFilename(name string) string {
	name = filepath.Base(strings.TrimSpace(name))
	if name == "" || name == "." || name == "/" {
		return "image"
	}
	return name
}
