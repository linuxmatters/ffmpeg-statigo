# Default recipe (shows available commands)
default:
    @just --list

# Clean build artifacts and downloads
clean:
    rm -rf temp/

# Build FFmpeg static library for current platform
build-ffmpeg:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "Building FFmpeg static library..."
    GOOS=$(go env GOOS)
    GOARCH=$(go env GOARCH)
    go run ./internal/builder "libffmpeg_${GOOS}_${GOARCH}.a"

# Build all Go packages
build:
    go build -v ./...

# Generate Go bindings from FFmpeg headers using libclang
generate:
    @go run ./internal/generator 2>&1 | grep -v "cgo-gcc-prolog\|deprecated" || true

# Run tests
test:
    go test -v ./...
