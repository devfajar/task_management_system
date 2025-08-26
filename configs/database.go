package configs

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Configs struct {
	HTTPAddr       string
	DatabaseURL    string
	MigrateOnStart bool
	GinMode        string // debug|release|test
}

func Load(paths ...string) (Configs, error) {
	if len(paths) == 0 {
		_ = godotenv.Load()
	} else {
		_ = godotenv.Load(paths...)
	}

	cfg := Configs{
		HTTPAddr:       getEnv("HTTP_ADDR", ":8080"),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		MigrateOnStart: getBool("MIGRATE_ON_START", true),
		GinMode:        strings.ToLower(getEnv("GIN_MODE", "release")),
	}

	if cfg.DatabaseURL == "" {
		return Configs{}, errors.New("missing DATABASE_URL")
	}
	switch cfg.GinMode {
	case "debug", "release", "test":
	default:
		cfg.GinMode = "release"
	}
	return cfg, nil
}

func getEnv(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return def
}

func getBool(k string, def bool) bool {
	if v := os.Getenv(k); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return def
}
