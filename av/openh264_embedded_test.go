//go:build embedded

package av

import (
	"path/filepath"
	"testing"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
	"github.com/stretchr/testify/require"
)

func TestEmbeddedOpenH264MP4(t *testing.T) {
	codec := ffmpeg.AVCodecFindEncoderByName(ffmpeg.GlobalCStr("libopenh264"))
	require.NotNil(t, codec)

	enc, err := NewEncoder(codec, func(ctx *ffmpeg.AVCodecContext) {
		ctx.SetWidth(fixtureWidth)
		ctx.SetHeight(fixtureHeight)
		ctx.SetPixFmt(ffmpeg.AVPixFmtYuv420P)
		ctx.SetTimeBase(ffmpeg.AVMakeQ(fixtureRateDen, fixtureRateNum))
		ctx.SetFramerate(ffmpeg.AVMakeQ(fixtureRateNum, fixtureRateDen))
		ctx.SetFlags(ctx.Flags() | ffmpeg.AVCodecFlagGlobalHeader)
	})
	require.NoError(t, err)
	defer enc.Close()

	path := filepath.Join(t.TempDir(), "openh264.mp4")
	out, err := CreateOutput(path)
	require.NoError(t, err)
	defer out.Close()
	stream, err := out.AddStream(enc)
	require.NoError(t, err)
	require.NoError(t, out.WriteHeader())

	writePacket := func(pkt *ffmpeg.AVPacket) error {
		pkt.SetStreamIndex(0)
		ffmpeg.AVPacketRescaleTs(pkt, enc.Raw().TimeBase(), stream.TimeBase())
		return out.WritePacket(pkt)
	}
	for i := range fixtureFrames {
		frame := newTestVideoFrame(t, i)
		err := enc.Encode(frame, writePacket)
		ffmpeg.AVFrameFree(&frame)
		require.NoError(t, err)
	}
	require.NoError(t, enc.Flush(writePacket))
	require.NoError(t, out.WriteTrailer())
	require.NoError(t, out.Close())

	in, err := Open(path)
	require.NoError(t, err)
	defer in.Close()
	video, err := in.BestStream(ffmpeg.AVMediaTypeVideo)
	require.NoError(t, err)
	require.Equal(t, ffmpeg.AVCodecIdH264, video.Codecpar().CodecId())
	require.Positive(t, video.Codecpar().ExtradataSize(), "MP4 must contain H.264 codec configuration")

	frames, width, height := videoStats(t, path)
	require.Positive(t, frames, "MP4 video must decode")
	require.Equal(t, fixtureWidth, width)
	require.Equal(t, fixtureHeight, height)
}
