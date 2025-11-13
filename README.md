# ffmpeg-statigo

**Real FFmpeg bindings for Go. Not a wrapper. Not a CLI tool. The actual libraries.**

Cross-platform, static FFmpeg libraries bundled directly into your Go binary.
Hardware acceleration included. Zero runtime dependencies. Ship it and forget it.

## Why This Exists

Every other Go ffmpeg projects wrap the `ffmpeg` command, ffmpeg-statigo gives you the actual FFmpeg C libraries with proper Go bindings.
Build once, deploy anywhere. No hunting for system FFmpeg. No version mismatches. Predictable codec support.

## Features

- **FFmpeg 8.0** - Latest release with AV1, H.265, H.264, VP8/9
- **Truly static** - Builds into your binary (just needs system `m`ath library)
- **Cross-platform** - Linux and macOS (arm64, amd64)
- **Hardware acceleration** - NVENC/NVDEC, VideoToolbox, Vulkan and QuickSync support
- **GPL build** - x264, x265, and all the good codecs included
- **Auto-generated** - Thin, predictable bindings directly from FFmpeg headers
- **Preserved documentation** - Original FFmpeg comments in your IDE

*Hard fork of the excellent [csnewman/ffmpeg-go](https://github.com/csnewman/ffmpeg-go), modernised with FFmpeg 8.0, Go 1.24, hardware acceleration and a 99.5% smaller git history.*

## Setup

```
go get github.com/linuxmatters/ffmpeg-statigo
```

## Example

![example.gif](example.gif)

(The asciiplayer demo playing [Big Buck Bunny](https://en.wikipedia.org/wiki/Big_Buck_Bunny), with some GIF artifacts)

The `examples` directory contains some ported and some custom examples based on the C docs.

An example of printing a file's metadata:

```go
package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/csnewman/ffmpeg-go"
)

func main() {
	slog.Info("Metadata")

	var ctx *ffmpeg.AVFormatContext

	url := ffmpeg.ToCStr(os.Args[1])
	defer url.Free()

	_, err := ffmpeg.AVFormatOpenInput(&ctx, url, nil, nil)
	if err != nil {
		log.Panicln(err)
	}

	defer ffmpeg.AVFormatFreeContext(ctx)

	if _, err := ffmpeg.AVFormatFindStreamInfo(ctx, nil); err != nil {
		log.Panicln(err)
	}

	ffmpeg.AVDumpFormat(ctx, 0, url, 0)

	streams := ctx.Streams()

	for i := uintptr(0); i < uintptr(ctx.NbStreams()); i++ {
		s := streams.Get(i)

		slog.Info("  Stream", "i", i, "id", s.Id(), "dur", s.Duration())

		meta := s.Metadata()

		var entry *ffmpeg.AVDictionaryEntry

		for {
			entry = ffmpeg.AVDictGet(meta, ffmpeg.GlobalCStr(""), entry, ffmpeg.AVDictIgnoreSuffix)
			if entry == nil {
				break
			}

			slog.Info("    Meta", "key", entry.Key(), "value", entry.Value())
		}
	}
}
```

```
2024/01/15 21:39:43 INFO Metadata
Input #0, mov,mp4,m4a,3gp,3g2,mj2, from './bbb.mp4':
  Metadata:
    major_brand     : isom
    minor_version   : 1
    compatible_brands: isomavc1
    creation_time   : 2013-12-16T17:59:32.000000Z
    title           : Big Buck Bunny, Sunflower version
    artist          : Blender Foundation 2008, Janus Bager Kristensen 2013
    comment         : Creative Commons Attribution 3.0 - http://bbb3d.renderfarming.net
    genre           : Animation
    composer        : Sacha Goedegebure
  Duration: 00:10:34.57, start: 0.000000, bitrate: 4486 kb/s
  Stream #0:0[0x1](und): Video: h264 (High) (avc1 / 0x31637661), yuv420p(progressive), 1920x1080 [SAR 1:1 DAR 16:9], 4001 kb/s, 60 fps, 60 tbr, 60k tbn (default)
    Metadata:
      creation_time   : 2013-12-16T17:59:32.000000Z
      handler_name    : GPAC ISO Video Handler
      vendor_id       : [0][0][0][0]
  Stream #0:1[0x2](und): Audio: mp3 (mp4a / 0x6134706D), 48000 Hz, stereo, fltp, 160 kb/s (default)
    Metadata:
      creation_time   : 2013-12-16T17:59:37.000000Z
      handler_name    : GPAC ISO Audio Handler
      vendor_id       : [0][0][0][0]
  Stream #0:2[0x3](und): Audio: ac3 (ac-3 / 0x332D6361), 48000 Hz, 5.1(side), fltp, 320 kb/s (default)
    Metadata:
      creation_time   : 2013-12-16T17:59:37.000000Z
      handler_name    : GPAC ISO Audio Handler
      vendor_id       : [0][0][0][0]
    Side data:
      audio service type: main
2024/01/15 21:39:43 INFO   Stream i=0 id=1 dur=38072000
2024/01/15 21:39:43 INFO     Meta key=creation_time value=2013-12-16T17:59:32.000000Z
2024/01/15 21:39:43 INFO     Meta key=language value=und
2024/01/15 21:39:43 INFO     Meta key=handler_name value="GPAC ISO Video Handler"
2024/01/15 21:39:43 INFO     Meta key=vendor_id value=[0][0][0][0]
2024/01/15 21:39:43 INFO   Stream i=1 id=2 dur=30441600
2024/01/15 21:39:43 INFO     Meta key=creation_time value=2013-12-16T17:59:37.000000Z
2024/01/15 21:39:43 INFO     Meta key=language value=und
2024/01/15 21:39:43 INFO     Meta key=handler_name value="GPAC ISO Audio Handler"
2024/01/15 21:39:43 INFO     Meta key=vendor_id value=[0][0][0][0]
2024/01/15 21:39:43 INFO   Stream i=2 id=3 dur=30438912
2024/01/15 21:39:43 INFO     Meta key=creation_time value=2013-12-16T17:59:37.000000Z
2024/01/15 21:39:43 INFO     Meta key=language value=und
2024/01/15 21:39:43 INFO     Meta key=handler_name value="GPAC ISO Audio Handler"
2024/01/15 21:39:43 INFO     Meta key=vendor_id value=[0][0][0][0]
```

## Library versions

| Library  | Version |
|----------|---------|
| FFmpeg   | 8.0     |
| ass      | 0.17.4  |
| brotli   | 1.2.0   |
| bz2      | 1.0.8   |
| freetype | 2.14.1  |
| fribidi  | 1.0.16  |
| harfbuzz | 12.2.0  |
| mp3lame  | 3.100   |
| ogg      | 1.3.6   |
| opus     | 1.5.2   |
| png      | 1.6.50  |
| speex    | 1.2.1   |
| theora   | 1.2.0   |
| unibreak | 6.1     |
| vpx      | 1.15.0  |
| x264     | head    |
| x265     | head    |
| dav1d    | 1.5.2   |
| rav1e    | 0.8.2   |
| zlib     | 1.3.1   |

### Codec Inclusion Policy

ffmpeg-statigo provides a curated FFmpeg static library with Go bindings for contemporary audio and video processing, encoding, conversion and streaming.

If you're working on a project in one of thes domains ffmpeg-go should be equipped with you need:
- Streaming platforms (live and on-demand)
- Content management systems
- Social media applications
- Video conferencing services
- Media transcoding pipelines
- Web-based media players
- Broadcasting tools
- Modern content creation workflows

If you feel a library or codec is missing that should be include, open an issue to discuss it's inclusion.

#### Why is codec `xyz` not enabled?

This is a purpose-built media toolkit for Go developers building modern applications—optimised for size, focused on current workflows, and opinionated about what matters in 2025 and beyond.

Codecs that support these use case are disabled to help keep the static library size reasonable:
- Video archivists preserving historical formats
- Retro gaming emulator developers
- Digital archaeology projects
- Legacy format conversion utilities
- Film restoration specialists

## Hardware Acceleration Support Matrix

| Codec          | NVENC (Linux)    | QuickSync (Linux) | VideoToolbox (macOS) | Vulkan Video (Cross-platform) |
|----------------|------------------|-------------------|----------------------|-------------------------------|
| **AV1**        | ✅ Encode/Decode | ✅ Encode/Decode  | ☑️ Decode            | ✅ Encode/Decode              |
| **H.266/VVC**  | ❌               | ☑️ Decode         | ❌                   | ☑️ Decode                     |
| **H.265/HEVC** | ✅ Encode/Decode | ✅ Encode/Decode  | ✅ Encode/Decode     | ✅ Encode/Decode              |
| **H.264/AVC**  | ✅ Encode/Decode | ✅ Encode/Decode  | ✅ Encode/Decode     | ✅ Encode/Decode              |
| **VP9**        | ✅ Encode/Decode | ✅ Encode/Decode  | ❌                   | ☑️ Decode                     |
| **VP8**        | ☑️ Dec️ode        | ☑️ Decode         | ❌                   | ❌                            |
| **MPEG-2**     | ☑️ Decode        | ✅ Encode/Decode  | ❌                   | ❌                            |
| **JPEG/MJPEG** | ☑️ Decode        | ✅ Encode/Decode  | ❌                   | ❌                            |

### Capabilities

- **NVENC/NVDEC**: [Most NVIDIA GPUs come with NVENC/NVDEC support](https://developer.nvidia.com/video-encode-decode-support-matrix) but some low-end and mobile models are exceptions.
  - Decoding & Encoding H.264 8-bit - Any NVIDIA GPU supporting NVENC/NVDEC
  - Decoding & Encoding HEVC 8-bit - Maxwell 2nd Gen (GM206) and newer
  - Decoding HEVC 10-bit - Maxwell 2nd Gen (GM206) and newer
  - Encoding HEVC 10-bit - Pascal and newer
  - Decoding AV1 8/10-bit - Ampere and newer
  - Encoding AV1 8/10-bit - Ada Lovelace and newer
- **QuickSync (QSV)**: Requires Intel CPU (6th gen Skylake+) or Intel Arc GPU. Uses libvpl/oneVPL dispatcher.
  - Decoding & Encoding H.264 8-bit - Any Intel GPU that supports Quick Sync Video
  - Decoding & Encoding HEVC 8-bit - Gen 9 Skylake (6th Gen Core) and newer
  - Decoding & Encoding HEVC 10-bit - Gen 9.5 Kaby Lake (7th Gen Core), Apollo Lake, Gemini Lake (Pentium and Celeron) and newer
  - Decoding AV1 8/10-bit - Gen 12 Tiger Lake (11th Gen Core) and newer
  - Encoding AV1 8/10-bit - Gen 12.5 DG2 / ARC A-series, Gen 12.7 Meteor Lake (14th Gen Core Mobile / 1st Gen Core Ultra) and newer
  - VP9 requires 7th gen Kaby Lake or newer
- **VideoToolbox**: Available on macOS with Apple Silicon or Intel Macs with hardware support.
  - Decoding & Encoding H.264 8-bit - Any VideoToolbox-supported Mac.
  - Decoding & Encoding HEVC 8/10-bit - Macs from 2017 and later
  - Decoding AV1 8/10-bit - Requires an M3 series Apple Silicon Mac
- **Vulkan Video**: Works with any GPU that has Vulkan 1.3+ drivers (NVIDIA, AMD, Intel).

## Licensing

The Go binding code is MIT licensed. However, the bundled FFmpeg libraries are compiled with GPL-licensed components like `x264` and `x265`.
Any project using ffmpeg-statigo inherits the GPL requirements from FFmpeg through this linking, making the combined work subject to GPLv3 licensing obligations.
