package av

import (
	"errors"
	"fmt"
	"io"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
)

// Encoder owns an encode [ffmpeg.AVCodecContext] and a scratch [ffmpeg.AVPacket].
// It implements [io.Closer]; Close is idempotent and frees both handles in
// dependency order. Reach any unsurfaced capability through Raw. Encoder is the
// mirror image of [Decoder]: it sends frames and drains packets.
type Encoder struct {
	ctx *ffmpeg.AVCodecContext
	pkt *ffmpeg.AVPacket
}

var _ io.Closer = (*Encoder)(nil)

// NewEncoder builds an encoder for codec. It allocates a context, hands it to
// configure so the caller sets every codec field (width/height/pix_fmt/
// time_base, sample_rate/ch_layout/sample_fmt, and any flags), opens the codec
// and allocates a scratch packet. On any error after the context is allocated
// it unwinds via Close so no handle leaks.
func NewEncoder(codec *ffmpeg.AVCodec, configure func(*ffmpeg.AVCodecContext)) (*Encoder, error) {
	if codec == nil {
		return nil, errors.New("nil encoder codec")
	}

	ctx := ffmpeg.AVCodecAllocContext3(codec)
	if ctx == nil {
		return nil, errors.New("alloc encoder context failed")
	}

	e := &Encoder{ctx: ctx}

	if configure != nil {
		configure(ctx)
	}

	if _, err := ffmpeg.AVCodecOpen2(ctx, codec, nil); err != nil {
		_ = e.Close()
		return nil, fmt.Errorf("open encoder: %w", err)
	}

	e.pkt = ffmpeg.AVPacketAlloc()
	if e.pkt == nil {
		_ = e.Close()
		return nil, errors.New("alloc encoder packet failed")
	}

	return e, nil
}

// NewEncoderByID is a convenience wrapper that looks up the encoder for id and
// calls NewEncoder. NewEncoder stays the primary entry point for callers that
// already hold an [ffmpeg.AVCodec].
func NewEncoderByID(id ffmpeg.AVCodecID, configure func(*ffmpeg.AVCodecContext)) (*Encoder, error) {
	codec := ffmpeg.AVCodecFindEncoder(id)
	if codec == nil {
		return nil, fmt.Errorf("find encoder (codec %v): %w", id, ffmpeg.WrapErr(ffmpeg.AVErrorEncoderNotFoundConst))
	}

	return NewEncoder(codec, configure)
}

// Encode sends frame to the encoder and calls fn for each packet it yields. The
// packet is valid only for the duration of fn; Encode unrefs the scratch packet
// after fn returns. Internal EAgain/EOF signals control the drain loop and never
// reach the caller; only fn's error or a real encode error propagates.
func (e *Encoder) Encode(frame *ffmpeg.AVFrame, fn func(*ffmpeg.AVPacket) error) error {
	if _, err := ffmpeg.AVCodecSendFrame(e.ctx, frame); err != nil {
		return fmt.Errorf("send frame: %w", err)
	}

	return e.drain(fn)
}

// Flush sends a nil frame to enter draining mode then drains every remaining
// packet through fn, mirroring Encode's callback and unref contract. Sending a
// nil frame and draining is safe whether or not the encoder advertises
// [ffmpeg.AVCodecCapDelay].
func (e *Encoder) Flush(fn func(*ffmpeg.AVPacket) error) error {
	if _, err := ffmpeg.AVCodecSendFrame(e.ctx, nil); err != nil {
		return fmt.Errorf("flush encoder: %w", err)
	}

	return e.drain(fn)
}

// drain receives packets into the scratch packet until EAgain/EOF, invoking fn
// for each and unreffing the scratch packet after every callback.
func (e *Encoder) drain(fn func(*ffmpeg.AVPacket) error) error {
	for {
		if _, err := ffmpeg.AVCodecReceivePacket(e.ctx, e.pkt); err != nil {
			if errors.Is(err, ffmpeg.EAgain) || errors.Is(err, ffmpeg.AVErrorEOF) {
				return nil
			}

			return fmt.Errorf("receive packet: %w", err)
		}

		err := fn(e.pkt)
		ffmpeg.AVPacketUnref(e.pkt)

		if err != nil {
			return err
		}
	}
}

// Raw returns the underlying encode context. Do not free it; Close owns it.
func (e *Encoder) Raw() *ffmpeg.AVCodecContext {
	return e.ctx
}

// Close frees the scratch packet and encode context and nils both handles. A
// second call is a no-op.
func (e *Encoder) Close() error {
	if e.pkt != nil {
		ffmpeg.AVPacketFree(&e.pkt)
		e.pkt = nil
	}

	if e.ctx != nil {
		ffmpeg.AVCodecFreeContext(&e.ctx)
		e.ctx = nil
	}

	return nil
}
