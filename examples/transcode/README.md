# transcode

A direct port of FFmpeg's `doc/examples/transcode.c`. Reads every stream from an input container, decodes audio and video through a passthrough filter graph, re-encodes to the same codec, and muxes into an output container. Subtitle and data streams are remuxed without transcoding.

Back: [examples/](../README.md)

> **Recommended alternative:** [`transcode-hl`](../transcode-hl/README.md) implements the same pipeline on the `av` package and is the recommended starting point for new code. See below for a known issue with this example against the current FFmpeg build.

## What it demonstrates

- Manual demux/decode/filter/encode/mux pipeline written directly against the raw bindings
- Allocating and opening decoder and encoder contexts with `AVCodecAllocContext3` / `AVCodecOpen2`
- Building separate filter graphs per stream (`null` passthrough for video, `anull` for audio) using `AVFilterGraphCreateFilter` / `AVFilterGraphParsePtr` / `AVFilterGraphConfig`
- Setting buffersink format constraints with `AVOptSetSlice` (see note below)
- The send/receive loops at both the decode (`AVCodecSendPacket` / `AVCodecReceiveFrame`) and encode (`AVCodecSendFrame` / `AVCodecReceivePacket`) sides
- Timestamp rescaling with `AVPacketRescaleTs` and `AVRescaleQ`
- Flushing decoders, filter graphs, and encoders at end of stream
- Managing `AVPacket` and `AVFrame` lifetimes manually throughout

## Run signature

```
transcode <input> <output>
```

The output container format is inferred from the file extension of `<output>`.

```bash
./transcode input.mp4 output.mkv
```

## Build

From the repo root, inside `nix develop`:

```bash
just build-examples
```

The binary is written to `examples/transcode/transcode`.

> The static libraries must be present first. Run `go run ./cmd/download-lib` if you have not done so already.

## Known issue

This example does not run successfully against the pinned FFmpeg 8.1 build. The code sets `pix_fmts` and `sample_fmts` options on the buffersink filter context after `AVFilterGraphCreateFilter` but before `AVFilterGraphConfig`. Current libavfilter rejects these as non-runtime options at that point in the initialisation sequence and returns an error.

The `av` package's `FilterGraph` wrapper (`transcode-hl`) avoids this by passing format constraints through the filter graph constructor before `AVFilterGraphConfig` is called. Use `transcode-hl` for runnable code.

## Expected output (if the issue is resolved)

`slog` lines describe each stream, codec, and frame rate detected:

```
2024/01/01 12:00:00 INFO Transcode
2024/01/01 12:00:00 INFO streams nb=2
2024/01/01 12:00:00 INFO Stream i=0
2024/01/01 12:00:00 INFO  >  codec=h264
2024/01/01 12:00:00 INFO  >  type=video
2024/01/01 12:00:00 INFO  >  fr={25 1}
...
2024/01/01 12:00:00 INFO End of file
```

The output file is written to `<output>` using the encoder settings inferred from the input streams.
