package ffmpeg_test

import (
	"testing"
	"unsafe"

	"github.com/linuxmatters/ffmpeg-statigo"
)

func TestAllocCStrBuffer(t *testing.T) {
	const size = 64
	str := ffmpeg.AllocCStr(size)
	defer str.Free()

	if str.RawPtr() == nil {
		t.Fatal("AllocCStr returned a nil buffer")
	}
	buf := unsafe.Slice((*byte)(str.RawPtr()), size)
	for i, value := range buf {
		if value != 0 {
			t.Fatalf("buffer byte %d = %d, want 0", i, value)
		}
	}

	copy(buf, "allocated")
	if got := str.String(); got != "allocated" {
		t.Fatalf("String() = %q, want %q", got, "allocated")
	}

	str.Free()
	if str.RawPtr() != nil {
		t.Fatal("Free did not clear the buffer pointer")
	}
}
