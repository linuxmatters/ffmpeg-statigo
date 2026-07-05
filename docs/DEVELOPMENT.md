# ffmpeg-statigo Developer Guide

This guide covers two audiences: **consumers** integrating ffmpeg-statigo into their own projects, and **contributors** working on ffmpeg-statigo itself.

## Why Manual Setup?

Go modules don't handle 100MB+ binary files. `go get` won't fetch static libraries from the module cache. The solution: git submodule + replace directive + library download.

## Prerequisites

- Go 1.26+
- GCC (CGO compilation)
- Git

## Consumer Integration

Projects using ffmpeg-statigo follow this pattern. Reference implementations with justfiles and GitHub workflows:

- **[Jive Encoder](https://github.com/linuxmatters/jive-encoder)** 🪩 - drop your podcast `.wav` in, get a shiny MP3 out with metadata, cover art, and all.
- **[Jive Visualiser](https://github.com/linuxmatters/jive-visualiser)** 🔥 - spin a `.wav` into a groovy MP4 visualiser with real-time audio frequencies.
- **[Jive Vocals](https://github.com/linuxmatters/jive-vocals)** 🕺 - turn raw microphone recordings into broadcast-ready audio in one command. No configuration, no surprises.

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

Wrap this in a `just setup` recipe for a consistent developer experience.

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
          go-version-file: go.mod

      - name: Download FFmpeg libraries
        run: cd third_party/ffmpeg-statigo && go run ./cmd/download-lib

      - name: Build
        run: go build -v ./...

      - name: Test
        run: go test -v ./...
```

### Cross-Compilation

Set `GOOS` and `GOARCH` before downloading: `GOOS=darwin GOARCH=arm64 go run ./cmd/download-lib`

## Working on ffmpeg-statigo

1. Clone the repository.
2. `go run ./cmd/download-lib` - downloads pre-built libraries.
3. `just build` - full rebuild from source and regenerate bindings.

## Build System

All builds go through `just`. Never use `go build` directly - the justfile handles CGO flags and build sequencing.

| Command | Purpose |
|---------|---------|
| `just build` | Build static library from source, regenerate bindings, compile |
| `just test` | Run tests |
| `just generate` | Regenerate Go bindings from headers |
| `just download-lib` | Download pre-built libraries |

### Pinned build dependencies

The from-source builder downloads ~20 third-party archives (x264, x265, dav1d, openssl, zlib, libvpx, rav1e, and more), compiles them, and links them into the shipped `libffmpeg.a`. Each archive is SHA-256 pinned in `internal/builder/digests.go`, keyed by download URL, closing the supply-chain hole of trusting unverified downloads (CWE-494).

Verification fails closed on three conditions:

- The pin is checked against the final archive bytes for **both** a fresh download and a download-cache hit, so a poisoned cache cannot be silently trusted.
- A digest mismatch deletes the archive so a re-run re-downloads rather than reusing poison, then aborts the build.
- A missing pin aborts with an actionable error instead of skipping verification. This fires automatically when a library's version or URL is bumped.

To bootstrap or refresh pins (trust-on-first-use), run in a trusted environment:

```sh
go run ./internal/builder --update-digests
```

This downloads every enabled library's archive, computes its SHA-256, and prints map entries for `archiveDigests`. Review each printed digest against the upstream-published checksum (openssl, opus, gnu.org, and similar publish official values that supersede the seeded ones), then paste the entries into `internal/builder/digests.go` and commit. Platform-gated libraries (e.g. `libiconv`, `vvenc`) are only fetched on their target platform, so run the command on each platform whose archives you need to pin.

## Project Layout

The root package has three source tiers. The generator skips C symbols it cannot safely express: variadics, fixed-size array parameters, anonymous structs, unions, and function-pointer fields. Every skip is recorded with a reason string and the total is regression-capped by `skipCeiling` in `internal/generator/main.go`. The skip summary annotates each skipped symbol that has a hand-written binding with a `(manual binding: <Name>)` note, making the covered-but-skipped set directly observable. The hand-written tier covers exactly these gaps.

### Tier 1 - generated bindings (never hand-edit; regenerate with `just generate`)

| File | Contents |
|------|----------|
| `constants.gen.go` | FFmpeg `#define` constants |
| `enums.gen.go` | C enum types and values |
| `structs.gen.go` | Struct wrappers and field accessors |
| `functions.gen.go` | Function wrappers |
| `callbacks.gen.go` | Callback typedefs |

### Tier 2 - core and foundation (hand-written)

| File | Purpose |
|------|---------|
| `ffmpeg.go` | CGO directives, platform linker flags, `AVError`/`WrapErr`, `CStr`, common type aliases |
| `array.go` | Generic `Array[T]` type, typed `To*Array` constructors, element-size constants |
| `arch_guard.go` | Compile-time 64-bit-only invariant (unsigned-underflow guard) |
| `doc.go` | Package documentation |

### Tier 3 - hand-written topic wrappers

| File(s) | Why hand-written | Contents |
|---------|-----------------|----------|
| `iterate.go` | `void **opaque` iterator pattern; generator emits no Go wrapper for this idiom | Registry iterators for codecs, muxers, demuxers, parsers, filters, bsfs; protocol enumeration; `AVChannelLayoutStandard` standard channel-layout iterator |
| `uuid.go` | Fixed `[16]uint8` array params; CGO cannot pass fixed arrays directly | `AVUUID` type, parse/format/compare |
| `display.go` | Fixed `int32[9]` matrix params; CGO cannot pass fixed arrays directly | `av_display_rotation_get`, `av_display_rotation_set`, `av_display_matrix_flip`; `av_exif_matrix_to_orientation`, `av_exif_orientation_to_matrix`; `AVDisplayMatrix` type |
| `streamgroup.go` | Anonymous C struct; CGO cannot reference it by name | `AVStreamGroupTileGridOffset` accessors |
| `opt.go` | Generic Go-slice parameter; no generator analogue | `AVOptSetSlice` - Go slice → C binary option setter |
| `image.go` | Fixed `[4]`-plane C arrays; cgo pointer rules prevent direct passing | `av_image_*` plane/linesize wrappers; shared fixed-plane helper |
| `samples.go` | Variable-length `uint8_t **` plane arrays; same shape constraint as `image.go` | `av_samples_*` audio sample-plane wrappers |
| `swscale.go` | Same fixed-plane constraint as `image.go` | `sws_*` software scaling and pixel-format conversion |
| `swresample.go` | Same plane-pointer constraint as `samples.go` | `swr_*` audio resampling |
| `avio.go` + `avio.c` | Function-pointer callbacks; requires `runtime/cgo.Handle` bridge | Custom-I/O `AVIOContext` (`AVIOAllocContext`, `AVBufferCreate`); each context gets its own handle, deleted on teardown |
| `log.go` + `log.c` | `av_log` callback requires cgo `//export` | `av_log` bridge to Go/`slog` |
| `log_format.go` | CGO cannot call C variadic functions | `AVLog`, `AVAsprintf`, `AVStrlcatf`, `AVBprintf`; formats on the Go side, passes through a fixed `"%s"` C shim (also neutralises format-string injection) |
| `fields.go` | Array and pointer-array fields the generator cannot express | Quant matrices (`intra_matrix`/`inter_matrix`/`chroma_intra_matrix`), `AVFrame.extended_data`, `AVPixFmtDescriptor.comp`, `AVMasteringDisplayMetadata.display_primaries`, `AVPanScan.position` |
| `helpers.go` | Small helpers with no generator analogue | `AVRational.String`, `ToAVHWFramesContext`, `AVRescaleDelta`, `AVSizeMult` |
| `parseutils.go` | Single out-params the generator cannot classify as in/out | `av_parse_ratio`, `av_parse_video_rate`, `av_codec_get_tag2` |

### Infrastructure and supporting paths

| Path | Description |
|------|-------------|
| `include/` | FFmpeg C headers |
| `lib/<os>_<arch>/` | Platform-specific static libraries (gitignored) |
| `internal/builder/` | Builds FFmpeg + 20 dependencies from source |
| `internal/generator/` | Generates Go bindings from headers using libclang (see [Generator](#generator) below) |
| `cmd/download-lib/` | Downloads pre-built libraries from GitHub Releases |
| `av/` | Optional high-level pipeline layer - owned `io.Closer` wrappers (Input/Decoder/Encoder/FilterGraph/Output); see [docs/PIPELINE.md](PIPELINE.md) |
| `examples/` | Working examples (transcode, metadata, etc.) |

## Generator

> [!IMPORTANT]
> Run `just generate` only inside `nix develop`. The generator links libclang 20 and relies on gcc 15 for system include discovery. Running it outside the Nix shell produces incorrect or incomplete output.

`internal/generator/` parses FFmpeg headers with libclang via the cgo binding `github.com/Newbluecake/bootstrap` and emits the `*.gen.go` files. Because it links libclang, regeneration requires a pinned clang (currently 20) and a working C toolchain, which is why it is only supported inside `nix develop`. The correctness gate is byte-identical output: after any generator change, run `just generate` then `git diff --stat -- '*.gen.go'` and confirm the diff is empty.

> **TODO (future consideration):** evaluate porting the generator to [`modernc.org/cc/v4`](https://pkg.go.dev/modernc.org/cc/v4), a pure-Go C99 frontend. This would drop the libclang/cgo dependency and the clang-version pinning churn, making regeneration toolchain-independent (no Nix shell or distro clang in CI). The risk is that `cc/v4`'s parse model differs from libclang, so the generated output and the existing libclang workarounds (unnamed-struct naming, `size_t`-reported-as-`int`) would need re-validating against the byte-identical gate before switching.

## Versioning

Two distinct version schemes:

- **Library releases** (`lib-8.1.2.x`): Static library builds, distributed via GitHub Releases
- **Module releases** (`v8.1.2.x`): Go module versions

The download tool automatically fetches the latest `lib-8.1.2.x` release matching the FFmpeg version in `lib/fetch.go`.

## Troubleshooting

### "cannot find -lffmpeg"

Libraries not downloaded. Run `go run ./cmd/download-lib` in the ffmpeg-statigo directory.

### CGO errors on NixOS

Enter the dev shell: `nix develop` (or let `direnv` activate automatically).

### Module cache issues

If `go get` was attempted without the replace directive: `go clean -modcache`
