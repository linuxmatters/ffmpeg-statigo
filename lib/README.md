# FFmpeg Static Libraries

This directory holds pre-built FFmpeg static libraries (`libffmpeg.a`) for linking into Go binaries. Static archives are gitignored and must be downloaded or built from source before running `go build` or `just build`.

## Download pre-built libraries

Run the downloader explicitly from the repository root:

```sh
# Default profile (GPLv3) -> lib/<os>_<arch>/libffmpeg.a
go run ./cmd/download-lib

# Embedded profile (LGPLv2.1, OpenH264) -> lib/embedded/<os>_<arch>/libffmpeg.a
GOFLAGS=-tags=embedded go run ./cmd/download-lib
```

Set `GOOS` and `GOARCH` to download a library for another target platform, for example:

```sh
go build -o bin/download-lib ./cmd/download-lib && GOOS=darwin GOARCH=arm64 ./bin/download-lib
```

The downloader verifies every archive against the `SHA256SUMS` file published on the matching `lib-X.Y.Z.N` GitHub release. If the checksum is missing or does not match, the download fails and the extracted files are deleted.

## Supported platforms and manual archives

If you download archives manually from [GitHub Releases](https://github.com/linuxmatters/ffmpeg-statigo/releases) (look for the `lib-X.Y.Z.N` tag that matches `LibRelease` in `fetch.go`), extract each archive into `lib/`:

| Platform | Default archive (`lib/<os>_<arch>/libffmpeg.a`) | Embedded archive (`lib/embedded/<os>_<arch>/libffmpeg.a`) |
| --- | --- | --- |
| `linux/amd64` | `ffmpeg-linux-amd64.tar.gz` | `ffmpeg-embedded-linux-amd64.tar.gz` |
| `linux/arm64` | `ffmpeg-linux-arm64.tar.gz` | `ffmpeg-embedded-linux-arm64.tar.gz` |
| `darwin/amd64` | `ffmpeg-darwin-amd64.tar.gz` | `ffmpeg-embedded-darwin-amd64.tar.gz` |
| `darwin/arm64` | `ffmpeg-darwin-arm64.tar.gz` | `ffmpeg-embedded-darwin-arm64.tar.gz` |
| `windows/amd64` | `ffmpeg-windows-amd64.tar.gz` | `ffmpeg-embedded-windows-amd64.tar.gz` |
| `windows/386` | `ffmpeg-windows-386.tar.gz` | `ffmpeg-embedded-windows-386.tar.gz` |

## Directory layout

```text
lib/
├── <os>_<arch>/
│   └── libffmpeg.a          # Default (GPLv3) static archive (gitignored)
├── embedded/
│   └── <os>_<arch>/
│       └── libffmpeg.a      # Embedded (LGPLv2.1) static archive (gitignored)
├── fetch.go                 # Release downloader and SHA-256 verification
├── profile_default.go       # Build-tag profile selection (!embedded)
├── profile_embedded.go      # Build-tag profile selection (embedded)
└── README.md
```

For consumer setup and CI caching, see [Go consumer integration](https://github.com/linuxmatters/ffmpeg-statigo/wiki/Go-Consumer-Integration). To compile static archives from source, see [Building from source](https://github.com/linuxmatters/ffmpeg-statigo/wiki/Building-From-Source).
