package logger

import (
	"fmt"
	"github.com/natefinch/lumberjack/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"main/internal/config"
	"os"
)

func toLumberjack(f *config.FileLoggerConfig) (*lumberjack.Roller, error) {
	roller, err := lumberjack.NewRoller(f.Filename, f.MaxSize, &lumberjack.Options{
		MaxAge:     f.MaxAge,
		MaxBackups: f.MaxBackups,
		LocalTime:  f.LocalTime,
		Compress:   f.Compress,
	})
	if err != nil {
		return nil, fmt.Errorf("initialisating roller: %w", err)
	}
	return roller, nil
}

func NewZapLogger(cfg config.LoggerConfig) (*zap.Logger, error) {
	level, err := zap.ParseAtomicLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("parsing level: %w", err)
	}
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	jsonEncoder := zapcore.NewJSONEncoder(encoderCfg)

	var cores []zapcore.Core

	stdout := zapcore.AddSync(os.Stdout)
	cores = append(cores, zapcore.NewCore(jsonEncoder, stdout, level))

	if cfg.File != nil {
		roller, err := toLumberjack(cfg.File)
		if err != nil {
			return nil, fmt.Errorf("initialisating lumberjack: %w", err)
		}
		file := zapcore.AddSync(roller)
		cores = append(cores, zapcore.NewCore(jsonEncoder, file, level))
	}

	core := zapcore.NewTee(cores...)
	return zap.New(core), nil
}

func ReplaceZap(cfg config.LoggerConfig) (func(), error) {
	logger, err := NewZapLogger(cfg)
	if err != nil {
		return nil, err
	}
	return zap.ReplaceGlobals(logger.WithOptions(zap.AddCaller())), nil
}
