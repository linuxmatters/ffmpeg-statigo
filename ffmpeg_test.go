package ffmpeg_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/linuxmatters/ffmpeg-statigo"
	"github.com/stretchr/testify/assert"
)

func TestVersions(t *testing.T) {
	// FFmpeg 8.0.x: libavcodec version 62.11.100 (0x3E0B64 = 4066148)
	assert.Equal(t, 4066148, int(ffmpeg.AVCodecVersion()), "AVCodec version should match expected")
	assert.Equal(t, ffmpeg.LIBAVCodecVersionInt, int(ffmpeg.AVCodecVersion()), "AVCodec version func and const should match")
}

// =============================================================================
// Test 2.1: GlobalCStr Thread Safety
// =============================================================================

// TestGlobalCStr_ConcurrentAccess verifies that GlobalCStr is race-safe under
// concurrent access. This test should be run with -race flag to detect data races.
func TestGlobalCStr_ConcurrentAccess(t *testing.T) {
	t.Run("concurrent_reads_same_key", func(t *testing.T) {
		// Pre-populate a key
		key := "concurrent_read_test_key"
		initial := ffmpeg.GlobalCStr(key)

		var wg sync.WaitGroup
		const numGoroutines = 100

		results := make([]*ffmpeg.CStr, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				results[idx] = ffmpeg.GlobalCStr(key)
			}(i)
		}

		wg.Wait()

		// All goroutines should get the same pointer
		for i, result := range results {
			if result != initial {
				t.Errorf("Goroutine %d got different CStr instance", i)
			}
		}
	})

	t.Run("concurrent_writes_different_keys", func(t *testing.T) {
		var wg sync.WaitGroup
		const numGoroutines = 100

		// Each goroutine writes a unique key
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				key := fmt.Sprintf("unique_key_%d", idx)
				result := ffmpeg.GlobalCStr(key)
				if result == nil {
					t.Errorf("GlobalCStr returned nil for key %s", key)
				}
				if result.String() != key {
					t.Errorf("GlobalCStr value mismatch: expected %s, got %s", key, result.String())
				}
			}(i)
		}

		wg.Wait()
	})

	t.Run("concurrent_mixed_reads_writes", func(t *testing.T) {
		var wg sync.WaitGroup
		const numGoroutines = 100
		const numKeys = 10

		// Multiple goroutines read and write overlapping keys
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				// Use modulo to create key overlap
				key := fmt.Sprintf("mixed_key_%d", idx%numKeys)
				result := ffmpeg.GlobalCStr(key)
				if result == nil {
					t.Errorf("GlobalCStr returned nil for key %s", key)
				}
			}(i)
		}

		wg.Wait()
	})

	t.Run("same_key_returns_same_instance", func(t *testing.T) {
		key := "identity_test_key"
		first := ffmpeg.GlobalCStr(key)
		second := ffmpeg.GlobalCStr(key)
		third := ffmpeg.GlobalCStr(key)

		if first != second || second != third {
			t.Error("GlobalCStr should return the same instance for the same key")
		}
	})
}

// =============================================================================
// Test 2.2: CStr Double-Free Protection
// =============================================================================

// TestCStr_DoubleFreeProtection verifies that CStr instances from GlobalCStr
// are protected from being freed (dontFree flag), preventing memory corruption.
func TestCStr_DoubleFreeProtection(t *testing.T) {
	t.Run("globalcstr_free_is_noop", func(t *testing.T) {
		key := "free_test_key"
		cstr := ffmpeg.GlobalCStr(key)

		// This should be a no-op due to dontFree flag
		cstr.Free()

		// Should still be accessible after Free()
		if cstr.String() != key {
			t.Errorf("GlobalCStr string changed after Free(): expected %s, got %s", key, cstr.String())
		}

		// Calling Free() multiple times should also be safe
		cstr.Free()
		cstr.Free()

		// Still accessible
		if cstr.String() != key {
			t.Errorf("GlobalCStr string corrupted after multiple Free() calls")
		}
	})

	t.Run("globalcstr_same_after_free_attempt", func(t *testing.T) {
		key := "persistence_test_key"
		cstr1 := ffmpeg.GlobalCStr(key)
		cstr1.Free() // Should be no-op

		// Getting the same key should return the same instance
		cstr2 := ffmpeg.GlobalCStr(key)
		if cstr1 != cstr2 {
			t.Error("GlobalCStr returned different instance after Free() attempt")
		}
	})

	t.Run("allocated_cstr_can_be_freed", func(t *testing.T) {
		// Regular allocated CStr should be freeable
		cstr := ffmpeg.AllocCStr(64)
		if cstr == nil {
			t.Fatal("AllocCStr returned nil")
		}

		// Should not panic
		cstr.Free()
	})

	t.Run("tocstr_can_be_freed", func(t *testing.T) {
		// ToCStr creates a freeable CStr
		cstr := ffmpeg.ToCStr("freeable_string")
		if cstr == nil {
			t.Fatal("ToCStr returned nil")
		}

		// Store the value before freeing
		val := cstr.String()
		if val != "freeable_string" {
			t.Errorf("ToCStr value mismatch: expected freeable_string, got %s", val)
		}

		// Should not panic
		cstr.Free()
	})
}

// TestCStr_BasicOperations verifies basic CStr functionality.
func TestCStr_BasicOperations(t *testing.T) {
	t.Run("alloc_creates_zeroed_buffer", func(t *testing.T) {
		cstr := ffmpeg.AllocCStr(10)
		defer cstr.Free()

		// Should be empty string (zeroed buffer)
		if cstr.String() != "" {
			t.Errorf("AllocCStr should create empty string, got %q", cstr.String())
		}
	})

	t.Run("tocstr_preserves_content", func(t *testing.T) {
		testCases := []string{
			"simple",
			"with spaces",
			"unicode: 日本語",
			"", // empty string
		}

		for _, tc := range testCases {
			cstr := ffmpeg.ToCStr(tc)
			defer cstr.Free()

			if cstr.String() != tc {
				t.Errorf("ToCStr content mismatch: expected %q, got %q", tc, cstr.String())
			}
		}
	})

	t.Run("globalcstr_preserves_content", func(t *testing.T) {
		testCases := []string{
			"global_simple",
			"global with spaces",
			"global_unicode: 日本語",
		}

		for _, tc := range testCases {
			cstr := ffmpeg.GlobalCStr(tc)
			if cstr.String() != tc {
				t.Errorf("GlobalCStr content mismatch: expected %q, got %q", tc, cstr.String())
			}
		}
	})

	t.Run("rawptr_not_nil", func(t *testing.T) {
		cstr := ffmpeg.ToCStr("test")
		defer cstr.Free()

		if cstr.RawPtr() == nil {
			t.Error("RawPtr should not return nil for valid CStr")
		}
	})
}
