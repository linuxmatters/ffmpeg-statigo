//go:build amd64 || arm64

package ffmpeg

import "unsafe"

// Compile-time guard: this package targets 64-bit only.
// Holders of the invariant:
//   - the size_t -> uint64 entry in the generator's type map (generator.go)
//   - the C.ulong(len) runtime cast in AllocCStr (ffmpeg.go)
//
// The supported target list lives in the per-platform #cgo ... LDFLAGS block in ffmpeg.go.
//
// On a 32-bit target sizeof(uintptr) == 4, so the subtraction underflows
// the unsigned domain and the compiler rejects the constant expression.
const _ = uint(unsafe.Sizeof(uintptr(0)) - 8)
