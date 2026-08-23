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

	if err := provider.Delete(ctx, "uploads/car.png"); err != nil {
		t.Fatalf("unexpected Delete error: %v", err)
	}
	if deletedPath != "/n/test-ns/b/test-bucket/o/uploads/car.png" {
		t.Errorf("expected deleted path /n/test-ns/b/test-bucket/o/uploads/car.png, got %s", deletedPath)
	}
}
