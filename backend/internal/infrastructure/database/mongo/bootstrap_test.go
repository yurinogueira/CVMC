package mongo

import (
	"errors"
	"testing"
)

func TestIsCollectionExistsError(t *testing.T) {
	tests := []struct {
		err      error
		expected bool
	}{
		{nil, false},
		{errors.New("collection already exists"), true},
		{errors.New("NamespaceExists (code 48): a collection already exists with that name"), true},
		{errors.New("connection timed out"), false},
	}

	for _, tt := range tests {
		result := isCollectionExistsError(tt.err)
		if result != tt.expected {
			t.Errorf("isCollectionExistsError(%v) = %v; want %v", tt.err, result, tt.expected)
		}
	}
}
