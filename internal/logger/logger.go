package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	level  zap.AtomicLevel
	logger *zap.Logger
}

func New(level string, debugMode bool, opts ...zap.Option) (*Logger, error) {
	var err error

	lg := &Logger{}

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "ts",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	lg.level = zap.NewAtomicLevel()
	lg.level.SetLevel(convertStringLevelToZap(level))

	config := zap.Config{
		Level:       lg.level,
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	if debugMode {
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
		config.Encoding = "console"
		config.Development = true
		config.Sampling = nil
	}

	lg.logger, err = config.Build(opts...)

	return lg, err
}

func (l *Logger) Logger() *zap.Logger {
	return l.logger
}

func convertStringLevelToZap(level string) zapcore.Level {
	switch level {
	case "ERROR":
		return zap.ErrorLevel
	case "INFO":
		return zap.InfoLevel
	case "DEBUG":
		return zap.DebugLevel
	}

	return zap.ErrorLevel
}
