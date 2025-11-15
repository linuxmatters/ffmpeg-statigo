# Default recipe (shows available commands)
default:
    @just --list

# Clean build artifacts and downloads
clean:
    rm -rf .build/
    rm examples/asciiplayer/asciiplayer
    rm examples/introspect/introspect
    rm examples/metadata/metadata
    rm examples/transcode/transcode

# Build FFmpeg static library
build-lib +args='':
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
    just build-lib ffmpeg --clean
    just build-lib
    just generate
    go build -v ./...
    just build-examples
    ./introspect

# Generate Go bindings from FFmpeg headers using libclang
generate:
    @go run ./internal/generator 2>&1 | grep -v "cgo-gcc-prolog\|deprecated" || true

# Run tests
test:
    go test -v ./...
