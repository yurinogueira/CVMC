package storage

import "context"

type File struct {
	Name        string
	Data        []byte
	ContentType string
}

type StoredObject struct {
	FileName string
	Size     int64
	Hash     string
}

type Provider interface {
	Save(ctx context.Context, path string, file File) (StoredObject, error)
	Get(ctx context.Context, path string) ([]byte, error)
	Delete(ctx context.Context, path string) error
}
