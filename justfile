# Default recipe (shows available commands)
default:
    @just --list

# Clean build artifacts and downloads
clean:
    rm -rf .build/

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
    go build -v ./examples/introspect/
    go build -v ./examples/metadata/
    go build -v ./examples/asciiplayer/
    go build -v ./examples/transcode/

# Build everything
build:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "→ Step 1: Building FFmpeg library..."
    just build-lib ffmpeg --clean
    just build-lib
    echo "→ Step 2: Regenerating Go bindings..."
    just generate
    echo "→ Step 3: Building Go packages..."
    go build -v ./...
    echo "→ Step 4: Building examples..."
    just build-examples
    echo "→ Step 5: Running introspection tool..."
    ./introspect

# Generate Go bindings from FFmpeg headers using libclang
generate:
    @go run ./internal/generator 2>&1 | grep -v "cgo-gcc-prolog\|deprecated" || true


# Run tests
test:
    go test -v ./...
