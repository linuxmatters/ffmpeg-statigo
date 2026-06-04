//go:build amd64 || arm64

package ffmpeg

import "unsafe"

// Compile-time guard: this package targets 64-bit only.
// Holders of the invariant:
//   - internal/generator/generator.go:30 (size_t -> uint64 codegen choice)
//   - ffmpeg.go:75 (C.ulong(len) runtime cast)
//
// The supported target list lives at ffmpeg.go:12.
//
// On a 32-bit target sizeof(uintptr) == 4, so the subtraction underflows
// the unsigned domain and the compiler rejects the constant expression.
const _ = uint(unsafe.Sizeof(uintptr(0)) - 8)
