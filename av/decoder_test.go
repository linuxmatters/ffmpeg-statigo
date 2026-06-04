package av

import (
	"errors"
	"testing"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
	"github.com/stretchr/testify/require"
)

func TestDecoder(t *testing.T) {
	in, err := Open(testFixture(t))
	require.NoError(t, err)
	defer in.Close()

	stream, err := in.BestStream(ffmpeg.AVMediaTypeVideo)
	require.NoError(t, err)

	dec, err := NewDecoder(stream)
	require.NoError(t, err)
	defer dec.Close()

	pkt := ffmpeg.AVPacketAlloc()
	require.NotNil(t, pkt)
	defer ffmpeg.AVPacketFree(&pkt)

	decoded := 0
	count := func(*ffmpeg.AVFrame) error {
		decoded++
		return nil
	}

	for {
		if err := in.ReadPacket(pkt); err != nil {
			require.ErrorIs(t, err, ffmpeg.AVErrorEOF)
			break
		}

		if pkt.StreamIndex() == stream.Index() {
			require.NoError(t, dec.Decode(pkt, count))
		}

		ffmpeg.AVPacketUnref(pkt)
	}

	require.NoError(t, dec.Flush(count))
	require.Greater(t, decoded, 0)
}

func TestDecoderCallbackErrorPropagates(t *testing.T) {
	in, err := Open(testFixture(t))
	require.NoError(t, err)
	defer in.Close()

	stream, err := in.BestStream(ffmpeg.AVMediaTypeVideo)
	require.NoError(t, err)

	dec, err := NewDecoder(stream)
	require.NoError(t, err)
	defer dec.Close()

	pkt := ffmpeg.AVPacketAlloc()
	require.NotNil(t, pkt)
	defer ffmpeg.AVPacketFree(&pkt)

	sentinel := errors.New("callback boom")
	fail := func(*ffmpeg.AVFrame) error { return sentinel }

	var got error
	for {
		if err := in.ReadPacket(pkt); err != nil {
			require.ErrorIs(t, err, ffmpeg.AVErrorEOF)
			break
		}

		if pkt.StreamIndex() == stream.Index() {
			if err := dec.Decode(pkt, fail); err != nil {
				got = err
				ffmpeg.AVPacketUnref(pkt)
				break
			}
		}

		ffmpeg.AVPacketUnref(pkt)
	}

	require.ErrorIs(t, got, sentinel)
}

func TestDecoderDoubleClose(t *testing.T) {
	in, err := Open(testFixture(t))
	require.NoError(t, err)
	defer in.Close()

	stream, err := in.BestStream(ffmpeg.AVMediaTypeVideo)
	require.NoError(t, err)

	dec, err := NewDecoder(stream)
	require.NoError(t, err)

	require.NoError(t, dec.Close())
	require.NoError(t, dec.Close())
}
