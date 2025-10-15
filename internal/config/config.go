package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	HTTPPort    string
	SMTPPort    string
	JWTSecret   string
	DomainName  string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", ""),
		HTTPPort:    getEnv("HTTP_PORT", "8080"),
		SMTPPort:    getEnv("SMTP_PORT", "1025"),
		JWTSecret:   getEnv("JWT_SECRET", "default_secret"),
		DomainName:  getEnv("DOMAIN_NAME", "localhost"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
