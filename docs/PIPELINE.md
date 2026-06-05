# Pipeline Guide

The `av` package (`github.com/linuxmatters/ffmpeg-statigo/av`) is an optional high-level layer over the root `ffmpeg` package. It adds Go-idiomatic ownership and lifetime management to the common demux, decode, filter, encode, mux pipeline.

Use `av` when you want resource safety and a smaller API surface. The root `ffmpeg` package remains the raw bridge to FFmpeg's C API and is fully usable on its own; drop to it whenever you need the full API. The two layers interoperate through `Raw`, so choosing `av` never traps you.

## Owned types

The package owns five types. Each wraps an FFmpeg handle, implements [`io.Closer`](https://pkg.go.dev/io#Closer), and exposes `Raw` for the underlying `*ffmpeg.AV*`.

### Input

Owns a demux `*ffmpeg.AVFormatContext` opened from a URL.

- `Open(url string) (*Input, error)` - opens the URL and reads stream information.
- `Streams() *ffmpeg.Array[*ffmpeg.AVStream]` - the demuxer's stream array.
- `NbStreams() uint` - the number of streams.
- `BestStream(mediaType ffmpeg.AVMediaType) (*ffmpeg.AVStream, error)` - wraps `av_find_best_stream`.
- `ReadPacket(pkt *ffmpeg.AVPacket) error` - reads the next packet; at end of stream the error satisfies `errors.Is(err, ffmpeg.AVErrorEOF)`.
- `Raw() *ffmpeg.AVFormatContext`
- `Close() error`

### Decoder

Owns a decode `*ffmpeg.AVCodecContext` and a scratch `*ffmpeg.AVFrame`.

- `NewDecoder(stream *ffmpeg.AVStream) (*Decoder, error)` - finds the decoder for the stream's codec, configures a context from the stream parameters, opens the codec, and allocates a scratch frame.
- `Decode(pkt *ffmpeg.AVPacket, fn func(*ffmpeg.AVFrame) error) error` - sends a packet and calls `fn` once per decoded frame.
- `Flush(fn func(*ffmpeg.AVFrame) error) error` - sends a nil packet and drains the remaining frames through `fn`.
- `Raw() *ffmpeg.AVCodecContext`
- `Close() error`

Internal `EAgain`/`AVErrorEOF` drain signals never reach the caller; only `fn`'s error or a real decode error propagates.

### FilterGraph

Owns an `*ffmpeg.AVFilterGraph` plus its single buffersrc and buffersink. It handles the single-input/single-output video and audio cases.

- `NewVideoFilterGraph(params VideoFilterParams, spec string) (*FilterGraph, error)`
- `NewAudioFilterGraph(params AudioFilterParams, spec string) (*FilterGraph, error)`
- `NewVideoFilterGraphFromContext(decCtx *ffmpeg.AVCodecContext, spec string, outPixFmt ffmpeg.AVPixelFormat) (*FilterGraph, error)`
- `NewAudioFilterGraphFromContext(decCtx *ffmpeg.AVCodecContext, spec string, outSampleFmt ffmpeg.AVSampleFormat, outSampleRate int, outChLayout *ffmpeg.AVChannelLayout) (*FilterGraph, error)`
- `Push(frame *ffmpeg.AVFrame) error` - sends a frame into the buffersrc. Push `nil` to flush.
- `Pull(fn func(*ffmpeg.AVFrame) error) error` - drains every filtered frame currently available, calling `fn` once per frame.
- `Raw() *ffmpeg.AVFilterGraph`
- `Close() error`

The `...FromContext` constructors read source parameters from a decoder context. Use `"null"` for a passthrough video graph and `"anull"` for audio.

### Encoder

Owns an encode `*ffmpeg.AVCodecContext` and a scratch `*ffmpeg.AVPacket`. It mirrors `Decoder`: it sends frames and drains packets.

- `NewEncoder(codec *ffmpeg.AVCodec, configure func(*ffmpeg.AVCodecContext)) (*Encoder, error)` - allocates a context, hands it to `configure` so you set every codec field before open, then opens the codec and allocates a scratch packet.
- `NewEncoderByID(id ffmpeg.AVCodecID, configure func(*ffmpeg.AVCodecContext)) (*Encoder, error)` - looks up the encoder for `id` then calls `NewEncoder`.
- `Encode(frame *ffmpeg.AVFrame, fn func(*ffmpeg.AVPacket) error) error` - sends a frame and calls `fn` once per encoded packet.
- `Flush(fn func(*ffmpeg.AVPacket) error) error` - sends a nil frame and drains the remaining packets through `fn`.
- `Raw() *ffmpeg.AVCodecContext`
- `Close() error`

Set `ffmpeg.AVCodecFlagGlobalHeader` in the `configure` callback when the muxer advertises `ffmpeg.AVFmtGlobalheader`; the flag cannot be set after open.

### Output

Owns a mux `*ffmpeg.AVFormatContext` and its IO. It is the muxing counterpart to `Input`.

- `CreateOutput(name string) (*Output, error)` - allocates a mux context, guessing the muxer from the filename. It does not open IO.
- `AddStream(enc *Encoder) (*ffmpeg.AVStream, error)` - adds a stream and copies the encoder's codec parameters and time base into it. It does not mutate the encoder.
- `WriteHeader() error` - opens IO when the format needs a file, then writes the header.
- `WritePacket(pkt *ffmpeg.AVPacket) error` - interleaves a packet. The caller owns `pkt` and must set its stream index and rescale its timestamps to the stream time base first.
- `WriteTrailer() error` - writes the trailer, finishing the file.
- `Raw() *ffmpeg.AVFormatContext`
- `Close() error`

## Ownership and Close contract

Every owned type implements `io.Closer`. `Close` frees the underlying FFmpeg handles in dependency order, so closing a type releases the resources it owns. `Close` is idempotent: a second call is a safe no-op and never double-frees.

Pair each constructor with a deferred `Close`:

```go
in, err := av.Open("input.mp4")
if err != nil {
    return err
}
defer in.Close()
```

## Callback frame lifetime

> [!WARNING]
> Frames and packets passed to `Decode`, `Encode`, and `Pull` callbacks are valid only for the duration of that callback. Each type unrefs its scratch buffer after the callback returns.

Copy or ref the data (for example with `ffmpeg.AVFrameRef`) if you need to keep it past the callback.

## The Raw escape hatch

Every owned type exposes `Raw`, returning the underlying `*ffmpeg.AV*` handle. Any FFmpeg capability the `av` layer does not surface stays reachable through `Raw`:

```go
globalHeader := output.Raw().Oformat().Flags()&ffmpeg.AVFmtGlobalheader != 0
```

Do not free a handle obtained from `Raw` yourself. Ownership stays with the `av` type and its `Close`.

## End-to-end

This sketch transcodes one video stream, adapted from [`examples/transcode-hl/main.go`](../examples/transcode-hl/main.go).

```go
input, err := av.Open(inPath)
if err != nil {
    return err
}
defer input.Close()

output, err := av.CreateOutput(outPath)
if err != nil {
    return err
}
defer output.Close()

globalHeader := output.Raw().Oformat().Flags()&ffmpeg.AVFmtGlobalheader != 0

inStream, err := input.BestStream(ffmpeg.AVMediaTypeVideo)
if err != nil {
    return err
}

decoder, err := av.NewDecoder(inStream)
if err != nil {
    return err
}
defer decoder.Close()
decCtx := decoder.Raw()

encoder, err := av.NewEncoderByID(decCtx.CodecId(), func(encCtx *ffmpeg.AVCodecContext) {
    encCtx.SetHeight(decCtx.Height())
    encCtx.SetWidth(decCtx.Width())
    encCtx.SetSampleAspectRatio(decCtx.SampleAspectRatio())
    encCtx.SetPixFmt(decCtx.PixFmt())
    encCtx.SetTimeBase(ffmpeg.AVInvQ(decCtx.Framerate()))
    if globalHeader {
        encCtx.SetFlags(encCtx.Flags() | ffmpeg.AVCodecFlagGlobalHeader)
    }
})
if err != nil {
    return err
}
defer encoder.Close()

outStream, err := output.AddStream(encoder)
if err != nil {
    return err
}

filter, err := av.NewVideoFilterGraphFromContext(decCtx, "null", encoder.Raw().PixFmt())
if err != nil {
    return err
}
defer filter.Close()

if err := output.WriteHeader(); err != nil {
    return err
}

packet := ffmpeg.AVPacketAlloc()
defer ffmpeg.AVPacketFree(&packet)

encCtx := encoder.Raw()
mux := func(pkt *ffmpeg.AVPacket) error {
    pkt.SetStreamIndex(outStream.Index())
    ffmpeg.AVPacketRescaleTs(pkt, encCtx.TimeBase(), outStream.TimeBase())
    return output.WritePacket(pkt)
}

for {
    if err := input.ReadPacket(packet); err != nil {
        if errors.Is(err, ffmpeg.AVErrorEOF) {
            break
        }
        return err
    }

    err := decoder.Decode(packet, func(frame *ffmpeg.AVFrame) error {
        frame.SetPts(frame.BestEffortTimestamp())
        if err := filter.Push(frame); err != nil {
            return err
        }
        return filter.Pull(func(filtered *ffmpeg.AVFrame) error {
            return encoder.Encode(filtered, mux)
        })
    })
    ffmpeg.AVPacketUnref(packet)
    if err != nil {
        return err
    }
}

// Flush the decoder, filter graph and encoder before the trailer so delayed
// tail frames are not dropped. See examples/transcode-hl/main.go for the full
// flush sequence (decoder.Flush -> filter Push(nil)/Pull -> encoder.Flush).

return output.WriteTrailer()
```

See [`examples/transcode-hl/main.go`](../examples/transcode-hl/main.go) for the full multi-stream version with decoder, filter, and encoder flushing at end of stream.
