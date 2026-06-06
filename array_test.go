package ffmpeg

import "testing"

// =============================================================================
// Test: Array Bounds Behaviour Documentation
// =============================================================================

// TestArray_BoundsBehaviour documents and verifies the bounds behaviour of the
// Array[T] type. This is a critical safety documentation test.
//
// IMPORTANT: The Array type has NO bounds checking, matching C behaviour.
// Out-of-bounds access leads to undefined behaviour (memory corruption, crashes,
// or silent data corruption). This is by design to match FFmpeg's C semantics.
//
// Users MUST track array lengths themselves, just as in C code.
func TestArray_BoundsBehaviour(t *testing.T) {
	t.Run("valid_indices_work_correctly", func(t *testing.T) {
		arr := AllocAVCodecIDArray(3)
		if arr == nil {
			t.Fatal("AllocAVCodecIDArray returned nil")
		}
		defer AVFree(arr.RawPtr())

		arr.Set(0, AVCodecIdH264)
		arr.Set(1, AVCodecIdHevc)
		arr.Set(2, AVCodecIdAV1)

		if arr.Get(0) != AVCodecIdH264 {
			t.Errorf("Get(0) = %v, want AVCodecIdH264", arr.Get(0))
		}
		if arr.Get(1) != AVCodecIdHevc {
			t.Errorf("Get(1) = %v, want AVCodecIdHevc", arr.Get(1))
		}
		if arr.Get(2) != AVCodecIdAV1 {
			t.Errorf("Get(2) = %v, want AVCodecIdAV1", arr.Get(2))
		}
	})

	t.Run("documents_no_bounds_checking", func(t *testing.T) {
		// This test documents that Array has NO bounds checking.
		// We cannot safely test out-of-bounds access as it causes undefined behaviour.
		//
		// DO NOT uncomment the following - it will cause memory corruption:
		//   arr := AllocAVCodecIDArray(3)
		//   arr.Get(100)  // UNDEFINED BEHAVIOUR - reads arbitrary memory
		//   arr.Set(100, AVCodecIdH264)  // UNDEFINED BEHAVIOUR - writes arbitrary memory
		//
		// Users must track array lengths themselves. This matches C semantics where
		// arrays are just pointers with no length information.

		t.Log("Array[T] has NO bounds checking - matches C semantics")
		t.Log("Out-of-bounds access causes undefined behaviour (crash or corruption)")
		t.Log("Users MUST track array lengths themselves")
	})

	t.Run("zero_index_is_valid", func(t *testing.T) {
		// Verify index 0 works for any non-empty array
		arr := AllocAVCodecIDArray(1)
		if arr == nil {
			t.Fatal("AllocAVCodecIDArray returned nil")
		}
		defer AVFree(arr.RawPtr())

		arr.Set(0, AVCodecIdMpeg2Video)
		if arr.Get(0) != AVCodecIdMpeg2Video {
			t.Errorf("Get(0) = %v, want AVCodecIdMpeg2Video", arr.Get(0))
		}
	})

	t.Run("last_valid_index", func(t *testing.T) {
		// Verify the last valid index (size-1) works correctly
		const size = 5
		arr := AllocAVCodecIDArray(size)
		if arr == nil {
			t.Fatal("AllocAVCodecIDArray returned nil")
		}
		defer AVFree(arr.RawPtr())

		arr.Set(size-1, AVCodecIdVp9)
		if arr.Get(size-1) != AVCodecIdVp9 {
			t.Errorf("Get(%d) = %v, want AVCodecIdVp9", size-1, arr.Get(size-1))
		}
	})

	t.Run("primitive_arrays_same_behaviour", func(t *testing.T) {
		// Verify primitive type arrays also have no bounds checking
		ptr := AVMalloc(uint64(3 * intSize))
		if ptr == nil {
			t.Fatal("AVMalloc returned nil")
		}
		defer AVFree(ptr)

		arr := ToIntArray(ptr)
		if arr == nil {
			t.Fatal("ToIntArray returned nil")
		}

		arr.Set(0, 42)
		arr.Set(1, 100)
		arr.Set(2, -1)

		if arr.Get(0) != 42 {
			t.Errorf("int array Get(0) = %v, want 42", arr.Get(0))
		}
		if arr.Get(1) != 100 {
			t.Errorf("int array Get(1) = %v, want 100", arr.Get(1))
		}
		if arr.Get(2) != -1 {
			t.Errorf("int array Get(2) = %v, want -1", arr.Get(2))
		}

		t.Log("Primitive arrays (ToIntArray, ToUint8Array, etc.) also have no bounds checking")
	})

	t.Run("nil_array_panics", func(t *testing.T) {
		// Nil array should panic on access - this IS checked
		var arr *Array[AVCodecID] = nil

		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic when accessing nil array, but none occurred")
			} else {
				t.Logf("Nil array access correctly panics: %v", r)
			}
		}()

		_ = arr.Get(0)
	})
}

// TestArray_SafeUsagePatterns documents safe patterns for using Array[T]
func TestArray_SafeUsagePatterns(t *testing.T) {
	t.Run("track_length_manually", func(t *testing.T) {
		// Safe pattern: always track length alongside array
		const length = 3
		arr := AllocAVCodecIDArray(length)
		if arr == nil {
			t.Fatal("AllocAVCodecIDArray returned nil")
		}
		defer AVFree(arr.RawPtr())

		// Safe iteration using tracked length
		values := []AVCodecID{AVCodecIdH264, AVCodecIdHevc, AVCodecIdAV1}
		for i := range length {
			arr.Set(uintptr(i), values[i])
		}

		// Safe read using tracked length
		for i := range length {
			if arr.Get(uintptr(i)) != values[i] {
				t.Errorf("Get(%d) mismatch", i)
			}
		}

		t.Log("Safe pattern: track array length in a separate variable")
	})

	t.Run("use_sentinel_values", func(t *testing.T) {
		// Safe pattern: use sentinel value to mark end (like null-terminated strings)
		// Many FFmpeg APIs use AVCodecIdNone or -1 as sentinel
		arr := AllocAVCodecIDArray(4)
		if arr == nil {
			t.Fatal("AllocAVCodecIDArray returned nil")
		}
		defer AVFree(arr.RawPtr())

		// Store values with sentinel
		arr.Set(0, AVCodecIdH264)
		arr.Set(1, AVCodecIdHevc)
		arr.Set(2, AVCodecIdAV1)
		arr.Set(3, AVCodecIdNone) // Sentinel

		// Safe iteration using sentinel
		count := 0
		for i := uintptr(0); ; i++ {
			val := arr.Get(i)
			if val == AVCodecIdNone {
				break
			}
			count++
			if count > 100 { // Safety limit
				t.Fatal("Sentinel not found within reasonable range")
			}
		}

		if count != 3 {
			t.Errorf("Found %d elements before sentinel, want 3", count)
		}

		t.Log("Safe pattern: use sentinel value (e.g., AVCodecIdNone) to mark end")
	})
}
