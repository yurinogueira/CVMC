package config

import (
	"log"
	"os"
	"strings"
)

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

	// Prevent deploying with insecure JWT secrets in production.
	if cfg.LogLevel != "debug" {
		if cfg.JWTSecret == "change-me" || cfg.JWTRefreshSecret == "change-me-too" {
			log.Fatal("FATAL: JWT_SECRET and JWT_REFRESH_SECRET must be set to secure values in production (LOG_LEVEL != debug)")
		}
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
