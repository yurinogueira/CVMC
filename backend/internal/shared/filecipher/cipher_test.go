package filecipher

import (
	"errors"
	"strings"
	"testing"
)

func TestFileCipher_EncryptDecryptPerOwner(t *testing.T) {
	fc := New("super-secret-master-key-with-at-least-32-bytes!!")
	sampleDataURL := "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII="

	encrypted, err := fc.EncryptForOwner("user-alice", sampleDataURL)
	if err != nil {
		t.Fatalf("unexpected encryption error: %v", err)
	}
	if !strings.HasPrefix(encrypted, EncryptedPrefix) {
		t.Fatalf("expected encrypted prefix %q, got %q", EncryptedPrefix, encrypted)
	}
	if strings.Contains(encrypted, "iVBORw0KGgo") {
		t.Fatalf("ciphertext leaked plaintext base64 content")
	}

	// Decrypt with rightful owner succeeds
	decrypted, err := fc.DecryptForOwner("user-alice", encrypted)
	if err != nil {
		t.Fatalf("expected decryption to succeed for owner, got %v", err)
	}
	if decrypted != sampleDataURL {
		t.Fatalf("decrypted mismatch: got %q, want %q", decrypted, sampleDataURL)
	}

	// Decrypt with another user fails (key + AAD mismatch)
	_, err = fc.DecryptForOwner("user-bob", encrypted)
	if !errors.Is(err, ErrDecryptionFailed) {
		t.Fatalf("expected ErrDecryptionFailed for wrong owner, got %v", err)
	}

	// External HTTP URLs remain untouched
	extURL := "https://example.com/car.jpg"
	encExt, err := fc.EncryptForOwner("user-alice", extURL)
	if err != nil || encExt != extURL {
		t.Fatalf("expected external URL to remain untouched, got %q (err=%v)", encExt, err)
	}

	// Legacy unencrypted data URL is returned transparently on DecryptForOwner
	legacy, err := fc.DecryptForOwner("user-alice", sampleDataURL)
	if err != nil || legacy != sampleDataURL {
		t.Fatalf("expected legacy data URL to pass through, got %q (err=%v)", legacy, err)
	}

	// EstimateDecodedBytes returns accurate byte count (68 bytes for the 1x1 PNG)
	estimated := EstimateDecodedBytes(sampleDataURL)
	if estimated != 68 {
		t.Fatalf("expected 68 bytes for 1x1 PNG, got %d", estimated)
	}
}
