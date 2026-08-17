package config

import "os"

type Config struct {
	Port               string
	JWTSecret          string
	JWTRefreshSecret   string
	MongoURI           string
	MongoDatabase      string
	UploadPath         string
	LogLevel           string
	StorageProvider    string
}

func Load() Config {
	return Config{
		Port:             getenv("PORT", "8080"),
		JWTSecret:        getenv("JWT_SECRET", "change-me"),
		JWTRefreshSecret:  getenv("JWT_REFRESH_SECRET", "change-me-too"),
		MongoURI:          getenv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase:     getenv("MONGO_DATABASE", "cvmc"),
		UploadPath:        getenv("UPLOAD_PATH", "./data/uploads"),
		LogLevel:          getenv("LOG_LEVEL", "info"),
		StorageProvider:   getenv("STORAGE_PROVIDER", "local"),
	}
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
