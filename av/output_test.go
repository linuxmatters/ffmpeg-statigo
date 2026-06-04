package av

import (
	"path/filepath"
	"testing"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
	"github.com/stretchr/testify/require"
)

// newTestEncoder builds an mpeg2video encoder at the fixture geometry. mpeg4 is
// absent from this build, so mpeg2video is the portable choice.
func newTestEncoder(t *testing.T) *Encoder {
	t.Helper()

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

	return enc
}

func TestOutputWriteAndReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.mp4")

	out, err := CreateOutput(path)
	require.NoError(t, err)
	defer out.Close()

	enc := newTestEncoder(t)
	defer enc.Close()

	stream, err := out.AddStream(enc)
	require.NoError(t, err)
	require.NotNil(t, stream)

	require.NoError(t, out.WriteHeader())

	for i := range fixtureFrames {
		frame := newTestVideoFrame(t, i)
		err := enc.Encode(frame, func(pkt *ffmpeg.AVPacket) error {
			pkt.SetStreamIndex(0)
			ffmpeg.AVPacketRescaleTs(pkt, enc.Raw().TimeBase(), stream.TimeBase())
			return out.WritePacket(pkt)
		})
		ffmpeg.AVFrameFree(&frame)
		require.NoError(t, err)
	}

	require.NoError(t, enc.Flush(func(pkt *ffmpeg.AVPacket) error {
		pkt.SetStreamIndex(0)
		ffmpeg.AVPacketRescaleTs(pkt, enc.Raw().TimeBase(), stream.TimeBase())
		return out.WritePacket(pkt)
	}))

	require.NoError(t, out.WriteTrailer())
	require.NoError(t, out.Close())

	in, err := Open(path)
	require.NoError(t, err)
	defer in.Close()

	require.Equal(t, uint(1), in.NbStreams())
}

func TestOutputDoubleClose(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.mp4")

	out, err := CreateOutput(path)
	require.NoError(t, err)

	enc := newTestEncoder(t)
	defer enc.Close()

	_, err = out.AddStream(enc)
	require.NoError(t, err)

	require.NoError(t, out.WriteHeader())
	require.NoError(t, out.WriteTrailer())

	require.NoError(t, out.Close())
	require.NoError(t, out.Close())
}
