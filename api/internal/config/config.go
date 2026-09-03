package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DatabaseURL    string
	AllowedOrigins []string
	Env            string
	LogLevel       string
	JWTSecret      string
	JWTTTL         string
}

func Load() Config {
	_ = godotenv.Load()
	_ = godotenv.Load(".env.local")

	return Config{
		Port:           envOr("PORT", "8080"),
		DatabaseURL:    envOr("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/employee_directory?sslmode=disable"),
		AllowedOrigins: splitCSV(envOr("CORS_ALLOWED_ORIGINS", "http://localhost:4321,http://localhost:4200,http://localhost:5173")),
		Env:            envOr("APP_ENV", "development"),
		LogLevel:       envOr("LOG_LEVEL", "info"),
		JWTSecret:      envOr("JWT_SECRET", "dev-insecure-secret-change-me"),
		JWTTTL:         envOr("JWT_TTL", "24h"),
	}
}

func envOr(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}

func EnvInt(k string, fallback int) int {
	if v := os.Getenv(k); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
