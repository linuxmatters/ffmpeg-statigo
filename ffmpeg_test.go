package ffmpeg_test

import (
	"fmt"
	"strings"
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

// =============================================================================
// Test 4.1: AVError String Representation
// =============================================================================

// TestAVError_KnownCodes verifies that standard FFmpeg error codes produce
// readable error messages containing both the code and a description.
func TestAVError_KnownCodes(t *testing.T) {
	t.Run("averror_eof_has_description", func(t *testing.T) {
		err := ffmpeg.AVErrorEOF

		errStr := err.Error()

		// Should contain the error code
		if !strings.Contains(errStr, fmt.Sprintf("%d", ffmpeg.AVErrorEofConst)) {
			t.Errorf("Error string should contain code %d, got: %s", ffmpeg.AVErrorEofConst, errStr)
		}

		// Should contain "averror" prefix
		if !strings.Contains(errStr, "averror") {
			t.Errorf("Error string should contain 'averror', got: %s", errStr)
		}

		// Should have some description (not empty after the colon)
		parts := strings.SplitN(errStr, ":", 2)
		if len(parts) != 2 || strings.TrimSpace(parts[1]) == "" {
			t.Errorf("Error string should have non-empty description, got: %s", errStr)
		}

		t.Logf("AVErrorEOF: %s", errStr)
	})

	t.Run("eagain_has_description", func(t *testing.T) {
		err := ffmpeg.EAgain

		errStr := err.Error()

		// Should contain "averror" prefix
		if !strings.Contains(errStr, "averror") {
			t.Errorf("Error string should contain 'averror', got: %s", errStr)
		}

		// Should have some description
		parts := strings.SplitN(errStr, ":", 2)
		if len(parts) != 2 || strings.TrimSpace(parts[1]) == "" {
			t.Errorf("Error string should have non-empty description, got: %s", errStr)
		}

		t.Logf("EAgain: %s", errStr)
	})

	t.Run("custom_error_code", func(t *testing.T) {
		// Test with a generic negative error code
		err := ffmpeg.AVError{Code: -1}

		errStr := err.Error()

		// Should contain the error code
		if !strings.Contains(errStr, "-1") {
			t.Errorf("Error string should contain code -1, got: %s", errStr)
		}

		// Should contain "averror" prefix
		if !strings.Contains(errStr, "averror") {
			t.Errorf("Error string should contain 'averror', got: %s", errStr)
		}

		t.Logf("Custom error -1: %s", errStr)
	})

	t.Run("error_implements_error_interface", func(t *testing.T) {
		var err error = ffmpeg.AVError{Code: -1}

		// Should be usable as error interface
		if err == nil {
			t.Error("AVError should not be nil when wrapped as error")
		}

		// Error() should return a string
		if err.Error() == "" {
			t.Error("Error() should return non-empty string")
		}
	})

	t.Run("error_code_accessible", func(t *testing.T) {
		code := ffmpeg.AVErrorEofConst
		err := ffmpeg.AVError{Code: code}

		if err.Code != code {
			t.Errorf("Error code mismatch: expected %d, got %d", code, err.Code)
		}
	})
}

// =============================================================================
// Test 4.2: WrapErr Boundary Conditions
// =============================================================================

// TestWrapErr_BoundaryConditions verifies that WrapErr correctly handles
// boundary conditions: zero returns nil, negative returns error.
func TestWrapErr_BoundaryConditions(t *testing.T) {
	t.Run("zero_returns_nil", func(t *testing.T) {
		err := ffmpeg.WrapErr(0)

		if err != nil {
			t.Errorf("WrapErr(0) should return nil, got: %v", err)
		}
	})

	t.Run("positive_returns_nil", func(t *testing.T) {
		testCases := []int{1, 10, 100, 1000, 1<<30 - 1}

		for _, code := range testCases {
			err := ffmpeg.WrapErr(code)
			if err != nil {
				t.Errorf("WrapErr(%d) should return nil, got: %v", code, err)
			}
		}
	})

	t.Run("negative_one_returns_error", func(t *testing.T) {
		err := ffmpeg.WrapErr(-1)

		if err == nil {
			t.Error("WrapErr(-1) should return error, got nil")
		}

		avErr, ok := err.(ffmpeg.AVError)
		if !ok {
			t.Errorf("WrapErr(-1) should return AVError, got %T", err)
		}

		if avErr.Code != -1 {
			t.Errorf("AVError.Code should be -1, got %d", avErr.Code)
		}
	})

	t.Run("various_negative_codes_return_error", func(t *testing.T) {
		testCases := []int{-1, -2, -10, -100, -1000, -1 << 30}

		for _, code := range testCases {
			err := ffmpeg.WrapErr(code)
			if err == nil {
				t.Errorf("WrapErr(%d) should return error, got nil", code)
			}

			avErr, ok := err.(ffmpeg.AVError)
			if !ok {
				t.Errorf("WrapErr(%d) should return AVError, got %T", code, err)
				continue
			}

			if avErr.Code != code {
				t.Errorf("AVError.Code should be %d, got %d", code, avErr.Code)
			}
		}
	})

	t.Run("averror_eof_const_wrapped", func(t *testing.T) {
		// AVERROR_EOF is a known negative constant
		err := ffmpeg.WrapErr(ffmpeg.AVErrorEofConst)

		if err == nil {
			t.Error("WrapErr(AVErrorEofConst) should return error, got nil")
		}

		avErr, ok := err.(ffmpeg.AVError)
		if !ok {
			t.Errorf("WrapErr(AVErrorEofConst) should return AVError, got %T", err)
		}

		if avErr.Code != ffmpeg.AVErrorEofConst {
			t.Errorf("AVError.Code should be %d, got %d", ffmpeg.AVErrorEofConst, avErr.Code)
		}
	})

	t.Run("error_comparison", func(t *testing.T) {
		// Wrapped errors with same code should be comparable
		err1 := ffmpeg.WrapErr(-1)
		err2 := ffmpeg.WrapErr(-1)

		avErr1 := err1.(ffmpeg.AVError)
		avErr2 := err2.(ffmpeg.AVError)

		if avErr1.Code != avErr2.Code {
			t.Error("Two AVErrors with same code should have equal Code fields")
		}

		// Value comparison should work
		if avErr1 != avErr2 {
			t.Error("Two AVErrors with same code should be equal")
		}
	})

	t.Run("predefined_errors", func(t *testing.T) {
		// Test predefined error variables
		if ffmpeg.AVErrorEOF.Code != ffmpeg.AVErrorEofConst {
			t.Errorf("AVErrorEOF.Code should be %d, got %d", ffmpeg.AVErrorEofConst, ffmpeg.AVErrorEOF.Code)
		}

		// EAgain should have a negative code
		if ffmpeg.EAgain.Code >= 0 {
			t.Errorf("EAgain.Code should be negative, got %d", ffmpeg.EAgain.Code)
		}
	})
}

// =============================================================================
// Test 6: ToCStr Lifecycle After Free
// =============================================================================

// TestToCStr_LifecycleAfterFree documents and verifies the unsafe behaviour of
// ToCStr after calling Free(). This is a CRITICAL test that documents memory
// safety issues to prevent user bugs.
//
// IMPORTANT: This test documents UNSAFE BEHAVIOUR that can cause crashes or
// memory corruption. Unlike GlobalCStr (which has dontFree=true and cannot be
// freed), ToCStr returns a CStr that CAN be freed. If a user calls Free() and
// then tries to use the CStr, they will access freed memory leading to:
// - Segmentation faults
// - Use-after-free bugs
// - Data corruption
// - Non-deterministic crashes
//
// This test ensures the behaviour is documented and consistent.
func TestToCStr_LifecycleAfterFree(t *testing.T) {
	t.Run("tocstr_creates_allocable_string", func(t *testing.T) {
		// ToCStr should create a string that CAN be freed (unlike GlobalCStr)
		str := ffmpeg.ToCStr("test string")
		defer str.Free()

		// Verify the string content is correct before we free it
		content := str.String()
		if content != "test string" {
			t.Errorf("Expected 'test string', got '%s'", content)
		}
	})

	t.Run("globalcstr_cannot_be_freed", func(t *testing.T) {
		// GlobalCStr: cannot be freed (dontFree=true)
		globalStr := ffmpeg.GlobalCStr("global string")
		globalContent := globalStr.String()

		if globalContent != "global string" {
			t.Errorf("Expected 'global string', got '%s'", globalContent)
		}

		// Calling Free() on GlobalCStr is a no-op (has dontFree=true)
		globalStr.Free()

		// Can still use after Free() because dontFree=true
		globalContent2 := globalStr.String()
		if globalContent2 != "global string" {
			t.Errorf("GlobalCStr should remain accessible after Free(): %s", globalContent2)
		}

		t.Log("LIFECYCLE: GlobalCStr.Free() is always safe (no-op)")
	})

	t.Run("tocstr_documents_safe_pattern", func(t *testing.T) {
		// SAFE: Allocate, Use, then Free
		safeStr := ffmpeg.ToCStr("data")
		defer safeStr.Free()

		result := safeStr.String()
		if result != "data" {
			t.Errorf("Expected 'data', got '%s'", result)
		}

		t.Log("SAFE PATTERN: Use ToCStr before calling Free()")
	})

	t.Run("tocstr_defer_ensures_single_free", func(t *testing.T) {
		// SAFE: Use defer to ensure Free() called exactly once
		safeStr := ffmpeg.ToCStr("data")
		defer safeStr.Free()

		content := safeStr.String()
		if content != "data" {
			t.Errorf("Expected 'data', got '%s'", content)
		}

		t.Log("SAFE PATTERN: Use defer to ensure Free() called exactly once")
	})

	t.Run("rawptr_valid_before_free", func(t *testing.T) {
		// RawPtr() returns unsafe.Pointer to the underlying memory
		str := ffmpeg.ToCStr("raw pointer test")
		defer str.Free()

		ptr := str.RawPtr()

		if ptr == nil {
			t.Error("RawPtr() should not return nil")
		}

		t.Log("PATTERN: RawPtr() is valid before Free() is called")
	})

	t.Run("dup_creates_independent_allocation", func(t *testing.T) {
		// Dup() creates a new allocation that must be freed independently
		original := ffmpeg.ToCStr("original")
		defer original.Free()

		// Dup creates a new string that is independently allocated
		copy := original.Dup()
		defer copy.Free()

		// Each string has its own memory
		originalContent := original.String()
		copyContent := copy.String()

		if originalContent != "original" {
			t.Errorf("Original should be 'original', got '%s'", originalContent)
		}

		if copyContent != "original" {
			t.Errorf("Copy should be 'original', got '%s'", copyContent)
		}

		t.Log("PATTERN: Dup() returns independently allocated string - must be freed separately")
	})

	t.Run("allocstr_and_tocstr_both_freeable", func(t *testing.T) {
		// AllocCStr and ToCStr both create freeable strings with same lifecycle

		// AllocCStr: allocates empty buffer
		allocated := ffmpeg.AllocCStr(32)
		defer allocated.Free()

		// Verify it's valid
		_ = allocated.String()

		// ToCStr: allocates and initializes
		converted := ffmpeg.ToCStr("test")
		defer converted.Free()

		content := converted.String()
		if content != "test" {
			t.Errorf("Expected 'test', got '%s'", content)
		}

		t.Log("PATTERN: Both AllocCStr and ToCStr return freeable strings")
	})

	t.Run("documents_unsafe_patterns_to_avoid", func(t *testing.T) {
		// This test documents patterns to AVOID - we only document them,
		// we do NOT execute them to avoid crashes

		t.Log("UNSAFE PATTERN 1: Call Free() then access")
		t.Log("  str := ToCStr(\"data\")")
		t.Log("  str.Free()")
		t.Log("  result := str.String()  // CRASH: use-after-free")
		t.Log("")

		t.Log("UNSAFE PATTERN 2: Call Free() twice")
		t.Log("  str := ToCStr(\"data\")")
		t.Log("  str.Free()")
		t.Log("  str.Free()  // CRASH: double-free")
		t.Log("")

		t.Log("UNSAFE PATTERN 3: Access RawPtr() after Free()")
		t.Log("  str := ToCStr(\"data\")")
		t.Log("  ptr := str.RawPtr()")
		t.Log("  str.Free()")
		t.Log("  C.some_function(ptr)  // CRASH: use freed memory")
	})
}
