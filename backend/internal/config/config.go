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
	ErrFileEncryptionInsecure   = errors.New("FILE_ENCRYPTION_KEY cannot use insecure default placeholder value")
	ErrJWTSecretTooShort        = fmt.Errorf("JWT_SECRET must be at least %d characters long", MinJWTSecretLength)
	ErrJWTRefreshSecretTooShort = fmt.Errorf("JWT_REFRESH_SECRET must be at least %d characters long", MinJWTSecretLength)
	ErrJWTSecretsIdentical      = errors.New("JWT_SECRET and JWT_REFRESH_SECRET must be different")
)

var defaultInsecureSecrets = map[string]struct{}{
	"change-me":     {},
	"change-me-too": {},
	"cvmc_dev_jwt_access_secret_key_at_least_32_chars!":                          {},
	"cvmc_dev_jwt_refresh_secret_key_at_least_32_chars!":                         {},
	"cvmc_dev_file_encryption_master_secret_key_32_chars!":                       {},
	"cvmc-default-fallback-file-encryption-secret-key":                           {},
	"substitua_por_uma_chave_secreta_jwt_forte":                                  {},
	"substitua_por_uma_chave_secreta_refresh_forte":                              {},
	"substitua_por_uma_chave_secreta_jwt_forte_com_no_minimo_32_chars":           {},
	"substitua_por_uma_chave_secreta_refresh_forte_com_no_minimo_32_chars":       {},
	"substitua_por_uma_chave_de_criptografia_de_arquivos_com_no_minimo_32_chars": {},
}

func isInsecureDefaultSecret(value string) bool {
	trimmed := strings.TrimSpace(value)
	if _, insecure := defaultInsecureSecrets[trimmed]; insecure {
		return true
	}
	lower := strings.ToLower(trimmed)
	return strings.HasPrefix(lower, "change-me") ||
		strings.HasPrefix(lower, "cvmc_dev_") ||
		strings.HasPrefix(lower, "substitua_por_")
}

func validateSecret(name, value string, emptyErr, insecureErr error) error {
	if value == "" {
		return emptyErr
	}
	if isInsecureDefaultSecret(value) {
		return insecureErr
	}
	if len(strings.TrimSpace(value)) < MinJWTSecretLength {
		return fmt.Errorf("%s must be at least %d characters long (got %d)", name, MinJWTSecretLength, len(strings.TrimSpace(value)))
	}
	return nil
}

// ValidateJWTSecrets ensures JWT secrets meet minimum entropy requirements and are not known insecure defaults.
func ValidateJWTSecrets(secret, refreshSecret string) error {
	if err := validateSecret("JWT_SECRET", secret, ErrJWTSecretEmpty, ErrJWTSecretInsecure); err != nil {
		return err
	}
	if err := validateSecret("JWT_REFRESH_SECRET", refreshSecret, ErrJWTRefreshSecretEmpty, ErrJWTRefreshSecretInsecure); err != nil {
		return err
	}
	if secret == refreshSecret {
		return ErrJWTSecretsIdentical
	}
	return nil
}

// ValidateFileEncryptionKey ensures an explicit FILE_ENCRYPTION_KEY is not a known default/placeholder and meets minimum length.
func ValidateFileEncryptionKey(key string) error {
	if key == "" {
		return nil
	}
	if isInsecureDefaultSecret(key) {
		return ErrFileEncryptionInsecure
	}
	if len(strings.TrimSpace(key)) < MinJWTSecretLength {
		return fmt.Errorf("FILE_ENCRYPTION_KEY must be at least %d characters long when provided", MinJWTSecretLength)
	}
	return nil
}

type Config struct {
	Port                        string
	JWTSecret                   string
	JWTRefreshSecret            string
	MongoURI                    string
	MongoDatabase               string
	UploadPath                  string
	LogLevel                    string
	StorageProvider             string
	OCIStorageBucket            string
	OCIStorageNamespace         string
	OCIStorageRegion            string
	OCIStorageEndpoint          string
	AllowedOrigins              []string
	TrustedProxies              []string
	CookieDomain                string
	CookieSecure                bool
	FIPEToken                   string
	FIPEBaseURL                 string
	AppBaseURL                  string
	SMTPHost                    string
	SMTPPort                    string
	SMTPUser                    string
	SMTPPass                    string
	EmailFrom                   string
	GoogleClientID              string
	FirebaseCredentialsJSON     string
	FirebaseCredentialsFile     string
	MaintenanceCronInterval     string
	ReportCleanupCronInterval   string
	MileageReminderCronInterval string
	FileEncryptionKey           string
}

func (c Config) Validate() error {
	if err := ValidateJWTSecrets(c.JWTSecret, c.JWTRefreshSecret); err != nil {
		return err
	}
	return ValidateFileEncryptionKey(c.FileEncryptionKey)
}

func Load() Config {
	cfg := Config{
		Port:                        getenv("PORT", "8080"),
		JWTSecret:                   getenv("JWT_SECRET", "change-me"),
		JWTRefreshSecret:            getenv("JWT_REFRESH_SECRET", "change-me-too"),
		MongoURI:                    getenv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase:               getenv("MONGO_DATABASE", "cvmc"),
		UploadPath:                  getenv("UPLOAD_PATH", "./data/uploads"),
		LogLevel:                    getenv("LOG_LEVEL", "info"),
		StorageProvider:             getenv("STORAGE_PROVIDER", "local"),
		OCIStorageBucket:            getenv("OCI_STORAGE_BUCKET", ""),
		OCIStorageNamespace:         getenv("OCI_STORAGE_NAMESPACE", ""),
		OCIStorageRegion:            getenv("OCI_STORAGE_REGION", "sa-saopaulo-1"),
		OCIStorageEndpoint:          getenv("OCI_STORAGE_ENDPOINT", ""),
		AllowedOrigins:              parseCommaSeparated(getenv("ALLOWED_ORIGINS", "http://localhost:5173,https://localhost,capacitor://localhost")),
		TrustedProxies:              parseCommaSeparated(getenv("TRUSTED_PROXIES", "127.0.0.1,::1")),
		CookieDomain:                getenv("COOKIE_DOMAIN", ""),
		CookieSecure:                getenv("COOKIE_SECURE", "true") == "true",
		FIPEToken:                   getenv("FIPE_API_TOKEN", ""),
		FIPEBaseURL:                 getenv("FIPE_BASE_URL", "https://fipe.parallelum.com.br/api/v2"),
		AppBaseURL:                  getenv("APP_BASE_URL", "http://localhost:5173"),
		SMTPHost:                    getenv("SMTP_HOST", ""),
		SMTPPort:                    getenv("SMTP_PORT", "587"),
		SMTPUser:                    getenv("SMTP_USER", ""),
		SMTPPass:                    getenv("SMTP_PASS", ""),
		EmailFrom:                   getenv("EMAIL_FROM", "no-reply@cvmc.com.br"),
		GoogleClientID:              getenv("GOOGLE_CLIENT_ID", ""),
		FirebaseCredentialsJSON:     getenv("FIREBASE_CREDENTIALS_JSON", ""),
		FirebaseCredentialsFile:     getenv("FIREBASE_CREDENTIALS_FILE", ""),
		MaintenanceCronInterval:     getenv("MAINTENANCE_CRON_INTERVAL", "24h"),
		ReportCleanupCronInterval:   getenv("REPORT_CLEANUP_CRON_INTERVAL", "24h"),
		MileageReminderCronInterval: getenv("MILEAGE_REMINDER_CRON_INTERVAL", "168h"),
		FileEncryptionKey:           getenv("FILE_ENCRYPTION_KEY", ""),
	}

	if cfg.FileEncryptionKey == "" {
		cfg.FileEncryptionKey = cfg.JWTSecret + ":cvmc-file-encryption-fallback"
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
