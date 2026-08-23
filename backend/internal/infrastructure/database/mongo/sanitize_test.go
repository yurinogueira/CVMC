package mongo

import "testing"

func TestSanitizeID(t *testing.T) {
	validIDs := []string{
		"60349d834e64df540e03b011",
		"user-123",
		"car_456-abc",
		"ABC123xyz",
	}
	for _, id := range validIDs {
		sanitized, err := SanitizeID(id)
		if err != nil {
			t.Errorf("SanitizeID(%q) returned unexpected error: %v", id, err)
		}
		if sanitized != id {
			t.Errorf("SanitizeID(%q) = %q; want %q", id, sanitized, id)
		}
	}

	invalidIDs := []string{
		"",
		" ",
		"user$name",
		"id;drop",
		"user/123",
		"{ $gt: '' }",
		`{"$ne": null}`,
	}
	for _, id := range invalidIDs {
		_, err := SanitizeID(id)
		if err == nil {
			t.Errorf("SanitizeID(%q) expected error, got nil", id)
		}
	}
}

func TestSanitizeEmail(t *testing.T) {
	validEmails := []string{
		"user@example.com",
		"test.name+tag@sub.domain.org",
		"ana@cvmc.com.br",
	}
	for _, email := range validEmails {
		sanitized, err := SanitizeEmail(email)
		if err != nil {
			t.Errorf("SanitizeEmail(%q) returned unexpected error: %v", email, err)
		}
		if sanitized == "" {
			t.Errorf("SanitizeEmail(%q) returned empty string", email)
		}
	}

	invalidEmails := []string{
		"",
		"   ",
		"invalid-email",
		"user$name@test.com",
		"user@test;drop.com",
		"{ $gt: '' }@example.com",
	}
	for _, email := range invalidEmails {
		_, err := SanitizeEmail(email)
		if err == nil {
			t.Errorf("SanitizeEmail(%q) expected error, got nil", email)
		}
	}
}
