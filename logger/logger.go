package logger

import (
	"log/slog"

	"github.com/voxtmault/mentoring/config"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	serverLogger *lumberjack.Logger
	errorLogger  *lumberjack.Logger
)

func InitLogger(conf *config.LoggingConfig) error {
	serverLogger = &lumberjack.Logger{
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
	errorLogger = &lumberjack.Logger{
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

	return nil
}

func GetServerLogger() *lumberjack.Logger {
	return serverLogger
}

func GetErrorLogger() *lumberjack.Logger {
	return errorLogger
}

func LogRequest(log []byte) error {
	if _, err := serverLogger.Write(log); err != nil {
		LogError(log, err)
		return err
	}

	serverLogger.Write([]byte("\n"))

	return nil
}

func LogError(log []byte, errMsg error) error {
	return nil
}

func CloseLogger() error {
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
