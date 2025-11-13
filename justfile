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

# Build all Go packages
build:
    go build -v ./...

# Build example programs
build-examples:
    go build -v ./examples/introspect/
    go build -v ./examples/metadata/
    go build -v ./examples/asciiplayer/
    go build -v ./examples/transcode/

# Build and run introspection tool
build-introspect:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "→ Step 1: Building FFmpeg library..."
    just build-lib ffmpeg --clean
    just build-lib
    echo "→ Step 2: Regenerating Go bindings from FFmpeg headers..."
    just generate
    echo "→ Step 3: Building Go packages..."
    just build
    echo "→ Step 4: Running introspection tool..."
    cd examples/introspect && go run main.go

# Generate Go bindings from FFmpeg headers using libclang
generate:
    @go run ./internal/generator 2>&1 | grep -v "cgo-gcc-prolog\|deprecated" || true


# Run tests
test:
    go test -v ./...
