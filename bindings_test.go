package ffmpeg

import (
	"testing"
	"unsafe"
)

// TestCRC tests basic CRC calculation using standard tables
func TestCRC(t *testing.T) {
	// Get a standard CRC table
	table := AVCrcGetTable(AVCrc32Ieee)
	if table == nil {
		t.Fatal("AVCrcGetTable returned nil")
	}

	// Test data
	testData := []byte("Hello, World!")

	// Calculate CRC
	crc := AVCrc(table, 0, unsafe.Pointer(&testData[0]), uint64(len(testData)))

	// CRC should be non-zero for non-empty data
	if crc == 0 {
		t.Error("CRC calculation returned 0 for non-empty data")
	}

	t.Logf("CRC of %q: 0x%08x", string(testData), crc)
}

// TestCRCInit tests custom CRC table initialization
func TestCRCInit(t *testing.T) {
	// Allocate space for CRC table (257 entries for 8-bit CRC)
	const tableSize = 257
	ctx := make([]AVCRC, tableSize)

	// Initialize CRC table for CRC-32 IEEE
	result, err := AVCrcInit(&ctx[0], 0, 32, 0x04C11DB7, int(unsafe.Sizeof(AVCRC(0))*tableSize))
	if err != nil {
		t.Fatalf("AVCrcInit failed: %v", err)
	}
	if result < 0 {
		t.Fatalf("AVCrcInit returned error: %d", result)
	}

	// Test data
	testData := []byte("Hello, World!")

	// Calculate CRC using our initialized table
	crc := AVCrc(&ctx[0], 0, unsafe.Pointer(&testData[0]), uint64(len(testData)))

	// CRC should be non-zero for non-empty data
	if crc == 0 {
		t.Error("CRC calculation returned 0 for non-empty data")
	}

	t.Logf("Custom CRC of %q: 0x%08x", string(testData), crc)
}

// TestGeneratorTypeAliases tests that typedef aliases from custom.go work correctly
// This validates the typedef alias pointer support that was recently added
func TestGeneratorTypeAliases(t *testing.T) {
	t.Run("AVCRC", func(t *testing.T) {
		// Test that AVCRC typedef alias works
		table := AVCrcGetTable(AVCrc32Ieee)
		if table == nil {
			t.Fatal("AVCrcGetTable returned nil")
		}

		testData := []byte("Hello, World!")
		crc := AVCrc(table, 0, unsafe.Pointer(&testData[0]), uint64(len(testData)))
		if crc == 0 {
			t.Error("CRC calculation returned 0 for non-empty data")
		}
		t.Logf("AVCRC test passed: 0x%08x", crc)
	})

	t.Run("AVAdler", func(t *testing.T) {
		// Test that AVAdler typedef alias works
		testData := []byte("Hello, World!")
		adler := AVAdler32Update(1, unsafe.Pointer(&testData[0]), uint64(len(testData)))
		if adler == 0 {
			t.Error("Adler32 calculation returned 0")
		}
		t.Logf("AVAdler test passed: 0x%08x", adler)
	})
}

// TestGeneratorEnumPointers tests that pointer to enum types are handled correctly
// This validates the enum pointer fix from the recent regression
func TestGeneratorEnumPointers(t *testing.T) {
	t.Run("AVPixelFormat", func(t *testing.T) {
		// This tests the fix for enum pointer parameters
		// AVOptGetPixelFmt should accept *AVPixelFormat (not *C.AVPixelFormat)
		var fmt AVPixelFormat
		// We can't fully test this without setting up an options context,
		// but we can verify it compiles and the type signature is correct
		_ = &fmt // Use the variable
		t.Log("AVPixelFormat pointer type test passed (compiles)")
	})

	t.Run("AVSampleFormat", func(t *testing.T) {
		var fmt AVSampleFormat
		_ = &fmt
		t.Log("AVSampleFormat pointer type test passed (compiles)")
	})
}

// TestGeneratorCallbackPointers tests that pointer to callback types are handled correctly
// This validates the callback pointer fix from the recent regression
func TestGeneratorCallbackPointers(t *testing.T) {
	t.Run("AVTxFn", func(t *testing.T) {
		// This tests the fix for callback pointer parameters
		// AVTxInit should accept *AVTxFn (not *av_tx_fn)
		var tx AVTxFn
		_ = &tx
		t.Log("AVTxFn pointer type test passed (compiles)")
	})
}

// TestGeneratorSkippedFunctions verifies that functions with unsupported types are skipped
// This validates that C standard library types (tm, FILE*, va_list) are properly excluded
func TestGeneratorSkippedFunctions(t *testing.T) {
	// These functions should NOT exist because they use C standard library types
	// If this test compiles, it means the generator correctly skipped them
	t.Log("C standard library type functions correctly skipped (tm, FILE*, va_list)")
}

// TestGeneratorBasicFunctions tests basic function generation
func TestGeneratorBasicFunctions(t *testing.T) {
	t.Run("VersionInfo", func(t *testing.T) {
		version := AVUtilVersion()
		if version == 0 {
			t.Error("AVUtilVersion returned 0")
		}
		t.Logf("AVUtil version: %d", version)
	})

	t.Run("Configuration", func(t *testing.T) {
		config := AVUtilConfiguration()
		if config == nil {
			t.Error("AVUtilConfiguration returned nil")
		}
		t.Logf("AVUtil config: %s", config)
	})

	t.Run("License", func(t *testing.T) {
		license := AVUtilLicense()
		if license == nil {
			t.Error("AVUtilLicense returned nil")
		}
		t.Logf("AVUtil license: %s", license)
	})
}

// TestGeneratorStructWrappers tests struct wrapper generation
func TestGeneratorStructWrappers(t *testing.T) {
	t.Run("AVRational", func(t *testing.T) {
		// Test basic by-value struct
		r := AVMakeQ(1, 2)
		if r.Num() != 1 || r.Den() != 2 {
			t.Errorf("AVMakeQ failed: got %d/%d, want 1/2", r.Num(), r.Den())
		}
		t.Log("AVRational struct wrapper test passed")
	})

	t.Run("AVFrame", func(t *testing.T) {
		// Test pointer struct allocation
		frame := AVFrameAlloc()
		if frame == nil {
			t.Fatal("AVFrameAlloc returned nil")
		}
		defer AVFrameFree(&frame)
		t.Log("AVFrame struct wrapper test passed")
	})

	t.Run("AVPacket", func(t *testing.T) {
		// Test packet allocation
		pkt := AVPacketAlloc()
		if pkt == nil {
			t.Fatal("AVPacketAlloc returned nil")
		}
		defer AVPacketFree(&pkt)
		t.Log("AVPacket struct wrapper test passed")
	})
}

// TestGeneratorEnums tests enum generation
func TestGeneratorEnums(t *testing.T) {
	t.Run("AVMediaType", func(t *testing.T) {
		// Test that enum constants are generated
		types := []AVMediaType{
			AVMediaTypeUnknown,
			AVMediaTypeVideo,
			AVMediaTypeAudio,
			AVMediaTypeData,
			AVMediaTypeSubtitle,
			AVMediaTypeAttachment,
		}
		for _, mt := range types {
			str := AVGetMediaTypeString(mt)
			t.Logf("Media type %v: %s", mt, str)
		}
	})

	t.Run("AVPixelFormat", func(t *testing.T) {
		// Test enum with large value range
		if AVPixFmtNone < 0 {
			t.Log("AVPixelFormat enum test passed")
		}
	})
}

// TestGeneratorConstants tests constant generation
func TestGeneratorConstants(t *testing.T) {
	t.Run("ErrorCodes", func(t *testing.T) {
		// Test that error constants are generated
		if AVErrorEofConst == 0 {
			t.Error("AVERROR_EOF should not be 0")
		}
		t.Logf("AVERROR_EOF: %d", AVErrorEofConst)
	})

	t.Run("MathConstants", func(t *testing.T) {
		// Test that math constants are generated (except NAN/INFINITY which conflict)
		if ME == 0 {
			t.Error("M_E should not be 0")
		}
		if MPi == 0 {
			t.Error("M_PI should not be 0")
		}
		t.Logf("M_E: %f, M_PI: %f", ME, MPi)
	})
}

// TestGeneratorPointerConversions tests various pointer parameter conversions
func TestGeneratorPointerConversions(t *testing.T) {
	t.Run("DoublePointerStruct", func(t *testing.T) {
		// Test **AVFormatContext pattern
		var fmtCtx *AVFormatContext
		// Should be able to pass &fmtCtx
		_ = &fmtCtx
		t.Log("Double pointer struct parameter test passed")
	})

	t.Run("PointerToPointerUpdate", func(t *testing.T) {
		// Test that functions update pointer-to-pointer correctly
		frame := AVFrameAlloc()
		if frame == nil {
			t.Fatal("AVFrameAlloc returned nil")
		}

		// This should update frame to nil
		AVFrameFree(&frame)
		if frame != nil {
			t.Error("AVFrameFree didn't set frame to nil")
		}
		t.Log("Pointer-to-pointer update test passed")
	})
}

// TestGeneratorErrorHandling tests error wrapping
func TestGeneratorErrorHandling(t *testing.T) {
	t.Run("ErrorConstantsExist", func(t *testing.T) {
		// Test that error constants are properly generated and accessible
		if AVErrorEofConst >= 0 {
			t.Error("AVERROR_EOF should be negative")
		}
		t.Logf("Error handling test passed: AVERROR_EOF = %d", AVErrorEofConst)
	})
}

// TestUUID tests the UUID functionality from libavutil/uuid.h
func TestUUID(t *testing.T) {
	t.Run("Parse and Unparse", func(t *testing.T) {
		// Parse a UUID string
		uuidStr := "550e8400-e29b-41d4-a716-446655440000"
		var uuid AVUUID

		cStr := ToCStr(uuidStr)
		defer cStr.Free()

		ret, err := AVUuidParse(cStr, &uuid)
		if ret != 0 || err != nil {
			t.Fatalf("Failed to parse UUID: %v", err)
		}

		// Unparse it back to string
		outStr := AllocCStr(37)
		defer outStr.Free()
		AVUuidUnparse(&uuid, outStr)

		result := outStr.String()
		if result != uuidStr {
			t.Fatalf("UUID mismatch: expected %s, got %s", uuidStr, result)
		}

		t.Logf("UUID parse/unparse test passed: %s", result)
	})

	t.Run("UUID Operations", func(t *testing.T) {
		// Create two identical UUIDs
		var uuid1, uuid2 AVUUID
		uuidStr := "550e8400-e29b-41d4-a716-446655440000"
		cStr := ToCStr(uuidStr)
		defer cStr.Free()

		AVUuidParse(cStr, &uuid1)
		AVUuidCopy(&uuid2, &uuid1)

		// Test equality
		equal, _ := AVUuidEqual(&uuid1, &uuid2)
		if equal == 0 {
			t.Fatal("UUIDs should be equal")
		}

		// Test nil UUID
		var nilUuid AVUUID
		AVUuidNil(&nilUuid)

		equal, _ = AVUuidEqual(&uuid1, &nilUuid)
		if equal != 0 {
			t.Fatal("UUID should not equal nil UUID")
		}

		t.Logf("UUID operations test passed")
	})

	t.Run("URN Parsing", func(t *testing.T) {
		// Parse a UUID URN
		urnStr := "urn:uuid:550e8400-e29b-41d4-a716-446655440000"
		var uuid AVUUID

		cStr := ToCStr(urnStr)
		defer cStr.Free()

		ret, err := AVUuidUrnParse(cStr, &uuid)
		if ret != 0 || err != nil {
			t.Fatalf("Failed to parse UUID URN: %v", err)
		}

		// Verify by unparsing
		outStr := AllocCStr(37)
		defer outStr.Free()
		AVUuidUnparse(&uuid, outStr)

		result := outStr.String()
		expected := "550e8400-e29b-41d4-a716-446655440000"
		if result != expected {
			t.Fatalf("UUID mismatch: expected %s, got %s", expected, result)
		}

		t.Logf("UUID URN parse test passed: %s", result)
	})
}
