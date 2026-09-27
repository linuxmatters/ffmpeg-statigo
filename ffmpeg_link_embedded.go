//go:build embedded

package ffmpeg

/*
#cgo linux,amd64 LDFLAGS: ${SRCDIR}/lib/embedded/linux_amd64/libffmpeg.a
#cgo linux,arm64 LDFLAGS: ${SRCDIR}/lib/embedded/linux_arm64/libffmpeg.a
#cgo darwin,amd64 LDFLAGS: ${SRCDIR}/lib/embedded/darwin_amd64/libffmpeg.a
#cgo darwin,arm64 LDFLAGS: ${SRCDIR}/lib/embedded/darwin_arm64/libffmpeg.a
#cgo windows,amd64 LDFLAGS: ${SRCDIR}/lib/embedded/windows_amd64/libffmpeg.a
#cgo windows,386 LDFLAGS: ${SRCDIR}/lib/embedded/windows_386/libffmpeg.a

#cgo linux LDFLAGS: -lm -ldl -lstdc++ -lpthread
#cgo darwin LDFLAGS: -lm -lc++ -lpthread
#cgo windows,amd64 LDFLAGS: -lstdc++ -lm -lws2_32 -lbcrypt -luser32
#cgo windows,386 LDFLAGS: -lstdc++ -lm -lws2_32 -lbcrypt -luser32
*/
import "C"
