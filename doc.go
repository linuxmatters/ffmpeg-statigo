// Package ffmpeg provides Go bindings to FFmpeg 8.0 libraries.
//
// This package includes static FFmpeg libraries with hardware acceleration
// support (NVENC, VideoToolbox, Vulkan, QuickSync) and contemporary codecs
// (H.264, H.265, VP8/9, AV1).
//
// The package provides direct access to FFmpeg's C API through CGO bindings.
// All FFmpeg functions are available with their original names prefixed with AV*.
//
// Example usage:
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
