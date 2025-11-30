# ffmpeg-statigo

**Real FFmpeg bindings for Go. Not a wrapper. Not a CLI tool. The actual libraries.**

Cross-platform, static FFmpeg libraries bundled directly into your Go binary.
Hardware acceleration included. Zero runtime dependencies. Ship it and forget it.

## Why This Exists

Every other Go ffmpeg projects wrap the `ffmpeg` command, ffmpeg-statigo gives you the actual FFmpeg C libraries with proper Go bindings.
Build once, deploy anywhere. No hunting for system FFmpeg. No version mismatches. Predictable codec support.

## Features

- **FFmpeg 8.0.1** - Latest release with AV1, H.265, H.264, VP8/9
- **Truly static** - Builds into your binary (just needs system `m` and `stdc++` libraries)
- **Cross-platform** - Linux and macOS (arm64, amd64)
- **Hardware acceleration** - NVENC/NVDEC, QuickSync, VA-API, VideoToolbox, and Vulkan support
- **GPL build** - x264, x265, and all the good codecs included
- **Auto-generated** - Thin, predictable bindings directly from FFmpeg headers
- **Preserved documentation** - Original FFmpeg comments in your IDE

*Hard fork of the excellent [csnewman/ffmpeg-go](https://github.com/csnewman/ffmpeg-go), modernised with FFmpeg 8.0.x, Go 1.24, hardware acceleration, ~85% FFmpeg API coverage, and a 99.5% smaller git history.*

## Installation

Static libraries (~100MB per platform) cannot be distributed via `go get`. Use as a git submodule with a replace directive:

1. Add submodule: `git submodule add https://github.com/linuxmatters/ffmpeg-statigo third_party/ffmpeg-statigo`
2. Add replace directive to `go.mod`: `replace github.com/linuxmatters/ffmpeg-statigo => ./third_party/ffmpeg-statigo`
3. Download libraries: `cd third_party/ffmpeg-statigo && go run ./cmd/download-lib`
4. Configure git for submodule-friendly pulls: `git config pull.ff only && git config submodule.recurse true`
5. Build: `go build ./...`

Step 4 prevents `git pull --rebase` from breaking submodule references. Fast-forward only pulls ensure submodule commits stay in sync with the parent repository.

Static libraries are gitignored—only the submodule reference is committed.

**See [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) for CI/CD integration, cross-compilation, and troubleshooting.**

## Codec Inclusion Policy

ffmpeg-statigo provides a **curated FFmpeg static library** focused on the core strengths of FFmpeg: **decoding, processing, and encoding** audio and video streams. ffmpeg-statigo is designed for **Go developers building modern streaming applications**. The pattern is:

1. **Generate content in Go** - Text, graphics, effects using Go's excellent libraries
2. **Feed frames to FFmpeg** - Use ffmpeg-statigo for encoding and stream processing
3. **Let FFmpeg handle codecs** - Hardware acceleration, format conversion, container muxing

#### What's Included

- **Decoders**: All contemporary formats (H.264, H.265, AV1, VP8/9, Opus, AAC, MP3)
- **Encoders**: Modern codecs for streaming and transcoding (x264, x265, dav1d, rav1e, vpx, lame, opus)
- **Hardware acceleration**: NVENC, QuickSync, VA-API, VideoToolbox, Vulkan Video
- **Containers**: MP4, MKV, WebM, DASH, HLS, and all major formats
- **Filters**: Video scaling, colour conversion, audio resampling, and processing filters
- **Streaming protocols**: RTMP, SRT, HLS, DASH

**Building something that deals with video?** You're probably covered:
- Streaming platform (your own Twitch/YouTube/Owncast)
- Content management system with media handling
- Social media app with video uploads
- Video conferencing service
- Transcoding pipeline for your media library
- Home media server ripping DVDs/Blu-rays
- Web-based media player
- Broadcasting tool or modern content creation workflow

#### What's Intentionally Absent

Some FFmpeg features commonly found in the `ffmpeg` CLI tool are **not included** because they're better implemented in Go:

- ❌ **`drawtext` filter** - Use Go's `image/draw` + `golang.org/x/image/font` packages instead
- ❌ **Font libraries** (freetype, harfbuzz, fontconfig) - Not needed when rendering in Go
- ❌ **`subtitles` and `ass` filters** - Use separate subtitle streams instead
- ❌ **libass** - Subtitle rendering library not needed
  - FFmpeg can still **copy, extract, and mux** subtitle streams without libass.

The excluded features can be added back if a compelling use case emerges. ffmpeg-statigo is a **living project** and the curated library evolves based on real-world needs.

If you need complete FFmpeg with all filters, use the official FFmpeg distribution. If you need modern streaming codecs with Go integration, ffmpeg-statigo is designed for you.

### Library versions

| Library          | Version     | Description                                                                        |
|------------------|-------------|------------------------------------------------------------------------------------|
| FFmpeg           | 8.0.1       | A complete, cross-platform solution to record, convert and stream audio and video  |
| dav1d            | 1.5.2       | AV1 cross-platform decoder, open-source, and focused on speed, size and correctness|
| glslang          | 15.4.0      | Khronos-reference front end for GLSL/ESSL and a SPIR-V generator                   |
| libdrm           | 2.4.129     | Direct Rendering Manager library and headers (*Linux only*)                        |
| libiconv         | 1.18        | A character set conversion library (*macOS only*)                                  |
| libsrt           | 1.5.5-rc.0a | A transport protocol for ultra low latency live video and audio streaming          |
| libva            | 2.22.0      | An implementation for VA-API (Video Acceleration API) (*Linux only*)               |
| libvpl           | 2.15.0      | Intel Video Processing Library (Intel VPL) API (*Linux only*)                      |
| libvpx           | 1.15.2      | High-quality, open video format for the web that's freely available to everyone    |
| libwebp          | 1.6.0       | A modern image format providing superior lossless and lossy compression            |
| libxml2          | 2.15.1      | An XML parser and toolkit implemented in C                                         |
| mp3lame          | 3.100       | A high quality MPEG Audio Layer III (MP3) encoder                                  |
| nv-codec-headers | 12.2.72.0   | Headers required to interface with Nvidias codec APIs (*Linux only*)               |
| openssl          | 3.6.0       | Open Source Toolkit for the TLS, DTLS, and QUIC protocols.                         |
| opus             | 1.5.2       | A totally open, royalty-free, highly versatile audio codec                         |
| rav1e            | 0.8.1       | The fastest and safest AV1 encoder.                                                |
| Vulkan-Headers   | 1.4.335     | Vulkan header files and API registry                                               |
| x264             | head        | H.264/MPEG-4 AVC compression format library for encoding video streams             |
| x265             | head        | H.265/MPEG-H HEVC compression format library for encoding video streams            |
| zimg             | 3.0.6       | Scaling, colorspace conversion, and dithering library                              |
| zlib             | 1.3.1       | A Massively Spiffy Yet Delicately Unobtrusive Compression Library                  |

VVenC 1.13.1 (Fraunhofer Versatile Video Encoder, a fast & efficient software H.266/VVC encoder) is in the build configuration, **but currently disabled**, as it adds ~25MB to the static ffmpeg library and is too slow for practical use.

### Enabled Codecs

Details of codecs, muxers and parsers available in enable in the static ffmpeg library that ffmpeg-statigo ships are documented in [`docs/CODECS.md`](docs/CODECS.md).

### Hardware Acceleration Support Matrix

| Codec          | NVENC (Linux)    | VA-API (Linux)    | QuickSync (Linux) | VideoToolbox (macOS) | Vulkan Video (Cross-platform) |
|----------------|------------------|-------------------|-------------------|----------------------|-------------------------------|
| **AV1**        | ✅ Encode/Decode | ☑️ Encode         | ✅ Encode/Decode  | ☑️ Decode            | ✅ Encode/Decode              |
| **H.266/VVC**  | ❌               | ❌                | ☑️ Decode         | ❌                   | ☑️ Decode                     |
| **H.265/HEVC** | ✅ Encode/Decode | ☑️ Encode         | ✅ Encode/Decode  | ✅ Encode/Decode     | ✅ Encode/Decode              |
| **H.264/AVC**  | ✅ Encode/Decode | ☑️ Encode         | ✅ Encode/Decode  | ✅ Encode/Decode     | ✅ Encode/Decode              |
| **VP9**        | ✅ Encode/Decode | ☑️ Encode         | ✅ Encode/Decode  | ❌                   | ☑️ Decode                     |
| **VP8**        | ☑️ Decode        | ☑️ Encode         | ☑️ Decode         | ❌                   | ❌                            |
| **MPEG-2**     | ☑️ Decode        | ☑️ Encode         | ✅ Encode/Decode  | ❌                   | ❌                            |
| **MJPEG**      | ☑️ Decode        | ☑️ Encode         | ✅ Encode/Decode  | ❌                   | ❌                            |

### Capabilities

- **NVENC/NVDEC**: [Most NVIDIA GPUs come with NVENC/NVDEC support](https://developer.nvidia.com/video-encode-decode-support-matrix) but some low-end and mobile models are exceptions.
  - Decoding & Encoding H.264 8-bit - Any NVIDIA GPU supporting NVENC/NVDEC
  - Decoding & Encoding HEVC 8-bit - Maxwell 2nd Gen (GM206) and newer
  - Decoding HEVC 10-bit - Maxwell 2nd Gen (GM206) and newer
  - Encoding HEVC 10-bit - Pascal and newer
  - Decoding AV1 8/10-bit - Ampere and newer
  - Encoding AV1 8/10-bit - Ada Lovelace and newer
- **VA-API**: Video Acceleration API for Intel, AMD, and NVIDIA (via `nvidia-vaapi-driver`) GPUs on Linux. Provides hardware encoding; decoding uses FFmpeg's hwaccel framework.
  - Encoding H.264, HEVC, AV1, VP8, VP9, MPEG-2, MJPEG - Hardware dependent
  - Intel: Use `iHD` driver (intel-media-driver) for broadest codec support
  - AMD: Use `radeonsi` driver (Mesa) for RDNA/GCN GPUs
  - NVIDIA: Use `nvidia` driver (nvidia-vaapi-driver) which translates to NVENC
- **QuickSync (QSV)**: Requires Intel CPU (11th gen Tiger Lake+) or Intel Arc GPU. Uses libvpl/oneVPL dispatcher. Older Intel CPUs (6th-10th gen) should use VA-API instead.
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
- **Vulkan Video**: Works with any GPU that has Vulkan 1.3+ drivers.
  - Decoding & Encoding H.264 8-bit
  - Decoding & Encoding HEVC 8/10-bit
  - Decoding & Encoding AV1 8/10-bit
  - **Works via MoltenVK on macOS when MoltenVK runtime is installed**

## Licensing

The Go binding code is MIT licensed. However, the bundled FFmpeg libraries are compiled with GPL-licensed components like `x264` and `x265`.
Any project using ffmpeg-statigo inherits the GPL requirements from FFmpeg through this linking, making the combined work subject to GPLv3 licensing obligations.
