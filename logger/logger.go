package logger

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	rootLogger *zap.Logger
)

func init() {
	rootLogger = zap.L()
}

func Configure(logLevel string) error {
	config := zap.NewProductionConfig()

	switch logLevel {
	case "err", "error":
		config.Level.SetLevel(zap.ErrorLevel)
	case "Wrn", "warn", "warning":
		config.Level.SetLevel(zap.WarnLevel)
	case "dbg", "debug":
		config.Level.SetLevel(zap.DebugLevel)
	case "info":
		config.Level.SetLevel(zap.InfoLevel)
	default:
		return fmt.Errorf("logger: invalid log level %q, must be one of error, warning, debug", logLevel)
	}

	// disable capturing callers in production
	config.DisableCaller = false

	// The default values are intended for automated logging
	config.Encoding = "console"
	config.EncoderConfig = zap.NewDevelopmentEncoderConfig()

	options := []zap.Option{
		// zap calls os.Exit on the error level by default
		zap.AddStacktrace(zapcore.DPanicLevel),
	}

	var err error
	if rootLogger, err = config.Build(options...); err != nil {
		return fmt.Errorf("logger: error initializing zap: %w", err)
	}

	if _, err := zap.RedirectStdLogAt(rootLogger, zap.DebugLevel); err != nil {
		return fmt.Errorf("logger: error redirecting stdlib logger: %w", err)
	}

	// Replace root logger to have all new loggers based off of above
	// settings
	zap.ReplaceGlobals(rootLogger)

	rootLogger.Info("logging initialized")

	return nil
}

func Named(name string) *zap.Logger {
	return rootLogger.Named(name)
}

func Sync(logger *zap.Logger) error {
	if err := logger.Sync(); err != nil {
		// zap generated an error when attempting to sync stdout/stderr
		// on linux:
		// https://github.com/uber-go/zap/issues/772
		// https://github.com/uber-go/zap/issues/370
		var pe *fs.PathError
		if errors.As(err, &pe) && pe.Op == "sync" && strings.HasPrefix(pe.Path, "/dev/") {
			return nil
		}
		return fmt.Errorf("logger: error syncing logger: %w", err)
	}
	return nil
}

func SyncRoot() {
	if err := Sync(rootLogger); err != nil {
		fmt.Fprintf(os.Stderr, "error syncing root logger: %v", err)
	}
}
