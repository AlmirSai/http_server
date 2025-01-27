package logging

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Environment string
	LogLevel    string
	FilePath    string
	MaxSize     int
	MaxBackups  int
	MaxAge      int
	Compress    bool
	Format      string
	Output      string
	TimeFormat  string
}

type Logger struct {
	*zap.Logger
}

func NewLogger(config *Config) (*Logger, error) {
	// Create logs directory if it doesn't exist
	logDir := filepath.Dir(config.FilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	// Configure log level
	level, err := zapcore.ParseLevel(config.LogLevel)
	if err != nil {
		level = zapcore.InfoLevel
	}

	// Configure encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Create JSON encoder for file output
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// Create console encoder for stdout
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	// Configure log rotation
	rotator := &lumberjack.Logger{
		Filename:   config.FilePath,
		MaxSize:    config.MaxSize,    // megabytes
		MaxBackups: config.MaxBackups, // number of backups
		MaxAge:     config.MaxAge,     // days
		Compress:   config.Compress,   // compress rotated files
	}

	// Create core with both file and console output
	core := zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, zapcore.AddSync(rotator), level),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
	)

	// Create logger
	logger := zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.Fields(
			zap.String("environment", config.Environment),
			zap.String("service", "auth-service"),
			zap.Time("boot_time", time.Now()),
		),
	)

	return &Logger{Logger: logger}, nil
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	if ctx == nil {
		return l
	}

	if traceID, ok := ctx.Value("trace_id").(string); ok {
		return &Logger{Logger: l.With(zap.String("trace_id", traceID))}
	}

	return l
}

// WithFields adds structured fields to the logger
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return &Logger{l.With(zapFields...)}
}

// Metrics logs metrics data in a format suitable for Grafana
func (l *Logger) Metrics(name string, value interface{}, tags map[string]string) {
	fields := make([]zap.Field, 0, len(tags)+2)
	fields = append(fields,
		zap.String("metric_name", name),
		zap.Any("metric_value", value),
	)

	for k, v := range tags {
		fields = append(fields, zap.String(fmt.Sprintf("tag_%s", k), v))
	}

	l.Info("metric", fields...)
}

// Error logs an error message with stack trace
func (l *Logger) Error(msg string, err error, fields ...zap.Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	l.Logger.Error(msg, fields...)
}

// Fatal logs a fatal error message and exits
func (l *Logger) Fatal(msg string, err error, fields ...zap.Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	l.Logger.Fatal(msg, fields...)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, fields...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, fields...)
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.Logger.Sync()
}
