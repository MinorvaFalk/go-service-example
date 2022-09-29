package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	// Sugared zap logger for easier usage
	Sugar *zap.SugaredLogger

	// Zap default logger
	Logger *zap.Logger
}

func NewLogger() *Logger {

	config := zap.NewProductionEncoderConfig()
	// Configure time encoder for logging
	config.EncodeTime = zapcore.RFC3339TimeEncoder

	// Configure encoding type
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	logLevel := zapcore.InfoLevel

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), logLevel),
	)

	// Added caller option for showing where function called
	// Added stacktrace option for error level
	// logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	logger := zap.New(core)

	return &Logger{
		Sugar:  logger.Sugar(),
		Logger: logger,
	}
}
