package filecipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"strings"
)

const (
	EncryptedPrefix = "enc:v1:"
	hkdfSalt        = "cvmc-file-cipher-v1-salt"
)

var (
	ErrDecryptionFailed = errors.New("failed to decrypt or authenticate file payload")
	ErrInvalidOwnerID   = errors.New("ownerID is required for file encryption")
)

type FileCipher struct {
	masterKey []byte
}

func New(masterSecret string) *FileCipher {
	clean := strings.TrimSpace(masterSecret)
	if clean == "" {
		clean = "cvmc-default-fallback-file-encryption-secret-key"
	}
	return &FileCipher{
		masterKey: []byte(clean),
	}
}

// deriveOwnerKey implements RFC 5869 HKDF-SHA256 (Extract + Expand for 32-byte AES-256 key)
func (c *FileCipher) deriveOwnerKey(ownerID string) []byte {
	// 1. HKDF-Extract(salt, IKM) -> PRK
	extractor := hmac.New(sha256.New, []byte(hkdfSalt))
	extractor.Write(c.masterKey)
	prk := extractor.Sum(nil)

	// 2. HKDF-Expand(PRK, info, 32) -> OKM (32 bytes = SHA-256 output length)
	expander := hmac.New(sha256.New, prk)
	expander.Write([]byte("cvmc-file-owner:" + ownerID))
	expander.Write([]byte{0x01})
	return expander.Sum(nil)
}

func IsEncrypted(value string) bool {
	return strings.HasPrefix(strings.TrimSpace(value), EncryptedPrefix)
}

func IsInlineFileData(value string) bool {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return false
	}
	lower := strings.ToLower(trimmed)
	if strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") {
		return false
	}
	return true
}

// EstimateDecodedBytes returns the actual size in bytes of a data URL or encrypted payload.
func EstimateDecodedBytes(value string) int64 {
	trimmed := strings.TrimSpace(value)
	if !IsInlineFileData(trimmed) {
		return 0
	}
	if strings.HasPrefix(trimmed, EncryptedPrefix) {
		rawB64 := strings.TrimPrefix(trimmed, EncryptedPrefix)
		decodedLen := int64(len(rawB64)) * 3 / 4
		// Subtract 12-byte nonce + 16-byte GCM tag + data URL header approximation
		if decodedLen > 28 {
			return (decodedLen - 28) * 3 / 4
		}
		return decodedLen
	}
	if idx := strings.Index(trimmed, ","); idx != -1 {
		b64Part := strings.TrimSpace(trimmed[idx+1:])
		padding := int64(0)
		if strings.HasSuffix(b64Part, "==") {
			padding = 2
		} else if strings.HasSuffix(b64Part, "=") {
			padding = 1
		}
		size := (int64(len(b64Part))*3)/4 - padding
		if size > 0 {
			return size
		}
	}
	return int64(len(trimmed))
}

func (c *FileCipher) EncryptForOwner(ownerID, plaintext string) (string, error) {
	if c == nil {
		return plaintext, nil
	}
	cleanOwner := strings.TrimSpace(ownerID)
	if cleanOwner == "" {
		return "", ErrInvalidOwnerID
	}
	trimmed := strings.TrimSpace(plaintext)
	if !IsInlineFileData(trimmed) || IsEncrypted(trimmed) {
		return trimmed, nil
	}

	key := c.deriveOwnerKey(cleanOwner)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	aad := []byte("owner:" + cleanOwner)
	sealed := gcm.Seal(nonce, nonce, []byte(trimmed), aad)
	return EncryptedPrefix + base64.StdEncoding.EncodeToString(sealed), nil
}

func (c *FileCipher) DecryptForOwner(ownerID, stored string) (string, error) {
	if c == nil {
		return stored, nil
	}
	trimmed := strings.TrimSpace(stored)
	if !IsEncrypted(trimmed) {
		return trimmed, nil
	}
	cleanOwner := strings.TrimSpace(ownerID)
	if cleanOwner == "" {
		return "", ErrInvalidOwnerID
	}

	rawB64 := strings.TrimPrefix(trimmed, EncryptedPrefix)
	data, err := base64.StdEncoding.DecodeString(rawB64)
	if err != nil {
		return "", ErrDecryptionFailed
	}

	key := c.deriveOwnerKey(cleanOwner)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize+gcm.Overhead() {
		return "", ErrDecryptionFailed
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	aad := []byte("owner:" + cleanOwner)
	plain, err := gcm.Open(nil, nonce, ciphertext, aad)
	if err != nil {
		return "", ErrDecryptionFailed
	}
	return string(plain), nil
}
