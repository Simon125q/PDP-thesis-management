package logging

import (
	"github.com/natefinch/lumberjack"
	slogmulti "github.com/samber/slog-multi"
	"log/slog"
	"os"
)

// SetupLogger configure global logger for application
func SetupLogger() {
	rotatingLogger := &lumberjack.Logger{
		Filename:   "application.log",
		MaxSize:    1,  // Max size of log file in MB
		MaxBackups: 3,  // Max number of copies
		MaxAge:     30, // Max number of days to store copies
	}

	fileHandler := slog.NewTextHandler(rotatingLogger, nil)
	terminalHandler := slog.NewTextHandler(os.Stdout, nil)

	multiHandler := slogmulti.Fanout(terminalHandler, fileHandler)

	logger := slog.New(multiHandler)
	slog.SetDefault(logger)

	slog.Info("Multi-logger is working")
}
