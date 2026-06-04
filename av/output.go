package av

import (
	"errors"
	"fmt"
	"io"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
)

// Output owns a mux [ffmpeg.AVFormatContext] and its IO. It implements
// [io.Closer]; Close is idempotent, mirrors any IO it opened, then frees the
// context. Reach any unsurfaced capability through Raw. Output is the muxing
// counterpart to [Input].
type Output struct {
	ctx    *ffmpeg.AVFormatContext
	name   string
	ioOpen bool
}

var _ io.Closer = (*Output)(nil)

// CreateOutput allocates a mux context for name, guessing the muxer from the
// filename. It does not open IO; WriteHeader opens it once the muxer flags are
// known. On any error after the context is allocated it unwinds via Close so no
// handle leaks.
func CreateOutput(name string) (*Output, error) {
	namePtr := ffmpeg.ToCStr(name)
	defer namePtr.Free()

	var ctx *ffmpeg.AVFormatContext
	if _, err := ffmpeg.AVFormatAllocOutputContext2(&ctx, nil, nil, namePtr); err != nil {
		return nil, fmt.Errorf("alloc output context %q: %w", name, err)
	}

	if ctx == nil {
		return nil, fmt.Errorf("alloc output context %q: nil context", name)
	}

	return &Output{ctx: ctx, name: name}, nil
}

// AddStream adds a stream to the muxer and copies enc's codec parameters into
// it, setting the stream time_base from the encoder. When the muxer advertises
// [ffmpeg.AVFmtGlobalheader], the encoder must have been opened with
// [ffmpeg.AVCodecFlagGlobalHeader]; that flag cannot be set after open, so the
// caller must set it in their NewEncoder configure callback. AddStream does not
// mutate the encoder.
func (o *Output) AddStream(enc *Encoder) (*ffmpeg.AVStream, error) {
	if enc == nil {
		return nil, errors.New("nil encoder")
	}

	stream := ffmpeg.AVFormatNewStream(o.ctx, nil)
	if stream == nil {
		return nil, errors.New("alloc output stream failed")
	}

	if _, err := ffmpeg.AVCodecParametersFromContext(stream.Codecpar(), enc.Raw()); err != nil {
		return nil, fmt.Errorf("copy codec parameters: %w", err)
	}

	stream.SetTimeBase(enc.Raw().TimeBase())

	return stream, nil
}

// WriteHeader opens the muxer IO when the format needs a file
// (not [ffmpeg.AVFmtNofile]) and writes the stream header. It records whether
// it opened IO so Close mirrors it.
func (o *Output) WriteHeader() error {
	if o.ctx.Oformat().Flags()&ffmpeg.AVFmtNofile == 0 {
		namePtr := ffmpeg.ToCStr(o.name)
		defer namePtr.Free()

		var pb *ffmpeg.AVIOContext
		if _, err := ffmpeg.AVIOOpen(&pb, namePtr, ffmpeg.AVIOFlagWrite); err != nil {
			return fmt.Errorf("open output IO %q: %w", o.name, err)
		}

		o.ctx.SetPb(pb)
		o.ioOpen = true
	}

	if _, err := ffmpeg.AVFormatWriteHeader(o.ctx, nil); err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	return nil
}

// WritePacket interleaves pkt into the muxer. The caller owns pkt and must set
// its stream index and rescale its timestamps to the stream time_base first.
func (o *Output) WritePacket(pkt *ffmpeg.AVPacket) error {
	if _, err := ffmpeg.AVInterleavedWriteFrame(o.ctx, pkt); err != nil {
		return fmt.Errorf("write packet: %w", err)
	}

	return nil
}

// WriteTrailer writes the stream trailer, finishing the file.
func (o *Output) WriteTrailer() error {
	if _, err := ffmpeg.AVWriteTrailer(o.ctx); err != nil {
		return fmt.Errorf("write trailer: %w", err)
	}

	return nil
}

// Raw returns the underlying mux context. Do not free it; Close owns it.
func (o *Output) Raw() *ffmpeg.AVFormatContext {
	return o.ctx
}

// Close mirrors any IO it opened, then frees the mux context and nils the
// handle. It closes IO only when WriteHeader opened it (the muxer is not
// [ffmpeg.AVFmtNofile]), closing the pb and clearing it before freeing the
// context. A second call is a no-op.
func (o *Output) Close() error {
	if o.ctx == nil {
		return nil
	}

	var ioErr error

	if o.ioOpen {
		if _, err := ffmpeg.AVIOClose(o.ctx.Pb()); err != nil {
			ioErr = fmt.Errorf("close output IO: %w", err)
		}

		o.ctx.SetPb(nil)
		o.ioOpen = false
	}

	ffmpeg.AVFormatFreeContext(o.ctx)
	o.ctx = nil

	return ioErr
}
