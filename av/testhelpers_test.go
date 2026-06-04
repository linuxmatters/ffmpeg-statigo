package av

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"unsafe"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
	"github.com/stretchr/testify/require"
)

const (
	fixtureWidth   = 160
	fixtureHeight  = 120
	fixtureFrames  = 15
	fixtureRateNum = 25
	fixtureRateDen = 1
)

var (
	fixtureOnce sync.Once
	fixturePath string
	fixtureErr  error
)

// testFixture synthesises a small, valid, video-only media file once per test
// process and returns its path. It uses the built-in, licence-clean mpeg2video
// encoder in an .mp4 container, written with the raw ffmpeg bindings. mpeg4 is
// preferred by the plan but is not compiled into this static library, so the
// helper falls back to mpeg2video (also built in, GPL-clean). The file is
// generated under a process-wide temp dir so every test reuses it.
func testFixture(t *testing.T) string {
	t.Helper()

	fixtureOnce.Do(func() {
		dir, err := os.MkdirTemp("", "av-fixture-")
		if err != nil {
			fixtureErr = fmt.Errorf("create temp dir: %w", err)
			return
		}

		path := filepath.Join(dir, "fixture.mp4")
		if err := writeFixture(path); err != nil {
			fixtureErr = err
			return
		}

		fixturePath = path
	})

	if fixtureErr != nil {
		t.Fatal(fixtureErr)
	}

	return fixturePath
}

// writeFixture encodes synthetic frames into path. Free ordering mirrors
// examples/transcode: trailer, then IO close (guarded by AVFmtNofile), packet,
// frame, codec context, format context.
func writeFixture(path string) (retErr error) {
	encoder := ffmpeg.AVCodecFindEncoder(ffmpeg.AVCodecIdMpeg2Video)
	if encoder == nil {
		return errors.New("mpeg2video encoder not found")
	}

	encCtx := ffmpeg.AVCodecAllocContext3(encoder)
	if encCtx == nil {
		return errors.New("alloc codec context failed")
	}
	defer ffmpeg.AVCodecFreeContext(&encCtx)

	encCtx.SetWidth(fixtureWidth)
	encCtx.SetHeight(fixtureHeight)
	encCtx.SetPixFmt(ffmpeg.AVPixFmtYuv420P)
	encCtx.SetTimeBase(ffmpeg.AVMakeQ(fixtureRateDen, fixtureRateNum))
	encCtx.SetFramerate(ffmpeg.AVMakeQ(fixtureRateNum, fixtureRateDen))
	encCtx.SetGopSize(10)
	encCtx.SetMaxBFrames(0)

	namePtr := ffmpeg.ToCStr(path)
	defer namePtr.Free()

	var ofmtCtx *ffmpeg.AVFormatContext
	if _, err := ffmpeg.AVFormatAllocOutputContext2(&ofmtCtx, nil, nil, namePtr); err != nil {
		return fmt.Errorf("alloc output context: %w", err)
	}
	defer ffmpeg.AVFormatFreeContext(ofmtCtx)

	if ofmtCtx.Oformat().Flags()&ffmpeg.AVFmtGlobalheader != 0 {
		encCtx.SetFlags(encCtx.Flags() | ffmpeg.AVCodecFlagGlobalHeader)
	}

	if _, err := ffmpeg.AVCodecOpen2(encCtx, encoder, nil); err != nil {
		return fmt.Errorf("open encoder: %w", err)
	}

	stream := ffmpeg.AVFormatNewStream(ofmtCtx, nil)
	if stream == nil {
		return errors.New("new stream failed")
	}

	if _, err := ffmpeg.AVCodecParametersFromContext(stream.Codecpar(), encCtx); err != nil {
		return fmt.Errorf("copy codec params: %w", err)
	}
	stream.SetTimeBase(encCtx.TimeBase())
	stream.SetAvgFrameRate(ffmpeg.AVMakeQ(fixtureRateNum, fixtureRateDen))
	stream.SetRFrameRate(ffmpeg.AVMakeQ(fixtureRateNum, fixtureRateDen))

	if ofmtCtx.Oformat().Flags()&ffmpeg.AVFmtNofile == 0 {
		var pb *ffmpeg.AVIOContext
		if _, err := ffmpeg.AVIOOpen(&pb, namePtr, ffmpeg.AVIOFlagWrite); err != nil {
			return fmt.Errorf("open io: %w", err)
		}
		ofmtCtx.SetPb(pb)

		defer func() {
			if _, err := ffmpeg.AVIOClose(ofmtCtx.Pb()); err != nil && retErr == nil {
				retErr = fmt.Errorf("close io: %w", err)
			}
			ofmtCtx.SetPb(nil)
		}()
	}

	if _, err := ffmpeg.AVFormatWriteHeader(ofmtCtx, nil); err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	frame := ffmpeg.AVFrameAlloc()
	if frame == nil {
		return errors.New("alloc frame failed")
	}
	defer ffmpeg.AVFrameFree(&frame)

	frame.SetWidth(fixtureWidth)
	frame.SetHeight(fixtureHeight)
	frame.SetFormat(int(ffmpeg.AVPixFmtYuv420P))

	if _, err := ffmpeg.AVFrameGetBuffer(frame, 0); err != nil {
		return fmt.Errorf("frame get buffer: %w", err)
	}

	pkt := ffmpeg.AVPacketAlloc()
	if pkt == nil {
		return errors.New("alloc packet failed")
	}
	defer ffmpeg.AVPacketFree(&pkt)

	for i := range fixtureFrames {
		fillFrame(frame, i)
		frame.SetPts(int64(i))

		if err := encodeFrame(encCtx, ofmtCtx, stream, pkt, frame); err != nil {
			return err
		}
	}

	if err := encodeFrame(encCtx, ofmtCtx, stream, pkt, nil); err != nil {
		return err
	}

	if _, err := ffmpeg.AVWriteTrailer(ofmtCtx); err != nil {
		return fmt.Errorf("write trailer: %w", err)
	}

	return nil
}

// encodeFrame sends frame (nil flushes) to the encoder and muxes every packet
// it yields, draining EAgain/EOF internally like examples/transcode.
func encodeFrame(encCtx *ffmpeg.AVCodecContext, ofmtCtx *ffmpeg.AVFormatContext, stream *ffmpeg.AVStream, pkt *ffmpeg.AVPacket, frame *ffmpeg.AVFrame) error {
	if _, err := ffmpeg.AVCodecSendFrame(encCtx, frame); err != nil {
		return fmt.Errorf("send frame: %w", err)
	}

	for {
		if _, err := ffmpeg.AVCodecReceivePacket(encCtx, pkt); err != nil {
			if errors.Is(err, ffmpeg.EAgain) || errors.Is(err, ffmpeg.AVErrorEOF) {
				return nil
			}

			return fmt.Errorf("receive packet: %w", err)
		}

		ffmpeg.AVPacketRescaleTs(pkt, encCtx.TimeBase(), stream.TimeBase())

		if _, err := ffmpeg.AVInterleavedWriteFrame(ofmtCtx, pkt); err != nil {
			ffmpeg.AVPacketUnref(pkt)
			return fmt.Errorf("write frame: %w", err)
		}

		ffmpeg.AVPacketUnref(pkt)
	}
}

// fillFrame writes a deterministic moving gradient into the YUV420P planes.
func fillFrame(frame *ffmpeg.AVFrame, n int) {
	data := frame.Data()
	linesize := frame.Linesize()

	yStride := linesize.Get(0)
	uStride := linesize.Get(1)
	vStride := linesize.Get(2)

	yPlane := planeSlice(data.Get(0), yStride*fixtureHeight)
	uPlane := planeSlice(data.Get(1), uStride*(fixtureHeight/2))
	vPlane := planeSlice(data.Get(2), vStride*(fixtureHeight/2))

	for y := range fixtureHeight {
		row := y * yStride
		for x := range fixtureWidth {
			yPlane[row+x] = byte((x + y + n*3) & 0xff)
		}
	}

	for y := range fixtureHeight / 2 {
		uRow := y * uStride
		vRow := y * vStride
		for x := range fixtureWidth / 2 {
			uPlane[uRow+x] = byte((128 + x + n) & 0xff)
			vPlane[vRow+x] = byte((64 + y + n) & 0xff)
		}
	}
}

func planeSlice(ptr unsafe.Pointer, length int) []byte {
	return unsafe.Slice((*byte)(ptr), length)
}

// TestFixtureOpens is a smoke test: it opens the synthesised fixture with the
// raw bindings, reads stream info and asserts at least one stream is present.
func TestFixtureOpens(t *testing.T) {
	path := testFixture(t)

	urlPtr := ffmpeg.ToCStr(path)
	defer urlPtr.Free()

	var fmtCtx *ffmpeg.AVFormatContext
	_, err := ffmpeg.AVFormatOpenInput(&fmtCtx, urlPtr, nil, nil)
	require.NoError(t, err)
	defer ffmpeg.AVFormatCloseInput(&fmtCtx)

	_, err = ffmpeg.AVFormatFindStreamInfo(fmtCtx, nil)
	require.NoError(t, err)

	require.GreaterOrEqual(t, fmtCtx.NbStreams(), uint(1))
}
