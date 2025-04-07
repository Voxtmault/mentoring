package logging

import (
	"log/slog"

	"github.com/voxtmault/mentoring/library-project/pkg/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(conf *config.LoggingConfig) (*lumberjack.Logger, *lumberjack.Logger) {
	serverLogger := &lumberjack.Logger{
		// Log path
		Filename: conf.ServerLogPath,
		// Log size MB
		MaxSize: conf.LogMaxSize,
		// Backup count
		MaxBackups: conf.LogMaxBackup,
		// expire days
		MaxAge: conf.LogMaxAge,
		// gzip compress
		Compress: conf.LogCompress,
	}
	errorLogger := &lumberjack.Logger{
		// Log path
		Filename: conf.ErrLogPath,
		// Log size MB
		MaxSize: conf.LogMaxSize,
		// Backup count
		MaxBackups: conf.LogMaxBackup,
		// expire days
		MaxAge: conf.LogMaxAge,
		// gzip compress
		Compress: conf.LogCompress,
	}

	return serverLogger, errorLogger
}

func CloseLogger(serverLogger, errorLogger *lumberjack.Logger) error {
	slog.Debug("closing server logger")
	if err := serverLogger.Close(); err != nil {
		return err
	}

	slog.Debug("closing error logger")
	if err := errorLogger.Close(); err != nil {
		return err
	}

	slog.Debug("logger successfully closed")
	return nil
}
