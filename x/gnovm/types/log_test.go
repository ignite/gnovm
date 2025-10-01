package types

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"cosmossdk.io/log"
	"gotest.tools/v3/assert"
)

func TestCosmosLoggerToSlogWrapper(t *testing.T) {
	tests := []struct {
		name     string
		logFunc  func(*slog.Logger)
		expected string
	}{
		{
			name: "debug message with attributes",
			logFunc: func(logger *slog.Logger) {
				logger.Debug("debug message", "key", "value", "number", 42)
			},
			expected: "debug message",
		},
		{
			name: "info message with attributes",
			logFunc: func(logger *slog.Logger) {
				logger.Info("info message", "user", "alice", "action", "login")
			},
			expected: "info message",
		},
		{
			name: "warn message with attributes",
			logFunc: func(logger *slog.Logger) {
				logger.Warn("warning message", "error", "connection timeout")
			},
			expected: "warning message",
		},
		{
			name: "error message with attributes",
			logFunc: func(logger *slog.Logger) {
				logger.Error("error message", "code", 500, "reason", "internal error")
			},
			expected: "error message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a buffer to capture log output
			var buf bytes.Buffer
			cosmosLogger := log.NewLogger(&buf)

			// Create slog wrapper
			slogLogger := NewSlogFromCosmosLogger(cosmosLogger)

			// Execute the log function
			tt.logFunc(slogLogger)

			// Verify the message was logged to the cosmos logger
			output := buf.String()
			assert.Assert(t, strings.Contains(output, tt.expected))
		})
	}
}

func TestSlogWrapperWithContext(t *testing.T) {
	var buf bytes.Buffer
	cosmosLogger := log.NewLogger(&buf)
	slogLogger := NewSlogFromCosmosLogger(cosmosLogger)

	// Test logging with context
	ctx := context.Background()
	slogLogger.InfoContext(ctx, "context message", "key", "value")

	output := buf.String()
	assert.Assert(t, strings.Contains(output, "context message"))
}

func TestSlogWrapperWith(t *testing.T) {
	var buf bytes.Buffer
	cosmosLogger := log.NewLogger(&buf)
	slogLogger := NewSlogFromCosmosLogger(cosmosLogger)

	// Create a child logger with additional attributes
	childLogger := slogLogger.With("module", "test", "version", "1.0")
	childLogger.Info("child logger message")

	output := buf.String()
	assert.Assert(t, strings.Contains(output, "child logger message"))
	assert.Assert(t, strings.Contains(output, "module"))
	assert.Assert(t, strings.Contains(output, "test"))
}

func TestSlogWrapperWithGroup(t *testing.T) {
	var buf bytes.Buffer
	cosmosLogger := log.NewLogger(&buf)
	slogLogger := NewSlogFromCosmosLogger(cosmosLogger)

	// Create a logger with a group
	groupLogger := slogLogger.WithGroup("database")
	groupLogger.Info("database operation", "table", "users", "operation", "select")

	output := buf.String()
	assert.Assert(t, strings.Contains(output, "database operation"))
}

func TestSlogWrapperOutputFormat(t *testing.T) {
	var buf bytes.Buffer
	cosmosLogger := log.NewLogger(&buf)
	slogLogger := NewSlogFromCosmosLogger(cosmosLogger)

	// Log a message to see the actual output format
	slogLogger.Info("test message", "module", "gnovm", "action", "demo", "count", 42)

	output := buf.String()
	t.Logf("Actual log output: %s", output)

	// Verify the message and attributes are present
	assert.Assert(t, strings.Contains(output, "test message"))
	assert.Assert(t, strings.Contains(output, "module"))
	assert.Assert(t, strings.Contains(output, "gnovm"))
	assert.Assert(t, strings.Contains(output, "action"))
	assert.Assert(t, strings.Contains(output, "demo"))
	assert.Assert(t, strings.Contains(output, "count"))
	assert.Assert(t, strings.Contains(output, "42"))
}
