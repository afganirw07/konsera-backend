package imageupload

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

var pngBytes = append([]byte("\x89PNG\r\n\x1a\n"), []byte("konsera-test-image")...)

func makeFileHeader(t *testing.T, filename string, data []byte) *multipart.FileHeader {
	t.Helper()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("image", filename)
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := part.Write(data); err != nil {
		t.Fatalf("write form file: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}

	reader := multipart.NewReader(bytes.NewReader(buf.Bytes()), writer.Boundary())
	form, err := reader.ReadForm(int64(buf.Len()))
	if err != nil {
		t.Fatalf("read form: %v", err)
	}
	t.Cleanup(func() { _ = form.RemoveAll() })

	files := form.File["image"]
	if len(files) != 1 {
		t.Fatalf("expected 1 uploaded file, got %d", len(files))
	}
	return files[0]
}

func TestUploadSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if got := r.FormValue("key"); got != "test-key" {
			t.Errorf("key = %q, want test-key", got)
		}

		file, header, err := r.FormFile("image")
		if err != nil {
			t.Errorf("missing image form file: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer file.Close()

		if header.Filename != "photo.png" {
			t.Errorf("filename = %q, want photo.png", header.Filename)
		}

		var received bytes.Buffer
		if _, err := received.ReadFrom(file); err != nil {
			t.Errorf("read image: %v", err)
		}
		if !bytes.Equal(received.Bytes(), pngBytes) {
			t.Errorf("image bytes were not forwarded intact")
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"success":true,"data":{"url":"https://i.ibb.co/example/photo.png"}}`)
	}))
	defer server.Close()

	service := NewService(Config{Path: server.URL + "/1/upload", APIKey: "test-key"})

	got, err := service.Upload(context.Background(), makeFileHeader(t, "photo.png", pngBytes))
	if err != nil {
		t.Fatalf("Upload() error = %v", err)
	}
	if want := "https://i.ibb.co/example/photo.png"; got != want {
		t.Fatalf("Upload() = %q, want %q", got, want)
	}
}

func TestUploadRequiresConfiguration(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{name: "missing path", config: Config{APIKey: "test-key"}},
		{name: "missing key", config: Config{Path: "https://api.imgbb.com/1/upload"}},
		{name: "missing both", config: Config{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewService(tt.config).Upload(context.Background(), makeFileHeader(t, "photo.png", pngBytes))
			if !errors.Is(err, ErrMissingConfig) {
				t.Fatalf("Upload() error = %v, want ErrMissingConfig", err)
			}
		})
	}
}

func TestUploadRejectsInvalidFilesBeforeRequest(t *testing.T) {
	hits := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		fmt.Fprint(w, `{"success":true,"data":{"url":"https://i.ibb.co/example/photo.png"}}`)
	}))
	defer server.Close()

	service := NewService(Config{Path: server.URL, APIKey: "test-key"})

	tests := []struct {
		name string
		file *multipart.FileHeader
	}{
		{name: "nil file", file: nil},
		{name: "empty file", file: &multipart.FileHeader{Filename: "empty.png", Size: 0}},
		{name: "too large", file: &multipart.FileHeader{Filename: "big.png", Size: MaxImageSize + 1}},
		{name: "unsupported content", file: makeFileHeader(t, "note.txt", []byte("just some plain text"))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := service.Upload(context.Background(), tt.file); !errors.Is(err, ErrInvalidImage) {
				t.Fatalf("Upload() error = %v, want ErrInvalidImage", err)
			}
		})
	}

	if hits != 0 {
		t.Fatalf("upstream received %d requests, want 0", hits)
	}
}

func TestUploadHandlesUpstreamFailures(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		body    string
		wantErr string
	}{
		{name: "non-2xx response", status: http.StatusInternalServerError, body: `{}`, wantErr: "status 500"},
		{name: "malformed json", status: http.StatusOK, body: `not-json`, wantErr: "decode image upload response"},
		{name: "rejected response", status: http.StatusOK, body: `{"success":false,"error":{"message":"Invalid API key"}}`, wantErr: "Invalid API key"},
		{name: "missing url", status: http.StatusOK, body: `{"success":true,"data":{}}`, wantErr: "did not contain a URL"},
		{name: "invalid url", status: http.StatusOK, body: `{"success":true,"data":{"url":"javascript:alert(1)"}}`, wantErr: "invalid URL"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				fmt.Fprint(w, tt.body)
			}))
			defer server.Close()

			service := NewService(Config{Path: server.URL, APIKey: "test-key"})
			_, err := service.Upload(context.Background(), makeFileHeader(t, "photo.png", pngBytes))
			if err == nil {
				t.Fatal("Upload() error = nil, want error")
			}
			if !bytes.Contains([]byte(err.Error()), []byte(tt.wantErr)) {
				t.Fatalf("Upload() error = %q, want it to contain %q", err.Error(), tt.wantErr)
			}
			if bytes.Contains([]byte(err.Error()), []byte("test-key")) {
				t.Fatalf("Upload() error leaks the API key: %q", err.Error())
			}
		})
	}
}
