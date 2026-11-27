# FFmpeg Statigo

Static FFmpeg 8.1.2 bindings for Go and standalone `ffmpeg-statigo` / `ffprobe-statigo` CLI binaries across Linux, macOS, and Windows.

FFmpeg Statigo bundles static FFmpeg libraries directly into your Go binary so you ship zero runtime FFmpeg installs. It also publishes standalone `ffmpeg-statigo` and `ffprobe-statigo` executables built from the same curated static toolchain.

> [!NOTE]
> FFmpeg Statigo is a hard fork of [csnewman/ffmpeg-go](https://github.com/csnewman/ffmpeg-go), updated with FFmpeg 8.1.2, Go 1.26, hardware acceleration, and ~90% FFmpeg API coverage.

## What is included

Every release ships two build profiles for Linux (`amd64`, `arm64`), macOS (`amd64`, `arm64`), and Windows (`amd64`, `386`).

| Area | `default` (GPLv3) | `embedded` (LGPLv2.1) |
| --- | --- | --- |
| Video | H.264 (`x264`), H.265 (`x265`), AV1 (`rav1e`, `dav1d`), VP8/VP9 (`libvpx`), WebP (`libwebp`) | H.264 (`OpenH264`), VP8/VP9 (`libvpx`), WebP (`libwebp`), built-in decoders (no AV1) |
| Audio | Opus, MP3 (`mp3lame`), AAC, FLAC, Vorbis, PCM | Opus, MP3 (`mp3lame`), AAC, FLAC, Vorbis, PCM |
| Hardware acceleration | NVENC/NVDEC, VA-API, QuickSync (Linux), VideoToolbox (macOS), Vulkan Video | None (software only) |

See the [Codec and hardware matrix](https://github.com/linuxmatters/ffmpeg-statigo/wiki/Codec-and-Hardware-Matrix) for library versions, enabled codecs, containers, filters, and GPU generation support.

## Install: Go library

Static libraries (~100 MB per platform) do not distribute via `go get`. Add FFmpeg Statigo as a git submodule with a `replace` directive:

1. Add the submodule: `git submodule add https://github.com/linuxmatters/ffmpeg-statigo third_party/ffmpeg-statigo`
2. Add a replace directive to `go.mod`: `replace github.com/linuxmatters/ffmpeg-statigo => ./third_party/ffmpeg-statigo`
3. Download the pre-built static library (verified against release `SHA256SUMS`): `cd third_party/ffmpeg-statigo && go run ./cmd/download-lib`
4. Configure git for fast-forward submodule pulls: `git config pull.ff only && git config submodule.recurse true`
5. Build your module: `go build ./...`

To use the LGPLv2.1 `embedded` profile instead of the default GPLv3 profile, pass `-tags=embedded` for both the library download and your Go build:

```sh
cd third_party/ffmpeg-statigo && GOFLAGS=-tags=embedded go run ./cmd/download-lib
cd ../.. && GOFLAGS=-tags=embedded go build ./...
```

### Quick start

The root package exposes raw C pointers that require manual `Free`, `Unref`, and `Close` calls. For automatic `io.Closer` cleanup, use the optional [`av`](av) package:

```go
package main

import (
	"fmt"
	"log"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
	"github.com/linuxmatters/ffmpeg-statigo/av"
)

func main() {
	in, err := av.Open("input.mp4")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	stream, err := in.BestStream(ffmpeg.AVMediaTypeVideo)
	if err != nil {
		log.Fatal(err)
	}

	dec, err := av.NewDecoder(stream)
	if err != nil {
		log.Fatal(err)
	}
	defer dec.Close()

	fmt.Printf("opened %d streams (video codec %d)\n", in.NbStreams(), stream.Codecpar().CodecId())
}
```

Browse runnable programs in [`examples/`](examples/), read the [High-level `av` package guide](https://github.com/linuxmatters/ffmpeg-statigo/wiki/High-Level-av-Package), or see [Go consumer integration](https://github.com/linuxmatters/ffmpeg-statigo/wiki/Go-Consumer-Integration) for CI setup, cross-compilation, and memory rules.

## Install: standalone CLI binaries

| | Platform | Repository |
| :-: | --- | --- |
| <img src="https://apt.wimpys.world/logo.svg" width="32" height="32" alt=""> | Debian and Ubuntu | [apt.wimpys.world](https://apt.wimpys.world/#ffmpeg-statigo) |
| <img src="https://rpm.wimpys.world/logo.svg" width="32" height="32" alt=""> | Fedora and openSUSE | [rpm.wimpys.world](https://rpm.wimpys.world/#ffmpeg-statigo) |
| <img src="https://pkg.wimpys.world/logo.svg" width="32" height="32" alt=""> | Arch Linux | [pkg.wimpys.world](https://pkg.wimpys.world/#ffmpeg-statigo) |
| <img src="https://brew.wimpys.world/logo.svg" width="32" height="32" alt=""> | macOS (Homebrew) | [brew.wimpys.world](https://brew.wimpys.world/#ffmpeg-statigo) |
| <img src="https://scoop.wimpys.world/logo.svg" width="32" height="32" alt=""> | Windows (Scoop) | [scoop.wimpys.world](https://scoop.wimpys.world/#ffmpeg-statigo) |
| <img src="https://nix.wimpys.world/logo.svg" width="32" height="32" alt=""> | Nix and NixOS | [nix.wimpys.world](https://nix.wimpys.world/#ffmpeg-statigo) |

APT, RPM, Arch, and Nix repositories are signed, and Homebrew and Scoop verify packages over HTTPS with SHA-256 checksums. Each page gives the setup steps and keeps ffmpeg-statigo up to date.

The Linux packages and binaries need glibc 2.39 or later, for example Ubuntu 24.04, Debian 13, or Fedora 40.

Packages install the `default` (GPLv3) profile as `ffmpeg-statigo` and `ffprobe-statigo` to coexist with a system `ffmpeg`, and Homebrew and Scoop also provide `ffprobe-statigo`.

### Any platform

[bin](https://github.com/marcosnils/bin) installs and updates the release binary:

```shell
bin install github.com/linuxmatters/ffmpeg-statigo
```

For direct `lib-X.Y.Z.N` release archives (`default` and `embedded`), SHA-256 checksum verification, and macOS Gatekeeper quarantine removal, see [Standalone CLI downloads](https://github.com/linuxmatters/ffmpeg-statigo/wiki/Standalone-CLI-Downloads).

## Built with FFmpeg Statigo

Tools built on these bindings for the [Linux Matters](https://linuxmatters.sh) podcast toolchain:

- **[Jive Encoder](https://github.com/linuxmatters/jive-encoder)** - converts podcast `.wav` recordings into tagged MP3 files with cover art.
- **[Jive Visualiser](https://github.com/linuxmatters/jive-visualiser)** - renders real-time audio frequency visualisations from `.wav` inputs to MP4 video.
- **[Jive Vocals](https://github.com/linuxmatters/jive-vocals)** - processes raw microphone recordings into broadcast-ready audio in one command.

If you ship a tool with FFmpeg Statigo, open a pull request to add your project to this list.

## Documentation

Full guides and reference pages live on the [FFmpeg Statigo GitHub Wiki](https://github.com/linuxmatters/ffmpeg-statigo/wiki).

### Users

- [Standalone CLI downloads](https://github.com/linuxmatters/ffmpeg-statigo/wiki/Standalone-CLI-Downloads) - direct `lib-X.Y.Z.N` archives, checksums, and macOS quarantine removal
- [Codec and hardware matrix](https://github.com/linuxmatters/ffmpeg-statigo/wiki/Codec-and-Hardware-Matrix) - bundled library versions, codecs, filters, and GPU support
- [High-level `av` package](https://github.com/linuxmatters/ffmpeg-statigo/wiki/High-Level-av-Package) - owned `io.Closer` pipeline wrappers and `Raw()` access
- [Go consumer integration](https://github.com/linuxmatters/ffmpeg-statigo/wiki/Go-Consumer-Integration) - submodules, CI/CD, cross-compilation, and CGO memory rules
- [Troubleshooting and FAQ](https://github.com/linuxmatters/ffmpeg-statigo/wiki/Troubleshooting-and-FAQ) - common build, linking, and runtime issues

### Contributors

- [Architecture and bindings](https://github.com/linuxmatters/ffmpeg-statigo/wiki/Architecture-and-Bindings) - three-tier binding layout and API coverage
- [Building from source](https://github.com/linuxmatters/ffmpeg-statigo/wiki/Building-From-Source) - `internal/builder`, profiles, flags, and pinned source digests
- [Binding generator and IR goldens](https://github.com/linuxmatters/ffmpeg-statigo/wiki/Binding-Generator-and-IR-Goldens) - hermetic header parsing and IR golden tests
- [Release process and versioning](https://github.com/linuxmatters/ffmpeg-statigo/wiki/Release-Process-and-Versioning) - `lib-X.Y.Z.N` and `vX.Y.Z.N` release workflows

## Licensing

The Go binding code is MIT licensed. The `default` static library and CLI archives include GPL-licensed `x264` and `x265`, so projects that link or redistribute them must meet GPLv3 requirements.

The `embedded` variant omits GPL build flags and uses BSD-licensed OpenH264, but FFmpeg's LGPLv2.1 obligations still apply. Static distribution requires a way for recipients to relink with a modified FFmpeg library, and the self-built OpenH264 archive does **not** qualify for Cisco's binary AVC patent royalty coverage.

Standalone CLI archives bundle their effective `LICENSE` (`COPYING.GPLv3` for `default`, `COPYING.LGPLv2.1` for `embedded`) and `NOTICE.txt` rather than the repository root MIT licence.
