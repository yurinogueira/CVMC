package mongo

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrInvalidID    = errors.New("invalid identifier format")
	ErrInvalidEmail = errors.New("invalid email format")

	idRegex    = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$`)
)

func SanitizeID(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" || len(trimmed) > 128 {
		return "", ErrInvalidID
	}
	if !idRegex.MatchString(trimmed) {
		return "", ErrInvalidID
	}
	return trimmed, nil
}

func SanitizeEmail(raw string) (string, error) {
	trimmed := strings.ToLower(strings.TrimSpace(raw))
	if trimmed == "" || len(trimmed) > 254 {
		return "", ErrInvalidEmail
	}
	if strings.ContainsAny(trimmed, "$\"'{};") {
		return "", ErrInvalidEmail
	}
	if !emailRegex.MatchString(trimmed) {
		return "", ErrInvalidEmail
	}
	return trimmed, nil
}
