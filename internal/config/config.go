package config

import (
	"os"

	"github.com/joho/godotenv"
)

type PSQLConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
	SSLmode  string
}

type Config struct {
	Port string
	PSQL PSQLConfig
}

func New() *Config {
	_ = godotenv.Load(".env")

	getEnvStr := func(key, defaultValue string) string {
		if value := os.Getenv(key); value != "" {
			return value
		}
		return defaultValue
	}

	return &Config{
		Port: ":" + getEnvStr("PORT", "8080"),
		PSQL: PSQLConfig{
			Host:     getEnvStr("DB_HOST", ""),
			Port:     getEnvStr("DB_PORT", "5432"),
			User:     getEnvStr("DB_USER", "postgres"),
			Password: getEnvStr("DB_PASS", "1234"),
			Name:     getEnvStr("DB_NAME", "mydb"),
			SSLmode:  getEnvStr("DB_SSLMODE", "disable"),
		},
	}
}
