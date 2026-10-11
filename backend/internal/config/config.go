package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	MinJWTSecretLength = 32
)

var (
	ErrJWTSecretEmpty           = errors.New("JWT_SECRET must not be empty")
	ErrJWTRefreshSecretEmpty    = errors.New("JWT_REFRESH_SECRET must not be empty")
	ErrJWTSecretInsecure        = errors.New("JWT_SECRET cannot use insecure default placeholder value ('change-me')")
	ErrJWTRefreshSecretInsecure = errors.New("JWT_REFRESH_SECRET cannot use insecure default placeholder value ('change-me-too')")
	ErrJWTSecretTooShort        = fmt.Errorf("JWT_SECRET must be at least %d characters long", MinJWTSecretLength)
	ErrJWTRefreshSecretTooShort = fmt.Errorf("JWT_REFRESH_SECRET must be at least %d characters long", MinJWTSecretLength)
	ErrJWTSecretsIdentical      = errors.New("JWT_SECRET and JWT_REFRESH_SECRET must be different")
)

var defaultInsecureSecrets = map[string]struct{}{
	"change-me":     {},
	"change-me-too": {},
}

// ValidateJWTSecrets ensures JWT secrets meet minimum entropy requirements and are not known insecure defaults.
func ValidateJWTSecrets(secret, refreshSecret string) error {
	if secret == "" {
		return ErrJWTSecretEmpty
	}
	if _, insecure := defaultInsecureSecrets[secret]; insecure {
		return ErrJWTSecretInsecure
	}
	if len(secret) < MinJWTSecretLength {
		return fmt.Errorf("JWT_SECRET must be at least %d characters long (got %d)", MinJWTSecretLength, len(secret))
	}

	if refreshSecret == "" {
		return ErrJWTRefreshSecretEmpty
	}
	if _, insecure := defaultInsecureSecrets[refreshSecret]; insecure {
		return ErrJWTRefreshSecretInsecure
	}
	if len(refreshSecret) < MinJWTSecretLength {
		return fmt.Errorf("JWT_REFRESH_SECRET must be at least %d characters long (got %d)", MinJWTSecretLength, len(refreshSecret))
	}

	if secret == refreshSecret {
		return ErrJWTSecretsIdentical
	}
	return nil
}

type Config struct {
	Port                string
	JWTSecret           string
	JWTRefreshSecret    string
	MongoURI            string
	MongoDatabase       string
	UploadPath          string
	LogLevel            string
	StorageProvider     string
	OCIStorageBucket    string
	OCIStorageNamespace string
	OCIStorageRegion    string
	OCIStorageEndpoint  string
	AllowedOrigins      []string
	TrustedProxies      []string
	CookieDomain        string
	CookieSecure        bool
	FIPEToken           string
	FIPEBaseURL         string
	AppBaseURL          string
	SMTPHost            string
	SMTPPort            string
	SMTPUser            string
	SMTPPass            string
	EmailFrom           string
}

func (c Config) Validate() error {
	return ValidateJWTSecrets(c.JWTSecret, c.JWTRefreshSecret)
}

func Load() Config {
	cfg := Config{
		Port:                getenv("PORT", "8080"),
		JWTSecret:           getenv("JWT_SECRET", "change-me"),
		JWTRefreshSecret:    getenv("JWT_REFRESH_SECRET", "change-me-too"),
		MongoURI:            getenv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase:       getenv("MONGO_DATABASE", "cvmc"),
		UploadPath:          getenv("UPLOAD_PATH", "./data/uploads"),
		LogLevel:            getenv("LOG_LEVEL", "info"),
		StorageProvider:     getenv("STORAGE_PROVIDER", "local"),
		OCIStorageBucket:    getenv("OCI_STORAGE_BUCKET", "project-files"),
		OCIStorageNamespace: getenv("OCI_STORAGE_NAMESPACE", ""),
		OCIStorageRegion:    getenv("OCI_STORAGE_REGION", "sa-saopaulo-1"),
		OCIStorageEndpoint:  getenv("OCI_STORAGE_ENDPOINT", ""),
		AllowedOrigins:      parseCommaSeparated(getenv("ALLOWED_ORIGINS", "http://localhost:5173")),
		TrustedProxies:      parseCommaSeparated(getenv("TRUSTED_PROXIES", "127.0.0.1,::1")),
		CookieDomain:        getenv("COOKIE_DOMAIN", ""),
		CookieSecure:        getenv("COOKIE_SECURE", "true") == "true",
		FIPEToken:           getenv("FIPE_API_TOKEN", ""),
		FIPEBaseURL:         getenv("FIPE_BASE_URL", "https://fipe.parallelum.com.br/api/v2"),
		AppBaseURL:          getenv("APP_BASE_URL", "http://localhost:5173"),
		SMTPHost:            getenv("SMTP_HOST", ""),
		SMTPPort:            getenv("SMTP_PORT", "587"),
		SMTPUser:            getenv("SMTP_USER", ""),
		SMTPPass:            getenv("SMTP_PASS", ""),
		EmailFrom:           getenv("EMAIL_FROM", "no-reply@cvmc.com.br"),
	}

	// Unconditionally reject default and low-entropy JWT secrets in all environments.
	if err := cfg.Validate(); err != nil {
		log.Fatalf("FATAL: invalid configuration: %v", err)
	}

	return cfg
}

func (c Config) HTTPAddress() string {
	return ":" + c.Port
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func parseCommaSeparated(raw string) []string {
	var items []string
	for _, item := range strings.Split(raw, ",") {
		item = strings.TrimSpace(item)
		if item != "" {
			items = append(items, item)
		}
	}
	return items
}
