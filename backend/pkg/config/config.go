package config

import (
	"log"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	RedisURL    string
	StripeKey   string
}

func Load() *Config {
	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),
		StripeKey:   getEnv("STRIPE_KEY", "sk_test_placeholder"),
	}

	return cfg
}

func getEnv(key string, fallback string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		if fallback == "" {
			log.Fatalf("Missing required env variable: %s", key)
		}
		return fallback
	}
	return val
}
