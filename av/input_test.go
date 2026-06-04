package av

import (
	"testing"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
	"github.com/stretchr/testify/require"
)

func TestOpenSuccess(t *testing.T) {
	in, err := Open(testFixture(t))
	require.NoError(t, err)
	defer in.Close()

	require.GreaterOrEqual(t, in.NbStreams(), uint(1))
}

func TestOpenFailure(t *testing.T) {
	in, err := Open("/nonexistent/path/does-not-exist.mp4")
	require.Error(t, err)
	require.Nil(t, in)
}

func TestBestStream(t *testing.T) {
	in, err := Open(testFixture(t))
	require.NoError(t, err)
	defer in.Close()

	stream, err := in.BestStream(ffmpeg.AVMediaTypeVideo)
	require.NoError(t, err)
	require.NotNil(t, stream)
}

func TestReadPacketToEOF(t *testing.T) {
	in, err := Open(testFixture(t))
	require.NoError(t, err)
	defer in.Close()

	pkt := ffmpeg.AVPacketAlloc()
	require.NotNil(t, pkt)
	defer ffmpeg.AVPacketFree(&pkt)

	var count int
	var readErr error
	for {
		readErr = in.ReadPacket(pkt)
		if readErr != nil {
			break
		}
		count++
		ffmpeg.AVPacketUnref(pkt)
	}

	require.ErrorIs(t, readErr, ffmpeg.AVErrorEOF)
	require.GreaterOrEqual(t, count, 1)
}

func TestCloseIdempotent(t *testing.T) {
	in, err := Open(testFixture(t))
	require.NoError(t, err)

	require.NoError(t, in.Close())
	require.NotPanics(t, func() {
		require.NoError(t, in.Close())
	})
}
