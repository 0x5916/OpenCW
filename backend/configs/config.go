package configs

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strings"
)

type Config struct {
	Port            string   `required:"true"`
	DBHost          string   `required:"true"`
	DBPort          string   `required:"true"`
	DBUser          string   `required:"true"`
	DBPassword      string   `required:"true"`
	DBName          string   `required:"true"`
	JWTSecret       []byte   `required:"true"`
	ResendAPIKey    string   `required:"true"`
	ResendFromEmail string   `required:"true"`
	GinMode         string   `required:"true"`
	CORSOrigins     []string // optional — no tag
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
		Port:            os.Getenv("PORT"),
		DBHost:          os.Getenv("DB_HOST"),
		DBPort:          os.Getenv("DB_PORT"),
		DBUser:          os.Getenv("DB_USER"),
		DBPassword:      os.Getenv("DB_PASSWORD"),
		DBName:          os.Getenv("DB_NAME"),
		ResendAPIKey:    os.Getenv("RESEND_API_KEY"),
		ResendFromEmail: os.Getenv("RESEND_FROM_EMAIL"),
		GinMode:         os.Getenv("GIN_MODE"),
		CORSOrigins:     corsOrigins,
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
	if err := App.Validate(); err != nil {
		slog.Error("Error validating config", "err", err)
		os.Exit(1)
	}
}

func (c Config) Validate() error {
	v := reflect.ValueOf(c)
	t := reflect.TypeOf(c)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if fieldType.Tag.Get("required") != "true" {
			continue
		}

		if isZero(field) {
			return fmt.Errorf("config field %q is required but empty", fieldType.Name)
		}
	}
	return nil
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Slice, reflect.Map, reflect.Array:
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	default:
		return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
	}
}

func GetGinMode() string {
	return os.Getenv("GIN_MODE")
}

func (c Config) IsRelease() bool {
	return c.GinMode == "release"
}
