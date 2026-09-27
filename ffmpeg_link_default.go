//go:build !embedded

package ffmpeg

/*
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/lib/linux_amd64
#cgo linux,arm64 LDFLAGS: -L${SRCDIR}/lib/linux_arm64
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/lib/darwin_amd64
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/lib/darwin_arm64

#cgo linux LDFLAGS: -lffmpeg -lm -ldl -lstdc++ -lpthread
#cgo darwin LDFLAGS: -lffmpeg -lstdc++ -lm -framework ApplicationServices -framework CoreVideo -framework CoreMedia -framework VideoToolbox -framework AudioToolbox
*/
import "C"
