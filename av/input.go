package av

import (
	"fmt"
	"io"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
)

// Input owns a demux [ffmpeg.AVFormatContext] opened from a URL. It implements
// [io.Closer]; Close is idempotent and frees the context. Reach any unsurfaced
// capability through Raw.
type Input struct {
	ctx *ffmpeg.AVFormatContext
}

var _ io.Closer = (*Input)(nil)

// Open opens url for demuxing and reads stream information. On any error after
// the context is allocated it unwinds via Close so no handle leaks.
func Open(url string) (*Input, error) {
	urlPtr := ffmpeg.ToCStr(url)
	defer urlPtr.Free()

	var ctx *ffmpeg.AVFormatContext
	if _, err := ffmpeg.AVFormatOpenInput(&ctx, urlPtr, nil, nil); err != nil {
		return nil, fmt.Errorf("open input %q: %w", url, err)
	}

	in := &Input{ctx: ctx}

	if _, err := ffmpeg.AVFormatFindStreamInfo(ctx, nil); err != nil {
		_ = in.Close()
		return nil, fmt.Errorf("find stream info %q: %w", url, err)
	}

	return in, nil
}

// Streams returns the demuxer's stream array. Ownership stays with the Input.
func (in *Input) Streams() *ffmpeg.Array[*ffmpeg.AVStream] {
	return in.ctx.Streams()
}

// NbStreams returns the number of streams in the input.
func (in *Input) NbStreams() uint {
	return in.ctx.NbStreams()
}

// BestStream returns the best stream of mediaType, wrapping av_find_best_stream.
// It returns a wrapped [ffmpeg.AVError] when no stream is found.
func (in *Input) BestStream(mediaType ffmpeg.AVMediaType) (*ffmpeg.AVStream, error) {
	idx, err := ffmpeg.AVFindBestStream(in.ctx, mediaType, -1, -1, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("find best stream (type %v): %w", mediaType, err)
	}

	return in.ctx.Streams().Get(uintptr(idx)), nil
}

// ReadPacket reads the next packet into pkt. At end of stream it returns an
// error satisfying errors.Is(err, [ffmpeg.AVErrorEOF]).
func (in *Input) ReadPacket(pkt *ffmpeg.AVPacket) error {
	_, err := ffmpeg.AVReadFrame(in.ctx, pkt)
	return err
}

// Raw returns the underlying demux context. Do not free it; Close owns it.
func (in *Input) Raw() *ffmpeg.AVFormatContext {
	return in.ctx
}

// Close frees the demux context and nils the handle. A second call is a no-op.
func (in *Input) Close() error {
	if in.ctx == nil {
		return nil
	}

	ffmpeg.AVFormatCloseInput(&in.ctx)
	in.ctx = nil

	return nil
}
