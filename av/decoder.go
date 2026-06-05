package av

import (
	"errors"
	"fmt"
	"io"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
)

// Decoder owns a decode [ffmpeg.AVCodecContext] and a scratch [ffmpeg.AVFrame].
// It implements [io.Closer]; Close is idempotent and frees both handles in
// dependency order. Reach any unsurfaced capability through Raw.
type Decoder struct {
	ctx   *ffmpeg.AVCodecContext
	frame *ffmpeg.AVFrame
}

var _ io.Closer = (*Decoder)(nil)

// NewDecoder builds a decoder for stream. It finds the decoder for the stream's
// codec, allocates and configures a context from the stream's parameters, sets
// the packet timebase, opens the codec and allocates a scratch frame. On any
// error after the context is allocated it unwinds via Close so no handle leaks.
func NewDecoder(stream *ffmpeg.AVStream) (*Decoder, error) {
	if stream == nil {
		return nil, errors.New("nil stream")
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

	if _, err := ffmpeg.AVCodecOpen2(ctx, codec, nil); err != nil {
		_ = d.Close()
		return nil, fmt.Errorf("open decoder: %w", err)
	}

	d.frame = ffmpeg.AVFrameAlloc()
	if d.frame == nil {
		_ = d.Close()
		return nil, errors.New("alloc decoder frame failed")
	}

	return d, nil
}

// Decode sends pkt to the decoder and calls fn for each frame it yields. The
// frame is valid only for the duration of fn; Decode unrefs the scratch frame
// after fn returns. Internal EAgain/EOF signals control the drain loop and
// never reach the caller; only fn's error or a real decode error propagates.
func (d *Decoder) Decode(pkt *ffmpeg.AVPacket, fn func(*ffmpeg.AVFrame) error) error {
	if _, err := ffmpeg.AVCodecSendPacket(d.ctx, pkt); err != nil {
		return fmt.Errorf("send packet: %w", err)
	}

	return d.drain(fn)
}

// Flush sends a nil packet to enter draining mode then drains every remaining
// frame through fn, mirroring Decode's callback and unref contract.
func (d *Decoder) Flush(fn func(*ffmpeg.AVFrame) error) error {
	if _, err := ffmpeg.AVCodecSendPacket(d.ctx, nil); err != nil {
		return fmt.Errorf("flush decoder: %w", err)
	}

	return d.drain(fn)
}

// drain receives frames into the scratch frame until EAgain/EOF, invoking fn
// for each and unreffing the scratch frame after every callback.
func (d *Decoder) drain(fn func(*ffmpeg.AVFrame) error) error {
	for {
		if _, err := ffmpeg.AVCodecReceiveFrame(d.ctx, d.frame); err != nil {
			if errors.Is(err, ffmpeg.EAgain) || errors.Is(err, ffmpeg.AVErrorEOF) {
				return nil
			}

			return fmt.Errorf("receive frame: %w", err)
		}

		err := fn(d.frame)
		ffmpeg.AVFrameUnref(d.frame)

		if err != nil {
			return err
		}
	}
}

// Raw returns the underlying decode context. Do not free it; Close owns it.
func (d *Decoder) Raw() *ffmpeg.AVCodecContext {
	return d.ctx
}

// Close frees the scratch frame and decode context and nils both handles. A
// second call is a no-op.
func (d *Decoder) Close() error {
	if d.frame != nil {
		ffmpeg.AVFrameFree(&d.frame)
		d.frame = nil
	}

	if d.ctx != nil {
		// Unregister any get_format callback before freeing: the registry is
		// keyed by context pointer, so a stale entry would mis-fire on a future
		// context reusing this address. Idempotent no-op for software decoders.
		d.ctx.ClearGetFormat()
		ffmpeg.AVCodecFreeContext(&d.ctx)
		d.ctx = nil
	}

	return nil
}
