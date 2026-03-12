package configs

import (
	"encoding/base64"
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	Port        string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	JWTSecret   []byte
	GinMode     string
	CORSOrigins []string
}

var App Config

func Load() {
	var corsOrigins []string
	if raw := os.Getenv("CORS_ORIGINS"); raw != "" {
		for _, o := range strings.Split(raw, ",") {
			if trimmed := strings.TrimSpace(o); trimmed != "" {
				corsOrigins = append(corsOrigins, trimmed)
			}
		}
	}

	App = Config{
		Port:        os.Getenv("PORT"),
		DBHost:      os.Getenv("DB_HOST"),
		DBPort:      os.Getenv("DB_PORT"),
		DBUser:      os.Getenv("DB_USER"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		DBName:      os.Getenv("DB_NAME"),
		GinMode:     os.Getenv("GIN_MODE"),
		CORSOrigins: corsOrigins,
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		slog.Error("JWT_SECRET is empty")
		os.Exit(1)
	}
	decoded, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		slog.Error("Cannot decode JWT_SECRET. Is it in base64 format?", "err", err)
		os.Exit(1)
	}
	if len(decoded) < 32 {
		slog.Warn("JWT_SECRET length is lower than 32 bytes. For security, consider a longer secret.", "length", len(decoded))
	}
	App.JWTSecret = decoded
}

func GetGinMode() string {
	return os.Getenv("GIN_MODE")
}

func (c Config) IsRelease() bool {
	return c.GinMode == "release"
}
