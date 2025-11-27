# ffmpeg-statigo Copilot Instructions

## Project Overview

Go bindings for FFmpeg 8.0 static libraries. Ships ~100MB platform-specific `libffmpeg.a` files that link directly into Go binaries—no system FFmpeg required.

## Critical Build Rules

**Always use `just build`**—never run `go build` directly. The justfile orchestrates:
1. Building FFmpeg static library from source (`internal/builder/`)
2. Regenerating Go bindings from headers (`internal/generator/`)
3. Compiling Go code with correct CGO flags

## Code Generation

Files matching `*.gen.go` are **auto-generated**—never edit directly:
- `constants.gen.go`, `enums.gen.go`, `structs.gen.go`, `functions.gen.go`, `callbacks.gen.go`

To regenerate after header changes: `just generate`

The generator (`internal/generator/`) parses FFmpeg headers in `include/` using libclang and outputs Go bindings.

## Architecture

| Component | Purpose |
|-----------|---------|
| `ffmpeg.go` | CGO directives, platform-specific linker flags, helper types (`CStr`, `AVError`) |
| `*.gen.go` | Generated FFmpeg API bindings |
| `include/` | FFmpeg C headers (libavcodec, libavformat, libavutil, etc.) |
| `lib/<os>_<arch>/` | Platform-specific static libraries (gitignored) |
| `internal/builder/` | Builds FFmpeg + 20 dependencies from source |
| `internal/generator/` | Generates Go bindings from headers |
| `cmd/download-lib/` | Downloads pre-built libraries from GitHub Releases |
| `examples/` | Working examples: transcode, metadata, asciiplayer, introspect |

## Internal Builder

The builder (`internal/builder/`) compiles FFmpeg and dependencies from source. Key files:

- `libraries.go` — Library definitions (URLs, versions, build configs)
- `buildsystems.go` — Autoconf, CMake, Meson, Cargo build implementations
- `main.go` — CLI entry point with `--clean` and `--list` flags

Build commands:
- `just build-static` — Build all libraries
- `just build-static ffmpeg --clean` — Rebuild FFmpeg only
- `go run ./internal/builder --list` — Show library versions

## Hardware Acceleration

Supports NVENC/NVDEC (Linux), QuickSync (Linux), VideoToolbox (macOS), and Vulkan Video. See `README.md` and `docs/CODECS.md` for codec and hardware details.

## Library Distribution

Static libraries distributed via GitHub Releases (`lib-8.0.0.x` tags), not git. Download with:
```
go run ./cmd/download-lib
```

## Platform Support

Linux (amd64, arm64) and macOS (amd64, arm64). CGO required (`CGO_ENABLED=1`).

Platform-specific code uses build tags in CGO directives—see `ffmpeg.go` lines 11-17.

## Testing

```
just test
```

Tests require downloaded libraries. See `ffmpeg_test.go` for version validation pattern.

## Key Patterns

- All FFmpeg types prefixed with `AV*` (e.g., `AVCodecContext`, `AVFrame`)
- Use `CStr` type for C string interop—call `.Free()` when done
- Wrap FFmpeg return codes with `WrapErr()` to convert to Go errors
- Check `AVErrorEOF` and `EAgain` for stream processing loops

## Environment
- NixOS development shell via `flake.nix`
- Fish shell for terminal commands
- CGO required (`CGO_ENABLED=1` in build)

## Further Reading

- `docs/DEVELOPMENT.md` — Consumer integration, project layout, troubleshooting
- `docs/CODECS.md` — Enabled codecs, muxers, parsers

## Reference Projects

See [jivedrop](https://github.com/linuxmatters/jivedrop) and [jivefire](https://github.com/linuxmatters/jivefire) for consumer integration patterns with justfiles and GitHub workflows.
