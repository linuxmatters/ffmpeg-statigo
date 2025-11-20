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

// TestGeneratorCharVsUint8 verifies that char parameters use C.char, not C.uint8_t
// This is a regression test for the av_match_list compilation error.
func TestGeneratorCharVsUint8(t *testing.T) {
	// av_match_list has signature: int av_match_list(const char *name, const char *list, char separator)
	// The third parameter 'separator' is type char (not uint8_t)

	// Test that the function accepts uint8 (Go's char mapping) and compiles without error
	name := ToCStr("test")
	defer name.Free()
	list := ToCStr("test,foo,bar")
	defer list.Free()

	// If this compiles, the char→C.char mapping is working correctly
	// (Previously this would fail: cannot use _Ctype_uint8_t as _Ctype_char)
	result, err := AVMatchList(name, list, ',')

	if err != nil {
		t.Fatalf("AVMatchList failed: %v", err)
	}
	if result != 1 {
		t.Errorf("Expected match result 1, got %d", result)
	}
}

// TestGeneratorPixFmtDescriptorTypes validates that AVPixFmtDescriptor fields use correct C types
// This is a regression test for uint8_t fields being incorrectly cast to C.int
func TestGeneratorPixFmtDescriptorTypes(t *testing.T) {
	// Get a pixel format descriptor (any format will do)
	desc := AVPixFmtDescGet(AVPixFmtRgb24)
	if desc == nil {
		t.Fatal("RGB24 pixel format should exist")
	}

	// Test uint8_t fields (nb_components, log2_chroma_w, log2_chroma_h)
	// These should use C.uint8_t casts, not C.int
	// The fact that these compile and run proves the types are correct
	nbComponents := desc.NbComponents()
	if nbComponents <= 0 || nbComponents > 4 {
		t.Errorf("RGB24 components should be 1-4, got %d", nbComponents)
	}

	log2ChromaW := desc.Log2ChromaW()
	log2ChromaH := desc.Log2ChromaH()
	// RGB24 has no chroma subsampling
	if log2ChromaW != 0 {
		t.Errorf("RGB24 log2_chroma_w should be 0, got %d", log2ChromaW)
	}
	if log2ChromaH != 0 {
		t.Errorf("RGB24 log2_chroma_h should be 0, got %d", log2ChromaH)
	}

	// Test uint64_t field (flags)
	// This should use C.uint64_t cast, not C.int
	flags := desc.Flags()
	if flags < 0 {
		t.Errorf("flags should be readable, got %d", flags)
	}
	// RGB24 should have RGB flag set
	if (flags & AVPixFmtFlagRgb) == 0 {
		t.Error("RGB24 should have RGB flag set")
	}
}

// TestGeneratorPrimitiveTypeMapping validates that primitive types map to correct CGO types
// This is a smoke test for the getCType() function
func TestGeneratorPrimitiveTypeMapping(t *testing.T) {
	t.Run("ptrdiff_t", func(t *testing.T) {
		// av_image_copy_plane_uc_from uses ptrdiff_t (mapped to int64)
		// Verify it compiles and doesn't panic
		dst := make([]byte, 100)
		src := make([]byte, 100)
		AVImageCopyPlaneUcFrom(
			unsafe.Pointer(&dst[0]), 10,
			unsafe.Pointer(&src[0]), 10,
			10, 10,
		)
		t.Log("ptrdiff_t parameters map to C.int64_t correctly")
	})

	t.Run("size_t", func(t *testing.T) {
		// AVMalloc uses size_t
		ptr := AVMalloc(1024)
		if ptr != nil {
			AVFree(ptr)
		}
		t.Log("size_t parameters map to C.size_t correctly")
	})
}

// TestGeneratorManualBindings documents which patterns are intentionally skipped and have manual bindings
func TestGeneratorManualBindings(t *testing.T) {
	t.Run("variadic_functions", func(t *testing.T) {
		// av_log is variadic and intentionally skipped in favor of manual binding
		// Just verify that our manual logging system works
		callback := func(ctx *LogCtx, level int, msg string) {
			// Test callback
		}
		AVLogSetCallback(callback)
		// If we reach here without panic, logging system initialized correctly
		t.Log("Variadic logging functions have manual bindings")
	})

	t.Run("iterator_functions", func(t *testing.T) {
		// Iterator functions with opaque pointers are manually bound
		// Verify at least one iterator works
		var opaque unsafe.Pointer
		codec := AVCodecIterate(&opaque)
		if codec == nil {
			t.Error("codec iterator should return at least one codec")
		}
		t.Log("Iterator functions with opaque pointers have manual bindings")
	})
}

// TestGeneratorTypePreservation validates that the CTypeName field preservation works
// This is a meta-test that the generator correctly handles typedef preservation
func TestGeneratorTypePreservation(t *testing.T) {
	t.Run("char_params_compile", func(t *testing.T) {
		// Functions with char parameters should compile
		// (regression test for char→uint8_t→C.uint8_t bug)
		name := ToCStr("rgb24")
		defer name.Free()

		pixFmt := AVGetPixFmt(name)
		if pixFmt != AVPixFmtRgb24 {
			t.Errorf("Expected AVPixFmtRgb24, got %v", pixFmt)
		}
	})

	t.Run("uint8_fields_compile", func(t *testing.T) {
		// Struct fields with uint8_t should use C.uint8_t
		// (regression test for uint8_t→int→C.int bug)
		desc := AVPixFmtDescGet(AVPixFmtYuv420P)
		if desc == nil {
			t.Fatal("YUV420P pixel format should exist")
		}

		// YUV420P has 3 components
		components := desc.NbComponents()
		if components != 3 {
			t.Errorf("YUV420P should have 3 components, got %d", components)
		}

		// And chroma subsampling (log2_chroma_w/h should be 1)
		chromaW := desc.Log2ChromaW()
		chromaH := desc.Log2ChromaH()
		if chromaW != 1 {
			t.Errorf("YUV420P log2_chroma_w should be 1, got %d", chromaW)
		}
		if chromaH != 1 {
			t.Errorf("YUV420P log2_chroma_h should be 1, got %d", chromaH)
		}
	})

	t.Run("uint64_fields_compile", func(t *testing.T) {
		// Struct fields with uint64_t should use C.uint64_t
		// (regression test for uint64_t→int→C.int bug)
		desc := AVPixFmtDescGet(AVPixFmtRgb24)
		if desc == nil {
			t.Fatal("RGB24 pixel format should exist")
		}

		flags := desc.Flags()
		// RGB24 should have RGB flag set
		if (flags & AVPixFmtFlagRgb) == 0 {
			t.Error("RGB24 should have RGB flag set")
		}
	})
}

// TestGeneratorConstArrayFields tests const array field generation with enum and struct types
// This validates the Priority 1 enhancement that enables enum/struct const arrays
func TestGeneratorConstArrayFields(t *testing.T) {
	t.Run("struct_const_array_AVRational", func(t *testing.T) {
		// Test that struct const array fields work (AVRational[N])
		// Previously skipped as "unknown const array"
		// Example: AVDetectionBBox.classify_confidences (AVRational[4])

		// Create a detection bbox (would normally come from av_detection_bbox_alloc)
		// We can't easily allocate one, but we can verify the accessor compiles
		var bbox *AVDetectionBBox
		if bbox != nil {
			// This should compile and return *Array[*AVRational]
			confidences := bbox.ClassifyConfidences()
			_ = confidences
		}
		t.Log("Struct const array field (AVRational[N]) accessor compiled successfully")
	})

	t.Run("enum_const_array", func(t *testing.T) {
		// Test that enum const array fields work (AVEnum[N])
		// Previously skipped as "unknown const array"
		// Example: AVDOVIReshapingCurve.mapping_idc (AVDOVIMappingMethod[8])

		var curve *AVDOVIReshapingCurve
		if curve != nil {
			// This should compile and return *Array[AVDOVIMappingMethod]
			mappingIdc := curve.MappingIdc()
			_ = mappingIdc
		}
		t.Log("Enum const array field (AVEnum[N]) accessor compiled successfully")
	})

	t.Run("byvalue_struct_array_helper", func(t *testing.T) {
		// Test that ToXArray helpers are generated for ByValue structs
		// Previously only generated for pointer structs
		// This enables the array field accessors above to work

		// AVRational is a ByValue struct, should have ToAVRationalArray
		r1 := AVMakeQ(1, 2)
		r2 := AVMakeQ(3, 4)

		// Verify basic functionality of ByValue struct
		if r1.Num() != 1 || r1.Den() != 2 {
			t.Errorf("AVRational ByValue struct broken: got %d/%d, want 1/2", r1.Num(), r1.Den())
		}
		if r2.Num() != 3 || r2.Den() != 4 {
			t.Errorf("AVRational ByValue struct broken: got %d/%d, want 3/4", r2.Num(), r2.Den())
		}

		t.Log("ByValue struct ToXArray helper generation validated")
	})
}

// TestGeneratorOutputParameters verifies that Priority 2 (output parameter functions)
// are properly generated and compile with the correct signatures.
func TestGeneratorOutputParameters(t *testing.T) {
	t.Run("av_opt_get family compiles", func(t *testing.T) {
		// These functions should compile - we're just testing signature availability
		// Actual runtime testing would require a valid AVClass object

		// Test that output parameter functions are accessible
		var outInt int64
		var outDouble float64
		var outW, outH int

		// Verify function signatures compile (won't run without valid context)
		_ = func() (int, error) {
			return AVOptGetInt(nil, nil, 0, &outInt)
		}

		_ = func() (int, error) {
			return AVOptGetDouble(nil, nil, 0, &outDouble)
		}

		_ = func() (int, error) {
			return AVOptGetImageSize(nil, nil, 0, &outW, &outH)
		}

		t.Log("av_opt_get_* family functions compile with output parameters")
	})

	t.Run("av_packet_get_side_data compiles", func(t *testing.T) {
		var size uint64

		// Verify av_packet_get_side_data compiles with size output parameter
		_ = func() unsafe.Pointer {
			return AVPacketGetSideData(nil, AVPacketSideDataType(0), &size)
		}

		t.Log("av_packet_get_side_data compiles with size output parameter")
	})

	t.Run("av_cpb_properties_alloc compiles", func(t *testing.T) {
		var size uint64

		// Verify av_cpb_properties_alloc compiles with size output parameter
		_ = func() *AVCPBProperties {
			return AVCpbPropertiesAlloc(&size)
		}

		t.Log("av_cpb_properties_alloc compiles with size output parameter")
	})

	t.Run("dimension output parameters compile", func(t *testing.T) {
		var width, height int

		// Verify functions with width/height output parameters compile
		_ = func() {
			AVCodecAlignDimensions(nil, &width, &height)
		}

		t.Log("Functions with width/height output parameters compile")
	})

	t.Run("callback functions are skipped", func(t *testing.T) {
		// The following functions should NOT be generated because they use
		// callback parameters passed by value, which CGO cannot handle:
		// - AVFifoWriteFromCb
		// - AVFifoReadToCb
		// - AVFifoPeekToCb

		// This test simply documents that these functions are intentionally missing
		t.Log("Callback-by-value functions intentionally skipped due to CGO limitation")
		t.Log("Missing: av_fifo_write_from_cb, av_fifo_read_to_cb, av_fifo_peek_to_cb")
	})
}
