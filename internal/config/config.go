package config

import (
	"os"
)

type Config struct {
	HTTPPort string
	DBDSN    string
}

func Load() Config {
	cfg := Config{
		HTTPPort: getenv("HTTP_PORT", "8080"),
		DBDSN:    getenv("DB_DSN", "postgres://postgres:postgres@localhost:5432/avito?sslmode=disable"),
	}

	return cfg
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
