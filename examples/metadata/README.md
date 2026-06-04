# metadata

Opens a media file, probes its streams, and dumps the container format info and per-stream metadata to stdout.

Back: [examples/](../README.md)

## What it demonstrates

- Opening a container with `AVFormatOpenInput` and probing streams with `AVFormatFindStreamInfo`
- Calling `AVDumpFormat` to print a human-readable format summary (the same output `ffmpeg -i` shows)
- Iterating over `AVStream` entries and reading `id`, `Duration`, and the metadata `AVDictionary`
- Walking an `AVDictionary` with `AVDictGet` using `AVDictIgnoreSuffix` to enumerate every key/value entry
- Using `GlobalCStr` for interned C string literals that do not need a `Free` call

## Run signature

```
metadata <file>
```

`<file>` accepts any URL that FFmpeg supports.

## Build

From the repo root, inside `nix develop`:

```bash
just build-examples
```

The binary is written to `examples/metadata/metadata`.

> The static libraries must be present first. Run `go run ./cmd/download-lib` if you have not done so already.

## Running

```bash
./examples/metadata/metadata /path/to/video.mp4
```

## Expected output

`AVDumpFormat` writes a format summary to FFmpeg's log (stderr by default). `slog` lines go to stdout. For a typical MP4 with one video and one audio stream you will see something like:

```
2024/01/01 12:00:00 INFO Metadata
Input #0, mov,mp4,m4a,3gp,3g2,mj2, from 'video.mp4':
  Duration: 00:01:23.45, start: 0.000000, bitrate: 2048 kb/s
    Stream #0:0(und): Video: h264, yuv420p, 1920x1080, ...
    Stream #0:1(und): Audio: aac, 48000 Hz, stereo, ...
2024/01/01 12:00:00 INFO   Stream i=0 id=1 dur=7512576
2024/01/01 12:00:00 INFO     Meta key=language value=und
2024/01/01 12:00:00 INFO   Stream i=1 id=2 dur=3993600
2024/01/01 12:00:00 INFO     Meta key=language value=und
2024/01/01 12:00:00 INFO     Meta key=handler_name value=SoundHandler
```

Streams with no metadata dictionary entries produce no `Meta` lines.
