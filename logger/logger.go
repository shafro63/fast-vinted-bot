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
	once.Do(func() { // Initialize once
		var handler slog.Handler
		config := LoadConfig()

		if config.FilePath != "" {
			file, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				panic("can't open log file " + err.Error())
			}
			handler = slog.NewJSONHandler(file, &slog.HandlerOptions{Level: config.Level})
		} else {
			handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: config.Level})
		}

		Log = slog.New(handler)
		slog.SetDefault(Log)
		Log.Info("Logger initialized !", "level", config.Level, "file", config.FilePath)
	})
}
