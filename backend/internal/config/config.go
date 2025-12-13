package config

import (
	"log"
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

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		defaultSecret := "your-secret-key-change-in-production"
		log.Printf("WARNING: JWT_SECRET environment variable is not set. Using insecure default value.")
		log.Printf("WARNING: This is a security risk! Set JWT_SECRET environment variable in production.")
		log.Printf("WARNING: Application will fail to start in production mode if JWT_SECRET is not set.")
		jwtSecret = defaultSecret
		
		// In production, fail if insecure default is used
		if os.Getenv("ENV") == "production" || os.Getenv("ENVIRONMENT") == "production" {
			log.Fatal("FATAL: JWT_SECRET must be set in production environment. Refusing to start with insecure default.")
		}
	} else if jwtSecret == "your-secret-key-change-in-production" {
		log.Printf("WARNING: JWT_SECRET is set to the insecure default value.")
		log.Printf("WARNING: Please set a secure random value for JWT_SECRET in production.")
	}

	return &Config{
		Port:         getEnv("PORT", "8080"),
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/logitrackpro?sslmode=disable"),
		OptimizerURL: getEnv("OPTIMIZER_URL", "http://localhost:8000"),
		JWTSecret:    jwtSecret,
		JWTExpiry:    jwtExpiry,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
