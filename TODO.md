# TODO

- [ ] Refactor the internal build tool
- [ ] Reorganise project structure
- [ ] Rebrand to ffmpeg-statigo
- [ ] How to embed/distribute the static FFmpeg library?
- [ ] Enable Dolby Vision/HDR10+ metadata handling via `--enable-libdovi` (if available)
- [ ] Add `libwebp` - web-optimised image format (part of WebM ecosystem)
- [ ] Add `libavif` - AVIF image support (AV1-based)
- [ ] Enable Whisper for transcription workflows (New in FFmpgeg in 8.0 ✨)
- [ ] Enable VA-API? - Still required now we have NVENC, QuickSync and Vulkan?

## More headers

### Adding Headers to the Generator

1. Modify the Generator Configuration

The header list is defined in a generator configuration file `generator/parser.go`. AI dd entries like:

```go
headers = append(headers,
    "libswscale/swscale.h",
    "libswresample/swresample.h",
    "libavformat/rtmp.h",
    "libavdevice/avdevice.h"
)
```

2. Regenerate Bindings

Run the generator with the updated clang (which you've already fixed for Linux compatibility):

```bash
just generate
```

### Critical Missing Headers for Streaming

#### Network Protocol Headers

- `libavformat/rtmp*.h` - RTMP/RTMPS protocol internals for custom handshaking and stream key management
- `libavformat/srt*.h` - SRT (Secure Reliable Transport) for low-latency streaming, increasingly required by platforms
- `libavformat/hls*.h` & libavformat/dash*.h - For generating adaptive bitrate streams with proper segment control

#### Real-time Processing Headers

- `libavfilter/buffersrc.h` & `libavfilter/buffersink.h` - Essential for building custom filter graphs for overlays, watermarks, and scene transitions
- `libavdevice/avdevice.h` - Capture device integration (webcams, capture cards, screen recording)
- `libavutil/fifo.h` - Thread-safe FIFO buffers for managing multiple output streams

#### Advanced Streaming Control

- `libavformat/avio_internal.h` - Custom I/O contexts for authentication and stream routing
- `libavcodec/bsf.h` - Bitstream filters for H.264/HEVC Annex B conversion (required by some platforms)
- `libavutil/threadmessage.h` - Inter-thread communication for multiple simultaneous outputs

### Missing Headers for Scaling and Resampling

- `libswresample/swresample.h` and `libswscale/swscale.h`
  - Creating adaptive bitrate ladders - scaling from a single 1080p input to 720p, 480p, 360p variants
  - Audio normalisation - resampling between 44.1kHz and 48kHz, converting 5.1 to stereo
  - Pixel format conversion - converting between YUV420P, NV12, and platform-specific requirements

## From the original author

- [ ] Expose more headers.
- [ ] Expose platform specific headers.
- [ ] Cleanup internal packages.

## Testing

- [ ] Create some test cases that exercise some of the FFmpeg API surface

### Basic Validation
# Rebuild static librarie for FFmpeg
```bash
just build-ffmpeg
```

**Regenerate bindings:**
```bash
just generate
```

**Run existing tests:**
```bash
just test
```

### Examples Validation
```bash
# Run transcode example
cd examples/transcode
go run main.go

# Run asciiplayer example
cd examples/asciiplayer
go run main.go

# Run metadata example
cd examples/metadata
go run main.go
```

# Onboard

The is my hard fork of the [ffmpeg-go](https://github.com/csnewman/ffmpeg-go) project called ffmpeg-statigo.

ffmpeg-statigo has been updated to Go 1.24, all dependencies uplifted to current versions, GitHub CI is fixed so the ffmpeg static libraries are built for all supported architectures and the generator has been updated to support current `clang` on Linux. Mostly critically this project has been updated from FFmpeg 6.1 to FFmpeg 8.0 including the addition os `x265`, `dav1d`, `rav1e` and hardware acceleration support for NVENC, QuickSync and Vulkan to complement the exist VideoToolbox support.

The git history has been purged all tags and historical commits of the static ffmpeg libraries as we are preparing this project to be launched under it's new name in a different GitHUb organisation. This brought the git repo down in size from 1.9GB to 9MB.

NixOS is the host development workstation. There is a `flake.nix` in the project to enable the required software in a Nix devevelopment shell, which is automatically activated via `direnv`. The `justfile` is used for all build and test commands. `fish` is the default shell.

Read the `README.md`, `TODO.md` and analyse the code to get a full understanding of the project. Let me know when you are ready to collaborate.
