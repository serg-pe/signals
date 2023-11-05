package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogLevels string

const (
	LogLevelDebug   LogLevels = "debug"
	LogLevelRelease LogLevels = "release"
)

func New(cfg LoggerConfig) (*zap.Logger, error) {
	var logLvl zapcore.Level
	switch cfg.Level {
	case string(LogLevelDebug):
		logLvl = zap.DebugLevel
	case string(LogLevelRelease):
		logLvl = zap.InfoLevel
	default:
		return nil, fmt.Errorf("log level not defined: allowed %s or %s, got '%s'", LogLevelDebug, LogLevelRelease, cfg.Level)
	}

	encCfg := zap.NewDevelopmentEncoderConfig()
	encCfg.EncodeTime = zapcore.RFC3339TimeEncoder

	logger := zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encCfg),
			zapcore.AddSync(os.Stdout),
			logLvl,
		),
		zap.AddCaller(),
	)

	if logLvl == zap.DebugLevel {
		logger.Warn("logger works in debug mode")
	}

	return logger, nil
}
