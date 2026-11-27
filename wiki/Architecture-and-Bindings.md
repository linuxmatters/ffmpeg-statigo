# Architecture and bindings

Repository architecture, raw C memory ownership rules, and FFmpeg API coverage for `ffmpeg-statigo`.

- [Three-tier repository layout](#three-tier-repository-layout)
- [Raw C memory ownership and error contract](#raw-c-memory-ownership-and-error-contract)
- [FFmpeg API coverage](#ffmpeg-api-coverage)

## Three-tier repository layout

`ffmpeg-statigo` separates raw C bindings, optional Go-managed pipelines, and build/generation tooling into three distinct layers.

### 1. Root package (`github.com/linuxmatters/ffmpeg-statigo`)

The root package exposes the FFmpeg C API directly over CGO (`CGO_ENABLED=1`) and is organised into three internal source tiers.

#### Tier 1 - Generated bindings

Never edit `*.gen.go` files by hand. Regenerate them from `include/` with `just generate`.

| File | Contents |
| --- | --- |
| `constants.gen.go` | FFmpeg `#define` constants |
| `enums.gen.go` | C enum types and values |
| `structs.gen.go` | Struct wrappers and field accessors |
| `functions.gen.go` | Exported function wrappers (`AV*`, `Swr*`, `Sws*`) |
| `callbacks.gen.go` | Callback typedefs |

#### Tier 2 - Core and foundation

| File | Purpose |
| --- | --- |
| `ffmpeg.go` | CGO `CFLAGS`, `AVError` and `WrapErr`, `EAgain`, `AVErrorEOF`, `CStr`, common typedef aliases |
| `ffmpeg_link_default.go`, `ffmpeg_link_embedded.go` | Platform `LDFLAGS` selected by the `!embedded` and `embedded` build tags |
| `array.go` | Generic unchecked `Array[T]` accessor, typed `To*Array` constructors, element-size constants |
| `arch_guard.go` | Compile-time 64-bit pointer-width guard (`//go:build amd64 \|\| arm64`) |
| `doc.go` | Root package documentation |

#### Tier 3 - Hand-written topic wrappers

The generator skips C symbols that CGO or the emitter cannot express directly: variadic functions, fixed-size array parameters, anonymous structs, unions, out-only primitive pointers, and C function-pointer callbacks. Each skip records a reason string, and `skipCeiling` in `internal/generator/main.go` caps the total. Hand-written topic files in the root package cover those gaps.

| File(s) | Why hand-written | Contents |
| --- | --- | --- |
| `iterate.go` | `void **opaque` iterator idiom has no generated wrapper | Registry iterators for codecs, muxers, demuxers, parsers, filters, and bitstream filters, plus protocol enumeration and `AVChannelLayoutStandard` |
| `uuid.go` | CGO cannot pass fixed `[16]uint8` arrays directly | `AVUUID` type with parse, format, and compare helpers |
| `display.go` | CGO cannot pass fixed `int32[9]` matrices directly | `AVDisplayMatrix`, `av_display_*` rotation/flip wrappers, and `av_exif_*` orientation wrappers |
| `streamgroup.go` | Anonymous C struct cannot be named in CGO | `AVStreamGroupTileGridOffset` accessors |
| `opt.go` | Generic Go-slice binary option helper | `AVOptSetSlice` |
| `image.go` | Fixed `[4]`-plane C arrays conflict with CGO pointer rules | Hybrid: 6 manual `av_image_*` plane/linesize wrappers and helpers (6 unwrapped siblings live in `functions.gen.go`) |
| `samples.go` | Variable-length `uint8_t **` sample-plane arrays | `av_samples_*` wrappers |
| `audio_fifo.go` | `void * const *data` plane-pointer parameters | `av_audio_fifo_*` write, read, peek, and `peek_at` wrappers (reuses `samplePointerArray`) |
| `swscale.go` | Fixed plane-pointer arrays | `sws_*` software scaling and pixel-format conversion |
| `swresample.go` | Sample plane-pointer arrays | `swr_*` audio resampling |
| `get_format.go` + `get_format.c` | Function-pointer callback trampoline (`//export`) | `AVCodecContext.get_format` callback bridge for pixel-format negotiation (including hardware formats) |
| `avio.go` + `avio.c` | Function-pointer callbacks via `runtime/cgo.Handle` | Custom-I/O `AVIOContext` (`AVIOAllocContext`, `AVBufferCreate`) with per-context handle cleanup |
| `log.go` + `log.c` | `av_log` callback via CGO `//export` | Bridge from FFmpeg logging to Go `log/slog` |
| `tx.go` + `tx.c` | CGO cannot invoke a C function pointer directly | `AVTxCall` C-side invoker for the `av_tx_fn` pointer returned by `av_tx_init` |
| `log_format.go` | CGO cannot call C variadic functions | `AVLog`, `AVAsprintf`, `AVStrlcatf`, `AVBprintf` (formats in Go and passes through a fixed `"%s"` C shim to prevent format-string injection) |
| `fields.go` | Array and pointer-array struct fields | Quantisation matrices (`intra_matrix`, `inter_matrix`, `chroma_intra_matrix`) and `AVFrame.extended_data` |
| `helpers.go` | Small cross-cutting helpers with no generator equivalent | `AVRational.String`, `ToAVHWFramesContext`, `AVRescaleDelta`, `AVSizeMult` |
| `parsers.go` | Primitive out-parameters the generator cannot classify | `av_ac3_parse_header`, `av_adts_header_parse`, `av_vorbis_parse_frame_flags` |
| `parseutils.go` | Single out-parameters the generator cannot classify | `av_parse_ratio`, `av_parse_video_rate`, `av_codec_get_tag2` |

### 2. High-level pipeline layer (`av/`)

The optional `av` package (`github.com/linuxmatters/ffmpeg-statigo/av`) adds owned Go resource management over the root bindings. It is hand-written and opt-in.

- **Six owned `io.Closer` types**: `Input` (demux `AVFormatContext`), `Decoder` (`AVCodecContext` + scratch `AVFrame`), `FilterGraph` (`AVFilterGraph` + single buffersrc/buffersink), `Encoder` (`AVCodecContext` + scratch `AVPacket`), `Output` (mux `AVFormatContext` + IO), and `HWDevice` (hardware `AVBufferRef` + `AVHWDeviceType`).
- **Idempotent teardown**: `Close` releases underlying FFmpeg handles in dependency order, and a second `Close` is a safe no-op.
- **Callback buffer lifetime**: Frames and packets passed to `Decode`, `Encode`, and `Pull` callbacks are valid only during that callback, because the wrapper unrefs its scratch buffer on return. Copy or reference (`ffmpeg.AVFrameRef`) any data you keep after the callback.
- **Raw escape hatch**: Every owned type exposes `Raw()` to return its `*ffmpeg.AV*` pointer so callers can mix high-level types with direct root-package calls. Do not free handles obtained from `Raw()` yourself.

### 3. Build, generation, and distribution tooling

| Path | Responsibility |
| --- | --- |
| `internal/builder/` | Builds FFmpeg and its static dependencies from source in topological order (`AllLibraries`) for `default` and `embedded` profiles across Linux, macOS, and Windows. Verifies every archive against SHA-256 pins in `internal/builder/digests.go` on both fresh downloads and cache hits, failing closed on mismatch or missing pin. |
| `internal/generator/` | Pure-Go hermetic header parser and binding emitter built on `modernc.org/cc/v4` (`CGO_ENABLED=0`, no libclang or host C compiler). Reads only `include/` and the 11 stub headers embedded from `internal/generator/sysinclude/`. |
| `cmd/download-lib/` and `lib/fetch.go` | Downloads pre-built `libffmpeg.a` archives from the pinned `LibRelease` GitHub Release and verifies each archive against `SHA256SUMS` before keeping extracted files. |
| `include/`, `lib/` | Committed FFmpeg C headers (`include/`) and gitignored static archives (`lib/<os>_<arch>/`, `lib/embedded/<os>_<arch>/`). |

#### Generator architecture and IR goldens

Header parsing sits behind the `HeaderParser` interface in `internal/generator/headerparser.go`. Its sole implementation is `ccParser` (`ccparser.go`), with `cc/v4` imports isolated to `ccparser.go`, `type.go`, `ctypename.go`, and `ccconfig.go` (`comments.go` holds comment extraction).

Correctness is enforced by three checks:

1. **Cross-target determinism**: `TestGeneratorIRIsIdenticalAcrossTargets` parses all four primary targets (`linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`) from any host and verifies identical output. CI verifies `git diff --stat -- '*.gen.go'` is empty after `just generate`.
2. **Skip regression ceiling**: `skipCeiling` in `internal/generator/main.go` (currently `203`) fails the generator if total skip markers increase unexpectedly.
3. **IR goldens (`internal/generator/testdata/ir/`)**: Running `go run ./internal/generator -dump-ir` emits the five `*.gen.go` files and updates three snapshots checked by `TestIRGoldensMatchFreshRun`:
   - `structure.txt`: parsed `Module` at the `HeaderParser` boundary before manual fixups.
   - `comments.txt`: post-`processComment` documentation text keyed by symbol.
   - `skips.txt`: sorted, deduplicated `(Symbol, Reason)` pairs with marker, symbol, and pair counts.

## Raw C memory ownership and error contract

> [!WARNING]
> The root package bindings are thin wrappers around FFmpeg's C API. They have no Go finalisers and no automatic memory cleanup.

### Strings (`CStr`)

`CStr` wraps a `*C.char` pointer in C memory (`ffmpeg.go`). Ownership depends on how the `CStr` is created:

| Constructor or source | Allocator | Caller calls `Free()`? | Rules |
| --- | --- | --- | --- |
| `ToCStr(val)` | `C.CString` | Yes | Frees via C `free` and clears the internal pointer. |
| `AllocCStr(n)` | `C.calloc` | Yes | Zero-initialises a buffer of `n` bytes; frees via C `free`. |
| `GlobalCStr(val)` | `C.CString` (interned) | No (`dontFree = true`) | Process-lifetime cache for compile-time constants such as codec or option names. Never pass dynamic or unbounded input. Call `Dup()` (`AVStrdup`) if FFmpeg takes ownership and frees the string. |
| FFmpeg-allocated return (`wrapAVCStr`) | `av_malloc` / `av_asprintf` | Yes | `Free()` routes through `C.av_free`. |
| Static or struct-owned return (`wrapStaticCStr`) | FFmpeg static or struct memory | No (`dontFree = true`) | `Free()` is a safe no-op. |

### Core FFmpeg object lifecycles

Every `AV*` struct in the root package wraps a raw C pointer. Follow FFmpeg's C ownership contract for allocation, reuse, and release:

| Type | Allocate or open | Per-iteration reset | Release |
| --- | --- | --- | --- |
| `AVFrame` | `AVFrameAlloc()` | `AVFrameUnref(frame)` | `AVFrameFree(&frame)` |
| `AVPacket` | `AVPacketAlloc()` | `AVPacketUnref(pkt)` | `AVPacketFree(&pkt)` |
| `AVCodecContext` | `AVCodecAllocContext3(codec)` + `AVCodecOpen2(...)` | `AVCodecFlushBuffers(ctx)` (on seek) | `AVCodecFreeContext(&ctx)` |
| `AVFormatContext` (input) | `AVFormatOpenInput(&ctx, ...)` | None | `AVFormatCloseInput(&ctx)` |
| `AVFormatContext` (output) | `AVFormatAllocOutputContext2(&ctx, ...)` | None | Close `AVIOContext` if opened, then `AVFormatFreeContext(ctx)` |
| `AVFilterGraph` | `AVFilterGraphAlloc()` | None | `AVFilterGraphFree(&graph)` |
| `Array[T]` | Typed `Alloc*Array` constructor | None | `arr.Free()` (only when owned) |

`Array[T]` carries no length and performs no bounds checks on `Get` or `Set`. Indexing past the valid C buffer length is undefined behaviour.

### Error handling and stream-loop sentinels

- **`WrapErr(code int) error`**: Returns `nil` when `code >= 0`, or `AVError{Code: code}` when `code < 0`. Generated functions that return an FFmpeg status code already wrap negative codes through `WrapErr`.
- **`AVError`**: Implements `Error()` via `AVStrerror` and implements `Is(target error) bool` by comparing `Code`, so `errors.Is` matches sentinel errors reliably.
- **`EAgain`** (`AVError{Code: -C.EAGAIN}`): The decoder, encoder, or filter needs more input before it can emit output, or you must drain output before sending more input.
- **`AVErrorEOF`** (`AVError{Code: AVErrorEofConst}`): End of input stream, or the decoder/encoder/filter has finished draining after a flush.

```go
decCtx := ffmpeg.AVCodecAllocContext3(codec)
defer ffmpeg.AVCodecFreeContext(&decCtx)

pkt := ffmpeg.AVPacketAlloc()
defer ffmpeg.AVPacketFree(&pkt)

frame := ffmpeg.AVFrameAlloc()
defer ffmpeg.AVFrameFree(&frame)

if _, err := ffmpeg.AVCodecSendPacket(decCtx, pkt); err != nil {
    return err
}

for {
    _, err := ffmpeg.AVCodecReceiveFrame(decCtx, frame)
    if errors.Is(err, ffmpeg.EAgain) || errors.Is(err, ffmpeg.AVErrorEOF) {
        break
    }
    if err != nil {
        return err
    }

    // Process frame here.
    ffmpeg.AVFrameUnref(frame)
}
```

## FFmpeg API coverage

As of FFmpeg 8.1.2 (measured 2026-11-25), `ffmpeg-statigo` binds **92.1%** (`978 / 1062`, or about 92%) of the public C functions in the parsed FFmpeg headers.

### Metric definition

Coverage measures bound functions over the parsed public function set:

```text
coverage = bound functions / parsed public functions
```

- **Parsed public functions**: Every public C function that the generator visits across the allowlisted headers in `internal/generator/parser.go` (generated wrappers plus generator function skips).
- **Bound functions**: Generated wrappers in `functions.gen.go` plus hand-written wrappers in the root topic files (`iterate.go`, `uuid.go`, `swscale.go`, `avio.go`, `log_format.go`, `audio_fifo.go`, `tx.go` + `tx.c`, and others) that re-bind a skipped function.

A function that the generator skips but a topic file re-binds counts as covered, not missing.

### Coverage counts (FFmpeg 8.1.2)

| Quantity | Value | Source |
| --- | --- | --- |
| Generated function wrappers | 925 | `grep -cE '^func [A-Z][A-Za-z0-9_]*\(' functions.gen.go` |
| Function skips (generator) | 137 | Skip summary, symbols without a `.` |
| Struct-field / other skips | 66 | Skip summary, symbols containing a `.` |
| Total skip markers | 203 | Skip summary header (`= skipCeiling`) |
| Re-bound by hand (covered) | 53 | Skipped functions with a manual binding in topic files |
| Genuinely missing functions | 84 | `137 - 53` |
| **Numerator (bound)** | **978** | `925 + 53` |
| **Denominator (parsed public functions)** | **1062** | `925 + 137` |
| **Coverage** | **92.1%** | `978 / 1062` |

The generated-only lower bound is **87.1%** (`925 / 1062`). Hand-written topic wrappers re-bind 53 skipped functions and raise total coverage to **92.1%**.

### Methodology notes

- **Fixed-size array generation**: Adding fixed-size array support to the generator cut total skips by 42 (from 245 to 203). Struct-field skips fell by 28 (from 94 to 66) as 11 inline struct-array fields and 17 multi-dimensional array fields moved into `structs.gen.go`, retiring 3 manual accessors in `fields.go` and lowering struct-field manual re-binds from 8 to 5. Function skips fell by 14 (from 151 to 137) as 14 unwrapped fixed-size array functions moved into `functions.gen.go`.
- **Counting scope**: Counting `^func [A-Z][A-Za-z0-9_]*\(` includes the 54 generated `Swr*` and `Sws*` wrappers alongside `AV*` wrappers, matching the 1062 parsed functions recorded in `internal/generator/testdata/ir/structure.txt`. Function re-binds total 53 because `AVTxCall` invokes the `av_tx_fn` function pointer rather than a skipped C function.

### Excluded headers

Headers excluded by design in `internal/generator/parser.go` sit outside the denominator:

- **16 platform-specific hardware-acceleration headers**: D3D11VA, DXVA2, VDPAU, VideoToolbox, QSV, VA-API, CUDA, Vulkan, OpenCL, AMF, MediaCodec/JNI, and OpenHarmony. These require vendor SDKs or risk link failures on static builds. Cross-platform hardware acceleration stays available through the device and frames-context APIs in `hwcontext.h`.
- **`smpte_436m.h`**: Symbols are absent from the FFmpeg 8.1.2 static library.
- **Non-API and union-only headers**: `attributes.h`, `intreadwrite.h`, and `refstruct.h`.

### Reproducing the counts

Run inside the Nix development shell (the generator is hermetic and does not probe host compilers):

```sh
# Print the skip summary (total markers and sorted unique symbols).
go run ./internal/generator 2>&1 >/dev/null

# Confirm the run did not change the committed bindings.
git diff --stat -- '*.gen.go'

# Classify skips: function skips have no dot, struct-field skips contain one.
go run ./internal/generator 2>/tmp/skips.txt >/dev/null
tail -n +3 /tmp/skips.txt | sed 's/^  //' > /tmp/sym.txt
grep -vc '\.' /tmp/sym.txt   # function skips (137)
grep -c  '\.' /tmp/sym.txt   # struct-field / other skips (66)

# Count skipped symbols annotated with a manual binding.
grep '(manual binding:' /tmp/skips.txt | grep -vc '\.'   # function re-binds (53)
grep '(manual binding:' /tmp/skips.txt | grep -c  '\.'   # struct-field re-binds (5)
```
