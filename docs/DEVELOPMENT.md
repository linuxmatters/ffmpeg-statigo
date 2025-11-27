# ffmpeg-statigo Developer Guide

## Why Manual Setup?

Go modules don't handle 100MB+ binary files. The `go get` command won't fetch static libraries from the module cache. Our solution: git submodule + replace directive + library download.

## Prerequisites

- Go 1.24+
- GCC (CGO compilation)
- Git

## Consumer Integration

Projects using ffmpeg-statigo follow this pattern. Reference implementations with justfiles and GitHub workflows:

- [jivedrop](https://github.com/linuxmatters/jivedrop) — *Drop your podcast .wav into a neat MP3, ship the show metadata, cover art, and all 🪩*
- [jivefire](https://github.com/linuxmatters/jivefire) — *Spin your podcast .wav into a groovy MP4 visualiser. Cava-inspired real-time audio frequencies 🔥*

### Project Structure

```
your-project/
├── go.mod                    # Contains replace directive
├── justfile                  # setup recipe downloads libraries
└── third_party/
    └── ffmpeg-statigo/       # Git submodule
```

### go.mod Replace Directive

```
go mod edit -replace github.com/linuxmatters/ffmpeg-statigo=./third_party/ffmpeg-statigo
```

### Setup Sequence

1. `git submodule update --init --recursive`
2. `cd third_party/ffmpeg-statigo && go run ./cmd/download-lib`

Wrap this in a `just setup` recipe for consistent developer experience.

### Example justfile

```just
# First-time setup: init submodule + download FFmpeg libraries
setup:
    git config pull.ff only
    git config submodule.recurse true
    git submodule update --init --recursive
    cd third_party/ffmpeg-statigo && go run ./cmd/download-lib

# Build the project
build:
    CGO_ENABLED=1 go build -o myapp ./cmd/myapp
```

### Example GitHub Workflow

Native builds per platform avoid cross-compilation complexity:

```yaml
name: Build
on: [push, pull_request]

jobs:
  build:
    name: Build ${{ matrix.os }} ${{ matrix.arch }}
    runs-on: ${{ matrix.runner }}
    strategy:
      matrix:
        include:
          - os: linux
            arch: amd64
            runner: ubuntu-24.04
          - os: linux
            arch: arm64
            runner: ubuntu-24.04-arm
          - os: darwin
            arch: amd64
            runner: macos-15-intel
          - os: darwin
            arch: arm64
            runner: macos-15
      fail-fast: false
    steps:
      - uses: actions/checkout@v6
        with:
          submodules: recursive

      - uses: actions/setup-go@v6
        with:
          go-version: '1.24'

      - name: Download FFmpeg libraries
        run: cd third_party/ffmpeg-statigo && go run ./cmd/download-lib

      - name: Build
        run: go build -v ./...

      - name: Test
        run: go test -v ./...
```

### Cross-Compilation

Set `GOOS` and `GOARCH` before downloading: `GOOS=darwin GOARCH=arm64 go run ./cmd/download-lib`

# ffmpeg-statigo Development

For working on ffmpeg-statigo itself:

1. Clone the repository
2. `go run ./cmd/download-lib` (downloads pre-built libraries)
3. `just build` (full rebuild from source + regenerate bindings)

## Build System

All builds go through `just`. Never use `go build` directly—the justfile handles CGO flags and build sequencing.

| Command | Purpose |
|---------|---------|
| `just build` | Build static library from source, regenerate bindings, compile |
| `just test` | Run tests |
| `just generate` | Regenerate Go bindings from headers |
| `just download-lib` | Download pre-built libraries |

## Project Layout

| Path | Description |
|------|-------------|
| `*.gen.go` | Auto-generated bindings (do not edit) |
| `ffmpeg.go` | Core CGO directives, helper types |
| `include/` | FFmpeg C headers |
| `lib/<os>_<arch>/` | Platform-specific static libraries |
| `internal/builder/` | Builds FFmpeg + dependencies from source |
| `internal/generator/` | Generates Go bindings from headers using libclang |
| `cmd/download-lib/` | Downloads pre-built libraries from GitHub Releases |
| `examples/` | Working examples (transcode, metadata, etc.) |

## Versioning

Two distinct version schemes:

- **Library releases** (`lib-8.0.0.x`): Static library builds, distributed via GitHub Releases
- **Module releases** (`v8.0.0.x`): Go module versions

The download tool automatically fetches the latest `lib-8.0.0.x` release matching the FFmpeg version in `lib/fetch.go`.

## Troubleshooting

### "cannot find -lffmpeg"

Libraries not downloaded. Run `go run ./cmd/download-lib` in the ffmpeg-statigo directory.

### CGO errors on NixOS

Enter the dev shell: `nix develop` (or let `direnv` activate automatically).

### Module cache issues

If `go get` was attempted without the replace directive: `go clean -modcache`
