package logger

import (
	"errors"
	"log"

	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Lg               *zap.Logger
	ErrWrongLogLevel = errors.New("log level should be debug or production")
)

func InitLogger() {
	var lvl zap.AtomicLevel
	switch config.Conf.Logger.Level {
	case config.Debug:
		lvl = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case config.Production:
		lvl = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	default:
		log.Fatal(ErrWrongLogLevel)
	}

	Lg, _ := zap.Config{
		Level:       lvl,
		Encoding:    "json",
		OutputPaths: []string{"/home/fomalhaut/log.out"},
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,
		},
	}.Build()
	zap.ReplaceGlobals(Lg)
	defer Lg.Sync()
}
