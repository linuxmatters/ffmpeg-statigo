# ffmpeg-statigo Developer Guide

## Why Manual Setup?

Go modules don't handle 100MB+ binary files. The `go get` command won't fetch static libraries from the module cache. Our solution: git submodule + replace directive + library download.

## Prerequisites

- Go 1.24+
- GCC (CGO compilation)
- Git

## Consumer Integration

Projects using ffmpeg-statigo follow this pattern. Reference implementations with justfiles and GitHub workflows:

- [jivedrop](https://github.com/linuxmatters/jivedrop) — *Drop your podcast .wav into a shiny MP3 with metadata, cover art, and all 🪩*
- [jivefire](https://github.com/linuxmatters/jivefire) — *Spin your podcast .wav into a groovy MP4 visualiser with Cava-inspired real-time audio frequencies 🔥*

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
# Check ffmpeg-statigo submodule is present
_check-submodule:
    #!/usr/bin/env bash
    if [ ! -f "third_party/ffmpeg-statigo/go.mod" ]; then
        echo "Error: ffmpeg-statigo submodule not initialised. Run 'just setup' first."
        exit 1
    fi
    if [ ! -f "third_party/ffmpeg-statigo/lib/$(go env GOOS)_$(go env GOARCH)/libffmpeg.a" ]; then
        echo "Error: ffmpeg-statigo library not downloaded. Run 'just setup' first."
        exit 1
    fi

# Get latest stable ffmpeg-statigo release tag from GitHub
_get-latest-tag:
    #!/usr/bin/env bash
    if command -v jq &> /dev/null; then
        curl -s https://api.github.com/repos/linuxmatters/ffmpeg-statigo/releases | \
            jq -r '[.[] | select(.prerelease == false and .draft == false and (.tag_name | startswith("v")))][0].tag_name'
    else
        curl -s https://api.github.com/repos/linuxmatters/ffmpeg-statigo/releases | \
            grep -B5 '"prerelease": false' | grep '"tag_name"' | grep -v 'lib-' | head -1 | cut -d'"' -f4
    fi

# Setup or update ffmpeg-statigo submodule and library
setup:
    #!/usr/bin/env bash
    set -e
    echo "Configuring git for submodule-friendly pulls..."
    git config pull.ff only
    git config submodule.recurse true

    # Get latest stable release tag
    TAG=$(just _get-latest-tag)
    if [ -z "$TAG" ] || [ "$TAG" = "null" ]; then
        echo "Error: Could not fetch latest release tag"
        exit 1
    fi

    # Initialise submodule if not already present
    if [ ! -f "third_party/ffmpeg-statigo/go.mod" ]; then
        echo "Initialising ffmpeg-statigo submodule..."
        git submodule update --init --recursive
    fi

    # Check current version
    cd third_party/ffmpeg-statigo
    git fetch --tags
    CURRENT=$(git describe --tags --exact-match 2>/dev/null || echo "")

    if [ "$CURRENT" = "$TAG" ]; then
        echo "ffmpeg-statigo already at latest version ($TAG)"
        cd ../..
    else
        if [ -n "$CURRENT" ]; then
            echo "Updating ffmpeg-statigo from $CURRENT to $TAG..."
        else
            echo "Setting up ffmpeg-statigo $TAG..."
        fi
        git checkout "$TAG"
        cd ../..

        # Remove old library to force re-download
        rm -f third_party/ffmpeg-statigo/lib/*/libffmpeg.a

        # Stage the submodule change
        git add third_party/ffmpeg-statigo
    fi

    # Download libraries (will skip if correct version already exists)
    echo "Checking ffmpeg-statigo libraries..."
    cd third_party/ffmpeg-statigo && go run ./cmd/download-lib
    cd ../..

    # Check if there are staged changes to commit
    if git diff --cached --quiet third_party/ffmpeg-statigo; then
        echo "Setup complete!"
    else
        echo ""
        echo "Setup complete! Submodule updated to $TAG"
        echo "Don't forget to commit: git commit -m 'chore: update ffmpeg-statigo to $TAG'"
    fi

# Build the project
build: _check-submodule
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
