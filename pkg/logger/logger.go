package logger

import (
	"errors"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(cfg LoggerConfig) (*zap.Logger, error) {
	var logLvl zapcore.Level
	switch cfg.Level {
	case "debug":
		logLvl = zap.DebugLevel
	case "release":
		logLvl = zap.InfoLevel
	default:
		return nil, errors.New("log level not defined: allowed 'debug' or 'release'")
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
