# Troubleshooting and FAQ

Solutions for common build, runtime, platform, and licensing issues in `ffmpeg-statigo`.

## Build and environment issues

### `cannot find -lffmpeg`

The pre-built static archive `lib/<os>_<arch>/libffmpeg.a` (or `lib/embedded/<os>_<arch>/libffmpeg.a`) is absent because static libraries are gitignored.

Run the downloader from the `ffmpeg-statigo` directory:

```sh
# Default (GPLv3) profile
go run ./cmd/download-lib

# Embedded (LGPLv2.1) profile
GOFLAGS=-tags=embedded go run ./cmd/download-lib
```

If GitHub API rate limits block `cmd/download-lib`, export `GITHUB_TOKEN` or `GH_TOKEN` before you run it. The downloader rejects any archive with a missing or mismatched SHA-256 checksum and has no override.

### CGO and GPU driver errors on NixOS

NixOS does not populate `/usr/lib` or standard global headers. Enter the flake shell so GCC, pkg-config, and GPU runtime paths are configured:

```sh
nix develop
```

If you use `direnv` and it blocks `.envrc`, inspect `.envrc` and run `direnv allow`. When you add or change files under `nix/`, stage them with `git add` and run `direnv reload`, because Nix flakes ignore untracked files.

On Linux, the dev shell sets `LD_LIBRARY_PATH`, `ONEVPL_SEARCH_PATH`, `LIBVA_DRIVERS_PATH`, and `LIBVA_DRIVER_NAME` for Vulkan, VA-API, and oneVPL hardware tests. Hardware tests fail when run outside a shell that exposes those runtime drivers.

### Go module cache errors

Running `go get github.com/linuxmatters/ffmpeg-statigo` without a `replace` directive pulls the Go sources into the module cache without the ~100MB static library archive, causing linker failures.

Use `ffmpeg-statigo` as a git submodule with a `replace` directive in `go.mod`, then clear any stale cached copy:

```sh
go mod edit -replace github.com/linuxmatters/ffmpeg-statigo=./third_party/ffmpeg-statigo
go clean -modcache
```

### `CGO_ENABLED=1` required

Every build and test of the root package and `av/` links against `libffmpeg.a` through cgo and requires `CGO_ENABLED=1`.

Use `just build` and `just test`, because `just/tailor.just` exports `CGO_ENABLED=1` automatically. If you invoke `go test` or `go build` manually in a consumer project, set `CGO_ENABLED=1` explicitly. The root package and `av/` require cgo, while pure-Go packages such as `internal/generator`, `lib/`, `cmd/download-lib`, and `internal/builder` also build with `CGO_ENABLED=0`.

## Platform-specific issues

### macOS Gatekeeper quarantine (`xattr`)

Standalone macOS CLI binaries (`ffmpeg` and `ffprobe`) are ad-hoc signed with `codesign --force --sign -`, not Apple Developer ID notarised. When you download a macOS CLI archive in a browser, Gatekeeper blocks execution until you clear the quarantine attribute on the extracted directory:

```sh
xattr -dr com.apple.quarantine <extracted-dir>
```

Replace `<extracted-dir>` with the unpacked directory name, for example `ffmpeg-bin-darwin-arm64`.

Two other macOS notes:

- Linker warnings about a macOS version mismatch are expected on macOS versions newer than 13.0, because libraries target `MACOSX_DEPLOYMENT_TARGET=13.0`.
- When building static libraries from source on macOS, ensure `LLVM_OBJCOPY` points to an `llvm-objcopy` binary so `combineMac` can strip Rust `__LLVM,__bitcode` and `__LLVM,__cmdline` sections from `libffmpeg.a`.

### Windows clean `PATH` and `386` stack alignment

Build and test on Windows inside the MSYS2 subsystem that matches `GOARCH`:

| Target | MSYS2 subsystem | Runtime DLL | Compiler triplet |
| --- | --- | --- | --- |
| `windows/amd64` | `UCRT64` (`/ucrt64/bin`) | `ucrtbase.dll` and `api-ms-win-crt-*.dll` | `x86_64-w64-mingw32` |
| `windows/386` | `MINGW32` (`/mingw32/bin`) | `msvcrt.dll` | `i686-w64-mingw32` |

Put the subsystem binary directory (`$MINGW_BIN`) ahead of other toolchains on `PATH` when building. Never mix `UCRT64` and `MINGW32` compilers or archives.

**Clean-`PATH` verification**: Standalone Windows CLI executables link statically with `-static -static-libgcc -static-libstdc++` and must run on machines without MSYS2 installed. Test them in PowerShell with `PATH` restricted so MSYS2 DLLs cannot mask an accidental dynamic dependency:

```powershell
$env:PATH = 'C:\Windows\System32;C:\Windows'
.\ffmpeg.exe -version
```

**Windows `386` stack alignment**: GCC 16 on 32-bit Windows can enter C functions with a 4-byte aligned stack, while FFmpeg's `swscale` SIMD code requires 16-byte stack alignment and crashes on misaligned access. `internal/builder` prevents this crash automatically:

- Passes `--extra-cflags=-mpreferred-stack-boundary=4` and `--extra-cflags=-mstackrealign` when configuring FFmpeg on `windows/386`.
- Builds `dav1d` with `-Dstack_alignment=4` on `windows/386`.
- Disables `x265` assembly (`-DENABLE_ASSEMBLY=OFF`) on `windows/386`.
- Filters zero-section COFF and `.dwo` members out of the merged archive via `internal/builder/coff_archive.go`.

If you customise FFmpeg C flags on `windows/386`, keep both `-mpreferred-stack-boundary=4` and `-mstackrealign` intact.

## Licensing FAQ

### Which licence applies to my project?

Three licences appear in this repository. Which one governs your distributed binary depends on what you ship:

| Component | Licence | Key obligations |
| --- | --- | --- |
| Go binding source code (`*.go`, `av/`, tooling) | MIT (`LICENSE` in repository root) | Preserve copyright and permission notices |
| `default` static library and CLI binaries (`x264`, `x265`, full codecs, hardware acceleration) | GPLv3 (`COPYING.GPLv3`) | Combined work inherits GPLv3, including Section 6 source-code availability |
| `embedded` static library and CLI binaries (software-only OpenH264, no GPL flags) | LGPLv2.1 (`COPYING.LGPLv2.1`) plus two-clause BSD for OpenH264 | Preserve copyright notices and satisfy LGPLv2.1 Section 6 relinking for static binaries |

### Does the `embedded` profile remove all copyleft obligations?

No. The `embedded` profile omits GPL build flags, `x264`, `x265`, `rav1e`, `dav1d`, and hardware acceleration, but FFmpeg itself remains under LGPLv2.1. If you distribute a statically linked binary built with `-tags=embedded`, LGPLv2.1 Section 6 requires you to give recipients a way to relink your application with a modified FFmpeg library.

### Does the self-built OpenH264 archive get Cisco's AVC patent coverage?

No. Cisco's AVC patent royalty coverage applies only to Cisco's separately downloaded binary. A static OpenH264 library compiled from source inside `lib/embedded/<os>_<arch>/libffmpeg.a` does not qualify for that coverage. Review the [OpenH264 FAQ](https://www.openh264.org/faq.html) and [FFmpeg legal guidance](https://www.ffmpeg.org/legal.html) before commercial distribution.

### Why do CLI archives carry a different `LICENSE` file from the repository root?

The repository root `LICENSE` is MIT and covers only the Go binding code. Each standalone CLI archive bundles compiled FFmpeg and its linked C/C++/Rust libraries, so `internal/builder/programs.go` stages FFmpeg's `COPYING.GPLv3` (`default`) or `COPYING.LGPLv2.1` (`embedded`) as `LICENSE` alongside `NOTICE.txt`.

## General FAQ

### Why does installation use a git submodule instead of `go get`?

Each pre-built `libffmpeg.a` archive is roughly 100MB per platform. Go's module proxy and git history are unsuited to storing gigabytes of binary archives across releases, so static archives ship on GitHub Releases (`lib-X.Y.Z.N`) and `cmd/download-lib` fetches only the archive for your target platform.

### Why are `drawtext`, `subtitles`, and `libass` excluded?

`ffmpeg-statigo` targets Go applications that draw text and graphics in Go (for example with `image/draw` and `golang.org/x/image/font`) and pass frames to FFmpeg for encoding and muxing. Omitting `freetype`, `harfbuzz`, `fontconfig`, and `libass` keeps the static archive smaller. FFmpeg can still copy, extract, and mux subtitle streams without `libass`.

### Why do Linux CLI packages require `glibc >= 2.39`?

Linux binaries build on `ubuntu-24.04` and `ubuntu-24.04-arm` runners, which ship glibc 2.39. `_build-libs.yml` verifies that no ELF symbol exceeds `GLIBC_2.39`, and `packaging/nfpm.yaml` encodes that floor for Debian, Ubuntu, Fedora, openSUSE, and Arch Linux packages. Supported distributions include Ubuntu 24.04, Debian 13, and Fedora 40 or later.
