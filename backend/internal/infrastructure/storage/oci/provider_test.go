package oci

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"cvmc/internal/application/ports/storage"
)

func TestOCIStorageProvider_SaveAndDelete(t *testing.T) {
	var savedBody []byte
	var deletedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			b, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			savedBody = b
			w.WriteHeader(http.StatusOK)
		case http.MethodDelete:
			deletedPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		case http.MethodGet:
			if r.URL.Path == "/n/test-ns/b/test-bucket/o/uploads/car.png" {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("sample file content"))
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	defer server.Close()

	provider := New(Config{
		Namespace: "test-ns",
		Bucket:    "test-bucket",
		Region:    "sa-saopaulo-1",
		Endpoint:  server.URL,
		Client:    server.Client(),
	})

	ctx := context.Background()
	testFile := storage.File{
		Name:        "car.png",
		Data:        []byte("sample file content"),
		ContentType: "image/png",
	}

	stored, err := provider.Save(ctx, "uploads/car.png", testFile)
	if err != nil {
		t.Fatalf("unexpected Save error: %v", err)
	}

	if stored.FileName != "uploads/car.png" {
		t.Errorf("expected FileName uploads/car.png, got %s", stored.FileName)
	}
	if stored.Size != int64(len(testFile.Data)) {
		t.Errorf("expected size %d, got %d", len(testFile.Data), stored.Size)
	}
	if string(savedBody) != string(testFile.Data) {
		t.Errorf("server received %q, expected %q", string(savedBody), string(testFile.Data))
	}

	gotData, err := provider.Get(ctx, "uploads/car.png")
	if err != nil {
		t.Fatalf("unexpected Get error: %v", err)
	}
	if string(gotData) != "sample file content" {
		t.Errorf("expected sample file content, got %s", string(gotData))
	}

	_, err = provider.Get(ctx, "../secret.txt")
	if err == nil {
		t.Error("expected error for path traversal in Get, got nil")
	}

	if err := provider.Delete(ctx, "uploads/car.png"); err != nil {
		t.Fatalf("unexpected Delete error: %v", err)
	}
	if deletedPath != "/n/test-ns/b/test-bucket/o/uploads/car.png" {
		t.Errorf("expected deleted path /n/test-ns/b/test-bucket/o/uploads/car.png, got %s", deletedPath)
	}

	// Test path traversal defense
	_, err = provider.Save(ctx, "../secret.txt", testFile)
	if err == nil {
		t.Error("expected error for path traversal in Save, got nil")
	}

	err = provider.Delete(ctx, "../secret.txt")
	if err == nil {
		t.Error("expected error for path traversal in Delete, got nil")
	}
}

func TestOCIStorageProvider_PreAuthenticatedRequest(t *testing.T) {
	ctx := context.Background()
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("par content"))
	}))
	defer server.Close()

	// PAR URL format: {server.URL}/p/some-par-token/n/test-ns/b/test-bucket/o/
	parEndpoint := server.URL + "/p/par-token/n/test-ns/b/test-bucket/o"
	provider := New(Config{
		Endpoint: parEndpoint,
		Client:   server.Client(),
	})

	stored, err := provider.Save(ctx, "reports/user_1/report.pdf", storage.File{
		Name: "report.pdf",
		Data: []byte("pdf-data"),
	})
	if err != nil {
		t.Fatalf("unexpected PAR Save error: %v", err)
	}
	if stored.FileName != "reports/user_1/report.pdf" {
		t.Errorf("unexpected StoredObject FileName: %s", stored.FileName)
	}
	if receivedPath != "/p/par-token/n/test-ns/b/test-bucket/o/reports/user_1/report.pdf" {
		t.Errorf("expected PAR path, got %s", receivedPath)
	}

	data, err := provider.Get(ctx, "reports/user_1/report.pdf")
	if err != nil {
		t.Fatalf("unexpected PAR Get error: %v", err)
	}
	if string(data) != "par content" {
		t.Errorf("expected 'par content', got %s", string(data))
	}
}
