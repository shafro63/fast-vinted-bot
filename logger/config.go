package logger

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

// Logger Config
type Config struct {
	Level    slog.Level
	FilePath string
}

func LoadConfig() *Config {
	// Charger le fichier .env
	_ = godotenv.Load()

	level := slog.LevelInfo
	switch os.Getenv("LOG_LEVEL") {
	case "DEBUG":
		level = slog.LevelDebug
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	}

	return &Config{
		Level:    level,
		FilePath: os.Getenv("LOG_FILE"),
	}
}
