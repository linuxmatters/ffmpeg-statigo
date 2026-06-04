# transcode-hl

The same decode/filter/encode/mux pipeline as [`transcode`](../transcode/README.md), rewritten on the `av` package (`github.com/linuxmatters/ffmpeg-statigo/av`). This is the recommended example for new code.

Back: [examples/](../README.md)

## What it demonstrates

- Using `av.Open` and `av.CreateOutput` for owned, `io.Closer`-managed demux/mux contexts
- Creating a `av.Decoder` from a stream with `av.NewDecoder` — finds, configures, and opens the codec in one call
- Building a passthrough filter graph with `av.NewVideoFilterGraphFromContext` (`"null"`) and `av.NewAudioFilterGraphFromContext` (`"anull"`), which allocate and configure the filter before `AVFilterGraphConfig` — avoiding the option-ordering issue present in [`transcode`](../transcode/README.md)
- The callback-based decode loop: `decoder.Decode(packet, func(frame) error { ... })`
- The callback-based filter drain: `filter.Pull(func(filtered) error { ... })`
- The callback-based encode loop: `encoder.Encode(frame, func(pkt) error { ... })`
- Adding an output stream with `output.AddStream(encoder)` — copies codec parameters and time base automatically
- Remuxing non-audio/video streams by dropping through to `output.Raw()` for the one operation (`AVFormatNewStream` + `AVCodecParametersCopy`) the `av` package does not wrap
- Flushing all stages at end of stream via `decoder.Flush`, `filter.Push(nil)`, `filter.Pull`, and `encoder.Flush`
- Using `Raw()` to reach the underlying `*ffmpeg.AV*` handle when the `av` layer does not surface a needed value (e.g. `output.Raw().Oformat().Flags()` for the global-header flag, `filter.Raw()` to read the buffersink time base)

For a full description of the `av` package API, see [docs/PIPELINE.md](../../docs/PIPELINE.md).

## Run signature

```
transcode-hl <input> <output>
```

The output container format is inferred from the file extension of `<output>`.

```bash
./transcode-hl input.mp4 output.mkv
```

## Build

From the repo root, inside `nix develop`:

```bash
just build-examples
```

The binary is written to `examples/transcode-hl/transcode-hl`.

> The static libraries must be present first. Run `go run ./cmd/download-lib` if you have not done so already.

## Expected output

`slog` lines log stream detection, packet processing, and end-of-file:

```
2024/01/01 12:00:00 INFO Transcode
2024/01/01 12:00:00 INFO End of file
```

The output file is written to `<output>`. Audio and video streams are transcoded to the same codec as the input; other stream types (subtitles, data) are remuxed without re-encoding.

## Comparison with transcode

| | `transcode` | `transcode-hl` |
|---|---|---|
| API layer | Raw `ffmpeg.*` bindings | `av` package wrappers |
| Lifetime management | Manual `Free`/`Unref` calls | `defer Close()` on each owned type |
| Filter graph setup | `AVFilterGraphCreateFilter` + `AVOptSetSlice` + `AVFilterGraphConfig` | `av.NewVideoFilterGraphFromContext` / `av.NewAudioFilterGraphFromContext` |
| Runnable against FFmpeg 8.1 | No (see [`transcode` known issue](../transcode/README.md#known-issue)) | Yes |
| Recommended for new code | No | Yes |
