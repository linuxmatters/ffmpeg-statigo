//go:build !embedded

package ffmpeg

/*
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/lib/linux_amd64
#cgo linux,arm64 LDFLAGS: -L${SRCDIR}/lib/linux_arm64
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/lib/darwin_amd64
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/lib/darwin_arm64
#cgo windows,amd64 LDFLAGS: ${SRCDIR}/lib/windows_amd64/libffmpeg.a
#cgo windows,386 LDFLAGS: ${SRCDIR}/lib/windows_386/libffmpeg.a

#cgo linux LDFLAGS: -lffmpeg -lm -ldl -lstdc++ -lpthread
#cgo darwin LDFLAGS: -lffmpeg -lstdc++ -lm -framework ApplicationServices -framework CoreVideo -framework CoreMedia -framework VideoToolbox -framework AudioToolbox
#cgo windows,amd64 LDFLAGS: -lstdc++ -lpthread -lm -lws2_32 -lbcrypt -luser32 -lntdll -luserenv -lcrypt32
#cgo windows,386 LDFLAGS: -lstdc++ -lpthread -lm -lws2_32 -lbcrypt -luser32 -lntdll -luserenv -lcrypt32
*/
import "C"
