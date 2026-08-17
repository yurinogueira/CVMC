package local

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"

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
	return os.Remove(filepath.Join(p.basePath, path))
}
