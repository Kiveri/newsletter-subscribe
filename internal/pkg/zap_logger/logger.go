package zap_logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger

	file *os.File
}

func NewLogger(config Config) (*Logger, error) {
	zapLevel := zap.NewAtomicLevel()
	if err := zapLevel.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, fmt.Errorf("unmarshal log level: %w", err)
	}

	if err := os.MkdirAll(config.Folder, 0755); err != nil {
		return nil, fmt.Errorf("os.MkdirAll: %w", err)
	}

	nowString := time.Now().UTC().Format("2006-01-02T15-04-05.000000")

	logFilePath := filepath.Join(config.Folder, fmt.Sprintf("%s.log", nowString))

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("os.OpenFile: %w", err)
	}

	zapConfig := zap.NewDevelopmentEncoderConfig()
	zapConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	zapEncoder := zapcore.NewConsoleEncoder(zapConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(zapEncoder, zapcore.AddSync(os.Stdout), zapLevel),
		zapcore.NewCore(zapEncoder, zapcore.AddSync(logFile), zapLevel),
	)

	zapLogger := zap.New(core, zap.AddCaller())

	return &Logger{zapLogger, logFile}, nil
}

func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		fmt.Printf("failed to close app zap_logger: %s", err)
	}
}

func (l *Logger) With(field ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(field...),
		file:   l.file,
	}
}

func FromCtx(ctx context.Context) *Logger {
	log, ok := ctx.Value("log").(*Logger)
	if !ok {
		panic("no zap_logger in context")
	}

	return log
}
