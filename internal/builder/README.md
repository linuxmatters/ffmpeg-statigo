# FFmpeg Builder

`internal/builder` downloads pinned third-party source archives, builds FFmpeg and its dependencies incrementally, and combines the static archives into `lib/<os>_<arch>/libffmpeg.a` (or `lib/embedded/<os>_<arch>/libffmpeg.a` with `--embedded`).
With `--programs`, it also builds, strips, signs on macOS, and stages standalone `ffmpeg` and `ffprobe` executables alongside `LICENSE` and `NOTICE.txt` into `bin/[embedded/]<os>_<arch>/`.

## File layout

| File | Responsibility |
| --- | --- |
| `main.go` | CLI entry point, build orchestration, `--clean`, `--list`, `--update-digests`, and platform archive merging (`combineMac`, `combineWindows`, `combineLinux`) |
| `library.go` | `Library` definition, incremental `Build` loop, configuration hashing, build environment setup, and `stripAndSignBinary` |
| `libraries.go` | Third-party library definitions, dependency order, platform filters, Windows overrides, `scrubProgramPkgConfig`, and `pruneStagingDylibs` |
| `profile.go` | Output paths, `default` and `embedded` profile selection, FFmpeg configure wrapper, and licence pre-staging |
| `programs.go` | CLI flag parsing (`parseBuilderArgs`), `LICENSE` staging, `NOTICE.txt` rendering, and `stagePrograms` |
| `ffmpeg_args.go` | Curated FFmpeg codec, format, filter, hardware acceleration, and `--enable-ffmpeg`/`--enable-ffprobe` configure flags |
| `buildsystems.go` | `AutoconfBuild`, `CMakeBuild`, `MesonBuild`, `CargoBuild`, and `MakefileBuild` implementations |
| `digests.go` | Pinned SHA-256 source archive manifest (`archiveDigests`) and fail-closed `verifyDigest` |
| `download.go` | Retryable `grab/v3` downloader and safe archive extraction |
| `state.go` | Per-library `.build-state` JSON tracking for incremental skips |
| `openh264.go` | OpenH264 library definition for the `embedded` profile |
| `coff_archive.go` | Windows archive filter that removes zero-section COFF and split DWARF (`.dwo`) members |
| `pointer.go` | `Disabled()` helper for disabled library entries such as `vvenc` |

## Usage and CLI flags

Run the builder from the repository root. Output paths derive automatically from `runtime.GOOS`, `runtime.GOARCH`, and `--embedded`, so the CLI takes no output path argument.

```sh
# Build default profile static library into lib/<os>_<arch>/libffmpeg.a
go run ./internal/builder

# Build default profile static library and standalone CLI binaries into bin/<os>_<arch>/
go run ./internal/builder --programs

# Build embedded LGPLv2.1 profile (OpenH264, no GPL flags) into lib/embedded/ and bin/embedded/
go run ./internal/builder --embedded --programs

# List libraries for a profile
go run ./internal/builder --list
go run ./internal/builder --embedded --list

# Rebuild only FFmpeg and re-stage CLI binaries
go run ./internal/builder --programs ffmpeg

# Clean build, source, and staged outputs for FFmpeg (or all profile libraries when none are named)
go run ./internal/builder --programs --clean ffmpeg
go run ./internal/builder --clean

# Print SHA-256 manifest entries for enabled archives on the current platform
go run ./internal/builder --update-digests
```

| Flag or argument | Behaviour |
| --- | --- |
| `--embedded` | Selects the LGPLv2.1/OpenH264 profile. Sets the build root to `.build/embedded` (`.build/embedded/windows_386_msvcrt` on `windows/386`), the static library path to `lib/embedded/<os>_<arch>/libffmpeg.a`, and the CLI directory to `bin/embedded/<os>_<arch>/`. |
| `--programs` | Builds `ffmpeg` and `ffprobe` (`ffmpeg.exe` and `ffprobe.exe` on Windows) alongside the static libraries, then stages the stripped binaries, `LICENSE`, and `NOTICE.txt` into `bin/[embedded/]<os>_<arch>/`. |
| `--clean` | Deletes `<buildRoot>/build/<lib>`, `<buildRoot>/src/<lib>`, staged `.a` files in `<buildRoot>/staging/lib/`, and staged executables in `<buildRoot>/staging/bin/` for the selected libraries (or all profile libraries), then exits. Cached downloads remain in `<buildRoot>/downloads/`. |
| `--list` | Prints a table of libraries, platforms, build systems, and link archives for the selected profile, then exits. |
| `--update-digests` | Downloads enabled archives for the current platform and prints sorted `archiveDigests` entries for review and inclusion in `digests.go`. |
| `<library> ...` | Positional library names restrict the build or `--clean` to those libraries. Unknown or profile-excluded names return an error. `combineLibraries` runs only on full builds with no positional filter, while `stagePrograms` runs on full builds or when `ffmpeg` is selected with `--programs`. |

Default profile builds use `.build` as the build root (`.build/windows_386_msvcrt` on `windows/386`).
Each library writes its incremental skip state to `<buildRoot>/build/<lib>/.build-state` (tracking URL, configuration hash, and required staged outputs) and streams command output to `<buildRoot>/build/<lib>/build.log`.

## Environment variables

| Variable | Default | Purpose |
| --- | --- | --- |
| `FFMPEG_RELEASE_TAG` | `GITHUB_REF_NAME` | Release tag written as `Release: <tag>` in `NOTICE.txt`. Overrides `GITHUB_REF_NAME` when both are set (for example during `workflow_dispatch` runs). |
| `MACOSX_DEPLOYMENT_TARGET` | `13.0` | Minimum macOS version set in the build environment on `darwin` when unset or empty. |
| `LLVM_OBJCOPY` | `llvm-objcopy` | Tool used on macOS by `combineMac` to strip `__LLVM,__bitcode` and `__LLVM,__cmdline` sections from `libffmpeg.a`. |
| `AR` | `ar` | Archive tool used on Windows by `combineWindows` for `ar -M` merging and `ar s` indexing. |
| `STRIP` | `strip` | Strip tool used by `combineWindows` (`--strip-unneeded`) and by `stripAndSignBinary` on all platforms. |
| `CODESIGN` | `codesign` | Signing tool used on macOS by `stripAndSignBinary` for ad-hoc binary signing. |

## How `--programs` builds and stages CLI binaries

1. **FFmpeg configure flags**: `ffmpegArgsCommonWithPrograms` replaces `--disable-programs` with `--enable-ffmpeg --enable-ffprobe --disable-ffplay`.
2. **Static runtime linking**:
   - **Linux**: FFmpeg's `PostExtract` hook runs `scrubProgramPkgConfig` to remove dynamic `-lgcc_s` flags from staged `.pc` files. Configure receives `--extra-ldexeflags=-static-libstdc++ -static-libgcc`.
   - **Windows**: FFmpeg's `PostExtract` hook also runs `scrubProgramPkgConfig` to remove dynamic `-lgcc_s` flags from staged `.pc` files. Configure receives `--extra-ldexeflags=-static -static-libgcc -static-libstdc++` so the executables depend only on built-in OS DLLs.
   - **macOS**: `libiconv` builds with `--disable-shared`, and FFmpeg's `PostExtract` hook runs `pruneStagingDylibs` to remove stale `*.dylib` files from `<buildRoot>/staging/lib/`. Executables link only against `/usr/lib/*` and `/System/Library/Frameworks/*` with `MACOSX_DEPLOYMENT_TARGET` defaulting to `13.0`.
3. **Stripping and macOS signing**: `stagePrograms` copies `ffmpeg` and `ffprobe` (`ffmpeg.exe` and `ffprobe.exe` on Windows) from `<buildRoot>/staging/bin/` into `bin/[embedded/]<os>_<arch>/`. It removes any existing target binary first so macOS allocates a fresh inode and avoids stale Mach-O vnode signature caches. `stripAndSignBinary` then runs `strip --strip-all <bin>` on Linux and Windows, or `strip <bin>` followed by `codesign --force --sign - <bin>` (ad-hoc signing) on macOS.
4. **macOS Gatekeeper quarantine**: Ad-hoc signed binaries are not Apple Developer ID notarised. When a user downloads a macOS CLI archive in a browser, they must remove the quarantine attribute after extraction:

   ```sh
   xattr -dr com.apple.quarantine <extracted-dir>
   ```

5. **Licence and notice staging**:
   - `stageLicense` writes `COPYING.GPLv3` (`default`) or `COPYING.LGPLv2.1` (`embedded`) to `bin/[embedded/]<os>_<arch>/LICENSE`. It checks `<buildRoot>/src/ffmpeg/`, `<buildRoot>/staging/share/ffmpeg/`, and finally the SHA-256-verified FFmpeg download archive in `<buildRoot>/downloads/`. It never copies the repository root MIT `LICENSE`.
   - `renderNotice` writes `bin/[embedded/]<os>_<arch>/NOTICE.txt` containing the profile, target, FFmpeg version tag, effective licence, repository URL, optional `Release:` line, GPLv3 Section 6 notice or LGPLv2.1 Section 4 relinking instructions, and every active bundled library with its SPDX licence, upstream URL, and pinned SHA-256 digest.

## Pinned source digests (`digests.go`)

Every third-party source archive URL is SHA-256 pinned in `digests.go` (`archiveDigests`) for fail-closed integrity verification (CWE-494).
`DownloadFile` runs `verifyDigest` on both fresh downloads and cached archives in `<buildRoot>/downloads/`:

- A missing pin aborts the build with an error naming the unpinned URL.
- A digest mismatch deletes the downloaded archive and aborts the build so a re-run downloads fresh bytes.

After bumping a library version or URL, run `go run ./internal/builder --update-digests` (and `go run ./internal/builder --embedded --update-digests` for OpenH264) in a trusted environment.
Verify each printed digest against upstream checksums before committing the updated entries to `digests.go`.
Platform-specific libraries download only on their target operating system, so update their pins on the corresponding platform.
