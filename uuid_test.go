package ffmpeg

import (
	"os"
	"strings"
	"testing"
)

// TestUUIDBindingsNotDuplicated pins the av_uuid_* skip in
// internal/generator/generator.go. The seven av_uuid_* symbols are manually
// wrapped in uuid.go because AVUUID is an array typedef that needs explicit
// pointer conversion in CGO. If a future libclang fixes the array-typedef
// handling, the generator might also emit these symbols, double-defining them.
// This assertion fires when any of the seven Go wrapper names also appears as
// a generated function, signalling that the skip site is now spurious.
func TestUUIDBindingsNotDuplicated(t *testing.T) {
	data, err := os.ReadFile("functions.gen.go")
	if err != nil {
		t.Fatalf("read functions.gen.go: %v", err)
	}
	src := string(data)
	// Seven manual wrappers from uuid.go.
	names := []string{
		"AVUuidParse",
		"AVUuidUrnParse",
		"AVUuidParseRange",
		"AVUuidUnparse",
		"AVUuidEqual",
		"AVUuidCopy",
		"AVUuidNil",
	}
	for _, n := range names {
		if strings.Contains(src, "func "+n+"(") {
			t.Errorf("%s is generated in functions.gen.go alongside the manual wrapper in uuid.go; the av_uuid skip is now spurious", n)
		}
	}
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

		if _, err := AVUuidParse(cStr, &uuid1); err != nil {
			t.Fatalf("AVUuidParse failed: %v", err)
		}
		AVUuidCopy(&uuid2, &uuid1)

		// Test equality
		equal, _ := AVUuidEqual(&uuid1, &uuid2)
		if equal == 0 {
			t.Fatal("UUIDs should be equal")
		}

		// Test nil UUID
		var nilUUID AVUUID
		AVUuidNil(&nilUUID)

		equal, _ = AVUuidEqual(&uuid1, &nilUUID)
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

// =============================================================================
// Test: UUID Parse Error Paths
// =============================================================================

// TestUUID_ParseErrorPaths tests error detection in UUID parsing functions.
// This validates that malformed UUIDs are properly rejected and return error codes,
// preventing silent corruption from accepting invalid data.
func TestUUID_ParseErrorPaths(t *testing.T) {
	t.Run("invalid_uuid_format_too_short", func(t *testing.T) {
		// UUID must be exactly 36 characters (plus NUL)
		tooShort := "550e8400-e29b-41d4-a716-44665544"
		var uuid AVUUID

		cStr := ToCStr(tooShort)
		defer cStr.Free()

		ret, err := AVUuidParse(cStr, &uuid)
		if ret == 0 || err == nil {
			t.Errorf("AVUuidParse should fail for too-short UUID, got ret=%d err=%v", ret, err)
		}

		t.Logf("Correctly rejected too-short UUID: ret=%d", ret)
	})

	t.Run("invalid_uuid_format_too_long", func(t *testing.T) {
		// Too many characters
		tooLong := "550e8400-e29b-41d4-a716-446655440000-extra"
		var uuid AVUUID

		cStr := ToCStr(tooLong)
		defer cStr.Free()

		ret, err := AVUuidParse(cStr, &uuid)
		if ret == 0 || err == nil {
			t.Errorf("AVUuidParse should fail for too-long UUID, got ret=%d err=%v", ret, err)
		}

		t.Logf("Correctly rejected too-long UUID: ret=%d", ret)
	})

	t.Run("invalid_uuid_missing_dashes", func(t *testing.T) {
		// Correct length but missing dashes
		noDashes := "550e8400e29b41d4a716446655440000"
		var uuid AVUUID

		cStr := ToCStr(noDashes)
		defer cStr.Free()

		ret, err := AVUuidParse(cStr, &uuid)
		if ret == 0 || err == nil {
			t.Errorf("AVUuidParse should fail for UUID without dashes, got ret=%d err=%v", ret, err)
		}

		t.Logf("Correctly rejected UUID without dashes: ret=%d", ret)
	})

	t.Run("invalid_uuid_wrong_dash_positions", func(t *testing.T) {
		// Dashes in wrong positions
		wrongDashes := "550e8400e29b-41d4-a716-446655440000"
		var uuid AVUUID

		cStr := ToCStr(wrongDashes)
		defer cStr.Free()

		ret, err := AVUuidParse(cStr, &uuid)
		if ret == 0 || err == nil {
			t.Errorf("AVUuidParse should fail for UUID with wrong dash positions, got ret=%d err=%v", ret, err)
		}

		t.Logf("Correctly rejected UUID with wrong dashes: ret=%d", ret)
	})

	t.Run("invalid_uuid_non_hex_characters", func(t *testing.T) {
		// Contains non-hex characters
		invalidChars := "550e8400-e29b-41d4-a716-44665544GGGG"
		var uuid AVUUID

		cStr := ToCStr(invalidChars)
		defer cStr.Free()

		ret, err := AVUuidParse(cStr, &uuid)
		if ret == 0 || err == nil {
			t.Errorf("AVUuidParse should fail for non-hex characters, got ret=%d err=%v", ret, err)
		}

		t.Logf("Correctly rejected UUID with non-hex chars: ret=%d", ret)
	})

	t.Run("invalid_uuid_empty_string", func(t *testing.T) {
		// Empty string
		var uuid AVUUID

		cStr := ToCStr("")
		defer cStr.Free()

		ret, err := AVUuidParse(cStr, &uuid)
		if ret == 0 || err == nil {
			t.Errorf("AVUuidParse should fail for empty string, got ret=%d err=%v", ret, err)
		}

		t.Logf("Correctly rejected empty UUID string: ret=%d", ret)
	})

	t.Run("valid_uuid_uppercase_accepted", func(t *testing.T) {
		// RFC 4122 specifies parsing should be case-insensitive
		uppercaseUUID := "550E8400-E29B-41D4-A716-446655440000"
		var uuid AVUUID

		cStr := ToCStr(uppercaseUUID)
		defer cStr.Free()

		ret, err := AVUuidParse(cStr, &uuid)
		if ret != 0 || err != nil {
			t.Fatalf("AVUuidParse should accept uppercase UUID, got ret=%d err=%v", ret, err)
		}

		// Verify by unparsing (should be lowercase)
		outStr := AllocCStr(37)
		defer outStr.Free()
		AVUuidUnparse(&uuid, outStr)

		result := outStr.String()
		expected := "550e8400-e29b-41d4-a716-446655440000"
		if result != expected {
			t.Errorf("Unparsed UUID mismatch: got %s, want %s", result, expected)
		}

		t.Logf("Uppercase UUID correctly parsed and normalized: %s -> %s", uppercaseUUID, result)
	})

	t.Run("urn_invalid_format_missing_prefix", func(t *testing.T) {
		// URN must start with "urn:uuid:"
		noPrefix := "550e8400-e29b-41d4-a716-446655440000"
		var uuid AVUUID

		cStr := ToCStr(noPrefix)
		defer cStr.Free()

		ret, err := AVUuidUrnParse(cStr, &uuid)
		if ret == 0 || err == nil {
			t.Errorf("AVUuidUrnParse should fail without urn:uuid: prefix, got ret=%d err=%v", ret, err)
		}

		t.Logf("Correctly rejected URN without prefix: ret=%d", ret)
	})

	t.Run("urn_invalid_format_wrong_prefix", func(t *testing.T) {
		// Wrong URN prefix
		wrongPrefix := "uuid:uuid:550e8400-e29b-41d4-a716-446655440000"
		var uuid AVUUID

		cStr := ToCStr(wrongPrefix)
		defer cStr.Free()

		ret, err := AVUuidUrnParse(cStr, &uuid)
		if ret == 0 || err == nil {
			t.Errorf("AVUuidUrnParse should fail with wrong prefix, got ret=%d err=%v", ret, err)
		}

		t.Logf("Correctly rejected URN with wrong prefix: ret=%d", ret)
	})

	t.Run("urn_invalid_uuid_part", func(t *testing.T) {
		// URN with invalid UUID part
		invalidURN := "urn:uuid:550e8400-e29b-41d4-a716-GGGGGGGGGGGG"
		var uuid AVUUID

		cStr := ToCStr(invalidURN)
		defer cStr.Free()

		ret, err := AVUuidUrnParse(cStr, &uuid)
		if ret == 0 || err == nil {
			t.Errorf("AVUuidUrnParse should fail with invalid UUID part, got ret=%d err=%v", ret, err)
		}

		t.Logf("Correctly rejected URN with invalid UUID: ret=%d", ret)
	})

	t.Run("urn_case_insensitive", func(t *testing.T) {
		// URN parsing should be case-insensitive for the "urn:uuid:" prefix
		mixedCase := "URN:UUID:550e8400-e29b-41d4-a716-446655440000"
		var uuid AVUUID

		cStr := ToCStr(mixedCase)
		defer cStr.Free()

		// This may or may not be accepted depending on FFmpeg's strictness
		// Document the actual behaviour
		ret, _ := AVUuidUrnParse(cStr, &uuid)

		if ret == 0 {
			t.Logf("Mixed-case URN prefix accepted: %s", mixedCase)
		} else {
			t.Logf("Mixed-case URN prefix rejected (ret=%d) - prefix is case-sensitive", ret)
		}
	})

	t.Run("parse_range_exact_length", func(t *testing.T) {
		// AVUuidParseRange requires exact 36-character range
		uuidStr := "550e8400-e29b-41d4-a716-446655440000"
		var uuid AVUUID

		// Create start and end pointers into the UUID string
		start := ToCStr(uuidStr)
		defer start.Free()

		// End pointer should be exactly 36 chars after start
		ret, err := AVUuidParseRange(start, start, &uuid)
		if ret == 0 || err == nil {
			t.Errorf("AVUuidParseRange should fail with zero-length range, got ret=%d err=%v", ret, err)
		}

		t.Logf("Correctly rejected zero-length range: ret=%d", ret)
	})

	t.Run("error_code_is_negative", func(t *testing.T) {
		// All UUID parse errors should return negative codes
		testCases := []string{
			"",                                     // empty
			"invalid",                              // too short
			"550e8400-e29b-41d4-a716-44665544XXXX", // invalid chars
		}

		for _, testStr := range testCases {
			var uuid AVUUID
			cStr := ToCStr(testStr)
			ret, err := AVUuidParse(cStr, &uuid)
			cStr.Free()

			if ret >= 0 {
				t.Errorf("AVUuidParse should fail with a negative error code, got %d for %q", ret, testStr)
			}
			if ret < 0 && err == nil {
				t.Errorf("Negative return code should have non-nil error, got nil for %q", testStr)
			}
		}

		t.Log("All UUID parse errors correctly return negative codes")
	})
}
