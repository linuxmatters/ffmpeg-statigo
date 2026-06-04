// Package av is an optional high-level layer over the root ffmpeg package.
//
// The root ffmpeg package is the raw bridge to FFmpeg's C API and stays fully
// usable on its own. This package adds Go-idiomatic ownership and lifetime
// management on top of it. Use av when you want resource safety and a smaller
// surface; drop to the root package whenever you need the full C API.
//
// # Ownership and Close contract
//
// Every owned type implements [io.Closer]. Close frees the underlying FFmpeg
// handles in dependency order, so closing a parent releases the resources it
// owns. Close is idempotent: a second call is a safe no-op and never double-frees.
// Pair each constructor with a deferred Close.
//
// # Raw escape hatch
//
// Every owned type exposes a Raw method returning the underlying *ffmpeg.AV*
// handle. Any FFmpeg capability this layer does not surface stays reachable
// through Raw, so the high-level types never trap you. Do not free a handle
// obtained from Raw yourself; ownership remains with the av type and its Close.
//
// Example usage:
//
//	import "github.com/linuxmatters/ffmpeg-statigo/av"
//
//	func main() {
//		in, err := av.Open("input.mp4")
//		if err != nil {
//			panic(err)
//		}
//		defer in.Close()
//
//		// Reach the raw handle for anything this layer does not surface.
//		ctx := in.Raw()
//		_ = ctx
//	}
package av
