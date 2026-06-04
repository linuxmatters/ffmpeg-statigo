package av

import (
	"testing"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
	"github.com/stretchr/testify/require"
)

func TestFilterGraphVideo(t *testing.T) {
	const (
		width  = 160
		height = 120
	)

	timeBase := ffmpeg.AVMakeQ(1, 25)
	sar := ffmpeg.AVMakeQ(1, 1)

	g, err := NewVideoFilterGraph(VideoFilterParams{
		Width:             width,
		Height:            height,
		PixFmt:            ffmpeg.AVPixFmtYuv420P,
		TimeBase:          timeBase,
		SampleAspectRatio: sar,
		OutPixFmt:         ffmpeg.AVPixFmtYuv420P,
	}, "null")
	require.NoError(t, err)
	defer g.Close()

	frame := ffmpeg.AVFrameAlloc()
	require.NotNil(t, frame)
	defer ffmpeg.AVFrameFree(&frame)

	frame.SetWidth(width)
	frame.SetHeight(height)
	frame.SetFormat(int(ffmpeg.AVPixFmtYuv420P))

	_, err = ffmpeg.AVFrameGetBuffer(frame, 0)
	require.NoError(t, err)

	frame.SetPts(0)

	require.NoError(t, g.Push(frame))

	var got int

	require.NoError(t, g.Pull(func(out *ffmpeg.AVFrame) error {
		got++
		require.Equal(t, width, out.Width())
		require.Equal(t, height, out.Height())
		return nil
	}))

	require.GreaterOrEqual(t, got, 1)
}

func TestFilterGraphAudio(t *testing.T) {
	const (
		sampleRate = 44100
		nbSamples  = 1024
		channels   = 2
	)

	timeBase := ffmpeg.AVMakeQ(1, sampleRate)

	frame := ffmpeg.AVFrameAlloc()
	require.NotNil(t, frame)
	defer ffmpeg.AVFrameFree(&frame)

	frame.SetNbSamples(nbSamples)
	frame.SetFormat(int(ffmpeg.AVSampleFmtFltp))
	frame.SetSampleRate(sampleRate)
	ffmpeg.AVChannelLayoutDefault(frame.ChLayout(), channels)

	g, err := NewAudioFilterGraph(AudioFilterParams{
		TimeBase:     timeBase,
		SampleRate:   sampleRate,
		SampleFmt:    ffmpeg.AVSampleFmtFltp,
		ChLayout:     frame.ChLayout(),
		OutSampleFmt: ffmpeg.AVSampleFmtFltp,
	}, "anull")
	require.NoError(t, err)
	defer g.Close()

	_, err = ffmpeg.AVFrameGetBuffer(frame, 0)
	require.NoError(t, err)

	frame.SetPts(0)

	require.NoError(t, g.Push(frame))

	var got int

	require.NoError(t, g.Pull(func(out *ffmpeg.AVFrame) error {
		got++
		require.Equal(t, nbSamples, out.NbSamples())
		return nil
	}))

	require.GreaterOrEqual(t, got, 1)
}

func TestFilterGraphDoubleClose(t *testing.T) {
	g, err := NewVideoFilterGraph(VideoFilterParams{
		Width:             160,
		Height:            120,
		PixFmt:            ffmpeg.AVPixFmtYuv420P,
		TimeBase:          ffmpeg.AVMakeQ(1, 25),
		SampleAspectRatio: ffmpeg.AVMakeQ(1, 1),
		OutPixFmt:         ffmpeg.AVPixFmtYuv420P,
	}, "null")
	require.NoError(t, err)

	require.NoError(t, g.Close())
	require.NoError(t, g.Close())
	require.Nil(t, g.Raw())
}
