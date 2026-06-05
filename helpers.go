package ffmpeg

/*
#include <libavutil/hwcontext.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func (s *AVRational) String() string {
	return fmt.Sprintf("%v/%v (%v)", s.Num(), s.Den(), s.Num()/s.Den())
}

// ToAVHWFramesContext converts an unsafe.Pointer (typically from AVBufferRef.Data())
// to an AVHWFramesContext wrapper. This is needed for configuring hardware frames
// contexts returned by AVHWFrameCtxAlloc().
func ToAVHWFramesContext(ptr unsafe.Pointer) *AVHWFramesContext {
	if ptr == nil {
		return nil
	}
	return &AVHWFramesContext{ptr: (*C.AVHWFramesContext)(ptr)}
}
