package user

import (
	"testing"
)

func TestSanitizeAvatarURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid https URL",
			input:    "https://lh3.googleusercontent.com/a/ACg8ocL...",
			expected: "https://lh3.googleusercontent.com/a/ACg8ocL...",
		},
		{
			name:     "valid https URL with surrounding whitespace",
			input:    "   https://lh3.googleusercontent.com/photo.jpg   \n",
			expected: "https://lh3.googleusercontent.com/photo.jpg",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "whitespace only",
			input:    "   ",
			expected: "",
		},
		{
			name:     "http insecure URL rejected",
			input:    "http://insecure.example.com/avatar.jpg",
			expected: "",
		},
		{
			name:     "javascript URI rejected",
			input:    "javascript:alert(1)",
			expected: "",
		},
		{
			name:     "data URI rejected",
			input:    "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAU...",
			expected: "",
		},
		{
			name:     "relative path rejected",
			input:    "/avatars/my-photo.jpg",
			expected: "",
		},
		{
			name:     "file scheme rejected",
			input:    "file:///etc/passwd",
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := SanitizeAvatarURL(tc.input)
			if got != tc.expected {
				t.Fatalf("expected SanitizeAvatarURL(%q) = %q, got %q", tc.input, tc.expected, got)
			}
		})
	}
}
