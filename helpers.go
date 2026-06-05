package ffmpeg

/*
#include <stdint.h>
#include <stddef.h>
#include <libavutil/hwcontext.h>
#include <libavutil/mathematics.h>
#include <libavutil/mem.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// String returns a human-readable representation of the rational as
// "num/den (num/den)", or "num/den (undefined)" when the denominator is zero.
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

// AVRescaleDelta rescales a timestamp while preserving known durations. It is
// designed to be called per audio packet to scale the input timestamp to a
// different time base. Compared to a simple AVRescaleQ call, it is robust
// against possible inconsistent frame durations.
//
// last is a state variable that must be preserved across all subsequent calls
// for the same stream; for the first call set *last to AVNoptsValue. The C side
// reads and updates it through the supplied pointer, so the caller observes the
// new state on return.
//
//	@param[in]     inTb     Input time base
//	@param[in]     inTs     Input timestamp
//	@param[in]     fsTb     Duration time base; typically finer-grained
//	                        (greater) than inTb and outTb
//	@param[in]     duration Duration till the next call (i.e. duration of the
//	                        current packet/frame, in samples not seconds)
//	@param[in,out] last     Pointer to a timestamp expressed in terms of fsTb,
//	                        acting as a state variable
//	@param[in]     outTb    Output time base
//	@return                 Timestamp expressed in terms of outTb
func AVRescaleDelta(inTb *AVRational, inTs int64, fsTb *AVRational, duration int, last *int64, outTb *AVRational) int64 {
	ret := C.av_rescale_delta(
		inTb.value, C.int64_t(inTs),
		fsTb.value, C.int(duration),
		(*C.int64_t)(unsafe.Pointer(last)),
		outTb.value,
	)
	return int64(ret)
}

// AVSizeMult multiplies two size_t values, checking for overflow. On success it
// returns the product and a nil error; on overflow it returns 0 and an AVError
// wrapping AVERROR(EINVAL).
//
//	@param[in] a Operand of multiplication
//	@param[in] b Operand of multiplication
//	@return      The product, or 0 on overflow with a non-nil error
func AVSizeMult(a, b uint) (uint, error) {
	var r C.size_t
	ret := C.av_size_mult(C.size_t(a), C.size_t(b), &r)
	return uint(r), WrapErr(int(ret))
}
