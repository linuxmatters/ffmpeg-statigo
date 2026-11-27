# Binding generator and IR goldens

`internal/generator/` parses the committed FFmpeg C headers under `include/` in pure Go and writes the five `*.gen.go` binding files in the repository root.

## Hermetic pure-Go parser

The generator uses [`modernc.org/cc/v4`](https://pkg.go.dev/modernc.org/cc/v4) as its C frontend. It needs no libclang, no C compiler, and no cgo, and it compiles under `CGO_ENABLED=0`.

Run the generator from any host with:

```sh
just generate
```

Pass `-v` to `go run ./internal/generator -v` to print the backend version and include-path trace.

Every input file comes from two repository-owned sources:

1. Committed FFmpeg headers under `include/`.
2. Eleven stub C library headers embedded from `internal/generator/sysinclude/`.

`ccconfig.go` configures `cc.Config` without running the host compiler (`cc -dM -E -`) or reading `/usr/include`. `ccPredefined` defines a fixed `gnu11` LP64 macro set directly in Go, and `newCCConfig(goos, goarch)` selects only the `cc.NewABI` table.

### Why system headers are stubbed

FFmpeg headers change their macro surface based on the host C library. Two real bugs motivated the hermetic stubs:

- `libavutil/mathematics.h` wraps 26 `M_*` constants in `#ifndef`. glibc defines all 26 under `_GNU_SOURCE`, while the Apple SDK defines only the 13 without the `f` suffix. Parsing against host headers emitted 13 extra constants on macOS.
- glibc's `aarch64` header `bits/math-vector.h` uses vector-math prototypes that `cc/v4` cannot parse, which broke translation on `linux/arm64`.

The eleven embedded stub headers in `internal/generator/sysinclude/` are `errno.h`, `inttypes.h`, `limits.h`, `math.h`, `stdarg.h`, `stddef.h`, `stdint.h`, `stdio.h`, `stdlib.h`, `string.h`, and `time.h`. `sysIncludeStubs` serves them from the sentinel search root `ffmpeg-statigo-sysinclude` through `cc.Config.FS` and returns `fs.ErrNotExist` for all other paths so `cc/v4` falls back to `include/` on disk.

When a future FFmpeg release includes a C library header that the stubs do not cover, the parse fails and names the missing header. Add a minimal stub under `internal/generator/sysinclude/` instead of reading host system headers.

## HeaderParser boundary

Parsing sits behind the `HeaderParser` interface in `internal/generator/headerparser.go`:

```go
type HeaderParser interface {
    Parse(skips *SkipCollector) *Module
    Version() string
}
```

`ccParser` in `ccparser.go` is the only implementation. Imports of `modernc.org/cc/v4` stay confined to four files (`ccconfig.go`, `ccparser.go`, `type.go`, and `ctypename.go`), while `comments.go` handles comment extraction without importing `cc/v4`, so a future parser backend requires no emitter changes:

| File | Role |
| --- | --- |
| `ccconfig.go` | Builds the hermetic `cc.Config`, predefined macros, and embedded `sysinclude/` filesystem |
| `ccparser.go` | Walks C declarations to determine symbol existence, names, and insertion order |
| `type.go` | Translates `cc/v4` types into internal type trees |
| `ctypename.go` | Formats `Field.CTypeName` strings |
| `comments.go` | Builds the comment table and attaches doc comments to declarations |

After `parser.Parse(skips)` returns `*Module`, `applyManualFixups` in `main.go` applies two pinned corrections before `Gen(m, skips)` emits Go code:

1. Sets `m.structs["AVRational"].ByValue = true` so `AVRational` passes by value rather than as a pointer aggregate.
2. Clears the stray doc comment on `m.enums["AVOptionType"]`.

If an FFmpeg header update removes or renames either target symbol, `applyManualFixups` returns an explicit error instead of panicking or emitting wrong bindings.

## IR goldens and `-dump-ir`

A byte diff across the five `*.gen.go` files spans tens of thousands of lines and hides the root cause of a parser change. Three committed golden files under `internal/generator/testdata/ir/` snapshot the intermediate representation (IR) directly.

The `-dump-ir` flag is additive on a normal generator run. It still writes all five `*.gen.go` files and refreshes all three golden files:

```sh
go run ./internal/generator -dump-ir
```

Inside the Nix shell, `dump_test.go` quotes the exact regeneration command `nix develop -c go run ./internal/generator -dump-ir`.

| Golden file | Snapshot point and contents | What a diff shows |
| --- | --- | --- |
| `structure.txt` | Parsed `*Module` at the `HeaderParser` boundary before `applyManualFixups`: all five insertion-order slices verbatim, then per symbol the expanded type tree, `CTypeName`, `BitWidth`, `Variadic`, `Typedefd`, `ByValue`, parameter names, and comment hash | The parsed declaration form changed |
| `comments.txt` | Post-`processComment` text keyed by symbol name | Comment extraction or formatting changed |
| `skips.txt` | Sorted, deduplicated `(Symbol, Reason)` pairs after `Gen` under a header line with marker, symbol, and pair counts (without working-directory-dependent `Manual` tags) | A skip decision or verbatim skip-reason string changed |

The IR dump format is an internal test artefact, not a public interface.

### Verification gates

Three checks guard generator determinism and correctness:

1. `TestIRGoldensMatchFreshRun` (`internal/generator/dump_test.go`) reparses `include/` in a temporary directory and fails on any byte difference against `testdata/ir/`, printing the first differing line and the regeneration command.
2. `TestGeneratorIRIsIdenticalAcrossTargets` (`internal/generator/crosstarget_test.go`) runs Linux and Darwin ABI targets from a single host so target-specific drift fails locally.
3. `.github/workflows/go-test.yml` runs `go run ./internal/generator` and `git diff --quiet -- ./*.gen.go` on `linux/amd64`, `linux/arm64`, `darwin/amd64`, and `darwin/arm64`. Any drift fails CI.

After any generator or header change, run `just generate`, check `git diff --stat -- '*.gen.go'`, refresh goldens with `-dump-ir` when the IR changed intentionally, and commit the updated goldens in the same commit.

## Skip tracking and `skipCeiling`

The generator skips C symbols that it cannot express safely in Go: variadic functions, fixed-size array parameters, anonymous structs, unions, and function-pointer fields. Every skip records a `(Symbol, Reason)` marker into `SkipCollector`, and the reason text appears verbatim in generated file comments.

At the end of each run, `printSkipSummary` prints the marker count and sorted symbol list to `stderr`. When a skipped symbol has a matching hand-written wrapper in the root package, the summary annotates it with `(manual binding: <Name>)`.

`const skipCeiling = 203` in `internal/generator/main.go` caps the total skip marker count. `enforceSkipCeiling` fails the generator whenever `skips.Total() > skipCeiling`. Because it is an upper bound rather than exact equality, it catches accidental regressions (such as a broken allowlist dropping dozens of symbols) without breaking when a new manual wrapper removes a skip.

Raise `skipCeiling` only as a deliberate curation decision alongside an FFmpeg header upgrade that introduces new unemittable C declarations. Never raise `skipCeiling` to silence a parser regression.

## Three-tier source layout and hand-written topic files

The root `ffmpeg` package separates generated code, core runtime glue, and curated wrappers into three tiers.

### Tier 1: generated bindings

Never hand-edit these five files. Regenerate them with `just generate`.

| File | Contents |
| --- | --- |
| `constants.gen.go` | FFmpeg `#define` constants |
| `enums.gen.go` | C enum types and values |
| `structs.gen.go` | Struct wrappers and field accessors |
| `functions.gen.go` | Function wrappers |
| `callbacks.gen.go` | Callback typedefs |

### Tier 2: core foundation

| File(s) | Purpose |
| --- | --- |
| `ffmpeg.go` | CGO `CFLAGS`, `AVError`, `WrapErr`, `CStr`, and common type aliases |
| `ffmpeg_link_default.go`, `ffmpeg_link_embedded.go` | Per-platform CGO `LDFLAGS` for the default (`!embedded`) and `embedded` build tags |
| `array.go` | Generic `Array[T]` wrapper, typed `To*Array` constructors, and element-size constants |
| `arch_guard.go` | Compile-time 64-bit pointer-size guard (`amd64` and `arm64`) |
| `doc.go` | Package documentation |

### Tier 3: hand-written topic wrappers

Group new manual wrappers by FFmpeg topic rather than adding to a single catch-all file. Pair C trampoline shims (`*.c`) with their matching Go file (`*.go`). Follow root package conventions: wrap FFmpeg error codes with `WrapErr`, manage C strings with `CStr` and `.Free()`, and prefix FFmpeg types with `AV*`.

| File(s) | Reason for manual binding | Contents |
| --- | --- | --- |
| `iterate.go` | `void **opaque` iterator idiom | Codec, muxer, demuxer, parser, filter, BSF, protocol, and `AVChannelLayoutStandard` iterators |
| `uuid.go` | Fixed `[16]uint8` array parameters | `AVUUID` type, parse, format, and compare helpers |
| `display.go` | Fixed `int32[9]` matrix parameters | `AVDisplayMatrix`, `av_display_*` rotation and flip helpers, and `av_exif_*` orientation wrappers |
| `streamgroup.go` | Anonymous C struct | `AVStreamGroupTileGridOffset` accessors |
| `opt.go` | Generic Go slice to C binary buffer | `AVOptSetSlice` |
| `image.go` | Fixed `[4]`-plane C arrays | Hybrid: six manual `av_image_*` plane wrappers (six unwrapped siblings stay in `functions.gen.go`) |
| `samples.go` | Variable-length `uint8_t **` plane arrays | `av_samples_*` audio sample-plane wrappers |
| `audio_fifo.go` | `void * const *data` plane pointers | `av_audio_fifo_*` write, read, peek, and `peek_at` wrappers (reuses `samplePointerArray`) |
| `swscale.go` | Fixed-plane C array parameters | `sws_*` scaling and pixel-format conversion |
| `swresample.go` | Plane-pointer parameters | `swr_*` audio resampling |
| `get_format.go` + `get_format.c` | C function-pointer callback | `AVCodecContext.get_format` `//export` trampoline for pixel-format selection |
| `avio.go` + `avio.c` | Function-pointer callbacks | Custom-I/O `AVIOContext` (`AVIOAllocContext`, `AVBufferCreate`) via `runtime/cgo.Handle` |
| `log.go` + `log.c` | `av_log` callback bridge | `//export` bridge forwarding FFmpeg logs to Go `slog` |
| `tx.go` + `tx.c` | Invoking C function pointer from Go | `AVTxCall` C-side invoker for `av_tx_fn` returned by `av_tx_init` |
| `log_format.go` | C variadic functions | `AVLog`, `AVAsprintf`, `AVStrlcatf`, and `AVBprintf` formatted in Go and passed through a fixed `"%s"` shim |
| `fields.go` | Array and pointer-array struct fields | Quant matrices (`intra_matrix`, `inter_matrix`, `chroma_intra_matrix`) and `AVFrame.extended_data` |
| `helpers.go` | Cross-cutting utilities | `AVRational.String`, `ToAVHWFramesContext`, `AVRescaleDelta`, `AVSizeMult` |
| `parsers.go` | Primitive out-parameter classification | `av_ac3_parse_header`, `av_adts_header_parse`, `av_vorbis_parse_frame_flags` |
| `parseutils.go` | Single out-parameter classification | `av_parse_ratio`, `av_parse_video_rate`, `av_codec_get_tag2` |

The optional high-level `av/` package (`Input`, `Decoder`, `Encoder`, `FilterGraph`, `Output`, `HWDevice`) is also hand-written and sits above the root package.

Next: read [Release Process and Versioning](Release-Process-and-Versioning) to ship updated headers, libraries, and Go module releases.
