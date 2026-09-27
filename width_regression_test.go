package ffmpeg_test

import (
	"strconv"
	"testing"

	"github.com/linuxmatters/ffmpeg-statigo"
)

func TestGeneratedSizeOutputWidth(t *testing.T) {
	pkt := ffmpeg.AVPacketAlloc()
	if pkt == nil {
		t.Fatal("AVPacketAlloc failed")
	}
	defer ffmpeg.AVPacketFree(&pkt)

	if ffmpeg.AVPacketNewSideData(pkt, ffmpeg.AVPktDataNewExtradata, 7) == nil {
		t.Fatal("AVPacketNewSideData failed")
	}
	size := ^uint64(0)
	if data := ffmpeg.AVPacketGetSideData(pkt, ffmpeg.AVPktDataNewExtradata, &size); data == nil || size != 7 {
		t.Fatalf("side data=%p size=%d, want non-nil and 7", data, size)
	}
	if data := ffmpeg.AVPacketGetSideData(pkt, ffmpeg.AVPktDataNewExtradata, nil); data == nil {
		t.Fatal("nil size output changed the result")
	}
	size = ^uint64(0)
	if data := ffmpeg.AVPacketGetSideData(pkt, ffmpeg.AVPktDataStereo3D, &size); data != nil || size != 0 {
		t.Fatalf("absent side data=%p size=%d, want nil and 0", data, size)
	}
}

func TestGeneratedSizeInputWidth(t *testing.T) {
	if strconv.IntSize != 32 {
		t.Skip("size_t narrowing requires a 32-bit target")
	}
	buf := ffmpeg.AVBufferAlloc(1)
	if buf == nil {
		t.Fatal("AVBufferAlloc failed")
	}
	defer ffmpeg.AVBufferUnref(&buf)

	for name, call := range map[string]func(){
		"scalar": func() { ffmpeg.AVFree(ffmpeg.AVMalloc(1 << 32)) },
		"setter": func() { buf.SetSize(1 << 32) },
	} {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Error("out-of-range size did not panic")
				}
			}()
			call()
		})
	}
	if got := buf.Size(); got != 1 {
		t.Fatalf("rejected setter changed size to %d", got)
	}
}
