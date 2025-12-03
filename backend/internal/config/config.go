package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port         string
	DatabaseURL  string
	OptimizerURL string
	JWTSecret    string
	JWTExpiry    int // hours
}

func Load() *Config {
	jwtExpiry := 24
	if exp := os.Getenv("JWT_EXPIRY_HOURS"); exp != "" {
		if val, err := strconv.Atoi(exp); err == nil {
			jwtExpiry = val
		}
	}

	return &Config{
		Port:         getEnv("PORT", "8080"),
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/logitrackpro?sslmode=disable"),
		OptimizerURL: getEnv("OPTIMIZER_URL", "http://localhost:8000"),
		JWTSecret:    getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpiry:    jwtExpiry,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
