package ffmpeg

/*
#include <stdint.h>
#include <libavutil/audio_fifo.h>
*/
import "C"

import "unsafe"

// AVAudioFifoWrite writes nbSamples samples per channel from data into af. The
// plane pointers must reference C-allocated memory (one plane per channel for
// planar formats, a single plane for packed).
//
// Returns the number of samples written, or a negative error code.
func AVAudioFifoWrite(af *AVAudioFifo, data []unsafe.Pointer, nbSamples, nbChannels int, sampleFmt AVSampleFormat) (int, error) {
	if len(data) < samplePlaneCount(nbChannels, sampleFmt) {
		return 0, WrapErr(AVErrorUnknownConst)
	}

	arr := newSamplePointerArray(data)
	defer arr.free()

	ret := C.av_audio_fifo_write(af.ptr, (*unsafe.Pointer)(unsafe.Pointer(arr.ptr)), C.int(nbSamples))
	return int(ret), WrapErr(int(ret))
}

// AVAudioFifoRead reads up to nbSamples samples per channel from af into data,
// removing them from the buffer. The plane pointers reference caller-owned
// destination memory C writes through directly (one plane per channel for planar
// formats, a single plane for packed).
//
// Returns the number of samples read, or a negative error code.
func AVAudioFifoRead(af *AVAudioFifo, data []unsafe.Pointer, nbSamples, nbChannels int, sampleFmt AVSampleFormat) (int, error) {
	if len(data) < samplePlaneCount(nbChannels, sampleFmt) {
		return 0, WrapErr(AVErrorUnknownConst)
	}

	arr := newSamplePointerArray(data)
	defer arr.free()

	ret := C.av_audio_fifo_read(af.ptr, (*unsafe.Pointer)(unsafe.Pointer(arr.ptr)), C.int(nbSamples))
	return int(ret), WrapErr(int(ret))
}

// AVAudioFifoPeek reads up to nbSamples samples per channel from af into data
// without removing them from the buffer. The plane pointers reference
// caller-owned destination memory C writes through directly (one plane per
// channel for planar formats, a single plane for packed).
//
// Returns the number of samples peeked, or a negative error code.
func AVAudioFifoPeek(af *AVAudioFifo, data []unsafe.Pointer, nbSamples, nbChannels int, sampleFmt AVSampleFormat) (int, error) {
	if len(data) < samplePlaneCount(nbChannels, sampleFmt) {
		return 0, WrapErr(AVErrorUnknownConst)
	}

	arr := newSamplePointerArray(data)
	defer arr.free()

	ret := C.av_audio_fifo_peek(af.ptr, (*unsafe.Pointer)(unsafe.Pointer(arr.ptr)), C.int(nbSamples))
	return int(ret), WrapErr(int(ret))
}

// AVAudioFifoPeekAt reads up to nbSamples samples per channel from af into data
// starting at offset, without removing them from the buffer. The plane pointers
// reference caller-owned destination memory C writes through directly (one plane
// per channel for planar formats, a single plane for packed).
//
// Returns the number of samples peeked, or a negative error code.
func AVAudioFifoPeekAt(af *AVAudioFifo, data []unsafe.Pointer, nbSamples, offset, nbChannels int, sampleFmt AVSampleFormat) (int, error) {
	if len(data) < samplePlaneCount(nbChannels, sampleFmt) {
		return 0, WrapErr(AVErrorUnknownConst)
	}

	arr := newSamplePointerArray(data)
	defer arr.free()

	ret := C.av_audio_fifo_peek_at(af.ptr, (*unsafe.Pointer)(unsafe.Pointer(arr.ptr)), C.int(nbSamples), C.int(offset))
	return int(ret), WrapErr(int(ret))
}
