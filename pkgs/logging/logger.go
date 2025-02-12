package logging

import (
	"log"
	"log/slog"
	"os"

	"github.com/natefinch/lumberjack"
	slogmulti "github.com/samber/slog-multi"
)

// SetupLogger configure global logger for application
func SetupLogger() {
	// Ensure logs directory exists
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		log.Fatalf("Failed to create logs directory: %v", err)
	}
	rotatingLogger := &lumberjack.Logger{
		Filename:   "logs/application.log",
		MaxSize:    6,  // Max size of log file in MB
		MaxBackups: 40, // Max number of copies
		MaxAge:     30, // Max number of days to store copies
	}

	fileHandler := slog.NewTextHandler(rotatingLogger, nil)
	// unncomment to enable loggoing in termianl
	// terminalHandler := slog.NewTextHandler(os.Stdout, nil)

	// unncomment to enable loggoing in termianl
	multiHandler := slogmulti.Fanout( /*terminalHandler,*/ fileHandler)

	logger := slog.New(multiHandler)
	slog.SetDefault(logger)

	slog.Info("Multi-logger is working")
}
