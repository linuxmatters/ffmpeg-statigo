package av

import (
	"testing"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
	"github.com/stretchr/testify/require"
)

func TestHWDeviceTypes(t *testing.T) {
	types := HWDeviceTypes()
	require.IsType(t, []string{}, types)
}

func TestNewHWDeviceUnknownType(t *testing.T) {
	dev, err := NewHWDevice("definitely-not-a-real-type")
	require.Error(t, err)
	require.Nil(t, dev)
}

func TestNewHWDeviceRealType(t *testing.T) {
	types := HWDeviceTypes()
	if len(types) == 0 {
		t.Skip("no hw device types compiled in")
	}

	// A real device type may or may not be usable on this host. Either outcome
	// is acceptable; the contract is that it never panics.
	dev, err := NewHWDevice(types[0])
	if err != nil {
		require.Nil(t, dev)
		return
	}

	require.NotNil(t, dev)
	require.NoError(t, dev.Close())
	require.NoError(t, dev.Close())
}

func TestTransferToSoftwareNil(t *testing.T) {
	frame, err := TransferToSoftware(nil)
	require.Error(t, err)
	require.Nil(t, frame)
}

func TestNewHWDecoderNilArgs(t *testing.T) {
	dev := &HWDevice{devType: ffmpeg.AVHWDeviceTypeNone}

	dec, err := NewHWDecoder(nil, dev)
	require.Error(t, err)
	require.Nil(t, dec)

	in, err := Open(testFixture(t))
	require.NoError(t, err)
	defer in.Close()

	stream, err := in.BestStream(ffmpeg.AVMediaTypeVideo)
	require.NoError(t, err)

	dec, err = NewHWDecoder(stream, nil)
	require.Error(t, err)
	require.Nil(t, dec)
}
