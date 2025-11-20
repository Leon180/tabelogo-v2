package logger

import (
	"bytes"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name    string
		level   string
		wantErr bool
	}{
		{
			name:    "debug level",
			level:   "debug",
			wantErr: false,
		},
		{
			name:    "info level",
			level:   "info",
			wantErr: false,
		},
		{
			name:    "warn level",
			level:   "warn",
			wantErr: false,
		},
		{
			name:    "error level",
			level:   "error",
			wantErr: false,
		},
		{
			name:    "invalid level defaults to info",
			level:   "invalid",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Init(tt.level)
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil {
				logger := GetLogger()
				if logger == nil {
					t.Error("GetLogger() returned nil")
				}
			}
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		level    string
		expected zapcore.Level
	}{
		{"debug", zapcore.DebugLevel},
		{"info", zapcore.InfoLevel},
		{"warn", zapcore.WarnLevel},
		{"error", zapcore.ErrorLevel},
		{"fatal", zapcore.FatalLevel},
		{"invalid", zapcore.InfoLevel}, // default
		{"", zapcore.InfoLevel},        // default
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			got := parseLevel(tt.level)
			if got != tt.expected {
				t.Errorf("parseLevel(%s) = %v, want %v", tt.level, got, tt.expected)
			}
		})
	}
}

func TestGetLogger(t *testing.T) {
	// Reset logger to test fallback behavior
	log = nil

	logger := GetLogger()
	if logger == nil {
		t.Error("GetLogger() should never return nil")
	}
}

func TestLoggerFunctions(t *testing.T) {
	// Create a buffer to capture logs
	var buf bytes.Buffer
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey: "msg",
		LevelKey:   "level",
	})
	core := zapcore.NewCore(encoder, zapcore.AddSync(&buf), zapcore.DebugLevel)
	log = zap.New(core)

	// Test logging functions
	Debug("test debug", zap.String("key", "value"))
	Info("test info")
	Warn("test warn")
	Error("test error")

	// Check that logs were written
	if buf.Len() == 0 {
		t.Error("No logs were written")
	}
}

func TestWith(t *testing.T) {
	err := Init("info")
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	childLogger := With(zap.String("service", "test"))
	if childLogger == nil {
		t.Error("With() returned nil")
	}
}
