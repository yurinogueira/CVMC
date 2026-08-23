package oci

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"cvmc/internal/application/ports/storage"
)

type Config struct {
	Namespace string
	Bucket    string
	Region    string
	Endpoint  string
	Client    *http.Client
}

type Provider struct {
	namespace string
	bucket    string
	region    string
	endpoint  string
	client    *http.Client
}

func New(cfg Config) *Provider {
	client := cfg.Client
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	endpoint := cfg.Endpoint
	if endpoint == "" {
		endpoint = fmt.Sprintf("https://objectstorage.%s.oraclecloud.com", cfg.Region)
	}
	return &Provider{
		namespace: cfg.Namespace,
		bucket:    cfg.Bucket,
		region:    cfg.Region,
		endpoint:  strings.TrimRight(endpoint, "/"),
		client:    client,
	}
}

func (p *Provider) objectURL(path string) string {
	escapedPath := url.PathEscape(strings.TrimPrefix(path, "/"))
	return fmt.Sprintf("%s/n/%s/b/%s/o/%s", p.endpoint, url.PathEscape(p.namespace), url.PathEscape(p.bucket), escapedPath)
}

func (p *Provider) Save(ctx context.Context, path string, file storage.File) (storage.StoredObject, error) {
	if p.namespace == "" || p.bucket == "" {
		return storage.StoredObject{}, errors.New("oci storage: namespace and bucket must be configured")
	}

	cleanPath := strings.TrimPrefix(path, "/")
	if strings.Contains(cleanPath, "..") {
		return storage.StoredObject{}, errors.New("invalid path: traversal detected")
	}

	targetURL := p.objectURL(cleanPath)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, targetURL, bytes.NewReader(file.Data))
	if err != nil {
		return storage.StoredObject{}, fmt.Errorf("oci storage: failed to create request: %w", err)
	}

	contentType := file.ContentType
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := p.client.Do(req)
	if err != nil {
		return storage.StoredObject{}, fmt.Errorf("oci storage: put request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return storage.StoredObject{}, fmt.Errorf("oci storage: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	hash := sha256.Sum256(file.Data)
	return storage.StoredObject{
		FileName: path,
		Size:     int64(len(file.Data)),
		Hash:     hex.EncodeToString(hash[:]),
	}, nil
}

func (p *Provider) Delete(ctx context.Context, path string) error {
	if p.namespace == "" || p.bucket == "" {
		return errors.New("oci storage: namespace and bucket must be configured")
	}

	cleanPath := strings.TrimPrefix(path, "/")
	if strings.Contains(cleanPath, "..") {
		return errors.New("invalid path: traversal detected")
	}

	targetURL := p.objectURL(cleanPath)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, targetURL, nil)
	if err != nil {
		return fmt.Errorf("oci storage: failed to create request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("oci storage: delete request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("oci storage: delete failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
