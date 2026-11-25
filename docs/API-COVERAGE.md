# API coverage

As of FFmpeg 8.1.2 (measured 2026-11-25), ffmpeg-statigo binds **~92%** of the
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

## Counts (FFmpeg 8.1.2)

| Quantity | Value | Source |
| --- | --- | --- |
| Generated function wrappers | 925 | `grep -cE '^func [A-Z][A-Za-z0-9_]*\(' functions.gen.go` |
| Function skips (generator) | 137 | skip summary, symbols without a `.` |
| Struct-field / other skips | 66 | skip summary, symbols containing a `.` |
| Total skip markers | 203 | skip summary header (`= skipCeiling`) |
| Re-bound by hand (covered) | 53 | skipped functions with a manual binding in topic files |
| Genuinely missing functions | 84 | 137 - 53 |
| **Numerator** (bound) | **978** | 925 + 53 |
| **Denominator** (parsed public funcs) | **1062** | 925 + 137 |
| **Coverage** | **92.1%** | 978 / 1062 |

The generated-only lower bound is 87.1% (925 / 1062). Hand-written wrappers in
the topic files re-bind 53 skipped functions, which raises coverage to 92.1%
(978 / 1062, or ~92%).

Generator support for fixed-size arrays reduced total skips by 42, from 245 to
203. Struct-field skips fell by 28, from 94 to 66. The generator now emits 11
inline struct-array fields and 17 multi-dimensional array fields in
`structs.gen.go`. That change retired 3 manual accessors from `fields.go` and
reduced struct-field manual re-binds from 8 to 5. Function skips fell by 14,
from 151 to 137, as 14 unwrapped fixed-size array functions moved into
`functions.gen.go`.

Earlier counts matched only `^func AV` in `functions.gen.go` (857 generated
wrappers and 1008 parsed functions). That pattern missed 54 generated `Swr*`
and `Sws*` wrappers even though the skip total included skipped `swr_*` and
`sws_*` functions. The table now counts all exported wrappers in
`functions.gen.go` (`^func [A-Z]`), so the denominator matches the 1062 parsed
functions in `internal/generator/testdata/ir/structure.txt`. The re-bind count
also matches the live generator summary (53, not 54), because `AVTxCall`
invokes the `av_tx_fn` function pointer rather than a skipped C function.

## Scope

The figure covers the **parsed public API only**. Headers excluded by design in
`internal/generator/parser.go` (commented with reasons) are out of denominator:

- 16 platform-specific hardware-acceleration headers (D3D11VA, DXVA2, VDPAU,
  VideoToolbox, QSV, VA-API, CUDA, Vulkan, OpenCL, AMF, MediaCodec/JNI,
  OpenHarmony) that need vendor SDKs or risk link failures on the static Linux
  and macOS builds.
- `smpte_436m.h` (symbols absent from the FFmpeg 8.1.2 static lib).
- Non-API or union-only headers (`attributes.h`, `intreadwrite.h`,
  `refstruct.h`).

Hardware acceleration is still available through the cross-platform device and
frames-context APIs (`hwcontext.h`); only the vendor-specific binding headers
are excluded.

## Reproduce

Run inside the Nix dev shell (the generator is hermetic: predefined macros and
C library stubs are embedded, so no host compiler is probed):

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
grep '(manual binding:' /tmp/skips.txt | grep -vc '\.'   # function rebinds (matches the 53 figure)
grep '(manual binding:' /tmp/skips.txt | grep -c  '\.'   # struct-field rebinds (5)
```
