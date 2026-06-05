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
	num, den := s.Num(), s.Den()
	if den == 0 {
		return fmt.Sprintf("%v/%v (undefined)", num, den)
	}
	return fmt.Sprintf("%v/%v (%v)", num, den, num/den)
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
