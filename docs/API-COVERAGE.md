# API coverage

As of FFmpeg 8.1.1 (measured 2026-06-08), ffmpeg-statigo binds **~90%** of the
public C functions in the parsed FFmpeg headers.

## Metric definition

Coverage is reported over the *function* domain:

```
coverage = bound functions / parsed public functions
```

- **Parsed public functions** is the set of public C functions the generator
  visits across the allowlisted headers in `internal/generator/parser.go`. It
  equals the generated wrappers plus every function the generator skips.
- **Bound functions** is the generated wrappers plus the hand-written wrappers
  that re-bind a skipped function (variadics, fixed-size arrays, anonymous
  structs, function-pointer bridges) in the topic files (`iterate.go`,
  `uuid.go`, `swscale.go`, `avio.go`, `log_format.go`, `audio_fifo.go`,
  `tx.go` + `tx.c`, and so on).

A function the generator skips but a topic file re-binds counts as covered, not
missing.

## Counts (FFmpeg 8.1.1)

| Quantity | Value | Source |
| --- | --- | --- |
| Generated function wrappers | 857 | `grep -cE '^func AV[A-Za-z0-9_]+\(' functions.gen.go` |
| Function skips (generator) | 151 | skip summary, symbols without a `.` |
| Struct-field / other skips | 94 | skip summary, symbols containing a `.` |
| Total skip markers | 245 | skip summary header (`= skipCeiling`) |
| Re-bound by hand (covered) | 54 | skipped functions re-exposed in topic files |
| Genuinely missing functions | 97 | 151 - 54 |
| **Numerator** (bound) | **911** | 857 + 54 |
| **Denominator** (parsed public funcs) | **1008** | 857 + 151 |
| **Coverage** | **90.4%** | 911 / 1008 |

The previous README figure of 85% matched the generated-only lower bound
(857 / 1008 = 85.0%); it ignored the hand-written re-binds, so it understated
real coverage. The recent header-promotion and hand-binding work added bindings,
lifting the figure from ~85% to ~90% (906 / 1008 = 89.9%). The five additional
hand-written bindings for the AVAudioFifo data path and av_tx invocation bring
the total to 911 / 1008 = 90.4%, still rounding to ~90%.

## Scope

The figure covers the **parsed public API only**. Headers excluded by design in
`internal/generator/parser.go` (commented with reasons) are out of denominator:

- 16 platform-specific hardware-acceleration headers (D3D11VA, DXVA2, VDPAU,
  VideoToolbox, QSV, VA-API, CUDA, Vulkan, OpenCL, AMF, MediaCodec/JNI,
  OpenHarmony) that need vendor SDKs or risk link failures on the static Linux
  and macOS builds.
- `smpte_436m.h` (symbols absent from the FFmpeg 8.1.1 static lib).
- Non-API or union-only headers (`attributes.h`, `intreadwrite.h`,
  `refstruct.h`).

Hardware acceleration is still available through the cross-platform device and
frames-context APIs (`hwcontext.h`); only the vendor-specific binding headers
are excluded.

## Reproduce

Run inside the Nix dev shell (libclang 20.1.8, gcc 15.2.0):

```sh
# Print the skip summary (total markers + sorted unique symbols).
go run ./internal/generator 2>&1 >/dev/null

# Confirm the run did not change the committed bindings.
git diff --stat -- '*.gen.go'   # must be empty

# Classify skips: function skips have no dot, struct-field skips contain one.
go run ./internal/generator 2>/tmp/skips.txt >/dev/null
tail -n +3 /tmp/skips.txt | sed 's/^  //' > /tmp/sym.txt
grep -vc '\.' /tmp/sym.txt   # function skips
grep -c  '\.' /tmp/sym.txt   # struct-field / other skips
```

Re-derive the numerator and denominator with the commands in the table above.
The skip ceiling lives in `internal/generator/main.go` (`skipCeiling`); a run
that exceeds it fails, so the skip total tracks API drift across FFmpeg upgrades.

The generator now annotates each skipped symbol that has a hand-written binding.
The re-bind count is directly readable from the skip summary; no manual
cross-reference required:

```sh
# Count skipped-but-rebound symbols. Function skips have no dot in the symbol;
# struct-field skips do, so split them to match the counts table.
go run ./internal/generator 2>/tmp/skips.txt >/dev/null
grep '(manual binding:' /tmp/skips.txt | grep -vc '\.'   # function rebinds (matches the 54 figure)
grep '(manual binding:' /tmp/skips.txt | grep -c  '\.'   # struct-field rebinds
```
