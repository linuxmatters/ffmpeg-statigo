# Go consumer integration

Add `ffmpeg-statigo` to a Go project as a Git submodule with a `go.mod` `replace` directive, then fetch the static library with `cmd/download-lib`.

Go modules do not distribute 100 MB static libraries through `go get`. The submodule tracks the exact Go binding version, and `cmd/download-lib` downloads the matching pre-built `libffmpeg.a` into the gitignored `lib/` directory.

## Prerequisites

- Go 1.26 or newer
- A C compiler for CGO (GCC on Linux, Clang on macOS, or MinGW-w64 on Windows)
- Git
- A 64-bit target architecture (`linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`, or `windows/amd64`). `arch_guard.go` rejects 32-bit Go builds at compile time.

## Initial setup

Organise your repository with `ffmpeg-statigo` under `third_party/`:

```text
your-project/
├── go.mod                    # Contains the replace directive
├── justfile                  # setup recipe downloads libraries
└── third_party/
    └── ffmpeg-statigo/       # Git submodule
```

Run these commands from your project root:

1. Add the submodule:

   ```sh
   git submodule add https://github.com/linuxmatters/ffmpeg-statigo third_party/ffmpeg-statigo
   ```

2. Point `go.mod` at the local submodule:

   ```sh
   go mod edit -require=github.com/linuxmatters/ffmpeg-statigo@v0.0.0
   go mod edit -replace github.com/linuxmatters/ffmpeg-statigo=./third_party/ffmpeg-statigo
   ```

3. Download the static library for your host platform:

   ```sh
   cd third_party/ffmpeg-statigo && go run ./cmd/download-lib
   ```

4. Configure Git for fast-forward pulls so `git pull --rebase` does not break submodule commits:

   ```sh
   git config pull.ff only && git config submodule.recurse true
   ```

5. Build your project:

   ```sh
   CGO_ENABLED=1 go build ./...
   ```

### Checksum verification and GitHub tokens

`cmd/download-lib` reads the pinned `LibRelease` tag in `lib/fetch.go`, fetches the release metadata from GitHub, and verifies the downloaded archive against its SHA-256 digest (falling back to `SHA256SUMS`). If the checksum is missing or mismatches, `cmd/download-lib` deletes the extracted files and exits with an error. There is no override.

Set `GITHUB_TOKEN` or `GH_TOKEN` in CI or shared networks so GitHub API requests do not hit anonymous rate limits.

## Consumer justfile

Wrap submodule checks, version updates, and library downloads in a `justfile` so every developer runs one setup command.

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
        rm -f third_party/ffmpeg-statigo/lib/*/libffmpeg.a third_party/ffmpeg-statigo/lib/embedded/*/libffmpeg.a

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

## GitHub Actions CI matrix

Use native runners per platform so CGO links against the host C toolchain without cross-compilation setup:

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
          go-version-file: go.mod

      - name: Download FFmpeg libraries
        env:
          GITHUB_TOKEN: ${{ github.token }}
        run: cd third_party/ffmpeg-statigo && go run ./cmd/download-lib

      - name: Build
        run: go build -v ./...

      - name: Test
        run: go test -v ./...
```

## Cross-compilation with `cmd/download-lib`

Set `GOOS` and `GOARCH` before running `cmd/download-lib` to fetch a static library for a different target platform:

```sh
cd third_party/ffmpeg-statigo
GOOS=darwin GOARCH=arm64 go run ./cmd/download-lib
```

`cmd/download-lib` accepts six `GOOS`/`GOARCH` pairs: `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`, `windows/amd64`, and `windows/386`. Remember that Go compilation of the root `ffmpeg` package enforces 64-bit targets via `arch_guard.go`.

When you cross-compile a CGO binary on a single host, you also need a matching C cross-compiler in `CC` and target system libraries. Prefer native CI runners when available.

## Embedded variant (`-tags=embedded`)

When an application does not need AV1, HEVC, hardware acceleration, or GPL codecs, select the `embedded` profile. It replaces `x264` and `x265` with OpenH264 and removes GPL build flags.

Pass `-tags=embedded` to both the library download and the Go build:

```sh
cd third_party/ffmpeg-statigo
GOFLAGS=-tags=embedded go run ./cmd/download-lib
cd ../..
GOFLAGS=-tags=embedded go build ./...
```

The tagged download fetches `ffmpeg-embedded-<os>-<arch>.tar.gz` and installs the static archive at `lib/embedded/<os>_<arch>/libffmpeg.a`. The untagged default build continues to use `ffmpeg-<os>-<arch>.tar.gz` at `lib/<os>_<arch>/libffmpeg.a`, so both profiles coexist in the same checkout.

If your consumer `justfile` targets the `embedded` profile:

1. Check `third_party/ffmpeg-statigo/lib/embedded/$(go env GOOS)_$(go env GOARCH)/libffmpeg.a` in `_check-submodule`.
2. Remove `third_party/ffmpeg-statigo/lib/embedded/*/libffmpeg.a` when bumping the submodule tag.
3. Pass `GOFLAGS=-tags=embedded` to both `go run ./cmd/download-lib` and `go build`.

### Embedded licensing notes

- OpenH264 uses a two-clause BSD licence and requires its copyright notice.
- A self-built static OpenH264 archive does not receive the AVC patent royalty coverage that Cisco provides for its separately downloaded binary ([OpenH264 FAQ](https://www.openh264.org/faq.html)).
- FFmpeg's LGPLv2.1 obligations still apply. Static distribution requires a way for recipients to relink with a modified FFmpeg library ([FFmpeg licence](https://www.ffmpeg.org/legal.html)).

See [Codec and hardware matrix](Codec-and-Hardware-Matrix) for the exact list of codecs excluded from the `embedded` profile.

## Reference projects

These open-source tools use this submodule and `justfile` pattern in production:

- [Jive Encoder](https://github.com/linuxmatters/jive-encoder)
- [Jive Visualiser](https://github.com/linuxmatters/jive-visualiser)
- [Jive Vocals](https://github.com/linuxmatters/jive-vocals)

## Next steps

- Read [High-level `av` package](High-Level-av-Package) to build pipelines with managed `io.Closer` types.
- Check [Troubleshooting and FAQ](Troubleshooting-and-FAQ) if CGO or linker errors occur.
