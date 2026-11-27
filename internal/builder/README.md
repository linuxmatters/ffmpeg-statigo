# FFmpeg Builder

`internal/builder` downloads SHA-256 pinned third-party sources, builds FFmpeg and its dependencies incrementally, and combines static archives into `lib/<os>_<arch>/libffmpeg.a` (or `lib/embedded/<os>_<arch>/libffmpeg.a` with `--embedded`). With `--programs`, it also builds, strips, signs on macOS, and stages `ffmpeg` and `ffprobe` alongside `LICENSE` and `NOTICE.txt` into `bin/[embedded/]<os>_<arch>/`.

For full CLI flag documentation, environment variables, `--programs` staging rules, and digest refresh steps, see [Building from source](https://github.com/linuxmatters/ffmpeg-statigo/wiki/Building-From-Source) on the GitHub Wiki.

## File layout

| File | Responsibility |
| --- | --- |
| `main.go` | CLI entry point, build orchestration, `--clean`, `--list`, `--update-digests`, and archive merging (`combineMac`, `combineWindows`, `combineLinux`) |
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

## Quick commands

Run all commands from the repository root:

```sh
# Build default (GPLv3) static library into lib/<os>_<arch>/libffmpeg.a
go run ./internal/builder

# Build default static library and standalone CLI binaries into bin/<os>_<arch>/
go run ./internal/builder --programs

# Build embedded (LGPLv2.1) static library and CLI binaries into lib/embedded/ and bin/embedded/
go run ./internal/builder --embedded --programs

# List libraries for a profile
go run ./internal/builder --list
go run ./internal/builder --embedded --list

# Rebuild FFmpeg only and re-stage CLI binaries
go run ./internal/builder --programs ffmpeg

# Clean FFmpeg build outputs, or all profile outputs when no library is named
go run ./internal/builder --clean ffmpeg
go run ./internal/builder --clean

# Print SHA-256 manifest entries for enabled archives on the current platform
go run ./internal/builder --update-digests
```
