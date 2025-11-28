# Default recipe (shows available commands)
default:
    @just --list

# Clean build artifacts and downloads
clean:
    @rm -rf .build/{build,src,staging} 2>/dev/null || true
    @rm examples/asciiplayer/asciiplayer 2>/dev/null || true
    @rm examples/introspect/introspect 2>/dev/null || true
    @rm examples/metadata/metadata 2>/dev/null || true
    @rm examples/transcode/transcode 2>/dev/null || true

# Build FFmpeg static library
build-static +args='':
    #!/usr/bin/env bash
    set -euo pipefail
    GOOS=$(go env GOOS)
    GOARCH=$(go env GOARCH)
    mkdir -p "lib/${GOOS}_${GOARCH}"
    go run ./internal/builder {{args}}

# Build example programs
build-examples:
    go build -v ./examples/asciiplayer/
    go build -v ./examples/introspect/
    go build -v ./examples/metadata/
    go build -v ./examples/transcode/

# Build everything
build:
    #!/usr/bin/env bash
    set -euo pipefail
    just build-static ffmpeg --clean
    just build-static
    go run ./internal/generator 2>&1 | grep -v "cgo-gcc-prolog\|deprecated" || true
    go build -a -v ./...
    just build-examples
    ./introspect

# Generate Go bindings
generate:
    go run ./internal/generator

# Run tests
test:
    go test -v ./...

# Download FFmpeg static libraries
download-lib:
    go run ./cmd/download-lib/

# Trigger FFmpeg library release
ffmpeg-release VERSION:
    #!/usr/bin/env bash
    set -euo pipefail
    if [[ ! "{{VERSION}}" =~ ^lib-[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        echo "Error: VERSION must start with 'lib-' and match format lib-X.Y.Z.N (e.g., lib-8.0.1.0)"
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
        echo "Error: VERSION must be in format X.Y.Z.N (e.g., 8.0.1.0)"
        exit 1
    fi
    gh workflow run go-release.yml -f version={{VERSION}}

# Check Go module release workflow status
go-release-status:
    gh run list --workflow=go-release.yml --limit 5
