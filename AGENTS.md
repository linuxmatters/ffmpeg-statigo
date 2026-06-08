# AGENTS.md

## Setup commands

- Enter development shell: `nix develop` (or let direnv activate automatically)
- Download FFmpeg libraries: `go run ./cmd/download-lib`
- Initialise submodules: `git submodule update --init --recursive`

## Build and test commands

- **Full build:** `just build` — builds FFmpeg from source, regenerates bindings, compiles all
- **Build FFmpeg only:** `just build-static ffmpeg --clean`
- **Build static libraries:** `just build-static` (uses current GOOS/GOARCH)
- **Regenerate bindings:** `just generate` or `go run ./internal/generator`
- **Build examples:** `just build-examples`
- **Run tests:** `just test`
- **Download libraries:** `go run ./cmd/download-lib`

## Code style

- **Auto-generated files:** Never edit `*.gen.go` files (constants, enums, structs, functions, callbacks) — regenerate with `just generate`
- **C string handling:** Use `CStr` type with `.Free()` for cleanup
- **Error handling:** Wrap FFmpeg return codes with `WrapErr()` function
- **Stream processing:** Check `AVErrorEOF` and `EAgain` in processing loops
- **Type naming:** All FFmpeg types prefixed with `AV*` (e.g., `AVCodecContext`, `AVFrame`)
- **CGO required:** All builds need `CGO_ENABLED=1`

## Testing instructions

- Run `just test` before committing
- Tests require downloaded libraries (`go run ./cmd/download-lib` first)
- See `ffmpeg_test.go` for version validation patterns

## PR/commit guidelines

- **Submodule workflow:** Configure git for fast-forward pulls only: `git config pull.ff only && git config submodule.recurse true`
- **Version schemes:** Two distinct versions — library releases (`lib-X.Y.Z.N`) and module releases (`vX.Y.Z.N`)
- **Release tags:** Library releases use `lib-` prefix; Go module releases use `v` prefix

## Development workflow

- **Never run `go build` directly** — always use `just build` for proper CGO flags and build sequencing
- **Cross-compilation:** Set `GOOS` and `GOARCH` before downloading: `GOOS=darwin GOARCH=arm64 go run ./cmd/download-lib`
- **Platform-specific builds:** Justfile auto-detects current platform, outputs to `lib/<os>_<arch>/`
- **Binding regeneration:** Required after FFmpeg header changes — run `just generate`
- **Nix-only regeneration:** Run `just generate` / `go run ./internal/generator` only inside `nix develop` (libclang 20.1.8, gcc 15.2.0), where system include discovery via `gcc -E -v` is guaranteed. Pass `-v` for toolchain and include-path trace output.

## Key architecture

- **Core:** `ffmpeg.go` contains CGO directives, platform linker flags, and base types (`AVError`/`WrapErr`, `CStr`); `array.go` holds the generic `Array[T]` type and its typed constructors; `arch_guard.go` enforces 64-bit-only at compile time
- **Generated bindings:** `*.gen.go` files in root directory — constants, enums, struct wrappers, function wrappers, callback typedefs; emitted by `internal/generator/` from FFmpeg headers; never hand-edit
- **Hand-written bindings:** topic files in the root package for symbols the generator skips (variadics, fixed-size array params, anonymous structs, function-pointer bridges); each skip is recorded with a reason, the skip summary notes when a skipped symbol has a hand-written binding, and the total is capped by `skipCeiling` in `internal/generator/main.go`
  - `iterate.go` — registry iterators (codec/muxer/demuxer/parser/filter/bsf) + protocol enumeration + `AVChannelLayoutStandard` standard channel-layout iterator
  - `uuid.go` — `AVUUID` type; `[16]uint8` array params CGO can't pass directly
  - `display.go` — `av_display_*` display-matrix + `av_exif_*` orientation wrappers; `int32[9]` matrix params CGO can't pass directly (`AVDisplayMatrix` type)
  - `streamgroup.go` — `AVStreamGroupTileGridOffset` accessors for anonymous C struct
  - `opt.go` — `AVOptSetSlice`; Go-slice → C binary option setter
  - `image.go` — `av_image_*` plane/linesize wrappers
  - `samples.go` — `av_samples_*` audio sample-plane wrappers
  - `audio_fifo.go` — `av_audio_fifo_*` data-path wrappers (write/read/peek/peek_at); `void * const *data` plane-pointer params CGO can't pass directly (reuses `samplePointerArray`)
  - `swscale.go` — `sws_*` software scaling / pixel-format conversion
  - `swresample.go` — `swr_*` audio resampling
  - `get_format.go` + `get_format.c` — `AVCodecContext.get_format` callback bridge (cgo `//export` trampoline) for selecting a decode pixel format, e.g. a hardware format
  - `avio.go` + `avio.c` — custom-I/O `AVIOContext` via `runtime/cgo.Handle` callback bridge
  - `log.go` + `log.c` — `av_log` callback bridge to Go/`slog` via cgo `//export`
  - `tx.go` + `tx.c` — `AVTxCall` forward-call invoker for the `av_tx_fn` pointer `av_tx_init` returns; CGO can't call a C function pointer from Go, so the shim invokes it C-side
  - `log_format.go` — variadic-format shims (`AVLog`, `AVAsprintf`, etc.); CGO can't call C varargs, so these format on the Go side and pass through a fixed `"%s"` C shim
  - `fields.go` — struct-field accessors the generator can't express (quant matrices, `AVFrame.extended_data`, pixel-format descriptor components, etc.)
  - `helpers.go` — small cross-cutting helpers (`AVRational.String`, `ToAVHWFramesContext`, `AVRescaleDelta`, `AVSizeMult`)
  - `parsers.go` — `av_ac3_parse_header` / `av_adts_header_parse` / `av_vorbis_parse_frame_flags`; primitive out-param parsers the generator can't classify as in/out
  - `parseutils.go` — `av_parse_ratio` / `av_parse_video_rate` / `av_codec_get_tag2`; out-param parse and codec-tag-lookup helpers
- **Headers:** `include/` contains FFmpeg C headers
- **Libraries:** `lib/<os>_<arch>/` contains platform-specific static libraries (gitignored)
- **Builder:** `internal/builder/` compiles FFmpeg + 20 dependencies from source
- **Generator:** `internal/generator/` parses headers using libclang, outputs Go bindings
- **Downloader:** `cmd/download-lib/` fetches pre-built libraries from GitHub Releases
- **Pipeline layer:** `av/` optional high-level layer over the root bindings — owned `io.Closer` resource wrappers (Input/Decoder/Encoder/FilterGraph/Output/HWDevice); not generated

## Hardware acceleration

Supported: NVENC/NVDEC (Linux), QuickSync (Linux), VideoToolbox (macOS), Vulkan Video (cross-platform). See `README.md` and `docs/CODECS.md` for codec matrix.

## Security considerations

- **GPL licensing:** Combined work inherits GPL requirements from FFmpeg/x264/x265
- **Static libraries gitignored:** Only submodule reference committed, not ~100MB binaries
- **Library distribution:** Use GitHub Releases for pre-built binaries, not git
- **Pinned build dependencies (CWE-494):** Every source archive the builder downloads is SHA-256 pinned in `internal/builder/digests.go`, keyed by download URL. The pin is verified on the final archive bytes for both a fresh download and a cache hit, so a poisoned download cache cannot be trusted. A mismatch deletes the archive and aborts; a missing pin (e.g. after a version or URL bump) fails closed with an actionable error rather than skipping verification. Bootstrap or refresh pins in a trusted environment with `go run ./internal/builder --update-digests`, verify the printed digests against upstream-published checksums, then commit `digests.go`. Pins seeded by `--update-digests` are trust-on-first-use; prefer upstream-published checksums where they exist.
