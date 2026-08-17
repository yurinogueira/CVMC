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
