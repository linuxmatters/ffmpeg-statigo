import 'just/loader.just'

# Default recipe (shows available commands)
default:
    @just --list

# Clean build artifacts and downloads
clean:
    @rm -rf .build/{build,src,staging} 2>/dev/null || true
    @rm examples/asciiplayer/asciiplayer 2>/dev/null || true
    @rm examples/hwdecode/hwdecode 2>/dev/null || true
    @rm examples/introspect/introspect 2>/dev/null || true
    @rm examples/metadata/metadata 2>/dev/null || true
    @rm examples/transcode/transcode 2>/dev/null || true
    @rm examples/transcode-hl/transcode-hl 2>/dev/null || true
    @rm asciiplayer 2>/dev/null || true
    @rm hwdecode 2>/dev/null || true
    @rm introspect 2>/dev/null || true
    @rm metadata 2>/dev/null || true
    @rm transcode 2>/dev/null || true
    @rm transcode-hl 2>/dev/null || true
    @rm generator 2>/dev/null || true
    @rm builder 2>/dev/null || true
    @if [ -d bin ] && [ ! -L bin ]; then for name in asciiplayer hwdecode introspect metadata transcode transcode-hl download-lib builder generator; do rm -f "bin/$name" "bin/$name.exe" 2>/dev/null || true; done; fi

# Build FFmpeg static library
build-static +args='':
    #!/usr/bin/env bash
    set -euo pipefail
    GOOS=$(go env GOOS)
    GOARCH=$(go env GOARCH)
    mkdir -p "lib/${GOOS}_${GOARCH}"
    go run ./internal/builder {{args}}

# Build embedded FFmpeg static library (then use GOFLAGS=-tags=embedded with just build or just test)
build-embedded +args='':
    #!/usr/bin/env bash
    set -euo pipefail
    GOOS=$(go env GOOS)
    GOARCH=$(go env GOARCH)
    mkdir -p "lib/embedded/${GOOS}_${GOARCH}"
    go run ./internal/builder --embedded {{args}}

# Build FFmpeg, regenerate bindings, build examples and inspect capabilities
build-ffmpeg:
    #!/usr/bin/env bash
    set -euo pipefail
    just build-static ffmpeg --clean
    just build-static
    go run ./internal/generator
    go build -a -v ./...
    just build
    ./bin/introspect

# Generate Go bindings
generate:
    go run ./internal/generator

# Download FFmpeg static libraries
download-lib:
    go run ./cmd/download-lib/

# Trigger FFmpeg library release
ffmpeg-release VERSION:
    #!/usr/bin/env bash
    set -euo pipefail
    if [[ ! "{{VERSION}}" =~ ^lib-[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        echo "Error: VERSION must start with 'lib-' and match format lib-X.Y.Z.N (e.g., lib-8.1.1.0)"
        exit 1
    fi
    gh workflow run ffmpeg-release.yml -f version={{VERSION}}

# Check library release workflow status
ffmpeg-release-status:
    gh run list --workflow=ffmpeg-release.yml --limit 5

# Trigger Go module release
go-release VERSION:
    #!/usr/bin/env bash
    set -euo pipefail
    if [[ ! "{{VERSION}}" =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        echo "Error: VERSION must be in format X.Y.Z.N (e.g., 8.1.1.0)"
        exit 1
    fi
    gh workflow run go-release.yml -f version={{VERSION}}

# Check Go module release workflow status
go-release-status:
    gh run list --workflow=go-release.yml --limit 5

# Check the static FFmpeg library is present (lint/vet need CGO to compile)
_check-lib:
    #!/usr/bin/env bash
    if [ ! -f "lib/$(go env GOOS)_$(go env GOARCH)/libffmpeg.a" ]; then
        echo "Error: static library missing. Run 'just download-lib' first."
        exit 1
    fi
