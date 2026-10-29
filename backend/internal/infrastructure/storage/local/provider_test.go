package local

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"cvmc/internal/application/ports/storage"
)

func TestSaveCreatesFile(t *testing.T) {
	tmp := t.TempDir()
	provider := New(tmp)

	result, err := provider.Save(context.Background(), "cars/a.txt", storage.File{Name: "a.txt", Data: []byte("hello"), ContentType: "text/plain"})
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}

	fullPath := filepath.Join(tmp, "cars/a.txt")
	if _, err := os.Stat(fullPath); err != nil {
		t.Fatalf("expected stored file: %v", err)
	}

	if result.Size != 5 {
		t.Fatalf("unexpected size: %d", result.Size)
	}
}

func TestSaveRejectsPathTraversal(t *testing.T) {
	tmp := t.TempDir()
	provider := New(tmp)

	_, err := provider.Save(context.Background(), "../outside.txt", storage.File{Name: "outside.txt", Data: []byte("malicious"), ContentType: "text/plain"})
	if err == nil {
		t.Fatalf("expected path traversal error, got nil")
	}
}

func TestDeleteRejectsPathTraversal(t *testing.T) {
	tmp := t.TempDir()
	provider := New(tmp)

	err := provider.Delete(context.Background(), "../../outside.txt")
	if err == nil {
		t.Fatalf("expected path traversal error, got nil")
	}
}

func TestGetFile(t *testing.T) {
	tmp := t.TempDir()
	provider := New(tmp)

	_, err := provider.Save(context.Background(), "reports/r.pdf", storage.File{Name: "r.pdf", Data: []byte("pdf-data"), ContentType: "application/pdf"})
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}

	data, err := provider.Get(context.Background(), "reports/r.pdf")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if string(data) != "pdf-data" {
		t.Fatalf("unexpected content: %s", string(data))
	}

	_, err = provider.Get(context.Background(), "../traversal.pdf")
	if err == nil {
		t.Fatalf("expected traversal error, got nil")
	}
}
