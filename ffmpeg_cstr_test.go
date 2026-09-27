package ffmpeg

import "testing"

func TestCStrDuplicateAllocator(t *testing.T) {
	const text = "duplicate string"
	source := ToCStr(text)
	defer source.Free()

	for _, tt := range []struct {
		name string
		dup  func() *CStr
		want string
	}{
		{"AVStrdup", func() *CStr { return AVStrdup(source) }, text},
		{"Dup", source.Dup, text},
		{"AVStrndup_truncated", func() *CStr { return AVStrndup(source, 9) }, "duplicate"},
		{"AVStrndup_full", func() *CStr { return AVStrndup(source, 100) }, text},
		{"AVStrndup_empty", func() *CStr { return AVStrndup(source, 0) }, ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.dup()
			if got == nil {
				t.Fatal("duplicate returned nil")
			}
			defer got.Free()
			if !got.avFree || got.dontFree {
				t.Error("duplicate must own its memory and use av_free")
			}
			if got.String() != tt.want {
				t.Errorf("duplicate = %q, want %q", got.String(), tt.want)
			}
			if got.RawPtr() == source.RawPtr() {
				t.Fatal("duplicate shares its allocation with the source")
			}
			got.Free()
			if got.RawPtr() != nil {
				t.Error("Free did not clear the pointer")
			}
			if source.String() != text {
				t.Error("Free changed the source string")
			}
		})
	}
}

func TestCStrWrapperOwnership(t *testing.T) {
	if wrapAVCStr(nil) != nil {
		t.Error("wrapAVCStr(nil) must return nil")
	}

	source := ToCStr("borrowed string")
	defer source.Free()
	if source.avFree {
		t.Error("ToCStr must use C.free")
	}
	borrowed := wrapStaticCStr(source.ptr)
	borrowed.Free()
	if borrowed.avFree || !borrowed.dontFree {
		t.Error("static wrapper must remain non-freeable")
	}
	if borrowed.RawPtr() != source.RawPtr() || borrowed.String() != "borrowed string" {
		t.Error("Free changed the borrowed string")
	}
}
