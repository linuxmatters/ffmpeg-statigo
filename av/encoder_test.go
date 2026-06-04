package av

import (
	"testing"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
	"github.com/stretchr/testify/require"
)

// newTestVideoFrame allocates a writable YUV420P frame at the fixture size and
// fills it with the deterministic gradient at index n.
func newTestVideoFrame(t *testing.T, n int) *ffmpeg.AVFrame {
	t.Helper()

	frame := ffmpeg.AVFrameAlloc()
	require.NotNil(t, frame)

	frame.SetWidth(fixtureWidth)
	frame.SetHeight(fixtureHeight)
	frame.SetFormat(int(ffmpeg.AVPixFmtYuv420P))

	_, err := ffmpeg.AVFrameGetBuffer(frame, 0)
	require.NoError(t, err)

	fillFrame(frame, n)
	frame.SetPts(int64(n))

	return frame
}

func TestEncoderEncodeAndFlush(t *testing.T) {
	codec := ffmpeg.AVCodecFindEncoder(ffmpeg.AVCodecIdMpeg2Video)
	require.NotNil(t, codec, "mpeg2video encoder not found")

	enc, err := NewEncoder(codec, func(ctx *ffmpeg.AVCodecContext) {
		ctx.SetWidth(fixtureWidth)
		ctx.SetHeight(fixtureHeight)
		ctx.SetPixFmt(ffmpeg.AVPixFmtYuv420P)
		ctx.SetTimeBase(ffmpeg.AVMakeQ(fixtureRateDen, fixtureRateNum))
		ctx.SetFramerate(ffmpeg.AVMakeQ(fixtureRateNum, fixtureRateDen))
		ctx.SetGopSize(10)
		ctx.SetMaxBFrames(0)
	})
	require.NoError(t, err)
	defer enc.Close()

	var packets int
	for i := range fixtureFrames {
		frame := newTestVideoFrame(t, i)
		err := enc.Encode(frame, func(*ffmpeg.AVPacket) error {
			packets++
			return nil
		})
		ffmpeg.AVFrameFree(&frame)
		require.NoError(t, err)
	}

	var flushed int
	err = enc.Flush(func(*ffmpeg.AVPacket) error {
		flushed++
		return nil
	})
	require.NoError(t, err)

	require.Positive(t, packets+flushed, "encoder emitted no packets")
	require.LessOrEqual(t, packets, fixtureFrames, "more packets than frames before flush")
}

func TestEncoderDoubleClose(t *testing.T) {
	codec := ffmpeg.AVCodecFindEncoder(ffmpeg.AVCodecIdMpeg2Video)
	require.NotNil(t, codec)

	enc, err := NewEncoder(codec, func(ctx *ffmpeg.AVCodecContext) {
		ctx.SetWidth(fixtureWidth)
		ctx.SetHeight(fixtureHeight)
		ctx.SetPixFmt(ffmpeg.AVPixFmtYuv420P)
		ctx.SetTimeBase(ffmpeg.AVMakeQ(fixtureRateDen, fixtureRateNum))
	})
	require.NoError(t, err)

	require.NoError(t, enc.Close())
	require.NoError(t, enc.Close())
}
