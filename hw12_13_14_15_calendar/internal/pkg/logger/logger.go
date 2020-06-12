package logger

import (
	"errors"
	"time"

	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	ErrWrongLogLevel = errors.New("log level should be one of debug/error/info/warn")
	DefaultLogFile   = "log.out"
)

func InitLogger() error {
	var lvl zap.AtomicLevel
	switch config.Conf.Logger.Level {
	case config.Debug:
		lvl = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case config.Error:
		lvl = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case config.Info:
		lvl = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case config.Warn:
		lvl = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	default:
		return ErrWrongLogLevel
	}

	var file string
	if config.Conf.Logger.File == "" {
		file = DefaultLogFile
	} else {
		file = config.Conf.Logger.File
	}

	Lg, _ := zap.Config{
		Level:       lvl,
		Encoding:    "json",
		OutputPaths: []string{file},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey: "time",
			EncodeTime: zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format(time.Stamp))
			}),
		},
	}.Build()
	zap.ReplaceGlobals(Lg)
	return nil
}
