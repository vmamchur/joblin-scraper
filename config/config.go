package config

import (
	"log"
	"os"
)

type Config struct {
	TgBotUrl string

	DB     DBConfig
	Djinni DjinniConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type DjinniConfig struct {
	Email    string
	Password string
}

func Load() Config {
	return Config{
		TgBotUrl: mustEnv("TG_BOT_URL"),
		DB: DBConfig{
			Host:     mustEnv("DB_HOST"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     mustEnv("DB_USER"),
			Password: mustEnv("DB_PASSWORD"),
			Name:     mustEnv("DB_NAME"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Djinni: DjinniConfig{
			Email:    mustEnv("DJINNI_EMAIL"),
			Password: mustEnv("DJINNI_PASSWORD"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func mustEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	log.Fatalf("Missing required env: %s\n", key)
	return ""
}
