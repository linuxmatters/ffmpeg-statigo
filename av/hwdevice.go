package av

import (
	"errors"
	"fmt"
	"io"
	"slices"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
)

// HWDevice owns a hardware device context [ffmpeg.AVBufferRef] together with the
// [ffmpeg.AVHWDeviceType] it was created for. It implements [io.Closer]; Close
// is idempotent and unrefs the device context. A decoder built from the same
// device through NewHWDecoder takes its own independent reference, so closing the
// HWDevice and the Decoder in any order is safe.
type HWDevice struct {
	ref     *ffmpeg.AVBufferRef
	devType ffmpeg.AVHWDeviceType
}

var _ io.Closer = (*HWDevice)(nil)

// NewHWDevice resolves typeName to an [ffmpeg.AVHWDeviceType] and creates a
// device context for it with default device and options. typeName must name a
// compiled-in device type (see HWDeviceTypes); an unknown name is an error. A
// missing driver or GPU surfaces as a normal wrapped error from device creation,
// not a panic.
func NewHWDevice(typeName string) (*HWDevice, error) {
	namePtr := ffmpeg.ToCStr(typeName)
	defer namePtr.Free()

	devType := ffmpeg.AVHWDeviceFindTypeByName(namePtr)
	if devType == ffmpeg.AVHWDeviceTypeNone {
		return nil, fmt.Errorf("unknown hw device type %q", typeName)
	}

	var ref *ffmpeg.AVBufferRef
	if _, err := ffmpeg.AVHWDeviceCtxCreate(&ref, devType, nil, nil, 0); err != nil {
		return nil, fmt.Errorf("create hw device %q: %w", typeName, err)
	}

	return &HWDevice{ref: ref, devType: devType}, nil
}

// HWDeviceTypes returns the names of every hardware device type compiled into
// the linked FFmpeg. It enumerates the registry and needs no GPU, so it is safe
// to call on hardware-free hosts. The result may be empty.
func HWDeviceTypes() []string {
	var types []string

	for dt := ffmpeg.AVHWDeviceIterateTypes(ffmpeg.AVHWDeviceTypeNone); dt != ffmpeg.AVHWDeviceTypeNone; dt = ffmpeg.AVHWDeviceIterateTypes(dt) {
		types = append(types, ffmpeg.AVHWDeviceGetTypeName(dt).String())
	}

	return types
}

// Type returns the device type this HWDevice was created for.
func (d *HWDevice) Type() ffmpeg.AVHWDeviceType {
	return d.devType
}

// Raw returns the underlying device context. Do not unref it; Close owns it.
// Pass it through [ffmpeg.AVBufferRef_] to obtain an independent reference for
// another consumer.
func (d *HWDevice) Raw() *ffmpeg.AVBufferRef {
	return d.ref
}

// Close unrefs the device context and nils the handle. A second call is a no-op.
func (d *HWDevice) Close() error {
	if d.ref != nil {
		ffmpeg.AVBufferUnref(&d.ref)
		d.ref = nil
	}

	return nil
}

// NewHWDecoder builds a Decoder for stream that decodes on dev. It mirrors
// NewDecoder's allocate, configure, open and unwind-via-Close flow, and before
// opening the codec it wires up hardware decode: it discovers the decoder's
// hardware pixel format for dev's type, attaches an independent reference to
// dev's device context, and installs a get_format callback that selects that
// hardware format. The returned Decoder yields frames in the hardware pixel
// format; read them back with TransferToSoftware.
func NewHWDecoder(stream *ffmpeg.AVStream, dev *HWDevice) (*Decoder, error) {
	if stream == nil {
		return nil, errors.New("nil stream")
	}

	if dev == nil {
		return nil, errors.New("nil hw device")
	}

	codec := ffmpeg.AVCodecFindDecoder(stream.Codecpar().CodecId())
	if codec == nil {
		return nil, fmt.Errorf("find decoder (codec %v): %w", stream.Codecpar().CodecId(), ffmpeg.WrapErr(ffmpeg.AVErrorDecoderNotFoundConst))
	}

	ctx := ffmpeg.AVCodecAllocContext3(codec)
	if ctx == nil {
		return nil, errors.New("alloc decoder context failed")
	}

	d := &Decoder{ctx: ctx}

	if _, err := ffmpeg.AVCodecParametersToContext(ctx, stream.Codecpar()); err != nil {
		_ = d.Close()
		return nil, fmt.Errorf("copy codec parameters to context: %w", err)
	}

	ctx.SetPktTimebase(stream.TimeBase())

	hwPixFmt := ffmpeg.AVPixFmtNone
	for i := 0; ; i++ {
		config := ffmpeg.AVCodecGetHWConfig(codec, i)
		if config == nil {
			break
		}

		if config.Methods()&ffmpeg.AVCodecHWConfigMethodHWDeviceCtx != 0 && config.DeviceType() == dev.Type() {
			hwPixFmt = config.PixFmt()
			break
		}
	}

	if hwPixFmt == ffmpeg.AVPixFmtNone {
		_ = d.Close()
		return nil, fmt.Errorf("decoder %s has no hw config for device type %v", ffmpeg.AVCodecGetName(stream.Codecpar().CodecId()).String(), dev.Type())
	}

	// Take an independent reference so the decoder owns its own device-context
	// reference; the HWDevice keeps its reference and both unref independently.
	ctx.SetHwDeviceCtx(ffmpeg.AVBufferRef_(dev.Raw()))

	ctx.SetGetFormat(func(_ *ffmpeg.AVCodecContext, formats []ffmpeg.AVPixelFormat) ffmpeg.AVPixelFormat {
		if slices.Contains(formats, hwPixFmt) {
			return hwPixFmt
		}

		return ffmpeg.AVPixFmtNone
	})

	if _, err := ffmpeg.AVCodecOpen2(ctx, codec, nil); err != nil {
		_ = d.Close()
		return nil, fmt.Errorf("open hw decoder: %w", err)
	}

	d.frame = ffmpeg.AVFrameAlloc()
	if d.frame == nil {
		_ = d.Close()
		return nil, errors.New("alloc decoder frame failed")
	}

	return d, nil
}

// TransferToSoftware copies hwFrame's data from device memory into a freshly
// allocated software frame and returns it. The caller owns the returned frame
// and must free it with [ffmpeg.AVFrameFree]. The presentation timestamp is
// carried over; on transfer error the scratch frame is freed and the error is
// returned.
func TransferToSoftware(hwFrame *ffmpeg.AVFrame) (*ffmpeg.AVFrame, error) {
	if hwFrame == nil {
		return nil, errors.New("nil hw frame")
	}

	sw := ffmpeg.AVFrameAlloc()
	if sw == nil {
		return nil, errors.New("alloc software frame failed")
	}

	if _, err := ffmpeg.AVHWFrameTransferData(sw, hwFrame, 0); err != nil {
		ffmpeg.AVFrameFree(&sw)
		return nil, fmt.Errorf("transfer hw frame to software: %w", err)
	}

	sw.SetPts(hwFrame.Pts())

	return sw, nil
}
