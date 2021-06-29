package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger represents the Logger
type Logger struct {
	level  zap.AtomicLevel
	logger *zap.Logger
}

// New creates new Logger
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
		Level:             lg.level,
		Development:       false,
		DisableCaller:     true,
		DisableStacktrace: true,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields:    nil,
	}

	if debugMode {
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
		config.Encoding = "console"
		config.Development = true
		config.Sampling = nil
		config.DisableCaller = false
		config.DisableStacktrace = false
	}

	lg.logger, err = config.Build(opts...)

	return lg, err
}

// Logger returns the ZAP logger instance
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
