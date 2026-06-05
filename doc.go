// Package ffmpeg provides Go bindings to FFmpeg 8.1.1 libraries.
//
// This package includes static FFmpeg libraries with hardware acceleration
// support (NVENC, VideoToolbox, Vulkan, QuickSync) and contemporary codecs
// (H.264, H.265, VP8/9, AV1).
//
// # Source layers
//
// The package is built in three tiers.
//
// Tier 1 — generated bindings: constants.gen.go, enums.gen.go, structs.gen.go,
// functions.gen.go, and callbacks.gen.go are emitted by internal/generator
// (libclang + dave/jennifer) from the FFmpeg headers. They cover the bulk of the
// API. Never edit these files; regenerate with just generate.
//
// Tier 2 — core and foundation: ffmpeg.go defines CGO directives, platform
// linker flags, and the base types every tier builds on (AVError/WrapErr, CStr).
// array.go holds the generic Array[T] type and its typed constructors.
// arch_guard.go enforces the 64-bit-only invariant at compile time.
//
// Tier 3 — hand-written topic wrappers: the generator deliberately skips C
// symbols it cannot safely express — function pointers, variadic functions,
// fixed-size array parameters, anonymous structs, and unions. Each skip is
// recorded with a reason, and the total is regression-capped by skipCeiling in
// internal/generator/main.go. The curated files cover: registry iterators
// (iterate.go), UUID types (uuid.go), anonymous-struct accessors (streamgroup.go,
// fields.go), generic option setter (opt.go), image planes (image.go), audio
// sample planes (samples.go), software scaling (swscale.go), audio resampling
// (swresample.go), custom I/O with a cgo.Handle callback bridge (avio.go +
// avio.c), av_log bridge to slog (log.go + log.c), variadic format shims
// (log_format.go), and cross-cutting helpers (helpers.go).
//
// # Optional high-level layer
//
// The subpackage github.com/linuxmatters/ffmpeg-statigo/av adds owned io.Closer
// pipeline wrappers (Input, Decoder, Encoder, FilterGraph, Output) for callers
// who prefer managed lifetimes over this raw bridge. See docs/PIPELINE.md.
//
// # Example usage
//
//	import "github.com/linuxmatters/ffmpeg-statigo"
//
//	func main() {
//		var ctx *ffmpeg.AVFormatContext
//		url := ffmpeg.ToCStr("input.mp4")
//		defer url.Free()
//
//		_, err := ffmpeg.AVFormatOpenInput(&ctx, url, nil, nil)
//		if err != nil {
//			panic(err)
//		}
//		defer ffmpeg.AVFormatFreeContext(ctx)
//	}
//
// For complete examples, see the examples/ directory.
package ffmpeg
