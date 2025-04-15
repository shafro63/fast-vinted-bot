package logger

import (
	"log/slog"
	"os"
	"sync"
)

// Global Logger
var (
	Log  *slog.Logger
	once sync.Once
)

func InitLogger() {
	once.Do(func() { // Assure que l'initialisation se fait une seule fois
		var handler slog.Handler
		config := LoadConfig()

		if config.FilePath != "" {
			file, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				panic("Impossible d'ouvrir le fichier de logs: " + err.Error())
			}
			handler = slog.NewJSONHandler(file, &slog.HandlerOptions{Level: config.Level})
		} else {
			handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: config.Level})
		}

		Log = slog.New(handler)
		slog.SetDefault(Log)
		Log.Debug("Logger initialisé", "niveau :", config.Level, "fichier :", config.FilePath)
	})
}
