# High-level `av` package

The `av` package (`github.com/linuxmatters/ffmpeg-statigo/av`) is an optional high-level layer over the root `ffmpeg` package. It adds Go-idiomatic ownership and lifetime management to the common demux, decode, filter, encode, and mux pipeline, including hardware-accelerated decoding.

Use `av` when you want resource safety and a smaller API surface. The root `ffmpeg` package remains the raw bridge to FFmpeg's C API and is usable on its own. Every `av` type exposes `Raw()`, so you can mix both layers in the same pipeline without losing access to any C API feature.

## Owned types

The package owns six types. Each type wraps an FFmpeg handle, implements [`io.Closer`](https://pkg.go.dev/io#Closer), and exposes `Raw()` for the underlying `*ffmpeg.AV*` pointer.

### `Input`

Owns a demux `*ffmpeg.AVFormatContext` opened from a file path or URL.

- `Open(url string) (*Input, error)` - opens the URL and reads stream information.
- `Streams() *ffmpeg.Array[*ffmpeg.AVStream]` - returns the demuxer's stream array.
- `NbStreams() uint` - returns the number of streams.
- `BestStream(mediaType ffmpeg.AVMediaType) (*ffmpeg.AVStream, error)` - wraps `av_find_best_stream`.
- `ReadPacket(pkt *ffmpeg.AVPacket) error` - reads the next packet. At end of stream, the error satisfies `errors.Is(err, ffmpeg.AVErrorEOF)`.
- `Raw() *ffmpeg.AVFormatContext`
- `Close() error`

### `Decoder`

Owns a decode `*ffmpeg.AVCodecContext` and a scratch `*ffmpeg.AVFrame`.

- `NewDecoder(stream *ffmpeg.AVStream) (*Decoder, error)` - finds the decoder for the stream's codec, configures a context from the stream parameters, sets the packet timebase, opens the codec, and allocates a scratch frame.
- `NewHWDecoder(stream *ffmpeg.AVStream, dev *HWDevice) (*Decoder, error)` - builds a hardware-accelerated decoder for `stream` on `dev`. See [Hardware-accelerated decoding](#hardware-accelerated-decoding).
- `Decode(pkt *ffmpeg.AVPacket, fn func(*ffmpeg.AVFrame) error) error` - sends a packet and calls `fn` once per decoded frame.
- `Flush(fn func(*ffmpeg.AVFrame) error) error` - sends a nil packet and drains the remaining frames through `fn`.
- `Raw() *ffmpeg.AVCodecContext`
- `Close() error`

Internal `EAgain` and `AVErrorEOF` drain signals never reach the caller. Only an error from `fn` or a real decode error propagates.

### `FilterGraph`

Owns an `*ffmpeg.AVFilterGraph` plus its single buffersrc and buffersink. It handles single-input and single-output video and audio graphs.

- `NewVideoFilterGraph(params VideoFilterParams, spec string) (*FilterGraph, error)`
- `NewAudioFilterGraph(params AudioFilterParams, spec string) (*FilterGraph, error)`
- `NewVideoFilterGraphFromContext(decCtx *ffmpeg.AVCodecContext, spec string, outPixFmt ffmpeg.AVPixelFormat) (*FilterGraph, error)`
- `NewAudioFilterGraphFromContext(decCtx *ffmpeg.AVCodecContext, spec string, outSampleFmt ffmpeg.AVSampleFormat, outSampleRate int, outChLayout *ffmpeg.AVChannelLayout) (*FilterGraph, error)`
- `Push(frame *ffmpeg.AVFrame) error` - sends a frame into the buffersrc. Push `nil` to flush.
- `Pull(fn func(*ffmpeg.AVFrame) error) error` - drains every filtered frame currently available, calling `fn` once per frame.
- `Raw() *ffmpeg.AVFilterGraph`
- `Close() error`

The `...FromContext` constructors read source parameters from a decoder context. Use `"null"` for a passthrough video graph and `"anull"` for audio.

### `Encoder`

Owns an encode `*ffmpeg.AVCodecContext` and a scratch `*ffmpeg.AVPacket`. It mirrors `Decoder`: it sends frames and drains packets.

- `NewEncoder(codec *ffmpeg.AVCodec, configure func(*ffmpeg.AVCodecContext)) (*Encoder, error)` - allocates a context, passes it to `configure` so you set codec fields before open, then opens the codec and allocates a scratch packet.
- `NewEncoderByID(id ffmpeg.AVCodecID, configure func(*ffmpeg.AVCodecContext)) (*Encoder, error)` - looks up the encoder for `id` and calls `NewEncoder`.
- `Encode(frame *ffmpeg.AVFrame, fn func(*ffmpeg.AVPacket) error) error` - sends a frame and calls `fn` once per encoded packet.
- `Flush(fn func(*ffmpeg.AVPacket) error) error` - sends a nil frame and drains the remaining packets through `fn`.
- `Raw() *ffmpeg.AVCodecContext`
- `Close() error`

Set `ffmpeg.AVCodecFlagGlobalHeader` in the `configure` callback when the muxer advertises `ffmpeg.AVFmtGlobalheader`. You cannot set that flag after the codec opens.

### `Output`

Owns a mux `*ffmpeg.AVFormatContext` and its IO. It is the muxing counterpart to `Input`.

- `CreateOutput(name string) (*Output, error)` - allocates a mux context and guesses the muxer from the filename. It does not open IO.
- `AddStream(enc *Encoder) (*ffmpeg.AVStream, error)` - adds a stream and copies the encoder's codec parameters and time base into it. It does not mutate the encoder.
- `WriteHeader() error` - opens IO when the format needs a file, then writes the header.
- `WritePacket(pkt *ffmpeg.AVPacket) error` - interleaves a packet. The caller owns `pkt` and must set its stream index and rescale its timestamps to the stream time base first.
- `WriteTrailer() error` - writes the trailer and finishes the file.
- `Raw() *ffmpeg.AVFormatContext`
- `Close() error`

### `HWDevice`

Owns a hardware device context (`*ffmpeg.AVBufferRef`) together with the `ffmpeg.AVHWDeviceType` that created it.

- `HWDeviceTypes() []string` - returns the names of every hardware device type compiled into the linked FFmpeg library.
- `NewHWDevice(typeName string) (*HWDevice, error)` - resolves `typeName` to an `ffmpeg.AVHWDeviceType` and creates a device context with the default device and options.
- `Type() ffmpeg.AVHWDeviceType` - returns the device type for this `HWDevice`.
- `Raw() *ffmpeg.AVBufferRef` - returns the underlying device context. Pass it through `ffmpeg.AVBufferRef_` if another consumer needs an independent reference.
- `Close() error` - unrefs the device context and sets the handle to nil.

## Hardware-accelerated decoding

Hardware decoding uses `HWDeviceTypes`, `NewHWDevice`, `NewHWDecoder`, and `TransferToSoftware` (defined in [`av/hwdevice.go`](https://github.com/linuxmatters/ffmpeg-statigo/blob/main/av/hwdevice.go)).

### Device discovery and creation

`HWDeviceTypes()` enumerates the FFmpeg hardware device registry and requires no GPU, so it is safe to call on hardware-free hosts. Under `-tags=embedded` or on builds without hardware backends, it returns an empty slice.

`NewHWDevice(typeName)` rejects unknown type names with an error. If a compiled-in backend has no matching driver or GPU on the host, `NewHWDevice` returns a wrapped error rather than panicking, so callers can try available types in order and fall back to software decoding.

### Creating a hardware decoder

`NewHWDecoder(stream, dev)` mirrors `NewDecoder` and wires hardware decoding before opening the codec:

1. Finds the decoder's `AVCodecHWConfig` entry that supports `AVCodecHWConfigMethodHWDeviceCtx` for `dev.Type()` and records its hardware pixel format.
2. Attaches an independent reference (`ffmpeg.AVBufferRef_(dev.Raw())`) to the decoder context, so closing `HWDevice` and `Decoder` in any order is safe.
3. Registers a `get_format` callback that selects the discovered hardware pixel format. `Decoder.Close()` clears this callback before freeing the context.

### Frame transfer to software memory

Frames yielded by a hardware `Decoder` live in device memory in the hardware pixel format. Call `TransferToSoftware(hwFrame)` to copy a hardware frame into a newly allocated software `*ffmpeg.AVFrame` in system memory.

- `TransferToSoftware` copies `hwFrame.Pts()` to the returned frame.
- On transfer failure, `TransferToSoftware` frees its scratch allocation and returns the error.
- The caller owns the returned software frame and must free it with `ffmpeg.AVFrameFree(&sw)`.

### Hardware decode example

This snippet probes compiled hardware device types, opens a hardware decoder with software fallback, and transfers hardware frames back to CPU memory. It is adapted from [`examples/hwdecode/main.go`](https://github.com/linuxmatters/ffmpeg-statigo/blob/main/examples/hwdecode/main.go).

```go
in, err := av.Open(inPath)
if err != nil {
    return err
}
defer in.Close()

stream, err := in.BestStream(ffmpeg.AVMediaTypeVideo)
if err != nil {
    return err
}

var dev *av.HWDevice
for _, name := range av.HWDeviceTypes() {
    if d, derr := av.NewHWDevice(name); derr == nil {
        dev = d
        break
    }
}

var dec *av.Decoder
if dev != nil {
    defer dev.Close()
    dec, err = av.NewHWDecoder(stream, dev)
}
if dec == nil {
    dec, err = av.NewDecoder(stream)
    if err != nil {
        return err
    }
}
defer dec.Close()

pkt := ffmpeg.AVPacketAlloc()
defer ffmpeg.AVPacketFree(&pkt)

emit := func(frame *ffmpeg.AVFrame) error {
    if dev != nil && frame.HwFramesCtx() != nil {
        sw, terr := av.TransferToSoftware(frame)
        if terr != nil {
            return terr
        }
        defer ffmpeg.AVFrameFree(&sw)
        // Process software frame sw here.
        return nil
    }
    // Process software frame directly.
    return nil
}

for {
    if err := in.ReadPacket(pkt); err != nil {
        if errors.Is(err, ffmpeg.AVErrorEOF) {
            break
        }
        return err
    }
    if pkt.StreamIndex() == stream.Index() {
        err := dec.Decode(pkt, emit)
        ffmpeg.AVPacketUnref(pkt)
        if err != nil {
            return err
        }
        continue
    }
    ffmpeg.AVPacketUnref(pkt)
}

return dec.Flush(emit)
```

See [`examples/hwdecode/main.go`](https://github.com/linuxmatters/ffmpeg-statigo/blob/main/examples/hwdecode/main.go) for the full program, including explicit device-type selection and matching the exact hardware pixel format.

## Ownership and `Close` contract

Every owned type implements `io.Closer`. `Close` frees the underlying FFmpeg handles in dependency order, so closing a type releases every resource it owns. `Close` is idempotent: a second call is a safe no-op and never double-frees.

Pair each constructor with a deferred `Close`:

```go
in, err := av.Open("input.mp4")
if err != nil {
    return err
}
defer in.Close()
```

If a constructor fails partway through setup, it unwinds through `Close` before returning the error so no C handle leaks.

## Callback frame lifetime

> [!WARNING]
> Frames and packets passed to `Decode`, `Encode`, and `Pull` callbacks are valid only for the duration of that callback. Each type unrefs its scratch buffer after the callback returns.

If you need to keep frame data past the callback, copy or reference it with `ffmpeg.AVFrameRef` (or `av.TransferToSoftware` for hardware frames) and free that copy yourself.

## The `Raw` escape hatch

Every owned type exposes `Raw()`, returning its underlying `*ffmpeg.AV*` handle. Any FFmpeg capability that `av` does not surface stays reachable through `Raw()`:

```go
globalHeader := output.Raw().Oformat().Flags()&ffmpeg.AVFmtGlobalheader != 0
```

Do not free or unref a handle obtained from `Raw()`. Ownership stays with the `av` type and its `Close()` method.

## End-to-end transcode example

This sketch transcodes one video stream, adapted from [`examples/transcode-hl/main.go`](https://github.com/linuxmatters/ffmpeg-statigo/blob/main/examples/transcode-hl/main.go).

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

// Flush the decoder, filter graph, and encoder before WriteTrailer so delayed
// tail frames are not dropped. See examples/transcode-hl/main.go for the full
// flush sequence (decoder.Flush -> filter Push(nil)/Pull -> encoder.Flush).

return output.WriteTrailer()
```

See [`examples/transcode-hl/main.go`](https://github.com/linuxmatters/ffmpeg-statigo/blob/main/examples/transcode-hl/main.go) for the full multi-stream version with decoder, filter, and encoder flushing at end of stream.

## Related pages

- [Codec and hardware matrix](Codec-and-Hardware-Matrix) for supported hardware acceleration backends by platform and codec.
- [Architecture and bindings](Architecture-and-Bindings) for raw binding memory rules and CGO conventions.
