package local

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"cvmc/internal/application/ports/storage"
)

type Provider struct {
	basePath string
}

func New(basePath string) *Provider {
	return &Provider{basePath: basePath}
}

func (p *Provider) Save(ctx context.Context, path string, file storage.File) (storage.StoredObject, error) {
	_ = ctx
	fullPath := filepath.Join(p.basePath, path)

	// Prevent path traversal
	absBase, _ := filepath.Abs(p.basePath)
	absFull, _ := filepath.Abs(fullPath)
	if !strings.HasPrefix(absFull, absBase+string(filepath.Separator)) && absFull != absBase {
		return storage.StoredObject{}, errors.New("invalid path: traversal detected")
	}

	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return storage.StoredObject{}, err
	}

	if err := os.WriteFile(fullPath, file.Data, 0o644); err != nil {
		return storage.StoredObject{}, err
	}

	hash := sha256.Sum256(file.Data)
	return storage.StoredObject{FileName: path, Size: int64(len(file.Data)), Hash: hex.EncodeToString(hash[:])}, nil
}

func (p *Provider) Delete(ctx context.Context, path string) error {
	_ = ctx
	fullPath := filepath.Join(p.basePath, path)

	// Prevent path traversal
	absBase, _ := filepath.Abs(p.basePath)
	absFull, _ := filepath.Abs(fullPath)
	if !strings.HasPrefix(absFull, absBase+string(filepath.Separator)) && absFull != absBase {
		return errors.New("invalid path: traversal detected")
	}

	return os.Remove(fullPath)
}
