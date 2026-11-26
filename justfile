import 'just/loader.just'

# Default recipe (shows available commands)
default:
    @just --list

# Clean build artifacts, preserving downloads and library archives
clean:
    @if [ ! -L .build ]; then for root in .build .build/embedded; do if [ ! -L "$root" ]; then rm -rf "$root/build" "$root/src" "$root/staging" 2>/dev/null || true; fi; done; fi
    @for name in asciiplayer hwdecode introspect metadata transcode transcode-hl; do rm -f "examples/$name/$name" "$name" 2>/dev/null || true; done
    @rm -f generator builder 2>/dev/null || true
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
    just generate
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

# Trigger Go module release
go-release VERSION:
    #!/usr/bin/env bash
    set -euo pipefail
    if [[ ! "{{VERSION}}" =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        echo "Error: VERSION must be in format X.Y.Z.N (e.g., 8.1.1.0)"
        exit 1
    fi
    gh workflow run go-release.yml -f version={{VERSION}}
