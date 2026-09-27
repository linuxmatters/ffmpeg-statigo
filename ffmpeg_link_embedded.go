//go:build embedded

package ffmpeg

/*
#cgo linux,amd64 LDFLAGS: ${SRCDIR}/lib/embedded/linux_amd64/libffmpeg.a
#cgo linux,arm64 LDFLAGS: ${SRCDIR}/lib/embedded/linux_arm64/libffmpeg.a
#cgo darwin,amd64 LDFLAGS: ${SRCDIR}/lib/embedded/darwin_amd64/libffmpeg.a
#cgo darwin,arm64 LDFLAGS: ${SRCDIR}/lib/embedded/darwin_arm64/libffmpeg.a

#cgo linux LDFLAGS: -lm -ldl -lstdc++ -lpthread
#cgo darwin LDFLAGS: -lm -lc++ -lpthread
*/
import "C"
