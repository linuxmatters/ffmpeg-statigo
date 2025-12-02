package ffmpeg_test

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"sync"
	"testing"

	"github.com/linuxmatters/ffmpeg-statigo"
)

// =============================================================================
// Test 3.1: SLog Adapter Message Assembly
// =============================================================================

// mockHandler captures log records for testing
type mockHandler struct {
	mu      sync.Mutex
	records []slog.Record
}

func (h *mockHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (h *mockHandler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.records = append(h.records, r)
	return nil
}

func (h *mockHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

func (h *mockHandler) WithGroup(_ string) slog.Handler {
	return h
}

func (h *mockHandler) getRecords() []slog.Record {
	h.mu.Lock()
	defer h.mu.Unlock()
	result := make([]slog.Record, len(h.records))
	copy(result, h.records)
	return result
}

func (h *mockHandler) clear() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.records = nil
}

// TestSLogAdapter_MultilineMessages verifies that multi-line log messages
// from FFmpeg are correctly assembled and logged as complete lines.
func TestSLogAdapter_MultilineMessages(t *testing.T) {
	t.Run("single_line_with_newline", func(t *testing.T) {
		handler := &mockHandler{}
		logger := slog.New(handler)
		callback := ffmpeg.SLogAdapter(logger)

		// Single complete line
		callback(nil, ffmpeg.AVLogInfo, "This is a complete message\n")

		records := handler.getRecords()
		if len(records) != 1 {
			t.Fatalf("Expected 1 record, got %d", len(records))
		}

		// Check the message contains our text
		var logValue string
		records[0].Attrs(func(a slog.Attr) bool {
			if a.Key == "log" {
				logValue = a.Value.String()
				return false
			}
			return true
		})

		if logValue != "This is a complete message" {
			t.Errorf("Expected 'This is a complete message', got %q", logValue)
		}
	})

	t.Run("message_split_across_calls", func(t *testing.T) {
		handler := &mockHandler{}
		logger := slog.New(handler)
		callback := ffmpeg.SLogAdapter(logger)

		// FFmpeg may split messages across multiple callback invocations
		callback(nil, ffmpeg.AVLogInfo, "First part ")
		callback(nil, ffmpeg.AVLogInfo, "second part ")
		callback(nil, ffmpeg.AVLogInfo, "final part\n")

		records := handler.getRecords()
		if len(records) != 1 {
			t.Fatalf("Expected 1 record, got %d", len(records))
		}

		var logValue string
		records[0].Attrs(func(a slog.Attr) bool {
			if a.Key == "log" {
				logValue = a.Value.String()
				return false
			}
			return true
		})

		if logValue != "First part second part final part" {
			t.Errorf("Expected assembled message, got %q", logValue)
		}
	})

	t.Run("multiple_lines_in_single_call", func(t *testing.T) {
		handler := &mockHandler{}
		logger := slog.New(handler)
		callback := ffmpeg.SLogAdapter(logger)

		// Multiple lines in one callback
		callback(nil, ffmpeg.AVLogInfo, "Line one\nLine two\nLine three\n")

		records := handler.getRecords()
		if len(records) != 3 {
			t.Fatalf("Expected 3 records, got %d", len(records))
		}

		expectedLines := []string{"Line one", "Line two", "Line three"}
		for i, rec := range records {
			var logValue string
			rec.Attrs(func(a slog.Attr) bool {
				if a.Key == "log" {
					logValue = a.Value.String()
					return false
				}
				return true
			})
			if logValue != expectedLines[i] {
				t.Errorf("Record %d: expected %q, got %q", i, expectedLines[i], logValue)
			}
		}
	})

	t.Run("incomplete_line_buffered", func(t *testing.T) {
		handler := &mockHandler{}
		logger := slog.New(handler)
		callback := ffmpeg.SLogAdapter(logger)

		// Send incomplete line (no newline)
		callback(nil, ffmpeg.AVLogInfo, "Incomplete line without newline")

		// Should not log yet
		records := handler.getRecords()
		if len(records) != 0 {
			t.Errorf("Expected 0 records for incomplete line, got %d", len(records))
		}

		// Now complete it
		callback(nil, ffmpeg.AVLogInfo, " - now complete\n")

		records = handler.getRecords()
		if len(records) != 1 {
			t.Fatalf("Expected 1 record after completion, got %d", len(records))
		}

		var logValue string
		records[0].Attrs(func(a slog.Attr) bool {
			if a.Key == "log" {
				logValue = a.Value.String()
				return false
			}
			return true
		})

		if logValue != "Incomplete line without newline - now complete" {
			t.Errorf("Expected assembled message, got %q", logValue)
		}
	})

	t.Run("empty_lines_handled", func(t *testing.T) {
		handler := &mockHandler{}
		logger := slog.New(handler)
		callback := ffmpeg.SLogAdapter(logger)

		// Empty lines between content
		callback(nil, ffmpeg.AVLogInfo, "Before\n\n\nAfter\n")

		records := handler.getRecords()
		// Should have: "Before", "", "", "After" = 4 records
		if len(records) != 4 {
			t.Errorf("Expected 4 records (including empty), got %d", len(records))
		}
	})
}

// =============================================================================
// Test 3.2: Log Callback Nil Context Handling
// =============================================================================

// TestSLogAdapter_NilContext verifies that the logging adapter handles nil
// context gracefully, falling back to "global" scope.
func TestSLogAdapter_NilContext(t *testing.T) {
	t.Run("nil_context_uses_global_scope", func(t *testing.T) {
		handler := &mockHandler{}
		logger := slog.New(handler)
		callback := ffmpeg.SLogAdapter(logger)

		// Should not panic with nil context
		callback(nil, ffmpeg.AVLogInfo, "Message with nil context\n")

		records := handler.getRecords()
		if len(records) != 1 {
			t.Fatalf("Expected 1 record, got %d", len(records))
		}

		var scopeValue string
		records[0].Attrs(func(a slog.Attr) bool {
			if a.Key == "scope" {
				scopeValue = a.Value.String()
				return false
			}
			return true
		})

		if scopeValue != "global" {
			t.Errorf("Expected scope 'global' for nil context, got %q", scopeValue)
		}
	})

	t.Run("nil_context_multiple_calls", func(t *testing.T) {
		handler := &mockHandler{}
		logger := slog.New(handler)
		callback := ffmpeg.SLogAdapter(logger)

		// Multiple calls with nil context should all work
		for i := 0; i < 10; i++ {
			callback(nil, ffmpeg.AVLogInfo, "Message\n")
		}

		records := handler.getRecords()
		if len(records) != 10 {
			t.Errorf("Expected 10 records, got %d", len(records))
		}

		// All should have global scope
		for i, rec := range records {
			var scopeValue string
			rec.Attrs(func(a slog.Attr) bool {
				if a.Key == "scope" {
					scopeValue = a.Value.String()
					return false
				}
				return true
			})
			if scopeValue != "global" {
				t.Errorf("Record %d: expected scope 'global', got %q", i, scopeValue)
			}
		}
	})
}

// =============================================================================
// Test 3.3: Log Level Mapping
// =============================================================================

// TestSLogAdapter_LogLevels verifies that FFmpeg log levels are correctly
// mapped to slog levels.
func TestSLogAdapter_LogLevels(t *testing.T) {
	testCases := []struct {
		name          string
		ffmpegLevel   int
		expectedLevel slog.Level
	}{
		{"error_level", ffmpeg.AVLogError, slog.LevelError},
		{"fatal_level", ffmpeg.AVLogFatal, slog.LevelError},
		{"warning_level", ffmpeg.AVLogWarning, slog.LevelWarn},
		{"info_level", ffmpeg.AVLogInfo, slog.LevelInfo},
		{"debug_level", ffmpeg.AVLogDebug, slog.LevelDebug},
		{"trace_level", ffmpeg.AVLogTrace, slog.LevelDebug},
		{"verbose_level", ffmpeg.AVLogVerbose, slog.LevelInfo}, // Falls to default
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := &mockHandler{}
			logger := slog.New(handler)
			callback := ffmpeg.SLogAdapter(logger)

			callback(nil, tc.ffmpegLevel, "Test message\n")

			records := handler.getRecords()
			if len(records) != 1 {
				t.Fatalf("Expected 1 record, got %d", len(records))
			}

			if records[0].Level != tc.expectedLevel {
				t.Errorf("Expected level %v, got %v", tc.expectedLevel, records[0].Level)
			}
		})
	}
}

// TestSLogAdapter_QuietLevelSilent verifies that AVLogQuiet produces no output.
func TestSLogAdapter_QuietLevelSilent(t *testing.T) {
	handler := &mockHandler{}
	logger := slog.New(handler)
	callback := ffmpeg.SLogAdapter(logger)

	callback(nil, ffmpeg.AVLogQuiet, "This should not be logged\n")

	records := handler.getRecords()
	if len(records) != 0 {
		t.Errorf("Expected 0 records for AVLogQuiet, got %d", len(records))
	}
}

// =============================================================================
// Test 3.4: Thread Safety
// =============================================================================

// TestSLogAdapter_ThreadSafety verifies that the logging adapter handles
// concurrent log messages correctly.
func TestSLogAdapter_ThreadSafety(t *testing.T) {
	handler := &mockHandler{}
	logger := slog.New(handler)
	callback := ffmpeg.SLogAdapter(logger)

	const numGoroutines = 50
	const messagesPerGoroutine = 10

	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < messagesPerGoroutine; j++ {
				// Each goroutine sends complete messages
				callback(nil, ffmpeg.AVLogInfo, "Message\n")
			}
		}(i)
	}

	wg.Wait()

	records := handler.getRecords()
	expectedCount := numGoroutines * messagesPerGoroutine
	if len(records) != expectedCount {
		t.Errorf("Expected %d records, got %d", expectedCount, len(records))
	}
}

// =============================================================================
// Test 3.5: Basic Logging Integration
// =============================================================================

// TestSLogAdapter_BasicIntegration verifies basic integration with slog.
func TestSLogAdapter_BasicIntegration(t *testing.T) {
	t.Run("logs_to_buffer", func(t *testing.T) {
		var buf bytes.Buffer
		logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
		callback := ffmpeg.SLogAdapter(logger)

		callback(nil, ffmpeg.AVLogInfo, "Test message for buffer\n")

		output := buf.String()
		if !strings.Contains(output, "Test message for buffer") {
			t.Errorf("Buffer should contain message, got: %s", output)
		}
		if !strings.Contains(output, "scope=global") {
			t.Errorf("Buffer should contain scope=global, got: %s", output)
		}
	})

	t.Run("message_attribute_present", func(t *testing.T) {
		handler := &mockHandler{}
		logger := slog.New(handler)
		callback := ffmpeg.SLogAdapter(logger)

		callback(nil, ffmpeg.AVLogInfo, "Check attributes\n")

		records := handler.getRecords()
		if len(records) != 1 {
			t.Fatalf("Expected 1 record, got %d", len(records))
		}

		foundLog := false
		foundScope := false
		records[0].Attrs(func(a slog.Attr) bool {
			if a.Key == "log" {
				foundLog = true
			}
			if a.Key == "scope" {
				foundScope = true
			}
			return true
		})

		if !foundLog {
			t.Error("Expected 'log' attribute to be present")
		}
		if !foundScope {
			t.Error("Expected 'scope' attribute to be present")
		}
	})
}
