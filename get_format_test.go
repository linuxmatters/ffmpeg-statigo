package ffmpeg_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
)

// TestSetGetFormatDecode opens a decoder over the synthesised fixture, installs a
// get_format callback, decodes a packet and asserts the trampoline fired with a
// non-empty offered list. libavcodec calls get_format even for software decode,
// so this exercises the bridge without a GPU. ClearGetFormat then unregisters,
// and a follow-up decode confirms the cleared path does not panic.
func TestSetGetFormatDecode(t *testing.T) {
	data := fixtureBytes(t)
	reader := bytes.NewReader(data)

	const bufSize = 4096

	avio := ffmpeg.AVIOAllocContext(bufSize, false,
		func(buf []byte) int {
			n, err := reader.Read(buf)
			if n == 0 && errors.Is(err, io.EOF) {
				return ffmpeg.AVErrorEofConst
			}
			return n
		},
		nil,
		func(offset int64, whence int) int64 {
			if whence&ffmpeg.AVSeekSize != 0 {
				return reader.Size()
			}
			pos, err := reader.Seek(offset, whence)
			if err != nil {
				return int64(ffmpeg.AVErrorEofConst)
			}
			return pos
		},
	)
	if avio == nil {
		t.Fatal("AVIOAllocContext returned nil")
	}
	defer avio.Close()

	fmtCtx := ffmpeg.AVFormatAllocContext()
	if fmtCtx == nil {
		t.Fatal("alloc format context failed")
	}

	fmtCtx.SetPb(avio.Context())
	fmtCtx.SetFlags(fmtCtx.Flags() | ffmpeg.AVFmtFlagCustomIo)

	if _, err := ffmpeg.AVFormatOpenInput(&fmtCtx, nil, nil, nil); err != nil {
		t.Fatalf("open input: %v", err)
	}
	defer ffmpeg.AVFormatCloseInput(&fmtCtx)

	if _, err := ffmpeg.AVFormatFindStreamInfo(fmtCtx, nil); err != nil {
		t.Fatalf("find stream info: %v", err)
	}

	if fmtCtx.NbStreams() < 1 {
		t.Fatal("expected at least one stream")
	}
	stream := fmtCtx.Streams().Get(0)

	codecpar := stream.Codecpar()
	decoder := ffmpeg.AVCodecFindDecoder(codecpar.CodecId())
	if decoder == nil {
		t.Skip("decoder for fixture codec not built into static library")
	}

	decCtx := ffmpeg.AVCodecAllocContext3(decoder)
	if decCtx == nil {
		t.Fatal("alloc decoder context failed")
	}
	defer ffmpeg.AVCodecFreeContext(&decCtx)

	if _, err := ffmpeg.AVCodecParametersToContext(decCtx, codecpar); err != nil {
		t.Fatalf("copy params to context: %v", err)
	}

	var (
		fired         bool
		offeredNonNil bool
		chosen        ffmpeg.AVPixelFormat
	)
	decCtx.SetGetFormat(func(_ *ffmpeg.AVCodecContext, formats []ffmpeg.AVPixelFormat) ffmpeg.AVPixelFormat {
		fired = true
		if len(formats) > 0 {
			offeredNonNil = true
			chosen = formats[0]
			return formats[0]
		}
		return ffmpeg.AVPixFmtNone
	})

	if _, err := ffmpeg.AVCodecOpen2(decCtx, decoder, nil); err != nil {
		decCtx.ClearGetFormat()
		t.Fatalf("open decoder: %v", err)
	}

	pkt := ffmpeg.AVPacketAlloc()
	if pkt == nil {
		t.Fatal("alloc packet failed")
	}
	defer ffmpeg.AVPacketFree(&pkt)

	frame := ffmpeg.AVFrameAlloc()
	if frame == nil {
		t.Fatal("alloc frame failed")
	}
	defer ffmpeg.AVFrameFree(&frame)

	decoded := false
	for !decoded {
		if _, err := ffmpeg.AVReadFrame(fmtCtx, pkt); err != nil {
			if errors.Is(err, ffmpeg.AVErrorEOF) {
				break
			}
			t.Fatalf("read frame: %v", err)
		}

		if pkt.StreamIndex() != 0 {
			ffmpeg.AVPacketUnref(pkt)
			continue
		}

		if _, err := ffmpeg.AVCodecSendPacket(decCtx, pkt); err != nil {
			ffmpeg.AVPacketUnref(pkt)
			t.Fatalf("send packet: %v", err)
		}
		ffmpeg.AVPacketUnref(pkt)

		for {
			if _, err := ffmpeg.AVCodecReceiveFrame(decCtx, frame); err != nil {
				if errors.Is(err, ffmpeg.EAgain) || errors.Is(err, ffmpeg.AVErrorEOF) {
					break
				}
				t.Fatalf("receive frame: %v", err)
			}
			decoded = true
			ffmpeg.AVFrameUnref(frame)
			break
		}
	}

	if !fired {
		t.Fatal("get_format callback did not fire")
	}
	if !offeredNonNil {
		t.Fatal("get_format offered an empty format list")
	}
	if chosen == ffmpeg.AVPixFmtNone {
		t.Fatal("get_format chose AV_PIX_FMT_NONE")
	}
	if !decoded {
		t.Fatal("decoder produced no frame")
	}

	// ClearGetFormat unregisters; a second call is idempotent and must not panic.
	decCtx.ClearGetFormat()
	decCtx.ClearGetFormat()
}

// TestSetGetFormatNilSafe confirms the register/unregister API tolerates nil
// receivers, nil pointers and a nil callback without panicking.
func TestSetGetFormatNilSafe(t *testing.T) {
	var nilCtx *ffmpeg.AVCodecContext
	nilCtx.SetGetFormat(func(_ *ffmpeg.AVCodecContext, _ []ffmpeg.AVPixelFormat) ffmpeg.AVPixelFormat {
		return ffmpeg.AVPixFmtNone
	})
	nilCtx.ClearGetFormat()

	codec := ffmpeg.AVCodecFindDecoder(ffmpeg.AVCodecIdMpeg2Video)
	if codec == nil {
		t.Skip("mpeg2video decoder not built into static library")
	}

	ctx := ffmpeg.AVCodecAllocContext3(codec)
	if ctx == nil {
		t.Fatal("alloc codec context failed")
	}
	defer ffmpeg.AVCodecFreeContext(&ctx)

	ctx.SetGetFormat(nil)
	ctx.ClearGetFormat()
}
