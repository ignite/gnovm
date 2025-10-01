package types

import (
	"context"
	"log/slog"

	"cosmossdk.io/log"
)

// cosmosLoggerHandler implements slog.Handler to bridge cosmos logger to slog.
// It converts slog logging calls to the corresponding cosmos logger methods,
// preserving structured logging attributes and group hierarchies.
type cosmosLoggerHandler struct {
	cosmosLogger log.Logger  // The underlying cosmos SDK logger
	attrs        []slog.Attr // Handler-level attributes that get added to all log records
	groups       []string    // Group names that form a hierarchical prefix for attribute keys
}

// NewSlogFromCosmosLogger creates an slog.Logger that wraps a cosmos logger.
// This allows code expecting an slog.Logger to work with cosmos SDK loggers,
// converting slog calls to the appropriate cosmos logger methods while
// preserving structured logging attributes.
//
// Example usage:
//
//	cosmosLogger := log.NewLogger(os.Stdout)
//	slogLogger := NewSlogFromCosmosLogger(cosmosLogger)
//	slogLogger.Info("message", "key", "value")
func NewSlogFromCosmosLogger(cosmosLogger log.Logger) *slog.Logger {
	handler := &cosmosLoggerHandler{
		cosmosLogger: cosmosLogger,
		attrs:        make([]slog.Attr, 0),
		groups:       make([]string, 0),
	}
	return slog.New(handler)
}

// Enabled implements slog.Handler.
// Since cosmos logger doesn't expose level checking capabilities,
// we assume all levels are enabled and let the cosmos logger handle filtering.
func (h *cosmosLoggerHandler) Enabled(_ context.Context, level slog.Level) bool {
	return true
}

// Handle implements slog.Handler.
// It converts an slog.Record to the appropriate cosmos logger method call,
// combining handler-level attributes with record-level attributes.
func (h *cosmosLoggerHandler) Handle(_ context.Context, record slog.Record) error {
	// Build key-value pairs from attributes and record attributes
	keyvals := make([]interface{}, 0, len(h.attrs)*2+record.NumAttrs()*2)

	// Add handler-level attributes
	for _, attr := range h.attrs {
		keyvals = append(keyvals, h.formatAttrKey(attr.Key), attr.Value.Any())
	}

	// Add record-level attributes
	record.Attrs(func(attr slog.Attr) bool {
		keyvals = append(keyvals, h.formatAttrKey(attr.Key), attr.Value.Any())
		return true
	})

	// Route to appropriate cosmos logger method based on level
	switch record.Level {
	case slog.LevelDebug:
		h.cosmosLogger.Debug(record.Message, keyvals...)
	case slog.LevelInfo:
		h.cosmosLogger.Info(record.Message, keyvals...)
	case slog.LevelWarn:
		h.cosmosLogger.Warn(record.Message, keyvals...)
	case slog.LevelError:
		h.cosmosLogger.Error(record.Message, keyvals...)
	default:
		h.cosmosLogger.Info(record.Message, keyvals...)
	}

	return nil
}

// WithAttrs implements slog.Handler.
// It returns a new handler that includes the given attributes in all log records.
// This creates a child logger with additional structured data.
func (h *cosmosLoggerHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	copy(newAttrs[len(h.attrs):], attrs)

	return &cosmosLoggerHandler{
		cosmosLogger: h.cosmosLogger,
		attrs:        newAttrs,
		groups:       h.groups,
	}
}

// WithGroup implements slog.Handler.
// It returns a new handler that prefixes all attribute keys with the group name,
// creating a hierarchical structure (e.g., "database.connection.timeout").
func (h *cosmosLoggerHandler) WithGroup(name string) slog.Handler {
	newGroups := make([]string, len(h.groups)+1)
	copy(newGroups, h.groups)
	newGroups[len(h.groups)] = name

	return &cosmosLoggerHandler{
		cosmosLogger: h.cosmosLogger,
		attrs:        h.attrs,
		groups:       newGroups,
	}
}

// formatAttrKey formats attribute keys with group prefixes.
// Groups are joined with dots to create hierarchical keys.
// For example, with groups ["db", "conn"] and key "timeout",
// this returns "db.conn.timeout".
func (h *cosmosLoggerHandler) formatAttrKey(key string) string {
	if len(h.groups) == 0 {
		return key
	}

	result := ""
	for i, group := range h.groups {
		if i > 0 {
			result += "."
		}
		result += group
	}
	return result + "." + key
}
