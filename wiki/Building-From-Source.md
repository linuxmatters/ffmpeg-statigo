# Building from source

Build the static `libffmpeg.a` archive and optional standalone `ffmpeg` and `ffprobe` binaries locally from pinned upstream sources.

## Development shell

Enter the Nix development shell before you run any build, test, or lint recipe:

```sh
nix develop
```

If you use `direnv`, run `direnv allow` after you inspect `.envrc`. Nix flakes exclude untracked files from evaluation. When you add or modify files under `nix/`, stage those files with `git add` and run `direnv reload` or start a new `nix develop` shell.

The flake in `flake.nix` provides Go 1.26, GCC, Clang tooling, `autoconf`, `automake`, `libtool`, `cmake`, `meson`, `ninja`, `nasm`, `yasm`, `gperf`, `pkg-config`, `python3`, `just`, `nfpm`, linters, and the Rust toolchain (`rustc`, `cargo`, `cargo-c` for `rav1e`). Platform shell hooks configure environment variables automatically:

- **macOS**: Sets `CGO_CFLAGS` to the active SDK sysroot, `LIBCXX_INCLUDE` to the Nix C++ headers for `zimg`, `MACOSX_DEPLOYMENT_TARGET=13.0`, and `LLVM_OBJCOPY` to strip Rust bitcode without shadowing Apple `ar` or `strip`.
- **Linux**: Prepends `/run/opengl-driver/lib`, `vulkan-loader`, `intel-media-driver`, and `vpl-gpu-rt` to `LD_LIBRARY_PATH`, sets `ONEVPL_SEARCH_PATH`, and detects `LIBVA_DRIVER_NAME` (`iHD`, `radeonsi`, or `nvidia`) from `lspci`.

## Just recipes

Run all build and quality commands through `just`. Tailor manages `build`, `test`, and `lint` in `just/tailor.just` with `CGO_ENABLED=1` exported. Do not run `go build` directly.

Before you run `just build`, `just test`, or `just lint`, provide compatible static archives with `just download-lib`, `just build-static`, or `just build-ffmpeg`. Managed quality commands never download archives or regenerate bindings automatically.

| Recipe | Purpose |
| --- | --- |
| `just build-ffmpeg` | Cleans and rebuilds FFmpeg, builds static archives, regenerates Go bindings, compiles all packages and examples, and runs `./bin/introspect` |
| `just build` | Compiles Go packages with CGO enabled and writes executables to `bin/` against existing archives |
| `just build-static [args]` | Runs `go run ./internal/builder` for the default (GPLv3) profile into `lib/<os>_<arch>/libffmpeg.a` |
| `just build-embedded [args]` | Runs `go run ./internal/builder --embedded` for the embedded (LGPLv2.1) profile into `lib/embedded/<os>_<arch>/libffmpeg.a` |
| `just test` | Runs `go vet ./...` and `go test -cover -covermode=atomic ./...` (some tests access real GPU hardware) |
| `just lint` or `just lint check` | Runs `actionlint`, `go mod tidy -diff`, Nix `vendorHash` checks, `golangci-lint`, `govulncheck`, `gocyclo`, and `ineffassign` without modifying files |
| `just lint correct` | Runs `go mod tidy`, formats Go source with `.golangci.yml`, updates `vendorHash` if present, and runs all lint checks |
| `just generate` | Runs `go run ./internal/generator` to regenerate the five `*.gen.go` binding files |
| `just download-lib` | Runs `go run ./cmd/download-lib/` to fetch pre-built static archives from GitHub Releases |
| `just clean` | Removes `.build/` and `.build/embedded/` build, source, and staging directories plus binaries in `bin/` while keeping cached downloads and `lib/` archives |

To build or test the embedded profile after `just build-embedded`, pass `GOFLAGS=-tags=embedded`:

```sh
GOFLAGS=-tags=embedded just build
GOFLAGS=-tags=embedded just test
```

## Builder file layout

The from-source builder lives in `internal/builder/`. Each file owns a single stage of archive fetching, compilation, packaging, or verification:

| File | Responsibility |
| --- | --- |
| `main.go` | CLI entry point, build orchestration, `--clean`, `--list`, `--update-digests`, and platform archive merging (`combineMac`, `combineWindows`, `combineLinux`) |
| `library.go` | `Library` struct, incremental `Build` loop, configuration hashing, build environment setup, and `stripAndSignBinary` |
| `libraries.go` | Third-party library definitions, dependency order, platform filters, Windows overrides, `scrubProgramPkgConfig`, and `pruneStagingDylibs` |
| `profile.go` | Output paths, `default` and `embedded` profile selection, FFmpeg configure wrapper, and licence pre-staging |
| `programs.go` | CLI flag parsing (`parseBuilderArgs`), `LICENSE` staging, `NOTICE.txt` rendering, and `stagePrograms` |
| `ffmpeg_args.go` | Curated FFmpeg codec, format, filter, hardware acceleration, and `--enable-ffmpeg` / `--enable-ffprobe` configure flags |
| `buildsystems.go` | `AutoconfBuild`, `CMakeBuild`, `MesonBuild`, `CargoBuild`, `OpenSSLBuild`, and `MakefileBuild` implementations |
| `digests.go` | Pinned SHA-256 source archive manifest (`archiveDigests`) and fail-closed `verifyDigest` |
| `download.go` | Retryable `grab/v3` downloader and safe archive extraction |
| `state.go` | Per-library `.build-state` JSON tracking for incremental skips |
| `openh264.go` | OpenH264 library definition for the `embedded` profile |
| `coff_archive.go` | Windows archive filter that removes zero-section COFF and split DWARF (`.dwo`) members |
| `pointer.go` | `Disabled()` helper for disabled library entries such as `vvenc` |

## Builder CLI flags

Run the builder from the repository root. Output paths derive automatically from `runtime.GOOS`, `runtime.GOARCH`, and `--embedded`, so the CLI takes no output path argument.

```sh
# Build the default static library
go run ./internal/builder

# Build the default static library and standalone CLI binaries
go run ./internal/builder --programs

# Build the embedded (LGPLv2.1) static library and CLI binaries
go run ./internal/builder --embedded --programs

# Rebuild only FFmpeg and re-stage CLI binaries
go run ./internal/builder --programs ffmpeg

# Clean only the FFmpeg build, source, and staged outputs
go run ./internal/builder --programs --clean ffmpeg

# List libraries for a profile
go run ./internal/builder --list
go run ./internal/builder --embedded --list

# Compute SHA-256 manifest entries for enabled archives on this host
go run ./internal/builder --update-digests
```

| Flag or argument | Behaviour |
| --- | --- |
| `--embedded` | Selects the LGPLv2.1 OpenH264 profile. Sets the build root to `.build/embedded` (`.build/embedded/windows_386_msvcrt` on `windows/386`), the static archive to `lib/embedded/<os>_<arch>/libffmpeg.a`, and the CLI directory to `bin/embedded/<os>_<arch>/`. |
| `--programs` | Builds `ffmpeg` and `ffprobe` (`ffmpeg.exe` and `ffprobe.exe` on Windows) alongside static libraries, then stages stripped binaries, `LICENSE`, and `NOTICE.txt` into `bin/[embedded/]<os>_<arch>/`. |
| `--clean` | Deletes `<buildRoot>/build/<lib>`, `<buildRoot>/src/<lib>`, staged `.a` files in `<buildRoot>/staging/lib/`, and staged executables in `<buildRoot>/staging/bin/` for the selected libraries, then exits. Cached downloads stay in `<buildRoot>/downloads/`. |
| `--list` | Prints a table of libraries, target platforms, build systems, and link archives for the selected profile, then exits. |
| `--update-digests` | Downloads enabled archives for the current platform and prints sorted `archiveDigests` entries for `digests.go`. |
| `<library> ...` | Restricts a build or `--clean` to the named libraries. Unknown or profile-excluded names return an error. `combineLibraries` runs only on full builds with no positional filter, and `stagePrograms` runs on full builds or when `ffmpeg` is named with `--programs`. |

Default profile builds use `.build` as the build root (`.build/windows_386_msvcrt` on `windows/386`). Each library stores incremental skip state in `<buildRoot>/build/<lib>/.build-state` (recording URL, configuration hash, and staged outputs) and writes build logs to `<buildRoot>/build/<lib>/build.log`.

## Environment variables

| Variable | Default | Purpose |
| --- | --- | --- |
| `FFMPEG_RELEASE_TAG` | `GITHUB_REF_NAME` | Tag written as `Release: <tag>` in `NOTICE.txt`. Overrides `GITHUB_REF_NAME` when both are set, for example during `workflow_dispatch` runs. |
| `MACOSX_DEPLOYMENT_TARGET` | `13.0` | Minimum macOS version exported in the build environment on `darwin` when unset or empty. |
| `LLVM_OBJCOPY` | `llvm-objcopy` | Tool that `combineMac` runs on macOS to strip `__LLVM,__bitcode` and `__LLVM,__cmdline` sections from `libffmpeg.a`. |
| `AR` | `ar` | Archive tool that `combineWindows` runs on Windows for `ar -M` merging and `ar s` indexing. |
| `STRIP` | `strip` | Strip tool that `combineWindows` (`--strip-unneeded`) and `stripAndSignBinary` run across platforms. |
| `CODESIGN` | `codesign` | Signing tool that `stripAndSignBinary` runs on macOS for ad-hoc binary signing. |

## Per-OS static runtime linking

When you pass `--programs`, `ffmpegArgsCommonWithPrograms` replaces `--disable-programs` with `--enable-ffmpeg --enable-ffprobe --disable-ffplay`. Each operating system applies specific steps so the resulting CLI binaries carry no non-system shared library dependencies:

- **Linux**: FFmpeg's `PostExtract` hook runs `scrubProgramPkgConfig`, which strips exact `-lgcc_s` tokens from `Libs` and `Libs.private` lines in staged `.pc` files. Configure receives `--extra-ldexeflags=-static-libstdc++ -static-libgcc`. After staging, `stripAndSignBinary` runs `strip --strip-all`. GPU hardware drivers load via `dlopen` at runtime.
- **Windows**: FFmpeg's `PostExtract` hook also runs `scrubProgramPkgConfig` to drop dynamic `-lgcc_s` flags. Configure receives `--extra-ldexeflags=-static -static-libgcc -static-libstdc++` so binaries import only built-in Windows DLLs (`ucrtbase.dll` on `amd64`, `msvcrt.dll` on `386`). On `windows/386`, configure adds `--extra-cflags=-mpreferred-stack-boundary=4` and `--extra-cflags=-mstackrealign` for `swscale` SIMD alignment, sets `dav1d` to `-Dstack_alignment=4`, and disables `x265` assembly (`-DENABLE_ASSEMBLY=OFF`). `combineWindows` filters zero-section COFF and `.dwo` members via `coff_archive.go`, and `stripAndSignBinary` runs `strip --strip-all`.
- **macOS**: `libiconv` builds with `--disable-shared`, and FFmpeg's `PostExtract` hook runs `pruneStagingDylibs` to remove any stale `*.dylib` files or symlinks from `<buildRoot>/staging/lib/`. Executables link only against `/usr/lib/*` and `/System/Library/Frameworks/*` with `MACOSX_DEPLOYMENT_TARGET` set to `13.0`. Before `stagePrograms` writes `ffmpeg` or `ffprobe`, it deletes any existing target file so macOS allocates a fresh inode and avoids stale Mach-O vnode signature caches. `stripAndSignBinary` then runs `strip` followed by `codesign --force --sign -` for ad-hoc signing.

For every `--programs` build, `stageLicense` writes `COPYING.GPLv3` (`default`) or `COPYING.LGPLv2.1` (`embedded`) to `bin/[embedded/]<os>_<arch>/LICENSE` from the FFmpeg source tree, staging share directory, or verified download archive. It never copies the repository root MIT `LICENSE`. `renderNotice` writes `bin/[embedded/]<os>_<arch>/NOTICE.txt` with the profile, target, FFmpeg version, licence terms, relinking instructions for LGPLv2.1, and every active bundled library's SPDX licence, upstream URL, and pinned SHA-256 digest.

## Pinned source digests (CWE-494)

Every third-party source archive URL in `internal/builder/libraries.go` and `internal/builder/openh264.go` is SHA-256 pinned in `internal/builder/digests.go` (`archiveDigests`). This design prevents unverified archive execution (CWE-494).

`DownloadFile` invokes `verifyDigest` fail-closed on both fresh downloads and cached files in `<buildRoot>/downloads/`:

1. A missing pin aborts the build with an error naming the unpinned URL.
2. A digest mismatch deletes the archive file and aborts the build so a re-run fetches fresh bytes rather than trusting a corrupted cache.

When you bump a library version or URL, refresh pins in a trusted environment:

```sh
go run ./internal/builder --update-digests
go run ./internal/builder --embedded --update-digests
```

Check each printed SHA-256 digest against the upstream project's official checksums before you paste it into `internal/builder/digests.go` and commit. Platform-gated libraries such as `libiconv` (macOS) or `libdrm` and `libva` (Linux) download only on their target operating system, so run `--update-digests` on each target platform whose pins changed.

Next: read [Binding Generator and IR Goldens](Binding-Generator-and-IR-Goldens) to regenerate Go bindings after an FFmpeg header update.
