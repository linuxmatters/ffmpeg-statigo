# Examples

Working examples for ffmpeg-statigo, the Go CGO bindings to FFmpeg 8.1.

| Example | Description | Run |
|---------|-------------|-----|
| [asciiplayer](asciiplayer/) | Decodes video and renders frames as ASCII in a tcell terminal | `./bin/asciiplayer <file>` |
| [introspect](introspect/) | Lists every codec, format, filter, protocol, BSF, and hardware accelerator compiled into the library; also generates FFmpeg configure flags | `./bin/introspect` or `./bin/introspect --enable <codec>` |
| [metadata](metadata/) | Opens a media file and dumps container and per-stream metadata | `./bin/metadata <file>` |
| [transcode](transcode/) | Port of `doc/examples/transcode.c`: full decode/filter/encode/mux pipeline written directly against the raw bindings | `./bin/transcode <input> <output>` |
| [transcode-hl](transcode-hl/) | The same pipeline rewritten on the `av` package (`Input`/`Decoder`/`FilterGraph`/`Encoder`/`Output`): the recommended starting point | `./bin/transcode-hl <input> <output>` |
| [hwdecode](hwdecode/) | Hardware-accelerated video decode via the `av` package; auto-selects the first usable hardware device, transfers frames to software memory, and falls back to software decoding when no hardware is available | `./bin/hwdecode <input> [device-type]` |

## Prerequisites

All examples require the FFmpeg static libraries and must be built inside the Nix development shell.

**1. Enter the development shell**

```bash
nix develop
```

Direnv activates this automatically if you have it configured.

**2. Download the static libraries** (first time, or after a clean)

```bash
go run ./cmd/download-lib
```

**3. Build all examples**

```bash
just build
```

The binaries are written to `bin/`, including all examples, `download-lib`, `builder`, and `generator`.

> All builds require `CGO_ENABLED=1`. The `just build` target sets the correct environment. Do not use `go build` directly.

## Running an example

From the repo root, after building:

```bash
./bin/metadata /path/to/video.mp4
./bin/transcode-hl input.mp4 output.mkv
```

Or change into the example directory first:

```bash
cd examples/transcode-hl
../../bin/transcode-hl input.mp4 output.mkv
```
