package ffmpeg

/*
#include <stdint.h>
#include <libavutil/rational.h>
#include <libavutil/parseutils.h>
#include <libavformat/avformat.h>
*/
import "C"

import "unsafe"

// Hand-written wrappers for libav lookup and parse helpers the generator skips.
// av_parse_ratio and av_parse_video_rate fill an AVRational out-param the
// generator cannot classify as in or out; av_codec_get_tag2 returns its codec
// tag through an unsigned int out-param. The out-params become the caller's
// supplied pointer plus a Go return value, with the status code folded into
// error via WrapErr.

// AVParseRatio parses str and stores the parsed ratio in q.
//
// A ratio with infinite (1/0) or negative value is considered valid, so check
// the stored ratio if you want to exclude those values. The undefined value can
// be expressed using the "0:0" string.
//
//	@param[in,out] q         AVRational which will contain the ratio
//	@param[in]     str       string to parse: a num:den pair, a float, or an
//	                         expression
//	@param[in]     max       maximum allowed numerator and denominator
//	@param[in]     logOffset log level offset applied to the log level of logCtx
//	@param[in]     logCtx    parent logging context, or nil
//	@return                  nil on success, an AVError otherwise
func AVParseRatio(q *AVRational, str *CStr, max, logOffset int, logCtx unsafe.Pointer) error {
	ret := C.av_parse_ratio(&q.value, cstrPtr(str), C.int(max), C.int(logOffset), logCtx)
	return WrapErr(int(ret))
}

// AVParseVideoRate parses str and stores the detected frame rate in rate.
//
//	@param[in,out] rate AVRational which will contain the detected frame rate
//	@param[in]     str  string to parse: a rate_num/rate_den pair, a float, or a
//	                    valid video rate abbreviation
//	@return             nil on success, an AVError otherwise
func AVParseVideoRate(rate *AVRational, str *CStr) error {
	ret := C.av_parse_video_rate(&rate.value, cstrPtr(str))
	return WrapErr(int(ret))
}

// AVCodecGetTag2 looks up the codec tag for id in the supplied tag table,
// writing the tag through the out-param and reporting whether a match was found.
//
// tags is the address of a single AVCodecTag table pointer, such as the result
// of AVFormatGetRiffVideoTags; the C side only reads it. The C function walks a
// NULL-terminated list of tables, so the wrapper NULL-terminates the single
// table internally before calling it.
//
//	@param[in]  tags address of an AVCodecTag table of codec_id-codec_tag pairs
//	@param[in]  id   codec ID to match to a codec tag
//	@return          the matched tag and true on success, 0 and false if no tag
//	                 was found for id
func AVCodecGetTag2(tags **AVCodecTag, id AVCodecID) (tag uint, found bool) {
	if tags == nil || *tags == nil {
		return 0, false
	}

	// av_codec_get_tag2 walks the outer list until it reaches a NULL table
	// pointer, so the list it receives must be NULL-terminated. The caller
	// supplies the address of a single table pointer, which is not terminated;
	// wrap it in a two-entry, NULL-terminated list so a missed lookup stops at
	// the terminator instead of reading past the table pointer.
	list := []*C.struct_AVCodecTag{(*tags).ptr, nil}

	var ctag C.uint
	ret := C.av_codec_get_tag2(&list[0], C.enum_AVCodecID(id), &ctag)
	return uint(ctag), ret != 0
}
