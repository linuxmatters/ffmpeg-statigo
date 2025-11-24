package ffmpeg

import "unsafe"

// #include <libavcodec/avcodec.h>
// #include <libavcodec/bsf.h>
// #include <libavcodec/codec.h>
// #include <libavcodec/codec_desc.h>
// #include <libavcodec/codec_id.h>
// #include <libavcodec/codec_par.h>
// #include <libavcodec/defs.h>
// #include <libavcodec/packet.h>
// #include <libavcodec/version.h>
// #include <libavcodec/version_major.h>
// #include <libavdevice/avdevice.h>
// #include <libavdevice/version.h>
// #include <libavdevice/version_major.h>
// #include <libavfilter/avfilter.h>
// #include <libavfilter/buffersink.h>
// #include <libavfilter/buffersrc.h>
// #include <libavfilter/version.h>
// #include <libavfilter/version_major.h>
// #include <libavformat/avformat.h>
// #include <libavformat/avio.h>
// #include <libavformat/version.h>
// #include <libavformat/version_major.h>
// #include <libavutil/adler32.h>
// #include <libavutil/aes.h>
// #include <libavutil/aes_ctr.h>
// #include <libavutil/ambient_viewing_environment.h>
// #include <libavutil/audio_fifo.h>
// #include <libavutil/avassert.h>
// #include <libavutil/avconfig.h>
// #include <libavutil/avstring.h>
// #include <libavutil/avutil.h>
// #include <libavutil/base64.h>
// #include <libavutil/blowfish.h>
// #include <libavutil/bprint.h>
// #include <libavutil/bswap.h>
// #include <libavutil/buffer.h>
// #include <libavutil/camellia.h>
// #include <libavutil/cast5.h>
// #include <libavutil/channel_layout.h>
// #include <libavutil/common.h>
// #include <libavutil/container_fifo.h>
// #include <libavutil/cpu.h>
// #include <libavutil/crc.h>
// #include <libavutil/csp.h>
// #include <libavutil/des.h>
// #include <libavutil/detection_bbox.h>
// #include <libavutil/dict.h>
// #include <libavutil/display.h>
// #include <libavutil/dovi_meta.h>
// #include <libavutil/downmix_info.h>
// #include <libavutil/encryption_info.h>
// #include <libavutil/error.h>
// #include <libavutil/eval.h>
// #include <libavutil/executor.h>
// #include <libavutil/ffversion.h>
// #include <libavutil/fifo.h>
// #include <libavutil/file.h>
// #include <libavutil/film_grain_params.h>
// #include <libavutil/frame.h>
// #include <libavutil/hash.h>
// #include <libavutil/hdr_dynamic_metadata.h>
// #include <libavutil/hdr_dynamic_vivid_metadata.h>
// #include <libavutil/hmac.h>
// #include <libavutil/hwcontext.h>
// #include <libavutil/iamf.h>
// #include <libavutil/imgutils.h>
// #include <libavutil/intfloat.h>
// #include <libavutil/lfg.h>
// #include <libavutil/log.h>
// #include <libavutil/lzo.h>
// #include <libavutil/macros.h>
// #include <libavutil/mastering_display_metadata.h>
// #include <libavutil/mathematics.h>
// #include <libavutil/mem.h>
// #include <libavutil/motion_vector.h>
// #include <libavutil/murmur3.h>
// #include <libavutil/opt.h>
// #include <libavutil/parseutils.h>
// #include <libavutil/pixdesc.h>
// #include <libavutil/pixelutils.h>
// #include <libavutil/pixfmt.h>
// #include <libavutil/random_seed.h>
// #include <libavutil/rational.h>
// #include <libavutil/rc4.h>
// #include <libavutil/replaygain.h>
// #include <libavutil/ripemd.h>
// #include <libavutil/samplefmt.h>
// #include <libavutil/sha.h>
// #include <libavutil/sha512.h>
// #include <libavutil/spherical.h>
// #include <libavutil/stereo3d.h>
// #include <libavutil/tdrdi.h>
// #include <libavutil/tea.h>
// #include <libavutil/threadmessage.h>
// #include <libavutil/time.h>
// #include <libavutil/timecode.h>
// #include <libavutil/timestamp.h>
// #include <libavutil/tree.h>
// #include <libavutil/twofish.h>
// #include <libavutil/tx.h>
// #include <libavutil/uuid.h>
// #include <libavutil/version.h>
// #include <libavutil/video_enc_params.h>
// #include <libavutil/video_hint.h>
// #include <libavutil/xtea.h>
// #include <libswresample/version.h>
// #include <libswresample/version_major.h>
// #include <libswresample/swresample.h>
// #include <libswscale/version.h>
// #include <libswscale/version_major.h>
// #include <libswscale/swscale.h>
import "C"

// --- Function avcodec_version ---

// AVCodecVersion wraps avcodec_version.
//
//	Return the LIBAVCODEC_VERSION_INT constant.
func AVCodecVersion() uint {
	ret := C.avcodec_version()
	return uint(ret)
}

// --- Function avcodec_configuration ---

// AVCodecConfiguration wraps avcodec_configuration.
//
//	Return the libavcodec build-time configuration.
func AVCodecConfiguration() *CStr {
	ret := C.avcodec_configuration()
	return wrapCStr(ret)
}

// --- Function avcodec_license ---

// AVCodecLicense wraps avcodec_license.
//
//	Return the libavcodec license.
func AVCodecLicense() *CStr {
	ret := C.avcodec_license()
	return wrapCStr(ret)
}

// --- Function avcodec_alloc_context3 ---

// AVCodecAllocContext3 wraps avcodec_alloc_context3.
/*
  Allocate an AVCodecContext and set its fields to default values. The
  resulting struct should be freed with avcodec_free_context().

  @param codec if non-NULL, allocate private data and initialize defaults
               for the given codec. It is illegal to then call avcodec_open2()
               with a different codec.
               If NULL, then the codec-specific defaults won't be initialized,
               which may result in suboptimal default settings (this is
               important mainly for encoders, e.g. libx264).

  @return An AVCodecContext filled with default values or NULL on failure.
*/
func AVCodecAllocContext3(codec *AVCodec) *AVCodecContext {
	var tmpcodec *C.AVCodec
	if codec != nil {
		tmpcodec = codec.ptr
	}
	ret := C.avcodec_alloc_context3(tmpcodec)
	var retMapped *AVCodecContext
	if ret != nil {
		retMapped = &AVCodecContext{ptr: ret}
	}
	return retMapped
}

// --- Function avcodec_free_context ---

// AVCodecFreeContext wraps avcodec_free_context.
/*
  Free the codec context and everything associated with it and write NULL to
  the provided pointer.
*/
func AVCodecFreeContext(avctx **AVCodecContext) {
	var ptravctx **C.AVCodecContext
	var tmpavctx *C.AVCodecContext
	var oldTmpavctx *C.AVCodecContext
	if avctx != nil {
		inneravctx := *avctx
		if inneravctx != nil {
			tmpavctx = inneravctx.ptr
			oldTmpavctx = tmpavctx
		}
		ptravctx = &tmpavctx
	}
	C.avcodec_free_context(ptravctx)
	if tmpavctx != oldTmpavctx && avctx != nil {
		if tmpavctx != nil {
			*avctx = &AVCodecContext{ptr: tmpavctx}
		} else {
			*avctx = nil
		}
	}
}

// --- Function avcodec_get_class ---

// AVCodecGetClass wraps avcodec_get_class.
/*
  Get the AVClass for AVCodecContext. It can be used in combination with
  AV_OPT_SEARCH_FAKE_OBJ for examining options.

  @see av_opt_find().
*/
func AVCodecGetClass() *AVClass {
	ret := C.avcodec_get_class()
	var retMapped *AVClass
	if ret != nil {
		retMapped = &AVClass{ptr: ret}
	}
	return retMapped
}

// --- Function avcodec_get_subtitle_rect_class ---

// AVCodecGetSubtitleRectClass wraps avcodec_get_subtitle_rect_class.
/*
  Get the AVClass for AVSubtitleRect. It can be used in combination with
  AV_OPT_SEARCH_FAKE_OBJ for examining options.

  @see av_opt_find().
*/
func AVCodecGetSubtitleRectClass() *AVClass {
	ret := C.avcodec_get_subtitle_rect_class()
	var retMapped *AVClass
	if ret != nil {
		retMapped = &AVClass{ptr: ret}
	}
	return retMapped
}

// --- Function avcodec_parameters_from_context ---

// AVCodecParametersFromContext wraps avcodec_parameters_from_context.
/*
  Fill the parameters struct based on the values from the supplied codec
  context. Any allocated fields in par are freed and replaced with duplicates
  of the corresponding fields in codec.

  @return >= 0 on success, a negative AVERROR code on failure
*/
func AVCodecParametersFromContext(par *AVCodecParameters, codec *AVCodecContext) (int, error) {
	var tmppar *C.AVCodecParameters
	if par != nil {
		tmppar = par.ptr
	}
	var tmpcodec *C.AVCodecContext
	if codec != nil {
		tmpcodec = codec.ptr
	}
	ret := C.avcodec_parameters_from_context(tmppar, tmpcodec)
	return int(ret), WrapErr(int(ret))
}

// --- Function avcodec_parameters_to_context ---

// AVCodecParametersToContext wraps avcodec_parameters_to_context.
/*
  Fill the codec context based on the values from the supplied codec
  parameters. Any allocated fields in codec that have a corresponding field in
  par are freed and replaced with duplicates of the corresponding field in par.
  Fields in codec that do not have a counterpart in par are not touched.

  @return >= 0 on success, a negative AVERROR code on failure.
*/
func AVCodecParametersToContext(codec *AVCodecContext, par *AVCodecParameters) (int, error) {
	var tmpcodec *C.AVCodecContext
	if codec != nil {
		tmpcodec = codec.ptr
	}
	var tmppar *C.AVCodecParameters
	if par != nil {
		tmppar = par.ptr
	}
	ret := C.avcodec_parameters_to_context(tmpcodec, tmppar)
	return int(ret), WrapErr(int(ret))
}

// --- Function avcodec_open2 ---

// AVCodecOpen2 wraps avcodec_open2.
/*
  Initialize the AVCodecContext to use the given AVCodec. Prior to using this
  function the context has to be allocated with avcodec_alloc_context3().

  The functions avcodec_find_decoder_by_name(), avcodec_find_encoder_by_name(),
  avcodec_find_decoder() and avcodec_find_encoder() provide an easy way for
  retrieving a codec.

  Depending on the codec, you might need to set options in the codec context
  also for decoding (e.g. width, height, or the pixel or audio sample format in
  the case the information is not available in the bitstream, as when decoding
  raw audio or video).

  Options in the codec context can be set either by setting them in the options
  AVDictionary, or by setting the values in the context itself, directly or by
  using the av_opt_set() API before calling this function.

  Example:
  @code
  av_dict_set(&opts, "b", "2.5M", 0);
  codec = avcodec_find_decoder(AV_CODEC_ID_H264);
  if (!codec)
      exit(1);

  context = avcodec_alloc_context3(codec);

  if (avcodec_open2(context, codec, opts) < 0)
      exit(1);
  @endcode

  In the case AVCodecParameters are available (e.g. when demuxing a stream
  using libavformat, and accessing the AVStream contained in the demuxer), the
  codec parameters can be copied to the codec context using
  avcodec_parameters_to_context(), as in the following example:

  @code
  AVStream *stream = ...;
  context = avcodec_alloc_context3(codec);
  if (avcodec_parameters_to_context(context, stream->codecpar) < 0)
      exit(1);
  if (avcodec_open2(context, codec, NULL) < 0)
      exit(1);
  @endcode

  @note Always call this function before using decoding routines (such as
  @ref avcodec_receive_frame()).

  @param avctx The context to initialize.
  @param codec The codec to open this context for. If a non-NULL codec has been
               previously passed to avcodec_alloc_context3() or
               for this context, then this parameter MUST be either NULL or
               equal to the previously passed codec.
  @param options A dictionary filled with AVCodecContext and codec-private
                 options, which are set on top of the options already set in
                 avctx, can be NULL. On return this object will be filled with
                 options that were not found in the avctx codec context.

  @return zero on success, a negative value on error
  @see avcodec_alloc_context3(), avcodec_find_decoder(), avcodec_find_encoder(),
       av_dict_set(), av_opt_set(), av_opt_find(), avcodec_parameters_to_context()
*/
func AVCodecOpen2(avctx *AVCodecContext, codec *AVCodec, options **AVDictionary) (int, error) {
	var tmpavctx *C.AVCodecContext
	if avctx != nil {
		tmpavctx = avctx.ptr
	}
	var tmpcodec *C.AVCodec
	if codec != nil {
		tmpcodec = codec.ptr
	}
	var ptroptions **C.AVDictionary
	var tmpoptions *C.AVDictionary
	var oldTmpoptions *C.AVDictionary
	if options != nil {
		inneroptions := *options
		if inneroptions != nil {
			tmpoptions = inneroptions.ptr
			oldTmpoptions = tmpoptions
		}
		ptroptions = &tmpoptions
	}
	ret := C.avcodec_open2(tmpavctx, tmpcodec, ptroptions)
	if tmpoptions != oldTmpoptions && options != nil {
		if tmpoptions != nil {
			*options = &AVDictionary{ptr: tmpoptions}
		} else {
			*options = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function avsubtitle_free ---

// AVSubtitleFree wraps avsubtitle_free.
/*
  Free all allocated data in the given subtitle struct.

  @param sub AVSubtitle to free.
*/
func AVSubtitleFree(sub *AVSubtitle) {
	var tmpsub *C.AVSubtitle
	if sub != nil {
		tmpsub = sub.ptr
	}
	C.avsubtitle_free(tmpsub)
}

// --- Function avcodec_default_get_buffer2 ---

// AVCodecDefaultGetBuffer2 wraps avcodec_default_get_buffer2.
/*
  The default callback for AVCodecContext.get_buffer2(). It is made public so
  it can be called by custom get_buffer2() implementations for decoders without
  AV_CODEC_CAP_DR1 set.
*/
func AVCodecDefaultGetBuffer2(s *AVCodecContext, frame *AVFrame, flags int) (int, error) {
	var tmps *C.AVCodecContext
	if s != nil {
		tmps = s.ptr
	}
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.avcodec_default_get_buffer2(tmps, tmpframe, C.int(flags))
	return int(ret), WrapErr(int(ret))
}

// --- Function avcodec_default_get_encode_buffer ---

// AVCodecDefaultGetEncodeBuffer wraps avcodec_default_get_encode_buffer.
/*
  The default callback for AVCodecContext.get_encode_buffer(). It is made public so
  it can be called by custom get_encode_buffer() implementations for encoders without
  AV_CODEC_CAP_DR1 set.
*/
func AVCodecDefaultGetEncodeBuffer(s *AVCodecContext, pkt *AVPacket, flags int) (int, error) {
	var tmps *C.AVCodecContext
	if s != nil {
		tmps = s.ptr
	}
	var tmppkt *C.AVPacket
	if pkt != nil {
		tmppkt = pkt.ptr
	}
	ret := C.avcodec_default_get_encode_buffer(tmps, tmppkt, C.int(flags))
	return int(ret), WrapErr(int(ret))
}

// --- Function avcodec_align_dimensions ---

// AVCodecAlignDimensions wraps avcodec_align_dimensions.
/*
  Modify width and height values so that they will result in a memory
  buffer that is acceptable for the codec if you do not use any horizontal
  padding.

  May only be used if a codec with AV_CODEC_CAP_DR1 has been opened.
*/
func AVCodecAlignDimensions(s *AVCodecContext, width *int, height *int) {
	var tmps *C.AVCodecContext
	if s != nil {
		tmps = s.ptr
	}
	C.avcodec_align_dimensions(tmps, (*C.int)(unsafe.Pointer(width)), (*C.int)(unsafe.Pointer(height)))
}

// --- Function avcodec_align_dimensions2 ---

// avcodec_align_dimensions2 skipped due to const array param linesizeAlign

// --- Function avcodec_decode_subtitle2 ---

// AVCodecDecodeSubtitle2 wraps avcodec_decode_subtitle2.
/*
  Decode a subtitle message.
  Return a negative value on error, otherwise return the number of bytes used.
  If no subtitle could be decompressed, got_sub_ptr is zero.
  Otherwise, the subtitle is stored in *sub.
  Note that AV_CODEC_CAP_DR1 is not available for subtitle codecs. This is for
  simplicity, because the performance difference is expected to be negligible
  and reusing a get_buffer written for video codecs would probably perform badly
  due to a potentially very different allocation pattern.

  Some decoders (those marked with AV_CODEC_CAP_DELAY) have a delay between input
  and output. This means that for some packets they will not immediately
  produce decoded output and need to be flushed at the end of decoding to get
  all the decoded data. Flushing is done by calling this function with packets
  with avpkt->data set to NULL and avpkt->size set to 0 until it stops
  returning subtitles. It is safe to flush even those decoders that are not
  marked with AV_CODEC_CAP_DELAY, then no subtitles will be returned.

  @note The AVCodecContext MUST have been opened with @ref avcodec_open2()
  before packets may be fed to the decoder.

  @param avctx the codec context
  @param[out] sub The preallocated AVSubtitle in which the decoded subtitle will be stored,
                  must be freed with avsubtitle_free if *got_sub_ptr is set.
  @param[in,out] got_sub_ptr Zero if no subtitle could be decompressed, otherwise, it is nonzero.
  @param[in] avpkt The input AVPacket containing the input buffer.
*/
func AVCodecDecodeSubtitle2(avctx *AVCodecContext, sub *AVSubtitle, gotSubPtr *int, avpkt *AVPacket) (int, error) {
	var tmpavctx *C.AVCodecContext
	if avctx != nil {
		tmpavctx = avctx.ptr
	}
	var tmpsub *C.AVSubtitle
	if sub != nil {
		tmpsub = sub.ptr
	}
	var tmpavpkt *C.AVPacket
	if avpkt != nil {
		tmpavpkt = avpkt.ptr
	}
	ret := C.avcodec_decode_subtitle2(tmpavctx, tmpsub, (*C.int)(unsafe.Pointer(gotSubPtr)), tmpavpkt)
	return int(ret), WrapErr(int(ret))
}

// --- Function avcodec_send_packet ---

// AVCodecSendPacket wraps avcodec_send_packet.
/*
  Supply raw packet data as input to a decoder.

  Internally, this call will copy relevant AVCodecContext fields, which can
  influence decoding per-packet, and apply them when the packet is actually
  decoded. (For example AVCodecContext.skip_frame, which might direct the
  decoder to drop the frame contained by the packet sent with this function.)

  @warning The input buffer, avpkt->data must be AV_INPUT_BUFFER_PADDING_SIZE
           larger than the actual read bytes because some optimized bitstream
           readers read 32 or 64 bits at once and could read over the end.

  @note The AVCodecContext MUST have been opened with @ref avcodec_open2()
        before packets may be fed to the decoder.

  @param avctx codec context
  @param[in] avpkt The input AVPacket. Usually, this will be a single video
                   frame, or several complete audio frames.
                   Ownership of the packet remains with the caller, and the
                   decoder will not write to the packet. The decoder may create
                   a reference to the packet data (or copy it if the packet is
                   not reference-counted).
                   Unlike with older APIs, the packet is always fully consumed,
                   and if it contains multiple frames (e.g. some audio codecs),
                   will require you to call avcodec_receive_frame() multiple
                   times afterwards before you can send a new packet.
                   It can be NULL (or an AVPacket with data set to NULL and
                   size set to 0); in this case, it is considered a flush
                   packet, which signals the end of the stream. Sending the
                   first flush packet will return success. Subsequent ones are
                   unnecessary and will return AVERROR_EOF. If the decoder
                   still has frames buffered, it will return them after sending
                   a flush packet.

  @retval 0                 success
  @retval AVERROR(EAGAIN)   input is not accepted in the current state - user
                            must read output with avcodec_receive_frame() (once
                            all output is read, the packet should be resent,
                            and the call will not fail with EAGAIN).
  @retval AVERROR_EOF       the decoder has been flushed, and no new packets can be
                            sent to it (also returned if more than 1 flush
                            packet is sent)
  @retval AVERROR(EINVAL)   codec not opened, it is an encoder, or requires flush
  @retval AVERROR(ENOMEM)   failed to add packet to internal queue, or similar
  @retval "another negative error code" legitimate decoding errors
*/
func AVCodecSendPacket(avctx *AVCodecContext, avpkt *AVPacket) (int, error) {
	var tmpavctx *C.AVCodecContext
	if avctx != nil {
		tmpavctx = avctx.ptr
	}
	var tmpavpkt *C.AVPacket
	if avpkt != nil {
		tmpavpkt = avpkt.ptr
	}
	ret := C.avcodec_send_packet(tmpavctx, tmpavpkt)
	return int(ret), WrapErr(int(ret))
}

// --- Function avcodec_receive_frame ---

// AVCodecReceiveFrame wraps avcodec_receive_frame.
/*
  Return decoded output data from a decoder or encoder (when the
  @ref AV_CODEC_FLAG_RECON_FRAME flag is used).

  @param avctx codec context
  @param frame This will be set to a reference-counted video or audio
               frame (depending on the decoder type) allocated by the
               codec. Note that the function will always call
               av_frame_unref(frame) before doing anything else.

  @retval 0                success, a frame was returned
  @retval AVERROR(EAGAIN)  output is not available in this state - user must
                           try to send new input
  @retval AVERROR_EOF      the codec has been fully flushed, and there will be
                           no more output frames
  @retval AVERROR(EINVAL)  codec not opened, or it is an encoder without the
                           @ref AV_CODEC_FLAG_RECON_FRAME flag enabled
  @retval "other negative error code" legitimate decoding errors
*/
func AVCodecReceiveFrame(avctx *AVCodecContext, frame *AVFrame) (int, error) {
	var tmpavctx *C.AVCodecContext
	if avctx != nil {
		tmpavctx = avctx.ptr
	}
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.avcodec_receive_frame(tmpavctx, tmpframe)
	return int(ret), WrapErr(int(ret))
}

// --- Function avcodec_send_frame ---

// AVCodecSendFrame wraps avcodec_send_frame.
/*
  Supply a raw video or audio frame to the encoder. Use avcodec_receive_packet()
  to retrieve buffered output packets.

  @param avctx     codec context
  @param[in] frame AVFrame containing the raw audio or video frame to be encoded.
                   Ownership of the frame remains with the caller, and the
                   encoder will not write to the frame. The encoder may create
                   a reference to the frame data (or copy it if the frame is
                   not reference-counted).
                   It can be NULL, in which case it is considered a flush
                   packet.  This signals the end of the stream. If the encoder
                   still has packets buffered, it will return them after this
                   call. Once flushing mode has been entered, additional flush
                   packets are ignored, and sending frames will return
                   AVERROR_EOF.

                   For audio:
                   If AV_CODEC_CAP_VARIABLE_FRAME_SIZE is set, then each frame
                   can have any number of samples.
                   If it is not set, frame->nb_samples must be equal to
                   avctx->frame_size for all frames except the last.
                   The final frame may be smaller than avctx->frame_size.
  @retval 0                 success
  @retval AVERROR(EAGAIN)   input is not accepted in the current state - user must
                            read output with avcodec_receive_packet() (once all
                            output is read, the packet should be resent, and the
                            call will not fail with EAGAIN).
  @retval AVERROR_EOF       the encoder has been flushed, and no new frames can
                            be sent to it
  @retval AVERROR(EINVAL)   codec not opened, it is a decoder, or requires flush
  @retval AVERROR(ENOMEM)   failed to add packet to internal queue, or similar
  @retval "another negative error code" legitimate encoding errors
*/
func AVCodecSendFrame(avctx *AVCodecContext, frame *AVFrame) (int, error) {
	var tmpavctx *C.AVCodecContext
	if avctx != nil {
		tmpavctx = avctx.ptr
	}
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.avcodec_send_frame(tmpavctx, tmpframe)
	return int(ret), WrapErr(int(ret))
}

// --- Function avcodec_receive_packet ---

// AVCodecReceivePacket wraps avcodec_receive_packet.
/*
  Read encoded data from the encoder.

  @param avctx codec context
  @param avpkt This will be set to a reference-counted packet allocated by the
               encoder. Note that the function will always call
               av_packet_unref(avpkt) before doing anything else.
  @retval 0               success
  @retval AVERROR(EAGAIN) output is not available in the current state - user must
                          try to send input
  @retval AVERROR_EOF     the encoder has been fully flushed, and there will be no
                          more output packets
  @retval AVERROR(EINVAL) codec not opened, or it is a decoder
  @retval "another negative error code" legitimate encoding errors
*/
func AVCodecReceivePacket(avctx *AVCodecContext, avpkt *AVPacket) (int, error) {
	var tmpavctx *C.AVCodecContext
	if avctx != nil {
		tmpavctx = avctx.ptr
	}
	var tmpavpkt *C.AVPacket
	if avpkt != nil {
		tmpavpkt = avpkt.ptr
	}
	ret := C.avcodec_receive_packet(tmpavctx, tmpavpkt)
	return int(ret), WrapErr(int(ret))
}

// --- Function avcodec_get_hw_frames_parameters ---

// AVCodecGetHWFramesParameters wraps avcodec_get_hw_frames_parameters.
/*
  Create and return a AVHWFramesContext with values adequate for hardware
  decoding. This is meant to get called from the get_format callback, and is
  a helper for preparing a AVHWFramesContext for AVCodecContext.hw_frames_ctx.
  This API is for decoding with certain hardware acceleration modes/APIs only.

  The returned AVHWFramesContext is not initialized. The caller must do this
  with av_hwframe_ctx_init().

  Calling this function is not a requirement, but makes it simpler to avoid
  codec or hardware API specific details when manually allocating frames.

  Alternatively to this, an API user can set AVCodecContext.hw_device_ctx,
  which sets up AVCodecContext.hw_frames_ctx fully automatically, and makes
  it unnecessary to call this function or having to care about
  AVHWFramesContext initialization at all.

  There are a number of requirements for calling this function:

  - It must be called from get_format with the same avctx parameter that was
    passed to get_format. Calling it outside of get_format is not allowed, and
    can trigger undefined behavior.
  - The function is not always supported (see description of return values).
    Even if this function returns successfully, hwaccel initialization could
    fail later. (The degree to which implementations check whether the stream
    is actually supported varies. Some do this check only after the user's
    get_format callback returns.)
  - The hw_pix_fmt must be one of the choices suggested by get_format. If the
    user decides to use a AVHWFramesContext prepared with this API function,
    the user must return the same hw_pix_fmt from get_format.
  - The device_ref passed to this function must support the given hw_pix_fmt.
  - After calling this API function, it is the user's responsibility to
    initialize the AVHWFramesContext (returned by the out_frames_ref parameter),
    and to set AVCodecContext.hw_frames_ctx to it. If done, this must be done
    before returning from get_format (this is implied by the normal
    AVCodecContext.hw_frames_ctx API rules).
  - The AVHWFramesContext parameters may change every time time get_format is
    called. Also, AVCodecContext.hw_frames_ctx is reset before get_format. So
    you are inherently required to go through this process again on every
    get_format call.
  - It is perfectly possible to call this function without actually using
    the resulting AVHWFramesContext. One use-case might be trying to reuse a
    previously initialized AVHWFramesContext, and calling this API function
    only to test whether the required frame parameters have changed.
  - Fields that use dynamically allocated values of any kind must not be set
    by the user unless setting them is explicitly allowed by the documentation.
    If the user sets AVHWFramesContext.free and AVHWFramesContext.user_opaque,
    the new free callback must call the potentially set previous free callback.
    This API call may set any dynamically allocated fields, including the free
    callback.

  The function will set at least the following fields on AVHWFramesContext
  (potentially more, depending on hwaccel API):

  - All fields set by av_hwframe_ctx_alloc().
  - Set the format field to hw_pix_fmt.
  - Set the sw_format field to the most suited and most versatile format. (An
    implication is that this will prefer generic formats over opaque formats
    with arbitrary restrictions, if possible.)
  - Set the width/height fields to the coded frame size, rounded up to the
    API-specific minimum alignment.
  - Only _if_ the hwaccel requires a pre-allocated pool: set the initial_pool_size
    field to the number of maximum reference surfaces possible with the codec,
    plus 1 surface for the user to work (meaning the user can safely reference
    at most 1 decoded surface at a time), plus additional buffering introduced
    by frame threading. If the hwaccel does not require pre-allocation, the
    field is left to 0, and the decoder will allocate new surfaces on demand
    during decoding.
  - Possibly AVHWFramesContext.hwctx fields, depending on the underlying
    hardware API.

  Essentially, out_frames_ref returns the same as av_hwframe_ctx_alloc(), but
  with basic frame parameters set.

  The function is stateless, and does not change the AVCodecContext or the
  device_ref AVHWDeviceContext.

  @param avctx The context which is currently calling get_format, and which
               implicitly contains all state needed for filling the returned
               AVHWFramesContext properly.
  @param device_ref A reference to the AVHWDeviceContext describing the device
                    which will be used by the hardware decoder.
  @param hw_pix_fmt The hwaccel format you are going to return from get_format.
  @param out_frames_ref On success, set to a reference to an _uninitialized_
                        AVHWFramesContext, created from the given device_ref.
                        Fields will be set to values required for decoding.
                        Not changed if an error is returned.
  @return zero on success, a negative value on error. The following error codes
          have special semantics:
       AVERROR(ENOENT): the decoder does not support this functionality. Setup
                        is always manual, or it is a decoder which does not
                        support setting AVCodecContext.hw_frames_ctx at all,
                        or it is a software format.
       AVERROR(EINVAL): it is known that hardware decoding is not supported for
                        this configuration, or the device_ref is not supported
                        for the hwaccel referenced by hw_pix_fmt.
*/
func AVCodecGetHWFramesParameters(avctx *AVCodecContext, deviceRef *AVBufferRef, hwPixFmt AVPixelFormat, outFramesRef **AVBufferRef) (int, error) {
	var tmpavctx *C.AVCodecContext
	if avctx != nil {
		tmpavctx = avctx.ptr
	}
	var tmpdeviceRef *C.AVBufferRef
	if deviceRef != nil {
		tmpdeviceRef = deviceRef.ptr
	}
	var ptroutFramesRef **C.AVBufferRef
	var tmpoutFramesRef *C.AVBufferRef
	var oldTmpoutFramesRef *C.AVBufferRef
	if outFramesRef != nil {
		inneroutFramesRef := *outFramesRef
		if inneroutFramesRef != nil {
			tmpoutFramesRef = inneroutFramesRef.ptr
			oldTmpoutFramesRef = tmpoutFramesRef
		}
		ptroutFramesRef = &tmpoutFramesRef
	}
	ret := C.avcodec_get_hw_frames_parameters(tmpavctx, tmpdeviceRef, C.enum_AVPixelFormat(hwPixFmt), ptroutFramesRef)
	if tmpoutFramesRef != oldTmpoutFramesRef && outFramesRef != nil {
		if tmpoutFramesRef != nil {
			*outFramesRef = &AVBufferRef{ptr: tmpoutFramesRef}
		} else {
			*outFramesRef = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function avcodec_get_supported_config ---

// avcodec_get_supported_config skipped due to outConfigs

// --- Function av_parser_iterate ---

// av_parser_iterate skipped due to opaque

// --- Function av_parser_init ---

// AVParserInit wraps av_parser_init.
func AVParserInit(codecId int) *AVCodecParserContext {
	ret := C.av_parser_init(C.int(codecId))
	var retMapped *AVCodecParserContext
	if ret != nil {
		retMapped = &AVCodecParserContext{ptr: ret}
	}
	return retMapped
}

// --- Function av_parser_parse2 ---

// av_parser_parse2 skipped due to poutbuf

// --- Function av_parser_close ---

// AVParserClose wraps av_parser_close.
func AVParserClose(s *AVCodecParserContext) {
	var tmps *C.AVCodecParserContext
	if s != nil {
		tmps = s.ptr
	}
	C.av_parser_close(tmps)
}

// --- Function avcodec_encode_subtitle ---

// AVCodecEncodeSubtitle wraps avcodec_encode_subtitle.
func AVCodecEncodeSubtitle(avctx *AVCodecContext, buf unsafe.Pointer, bufSize int, sub *AVSubtitle) (int, error) {
	var tmpavctx *C.AVCodecContext
	if avctx != nil {
		tmpavctx = avctx.ptr
	}
	var tmpsub *C.AVSubtitle
	if sub != nil {
		tmpsub = sub.ptr
	}
	ret := C.avcodec_encode_subtitle(tmpavctx, (*C.uint8_t)(buf), C.int(bufSize), tmpsub)
	return int(ret), WrapErr(int(ret))
}

// --- Function avcodec_pix_fmt_to_codec_tag ---

// AVCodecPixFmtToCodecTag wraps avcodec_pix_fmt_to_codec_tag.
/*
  Return a value representing the fourCC code associated to the
  pixel format pix_fmt, or 0 if no associated fourCC code can be
  found.
*/
func AVCodecPixFmtToCodecTag(pixFmt AVPixelFormat) uint {
	ret := C.avcodec_pix_fmt_to_codec_tag(C.enum_AVPixelFormat(pixFmt))
	return uint(ret)
}

// --- Function avcodec_find_best_pix_fmt_of_list ---

// AVCodecFindBestPixFmtOfList wraps avcodec_find_best_pix_fmt_of_list.
/*
  Find the best pixel format to convert to given a certain source pixel
  format.  When converting from one pixel format to another, information loss
  may occur.  For example, when converting from RGB24 to GRAY, the color
  information will be lost. Similarly, other losses occur when converting from
  some formats to other formats. avcodec_find_best_pix_fmt_of_2() searches which of
  the given pixel formats should be used to suffer the least amount of loss.
  The pixel formats from which it chooses one, are determined by the
  pix_fmt_list parameter.


  @param[in] pix_fmt_list AV_PIX_FMT_NONE terminated array of pixel formats to choose from
  @param[in] src_pix_fmt source pixel format
  @param[in] has_alpha Whether the source pixel format alpha channel is used.
  @param[out] loss_ptr Combination of flags informing you what kind of losses will occur.
  @return The best pixel format to convert to or -1 if none was found.
*/
func AVCodecFindBestPixFmtOfList(pixFmtList *AVPixelFormat, srcPixFmt AVPixelFormat, hasAlpha int, lossPtr *int) AVPixelFormat {
	var tmppixFmtList *C.enum_AVPixelFormat
	if pixFmtList != nil {
		tmppixFmtList = (*C.enum_AVPixelFormat)(unsafe.Pointer(pixFmtList))
	}
	ret := C.avcodec_find_best_pix_fmt_of_list(tmppixFmtList, C.enum_AVPixelFormat(srcPixFmt), C.int(hasAlpha), (*C.int)(unsafe.Pointer(lossPtr)))
	return AVPixelFormat(ret)
}

// --- Function avcodec_default_get_format ---

// AVCodecDefaultGetFormat wraps avcodec_default_get_format.
func AVCodecDefaultGetFormat(s *AVCodecContext, fmt *AVPixelFormat) AVPixelFormat {
	var tmps *C.AVCodecContext
	if s != nil {
		tmps = s.ptr
	}
	var tmpfmt *C.enum_AVPixelFormat
	if fmt != nil {
		tmpfmt = (*C.enum_AVPixelFormat)(unsafe.Pointer(fmt))
	}
	ret := C.avcodec_default_get_format(tmps, tmpfmt)
	return AVPixelFormat(ret)
}

// --- Function avcodec_string ---

// AVCodecString wraps avcodec_string.
func AVCodecString(buf *CStr, bufSize int, enc *AVCodecContext, encode int) {
	var tmpbuf *C.char
	if buf != nil {
		tmpbuf = buf.ptr
	}
	var tmpenc *C.AVCodecContext
	if enc != nil {
		tmpenc = enc.ptr
	}
	C.avcodec_string(tmpbuf, C.int(bufSize), tmpenc, C.int(encode))
}

// --- Function avcodec_default_execute ---

// avcodec_default_execute skipped due to func.

// --- Function avcodec_default_execute2 ---

// avcodec_default_execute2 skipped due to func.

// --- Function avcodec_fill_audio_frame ---

// AVCodecFillAudioFrame wraps avcodec_fill_audio_frame.
/*
  Fill AVFrame audio data and linesize pointers.

  The buffer buf must be a preallocated buffer with a size big enough
  to contain the specified samples amount. The filled AVFrame data
  pointers will point to this buffer.

  AVFrame extended_data channel pointers are allocated if necessary for
  planar audio.

  @param frame       the AVFrame
                     frame->nb_samples must be set prior to calling the
                     function. This function fills in frame->data,
                     frame->extended_data, frame->linesize[0].
  @param nb_channels channel count
  @param sample_fmt  sample format
  @param buf         buffer to use for frame data
  @param buf_size    size of buffer
  @param align       plane size sample alignment (0 = default)
  @return            >=0 on success, negative error code on failure
  @todo return the size in bytes required to store the samples in
  case of success, at the next libavutil bump
*/
func AVCodecFillAudioFrame(frame *AVFrame, nbChannels int, sampleFmt AVSampleFormat, buf unsafe.Pointer, bufSize int, align int) (int, error) {
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.avcodec_fill_audio_frame(tmpframe, C.int(nbChannels), C.enum_AVSampleFormat(sampleFmt), (*C.uint8_t)(buf), C.int(bufSize), C.int(align))
	return int(ret), WrapErr(int(ret))
}

// --- Function avcodec_flush_buffers ---

// AVCodecFlushBuffers wraps avcodec_flush_buffers.
/*
  Reset the internal codec state / flush internal buffers. Should be called
  e.g. when seeking or when switching to a different stream.

  @note for decoders, this function just releases any references the decoder
  might keep internally, but the caller's references remain valid.

  @note for encoders, this function will only do something if the encoder
  declares support for AV_CODEC_CAP_ENCODER_FLUSH. When called, the encoder
  will drain any remaining packets, and can then be reused for a different
  stream (as opposed to sending a null frame which will leave the encoder
  in a permanent EOF state after draining). This can be desirable if the
  cost of tearing down and replacing the encoder instance is high.
*/
func AVCodecFlushBuffers(avctx *AVCodecContext) {
	var tmpavctx *C.AVCodecContext
	if avctx != nil {
		tmpavctx = avctx.ptr
	}
	C.avcodec_flush_buffers(tmpavctx)
}

// --- Function av_get_audio_frame_duration ---

// AVGetAudioFrameDuration wraps av_get_audio_frame_duration.
/*
  Return audio frame duration.

  @param avctx        codec context
  @param frame_bytes  size of the frame, or 0 if unknown
  @return             frame duration, in samples, if known. 0 if not able to
                      determine.
*/
func AVGetAudioFrameDuration(avctx *AVCodecContext, frameBytes int) (int, error) {
	var tmpavctx *C.AVCodecContext
	if avctx != nil {
		tmpavctx = avctx.ptr
	}
	ret := C.av_get_audio_frame_duration(tmpavctx, C.int(frameBytes))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_fast_padded_malloc ---

// AVFastPaddedMalloc wraps av_fast_padded_malloc.
/*
  Same behaviour av_fast_malloc but the buffer has additional
  AV_INPUT_BUFFER_PADDING_SIZE at the end which will always be 0.

  In addition the whole buffer will initially and after resizes
  be 0-initialized so that no uninitialized data will ever appear.
*/
func AVFastPaddedMalloc(ptr unsafe.Pointer, size *uint, minSize uint64) {
	C.av_fast_padded_malloc(ptr, (*C.uint)(unsafe.Pointer(size)), C.size_t(minSize))
}

// --- Function av_fast_padded_mallocz ---

// AVFastPaddedMallocz wraps av_fast_padded_mallocz.
/*
  Same behaviour av_fast_padded_malloc except that buffer will always
  be 0-initialized after call.
*/
func AVFastPaddedMallocz(ptr unsafe.Pointer, size *uint, minSize uint64) {
	C.av_fast_padded_mallocz(ptr, (*C.uint)(unsafe.Pointer(size)), C.size_t(minSize))
}

// --- Function avcodec_is_open ---

// AVCodecIsOpen wraps avcodec_is_open.
/*
  @return a positive value if s is open (i.e. avcodec_open2() was called on it),
  0 otherwise.
*/
func AVCodecIsOpen(s *AVCodecContext) (int, error) {
	var tmps *C.AVCodecContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.avcodec_is_open(tmps)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_bsf_get_by_name ---

// AVBsfGetByName wraps av_bsf_get_by_name.
/*
  @return a bitstream filter with the specified name or NULL if no such
          bitstream filter exists.
*/
func AVBsfGetByName(name *CStr) *AVBitStreamFilter {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_bsf_get_by_name(tmpname)
	var retMapped *AVBitStreamFilter
	if ret != nil {
		retMapped = &AVBitStreamFilter{ptr: ret}
	}
	return retMapped
}

// --- Function av_bsf_iterate ---

// av_bsf_iterate skipped due to opaque

// --- Function av_bsf_alloc ---

// AVBsfAlloc wraps av_bsf_alloc.
/*
  Allocate a context for a given bitstream filter. The caller must fill in the
  context parameters as described in the documentation and then call
  av_bsf_init() before sending any data to the filter.

  @param filter the filter for which to allocate an instance.
  @param[out] ctx a pointer into which the pointer to the newly-allocated context
                  will be written. It must be freed with av_bsf_free() after the
                  filtering is done.

  @return 0 on success, a negative AVERROR code on failure
*/
func AVBsfAlloc(filter *AVBitStreamFilter, ctx **AVBSFContext) (int, error) {
	var tmpfilter *C.AVBitStreamFilter
	if filter != nil {
		tmpfilter = filter.ptr
	}
	var ptrctx **C.AVBSFContext
	var tmpctx *C.AVBSFContext
	var oldTmpctx *C.AVBSFContext
	if ctx != nil {
		innerctx := *ctx
		if innerctx != nil {
			tmpctx = innerctx.ptr
			oldTmpctx = tmpctx
		}
		ptrctx = &tmpctx
	}
	ret := C.av_bsf_alloc(tmpfilter, ptrctx)
	if tmpctx != oldTmpctx && ctx != nil {
		if tmpctx != nil {
			*ctx = &AVBSFContext{ptr: tmpctx}
		} else {
			*ctx = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function av_bsf_init ---

// AVBsfInit wraps av_bsf_init.
/*
  Prepare the filter for use, after all the parameters and options have been
  set.

  @param ctx a AVBSFContext previously allocated with av_bsf_alloc()
*/
func AVBsfInit(ctx *AVBSFContext) (int, error) {
	var tmpctx *C.AVBSFContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	ret := C.av_bsf_init(tmpctx)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_bsf_send_packet ---

// AVBsfSendPacket wraps av_bsf_send_packet.
/*
  Submit a packet for filtering.

  After sending each packet, the filter must be completely drained by calling
  av_bsf_receive_packet() repeatedly until it returns AVERROR(EAGAIN) or
  AVERROR_EOF.

  @param ctx an initialized AVBSFContext
  @param pkt the packet to filter. The bitstream filter will take ownership of
  the packet and reset the contents of pkt. pkt is not touched if an error occurs.
  If pkt is empty (i.e. NULL, or pkt->data is NULL and pkt->side_data_elems zero),
  it signals the end of the stream (i.e. no more non-empty packets will be sent;
  sending more empty packets does nothing) and will cause the filter to output
  any packets it may have buffered internally.

  @return
   - 0 on success.
   - AVERROR(EAGAIN) if packets need to be retrieved from the filter (using
     av_bsf_receive_packet()) before new input can be consumed.
   - Another negative AVERROR value if an error occurs.
*/
func AVBsfSendPacket(ctx *AVBSFContext, pkt *AVPacket) (int, error) {
	var tmpctx *C.AVBSFContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	var tmppkt *C.AVPacket
	if pkt != nil {
		tmppkt = pkt.ptr
	}
	ret := C.av_bsf_send_packet(tmpctx, tmppkt)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_bsf_receive_packet ---

// AVBsfReceivePacket wraps av_bsf_receive_packet.
/*
  Retrieve a filtered packet.

  @param ctx an initialized AVBSFContext
  @param[out] pkt this struct will be filled with the contents of the filtered
                  packet. It is owned by the caller and must be freed using
                  av_packet_unref() when it is no longer needed.
                  This parameter should be "clean" (i.e. freshly allocated
                  with av_packet_alloc() or unreffed with av_packet_unref())
                  when this function is called. If this function returns
                  successfully, the contents of pkt will be completely
                  overwritten by the returned data. On failure, pkt is not
                  touched.

  @return
   - 0 on success.
   - AVERROR(EAGAIN) if more packets need to be sent to the filter (using
     av_bsf_send_packet()) to get more output.
   - AVERROR_EOF if there will be no further output from the filter.
   - Another negative AVERROR value if an error occurs.

  @note one input packet may result in several output packets, so after sending
  a packet with av_bsf_send_packet(), this function needs to be called
  repeatedly until it stops returning 0. It is also possible for a filter to
  output fewer packets than were sent to it, so this function may return
  AVERROR(EAGAIN) immediately after a successful av_bsf_send_packet() call.
*/
func AVBsfReceivePacket(ctx *AVBSFContext, pkt *AVPacket) (int, error) {
	var tmpctx *C.AVBSFContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	var tmppkt *C.AVPacket
	if pkt != nil {
		tmppkt = pkt.ptr
	}
	ret := C.av_bsf_receive_packet(tmpctx, tmppkt)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_bsf_flush ---

// AVBsfFlush wraps av_bsf_flush.
//
//	Reset the internal bitstream filter state. Should be called e.g. when seeking.
func AVBsfFlush(ctx *AVBSFContext) {
	var tmpctx *C.AVBSFContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	C.av_bsf_flush(tmpctx)
}

// --- Function av_bsf_free ---

// AVBsfFree wraps av_bsf_free.
/*
  Free a bitstream filter context and everything associated with it; write NULL
  into the supplied pointer.
*/
func AVBsfFree(ctx **AVBSFContext) {
	var ptrctx **C.AVBSFContext
	var tmpctx *C.AVBSFContext
	var oldTmpctx *C.AVBSFContext
	if ctx != nil {
		innerctx := *ctx
		if innerctx != nil {
			tmpctx = innerctx.ptr
			oldTmpctx = tmpctx
		}
		ptrctx = &tmpctx
	}
	C.av_bsf_free(ptrctx)
	if tmpctx != oldTmpctx && ctx != nil {
		if tmpctx != nil {
			*ctx = &AVBSFContext{ptr: tmpctx}
		} else {
			*ctx = nil
		}
	}
}

// --- Function av_bsf_get_class ---

// AVBsfGetClass wraps av_bsf_get_class.
/*
  Get the AVClass for AVBSFContext. It can be used in combination with
  AV_OPT_SEARCH_FAKE_OBJ for examining options.

  @see av_opt_find().
*/
func AVBsfGetClass() *AVClass {
	ret := C.av_bsf_get_class()
	var retMapped *AVClass
	if ret != nil {
		retMapped = &AVClass{ptr: ret}
	}
	return retMapped
}

// --- Function av_bsf_list_alloc ---

// AVBsfListAlloc wraps av_bsf_list_alloc.
/*
  Allocate empty list of bitstream filters.
  The list must be later freed by av_bsf_list_free()
  or finalized by av_bsf_list_finalize().

  @return Pointer to @ref AVBSFList on success, NULL in case of failure
*/
func AVBsfListAlloc() *AVBSFList {
	ret := C.av_bsf_list_alloc()
	var retMapped *AVBSFList
	if ret != nil {
		retMapped = &AVBSFList{ptr: ret}
	}
	return retMapped
}

// --- Function av_bsf_list_free ---

// AVBsfListFree wraps av_bsf_list_free.
/*
  Free list of bitstream filters.

  @param lst Pointer to pointer returned by av_bsf_list_alloc()
*/
func AVBsfListFree(lst **AVBSFList) {
	var ptrlst **C.AVBSFList
	var tmplst *C.AVBSFList
	var oldTmplst *C.AVBSFList
	if lst != nil {
		innerlst := *lst
		if innerlst != nil {
			tmplst = innerlst.ptr
			oldTmplst = tmplst
		}
		ptrlst = &tmplst
	}
	C.av_bsf_list_free(ptrlst)
	if tmplst != oldTmplst && lst != nil {
		if tmplst != nil {
			*lst = &AVBSFList{ptr: tmplst}
		} else {
			*lst = nil
		}
	}
}

// --- Function av_bsf_list_append ---

// AVBsfListAppend wraps av_bsf_list_append.
/*
  Append bitstream filter to the list of bitstream filters.

  @param lst List to append to
  @param bsf Filter context to be appended

  @return >=0 on success, negative AVERROR in case of failure
*/
func AVBsfListAppend(lst *AVBSFList, bsf *AVBSFContext) (int, error) {
	var tmplst *C.AVBSFList
	if lst != nil {
		tmplst = lst.ptr
	}
	var tmpbsf *C.AVBSFContext
	if bsf != nil {
		tmpbsf = bsf.ptr
	}
	ret := C.av_bsf_list_append(tmplst, tmpbsf)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_bsf_list_append2 ---

// AVBsfListAppend2 wraps av_bsf_list_append2.
/*
  Construct new bitstream filter context given it's name and options
  and append it to the list of bitstream filters.

  @param lst      List to append to
  @param bsf_name Name of the bitstream filter
  @param options  Options for the bitstream filter, can be set to NULL

  @return >=0 on success, negative AVERROR in case of failure
*/
func AVBsfListAppend2(lst *AVBSFList, bsfName *CStr, options **AVDictionary) (int, error) {
	var tmplst *C.AVBSFList
	if lst != nil {
		tmplst = lst.ptr
	}
	var tmpbsfName *C.char
	if bsfName != nil {
		tmpbsfName = bsfName.ptr
	}
	var ptroptions **C.AVDictionary
	var tmpoptions *C.AVDictionary
	var oldTmpoptions *C.AVDictionary
	if options != nil {
		inneroptions := *options
		if inneroptions != nil {
			tmpoptions = inneroptions.ptr
			oldTmpoptions = tmpoptions
		}
		ptroptions = &tmpoptions
	}
	ret := C.av_bsf_list_append2(tmplst, tmpbsfName, ptroptions)
	if tmpoptions != oldTmpoptions && options != nil {
		if tmpoptions != nil {
			*options = &AVDictionary{ptr: tmpoptions}
		} else {
			*options = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function av_bsf_list_finalize ---

// AVBsfListFinalize wraps av_bsf_list_finalize.
/*
  Finalize list of bitstream filters.

  This function will transform @ref AVBSFList to single @ref AVBSFContext,
  so the whole chain of bitstream filters can be treated as single filter
  freshly allocated by av_bsf_alloc().
  If the call is successful, @ref AVBSFList structure is freed and lst
  will be set to NULL. In case of failure, caller is responsible for
  freeing the structure by av_bsf_list_free()

  @param      lst Filter list structure to be transformed
  @param[out] bsf Pointer to be set to newly created @ref AVBSFContext structure
                  representing the chain of bitstream filters

  @return >=0 on success, negative AVERROR in case of failure
*/
func AVBsfListFinalize(lst **AVBSFList, bsf **AVBSFContext) (int, error) {
	var ptrlst **C.AVBSFList
	var tmplst *C.AVBSFList
	var oldTmplst *C.AVBSFList
	if lst != nil {
		innerlst := *lst
		if innerlst != nil {
			tmplst = innerlst.ptr
			oldTmplst = tmplst
		}
		ptrlst = &tmplst
	}
	var ptrbsf **C.AVBSFContext
	var tmpbsf *C.AVBSFContext
	var oldTmpbsf *C.AVBSFContext
	if bsf != nil {
		innerbsf := *bsf
		if innerbsf != nil {
			tmpbsf = innerbsf.ptr
			oldTmpbsf = tmpbsf
		}
		ptrbsf = &tmpbsf
	}
	ret := C.av_bsf_list_finalize(ptrlst, ptrbsf)
	if tmplst != oldTmplst && lst != nil {
		if tmplst != nil {
			*lst = &AVBSFList{ptr: tmplst}
		} else {
			*lst = nil
		}
	}
	if tmpbsf != oldTmpbsf && bsf != nil {
		if tmpbsf != nil {
			*bsf = &AVBSFContext{ptr: tmpbsf}
		} else {
			*bsf = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function av_bsf_list_parse_str ---

// AVBsfListParseStr wraps av_bsf_list_parse_str.
/*
  Parse string describing list of bitstream filters and create single
  @ref AVBSFContext describing the whole chain of bitstream filters.
  Resulting @ref AVBSFContext can be treated as any other @ref AVBSFContext freshly
  allocated by av_bsf_alloc().

  @param      str String describing chain of bitstream filters in format
                  `bsf1[=opt1=val1:opt2=val2][,bsf2]`
  @param[out] bsf Pointer to be set to newly created @ref AVBSFContext structure
                  representing the chain of bitstream filters

  @return >=0 on success, negative AVERROR in case of failure
*/
func AVBsfListParseStr(str *CStr, bsf **AVBSFContext) (int, error) {
	var tmpstr *C.char
	if str != nil {
		tmpstr = str.ptr
	}
	var ptrbsf **C.AVBSFContext
	var tmpbsf *C.AVBSFContext
	var oldTmpbsf *C.AVBSFContext
	if bsf != nil {
		innerbsf := *bsf
		if innerbsf != nil {
			tmpbsf = innerbsf.ptr
			oldTmpbsf = tmpbsf
		}
		ptrbsf = &tmpbsf
	}
	ret := C.av_bsf_list_parse_str(tmpstr, ptrbsf)
	if tmpbsf != oldTmpbsf && bsf != nil {
		if tmpbsf != nil {
			*bsf = &AVBSFContext{ptr: tmpbsf}
		} else {
			*bsf = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function av_bsf_get_null_filter ---

// AVBsfGetNullFilter wraps av_bsf_get_null_filter.
/*
  Get null/pass-through bitstream filter.

  @param[out] bsf Pointer to be set to new instance of pass-through bitstream filter

  @return
*/
func AVBsfGetNullFilter(bsf **AVBSFContext) (int, error) {
	var ptrbsf **C.AVBSFContext
	var tmpbsf *C.AVBSFContext
	var oldTmpbsf *C.AVBSFContext
	if bsf != nil {
		innerbsf := *bsf
		if innerbsf != nil {
			tmpbsf = innerbsf.ptr
			oldTmpbsf = tmpbsf
		}
		ptrbsf = &tmpbsf
	}
	ret := C.av_bsf_get_null_filter(ptrbsf)
	if tmpbsf != oldTmpbsf && bsf != nil {
		if tmpbsf != nil {
			*bsf = &AVBSFContext{ptr: tmpbsf}
		} else {
			*bsf = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function av_codec_iterate ---

// av_codec_iterate skipped due to opaque

// --- Function avcodec_find_decoder ---

// AVCodecFindDecoder wraps avcodec_find_decoder.
/*
  Find a registered decoder with a matching codec ID.

  @param id AVCodecID of the requested decoder
  @return A decoder if one was found, NULL otherwise.
*/
func AVCodecFindDecoder(id AVCodecID) *AVCodec {
	ret := C.avcodec_find_decoder(C.enum_AVCodecID(id))
	var retMapped *AVCodec
	if ret != nil {
		retMapped = &AVCodec{ptr: ret}
	}
	return retMapped
}

// --- Function avcodec_find_decoder_by_name ---

// AVCodecFindDecoderByName wraps avcodec_find_decoder_by_name.
/*
  Find a registered decoder with the specified name.

  @param name name of the requested decoder
  @return A decoder if one was found, NULL otherwise.
*/
func AVCodecFindDecoderByName(name *CStr) *AVCodec {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.avcodec_find_decoder_by_name(tmpname)
	var retMapped *AVCodec
	if ret != nil {
		retMapped = &AVCodec{ptr: ret}
	}
	return retMapped
}

// --- Function avcodec_find_encoder ---

// AVCodecFindEncoder wraps avcodec_find_encoder.
/*
  Find a registered encoder with a matching codec ID.

  @param id AVCodecID of the requested encoder
  @return An encoder if one was found, NULL otherwise.
*/
func AVCodecFindEncoder(id AVCodecID) *AVCodec {
	ret := C.avcodec_find_encoder(C.enum_AVCodecID(id))
	var retMapped *AVCodec
	if ret != nil {
		retMapped = &AVCodec{ptr: ret}
	}
	return retMapped
}

// --- Function avcodec_find_encoder_by_name ---

// AVCodecFindEncoderByName wraps avcodec_find_encoder_by_name.
/*
  Find a registered encoder with the specified name.

  @param name name of the requested encoder
  @return An encoder if one was found, NULL otherwise.
*/
func AVCodecFindEncoderByName(name *CStr) *AVCodec {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.avcodec_find_encoder_by_name(tmpname)
	var retMapped *AVCodec
	if ret != nil {
		retMapped = &AVCodec{ptr: ret}
	}
	return retMapped
}

// --- Function av_codec_is_encoder ---

// AVCodecIsEncoder wraps av_codec_is_encoder.
//
//	@return a non-zero number if codec is an encoder, zero otherwise
func AVCodecIsEncoder(codec *AVCodec) (int, error) {
	var tmpcodec *C.AVCodec
	if codec != nil {
		tmpcodec = codec.ptr
	}
	ret := C.av_codec_is_encoder(tmpcodec)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_codec_is_decoder ---

// AVCodecIsDecoder wraps av_codec_is_decoder.
//
//	@return a non-zero number if codec is a decoder, zero otherwise
func AVCodecIsDecoder(codec *AVCodec) (int, error) {
	var tmpcodec *C.AVCodec
	if codec != nil {
		tmpcodec = codec.ptr
	}
	ret := C.av_codec_is_decoder(tmpcodec)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_get_profile_name ---

// AVGetProfileName wraps av_get_profile_name.
/*
  Return a name for the specified profile, if available.

  @param codec the codec that is searched for the given profile
  @param profile the profile value for which a name is requested
  @return A name for the profile if found, NULL otherwise.
*/
func AVGetProfileName(codec *AVCodec, profile int) *CStr {
	var tmpcodec *C.AVCodec
	if codec != nil {
		tmpcodec = codec.ptr
	}
	ret := C.av_get_profile_name(tmpcodec, C.int(profile))
	return wrapCStr(ret)
}

// --- Function avcodec_get_hw_config ---

// AVCodecGetHWConfig wraps avcodec_get_hw_config.
/*
  Retrieve supported hardware configurations for a codec.

  Values of index from zero to some maximum return the indexed configuration
  descriptor; all other values return NULL.  If the codec does not support
  any hardware configurations then it will always return NULL.
*/
func AVCodecGetHWConfig(codec *AVCodec, index int) *AVCodecHWConfig {
	var tmpcodec *C.AVCodec
	if codec != nil {
		tmpcodec = codec.ptr
	}
	ret := C.avcodec_get_hw_config(tmpcodec, C.int(index))
	var retMapped *AVCodecHWConfig
	if ret != nil {
		retMapped = &AVCodecHWConfig{ptr: ret}
	}
	return retMapped
}

// --- Function avcodec_descriptor_get ---

// AVCodecDescriptorGet wraps avcodec_descriptor_get.
//
//	@return descriptor for given codec ID or NULL if no descriptor exists.
func AVCodecDescriptorGet(id AVCodecID) *AVCodecDescriptor {
	ret := C.avcodec_descriptor_get(C.enum_AVCodecID(id))
	var retMapped *AVCodecDescriptor
	if ret != nil {
		retMapped = &AVCodecDescriptor{ptr: ret}
	}
	return retMapped
}

// --- Function avcodec_descriptor_next ---

// AVCodecDescriptorNext wraps avcodec_descriptor_next.
/*
  Iterate over all codec descriptors known to libavcodec.

  @param prev previous descriptor. NULL to get the first descriptor.

  @return next descriptor or NULL after the last descriptor
*/
func AVCodecDescriptorNext(prev *AVCodecDescriptor) *AVCodecDescriptor {
	var tmpprev *C.AVCodecDescriptor
	if prev != nil {
		tmpprev = prev.ptr
	}
	ret := C.avcodec_descriptor_next(tmpprev)
	var retMapped *AVCodecDescriptor
	if ret != nil {
		retMapped = &AVCodecDescriptor{ptr: ret}
	}
	return retMapped
}

// --- Function avcodec_descriptor_get_by_name ---

// AVCodecDescriptorGetByName wraps avcodec_descriptor_get_by_name.
/*
  @return codec descriptor with the given name or NULL if no such descriptor
          exists.
*/
func AVCodecDescriptorGetByName(name *CStr) *AVCodecDescriptor {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.avcodec_descriptor_get_by_name(tmpname)
	var retMapped *AVCodecDescriptor
	if ret != nil {
		retMapped = &AVCodecDescriptor{ptr: ret}
	}
	return retMapped
}

// --- Function avcodec_get_type ---

// AVCodecGetType wraps avcodec_get_type.
//
//	Get the type of the given codec.
func AVCodecGetType(codecId AVCodecID) AVMediaType {
	ret := C.avcodec_get_type(C.enum_AVCodecID(codecId))
	return AVMediaType(ret)
}

// --- Function avcodec_get_name ---

// AVCodecGetName wraps avcodec_get_name.
/*
  Get the name of a codec.
  @return  a static string identifying the codec; never NULL
*/
func AVCodecGetName(id AVCodecID) *CStr {
	ret := C.avcodec_get_name(C.enum_AVCodecID(id))
	return wrapCStr(ret)
}

// --- Function av_get_bits_per_sample ---

// AVGetBitsPerSample wraps av_get_bits_per_sample.
/*
  Return codec bits per sample.

  @param[in] codec_id the codec
  @return Number of bits per sample or zero if unknown for the given codec.
*/
func AVGetBitsPerSample(codecId AVCodecID) (int, error) {
	ret := C.av_get_bits_per_sample(C.enum_AVCodecID(codecId))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_get_exact_bits_per_sample ---

// AVGetExactBitsPerSample wraps av_get_exact_bits_per_sample.
/*
  Return codec bits per sample.
  Only return non-zero if the bits per sample is exactly correct, not an
  approximation.

  @param[in] codec_id the codec
  @return Number of bits per sample or zero if unknown for the given codec.
*/
func AVGetExactBitsPerSample(codecId AVCodecID) (int, error) {
	ret := C.av_get_exact_bits_per_sample(C.enum_AVCodecID(codecId))
	return int(ret), WrapErr(int(ret))
}

// --- Function avcodec_profile_name ---

// AVCodecProfileName wraps avcodec_profile_name.
/*
  Return a name for the specified profile, if available.

  @param codec_id the ID of the codec to which the requested profile belongs
  @param profile the profile value for which a name is requested
  @return A name for the profile if found, NULL otherwise.

  @note unlike av_get_profile_name(), which searches a list of profiles
        supported by a specific decoder or encoder implementation, this
        function searches the list of profiles from the AVCodecDescriptor
*/
func AVCodecProfileName(codecId AVCodecID, profile int) *CStr {
	ret := C.avcodec_profile_name(C.enum_AVCodecID(codecId), C.int(profile))
	return wrapCStr(ret)
}

// --- Function av_get_pcm_codec ---

// AVGetPcmCodec wraps av_get_pcm_codec.
/*
  Return the PCM codec associated with a sample format.
  @param be  endianness, 0 for little, 1 for big,
             -1 (or anything else) for native
  @return  AV_CODEC_ID_PCM_* or AV_CODEC_ID_NONE
*/
func AVGetPcmCodec(fmt AVSampleFormat, be int) AVCodecID {
	ret := C.av_get_pcm_codec(C.enum_AVSampleFormat(fmt), C.int(be))
	return AVCodecID(ret)
}

// --- Function avcodec_parameters_alloc ---

// AVCodecParametersAlloc wraps avcodec_parameters_alloc.
/*
  Allocate a new AVCodecParameters and set its fields to default values
  (unknown/invalid/0). The returned struct must be freed with
  avcodec_parameters_free().
*/
func AVCodecParametersAlloc() *AVCodecParameters {
	ret := C.avcodec_parameters_alloc()
	var retMapped *AVCodecParameters
	if ret != nil {
		retMapped = &AVCodecParameters{ptr: ret}
	}
	return retMapped
}

// --- Function avcodec_parameters_free ---

// AVCodecParametersFree wraps avcodec_parameters_free.
/*
  Free an AVCodecParameters instance and everything associated with it and
  write NULL to the supplied pointer.
*/
func AVCodecParametersFree(par **AVCodecParameters) {
	var ptrpar **C.AVCodecParameters
	var tmppar *C.AVCodecParameters
	var oldTmppar *C.AVCodecParameters
	if par != nil {
		innerpar := *par
		if innerpar != nil {
			tmppar = innerpar.ptr
			oldTmppar = tmppar
		}
		ptrpar = &tmppar
	}
	C.avcodec_parameters_free(ptrpar)
	if tmppar != oldTmppar && par != nil {
		if tmppar != nil {
			*par = &AVCodecParameters{ptr: tmppar}
		} else {
			*par = nil
		}
	}
}

// --- Function avcodec_parameters_copy ---

// AVCodecParametersCopy wraps avcodec_parameters_copy.
/*
  Copy the contents of src to dst. Any allocated fields in dst are freed and
  replaced with newly allocated duplicates of the corresponding fields in src.

  @return >= 0 on success, a negative AVERROR code on failure.
*/
func AVCodecParametersCopy(dst *AVCodecParameters, src *AVCodecParameters) (int, error) {
	var tmpdst *C.AVCodecParameters
	if dst != nil {
		tmpdst = dst.ptr
	}
	var tmpsrc *C.AVCodecParameters
	if src != nil {
		tmpsrc = src.ptr
	}
	ret := C.avcodec_parameters_copy(tmpdst, tmpsrc)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_get_audio_frame_duration2 ---

// AVGetAudioFrameDuration2 wraps av_get_audio_frame_duration2.
/*
  This function is the same as av_get_audio_frame_duration(), except it works
  with AVCodecParameters instead of an AVCodecContext.
*/
func AVGetAudioFrameDuration2(par *AVCodecParameters, frameBytes int) (int, error) {
	var tmppar *C.AVCodecParameters
	if par != nil {
		tmppar = par.ptr
	}
	ret := C.av_get_audio_frame_duration2(tmppar, C.int(frameBytes))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_cpb_properties_alloc ---

// AVCpbPropertiesAlloc wraps av_cpb_properties_alloc.
/*
  Allocate a CPB properties structure and initialize its fields to default
  values.

  @param size if non-NULL, the size of the allocated struct will be written
              here. This is useful for embedding it in side data.

  @return the newly allocated struct or NULL on failure
*/
func AVCpbPropertiesAlloc(size *uint64) *AVCPBProperties {
	ret := C.av_cpb_properties_alloc((*C.size_t)(unsafe.Pointer(size)))
	var retMapped *AVCPBProperties
	if ret != nil {
		retMapped = &AVCPBProperties{ptr: ret}
	}
	return retMapped
}

// --- Function av_xiphlacing ---

// AVXiphlacing wraps av_xiphlacing.
/*
  Encode extradata length to a buffer. Used by xiph codecs.

  @param s buffer to write to; must be at least (v/255+1) bytes long
  @param v size of extradata in bytes
  @return number of bytes written to the buffer.
*/
func AVXiphlacing(s unsafe.Pointer, v uint) uint {
	ret := C.av_xiphlacing((*C.uchar)(s), C.uint(v))
	return uint(ret)
}

// --- Function av_packet_side_data_new ---

// AVPacketSideDataNew wraps av_packet_side_data_new.
/*
  Allocate a new packet side data.

  @param sd    pointer to an array of side data to which the side data should
               be added. *sd may be NULL, in which case the array will be
               initialized.
  @param nb_sd pointer to an integer containing the number of entries in
               the array. The integer value will be increased by 1 on success.
  @param type  side data type
  @param size  desired side data size
  @param flags currently unused. Must be zero

  @return pointer to freshly allocated side data on success, or NULL otherwise.
*/
func AVPacketSideDataNew(psd **AVPacketSideData, pnbSd *int, _type AVPacketSideDataType, size uint64, flags int) *AVPacketSideData {
	var ptrpsd **C.AVPacketSideData
	var tmppsd *C.AVPacketSideData
	var oldTmppsd *C.AVPacketSideData
	if psd != nil {
		innerpsd := *psd
		if innerpsd != nil {
			tmppsd = innerpsd.ptr
			oldTmppsd = tmppsd
		}
		ptrpsd = &tmppsd
	}
	ret := C.av_packet_side_data_new(ptrpsd, (*C.int)(unsafe.Pointer(pnbSd)), C.enum_AVPacketSideDataType(_type), C.size_t(size), C.int(flags))
	if tmppsd != oldTmppsd && psd != nil {
		if tmppsd != nil {
			*psd = &AVPacketSideData{ptr: tmppsd}
		} else {
			*psd = nil
		}
	}
	var retMapped *AVPacketSideData
	if ret != nil {
		retMapped = &AVPacketSideData{ptr: ret}
	}
	return retMapped
}

// --- Function av_packet_side_data_add ---

// AVPacketSideDataAdd wraps av_packet_side_data_add.
/*
  Wrap existing data as packet side data.

  @param sd    pointer to an array of side data to which the side data should
               be added. *sd may be NULL, in which case the array will be
               initialized
  @param nb_sd pointer to an integer containing the number of entries in
               the array. The integer value will be increased by 1 on success.
  @param type  side data type
  @param data  a data array. It must be allocated with the av_malloc() family
               of functions. The ownership of the data is transferred to the
               side data array on success
  @param size  size of the data array
  @param flags currently unused. Must be zero

  @return pointer to freshly allocated side data on success, or NULL otherwise
          On failure, the side data array is unchanged and the data remains
          owned by the caller.
*/
func AVPacketSideDataAdd(sd **AVPacketSideData, nbSd *int, _type AVPacketSideDataType, data unsafe.Pointer, size uint64, flags int) *AVPacketSideData {
	var ptrsd **C.AVPacketSideData
	var tmpsd *C.AVPacketSideData
	var oldTmpsd *C.AVPacketSideData
	if sd != nil {
		innersd := *sd
		if innersd != nil {
			tmpsd = innersd.ptr
			oldTmpsd = tmpsd
		}
		ptrsd = &tmpsd
	}
	ret := C.av_packet_side_data_add(ptrsd, (*C.int)(unsafe.Pointer(nbSd)), C.enum_AVPacketSideDataType(_type), data, C.size_t(size), C.int(flags))
	if tmpsd != oldTmpsd && sd != nil {
		if tmpsd != nil {
			*sd = &AVPacketSideData{ptr: tmpsd}
		} else {
			*sd = nil
		}
	}
	var retMapped *AVPacketSideData
	if ret != nil {
		retMapped = &AVPacketSideData{ptr: ret}
	}
	return retMapped
}

// --- Function av_packet_side_data_get ---

// AVPacketSideDataGet wraps av_packet_side_data_get.
/*
  Get side information from a side data array.

  @param sd    the array from which the side data should be fetched
  @param nb_sd value containing the number of entries in the array.
  @param type  desired side information type

  @return pointer to side data if present or NULL otherwise
*/
func AVPacketSideDataGet(sd *AVPacketSideData, nbSd int, _type AVPacketSideDataType) *AVPacketSideData {
	var tmpsd *C.AVPacketSideData
	if sd != nil {
		tmpsd = sd.ptr
	}
	ret := C.av_packet_side_data_get(tmpsd, C.int(nbSd), C.enum_AVPacketSideDataType(_type))
	var retMapped *AVPacketSideData
	if ret != nil {
		retMapped = &AVPacketSideData{ptr: ret}
	}
	return retMapped
}

// --- Function av_packet_side_data_remove ---

// AVPacketSideDataRemove wraps av_packet_side_data_remove.
/*
  Remove side data of the given type from a side data array.

  @param sd    the array from which the side data should be removed
  @param nb_sd pointer to an integer containing the number of entries in
               the array. Will be reduced by the amount of entries removed
               upon return
  @param type  side information type
*/
func AVPacketSideDataRemove(sd *AVPacketSideData, nbSd *int, _type AVPacketSideDataType) {
	var tmpsd *C.AVPacketSideData
	if sd != nil {
		tmpsd = sd.ptr
	}
	C.av_packet_side_data_remove(tmpsd, (*C.int)(unsafe.Pointer(nbSd)), C.enum_AVPacketSideDataType(_type))
}

// --- Function av_packet_side_data_free ---

// AVPacketSideDataFree wraps av_packet_side_data_free.
/*
  Convenience function to free all the side data stored in an array, and
  the array itself.

  @param sd    pointer to array of side data to free. Will be set to NULL
               upon return.
  @param nb_sd pointer to an integer containing the number of entries in
               the array. Will be set to 0 upon return.
*/
func AVPacketSideDataFree(sd **AVPacketSideData, nbSd *int) {
	var ptrsd **C.AVPacketSideData
	var tmpsd *C.AVPacketSideData
	var oldTmpsd *C.AVPacketSideData
	if sd != nil {
		innersd := *sd
		if innersd != nil {
			tmpsd = innersd.ptr
			oldTmpsd = tmpsd
		}
		ptrsd = &tmpsd
	}
	C.av_packet_side_data_free(ptrsd, (*C.int)(unsafe.Pointer(nbSd)))
	if tmpsd != oldTmpsd && sd != nil {
		if tmpsd != nil {
			*sd = &AVPacketSideData{ptr: tmpsd}
		} else {
			*sd = nil
		}
	}
}

// --- Function av_packet_side_data_name ---

// AVPacketSideDataName wraps av_packet_side_data_name.
func AVPacketSideDataName(_type AVPacketSideDataType) *CStr {
	ret := C.av_packet_side_data_name(C.enum_AVPacketSideDataType(_type))
	return wrapCStr(ret)
}

// --- Function av_packet_alloc ---

// AVPacketAlloc wraps av_packet_alloc.
/*
  Allocate an AVPacket and set its fields to default values.  The resulting
  struct must be freed using av_packet_free().

  @return An AVPacket filled with default values or NULL on failure.

  @note this only allocates the AVPacket itself, not the data buffers. Those
  must be allocated through other means such as av_new_packet.

  @see av_new_packet
*/
func AVPacketAlloc() *AVPacket {
	ret := C.av_packet_alloc()
	var retMapped *AVPacket
	if ret != nil {
		retMapped = &AVPacket{ptr: ret}
	}
	return retMapped
}

// --- Function av_packet_clone ---

// AVPacketClone wraps av_packet_clone.
/*
  Create a new packet that references the same data as src.

  This is a shortcut for av_packet_alloc()+av_packet_ref().

  @return newly created AVPacket on success, NULL on error.

  @see av_packet_alloc
  @see av_packet_ref
*/
func AVPacketClone(src *AVPacket) *AVPacket {
	var tmpsrc *C.AVPacket
	if src != nil {
		tmpsrc = src.ptr
	}
	ret := C.av_packet_clone(tmpsrc)
	var retMapped *AVPacket
	if ret != nil {
		retMapped = &AVPacket{ptr: ret}
	}
	return retMapped
}

// --- Function av_packet_free ---

// AVPacketFree wraps av_packet_free.
/*
  Free the packet, if the packet is reference counted, it will be
  unreferenced first.

  @param pkt packet to be freed. The pointer will be set to NULL.
  @note passing NULL is a no-op.
*/
func AVPacketFree(pkt **AVPacket) {
	var ptrpkt **C.AVPacket
	var tmppkt *C.AVPacket
	var oldTmppkt *C.AVPacket
	if pkt != nil {
		innerpkt := *pkt
		if innerpkt != nil {
			tmppkt = innerpkt.ptr
			oldTmppkt = tmppkt
		}
		ptrpkt = &tmppkt
	}
	C.av_packet_free(ptrpkt)
	if tmppkt != oldTmppkt && pkt != nil {
		if tmppkt != nil {
			*pkt = &AVPacket{ptr: tmppkt}
		} else {
			*pkt = nil
		}
	}
}

// --- Function av_init_packet ---

// AVInitPacket wraps av_init_packet.
/*
  Initialize optional fields of a packet with default values.

  Note, this does not touch the data and size members, which have to be
  initialized separately.

  @param pkt packet

  @see av_packet_alloc
  @see av_packet_unref

  @deprecated This function is deprecated. Once it's removed,
  sizeof(AVPacket) will not be a part of the ABI anymore.
*/
func AVInitPacket(pkt *AVPacket) {
	var tmppkt *C.AVPacket
	if pkt != nil {
		tmppkt = pkt.ptr
	}
	C.av_init_packet(tmppkt)
}

// --- Function av_new_packet ---

// AVNewPacket wraps av_new_packet.
/*
  Allocate the payload of a packet and initialize its fields with
  default values.

  @param pkt packet
  @param size wanted payload size
  @return 0 if OK, AVERROR_xxx otherwise
*/
func AVNewPacket(pkt *AVPacket, size int) (int, error) {
	var tmppkt *C.AVPacket
	if pkt != nil {
		tmppkt = pkt.ptr
	}
	ret := C.av_new_packet(tmppkt, C.int(size))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_shrink_packet ---

// AVShrinkPacket wraps av_shrink_packet.
/*
  Reduce packet size, correctly zeroing padding

  @param pkt packet
  @param size new size
*/
func AVShrinkPacket(pkt *AVPacket, size int) {
	var tmppkt *C.AVPacket
	if pkt != nil {
		tmppkt = pkt.ptr
	}
	C.av_shrink_packet(tmppkt, C.int(size))
}

// --- Function av_grow_packet ---

// AVGrowPacket wraps av_grow_packet.
/*
  Increase packet size, correctly zeroing padding

  @param pkt packet
  @param grow_by number of bytes by which to increase the size of the packet
*/
func AVGrowPacket(pkt *AVPacket, growBy int) (int, error) {
	var tmppkt *C.AVPacket
	if pkt != nil {
		tmppkt = pkt.ptr
	}
	ret := C.av_grow_packet(tmppkt, C.int(growBy))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_packet_from_data ---

// AVPacketFromData wraps av_packet_from_data.
/*
  Initialize a reference-counted packet from av_malloc()ed data.

  @param pkt packet to be initialized. This function will set the data, size,
         and buf fields, all others are left untouched.
  @param data Data allocated by av_malloc() to be used as packet data. If this
         function returns successfully, the data is owned by the underlying AVBuffer.
         The caller may not access the data through other means.
  @param size size of data in bytes, without the padding. I.e. the full buffer
         size is assumed to be size + AV_INPUT_BUFFER_PADDING_SIZE.

  @return 0 on success, a negative AVERROR on error
*/
func AVPacketFromData(pkt *AVPacket, data unsafe.Pointer, size int) (int, error) {
	var tmppkt *C.AVPacket
	if pkt != nil {
		tmppkt = pkt.ptr
	}
	ret := C.av_packet_from_data(tmppkt, (*C.uint8_t)(data), C.int(size))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_packet_new_side_data ---

// AVPacketNewSideData wraps av_packet_new_side_data.
/*
  Allocate new information of a packet.

  @param pkt packet
  @param type side information type
  @param size side information size
  @return pointer to fresh allocated data or NULL otherwise
*/
func AVPacketNewSideData(pkt *AVPacket, _type AVPacketSideDataType, size uint64) unsafe.Pointer {
	var tmppkt *C.AVPacket
	if pkt != nil {
		tmppkt = pkt.ptr
	}
	ret := C.av_packet_new_side_data(tmppkt, C.enum_AVPacketSideDataType(_type), C.size_t(size))
	return unsafe.Pointer(ret)
}

// --- Function av_packet_add_side_data ---

// AVPacketAddSideData wraps av_packet_add_side_data.
/*
  Wrap an existing array as a packet side data.

  @param pkt packet
  @param type side information type
  @param data the side data array. It must be allocated with the av_malloc()
              family of functions. The ownership of the data is transferred to
              pkt.
  @param size side information size
  @return a non-negative number on success, a negative AVERROR code on
          failure. On failure, the packet is unchanged and the data remains
          owned by the caller.
*/
func AVPacketAddSideData(pkt *AVPacket, _type AVPacketSideDataType, data unsafe.Pointer, size uint64) (int, error) {
	var tmppkt *C.AVPacket
	if pkt != nil {
		tmppkt = pkt.ptr
	}
	ret := C.av_packet_add_side_data(tmppkt, C.enum_AVPacketSideDataType(_type), (*C.uint8_t)(data), C.size_t(size))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_packet_shrink_side_data ---

// AVPacketShrinkSideData wraps av_packet_shrink_side_data.
/*
  Shrink the already allocated side data buffer

  @param pkt packet
  @param type side information type
  @param size new side information size
  @return 0 on success, < 0 on failure
*/
func AVPacketShrinkSideData(pkt *AVPacket, _type AVPacketSideDataType, size uint64) (int, error) {
	var tmppkt *C.AVPacket
	if pkt != nil {
		tmppkt = pkt.ptr
	}
	ret := C.av_packet_shrink_side_data(tmppkt, C.enum_AVPacketSideDataType(_type), C.size_t(size))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_packet_get_side_data ---

// AVPacketGetSideData wraps av_packet_get_side_data.
/*
  Get side information from packet.

  @param pkt packet
  @param type desired side information type
  @param size If supplied, *size will be set to the size of the side data
              or to zero if the desired side data is not present.
  @return pointer to data if present or NULL otherwise
*/
func AVPacketGetSideData(pkt *AVPacket, _type AVPacketSideDataType, size *uint64) unsafe.Pointer {
	var tmppkt *C.AVPacket
	if pkt != nil {
		tmppkt = pkt.ptr
	}
	ret := C.av_packet_get_side_data(tmppkt, C.enum_AVPacketSideDataType(_type), (*C.size_t)(unsafe.Pointer(size)))
	return unsafe.Pointer(ret)
}

// --- Function av_packet_pack_dictionary ---

// AVPacketPackDictionary wraps av_packet_pack_dictionary.
/*
  Pack a dictionary for use in side_data.

  @param dict The dictionary to pack.
  @param size pointer to store the size of the returned data
  @return pointer to data if successful, NULL otherwise
*/
func AVPacketPackDictionary(dict *AVDictionary, size *uint64) unsafe.Pointer {
	var tmpdict *C.AVDictionary
	if dict != nil {
		tmpdict = dict.ptr
	}
	ret := C.av_packet_pack_dictionary(tmpdict, (*C.size_t)(unsafe.Pointer(size)))
	return unsafe.Pointer(ret)
}

// --- Function av_packet_unpack_dictionary ---

// AVPacketUnpackDictionary wraps av_packet_unpack_dictionary.
/*
  Unpack a dictionary from side_data.

  @param data data from side_data
  @param size size of the data
  @param dict the metadata storage dictionary
  @return 0 on success, < 0 on failure
*/
func AVPacketUnpackDictionary(data unsafe.Pointer, size uint64, dict **AVDictionary) (int, error) {
	var ptrdict **C.AVDictionary
	var tmpdict *C.AVDictionary
	var oldTmpdict *C.AVDictionary
	if dict != nil {
		innerdict := *dict
		if innerdict != nil {
			tmpdict = innerdict.ptr
			oldTmpdict = tmpdict
		}
		ptrdict = &tmpdict
	}
	ret := C.av_packet_unpack_dictionary((*C.uint8_t)(data), C.size_t(size), ptrdict)
	if tmpdict != oldTmpdict && dict != nil {
		if tmpdict != nil {
			*dict = &AVDictionary{ptr: tmpdict}
		} else {
			*dict = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function av_packet_free_side_data ---

// AVPacketFreeSideData wraps av_packet_free_side_data.
/*
  Convenience function to free all the side data stored.
  All the other fields stay untouched.

  @param pkt packet
*/
func AVPacketFreeSideData(pkt *AVPacket) {
	var tmppkt *C.AVPacket
	if pkt != nil {
		tmppkt = pkt.ptr
	}
	C.av_packet_free_side_data(tmppkt)
}

// --- Function av_packet_ref ---

// AVPacketRef wraps av_packet_ref.
/*
  Setup a new reference to the data described by a given packet

  If src is reference-counted, setup dst as a new reference to the
  buffer in src. Otherwise allocate a new buffer in dst and copy the
  data from src into it.

  All the other fields are copied from src.

  @see av_packet_unref

  @param dst Destination packet. Will be completely overwritten.
  @param src Source packet

  @return 0 on success, a negative AVERROR on error. On error, dst
          will be blank (as if returned by av_packet_alloc()).
*/
func AVPacketRef(dst *AVPacket, src *AVPacket) (int, error) {
	var tmpdst *C.AVPacket
	if dst != nil {
		tmpdst = dst.ptr
	}
	var tmpsrc *C.AVPacket
	if src != nil {
		tmpsrc = src.ptr
	}
	ret := C.av_packet_ref(tmpdst, tmpsrc)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_packet_unref ---

// AVPacketUnref wraps av_packet_unref.
/*
  Wipe the packet.

  Unreference the buffer referenced by the packet and reset the
  remaining packet fields to their default values.

  @param pkt The packet to be unreferenced.
*/
func AVPacketUnref(pkt *AVPacket) {
	var tmppkt *C.AVPacket
	if pkt != nil {
		tmppkt = pkt.ptr
	}
	C.av_packet_unref(tmppkt)
}

// --- Function av_packet_move_ref ---

// AVPacketMoveRef wraps av_packet_move_ref.
/*
  Move every field in src to dst and reset src.

  @see av_packet_unref

  @param src Source packet, will be reset
  @param dst Destination packet
*/
func AVPacketMoveRef(dst *AVPacket, src *AVPacket) {
	var tmpdst *C.AVPacket
	if dst != nil {
		tmpdst = dst.ptr
	}
	var tmpsrc *C.AVPacket
	if src != nil {
		tmpsrc = src.ptr
	}
	C.av_packet_move_ref(tmpdst, tmpsrc)
}

// --- Function av_packet_copy_props ---

// AVPacketCopyProps wraps av_packet_copy_props.
/*
  Copy only "properties" fields from src to dst.

  Properties for the purpose of this function are all the fields
  beside those related to the packet data (buf, data, size)

  @param dst Destination packet
  @param src Source packet

  @return 0 on success AVERROR on failure.
*/
func AVPacketCopyProps(dst *AVPacket, src *AVPacket) (int, error) {
	var tmpdst *C.AVPacket
	if dst != nil {
		tmpdst = dst.ptr
	}
	var tmpsrc *C.AVPacket
	if src != nil {
		tmpsrc = src.ptr
	}
	ret := C.av_packet_copy_props(tmpdst, tmpsrc)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_packet_make_refcounted ---

// AVPacketMakeRefcounted wraps av_packet_make_refcounted.
/*
  Ensure the data described by a given packet is reference counted.

  @note This function does not ensure that the reference will be writable.
        Use av_packet_make_writable instead for that purpose.

  @see av_packet_ref
  @see av_packet_make_writable

  @param pkt packet whose data should be made reference counted.

  @return 0 on success, a negative AVERROR on error. On failure, the
          packet is unchanged.
*/
func AVPacketMakeRefcounted(pkt *AVPacket) (int, error) {
	var tmppkt *C.AVPacket
	if pkt != nil {
		tmppkt = pkt.ptr
	}
	ret := C.av_packet_make_refcounted(tmppkt)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_packet_make_writable ---

// AVPacketMakeWritable wraps av_packet_make_writable.
/*
  Create a writable reference for the data described by a given packet,
  avoiding data copy if possible.

  @param pkt Packet whose data should be made writable.

  @return 0 on success, a negative AVERROR on failure. On failure, the
          packet is unchanged.
*/
func AVPacketMakeWritable(pkt *AVPacket) (int, error) {
	var tmppkt *C.AVPacket
	if pkt != nil {
		tmppkt = pkt.ptr
	}
	ret := C.av_packet_make_writable(tmppkt)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_packet_rescale_ts ---

// AVPacketRescaleTs wraps av_packet_rescale_ts.
/*
  Convert valid timing fields (timestamps / durations) in a packet from one
  timebase to another. Timestamps with unknown values (AV_NOPTS_VALUE) will be
  ignored.

  @param pkt packet on which the conversion will be performed
  @param tb_src source timebase, in which the timing fields in pkt are
                expressed
  @param tb_dst destination timebase, to which the timing fields will be
                converted
*/
func AVPacketRescaleTs(pkt *AVPacket, tbSrc *AVRational, tbDst *AVRational) {
	var tmppkt *C.AVPacket
	if pkt != nil {
		tmppkt = pkt.ptr
	}
	C.av_packet_rescale_ts(tmppkt, tbSrc.value, tbDst.value)
}

// --- Function av_container_fifo_alloc_avpacket ---

// AVContainerFifoAllocAVPacket wraps av_container_fifo_alloc_avpacket.
/*
  Allocate an AVContainerFifo instance for AVPacket.

  @param flags currently unused
*/
func AVContainerFifoAllocAVPacket(flags uint) *AVContainerFifo {
	ret := C.av_container_fifo_alloc_avpacket(C.uint(flags))
	var retMapped *AVContainerFifo
	if ret != nil {
		retMapped = &AVContainerFifo{ptr: ret}
	}
	return retMapped
}

// --- Function avdevice_version ---

// AVDeviceVersion wraps avdevice_version.
//
//	Return the LIBAVDEVICE_VERSION_INT constant.
func AVDeviceVersion() uint {
	ret := C.avdevice_version()
	return uint(ret)
}

// --- Function avdevice_configuration ---

// AVDeviceConfiguration wraps avdevice_configuration.
//
//	Return the libavdevice build-time configuration.
func AVDeviceConfiguration() *CStr {
	ret := C.avdevice_configuration()
	return wrapCStr(ret)
}

// --- Function avdevice_license ---

// AVDeviceLicense wraps avdevice_license.
//
//	Return the libavdevice license.
func AVDeviceLicense() *CStr {
	ret := C.avdevice_license()
	return wrapCStr(ret)
}

// --- Function avdevice_register_all ---

// AVDeviceRegisterAll wraps avdevice_register_all.
//
//	Initialize libavdevice and register all the input and output devices.
func AVDeviceRegisterAll() {
	C.avdevice_register_all()
}

// --- Function av_input_audio_device_next ---

// AVInputAudioDeviceNext wraps av_input_audio_device_next.
/*
  Audio input devices iterator.

  If d is NULL, returns the first registered input audio/video device,
  if d is non-NULL, returns the next registered input audio/video device after d
  or NULL if d is the last one.
*/
func AVInputAudioDeviceNext(d *AVInputFormat) *AVInputFormat {
	var tmpd *C.AVInputFormat
	if d != nil {
		tmpd = d.ptr
	}
	ret := C.av_input_audio_device_next(tmpd)
	var retMapped *AVInputFormat
	if ret != nil {
		retMapped = &AVInputFormat{ptr: ret}
	}
	return retMapped
}

// --- Function av_input_video_device_next ---

// AVInputVideoDeviceNext wraps av_input_video_device_next.
/*
  Video input devices iterator.

  If d is NULL, returns the first registered input audio/video device,
  if d is non-NULL, returns the next registered input audio/video device after d
  or NULL if d is the last one.
*/
func AVInputVideoDeviceNext(d *AVInputFormat) *AVInputFormat {
	var tmpd *C.AVInputFormat
	if d != nil {
		tmpd = d.ptr
	}
	ret := C.av_input_video_device_next(tmpd)
	var retMapped *AVInputFormat
	if ret != nil {
		retMapped = &AVInputFormat{ptr: ret}
	}
	return retMapped
}

// --- Function av_output_audio_device_next ---

// AVOutputAudioDeviceNext wraps av_output_audio_device_next.
/*
  Audio output devices iterator.

  If d is NULL, returns the first registered output audio/video device,
  if d is non-NULL, returns the next registered output audio/video device after d
  or NULL if d is the last one.
*/
func AVOutputAudioDeviceNext(d *AVOutputFormat) *AVOutputFormat {
	var tmpd *C.AVOutputFormat
	if d != nil {
		tmpd = d.ptr
	}
	ret := C.av_output_audio_device_next(tmpd)
	var retMapped *AVOutputFormat
	if ret != nil {
		retMapped = &AVOutputFormat{ptr: ret}
	}
	return retMapped
}

// --- Function av_output_video_device_next ---

// AVOutputVideoDeviceNext wraps av_output_video_device_next.
/*
  Video output devices iterator.

  If d is NULL, returns the first registered output audio/video device,
  if d is non-NULL, returns the next registered output audio/video device after d
  or NULL if d is the last one.
*/
func AVOutputVideoDeviceNext(d *AVOutputFormat) *AVOutputFormat {
	var tmpd *C.AVOutputFormat
	if d != nil {
		tmpd = d.ptr
	}
	ret := C.av_output_video_device_next(tmpd)
	var retMapped *AVOutputFormat
	if ret != nil {
		retMapped = &AVOutputFormat{ptr: ret}
	}
	return retMapped
}

// --- Function avdevice_app_to_dev_control_message ---

// AVDeviceAppToDevControlMessage wraps avdevice_app_to_dev_control_message.
/*
  Send control message from application to device.

  @param s         device context.
  @param type      message type.
  @param data      message data. Exact type depends on message type.
  @param data_size size of message data.
  @return >= 0 on success, negative on error.
          AVERROR(ENOSYS) when device doesn't implement handler of the message.
*/
func AVDeviceAppToDevControlMessage(s *AVFormatContext, _type AVAppToDevMessageType, data unsafe.Pointer, dataSize uint64) (int, error) {
	var tmps *C.AVFormatContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.avdevice_app_to_dev_control_message(tmps, C.enum_AVAppToDevMessageType(_type), data, C.size_t(dataSize))
	return int(ret), WrapErr(int(ret))
}

// --- Function avdevice_dev_to_app_control_message ---

// AVDeviceDevToAppControlMessage wraps avdevice_dev_to_app_control_message.
/*
  Send control message from device to application.

  @param s         device context.
  @param type      message type.
  @param data      message data. Can be NULL.
  @param data_size size of message data.
  @return >= 0 on success, negative on error.
          AVERROR(ENOSYS) when application doesn't implement handler of the message.
*/
func AVDeviceDevToAppControlMessage(s *AVFormatContext, _type AVDevToAppMessageType, data unsafe.Pointer, dataSize uint64) (int, error) {
	var tmps *C.AVFormatContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.avdevice_dev_to_app_control_message(tmps, C.enum_AVDevToAppMessageType(_type), data, C.size_t(dataSize))
	return int(ret), WrapErr(int(ret))
}

// --- Function avdevice_list_devices ---

// AVDeviceListDevices wraps avdevice_list_devices.
/*
  List devices.

  Returns available device names and their parameters.

  @note: Some devices may accept system-dependent device names that cannot be
         autodetected. The list returned by this function cannot be assumed to
         be always completed.

  @param s                device context.
  @param[out] device_list list of autodetected devices.
  @return count of autodetected devices, negative on error.
*/
func AVDeviceListDevices(s *AVFormatContext, deviceList **AVDeviceInfoList) (int, error) {
	var tmps *C.AVFormatContext
	if s != nil {
		tmps = s.ptr
	}
	var ptrdeviceList **C.AVDeviceInfoList
	var tmpdeviceList *C.AVDeviceInfoList
	var oldTmpdeviceList *C.AVDeviceInfoList
	if deviceList != nil {
		innerdeviceList := *deviceList
		if innerdeviceList != nil {
			tmpdeviceList = innerdeviceList.ptr
			oldTmpdeviceList = tmpdeviceList
		}
		ptrdeviceList = &tmpdeviceList
	}
	ret := C.avdevice_list_devices(tmps, ptrdeviceList)
	if tmpdeviceList != oldTmpdeviceList && deviceList != nil {
		if tmpdeviceList != nil {
			*deviceList = &AVDeviceInfoList{ptr: tmpdeviceList}
		} else {
			*deviceList = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function avdevice_free_list_devices ---

// AVDeviceFreeListDevices wraps avdevice_free_list_devices.
/*
  Convenient function to free result of avdevice_list_devices().

  @param device_list device list to be freed.
*/
func AVDeviceFreeListDevices(deviceList **AVDeviceInfoList) {
	var ptrdeviceList **C.AVDeviceInfoList
	var tmpdeviceList *C.AVDeviceInfoList
	var oldTmpdeviceList *C.AVDeviceInfoList
	if deviceList != nil {
		innerdeviceList := *deviceList
		if innerdeviceList != nil {
			tmpdeviceList = innerdeviceList.ptr
			oldTmpdeviceList = tmpdeviceList
		}
		ptrdeviceList = &tmpdeviceList
	}
	C.avdevice_free_list_devices(ptrdeviceList)
	if tmpdeviceList != oldTmpdeviceList && deviceList != nil {
		if tmpdeviceList != nil {
			*deviceList = &AVDeviceInfoList{ptr: tmpdeviceList}
		} else {
			*deviceList = nil
		}
	}
}

// --- Function avdevice_list_input_sources ---

// AVDeviceListInputSources wraps avdevice_list_input_sources.
/*
  List devices.

  Returns available device names and their parameters.
  These are convenient wrappers for avdevice_list_devices().
  Device context is allocated and deallocated internally.

  @param device           device format. May be NULL if device name is set.
  @param device_name      device name. May be NULL if device format is set.
  @param device_options   An AVDictionary filled with device-private options. May be NULL.
                          The same options must be passed later to avformat_write_header() for output
                          devices or avformat_open_input() for input devices, or at any other place
                          that affects device-private options.
  @param[out] device_list list of autodetected devices
  @return count of autodetected devices, negative on error.
  @note device argument takes precedence over device_name when both are set.
*/
func AVDeviceListInputSources(device *AVInputFormat, deviceName *CStr, deviceOptions *AVDictionary, deviceList **AVDeviceInfoList) (int, error) {
	var tmpdevice *C.AVInputFormat
	if device != nil {
		tmpdevice = device.ptr
	}
	var tmpdeviceName *C.char
	if deviceName != nil {
		tmpdeviceName = deviceName.ptr
	}
	var tmpdeviceOptions *C.AVDictionary
	if deviceOptions != nil {
		tmpdeviceOptions = deviceOptions.ptr
	}
	var ptrdeviceList **C.AVDeviceInfoList
	var tmpdeviceList *C.AVDeviceInfoList
	var oldTmpdeviceList *C.AVDeviceInfoList
	if deviceList != nil {
		innerdeviceList := *deviceList
		if innerdeviceList != nil {
			tmpdeviceList = innerdeviceList.ptr
			oldTmpdeviceList = tmpdeviceList
		}
		ptrdeviceList = &tmpdeviceList
	}
	ret := C.avdevice_list_input_sources(tmpdevice, tmpdeviceName, tmpdeviceOptions, ptrdeviceList)
	if tmpdeviceList != oldTmpdeviceList && deviceList != nil {
		if tmpdeviceList != nil {
			*deviceList = &AVDeviceInfoList{ptr: tmpdeviceList}
		} else {
			*deviceList = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function avdevice_list_output_sinks ---

// AVDeviceListOutputSinks wraps avdevice_list_output_sinks.
func AVDeviceListOutputSinks(device *AVOutputFormat, deviceName *CStr, deviceOptions *AVDictionary, deviceList **AVDeviceInfoList) (int, error) {
	var tmpdevice *C.AVOutputFormat
	if device != nil {
		tmpdevice = device.ptr
	}
	var tmpdeviceName *C.char
	if deviceName != nil {
		tmpdeviceName = deviceName.ptr
	}
	var tmpdeviceOptions *C.AVDictionary
	if deviceOptions != nil {
		tmpdeviceOptions = deviceOptions.ptr
	}
	var ptrdeviceList **C.AVDeviceInfoList
	var tmpdeviceList *C.AVDeviceInfoList
	var oldTmpdeviceList *C.AVDeviceInfoList
	if deviceList != nil {
		innerdeviceList := *deviceList
		if innerdeviceList != nil {
			tmpdeviceList = innerdeviceList.ptr
			oldTmpdeviceList = tmpdeviceList
		}
		ptrdeviceList = &tmpdeviceList
	}
	ret := C.avdevice_list_output_sinks(tmpdevice, tmpdeviceName, tmpdeviceOptions, ptrdeviceList)
	if tmpdeviceList != oldTmpdeviceList && deviceList != nil {
		if tmpdeviceList != nil {
			*deviceList = &AVDeviceInfoList{ptr: tmpdeviceList}
		} else {
			*deviceList = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function avfilter_version ---

// AVFilterVersion wraps avfilter_version.
//
//	Return the LIBAVFILTER_VERSION_INT constant.
func AVFilterVersion() uint {
	ret := C.avfilter_version()
	return uint(ret)
}

// --- Function avfilter_configuration ---

// AVFilterConfiguration wraps avfilter_configuration.
//
//	Return the libavfilter build-time configuration.
func AVFilterConfiguration() *CStr {
	ret := C.avfilter_configuration()
	return wrapCStr(ret)
}

// --- Function avfilter_license ---

// AVFilterLicense wraps avfilter_license.
//
//	Return the libavfilter license.
func AVFilterLicense() *CStr {
	ret := C.avfilter_license()
	return wrapCStr(ret)
}

// --- Function avfilter_pad_get_name ---

// AVFilterPadGetName wraps avfilter_pad_get_name.
/*
  Get the name of an AVFilterPad.

  @param pads an array of AVFilterPads
  @param pad_idx index of the pad in the array; it is the caller's
                 responsibility to ensure the index is valid

  @return name of the pad_idx'th pad in pads
*/
func AVFilterPadGetName(pads *AVFilterPad, padIdx int) *CStr {
	var tmppads *C.AVFilterPad
	if pads != nil {
		tmppads = pads.ptr
	}
	ret := C.avfilter_pad_get_name(tmppads, C.int(padIdx))
	return wrapCStr(ret)
}

// --- Function avfilter_pad_get_type ---

// AVFilterPadGetType wraps avfilter_pad_get_type.
/*
  Get the type of an AVFilterPad.

  @param pads an array of AVFilterPads
  @param pad_idx index of the pad in the array; it is the caller's
                 responsibility to ensure the index is valid

  @return type of the pad_idx'th pad in pads
*/
func AVFilterPadGetType(pads *AVFilterPad, padIdx int) AVMediaType {
	var tmppads *C.AVFilterPad
	if pads != nil {
		tmppads = pads.ptr
	}
	ret := C.avfilter_pad_get_type(tmppads, C.int(padIdx))
	return AVMediaType(ret)
}

// --- Function avfilter_link_get_hw_frames_ctx ---

// AVFilterLinkGetHWFramesCtx wraps avfilter_link_get_hw_frames_ctx.
/*
  Get the hardware frames context of a filter link.

  @param link an AVFilterLink

  @return a ref-counted copy of the link's hw_frames_ctx field if there is
          a hardware frames context associated with the link or NULL otherwise.
          The returned AVBufferRef needs to be released with av_buffer_unref()
          when it is no longer used.
*/
func AVFilterLinkGetHWFramesCtx(link *AVFilterLink) *AVBufferRef {
	var tmplink *C.AVFilterLink
	if link != nil {
		tmplink = link.ptr
	}
	ret := C.avfilter_link_get_hw_frames_ctx(tmplink)
	var retMapped *AVBufferRef
	if ret != nil {
		retMapped = &AVBufferRef{ptr: ret}
	}
	return retMapped
}

// --- Function avfilter_filter_pad_count ---

// AVFilterFilterPadCount wraps avfilter_filter_pad_count.
//
//	Get the number of elements in an AVFilter's inputs or outputs array.
func AVFilterFilterPadCount(filter *AVFilter, isOutput int) uint {
	var tmpfilter *C.AVFilter
	if filter != nil {
		tmpfilter = filter.ptr
	}
	ret := C.avfilter_filter_pad_count(tmpfilter, C.int(isOutput))
	return uint(ret)
}

// --- Function avfilter_link ---

// AVFilterLink_ wraps avfilter_link.
/*
  Link two filters together.

  @param src    the source filter
  @param srcpad index of the output pad on the source filter
  @param dst    the destination filter
  @param dstpad index of the input pad on the destination filter
  @return       zero on success
*/
func AVFilterLink_(src *AVFilterContext, srcpad uint, dst *AVFilterContext, dstpad uint) (int, error) {
	var tmpsrc *C.AVFilterContext
	if src != nil {
		tmpsrc = src.ptr
	}
	var tmpdst *C.AVFilterContext
	if dst != nil {
		tmpdst = dst.ptr
	}
	ret := C.avfilter_link(tmpsrc, C.uint(srcpad), tmpdst, C.uint(dstpad))
	return int(ret), WrapErr(int(ret))
}

// --- Function avfilter_process_command ---

// AVFilterProcessCommand wraps avfilter_process_command.
/*
  Make the filter instance process a command.
  It is recommended to use avfilter_graph_send_command().
*/
func AVFilterProcessCommand(filter *AVFilterContext, cmd *CStr, arg *CStr, res *CStr, resLen int, flags int) (int, error) {
	var tmpfilter *C.AVFilterContext
	if filter != nil {
		tmpfilter = filter.ptr
	}
	var tmpcmd *C.char
	if cmd != nil {
		tmpcmd = cmd.ptr
	}
	var tmparg *C.char
	if arg != nil {
		tmparg = arg.ptr
	}
	var tmpres *C.char
	if res != nil {
		tmpres = res.ptr
	}
	ret := C.avfilter_process_command(tmpfilter, tmpcmd, tmparg, tmpres, C.int(resLen), C.int(flags))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_filter_iterate ---

// av_filter_iterate skipped due to opaque

// --- Function avfilter_get_by_name ---

// AVFilterGetByName wraps avfilter_get_by_name.
/*
  Get a filter definition matching the given name.

  @param name the filter name to find
  @return     the filter definition, if any matching one is registered.
              NULL if none found.
*/
func AVFilterGetByName(name *CStr) *AVFilter {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.avfilter_get_by_name(tmpname)
	var retMapped *AVFilter
	if ret != nil {
		retMapped = &AVFilter{ptr: ret}
	}
	return retMapped
}

// --- Function avfilter_init_str ---

// AVFilterInitStr wraps avfilter_init_str.
/*
  Initialize a filter with the supplied parameters.

  @param ctx  uninitialized filter context to initialize
  @param args Options to initialize the filter with. This must be a
              ':'-separated list of options in the 'key=value' form.
              May be NULL if the options have been set directly using the
              AVOptions API or there are no options that need to be set.
  @return 0 on success, a negative AVERROR on failure
*/
func AVFilterInitStr(ctx *AVFilterContext, args *CStr) (int, error) {
	var tmpctx *C.AVFilterContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	var tmpargs *C.char
	if args != nil {
		tmpargs = args.ptr
	}
	ret := C.avfilter_init_str(tmpctx, tmpargs)
	return int(ret), WrapErr(int(ret))
}

// --- Function avfilter_init_dict ---

// AVFilterInitDict wraps avfilter_init_dict.
/*
  Initialize a filter with the supplied dictionary of options.

  @param ctx     uninitialized filter context to initialize
  @param options An AVDictionary filled with options for this filter. On
                 return this parameter will be destroyed and replaced with
                 a dict containing options that were not found. This dictionary
                 must be freed by the caller.
                 May be NULL, then this function is equivalent to
                 avfilter_init_str() with the second parameter set to NULL.
  @return 0 on success, a negative AVERROR on failure

  @note This function and avfilter_init_str() do essentially the same thing,
  the difference is in manner in which the options are passed. It is up to the
  calling code to choose whichever is more preferable. The two functions also
  behave differently when some of the provided options are not declared as
  supported by the filter. In such a case, avfilter_init_str() will fail, but
  this function will leave those extra options in the options AVDictionary and
  continue as usual.
*/
func AVFilterInitDict(ctx *AVFilterContext, options **AVDictionary) (int, error) {
	var tmpctx *C.AVFilterContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	var ptroptions **C.AVDictionary
	var tmpoptions *C.AVDictionary
	var oldTmpoptions *C.AVDictionary
	if options != nil {
		inneroptions := *options
		if inneroptions != nil {
			tmpoptions = inneroptions.ptr
			oldTmpoptions = tmpoptions
		}
		ptroptions = &tmpoptions
	}
	ret := C.avfilter_init_dict(tmpctx, ptroptions)
	if tmpoptions != oldTmpoptions && options != nil {
		if tmpoptions != nil {
			*options = &AVDictionary{ptr: tmpoptions}
		} else {
			*options = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function avfilter_free ---

// AVFilterFree wraps avfilter_free.
/*
  Free a filter context. This will also remove the filter from its
  filtergraph's list of filters.

  @param filter the filter to free
*/
func AVFilterFree(filter *AVFilterContext) {
	var tmpfilter *C.AVFilterContext
	if filter != nil {
		tmpfilter = filter.ptr
	}
	C.avfilter_free(tmpfilter)
}

// --- Function avfilter_insert_filter ---

// AVFilterInsertFilter wraps avfilter_insert_filter.
/*
  Insert a filter in the middle of an existing link.

  @param link the link into which the filter should be inserted
  @param filt the filter to be inserted
  @param filt_srcpad_idx the input pad on the filter to connect
  @param filt_dstpad_idx the output pad on the filter to connect
  @return     zero on success
*/
func AVFilterInsertFilter(link *AVFilterLink, filt *AVFilterContext, filtSrcpadIdx uint, filtDstpadIdx uint) (int, error) {
	var tmplink *C.AVFilterLink
	if link != nil {
		tmplink = link.ptr
	}
	var tmpfilt *C.AVFilterContext
	if filt != nil {
		tmpfilt = filt.ptr
	}
	ret := C.avfilter_insert_filter(tmplink, tmpfilt, C.uint(filtSrcpadIdx), C.uint(filtDstpadIdx))
	return int(ret), WrapErr(int(ret))
}

// --- Function avfilter_get_class ---

// AVFilterGetClass wraps avfilter_get_class.
/*
  @return AVClass for AVFilterContext.

  @see av_opt_find().
*/
func AVFilterGetClass() *AVClass {
	ret := C.avfilter_get_class()
	var retMapped *AVClass
	if ret != nil {
		retMapped = &AVClass{ptr: ret}
	}
	return retMapped
}

// --- Function avfilter_graph_alloc ---

// AVFilterGraphAlloc wraps avfilter_graph_alloc.
/*
  Allocate a filter graph.

  @return the allocated filter graph on success or NULL.
*/
func AVFilterGraphAlloc() *AVFilterGraph {
	ret := C.avfilter_graph_alloc()
	var retMapped *AVFilterGraph
	if ret != nil {
		retMapped = &AVFilterGraph{ptr: ret}
	}
	return retMapped
}

// --- Function avfilter_graph_alloc_filter ---

// AVFilterGraphAllocFilter wraps avfilter_graph_alloc_filter.
/*
  Create a new filter instance in a filter graph.

  @param graph graph in which the new filter will be used
  @param filter the filter to create an instance of
  @param name Name to give to the new instance (will be copied to
              AVFilterContext.name). This may be used by the caller to identify
              different filters, libavfilter itself assigns no semantics to
              this parameter. May be NULL.

  @return the context of the newly created filter instance (note that it is
          also retrievable directly through AVFilterGraph.filters or with
          avfilter_graph_get_filter()) on success or NULL on failure.
*/
func AVFilterGraphAllocFilter(graph *AVFilterGraph, filter *AVFilter, name *CStr) *AVFilterContext {
	var tmpgraph *C.AVFilterGraph
	if graph != nil {
		tmpgraph = graph.ptr
	}
	var tmpfilter *C.AVFilter
	if filter != nil {
		tmpfilter = filter.ptr
	}
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.avfilter_graph_alloc_filter(tmpgraph, tmpfilter, tmpname)
	var retMapped *AVFilterContext
	if ret != nil {
		retMapped = &AVFilterContext{ptr: ret}
	}
	return retMapped
}

// --- Function avfilter_graph_get_filter ---

// AVFilterGraphGetFilter wraps avfilter_graph_get_filter.
/*
  Get a filter instance identified by instance name from graph.

  @param graph filter graph to search through.
  @param name filter instance name (should be unique in the graph).
  @return the pointer to the found filter instance or NULL if it
  cannot be found.
*/
func AVFilterGraphGetFilter(graph *AVFilterGraph, name *CStr) *AVFilterContext {
	var tmpgraph *C.AVFilterGraph
	if graph != nil {
		tmpgraph = graph.ptr
	}
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.avfilter_graph_get_filter(tmpgraph, tmpname)
	var retMapped *AVFilterContext
	if ret != nil {
		retMapped = &AVFilterContext{ptr: ret}
	}
	return retMapped
}

// --- Function avfilter_graph_create_filter ---

// AVFilterGraphCreateFilter wraps avfilter_graph_create_filter.
/*
  A convenience wrapper that allocates and initializes a filter in a single
  step. The filter instance is created from the filter filt and inited with the
  parameter args. opaque is currently ignored.

  In case of success put in *filt_ctx the pointer to the created
  filter instance, otherwise set *filt_ctx to NULL.

  @param name the instance name to give to the created filter instance
  @param graph_ctx the filter graph
  @return a negative AVERROR error code in case of failure, a non
  negative value otherwise

  @warning Since the filter is initialized after this function successfully
           returns, you MUST NOT set any further options on it. If you need to
           do that, call ::avfilter_graph_alloc_filter(), followed by setting
           the options, followed by ::avfilter_init_dict() instead of this
           function.
*/
func AVFilterGraphCreateFilter(filtCtx **AVFilterContext, filt *AVFilter, name *CStr, args *CStr, opaque unsafe.Pointer, graphCtx *AVFilterGraph) (int, error) {
	var ptrfiltCtx **C.AVFilterContext
	var tmpfiltCtx *C.AVFilterContext
	var oldTmpfiltCtx *C.AVFilterContext
	if filtCtx != nil {
		innerfiltCtx := *filtCtx
		if innerfiltCtx != nil {
			tmpfiltCtx = innerfiltCtx.ptr
			oldTmpfiltCtx = tmpfiltCtx
		}
		ptrfiltCtx = &tmpfiltCtx
	}
	var tmpfilt *C.AVFilter
	if filt != nil {
		tmpfilt = filt.ptr
	}
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	var tmpargs *C.char
	if args != nil {
		tmpargs = args.ptr
	}
	var tmpgraphCtx *C.AVFilterGraph
	if graphCtx != nil {
		tmpgraphCtx = graphCtx.ptr
	}
	ret := C.avfilter_graph_create_filter(ptrfiltCtx, tmpfilt, tmpname, tmpargs, opaque, tmpgraphCtx)
	if tmpfiltCtx != oldTmpfiltCtx && filtCtx != nil {
		if tmpfiltCtx != nil {
			*filtCtx = &AVFilterContext{ptr: tmpfiltCtx}
		} else {
			*filtCtx = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function avfilter_graph_set_auto_convert ---

// AVFilterGraphSetAutoConvert wraps avfilter_graph_set_auto_convert.
/*
  Enable or disable automatic format conversion inside the graph.

  Note that format conversion can still happen inside explicitly inserted
  scale and aresample filters.

  @param flags  any of the AVFILTER_AUTO_CONVERT_* constants
*/
func AVFilterGraphSetAutoConvert(graph *AVFilterGraph, flags uint) {
	var tmpgraph *C.AVFilterGraph
	if graph != nil {
		tmpgraph = graph.ptr
	}
	C.avfilter_graph_set_auto_convert(tmpgraph, C.uint(flags))
}

// --- Function avfilter_graph_config ---

// AVFilterGraphConfig wraps avfilter_graph_config.
/*
  Check validity and configure all the links and formats in the graph.

  @param graphctx the filter graph
  @param log_ctx context used for logging
  @return >= 0 in case of success, a negative AVERROR code otherwise
*/
func AVFilterGraphConfig(graphctx *AVFilterGraph, logCtx unsafe.Pointer) (int, error) {
	var tmpgraphctx *C.AVFilterGraph
	if graphctx != nil {
		tmpgraphctx = graphctx.ptr
	}
	ret := C.avfilter_graph_config(tmpgraphctx, logCtx)
	return int(ret), WrapErr(int(ret))
}

// --- Function avfilter_graph_free ---

// AVFilterGraphFree wraps avfilter_graph_free.
/*
  Free a graph, destroy its links, and set *graph to NULL.
  If *graph is NULL, do nothing.
*/
func AVFilterGraphFree(graph **AVFilterGraph) {
	var ptrgraph **C.AVFilterGraph
	var tmpgraph *C.AVFilterGraph
	var oldTmpgraph *C.AVFilterGraph
	if graph != nil {
		innergraph := *graph
		if innergraph != nil {
			tmpgraph = innergraph.ptr
			oldTmpgraph = tmpgraph
		}
		ptrgraph = &tmpgraph
	}
	C.avfilter_graph_free(ptrgraph)
	if tmpgraph != oldTmpgraph && graph != nil {
		if tmpgraph != nil {
			*graph = &AVFilterGraph{ptr: tmpgraph}
		} else {
			*graph = nil
		}
	}
}

// --- Function avfilter_inout_alloc ---

// AVFilterInoutAlloc wraps avfilter_inout_alloc.
/*
  Allocate a single AVFilterInOut entry.
  Must be freed with avfilter_inout_free().
  @return allocated AVFilterInOut on success, NULL on failure.
*/
func AVFilterInoutAlloc() *AVFilterInOut {
	ret := C.avfilter_inout_alloc()
	var retMapped *AVFilterInOut
	if ret != nil {
		retMapped = &AVFilterInOut{ptr: ret}
	}
	return retMapped
}

// --- Function avfilter_inout_free ---

// AVFilterInoutFree wraps avfilter_inout_free.
/*
  Free the supplied list of AVFilterInOut and set *inout to NULL.
  If *inout is NULL, do nothing.
*/
func AVFilterInoutFree(inout **AVFilterInOut) {
	var ptrinout **C.AVFilterInOut
	var tmpinout *C.AVFilterInOut
	var oldTmpinout *C.AVFilterInOut
	if inout != nil {
		innerinout := *inout
		if innerinout != nil {
			tmpinout = innerinout.ptr
			oldTmpinout = tmpinout
		}
		ptrinout = &tmpinout
	}
	C.avfilter_inout_free(ptrinout)
	if tmpinout != oldTmpinout && inout != nil {
		if tmpinout != nil {
			*inout = &AVFilterInOut{ptr: tmpinout}
		} else {
			*inout = nil
		}
	}
}

// --- Function avfilter_graph_parse ---

// AVFilterGraphParse wraps avfilter_graph_parse.
/*
  Add a graph described by a string to a graph.

  @note The caller must provide the lists of inputs and outputs,
  which therefore must be known before calling the function.

  @note The inputs parameter describes inputs of the already existing
  part of the graph; i.e. from the point of view of the newly created
  part, they are outputs. Similarly the outputs parameter describes
  outputs of the already existing filters, which are provided as
  inputs to the parsed filters.

  @param graph   the filter graph where to link the parsed graph context
  @param filters string to be parsed
  @param inputs  linked list to the inputs of the graph
  @param outputs linked list to the outputs of the graph
  @return zero on success, a negative AVERROR code on error
*/
func AVFilterGraphParse(graph *AVFilterGraph, filters *CStr, inputs *AVFilterInOut, outputs *AVFilterInOut, logCtx unsafe.Pointer) (int, error) {
	var tmpgraph *C.AVFilterGraph
	if graph != nil {
		tmpgraph = graph.ptr
	}
	var tmpfilters *C.char
	if filters != nil {
		tmpfilters = filters.ptr
	}
	var tmpinputs *C.AVFilterInOut
	if inputs != nil {
		tmpinputs = inputs.ptr
	}
	var tmpoutputs *C.AVFilterInOut
	if outputs != nil {
		tmpoutputs = outputs.ptr
	}
	ret := C.avfilter_graph_parse(tmpgraph, tmpfilters, tmpinputs, tmpoutputs, logCtx)
	return int(ret), WrapErr(int(ret))
}

// --- Function avfilter_graph_parse_ptr ---

// AVFilterGraphParsePtr wraps avfilter_graph_parse_ptr.
/*
  Add a graph described by a string to a graph.

  In the graph filters description, if the input label of the first
  filter is not specified, "in" is assumed; if the output label of
  the last filter is not specified, "out" is assumed.

  @param graph   the filter graph where to link the parsed graph context
  @param filters string to be parsed
  @param inputs  pointer to a linked list to the inputs of the graph, may be NULL.
                 If non-NULL, *inputs is updated to contain the list of open inputs
                 after the parsing, should be freed with avfilter_inout_free().
  @param outputs pointer to a linked list to the outputs of the graph, may be NULL.
                 If non-NULL, *outputs is updated to contain the list of open outputs
                 after the parsing, should be freed with avfilter_inout_free().
  @return non negative on success, a negative AVERROR code on error
*/
func AVFilterGraphParsePtr(graph *AVFilterGraph, filters *CStr, inputs **AVFilterInOut, outputs **AVFilterInOut, logCtx unsafe.Pointer) (int, error) {
	var tmpgraph *C.AVFilterGraph
	if graph != nil {
		tmpgraph = graph.ptr
	}
	var tmpfilters *C.char
	if filters != nil {
		tmpfilters = filters.ptr
	}
	var ptrinputs **C.AVFilterInOut
	var tmpinputs *C.AVFilterInOut
	var oldTmpinputs *C.AVFilterInOut
	if inputs != nil {
		innerinputs := *inputs
		if innerinputs != nil {
			tmpinputs = innerinputs.ptr
			oldTmpinputs = tmpinputs
		}
		ptrinputs = &tmpinputs
	}
	var ptroutputs **C.AVFilterInOut
	var tmpoutputs *C.AVFilterInOut
	var oldTmpoutputs *C.AVFilterInOut
	if outputs != nil {
		inneroutputs := *outputs
		if inneroutputs != nil {
			tmpoutputs = inneroutputs.ptr
			oldTmpoutputs = tmpoutputs
		}
		ptroutputs = &tmpoutputs
	}
	ret := C.avfilter_graph_parse_ptr(tmpgraph, tmpfilters, ptrinputs, ptroutputs, logCtx)
	if tmpinputs != oldTmpinputs && inputs != nil {
		if tmpinputs != nil {
			*inputs = &AVFilterInOut{ptr: tmpinputs}
		} else {
			*inputs = nil
		}
	}
	if tmpoutputs != oldTmpoutputs && outputs != nil {
		if tmpoutputs != nil {
			*outputs = &AVFilterInOut{ptr: tmpoutputs}
		} else {
			*outputs = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function avfilter_graph_parse2 ---

// AVFilterGraphParse2 wraps avfilter_graph_parse2.
/*
  Add a graph described by a string to a graph.

  @param[in]  graph   the filter graph where to link the parsed graph context
  @param[in]  filters string to be parsed
  @param[out] inputs  a linked list of all free (unlinked) inputs of the
                      parsed graph will be returned here. It is to be freed
                      by the caller using avfilter_inout_free().
  @param[out] outputs a linked list of all free (unlinked) outputs of the
                      parsed graph will be returned here. It is to be freed by the
                      caller using avfilter_inout_free().
  @return zero on success, a negative AVERROR code on error

  @note This function returns the inputs and outputs that are left
  unlinked after parsing the graph and the caller then deals with
  them.
  @note This function makes no reference whatsoever to already
  existing parts of the graph and the inputs parameter will on return
  contain inputs of the newly parsed part of the graph.  Analogously
  the outputs parameter will contain outputs of the newly created
  filters.
*/
func AVFilterGraphParse2(graph *AVFilterGraph, filters *CStr, inputs **AVFilterInOut, outputs **AVFilterInOut) (int, error) {
	var tmpgraph *C.AVFilterGraph
	if graph != nil {
		tmpgraph = graph.ptr
	}
	var tmpfilters *C.char
	if filters != nil {
		tmpfilters = filters.ptr
	}
	var ptrinputs **C.AVFilterInOut
	var tmpinputs *C.AVFilterInOut
	var oldTmpinputs *C.AVFilterInOut
	if inputs != nil {
		innerinputs := *inputs
		if innerinputs != nil {
			tmpinputs = innerinputs.ptr
			oldTmpinputs = tmpinputs
		}
		ptrinputs = &tmpinputs
	}
	var ptroutputs **C.AVFilterInOut
	var tmpoutputs *C.AVFilterInOut
	var oldTmpoutputs *C.AVFilterInOut
	if outputs != nil {
		inneroutputs := *outputs
		if inneroutputs != nil {
			tmpoutputs = inneroutputs.ptr
			oldTmpoutputs = tmpoutputs
		}
		ptroutputs = &tmpoutputs
	}
	ret := C.avfilter_graph_parse2(tmpgraph, tmpfilters, ptrinputs, ptroutputs)
	if tmpinputs != oldTmpinputs && inputs != nil {
		if tmpinputs != nil {
			*inputs = &AVFilterInOut{ptr: tmpinputs}
		} else {
			*inputs = nil
		}
	}
	if tmpoutputs != oldTmpoutputs && outputs != nil {
		if tmpoutputs != nil {
			*outputs = &AVFilterInOut{ptr: tmpoutputs}
		} else {
			*outputs = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function avfilter_graph_segment_parse ---

// AVFilterGraphSegmentParse wraps avfilter_graph_segment_parse.
/*
  Parse a textual filtergraph description into an intermediate form.

  This intermediate representation is intended to be modified by the caller as
  described in the documentation of AVFilterGraphSegment and its children, and
  then applied to the graph either manually or with other
  avfilter_graph_segment_*() functions. See the documentation for
  avfilter_graph_segment_apply() for the canonical way to apply
  AVFilterGraphSegment.

  @param graph Filter graph the parsed segment is associated with. Will only be
               used for logging and similar auxiliary purposes. The graph will
               not be actually modified by this function - the parsing results
               are instead stored in seg for further processing.
  @param graph_str a string describing the filtergraph segment
  @param flags reserved for future use, caller must set to 0 for now
  @param seg A pointer to the newly-created AVFilterGraphSegment is written
             here on success. The graph segment is owned by the caller and must
             be freed with avfilter_graph_segment_free() before graph itself is
             freed.

  @retval "non-negative number" success
  @retval "negative error code" failure
*/
func AVFilterGraphSegmentParse(graph *AVFilterGraph, graphStr *CStr, flags int, seg **AVFilterGraphSegment) (int, error) {
	var tmpgraph *C.AVFilterGraph
	if graph != nil {
		tmpgraph = graph.ptr
	}
	var tmpgraphStr *C.char
	if graphStr != nil {
		tmpgraphStr = graphStr.ptr
	}
	var ptrseg **C.AVFilterGraphSegment
	var tmpseg *C.AVFilterGraphSegment
	var oldTmpseg *C.AVFilterGraphSegment
	if seg != nil {
		innerseg := *seg
		if innerseg != nil {
			tmpseg = innerseg.ptr
			oldTmpseg = tmpseg
		}
		ptrseg = &tmpseg
	}
	ret := C.avfilter_graph_segment_parse(tmpgraph, tmpgraphStr, C.int(flags), ptrseg)
	if tmpseg != oldTmpseg && seg != nil {
		if tmpseg != nil {
			*seg = &AVFilterGraphSegment{ptr: tmpseg}
		} else {
			*seg = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function avfilter_graph_segment_create_filters ---

// AVFilterGraphSegmentCreateFilters wraps avfilter_graph_segment_create_filters.
/*
  Create filters specified in a graph segment.

  Walk through the creation-pending AVFilterParams in the segment and create
  new filter instances for them.
  Creation-pending params are those where AVFilterParams.filter_name is
  non-NULL (and hence AVFilterParams.filter is NULL). All other AVFilterParams
  instances are ignored.

  For any filter created by this function, the corresponding
  AVFilterParams.filter is set to the newly-created filter context,
  AVFilterParams.filter_name and AVFilterParams.instance_name are freed and set
  to NULL.

  @param seg the filtergraph segment to process
  @param flags reserved for future use, caller must set to 0 for now

  @retval "non-negative number" Success, all creation-pending filters were
                                successfully created
  @retval AVERROR_FILTER_NOT_FOUND some filter's name did not correspond to a
                                   known filter
  @retval "another negative error code" other failures

  @note Calling this function multiple times is safe, as it is idempotent.
*/
func AVFilterGraphSegmentCreateFilters(seg *AVFilterGraphSegment, flags int) (int, error) {
	var tmpseg *C.AVFilterGraphSegment
	if seg != nil {
		tmpseg = seg.ptr
	}
	ret := C.avfilter_graph_segment_create_filters(tmpseg, C.int(flags))
	return int(ret), WrapErr(int(ret))
}

// --- Function avfilter_graph_segment_apply_opts ---

// AVFilterGraphSegmentApplyOpts wraps avfilter_graph_segment_apply_opts.
/*
  Apply parsed options to filter instances in a graph segment.

  Walk through all filter instances in the graph segment that have option
  dictionaries associated with them and apply those options with
  av_opt_set_dict2(..., AV_OPT_SEARCH_CHILDREN). AVFilterParams.opts is
  replaced by the dictionary output by av_opt_set_dict2(), which should be
  empty (NULL) if all options were successfully applied.

  If any options could not be found, this function will continue processing all
  other filters and finally return AVERROR_OPTION_NOT_FOUND (unless another
  error happens). The calling program may then deal with unapplied options as
  it wishes.

  Any creation-pending filters (see avfilter_graph_segment_create_filters())
  present in the segment will cause this function to fail. AVFilterParams with
  no associated filter context are simply skipped.

  @param seg the filtergraph segment to process
  @param flags reserved for future use, caller must set to 0 for now

  @retval "non-negative number" Success, all options were successfully applied.
  @retval AVERROR_OPTION_NOT_FOUND some options were not found in a filter
  @retval "another negative error code" other failures

  @note Calling this function multiple times is safe, as it is idempotent.
*/
func AVFilterGraphSegmentApplyOpts(seg *AVFilterGraphSegment, flags int) (int, error) {
	var tmpseg *C.AVFilterGraphSegment
	if seg != nil {
		tmpseg = seg.ptr
	}
	ret := C.avfilter_graph_segment_apply_opts(tmpseg, C.int(flags))
	return int(ret), WrapErr(int(ret))
}

// --- Function avfilter_graph_segment_init ---

// AVFilterGraphSegmentInit wraps avfilter_graph_segment_init.
/*
  Initialize all filter instances in a graph segment.

  Walk through all filter instances in the graph segment and call
  avfilter_init_dict(..., NULL) on those that have not been initialized yet.

  Any creation-pending filters (see avfilter_graph_segment_create_filters())
  present in the segment will cause this function to fail. AVFilterParams with
  no associated filter context or whose filter context is already initialized,
  are simply skipped.

  @param seg the filtergraph segment to process
  @param flags reserved for future use, caller must set to 0 for now

  @retval "non-negative number" Success, all filter instances were successfully
                                initialized
  @retval "negative error code" failure

  @note Calling this function multiple times is safe, as it is idempotent.
*/
func AVFilterGraphSegmentInit(seg *AVFilterGraphSegment, flags int) (int, error) {
	var tmpseg *C.AVFilterGraphSegment
	if seg != nil {
		tmpseg = seg.ptr
	}
	ret := C.avfilter_graph_segment_init(tmpseg, C.int(flags))
	return int(ret), WrapErr(int(ret))
}

// --- Function avfilter_graph_segment_link ---

// AVFilterGraphSegmentLink wraps avfilter_graph_segment_link.
/*
  Link filters in a graph segment.

  Walk through all filter instances in the graph segment and try to link all
  unlinked input and output pads. Any creation-pending filters (see
  avfilter_graph_segment_create_filters()) present in the segment will cause
  this function to fail. Disabled filters and already linked pads are skipped.

  Every filter output pad that has a corresponding AVFilterPadParams with a
  non-NULL label is
  - linked to the input with the matching label, if one exists;
  - exported in the outputs linked list otherwise, with the label preserved.
  Unlabeled outputs are
  - linked to the first unlinked unlabeled input in the next non-disabled
    filter in the chain, if one exists
  - exported in the outputs linked list otherwise, with NULL label

  Similarly, unlinked input pads are exported in the inputs linked list.

  @param seg the filtergraph segment to process
  @param flags reserved for future use, caller must set to 0 for now
  @param[out] inputs  a linked list of all free (unlinked) inputs of the
                      filters in this graph segment will be returned here. It
                      is to be freed by the caller using avfilter_inout_free().
  @param[out] outputs a linked list of all free (unlinked) outputs of the
                      filters in this graph segment will be returned here. It
                      is to be freed by the caller using avfilter_inout_free().

  @retval "non-negative number" success
  @retval "negative error code" failure

  @note Calling this function multiple times is safe, as it is idempotent.
*/
func AVFilterGraphSegmentLink(seg *AVFilterGraphSegment, flags int, inputs **AVFilterInOut, outputs **AVFilterInOut) (int, error) {
	var tmpseg *C.AVFilterGraphSegment
	if seg != nil {
		tmpseg = seg.ptr
	}
	var ptrinputs **C.AVFilterInOut
	var tmpinputs *C.AVFilterInOut
	var oldTmpinputs *C.AVFilterInOut
	if inputs != nil {
		innerinputs := *inputs
		if innerinputs != nil {
			tmpinputs = innerinputs.ptr
			oldTmpinputs = tmpinputs
		}
		ptrinputs = &tmpinputs
	}
	var ptroutputs **C.AVFilterInOut
	var tmpoutputs *C.AVFilterInOut
	var oldTmpoutputs *C.AVFilterInOut
	if outputs != nil {
		inneroutputs := *outputs
		if inneroutputs != nil {
			tmpoutputs = inneroutputs.ptr
			oldTmpoutputs = tmpoutputs
		}
		ptroutputs = &tmpoutputs
	}
	ret := C.avfilter_graph_segment_link(tmpseg, C.int(flags), ptrinputs, ptroutputs)
	if tmpinputs != oldTmpinputs && inputs != nil {
		if tmpinputs != nil {
			*inputs = &AVFilterInOut{ptr: tmpinputs}
		} else {
			*inputs = nil
		}
	}
	if tmpoutputs != oldTmpoutputs && outputs != nil {
		if tmpoutputs != nil {
			*outputs = &AVFilterInOut{ptr: tmpoutputs}
		} else {
			*outputs = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function avfilter_graph_segment_apply ---

// AVFilterGraphSegmentApply wraps avfilter_graph_segment_apply.
/*
  Apply all filter/link descriptions from a graph segment to the associated filtergraph.

  This functions is currently equivalent to calling the following in sequence:
  - avfilter_graph_segment_create_filters();
  - avfilter_graph_segment_apply_opts();
  - avfilter_graph_segment_init();
  - avfilter_graph_segment_link();
  failing if any of them fails. This list may be extended in the future.

  Since the above functions are idempotent, the caller may call some of them
  manually, then do some custom processing on the filtergraph, then call this
  function to do the rest.

  @param seg the filtergraph segment to process
  @param flags reserved for future use, caller must set to 0 for now
  @param[out] inputs passed to avfilter_graph_segment_link()
  @param[out] outputs passed to avfilter_graph_segment_link()

  @retval "non-negative number" success
  @retval "negative error code" failure

  @note Calling this function multiple times is safe, as it is idempotent.
*/
func AVFilterGraphSegmentApply(seg *AVFilterGraphSegment, flags int, inputs **AVFilterInOut, outputs **AVFilterInOut) (int, error) {
	var tmpseg *C.AVFilterGraphSegment
	if seg != nil {
		tmpseg = seg.ptr
	}
	var ptrinputs **C.AVFilterInOut
	var tmpinputs *C.AVFilterInOut
	var oldTmpinputs *C.AVFilterInOut
	if inputs != nil {
		innerinputs := *inputs
		if innerinputs != nil {
			tmpinputs = innerinputs.ptr
			oldTmpinputs = tmpinputs
		}
		ptrinputs = &tmpinputs
	}
	var ptroutputs **C.AVFilterInOut
	var tmpoutputs *C.AVFilterInOut
	var oldTmpoutputs *C.AVFilterInOut
	if outputs != nil {
		inneroutputs := *outputs
		if inneroutputs != nil {
			tmpoutputs = inneroutputs.ptr
			oldTmpoutputs = tmpoutputs
		}
		ptroutputs = &tmpoutputs
	}
	ret := C.avfilter_graph_segment_apply(tmpseg, C.int(flags), ptrinputs, ptroutputs)
	if tmpinputs != oldTmpinputs && inputs != nil {
		if tmpinputs != nil {
			*inputs = &AVFilterInOut{ptr: tmpinputs}
		} else {
			*inputs = nil
		}
	}
	if tmpoutputs != oldTmpoutputs && outputs != nil {
		if tmpoutputs != nil {
			*outputs = &AVFilterInOut{ptr: tmpoutputs}
		} else {
			*outputs = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function avfilter_graph_segment_free ---

// AVFilterGraphSegmentFree wraps avfilter_graph_segment_free.
/*
  Free the provided AVFilterGraphSegment and everything associated with it.

  @param seg double pointer to the AVFilterGraphSegment to be freed. NULL will
  be written to this pointer on exit from this function.

  @note
  The filter contexts (AVFilterParams.filter) are owned by AVFilterGraph rather
  than AVFilterGraphSegment, so they are not freed.
*/
func AVFilterGraphSegmentFree(seg **AVFilterGraphSegment) {
	var ptrseg **C.AVFilterGraphSegment
	var tmpseg *C.AVFilterGraphSegment
	var oldTmpseg *C.AVFilterGraphSegment
	if seg != nil {
		innerseg := *seg
		if innerseg != nil {
			tmpseg = innerseg.ptr
			oldTmpseg = tmpseg
		}
		ptrseg = &tmpseg
	}
	C.avfilter_graph_segment_free(ptrseg)
	if tmpseg != oldTmpseg && seg != nil {
		if tmpseg != nil {
			*seg = &AVFilterGraphSegment{ptr: tmpseg}
		} else {
			*seg = nil
		}
	}
}

// --- Function avfilter_graph_send_command ---

// AVFilterGraphSendCommand wraps avfilter_graph_send_command.
/*
  Send a command to one or more filter instances.

  @param graph  the filter graph
  @param target the filter(s) to which the command should be sent
                "all" sends to all filters
                otherwise it can be a filter or filter instance name
                which will send the command to all matching filters.
  @param cmd    the command to send, for handling simplicity all commands must be alphanumeric only
  @param arg    the argument for the command
  @param res    a buffer with size res_size where the filter(s) can return a response.

  @returns >=0 on success otherwise an error code.
               AVERROR(ENOSYS) on unsupported commands
*/
func AVFilterGraphSendCommand(graph *AVFilterGraph, target *CStr, cmd *CStr, arg *CStr, res *CStr, resLen int, flags int) (int, error) {
	var tmpgraph *C.AVFilterGraph
	if graph != nil {
		tmpgraph = graph.ptr
	}
	var tmptarget *C.char
	if target != nil {
		tmptarget = target.ptr
	}
	var tmpcmd *C.char
	if cmd != nil {
		tmpcmd = cmd.ptr
	}
	var tmparg *C.char
	if arg != nil {
		tmparg = arg.ptr
	}
	var tmpres *C.char
	if res != nil {
		tmpres = res.ptr
	}
	ret := C.avfilter_graph_send_command(tmpgraph, tmptarget, tmpcmd, tmparg, tmpres, C.int(resLen), C.int(flags))
	return int(ret), WrapErr(int(ret))
}

// --- Function avfilter_graph_queue_command ---

// AVFilterGraphQueueCommand wraps avfilter_graph_queue_command.
/*
  Queue a command for one or more filter instances.

  @param graph  the filter graph
  @param target the filter(s) to which the command should be sent
                "all" sends to all filters
                otherwise it can be a filter or filter instance name
                which will send the command to all matching filters.
  @param cmd    the command to sent, for handling simplicity all commands must be alphanumeric only
  @param arg    the argument for the command
  @param ts     time at which the command should be sent to the filter

  @note As this executes commands after this function returns, no return code
        from the filter is provided, also AVFILTER_CMD_FLAG_ONE is not supported.
*/
func AVFilterGraphQueueCommand(graph *AVFilterGraph, target *CStr, cmd *CStr, arg *CStr, flags int, ts float64) (int, error) {
	var tmpgraph *C.AVFilterGraph
	if graph != nil {
		tmpgraph = graph.ptr
	}
	var tmptarget *C.char
	if target != nil {
		tmptarget = target.ptr
	}
	var tmpcmd *C.char
	if cmd != nil {
		tmpcmd = cmd.ptr
	}
	var tmparg *C.char
	if arg != nil {
		tmparg = arg.ptr
	}
	ret := C.avfilter_graph_queue_command(tmpgraph, tmptarget, tmpcmd, tmparg, C.int(flags), C.double(ts))
	return int(ret), WrapErr(int(ret))
}

// --- Function avfilter_graph_dump ---

// AVFilterGraphDump wraps avfilter_graph_dump.
/*
  Dump a graph into a human-readable string representation.

  @param graph    the graph to dump
  @param options  formatting options; currently ignored
  @return  a string, or NULL in case of memory allocation failure;
           the string must be freed using av_free
*/
func AVFilterGraphDump(graph *AVFilterGraph, options *CStr) *CStr {
	var tmpgraph *C.AVFilterGraph
	if graph != nil {
		tmpgraph = graph.ptr
	}
	var tmpoptions *C.char
	if options != nil {
		tmpoptions = options.ptr
	}
	ret := C.avfilter_graph_dump(tmpgraph, tmpoptions)
	return wrapCStr(ret)
}

// --- Function avfilter_graph_request_oldest ---

// AVFilterGraphRequestOldest wraps avfilter_graph_request_oldest.
/*
  Request a frame on the oldest sink link.

  If the request returns AVERROR_EOF, try the next.

  Note that this function is not meant to be the sole scheduling mechanism
  of a filtergraph, only a convenience function to help drain a filtergraph
  in a balanced way under normal circumstances.

  Also note that AVERROR_EOF does not mean that frames did not arrive on
  some of the sinks during the process.
  When there are multiple sink links, in case the requested link
  returns an EOF, this may cause a filter to flush pending frames
  which are sent to another sink link, although unrequested.

  @return  the return value of ff_request_frame(),
           or AVERROR_EOF if all links returned AVERROR_EOF
*/
func AVFilterGraphRequestOldest(graph *AVFilterGraph) (int, error) {
	var tmpgraph *C.AVFilterGraph
	if graph != nil {
		tmpgraph = graph.ptr
	}
	ret := C.avfilter_graph_request_oldest(tmpgraph)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_buffersink_get_frame_flags ---

// AVBuffersinkGetFrameFlags wraps av_buffersink_get_frame_flags.
/*
  Get a frame with filtered data from sink and put it in frame.

  @param ctx    pointer to a buffersink or abuffersink filter context.
  @param frame  pointer to an allocated frame that will be filled with data.
                The data must be freed using av_frame_unref() / av_frame_free()
  @param flags  a combination of AV_BUFFERSINK_FLAG_* flags

  @return  >= 0 in for success, a negative AVERROR code for failure.
*/
func AVBuffersinkGetFrameFlags(ctx *AVFilterContext, frame *AVFrame, flags int) (int, error) {
	var tmpctx *C.AVFilterContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.av_buffersink_get_frame_flags(tmpctx, tmpframe, C.int(flags))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_buffersink_set_frame_size ---

// AVBuffersinkSetFrameSize wraps av_buffersink_set_frame_size.
/*
  Set the frame size for an audio buffer sink.

  All calls to av_buffersink_get_buffer_ref will return a buffer with
  exactly the specified number of samples, or AVERROR(EAGAIN) if there is
  not enough. The last buffer at EOF will be padded with 0.
*/
func AVBuffersinkSetFrameSize(ctx *AVFilterContext, frameSize uint) {
	var tmpctx *C.AVFilterContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	C.av_buffersink_set_frame_size(tmpctx, C.uint(frameSize))
}

// --- Function av_buffersink_get_type ---

// AVBuffersinkGetType wraps av_buffersink_get_type.
func AVBuffersinkGetType(ctx *AVFilterContext) AVMediaType {
	var tmpctx *C.AVFilterContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	ret := C.av_buffersink_get_type(tmpctx)
	return AVMediaType(ret)
}

// --- Function av_buffersink_get_time_base ---

// AVBuffersinkGetTimeBase wraps av_buffersink_get_time_base.
func AVBuffersinkGetTimeBase(ctx *AVFilterContext) *AVRational {
	var tmpctx *C.AVFilterContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	ret := C.av_buffersink_get_time_base(tmpctx)
	return &AVRational{value: ret}
}

// --- Function av_buffersink_get_format ---

// AVBuffersinkGetFormat wraps av_buffersink_get_format.
func AVBuffersinkGetFormat(ctx *AVFilterContext) (int, error) {
	var tmpctx *C.AVFilterContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	ret := C.av_buffersink_get_format(tmpctx)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_buffersink_get_frame_rate ---

// AVBuffersinkGetFrameRate wraps av_buffersink_get_frame_rate.
func AVBuffersinkGetFrameRate(ctx *AVFilterContext) *AVRational {
	var tmpctx *C.AVFilterContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	ret := C.av_buffersink_get_frame_rate(tmpctx)
	return &AVRational{value: ret}
}

// --- Function av_buffersink_get_w ---

// AVBuffersinkGetW wraps av_buffersink_get_w.
func AVBuffersinkGetW(ctx *AVFilterContext) (int, error) {
	var tmpctx *C.AVFilterContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	ret := C.av_buffersink_get_w(tmpctx)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_buffersink_get_h ---

// AVBuffersinkGetH wraps av_buffersink_get_h.
func AVBuffersinkGetH(ctx *AVFilterContext) (int, error) {
	var tmpctx *C.AVFilterContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	ret := C.av_buffersink_get_h(tmpctx)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_buffersink_get_sample_aspect_ratio ---

// AVBuffersinkGetSampleAspectRatio wraps av_buffersink_get_sample_aspect_ratio.
func AVBuffersinkGetSampleAspectRatio(ctx *AVFilterContext) *AVRational {
	var tmpctx *C.AVFilterContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	ret := C.av_buffersink_get_sample_aspect_ratio(tmpctx)
	return &AVRational{value: ret}
}

// --- Function av_buffersink_get_colorspace ---

// AVBuffersinkGetColorspace wraps av_buffersink_get_colorspace.
func AVBuffersinkGetColorspace(ctx *AVFilterContext) AVColorSpace {
	var tmpctx *C.AVFilterContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	ret := C.av_buffersink_get_colorspace(tmpctx)
	return AVColorSpace(ret)
}

// --- Function av_buffersink_get_color_range ---

// AVBuffersinkGetColorRange wraps av_buffersink_get_color_range.
func AVBuffersinkGetColorRange(ctx *AVFilterContext) AVColorRange {
	var tmpctx *C.AVFilterContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	ret := C.av_buffersink_get_color_range(tmpctx)
	return AVColorRange(ret)
}

// --- Function av_buffersink_get_channels ---

// AVBuffersinkGetChannels wraps av_buffersink_get_channels.
func AVBuffersinkGetChannels(ctx *AVFilterContext) (int, error) {
	var tmpctx *C.AVFilterContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	ret := C.av_buffersink_get_channels(tmpctx)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_buffersink_get_ch_layout ---

// AVBuffersinkGetChLayout wraps av_buffersink_get_ch_layout.
func AVBuffersinkGetChLayout(ctx *AVFilterContext, chLayout *AVChannelLayout) (int, error) {
	var tmpctx *C.AVFilterContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	var tmpchLayout *C.AVChannelLayout
	if chLayout != nil {
		tmpchLayout = chLayout.ptr
	}
	ret := C.av_buffersink_get_ch_layout(tmpctx, tmpchLayout)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_buffersink_get_sample_rate ---

// AVBuffersinkGetSampleRate wraps av_buffersink_get_sample_rate.
func AVBuffersinkGetSampleRate(ctx *AVFilterContext) (int, error) {
	var tmpctx *C.AVFilterContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	ret := C.av_buffersink_get_sample_rate(tmpctx)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_buffersink_get_hw_frames_ctx ---

// AVBuffersinkGetHWFramesCtx wraps av_buffersink_get_hw_frames_ctx.
func AVBuffersinkGetHWFramesCtx(ctx *AVFilterContext) *AVBufferRef {
	var tmpctx *C.AVFilterContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	ret := C.av_buffersink_get_hw_frames_ctx(tmpctx)
	var retMapped *AVBufferRef
	if ret != nil {
		retMapped = &AVBufferRef{ptr: ret}
	}
	return retMapped
}

// --- Function av_buffersink_get_side_data ---

// av_buffersink_get_side_data skipped due to pointer-to-pointer return type

// --- Function av_buffersink_get_frame ---

// AVBuffersinkGetFrame wraps av_buffersink_get_frame.
/*
  Get a frame with filtered data from sink and put it in frame.

  @param ctx pointer to a context of a buffersink or abuffersink AVFilter.
  @param frame pointer to an allocated frame that will be filled with data.
               The data must be freed using av_frame_unref() / av_frame_free()

  @return
          - >= 0 if a frame was successfully returned.
          - AVERROR(EAGAIN) if no frames are available at this point; more
            input frames must be added to the filtergraph to get more output.
          - AVERROR_EOF if there will be no more output frames on this sink.
          - A different negative AVERROR code in other failure cases.
*/
func AVBuffersinkGetFrame(ctx *AVFilterContext, frame *AVFrame) (int, error) {
	var tmpctx *C.AVFilterContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.av_buffersink_get_frame(tmpctx, tmpframe)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_buffersink_get_samples ---

// AVBuffersinkGetSamples wraps av_buffersink_get_samples.
/*
  Same as av_buffersink_get_frame(), but with the ability to specify the number
  of samples read. This function is less efficient than
  av_buffersink_get_frame(), because it copies the data around.

  @param ctx pointer to a context of the abuffersink AVFilter.
  @param frame pointer to an allocated frame that will be filled with data.
               The data must be freed using av_frame_unref() / av_frame_free()
               frame will contain exactly nb_samples audio samples, except at
               the end of stream, when it can contain less than nb_samples.

  @return The return codes have the same meaning as for
          av_buffersink_get_frame().

  @warning do not mix this function with av_buffersink_get_frame(). Use only one or
  the other with a single sink, not both.
*/
func AVBuffersinkGetSamples(ctx *AVFilterContext, frame *AVFrame, nbSamples int) (int, error) {
	var tmpctx *C.AVFilterContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.av_buffersink_get_samples(tmpctx, tmpframe, C.int(nbSamples))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_buffersrc_get_nb_failed_requests ---

// AVBuffersrcGetNbFailedRequests wraps av_buffersrc_get_nb_failed_requests.
/*
  Get the number of failed requests.

  A failed request is when the request_frame method is called while no
  frame is present in the buffer.
  The number is reset when a frame is added.
*/
func AVBuffersrcGetNbFailedRequests(bufferSrc *AVFilterContext) uint {
	var tmpbufferSrc *C.AVFilterContext
	if bufferSrc != nil {
		tmpbufferSrc = bufferSrc.ptr
	}
	ret := C.av_buffersrc_get_nb_failed_requests(tmpbufferSrc)
	return uint(ret)
}

// --- Function av_buffersrc_parameters_alloc ---

// AVBuffersrcParametersAlloc wraps av_buffersrc_parameters_alloc.
/*
  Allocate a new AVBufferSrcParameters instance. It should be freed by the
  caller with av_free().
*/
func AVBuffersrcParametersAlloc() *AVBufferSrcParameters {
	ret := C.av_buffersrc_parameters_alloc()
	var retMapped *AVBufferSrcParameters
	if ret != nil {
		retMapped = &AVBufferSrcParameters{ptr: ret}
	}
	return retMapped
}

// --- Function av_buffersrc_parameters_set ---

// AVBuffersrcParametersSet wraps av_buffersrc_parameters_set.
/*
  Initialize the buffersrc or abuffersrc filter with the provided parameters.
  This function may be called multiple times, the later calls override the
  previous ones. Some of the parameters may also be set through AVOptions, then
  whatever method is used last takes precedence.

  @param ctx an instance of the buffersrc or abuffersrc filter
  @param param the stream parameters. The frames later passed to this filter
               must conform to those parameters. All the allocated fields in
               param remain owned by the caller, libavfilter will make internal
               copies or references when necessary.
  @return 0 on success, a negative AVERROR code on failure.
*/
func AVBuffersrcParametersSet(ctx *AVFilterContext, param *AVBufferSrcParameters) (int, error) {
	var tmpctx *C.AVFilterContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	var tmpparam *C.AVBufferSrcParameters
	if param != nil {
		tmpparam = param.ptr
	}
	ret := C.av_buffersrc_parameters_set(tmpctx, tmpparam)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_buffersrc_write_frame ---

// AVBuffersrcWriteFrame wraps av_buffersrc_write_frame.
/*
  Add a frame to the buffer source.

  @param ctx   an instance of the buffersrc filter
  @param frame frame to be added. If the frame is reference counted, this
  function will make a new reference to it. Otherwise the frame data will be
  copied.

  @return 0 on success, a negative AVERROR on error

  This function is equivalent to av_buffersrc_add_frame_flags() with the
  AV_BUFFERSRC_FLAG_KEEP_REF flag.
*/
func AVBuffersrcWriteFrame(ctx *AVFilterContext, frame *AVFrame) (int, error) {
	var tmpctx *C.AVFilterContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.av_buffersrc_write_frame(tmpctx, tmpframe)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_buffersrc_add_frame ---

// AVBuffersrcAddFrame wraps av_buffersrc_add_frame.
/*
  Add a frame to the buffer source.

  @param ctx   an instance of the buffersrc filter
  @param frame frame to be added. If the frame is reference counted, this
  function will take ownership of the reference(s) and reset the frame.
  Otherwise the frame data will be copied. If this function returns an error,
  the input frame is not touched.

  @return 0 on success, a negative AVERROR on error.

  @note the difference between this function and av_buffersrc_write_frame() is
  that av_buffersrc_write_frame() creates a new reference to the input frame,
  while this function takes ownership of the reference passed to it.

  This function is equivalent to av_buffersrc_add_frame_flags() without the
  AV_BUFFERSRC_FLAG_KEEP_REF flag.
*/
func AVBuffersrcAddFrame(ctx *AVFilterContext, frame *AVFrame) (int, error) {
	var tmpctx *C.AVFilterContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.av_buffersrc_add_frame(tmpctx, tmpframe)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_buffersrc_add_frame_flags ---

// AVBuffersrcAddFrameFlags wraps av_buffersrc_add_frame_flags.
/*
  Add a frame to the buffer source.

  By default, if the frame is reference-counted, this function will take
  ownership of the reference(s) and reset the frame. This can be controlled
  using the flags.

  If this function returns an error, the input frame is not touched.

  @param buffer_src  pointer to a buffer source context
  @param frame       a frame, or NULL to mark EOF
  @param flags       a combination of AV_BUFFERSRC_FLAG_*
  @return            >= 0 in case of success, a negative AVERROR code
                     in case of failure
*/
func AVBuffersrcAddFrameFlags(bufferSrc *AVFilterContext, frame *AVFrame, flags int) (int, error) {
	var tmpbufferSrc *C.AVFilterContext
	if bufferSrc != nil {
		tmpbufferSrc = bufferSrc.ptr
	}
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.av_buffersrc_add_frame_flags(tmpbufferSrc, tmpframe, C.int(flags))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_buffersrc_close ---

// AVBuffersrcClose wraps av_buffersrc_close.
/*
  Close the buffer source after EOF.

  This is similar to passing NULL to av_buffersrc_add_frame_flags()
  except it takes the timestamp of the EOF, i.e. the timestamp of the end
  of the last frame.
*/
func AVBuffersrcClose(ctx *AVFilterContext, pts int64, flags uint) (int, error) {
	var tmpctx *C.AVFilterContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	ret := C.av_buffersrc_close(tmpctx, C.int64_t(pts), C.uint(flags))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_get_packet ---

// AVGetPacket wraps av_get_packet.
/*
  Allocate and read the payload of a packet and initialize its
  fields with default values.

  @param s    associated IO context
  @param pkt packet
  @param size desired payload size
  @return >0 (read size) if OK, AVERROR_xxx otherwise
*/
func AVGetPacket(s *AVIOContext, pkt *AVPacket, size int) (int, error) {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	var tmppkt *C.AVPacket
	if pkt != nil {
		tmppkt = pkt.ptr
	}
	ret := C.av_get_packet(tmps, tmppkt, C.int(size))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_append_packet ---

// AVAppendPacket wraps av_append_packet.
/*
  Read data and append it to the current content of the AVPacket.
  If pkt->size is 0 this is identical to av_get_packet.
  Note that this uses av_grow_packet and thus involves a realloc
  which is inefficient. Thus this function should only be used
  when there is no reasonable way to know (an upper bound of)
  the final size.

  @param s    associated IO context
  @param pkt packet
  @param size amount of data to read
  @return >0 (read size) if OK, AVERROR_xxx otherwise, previous data
          will not be lost even if an error occurs.
*/
func AVAppendPacket(s *AVIOContext, pkt *AVPacket, size int) (int, error) {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	var tmppkt *C.AVPacket
	if pkt != nil {
		tmppkt = pkt.ptr
	}
	ret := C.av_append_packet(tmps, tmppkt, C.int(size))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_disposition_from_string ---

// AVDispositionFromString wraps av_disposition_from_string.
/*
  @return The AV_DISPOSITION_* flag corresponding to disp or a negative error
          code if disp does not correspond to a known stream disposition.
*/
func AVDispositionFromString(disp *CStr) (int, error) {
	var tmpdisp *C.char
	if disp != nil {
		tmpdisp = disp.ptr
	}
	ret := C.av_disposition_from_string(tmpdisp)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_disposition_to_string ---

// AVDispositionToString wraps av_disposition_to_string.
/*
  @param disposition a combination of AV_DISPOSITION_* values
  @return The string description corresponding to the lowest set bit in
          disposition. NULL when the lowest set bit does not correspond
          to a known disposition or when disposition is 0.
*/
func AVDispositionToString(disposition int) *CStr {
	ret := C.av_disposition_to_string(C.int(disposition))
	return wrapCStr(ret)
}

// --- Function av_stream_get_parser ---

// AVStreamGetParser wraps av_stream_get_parser.
func AVStreamGetParser(s *AVStream) *AVCodecParserContext {
	var tmps *C.AVStream
	if s != nil {
		tmps = s.ptr
	}
	ret := C.av_stream_get_parser(tmps)
	var retMapped *AVCodecParserContext
	if ret != nil {
		retMapped = &AVCodecParserContext{ptr: ret}
	}
	return retMapped
}

// --- Function avformat_version ---

// AVFormatVersion wraps avformat_version.
//
//	Return the LIBAVFORMAT_VERSION_INT constant.
func AVFormatVersion() uint {
	ret := C.avformat_version()
	return uint(ret)
}

// --- Function avformat_configuration ---

// AVFormatConfiguration wraps avformat_configuration.
//
//	Return the libavformat build-time configuration.
func AVFormatConfiguration() *CStr {
	ret := C.avformat_configuration()
	return wrapCStr(ret)
}

// --- Function avformat_license ---

// AVFormatLicense wraps avformat_license.
//
//	Return the libavformat license.
func AVFormatLicense() *CStr {
	ret := C.avformat_license()
	return wrapCStr(ret)
}

// --- Function avformat_network_init ---

// AVFormatNetworkInit wraps avformat_network_init.
/*
  Do global initialization of network libraries. This is optional,
  and not recommended anymore.

  This functions only exists to work around thread-safety issues
  with older GnuTLS or OpenSSL libraries. If libavformat is linked
  to newer versions of those libraries, or if you do not use them,
  calling this function is unnecessary. Otherwise, you need to call
  this function before any other threads using them are started.

  This function will be deprecated once support for older GnuTLS and
  OpenSSL libraries is removed, and this function has no purpose
  anymore.
*/
func AVFormatNetworkInit() (int, error) {
	ret := C.avformat_network_init()
	return int(ret), WrapErr(int(ret))
}

// --- Function avformat_network_deinit ---

// AVFormatNetworkDeinit wraps avformat_network_deinit.
/*
  Undo the initialization done by avformat_network_init. Call it only
  once for each time you called avformat_network_init.
*/
func AVFormatNetworkDeinit() (int, error) {
	ret := C.avformat_network_deinit()
	return int(ret), WrapErr(int(ret))
}

// --- Function av_muxer_iterate ---

// av_muxer_iterate skipped due to opaque

// --- Function av_demuxer_iterate ---

// av_demuxer_iterate skipped due to opaque

// --- Function avformat_alloc_context ---

// AVFormatAllocContext wraps avformat_alloc_context.
/*
  Allocate an AVFormatContext.
  avformat_free_context() can be used to free the context and everything
  allocated by the framework within it.
*/
func AVFormatAllocContext() *AVFormatContext {
	ret := C.avformat_alloc_context()
	var retMapped *AVFormatContext
	if ret != nil {
		retMapped = &AVFormatContext{ptr: ret}
	}
	return retMapped
}

// --- Function avformat_free_context ---

// AVFormatFreeContext wraps avformat_free_context.
/*
  Free an AVFormatContext and all its streams.
  @param s context to free
*/
func AVFormatFreeContext(s *AVFormatContext) {
	var tmps *C.AVFormatContext
	if s != nil {
		tmps = s.ptr
	}
	C.avformat_free_context(tmps)
}

// --- Function avformat_get_class ---

// AVFormatGetClass wraps avformat_get_class.
/*
  Get the AVClass for AVFormatContext. It can be used in combination with
  AV_OPT_SEARCH_FAKE_OBJ for examining options.

  @see av_opt_find().
*/
func AVFormatGetClass() *AVClass {
	ret := C.avformat_get_class()
	var retMapped *AVClass
	if ret != nil {
		retMapped = &AVClass{ptr: ret}
	}
	return retMapped
}

// --- Function av_stream_get_class ---

// AVStreamGetClass wraps av_stream_get_class.
/*
  Get the AVClass for AVStream. It can be used in combination with
  AV_OPT_SEARCH_FAKE_OBJ for examining options.

  @see av_opt_find().
*/
func AVStreamGetClass() *AVClass {
	ret := C.av_stream_get_class()
	var retMapped *AVClass
	if ret != nil {
		retMapped = &AVClass{ptr: ret}
	}
	return retMapped
}

// --- Function av_stream_group_get_class ---

// AVStreamGroupGetClass wraps av_stream_group_get_class.
/*
  Get the AVClass for AVStreamGroup. It can be used in combination with
  AV_OPT_SEARCH_FAKE_OBJ for examining options.

  @see av_opt_find().
*/
func AVStreamGroupGetClass() *AVClass {
	ret := C.av_stream_group_get_class()
	var retMapped *AVClass
	if ret != nil {
		retMapped = &AVClass{ptr: ret}
	}
	return retMapped
}

// --- Function avformat_stream_group_name ---

// AVFormatStreamGroupName wraps avformat_stream_group_name.
//
//	@return a string identifying the stream group type, or NULL if unknown
func AVFormatStreamGroupName(_type AVStreamGroupParamsType) *CStr {
	ret := C.avformat_stream_group_name(C.enum_AVStreamGroupParamsType(_type))
	return wrapCStr(ret)
}

// --- Function avformat_stream_group_create ---

// AVFormatStreamGroupCreate wraps avformat_stream_group_create.
/*
  Add a new empty stream group to a media file.

  When demuxing, it may be called by the demuxer in read_header(). If the
  flag AVFMTCTX_NOHEADER is set in s.ctx_flags, then it may also
  be called in read_packet().

  When muxing, may be called by the user before avformat_write_header().

  User is required to call avformat_free_context() to clean up the allocation
  by avformat_stream_group_create().

  New streams can be added to the group with avformat_stream_group_add_stream().

  @param s media file handle

  @return newly created group or NULL on error.
  @see avformat_new_stream, avformat_stream_group_add_stream.
*/
func AVFormatStreamGroupCreate(s *AVFormatContext, _type AVStreamGroupParamsType, options **AVDictionary) *AVStreamGroup {
	var tmps *C.AVFormatContext
	if s != nil {
		tmps = s.ptr
	}
	var ptroptions **C.AVDictionary
	var tmpoptions *C.AVDictionary
	var oldTmpoptions *C.AVDictionary
	if options != nil {
		inneroptions := *options
		if inneroptions != nil {
			tmpoptions = inneroptions.ptr
			oldTmpoptions = tmpoptions
		}
		ptroptions = &tmpoptions
	}
	ret := C.avformat_stream_group_create(tmps, C.enum_AVStreamGroupParamsType(_type), ptroptions)
	if tmpoptions != oldTmpoptions && options != nil {
		if tmpoptions != nil {
			*options = &AVDictionary{ptr: tmpoptions}
		} else {
			*options = nil
		}
	}
	var retMapped *AVStreamGroup
	if ret != nil {
		retMapped = &AVStreamGroup{ptr: ret}
	}
	return retMapped
}

// --- Function avformat_new_stream ---

// AVFormatNewStream wraps avformat_new_stream.
/*
  Add a new stream to a media file.

  When demuxing, it is called by the demuxer in read_header(). If the
  flag AVFMTCTX_NOHEADER is set in s.ctx_flags, then it may also
  be called in read_packet().

  When muxing, should be called by the user before avformat_write_header().

  User is required to call avformat_free_context() to clean up the allocation
  by avformat_new_stream().

  @param s media file handle
  @param c unused, does nothing

  @return newly created stream or NULL on error.
*/
func AVFormatNewStream(s *AVFormatContext, c *AVCodec) *AVStream {
	var tmps *C.AVFormatContext
	if s != nil {
		tmps = s.ptr
	}
	var tmpc *C.AVCodec
	if c != nil {
		tmpc = c.ptr
	}
	ret := C.avformat_new_stream(tmps, tmpc)
	var retMapped *AVStream
	if ret != nil {
		retMapped = &AVStream{ptr: ret}
	}
	return retMapped
}

// --- Function avformat_stream_group_add_stream ---

// AVFormatStreamGroupAddStream wraps avformat_stream_group_add_stream.
/*
  Add an already allocated stream to a stream group.

  When demuxing, it may be called by the demuxer in read_header(). If the
  flag AVFMTCTX_NOHEADER is set in s.ctx_flags, then it may also
  be called in read_packet().

  When muxing, may be called by the user before avformat_write_header() after
  having allocated a new group with avformat_stream_group_create() and stream with
  avformat_new_stream().

  User is required to call avformat_free_context() to clean up the allocation
  by avformat_stream_group_add_stream().

  @param stg stream group belonging to a media file.
  @param st  stream in the media file to add to the group.

  @retval 0                 success
  @retval AVERROR(EEXIST)   the stream was already in the group
  @retval "another negative error code" legitimate errors

  @see avformat_new_stream, avformat_stream_group_create.
*/
func AVFormatStreamGroupAddStream(stg *AVStreamGroup, st *AVStream) (int, error) {
	var tmpstg *C.AVStreamGroup
	if stg != nil {
		tmpstg = stg.ptr
	}
	var tmpst *C.AVStream
	if st != nil {
		tmpst = st.ptr
	}
	ret := C.avformat_stream_group_add_stream(tmpstg, tmpst)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_new_program ---

// AVNewProgram wraps av_new_program.
func AVNewProgram(s *AVFormatContext, id int) *AVProgram {
	var tmps *C.AVFormatContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.av_new_program(tmps, C.int(id))
	var retMapped *AVProgram
	if ret != nil {
		retMapped = &AVProgram{ptr: ret}
	}
	return retMapped
}

// --- Function avformat_alloc_output_context2 ---

// AVFormatAllocOutputContext2 wraps avformat_alloc_output_context2.
/*
  Allocate an AVFormatContext for an output format.
  avformat_free_context() can be used to free the context and
  everything allocated by the framework within it.

  @param ctx           pointee is set to the created format context,
                       or to NULL in case of failure
  @param oformat       format to use for allocating the context, if NULL
                       format_name and filename are used instead
  @param format_name   the name of output format to use for allocating the
                       context, if NULL filename is used instead
  @param filename      the name of the filename to use for allocating the
                       context, may be NULL

  @return  >= 0 in case of success, a negative AVERROR code in case of
           failure
*/
func AVFormatAllocOutputContext2(ctx **AVFormatContext, oformat *AVOutputFormat, formatName *CStr, filename *CStr) (int, error) {
	var ptrctx **C.AVFormatContext
	var tmpctx *C.AVFormatContext
	var oldTmpctx *C.AVFormatContext
	if ctx != nil {
		innerctx := *ctx
		if innerctx != nil {
			tmpctx = innerctx.ptr
			oldTmpctx = tmpctx
		}
		ptrctx = &tmpctx
	}
	var tmpoformat *C.AVOutputFormat
	if oformat != nil {
		tmpoformat = oformat.ptr
	}
	var tmpformatName *C.char
	if formatName != nil {
		tmpformatName = formatName.ptr
	}
	var tmpfilename *C.char
	if filename != nil {
		tmpfilename = filename.ptr
	}
	ret := C.avformat_alloc_output_context2(ptrctx, tmpoformat, tmpformatName, tmpfilename)
	if tmpctx != oldTmpctx && ctx != nil {
		if tmpctx != nil {
			*ctx = &AVFormatContext{ptr: tmpctx}
		} else {
			*ctx = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function av_find_input_format ---

// AVFindInputFormat wraps av_find_input_format.
//
//	Find AVInputFormat based on the short name of the input format.
func AVFindInputFormat(shortName *CStr) *AVInputFormat {
	var tmpshortName *C.char
	if shortName != nil {
		tmpshortName = shortName.ptr
	}
	ret := C.av_find_input_format(tmpshortName)
	var retMapped *AVInputFormat
	if ret != nil {
		retMapped = &AVInputFormat{ptr: ret}
	}
	return retMapped
}

// --- Function av_probe_input_format ---

// AVProbeInputFormat wraps av_probe_input_format.
/*
  Guess the file format.

  @param pd        data to be probed
  @param is_opened Whether the file is already opened; determines whether
                   demuxers with or without AVFMT_NOFILE are probed.
*/
func AVProbeInputFormat(pd *AVProbeData, isOpened int) *AVInputFormat {
	var tmppd *C.AVProbeData
	if pd != nil {
		tmppd = pd.ptr
	}
	ret := C.av_probe_input_format(tmppd, C.int(isOpened))
	var retMapped *AVInputFormat
	if ret != nil {
		retMapped = &AVInputFormat{ptr: ret}
	}
	return retMapped
}

// --- Function av_probe_input_format2 ---

// av_probe_input_format2 skipped due to scoreMax (non-output primitive pointer)

// --- Function av_probe_input_format3 ---

// AVProbeInputFormat3 wraps av_probe_input_format3.
/*
  Guess the file format.

  @param is_opened Whether the file is already opened; determines whether
                   demuxers with or without AVFMT_NOFILE are probed.
  @param score_ret The score of the best detection.
*/
func AVProbeInputFormat3(pd *AVProbeData, isOpened int, scoreRet *int) *AVInputFormat {
	var tmppd *C.AVProbeData
	if pd != nil {
		tmppd = pd.ptr
	}
	ret := C.av_probe_input_format3(tmppd, C.int(isOpened), (*C.int)(unsafe.Pointer(scoreRet)))
	var retMapped *AVInputFormat
	if ret != nil {
		retMapped = &AVInputFormat{ptr: ret}
	}
	return retMapped
}

// --- Function av_probe_input_buffer2 ---

// AVProbeInputBuffer2 wraps av_probe_input_buffer2.
/*
  Probe a bytestream to determine the input format. Each time a probe returns
  with a score that is too low, the probe buffer size is increased and another
  attempt is made. When the maximum probe size is reached, the input format
  with the highest score is returned.

  @param pb             the bytestream to probe
  @param fmt            the input format is put here
  @param url            the url of the stream
  @param logctx         the log context
  @param offset         the offset within the bytestream to probe from
  @param max_probe_size the maximum probe buffer size (zero for default)

  @return the score in case of success, a negative value corresponding to an
          the maximal score is AVPROBE_SCORE_MAX
          AVERROR code otherwise
*/
func AVProbeInputBuffer2(pb *AVIOContext, fmt **AVInputFormat, url *CStr, logctx unsafe.Pointer, offset uint, maxProbeSize uint) (int, error) {
	var tmppb *C.AVIOContext
	if pb != nil {
		tmppb = pb.ptr
	}
	var ptrfmt **C.AVInputFormat
	var tmpfmt *C.AVInputFormat
	var oldTmpfmt *C.AVInputFormat
	if fmt != nil {
		innerfmt := *fmt
		if innerfmt != nil {
			tmpfmt = innerfmt.ptr
			oldTmpfmt = tmpfmt
		}
		ptrfmt = &tmpfmt
	}
	var tmpurl *C.char
	if url != nil {
		tmpurl = url.ptr
	}
	ret := C.av_probe_input_buffer2(tmppb, ptrfmt, tmpurl, logctx, C.uint(offset), C.uint(maxProbeSize))
	if tmpfmt != oldTmpfmt && fmt != nil {
		if tmpfmt != nil {
			*fmt = &AVInputFormat{ptr: tmpfmt}
		} else {
			*fmt = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function av_probe_input_buffer ---

// AVProbeInputBuffer wraps av_probe_input_buffer.
//
//	Like av_probe_input_buffer2() but returns 0 on success
func AVProbeInputBuffer(pb *AVIOContext, fmt **AVInputFormat, url *CStr, logctx unsafe.Pointer, offset uint, maxProbeSize uint) (int, error) {
	var tmppb *C.AVIOContext
	if pb != nil {
		tmppb = pb.ptr
	}
	var ptrfmt **C.AVInputFormat
	var tmpfmt *C.AVInputFormat
	var oldTmpfmt *C.AVInputFormat
	if fmt != nil {
		innerfmt := *fmt
		if innerfmt != nil {
			tmpfmt = innerfmt.ptr
			oldTmpfmt = tmpfmt
		}
		ptrfmt = &tmpfmt
	}
	var tmpurl *C.char
	if url != nil {
		tmpurl = url.ptr
	}
	ret := C.av_probe_input_buffer(tmppb, ptrfmt, tmpurl, logctx, C.uint(offset), C.uint(maxProbeSize))
	if tmpfmt != oldTmpfmt && fmt != nil {
		if tmpfmt != nil {
			*fmt = &AVInputFormat{ptr: tmpfmt}
		} else {
			*fmt = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function avformat_open_input ---

// AVFormatOpenInput wraps avformat_open_input.
/*
  Open an input stream and read the header. The codecs are not opened.
  The stream must be closed with avformat_close_input().

  @param ps       Pointer to user-supplied AVFormatContext (allocated by
                  avformat_alloc_context). May be a pointer to NULL, in
                  which case an AVFormatContext is allocated by this
                  function and written into ps.
                  Note that a user-supplied AVFormatContext will be freed
                  on failure and its pointer set to NULL.
  @param url      URL of the stream to open.
  @param fmt      If non-NULL, this parameter forces a specific input format.
                  Otherwise the format is autodetected.
  @param options  A dictionary filled with AVFormatContext and demuxer-private
                  options.
                  On return this parameter will be destroyed and replaced with
                  a dict containing options that were not found. May be NULL.

  @return 0 on success; on failure: frees ps, sets its pointer to NULL,
          and returns a negative AVERROR.

  @note If you want to use custom IO, preallocate the format context and set its pb field.
*/
func AVFormatOpenInput(ps **AVFormatContext, url *CStr, fmt *AVInputFormat, options **AVDictionary) (int, error) {
	var ptrps **C.AVFormatContext
	var tmpps *C.AVFormatContext
	var oldTmpps *C.AVFormatContext
	if ps != nil {
		innerps := *ps
		if innerps != nil {
			tmpps = innerps.ptr
			oldTmpps = tmpps
		}
		ptrps = &tmpps
	}
	var tmpurl *C.char
	if url != nil {
		tmpurl = url.ptr
	}
	var tmpfmt *C.AVInputFormat
	if fmt != nil {
		tmpfmt = fmt.ptr
	}
	var ptroptions **C.AVDictionary
	var tmpoptions *C.AVDictionary
	var oldTmpoptions *C.AVDictionary
	if options != nil {
		inneroptions := *options
		if inneroptions != nil {
			tmpoptions = inneroptions.ptr
			oldTmpoptions = tmpoptions
		}
		ptroptions = &tmpoptions
	}
	ret := C.avformat_open_input(ptrps, tmpurl, tmpfmt, ptroptions)
	if tmpps != oldTmpps && ps != nil {
		if tmpps != nil {
			*ps = &AVFormatContext{ptr: tmpps}
		} else {
			*ps = nil
		}
	}
	if tmpoptions != oldTmpoptions && options != nil {
		if tmpoptions != nil {
			*options = &AVDictionary{ptr: tmpoptions}
		} else {
			*options = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function avformat_find_stream_info ---

// AVFormatFindStreamInfo wraps avformat_find_stream_info.
/*
  Read packets of a media file to get stream information. This
  is useful for file formats with no headers such as MPEG. This
  function also computes the real framerate in case of MPEG-2 repeat
  frame mode.
  The logical file position is not changed by this function;
  examined packets may be buffered for later processing.

  @param ic media file handle
  @param options  If non-NULL, an ic.nb_streams long array of pointers to
                  dictionaries, where i-th member contains options for
                  codec corresponding to i-th stream.
                  On return each dictionary will be filled with options that were not found.
  @return >=0 if OK, AVERROR_xxx on error

  @note this function isn't guaranteed to open all the codecs, so
        options being non-empty at return is a perfectly normal behavior.

  @todo Let the user decide somehow what information is needed so that
        we do not waste time getting stuff the user does not need.
*/
func AVFormatFindStreamInfo(ic *AVFormatContext, options **AVDictionary) (int, error) {
	var tmpic *C.AVFormatContext
	if ic != nil {
		tmpic = ic.ptr
	}
	var ptroptions **C.AVDictionary
	var tmpoptions *C.AVDictionary
	var oldTmpoptions *C.AVDictionary
	if options != nil {
		inneroptions := *options
		if inneroptions != nil {
			tmpoptions = inneroptions.ptr
			oldTmpoptions = tmpoptions
		}
		ptroptions = &tmpoptions
	}
	ret := C.avformat_find_stream_info(tmpic, ptroptions)
	if tmpoptions != oldTmpoptions && options != nil {
		if tmpoptions != nil {
			*options = &AVDictionary{ptr: tmpoptions}
		} else {
			*options = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function av_find_program_from_stream ---

// AVFindProgramFromStream wraps av_find_program_from_stream.
/*
  Find the programs which belong to a given stream.

  @param ic    media file handle
  @param last  the last found program, the search will start after this
               program, or from the beginning if it is NULL
  @param s     stream index

  @return the next program which belongs to s, NULL if no program is found or
          the last program is not among the programs of ic.
*/
func AVFindProgramFromStream(ic *AVFormatContext, last *AVProgram, s int) *AVProgram {
	var tmpic *C.AVFormatContext
	if ic != nil {
		tmpic = ic.ptr
	}
	var tmplast *C.AVProgram
	if last != nil {
		tmplast = last.ptr
	}
	ret := C.av_find_program_from_stream(tmpic, tmplast, C.int(s))
	var retMapped *AVProgram
	if ret != nil {
		retMapped = &AVProgram{ptr: ret}
	}
	return retMapped
}

// --- Function av_program_add_stream_index ---

// AVProgramAddStreamIndex wraps av_program_add_stream_index.
func AVProgramAddStreamIndex(ac *AVFormatContext, progid int, idx uint) {
	var tmpac *C.AVFormatContext
	if ac != nil {
		tmpac = ac.ptr
	}
	C.av_program_add_stream_index(tmpac, C.int(progid), C.uint(idx))
}

// --- Function av_find_best_stream ---

// AVFindBestStream wraps av_find_best_stream.
/*
  Find the "best" stream in the file.
  The best stream is determined according to various heuristics as the most
  likely to be what the user expects.
  If the decoder parameter is non-NULL, av_find_best_stream will find the
  default decoder for the stream's codec; streams for which no decoder can
  be found are ignored.

  @param ic                media file handle
  @param type              stream type: video, audio, subtitles, etc.
  @param wanted_stream_nb  user-requested stream number,
                           or -1 for automatic selection
  @param related_stream    try to find a stream related (eg. in the same
                           program) to this one, or -1 if none
  @param decoder_ret       if non-NULL, returns the decoder for the
                           selected stream
  @param flags             flags; none are currently defined

  @return  the non-negative stream number in case of success,
           AVERROR_STREAM_NOT_FOUND if no stream with the requested type
           could be found,
           AVERROR_DECODER_NOT_FOUND if streams were found but no decoder

  @note  If av_find_best_stream returns successfully and decoder_ret is not
         NULL, then *decoder_ret is guaranteed to be set to a valid AVCodec.
*/
func AVFindBestStream(ic *AVFormatContext, _type AVMediaType, wantedStreamNb int, relatedStream int, decoderRet **AVCodec, flags int) (int, error) {
	var tmpic *C.AVFormatContext
	if ic != nil {
		tmpic = ic.ptr
	}
	var ptrdecoderRet **C.AVCodec
	var tmpdecoderRet *C.AVCodec
	var oldTmpdecoderRet *C.AVCodec
	if decoderRet != nil {
		innerdecoderRet := *decoderRet
		if innerdecoderRet != nil {
			tmpdecoderRet = innerdecoderRet.ptr
			oldTmpdecoderRet = tmpdecoderRet
		}
		ptrdecoderRet = &tmpdecoderRet
	}
	ret := C.av_find_best_stream(tmpic, C.enum_AVMediaType(_type), C.int(wantedStreamNb), C.int(relatedStream), ptrdecoderRet, C.int(flags))
	if tmpdecoderRet != oldTmpdecoderRet && decoderRet != nil {
		if tmpdecoderRet != nil {
			*decoderRet = &AVCodec{ptr: tmpdecoderRet}
		} else {
			*decoderRet = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function av_read_frame ---

// AVReadFrame wraps av_read_frame.
/*
  Return the next frame of a stream.
  This function returns what is stored in the file, and does not validate
  that what is there are valid frames for the decoder. It will split what is
  stored in the file into frames and return one for each call. It will not
  omit invalid data between valid frames so as to give the decoder the maximum
  information possible for decoding.

  On success, the returned packet is reference-counted (pkt->buf is set) and
  valid indefinitely. The packet must be freed with av_packet_unref() when
  it is no longer needed. For video, the packet contains exactly one frame.
  For audio, it contains an integer number of frames if each frame has
  a known fixed size (e.g. PCM or ADPCM data). If the audio frames have
  a variable size (e.g. MPEG audio), then it contains one frame.

  pkt->pts, pkt->dts and pkt->duration are always set to correct
  values in AVStream.time_base units (and guessed if the format cannot
  provide them). pkt->pts can be AV_NOPTS_VALUE if the video format
  has B-frames, so it is better to rely on pkt->dts if you do not
  decompress the payload.

  @return 0 if OK, < 0 on error or end of file. On error, pkt will be blank
          (as if it came from av_packet_alloc()).

  @note pkt will be initialized, so it may be uninitialized, but it must not
        contain data that needs to be freed.
*/
func AVReadFrame(s *AVFormatContext, pkt *AVPacket) (int, error) {
	var tmps *C.AVFormatContext
	if s != nil {
		tmps = s.ptr
	}
	var tmppkt *C.AVPacket
	if pkt != nil {
		tmppkt = pkt.ptr
	}
	ret := C.av_read_frame(tmps, tmppkt)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_seek_frame ---

// AVSeekFrame wraps av_seek_frame.
/*
  Seek to the keyframe at timestamp.
  'timestamp' in 'stream_index'.

  @param s            media file handle
  @param stream_index If stream_index is (-1), a default stream is selected,
                      and timestamp is automatically converted from
                      AV_TIME_BASE units to the stream specific time_base.
  @param timestamp    Timestamp in AVStream.time_base units or, if no stream
                      is specified, in AV_TIME_BASE units.
  @param flags        flags which select direction and seeking mode

  @return >= 0 on success
*/
func AVSeekFrame(s *AVFormatContext, streamIndex int, timestamp int64, flags int) (int, error) {
	var tmps *C.AVFormatContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.av_seek_frame(tmps, C.int(streamIndex), C.int64_t(timestamp), C.int(flags))
	return int(ret), WrapErr(int(ret))
}

// --- Function avformat_seek_file ---

// AVFormatSeekFile wraps avformat_seek_file.
/*
  Seek to timestamp ts.
  Seeking will be done so that the point from which all active streams
  can be presented successfully will be closest to ts and within min/max_ts.
  Active streams are all streams that have AVStream.discard < AVDISCARD_ALL.

  If flags contain AVSEEK_FLAG_BYTE, then all timestamps are in bytes and
  are the file position (this may not be supported by all demuxers).
  If flags contain AVSEEK_FLAG_FRAME, then all timestamps are in frames
  in the stream with stream_index (this may not be supported by all demuxers).
  Otherwise all timestamps are in units of the stream selected by stream_index
  or if stream_index is -1, in AV_TIME_BASE units.
  If flags contain AVSEEK_FLAG_ANY, then non-keyframes are treated as
  keyframes (this may not be supported by all demuxers).
  If flags contain AVSEEK_FLAG_BACKWARD, it is ignored.

  @param s            media file handle
  @param stream_index index of the stream which is used as time base reference
  @param min_ts       smallest acceptable timestamp
  @param ts           target timestamp
  @param max_ts       largest acceptable timestamp
  @param flags        flags
  @return >=0 on success, error code otherwise

  @note This is part of the new seek API which is still under construction.
*/
func AVFormatSeekFile(s *AVFormatContext, streamIndex int, minTs int64, ts int64, maxTs int64, flags int) (int, error) {
	var tmps *C.AVFormatContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.avformat_seek_file(tmps, C.int(streamIndex), C.int64_t(minTs), C.int64_t(ts), C.int64_t(maxTs), C.int(flags))
	return int(ret), WrapErr(int(ret))
}

// --- Function avformat_flush ---

// AVFormatFlush wraps avformat_flush.
/*
  Discard all internally buffered data. This can be useful when dealing with
  discontinuities in the byte stream. Generally works only with formats that
  can resync. This includes headerless formats like MPEG-TS/TS but should also
  work with NUT, Ogg and in a limited way AVI for example.

  The set of streams, the detected duration, stream parameters and codecs do
  not change when calling this function. If you want a complete reset, it's
  better to open a new AVFormatContext.

  This does not flush the AVIOContext (s->pb). If necessary, call
  avio_flush(s->pb) before calling this function.

  @param s media file handle
  @return >=0 on success, error code otherwise
*/
func AVFormatFlush(s *AVFormatContext) (int, error) {
	var tmps *C.AVFormatContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.avformat_flush(tmps)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_read_play ---

// AVReadPlay wraps av_read_play.
/*
  Start playing a network-based stream (e.g. RTSP stream) at the
  current position.
*/
func AVReadPlay(s *AVFormatContext) (int, error) {
	var tmps *C.AVFormatContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.av_read_play(tmps)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_read_pause ---

// AVReadPause wraps av_read_pause.
/*
  Pause a network-based stream (e.g. RTSP stream).

  Use av_read_play() to resume it.
*/
func AVReadPause(s *AVFormatContext) (int, error) {
	var tmps *C.AVFormatContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.av_read_pause(tmps)
	return int(ret), WrapErr(int(ret))
}

// --- Function avformat_close_input ---

// AVFormatCloseInput wraps avformat_close_input.
/*
  Close an opened input AVFormatContext. Free it and all its contents
  and set *s to NULL.
*/
func AVFormatCloseInput(s **AVFormatContext) {
	var ptrs **C.AVFormatContext
	var tmps *C.AVFormatContext
	var oldTmps *C.AVFormatContext
	if s != nil {
		inners := *s
		if inners != nil {
			tmps = inners.ptr
			oldTmps = tmps
		}
		ptrs = &tmps
	}
	C.avformat_close_input(ptrs)
	if tmps != oldTmps && s != nil {
		if tmps != nil {
			*s = &AVFormatContext{ptr: tmps}
		} else {
			*s = nil
		}
	}
}

// --- Function avformat_write_header ---

// AVFormatWriteHeader wraps avformat_write_header.
/*
  Allocate the stream private data and write the stream header to
  an output media file.

  @param s        Media file handle, must be allocated with
                  avformat_alloc_context().
                  Its \ref AVFormatContext.oformat "oformat" field must be set
                  to the desired output format;
                  Its \ref AVFormatContext.pb "pb" field must be set to an
                  already opened ::AVIOContext.
  @param options  An ::AVDictionary filled with AVFormatContext and
                  muxer-private options.
                  On return this parameter will be destroyed and replaced with
                  a dict containing options that were not found. May be NULL.

  @retval AVSTREAM_INIT_IN_WRITE_HEADER On success, if the codec had not already been
                                        fully initialized in avformat_init_output().
  @retval AVSTREAM_INIT_IN_INIT_OUTPUT  On success, if the codec had already been fully
                                        initialized in avformat_init_output().
  @retval AVERROR                       A negative AVERROR on failure.

  @see av_opt_find, av_dict_set, avio_open, av_oformat_next, avformat_init_output.
*/
func AVFormatWriteHeader(s *AVFormatContext, options **AVDictionary) (int, error) {
	var tmps *C.AVFormatContext
	if s != nil {
		tmps = s.ptr
	}
	var ptroptions **C.AVDictionary
	var tmpoptions *C.AVDictionary
	var oldTmpoptions *C.AVDictionary
	if options != nil {
		inneroptions := *options
		if inneroptions != nil {
			tmpoptions = inneroptions.ptr
			oldTmpoptions = tmpoptions
		}
		ptroptions = &tmpoptions
	}
	ret := C.avformat_write_header(tmps, ptroptions)
	if tmpoptions != oldTmpoptions && options != nil {
		if tmpoptions != nil {
			*options = &AVDictionary{ptr: tmpoptions}
		} else {
			*options = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function avformat_init_output ---

// AVFormatInitOutput wraps avformat_init_output.
/*
  Allocate the stream private data and initialize the codec, but do not write the header.
  May optionally be used before avformat_write_header() to initialize stream parameters
  before actually writing the header.
  If using this function, do not pass the same options to avformat_write_header().

  @param s        Media file handle, must be allocated with
                  avformat_alloc_context().
                  Its \ref AVFormatContext.oformat "oformat" field must be set
                  to the desired output format;
                  Its \ref AVFormatContext.pb "pb" field must be set to an
                  already opened ::AVIOContext.
  @param options  An ::AVDictionary filled with AVFormatContext and
                  muxer-private options.
                  On return this parameter will be destroyed and replaced with
                  a dict containing options that were not found. May be NULL.

  @retval AVSTREAM_INIT_IN_WRITE_HEADER On success, if the codec requires
                                        avformat_write_header to fully initialize.
  @retval AVSTREAM_INIT_IN_INIT_OUTPUT  On success, if the codec has been fully
                                        initialized.
  @retval AVERROR                       Anegative AVERROR on failure.

  @see av_opt_find, av_dict_set, avio_open, av_oformat_next, avformat_write_header.
*/
func AVFormatInitOutput(s *AVFormatContext, options **AVDictionary) (int, error) {
	var tmps *C.AVFormatContext
	if s != nil {
		tmps = s.ptr
	}
	var ptroptions **C.AVDictionary
	var tmpoptions *C.AVDictionary
	var oldTmpoptions *C.AVDictionary
	if options != nil {
		inneroptions := *options
		if inneroptions != nil {
			tmpoptions = inneroptions.ptr
			oldTmpoptions = tmpoptions
		}
		ptroptions = &tmpoptions
	}
	ret := C.avformat_init_output(tmps, ptroptions)
	if tmpoptions != oldTmpoptions && options != nil {
		if tmpoptions != nil {
			*options = &AVDictionary{ptr: tmpoptions}
		} else {
			*options = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function av_write_frame ---

// AVWriteFrame wraps av_write_frame.
/*
  Write a packet to an output media file.

  This function passes the packet directly to the muxer, without any buffering
  or reordering. The caller is responsible for correctly interleaving the
  packets if the format requires it. Callers that want libavformat to handle
  the interleaving should call av_interleaved_write_frame() instead of this
  function.

  @param s media file handle
  @param pkt The packet containing the data to be written. Note that unlike
             av_interleaved_write_frame(), this function does not take
             ownership of the packet passed to it (though some muxers may make
             an internal reference to the input packet).
             <br>
             This parameter can be NULL (at any time, not just at the end), in
             order to immediately flush data buffered within the muxer, for
             muxers that buffer up data internally before writing it to the
             output.
             <br>
             Packet's @ref AVPacket.stream_index "stream_index" field must be
             set to the index of the corresponding stream in @ref
             AVFormatContext.streams "s->streams".
             <br>
             The timestamps (@ref AVPacket.pts "pts", @ref AVPacket.dts "dts")
             must be set to correct values in the stream's timebase (unless the
             output format is flagged with the AVFMT_NOTIMESTAMPS flag, then
             they can be set to AV_NOPTS_VALUE).
             The dts for subsequent packets passed to this function must be strictly
             increasing when compared in their respective timebases (unless the
             output format is flagged with the AVFMT_TS_NONSTRICT, then they
             merely have to be nondecreasing).  @ref AVPacket.duration
             "duration") should also be set if known.
  @return < 0 on error, = 0 if OK, 1 if flushed and there is no more data to flush

  @see av_interleaved_write_frame()
*/
func AVWriteFrame(s *AVFormatContext, pkt *AVPacket) (int, error) {
	var tmps *C.AVFormatContext
	if s != nil {
		tmps = s.ptr
	}
	var tmppkt *C.AVPacket
	if pkt != nil {
		tmppkt = pkt.ptr
	}
	ret := C.av_write_frame(tmps, tmppkt)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_interleaved_write_frame ---

// AVInterleavedWriteFrame wraps av_interleaved_write_frame.
/*
  Write a packet to an output media file ensuring correct interleaving.

  This function will buffer the packets internally as needed to make sure the
  packets in the output file are properly interleaved, usually ordered by
  increasing dts. Callers doing their own interleaving should call
  av_write_frame() instead of this function.

  Using this function instead of av_write_frame() can give muxers advance
  knowledge of future packets, improving e.g. the behaviour of the mp4
  muxer for VFR content in fragmenting mode.

  @param s media file handle
  @param pkt The packet containing the data to be written.
             <br>
             If the packet is reference-counted, this function will take
             ownership of this reference and unreference it later when it sees
             fit. If the packet is not reference-counted, libavformat will
             make a copy.
             The returned packet will be blank (as if returned from
             av_packet_alloc()), even on error.
             <br>
             This parameter can be NULL (at any time, not just at the end), to
             flush the interleaving queues.
             <br>
             Packet's @ref AVPacket.stream_index "stream_index" field must be
             set to the index of the corresponding stream in @ref
             AVFormatContext.streams "s->streams".
             <br>
             The timestamps (@ref AVPacket.pts "pts", @ref AVPacket.dts "dts")
             must be set to correct values in the stream's timebase (unless the
             output format is flagged with the AVFMT_NOTIMESTAMPS flag, then
             they can be set to AV_NOPTS_VALUE).
             The dts for subsequent packets in one stream must be strictly
             increasing (unless the output format is flagged with the
             AVFMT_TS_NONSTRICT, then they merely have to be nondecreasing).
             @ref AVPacket.duration "duration" should also be set if known.

  @return 0 on success, a negative AVERROR on error.

  @see av_write_frame(), AVFormatContext.max_interleave_delta
*/
func AVInterleavedWriteFrame(s *AVFormatContext, pkt *AVPacket) (int, error) {
	var tmps *C.AVFormatContext
	if s != nil {
		tmps = s.ptr
	}
	var tmppkt *C.AVPacket
	if pkt != nil {
		tmppkt = pkt.ptr
	}
	ret := C.av_interleaved_write_frame(tmps, tmppkt)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_write_uncoded_frame ---

// AVWriteUncodedFrame wraps av_write_uncoded_frame.
/*
  Write an uncoded frame to an output media file.

  The frame must be correctly interleaved according to the container
  specification; if not, av_interleaved_write_uncoded_frame() must be used.

  See av_interleaved_write_uncoded_frame() for details.
*/
func AVWriteUncodedFrame(s *AVFormatContext, streamIndex int, frame *AVFrame) (int, error) {
	var tmps *C.AVFormatContext
	if s != nil {
		tmps = s.ptr
	}
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.av_write_uncoded_frame(tmps, C.int(streamIndex), tmpframe)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_interleaved_write_uncoded_frame ---

// AVInterleavedWriteUncodedFrame wraps av_interleaved_write_uncoded_frame.
/*
  Write an uncoded frame to an output media file.

  If the muxer supports it, this function makes it possible to write an AVFrame
  structure directly, without encoding it into a packet.
  It is mostly useful for devices and similar special muxers that use raw
  video or PCM data and will not serialize it into a byte stream.

  To test whether it is possible to use it with a given muxer and stream,
  use av_write_uncoded_frame_query().

  The caller gives up ownership of the frame and must not access it
  afterwards.

  @return  >=0 for success, a negative code on error
*/
func AVInterleavedWriteUncodedFrame(s *AVFormatContext, streamIndex int, frame *AVFrame) (int, error) {
	var tmps *C.AVFormatContext
	if s != nil {
		tmps = s.ptr
	}
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.av_interleaved_write_uncoded_frame(tmps, C.int(streamIndex), tmpframe)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_write_uncoded_frame_query ---

// AVWriteUncodedFrameQuery wraps av_write_uncoded_frame_query.
/*
  Test whether a muxer supports uncoded frame.

  @return  >=0 if an uncoded frame can be written to that muxer and stream,
           <0 if not
*/
func AVWriteUncodedFrameQuery(s *AVFormatContext, streamIndex int) (int, error) {
	var tmps *C.AVFormatContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.av_write_uncoded_frame_query(tmps, C.int(streamIndex))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_write_trailer ---

// AVWriteTrailer wraps av_write_trailer.
/*
  Write the stream trailer to an output media file and free the
  file private data.

  May only be called after a successful call to avformat_write_header.

  @param s media file handle
  @return 0 if OK, AVERROR_xxx on error
*/
func AVWriteTrailer(s *AVFormatContext) (int, error) {
	var tmps *C.AVFormatContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.av_write_trailer(tmps)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_guess_format ---

// AVGuessFormat wraps av_guess_format.
/*
  Return the output format in the list of registered output formats
  which best matches the provided parameters, or return NULL if
  there is no match.

  @param short_name if non-NULL checks if short_name matches with the
                    names of the registered formats
  @param filename   if non-NULL checks if filename terminates with the
                    extensions of the registered formats
  @param mime_type  if non-NULL checks if mime_type matches with the
                    MIME type of the registered formats
*/
func AVGuessFormat(shortName *CStr, filename *CStr, mimeType *CStr) *AVOutputFormat {
	var tmpshortName *C.char
	if shortName != nil {
		tmpshortName = shortName.ptr
	}
	var tmpfilename *C.char
	if filename != nil {
		tmpfilename = filename.ptr
	}
	var tmpmimeType *C.char
	if mimeType != nil {
		tmpmimeType = mimeType.ptr
	}
	ret := C.av_guess_format(tmpshortName, tmpfilename, tmpmimeType)
	var retMapped *AVOutputFormat
	if ret != nil {
		retMapped = &AVOutputFormat{ptr: ret}
	}
	return retMapped
}

// --- Function av_guess_codec ---

// AVGuessCodec wraps av_guess_codec.
//
//	Guess the codec ID based upon muxer and filename.
func AVGuessCodec(fmt *AVOutputFormat, shortName *CStr, filename *CStr, mimeType *CStr, _type AVMediaType) AVCodecID {
	var tmpfmt *C.AVOutputFormat
	if fmt != nil {
		tmpfmt = fmt.ptr
	}
	var tmpshortName *C.char
	if shortName != nil {
		tmpshortName = shortName.ptr
	}
	var tmpfilename *C.char
	if filename != nil {
		tmpfilename = filename.ptr
	}
	var tmpmimeType *C.char
	if mimeType != nil {
		tmpmimeType = mimeType.ptr
	}
	ret := C.av_guess_codec(tmpfmt, tmpshortName, tmpfilename, tmpmimeType, C.enum_AVMediaType(_type))
	return AVCodecID(ret)
}

// --- Function av_get_output_timestamp ---

// av_get_output_timestamp skipped due to dts (non-output primitive pointer)

// --- Function av_hex_dump ---

// av_hex_dump skipped due to f (non-output primitive pointer)

// --- Function av_hex_dump_log ---

// AVHexDumpLog wraps av_hex_dump_log.
/*
  Send a nice hexadecimal dump of a buffer to the log.

  @param avcl A pointer to an arbitrary struct of which the first field is a
  pointer to an AVClass struct.
  @param level The importance level of the message, lower values signifying
  higher importance.
  @param buf buffer
  @param size buffer size

  @see av_hex_dump, av_pkt_dump2, av_pkt_dump_log2
*/
func AVHexDumpLog(avcl unsafe.Pointer, level int, buf unsafe.Pointer, size int) {
	C.av_hex_dump_log(avcl, C.int(level), (*C.uint8_t)(buf), C.int(size))
}

// --- Function av_pkt_dump2 ---

// av_pkt_dump2 skipped due to f (non-output primitive pointer)

// --- Function av_pkt_dump_log2 ---

// AVPktDumpLog2 wraps av_pkt_dump_log2.
/*
  Send a nice dump of a packet to the log.

  @param avcl A pointer to an arbitrary struct of which the first field is a
  pointer to an AVClass struct.
  @param level The importance level of the message, lower values signifying
  higher importance.
  @param pkt packet to dump
  @param dump_payload True if the payload must be displayed, too.
  @param st AVStream that the packet belongs to
*/
func AVPktDumpLog2(avcl unsafe.Pointer, level int, pkt *AVPacket, dumpPayload int, st *AVStream) {
	var tmppkt *C.AVPacket
	if pkt != nil {
		tmppkt = pkt.ptr
	}
	var tmpst *C.AVStream
	if st != nil {
		tmpst = st.ptr
	}
	C.av_pkt_dump_log2(avcl, C.int(level), tmppkt, C.int(dumpPayload), tmpst)
}

// --- Function av_codec_get_id ---

// AVCodecGetId wraps av_codec_get_id.
/*
  Get the AVCodecID for the given codec tag tag.
  If no codec id is found returns AV_CODEC_ID_NONE.

  @param tags list of supported codec_id-codec_tag pairs, as stored
  in AVInputFormat.codec_tag and AVOutputFormat.codec_tag
  @param tag  codec tag to match to a codec ID
*/
func AVCodecGetId(tags **AVCodecTag, tag uint) AVCodecID {
	var ptrtags **C.struct_AVCodecTag
	var tmptags *C.struct_AVCodecTag
	var oldTmptags *C.struct_AVCodecTag
	if tags != nil {
		innertags := *tags
		if innertags != nil {
			tmptags = innertags.ptr
			oldTmptags = tmptags
		}
		ptrtags = &tmptags
	}
	ret := C.av_codec_get_id(ptrtags, C.uint(tag))
	if tmptags != oldTmptags && tags != nil {
		if tmptags != nil {
			*tags = &AVCodecTag{ptr: tmptags}
		} else {
			*tags = nil
		}
	}
	return AVCodecID(ret)
}

// --- Function av_codec_get_tag ---

// AVCodecGetTag wraps av_codec_get_tag.
/*
  Get the codec tag for the given codec id id.
  If no codec tag is found returns 0.

  @param tags list of supported codec_id-codec_tag pairs, as stored
  in AVInputFormat.codec_tag and AVOutputFormat.codec_tag
  @param id   codec ID to match to a codec tag
*/
func AVCodecGetTag(tags **AVCodecTag, id AVCodecID) uint {
	var ptrtags **C.struct_AVCodecTag
	var tmptags *C.struct_AVCodecTag
	var oldTmptags *C.struct_AVCodecTag
	if tags != nil {
		innertags := *tags
		if innertags != nil {
			tmptags = innertags.ptr
			oldTmptags = tmptags
		}
		ptrtags = &tmptags
	}
	ret := C.av_codec_get_tag(ptrtags, C.enum_AVCodecID(id))
	if tmptags != oldTmptags && tags != nil {
		if tmptags != nil {
			*tags = &AVCodecTag{ptr: tmptags}
		} else {
			*tags = nil
		}
	}
	return uint(ret)
}

// --- Function av_codec_get_tag2 ---

// av_codec_get_tag2 skipped due to tag (non-output primitive pointer)

// --- Function av_find_default_stream_index ---

// AVFindDefaultStreamIndex wraps av_find_default_stream_index.
func AVFindDefaultStreamIndex(s *AVFormatContext) (int, error) {
	var tmps *C.AVFormatContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.av_find_default_stream_index(tmps)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_index_search_timestamp ---

// AVIndexSearchTimestamp wraps av_index_search_timestamp.
/*
  Get the index for a specific timestamp.

  @param st        stream that the timestamp belongs to
  @param timestamp timestamp to retrieve the index for
  @param flags if AVSEEK_FLAG_BACKWARD then the returned index will correspond
                  to the timestamp which is <= the requested one, if backward
                  is 0, then it will be >=
               if AVSEEK_FLAG_ANY seek to any frame, only keyframes otherwise
  @return < 0 if no such timestamp could be found
*/
func AVIndexSearchTimestamp(st *AVStream, timestamp int64, flags int) (int, error) {
	var tmpst *C.AVStream
	if st != nil {
		tmpst = st.ptr
	}
	ret := C.av_index_search_timestamp(tmpst, C.int64_t(timestamp), C.int(flags))
	return int(ret), WrapErr(int(ret))
}

// --- Function avformat_index_get_entries_count ---

// AVFormatIndexGetEntriesCount wraps avformat_index_get_entries_count.
/*
  Get the index entry count for the given AVStream.

  @param st stream
  @return the number of index entries in the stream
*/
func AVFormatIndexGetEntriesCount(st *AVStream) (int, error) {
	var tmpst *C.AVStream
	if st != nil {
		tmpst = st.ptr
	}
	ret := C.avformat_index_get_entries_count(tmpst)
	return int(ret), WrapErr(int(ret))
}

// --- Function avformat_index_get_entry ---

// AVFormatIndexGetEntry wraps avformat_index_get_entry.
/*
  Get the AVIndexEntry corresponding to the given index.

  @param st          Stream containing the requested AVIndexEntry.
  @param idx         The desired index.
  @return A pointer to the requested AVIndexEntry if it exists, NULL otherwise.

  @note The pointer returned by this function is only guaranteed to be valid
        until any function that takes the stream or the parent AVFormatContext
        as input argument is called.
*/
func AVFormatIndexGetEntry(st *AVStream, idx int) *AVIndexEntry {
	var tmpst *C.AVStream
	if st != nil {
		tmpst = st.ptr
	}
	ret := C.avformat_index_get_entry(tmpst, C.int(idx))
	var retMapped *AVIndexEntry
	if ret != nil {
		retMapped = &AVIndexEntry{ptr: ret}
	}
	return retMapped
}

// --- Function avformat_index_get_entry_from_timestamp ---

// AVFormatIndexGetEntryFromTimestamp wraps avformat_index_get_entry_from_timestamp.
/*
  Get the AVIndexEntry corresponding to the given timestamp.

  @param st          Stream containing the requested AVIndexEntry.
  @param wanted_timestamp   Timestamp to retrieve the index entry for.
  @param flags       If AVSEEK_FLAG_BACKWARD then the returned entry will correspond
                     to the timestamp which is <= the requested one, if backward
                     is 0, then it will be >=
                     if AVSEEK_FLAG_ANY seek to any frame, only keyframes otherwise.
  @return A pointer to the requested AVIndexEntry if it exists, NULL otherwise.

  @note The pointer returned by this function is only guaranteed to be valid
        until any function that takes the stream or the parent AVFormatContext
        as input argument is called.
*/
func AVFormatIndexGetEntryFromTimestamp(st *AVStream, wantedTimestamp int64, flags int) *AVIndexEntry {
	var tmpst *C.AVStream
	if st != nil {
		tmpst = st.ptr
	}
	ret := C.avformat_index_get_entry_from_timestamp(tmpst, C.int64_t(wantedTimestamp), C.int(flags))
	var retMapped *AVIndexEntry
	if ret != nil {
		retMapped = &AVIndexEntry{ptr: ret}
	}
	return retMapped
}

// --- Function av_add_index_entry ---

// AVAddIndexEntry wraps av_add_index_entry.
/*
  Add an index entry into a sorted list. Update the entry if the list
  already contains it.

  @param timestamp timestamp in the time base of the given stream
*/
func AVAddIndexEntry(st *AVStream, pos int64, timestamp int64, size int, distance int, flags int) (int, error) {
	var tmpst *C.AVStream
	if st != nil {
		tmpst = st.ptr
	}
	ret := C.av_add_index_entry(tmpst, C.int64_t(pos), C.int64_t(timestamp), C.int(size), C.int(distance), C.int(flags))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_url_split ---

// AVUrlSplit wraps av_url_split.
/*
  Split a URL string into components.

  The pointers to buffers for storing individual components may be null,
  in order to ignore that component. Buffers for components not found are
  set to empty strings. If the port is not found, it is set to a negative
  value.

  @param proto the buffer for the protocol
  @param proto_size the size of the proto buffer
  @param authorization the buffer for the authorization
  @param authorization_size the size of the authorization buffer
  @param hostname the buffer for the host name
  @param hostname_size the size of the hostname buffer
  @param port_ptr a pointer to store the port number in
  @param path the buffer for the path
  @param path_size the size of the path buffer
  @param url the URL to split
*/
func AVUrlSplit(proto *CStr, protoSize int, authorization *CStr, authorizationSize int, hostname *CStr, hostnameSize int, portPtr *int, path *CStr, pathSize int, url *CStr) {
	var tmpproto *C.char
	if proto != nil {
		tmpproto = proto.ptr
	}
	var tmpauthorization *C.char
	if authorization != nil {
		tmpauthorization = authorization.ptr
	}
	var tmphostname *C.char
	if hostname != nil {
		tmphostname = hostname.ptr
	}
	var tmppath *C.char
	if path != nil {
		tmppath = path.ptr
	}
	var tmpurl *C.char
	if url != nil {
		tmpurl = url.ptr
	}
	C.av_url_split(tmpproto, C.int(protoSize), tmpauthorization, C.int(authorizationSize), tmphostname, C.int(hostnameSize), (*C.int)(unsafe.Pointer(portPtr)), tmppath, C.int(pathSize), tmpurl)
}

// --- Function av_dump_format ---

// AVDumpFormat wraps av_dump_format.
/*
  Print detailed information about the input or output format, such as
  duration, bitrate, streams, container, programs, metadata, side data,
  codec and time base.

  @param ic        the context to analyze
  @param index     index of the stream to dump information about
  @param url       the URL to print, such as source or destination file
  @param is_output Select whether the specified context is an input(0) or output(1)
*/
func AVDumpFormat(ic *AVFormatContext, index int, url *CStr, isOutput int) {
	var tmpic *C.AVFormatContext
	if ic != nil {
		tmpic = ic.ptr
	}
	var tmpurl *C.char
	if url != nil {
		tmpurl = url.ptr
	}
	C.av_dump_format(tmpic, C.int(index), tmpurl, C.int(isOutput))
}

// --- Function av_get_frame_filename2 ---

// AVGetFrameFilename2 wraps av_get_frame_filename2.
/*
  Return in 'buf' the path with '%d' replaced by a number.

  Also handles the '%0nd' format where 'n' is the total number
  of digits and '%%'.

  @param buf destination buffer
  @param buf_size destination buffer size
  @param path numbered sequence string
  @param number frame number
  @param flags AV_FRAME_FILENAME_FLAGS_*
  @return 0 if OK, -1 on format error
*/
func AVGetFrameFilename2(buf *CStr, bufSize int, path *CStr, number int, flags int) (int, error) {
	var tmpbuf *C.char
	if buf != nil {
		tmpbuf = buf.ptr
	}
	var tmppath *C.char
	if path != nil {
		tmppath = path.ptr
	}
	ret := C.av_get_frame_filename2(tmpbuf, C.int(bufSize), tmppath, C.int(number), C.int(flags))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_get_frame_filename ---

// AVGetFrameFilename wraps av_get_frame_filename.
func AVGetFrameFilename(buf *CStr, bufSize int, path *CStr, number int) (int, error) {
	var tmpbuf *C.char
	if buf != nil {
		tmpbuf = buf.ptr
	}
	var tmppath *C.char
	if path != nil {
		tmppath = path.ptr
	}
	ret := C.av_get_frame_filename(tmpbuf, C.int(bufSize), tmppath, C.int(number))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_filename_number_test ---

// AVFilenameNumberTest wraps av_filename_number_test.
/*
  Check whether filename actually is a numbered sequence generator.

  @param filename possible numbered sequence string
  @return 1 if a valid numbered sequence string, 0 otherwise
*/
func AVFilenameNumberTest(filename *CStr) (int, error) {
	var tmpfilename *C.char
	if filename != nil {
		tmpfilename = filename.ptr
	}
	ret := C.av_filename_number_test(tmpfilename)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_sdp_create ---

// av_sdp_create skipped due to ac

// --- Function av_match_ext ---

// AVMatchExt wraps av_match_ext.
/*
  Return a positive value if the given filename has one of the given
  extensions, 0 otherwise.

  @param filename   file name to check against the given extensions
  @param extensions a comma-separated list of filename extensions
*/
func AVMatchExt(filename *CStr, extensions *CStr) (int, error) {
	var tmpfilename *C.char
	if filename != nil {
		tmpfilename = filename.ptr
	}
	var tmpextensions *C.char
	if extensions != nil {
		tmpextensions = extensions.ptr
	}
	ret := C.av_match_ext(tmpfilename, tmpextensions)
	return int(ret), WrapErr(int(ret))
}

// --- Function avformat_query_codec ---

// AVFormatQueryCodec wraps avformat_query_codec.
/*
  Test if the given container can store a codec.

  @param ofmt           container to check for compatibility
  @param codec_id       codec to potentially store in container
  @param std_compliance standards compliance level, one of FF_COMPLIANCE_*

  @return 1 if codec with ID codec_id can be stored in ofmt, 0 if it cannot.
          A negative number if this information is not available.
*/
func AVFormatQueryCodec(ofmt *AVOutputFormat, codecId AVCodecID, stdCompliance int) (int, error) {
	var tmpofmt *C.AVOutputFormat
	if ofmt != nil {
		tmpofmt = ofmt.ptr
	}
	ret := C.avformat_query_codec(tmpofmt, C.enum_AVCodecID(codecId), C.int(stdCompliance))
	return int(ret), WrapErr(int(ret))
}

// --- Function avformat_get_riff_video_tags ---

// AVFormatGetRiffVideoTags wraps avformat_get_riff_video_tags.
//
//	@return the table mapping RIFF FourCCs for video to libavcodec AVCodecID.
func AVFormatGetRiffVideoTags() *AVCodecTag {
	ret := C.avformat_get_riff_video_tags()
	var retMapped *AVCodecTag
	if ret != nil {
		retMapped = &AVCodecTag{ptr: ret}
	}
	return retMapped
}

// --- Function avformat_get_riff_audio_tags ---

// AVFormatGetRiffAudioTags wraps avformat_get_riff_audio_tags.
//
//	@return the table mapping RIFF FourCCs for audio to AVCodecID.
func AVFormatGetRiffAudioTags() *AVCodecTag {
	ret := C.avformat_get_riff_audio_tags()
	var retMapped *AVCodecTag
	if ret != nil {
		retMapped = &AVCodecTag{ptr: ret}
	}
	return retMapped
}

// --- Function avformat_get_mov_video_tags ---

// AVFormatGetMovVideoTags wraps avformat_get_mov_video_tags.
//
//	@return the table mapping MOV FourCCs for video to libavcodec AVCodecID.
func AVFormatGetMovVideoTags() *AVCodecTag {
	ret := C.avformat_get_mov_video_tags()
	var retMapped *AVCodecTag
	if ret != nil {
		retMapped = &AVCodecTag{ptr: ret}
	}
	return retMapped
}

// --- Function avformat_get_mov_audio_tags ---

// AVFormatGetMovAudioTags wraps avformat_get_mov_audio_tags.
//
//	@return the table mapping MOV FourCCs for audio to AVCodecID.
func AVFormatGetMovAudioTags() *AVCodecTag {
	ret := C.avformat_get_mov_audio_tags()
	var retMapped *AVCodecTag
	if ret != nil {
		retMapped = &AVCodecTag{ptr: ret}
	}
	return retMapped
}

// --- Function av_guess_sample_aspect_ratio ---

// AVGuessSampleAspectRatio wraps av_guess_sample_aspect_ratio.
/*
  Guess the sample aspect ratio of a frame, based on both the stream and the
  frame aspect ratio.

  Since the frame aspect ratio is set by the codec but the stream aspect ratio
  is set by the demuxer, these two may not be equal. This function tries to
  return the value that you should use if you would like to display the frame.

  Basic logic is to use the stream aspect ratio if it is set to something sane
  otherwise use the frame aspect ratio. This way a container setting, which is
  usually easy to modify can override the coded value in the frames.

  @param format the format context which the stream is part of
  @param stream the stream which the frame is part of
  @param frame the frame with the aspect ratio to be determined
  @return the guessed (valid) sample_aspect_ratio, 0/1 if no idea
*/
func AVGuessSampleAspectRatio(format *AVFormatContext, stream *AVStream, frame *AVFrame) *AVRational {
	var tmpformat *C.AVFormatContext
	if format != nil {
		tmpformat = format.ptr
	}
	var tmpstream *C.AVStream
	if stream != nil {
		tmpstream = stream.ptr
	}
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.av_guess_sample_aspect_ratio(tmpformat, tmpstream, tmpframe)
	return &AVRational{value: ret}
}

// --- Function av_guess_frame_rate ---

// AVGuessFrameRate wraps av_guess_frame_rate.
/*
  Guess the frame rate, based on both the container and codec information.

  @param ctx the format context which the stream is part of
  @param stream the stream which the frame is part of
  @param frame the frame for which the frame rate should be determined, may be NULL
  @return the guessed (valid) frame rate, 0/1 if no idea
*/
func AVGuessFrameRate(ctx *AVFormatContext, stream *AVStream, frame *AVFrame) *AVRational {
	var tmpctx *C.AVFormatContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	var tmpstream *C.AVStream
	if stream != nil {
		tmpstream = stream.ptr
	}
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.av_guess_frame_rate(tmpctx, tmpstream, tmpframe)
	return &AVRational{value: ret}
}

// --- Function avformat_match_stream_specifier ---

// AVFormatMatchStreamSpecifier wraps avformat_match_stream_specifier.
/*
  Check if the stream st contained in s is matched by the stream specifier
  spec.

  See the "stream specifiers" chapter in the documentation for the syntax
  of spec.

  @return  >0 if st is matched by spec;
           0  if st is not matched by spec;
           AVERROR code if spec is invalid

  @note  A stream specifier can match several streams in the format.
*/
func AVFormatMatchStreamSpecifier(s *AVFormatContext, st *AVStream, spec *CStr) (int, error) {
	var tmps *C.AVFormatContext
	if s != nil {
		tmps = s.ptr
	}
	var tmpst *C.AVStream
	if st != nil {
		tmpst = st.ptr
	}
	var tmpspec *C.char
	if spec != nil {
		tmpspec = spec.ptr
	}
	ret := C.avformat_match_stream_specifier(tmps, tmpst, tmpspec)
	return int(ret), WrapErr(int(ret))
}

// --- Function avformat_queue_attached_pictures ---

// AVFormatQueueAttachedPictures wraps avformat_queue_attached_pictures.
func AVFormatQueueAttachedPictures(s *AVFormatContext) (int, error) {
	var tmps *C.AVFormatContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.avformat_queue_attached_pictures(tmps)
	return int(ret), WrapErr(int(ret))
}

// --- Function avformat_transfer_internal_stream_timing_info ---

// AVFormatTransferInternalStreamTimingInfo wraps avformat_transfer_internal_stream_timing_info.
//
//	@deprecated do not call this function
func AVFormatTransferInternalStreamTimingInfo(ofmt *AVOutputFormat, ost *AVStream, ist *AVStream, copyTb AVTimebaseSource) (int, error) {
	var tmpofmt *C.AVOutputFormat
	if ofmt != nil {
		tmpofmt = ofmt.ptr
	}
	var tmpost *C.AVStream
	if ost != nil {
		tmpost = ost.ptr
	}
	var tmpist *C.AVStream
	if ist != nil {
		tmpist = ist.ptr
	}
	ret := C.avformat_transfer_internal_stream_timing_info(tmpofmt, tmpost, tmpist, C.enum_AVTimebaseSource(copyTb))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_stream_get_codec_timebase ---

// AVStreamGetCodecTimebase wraps av_stream_get_codec_timebase.
//
//	@deprecated do not call this function
func AVStreamGetCodecTimebase(st *AVStream) *AVRational {
	var tmpst *C.AVStream
	if st != nil {
		tmpst = st.ptr
	}
	ret := C.av_stream_get_codec_timebase(tmpst)
	return &AVRational{value: ret}
}

// --- Function avio_find_protocol_name ---

// AVIOFindProtocolName wraps avio_find_protocol_name.
/*
  Return the name of the protocol that will handle the passed URL.

  NULL is returned if no protocol could be found for the given URL.

  @return Name of the protocol or NULL.
*/
func AVIOFindProtocolName(url *CStr) *CStr {
	var tmpurl *C.char
	if url != nil {
		tmpurl = url.ptr
	}
	ret := C.avio_find_protocol_name(tmpurl)
	return wrapCStr(ret)
}

// --- Function avio_check ---

// AVIOCheck wraps avio_check.
/*
  Return AVIO_FLAG_* access flags corresponding to the access permissions
  of the resource in url, or a negative value corresponding to an
  AVERROR code in case of failure. The returned access flags are
  masked by the value in flags.

  @note This function is intrinsically unsafe, in the sense that the
  checked resource may change its existence or permission status from
  one call to another. Thus you should not trust the returned value,
  unless you are sure that no other processes are accessing the
  checked resource.
*/
func AVIOCheck(url *CStr, flags int) (int, error) {
	var tmpurl *C.char
	if url != nil {
		tmpurl = url.ptr
	}
	ret := C.avio_check(tmpurl, C.int(flags))
	return int(ret), WrapErr(int(ret))
}

// --- Function avio_open_dir ---

// AVIOOpenDir wraps avio_open_dir.
/*
  Open directory for reading.

  @param s       directory read context. Pointer to a NULL pointer must be passed.
  @param url     directory to be listed.
  @param options A dictionary filled with protocol-private options. On return
                 this parameter will be destroyed and replaced with a dictionary
                 containing options that were not found. May be NULL.
  @return >=0 on success or negative on error.
*/
func AVIOOpenDir(s **AVIODirContext, url *CStr, options **AVDictionary) (int, error) {
	var ptrs **C.AVIODirContext
	var tmps *C.AVIODirContext
	var oldTmps *C.AVIODirContext
	if s != nil {
		inners := *s
		if inners != nil {
			tmps = inners.ptr
			oldTmps = tmps
		}
		ptrs = &tmps
	}
	var tmpurl *C.char
	if url != nil {
		tmpurl = url.ptr
	}
	var ptroptions **C.AVDictionary
	var tmpoptions *C.AVDictionary
	var oldTmpoptions *C.AVDictionary
	if options != nil {
		inneroptions := *options
		if inneroptions != nil {
			tmpoptions = inneroptions.ptr
			oldTmpoptions = tmpoptions
		}
		ptroptions = &tmpoptions
	}
	ret := C.avio_open_dir(ptrs, tmpurl, ptroptions)
	if tmps != oldTmps && s != nil {
		if tmps != nil {
			*s = &AVIODirContext{ptr: tmps}
		} else {
			*s = nil
		}
	}
	if tmpoptions != oldTmpoptions && options != nil {
		if tmpoptions != nil {
			*options = &AVDictionary{ptr: tmpoptions}
		} else {
			*options = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function avio_read_dir ---

// AVIOReadDir wraps avio_read_dir.
/*
  Get next directory entry.

  Returned entry must be freed with avio_free_directory_entry(). In particular
  it may outlive AVIODirContext.

  @param s         directory read context.
  @param[out] next next entry or NULL when no more entries.
  @return >=0 on success or negative on error. End of list is not considered an
              error.
*/
func AVIOReadDir(s *AVIODirContext, next **AVIODirEntry) (int, error) {
	var tmps *C.AVIODirContext
	if s != nil {
		tmps = s.ptr
	}
	var ptrnext **C.AVIODirEntry
	var tmpnext *C.AVIODirEntry
	var oldTmpnext *C.AVIODirEntry
	if next != nil {
		innernext := *next
		if innernext != nil {
			tmpnext = innernext.ptr
			oldTmpnext = tmpnext
		}
		ptrnext = &tmpnext
	}
	ret := C.avio_read_dir(tmps, ptrnext)
	if tmpnext != oldTmpnext && next != nil {
		if tmpnext != nil {
			*next = &AVIODirEntry{ptr: tmpnext}
		} else {
			*next = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function avio_close_dir ---

// AVIOCloseDir wraps avio_close_dir.
/*
  Close directory.

  @note Entries created using avio_read_dir() are not deleted and must be
  freeded with avio_free_directory_entry().

  @param s         directory read context.
  @return >=0 on success or negative on error.
*/
func AVIOCloseDir(s **AVIODirContext) (int, error) {
	var ptrs **C.AVIODirContext
	var tmps *C.AVIODirContext
	var oldTmps *C.AVIODirContext
	if s != nil {
		inners := *s
		if inners != nil {
			tmps = inners.ptr
			oldTmps = tmps
		}
		ptrs = &tmps
	}
	ret := C.avio_close_dir(ptrs)
	if tmps != oldTmps && s != nil {
		if tmps != nil {
			*s = &AVIODirContext{ptr: tmps}
		} else {
			*s = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function avio_free_directory_entry ---

// AVIOFreeDirectoryEntry wraps avio_free_directory_entry.
/*
  Free entry allocated by avio_read_dir().

  @param entry entry to be freed.
*/
func AVIOFreeDirectoryEntry(entry **AVIODirEntry) {
	var ptrentry **C.AVIODirEntry
	var tmpentry *C.AVIODirEntry
	var oldTmpentry *C.AVIODirEntry
	if entry != nil {
		innerentry := *entry
		if innerentry != nil {
			tmpentry = innerentry.ptr
			oldTmpentry = tmpentry
		}
		ptrentry = &tmpentry
	}
	C.avio_free_directory_entry(ptrentry)
	if tmpentry != oldTmpentry && entry != nil {
		if tmpentry != nil {
			*entry = &AVIODirEntry{ptr: tmpentry}
		} else {
			*entry = nil
		}
	}
}

// --- Function avio_alloc_context ---

// avio_alloc_context skipped due to read_packet.

// --- Function avio_context_free ---

// AVIOContextFree wraps avio_context_free.
/*
  Free the supplied IO context and everything associated with it.

  @param s Double pointer to the IO context. This function will write NULL
  into s.
*/
func AVIOContextFree(s **AVIOContext) {
	var ptrs **C.AVIOContext
	var tmps *C.AVIOContext
	var oldTmps *C.AVIOContext
	if s != nil {
		inners := *s
		if inners != nil {
			tmps = inners.ptr
			oldTmps = tmps
		}
		ptrs = &tmps
	}
	C.avio_context_free(ptrs)
	if tmps != oldTmps && s != nil {
		if tmps != nil {
			*s = &AVIOContext{ptr: tmps}
		} else {
			*s = nil
		}
	}
}

// --- Function avio_w8 ---

// AVIOW8 wraps avio_w8.
func AVIOW8(s *AVIOContext, b int) {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	C.avio_w8(tmps, C.int(b))
}

// --- Function avio_write ---

// AVIOWrite wraps avio_write.
func AVIOWrite(s *AVIOContext, buf unsafe.Pointer, size int) {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	C.avio_write(tmps, (*C.uchar)(buf), C.int(size))
}

// --- Function avio_wl64 ---

// AVIOWl64 wraps avio_wl64.
func AVIOWl64(s *AVIOContext, val uint64) {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	C.avio_wl64(tmps, C.uint64_t(val))
}

// --- Function avio_wb64 ---

// AVIOWb64 wraps avio_wb64.
func AVIOWb64(s *AVIOContext, val uint64) {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	C.avio_wb64(tmps, C.uint64_t(val))
}

// --- Function avio_wl32 ---

// AVIOWl32 wraps avio_wl32.
func AVIOWl32(s *AVIOContext, val uint) {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	C.avio_wl32(tmps, C.uint(val))
}

// --- Function avio_wb32 ---

// AVIOWb32 wraps avio_wb32.
func AVIOWb32(s *AVIOContext, val uint) {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	C.avio_wb32(tmps, C.uint(val))
}

// --- Function avio_wl24 ---

// AVIOWl24 wraps avio_wl24.
func AVIOWl24(s *AVIOContext, val uint) {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	C.avio_wl24(tmps, C.uint(val))
}

// --- Function avio_wb24 ---

// AVIOWb24 wraps avio_wb24.
func AVIOWb24(s *AVIOContext, val uint) {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	C.avio_wb24(tmps, C.uint(val))
}

// --- Function avio_wl16 ---

// AVIOWl16 wraps avio_wl16.
func AVIOWl16(s *AVIOContext, val uint) {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	C.avio_wl16(tmps, C.uint(val))
}

// --- Function avio_wb16 ---

// AVIOWb16 wraps avio_wb16.
func AVIOWb16(s *AVIOContext, val uint) {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	C.avio_wb16(tmps, C.uint(val))
}

// --- Function avio_put_str ---

// AVIOPutStr wraps avio_put_str.
/*
  Write a NULL-terminated string.
  @return number of bytes written.
*/
func AVIOPutStr(s *AVIOContext, str *CStr) (int, error) {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	var tmpstr *C.char
	if str != nil {
		tmpstr = str.ptr
	}
	ret := C.avio_put_str(tmps, tmpstr)
	return int(ret), WrapErr(int(ret))
}

// --- Function avio_put_str16le ---

// AVIOPutStr16Le wraps avio_put_str16le.
/*
  Convert an UTF-8 string to UTF-16LE and write it.
  @param s the AVIOContext
  @param str NULL-terminated UTF-8 string

  @return number of bytes written.
*/
func AVIOPutStr16Le(s *AVIOContext, str *CStr) (int, error) {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	var tmpstr *C.char
	if str != nil {
		tmpstr = str.ptr
	}
	ret := C.avio_put_str16le(tmps, tmpstr)
	return int(ret), WrapErr(int(ret))
}

// --- Function avio_put_str16be ---

// AVIOPutStr16Be wraps avio_put_str16be.
/*
  Convert an UTF-8 string to UTF-16BE and write it.
  @param s the AVIOContext
  @param str NULL-terminated UTF-8 string

  @return number of bytes written.
*/
func AVIOPutStr16Be(s *AVIOContext, str *CStr) (int, error) {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	var tmpstr *C.char
	if str != nil {
		tmpstr = str.ptr
	}
	ret := C.avio_put_str16be(tmps, tmpstr)
	return int(ret), WrapErr(int(ret))
}

// --- Function avio_write_marker ---

// AVIOWriteMarker wraps avio_write_marker.
/*
  Mark the written bytestream as a specific type.

  Zero-length ranges are omitted from the output.

  @param s    the AVIOContext
  @param time the stream time the current bytestream pos corresponds to
              (in AV_TIME_BASE units), or AV_NOPTS_VALUE if unknown or not
              applicable
  @param type the kind of data written starting at the current pos
*/
func AVIOWriteMarker(s *AVIOContext, time int64, _type AVIODataMarkerType) {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	C.avio_write_marker(tmps, C.int64_t(time), C.enum_AVIODataMarkerType(_type))
}

// --- Function avio_seek ---

// AVIOSeek wraps avio_seek.
/*
  fseek() equivalent for AVIOContext.
  @return new position or AVERROR.
*/
func AVIOSeek(s *AVIOContext, offset int64, whence int) int64 {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.avio_seek(tmps, C.int64_t(offset), C.int(whence))
	return int64(ret)
}

// --- Function avio_skip ---

// AVIOSkip wraps avio_skip.
/*
  Skip given number of bytes forward
  @return new position or AVERROR.
*/
func AVIOSkip(s *AVIOContext, offset int64) int64 {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.avio_skip(tmps, C.int64_t(offset))
	return int64(ret)
}

// --- Function avio_tell ---

// AVIOTell wraps avio_tell.
/*
  ftell() equivalent for AVIOContext.
  @return position or AVERROR.
*/
func AVIOTell(s *AVIOContext) int64 {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.avio_tell(tmps)
	return int64(ret)
}

// --- Function avio_size ---

// AVIOSize wraps avio_size.
/*
  Get the filesize.
  @return filesize or AVERROR
*/
func AVIOSize(s *AVIOContext) int64 {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.avio_size(tmps)
	return int64(ret)
}

// --- Function avio_feof ---

// AVIOFeof wraps avio_feof.
/*
  Similar to feof() but also returns nonzero on read errors.
  @return non zero if and only if at end of file or a read error happened when reading.
*/
func AVIOFeof(s *AVIOContext) (int, error) {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.avio_feof(tmps)
	return int(ret), WrapErr(int(ret))
}

// --- Function avio_vprintf ---

// avio_vprintf skipped due to ap.

// --- Function avio_printf ---

// avio_printf skipped due to variadic arg.

// --- Function avio_print_string_array ---

// avio_print_string_array skipped due to strings

// --- Function avio_flush ---

// AVIOFlush wraps avio_flush.
/*
  Force flushing of buffered data.

  For write streams, force the buffered data to be immediately written to the output,
  without to wait to fill the internal buffer.

  For read streams, discard all currently buffered data, and advance the
  reported file position to that of the underlying stream. This does not
  read new data, and does not perform any seeks.
*/
func AVIOFlush(s *AVIOContext) {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	C.avio_flush(tmps)
}

// --- Function avio_read ---

// AVIORead wraps avio_read.
/*
  Read size bytes from AVIOContext into buf.
  @return number of bytes read or AVERROR
*/
func AVIORead(s *AVIOContext, buf unsafe.Pointer, size int) (int, error) {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.avio_read(tmps, (*C.uchar)(buf), C.int(size))
	return int(ret), WrapErr(int(ret))
}

// --- Function avio_read_partial ---

// AVIOReadPartial wraps avio_read_partial.
/*
  Read size bytes from AVIOContext into buf. Unlike avio_read(), this is allowed
  to read fewer bytes than requested. The missing bytes can be read in the next
  call. This always tries to read at least 1 byte.
  Useful to reduce latency in certain cases.
  @return number of bytes read or AVERROR
*/
func AVIOReadPartial(s *AVIOContext, buf unsafe.Pointer, size int) (int, error) {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.avio_read_partial(tmps, (*C.uchar)(buf), C.int(size))
	return int(ret), WrapErr(int(ret))
}

// --- Function avio_r8 ---

// AVIOR8 wraps avio_r8.
/*

  @note return 0 if EOF, so you cannot use it if EOF handling is
        necessary
*/
func AVIOR8(s *AVIOContext) (int, error) {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.avio_r8(tmps)
	return int(ret), WrapErr(int(ret))
}

// --- Function avio_rl16 ---

// AVIORl16 wraps avio_rl16.
func AVIORl16(s *AVIOContext) uint {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.avio_rl16(tmps)
	return uint(ret)
}

// --- Function avio_rl24 ---

// AVIORl24 wraps avio_rl24.
func AVIORl24(s *AVIOContext) uint {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.avio_rl24(tmps)
	return uint(ret)
}

// --- Function avio_rl32 ---

// AVIORl32 wraps avio_rl32.
func AVIORl32(s *AVIOContext) uint {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.avio_rl32(tmps)
	return uint(ret)
}

// --- Function avio_rl64 ---

// AVIORl64 wraps avio_rl64.
func AVIORl64(s *AVIOContext) uint64 {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.avio_rl64(tmps)
	return uint64(ret)
}

// --- Function avio_rb16 ---

// AVIORb16 wraps avio_rb16.
func AVIORb16(s *AVIOContext) uint {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.avio_rb16(tmps)
	return uint(ret)
}

// --- Function avio_rb24 ---

// AVIORb24 wraps avio_rb24.
func AVIORb24(s *AVIOContext) uint {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.avio_rb24(tmps)
	return uint(ret)
}

// --- Function avio_rb32 ---

// AVIORb32 wraps avio_rb32.
func AVIORb32(s *AVIOContext) uint {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.avio_rb32(tmps)
	return uint(ret)
}

// --- Function avio_rb64 ---

// AVIORb64 wraps avio_rb64.
func AVIORb64(s *AVIOContext) uint64 {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.avio_rb64(tmps)
	return uint64(ret)
}

// --- Function avio_get_str ---

// AVIOGetStr wraps avio_get_str.
/*
  Read a string from pb into buf. The reading will terminate when either
  a NULL character was encountered, maxlen bytes have been read, or nothing
  more can be read from pb. The result is guaranteed to be NULL-terminated, it
  will be truncated if buf is too small.
  Note that the string is not interpreted or validated in any way, it
  might get truncated in the middle of a sequence for multi-byte encodings.

  @return number of bytes read (is always <= maxlen).
  If reading ends on EOF or error, the return value will be one more than
  bytes actually read.
*/
func AVIOGetStr(pb *AVIOContext, maxlen int, buf *CStr, buflen int) (int, error) {
	var tmppb *C.AVIOContext
	if pb != nil {
		tmppb = pb.ptr
	}
	var tmpbuf *C.char
	if buf != nil {
		tmpbuf = buf.ptr
	}
	ret := C.avio_get_str(tmppb, C.int(maxlen), tmpbuf, C.int(buflen))
	return int(ret), WrapErr(int(ret))
}

// --- Function avio_get_str16le ---

// AVIOGetStr16Le wraps avio_get_str16le.
/*
  Read a UTF-16 string from pb and convert it to UTF-8.
  The reading will terminate when either a null or invalid character was
  encountered or maxlen bytes have been read.
  @return number of bytes read (is always <= maxlen)
*/
func AVIOGetStr16Le(pb *AVIOContext, maxlen int, buf *CStr, buflen int) (int, error) {
	var tmppb *C.AVIOContext
	if pb != nil {
		tmppb = pb.ptr
	}
	var tmpbuf *C.char
	if buf != nil {
		tmpbuf = buf.ptr
	}
	ret := C.avio_get_str16le(tmppb, C.int(maxlen), tmpbuf, C.int(buflen))
	return int(ret), WrapErr(int(ret))
}

// --- Function avio_get_str16be ---

// AVIOGetStr16Be wraps avio_get_str16be.
func AVIOGetStr16Be(pb *AVIOContext, maxlen int, buf *CStr, buflen int) (int, error) {
	var tmppb *C.AVIOContext
	if pb != nil {
		tmppb = pb.ptr
	}
	var tmpbuf *C.char
	if buf != nil {
		tmpbuf = buf.ptr
	}
	ret := C.avio_get_str16be(tmppb, C.int(maxlen), tmpbuf, C.int(buflen))
	return int(ret), WrapErr(int(ret))
}

// --- Function avio_open ---

// AVIOOpen wraps avio_open.
/*
  Create and initialize a AVIOContext for accessing the
  resource indicated by url.
  @note When the resource indicated by url has been opened in
  read+write mode, the AVIOContext can be used only for writing.

  @param s Used to return the pointer to the created AVIOContext.
  In case of failure the pointed to value is set to NULL.
  @param url resource to access
  @param flags flags which control how the resource indicated by url
  is to be opened
  @return >= 0 in case of success, a negative value corresponding to an
  AVERROR code in case of failure
*/
func AVIOOpen(s **AVIOContext, url *CStr, flags int) (int, error) {
	var ptrs **C.AVIOContext
	var tmps *C.AVIOContext
	var oldTmps *C.AVIOContext
	if s != nil {
		inners := *s
		if inners != nil {
			tmps = inners.ptr
			oldTmps = tmps
		}
		ptrs = &tmps
	}
	var tmpurl *C.char
	if url != nil {
		tmpurl = url.ptr
	}
	ret := C.avio_open(ptrs, tmpurl, C.int(flags))
	if tmps != oldTmps && s != nil {
		if tmps != nil {
			*s = &AVIOContext{ptr: tmps}
		} else {
			*s = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function avio_open2 ---

// AVIOOpen2 wraps avio_open2.
/*
  Create and initialize a AVIOContext for accessing the
  resource indicated by url.
  @note When the resource indicated by url has been opened in
  read+write mode, the AVIOContext can be used only for writing.

  @param s Used to return the pointer to the created AVIOContext.
  In case of failure the pointed to value is set to NULL.
  @param url resource to access
  @param flags flags which control how the resource indicated by url
  is to be opened
  @param int_cb an interrupt callback to be used at the protocols level
  @param options  A dictionary filled with protocol-private options. On return
  this parameter will be destroyed and replaced with a dict containing options
  that were not found. May be NULL.
  @return >= 0 in case of success, a negative value corresponding to an
  AVERROR code in case of failure
*/
func AVIOOpen2(s **AVIOContext, url *CStr, flags int, intCb *AVIOInterruptCB, options **AVDictionary) (int, error) {
	var ptrs **C.AVIOContext
	var tmps *C.AVIOContext
	var oldTmps *C.AVIOContext
	if s != nil {
		inners := *s
		if inners != nil {
			tmps = inners.ptr
			oldTmps = tmps
		}
		ptrs = &tmps
	}
	var tmpurl *C.char
	if url != nil {
		tmpurl = url.ptr
	}
	var tmpintCb *C.AVIOInterruptCB
	if intCb != nil {
		tmpintCb = intCb.ptr
	}
	var ptroptions **C.AVDictionary
	var tmpoptions *C.AVDictionary
	var oldTmpoptions *C.AVDictionary
	if options != nil {
		inneroptions := *options
		if inneroptions != nil {
			tmpoptions = inneroptions.ptr
			oldTmpoptions = tmpoptions
		}
		ptroptions = &tmpoptions
	}
	ret := C.avio_open2(ptrs, tmpurl, C.int(flags), tmpintCb, ptroptions)
	if tmps != oldTmps && s != nil {
		if tmps != nil {
			*s = &AVIOContext{ptr: tmps}
		} else {
			*s = nil
		}
	}
	if tmpoptions != oldTmpoptions && options != nil {
		if tmpoptions != nil {
			*options = &AVDictionary{ptr: tmpoptions}
		} else {
			*options = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function avio_close ---

// AVIOClose wraps avio_close.
/*
  Close the resource accessed by the AVIOContext s and free it.
  This function can only be used if s was opened by avio_open().

  The internal buffer is automatically flushed before closing the
  resource.

  @return 0 on success, an AVERROR < 0 on error.
  @see avio_closep
*/
func AVIOClose(s *AVIOContext) (int, error) {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.avio_close(tmps)
	return int(ret), WrapErr(int(ret))
}

// --- Function avio_closep ---

// AVIOClosep wraps avio_closep.
/*
  Close the resource accessed by the AVIOContext *s, free it
  and set the pointer pointing to it to NULL.
  This function can only be used if s was opened by avio_open().

  The internal buffer is automatically flushed before closing the
  resource.

  @return 0 on success, an AVERROR < 0 on error.
  @see avio_close
*/
func AVIOClosep(s **AVIOContext) (int, error) {
	var ptrs **C.AVIOContext
	var tmps *C.AVIOContext
	var oldTmps *C.AVIOContext
	if s != nil {
		inners := *s
		if inners != nil {
			tmps = inners.ptr
			oldTmps = tmps
		}
		ptrs = &tmps
	}
	ret := C.avio_closep(ptrs)
	if tmps != oldTmps && s != nil {
		if tmps != nil {
			*s = &AVIOContext{ptr: tmps}
		} else {
			*s = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function avio_open_dyn_buf ---

// AVIOOpenDynBuf wraps avio_open_dyn_buf.
/*
  Open a write only memory stream.

  @param s new IO context
  @return zero if no error.
*/
func AVIOOpenDynBuf(s **AVIOContext) (int, error) {
	var ptrs **C.AVIOContext
	var tmps *C.AVIOContext
	var oldTmps *C.AVIOContext
	if s != nil {
		inners := *s
		if inners != nil {
			tmps = inners.ptr
			oldTmps = tmps
		}
		ptrs = &tmps
	}
	ret := C.avio_open_dyn_buf(ptrs)
	if tmps != oldTmps && s != nil {
		if tmps != nil {
			*s = &AVIOContext{ptr: tmps}
		} else {
			*s = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function avio_get_dyn_buf ---

// avio_get_dyn_buf skipped due to pbuffer

// --- Function avio_close_dyn_buf ---

// avio_close_dyn_buf skipped due to pbuffer

// --- Function avio_enum_protocols ---

// avio_enum_protocols skipped due to opaque

// --- Function avio_protocol_get_class ---

// AVIOProtocolGetClass wraps avio_protocol_get_class.
/*
  Get AVClass by names of available protocols.

  @return A AVClass of input protocol name or NULL
*/
func AVIOProtocolGetClass(name *CStr) *AVClass {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.avio_protocol_get_class(tmpname)
	var retMapped *AVClass
	if ret != nil {
		retMapped = &AVClass{ptr: ret}
	}
	return retMapped
}

// --- Function avio_pause ---

// AVIOPause wraps avio_pause.
/*
  Pause and resume playing - only meaningful if using a network streaming
  protocol (e.g. MMS).

  @param h     IO context from which to call the read_pause function pointer
  @param pause 1 for pause, 0 for resume
*/
func AVIOPause(h *AVIOContext, pause int) (int, error) {
	var tmph *C.AVIOContext
	if h != nil {
		tmph = h.ptr
	}
	ret := C.avio_pause(tmph, C.int(pause))
	return int(ret), WrapErr(int(ret))
}

// --- Function avio_seek_time ---

// AVIOSeekTime wraps avio_seek_time.
/*
  Seek to a given timestamp relative to some component stream.
  Only meaningful if using a network streaming protocol (e.g. MMS.).

  @param h IO context from which to call the seek function pointers
  @param stream_index The stream index that the timestamp is relative to.
         If stream_index is (-1) the timestamp should be in AV_TIME_BASE
         units from the beginning of the presentation.
         If a stream_index >= 0 is used and the protocol does not support
         seeking based on component streams, the call will fail.
  @param timestamp timestamp in AVStream.time_base units
         or if there is no stream specified then in AV_TIME_BASE units.
  @param flags Optional combination of AVSEEK_FLAG_BACKWARD, AVSEEK_FLAG_BYTE
         and AVSEEK_FLAG_ANY. The protocol may silently ignore
         AVSEEK_FLAG_BACKWARD and AVSEEK_FLAG_ANY, but AVSEEK_FLAG_BYTE will
         fail if used and not supported.
  @return >= 0 on success
  @see AVInputFormat::read_seek
*/
func AVIOSeekTime(h *AVIOContext, streamIndex int, timestamp int64, flags int) int64 {
	var tmph *C.AVIOContext
	if h != nil {
		tmph = h.ptr
	}
	ret := C.avio_seek_time(tmph, C.int(streamIndex), C.int64_t(timestamp), C.int(flags))
	return int64(ret)
}

// --- Function avio_read_to_bprint ---

// AVIOReadToBprint wraps avio_read_to_bprint.
/*
  Read contents of h into print buffer, up to max_size bytes, or up to EOF.

  @return 0 for success (max_size bytes read or EOF reached), negative error
  code otherwise
*/
func AVIOReadToBprint(h *AVIOContext, pb *AVBPrint, maxSize uint64) (int, error) {
	var tmph *C.AVIOContext
	if h != nil {
		tmph = h.ptr
	}
	var tmppb *C.AVBPrint
	if pb != nil {
		tmppb = pb.ptr
	}
	ret := C.avio_read_to_bprint(tmph, tmppb, C.size_t(maxSize))
	return int(ret), WrapErr(int(ret))
}

// --- Function avio_accept ---

// AVIOAccept wraps avio_accept.
/*
  Accept and allocate a client context on a server context.
  @param  s the server context
  @param  c the client context, must be unallocated
  @return   >= 0 on success or a negative value corresponding
            to an AVERROR on failure
*/
func AVIOAccept(s *AVIOContext, c **AVIOContext) (int, error) {
	var tmps *C.AVIOContext
	if s != nil {
		tmps = s.ptr
	}
	var ptrc **C.AVIOContext
	var tmpc *C.AVIOContext
	var oldTmpc *C.AVIOContext
	if c != nil {
		innerc := *c
		if innerc != nil {
			tmpc = innerc.ptr
			oldTmpc = tmpc
		}
		ptrc = &tmpc
	}
	ret := C.avio_accept(tmps, ptrc)
	if tmpc != oldTmpc && c != nil {
		if tmpc != nil {
			*c = &AVIOContext{ptr: tmpc}
		} else {
			*c = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function avio_handshake ---

// AVIOHandshake wraps avio_handshake.
/*
  Perform one step of the protocol handshake to accept a new client.
  This function must be called on a client returned by avio_accept() before
  using it as a read/write context.
  It is separate from avio_accept() because it may block.
  A step of the handshake is defined by places where the application may
  decide to change the proceedings.
  For example, on a protocol with a request header and a reply header, each
  one can constitute a step because the application may use the parameters
  from the request to change parameters in the reply; or each individual
  chunk of the request can constitute a step.
  If the handshake is already finished, avio_handshake() does nothing and
  returns 0 immediately.

  @param  c the client context to perform the handshake on
  @return   0   on a complete and successful handshake
            > 0 if the handshake progressed, but is not complete
            < 0 for an AVERROR code
*/
func AVIOHandshake(c *AVIOContext) (int, error) {
	var tmpc *C.AVIOContext
	if c != nil {
		tmpc = c.ptr
	}
	ret := C.avio_handshake(tmpc)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_adler32_update ---

// AVAdler32Update wraps av_adler32_update.
/*
  Calculate the Adler32 checksum of a buffer.

  Passing the return value to a subsequent av_adler32_update() call
  allows the checksum of multiple buffers to be calculated as though
  they were concatenated.

  @param adler initial checksum value
  @param buf   pointer to input buffer
  @param len   size of input buffer
  @return      updated checksum
*/
func AVAdler32Update(adler AVAdler, buf unsafe.Pointer, len uint64) AVAdler {
	ret := C.av_adler32_update(C.AVAdler(adler), (*C.uint8_t)(buf), C.size_t(len))
	return AVAdler(ret)
}

// --- Function av_aes_alloc ---

// AVAesAlloc wraps av_aes_alloc.
//
//	Allocate an AVAES context.
func AVAesAlloc() *AVAES {
	ret := C.av_aes_alloc()
	var retMapped *AVAES
	if ret != nil {
		retMapped = &AVAES{ptr: ret}
	}
	return retMapped
}

// --- Function av_aes_init ---

// AVAesInit wraps av_aes_init.
/*
  Initialize an AVAES context.

  @param a The AVAES context
  @param key Pointer to the key
  @param key_bits 128, 192 or 256
  @param decrypt 0 for encryption, 1 for decryption
*/
func AVAesInit(a *AVAES, key unsafe.Pointer, keyBits int, decrypt int) (int, error) {
	var tmpa *C.struct_AVAES
	if a != nil {
		tmpa = a.ptr
	}
	ret := C.av_aes_init(tmpa, (*C.uint8_t)(key), C.int(keyBits), C.int(decrypt))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_aes_crypt ---

// AVAesCrypt wraps av_aes_crypt.
/*
  Encrypt or decrypt a buffer using a previously initialized context.

  @param a The AVAES context
  @param dst destination array, can be equal to src
  @param src source array, can be equal to dst
  @param count number of 16 byte blocks
  @param iv initialization vector for CBC mode, if NULL then ECB will be used
  @param decrypt 0 for encryption, 1 for decryption
*/
func AVAesCrypt(a *AVAES, dst unsafe.Pointer, src unsafe.Pointer, count int, iv unsafe.Pointer, decrypt int) {
	var tmpa *C.struct_AVAES
	if a != nil {
		tmpa = a.ptr
	}
	C.av_aes_crypt(tmpa, (*C.uint8_t)(dst), (*C.uint8_t)(src), C.int(count), (*C.uint8_t)(iv), C.int(decrypt))
}

// --- Function av_aes_ctr_alloc ---

// AVAesCtrAlloc wraps av_aes_ctr_alloc.
//
//	Allocate an AVAESCTR context.
func AVAesCtrAlloc() *AVAESCTR {
	ret := C.av_aes_ctr_alloc()
	var retMapped *AVAESCTR
	if ret != nil {
		retMapped = &AVAESCTR{ptr: ret}
	}
	return retMapped
}

// --- Function av_aes_ctr_init ---

// AVAesCtrInit wraps av_aes_ctr_init.
/*
  Initialize an AVAESCTR context.

  @param a The AVAESCTR context to initialize
  @param key encryption key, must have a length of AES_CTR_KEY_SIZE
*/
func AVAesCtrInit(a *AVAESCTR, key unsafe.Pointer) (int, error) {
	var tmpa *C.struct_AVAESCTR
	if a != nil {
		tmpa = a.ptr
	}
	ret := C.av_aes_ctr_init(tmpa, (*C.uint8_t)(key))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_aes_ctr_free ---

// AVAesCtrFree wraps av_aes_ctr_free.
/*
  Release an AVAESCTR context.

  @param a The AVAESCTR context
*/
func AVAesCtrFree(a *AVAESCTR) {
	var tmpa *C.struct_AVAESCTR
	if a != nil {
		tmpa = a.ptr
	}
	C.av_aes_ctr_free(tmpa)
}

// --- Function av_aes_ctr_crypt ---

// AVAesCtrCrypt wraps av_aes_ctr_crypt.
/*
  Process a buffer using a previously initialized context.

  @param a The AVAESCTR context
  @param dst destination array, can be equal to src
  @param src source array, can be equal to dst
  @param size the size of src and dst
*/
func AVAesCtrCrypt(a *AVAESCTR, dst unsafe.Pointer, src unsafe.Pointer, size int) {
	var tmpa *C.struct_AVAESCTR
	if a != nil {
		tmpa = a.ptr
	}
	C.av_aes_ctr_crypt(tmpa, (*C.uint8_t)(dst), (*C.uint8_t)(src), C.int(size))
}

// --- Function av_aes_ctr_get_iv ---

// AVAesCtrGetIv wraps av_aes_ctr_get_iv.
//
//	Get the current iv
func AVAesCtrGetIv(a *AVAESCTR) unsafe.Pointer {
	var tmpa *C.struct_AVAESCTR
	if a != nil {
		tmpa = a.ptr
	}
	ret := C.av_aes_ctr_get_iv(tmpa)
	return unsafe.Pointer(ret)
}

// --- Function av_aes_ctr_set_random_iv ---

// AVAesCtrSetRandomIv wraps av_aes_ctr_set_random_iv.
//
//	Generate a random iv
func AVAesCtrSetRandomIv(a *AVAESCTR) {
	var tmpa *C.struct_AVAESCTR
	if a != nil {
		tmpa = a.ptr
	}
	C.av_aes_ctr_set_random_iv(tmpa)
}

// --- Function av_aes_ctr_set_iv ---

// AVAesCtrSetIv wraps av_aes_ctr_set_iv.
//
//	Forcefully change the 8-byte iv
func AVAesCtrSetIv(a *AVAESCTR, iv unsafe.Pointer) {
	var tmpa *C.struct_AVAESCTR
	if a != nil {
		tmpa = a.ptr
	}
	C.av_aes_ctr_set_iv(tmpa, (*C.uint8_t)(iv))
}

// --- Function av_aes_ctr_set_full_iv ---

// AVAesCtrSetFullIv wraps av_aes_ctr_set_full_iv.
//
//	Forcefully change the "full" 16-byte iv, including the counter
func AVAesCtrSetFullIv(a *AVAESCTR, iv unsafe.Pointer) {
	var tmpa *C.struct_AVAESCTR
	if a != nil {
		tmpa = a.ptr
	}
	C.av_aes_ctr_set_full_iv(tmpa, (*C.uint8_t)(iv))
}

// --- Function av_aes_ctr_increment_iv ---

// AVAesCtrIncrementIv wraps av_aes_ctr_increment_iv.
//
//	Increment the top 64 bit of the iv (performed after each frame)
func AVAesCtrIncrementIv(a *AVAESCTR) {
	var tmpa *C.struct_AVAESCTR
	if a != nil {
		tmpa = a.ptr
	}
	C.av_aes_ctr_increment_iv(tmpa)
}

// --- Function av_ambient_viewing_environment_alloc ---

// AVAmbientViewingEnvironmentAlloc wraps av_ambient_viewing_environment_alloc.
/*
  Allocate an AVAmbientViewingEnvironment structure.

  @return the newly allocated struct or NULL on failure
*/
func AVAmbientViewingEnvironmentAlloc(size *uint64) *AVAmbientViewingEnvironment {
	ret := C.av_ambient_viewing_environment_alloc((*C.size_t)(unsafe.Pointer(size)))
	var retMapped *AVAmbientViewingEnvironment
	if ret != nil {
		retMapped = &AVAmbientViewingEnvironment{ptr: ret}
	}
	return retMapped
}

// --- Function av_ambient_viewing_environment_create_side_data ---

// AVAmbientViewingEnvironmentCreateSideData wraps av_ambient_viewing_environment_create_side_data.
/*
  Allocate and add an AVAmbientViewingEnvironment structure to an existing
  AVFrame as side data.

  @return the newly allocated struct, or NULL on failure
*/
func AVAmbientViewingEnvironmentCreateSideData(frame *AVFrame) *AVAmbientViewingEnvironment {
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.av_ambient_viewing_environment_create_side_data(tmpframe)
	var retMapped *AVAmbientViewingEnvironment
	if ret != nil {
		retMapped = &AVAmbientViewingEnvironment{ptr: ret}
	}
	return retMapped
}

// --- Function av_audio_fifo_free ---

// AVAudioFifoFree wraps av_audio_fifo_free.
/*
  Free an AVAudioFifo.

  @param af  AVAudioFifo to free
*/
func AVAudioFifoFree(af *AVAudioFifo) {
	var tmpaf *C.AVAudioFifo
	if af != nil {
		tmpaf = af.ptr
	}
	C.av_audio_fifo_free(tmpaf)
}

// --- Function av_audio_fifo_alloc ---

// AVAudioFifoAlloc wraps av_audio_fifo_alloc.
/*
  Allocate an AVAudioFifo.

  @param sample_fmt  sample format
  @param channels    number of channels
  @param nb_samples  initial allocation size, in samples
  @return            newly allocated AVAudioFifo, or NULL on error
*/
func AVAudioFifoAlloc(sampleFmt AVSampleFormat, channels int, nbSamples int) *AVAudioFifo {
	ret := C.av_audio_fifo_alloc(C.enum_AVSampleFormat(sampleFmt), C.int(channels), C.int(nbSamples))
	var retMapped *AVAudioFifo
	if ret != nil {
		retMapped = &AVAudioFifo{ptr: ret}
	}
	return retMapped
}

// --- Function av_audio_fifo_realloc ---

// AVAudioFifoRealloc wraps av_audio_fifo_realloc.
/*
  Reallocate an AVAudioFifo.

  @param af          AVAudioFifo to reallocate
  @param nb_samples  new allocation size, in samples
  @return            0 if OK, or negative AVERROR code on failure
*/
func AVAudioFifoRealloc(af *AVAudioFifo, nbSamples int) (int, error) {
	var tmpaf *C.AVAudioFifo
	if af != nil {
		tmpaf = af.ptr
	}
	ret := C.av_audio_fifo_realloc(tmpaf, C.int(nbSamples))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_audio_fifo_write ---

// av_audio_fifo_write skipped due to data

// --- Function av_audio_fifo_peek ---

// av_audio_fifo_peek skipped due to data

// --- Function av_audio_fifo_peek_at ---

// av_audio_fifo_peek_at skipped due to data

// --- Function av_audio_fifo_read ---

// av_audio_fifo_read skipped due to data

// --- Function av_audio_fifo_drain ---

// AVAudioFifoDrain wraps av_audio_fifo_drain.
/*
  Drain data from an AVAudioFifo.

  Removes the data without reading it.

  @param af          AVAudioFifo to drain
  @param nb_samples  number of samples to drain
  @return            0 if OK, or negative AVERROR code on failure
*/
func AVAudioFifoDrain(af *AVAudioFifo, nbSamples int) (int, error) {
	var tmpaf *C.AVAudioFifo
	if af != nil {
		tmpaf = af.ptr
	}
	ret := C.av_audio_fifo_drain(tmpaf, C.int(nbSamples))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_audio_fifo_reset ---

// AVAudioFifoReset wraps av_audio_fifo_reset.
/*
  Reset the AVAudioFifo buffer.

  This empties all data in the buffer.

  @param af  AVAudioFifo to reset
*/
func AVAudioFifoReset(af *AVAudioFifo) {
	var tmpaf *C.AVAudioFifo
	if af != nil {
		tmpaf = af.ptr
	}
	C.av_audio_fifo_reset(tmpaf)
}

// --- Function av_audio_fifo_size ---

// AVAudioFifoSize wraps av_audio_fifo_size.
/*
  Get the current number of samples in the AVAudioFifo available for reading.

  @param af  the AVAudioFifo to query
  @return    number of samples available for reading
*/
func AVAudioFifoSize(af *AVAudioFifo) (int, error) {
	var tmpaf *C.AVAudioFifo
	if af != nil {
		tmpaf = af.ptr
	}
	ret := C.av_audio_fifo_size(tmpaf)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_audio_fifo_space ---

// AVAudioFifoSpace wraps av_audio_fifo_space.
/*
  Get the current number of samples in the AVAudioFifo available for writing.

  @param af  the AVAudioFifo to query
  @return    number of samples available for writing
*/
func AVAudioFifoSpace(af *AVAudioFifo) (int, error) {
	var tmpaf *C.AVAudioFifo
	if af != nil {
		tmpaf = af.ptr
	}
	ret := C.av_audio_fifo_space(tmpaf)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_assert0_fpu ---

// AVAssert0Fpu wraps av_assert0_fpu.
/*
  Assert that floating point operations can be executed.

  This will av_assert0() that the cpu is not in MMX state on X86
*/
func AVAssert0Fpu() {
	C.av_assert0_fpu()
}

// --- Function av_strstart ---

// av_strstart skipped due to ptr

// --- Function av_stristart ---

// av_stristart skipped due to ptr

// --- Function av_stristr ---

// AVStristr wraps av_stristr.
/*
  Locate the first case-independent occurrence in the string haystack
  of the string needle.  A zero-length string needle is considered to
  match at the start of haystack.

  This function is a case-insensitive version of the standard strstr().

  @param haystack string to search in
  @param needle   string to search for
  @return         pointer to the located match within haystack
                  or a null pointer if no match
*/
func AVStristr(haystack *CStr, needle *CStr) *CStr {
	var tmphaystack *C.char
	if haystack != nil {
		tmphaystack = haystack.ptr
	}
	var tmpneedle *C.char
	if needle != nil {
		tmpneedle = needle.ptr
	}
	ret := C.av_stristr(tmphaystack, tmpneedle)
	return wrapCStr(ret)
}

// --- Function av_strnstr ---

// AVStrnstr wraps av_strnstr.
/*
  Locate the first occurrence of the string needle in the string haystack
  where not more than hay_length characters are searched. A zero-length
  string needle is considered to match at the start of haystack.

  This function is a length-limited version of the standard strstr().

  @param haystack   string to search in
  @param needle     string to search for
  @param hay_length length of string to search in
  @return           pointer to the located match within haystack
                    or a null pointer if no match
*/
func AVStrnstr(haystack *CStr, needle *CStr, hayLength uint64) *CStr {
	var tmphaystack *C.char
	if haystack != nil {
		tmphaystack = haystack.ptr
	}
	var tmpneedle *C.char
	if needle != nil {
		tmpneedle = needle.ptr
	}
	ret := C.av_strnstr(tmphaystack, tmpneedle, C.size_t(hayLength))
	return wrapCStr(ret)
}

// --- Function av_strlcpy ---

// AVStrlcpy wraps av_strlcpy.
/*
  Copy the string src to dst, but no more than size - 1 bytes, and
  null-terminate dst.

  This function is the same as BSD strlcpy().

  @param dst destination buffer
  @param src source string
  @param size size of destination buffer
  @return the length of src

  @warning since the return value is the length of src, src absolutely
  _must_ be a properly 0-terminated string, otherwise this will read beyond
  the end of the buffer and possibly crash.
*/
func AVStrlcpy(dst *CStr, src *CStr, size uint64) uint64 {
	var tmpdst *C.char
	if dst != nil {
		tmpdst = dst.ptr
	}
	var tmpsrc *C.char
	if src != nil {
		tmpsrc = src.ptr
	}
	ret := C.av_strlcpy(tmpdst, tmpsrc, C.size_t(size))
	return uint64(ret)
}

// --- Function av_strlcat ---

// AVStrlcat wraps av_strlcat.
/*
  Append the string src to the string dst, but to a total length of
  no more than size - 1 bytes, and null-terminate dst.

  This function is similar to BSD strlcat(), but differs when
  size <= strlen(dst).

  @param dst destination buffer
  @param src source string
  @param size size of destination buffer
  @return the total length of src and dst

  @warning since the return value use the length of src and dst, these
  absolutely _must_ be a properly 0-terminated strings, otherwise this
  will read beyond the end of the buffer and possibly crash.
*/
func AVStrlcat(dst *CStr, src *CStr, size uint64) uint64 {
	var tmpdst *C.char
	if dst != nil {
		tmpdst = dst.ptr
	}
	var tmpsrc *C.char
	if src != nil {
		tmpsrc = src.ptr
	}
	ret := C.av_strlcat(tmpdst, tmpsrc, C.size_t(size))
	return uint64(ret)
}

// --- Function av_strlcatf ---

// av_strlcatf skipped due to variadic arg.

// --- Function av_strnlen ---

// AVStrnlen wraps av_strnlen.
/*
  Get the count of continuous non zero chars starting from the beginning.

  @param s   the string whose length to count
  @param len maximum number of characters to check in the string, that
             is the maximum value which is returned by the function
*/
func AVStrnlen(s *CStr, len uint64) uint64 {
	var tmps *C.char
	if s != nil {
		tmps = s.ptr
	}
	ret := C.av_strnlen(tmps, C.size_t(len))
	return uint64(ret)
}

// --- Function av_asprintf ---

// av_asprintf skipped due to variadic arg.

// --- Function av_get_token ---

// av_get_token skipped due to buf

// --- Function av_strtok ---

// av_strtok skipped due to saveptr

// --- Function av_isdigit ---

// AVIsdigit wraps av_isdigit.
//
//	Locale-independent conversion of ASCII isdigit.
func AVIsdigit(c int) (int, error) {
	ret := C.av_isdigit(C.int(c))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_isgraph ---

// AVIsgraph wraps av_isgraph.
//
//	Locale-independent conversion of ASCII isgraph.
func AVIsgraph(c int) (int, error) {
	ret := C.av_isgraph(C.int(c))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_isspace ---

// AVIsspace wraps av_isspace.
//
//	Locale-independent conversion of ASCII isspace.
func AVIsspace(c int) (int, error) {
	ret := C.av_isspace(C.int(c))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_toupper ---

// AVToupper wraps av_toupper.
//
//	Locale-independent conversion of ASCII characters to uppercase.
func AVToupper(c int) (int, error) {
	ret := C.av_toupper(C.int(c))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_tolower ---

// AVTolower wraps av_tolower.
//
//	Locale-independent conversion of ASCII characters to lowercase.
func AVTolower(c int) (int, error) {
	ret := C.av_tolower(C.int(c))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_isxdigit ---

// AVIsxdigit wraps av_isxdigit.
//
//	Locale-independent conversion of ASCII isxdigit.
func AVIsxdigit(c int) (int, error) {
	ret := C.av_isxdigit(C.int(c))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_strcasecmp ---

// AVStrcasecmp wraps av_strcasecmp.
/*
  Locale-independent case-insensitive compare.
  @note This means only ASCII-range characters are case-insensitive
*/
func AVStrcasecmp(a *CStr, b *CStr) (int, error) {
	var tmpa *C.char
	if a != nil {
		tmpa = a.ptr
	}
	var tmpb *C.char
	if b != nil {
		tmpb = b.ptr
	}
	ret := C.av_strcasecmp(tmpa, tmpb)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_strncasecmp ---

// AVStrncasecmp wraps av_strncasecmp.
/*
  Locale-independent case-insensitive compare.
  @note This means only ASCII-range characters are case-insensitive
*/
func AVStrncasecmp(a *CStr, b *CStr, n uint64) (int, error) {
	var tmpa *C.char
	if a != nil {
		tmpa = a.ptr
	}
	var tmpb *C.char
	if b != nil {
		tmpb = b.ptr
	}
	ret := C.av_strncasecmp(tmpa, tmpb, C.size_t(n))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_strireplace ---

// AVStrireplace wraps av_strireplace.
/*
  Locale-independent strings replace.
  @note This means only ASCII-range characters are replaced.
*/
func AVStrireplace(str *CStr, from *CStr, to *CStr) *CStr {
	var tmpstr *C.char
	if str != nil {
		tmpstr = str.ptr
	}
	var tmpfrom *C.char
	if from != nil {
		tmpfrom = from.ptr
	}
	var tmpto *C.char
	if to != nil {
		tmpto = to.ptr
	}
	ret := C.av_strireplace(tmpstr, tmpfrom, tmpto)
	return wrapCStr(ret)
}

// --- Function av_basename ---

// AVBasename wraps av_basename.
/*
  Thread safe basename.
  @param path the string to parse, on DOS both \ and / are considered separators.
  @return pointer to the basename substring.
  If path does not contain a slash, the function returns a copy of path.
  If path is a NULL pointer or points to an empty string, a pointer
  to a string "." is returned.
*/
func AVBasename(path *CStr) *CStr {
	var tmppath *C.char
	if path != nil {
		tmppath = path.ptr
	}
	ret := C.av_basename(tmppath)
	return wrapCStr(ret)
}

// --- Function av_dirname ---

// AVDirname wraps av_dirname.
/*
  Thread safe dirname.
  @param path the string to parse, on DOS both \ and / are considered separators.
  @return A pointer to a string that's the parent directory of path.
  If path is a NULL pointer or points to an empty string, a pointer
  to a string "." is returned.
  @note the function may modify the contents of the path, so copies should be passed.
*/
func AVDirname(path *CStr) *CStr {
	var tmppath *C.char
	if path != nil {
		tmppath = path.ptr
	}
	ret := C.av_dirname(tmppath)
	return wrapCStr(ret)
}

// --- Function av_match_name ---

// AVMatchName wraps av_match_name.
/*
  Match instances of a name in a comma-separated list of names.
  List entries are checked from the start to the end of the names list,
  the first match ends further processing. If an entry prefixed with '-'
  matches, then 0 is returned. The "ALL" list entry is considered to
  match all names.

  @param name  Name to look for.
  @param names List of names.
  @return 1 on match, 0 otherwise.
*/
func AVMatchName(name *CStr, names *CStr) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	var tmpnames *C.char
	if names != nil {
		tmpnames = names.ptr
	}
	ret := C.av_match_name(tmpname, tmpnames)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_append_path_component ---

// AVAppendPathComponent wraps av_append_path_component.
/*
  Append path component to the existing path.
  Path separator '/' is placed between when needed.
  Resulting string have to be freed with av_free().
  @param path      base path
  @param component component to be appended
  @return new path or NULL on error.
*/
func AVAppendPathComponent(path *CStr, component *CStr) *CStr {
	var tmppath *C.char
	if path != nil {
		tmppath = path.ptr
	}
	var tmpcomponent *C.char
	if component != nil {
		tmpcomponent = component.ptr
	}
	ret := C.av_append_path_component(tmppath, tmpcomponent)
	return wrapCStr(ret)
}

// --- Function av_escape ---

// av_escape skipped due to dst

// --- Function av_utf8_decode ---

// av_utf8_decode skipped due to codep (non-output primitive pointer)

// --- Function av_match_list ---

// AVMatchList wraps av_match_list.
/*
  Check if a name is in a list.
  @returns 0 if not found, or the 1 based index where it has been found in the
             list.
*/
func AVMatchList(name *CStr, list *CStr, separator uint8) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	var tmplist *C.char
	if list != nil {
		tmplist = list.ptr
	}
	ret := C.av_match_list(tmpname, tmplist, C.char(separator))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_sscanf ---

// av_sscanf skipped due to variadic arg.

// --- Function avutil_version ---

// AVUtilVersion wraps avutil_version.
//
//	Return the LIBAVUTIL_VERSION_INT constant.
func AVUtilVersion() uint {
	ret := C.avutil_version()
	return uint(ret)
}

// --- Function av_version_info ---

// AVVersionInfo wraps av_version_info.
/*
  Return an informative version string. This usually is the actual release
  version number or a git commit description. This string has no fixed format
  and can change any time. It should never be parsed by code.
*/
func AVVersionInfo() *CStr {
	ret := C.av_version_info()
	return wrapCStr(ret)
}

// --- Function avutil_configuration ---

// AVUtilConfiguration wraps avutil_configuration.
//
//	Return the libavutil build-time configuration.
func AVUtilConfiguration() *CStr {
	ret := C.avutil_configuration()
	return wrapCStr(ret)
}

// --- Function avutil_license ---

// AVUtilLicense wraps avutil_license.
//
//	Return the libavutil license.
func AVUtilLicense() *CStr {
	ret := C.avutil_license()
	return wrapCStr(ret)
}

// --- Function av_get_media_type_string ---

// AVGetMediaTypeString wraps av_get_media_type_string.
/*
  Return a string describing the media_type enum, NULL if media_type
  is unknown.
*/
func AVGetMediaTypeString(mediaType AVMediaType) *CStr {
	ret := C.av_get_media_type_string(C.enum_AVMediaType(mediaType))
	return wrapCStr(ret)
}

// --- Function av_get_picture_type_char ---

// AVGetPictureTypeChar wraps av_get_picture_type_char.
/*
  Return a single letter to describe the given picture type
  pict_type.

  @param[in] pict_type the picture type @return a single character
  representing the picture type, '?' if pict_type is unknown
*/
func AVGetPictureTypeChar(pictType AVPictureType) uint8 {
	ret := C.av_get_picture_type_char(C.enum_AVPictureType(pictType))
	return uint8(ret)
}

// --- Function av_x_if_null ---

// AVXIfNull wraps av_x_if_null.
//
//	Return x default pointer in case p is NULL.
func AVXIfNull(p unsafe.Pointer, x unsafe.Pointer) unsafe.Pointer {
	ret := C.av_x_if_null(p, x)
	return ret
}

// --- Function av_int_list_length_for_size ---

// AVIntListLengthForSize wraps av_int_list_length_for_size.
/*
  Compute the length of an integer list.

  @param elsize  size in bytes of each list element (only 1, 2, 4 or 8)
  @param term    list terminator (usually 0 or -1)
  @param list    pointer to the list
  @return  length of the list, in elements, not counting the terminator
*/
func AVIntListLengthForSize(elsize uint, list unsafe.Pointer, term uint64) uint {
	ret := C.av_int_list_length_for_size(C.uint(elsize), list, C.uint64_t(term))
	return uint(ret)
}

// --- Function av_get_time_base_q ---

// AVGetTimeBaseQ wraps av_get_time_base_q.
//
//	Return the fractional representation of the internal time base.
func AVGetTimeBaseQ() *AVRational {
	ret := C.av_get_time_base_q()
	return &AVRational{value: ret}
}

// --- Function av_fourcc_make_string ---

// AVFourccMakeString wraps av_fourcc_make_string.
/*
  Fill the provided buffer with a string containing a FourCC (four-character
  code) representation.

  @param buf    a buffer with size in bytes of at least AV_FOURCC_MAX_STRING_SIZE
  @param fourcc the fourcc to represent
  @return the buffer in input
*/
func AVFourccMakeString(buf *CStr, fourcc uint32) *CStr {
	var tmpbuf *C.char
	if buf != nil {
		tmpbuf = buf.ptr
	}
	ret := C.av_fourcc_make_string(tmpbuf, C.uint32_t(fourcc))
	return wrapCStr(ret)
}

// --- Function av_base64_decode ---

// AVBase64Decode wraps av_base64_decode.
/*
  Decode a base64-encoded string.

  @param out      buffer for decoded data
  @param in       null-terminated input string
  @param out_size size in bytes of the out buffer, must be at
                  least 3/4 of the length of in, that is AV_BASE64_DECODE_SIZE(strlen(in))
  @return         number of bytes written, or a negative value in case of
                  invalid input
*/
func AVBase64Decode(out unsafe.Pointer, in *CStr, outSize int) (int, error) {
	var tmpin *C.char
	if in != nil {
		tmpin = in.ptr
	}
	ret := C.av_base64_decode((*C.uint8_t)(out), tmpin, C.int(outSize))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_base64_encode ---

// AVBase64Encode wraps av_base64_encode.
/*
  Encode data to base64 and null-terminate.

  @param out      buffer for encoded data
  @param out_size size in bytes of the out buffer (including the
                  null terminator), must be at least AV_BASE64_SIZE(in_size)
  @param in       input buffer containing the data to encode
  @param in_size  size in bytes of the in buffer
  @return         out or NULL in case of error
*/
func AVBase64Encode(out *CStr, outSize int, in unsafe.Pointer, inSize int) *CStr {
	var tmpout *C.char
	if out != nil {
		tmpout = out.ptr
	}
	ret := C.av_base64_encode(tmpout, C.int(outSize), (*C.uint8_t)(in), C.int(inSize))
	return wrapCStr(ret)
}

// --- Function av_blowfish_alloc ---

// AVBlowfishAlloc wraps av_blowfish_alloc.
//
//	Allocate an AVBlowfish context.
func AVBlowfishAlloc() *AVBlowfish {
	ret := C.av_blowfish_alloc()
	var retMapped *AVBlowfish
	if ret != nil {
		retMapped = &AVBlowfish{ptr: ret}
	}
	return retMapped
}

// --- Function av_blowfish_init ---

// AVBlowfishInit wraps av_blowfish_init.
/*
  Initialize an AVBlowfish context.

  @param ctx an AVBlowfish context
  @param key a key
  @param key_len length of the key
*/
func AVBlowfishInit(ctx *AVBlowfish, key unsafe.Pointer, keyLen int) {
	var tmpctx *C.AVBlowfish
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	C.av_blowfish_init(tmpctx, (*C.uint8_t)(key), C.int(keyLen))
}

// --- Function av_blowfish_crypt_ecb ---

// av_blowfish_crypt_ecb skipped due to xl (non-output primitive pointer)

// --- Function av_blowfish_crypt ---

// AVBlowfishCrypt wraps av_blowfish_crypt.
/*
  Encrypt or decrypt a buffer using a previously initialized context.

  @param ctx an AVBlowfish context
  @param dst destination array, can be equal to src
  @param src source array, can be equal to dst
  @param count number of 8 byte blocks
  @param iv initialization vector for CBC mode, if NULL ECB will be used
  @param decrypt 0 for encryption, 1 for decryption
*/
func AVBlowfishCrypt(ctx *AVBlowfish, dst unsafe.Pointer, src unsafe.Pointer, count int, iv unsafe.Pointer, decrypt int) {
	var tmpctx *C.AVBlowfish
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	C.av_blowfish_crypt(tmpctx, (*C.uint8_t)(dst), (*C.uint8_t)(src), C.int(count), (*C.uint8_t)(iv), C.int(decrypt))
}

// --- Function av_bprint_init ---

// AVBprintInit wraps av_bprint_init.
/*
  Init a print buffer.

  @param buf        buffer to init
  @param size_init  initial size (including the final 0)
  @param size_max   maximum size;
                    - `0` means do not write anything, just count the length
                    - `1` is replaced by the maximum value for automatic storage
                        any large value means that the internal buffer will be
                        reallocated as needed up to that limit
                    - `-1` is converted to `UINT_MAX`, the largest limit possible.
                    Check also `AV_BPRINT_SIZE_*` macros.
*/
func AVBprintInit(buf *AVBPrint, sizeInit uint, sizeMax uint) {
	var tmpbuf *C.AVBPrint
	if buf != nil {
		tmpbuf = buf.ptr
	}
	C.av_bprint_init(tmpbuf, C.uint(sizeInit), C.uint(sizeMax))
}

// --- Function av_bprint_init_for_buffer ---

// AVBprintInitForBuffer wraps av_bprint_init_for_buffer.
/*
  Init a print buffer using a pre-existing buffer.

  The buffer will not be reallocated.
  In case size equals zero, the AVBPrint will be initialized to use
  the internal buffer as if using AV_BPRINT_SIZE_COUNT_ONLY with
  av_bprint_init().

  @param buf     buffer structure to init
  @param buffer  byte buffer to use for the string data
  @param size    size of buffer
*/
func AVBprintInitForBuffer(buf *AVBPrint, buffer *CStr, size uint) {
	var tmpbuf *C.AVBPrint
	if buf != nil {
		tmpbuf = buf.ptr
	}
	var tmpbuffer *C.char
	if buffer != nil {
		tmpbuffer = buffer.ptr
	}
	C.av_bprint_init_for_buffer(tmpbuf, tmpbuffer, C.uint(size))
}

// --- Function av_bprintf ---

// av_bprintf skipped due to variadic arg.

// --- Function av_vbprintf ---

// av_vbprintf skipped due to vl_arg.

// --- Function av_bprint_chars ---

// AVBprintChars wraps av_bprint_chars.
//
//	Append char c n times to a print buffer.
func AVBprintChars(buf *AVBPrint, c uint8, n uint) {
	var tmpbuf *C.AVBPrint
	if buf != nil {
		tmpbuf = buf.ptr
	}
	C.av_bprint_chars(tmpbuf, C.char(c), C.uint(n))
}

// --- Function av_bprint_append_data ---

// AVBprintAppendData wraps av_bprint_append_data.
/*
  Append data to a print buffer.

  @param buf  bprint buffer to use
  @param data pointer to data
  @param size size of data
*/
func AVBprintAppendData(buf *AVBPrint, data *CStr, size uint) {
	var tmpbuf *C.AVBPrint
	if buf != nil {
		tmpbuf = buf.ptr
	}
	var tmpdata *C.char
	if data != nil {
		tmpdata = data.ptr
	}
	C.av_bprint_append_data(tmpbuf, tmpdata, C.uint(size))
}

// --- Function av_bprint_strftime ---

// av_bprint_strftime skipped due to tm.

// --- Function av_bprint_get_buffer ---

// av_bprint_get_buffer skipped due to mem

// --- Function av_bprint_clear ---

// AVBprintClear wraps av_bprint_clear.
//
//	Reset the string to "" but keep internal allocated data.
func AVBprintClear(buf *AVBPrint) {
	var tmpbuf *C.AVBPrint
	if buf != nil {
		tmpbuf = buf.ptr
	}
	C.av_bprint_clear(tmpbuf)
}

// --- Function av_bprint_is_complete ---

// AVBprintIsComplete wraps av_bprint_is_complete.
/*
  Test if the print buffer is complete (not truncated).

  It may have been truncated due to a memory allocation failure
  or the size_max limit (compare size and size_max if necessary).
*/
func AVBprintIsComplete(buf *AVBPrint) (int, error) {
	var tmpbuf *C.AVBPrint
	if buf != nil {
		tmpbuf = buf.ptr
	}
	ret := C.av_bprint_is_complete(tmpbuf)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_bprint_finalize ---

// av_bprint_finalize skipped due to retStr

// --- Function av_bprint_escape ---

// AVBprintEscape wraps av_bprint_escape.
/*
  Escape the content in src and append it to dstbuf.

  @param dstbuf        already inited destination bprint buffer
  @param src           string containing the text to escape
  @param special_chars string containing the special characters which
                       need to be escaped, can be NULL
  @param mode          escape mode to employ, see AV_ESCAPE_MODE_* macros.
                       Any unknown value for mode will be considered equivalent to
                       AV_ESCAPE_MODE_BACKSLASH, but this behaviour can change without
                       notice.
  @param flags         flags which control how to escape, see AV_ESCAPE_FLAG_* macros
*/
func AVBprintEscape(dstbuf *AVBPrint, src *CStr, specialChars *CStr, mode AVEscapeMode, flags int) {
	var tmpdstbuf *C.AVBPrint
	if dstbuf != nil {
		tmpdstbuf = dstbuf.ptr
	}
	var tmpsrc *C.char
	if src != nil {
		tmpsrc = src.ptr
	}
	var tmpspecialChars *C.char
	if specialChars != nil {
		tmpspecialChars = specialChars.ptr
	}
	C.av_bprint_escape(tmpdstbuf, tmpsrc, tmpspecialChars, C.enum_AVEscapeMode(mode), C.int(flags))
}

// --- Function av_bswap16 ---

// AVBswap16 wraps av_bswap16.
func AVBswap16(x uint16) uint16 {
	ret := C.av_bswap16(C.uint16_t(x))
	return uint16(ret)
}

// --- Function av_bswap32 ---

// AVBswap32 wraps av_bswap32.
func AVBswap32(x uint32) uint32 {
	ret := C.av_bswap32(C.uint32_t(x))
	return uint32(ret)
}

// --- Function av_bswap64 ---

// AVBswap64 wraps av_bswap64.
func AVBswap64(x uint64) uint64 {
	ret := C.av_bswap64(C.uint64_t(x))
	return uint64(ret)
}

// --- Function av_buffer_alloc ---

// AVBufferAlloc wraps av_buffer_alloc.
/*
  Allocate an AVBuffer of the given size using av_malloc().

  @return an AVBufferRef of given size or NULL when out of memory
*/
func AVBufferAlloc(size uint64) *AVBufferRef {
	ret := C.av_buffer_alloc(C.size_t(size))
	var retMapped *AVBufferRef
	if ret != nil {
		retMapped = &AVBufferRef{ptr: ret}
	}
	return retMapped
}

// --- Function av_buffer_allocz ---

// AVBufferAllocz wraps av_buffer_allocz.
/*
  Same as av_buffer_alloc(), except the returned buffer will be initialized
  to zero.
*/
func AVBufferAllocz(size uint64) *AVBufferRef {
	ret := C.av_buffer_allocz(C.size_t(size))
	var retMapped *AVBufferRef
	if ret != nil {
		retMapped = &AVBufferRef{ptr: ret}
	}
	return retMapped
}

// --- Function av_buffer_create ---

// av_buffer_create skipped due to free.

// --- Function av_buffer_default_free ---

// AVBufferDefaultFree wraps av_buffer_default_free.
/*
  Default free callback, which calls av_free() on the buffer data.
  This function is meant to be passed to av_buffer_create(), not called
  directly.
*/
func AVBufferDefaultFree(opaque unsafe.Pointer, data unsafe.Pointer) {
	C.av_buffer_default_free(opaque, (*C.uint8_t)(data))
}

// --- Function av_buffer_ref ---

// AVBufferRef_ wraps av_buffer_ref.
/*
  Create a new reference to an AVBuffer.

  @return a new AVBufferRef referring to the same AVBuffer as buf or NULL on
  failure.
*/
func AVBufferRef_(buf *AVBufferRef) *AVBufferRef {
	var tmpbuf *C.AVBufferRef
	if buf != nil {
		tmpbuf = buf.ptr
	}
	ret := C.av_buffer_ref(tmpbuf)
	var retMapped *AVBufferRef
	if ret != nil {
		retMapped = &AVBufferRef{ptr: ret}
	}
	return retMapped
}

// --- Function av_buffer_unref ---

// AVBufferUnref wraps av_buffer_unref.
/*
  Free a given reference and automatically free the buffer if there are no more
  references to it.

  @param buf the reference to be freed. The pointer is set to NULL on return.
*/
func AVBufferUnref(buf **AVBufferRef) {
	var ptrbuf **C.AVBufferRef
	var tmpbuf *C.AVBufferRef
	var oldTmpbuf *C.AVBufferRef
	if buf != nil {
		innerbuf := *buf
		if innerbuf != nil {
			tmpbuf = innerbuf.ptr
			oldTmpbuf = tmpbuf
		}
		ptrbuf = &tmpbuf
	}
	C.av_buffer_unref(ptrbuf)
	if tmpbuf != oldTmpbuf && buf != nil {
		if tmpbuf != nil {
			*buf = &AVBufferRef{ptr: tmpbuf}
		} else {
			*buf = nil
		}
	}
}

// --- Function av_buffer_is_writable ---

// AVBufferIsWritable wraps av_buffer_is_writable.
/*
  @return 1 if the caller may write to the data referred to by buf (which is
  true if and only if buf is the only reference to the underlying AVBuffer).
  Return 0 otherwise.
  A positive answer is valid until av_buffer_ref() is called on buf.
*/
func AVBufferIsWritable(buf *AVBufferRef) (int, error) {
	var tmpbuf *C.AVBufferRef
	if buf != nil {
		tmpbuf = buf.ptr
	}
	ret := C.av_buffer_is_writable(tmpbuf)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_buffer_get_opaque ---

// AVBufferGetOpaque wraps av_buffer_get_opaque.
//
//	@return the opaque parameter set by av_buffer_create.
func AVBufferGetOpaque(buf *AVBufferRef) unsafe.Pointer {
	var tmpbuf *C.AVBufferRef
	if buf != nil {
		tmpbuf = buf.ptr
	}
	ret := C.av_buffer_get_opaque(tmpbuf)
	return ret
}

// --- Function av_buffer_get_ref_count ---

// AVBufferGetRefCount wraps av_buffer_get_ref_count.
func AVBufferGetRefCount(buf *AVBufferRef) (int, error) {
	var tmpbuf *C.AVBufferRef
	if buf != nil {
		tmpbuf = buf.ptr
	}
	ret := C.av_buffer_get_ref_count(tmpbuf)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_buffer_make_writable ---

// AVBufferMakeWritable wraps av_buffer_make_writable.
/*
  Create a writable reference from a given buffer reference, avoiding data copy
  if possible.

  @param buf buffer reference to make writable. On success, buf is either left
             untouched, or it is unreferenced and a new writable AVBufferRef is
             written in its place. On failure, buf is left untouched.
  @return 0 on success, a negative AVERROR on failure.
*/
func AVBufferMakeWritable(buf **AVBufferRef) (int, error) {
	var ptrbuf **C.AVBufferRef
	var tmpbuf *C.AVBufferRef
	var oldTmpbuf *C.AVBufferRef
	if buf != nil {
		innerbuf := *buf
		if innerbuf != nil {
			tmpbuf = innerbuf.ptr
			oldTmpbuf = tmpbuf
		}
		ptrbuf = &tmpbuf
	}
	ret := C.av_buffer_make_writable(ptrbuf)
	if tmpbuf != oldTmpbuf && buf != nil {
		if tmpbuf != nil {
			*buf = &AVBufferRef{ptr: tmpbuf}
		} else {
			*buf = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function av_buffer_realloc ---

// AVBufferRealloc wraps av_buffer_realloc.
/*
  Reallocate a given buffer.

  @param buf  a buffer reference to reallocate. On success, buf will be
              unreferenced and a new reference with the required size will be
              written in its place. On failure buf will be left untouched. *buf
              may be NULL, then a new buffer is allocated.
  @param size required new buffer size.
  @return 0 on success, a negative AVERROR on failure.

  @note the buffer is actually reallocated with av_realloc() only if it was
  initially allocated through av_buffer_realloc(NULL) and there is only one
  reference to it (i.e. the one passed to this function). In all other cases
  a new buffer is allocated and the data is copied.
*/
func AVBufferRealloc(buf **AVBufferRef, size uint64) (int, error) {
	var ptrbuf **C.AVBufferRef
	var tmpbuf *C.AVBufferRef
	var oldTmpbuf *C.AVBufferRef
	if buf != nil {
		innerbuf := *buf
		if innerbuf != nil {
			tmpbuf = innerbuf.ptr
			oldTmpbuf = tmpbuf
		}
		ptrbuf = &tmpbuf
	}
	ret := C.av_buffer_realloc(ptrbuf, C.size_t(size))
	if tmpbuf != oldTmpbuf && buf != nil {
		if tmpbuf != nil {
			*buf = &AVBufferRef{ptr: tmpbuf}
		} else {
			*buf = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function av_buffer_replace ---

// AVBufferReplace wraps av_buffer_replace.
/*
  Ensure dst refers to the same data as src.

  When *dst is already equivalent to src, do nothing. Otherwise unreference dst
  and replace it with a new reference to src.

  @param dst Pointer to either a valid buffer reference or NULL. On success,
             this will point to a buffer reference equivalent to src. On
             failure, dst will be left untouched.
  @param src A buffer reference to replace dst with. May be NULL, then this
             function is equivalent to av_buffer_unref(dst).
  @return 0 on success
          AVERROR(ENOMEM) on memory allocation failure.
*/
func AVBufferReplace(dst **AVBufferRef, src *AVBufferRef) (int, error) {
	var ptrdst **C.AVBufferRef
	var tmpdst *C.AVBufferRef
	var oldTmpdst *C.AVBufferRef
	if dst != nil {
		innerdst := *dst
		if innerdst != nil {
			tmpdst = innerdst.ptr
			oldTmpdst = tmpdst
		}
		ptrdst = &tmpdst
	}
	var tmpsrc *C.AVBufferRef
	if src != nil {
		tmpsrc = src.ptr
	}
	ret := C.av_buffer_replace(ptrdst, tmpsrc)
	if tmpdst != oldTmpdst && dst != nil {
		if tmpdst != nil {
			*dst = &AVBufferRef{ptr: tmpdst}
		} else {
			*dst = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function av_buffer_pool_init ---

// av_buffer_pool_init skipped due to alloc.

// --- Function av_buffer_pool_init2 ---

// av_buffer_pool_init2 skipped due to alloc.

// --- Function av_buffer_pool_uninit ---

// AVBufferPoolUninit wraps av_buffer_pool_uninit.
/*
  Mark the pool as being available for freeing. It will actually be freed only
  once all the allocated buffers associated with the pool are released. Thus it
  is safe to call this function while some of the allocated buffers are still
  in use.

  @param pool pointer to the pool to be freed. It will be set to NULL.
*/
func AVBufferPoolUninit(pool **AVBufferPool) {
	var ptrpool **C.AVBufferPool
	var tmppool *C.AVBufferPool
	var oldTmppool *C.AVBufferPool
	if pool != nil {
		innerpool := *pool
		if innerpool != nil {
			tmppool = innerpool.ptr
			oldTmppool = tmppool
		}
		ptrpool = &tmppool
	}
	C.av_buffer_pool_uninit(ptrpool)
	if tmppool != oldTmppool && pool != nil {
		if tmppool != nil {
			*pool = &AVBufferPool{ptr: tmppool}
		} else {
			*pool = nil
		}
	}
}

// --- Function av_buffer_pool_get ---

// AVBufferPoolGet wraps av_buffer_pool_get.
/*
  Allocate a new AVBuffer, reusing an old buffer from the pool when available.
  This function may be called simultaneously from multiple threads.

  @return a reference to the new buffer on success, NULL on error.
*/
func AVBufferPoolGet(pool *AVBufferPool) *AVBufferRef {
	var tmppool *C.AVBufferPool
	if pool != nil {
		tmppool = pool.ptr
	}
	ret := C.av_buffer_pool_get(tmppool)
	var retMapped *AVBufferRef
	if ret != nil {
		retMapped = &AVBufferRef{ptr: ret}
	}
	return retMapped
}

// --- Function av_buffer_pool_buffer_get_opaque ---

// AVBufferPoolBufferGetOpaque wraps av_buffer_pool_buffer_get_opaque.
/*
  Query the original opaque parameter of an allocated buffer in the pool.

  @param ref a buffer reference to a buffer returned by av_buffer_pool_get.
  @return the opaque parameter set by the buffer allocator function of the
          buffer pool.

  @note the opaque parameter of ref is used by the buffer pool implementation,
  therefore you have to use this function to access the original opaque
  parameter of an allocated buffer.
*/
func AVBufferPoolBufferGetOpaque(ref *AVBufferRef) unsafe.Pointer {
	var tmpref *C.AVBufferRef
	if ref != nil {
		tmpref = ref.ptr
	}
	ret := C.av_buffer_pool_buffer_get_opaque(tmpref)
	return ret
}

// --- Function av_camellia_alloc ---

// AVCamelliaAlloc wraps av_camellia_alloc.
/*
  Allocate an AVCAMELLIA context
  To free the struct: av_free(ptr)
*/
func AVCamelliaAlloc() *AVCAMELLIA {
	ret := C.av_camellia_alloc()
	var retMapped *AVCAMELLIA
	if ret != nil {
		retMapped = &AVCAMELLIA{ptr: ret}
	}
	return retMapped
}

// --- Function av_camellia_init ---

// AVCamelliaInit wraps av_camellia_init.
/*
  Initialize an AVCAMELLIA context.

  @param ctx an AVCAMELLIA context
  @param key a key of 16, 24, 32 bytes used for encryption/decryption
  @param key_bits number of keybits: possible are 128, 192, 256
*/
func AVCamelliaInit(ctx *AVCAMELLIA, key unsafe.Pointer, keyBits int) (int, error) {
	var tmpctx *C.struct_AVCAMELLIA
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	ret := C.av_camellia_init(tmpctx, (*C.uint8_t)(key), C.int(keyBits))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_camellia_crypt ---

// AVCamelliaCrypt wraps av_camellia_crypt.
/*
  Encrypt or decrypt a buffer using a previously initialized context

  @param ctx an AVCAMELLIA context
  @param dst destination array, can be equal to src
  @param src source array, can be equal to dst
  @param count number of 16 byte blocks
  @param iv initialization vector for CBC mode, NULL for ECB mode
  @param decrypt 0 for encryption, 1 for decryption
*/
func AVCamelliaCrypt(ctx *AVCAMELLIA, dst unsafe.Pointer, src unsafe.Pointer, count int, iv unsafe.Pointer, decrypt int) {
	var tmpctx *C.struct_AVCAMELLIA
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	C.av_camellia_crypt(tmpctx, (*C.uint8_t)(dst), (*C.uint8_t)(src), C.int(count), (*C.uint8_t)(iv), C.int(decrypt))
}

// --- Function av_cast5_alloc ---

// AVCast5Alloc wraps av_cast5_alloc.
/*
  Allocate an AVCAST5 context
  To free the struct: av_free(ptr)
*/
func AVCast5Alloc() *AVCAST5 {
	ret := C.av_cast5_alloc()
	var retMapped *AVCAST5
	if ret != nil {
		retMapped = &AVCAST5{ptr: ret}
	}
	return retMapped
}

// --- Function av_cast5_init ---

// AVCast5Init wraps av_cast5_init.
/*
  Initialize an AVCAST5 context.

  @param ctx an AVCAST5 context
  @param key a key of 5,6,...16 bytes used for encryption/decryption
  @param key_bits number of keybits: possible are 40,48,...,128
  @return 0 on success, less than 0 on failure
*/
func AVCast5Init(ctx *AVCAST5, key unsafe.Pointer, keyBits int) (int, error) {
	var tmpctx *C.struct_AVCAST5
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	ret := C.av_cast5_init(tmpctx, (*C.uint8_t)(key), C.int(keyBits))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_cast5_crypt ---

// AVCast5Crypt wraps av_cast5_crypt.
/*
  Encrypt or decrypt a buffer using a previously initialized context, ECB mode only

  @param ctx an AVCAST5 context
  @param dst destination array, can be equal to src
  @param src source array, can be equal to dst
  @param count number of 8 byte blocks
  @param decrypt 0 for encryption, 1 for decryption
*/
func AVCast5Crypt(ctx *AVCAST5, dst unsafe.Pointer, src unsafe.Pointer, count int, decrypt int) {
	var tmpctx *C.struct_AVCAST5
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	C.av_cast5_crypt(tmpctx, (*C.uint8_t)(dst), (*C.uint8_t)(src), C.int(count), C.int(decrypt))
}

// --- Function av_cast5_crypt2 ---

// AVCast5Crypt2 wraps av_cast5_crypt2.
/*
  Encrypt or decrypt a buffer using a previously initialized context

  @param ctx an AVCAST5 context
  @param dst destination array, can be equal to src
  @param src source array, can be equal to dst
  @param count number of 8 byte blocks
  @param iv initialization vector for CBC mode, NULL for ECB mode
  @param decrypt 0 for encryption, 1 for decryption
*/
func AVCast5Crypt2(ctx *AVCAST5, dst unsafe.Pointer, src unsafe.Pointer, count int, iv unsafe.Pointer, decrypt int) {
	var tmpctx *C.struct_AVCAST5
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	C.av_cast5_crypt2(tmpctx, (*C.uint8_t)(dst), (*C.uint8_t)(src), C.int(count), (*C.uint8_t)(iv), C.int(decrypt))
}

// --- Function av_channel_name ---

// AVChannelName wraps av_channel_name.
/*
  Get a human readable string in an abbreviated form describing a given channel.
  This is the inverse function of @ref av_channel_from_string().

  @param buf pre-allocated buffer where to put the generated string
  @param buf_size size in bytes of the buffer.
  @param channel the AVChannel whose name to get
  @return amount of bytes needed to hold the output string, or a negative AVERROR
          on failure. If the returned value is bigger than buf_size, then the
          string was truncated.
*/
func AVChannelName(buf *CStr, bufSize uint64, channel AVChannel) (int, error) {
	var tmpbuf *C.char
	if buf != nil {
		tmpbuf = buf.ptr
	}
	ret := C.av_channel_name(tmpbuf, C.size_t(bufSize), C.enum_AVChannel(channel))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_channel_name_bprint ---

// AVChannelNameBprint wraps av_channel_name_bprint.
/*
  bprint variant of av_channel_name().

  @note the string will be appended to the bprint buffer.
*/
func AVChannelNameBprint(bp *AVBPrint, channelId AVChannel) {
	var tmpbp *C.AVBPrint
	if bp != nil {
		tmpbp = bp.ptr
	}
	C.av_channel_name_bprint(tmpbp, C.enum_AVChannel(channelId))
}

// --- Function av_channel_description ---

// AVChannelDescription wraps av_channel_description.
/*
  Get a human readable string describing a given channel.

  @param buf pre-allocated buffer where to put the generated string
  @param buf_size size in bytes of the buffer.
  @param channel the AVChannel whose description to get
  @return amount of bytes needed to hold the output string, or a negative AVERROR
          on failure. If the returned value is bigger than buf_size, then the
          string was truncated.
*/
func AVChannelDescription(buf *CStr, bufSize uint64, channel AVChannel) (int, error) {
	var tmpbuf *C.char
	if buf != nil {
		tmpbuf = buf.ptr
	}
	ret := C.av_channel_description(tmpbuf, C.size_t(bufSize), C.enum_AVChannel(channel))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_channel_description_bprint ---

// AVChannelDescriptionBprint wraps av_channel_description_bprint.
/*
  bprint variant of av_channel_description().

  @note the string will be appended to the bprint buffer.
*/
func AVChannelDescriptionBprint(bp *AVBPrint, channelId AVChannel) {
	var tmpbp *C.AVBPrint
	if bp != nil {
		tmpbp = bp.ptr
	}
	C.av_channel_description_bprint(tmpbp, C.enum_AVChannel(channelId))
}

// --- Function av_channel_from_string ---

// AVChannelFromString wraps av_channel_from_string.
/*
  This is the inverse function of @ref av_channel_name().

  @return the channel with the given name
          AV_CHAN_NONE when name does not identify a known channel
*/
func AVChannelFromString(name *CStr) AVChannel {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_channel_from_string(tmpname)
	return AVChannel(ret)
}

// --- Function av_channel_layout_custom_init ---

// AVChannelLayoutCustomInit wraps av_channel_layout_custom_init.
/*
  Initialize a custom channel layout with the specified number of channels.
  The channel map will be allocated and the designation of all channels will
  be set to AV_CHAN_UNKNOWN.

  This is only a convenience helper function, a custom channel layout can also
  be constructed without using this.

  @param channel_layout the layout structure to be initialized
  @param nb_channels the number of channels

  @return 0 on success
          AVERROR(EINVAL) if the number of channels <= 0
          AVERROR(ENOMEM) if the channel map could not be allocated
*/
func AVChannelLayoutCustomInit(channelLayout *AVChannelLayout, nbChannels int) (int, error) {
	var tmpchannelLayout *C.AVChannelLayout
	if channelLayout != nil {
		tmpchannelLayout = channelLayout.ptr
	}
	ret := C.av_channel_layout_custom_init(tmpchannelLayout, C.int(nbChannels))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_channel_layout_from_mask ---

// AVChannelLayoutFromMask wraps av_channel_layout_from_mask.
/*
  Initialize a native channel layout from a bitmask indicating which channels
  are present.

  @param channel_layout the layout structure to be initialized
  @param mask bitmask describing the channel layout

  @return 0 on success
          AVERROR(EINVAL) for invalid mask values
*/
func AVChannelLayoutFromMask(channelLayout *AVChannelLayout, mask uint64) (int, error) {
	var tmpchannelLayout *C.AVChannelLayout
	if channelLayout != nil {
		tmpchannelLayout = channelLayout.ptr
	}
	ret := C.av_channel_layout_from_mask(tmpchannelLayout, C.uint64_t(mask))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_channel_layout_from_string ---

// AVChannelLayoutFromString wraps av_channel_layout_from_string.
/*
  Initialize a channel layout from a given string description.
  The input string can be represented by:
   - the formal channel layout name (returned by av_channel_layout_describe())
   - single or multiple channel names (returned by av_channel_name(), eg. "FL",
     or concatenated with "+", each optionally containing a custom name after
     a "@", eg. "FL@Left+FR@Right+LFE")
   - a decimal or hexadecimal value of a native channel layout (eg. "4" or "0x4")
   - the number of channels with default layout (eg. "4c")
   - the number of unordered channels (eg. "4C" or "4 channels")
   - the ambisonic order followed by optional non-diegetic channels (eg.
     "ambisonic 2+stereo")
  On error, the channel layout will remain uninitialized, but not necessarily
  untouched.

  @param channel_layout uninitialized channel layout for the result
  @param str string describing the channel layout
  @return 0 on success parsing the channel layout
          AVERROR(EINVAL) if an invalid channel layout string was provided
          AVERROR(ENOMEM) if there was not enough memory
*/
func AVChannelLayoutFromString(channelLayout *AVChannelLayout, str *CStr) (int, error) {
	var tmpchannelLayout *C.AVChannelLayout
	if channelLayout != nil {
		tmpchannelLayout = channelLayout.ptr
	}
	var tmpstr *C.char
	if str != nil {
		tmpstr = str.ptr
	}
	ret := C.av_channel_layout_from_string(tmpchannelLayout, tmpstr)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_channel_layout_default ---

// AVChannelLayoutDefault wraps av_channel_layout_default.
/*
  Get the default channel layout for a given number of channels.

  @param ch_layout the layout structure to be initialized
  @param nb_channels number of channels
*/
func AVChannelLayoutDefault(chLayout *AVChannelLayout, nbChannels int) {
	var tmpchLayout *C.AVChannelLayout
	if chLayout != nil {
		tmpchLayout = chLayout.ptr
	}
	C.av_channel_layout_default(tmpchLayout, C.int(nbChannels))
}

// --- Function av_channel_layout_standard ---

// av_channel_layout_standard skipped due to opaque

// --- Function av_channel_layout_uninit ---

// AVChannelLayoutUninit wraps av_channel_layout_uninit.
/*
  Free any allocated data in the channel layout and reset the channel
  count to 0.

  @param channel_layout the layout structure to be uninitialized
*/
func AVChannelLayoutUninit(channelLayout *AVChannelLayout) {
	var tmpchannelLayout *C.AVChannelLayout
	if channelLayout != nil {
		tmpchannelLayout = channelLayout.ptr
	}
	C.av_channel_layout_uninit(tmpchannelLayout)
}

// --- Function av_channel_layout_copy ---

// AVChannelLayoutCopy wraps av_channel_layout_copy.
/*
  Make a copy of a channel layout. This differs from just assigning src to dst
  in that it allocates and copies the map for AV_CHANNEL_ORDER_CUSTOM.

  @note the destination channel_layout will be always uninitialized before copy.

  @param dst destination channel layout
  @param src source channel layout
  @return 0 on success, a negative AVERROR on error.
*/
func AVChannelLayoutCopy(dst *AVChannelLayout, src *AVChannelLayout) (int, error) {
	var tmpdst *C.AVChannelLayout
	if dst != nil {
		tmpdst = dst.ptr
	}
	var tmpsrc *C.AVChannelLayout
	if src != nil {
		tmpsrc = src.ptr
	}
	ret := C.av_channel_layout_copy(tmpdst, tmpsrc)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_channel_layout_describe ---

// AVChannelLayoutDescribe wraps av_channel_layout_describe.
/*
  Get a human-readable string describing the channel layout properties.
  The string will be in the same format that is accepted by
  @ref av_channel_layout_from_string(), allowing to rebuild the same
  channel layout, except for opaque pointers.

  @param channel_layout channel layout to be described
  @param buf pre-allocated buffer where to put the generated string
  @param buf_size size in bytes of the buffer.
  @return amount of bytes needed to hold the output string, or a negative AVERROR
          on failure. If the returned value is bigger than buf_size, then the
          string was truncated.
*/
func AVChannelLayoutDescribe(channelLayout *AVChannelLayout, buf *CStr, bufSize uint64) (int, error) {
	var tmpchannelLayout *C.AVChannelLayout
	if channelLayout != nil {
		tmpchannelLayout = channelLayout.ptr
	}
	var tmpbuf *C.char
	if buf != nil {
		tmpbuf = buf.ptr
	}
	ret := C.av_channel_layout_describe(tmpchannelLayout, tmpbuf, C.size_t(bufSize))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_channel_layout_describe_bprint ---

// AVChannelLayoutDescribeBprint wraps av_channel_layout_describe_bprint.
/*
  bprint variant of av_channel_layout_describe().

  @note the string will be appended to the bprint buffer.
  @return 0 on success, or a negative AVERROR value on failure.
*/
func AVChannelLayoutDescribeBprint(channelLayout *AVChannelLayout, bp *AVBPrint) (int, error) {
	var tmpchannelLayout *C.AVChannelLayout
	if channelLayout != nil {
		tmpchannelLayout = channelLayout.ptr
	}
	var tmpbp *C.AVBPrint
	if bp != nil {
		tmpbp = bp.ptr
	}
	ret := C.av_channel_layout_describe_bprint(tmpchannelLayout, tmpbp)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_channel_layout_channel_from_index ---

// AVChannelLayoutChannelFromIndex wraps av_channel_layout_channel_from_index.
/*
  Get the channel with the given index in a channel layout.

  @param channel_layout input channel layout
  @param idx index of the channel
  @return channel with the index idx in channel_layout on success or
          AV_CHAN_NONE on failure (if idx is not valid or the channel order is
          unspecified)
*/
func AVChannelLayoutChannelFromIndex(channelLayout *AVChannelLayout, idx uint) AVChannel {
	var tmpchannelLayout *C.AVChannelLayout
	if channelLayout != nil {
		tmpchannelLayout = channelLayout.ptr
	}
	ret := C.av_channel_layout_channel_from_index(tmpchannelLayout, C.uint(idx))
	return AVChannel(ret)
}

// --- Function av_channel_layout_index_from_channel ---

// AVChannelLayoutIndexFromChannel wraps av_channel_layout_index_from_channel.
/*
  Get the index of a given channel in a channel layout. In case multiple
  channels are found, only the first match will be returned.

  @param channel_layout input channel layout
  @param channel the channel whose index to obtain
  @return index of channel in channel_layout on success or a negative number if
          channel is not present in channel_layout.
*/
func AVChannelLayoutIndexFromChannel(channelLayout *AVChannelLayout, channel AVChannel) (int, error) {
	var tmpchannelLayout *C.AVChannelLayout
	if channelLayout != nil {
		tmpchannelLayout = channelLayout.ptr
	}
	ret := C.av_channel_layout_index_from_channel(tmpchannelLayout, C.enum_AVChannel(channel))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_channel_layout_index_from_string ---

// AVChannelLayoutIndexFromString wraps av_channel_layout_index_from_string.
/*
  Get the index in a channel layout of a channel described by the given string.
  In case multiple channels are found, only the first match will be returned.

  This function accepts channel names in the same format as
  @ref av_channel_from_string().

  @param channel_layout input channel layout
  @param name string describing the channel whose index to obtain
  @return a channel index described by the given string, or a negative AVERROR
          value.
*/
func AVChannelLayoutIndexFromString(channelLayout *AVChannelLayout, name *CStr) (int, error) {
	var tmpchannelLayout *C.AVChannelLayout
	if channelLayout != nil {
		tmpchannelLayout = channelLayout.ptr
	}
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_channel_layout_index_from_string(tmpchannelLayout, tmpname)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_channel_layout_channel_from_string ---

// AVChannelLayoutChannelFromString wraps av_channel_layout_channel_from_string.
/*
  Get a channel described by the given string.

  This function accepts channel names in the same format as
  @ref av_channel_from_string().

  @param channel_layout input channel layout
  @param name string describing the channel to obtain
  @return a channel described by the given string in channel_layout on success
          or AV_CHAN_NONE on failure (if the string is not valid or the channel
          order is unspecified)
*/
func AVChannelLayoutChannelFromString(channelLayout *AVChannelLayout, name *CStr) AVChannel {
	var tmpchannelLayout *C.AVChannelLayout
	if channelLayout != nil {
		tmpchannelLayout = channelLayout.ptr
	}
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_channel_layout_channel_from_string(tmpchannelLayout, tmpname)
	return AVChannel(ret)
}

// --- Function av_channel_layout_subset ---

// AVChannelLayoutSubset wraps av_channel_layout_subset.
/*
  Find out what channels from a given set are present in a channel layout,
  without regard for their positions.

  @param channel_layout input channel layout
  @param mask a combination of AV_CH_* representing a set of channels
  @return a bitfield representing all the channels from mask that are present
          in channel_layout
*/
func AVChannelLayoutSubset(channelLayout *AVChannelLayout, mask uint64) uint64 {
	var tmpchannelLayout *C.AVChannelLayout
	if channelLayout != nil {
		tmpchannelLayout = channelLayout.ptr
	}
	ret := C.av_channel_layout_subset(tmpchannelLayout, C.uint64_t(mask))
	return uint64(ret)
}

// --- Function av_channel_layout_check ---

// AVChannelLayoutCheck wraps av_channel_layout_check.
/*
  Check whether a channel layout is valid, i.e. can possibly describe audio
  data.

  @param channel_layout input channel layout
  @return 1 if channel_layout is valid, 0 otherwise.
*/
func AVChannelLayoutCheck(channelLayout *AVChannelLayout) (int, error) {
	var tmpchannelLayout *C.AVChannelLayout
	if channelLayout != nil {
		tmpchannelLayout = channelLayout.ptr
	}
	ret := C.av_channel_layout_check(tmpchannelLayout)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_channel_layout_compare ---

// AVChannelLayoutCompare wraps av_channel_layout_compare.
/*
  Check whether two channel layouts are semantically the same, i.e. the same
  channels are present on the same positions in both.

  If one of the channel layouts is AV_CHANNEL_ORDER_UNSPEC, while the other is
  not, they are considered to be unequal. If both are AV_CHANNEL_ORDER_UNSPEC,
  they are considered equal iff the channel counts are the same in both.

  @param chl input channel layout
  @param chl1 input channel layout
  @return 0 if chl and chl1 are equal, 1 if they are not equal. A negative
          AVERROR code if one or both are invalid.
*/
func AVChannelLayoutCompare(chl *AVChannelLayout, chl1 *AVChannelLayout) (int, error) {
	var tmpchl *C.AVChannelLayout
	if chl != nil {
		tmpchl = chl.ptr
	}
	var tmpchl1 *C.AVChannelLayout
	if chl1 != nil {
		tmpchl1 = chl1.ptr
	}
	ret := C.av_channel_layout_compare(tmpchl, tmpchl1)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_channel_layout_ambisonic_order ---

// AVChannelLayoutAmbisonicOrder wraps av_channel_layout_ambisonic_order.
/*
  Return the order if the layout is n-th order standard-order ambisonic.
  The presence of optional extra non-diegetic channels at the end is not taken
  into account.

  @param channel_layout input channel layout
  @return the order of the layout, a negative error code otherwise.
*/
func AVChannelLayoutAmbisonicOrder(channelLayout *AVChannelLayout) (int, error) {
	var tmpchannelLayout *C.AVChannelLayout
	if channelLayout != nil {
		tmpchannelLayout = channelLayout.ptr
	}
	ret := C.av_channel_layout_ambisonic_order(tmpchannelLayout)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_channel_layout_retype ---

// AVChannelLayoutRetype wraps av_channel_layout_retype.
/*
  Change the AVChannelOrder of a channel layout.

  Change of AVChannelOrder can be either lossless or lossy. In case of a
  lossless conversion all the channel designations and the associated channel
  names (if any) are kept. On a lossy conversion the channel names and channel
  designations might be lost depending on the capabilities of the desired
  AVChannelOrder. Note that some conversions are simply not possible in which
  case this function returns AVERROR(ENOSYS).

  The following conversions are supported:

  Any       -> Custom     : Always possible, always lossless.
  Any       -> Unspecified: Always possible, lossless if channel designations
    are all unknown and channel names are not used, lossy otherwise.
  Custom    -> Ambisonic  : Possible if it contains ambisonic channels with
    optional non-diegetic channels in the end. Lossy if the channels have
    custom names, lossless otherwise.
  Custom    -> Native     : Possible if it contains native channels in native
      order. Lossy if the channels have custom names, lossless otherwise.

  On error this function keeps the original channel layout untouched.

  @param channel_layout channel layout which will be changed
  @param order the desired channel layout order
  @param flags a combination of AV_CHANNEL_LAYOUT_RETYPE_FLAG_* constants
  @return 0 if the conversion was successful and lossless or if the channel
            layout was already in the desired order
          >0 if the conversion was successful but lossy
          AVERROR(ENOSYS) if the conversion was not possible (or would be
            lossy and AV_CHANNEL_LAYOUT_RETYPE_FLAG_LOSSLESS was specified)
          AVERROR(EINVAL), AVERROR(ENOMEM) on error
*/
func AVChannelLayoutRetype(channelLayout *AVChannelLayout, order AVChannelOrder, flags int) (int, error) {
	var tmpchannelLayout *C.AVChannelLayout
	if channelLayout != nil {
		tmpchannelLayout = channelLayout.ptr
	}
	ret := C.av_channel_layout_retype(tmpchannelLayout, C.enum_AVChannelOrder(order), C.int(flags))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_log2 ---

// AVLog2 wraps av_log2.
func AVLog2(v uint) (int, error) {
	ret := C.av_log2(C.uint(v))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_log2_16bit ---

// AVLog216Bit wraps av_log2_16bit.
func AVLog216Bit(v uint) (int, error) {
	ret := C.av_log2_16bit(C.uint(v))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_clip_c ---

// AVClipC wraps av_clip_c.
/*
  Clip a signed integer value into the amin-amax range.
  @param a value to clip
  @param amin minimum value of the clip range
  @param amax maximum value of the clip range
  @return clipped value
*/
func AVClipC(a int, amin int, amax int) (int, error) {
	ret := C.av_clip_c(C.int(a), C.int(amin), C.int(amax))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_clip64_c ---

// AVClip64C wraps av_clip64_c.
/*
  Clip a signed 64bit integer value into the amin-amax range.
  @param a value to clip
  @param amin minimum value of the clip range
  @param amax maximum value of the clip range
  @return clipped value
*/
func AVClip64C(a int64, amin int64, amax int64) int64 {
	ret := C.av_clip64_c(C.int64_t(a), C.int64_t(amin), C.int64_t(amax))
	return int64(ret)
}

// --- Function av_clip_uint8_c ---

// AVClipUint8C wraps av_clip_uint8_c.
/*
  Clip a signed integer value into the 0-255 range.
  @param a value to clip
  @return clipped value
*/
func AVClipUint8C(a int) uint8 {
	ret := C.av_clip_uint8_c(C.int(a))
	return uint8(ret)
}

// --- Function av_clip_int8_c ---

// AVClipInt8C wraps av_clip_int8_c.
/*
  Clip a signed integer value into the -128,127 range.
  @param a value to clip
  @return clipped value
*/
func AVClipInt8C(a int) int8 {
	ret := C.av_clip_int8_c(C.int(a))
	return int8(ret)
}

// --- Function av_clip_uint16_c ---

// AVClipUint16C wraps av_clip_uint16_c.
/*
  Clip a signed integer value into the 0-65535 range.
  @param a value to clip
  @return clipped value
*/
func AVClipUint16C(a int) uint16 {
	ret := C.av_clip_uint16_c(C.int(a))
	return uint16(ret)
}

// --- Function av_clip_int16_c ---

// AVClipInt16C wraps av_clip_int16_c.
/*
  Clip a signed integer value into the -32768,32767 range.
  @param a value to clip
  @return clipped value
*/
func AVClipInt16C(a int) int16 {
	ret := C.av_clip_int16_c(C.int(a))
	return int16(ret)
}

// --- Function av_clipl_int32_c ---

// AVCliplInt32C wraps av_clipl_int32_c.
/*
  Clip a signed 64-bit integer value into the -2147483648,2147483647 range.
  @param a value to clip
  @return clipped value
*/
func AVCliplInt32C(a int64) int32 {
	ret := C.av_clipl_int32_c(C.int64_t(a))
	return int32(ret)
}

// --- Function av_clip_intp2_c ---

// AVClipIntp2C wraps av_clip_intp2_c.
/*
  Clip a signed integer into the -(2^p),(2^p-1) range.
  @param  a value to clip
  @param  p bit position to clip at
  @return clipped value
*/
func AVClipIntp2C(a int, p int) (int, error) {
	ret := C.av_clip_intp2_c(C.int(a), C.int(p))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_clip_uintp2_c ---

// AVClipUintp2C wraps av_clip_uintp2_c.
/*
  Clip a signed integer to an unsigned power of two range.
  @param  a value to clip
  @param  p bit position to clip at
  @return clipped value
*/
func AVClipUintp2C(a int, p int) uint {
	ret := C.av_clip_uintp2_c(C.int(a), C.int(p))
	return uint(ret)
}

// --- Function av_zero_extend_c ---

// AVZeroExtendC wraps av_zero_extend_c.
/*
  Clear high bits from an unsigned integer starting with specific bit position
  @param  a value to clip
  @param  p bit position to clip at. Must be between 0 and 31.
  @return clipped value
*/
func AVZeroExtendC(a uint, p uint) uint {
	ret := C.av_zero_extend_c(C.uint(a), C.uint(p))
	return uint(ret)
}

// --- Function av_mod_uintp2_c ---

// AVModUintp2C wraps av_mod_uintp2_c.
func AVModUintp2C(a uint, p uint) uint {
	ret := C.av_mod_uintp2_c(C.uint(a), C.uint(p))
	return uint(ret)
}

// --- Function av_sat_add32_c ---

// AVSatAdd32C wraps av_sat_add32_c.
/*
  Add two signed 32-bit values with saturation.

  @param  a one value
  @param  b another value
  @return sum with signed saturation
*/
func AVSatAdd32C(a int, b int) (int, error) {
	ret := C.av_sat_add32_c(C.int(a), C.int(b))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_sat_dadd32_c ---

// AVSatDadd32C wraps av_sat_dadd32_c.
/*
  Add a doubled value to another value with saturation at both stages.

  @param  a first value
  @param  b value doubled and added to a
  @return sum sat(a + sat(2*b)) with signed saturation
*/
func AVSatDadd32C(a int, b int) (int, error) {
	ret := C.av_sat_dadd32_c(C.int(a), C.int(b))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_sat_sub32_c ---

// AVSatSub32C wraps av_sat_sub32_c.
/*
  Subtract two signed 32-bit values with saturation.

  @param  a one value
  @param  b another value
  @return difference with signed saturation
*/
func AVSatSub32C(a int, b int) (int, error) {
	ret := C.av_sat_sub32_c(C.int(a), C.int(b))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_sat_dsub32_c ---

// AVSatDsub32C wraps av_sat_dsub32_c.
/*
  Subtract a doubled value from another value with saturation at both stages.

  @param  a first value
  @param  b value doubled and subtracted from a
  @return difference sat(a - sat(2*b)) with signed saturation
*/
func AVSatDsub32C(a int, b int) (int, error) {
	ret := C.av_sat_dsub32_c(C.int(a), C.int(b))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_sat_add64_c ---

// AVSatAdd64C wraps av_sat_add64_c.
/*
  Add two signed 64-bit values with saturation.

  @param  a one value
  @param  b another value
  @return sum with signed saturation
*/
func AVSatAdd64C(a int64, b int64) int64 {
	ret := C.av_sat_add64_c(C.int64_t(a), C.int64_t(b))
	return int64(ret)
}

// --- Function av_sat_sub64_c ---

// AVSatSub64C wraps av_sat_sub64_c.
/*
  Subtract two signed 64-bit values with saturation.

  @param  a one value
  @param  b another value
  @return difference with signed saturation
*/
func AVSatSub64C(a int64, b int64) int64 {
	ret := C.av_sat_sub64_c(C.int64_t(a), C.int64_t(b))
	return int64(ret)
}

// --- Function av_clipf_c ---

// AVClipfC wraps av_clipf_c.
/*
  Clip a float value into the amin-amax range.
  If a is nan or -inf amin will be returned.
  If a is +inf amax will be returned.
  @param a value to clip
  @param amin minimum value of the clip range
  @param amax maximum value of the clip range
  @return clipped value
*/
func AVClipfC(a float32, amin float32, amax float32) float32 {
	ret := C.av_clipf_c(C.float(a), C.float(amin), C.float(amax))
	return float32(ret)
}

// --- Function av_clipd_c ---

// AVClipdC wraps av_clipd_c.
/*
  Clip a double value into the amin-amax range.
  If a is nan or -inf amin will be returned.
  If a is +inf amax will be returned.
  @param a value to clip
  @param amin minimum value of the clip range
  @param amax maximum value of the clip range
  @return clipped value
*/
func AVClipdC(a float64, amin float64, amax float64) float64 {
	ret := C.av_clipd_c(C.double(a), C.double(amin), C.double(amax))
	return float64(ret)
}

// --- Function av_ceil_log2_c ---

// AVCeilLog2C wraps av_ceil_log2_c.
/*
  Compute ceil(log2(x)).
  @param x value used to compute ceil(log2(x))
  @return computed ceiling of log2(x)
*/
func AVCeilLog2C(x int) (int, error) {
	ret := C.av_ceil_log2_c(C.int(x))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_popcount_c ---

// AVPopcountC wraps av_popcount_c.
/*
  Count number of bits set to one in x
  @param x value to count bits of
  @return the number of bits set to one in x
*/
func AVPopcountC(x uint32) (int, error) {
	ret := C.av_popcount_c(C.uint32_t(x))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_popcount64_c ---

// AVPopcount64C wraps av_popcount64_c.
/*
  Count number of bits set to one in x
  @param x value to count bits of
  @return the number of bits set to one in x
*/
func AVPopcount64C(x uint64) (int, error) {
	ret := C.av_popcount64_c(C.uint64_t(x))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_parity_c ---

// AVParityC wraps av_parity_c.
func AVParityC(v uint32) (int, error) {
	ret := C.av_parity_c(C.uint32_t(v))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_container_fifo_alloc ---

// av_container_fifo_alloc skipped due to container_alloc.

// --- Function av_container_fifo_alloc_avframe ---

// AVContainerFifoAllocAVFrame wraps av_container_fifo_alloc_avframe.
/*
  Allocate an AVContainerFifo instance for AVFrames.

  @param flags currently unused
*/
func AVContainerFifoAllocAVFrame(flags uint) *AVContainerFifo {
	ret := C.av_container_fifo_alloc_avframe(C.uint(flags))
	var retMapped *AVContainerFifo
	if ret != nil {
		retMapped = &AVContainerFifo{ptr: ret}
	}
	return retMapped
}

// --- Function av_container_fifo_free ---

// AVContainerFifoFree wraps av_container_fifo_free.
//
//	Free a AVContainerFifo and everything in it.
func AVContainerFifoFree(cf **AVContainerFifo) {
	var ptrcf **C.AVContainerFifo
	var tmpcf *C.AVContainerFifo
	var oldTmpcf *C.AVContainerFifo
	if cf != nil {
		innercf := *cf
		if innercf != nil {
			tmpcf = innercf.ptr
			oldTmpcf = tmpcf
		}
		ptrcf = &tmpcf
	}
	C.av_container_fifo_free(ptrcf)
	if tmpcf != oldTmpcf && cf != nil {
		if tmpcf != nil {
			*cf = &AVContainerFifo{ptr: tmpcf}
		} else {
			*cf = nil
		}
	}
}

// --- Function av_container_fifo_write ---

// AVContainerFifoWrite wraps av_container_fifo_write.
/*
  Write the contents of obj to the FIFO.

  The fifo_transfer() callback previously provided to av_container_fifo_alloc()
  will be called with obj as src in order to perform the actual transfer.
*/
func AVContainerFifoWrite(cf *AVContainerFifo, obj unsafe.Pointer, flags uint) (int, error) {
	var tmpcf *C.AVContainerFifo
	if cf != nil {
		tmpcf = cf.ptr
	}
	ret := C.av_container_fifo_write(tmpcf, obj, C.uint(flags))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_container_fifo_read ---

// AVContainerFifoRead wraps av_container_fifo_read.
/*
  Read the next available object from the FIFO into obj.

  The fifo_read() callback previously provided to av_container_fifo_alloc()
  will be called with obj as dst in order to perform the actual transfer.
*/
func AVContainerFifoRead(cf *AVContainerFifo, obj unsafe.Pointer, flags uint) (int, error) {
	var tmpcf *C.AVContainerFifo
	if cf != nil {
		tmpcf = cf.ptr
	}
	ret := C.av_container_fifo_read(tmpcf, obj, C.uint(flags))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_container_fifo_peek ---

// av_container_fifo_peek skipped due to pobj

// --- Function av_container_fifo_drain ---

// AVContainerFifoDrain wraps av_container_fifo_drain.
/*
  Discard the specified number of elements from the FIFO.

  @param nb_elems number of elements to discard, MUST NOT be larger than
                  av_fifo_can_read(f)
*/
func AVContainerFifoDrain(cf *AVContainerFifo, nbElems uint64) {
	var tmpcf *C.AVContainerFifo
	if cf != nil {
		tmpcf = cf.ptr
	}
	C.av_container_fifo_drain(tmpcf, C.size_t(nbElems))
}

// --- Function av_container_fifo_can_read ---

// AVContainerFifoCanRead wraps av_container_fifo_can_read.
//
//	@return number of objects available for reading
func AVContainerFifoCanRead(cf *AVContainerFifo) uint64 {
	var tmpcf *C.AVContainerFifo
	if cf != nil {
		tmpcf = cf.ptr
	}
	ret := C.av_container_fifo_can_read(tmpcf)
	return uint64(ret)
}

// --- Function av_get_cpu_flags ---

// AVGetCpuFlags wraps av_get_cpu_flags.
/*
  Return the flags which specify extensions supported by the CPU.
  The returned value is affected by av_force_cpu_flags() if that was used
  before. So av_get_cpu_flags() can easily be used in an application to
  detect the enabled cpu flags.
*/
func AVGetCpuFlags() (int, error) {
	ret := C.av_get_cpu_flags()
	return int(ret), WrapErr(int(ret))
}

// --- Function av_force_cpu_flags ---

// AVForceCpuFlags wraps av_force_cpu_flags.
/*
  Disables cpu detection and forces the specified flags.
  -1 is a special case that disables forcing of specific flags.
*/
func AVForceCpuFlags(flags int) {
	C.av_force_cpu_flags(C.int(flags))
}

// --- Function av_parse_cpu_caps ---

// av_parse_cpu_caps skipped due to flags (non-output primitive pointer)

// --- Function av_cpu_count ---

// AVCpuCount wraps av_cpu_count.
//
//	@return the number of logical CPU cores present.
func AVCpuCount() (int, error) {
	ret := C.av_cpu_count()
	return int(ret), WrapErr(int(ret))
}

// --- Function av_cpu_force_count ---

// AVCpuForceCount wraps av_cpu_force_count.
/*
  Overrides cpu count detection and forces the specified count.
  Count < 1 disables forcing of specific count.
*/
func AVCpuForceCount(count int) {
	C.av_cpu_force_count(C.int(count))
}

// --- Function av_cpu_max_align ---

// AVCpuMaxAlign wraps av_cpu_max_align.
/*
  Get the maximum data alignment that may be required by FFmpeg.

  Note that this is affected by the build configuration and the CPU flags mask,
  so e.g. if the CPU supports AVX, but libavutil has been built with
  --disable-avx or the AV_CPU_FLAG_AVX flag has been disabled through
   av_set_cpu_flags_mask(), then this function will behave as if AVX is not
   present.
*/
func AVCpuMaxAlign() uint64 {
	ret := C.av_cpu_max_align()
	return uint64(ret)
}

// --- Function av_crc_init ---

// AVCrcInit wraps av_crc_init.
/*
  Initialize a CRC table.
  @param ctx must be an array of size sizeof(AVCRC)*257 or sizeof(AVCRC)*1024
  @param le If 1, the lowest bit represents the coefficient for the highest
            exponent of the corresponding polynomial (both for poly and
            actual CRC).
            If 0, you must swap the CRC parameter and the result of av_crc
            if you need the standard representation (can be simplified in
            most cases to e.g. bswap16):
            av_bswap32(crc << (32-bits))
  @param bits number of bits for the CRC
  @param poly generator polynomial without the x**bits coefficient, in the
              representation as specified by le
  @param ctx_size size of ctx in bytes
  @return <0 on failure
*/
func AVCrcInit(ctx *AVCRC, le int, bits int, poly uint32, ctxSize int) (int, error) {
	var tmpctx *C.AVCRC
	if ctx != nil {
		tmpctx = (*C.AVCRC)(unsafe.Pointer(ctx))
	}
	ret := C.av_crc_init(tmpctx, C.int(le), C.int(bits), C.uint32_t(poly), C.int(ctxSize))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_crc_get_table ---

// AVCrcGetTable wraps av_crc_get_table.
/*
  Get an initialized standard CRC table.
  @param crc_id ID of a standard CRC
  @return a pointer to the CRC table or NULL on failure
*/
func AVCrcGetTable(crcId AVCRCId) *AVCRC {
	ret := C.av_crc_get_table(C.AVCRCId(crcId))
	return (*AVCRC)(unsafe.Pointer(ret))
}

// --- Function av_crc ---

// AVCrc wraps av_crc.
/*
  Calculate the CRC of a block.
  @param ctx initialized AVCRC array (see av_crc_init())
  @param crc CRC of previous blocks if any or initial value for CRC
  @param buffer buffer whose CRC to calculate
  @param length length of the buffer
  @return CRC updated with the data from the given block

  @see av_crc_init() "le" parameter
*/
func AVCrc(ctx *AVCRC, crc uint32, buffer unsafe.Pointer, length uint64) uint32 {
	var tmpctx *C.AVCRC
	if ctx != nil {
		tmpctx = (*C.AVCRC)(unsafe.Pointer(ctx))
	}
	ret := C.av_crc(tmpctx, C.uint32_t(crc), (*C.uint8_t)(buffer), C.size_t(length))
	return uint32(ret)
}

// --- Function av_csp_luma_coeffs_from_avcsp ---

// AVCspLumaCoeffsFromAVCsp wraps av_csp_luma_coeffs_from_avcsp.
/*
  Retrieves the Luma coefficients necessary to construct a conversion matrix
  from an enum constant describing the colorspace.
  @param csp An enum constant indicating YUV or similar colorspace.
  @return The Luma coefficients associated with that colorspace, or NULL
      if the constant is unknown to libavutil.
*/
func AVCspLumaCoeffsFromAVCsp(csp AVColorSpace) *AVLumaCoefficients {
	ret := C.av_csp_luma_coeffs_from_avcsp(C.enum_AVColorSpace(csp))
	var retMapped *AVLumaCoefficients
	if ret != nil {
		retMapped = &AVLumaCoefficients{ptr: ret}
	}
	return retMapped
}

// --- Function av_csp_primaries_desc_from_id ---

// AVCspPrimariesDescFromId wraps av_csp_primaries_desc_from_id.
/*
  Retrieves a complete gamut description from an enum constant describing the
  color primaries.
  @param prm An enum constant indicating primaries
  @return A description of the colorspace gamut associated with that enum
      constant, or NULL if the constant is unknown to libavutil.
*/
func AVCspPrimariesDescFromId(prm AVColorPrimaries) *AVColorPrimariesDesc {
	ret := C.av_csp_primaries_desc_from_id(C.enum_AVColorPrimaries(prm))
	var retMapped *AVColorPrimariesDesc
	if ret != nil {
		retMapped = &AVColorPrimariesDesc{ptr: ret}
	}
	return retMapped
}

// --- Function av_csp_primaries_id_from_desc ---

// AVCspPrimariesIdFromDesc wraps av_csp_primaries_id_from_desc.
/*
  Detects which enum AVColorPrimaries constant corresponds to the given complete
  gamut description.
  @see enum AVColorPrimaries
  @param prm A description of the colorspace gamut
  @return The enum constant associated with this gamut, or
      AVCOL_PRI_UNSPECIFIED if no clear match can be identified.
*/
func AVCspPrimariesIdFromDesc(prm *AVColorPrimariesDesc) AVColorPrimaries {
	var tmpprm *C.AVColorPrimariesDesc
	if prm != nil {
		tmpprm = prm.ptr
	}
	ret := C.av_csp_primaries_id_from_desc(tmpprm)
	return AVColorPrimaries(ret)
}

// --- Function av_csp_approximate_trc_gamma ---

// AVCspApproximateTrcGamma wraps av_csp_approximate_trc_gamma.
/*
  Determine a suitable 'gamma' value to match the supplied
  AVColorTransferCharacteristic.

  See Apple Technical Note TN2257 (https://developer.apple.com/library/mac/technotes/tn2257/_index.html)

  This function returns the gamma exponent for the OETF. For example, sRGB is approximated
  by gamma 2.2, not by gamma 0.45455.

  @return Will return an approximation to the simple gamma function matching
          the supplied Transfer Characteristic, Will return 0.0 for any
          we cannot reasonably match against.
*/
func AVCspApproximateTrcGamma(trc AVColorTransferCharacteristic) float64 {
	ret := C.av_csp_approximate_trc_gamma(C.enum_AVColorTransferCharacteristic(trc))
	return float64(ret)
}

// --- Function av_csp_trc_func_from_id ---

// AVCspTrcFuncFromId wraps av_csp_trc_func_from_id.
/*
  Determine the function needed to apply the given
  AVColorTransferCharacteristic to linear input.

  The function returned should expect a nominal domain and range of [0.0-1.0]
  values outside of this range maybe valid depending on the chosen
  characteristic function.

  @return Will return pointer to the function matching the
          supplied Transfer Characteristic. If unspecified will
          return NULL:
*/
func AVCspTrcFuncFromId(trc AVColorTransferCharacteristic) AVCspTrcFunction {
	ret := C.av_csp_trc_func_from_id(C.enum_AVColorTransferCharacteristic(trc))
	return AVCspTrcFunction(ret)
}

// --- Function av_csp_trc_func_inv_from_id ---

// AVCspTrcFuncInvFromId wraps av_csp_trc_func_inv_from_id.
//
//	Returns the mathematical inverse of the corresponding TRC function.
func AVCspTrcFuncInvFromId(trc AVColorTransferCharacteristic) AVCspTrcFunction {
	ret := C.av_csp_trc_func_inv_from_id(C.enum_AVColorTransferCharacteristic(trc))
	return AVCspTrcFunction(ret)
}

// --- Function av_csp_itu_eotf ---

// AVCspItuEotf wraps av_csp_itu_eotf.
/*
  Returns the ITU EOTF corresponding to a given TRC. This converts from the
  signal level [0,1] to the raw output display luminance in nits (cd/m^2).
  This is done per channel in RGB space, except for AVCOL_TRC_SMPTE428, which
  assumes CIE XYZ in- and output.

  @return A pointer to the function implementing the given TRC, or NULL if no
          such function is defined.

  @note In general, the resulting function is defined (wherever possible) for
        out-of-range values, even though these values do not have a physical
        meaning on the given display. Users should clamp inputs (or outputs)
        if this behavior is not desired.

        This is also the case for functions like PQ, which are defined over an
        absolute signal range independent of the target display capabilities.
*/
func AVCspItuEotf(trc AVColorTransferCharacteristic) AVCspEotfFunction {
	ret := C.av_csp_itu_eotf(C.enum_AVColorTransferCharacteristic(trc))
	return AVCspEotfFunction(ret)
}

// --- Function av_csp_itu_eotf_inv ---

// AVCspItuEotfInv wraps av_csp_itu_eotf_inv.
//
//	Returns the mathematical inverse of the corresponding EOTF.
func AVCspItuEotfInv(trc AVColorTransferCharacteristic) AVCspEotfFunction {
	ret := C.av_csp_itu_eotf_inv(C.enum_AVColorTransferCharacteristic(trc))
	return AVCspEotfFunction(ret)
}

// --- Function av_des_alloc ---

// AVDesAlloc wraps av_des_alloc.
//
//	Allocate an AVDES context.
func AVDesAlloc() *AVDES {
	ret := C.av_des_alloc()
	var retMapped *AVDES
	if ret != nil {
		retMapped = &AVDES{ptr: ret}
	}
	return retMapped
}

// --- Function av_des_init ---

// AVDesInit wraps av_des_init.
/*
  @brief Initializes an AVDES context.

  @param d pointer to a AVDES structure to initialize
  @param key pointer to the key to use
  @param key_bits must be 64 or 192
  @param decrypt 0 for encryption/CBC-MAC, 1 for decryption
  @return zero on success, negative value otherwise
*/
func AVDesInit(d *AVDES, key unsafe.Pointer, keyBits int, decrypt int) (int, error) {
	var tmpd *C.AVDES
	if d != nil {
		tmpd = d.ptr
	}
	ret := C.av_des_init(tmpd, (*C.uint8_t)(key), C.int(keyBits), C.int(decrypt))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_des_crypt ---

// AVDesCrypt wraps av_des_crypt.
/*
  @brief Encrypts / decrypts using the DES algorithm.

  @param d pointer to the AVDES structure
  @param dst destination array, can be equal to src, must be 8-byte aligned
  @param src source array, can be equal to dst, must be 8-byte aligned, may be NULL
  @param count number of 8 byte blocks
  @param iv initialization vector for CBC mode, if NULL then ECB will be used,
            must be 8-byte aligned
  @param decrypt 0 for encryption, 1 for decryption
*/
func AVDesCrypt(d *AVDES, dst unsafe.Pointer, src unsafe.Pointer, count int, iv unsafe.Pointer, decrypt int) {
	var tmpd *C.AVDES
	if d != nil {
		tmpd = d.ptr
	}
	C.av_des_crypt(tmpd, (*C.uint8_t)(dst), (*C.uint8_t)(src), C.int(count), (*C.uint8_t)(iv), C.int(decrypt))
}

// --- Function av_des_mac ---

// AVDesMac wraps av_des_mac.
/*
  @brief Calculates CBC-MAC using the DES algorithm.

  @param d pointer to the AVDES structure
  @param dst destination array, can be equal to src, must be 8-byte aligned
  @param src source array, can be equal to dst, must be 8-byte aligned, may be NULL
  @param count number of 8 byte blocks
*/
func AVDesMac(d *AVDES, dst unsafe.Pointer, src unsafe.Pointer, count int) {
	var tmpd *C.AVDES
	if d != nil {
		tmpd = d.ptr
	}
	C.av_des_mac(tmpd, (*C.uint8_t)(dst), (*C.uint8_t)(src), C.int(count))
}

// --- Function av_get_detection_bbox ---

// AVGetDetectionBbox wraps av_get_detection_bbox.
//
//	Get the bounding box at the specified {@code idx}. Must be between 0 and nb_bboxes.
func AVGetDetectionBbox(header *AVDetectionBBoxHeader, idx uint) *AVDetectionBBox {
	var tmpheader *C.AVDetectionBBoxHeader
	if header != nil {
		tmpheader = header.ptr
	}
	ret := C.av_get_detection_bbox(tmpheader, C.uint(idx))
	var retMapped *AVDetectionBBox
	if ret != nil {
		retMapped = &AVDetectionBBox{ptr: ret}
	}
	return retMapped
}

// --- Function av_detection_bbox_alloc ---

// AVDetectionBboxAlloc wraps av_detection_bbox_alloc.
/*
  Allocates memory for AVDetectionBBoxHeader, plus an array of {@code nb_bboxes}
  AVDetectionBBox, and initializes the variables.
  Can be freed with a normal av_free() call.

  @param nb_bboxes number of AVDetectionBBox structures to allocate
  @param out_size if non-NULL, the size in bytes of the resulting data array is
  written here.
*/
func AVDetectionBboxAlloc(nbBboxes uint32, outSize *uint64) *AVDetectionBBoxHeader {
	ret := C.av_detection_bbox_alloc(C.uint32_t(nbBboxes), (*C.size_t)(unsafe.Pointer(outSize)))
	var retMapped *AVDetectionBBoxHeader
	if ret != nil {
		retMapped = &AVDetectionBBoxHeader{ptr: ret}
	}
	return retMapped
}

// --- Function av_detection_bbox_create_side_data ---

// AVDetectionBboxCreateSideData wraps av_detection_bbox_create_side_data.
/*
  Allocates memory for AVDetectionBBoxHeader, plus an array of {@code nb_bboxes}
  AVDetectionBBox, in the given AVFrame {@code frame} as AVFrameSideData of type
  AV_FRAME_DATA_DETECTION_BBOXES and initializes the variables.
*/
func AVDetectionBboxCreateSideData(frame *AVFrame, nbBboxes uint32) *AVDetectionBBoxHeader {
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.av_detection_bbox_create_side_data(tmpframe, C.uint32_t(nbBboxes))
	var retMapped *AVDetectionBBoxHeader
	if ret != nil {
		retMapped = &AVDetectionBBoxHeader{ptr: ret}
	}
	return retMapped
}

// --- Function av_dict_get ---

// AVDictGet wraps av_dict_get.
/*
  Get a dictionary entry with matching key.

  The returned entry key or value must not be changed, or it will
  cause undefined behavior.

  @param prev  Set to the previous matching element to find the next.
               If set to NULL the first matching element is returned.
  @param key   Matching key
  @param flags A collection of AV_DICT_* flags controlling how the
               entry is retrieved

  @return      Found entry or NULL in case no matching entry was found in the dictionary
*/
func AVDictGet(m *AVDictionary, key *CStr, prev *AVDictionaryEntry, flags int) *AVDictionaryEntry {
	var tmpm *C.AVDictionary
	if m != nil {
		tmpm = m.ptr
	}
	var tmpkey *C.char
	if key != nil {
		tmpkey = key.ptr
	}
	var tmpprev *C.AVDictionaryEntry
	if prev != nil {
		tmpprev = prev.ptr
	}
	ret := C.av_dict_get(tmpm, tmpkey, tmpprev, C.int(flags))
	var retMapped *AVDictionaryEntry
	if ret != nil {
		retMapped = &AVDictionaryEntry{ptr: ret}
	}
	return retMapped
}

// --- Function av_dict_iterate ---

// AVDictIterate wraps av_dict_iterate.
/*
  Iterate over a dictionary

  Iterates through all entries in the dictionary.

  @warning The returned AVDictionaryEntry key/value must not be changed.

  @warning As av_dict_set() invalidates all previous entries returned
  by this function, it must not be called while iterating over the dict.

  Typical usage:
  @code
  const AVDictionaryEntry *e = NULL;
  while ((e = av_dict_iterate(m, e))) {
      // ...
  }
  @endcode

  @param m     The dictionary to iterate over
  @param prev  Pointer to the previous AVDictionaryEntry, NULL initially

  @retval AVDictionaryEntry* The next element in the dictionary
  @retval NULL               No more elements in the dictionary
*/
func AVDictIterate(m *AVDictionary, prev *AVDictionaryEntry) *AVDictionaryEntry {
	var tmpm *C.AVDictionary
	if m != nil {
		tmpm = m.ptr
	}
	var tmpprev *C.AVDictionaryEntry
	if prev != nil {
		tmpprev = prev.ptr
	}
	ret := C.av_dict_iterate(tmpm, tmpprev)
	var retMapped *AVDictionaryEntry
	if ret != nil {
		retMapped = &AVDictionaryEntry{ptr: ret}
	}
	return retMapped
}

// --- Function av_dict_count ---

// AVDictCount wraps av_dict_count.
/*
  Get number of entries in dictionary.

  @param m dictionary
  @return  number of entries in dictionary
*/
func AVDictCount(m *AVDictionary) (int, error) {
	var tmpm *C.AVDictionary
	if m != nil {
		tmpm = m.ptr
	}
	ret := C.av_dict_count(tmpm)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_dict_set ---

// AVDictSet wraps av_dict_set.
/*
  Set the given entry in *pm, overwriting an existing entry.

  Note: If AV_DICT_DONT_STRDUP_KEY or AV_DICT_DONT_STRDUP_VAL is set,
  these arguments will be freed on error.

  @warning Adding a new entry to a dictionary invalidates all existing entries
  previously returned with av_dict_get() or av_dict_iterate().

  @param pm        Pointer to a pointer to a dictionary struct. If *pm is NULL
                   a dictionary struct is allocated and put in *pm.
  @param key       Entry key to add to *pm (will either be av_strduped or added as a new key depending on flags)
  @param value     Entry value to add to *pm (will be av_strduped or added as a new key depending on flags).
                   Passing a NULL value will cause an existing entry to be deleted.

  @return          >= 0 on success otherwise an error code <0
*/
func AVDictSet(pm **AVDictionary, key *CStr, value *CStr, flags int) (int, error) {
	var ptrpm **C.AVDictionary
	var tmppm *C.AVDictionary
	var oldTmppm *C.AVDictionary
	if pm != nil {
		innerpm := *pm
		if innerpm != nil {
			tmppm = innerpm.ptr
			oldTmppm = tmppm
		}
		ptrpm = &tmppm
	}
	var tmpkey *C.char
	if key != nil {
		tmpkey = key.ptr
	}
	var tmpvalue *C.char
	if value != nil {
		tmpvalue = value.ptr
	}
	ret := C.av_dict_set(ptrpm, tmpkey, tmpvalue, C.int(flags))
	if tmppm != oldTmppm && pm != nil {
		if tmppm != nil {
			*pm = &AVDictionary{ptr: tmppm}
		} else {
			*pm = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function av_dict_set_int ---

// AVDictSetInt wraps av_dict_set_int.
/*
  Convenience wrapper for av_dict_set() that converts the value to a string
  and stores it.

  Note: If ::AV_DICT_DONT_STRDUP_KEY is set, key will be freed on error.
*/
func AVDictSetInt(pm **AVDictionary, key *CStr, value int64, flags int) (int, error) {
	var ptrpm **C.AVDictionary
	var tmppm *C.AVDictionary
	var oldTmppm *C.AVDictionary
	if pm != nil {
		innerpm := *pm
		if innerpm != nil {
			tmppm = innerpm.ptr
			oldTmppm = tmppm
		}
		ptrpm = &tmppm
	}
	var tmpkey *C.char
	if key != nil {
		tmpkey = key.ptr
	}
	ret := C.av_dict_set_int(ptrpm, tmpkey, C.int64_t(value), C.int(flags))
	if tmppm != oldTmppm && pm != nil {
		if tmppm != nil {
			*pm = &AVDictionary{ptr: tmppm}
		} else {
			*pm = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function av_dict_parse_string ---

// AVDictParseString wraps av_dict_parse_string.
/*
  Parse the key/value pairs list and add the parsed entries to a dictionary.

  In case of failure, all the successfully set entries are stored in
  *pm. You may need to manually free the created dictionary.

  @param key_val_sep  A 0-terminated list of characters used to separate
                      key from value
  @param pairs_sep    A 0-terminated list of characters used to separate
                      two pairs from each other
  @param flags        Flags to use when adding to the dictionary.
                      ::AV_DICT_DONT_STRDUP_KEY and ::AV_DICT_DONT_STRDUP_VAL
                      are ignored since the key/value tokens will always
                      be duplicated.

  @return             0 on success, negative AVERROR code on failure
*/
func AVDictParseString(pm **AVDictionary, str *CStr, keyValSep *CStr, pairsSep *CStr, flags int) (int, error) {
	var ptrpm **C.AVDictionary
	var tmppm *C.AVDictionary
	var oldTmppm *C.AVDictionary
	if pm != nil {
		innerpm := *pm
		if innerpm != nil {
			tmppm = innerpm.ptr
			oldTmppm = tmppm
		}
		ptrpm = &tmppm
	}
	var tmpstr *C.char
	if str != nil {
		tmpstr = str.ptr
	}
	var tmpkeyValSep *C.char
	if keyValSep != nil {
		tmpkeyValSep = keyValSep.ptr
	}
	var tmppairsSep *C.char
	if pairsSep != nil {
		tmppairsSep = pairsSep.ptr
	}
	ret := C.av_dict_parse_string(ptrpm, tmpstr, tmpkeyValSep, tmppairsSep, C.int(flags))
	if tmppm != oldTmppm && pm != nil {
		if tmppm != nil {
			*pm = &AVDictionary{ptr: tmppm}
		} else {
			*pm = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function av_dict_copy ---

// AVDictCopy wraps av_dict_copy.
/*
  Copy entries from one AVDictionary struct into another.

  @note Metadata is read using the ::AV_DICT_IGNORE_SUFFIX flag

  @param dst   Pointer to a pointer to a AVDictionary struct to copy into. If *dst is NULL,
               this function will allocate a struct for you and put it in *dst
  @param src   Pointer to the source AVDictionary struct to copy items from.
  @param flags Flags to use when setting entries in *dst

  @return 0 on success, negative AVERROR code on failure. If dst was allocated
            by this function, callers should free the associated memory.
*/
func AVDictCopy(dst **AVDictionary, src *AVDictionary, flags int) (int, error) {
	var ptrdst **C.AVDictionary
	var tmpdst *C.AVDictionary
	var oldTmpdst *C.AVDictionary
	if dst != nil {
		innerdst := *dst
		if innerdst != nil {
			tmpdst = innerdst.ptr
			oldTmpdst = tmpdst
		}
		ptrdst = &tmpdst
	}
	var tmpsrc *C.AVDictionary
	if src != nil {
		tmpsrc = src.ptr
	}
	ret := C.av_dict_copy(ptrdst, tmpsrc, C.int(flags))
	if tmpdst != oldTmpdst && dst != nil {
		if tmpdst != nil {
			*dst = &AVDictionary{ptr: tmpdst}
		} else {
			*dst = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function av_dict_free ---

// AVDictFree wraps av_dict_free.
/*
  Free all the memory allocated for an AVDictionary struct
  and all keys and values.
*/
func AVDictFree(m **AVDictionary) {
	var ptrm **C.AVDictionary
	var tmpm *C.AVDictionary
	var oldTmpm *C.AVDictionary
	if m != nil {
		innerm := *m
		if innerm != nil {
			tmpm = innerm.ptr
			oldTmpm = tmpm
		}
		ptrm = &tmpm
	}
	C.av_dict_free(ptrm)
	if tmpm != oldTmpm && m != nil {
		if tmpm != nil {
			*m = &AVDictionary{ptr: tmpm}
		} else {
			*m = nil
		}
	}
}

// --- Function av_dict_get_string ---

// av_dict_get_string skipped due to buffer

// --- Function av_display_rotation_get ---

// av_display_rotation_get skipped due to const array param matrix

// --- Function av_display_rotation_set ---

// av_display_rotation_set skipped due to const array param matrix

// --- Function av_display_matrix_flip ---

// av_display_matrix_flip skipped due to const array param matrix

// --- Function av_dovi_alloc ---

// AVDoviAlloc wraps av_dovi_alloc.
/*
  Allocate a AVDOVIDecoderConfigurationRecord structure and initialize its
  fields to default values.

  @return the newly allocated struct or NULL on failure
*/
func AVDoviAlloc(size *uint64) *AVDOVIDecoderConfigurationRecord {
	ret := C.av_dovi_alloc((*C.size_t)(unsafe.Pointer(size)))
	var retMapped *AVDOVIDecoderConfigurationRecord
	if ret != nil {
		retMapped = &AVDOVIDecoderConfigurationRecord{ptr: ret}
	}
	return retMapped
}

// --- Function av_dovi_get_header ---

// AVDoviGetHeader wraps av_dovi_get_header.
func AVDoviGetHeader(data *AVDOVIMetadata) *AVDOVIRpuDataHeader {
	var tmpdata *C.AVDOVIMetadata
	if data != nil {
		tmpdata = data.ptr
	}
	ret := C.av_dovi_get_header(tmpdata)
	var retMapped *AVDOVIRpuDataHeader
	if ret != nil {
		retMapped = &AVDOVIRpuDataHeader{ptr: ret}
	}
	return retMapped
}

// --- Function av_dovi_get_mapping ---

// AVDoviGetMapping wraps av_dovi_get_mapping.
func AVDoviGetMapping(data *AVDOVIMetadata) *AVDOVIDataMapping {
	var tmpdata *C.AVDOVIMetadata
	if data != nil {
		tmpdata = data.ptr
	}
	ret := C.av_dovi_get_mapping(tmpdata)
	var retMapped *AVDOVIDataMapping
	if ret != nil {
		retMapped = &AVDOVIDataMapping{ptr: ret}
	}
	return retMapped
}

// --- Function av_dovi_get_color ---

// AVDoviGetColor wraps av_dovi_get_color.
func AVDoviGetColor(data *AVDOVIMetadata) *AVDOVIColorMetadata {
	var tmpdata *C.AVDOVIMetadata
	if data != nil {
		tmpdata = data.ptr
	}
	ret := C.av_dovi_get_color(tmpdata)
	var retMapped *AVDOVIColorMetadata
	if ret != nil {
		retMapped = &AVDOVIColorMetadata{ptr: ret}
	}
	return retMapped
}

// --- Function av_dovi_get_ext ---

// AVDoviGetExt wraps av_dovi_get_ext.
func AVDoviGetExt(data *AVDOVIMetadata, index int) *AVDOVIDmData {
	var tmpdata *C.AVDOVIMetadata
	if data != nil {
		tmpdata = data.ptr
	}
	ret := C.av_dovi_get_ext(tmpdata, C.int(index))
	var retMapped *AVDOVIDmData
	if ret != nil {
		retMapped = &AVDOVIDmData{ptr: ret}
	}
	return retMapped
}

// --- Function av_dovi_find_level ---

// AVDoviFindLevel wraps av_dovi_find_level.
/*
  Find an extension block with a given level, or NULL. In the case of
  multiple extension blocks, only the first is returned.
*/
func AVDoviFindLevel(data *AVDOVIMetadata, level uint8) *AVDOVIDmData {
	var tmpdata *C.AVDOVIMetadata
	if data != nil {
		tmpdata = data.ptr
	}
	ret := C.av_dovi_find_level(tmpdata, C.uint8_t(level))
	var retMapped *AVDOVIDmData
	if ret != nil {
		retMapped = &AVDOVIDmData{ptr: ret}
	}
	return retMapped
}

// --- Function av_dovi_metadata_alloc ---

// AVDoviMetadataAlloc wraps av_dovi_metadata_alloc.
/*
  Allocate an AVDOVIMetadata structure and initialize its
  fields to default values.

  @param size If this parameter is non-NULL, the size in bytes of the
              allocated struct will be written here on success

  @return the newly allocated struct or NULL on failure
*/
func AVDoviMetadataAlloc(size *uint64) *AVDOVIMetadata {
	ret := C.av_dovi_metadata_alloc((*C.size_t)(unsafe.Pointer(size)))
	var retMapped *AVDOVIMetadata
	if ret != nil {
		retMapped = &AVDOVIMetadata{ptr: ret}
	}
	return retMapped
}

// --- Function av_downmix_info_update_side_data ---

// AVDownmixInfoUpdateSideData wraps av_downmix_info_update_side_data.
/*
  Get a frame's AV_FRAME_DATA_DOWNMIX_INFO side data for editing.

  If the side data is absent, it is created and added to the frame.

  @param frame the frame for which the side data is to be obtained or created

  @return the AVDownmixInfo structure to be edited by the caller, or NULL if
          the structure cannot be allocated.
*/
func AVDownmixInfoUpdateSideData(frame *AVFrame) *AVDownmixInfo {
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.av_downmix_info_update_side_data(tmpframe)
	var retMapped *AVDownmixInfo
	if ret != nil {
		retMapped = &AVDownmixInfo{ptr: ret}
	}
	return retMapped
}

// --- Function av_encryption_info_alloc ---

// AVEncryptionInfoAlloc wraps av_encryption_info_alloc.
/*
  Allocates an AVEncryptionInfo structure and sub-pointers to hold the given
  number of subsamples.  This will allocate pointers for the key ID, IV,
  and subsample entries, set the size members, and zero-initialize the rest.

  @param subsample_count The number of subsamples.
  @param key_id_size The number of bytes in the key ID, should be 16.
  @param iv_size The number of bytes in the IV, should be 16.

  @return The new AVEncryptionInfo structure, or NULL on error.
*/
func AVEncryptionInfoAlloc(subsampleCount uint32, keyIdSize uint32, ivSize uint32) *AVEncryptionInfo {
	ret := C.av_encryption_info_alloc(C.uint32_t(subsampleCount), C.uint32_t(keyIdSize), C.uint32_t(ivSize))
	var retMapped *AVEncryptionInfo
	if ret != nil {
		retMapped = &AVEncryptionInfo{ptr: ret}
	}
	return retMapped
}

// --- Function av_encryption_info_clone ---

// AVEncryptionInfoClone wraps av_encryption_info_clone.
/*
  Allocates an AVEncryptionInfo structure with a copy of the given data.
  @return The new AVEncryptionInfo structure, or NULL on error.
*/
func AVEncryptionInfoClone(info *AVEncryptionInfo) *AVEncryptionInfo {
	var tmpinfo *C.AVEncryptionInfo
	if info != nil {
		tmpinfo = info.ptr
	}
	ret := C.av_encryption_info_clone(tmpinfo)
	var retMapped *AVEncryptionInfo
	if ret != nil {
		retMapped = &AVEncryptionInfo{ptr: ret}
	}
	return retMapped
}

// --- Function av_encryption_info_free ---

// AVEncryptionInfoFree wraps av_encryption_info_free.
/*
  Frees the given encryption info object.  This MUST NOT be used to free the
  side-data data pointer, that should use normal side-data methods.
*/
func AVEncryptionInfoFree(info *AVEncryptionInfo) {
	var tmpinfo *C.AVEncryptionInfo
	if info != nil {
		tmpinfo = info.ptr
	}
	C.av_encryption_info_free(tmpinfo)
}

// --- Function av_encryption_info_get_side_data ---

// AVEncryptionInfoGetSideData wraps av_encryption_info_get_side_data.
/*
  Creates a copy of the AVEncryptionInfo that is contained in the given side
  data.  The resulting object should be passed to av_encryption_info_free()
  when done.

  @return The new AVEncryptionInfo structure, or NULL on error.
*/
func AVEncryptionInfoGetSideData(sideData unsafe.Pointer, sideDataSize uint64) *AVEncryptionInfo {
	ret := C.av_encryption_info_get_side_data((*C.uint8_t)(sideData), C.size_t(sideDataSize))
	var retMapped *AVEncryptionInfo
	if ret != nil {
		retMapped = &AVEncryptionInfo{ptr: ret}
	}
	return retMapped
}

// --- Function av_encryption_info_add_side_data ---

// AVEncryptionInfoAddSideData wraps av_encryption_info_add_side_data.
/*
  Allocates and initializes side data that holds a copy of the given encryption
  info.  The resulting pointer should be either freed using av_free or given
  to av_packet_add_side_data().

  @return The new side-data pointer, or NULL.
*/
func AVEncryptionInfoAddSideData(info *AVEncryptionInfo, sideDataSize *uint64) unsafe.Pointer {
	var tmpinfo *C.AVEncryptionInfo
	if info != nil {
		tmpinfo = info.ptr
	}
	ret := C.av_encryption_info_add_side_data(tmpinfo, (*C.size_t)(unsafe.Pointer(sideDataSize)))
	return unsafe.Pointer(ret)
}

// --- Function av_encryption_init_info_alloc ---

// AVEncryptionInitInfoAlloc wraps av_encryption_init_info_alloc.
/*
  Allocates an AVEncryptionInitInfo structure and sub-pointers to hold the
  given sizes.  This will allocate pointers and set all the fields.

  @return The new AVEncryptionInitInfo structure, or NULL on error.
*/
func AVEncryptionInitInfoAlloc(systemIdSize uint32, numKeyIds uint32, keyIdSize uint32, dataSize uint32) *AVEncryptionInitInfo {
	ret := C.av_encryption_init_info_alloc(C.uint32_t(systemIdSize), C.uint32_t(numKeyIds), C.uint32_t(keyIdSize), C.uint32_t(dataSize))
	var retMapped *AVEncryptionInitInfo
	if ret != nil {
		retMapped = &AVEncryptionInitInfo{ptr: ret}
	}
	return retMapped
}

// --- Function av_encryption_init_info_free ---

// AVEncryptionInitInfoFree wraps av_encryption_init_info_free.
/*
  Frees the given encryption init info object.  This MUST NOT be used to free
  the side-data data pointer, that should use normal side-data methods.
*/
func AVEncryptionInitInfoFree(info *AVEncryptionInitInfo) {
	var tmpinfo *C.AVEncryptionInitInfo
	if info != nil {
		tmpinfo = info.ptr
	}
	C.av_encryption_init_info_free(tmpinfo)
}

// --- Function av_encryption_init_info_get_side_data ---

// AVEncryptionInitInfoGetSideData wraps av_encryption_init_info_get_side_data.
/*
  Creates a copy of the AVEncryptionInitInfo that is contained in the given
  side data.  The resulting object should be passed to
  av_encryption_init_info_free() when done.

  @return The new AVEncryptionInitInfo structure, or NULL on error.
*/
func AVEncryptionInitInfoGetSideData(sideData unsafe.Pointer, sideDataSize uint64) *AVEncryptionInitInfo {
	ret := C.av_encryption_init_info_get_side_data((*C.uint8_t)(sideData), C.size_t(sideDataSize))
	var retMapped *AVEncryptionInitInfo
	if ret != nil {
		retMapped = &AVEncryptionInitInfo{ptr: ret}
	}
	return retMapped
}

// --- Function av_encryption_init_info_add_side_data ---

// AVEncryptionInitInfoAddSideData wraps av_encryption_init_info_add_side_data.
/*
  Allocates and initializes side data that holds a copy of the given encryption
  init info.  The resulting pointer should be either freed using av_free or
  given to av_packet_add_side_data().

  @return The new side-data pointer, or NULL.
*/
func AVEncryptionInitInfoAddSideData(info *AVEncryptionInitInfo, sideDataSize *uint64) unsafe.Pointer {
	var tmpinfo *C.AVEncryptionInitInfo
	if info != nil {
		tmpinfo = info.ptr
	}
	ret := C.av_encryption_init_info_add_side_data(tmpinfo, (*C.size_t)(unsafe.Pointer(sideDataSize)))
	return unsafe.Pointer(ret)
}

// --- Function av_strerror ---

// AVStrerror wraps av_strerror.
/*
  Put a description of the AVERROR code errnum in errbuf.
  In case of failure the global variable errno is set to indicate the
  error. Even in case of failure av_strerror() will print a generic
  error message indicating the errnum provided to errbuf.

  @param errnum      error code to describe
  @param errbuf      buffer to which description is written
  @param errbuf_size the size in bytes of errbuf
  @return 0 on success, a negative value if a description for errnum
  cannot be found
*/
func AVStrerror(errnum int, errbuf *CStr, errbufSize uint64) (int, error) {
	var tmperrbuf *C.char
	if errbuf != nil {
		tmperrbuf = errbuf.ptr
	}
	ret := C.av_strerror(C.int(errnum), tmperrbuf, C.size_t(errbufSize))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_make_error_string ---

// AVMakeErrorString wraps av_make_error_string.
/*
  Fill the provided buffer with a string containing an error string
  corresponding to the AVERROR code errnum.

  @param errbuf         a buffer
  @param errbuf_size    size in bytes of errbuf
  @param errnum         error code to describe
  @return the buffer in input, filled with the error description
  @see av_strerror()
*/
func AVMakeErrorString(errbuf *CStr, errbufSize uint64, errnum int) *CStr {
	var tmperrbuf *C.char
	if errbuf != nil {
		tmperrbuf = errbuf.ptr
	}
	ret := C.av_make_error_string(tmperrbuf, C.size_t(errbufSize), C.int(errnum))
	return wrapCStr(ret)
}

// --- Function av_expr_parse_and_eval ---

// av_expr_parse_and_eval skipped due to res (non-output primitive pointer)

// --- Function av_expr_parse ---

// av_expr_parse skipped due to constNames

// --- Function av_expr_eval ---

// av_expr_eval skipped due to constValues (non-output primitive pointer)

// --- Function av_expr_count_vars ---

// AVExprCountVars wraps av_expr_count_vars.
/*
  Track the presence of variables and their number of occurrences in a parsed expression

  @param e the AVExpr to track variables in
  @param counter a zero-initialized array where the count of each variable will be stored
  @param size size of array
  @return 0 on success, a negative value indicates that no expression or array was passed
  or size was zero
*/
func AVExprCountVars(e *AVExpr, counter *uint, size int) (int, error) {
	var tmpe *C.AVExpr
	if e != nil {
		tmpe = e.ptr
	}
	ret := C.av_expr_count_vars(tmpe, (*C.uint)(unsafe.Pointer(counter)), C.int(size))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_expr_count_func ---

// AVExprCountFunc wraps av_expr_count_func.
/*
  Track the presence of user provided functions and their number of occurrences
  in a parsed expression.

  @param e the AVExpr to track user provided functions in
  @param counter a zero-initialized array where the count of each function will be stored
                 if you passed 5 functions with 2 arguments to av_expr_parse()
                 then for arg=2 this will use up to 5 entries.
  @param size size of array
  @param arg number of arguments the counted functions have
  @return 0 on success, a negative value indicates that no expression or array was passed
  or size was zero
*/
func AVExprCountFunc(e *AVExpr, counter *uint, size int, arg int) (int, error) {
	var tmpe *C.AVExpr
	if e != nil {
		tmpe = e.ptr
	}
	ret := C.av_expr_count_func(tmpe, (*C.uint)(unsafe.Pointer(counter)), C.int(size), C.int(arg))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_expr_free ---

// AVExprFree wraps av_expr_free.
//
//	Free a parsed expression previously created with av_expr_parse().
func AVExprFree(e *AVExpr) {
	var tmpe *C.AVExpr
	if e != nil {
		tmpe = e.ptr
	}
	C.av_expr_free(tmpe)
}

// --- Function av_strtod ---

// av_strtod skipped due to tail

// --- Function av_executor_alloc ---

// AVExecutorAlloc wraps av_executor_alloc.
/*
  Alloc executor
  @param callbacks callback structure for executor
  @param thread_count worker thread number, 0 for run on caller's thread directly
  @return return the executor
*/
func AVExecutorAlloc(callbacks *AVTaskCallbacks, threadCount int) *AVExecutor {
	var tmpcallbacks *C.AVTaskCallbacks
	if callbacks != nil {
		tmpcallbacks = callbacks.ptr
	}
	ret := C.av_executor_alloc(tmpcallbacks, C.int(threadCount))
	var retMapped *AVExecutor
	if ret != nil {
		retMapped = &AVExecutor{ptr: ret}
	}
	return retMapped
}

// --- Function av_executor_free ---

// AVExecutorFree wraps av_executor_free.
/*
  Free executor
  @param e  pointer to executor
*/
func AVExecutorFree(e **AVExecutor) {
	var ptre **C.AVExecutor
	var tmpe *C.AVExecutor
	var oldTmpe *C.AVExecutor
	if e != nil {
		innere := *e
		if innere != nil {
			tmpe = innere.ptr
			oldTmpe = tmpe
		}
		ptre = &tmpe
	}
	C.av_executor_free(ptre)
	if tmpe != oldTmpe && e != nil {
		if tmpe != nil {
			*e = &AVExecutor{ptr: tmpe}
		} else {
			*e = nil
		}
	}
}

// --- Function av_executor_execute ---

// AVExecutorExecute wraps av_executor_execute.
/*
  Add task to executor
  @param e pointer to executor
  @param t pointer to task. If NULL, it will wakeup one work thread
*/
func AVExecutorExecute(e *AVExecutor, t *AVTask) {
	var tmpe *C.AVExecutor
	if e != nil {
		tmpe = e.ptr
	}
	var tmpt *C.AVTask
	if t != nil {
		tmpt = t.ptr
	}
	C.av_executor_execute(tmpe, tmpt)
}

// --- Function av_fifo_alloc2 ---

// AVFifoAlloc2 wraps av_fifo_alloc2.
/*
  Allocate and initialize an AVFifo with a given element size.

  @param elems     initial number of elements that can be stored in the FIFO
  @param elem_size Size in bytes of a single element. Further operations on
                   the returned FIFO will implicitly use this element size.
  @param flags a combination of AV_FIFO_FLAG_*

  @return newly-allocated AVFifo on success, a negative error code on failure
*/
func AVFifoAlloc2(elems uint64, elemSize uint64, flags uint) *AVFifo {
	ret := C.av_fifo_alloc2(C.size_t(elems), C.size_t(elemSize), C.uint(flags))
	var retMapped *AVFifo
	if ret != nil {
		retMapped = &AVFifo{ptr: ret}
	}
	return retMapped
}

// --- Function av_fifo_elem_size ---

// AVFifoElemSize wraps av_fifo_elem_size.
/*
  @return Element size for FIFO operations. This element size is set at
          FIFO allocation and remains constant during its lifetime
*/
func AVFifoElemSize(f *AVFifo) uint64 {
	var tmpf *C.AVFifo
	if f != nil {
		tmpf = f.ptr
	}
	ret := C.av_fifo_elem_size(tmpf)
	return uint64(ret)
}

// --- Function av_fifo_auto_grow_limit ---

// AVFifoAutoGrowLimit wraps av_fifo_auto_grow_limit.
/*
  Set the maximum size (in elements) to which the FIFO can be resized
  automatically. Has no effect unless AV_FIFO_FLAG_AUTO_GROW is used.
*/
func AVFifoAutoGrowLimit(f *AVFifo, maxElems uint64) {
	var tmpf *C.AVFifo
	if f != nil {
		tmpf = f.ptr
	}
	C.av_fifo_auto_grow_limit(tmpf, C.size_t(maxElems))
}

// --- Function av_fifo_can_read ---

// AVFifoCanRead wraps av_fifo_can_read.
//
//	@return number of elements available for reading from the given FIFO.
func AVFifoCanRead(f *AVFifo) uint64 {
	var tmpf *C.AVFifo
	if f != nil {
		tmpf = f.ptr
	}
	ret := C.av_fifo_can_read(tmpf)
	return uint64(ret)
}

// --- Function av_fifo_can_write ---

// AVFifoCanWrite wraps av_fifo_can_write.
/*
  @return Number of elements that can be written into the given FIFO without
          growing it.

          In other words, this number of elements or less is guaranteed to fit
          into the FIFO. More data may be written when the
          AV_FIFO_FLAG_AUTO_GROW flag was specified at FIFO creation, but this
          may involve memory allocation, which can fail.
*/
func AVFifoCanWrite(f *AVFifo) uint64 {
	var tmpf *C.AVFifo
	if f != nil {
		tmpf = f.ptr
	}
	ret := C.av_fifo_can_write(tmpf)
	return uint64(ret)
}

// --- Function av_fifo_grow2 ---

// AVFifoGrow2 wraps av_fifo_grow2.
/*
  Enlarge an AVFifo.

  On success, the FIFO will be large enough to hold exactly
  inc + av_fifo_can_read() + av_fifo_can_write()
  elements. In case of failure, the old FIFO is kept unchanged.

  @param f AVFifo to resize
  @param inc number of elements to allocate for, in addition to the current
             allocated size
  @return a non-negative number on success, a negative error code on failure
*/
func AVFifoGrow2(f *AVFifo, inc uint64) (int, error) {
	var tmpf *C.AVFifo
	if f != nil {
		tmpf = f.ptr
	}
	ret := C.av_fifo_grow2(tmpf, C.size_t(inc))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_fifo_write ---

// AVFifoWrite wraps av_fifo_write.
/*
  Write data into a FIFO.

  In case nb_elems > av_fifo_can_write(f) and the AV_FIFO_FLAG_AUTO_GROW flag
  was not specified at FIFO creation, nothing is written and an error
  is returned.

  Calling function is guaranteed to succeed if nb_elems <= av_fifo_can_write(f).

  @param f the FIFO buffer
  @param buf Data to be written. nb_elems * av_fifo_elem_size(f) bytes will be
             read from buf on success.
  @param nb_elems number of elements to write into FIFO

  @return a non-negative number on success, a negative error code on failure
*/
func AVFifoWrite(f *AVFifo, buf unsafe.Pointer, nbElems uint64) (int, error) {
	var tmpf *C.AVFifo
	if f != nil {
		tmpf = f.ptr
	}
	ret := C.av_fifo_write(tmpf, buf, C.size_t(nbElems))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_fifo_write_from_cb ---

// av_fifo_write_from_cb skipped due to readCb (callback by value)

// --- Function av_fifo_read ---

// AVFifoRead wraps av_fifo_read.
/*
  Read data from a FIFO.

  In case nb_elems > av_fifo_can_read(f), nothing is read and an error
  is returned.

  @param f the FIFO buffer
  @param buf Buffer to store the data. nb_elems * av_fifo_elem_size(f) bytes
             will be written into buf on success.
  @param nb_elems number of elements to read from FIFO

  @return a non-negative number on success, a negative error code on failure
*/
func AVFifoRead(f *AVFifo, buf unsafe.Pointer, nbElems uint64) (int, error) {
	var tmpf *C.AVFifo
	if f != nil {
		tmpf = f.ptr
	}
	ret := C.av_fifo_read(tmpf, buf, C.size_t(nbElems))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_fifo_read_to_cb ---

// av_fifo_read_to_cb skipped due to writeCb (callback by value)

// --- Function av_fifo_peek ---

// AVFifoPeek wraps av_fifo_peek.
/*
  Read data from a FIFO without modifying FIFO state.

  Returns an error if an attempt is made to peek to nonexistent elements
  (i.e. if offset + nb_elems is larger than av_fifo_can_read(f)).

  @param f the FIFO buffer
  @param buf Buffer to store the data. nb_elems * av_fifo_elem_size(f) bytes
             will be written into buf.
  @param nb_elems number of elements to read from FIFO
  @param offset number of initial elements to skip.

  @return a non-negative number on success, a negative error code on failure
*/
func AVFifoPeek(f *AVFifo, buf unsafe.Pointer, nbElems uint64, offset uint64) (int, error) {
	var tmpf *C.AVFifo
	if f != nil {
		tmpf = f.ptr
	}
	ret := C.av_fifo_peek(tmpf, buf, C.size_t(nbElems), C.size_t(offset))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_fifo_peek_to_cb ---

// av_fifo_peek_to_cb skipped due to writeCb (callback by value)

// --- Function av_fifo_drain2 ---

// AVFifoDrain2 wraps av_fifo_drain2.
/*
  Discard the specified amount of data from an AVFifo.
  @param size number of elements to discard, MUST NOT be larger than
              av_fifo_can_read(f)
*/
func AVFifoDrain2(f *AVFifo, size uint64) {
	var tmpf *C.AVFifo
	if f != nil {
		tmpf = f.ptr
	}
	C.av_fifo_drain2(tmpf, C.size_t(size))
}

// --- Function av_fifo_reset2 ---

// AVFifoReset2 wraps av_fifo_reset2.
/*
  Empty the AVFifo.
  @param f AVFifo to reset
*/
func AVFifoReset2(f *AVFifo) {
	var tmpf *C.AVFifo
	if f != nil {
		tmpf = f.ptr
	}
	C.av_fifo_reset2(tmpf)
}

// --- Function av_fifo_freep2 ---

// AVFifoFreep2 wraps av_fifo_freep2.
/*
  Free an AVFifo and reset pointer to NULL.
  @param f Pointer to an AVFifo to free. *f == NULL is allowed.
*/
func AVFifoFreep2(f **AVFifo) {
	var ptrf **C.AVFifo
	var tmpf *C.AVFifo
	var oldTmpf *C.AVFifo
	if f != nil {
		innerf := *f
		if innerf != nil {
			tmpf = innerf.ptr
			oldTmpf = tmpf
		}
		ptrf = &tmpf
	}
	C.av_fifo_freep2(ptrf)
	if tmpf != oldTmpf && f != nil {
		if tmpf != nil {
			*f = &AVFifo{ptr: tmpf}
		} else {
			*f = nil
		}
	}
}

// --- Function av_file_map ---

// av_file_map skipped due to bufptr

// --- Function av_file_unmap ---

// AVFileUnmap wraps av_file_unmap.
/*
  Unmap or free the buffer bufptr created by av_file_map().

  @param bufptr the buffer previously created with av_file_map()
  @param size size in bytes of bufptr, must be the same as returned
  by av_file_map()
*/
func AVFileUnmap(bufptr unsafe.Pointer, size uint64) {
	C.av_file_unmap((*C.uint8_t)(bufptr), C.size_t(size))
}

// --- Function av_film_grain_params_alloc ---

// AVFilmGrainParamsAlloc wraps av_film_grain_params_alloc.
/*
  Allocate an AVFilmGrainParams structure and set its fields to
  default values. The resulting struct can be freed using av_freep().
  If size is not NULL it will be set to the number of bytes allocated.

  @return An AVFilmGrainParams filled with default values or NULL
          on failure.
*/
func AVFilmGrainParamsAlloc(size *uint64) *AVFilmGrainParams {
	ret := C.av_film_grain_params_alloc((*C.size_t)(unsafe.Pointer(size)))
	var retMapped *AVFilmGrainParams
	if ret != nil {
		retMapped = &AVFilmGrainParams{ptr: ret}
	}
	return retMapped
}

// --- Function av_film_grain_params_create_side_data ---

// AVFilmGrainParamsCreateSideData wraps av_film_grain_params_create_side_data.
/*
  Allocate a complete AVFilmGrainParams and add it to the frame.

  @param frame The frame which side data is added to.

  @return The AVFilmGrainParams structure to be filled by caller.
*/
func AVFilmGrainParamsCreateSideData(frame *AVFrame) *AVFilmGrainParams {
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.av_film_grain_params_create_side_data(tmpframe)
	var retMapped *AVFilmGrainParams
	if ret != nil {
		retMapped = &AVFilmGrainParams{ptr: ret}
	}
	return retMapped
}

// --- Function av_film_grain_params_select ---

// AVFilmGrainParamsSelect wraps av_film_grain_params_select.
/*
  Select the most appropriate film grain parameters set for the frame,
  taking into account the frame's format, resolution and video signal
  characteristics.

  @note, for H.274, this may select a film grain parameter set with
  greater chroma resolution than the frame. Users should take care to
  correctly adjust the chroma grain frequency to the frame.
*/
func AVFilmGrainParamsSelect(frame *AVFrame) *AVFilmGrainParams {
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.av_film_grain_params_select(tmpframe)
	var retMapped *AVFilmGrainParams
	if ret != nil {
		retMapped = &AVFilmGrainParams{ptr: ret}
	}
	return retMapped
}

// --- Function av_frame_alloc ---

// AVFrameAlloc wraps av_frame_alloc.
/*
  Allocate an AVFrame and set its fields to default values.  The resulting
  struct must be freed using av_frame_free().

  @return An AVFrame filled with default values or NULL on failure.

  @note this only allocates the AVFrame itself, not the data buffers. Those
  must be allocated through other means, e.g. with av_frame_get_buffer() or
  manually.
*/
func AVFrameAlloc() *AVFrame {
	ret := C.av_frame_alloc()
	var retMapped *AVFrame
	if ret != nil {
		retMapped = &AVFrame{ptr: ret}
	}
	return retMapped
}

// --- Function av_frame_free ---

// AVFrameFree wraps av_frame_free.
/*
  Free the frame and any dynamically allocated objects in it,
  e.g. extended_data. If the frame is reference counted, it will be
  unreferenced first.

  @param frame frame to be freed. The pointer will be set to NULL.
*/
func AVFrameFree(frame **AVFrame) {
	var ptrframe **C.AVFrame
	var tmpframe *C.AVFrame
	var oldTmpframe *C.AVFrame
	if frame != nil {
		innerframe := *frame
		if innerframe != nil {
			tmpframe = innerframe.ptr
			oldTmpframe = tmpframe
		}
		ptrframe = &tmpframe
	}
	C.av_frame_free(ptrframe)
	if tmpframe != oldTmpframe && frame != nil {
		if tmpframe != nil {
			*frame = &AVFrame{ptr: tmpframe}
		} else {
			*frame = nil
		}
	}
}

// --- Function av_frame_ref ---

// AVFrameRef wraps av_frame_ref.
/*
  Set up a new reference to the data described by the source frame.

  Copy frame properties from src to dst and create a new reference for each
  AVBufferRef from src.

  If src is not reference counted, new buffers are allocated and the data is
  copied.

  @warning: dst MUST have been either unreferenced with av_frame_unref(dst),
            or newly allocated with av_frame_alloc() before calling this
            function, or undefined behavior will occur.

  @return 0 on success, a negative AVERROR on error
*/
func AVFrameRef(dst *AVFrame, src *AVFrame) (int, error) {
	var tmpdst *C.AVFrame
	if dst != nil {
		tmpdst = dst.ptr
	}
	var tmpsrc *C.AVFrame
	if src != nil {
		tmpsrc = src.ptr
	}
	ret := C.av_frame_ref(tmpdst, tmpsrc)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_frame_replace ---

// AVFrameReplace wraps av_frame_replace.
/*
  Ensure the destination frame refers to the same data described by the source
  frame, either by creating a new reference for each AVBufferRef from src if
  they differ from those in dst, by allocating new buffers and copying data if
  src is not reference counted, or by unrefencing it if src is empty.

  Frame properties on dst will be replaced by those from src.

  @return 0 on success, a negative AVERROR on error. On error, dst is
          unreferenced.
*/
func AVFrameReplace(dst *AVFrame, src *AVFrame) (int, error) {
	var tmpdst *C.AVFrame
	if dst != nil {
		tmpdst = dst.ptr
	}
	var tmpsrc *C.AVFrame
	if src != nil {
		tmpsrc = src.ptr
	}
	ret := C.av_frame_replace(tmpdst, tmpsrc)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_frame_clone ---

// AVFrameClone wraps av_frame_clone.
/*
  Create a new frame that references the same data as src.

  This is a shortcut for av_frame_alloc()+av_frame_ref().

  @return newly created AVFrame on success, NULL on error.
*/
func AVFrameClone(src *AVFrame) *AVFrame {
	var tmpsrc *C.AVFrame
	if src != nil {
		tmpsrc = src.ptr
	}
	ret := C.av_frame_clone(tmpsrc)
	var retMapped *AVFrame
	if ret != nil {
		retMapped = &AVFrame{ptr: ret}
	}
	return retMapped
}

// --- Function av_frame_unref ---

// AVFrameUnref wraps av_frame_unref.
//
//	Unreference all the buffers referenced by frame and reset the frame fields.
func AVFrameUnref(frame *AVFrame) {
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	C.av_frame_unref(tmpframe)
}

// --- Function av_frame_move_ref ---

// AVFrameMoveRef wraps av_frame_move_ref.
/*
  Move everything contained in src to dst and reset src.

  @warning: dst is not unreferenced, but directly overwritten without reading
            or deallocating its contents. Call av_frame_unref(dst) manually
            before calling this function to ensure that no memory is leaked.
*/
func AVFrameMoveRef(dst *AVFrame, src *AVFrame) {
	var tmpdst *C.AVFrame
	if dst != nil {
		tmpdst = dst.ptr
	}
	var tmpsrc *C.AVFrame
	if src != nil {
		tmpsrc = src.ptr
	}
	C.av_frame_move_ref(tmpdst, tmpsrc)
}

// --- Function av_frame_get_buffer ---

// AVFrameGetBuffer wraps av_frame_get_buffer.
/*
  Allocate new buffer(s) for audio or video data.

  The following fields must be set on frame before calling this function:
  - format (pixel format for video, sample format for audio)
  - width and height for video
  - nb_samples and ch_layout for audio

  This function will fill AVFrame.data and AVFrame.buf arrays and, if
  necessary, allocate and fill AVFrame.extended_data and AVFrame.extended_buf.
  For planar formats, one buffer will be allocated for each plane.

  @warning: if frame already has been allocated, calling this function will
            leak memory. In addition, undefined behavior can occur in certain
            cases.

  @param frame frame in which to store the new buffers.
  @param align Required buffer size and data pointer alignment. If equal to 0,
               alignment will be chosen automatically for the current CPU.
               It is highly recommended to pass 0 here unless you know what
               you are doing.

  @return 0 on success, a negative AVERROR on error.
*/
func AVFrameGetBuffer(frame *AVFrame, align int) (int, error) {
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.av_frame_get_buffer(tmpframe, C.int(align))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_frame_is_writable ---

// AVFrameIsWritable wraps av_frame_is_writable.
/*
  Check if the frame data is writable.

  @return A positive value if the frame data is writable (which is true if and
  only if each of the underlying buffers has only one reference, namely the one
  stored in this frame). Return 0 otherwise.

  If 1 is returned the answer is valid until av_buffer_ref() is called on any
  of the underlying AVBufferRefs (e.g. through av_frame_ref() or directly).

  @see av_frame_make_writable(), av_buffer_is_writable()
*/
func AVFrameIsWritable(frame *AVFrame) (int, error) {
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.av_frame_is_writable(tmpframe)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_frame_make_writable ---

// AVFrameMakeWritable wraps av_frame_make_writable.
/*
  Ensure that the frame data is writable, avoiding data copy if possible.

  Do nothing if the frame is writable, allocate new buffers and copy the data
  if it is not. Non-refcounted frames behave as non-writable, i.e. a copy
  is always made.

  @return 0 on success, a negative AVERROR on error.

  @see av_frame_is_writable(), av_buffer_is_writable(),
  av_buffer_make_writable()
*/
func AVFrameMakeWritable(frame *AVFrame) (int, error) {
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.av_frame_make_writable(tmpframe)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_frame_copy ---

// AVFrameCopy wraps av_frame_copy.
/*
  Copy the frame data from src to dst.

  This function does not allocate anything, dst must be already initialized and
  allocated with the same parameters as src.

  This function only copies the frame data (i.e. the contents of the data /
  extended data arrays), not any other properties.

  @return >= 0 on success, a negative AVERROR on error.
*/
func AVFrameCopy(dst *AVFrame, src *AVFrame) (int, error) {
	var tmpdst *C.AVFrame
	if dst != nil {
		tmpdst = dst.ptr
	}
	var tmpsrc *C.AVFrame
	if src != nil {
		tmpsrc = src.ptr
	}
	ret := C.av_frame_copy(tmpdst, tmpsrc)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_frame_copy_props ---

// AVFrameCopyProps wraps av_frame_copy_props.
/*
  Copy only "metadata" fields from src to dst.

  Metadata for the purpose of this function are those fields that do not affect
  the data layout in the buffers.  E.g. pts, sample rate (for audio) or sample
  aspect ratio (for video), but not width/height or channel layout.
  Side data is also copied.
*/
func AVFrameCopyProps(dst *AVFrame, src *AVFrame) (int, error) {
	var tmpdst *C.AVFrame
	if dst != nil {
		tmpdst = dst.ptr
	}
	var tmpsrc *C.AVFrame
	if src != nil {
		tmpsrc = src.ptr
	}
	ret := C.av_frame_copy_props(tmpdst, tmpsrc)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_frame_get_plane_buffer ---

// AVFrameGetPlaneBuffer wraps av_frame_get_plane_buffer.
/*
  Get the buffer reference a given data plane is stored in.

  @param frame the frame to get the plane's buffer from
  @param plane index of the data plane of interest in frame->extended_data.

  @return the buffer reference that contains the plane or NULL if the input
  frame is not valid.
*/
func AVFrameGetPlaneBuffer(frame *AVFrame, plane int) *AVBufferRef {
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.av_frame_get_plane_buffer(tmpframe, C.int(plane))
	var retMapped *AVBufferRef
	if ret != nil {
		retMapped = &AVBufferRef{ptr: ret}
	}
	return retMapped
}

// --- Function av_frame_new_side_data ---

// AVFrameNewSideData wraps av_frame_new_side_data.
/*
  Add a new side data to a frame.

  @param frame a frame to which the side data should be added
  @param type type of the added side data
  @param size size of the side data

  @return newly added side data on success, NULL on error
*/
func AVFrameNewSideData(frame *AVFrame, _type AVFrameSideDataType, size uint64) *AVFrameSideData {
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.av_frame_new_side_data(tmpframe, C.enum_AVFrameSideDataType(_type), C.size_t(size))
	var retMapped *AVFrameSideData
	if ret != nil {
		retMapped = &AVFrameSideData{ptr: ret}
	}
	return retMapped
}

// --- Function av_frame_new_side_data_from_buf ---

// AVFrameNewSideDataFromBuf wraps av_frame_new_side_data_from_buf.
/*
  Add a new side data to a frame from an existing AVBufferRef

  @param frame a frame to which the side data should be added
  @param type  the type of the added side data
  @param buf   an AVBufferRef to add as side data. The ownership of
               the reference is transferred to the frame.

  @return newly added side data on success, NULL on error. On failure
          the frame is unchanged and the AVBufferRef remains owned by
          the caller.
*/
func AVFrameNewSideDataFromBuf(frame *AVFrame, _type AVFrameSideDataType, buf *AVBufferRef) *AVFrameSideData {
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	var tmpbuf *C.AVBufferRef
	if buf != nil {
		tmpbuf = buf.ptr
	}
	ret := C.av_frame_new_side_data_from_buf(tmpframe, C.enum_AVFrameSideDataType(_type), tmpbuf)
	var retMapped *AVFrameSideData
	if ret != nil {
		retMapped = &AVFrameSideData{ptr: ret}
	}
	return retMapped
}

// --- Function av_frame_get_side_data ---

// AVFrameGetSideData wraps av_frame_get_side_data.
/*
  @return a pointer to the side data of a given type on success, NULL if there
  is no side data with such type in this frame.
*/
func AVFrameGetSideData(frame *AVFrame, _type AVFrameSideDataType) *AVFrameSideData {
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.av_frame_get_side_data(tmpframe, C.enum_AVFrameSideDataType(_type))
	var retMapped *AVFrameSideData
	if ret != nil {
		retMapped = &AVFrameSideData{ptr: ret}
	}
	return retMapped
}

// --- Function av_frame_remove_side_data ---

// AVFrameRemoveSideData wraps av_frame_remove_side_data.
//
//	Remove and free all side data instances of the given type.
func AVFrameRemoveSideData(frame *AVFrame, _type AVFrameSideDataType) {
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	C.av_frame_remove_side_data(tmpframe, C.enum_AVFrameSideDataType(_type))
}

// --- Function av_frame_apply_cropping ---

// AVFrameApplyCropping wraps av_frame_apply_cropping.
/*
  Crop the given video AVFrame according to its crop_left/crop_top/crop_right/
  crop_bottom fields. If cropping is successful, the function will adjust the
  data pointers and the width/height fields, and set the crop fields to 0.

  In all cases, the cropping boundaries will be rounded to the inherent
  alignment of the pixel format. In some cases, such as for opaque hwaccel
  formats, the left/top cropping is ignored. The crop fields are set to 0 even
  if the cropping was rounded or ignored.

  @param frame the frame which should be cropped
  @param flags Some combination of AV_FRAME_CROP_* flags, or 0.

  @return >= 0 on success, a negative AVERROR on error. If the cropping fields
  were invalid, AVERROR(ERANGE) is returned, and nothing is changed.
*/
func AVFrameApplyCropping(frame *AVFrame, flags int) (int, error) {
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.av_frame_apply_cropping(tmpframe, C.int(flags))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_frame_side_data_name ---

// AVFrameSideDataName wraps av_frame_side_data_name.
//
//	@return a string identifying the side data type
func AVFrameSideDataName(_type AVFrameSideDataType) *CStr {
	ret := C.av_frame_side_data_name(C.enum_AVFrameSideDataType(_type))
	return wrapCStr(ret)
}

// --- Function av_frame_side_data_desc ---

// AVFrameSideDataDesc wraps av_frame_side_data_desc.
/*
  @return side data descriptor corresponding to a given side data type, NULL
          when not available.
*/
func AVFrameSideDataDesc(_type AVFrameSideDataType) *AVSideDataDescriptor {
	ret := C.av_frame_side_data_desc(C.enum_AVFrameSideDataType(_type))
	var retMapped *AVSideDataDescriptor
	if ret != nil {
		retMapped = &AVSideDataDescriptor{ptr: ret}
	}
	return retMapped
}

// --- Function av_frame_side_data_free ---

// av_frame_side_data_free skipped due to sd

// --- Function av_frame_side_data_new ---

// av_frame_side_data_new skipped due to sd

// --- Function av_frame_side_data_add ---

// av_frame_side_data_add skipped due to sd

// --- Function av_frame_side_data_clone ---

// av_frame_side_data_clone skipped due to sd

// --- Function av_frame_side_data_get_c ---

// AVFrameSideDataGetC wraps av_frame_side_data_get_c.
/*
  Get a side data entry of a specific type from an array.

  @param sd    array of side data.
  @param nb_sd integer containing the number of entries in the array.
  @param type  type of side data to be queried

  @return a pointer to the side data of a given type on success, NULL if there
          is no side data with such type in this set.
*/
func AVFrameSideDataGetC(sd **AVFrameSideData, nbSd int, _type AVFrameSideDataType) *AVFrameSideData {
	var ptrsd **C.AVFrameSideData
	var tmpsd *C.AVFrameSideData
	var oldTmpsd *C.AVFrameSideData
	if sd != nil {
		innersd := *sd
		if innersd != nil {
			tmpsd = innersd.ptr
			oldTmpsd = tmpsd
		}
		ptrsd = &tmpsd
	}
	ret := C.av_frame_side_data_get_c(ptrsd, C.int(nbSd), C.enum_AVFrameSideDataType(_type))
	if tmpsd != oldTmpsd && sd != nil {
		if tmpsd != nil {
			*sd = &AVFrameSideData{ptr: tmpsd}
		} else {
			*sd = nil
		}
	}
	var retMapped *AVFrameSideData
	if ret != nil {
		retMapped = &AVFrameSideData{ptr: ret}
	}
	return retMapped
}

// --- Function av_frame_side_data_get ---

// AVFrameSideDataGet wraps av_frame_side_data_get.
/*
  Wrapper around av_frame_side_data_get_c() to workaround the limitation
  that for any type T the conversion from T * const * to const T * const *
  is not performed automatically in C.
  @see av_frame_side_data_get_c()
*/
func AVFrameSideDataGet(sd **AVFrameSideData, nbSd int, _type AVFrameSideDataType) *AVFrameSideData {
	var ptrsd **C.AVFrameSideData
	var tmpsd *C.AVFrameSideData
	var oldTmpsd *C.AVFrameSideData
	if sd != nil {
		innersd := *sd
		if innersd != nil {
			tmpsd = innersd.ptr
			oldTmpsd = tmpsd
		}
		ptrsd = &tmpsd
	}
	ret := C.av_frame_side_data_get(ptrsd, C.int(nbSd), C.enum_AVFrameSideDataType(_type))
	if tmpsd != oldTmpsd && sd != nil {
		if tmpsd != nil {
			*sd = &AVFrameSideData{ptr: tmpsd}
		} else {
			*sd = nil
		}
	}
	var retMapped *AVFrameSideData
	if ret != nil {
		retMapped = &AVFrameSideData{ptr: ret}
	}
	return retMapped
}

// --- Function av_frame_side_data_remove ---

// av_frame_side_data_remove skipped due to sd

// --- Function av_frame_side_data_remove_by_props ---

// av_frame_side_data_remove_by_props skipped due to sd

// --- Function av_hash_alloc ---

// AVHashAlloc wraps av_hash_alloc.
/*
  Allocate a hash context for the algorithm specified by name.

  @return  >= 0 for success, a negative error code for failure

  @note The context is not initialized after a call to this function; you must
  call av_hash_init() to do so.
*/
func AVHashAlloc(ctx **AVHashContext, name *CStr) (int, error) {
	var ptrctx **C.struct_AVHashContext
	var tmpctx *C.struct_AVHashContext
	var oldTmpctx *C.struct_AVHashContext
	if ctx != nil {
		innerctx := *ctx
		if innerctx != nil {
			tmpctx = innerctx.ptr
			oldTmpctx = tmpctx
		}
		ptrctx = &tmpctx
	}
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_hash_alloc(ptrctx, tmpname)
	if tmpctx != oldTmpctx && ctx != nil {
		if tmpctx != nil {
			*ctx = &AVHashContext{ptr: tmpctx}
		} else {
			*ctx = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function av_hash_names ---

// AVHashNames wraps av_hash_names.
/*
  Get the names of available hash algorithms.

  This function can be used to enumerate the algorithms.

  @param[in] i  Index of the hash algorithm, starting from 0
  @return       Pointer to a static string or `NULL` if `i` is out of range
*/
func AVHashNames(i int) *CStr {
	ret := C.av_hash_names(C.int(i))
	return wrapCStr(ret)
}

// --- Function av_hash_get_name ---

// AVHashGetName wraps av_hash_get_name.
//
//	Get the name of the algorithm corresponding to the given hash context.
func AVHashGetName(ctx *AVHashContext) *CStr {
	var tmpctx *C.struct_AVHashContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	ret := C.av_hash_get_name(tmpctx)
	return wrapCStr(ret)
}

// --- Function av_hash_get_size ---

// AVHashGetSize wraps av_hash_get_size.
/*
  Get the size of the resulting hash value in bytes.

  The maximum value this function will currently return is available as macro
  #AV_HASH_MAX_SIZE.

  @param[in]     ctx Hash context
  @return            Size of the hash value in bytes
*/
func AVHashGetSize(ctx *AVHashContext) (int, error) {
	var tmpctx *C.struct_AVHashContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	ret := C.av_hash_get_size(tmpctx)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_hash_init ---

// AVHashInit wraps av_hash_init.
/*
  Initialize or reset a hash context.

  @param[in,out] ctx Hash context
*/
func AVHashInit(ctx *AVHashContext) {
	var tmpctx *C.struct_AVHashContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	C.av_hash_init(tmpctx)
}

// --- Function av_hash_update ---

// AVHashUpdate wraps av_hash_update.
/*
  Update a hash context with additional data.

  @param[in,out] ctx Hash context
  @param[in]     src Data to be added to the hash context
  @param[in]     len Size of the additional data
*/
func AVHashUpdate(ctx *AVHashContext, src unsafe.Pointer, len uint64) {
	var tmpctx *C.struct_AVHashContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	C.av_hash_update(tmpctx, (*C.uint8_t)(src), C.size_t(len))
}

// --- Function av_hash_final ---

// AVHashFinal wraps av_hash_final.
/*
  Finalize a hash context and compute the actual hash value.

  The minimum size of `dst` buffer is given by av_hash_get_size() or
  #AV_HASH_MAX_SIZE. The use of the latter macro is discouraged.

  It is not safe to update or finalize a hash context again, if it has already
  been finalized.

  @param[in,out] ctx Hash context
  @param[out]    dst Where the final hash value will be stored

  @see av_hash_final_bin() provides an alternative API
*/
func AVHashFinal(ctx *AVHashContext, dst unsafe.Pointer) {
	var tmpctx *C.struct_AVHashContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	C.av_hash_final(tmpctx, (*C.uint8_t)(dst))
}

// --- Function av_hash_final_bin ---

// AVHashFinalBin wraps av_hash_final_bin.
/*
  Finalize a hash context and store the actual hash value in a buffer.

  It is not safe to update or finalize a hash context again, if it has already
  been finalized.

  If `size` is smaller than the hash size (given by av_hash_get_size()), the
  hash is truncated; if size is larger, the buffer is padded with 0.

  @param[in,out] ctx  Hash context
  @param[out]    dst  Where the final hash value will be stored
  @param[in]     size Number of bytes to write to `dst`
*/
func AVHashFinalBin(ctx *AVHashContext, dst unsafe.Pointer, size int) {
	var tmpctx *C.struct_AVHashContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	C.av_hash_final_bin(tmpctx, (*C.uint8_t)(dst), C.int(size))
}

// --- Function av_hash_final_hex ---

// AVHashFinalHex wraps av_hash_final_hex.
/*
  Finalize a hash context and store the hexadecimal representation of the
  actual hash value as a string.

  It is not safe to update or finalize a hash context again, if it has already
  been finalized.

  The string is always 0-terminated.

  If `size` is smaller than `2 * hash_size + 1`, where `hash_size` is the
  value returned by av_hash_get_size(), the string will be truncated.

  @param[in,out] ctx  Hash context
  @param[out]    dst  Where the string will be stored
  @param[in]     size Maximum number of bytes to write to `dst`
*/
func AVHashFinalHex(ctx *AVHashContext, dst unsafe.Pointer, size int) {
	var tmpctx *C.struct_AVHashContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	C.av_hash_final_hex(tmpctx, (*C.uint8_t)(dst), C.int(size))
}

// --- Function av_hash_final_b64 ---

// AVHashFinalB64 wraps av_hash_final_b64.
/*
  Finalize a hash context and store the Base64 representation of the
  actual hash value as a string.

  It is not safe to update or finalize a hash context again, if it has already
  been finalized.

  The string is always 0-terminated.

  If `size` is smaller than AV_BASE64_SIZE(hash_size), where `hash_size` is
  the value returned by av_hash_get_size(), the string will be truncated.

  @param[in,out] ctx  Hash context
  @param[out]    dst  Where the final hash value will be stored
  @param[in]     size Maximum number of bytes to write to `dst`
*/
func AVHashFinalB64(ctx *AVHashContext, dst unsafe.Pointer, size int) {
	var tmpctx *C.struct_AVHashContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	C.av_hash_final_b64(tmpctx, (*C.uint8_t)(dst), C.int(size))
}

// --- Function av_hash_freep ---

// AVHashFreep wraps av_hash_freep.
/*
  Free hash context and set hash context pointer to `NULL`.

  @param[in,out] ctx  Pointer to hash context
*/
func AVHashFreep(ctx **AVHashContext) {
	var ptrctx **C.struct_AVHashContext
	var tmpctx *C.struct_AVHashContext
	var oldTmpctx *C.struct_AVHashContext
	if ctx != nil {
		innerctx := *ctx
		if innerctx != nil {
			tmpctx = innerctx.ptr
			oldTmpctx = tmpctx
		}
		ptrctx = &tmpctx
	}
	C.av_hash_freep(ptrctx)
	if tmpctx != oldTmpctx && ctx != nil {
		if tmpctx != nil {
			*ctx = &AVHashContext{ptr: tmpctx}
		} else {
			*ctx = nil
		}
	}
}

// --- Function av_dynamic_hdr_plus_alloc ---

// AVDynamicHdrPlusAlloc wraps av_dynamic_hdr_plus_alloc.
/*
  Allocate an AVDynamicHDRPlus structure and set its fields to
  default values. The resulting struct can be freed using av_freep().

  @return An AVDynamicHDRPlus filled with default values or NULL
          on failure.
*/
func AVDynamicHdrPlusAlloc(size *uint64) *AVDynamicHDRPlus {
	ret := C.av_dynamic_hdr_plus_alloc((*C.size_t)(unsafe.Pointer(size)))
	var retMapped *AVDynamicHDRPlus
	if ret != nil {
		retMapped = &AVDynamicHDRPlus{ptr: ret}
	}
	return retMapped
}

// --- Function av_dynamic_hdr_plus_create_side_data ---

// AVDynamicHdrPlusCreateSideData wraps av_dynamic_hdr_plus_create_side_data.
/*
  Allocate a complete AVDynamicHDRPlus and add it to the frame.
  @param frame The frame which side data is added to.

  @return The AVDynamicHDRPlus structure to be filled by caller or NULL
          on failure.
*/
func AVDynamicHdrPlusCreateSideData(frame *AVFrame) *AVDynamicHDRPlus {
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.av_dynamic_hdr_plus_create_side_data(tmpframe)
	var retMapped *AVDynamicHDRPlus
	if ret != nil {
		retMapped = &AVDynamicHDRPlus{ptr: ret}
	}
	return retMapped
}

// --- Function av_dynamic_hdr_plus_from_t35 ---

// AVDynamicHdrPlusFromT35 wraps av_dynamic_hdr_plus_from_t35.
/*
  Parse the user data registered ITU-T T.35 to AVbuffer (AVDynamicHDRPlus).
  The T.35 buffer must begin with the application mode, skipping the
  country code, terminal provider codes, and application identifier.
  @param s A pointer containing the decoded AVDynamicHDRPlus structure.
  @param data The byte array containing the raw ITU-T T.35 data.
  @param size Size of the data array in bytes.

  @return >= 0 on success. Otherwise, returns the appropriate AVERROR.
*/
func AVDynamicHdrPlusFromT35(s *AVDynamicHDRPlus, data unsafe.Pointer, size uint64) (int, error) {
	var tmps *C.AVDynamicHDRPlus
	if s != nil {
		tmps = s.ptr
	}
	ret := C.av_dynamic_hdr_plus_from_t35(tmps, (*C.uint8_t)(data), C.size_t(size))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_dynamic_hdr_plus_to_t35 ---

// av_dynamic_hdr_plus_to_t35 skipped due to data

// --- Function av_dynamic_hdr_vivid_alloc ---

// AVDynamicHdrVividAlloc wraps av_dynamic_hdr_vivid_alloc.
/*
  Allocate an AVDynamicHDRVivid structure and set its fields to
  default values. The resulting struct can be freed using av_freep().

  @return An AVDynamicHDRVivid filled with default values or NULL
          on failure.
*/
func AVDynamicHdrVividAlloc(size *uint64) *AVDynamicHDRVivid {
	ret := C.av_dynamic_hdr_vivid_alloc((*C.size_t)(unsafe.Pointer(size)))
	var retMapped *AVDynamicHDRVivid
	if ret != nil {
		retMapped = &AVDynamicHDRVivid{ptr: ret}
	}
	return retMapped
}

// --- Function av_dynamic_hdr_vivid_create_side_data ---

// AVDynamicHdrVividCreateSideData wraps av_dynamic_hdr_vivid_create_side_data.
/*
  Allocate a complete AVDynamicHDRVivid and add it to the frame.
  @param frame The frame which side data is added to.

  @return The AVDynamicHDRVivid structure to be filled by caller or NULL
          on failure.
*/
func AVDynamicHdrVividCreateSideData(frame *AVFrame) *AVDynamicHDRVivid {
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.av_dynamic_hdr_vivid_create_side_data(tmpframe)
	var retMapped *AVDynamicHDRVivid
	if ret != nil {
		retMapped = &AVDynamicHDRVivid{ptr: ret}
	}
	return retMapped
}

// --- Function av_hmac_alloc ---

// AVHmacAlloc wraps av_hmac_alloc.
/*
  Allocate an AVHMAC context.
  @param type The hash function used for the HMAC.
*/
func AVHmacAlloc(_type AVHMACType) *AVHMAC {
	ret := C.av_hmac_alloc(C.enum_AVHMACType(_type))
	var retMapped *AVHMAC
	if ret != nil {
		retMapped = &AVHMAC{ptr: ret}
	}
	return retMapped
}

// --- Function av_hmac_free ---

// AVHmacFree wraps av_hmac_free.
/*
  Free an AVHMAC context.
  @param ctx The context to free, may be NULL
*/
func AVHmacFree(ctx *AVHMAC) {
	var tmpctx *C.AVHMAC
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	C.av_hmac_free(tmpctx)
}

// --- Function av_hmac_init ---

// AVHmacInit wraps av_hmac_init.
/*
  Initialize an AVHMAC context with an authentication key.
  @param ctx    The HMAC context
  @param key    The authentication key
  @param keylen The length of the key, in bytes
*/
func AVHmacInit(ctx *AVHMAC, key unsafe.Pointer, keylen uint) {
	var tmpctx *C.AVHMAC
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	C.av_hmac_init(tmpctx, (*C.uint8_t)(key), C.uint(keylen))
}

// --- Function av_hmac_update ---

// AVHmacUpdate wraps av_hmac_update.
/*
  Hash data with the HMAC.
  @param ctx  The HMAC context
  @param data The data to hash
  @param len  The length of the data, in bytes
*/
func AVHmacUpdate(ctx *AVHMAC, data unsafe.Pointer, len uint) {
	var tmpctx *C.AVHMAC
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	C.av_hmac_update(tmpctx, (*C.uint8_t)(data), C.uint(len))
}

// --- Function av_hmac_final ---

// AVHmacFinal wraps av_hmac_final.
/*
  Finish hashing and output the HMAC digest.
  @param ctx    The HMAC context
  @param out    The output buffer to write the digest into
  @param outlen The length of the out buffer, in bytes
  @return       The number of bytes written to out, or a negative error code.
*/
func AVHmacFinal(ctx *AVHMAC, out unsafe.Pointer, outlen uint) (int, error) {
	var tmpctx *C.AVHMAC
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	ret := C.av_hmac_final(tmpctx, (*C.uint8_t)(out), C.uint(outlen))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_hmac_calc ---

// AVHmacCalc wraps av_hmac_calc.
/*
  Hash an array of data with a key.
  @param ctx    The HMAC context
  @param data   The data to hash
  @param len    The length of the data, in bytes
  @param key    The authentication key
  @param keylen The length of the key, in bytes
  @param out    The output buffer to write the digest into
  @param outlen The length of the out buffer, in bytes
  @return       The number of bytes written to out, or a negative error code.
*/
func AVHmacCalc(ctx *AVHMAC, data unsafe.Pointer, len uint, key unsafe.Pointer, keylen uint, out unsafe.Pointer, outlen uint) (int, error) {
	var tmpctx *C.AVHMAC
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	ret := C.av_hmac_calc(tmpctx, (*C.uint8_t)(data), C.uint(len), (*C.uint8_t)(key), C.uint(keylen), (*C.uint8_t)(out), C.uint(outlen))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_hwdevice_find_type_by_name ---

// AVHWDeviceFindTypeByName wraps av_hwdevice_find_type_by_name.
/*
  Look up an AVHWDeviceType by name.

  @param name String name of the device type (case-insensitive).
  @return The type from enum AVHWDeviceType, or AV_HWDEVICE_TYPE_NONE if
          not found.
*/
func AVHWDeviceFindTypeByName(name *CStr) AVHWDeviceType {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_hwdevice_find_type_by_name(tmpname)
	return AVHWDeviceType(ret)
}

// --- Function av_hwdevice_get_type_name ---

// AVHWDeviceGetTypeName wraps av_hwdevice_get_type_name.
/*
  Get the string name of an AVHWDeviceType.

  @param type Type from enum AVHWDeviceType.
  @return Pointer to a static string containing the name, or NULL if the type
          is not valid.
*/
func AVHWDeviceGetTypeName(_type AVHWDeviceType) *CStr {
	ret := C.av_hwdevice_get_type_name(C.enum_AVHWDeviceType(_type))
	return wrapCStr(ret)
}

// --- Function av_hwdevice_iterate_types ---

// AVHWDeviceIterateTypes wraps av_hwdevice_iterate_types.
/*
  Iterate over supported device types.

  @param prev AV_HWDEVICE_TYPE_NONE initially, then the previous type
              returned by this function in subsequent iterations.
  @return The next usable device type from enum AVHWDeviceType, or
          AV_HWDEVICE_TYPE_NONE if there are no more.
*/
func AVHWDeviceIterateTypes(prev AVHWDeviceType) AVHWDeviceType {
	ret := C.av_hwdevice_iterate_types(C.enum_AVHWDeviceType(prev))
	return AVHWDeviceType(ret)
}

// --- Function av_hwdevice_ctx_alloc ---

// AVHWDeviceCtxAlloc wraps av_hwdevice_ctx_alloc.
/*
  Allocate an AVHWDeviceContext for a given hardware type.

  @param type the type of the hardware device to allocate.
  @return a reference to the newly created AVHWDeviceContext on success or NULL
          on failure.
*/
func AVHWDeviceCtxAlloc(_type AVHWDeviceType) *AVBufferRef {
	ret := C.av_hwdevice_ctx_alloc(C.enum_AVHWDeviceType(_type))
	var retMapped *AVBufferRef
	if ret != nil {
		retMapped = &AVBufferRef{ptr: ret}
	}
	return retMapped
}

// --- Function av_hwdevice_ctx_init ---

// AVHWDeviceCtxInit wraps av_hwdevice_ctx_init.
/*
  Finalize the device context before use. This function must be called after
  the context is filled with all the required information and before it is
  used in any way.

  @param ref a reference to the AVHWDeviceContext
  @return 0 on success, a negative AVERROR code on failure
*/
func AVHWDeviceCtxInit(ref *AVBufferRef) (int, error) {
	var tmpref *C.AVBufferRef
	if ref != nil {
		tmpref = ref.ptr
	}
	ret := C.av_hwdevice_ctx_init(tmpref)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_hwdevice_ctx_create ---

// AVHWDeviceCtxCreate wraps av_hwdevice_ctx_create.
/*
  Open a device of the specified type and create an AVHWDeviceContext for it.

  This is a convenience function intended to cover the simple cases. Callers
  who need to fine-tune device creation/management should open the device
  manually and then wrap it in an AVHWDeviceContext using
  av_hwdevice_ctx_alloc()/av_hwdevice_ctx_init().

  The returned context is already initialized and ready for use, the caller
  should not call av_hwdevice_ctx_init() on it. The user_opaque/free fields of
  the created AVHWDeviceContext are set by this function and should not be
  touched by the caller.

  @param device_ctx On success, a reference to the newly-created device context
                    will be written here. The reference is owned by the caller
                    and must be released with av_buffer_unref() when no longer
                    needed. On failure, NULL will be written to this pointer.
  @param type The type of the device to create.
  @param device A type-specific string identifying the device to open.
  @param opts A dictionary of additional (type-specific) options to use in
              opening the device. The dictionary remains owned by the caller.
  @param flags currently unused

  @return 0 on success, a negative AVERROR code on failure.
*/
func AVHWDeviceCtxCreate(deviceCtx **AVBufferRef, _type AVHWDeviceType, device *CStr, opts *AVDictionary, flags int) (int, error) {
	var ptrdeviceCtx **C.AVBufferRef
	var tmpdeviceCtx *C.AVBufferRef
	var oldTmpdeviceCtx *C.AVBufferRef
	if deviceCtx != nil {
		innerdeviceCtx := *deviceCtx
		if innerdeviceCtx != nil {
			tmpdeviceCtx = innerdeviceCtx.ptr
			oldTmpdeviceCtx = tmpdeviceCtx
		}
		ptrdeviceCtx = &tmpdeviceCtx
	}
	var tmpdevice *C.char
	if device != nil {
		tmpdevice = device.ptr
	}
	var tmpopts *C.AVDictionary
	if opts != nil {
		tmpopts = opts.ptr
	}
	ret := C.av_hwdevice_ctx_create(ptrdeviceCtx, C.enum_AVHWDeviceType(_type), tmpdevice, tmpopts, C.int(flags))
	if tmpdeviceCtx != oldTmpdeviceCtx && deviceCtx != nil {
		if tmpdeviceCtx != nil {
			*deviceCtx = &AVBufferRef{ptr: tmpdeviceCtx}
		} else {
			*deviceCtx = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function av_hwdevice_ctx_create_derived ---

// AVHWDeviceCtxCreateDerived wraps av_hwdevice_ctx_create_derived.
/*
  Create a new device of the specified type from an existing device.

  If the source device is a device of the target type or was originally
  derived from such a device (possibly through one or more intermediate
  devices of other types), then this will return a reference to the
  existing device of the same type as is requested.

  Otherwise, it will attempt to derive a new device from the given source
  device.  If direct derivation to the new type is not implemented, it will
  attempt the same derivation from each ancestor of the source device in
  turn looking for an implemented derivation method.

  @param dst_ctx On success, a reference to the newly-created
                 AVHWDeviceContext.
  @param type    The type of the new device to create.
  @param src_ctx A reference to an existing AVHWDeviceContext which will be
                 used to create the new device.
  @param flags   Currently unused; should be set to zero.
  @return        Zero on success, a negative AVERROR code on failure.
*/
func AVHWDeviceCtxCreateDerived(dstCtx **AVBufferRef, _type AVHWDeviceType, srcCtx *AVBufferRef, flags int) (int, error) {
	var ptrdstCtx **C.AVBufferRef
	var tmpdstCtx *C.AVBufferRef
	var oldTmpdstCtx *C.AVBufferRef
	if dstCtx != nil {
		innerdstCtx := *dstCtx
		if innerdstCtx != nil {
			tmpdstCtx = innerdstCtx.ptr
			oldTmpdstCtx = tmpdstCtx
		}
		ptrdstCtx = &tmpdstCtx
	}
	var tmpsrcCtx *C.AVBufferRef
	if srcCtx != nil {
		tmpsrcCtx = srcCtx.ptr
	}
	ret := C.av_hwdevice_ctx_create_derived(ptrdstCtx, C.enum_AVHWDeviceType(_type), tmpsrcCtx, C.int(flags))
	if tmpdstCtx != oldTmpdstCtx && dstCtx != nil {
		if tmpdstCtx != nil {
			*dstCtx = &AVBufferRef{ptr: tmpdstCtx}
		} else {
			*dstCtx = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function av_hwdevice_ctx_create_derived_opts ---

// AVHWDeviceCtxCreateDerivedOpts wraps av_hwdevice_ctx_create_derived_opts.
/*
  Create a new device of the specified type from an existing device.

  This function performs the same action as av_hwdevice_ctx_create_derived,
  however, it is able to set options for the new device to be derived.

  @param dst_ctx On success, a reference to the newly-created
                 AVHWDeviceContext.
  @param type    The type of the new device to create.
  @param src_ctx A reference to an existing AVHWDeviceContext which will be
                 used to create the new device.
  @param options Options for the new device to create, same format as in
                 av_hwdevice_ctx_create.
  @param flags   Currently unused; should be set to zero.
  @return        Zero on success, a negative AVERROR code on failure.
*/
func AVHWDeviceCtxCreateDerivedOpts(dstCtx **AVBufferRef, _type AVHWDeviceType, srcCtx *AVBufferRef, options *AVDictionary, flags int) (int, error) {
	var ptrdstCtx **C.AVBufferRef
	var tmpdstCtx *C.AVBufferRef
	var oldTmpdstCtx *C.AVBufferRef
	if dstCtx != nil {
		innerdstCtx := *dstCtx
		if innerdstCtx != nil {
			tmpdstCtx = innerdstCtx.ptr
			oldTmpdstCtx = tmpdstCtx
		}
		ptrdstCtx = &tmpdstCtx
	}
	var tmpsrcCtx *C.AVBufferRef
	if srcCtx != nil {
		tmpsrcCtx = srcCtx.ptr
	}
	var tmpoptions *C.AVDictionary
	if options != nil {
		tmpoptions = options.ptr
	}
	ret := C.av_hwdevice_ctx_create_derived_opts(ptrdstCtx, C.enum_AVHWDeviceType(_type), tmpsrcCtx, tmpoptions, C.int(flags))
	if tmpdstCtx != oldTmpdstCtx && dstCtx != nil {
		if tmpdstCtx != nil {
			*dstCtx = &AVBufferRef{ptr: tmpdstCtx}
		} else {
			*dstCtx = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function av_hwframe_ctx_alloc ---

// AVHWFrameCtxAlloc wraps av_hwframe_ctx_alloc.
/*
  Allocate an AVHWFramesContext tied to a given device context.

  @param device_ctx a reference to a AVHWDeviceContext. This function will make
                    a new reference for internal use, the one passed to the
                    function remains owned by the caller.
  @return a reference to the newly created AVHWFramesContext on success or NULL
          on failure.
*/
func AVHWFrameCtxAlloc(deviceCtx *AVBufferRef) *AVBufferRef {
	var tmpdeviceCtx *C.AVBufferRef
	if deviceCtx != nil {
		tmpdeviceCtx = deviceCtx.ptr
	}
	ret := C.av_hwframe_ctx_alloc(tmpdeviceCtx)
	var retMapped *AVBufferRef
	if ret != nil {
		retMapped = &AVBufferRef{ptr: ret}
	}
	return retMapped
}

// --- Function av_hwframe_ctx_init ---

// AVHWFrameCtxInit wraps av_hwframe_ctx_init.
/*
  Finalize the context before use. This function must be called after the
  context is filled with all the required information and before it is attached
  to any frames.

  @param ref a reference to the AVHWFramesContext
  @return 0 on success, a negative AVERROR code on failure
*/
func AVHWFrameCtxInit(ref *AVBufferRef) (int, error) {
	var tmpref *C.AVBufferRef
	if ref != nil {
		tmpref = ref.ptr
	}
	ret := C.av_hwframe_ctx_init(tmpref)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_hwframe_get_buffer ---

// AVHWFrameGetBuffer wraps av_hwframe_get_buffer.
/*
  Allocate a new frame attached to the given AVHWFramesContext.

  @param hwframe_ctx a reference to an AVHWFramesContext
  @param frame an empty (freshly allocated or unreffed) frame to be filled with
               newly allocated buffers.
  @param flags currently unused, should be set to zero
  @return 0 on success, a negative AVERROR code on failure
*/
func AVHWFrameGetBuffer(hwframeCtx *AVBufferRef, frame *AVFrame, flags int) (int, error) {
	var tmphwframeCtx *C.AVBufferRef
	if hwframeCtx != nil {
		tmphwframeCtx = hwframeCtx.ptr
	}
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.av_hwframe_get_buffer(tmphwframeCtx, tmpframe, C.int(flags))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_hwframe_transfer_data ---

// AVHWFrameTransferData wraps av_hwframe_transfer_data.
/*
  Copy data to or from a hw surface. At least one of dst/src must have an
  AVHWFramesContext attached.

  If src has an AVHWFramesContext attached, then the format of dst (if set)
  must use one of the formats returned by av_hwframe_transfer_get_formats(src,
  AV_HWFRAME_TRANSFER_DIRECTION_FROM).
  If dst has an AVHWFramesContext attached, then the format of src must use one
  of the formats returned by av_hwframe_transfer_get_formats(dst,
  AV_HWFRAME_TRANSFER_DIRECTION_TO)

  dst may be "clean" (i.e. with data/buf pointers unset), in which case the
  data buffers will be allocated by this function using av_frame_get_buffer().
  If dst->format is set, then this format will be used, otherwise (when
  dst->format is AV_PIX_FMT_NONE) the first acceptable format will be chosen.

  The two frames must have matching allocated dimensions (i.e. equal to
  AVHWFramesContext.width/height), since not all device types support
  transferring a sub-rectangle of the whole surface. The display dimensions
  (i.e. AVFrame.width/height) may be smaller than the allocated dimensions, but
  also have to be equal for both frames. When the display dimensions are
  smaller than the allocated dimensions, the content of the padding in the
  destination frame is unspecified.

  @param dst the destination frame. dst is not touched on failure.
  @param src the source frame.
  @param flags currently unused, should be set to zero
  @return 0 on success, a negative AVERROR error code on failure.
*/
func AVHWFrameTransferData(dst *AVFrame, src *AVFrame, flags int) (int, error) {
	var tmpdst *C.AVFrame
	if dst != nil {
		tmpdst = dst.ptr
	}
	var tmpsrc *C.AVFrame
	if src != nil {
		tmpsrc = src.ptr
	}
	ret := C.av_hwframe_transfer_data(tmpdst, tmpsrc, C.int(flags))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_hwframe_transfer_get_formats ---

// av_hwframe_transfer_get_formats skipped due to formats

// --- Function av_hwdevice_hwconfig_alloc ---

// AVHWDeviceHWConfigAlloc wraps av_hwdevice_hwconfig_alloc.
/*
  Allocate a HW-specific configuration structure for a given HW device.
  After use, the user must free all members as required by the specific
  hardware structure being used, then free the structure itself with
  av_free().

  @param device_ctx a reference to the associated AVHWDeviceContext.
  @return The newly created HW-specific configuration structure on
          success or NULL on failure.
*/
func AVHWDeviceHWConfigAlloc(deviceCtx *AVBufferRef) unsafe.Pointer {
	var tmpdeviceCtx *C.AVBufferRef
	if deviceCtx != nil {
		tmpdeviceCtx = deviceCtx.ptr
	}
	ret := C.av_hwdevice_hwconfig_alloc(tmpdeviceCtx)
	return ret
}

// --- Function av_hwdevice_get_hwframe_constraints ---

// AVHWDeviceGetHWFrameConstraints wraps av_hwdevice_get_hwframe_constraints.
/*
  Get the constraints on HW frames given a device and the HW-specific
  configuration to be used with that device.  If no HW-specific
  configuration is provided, returns the maximum possible capabilities
  of the device.

  @param ref a reference to the associated AVHWDeviceContext.
  @param hwconfig a filled HW-specific configuration structure, or NULL
         to return the maximum possible capabilities of the device.
  @return AVHWFramesConstraints structure describing the constraints
          on the device, or NULL if not available.
*/
func AVHWDeviceGetHWFrameConstraints(ref *AVBufferRef, hwconfig unsafe.Pointer) *AVHWFramesConstraints {
	var tmpref *C.AVBufferRef
	if ref != nil {
		tmpref = ref.ptr
	}
	ret := C.av_hwdevice_get_hwframe_constraints(tmpref, hwconfig)
	var retMapped *AVHWFramesConstraints
	if ret != nil {
		retMapped = &AVHWFramesConstraints{ptr: ret}
	}
	return retMapped
}

// --- Function av_hwframe_constraints_free ---

// AVHWFrameConstraintsFree wraps av_hwframe_constraints_free.
/*
  Free an AVHWFrameConstraints structure.

  @param constraints The (filled or unfilled) AVHWFrameConstraints structure.
*/
func AVHWFrameConstraintsFree(constraints **AVHWFramesConstraints) {
	var ptrconstraints **C.AVHWFramesConstraints
	var tmpconstraints *C.AVHWFramesConstraints
	var oldTmpconstraints *C.AVHWFramesConstraints
	if constraints != nil {
		innerconstraints := *constraints
		if innerconstraints != nil {
			tmpconstraints = innerconstraints.ptr
			oldTmpconstraints = tmpconstraints
		}
		ptrconstraints = &tmpconstraints
	}
	C.av_hwframe_constraints_free(ptrconstraints)
	if tmpconstraints != oldTmpconstraints && constraints != nil {
		if tmpconstraints != nil {
			*constraints = &AVHWFramesConstraints{ptr: tmpconstraints}
		} else {
			*constraints = nil
		}
	}
}

// --- Function av_hwframe_map ---

// AVHWFrameMap wraps av_hwframe_map.
/*
  Map a hardware frame.

  This has a number of different possible effects, depending on the format
  and origin of the src and dst frames.  On input, src should be a usable
  frame with valid buffers and dst should be blank (typically as just created
  by av_frame_alloc()).  src should have an associated hwframe context, and
  dst may optionally have a format and associated hwframe context.

  If src was created by mapping a frame from the hwframe context of dst,
  then this function undoes the mapping - dst is replaced by a reference to
  the frame that src was originally mapped from.

  If both src and dst have an associated hwframe context, then this function
  attempts to map the src frame from its hardware context to that of dst and
  then fill dst with appropriate data to be usable there.  This will only be
  possible if the hwframe contexts and associated devices are compatible -
  given compatible devices, av_hwframe_ctx_create_derived() can be used to
  create a hwframe context for dst in which mapping should be possible.

  If src has a hwframe context but dst does not, then the src frame is
  mapped to normal memory and should thereafter be usable as a normal frame.
  If the format is set on dst, then the mapping will attempt to create dst
  with that format and fail if it is not possible.  If format is unset (is
  AV_PIX_FMT_NONE) then dst will be mapped with whatever the most appropriate
  format to use is (probably the sw_format of the src hwframe context).

  A return value of AVERROR(ENOSYS) indicates that the mapping is not
  possible with the given arguments and hwframe setup, while other return
  values indicate that it failed somehow.

  On failure, the destination frame will be left blank, except for the
  hw_frames_ctx/format fields they may have been set by the caller - those will
  be preserved as they were.

  @param dst Destination frame, to contain the mapping.
  @param src Source frame, to be mapped.
  @param flags Some combination of AV_HWFRAME_MAP_* flags.
  @return Zero on success, negative AVERROR code on failure.
*/
func AVHWFrameMap(dst *AVFrame, src *AVFrame, flags int) (int, error) {
	var tmpdst *C.AVFrame
	if dst != nil {
		tmpdst = dst.ptr
	}
	var tmpsrc *C.AVFrame
	if src != nil {
		tmpsrc = src.ptr
	}
	ret := C.av_hwframe_map(tmpdst, tmpsrc, C.int(flags))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_hwframe_ctx_create_derived ---

// AVHWFrameCtxCreateDerived wraps av_hwframe_ctx_create_derived.
/*
  Create and initialise an AVHWFramesContext as a mapping of another existing
  AVHWFramesContext on a different device.

  av_hwframe_ctx_init() should not be called after this.

  @param derived_frame_ctx  On success, a reference to the newly created
                            AVHWFramesContext.
  @param format             The AVPixelFormat for the derived context.
  @param derived_device_ctx A reference to the device to create the new
                            AVHWFramesContext on.
  @param source_frame_ctx   A reference to an existing AVHWFramesContext
                            which will be mapped to the derived context.
  @param flags  Some combination of AV_HWFRAME_MAP_* flags, defining the
                mapping parameters to apply to frames which are allocated
                in the derived device.
  @return       Zero on success, negative AVERROR code on failure.
*/
func AVHWFrameCtxCreateDerived(derivedFrameCtx **AVBufferRef, format AVPixelFormat, derivedDeviceCtx *AVBufferRef, sourceFrameCtx *AVBufferRef, flags int) (int, error) {
	var ptrderivedFrameCtx **C.AVBufferRef
	var tmpderivedFrameCtx *C.AVBufferRef
	var oldTmpderivedFrameCtx *C.AVBufferRef
	if derivedFrameCtx != nil {
		innerderivedFrameCtx := *derivedFrameCtx
		if innerderivedFrameCtx != nil {
			tmpderivedFrameCtx = innerderivedFrameCtx.ptr
			oldTmpderivedFrameCtx = tmpderivedFrameCtx
		}
		ptrderivedFrameCtx = &tmpderivedFrameCtx
	}
	var tmpderivedDeviceCtx *C.AVBufferRef
	if derivedDeviceCtx != nil {
		tmpderivedDeviceCtx = derivedDeviceCtx.ptr
	}
	var tmpsourceFrameCtx *C.AVBufferRef
	if sourceFrameCtx != nil {
		tmpsourceFrameCtx = sourceFrameCtx.ptr
	}
	ret := C.av_hwframe_ctx_create_derived(ptrderivedFrameCtx, C.enum_AVPixelFormat(format), tmpderivedDeviceCtx, tmpsourceFrameCtx, C.int(flags))
	if tmpderivedFrameCtx != oldTmpderivedFrameCtx && derivedFrameCtx != nil {
		if tmpderivedFrameCtx != nil {
			*derivedFrameCtx = &AVBufferRef{ptr: tmpderivedFrameCtx}
		} else {
			*derivedFrameCtx = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function av_iamf_param_definition_get_class ---

// AVIamfParamDefinitionGetClass wraps av_iamf_param_definition_get_class.
func AVIamfParamDefinitionGetClass() *AVClass {
	ret := C.av_iamf_param_definition_get_class()
	var retMapped *AVClass
	if ret != nil {
		retMapped = &AVClass{ptr: ret}
	}
	return retMapped
}

// --- Function av_iamf_param_definition_alloc ---

// AVIamfParamDefinitionAlloc wraps av_iamf_param_definition_alloc.
/*
  Allocates memory for AVIAMFParamDefinition, plus an array of {@code nb_subblocks}
  amount of subblocks of the given type and initializes the variables. Can be
  freed with a normal av_free() call.

  @param size if non-NULL, the size in bytes of the resulting data array is written here.
*/
func AVIamfParamDefinitionAlloc(_type AVIAMFParamDefinitionType, nbSubblocks uint, size *uint64) *AVIAMFParamDefinition {
	ret := C.av_iamf_param_definition_alloc(C.enum_AVIAMFParamDefinitionType(_type), C.uint(nbSubblocks), (*C.size_t)(unsafe.Pointer(size)))
	var retMapped *AVIAMFParamDefinition
	if ret != nil {
		retMapped = &AVIAMFParamDefinition{ptr: ret}
	}
	return retMapped
}

// --- Function av_iamf_param_definition_get_subblock ---

// AVIamfParamDefinitionGetSubblock wraps av_iamf_param_definition_get_subblock.
/*
  Get the subblock at the specified {@code idx}. Must be between 0 and nb_subblocks - 1.

  The @ref AVIAMFParamDefinition.type "param definition type" defines
  the struct type of the returned pointer.
*/
func AVIamfParamDefinitionGetSubblock(par *AVIAMFParamDefinition, idx uint) unsafe.Pointer {
	var tmppar *C.AVIAMFParamDefinition
	if par != nil {
		tmppar = par.ptr
	}
	ret := C.av_iamf_param_definition_get_subblock(tmppar, C.uint(idx))
	return ret
}

// --- Function av_iamf_audio_element_get_class ---

// AVIamfAudioElementGetClass wraps av_iamf_audio_element_get_class.
func AVIamfAudioElementGetClass() *AVClass {
	ret := C.av_iamf_audio_element_get_class()
	var retMapped *AVClass
	if ret != nil {
		retMapped = &AVClass{ptr: ret}
	}
	return retMapped
}

// --- Function av_iamf_audio_element_alloc ---

// AVIamfAudioElementAlloc wraps av_iamf_audio_element_alloc.
/*
  Allocates a AVIAMFAudioElement, and initializes its fields with default values.
  No layers are allocated. Must be freed with av_iamf_audio_element_free().

  @see av_iamf_audio_element_add_layer()
*/
func AVIamfAudioElementAlloc() *AVIAMFAudioElement {
	ret := C.av_iamf_audio_element_alloc()
	var retMapped *AVIAMFAudioElement
	if ret != nil {
		retMapped = &AVIAMFAudioElement{ptr: ret}
	}
	return retMapped
}

// --- Function av_iamf_audio_element_add_layer ---

// AVIamfAudioElementAddLayer wraps av_iamf_audio_element_add_layer.
/*
  Allocate a layer and add it to a given AVIAMFAudioElement.
  It is freed by av_iamf_audio_element_free() alongside the rest of the parent
  AVIAMFAudioElement.

  @return a pointer to the allocated layer.
*/
func AVIamfAudioElementAddLayer(audioElement *AVIAMFAudioElement) *AVIAMFLayer {
	var tmpaudioElement *C.AVIAMFAudioElement
	if audioElement != nil {
		tmpaudioElement = audioElement.ptr
	}
	ret := C.av_iamf_audio_element_add_layer(tmpaudioElement)
	var retMapped *AVIAMFLayer
	if ret != nil {
		retMapped = &AVIAMFLayer{ptr: ret}
	}
	return retMapped
}

// --- Function av_iamf_audio_element_free ---

// AVIamfAudioElementFree wraps av_iamf_audio_element_free.
/*
  Free an AVIAMFAudioElement and all its contents.

  @param audio_element pointer to pointer to an allocated AVIAMFAudioElement.
                       upon return, *audio_element will be set to NULL.
*/
func AVIamfAudioElementFree(audioElement **AVIAMFAudioElement) {
	var ptraudioElement **C.AVIAMFAudioElement
	var tmpaudioElement *C.AVIAMFAudioElement
	var oldTmpaudioElement *C.AVIAMFAudioElement
	if audioElement != nil {
		inneraudioElement := *audioElement
		if inneraudioElement != nil {
			tmpaudioElement = inneraudioElement.ptr
			oldTmpaudioElement = tmpaudioElement
		}
		ptraudioElement = &tmpaudioElement
	}
	C.av_iamf_audio_element_free(ptraudioElement)
	if tmpaudioElement != oldTmpaudioElement && audioElement != nil {
		if tmpaudioElement != nil {
			*audioElement = &AVIAMFAudioElement{ptr: tmpaudioElement}
		} else {
			*audioElement = nil
		}
	}
}

// --- Function av_iamf_mix_presentation_get_class ---

// AVIamfMixPresentationGetClass wraps av_iamf_mix_presentation_get_class.
func AVIamfMixPresentationGetClass() *AVClass {
	ret := C.av_iamf_mix_presentation_get_class()
	var retMapped *AVClass
	if ret != nil {
		retMapped = &AVClass{ptr: ret}
	}
	return retMapped
}

// --- Function av_iamf_mix_presentation_alloc ---

// AVIamfMixPresentationAlloc wraps av_iamf_mix_presentation_alloc.
/*
  Allocates a AVIAMFMixPresentation, and initializes its fields with default
  values. No submixes are allocated.
  Must be freed with av_iamf_mix_presentation_free().

  @see av_iamf_mix_presentation_add_submix()
*/
func AVIamfMixPresentationAlloc() *AVIAMFMixPresentation {
	ret := C.av_iamf_mix_presentation_alloc()
	var retMapped *AVIAMFMixPresentation
	if ret != nil {
		retMapped = &AVIAMFMixPresentation{ptr: ret}
	}
	return retMapped
}

// --- Function av_iamf_mix_presentation_add_submix ---

// AVIamfMixPresentationAddSubmix wraps av_iamf_mix_presentation_add_submix.
/*
  Allocate a submix and add it to a given AVIAMFMixPresentation.
  It is freed by av_iamf_mix_presentation_free() alongside the rest of the
  parent AVIAMFMixPresentation.

  @return a pointer to the allocated submix.
*/
func AVIamfMixPresentationAddSubmix(mixPresentation *AVIAMFMixPresentation) *AVIAMFSubmix {
	var tmpmixPresentation *C.AVIAMFMixPresentation
	if mixPresentation != nil {
		tmpmixPresentation = mixPresentation.ptr
	}
	ret := C.av_iamf_mix_presentation_add_submix(tmpmixPresentation)
	var retMapped *AVIAMFSubmix
	if ret != nil {
		retMapped = &AVIAMFSubmix{ptr: ret}
	}
	return retMapped
}

// --- Function av_iamf_submix_add_element ---

// AVIamfSubmixAddElement wraps av_iamf_submix_add_element.
/*
  Allocate a submix element and add it to a given AVIAMFSubmix.
  It is freed by av_iamf_mix_presentation_free() alongside the rest of the
  parent AVIAMFSubmix.

  @return a pointer to the allocated submix.
*/
func AVIamfSubmixAddElement(submix *AVIAMFSubmix) *AVIAMFSubmixElement {
	var tmpsubmix *C.AVIAMFSubmix
	if submix != nil {
		tmpsubmix = submix.ptr
	}
	ret := C.av_iamf_submix_add_element(tmpsubmix)
	var retMapped *AVIAMFSubmixElement
	if ret != nil {
		retMapped = &AVIAMFSubmixElement{ptr: ret}
	}
	return retMapped
}

// --- Function av_iamf_submix_add_layout ---

// AVIamfSubmixAddLayout wraps av_iamf_submix_add_layout.
/*
  Allocate a submix layout and add it to a given AVIAMFSubmix.
  It is freed by av_iamf_mix_presentation_free() alongside the rest of the
  parent AVIAMFSubmix.

  @return a pointer to the allocated submix.
*/
func AVIamfSubmixAddLayout(submix *AVIAMFSubmix) *AVIAMFSubmixLayout {
	var tmpsubmix *C.AVIAMFSubmix
	if submix != nil {
		tmpsubmix = submix.ptr
	}
	ret := C.av_iamf_submix_add_layout(tmpsubmix)
	var retMapped *AVIAMFSubmixLayout
	if ret != nil {
		retMapped = &AVIAMFSubmixLayout{ptr: ret}
	}
	return retMapped
}

// --- Function av_iamf_mix_presentation_free ---

// AVIamfMixPresentationFree wraps av_iamf_mix_presentation_free.
/*
  Free an AVIAMFMixPresentation and all its contents.

  @param mix_presentation pointer to pointer to an allocated AVIAMFMixPresentation.
                          upon return, *mix_presentation will be set to NULL.
*/
func AVIamfMixPresentationFree(mixPresentation **AVIAMFMixPresentation) {
	var ptrmixPresentation **C.AVIAMFMixPresentation
	var tmpmixPresentation *C.AVIAMFMixPresentation
	var oldTmpmixPresentation *C.AVIAMFMixPresentation
	if mixPresentation != nil {
		innermixPresentation := *mixPresentation
		if innermixPresentation != nil {
			tmpmixPresentation = innermixPresentation.ptr
			oldTmpmixPresentation = tmpmixPresentation
		}
		ptrmixPresentation = &tmpmixPresentation
	}
	C.av_iamf_mix_presentation_free(ptrmixPresentation)
	if tmpmixPresentation != oldTmpmixPresentation && mixPresentation != nil {
		if tmpmixPresentation != nil {
			*mixPresentation = &AVIAMFMixPresentation{ptr: tmpmixPresentation}
		} else {
			*mixPresentation = nil
		}
	}
}

// --- Function av_image_fill_max_pixsteps ---

// av_image_fill_max_pixsteps skipped due to const array param maxPixsteps

// --- Function av_image_get_linesize ---

// AVImageGetLinesize wraps av_image_get_linesize.
/*
  Compute the size of an image line with format pix_fmt and width
  width for the plane plane.

  @return the computed size in bytes
*/
func AVImageGetLinesize(pixFmt AVPixelFormat, width int, plane int) (int, error) {
	ret := C.av_image_get_linesize(C.enum_AVPixelFormat(pixFmt), C.int(width), C.int(plane))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_image_fill_linesizes ---

// av_image_fill_linesizes skipped due to const array param linesizes

// --- Function av_image_fill_plane_sizes ---

// av_image_fill_plane_sizes skipped due to const array param size

// --- Function av_image_fill_pointers ---

// av_image_fill_pointers skipped due to const array param data

// --- Function av_image_alloc ---

// av_image_alloc skipped due to const array param pointers

// --- Function av_image_copy_plane ---

// AVImageCopyPlane wraps av_image_copy_plane.
/*
  Copy image plane from src to dst.
  That is, copy "height" number of lines of "bytewidth" bytes each.
  The first byte of each successive line is separated by *_linesize
  bytes.

  bytewidth must be contained by both absolute values of dst_linesize
  and src_linesize, otherwise the function behavior is undefined.

  @param dst          destination plane to copy to
  @param dst_linesize linesize for the image plane in dst
  @param src          source plane to copy from
  @param src_linesize linesize for the image plane in src
  @param height       height (number of lines) of the plane
*/
func AVImageCopyPlane(dst unsafe.Pointer, dstLinesize int, src unsafe.Pointer, srcLinesize int, bytewidth int, height int) {
	C.av_image_copy_plane((*C.uint8_t)(dst), C.int(dstLinesize), (*C.uint8_t)(src), C.int(srcLinesize), C.int(bytewidth), C.int(height))
}

// --- Function av_image_copy_plane_uc_from ---

// AVImageCopyPlaneUcFrom wraps av_image_copy_plane_uc_from.
/*
  Copy image data located in uncacheable (e.g. GPU mapped) memory. Where
  available, this function will use special functionality for reading from such
  memory, which may result in greatly improved performance compared to plain
  av_image_copy_plane().

  bytewidth must be contained by both absolute values of dst_linesize
  and src_linesize, otherwise the function behavior is undefined.

  @note The linesize parameters have the type ptrdiff_t here, while they are
        int for av_image_copy_plane().
  @note On x86, the linesizes currently need to be aligned to the cacheline
        size (i.e. 64) to get improved performance.
*/
func AVImageCopyPlaneUcFrom(dst unsafe.Pointer, dstLinesize int64, src unsafe.Pointer, srcLinesize int64, bytewidth int64, height int) {
	C.av_image_copy_plane_uc_from((*C.uint8_t)(dst), C.ptrdiff_t(dstLinesize), (*C.uint8_t)(src), C.ptrdiff_t(srcLinesize), C.ptrdiff_t(bytewidth), C.int(height))
}

// --- Function av_image_copy ---

// av_image_copy skipped due to const array param dstData

// --- Function av_image_copy2 ---

// av_image_copy2 skipped due to const array param dstData

// --- Function av_image_copy_uc_from ---

// av_image_copy_uc_from skipped due to const array param dstData

// --- Function av_image_fill_arrays ---

// av_image_fill_arrays skipped due to const array param dstData

// --- Function av_image_get_buffer_size ---

// AVImageGetBufferSize wraps av_image_get_buffer_size.
/*
  Return the size in bytes of the amount of data required to store an
  image with the given parameters.

  @param pix_fmt  the pixel format of the image
  @param width    the width of the image in pixels
  @param height   the height of the image in pixels
  @param align    the assumed linesize alignment
  @return the buffer size in bytes, a negative error code in case of failure
*/
func AVImageGetBufferSize(pixFmt AVPixelFormat, width int, height int, align int) (int, error) {
	ret := C.av_image_get_buffer_size(C.enum_AVPixelFormat(pixFmt), C.int(width), C.int(height), C.int(align))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_image_copy_to_buffer ---

// av_image_copy_to_buffer skipped due to const array param srcData

// --- Function av_image_check_size ---

// AVImageCheckSize wraps av_image_check_size.
/*
  Check if the given dimension of an image is valid, meaning that all
  bytes of the image can be addressed with a signed int.

  @param w the width of the picture
  @param h the height of the picture
  @param log_offset the offset to sum to the log level for logging with log_ctx
  @param log_ctx the parent logging context, it may be NULL
  @return >= 0 if valid, a negative error code otherwise
*/
func AVImageCheckSize(w uint, h uint, logOffset int, logCtx unsafe.Pointer) (int, error) {
	ret := C.av_image_check_size(C.uint(w), C.uint(h), C.int(logOffset), logCtx)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_image_check_size2 ---

// AVImageCheckSize2 wraps av_image_check_size2.
/*
  Check if the given dimension of an image is valid, meaning that all
  bytes of a plane of an image with the specified pix_fmt can be addressed
  with a signed int.

  @param w the width of the picture
  @param h the height of the picture
  @param max_pixels the maximum number of pixels the user wants to accept
  @param pix_fmt the pixel format, can be AV_PIX_FMT_NONE if unknown.
  @param log_offset the offset to sum to the log level for logging with log_ctx
  @param log_ctx the parent logging context, it may be NULL
  @return >= 0 if valid, a negative error code otherwise
*/
func AVImageCheckSize2(w uint, h uint, maxPixels int64, pixFmt AVPixelFormat, logOffset int, logCtx unsafe.Pointer) (int, error) {
	ret := C.av_image_check_size2(C.uint(w), C.uint(h), C.int64_t(maxPixels), C.enum_AVPixelFormat(pixFmt), C.int(logOffset), logCtx)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_image_check_sar ---

// AVImageCheckSar wraps av_image_check_sar.
/*
  Check if the given sample aspect ratio of an image is valid.

  It is considered invalid if the denominator is 0 or if applying the ratio
  to the image size would make the smaller dimension less than 1. If the
  sar numerator is 0, it is considered unknown and will return as valid.

  @param w width of the image
  @param h height of the image
  @param sar sample aspect ratio of the image
  @return 0 if valid, a negative AVERROR code otherwise
*/
func AVImageCheckSar(w uint, h uint, sar *AVRational) (int, error) {
	ret := C.av_image_check_sar(C.uint(w), C.uint(h), sar.value)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_image_fill_black ---

// av_image_fill_black skipped due to const array param dstData

// --- Function av_image_fill_color ---

// av_image_fill_color skipped due to const array param dstData

// --- Function av_int2float ---

// AVInt2Float wraps av_int2float.
//
//	Reinterpret a 32-bit integer as a float.
func AVInt2Float(i uint32) float32 {
	ret := C.av_int2float(C.uint32_t(i))
	return float32(ret)
}

// --- Function av_float2int ---

// AVFloat2Int wraps av_float2int.
//
//	Reinterpret a float as a 32-bit integer.
func AVFloat2Int(f float32) uint32 {
	ret := C.av_float2int(C.float(f))
	return uint32(ret)
}

// --- Function av_int2double ---

// AVInt2Double wraps av_int2double.
//
//	Reinterpret a 64-bit integer as a double.
func AVInt2Double(i uint64) float64 {
	ret := C.av_int2double(C.uint64_t(i))
	return float64(ret)
}

// --- Function av_double2int ---

// AVDouble2Int wraps av_double2int.
//
//	Reinterpret a double as a 64-bit integer.
func AVDouble2Int(f float64) uint64 {
	ret := C.av_double2int(C.double(f))
	return uint64(ret)
}

// --- Function av_lfg_init ---

// AVLfgInit wraps av_lfg_init.
func AVLfgInit(c *AVLFG, seed uint) {
	var tmpc *C.AVLFG
	if c != nil {
		tmpc = c.ptr
	}
	C.av_lfg_init(tmpc, C.uint(seed))
}

// --- Function av_lfg_init_from_data ---

// AVLfgInitFromData wraps av_lfg_init_from_data.
/*
  Seed the state of the ALFG using binary data.

  @return 0 on success, negative value (AVERROR) on failure.
*/
func AVLfgInitFromData(c *AVLFG, data unsafe.Pointer, length uint) (int, error) {
	var tmpc *C.AVLFG
	if c != nil {
		tmpc = c.ptr
	}
	ret := C.av_lfg_init_from_data(tmpc, (*C.uint8_t)(data), C.uint(length))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_lfg_get ---

// AVLfgGet wraps av_lfg_get.
/*
  Get the next random unsigned 32-bit number using an ALFG.

  Please also consider a simple LCG like state= state*1664525+1013904223,
  it may be good enough and faster for your specific use case.
*/
func AVLfgGet(c *AVLFG) uint {
	var tmpc *C.AVLFG
	if c != nil {
		tmpc = c.ptr
	}
	ret := C.av_lfg_get(tmpc)
	return uint(ret)
}

// --- Function av_mlfg_get ---

// AVMlfgGet wraps av_mlfg_get.
/*
  Get the next random unsigned 32-bit number using a MLFG.

  Please also consider av_lfg_get() above, it is faster.
*/
func AVMlfgGet(c *AVLFG) uint {
	var tmpc *C.AVLFG
	if c != nil {
		tmpc = c.ptr
	}
	ret := C.av_mlfg_get(tmpc)
	return uint(ret)
}

// --- Function av_bmg_get ---

// av_bmg_get skipped due to const array param out

// --- Function av_log ---

// av_log skipped due to variadic arg.

// --- Function av_log_once ---

// av_log_once skipped due to variadic arg.

// --- Function av_vlog ---

// av_vlog skipped due to vl.

// --- Function av_log_get_level ---

// AVLogGetLevel wraps av_log_get_level.
/*
  Get the current log level

  @see lavu_log_constants

  @return Current log level
*/
func AVLogGetLevel() (int, error) {
	ret := C.av_log_get_level()
	return int(ret), WrapErr(int(ret))
}

// --- Function av_log_set_level ---

// AVLogSetLevel wraps av_log_set_level.
/*
  Set the log level

  @see lavu_log_constants

  @param level Logging level
*/
func AVLogSetLevel(level int) {
	C.av_log_set_level(C.int(level))
}

// --- Function av_log_set_callback ---

// av_log_set_callback skipped due to callback.

// --- Function av_log_default_callback ---

// av_log_default_callback skipped due to vl.

// --- Function av_default_item_name ---

// AVDefaultItemName wraps av_default_item_name.
/*
  Return the context name

  @param  ctx The AVClass context

  @return The AVClass class_name
*/
func AVDefaultItemName(ctx unsafe.Pointer) *CStr {
	ret := C.av_default_item_name(ctx)
	return wrapCStr(ret)
}

// --- Function av_default_get_category ---

// AVDefaultGetCategory wraps av_default_get_category.
func AVDefaultGetCategory(ptr unsafe.Pointer) AVClassCategory {
	ret := C.av_default_get_category(ptr)
	return AVClassCategory(ret)
}

// --- Function av_log_format_line ---

// av_log_format_line skipped due to vl.

// --- Function av_log_format_line2 ---

// av_log_format_line2 skipped due to vl.

// --- Function av_log_set_flags ---

// AVLogSetFlags wraps av_log_set_flags.
func AVLogSetFlags(arg int) {
	C.av_log_set_flags(C.int(arg))
}

// --- Function av_log_get_flags ---

// AVLogGetFlags wraps av_log_get_flags.
func AVLogGetFlags() (int, error) {
	ret := C.av_log_get_flags()
	return int(ret), WrapErr(int(ret))
}

// --- Function av_lzo1x_decode ---

// av_lzo1x_decode skipped due to inlen (non-output primitive pointer)

// --- Function av_mastering_display_metadata_alloc ---

// AVMasteringDisplayMetadataAlloc wraps av_mastering_display_metadata_alloc.
/*
  Allocate an AVMasteringDisplayMetadata structure and set its fields to
  default values. The resulting struct can be freed using av_freep().

  @return An AVMasteringDisplayMetadata filled with default values or NULL
          on failure.
*/
func AVMasteringDisplayMetadataAlloc() *AVMasteringDisplayMetadata {
	ret := C.av_mastering_display_metadata_alloc()
	var retMapped *AVMasteringDisplayMetadata
	if ret != nil {
		retMapped = &AVMasteringDisplayMetadata{ptr: ret}
	}
	return retMapped
}

// --- Function av_mastering_display_metadata_alloc_size ---

// AVMasteringDisplayMetadataAllocSize wraps av_mastering_display_metadata_alloc_size.
/*
  Allocate an AVMasteringDisplayMetadata structure and set its fields to
  default values. The resulting struct can be freed using av_freep().

  @return An AVMasteringDisplayMetadata filled with default values or NULL
          on failure.
*/
func AVMasteringDisplayMetadataAllocSize(size *uint64) *AVMasteringDisplayMetadata {
	ret := C.av_mastering_display_metadata_alloc_size((*C.size_t)(unsafe.Pointer(size)))
	var retMapped *AVMasteringDisplayMetadata
	if ret != nil {
		retMapped = &AVMasteringDisplayMetadata{ptr: ret}
	}
	return retMapped
}

// --- Function av_mastering_display_metadata_create_side_data ---

// AVMasteringDisplayMetadataCreateSideData wraps av_mastering_display_metadata_create_side_data.
/*
  Allocate a complete AVMasteringDisplayMetadata and add it to the frame.

  @param frame The frame which side data is added to.

  @return The AVMasteringDisplayMetadata structure to be filled by caller.
*/
func AVMasteringDisplayMetadataCreateSideData(frame *AVFrame) *AVMasteringDisplayMetadata {
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.av_mastering_display_metadata_create_side_data(tmpframe)
	var retMapped *AVMasteringDisplayMetadata
	if ret != nil {
		retMapped = &AVMasteringDisplayMetadata{ptr: ret}
	}
	return retMapped
}

// --- Function av_content_light_metadata_alloc ---

// AVContentLightMetadataAlloc wraps av_content_light_metadata_alloc.
/*
  Allocate an AVContentLightMetadata structure and set its fields to
  default values. The resulting struct can be freed using av_freep().

  @return An AVContentLightMetadata filled with default values or NULL
          on failure.
*/
func AVContentLightMetadataAlloc(size *uint64) *AVContentLightMetadata {
	ret := C.av_content_light_metadata_alloc((*C.size_t)(unsafe.Pointer(size)))
	var retMapped *AVContentLightMetadata
	if ret != nil {
		retMapped = &AVContentLightMetadata{ptr: ret}
	}
	return retMapped
}

// --- Function av_content_light_metadata_create_side_data ---

// AVContentLightMetadataCreateSideData wraps av_content_light_metadata_create_side_data.
/*
  Allocate a complete AVContentLightMetadata and add it to the frame.

  @param frame The frame which side data is added to.

  @return The AVContentLightMetadata structure to be filled by caller.
*/
func AVContentLightMetadataCreateSideData(frame *AVFrame) *AVContentLightMetadata {
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.av_content_light_metadata_create_side_data(tmpframe)
	var retMapped *AVContentLightMetadata
	if ret != nil {
		retMapped = &AVContentLightMetadata{ptr: ret}
	}
	return retMapped
}

// --- Function av_gcd ---

// AVGcd wraps av_gcd.
/*
  Compute the greatest common divisor of two integer operands.

  @param a Operand
  @param b Operand
  @return GCD of a and b up to sign; if a >= 0 and b >= 0, return value is >= 0;
  if a == 0 and b == 0, returns 0.
*/
func AVGcd(a int64, b int64) int64 {
	ret := C.av_gcd(C.int64_t(a), C.int64_t(b))
	return int64(ret)
}

// --- Function av_rescale ---

// AVRescale wraps av_rescale.
/*
  Rescale a 64-bit integer with rounding to nearest.

  The operation is mathematically equivalent to `a * b / c`, but writing that
  directly can overflow.

  This function is equivalent to av_rescale_rnd() with #AV_ROUND_NEAR_INF.

  @see av_rescale_rnd(), av_rescale_q(), av_rescale_q_rnd()
*/
func AVRescale(a int64, b int64, c int64) int64 {
	ret := C.av_rescale(C.int64_t(a), C.int64_t(b), C.int64_t(c))
	return int64(ret)
}

// --- Function av_rescale_rnd ---

// AVRescaleRnd wraps av_rescale_rnd.
/*
  Rescale a 64-bit integer with specified rounding.

  The operation is mathematically equivalent to `a * b / c`, but writing that
  directly can overflow, and does not support different rounding methods.
  If the result is not representable then INT64_MIN is returned.

  @see av_rescale(), av_rescale_q(), av_rescale_q_rnd()
*/
func AVRescaleRnd(a int64, b int64, c int64, rnd AVRounding) int64 {
	ret := C.av_rescale_rnd(C.int64_t(a), C.int64_t(b), C.int64_t(c), C.enum_AVRounding(rnd))
	return int64(ret)
}

// --- Function av_rescale_q ---

// AVRescaleQ wraps av_rescale_q.
/*
  Rescale a 64-bit integer by 2 rational numbers.

  The operation is mathematically equivalent to `a * bq / cq`.

  This function is equivalent to av_rescale_q_rnd() with #AV_ROUND_NEAR_INF.

  @see av_rescale(), av_rescale_rnd(), av_rescale_q_rnd()
*/
func AVRescaleQ(a int64, bq *AVRational, cq *AVRational) int64 {
	ret := C.av_rescale_q(C.int64_t(a), bq.value, cq.value)
	return int64(ret)
}

// --- Function av_rescale_q_rnd ---

// AVRescaleQRnd wraps av_rescale_q_rnd.
/*
  Rescale a 64-bit integer by 2 rational numbers with specified rounding.

  The operation is mathematically equivalent to `a * bq / cq`.

  @see av_rescale(), av_rescale_rnd(), av_rescale_q()
*/
func AVRescaleQRnd(a int64, bq *AVRational, cq *AVRational, rnd AVRounding) int64 {
	ret := C.av_rescale_q_rnd(C.int64_t(a), bq.value, cq.value, C.enum_AVRounding(rnd))
	return int64(ret)
}

// --- Function av_compare_ts ---

// AVCompareTs wraps av_compare_ts.
/*
  Compare two timestamps each in its own time base.

  @return One of the following values:
          - -1 if `ts_a` is before `ts_b`
          - 1 if `ts_a` is after `ts_b`
          - 0 if they represent the same position

  @warning
  The result of the function is undefined if one of the timestamps is outside
  the `int64_t` range when represented in the other's timebase.
*/
func AVCompareTs(tsA int64, tbA *AVRational, tsB int64, tbB *AVRational) (int, error) {
	ret := C.av_compare_ts(C.int64_t(tsA), tbA.value, C.int64_t(tsB), tbB.value)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_compare_mod ---

// AVCompareMod wraps av_compare_mod.
/*
  Compare the remainders of two integer operands divided by a common divisor.

  In other words, compare the least significant `log2(mod)` bits of integers
  `a` and `b`.

  @code{.c}
  av_compare_mod(0x11, 0x02, 0x10) < 0 // since 0x11 % 0x10  (0x1) < 0x02 % 0x10  (0x2)
  av_compare_mod(0x11, 0x02, 0x20) > 0 // since 0x11 % 0x20 (0x11) > 0x02 % 0x20 (0x02)
  @endcode

  @param a Operand
  @param b Operand
  @param mod Divisor; must be a power of 2
  @return
          - a negative value if `a % mod < b % mod`
          - a positive value if `a % mod > b % mod`
          - zero             if `a % mod == b % mod`
*/
func AVCompareMod(a uint64, b uint64, mod uint64) int64 {
	ret := C.av_compare_mod(C.uint64_t(a), C.uint64_t(b), C.uint64_t(mod))
	return int64(ret)
}

// --- Function av_rescale_delta ---

// av_rescale_delta skipped due to last (non-output primitive pointer)

// --- Function av_add_stable ---

// AVAddStable wraps av_add_stable.
/*
  Add a value to a timestamp.

  This function guarantees that when the same value is repeatedly added that
  no accumulation of rounding errors occurs.

  @param[in] ts     Input timestamp
  @param[in] ts_tb  Input timestamp time base
  @param[in] inc    Value to be added
  @param[in] inc_tb Time base of `inc`
*/
func AVAddStable(tsTb *AVRational, ts int64, incTb *AVRational, inc int64) int64 {
	ret := C.av_add_stable(tsTb.value, C.int64_t(ts), incTb.value, C.int64_t(inc))
	return int64(ret)
}

// --- Function av_bessel_i0 ---

// AVBesselI0 wraps av_bessel_i0.
//
//	0th order modified bessel function of the first kind.
func AVBesselI0(x float64) float64 {
	ret := C.av_bessel_i0(C.double(x))
	return float64(ret)
}

// --- Function av_malloc ---

// AVMalloc wraps av_malloc.
/*
  Allocate a memory block with alignment suitable for all memory accesses
  (including vectors if available on the CPU).

  @param size Size in bytes for the memory block to be allocated
  @return Pointer to the allocated block, or `NULL` if the block cannot
          be allocated
  @see av_mallocz()
*/
func AVMalloc(size uint64) unsafe.Pointer {
	ret := C.av_malloc(C.size_t(size))
	return ret
}

// --- Function av_mallocz ---

// AVMallocz wraps av_mallocz.
/*
  Allocate a memory block with alignment suitable for all memory accesses
  (including vectors if available on the CPU) and zero all the bytes of the
  block.

  @param size Size in bytes for the memory block to be allocated
  @return Pointer to the allocated block, or `NULL` if it cannot be allocated
  @see av_malloc()
*/
func AVMallocz(size uint64) unsafe.Pointer {
	ret := C.av_mallocz(C.size_t(size))
	return ret
}

// --- Function av_malloc_array ---

// AVMallocArray wraps av_malloc_array.
/*
  Allocate a memory block for an array with av_malloc().

  The allocated memory will have size `size * nmemb` bytes.

  @param nmemb Number of element
  @param size  Size of a single element
  @return Pointer to the allocated block, or `NULL` if the block cannot
          be allocated
  @see av_malloc()
*/
func AVMallocArray(nmemb uint64, size uint64) unsafe.Pointer {
	ret := C.av_malloc_array(C.size_t(nmemb), C.size_t(size))
	return ret
}

// --- Function av_calloc ---

// AVCalloc wraps av_calloc.
/*
  Allocate a memory block for an array with av_mallocz().

  The allocated memory will have size `size * nmemb` bytes.

  @param nmemb Number of elements
  @param size  Size of the single element
  @return Pointer to the allocated block, or `NULL` if the block cannot
          be allocated

  @see av_mallocz()
  @see av_malloc_array()
*/
func AVCalloc(nmemb uint64, size uint64) unsafe.Pointer {
	ret := C.av_calloc(C.size_t(nmemb), C.size_t(size))
	return ret
}

// --- Function av_realloc ---

// AVRealloc wraps av_realloc.
/*
  Allocate, reallocate, or free a block of memory.

  If `ptr` is `NULL` and `size` > 0, allocate a new block. Otherwise, expand or
  shrink that block of memory according to `size`.

  @param ptr  Pointer to a memory block already allocated with
              av_realloc() or `NULL`
  @param size Size in bytes of the memory block to be allocated or
              reallocated

  @return Pointer to a newly-reallocated block or `NULL` if the block
          cannot be reallocated

  @warning Unlike av_malloc(), the returned pointer is not guaranteed to be
           correctly aligned. The returned pointer must be freed after even
           if size is zero.
  @see av_fast_realloc()
  @see av_reallocp()
*/
func AVRealloc(ptr unsafe.Pointer, size uint64) unsafe.Pointer {
	ret := C.av_realloc(ptr, C.size_t(size))
	return ret
}

// --- Function av_reallocp ---

// AVReallocp wraps av_reallocp.
/*
  Allocate, reallocate, or free a block of memory through a pointer to a
  pointer.

  If `*ptr` is `NULL` and `size` > 0, allocate a new block. If `size` is
  zero, free the memory block pointed to by `*ptr`. Otherwise, expand or
  shrink that block of memory according to `size`.

  @param[in,out] ptr  Pointer to a pointer to a memory block already allocated
                      with av_realloc(), or a pointer to `NULL`. The pointer
                      is updated on success, or freed on failure.
  @param[in]     size Size in bytes for the memory block to be allocated or
                      reallocated

  @return Zero on success, an AVERROR error code on failure

  @warning Unlike av_malloc(), the allocated memory is not guaranteed to be
           correctly aligned.
*/
func AVReallocp(ptr unsafe.Pointer, size uint64) (int, error) {
	ret := C.av_reallocp(ptr, C.size_t(size))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_realloc_f ---

// AVReallocF wraps av_realloc_f.
/*
  Allocate, reallocate, or free a block of memory.

  This function does the same thing as av_realloc(), except:
  - It takes two size arguments and allocates `nelem * elsize` bytes,
    after checking the result of the multiplication for integer overflow.
  - It frees the input block in case of failure, thus avoiding the memory
    leak with the classic
    @code{.c}
    buf = realloc(buf);
    if (!buf)
        return -1;
    @endcode
    pattern.
*/
func AVReallocF(ptr unsafe.Pointer, nelem uint64, elsize uint64) unsafe.Pointer {
	ret := C.av_realloc_f(ptr, C.size_t(nelem), C.size_t(elsize))
	return ret
}

// --- Function av_realloc_array ---

// AVReallocArray wraps av_realloc_array.
/*
  Allocate, reallocate, or free an array.

  If `ptr` is `NULL` and `nmemb` > 0, allocate a new block.

  @param ptr   Pointer to a memory block already allocated with
               av_realloc() or `NULL`
  @param nmemb Number of elements in the array
  @param size  Size of the single element of the array

  @return Pointer to a newly-reallocated block or NULL if the block
          cannot be reallocated

  @warning Unlike av_malloc(), the allocated memory is not guaranteed to be
           correctly aligned. The returned pointer must be freed after even if
           nmemb is zero.
  @see av_reallocp_array()
*/
func AVReallocArray(ptr unsafe.Pointer, nmemb uint64, size uint64) unsafe.Pointer {
	ret := C.av_realloc_array(ptr, C.size_t(nmemb), C.size_t(size))
	return ret
}

// --- Function av_reallocp_array ---

// AVReallocpArray wraps av_reallocp_array.
/*
  Allocate, reallocate an array through a pointer to a pointer.

  If `*ptr` is `NULL` and `nmemb` > 0, allocate a new block.

  @param[in,out] ptr   Pointer to a pointer to a memory block already
                       allocated with av_realloc(), or a pointer to `NULL`.
                       The pointer is updated on success, or freed on failure.
  @param[in]     nmemb Number of elements
  @param[in]     size  Size of the single element

  @return Zero on success, an AVERROR error code on failure

  @warning Unlike av_malloc(), the allocated memory is not guaranteed to be
           correctly aligned. *ptr must be freed after even if nmemb is zero.
*/
func AVReallocpArray(ptr unsafe.Pointer, nmemb uint64, size uint64) (int, error) {
	ret := C.av_reallocp_array(ptr, C.size_t(nmemb), C.size_t(size))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_fast_realloc ---

// AVFastRealloc wraps av_fast_realloc.
/*
  Reallocate the given buffer if it is not large enough, otherwise do nothing.

  If the given buffer is `NULL`, then a new uninitialized buffer is allocated.

  If the given buffer is not large enough, and reallocation fails, `NULL` is
  returned and `*size` is set to 0, but the original buffer is not changed or
  freed.

  A typical use pattern follows:

  @code{.c}
  uint8_t *buf = ...;
  uint8_t *new_buf = av_fast_realloc(buf, &current_size, size_needed);
  if (!new_buf) {
      // Allocation failed; clean up original buffer
      av_freep(&buf);
      return AVERROR(ENOMEM);
  }
  @endcode

  @param[in,out] ptr      Already allocated buffer, or `NULL`
  @param[in,out] size     Pointer to the size of buffer `ptr`. `*size` is
                          updated to the new allocated size, in particular 0
                          in case of failure.
  @param[in]     min_size Desired minimal size of buffer `ptr`
  @return `ptr` if the buffer is large enough, a pointer to newly reallocated
          buffer if the buffer was not large enough, or `NULL` in case of
          error
  @see av_realloc()
  @see av_fast_malloc()
*/
func AVFastRealloc(ptr unsafe.Pointer, size *uint, minSize uint64) unsafe.Pointer {
	ret := C.av_fast_realloc(ptr, (*C.uint)(unsafe.Pointer(size)), C.size_t(minSize))
	return ret
}

// --- Function av_fast_malloc ---

// AVFastMalloc wraps av_fast_malloc.
/*
  Allocate a buffer, reusing the given one if large enough.

  Contrary to av_fast_realloc(), the current buffer contents might not be
  preserved and on error the old buffer is freed, thus no special handling to
  avoid memleaks is necessary.

  `*ptr` is allowed to be `NULL`, in which case allocation always happens if
  `size_needed` is greater than 0.

  @code{.c}
  uint8_t *buf = ...;
  av_fast_malloc(&buf, &current_size, size_needed);
  if (!buf) {
      // Allocation failed; buf already freed
      return AVERROR(ENOMEM);
  }
  @endcode

  @param[in,out] ptr      Pointer to pointer to an already allocated buffer.
                          `*ptr` will be overwritten with pointer to new
                          buffer on success or `NULL` on failure
  @param[in,out] size     Pointer to the size of buffer `*ptr`. `*size` is
                          updated to the new allocated size, in particular 0
                          in case of failure.
  @param[in]     min_size Desired minimal size of buffer `*ptr`
  @see av_realloc()
  @see av_fast_mallocz()
*/
func AVFastMalloc(ptr unsafe.Pointer, size *uint, minSize uint64) {
	C.av_fast_malloc(ptr, (*C.uint)(unsafe.Pointer(size)), C.size_t(minSize))
}

// --- Function av_fast_mallocz ---

// AVFastMallocz wraps av_fast_mallocz.
/*
  Allocate and clear a buffer, reusing the given one if large enough.

  Like av_fast_malloc(), but all newly allocated space is initially cleared.
  Reused buffer is not cleared.

  `*ptr` is allowed to be `NULL`, in which case allocation always happens if
  `size_needed` is greater than 0.

  @param[in,out] ptr      Pointer to pointer to an already allocated buffer.
                          `*ptr` will be overwritten with pointer to new
                          buffer on success or `NULL` on failure
  @param[in,out] size     Pointer to the size of buffer `*ptr`. `*size` is
                          updated to the new allocated size, in particular 0
                          in case of failure.
  @param[in]     min_size Desired minimal size of buffer `*ptr`
  @see av_fast_malloc()
*/
func AVFastMallocz(ptr unsafe.Pointer, size *uint, minSize uint64) {
	C.av_fast_mallocz(ptr, (*C.uint)(unsafe.Pointer(size)), C.size_t(minSize))
}

// --- Function av_free ---

// AVFree wraps av_free.
/*
  Free a memory block which has been allocated with a function of av_malloc()
  or av_realloc() family.

  @param ptr Pointer to the memory block which should be freed.

  @note `ptr = NULL` is explicitly allowed.
  @note It is recommended that you use av_freep() instead, to prevent leaving
        behind dangling pointers.
  @see av_freep()
*/
func AVFree(ptr unsafe.Pointer) {
	C.av_free(ptr)
}

// --- Function av_freep ---

// AVFreep wraps av_freep.
/*
  Free a memory block which has been allocated with a function of av_malloc()
  or av_realloc() family, and set the pointer pointing to it to `NULL`.

  @code{.c}
  uint8_t *buf = av_malloc(16);
  av_free(buf);
  buf now contains a dangling pointer to freed memory, and accidental
  dereference of buf will result in a use-after-free, which may be a
  security risk.

  uint8_t *buf = av_malloc(16);
  av_freep(&buf);
  buf is now NULL, and accidental dereference will only result in a
  NULL-pointer dereference.
  @endcode

  @param ptr Pointer to the pointer to the memory block which should be freed
  @note `*ptr = NULL` is safe and leads to no action.
  @see av_free()
*/
func AVFreep(ptr unsafe.Pointer) {
	C.av_freep(ptr)
}

// --- Function av_strdup ---

// AVStrdup wraps av_strdup.
/*
  Duplicate a string.

  @param s String to be duplicated
  @return Pointer to a newly-allocated string containing a
          copy of `s` or `NULL` if the string cannot be allocated
  @see av_strndup()
*/
func AVStrdup(s *CStr) *CStr {
	var tmps *C.char
	if s != nil {
		tmps = s.ptr
	}
	ret := C.av_strdup(tmps)
	return wrapCStr(ret)
}

// --- Function av_strndup ---

// AVStrndup wraps av_strndup.
/*
  Duplicate a substring of a string.

  @param s   String to be duplicated
  @param len Maximum length of the resulting string (not counting the
             terminating byte)
  @return Pointer to a newly-allocated string containing a
          substring of `s` or `NULL` if the string cannot be allocated
*/
func AVStrndup(s *CStr, len uint64) *CStr {
	var tmps *C.char
	if s != nil {
		tmps = s.ptr
	}
	ret := C.av_strndup(tmps, C.size_t(len))
	return wrapCStr(ret)
}

// --- Function av_memdup ---

// AVMemdup wraps av_memdup.
/*
  Duplicate a buffer with av_malloc().

  @param p    Buffer to be duplicated
  @param size Size in bytes of the buffer copied
  @return Pointer to a newly allocated buffer containing a
          copy of `p` or `NULL` if the buffer cannot be allocated
*/
func AVMemdup(p unsafe.Pointer, size uint64) unsafe.Pointer {
	ret := C.av_memdup(p, C.size_t(size))
	return ret
}

// --- Function av_memcpy_backptr ---

// AVMemcpyBackptr wraps av_memcpy_backptr.
/*
  Overlapping memcpy() implementation.

  @param dst  Destination buffer
  @param back Number of bytes back to start copying (i.e. the initial size of
              the overlapping window); must be > 0
  @param cnt  Number of bytes to copy; must be >= 0

  @note `cnt > back` is valid, this will copy the bytes we just copied,
        thus creating a repeating pattern with a period length of `back`.
*/
func AVMemcpyBackptr(dst unsafe.Pointer, back int, cnt int) {
	C.av_memcpy_backptr((*C.uint8_t)(dst), C.int(back), C.int(cnt))
}

// --- Function av_dynarray_add ---

// AVDynarrayAdd wraps av_dynarray_add.
/*
  Add the pointer to an element to a dynamic array.

  The array to grow is supposed to be an array of pointers to
  structures, and the element to add must be a pointer to an already
  allocated structure.

  The array is reallocated when its size reaches powers of 2.
  Therefore, the amortized cost of adding an element is constant.

  In case of success, the pointer to the array is updated in order to
  point to the new grown array, and the number pointed to by `nb_ptr`
  is incremented.
  In case of failure, the array is freed, `*tab_ptr` is set to `NULL` and
  `*nb_ptr` is set to 0.

  @param[in,out] tab_ptr Pointer to the array to grow
  @param[in,out] nb_ptr  Pointer to the number of elements in the array
  @param[in]     elem    Element to add
  @see av_dynarray_add_nofree(), av_dynarray2_add()
*/
func AVDynarrayAdd(tabPtr unsafe.Pointer, nbPtr *int, elem unsafe.Pointer) {
	C.av_dynarray_add(tabPtr, (*C.int)(unsafe.Pointer(nbPtr)), elem)
}

// --- Function av_dynarray_add_nofree ---

// AVDynarrayAddNofree wraps av_dynarray_add_nofree.
/*
  Add an element to a dynamic array.

  Function has the same functionality as av_dynarray_add(),
  but it doesn't free memory on fails. It returns error code
  instead and leave current buffer untouched.

  @return >=0 on success, negative otherwise
  @see av_dynarray_add(), av_dynarray2_add()
*/
func AVDynarrayAddNofree(tabPtr unsafe.Pointer, nbPtr *int, elem unsafe.Pointer) (int, error) {
	ret := C.av_dynarray_add_nofree(tabPtr, (*C.int)(unsafe.Pointer(nbPtr)), elem)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_dynarray2_add ---

// av_dynarray2_add skipped due to tabPtr

// --- Function av_size_mult ---

// av_size_mult skipped due to r (non-output primitive pointer)

// --- Function av_max_alloc ---

// AVMaxAlloc wraps av_max_alloc.
/*
  Set the maximum size that may be allocated in one block.

  The value specified with this function is effective for all libavutil's @ref
  lavu_mem_funcs "heap management functions."

  By default, the max value is defined as `INT_MAX`.

  @param max Value to be set as the new maximum size

  @warning Exercise extreme caution when using this function. Don't touch
           this if you do not understand the full consequence of doing so.
*/
func AVMaxAlloc(max uint64) {
	C.av_max_alloc(C.size_t(max))
}

// --- Function av_murmur3_alloc ---

// AVMurmur3Alloc wraps av_murmur3_alloc.
/*
  Allocate an AVMurMur3 hash context.

  @return Uninitialized hash context or `NULL` in case of error
*/
func AVMurmur3Alloc() *AVMurMur3 {
	ret := C.av_murmur3_alloc()
	var retMapped *AVMurMur3
	if ret != nil {
		retMapped = &AVMurMur3{ptr: ret}
	}
	return retMapped
}

// --- Function av_murmur3_init_seeded ---

// AVMurmur3InitSeeded wraps av_murmur3_init_seeded.
/*
  Initialize or reinitialize an AVMurMur3 hash context with a seed.

  @param[out] c    Hash context
  @param[in]  seed Random seed

  @see av_murmur3_init()
  @see @ref lavu_murmur3_seedinfo "Detailed description" on a discussion of
  seeds for MurmurHash3.
*/
func AVMurmur3InitSeeded(c *AVMurMur3, seed uint64) {
	var tmpc *C.struct_AVMurMur3
	if c != nil {
		tmpc = c.ptr
	}
	C.av_murmur3_init_seeded(tmpc, C.uint64_t(seed))
}

// --- Function av_murmur3_init ---

// AVMurmur3Init wraps av_murmur3_init.
/*
  Initialize or reinitialize an AVMurMur3 hash context.

  Equivalent to av_murmur3_init_seeded() with a built-in seed.

  @param[out] c    Hash context

  @see av_murmur3_init_seeded()
  @see @ref lavu_murmur3_seedinfo "Detailed description" on a discussion of
  seeds for MurmurHash3.
*/
func AVMurmur3Init(c *AVMurMur3) {
	var tmpc *C.struct_AVMurMur3
	if c != nil {
		tmpc = c.ptr
	}
	C.av_murmur3_init(tmpc)
}

// --- Function av_murmur3_update ---

// AVMurmur3Update wraps av_murmur3_update.
/*
  Update hash context with new data.

  @param[out] c    Hash context
  @param[in]  src  Input data to update hash with
  @param[in]  len  Number of bytes to read from `src`
*/
func AVMurmur3Update(c *AVMurMur3, src unsafe.Pointer, len uint64) {
	var tmpc *C.struct_AVMurMur3
	if c != nil {
		tmpc = c.ptr
	}
	C.av_murmur3_update(tmpc, (*C.uint8_t)(src), C.size_t(len))
}

// --- Function av_murmur3_final ---

// av_murmur3_final skipped due to const array param dst

// --- Function av_opt_set_defaults ---

// AVOptSetDefaults wraps av_opt_set_defaults.
/*
  Set the values of all AVOption fields to their default values.

  @param s an AVOption-enabled struct (its first member must be a pointer to AVClass)
*/
func AVOptSetDefaults(s unsafe.Pointer) {
	C.av_opt_set_defaults(s)
}

// --- Function av_opt_set_defaults2 ---

// AVOptSetDefaults2 wraps av_opt_set_defaults2.
/*
  Set the values of all AVOption fields to their default values. Only these
  AVOption fields for which (opt->flags & mask) == flags will have their
  default applied to s.

  @param s an AVOption-enabled struct (its first member must be a pointer to AVClass)
  @param mask combination of AV_OPT_FLAG_*
  @param flags combination of AV_OPT_FLAG_*
*/
func AVOptSetDefaults2(s unsafe.Pointer, mask int, flags int) {
	C.av_opt_set_defaults2(s, C.int(mask), C.int(flags))
}

// --- Function av_opt_free ---

// AVOptFree wraps av_opt_free.
//
//	Free all allocated objects in obj.
func AVOptFree(obj unsafe.Pointer) {
	C.av_opt_free(obj)
}

// --- Function av_opt_next ---

// AVOptNext wraps av_opt_next.
/*
  Iterate over all AVOptions belonging to obj.

  @param obj an AVOptions-enabled struct or a double pointer to an
             AVClass describing it.
  @param prev result of the previous call to av_opt_next() on this object
              or NULL
  @return next AVOption or NULL
*/
func AVOptNext(obj unsafe.Pointer, prev *AVOption) *AVOption {
	var tmpprev *C.AVOption
	if prev != nil {
		tmpprev = prev.ptr
	}
	ret := C.av_opt_next(obj, tmpprev)
	var retMapped *AVOption
	if ret != nil {
		retMapped = &AVOption{ptr: ret}
	}
	return retMapped
}

// --- Function av_opt_child_next ---

// AVOptChildNext wraps av_opt_child_next.
/*
  Iterate over AVOptions-enabled children of obj.

  @param prev result of a previous call to this function or NULL
  @return next AVOptions-enabled child or NULL
*/
func AVOptChildNext(obj unsafe.Pointer, prev unsafe.Pointer) unsafe.Pointer {
	ret := C.av_opt_child_next(obj, prev)
	return ret
}

// --- Function av_opt_child_class_iterate ---

// av_opt_child_class_iterate skipped due to iter

// --- Function av_opt_find ---

// AVOptFind wraps av_opt_find.
/*
  Look for an option in an object. Consider only options which
  have all the specified flags set.

  @param[in] obj A pointer to a struct whose first element is a
                 pointer to an AVClass.
                 Alternatively a double pointer to an AVClass, if
                 AV_OPT_SEARCH_FAKE_OBJ search flag is set.
  @param[in] name The name of the option to look for.
  @param[in] unit When searching for named constants, name of the unit
                  it belongs to.
  @param opt_flags Find only options with all the specified flags set (AV_OPT_FLAG).
  @param search_flags A combination of AV_OPT_SEARCH_*.

  @return A pointer to the option found, or NULL if no option
          was found.

  @note Options found with AV_OPT_SEARCH_CHILDREN flag may not be settable
  directly with av_opt_set(). Use special calls which take an options
  AVDictionary (e.g. avformat_open_input()) to set options found with this
  flag.
*/
func AVOptFind(obj unsafe.Pointer, name *CStr, unit *CStr, optFlags int, searchFlags int) *AVOption {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	var tmpunit *C.char
	if unit != nil {
		tmpunit = unit.ptr
	}
	ret := C.av_opt_find(obj, tmpname, tmpunit, C.int(optFlags), C.int(searchFlags))
	var retMapped *AVOption
	if ret != nil {
		retMapped = &AVOption{ptr: ret}
	}
	return retMapped
}

// --- Function av_opt_find2 ---

// av_opt_find2 skipped due to targetObj

// --- Function av_opt_show2 ---

// AVOptShow2 wraps av_opt_show2.
/*
  Show the obj options.

  @param req_flags requested flags for the options to show. Show only the
  options for which it is opt->flags & req_flags.
  @param rej_flags rejected flags for the options to show. Show only the
  options for which it is !(opt->flags & req_flags).
  @param av_log_obj log context to use for showing the options
*/
func AVOptShow2(obj unsafe.Pointer, avLogObj unsafe.Pointer, reqFlags int, rejFlags int) (int, error) {
	ret := C.av_opt_show2(obj, avLogObj, C.int(reqFlags), C.int(rejFlags))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_get_key_value ---

// av_opt_get_key_value skipped due to ropts

// --- Function av_set_options_string ---

// AVSetOptionsString wraps av_set_options_string.
/*
  Parse the key/value pairs list in opts. For each key/value pair
  found, stores the value in the field in ctx that is named like the
  key. ctx must be an AVClass context, storing is done using
  AVOptions.

  @param opts options string to parse, may be NULL
  @param key_val_sep a 0-terminated list of characters used to
  separate key from value
  @param pairs_sep a 0-terminated list of characters used to separate
  two pairs from each other
  @return the number of successfully set key/value pairs, or a negative
  value corresponding to an AVERROR code in case of error:
  AVERROR(EINVAL) if opts cannot be parsed,
  the error code issued by av_opt_set() if a key/value pair
  cannot be set
*/
func AVSetOptionsString(ctx unsafe.Pointer, opts *CStr, keyValSep *CStr, pairsSep *CStr) (int, error) {
	var tmpopts *C.char
	if opts != nil {
		tmpopts = opts.ptr
	}
	var tmpkeyValSep *C.char
	if keyValSep != nil {
		tmpkeyValSep = keyValSep.ptr
	}
	var tmppairsSep *C.char
	if pairsSep != nil {
		tmppairsSep = pairsSep.ptr
	}
	ret := C.av_set_options_string(ctx, tmpopts, tmpkeyValSep, tmppairsSep)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_set_from_string ---

// av_opt_set_from_string skipped due to shorthand

// --- Function av_opt_set_dict ---

// AVOptSetDict wraps av_opt_set_dict.
/*
  Set all the options from a given dictionary on an object.

  @param obj a struct whose first element is a pointer to AVClass
  @param options options to process. This dictionary will be freed and replaced
                 by a new one containing all options not found in obj.
                 Of course this new dictionary needs to be freed by caller
                 with av_dict_free().

  @return 0 on success, a negative AVERROR if some option was found in obj,
          but could not be set.

  @see av_dict_copy()
*/
func AVOptSetDict(obj unsafe.Pointer, options **AVDictionary) (int, error) {
	var ptroptions **C.AVDictionary
	var tmpoptions *C.AVDictionary
	var oldTmpoptions *C.AVDictionary
	if options != nil {
		inneroptions := *options
		if inneroptions != nil {
			tmpoptions = inneroptions.ptr
			oldTmpoptions = tmpoptions
		}
		ptroptions = &tmpoptions
	}
	ret := C.av_opt_set_dict(obj, ptroptions)
	if tmpoptions != oldTmpoptions && options != nil {
		if tmpoptions != nil {
			*options = &AVDictionary{ptr: tmpoptions}
		} else {
			*options = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_set_dict2 ---

// AVOptSetDict2 wraps av_opt_set_dict2.
/*
  Set all the options from a given dictionary on an object.

  @param obj a struct whose first element is a pointer to AVClass
  @param options options to process. This dictionary will be freed and replaced
                 by a new one containing all options not found in obj.
                 Of course this new dictionary needs to be freed by caller
                 with av_dict_free().
  @param search_flags A combination of AV_OPT_SEARCH_*.

  @return 0 on success, a negative AVERROR if some option was found in obj,
          but could not be set.

  @see av_dict_copy()
*/
func AVOptSetDict2(obj unsafe.Pointer, options **AVDictionary, searchFlags int) (int, error) {
	var ptroptions **C.AVDictionary
	var tmpoptions *C.AVDictionary
	var oldTmpoptions *C.AVDictionary
	if options != nil {
		inneroptions := *options
		if inneroptions != nil {
			tmpoptions = inneroptions.ptr
			oldTmpoptions = tmpoptions
		}
		ptroptions = &tmpoptions
	}
	ret := C.av_opt_set_dict2(obj, ptroptions, C.int(searchFlags))
	if tmpoptions != oldTmpoptions && options != nil {
		if tmpoptions != nil {
			*options = &AVDictionary{ptr: tmpoptions}
		} else {
			*options = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_copy ---

// AVOptCopy wraps av_opt_copy.
/*
  Copy options from src object into dest object.

  The underlying AVClass of both src and dest must coincide. The guarantee
  below does not apply if this is not fulfilled.

  Options that require memory allocation (e.g. string or binary) are malloc'ed in dest object.
  Original memory allocated for such options is freed unless both src and dest options points to the same memory.

  Even on error it is guaranteed that allocated options from src and dest
  no longer alias each other afterwards; in particular calling av_opt_free()
  on both src and dest is safe afterwards if dest has been memdup'ed from src.

  @param dest Object to copy from
  @param src  Object to copy into
  @return 0 on success, negative on error
*/
func AVOptCopy(dest unsafe.Pointer, src unsafe.Pointer) (int, error) {
	ret := C.av_opt_copy(dest, src)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_set ---

// AVOptSet wraps av_opt_set.
/*
  Those functions set the field of obj with the given name to value.

  @param[in] obj A struct whose first element is a pointer to an AVClass.
  @param[in] name the name of the field to set
  @param[in] val The value to set. In case of av_opt_set() if the field is not
  of a string type, then the given string is parsed.
  SI postfixes and some named scalars are supported.
  If the field is of a numeric type, it has to be a numeric or named
  scalar. Behavior with more than one scalar and +- infix operators
  is undefined.
  If the field is of a flags type, it has to be a sequence of numeric
  scalars or named flags separated by '+' or '-'. Prefixing a flag
  with '+' causes it to be set without affecting the other flags;
  similarly, '-' unsets a flag.
  If the field is of a dictionary type, it has to be a ':' separated list of
  key=value parameters. Values containing ':' special characters must be
  escaped.
  @param search_flags flags passed to av_opt_find2. I.e. if AV_OPT_SEARCH_CHILDREN
  is passed here, then the option may be set on a child of obj.

  @return 0 if the value has been set, or an AVERROR code in case of
  error:
  AVERROR_OPTION_NOT_FOUND if no matching option exists
  AVERROR(ERANGE) if the value is out of range
  AVERROR(EINVAL) if the value is not valid
*/
func AVOptSet(obj unsafe.Pointer, name *CStr, val *CStr, searchFlags int) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	var tmpval *C.char
	if val != nil {
		tmpval = val.ptr
	}
	ret := C.av_opt_set(obj, tmpname, tmpval, C.int(searchFlags))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_set_int ---

// AVOptSetInt wraps av_opt_set_int.
func AVOptSetInt(obj unsafe.Pointer, name *CStr, val int64, searchFlags int) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_opt_set_int(obj, tmpname, C.int64_t(val), C.int(searchFlags))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_set_double ---

// AVOptSetDouble wraps av_opt_set_double.
func AVOptSetDouble(obj unsafe.Pointer, name *CStr, val float64, searchFlags int) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_opt_set_double(obj, tmpname, C.double(val), C.int(searchFlags))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_set_q ---

// AVOptSetQ wraps av_opt_set_q.
func AVOptSetQ(obj unsafe.Pointer, name *CStr, val *AVRational, searchFlags int) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_opt_set_q(obj, tmpname, val.value, C.int(searchFlags))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_set_bin ---

// AVOptSetBin wraps av_opt_set_bin.
func AVOptSetBin(obj unsafe.Pointer, name *CStr, val unsafe.Pointer, size int, searchFlags int) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_opt_set_bin(obj, tmpname, (*C.uint8_t)(val), C.int(size), C.int(searchFlags))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_set_image_size ---

// AVOptSetImageSize wraps av_opt_set_image_size.
func AVOptSetImageSize(obj unsafe.Pointer, name *CStr, w int, h int, searchFlags int) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_opt_set_image_size(obj, tmpname, C.int(w), C.int(h), C.int(searchFlags))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_set_pixel_fmt ---

// AVOptSetPixelFmt wraps av_opt_set_pixel_fmt.
func AVOptSetPixelFmt(obj unsafe.Pointer, name *CStr, fmt AVPixelFormat, searchFlags int) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_opt_set_pixel_fmt(obj, tmpname, C.enum_AVPixelFormat(fmt), C.int(searchFlags))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_set_sample_fmt ---

// AVOptSetSampleFmt wraps av_opt_set_sample_fmt.
func AVOptSetSampleFmt(obj unsafe.Pointer, name *CStr, fmt AVSampleFormat, searchFlags int) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_opt_set_sample_fmt(obj, tmpname, C.enum_AVSampleFormat(fmt), C.int(searchFlags))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_set_video_rate ---

// AVOptSetVideoRate wraps av_opt_set_video_rate.
func AVOptSetVideoRate(obj unsafe.Pointer, name *CStr, val *AVRational, searchFlags int) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_opt_set_video_rate(obj, tmpname, val.value, C.int(searchFlags))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_set_chlayout ---

// AVOptSetChlayout wraps av_opt_set_chlayout.
/*
  @note Any old chlayout present is discarded and replaced with a copy of the new one. The
  caller still owns layout and is responsible for uninitializing it.
*/
func AVOptSetChlayout(obj unsafe.Pointer, name *CStr, layout *AVChannelLayout, searchFlags int) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	var tmplayout *C.AVChannelLayout
	if layout != nil {
		tmplayout = layout.ptr
	}
	ret := C.av_opt_set_chlayout(obj, tmpname, tmplayout, C.int(searchFlags))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_set_dict_val ---

// AVOptSetDictVal wraps av_opt_set_dict_val.
/*
  @note Any old dictionary present is discarded and replaced with a copy of the new one. The
  caller still owns val is and responsible for freeing it.
*/
func AVOptSetDictVal(obj unsafe.Pointer, name *CStr, val *AVDictionary, searchFlags int) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	var tmpval *C.AVDictionary
	if val != nil {
		tmpval = val.ptr
	}
	ret := C.av_opt_set_dict_val(obj, tmpname, tmpval, C.int(searchFlags))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_set_array ---

// AVOptSetArray wraps av_opt_set_array.
/*
  Add, replace, or remove elements for an array option. Which of these
  operations is performed depends on the values of val and search_flags.

  @param start_elem Index of the first array element to modify; must not be
                    larger than array size as returned by
                    av_opt_get_array_size().
  @param nb_elems number of array elements to modify; when val is NULL,
                  start_elem+nb_elems must not be larger than array size as
                  returned by av_opt_get_array_size()

  @param val_type Option type corresponding to the type of val, ignored when val is
                  NULL.

                  The effect of this function will will be as if av_opt_setX()
                  was called for each element, where X is specified by type.
                  E.g. AV_OPT_TYPE_STRING corresponds to av_opt_set().

                  Typically this should be the same as the scalarized type of
                  the AVOption being set, but certain conversions are also
                  possible - the same as those done by the corresponding
                  av_opt_set*() function. E.g. any option type can be set from
                  a string, numeric types can be set from int64, double, or
                  rational, etc.

  @param val Array with nb_elems elements or NULL.

             When NULL, nb_elems array elements starting at start_elem are
             removed from the array. Any array elements remaining at the end
             are shifted by nb_elems towards the first element in order to keep
             the array contiguous.

             Otherwise (val is non-NULL), the type of val must match the
             underlying C type as documented for val_type.

             When AV_OPT_ARRAY_REPLACE is not set in search_flags, the array is
             enlarged by nb_elems, and the contents of val are inserted at
             start_elem. Previously existing array elements from start_elem
             onwards (if present) are shifted by nb_elems away from the first
             element in order to make space for the new elements.

             When AV_OPT_ARRAY_REPLACE is set in search_flags, the contents
             of val replace existing array elements from start_elem to
             start_elem+nb_elems (if present). New array size is
             max(start_elem + nb_elems, old array size).
*/
func AVOptSetArray(obj unsafe.Pointer, name *CStr, searchFlags int, startElem uint, nbElems uint, valType AVOptionType, val unsafe.Pointer) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_opt_set_array(obj, tmpname, C.int(searchFlags), C.uint(startElem), C.uint(nbElems), C.enum_AVOptionType(valType), val)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_get ---

// av_opt_get skipped due to outVal

// --- Function av_opt_get_int ---

// AVOptGetInt wraps av_opt_get_int.
func AVOptGetInt(obj unsafe.Pointer, name *CStr, searchFlags int, outVal *int64) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_opt_get_int(obj, tmpname, C.int(searchFlags), (*C.int64_t)(unsafe.Pointer(outVal)))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_get_double ---

// AVOptGetDouble wraps av_opt_get_double.
func AVOptGetDouble(obj unsafe.Pointer, name *CStr, searchFlags int, outVal *float64) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_opt_get_double(obj, tmpname, C.int(searchFlags), (*C.double)(unsafe.Pointer(outVal)))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_get_q ---

// av_opt_get_q skipped due to outVal

// --- Function av_opt_get_image_size ---

// AVOptGetImageSize wraps av_opt_get_image_size.
func AVOptGetImageSize(obj unsafe.Pointer, name *CStr, searchFlags int, wOut *int, hOut *int) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_opt_get_image_size(obj, tmpname, C.int(searchFlags), (*C.int)(unsafe.Pointer(wOut)), (*C.int)(unsafe.Pointer(hOut)))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_get_pixel_fmt ---

// AVOptGetPixelFmt wraps av_opt_get_pixel_fmt.
func AVOptGetPixelFmt(obj unsafe.Pointer, name *CStr, searchFlags int, outFmt *AVPixelFormat) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	var tmpoutFmt *C.enum_AVPixelFormat
	if outFmt != nil {
		tmpoutFmt = (*C.enum_AVPixelFormat)(unsafe.Pointer(outFmt))
	}
	ret := C.av_opt_get_pixel_fmt(obj, tmpname, C.int(searchFlags), tmpoutFmt)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_get_sample_fmt ---

// AVOptGetSampleFmt wraps av_opt_get_sample_fmt.
func AVOptGetSampleFmt(obj unsafe.Pointer, name *CStr, searchFlags int, outFmt *AVSampleFormat) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	var tmpoutFmt *C.enum_AVSampleFormat
	if outFmt != nil {
		tmpoutFmt = (*C.enum_AVSampleFormat)(unsafe.Pointer(outFmt))
	}
	ret := C.av_opt_get_sample_fmt(obj, tmpname, C.int(searchFlags), tmpoutFmt)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_get_video_rate ---

// av_opt_get_video_rate skipped due to outVal

// --- Function av_opt_get_chlayout ---

// AVOptGetChlayout wraps av_opt_get_chlayout.
/*
  @param[out] layout The returned layout is a copy of the actual value and must
  be freed with av_channel_layout_uninit() by the caller
*/
func AVOptGetChlayout(obj unsafe.Pointer, name *CStr, searchFlags int, layout *AVChannelLayout) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	var tmplayout *C.AVChannelLayout
	if layout != nil {
		tmplayout = layout.ptr
	}
	ret := C.av_opt_get_chlayout(obj, tmpname, C.int(searchFlags), tmplayout)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_get_dict_val ---

// AVOptGetDictVal wraps av_opt_get_dict_val.
/*
  @param[out] out_val The returned dictionary is a copy of the actual value and must
  be freed with av_dict_free() by the caller
*/
func AVOptGetDictVal(obj unsafe.Pointer, name *CStr, searchFlags int, outVal **AVDictionary) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	var ptroutVal **C.AVDictionary
	var tmpoutVal *C.AVDictionary
	var oldTmpoutVal *C.AVDictionary
	if outVal != nil {
		inneroutVal := *outVal
		if inneroutVal != nil {
			tmpoutVal = inneroutVal.ptr
			oldTmpoutVal = tmpoutVal
		}
		ptroutVal = &tmpoutVal
	}
	ret := C.av_opt_get_dict_val(obj, tmpname, C.int(searchFlags), ptroutVal)
	if tmpoutVal != oldTmpoutVal && outVal != nil {
		if tmpoutVal != nil {
			*outVal = &AVDictionary{ptr: tmpoutVal}
		} else {
			*outVal = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_get_array_size ---

// AVOptGetArraySize wraps av_opt_get_array_size.
//
//	For an array-type option, get the number of elements in the array.
func AVOptGetArraySize(obj unsafe.Pointer, name *CStr, searchFlags int, outVal *uint) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_opt_get_array_size(obj, tmpname, C.int(searchFlags), (*C.uint)(unsafe.Pointer(outVal)))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_get_array ---

// AVOptGetArray wraps av_opt_get_array.
/*
  For an array-type option, retrieve the values of one or more array elements.

  @param start_elem index of the first array element to retrieve
  @param nb_elems number of array elements to retrieve; start_elem+nb_elems
                  must not be larger than array size as returned by
                  av_opt_get_array_size()

  @param out_type Option type corresponding to the desired output.

                  The array elements produced by this function will
                  will be as if av_opt_getX() was called for each element,
                  where X is specified by out_type. E.g. AV_OPT_TYPE_STRING
                  corresponds to av_opt_get().

                  Typically this should be the same as the scalarized type of
                  the AVOption being retrieved, but certain conversions are
                  also possible - the same as those done by the corresponding
                  av_opt_get*() function. E.g. any option type can be retrieved
                  as a string, numeric types can be retrieved as int64, double,
                  or rational, etc.

  @param out_val  Array with nb_elems members into which the output will be
                  written. The array type must match the underlying C type as
                  documented for out_type, and be zeroed on entry to this
                  function.

                  For dynamically allocated types (strings, binary, dicts,
                  etc.), the result is owned and freed by the caller.
*/
func AVOptGetArray(obj unsafe.Pointer, name *CStr, searchFlags int, startElem uint, nbElems uint, outType AVOptionType, outVal unsafe.Pointer) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_opt_get_array(obj, tmpname, C.int(searchFlags), C.uint(startElem), C.uint(nbElems), C.enum_AVOptionType(outType), outVal)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_eval_flags ---

// AVOptEvalFlags wraps av_opt_eval_flags.
/*
  This group of functions can be used to evaluate option strings
  and get numbers out of them. They do the same thing as av_opt_set(),
  except the result is written into the caller-supplied pointer.

  @param obj a struct whose first element is a pointer to AVClass.
  @param o an option for which the string is to be evaluated.
  @param val string to be evaluated.
  @param *_out value of the string will be written here.

  @return 0 on success, a negative number on failure.
*/
func AVOptEvalFlags(obj unsafe.Pointer, o *AVOption, val *CStr, flagsOut *int) (int, error) {
	var tmpo *C.AVOption
	if o != nil {
		tmpo = o.ptr
	}
	var tmpval *C.char
	if val != nil {
		tmpval = val.ptr
	}
	ret := C.av_opt_eval_flags(obj, tmpo, tmpval, (*C.int)(unsafe.Pointer(flagsOut)))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_eval_int ---

// AVOptEvalInt wraps av_opt_eval_int.
func AVOptEvalInt(obj unsafe.Pointer, o *AVOption, val *CStr, intOut *int) (int, error) {
	var tmpo *C.AVOption
	if o != nil {
		tmpo = o.ptr
	}
	var tmpval *C.char
	if val != nil {
		tmpval = val.ptr
	}
	ret := C.av_opt_eval_int(obj, tmpo, tmpval, (*C.int)(unsafe.Pointer(intOut)))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_eval_uint ---

// AVOptEvalUint wraps av_opt_eval_uint.
func AVOptEvalUint(obj unsafe.Pointer, o *AVOption, val *CStr, uintOut *uint) (int, error) {
	var tmpo *C.AVOption
	if o != nil {
		tmpo = o.ptr
	}
	var tmpval *C.char
	if val != nil {
		tmpval = val.ptr
	}
	ret := C.av_opt_eval_uint(obj, tmpo, tmpval, (*C.uint)(unsafe.Pointer(uintOut)))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_eval_int64 ---

// AVOptEvalInt64 wraps av_opt_eval_int64.
func AVOptEvalInt64(obj unsafe.Pointer, o *AVOption, val *CStr, int64Out *int64) (int, error) {
	var tmpo *C.AVOption
	if o != nil {
		tmpo = o.ptr
	}
	var tmpval *C.char
	if val != nil {
		tmpval = val.ptr
	}
	ret := C.av_opt_eval_int64(obj, tmpo, tmpval, (*C.int64_t)(unsafe.Pointer(int64Out)))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_eval_float ---

// AVOptEvalFloat wraps av_opt_eval_float.
func AVOptEvalFloat(obj unsafe.Pointer, o *AVOption, val *CStr, floatOut *float32) (int, error) {
	var tmpo *C.AVOption
	if o != nil {
		tmpo = o.ptr
	}
	var tmpval *C.char
	if val != nil {
		tmpval = val.ptr
	}
	ret := C.av_opt_eval_float(obj, tmpo, tmpval, (*C.float)(unsafe.Pointer(floatOut)))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_eval_double ---

// AVOptEvalDouble wraps av_opt_eval_double.
func AVOptEvalDouble(obj unsafe.Pointer, o *AVOption, val *CStr, doubleOut *float64) (int, error) {
	var tmpo *C.AVOption
	if o != nil {
		tmpo = o.ptr
	}
	var tmpval *C.char
	if val != nil {
		tmpval = val.ptr
	}
	ret := C.av_opt_eval_double(obj, tmpo, tmpval, (*C.double)(unsafe.Pointer(doubleOut)))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_eval_q ---

// av_opt_eval_q skipped due to qOut

// --- Function av_opt_ptr ---

// AVOptPtr wraps av_opt_ptr.
/*
  Gets a pointer to the requested field in a struct.
  This function allows accessing a struct even when its fields are moved or
  renamed since the application making the access has been compiled,

  @returns a pointer to the field, it can be cast to the correct type and read
           or written to.

  @deprecated direct access to AVOption-exported fields is not supported
*/
func AVOptPtr(avclass *AVClass, obj unsafe.Pointer, name *CStr) unsafe.Pointer {
	var tmpavclass *C.AVClass
	if avclass != nil {
		tmpavclass = avclass.ptr
	}
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_opt_ptr(tmpavclass, obj, tmpname)
	return ret
}

// --- Function av_opt_is_set_to_default ---

// AVOptIsSetToDefault wraps av_opt_is_set_to_default.
/*
  Check if given option is set to its default value.

  Options o must belong to the obj. This function must not be called to check child's options state.
  @see av_opt_is_set_to_default_by_name().

  @param obj  AVClass object to check option on
  @param o    option to be checked
  @return     >0 when option is set to its default,
               0 when option is not set its default,
              <0 on error
*/
func AVOptIsSetToDefault(obj unsafe.Pointer, o *AVOption) (int, error) {
	var tmpo *C.AVOption
	if o != nil {
		tmpo = o.ptr
	}
	ret := C.av_opt_is_set_to_default(obj, tmpo)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_is_set_to_default_by_name ---

// AVOptIsSetToDefaultByName wraps av_opt_is_set_to_default_by_name.
/*
  Check if given option is set to its default value.

  @param obj          AVClass object to check option on
  @param name         option name
  @param search_flags combination of AV_OPT_SEARCH_*
  @return             >0 when option is set to its default,
                      0 when option is not set its default,
                      <0 on error
*/
func AVOptIsSetToDefaultByName(obj unsafe.Pointer, name *CStr, searchFlags int) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_opt_is_set_to_default_by_name(obj, tmpname, C.int(searchFlags))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_flag_is_set ---

// AVOptFlagIsSet wraps av_opt_flag_is_set.
/*
  Check whether a particular flag is set in a flags field.

  @param field_name the name of the flag field option
  @param flag_name the name of the flag to check
  @return non-zero if the flag is set, zero if the flag isn't set,
          isn't of the right type, or the flags field doesn't exist.
*/
func AVOptFlagIsSet(obj unsafe.Pointer, fieldName *CStr, flagName *CStr) (int, error) {
	var tmpfieldName *C.char
	if fieldName != nil {
		tmpfieldName = fieldName.ptr
	}
	var tmpflagName *C.char
	if flagName != nil {
		tmpflagName = flagName.ptr
	}
	ret := C.av_opt_flag_is_set(obj, tmpfieldName, tmpflagName)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_serialize ---

// av_opt_serialize skipped due to buffer

// --- Function av_opt_freep_ranges ---

// AVOptFreepRanges wraps av_opt_freep_ranges.
//
//	Free an AVOptionRanges struct and set it to NULL.
func AVOptFreepRanges(ranges **AVOptionRanges) {
	var ptrranges **C.AVOptionRanges
	var tmpranges *C.AVOptionRanges
	var oldTmpranges *C.AVOptionRanges
	if ranges != nil {
		innerranges := *ranges
		if innerranges != nil {
			tmpranges = innerranges.ptr
			oldTmpranges = tmpranges
		}
		ptrranges = &tmpranges
	}
	C.av_opt_freep_ranges(ptrranges)
	if tmpranges != oldTmpranges && ranges != nil {
		if tmpranges != nil {
			*ranges = &AVOptionRanges{ptr: tmpranges}
		} else {
			*ranges = nil
		}
	}
}

// --- Function av_opt_query_ranges ---

// AVOptQueryRanges wraps av_opt_query_ranges.
/*
  Get a list of allowed ranges for the given option.

  The returned list may depend on other fields in obj like for example profile.

  @param flags is a bitmask of flags, undefined flags should not be set and should be ignored
               AV_OPT_SEARCH_FAKE_OBJ indicates that the obj is a double pointer to a AVClass instead of a full instance
               AV_OPT_MULTI_COMPONENT_RANGE indicates that function may return more than one component, @see AVOptionRanges

  The result must be freed with av_opt_freep_ranges.

  @return number of components returned on success, a negative error code otherwise
*/
func AVOptQueryRanges(param0 **AVOptionRanges, obj unsafe.Pointer, key *CStr, flags int) (int, error) {
	var ptrparam0 **C.AVOptionRanges
	var tmpparam0 *C.AVOptionRanges
	var oldTmpparam0 *C.AVOptionRanges
	if param0 != nil {
		innerparam0 := *param0
		if innerparam0 != nil {
			tmpparam0 = innerparam0.ptr
			oldTmpparam0 = tmpparam0
		}
		ptrparam0 = &tmpparam0
	}
	var tmpkey *C.char
	if key != nil {
		tmpkey = key.ptr
	}
	ret := C.av_opt_query_ranges(ptrparam0, obj, tmpkey, C.int(flags))
	if tmpparam0 != oldTmpparam0 && param0 != nil {
		if tmpparam0 != nil {
			*param0 = &AVOptionRanges{ptr: tmpparam0}
		} else {
			*param0 = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function av_opt_query_ranges_default ---

// AVOptQueryRangesDefault wraps av_opt_query_ranges_default.
/*
  Get a default list of allowed ranges for the given option.

  This list is constructed without using the AVClass.query_ranges() callback
  and can be used as fallback from within the callback.

  @param flags is a bitmask of flags, undefined flags should not be set and should be ignored
               AV_OPT_SEARCH_FAKE_OBJ indicates that the obj is a double pointer to a AVClass instead of a full instance
               AV_OPT_MULTI_COMPONENT_RANGE indicates that function may return more than one component, @see AVOptionRanges

  The result must be freed with av_opt_free_ranges.

  @return number of components returned on success, a negative error code otherwise
*/
func AVOptQueryRangesDefault(param0 **AVOptionRanges, obj unsafe.Pointer, key *CStr, flags int) (int, error) {
	var ptrparam0 **C.AVOptionRanges
	var tmpparam0 *C.AVOptionRanges
	var oldTmpparam0 *C.AVOptionRanges
	if param0 != nil {
		innerparam0 := *param0
		if innerparam0 != nil {
			tmpparam0 = innerparam0.ptr
			oldTmpparam0 = tmpparam0
		}
		ptrparam0 = &tmpparam0
	}
	var tmpkey *C.char
	if key != nil {
		tmpkey = key.ptr
	}
	ret := C.av_opt_query_ranges_default(ptrparam0, obj, tmpkey, C.int(flags))
	if tmpparam0 != oldTmpparam0 && param0 != nil {
		if tmpparam0 != nil {
			*param0 = &AVOptionRanges{ptr: tmpparam0}
		} else {
			*param0 = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function av_parse_ratio ---

// av_parse_ratio skipped due to q

// --- Function av_parse_video_size ---

// AVParseVideoSize wraps av_parse_video_size.
/*
  Parse str and put in width_ptr and height_ptr the detected values.

  @param[in,out] width_ptr pointer to the variable which will contain the detected
  width value
  @param[in,out] height_ptr pointer to the variable which will contain the detected
  height value
  @param[in] str the string to parse: it has to be a string in the format
  width x height or a valid video size abbreviation.
  @return >= 0 on success, a negative error code otherwise
*/
func AVParseVideoSize(widthPtr *int, heightPtr *int, str *CStr) (int, error) {
	var tmpstr *C.char
	if str != nil {
		tmpstr = str.ptr
	}
	ret := C.av_parse_video_size((*C.int)(unsafe.Pointer(widthPtr)), (*C.int)(unsafe.Pointer(heightPtr)), tmpstr)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_parse_video_rate ---

// av_parse_video_rate skipped due to rate

// --- Function av_parse_color ---

// AVParseColor wraps av_parse_color.
/*
  Put the RGBA values that correspond to color_string in rgba_color.

  @param rgba_color 4-elements array of uint8_t values, where the respective
  red, green, blue and alpha component values are written.
  @param color_string a string specifying a color. It can be the name of
  a color (case insensitive match) or a [0x|#]RRGGBB[AA] sequence,
  possibly followed by "@" and a string representing the alpha
  component.
  The alpha component may be a string composed by "0x" followed by an
  hexadecimal number or a decimal number between 0.0 and 1.0, which
  represents the opacity value (0x00/0.0 means completely transparent,
  0xff/1.0 completely opaque).
  If the alpha component is not specified then 0xff is assumed.
  The string "random" will result in a random color.
  @param slen length of the initial part of color_string containing the
  color. It can be set to -1 if color_string is a null terminated string
  containing nothing else than the color.
  @param log_ctx a pointer to an arbitrary struct of which the first field
  is a pointer to an AVClass struct (used for av_log()). Can be NULL.
  @return >= 0 in case of success, a negative value in case of
  failure (for example if color_string cannot be parsed).
*/
func AVParseColor(rgbaColor unsafe.Pointer, colorString *CStr, slen int, logCtx unsafe.Pointer) (int, error) {
	var tmpcolorString *C.char
	if colorString != nil {
		tmpcolorString = colorString.ptr
	}
	ret := C.av_parse_color((*C.uint8_t)(rgbaColor), tmpcolorString, C.int(slen), logCtx)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_get_known_color_name ---

// av_get_known_color_name skipped due to rgb

// --- Function av_parse_time ---

// av_parse_time skipped due to timeval (non-output primitive pointer)

// --- Function av_find_info_tag ---

// AVFindInfoTag wraps av_find_info_tag.
/*
  Attempt to find a specific tag in a URL.

  syntax: '?tag1=val1&tag2=val2...'. Little URL decoding is done.
  Return 1 if found.
*/
func AVFindInfoTag(arg *CStr, argSize int, tag1 *CStr, info *CStr) (int, error) {
	var tmparg *C.char
	if arg != nil {
		tmparg = arg.ptr
	}
	var tmptag1 *C.char
	if tag1 != nil {
		tmptag1 = tag1.ptr
	}
	var tmpinfo *C.char
	if info != nil {
		tmpinfo = info.ptr
	}
	ret := C.av_find_info_tag(tmparg, C.int(argSize), tmptag1, tmpinfo)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_small_strptime ---

// av_small_strptime skipped due to dt.

// --- Function av_timegm ---

// av_timegm skipped due to tm.

// --- Function av_get_bits_per_pixel ---

// AVGetBitsPerPixel wraps av_get_bits_per_pixel.
/*
  Return the number of bits per pixel used by the pixel format
  described by pixdesc. Note that this is not the same as the number
  of bits per sample.

  The returned number of bits refers to the number of bits actually
  used for storing the pixel information, that is padding bits are
  not counted.
*/
func AVGetBitsPerPixel(pixdesc *AVPixFmtDescriptor) (int, error) {
	var tmppixdesc *C.AVPixFmtDescriptor
	if pixdesc != nil {
		tmppixdesc = pixdesc.ptr
	}
	ret := C.av_get_bits_per_pixel(tmppixdesc)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_get_padded_bits_per_pixel ---

// AVGetPaddedBitsPerPixel wraps av_get_padded_bits_per_pixel.
/*
  Return the number of bits per pixel for the pixel format
  described by pixdesc, including any padding or unused bits.
*/
func AVGetPaddedBitsPerPixel(pixdesc *AVPixFmtDescriptor) (int, error) {
	var tmppixdesc *C.AVPixFmtDescriptor
	if pixdesc != nil {
		tmppixdesc = pixdesc.ptr
	}
	ret := C.av_get_padded_bits_per_pixel(tmppixdesc)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_pix_fmt_desc_get ---

// AVPixFmtDescGet wraps av_pix_fmt_desc_get.
/*
  @return a pixel format descriptor for provided pixel format or NULL if
  this pixel format is unknown.
*/
func AVPixFmtDescGet(pixFmt AVPixelFormat) *AVPixFmtDescriptor {
	ret := C.av_pix_fmt_desc_get(C.enum_AVPixelFormat(pixFmt))
	var retMapped *AVPixFmtDescriptor
	if ret != nil {
		retMapped = &AVPixFmtDescriptor{ptr: ret}
	}
	return retMapped
}

// --- Function av_pix_fmt_desc_next ---

// AVPixFmtDescNext wraps av_pix_fmt_desc_next.
/*
  Iterate over all pixel format descriptors known to libavutil.

  @param prev previous descriptor. NULL to get the first descriptor.

  @return next descriptor or NULL after the last descriptor
*/
func AVPixFmtDescNext(prev *AVPixFmtDescriptor) *AVPixFmtDescriptor {
	var tmpprev *C.AVPixFmtDescriptor
	if prev != nil {
		tmpprev = prev.ptr
	}
	ret := C.av_pix_fmt_desc_next(tmpprev)
	var retMapped *AVPixFmtDescriptor
	if ret != nil {
		retMapped = &AVPixFmtDescriptor{ptr: ret}
	}
	return retMapped
}

// --- Function av_pix_fmt_desc_get_id ---

// AVPixFmtDescGetId wraps av_pix_fmt_desc_get_id.
/*
  @return an AVPixelFormat id described by desc, or AV_PIX_FMT_NONE if desc
  is not a valid pointer to a pixel format descriptor.
*/
func AVPixFmtDescGetId(desc *AVPixFmtDescriptor) AVPixelFormat {
	var tmpdesc *C.AVPixFmtDescriptor
	if desc != nil {
		tmpdesc = desc.ptr
	}
	ret := C.av_pix_fmt_desc_get_id(tmpdesc)
	return AVPixelFormat(ret)
}

// --- Function av_pix_fmt_get_chroma_sub_sample ---

// AVPixFmtGetChromaSubSample wraps av_pix_fmt_get_chroma_sub_sample.
/*
  Utility function to access log2_chroma_w log2_chroma_h from
  the pixel format AVPixFmtDescriptor.

  @param[in]  pix_fmt the pixel format
  @param[out] h_shift store log2_chroma_w (horizontal/width shift)
  @param[out] v_shift store log2_chroma_h (vertical/height shift)

  @return 0 on success, AVERROR(ENOSYS) on invalid or unknown pixel format
*/
func AVPixFmtGetChromaSubSample(pixFmt AVPixelFormat, hShift *int, vShift *int) (int, error) {
	ret := C.av_pix_fmt_get_chroma_sub_sample(C.enum_AVPixelFormat(pixFmt), (*C.int)(unsafe.Pointer(hShift)), (*C.int)(unsafe.Pointer(vShift)))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_pix_fmt_count_planes ---

// AVPixFmtCountPlanes wraps av_pix_fmt_count_planes.
/*
  @return number of planes in pix_fmt, a negative AVERROR if pix_fmt is not a
  valid pixel format.
*/
func AVPixFmtCountPlanes(pixFmt AVPixelFormat) (int, error) {
	ret := C.av_pix_fmt_count_planes(C.enum_AVPixelFormat(pixFmt))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_color_range_name ---

// AVColorRangeName wraps av_color_range_name.
//
//	@return the name for provided color range or NULL if unknown.
func AVColorRangeName(_range AVColorRange) *CStr {
	ret := C.av_color_range_name(C.enum_AVColorRange(_range))
	return wrapCStr(ret)
}

// --- Function av_color_range_from_name ---

// AVColorRangeFromName wraps av_color_range_from_name.
//
//	@return the AVColorRange value for name or an AVError if not found.
func AVColorRangeFromName(name *CStr) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_color_range_from_name(tmpname)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_color_primaries_name ---

// AVColorPrimariesName wraps av_color_primaries_name.
//
//	@return the name for provided color primaries or NULL if unknown.
func AVColorPrimariesName(primaries AVColorPrimaries) *CStr {
	ret := C.av_color_primaries_name(C.enum_AVColorPrimaries(primaries))
	return wrapCStr(ret)
}

// --- Function av_color_primaries_from_name ---

// AVColorPrimariesFromName wraps av_color_primaries_from_name.
//
//	@return the AVColorPrimaries value for name or an AVError if not found.
func AVColorPrimariesFromName(name *CStr) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_color_primaries_from_name(tmpname)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_color_transfer_name ---

// AVColorTransferName wraps av_color_transfer_name.
//
//	@return the name for provided color transfer or NULL if unknown.
func AVColorTransferName(transfer AVColorTransferCharacteristic) *CStr {
	ret := C.av_color_transfer_name(C.enum_AVColorTransferCharacteristic(transfer))
	return wrapCStr(ret)
}

// --- Function av_color_transfer_from_name ---

// AVColorTransferFromName wraps av_color_transfer_from_name.
//
//	@return the AVColorTransferCharacteristic value for name or an AVError if not found.
func AVColorTransferFromName(name *CStr) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_color_transfer_from_name(tmpname)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_color_space_name ---

// AVColorSpaceName wraps av_color_space_name.
//
//	@return the name for provided color space or NULL if unknown.
func AVColorSpaceName(space AVColorSpace) *CStr {
	ret := C.av_color_space_name(C.enum_AVColorSpace(space))
	return wrapCStr(ret)
}

// --- Function av_color_space_from_name ---

// AVColorSpaceFromName wraps av_color_space_from_name.
//
//	@return the AVColorSpace value for name or an AVError if not found.
func AVColorSpaceFromName(name *CStr) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_color_space_from_name(tmpname)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_chroma_location_name ---

// AVChromaLocationName wraps av_chroma_location_name.
//
//	@return the name for provided chroma location or NULL if unknown.
func AVChromaLocationName(location AVChromaLocation) *CStr {
	ret := C.av_chroma_location_name(C.enum_AVChromaLocation(location))
	return wrapCStr(ret)
}

// --- Function av_chroma_location_from_name ---

// AVChromaLocationFromName wraps av_chroma_location_from_name.
//
//	@return the AVChromaLocation value for name or an AVError if not found.
func AVChromaLocationFromName(name *CStr) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_chroma_location_from_name(tmpname)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_chroma_location_enum_to_pos ---

// av_chroma_location_enum_to_pos skipped due to xpos (non-output primitive pointer)

// --- Function av_chroma_location_pos_to_enum ---

// AVChromaLocationPosToEnum wraps av_chroma_location_pos_to_enum.
/*
  Converts swscale x/y chroma position to AVChromaLocation.

  The positions represent the chroma (0,0) position in a coordinates system
  with luma (0,0) representing the origin and luma(1,1) representing 256,256

  @param xpos  horizontal chroma sample position
  @param ypos  vertical   chroma sample position
*/
func AVChromaLocationPosToEnum(xpos int, ypos int) AVChromaLocation {
	ret := C.av_chroma_location_pos_to_enum(C.int(xpos), C.int(ypos))
	return AVChromaLocation(ret)
}

// --- Function av_get_pix_fmt ---

// AVGetPixFmt wraps av_get_pix_fmt.
/*
  Return the pixel format corresponding to name.

  If there is no pixel format with name name, then looks for a
  pixel format with the name corresponding to the native endian
  format of name.
  For example in a little-endian system, first looks for "gray16",
  then for "gray16le".

  Finally if no pixel format has been found, returns AV_PIX_FMT_NONE.
*/
func AVGetPixFmt(name *CStr) AVPixelFormat {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_get_pix_fmt(tmpname)
	return AVPixelFormat(ret)
}

// --- Function av_get_pix_fmt_name ---

// AVGetPixFmtName wraps av_get_pix_fmt_name.
/*
  Return the short name for a pixel format, NULL in case pix_fmt is
  unknown.

  @see av_get_pix_fmt(), av_get_pix_fmt_string()
*/
func AVGetPixFmtName(pixFmt AVPixelFormat) *CStr {
	ret := C.av_get_pix_fmt_name(C.enum_AVPixelFormat(pixFmt))
	return wrapCStr(ret)
}

// --- Function av_get_pix_fmt_string ---

// AVGetPixFmtString wraps av_get_pix_fmt_string.
/*
  Print in buf the string corresponding to the pixel format with
  number pix_fmt, or a header if pix_fmt is negative.

  @param buf the buffer where to write the string
  @param buf_size the size of buf
  @param pix_fmt the number of the pixel format to print the
  corresponding info string, or a negative value to print the
  corresponding header.
*/
func AVGetPixFmtString(buf *CStr, bufSize int, pixFmt AVPixelFormat) *CStr {
	var tmpbuf *C.char
	if buf != nil {
		tmpbuf = buf.ptr
	}
	ret := C.av_get_pix_fmt_string(tmpbuf, C.int(bufSize), C.enum_AVPixelFormat(pixFmt))
	return wrapCStr(ret)
}

// --- Function av_read_image_line2 ---

// av_read_image_line2 skipped due to const array param data

// --- Function av_read_image_line ---

// av_read_image_line skipped due to dst (non-output primitive pointer)

// --- Function av_write_image_line2 ---

// av_write_image_line2 skipped due to const array param data

// --- Function av_write_image_line ---

// av_write_image_line skipped due to src (non-output primitive pointer)

// --- Function av_pix_fmt_swap_endianness ---

// AVPixFmtSwapEndianness wraps av_pix_fmt_swap_endianness.
/*
  Utility function to swap the endianness of a pixel format.

  @param[in]  pix_fmt the pixel format

  @return pixel format with swapped endianness if it exists,
  otherwise AV_PIX_FMT_NONE
*/
func AVPixFmtSwapEndianness(pixFmt AVPixelFormat) AVPixelFormat {
	ret := C.av_pix_fmt_swap_endianness(C.enum_AVPixelFormat(pixFmt))
	return AVPixelFormat(ret)
}

// --- Function av_get_pix_fmt_loss ---

// AVGetPixFmtLoss wraps av_get_pix_fmt_loss.
/*
  Compute what kind of losses will occur when converting from one specific
  pixel format to another.
  When converting from one pixel format to another, information loss may occur.
  For example, when converting from RGB24 to GRAY, the color information will
  be lost. Similarly, other losses occur when converting from some formats to
  other formats. These losses can involve loss of chroma, but also loss of
  resolution, loss of color depth, loss due to the color space conversion, loss
  of the alpha bits or loss due to color quantization.
  av_get_fix_fmt_loss() informs you about the various types of losses
  which will occur when converting from one pixel format to another.

  @param[in] dst_pix_fmt destination pixel format
  @param[in] src_pix_fmt source pixel format
  @param[in] has_alpha Whether the source pixel format alpha channel is used.
  @return Combination of flags informing you what kind of losses will occur
  (maximum loss for an invalid dst_pix_fmt).
*/
func AVGetPixFmtLoss(dstPixFmt AVPixelFormat, srcPixFmt AVPixelFormat, hasAlpha int) (int, error) {
	ret := C.av_get_pix_fmt_loss(C.enum_AVPixelFormat(dstPixFmt), C.enum_AVPixelFormat(srcPixFmt), C.int(hasAlpha))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_find_best_pix_fmt_of_2 ---

// AVFindBestPixFmtOf2 wraps av_find_best_pix_fmt_of_2.
/*
  Compute what kind of losses will occur when converting from one specific
  pixel format to another.
  When converting from one pixel format to another, information loss may occur.
  For example, when converting from RGB24 to GRAY, the color information will
  be lost. Similarly, other losses occur when converting from some formats to
  other formats. These losses can involve loss of chroma, but also loss of
  resolution, loss of color depth, loss due to the color space conversion, loss
  of the alpha bits or loss due to color quantization.
  av_get_fix_fmt_loss() informs you about the various types of losses
  which will occur when converting from one pixel format to another.

  @param[in] dst_pix_fmt destination pixel format
  @param[in] src_pix_fmt source pixel format
  @param[in] has_alpha Whether the source pixel format alpha channel is used.
  @return Combination of flags informing you what kind of losses will occur
  (maximum loss for an invalid dst_pix_fmt).
*/
func AVFindBestPixFmtOf2(dstPixFmt1 AVPixelFormat, dstPixFmt2 AVPixelFormat, srcPixFmt AVPixelFormat, hasAlpha int, lossPtr *int) AVPixelFormat {
	ret := C.av_find_best_pix_fmt_of_2(C.enum_AVPixelFormat(dstPixFmt1), C.enum_AVPixelFormat(dstPixFmt2), C.enum_AVPixelFormat(srcPixFmt), C.int(hasAlpha), (*C.int)(unsafe.Pointer(lossPtr)))
	return AVPixelFormat(ret)
}

// --- Function av_pixelutils_get_sad_fn ---

// AVPixelutilsGetSadFn wraps av_pixelutils_get_sad_fn.
/*
  Get a potentially optimized pointer to a Sum-of-absolute-differences
  function (see the av_pixelutils_sad_fn prototype).

  @param w_bits  1<<w_bits is the requested width of the block size
  @param h_bits  1<<h_bits is the requested height of the block size
  @param aligned If set to 2, the returned sad function will assume src1 and
                 src2 addresses are aligned on the block size.
                 If set to 1, the returned sad function will assume src1 is
                 aligned on the block size.
                 If set to 0, the returned sad function assume no particular
                 alignment.
  @param log_ctx context used for logging, can be NULL

  @return a pointer to the SAD function or NULL in case of error (because of
          invalid parameters)
*/
func AVPixelutilsGetSadFn(wBits int, hBits int, aligned int, logCtx unsafe.Pointer) AVPixelutilsSadFn {
	ret := C.av_pixelutils_get_sad_fn(C.int(wBits), C.int(hBits), C.int(aligned), logCtx)
	return AVPixelutilsSadFn(ret)
}

// --- Function av_get_random_seed ---

// AVGetRandomSeed wraps av_get_random_seed.
/*
  Get a seed to use in conjunction with random functions.
  This function tries to provide a good seed at a best effort bases.
  Its possible to call this function multiple times if more bits are needed.
  It can be quite slow, which is why it should only be used as seed for a faster
  PRNG. The quality of the seed depends on the platform.
*/
func AVGetRandomSeed() uint32 {
	ret := C.av_get_random_seed()
	return uint32(ret)
}

// --- Function av_random_bytes ---

// AVRandomBytes wraps av_random_bytes.
/*
  Generate cryptographically secure random data, i.e. suitable for use as
  encryption keys and similar.

  @param buf buffer into which the random data will be written
  @param len size of buf in bytes

  @retval 0                         success, len bytes of random data was written
                                    into buf
  @retval "a negative AVERROR code" random data could not be generated
*/
func AVRandomBytes(buf unsafe.Pointer, len uint64) (int, error) {
	ret := C.av_random_bytes((*C.uint8_t)(buf), C.size_t(len))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_make_q ---

// AVMakeQ wraps av_make_q.
/*
  Create an AVRational.

  Useful for compilers that do not support compound literals.

  @note The return value is not reduced.
  @see av_reduce()
*/
func AVMakeQ(num int, den int) *AVRational {
	ret := C.av_make_q(C.int(num), C.int(den))
	return &AVRational{value: ret}
}

// --- Function av_cmp_q ---

// AVCmpQ wraps av_cmp_q.
/*
  Compare two rationals.

  @param a First rational
  @param b Second rational

  @return One of the following values:
          - 0 if `a == b`
          - 1 if `a > b`
          - -1 if `a < b`
          - `INT_MIN` if one of the values is of the form `0 / 0`
*/
func AVCmpQ(a *AVRational, b *AVRational) (int, error) {
	ret := C.av_cmp_q(a.value, b.value)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_q2d ---

// AVQ2D wraps av_q2d.
/*
  Convert an AVRational to a `double`.
  @param a AVRational to convert
  @return `a` in floating-point form
  @see av_d2q()
*/
func AVQ2D(a *AVRational) float64 {
	ret := C.av_q2d(a.value)
	return float64(ret)
}

// --- Function av_reduce ---

// AVReduce wraps av_reduce.
/*
  Reduce a fraction.

  This is useful for framerate calculations.

  @param[out] dst_num Destination numerator
  @param[out] dst_den Destination denominator
  @param[in]      num Source numerator
  @param[in]      den Source denominator
  @param[in]      max Maximum allowed values for `dst_num` & `dst_den`
  @return 1 if the operation is exact, 0 otherwise
*/
func AVReduce(dstNum *int, dstDen *int, num int64, den int64, max int64) (int, error) {
	ret := C.av_reduce((*C.int)(unsafe.Pointer(dstNum)), (*C.int)(unsafe.Pointer(dstDen)), C.int64_t(num), C.int64_t(den), C.int64_t(max))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_mul_q ---

// AVMulQ wraps av_mul_q.
/*
  Multiply two rationals.
  @param b First rational
  @param c Second rational
  @return b*c
*/
func AVMulQ(b *AVRational, c *AVRational) *AVRational {
	ret := C.av_mul_q(b.value, c.value)
	return &AVRational{value: ret}
}

// --- Function av_div_q ---

// AVDivQ wraps av_div_q.
/*
  Divide one rational by another.
  @param b First rational
  @param c Second rational
  @return b/c
*/
func AVDivQ(b *AVRational, c *AVRational) *AVRational {
	ret := C.av_div_q(b.value, c.value)
	return &AVRational{value: ret}
}

// --- Function av_add_q ---

// AVAddQ wraps av_add_q.
/*
  Add two rationals.
  @param b First rational
  @param c Second rational
  @return b+c
*/
func AVAddQ(b *AVRational, c *AVRational) *AVRational {
	ret := C.av_add_q(b.value, c.value)
	return &AVRational{value: ret}
}

// --- Function av_sub_q ---

// AVSubQ wraps av_sub_q.
/*
  Subtract one rational from another.
  @param b First rational
  @param c Second rational
  @return b-c
*/
func AVSubQ(b *AVRational, c *AVRational) *AVRational {
	ret := C.av_sub_q(b.value, c.value)
	return &AVRational{value: ret}
}

// --- Function av_inv_q ---

// AVInvQ wraps av_inv_q.
/*
  Invert a rational.
  @param q value
  @return 1 / q
*/
func AVInvQ(q *AVRational) *AVRational {
	ret := C.av_inv_q(q.value)
	return &AVRational{value: ret}
}

// --- Function av_d2q ---

// AVD2Q wraps av_d2q.
/*
  Convert a double precision floating point number to a rational.

  In case of infinity, the returned value is expressed as `{1, 0}` or
  `{-1, 0}` depending on the sign.

  In general rational numbers with |num| <= 1<<26 && |den| <= 1<<26
  can be recovered exactly from their double representation.
  (no exceptions were found within 1B random ones)

  @param d   `double` to convert
  @param max Maximum allowed numerator and denominator
  @return `d` in AVRational form
  @see av_q2d()
*/
func AVD2Q(d float64, max int) *AVRational {
	ret := C.av_d2q(C.double(d), C.int(max))
	return &AVRational{value: ret}
}

// --- Function av_nearer_q ---

// AVNearerQ wraps av_nearer_q.
/*
  Find which of the two rationals is closer to another rational.

  @param q     Rational to be compared against
  @param q1    Rational to be tested
  @param q2    Rational to be tested
  @return One of the following values:
          - 1 if `q1` is nearer to `q` than `q2`
          - -1 if `q2` is nearer to `q` than `q1`
          - 0 if they have the same distance
*/
func AVNearerQ(q *AVRational, q1 *AVRational, q2 *AVRational) (int, error) {
	ret := C.av_nearer_q(q.value, q1.value, q2.value)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_find_nearest_q_idx ---

// av_find_nearest_q_idx skipped due to qList

// --- Function av_q2intfloat ---

// AVQ2Intfloat wraps av_q2intfloat.
/*
  Convert an AVRational to a IEEE 32-bit `float` expressed in fixed-point
  format.

  @param q Rational to be converted
  @return Equivalent floating-point value, expressed as an unsigned 32-bit
          integer.
  @note The returned value is platform-indepedant.
*/
func AVQ2Intfloat(q *AVRational) uint32 {
	ret := C.av_q2intfloat(q.value)
	return uint32(ret)
}

// --- Function av_gcd_q ---

// AVGcdQ wraps av_gcd_q.
/*
  Return the best rational so that a and b are multiple of it.
  If the resulting denominator is larger than max_den, return def.
*/
func AVGcdQ(a *AVRational, b *AVRational, maxDen int, def *AVRational) *AVRational {
	ret := C.av_gcd_q(a.value, b.value, C.int(maxDen), def.value)
	return &AVRational{value: ret}
}

// --- Function av_rc4_alloc ---

// AVRc4Alloc wraps av_rc4_alloc.
//
//	Allocate an AVRC4 context.
func AVRc4Alloc() *AVRC4 {
	ret := C.av_rc4_alloc()
	var retMapped *AVRC4
	if ret != nil {
		retMapped = &AVRC4{ptr: ret}
	}
	return retMapped
}

// --- Function av_rc4_init ---

// AVRc4Init wraps av_rc4_init.
/*
  @brief Initializes an AVRC4 context.

  @param d pointer to the AVRC4 context
  @param key buffer containing the key
  @param key_bits must be a multiple of 8
  @param decrypt 0 for encryption, 1 for decryption, currently has no effect
  @return zero on success, negative value otherwise
*/
func AVRc4Init(d *AVRC4, key unsafe.Pointer, keyBits int, decrypt int) (int, error) {
	var tmpd *C.AVRC4
	if d != nil {
		tmpd = d.ptr
	}
	ret := C.av_rc4_init(tmpd, (*C.uint8_t)(key), C.int(keyBits), C.int(decrypt))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_rc4_crypt ---

// AVRc4Crypt wraps av_rc4_crypt.
/*
  @brief Encrypts / decrypts using the RC4 algorithm.

  @param d pointer to the AVRC4 context
  @param count number of bytes
  @param dst destination array, can be equal to src
  @param src source array, can be equal to dst, may be NULL
  @param iv not (yet) used for RC4, should be NULL
  @param decrypt 0 for encryption, 1 for decryption, not (yet) used
*/
func AVRc4Crypt(d *AVRC4, dst unsafe.Pointer, src unsafe.Pointer, count int, iv unsafe.Pointer, decrypt int) {
	var tmpd *C.AVRC4
	if d != nil {
		tmpd = d.ptr
	}
	C.av_rc4_crypt(tmpd, (*C.uint8_t)(dst), (*C.uint8_t)(src), C.int(count), (*C.uint8_t)(iv), C.int(decrypt))
}

// --- Function av_ripemd_alloc ---

// AVRipemdAlloc wraps av_ripemd_alloc.
//
//	Allocate an AVRIPEMD context.
func AVRipemdAlloc() *AVRIPEMD {
	ret := C.av_ripemd_alloc()
	var retMapped *AVRIPEMD
	if ret != nil {
		retMapped = &AVRIPEMD{ptr: ret}
	}
	return retMapped
}

// --- Function av_ripemd_init ---

// AVRipemdInit wraps av_ripemd_init.
/*
  Initialize RIPEMD hashing.

  @param context pointer to the function context (of size av_ripemd_size)
  @param bits    number of bits in digest (128, 160, 256 or 320 bits)
  @return        zero if initialization succeeded, -1 otherwise
*/
func AVRipemdInit(context *AVRIPEMD, bits int) (int, error) {
	var tmpcontext *C.struct_AVRIPEMD
	if context != nil {
		tmpcontext = context.ptr
	}
	ret := C.av_ripemd_init(tmpcontext, C.int(bits))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_ripemd_update ---

// AVRipemdUpdate wraps av_ripemd_update.
/*
  Update hash value.

  @param context hash function context
  @param data    input data to update hash with
  @param len     input data length
*/
func AVRipemdUpdate(context *AVRIPEMD, data unsafe.Pointer, len uint64) {
	var tmpcontext *C.struct_AVRIPEMD
	if context != nil {
		tmpcontext = context.ptr
	}
	C.av_ripemd_update(tmpcontext, (*C.uint8_t)(data), C.size_t(len))
}

// --- Function av_ripemd_final ---

// AVRipemdFinal wraps av_ripemd_final.
/*
  Finish hashing and output digest value.

  @param context hash function context
  @param digest  buffer where output digest value is stored
*/
func AVRipemdFinal(context *AVRIPEMD, digest unsafe.Pointer) {
	var tmpcontext *C.struct_AVRIPEMD
	if context != nil {
		tmpcontext = context.ptr
	}
	C.av_ripemd_final(tmpcontext, (*C.uint8_t)(digest))
}

// --- Function av_get_sample_fmt_name ---

// AVGetSampleFmtName wraps av_get_sample_fmt_name.
/*
  Return the name of sample_fmt, or NULL if sample_fmt is not
  recognized.
*/
func AVGetSampleFmtName(sampleFmt AVSampleFormat) *CStr {
	ret := C.av_get_sample_fmt_name(C.enum_AVSampleFormat(sampleFmt))
	return wrapCStr(ret)
}

// --- Function av_get_sample_fmt ---

// AVGetSampleFmt wraps av_get_sample_fmt.
/*
  Return a sample format corresponding to name, or AV_SAMPLE_FMT_NONE
  on error.
*/
func AVGetSampleFmt(name *CStr) AVSampleFormat {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_get_sample_fmt(tmpname)
	return AVSampleFormat(ret)
}

// --- Function av_get_alt_sample_fmt ---

// AVGetAltSampleFmt wraps av_get_alt_sample_fmt.
/*
  Return the planar<->packed alternative form of the given sample format, or
  AV_SAMPLE_FMT_NONE on error. If the passed sample_fmt is already in the
  requested planar/packed format, the format returned is the same as the
  input.
*/
func AVGetAltSampleFmt(sampleFmt AVSampleFormat, planar int) AVSampleFormat {
	ret := C.av_get_alt_sample_fmt(C.enum_AVSampleFormat(sampleFmt), C.int(planar))
	return AVSampleFormat(ret)
}

// --- Function av_get_packed_sample_fmt ---

// AVGetPackedSampleFmt wraps av_get_packed_sample_fmt.
/*
  Get the packed alternative form of the given sample format.

  If the passed sample_fmt is already in packed format, the format returned is
  the same as the input.

  @return  the packed alternative form of the given sample format or
  AV_SAMPLE_FMT_NONE on error.
*/
func AVGetPackedSampleFmt(sampleFmt AVSampleFormat) AVSampleFormat {
	ret := C.av_get_packed_sample_fmt(C.enum_AVSampleFormat(sampleFmt))
	return AVSampleFormat(ret)
}

// --- Function av_get_planar_sample_fmt ---

// AVGetPlanarSampleFmt wraps av_get_planar_sample_fmt.
/*
  Get the planar alternative form of the given sample format.

  If the passed sample_fmt is already in planar format, the format returned is
  the same as the input.

  @return  the planar alternative form of the given sample format or
  AV_SAMPLE_FMT_NONE on error.
*/
func AVGetPlanarSampleFmt(sampleFmt AVSampleFormat) AVSampleFormat {
	ret := C.av_get_planar_sample_fmt(C.enum_AVSampleFormat(sampleFmt))
	return AVSampleFormat(ret)
}

// --- Function av_get_sample_fmt_string ---

// AVGetSampleFmtString wraps av_get_sample_fmt_string.
/*
  Generate a string corresponding to the sample format with
  sample_fmt, or a header if sample_fmt is negative.

  @param buf the buffer where to write the string
  @param buf_size the size of buf
  @param sample_fmt the number of the sample format to print the
  corresponding info string, or a negative value to print the
  corresponding header.
  @return the pointer to the filled buffer or NULL if sample_fmt is
  unknown or in case of other errors
*/
func AVGetSampleFmtString(buf *CStr, bufSize int, sampleFmt AVSampleFormat) *CStr {
	var tmpbuf *C.char
	if buf != nil {
		tmpbuf = buf.ptr
	}
	ret := C.av_get_sample_fmt_string(tmpbuf, C.int(bufSize), C.enum_AVSampleFormat(sampleFmt))
	return wrapCStr(ret)
}

// --- Function av_get_bytes_per_sample ---

// AVGetBytesPerSample wraps av_get_bytes_per_sample.
/*
  Return number of bytes per sample.

  @param sample_fmt the sample format
  @return number of bytes per sample or zero if unknown for the given
  sample format
*/
func AVGetBytesPerSample(sampleFmt AVSampleFormat) (int, error) {
	ret := C.av_get_bytes_per_sample(C.enum_AVSampleFormat(sampleFmt))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_sample_fmt_is_planar ---

// AVSampleFmtIsPlanar wraps av_sample_fmt_is_planar.
/*
  Check if the sample format is planar.

  @param sample_fmt the sample format to inspect
  @return 1 if the sample format is planar, 0 if it is interleaved
*/
func AVSampleFmtIsPlanar(sampleFmt AVSampleFormat) (int, error) {
	ret := C.av_sample_fmt_is_planar(C.enum_AVSampleFormat(sampleFmt))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_samples_get_buffer_size ---

// AVSamplesGetBufferSize wraps av_samples_get_buffer_size.
/*
  Get the required buffer size for the given audio parameters.

  @param[out] linesize calculated linesize, may be NULL
  @param nb_channels   the number of channels
  @param nb_samples    the number of samples in a single channel
  @param sample_fmt    the sample format
  @param align         buffer size alignment (0 = default, 1 = no alignment)
  @return              required buffer size, or negative error code on failure
*/
func AVSamplesGetBufferSize(linesize *int, nbChannels int, nbSamples int, sampleFmt AVSampleFormat, align int) (int, error) {
	ret := C.av_samples_get_buffer_size((*C.int)(unsafe.Pointer(linesize)), C.int(nbChannels), C.int(nbSamples), C.enum_AVSampleFormat(sampleFmt), C.int(align))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_samples_fill_arrays ---

// av_samples_fill_arrays skipped due to audioData

// --- Function av_samples_alloc ---

// av_samples_alloc skipped due to audioData

// --- Function av_samples_alloc_array_and_samples ---

// av_samples_alloc_array_and_samples skipped due to audioData

// --- Function av_samples_copy ---

// av_samples_copy skipped due to dst

// --- Function av_samples_set_silence ---

// av_samples_set_silence skipped due to audioData

// --- Function av_sha_alloc ---

// AVShaAlloc wraps av_sha_alloc.
//
//	Allocate an AVSHA context.
func AVShaAlloc() *AVSHA {
	ret := C.av_sha_alloc()
	var retMapped *AVSHA
	if ret != nil {
		retMapped = &AVSHA{ptr: ret}
	}
	return retMapped
}

// --- Function av_sha_init ---

// AVShaInit wraps av_sha_init.
/*
  Initialize SHA-1 or SHA-2 hashing.

  @param context pointer to the function context (of size av_sha_size)
  @param bits    number of bits in digest (SHA-1 - 160 bits, SHA-2 224 or 256 bits)
  @return        zero if initialization succeeded, -1 otherwise
*/
func AVShaInit(context *AVSHA, bits int) (int, error) {
	var tmpcontext *C.struct_AVSHA
	if context != nil {
		tmpcontext = context.ptr
	}
	ret := C.av_sha_init(tmpcontext, C.int(bits))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_sha_update ---

// AVShaUpdate wraps av_sha_update.
/*
  Update hash value.

  @param ctx     hash function context
  @param data    input data to update hash with
  @param len     input data length
*/
func AVShaUpdate(ctx *AVSHA, data unsafe.Pointer, len uint64) {
	var tmpctx *C.struct_AVSHA
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	C.av_sha_update(tmpctx, (*C.uint8_t)(data), C.size_t(len))
}

// --- Function av_sha_final ---

// AVShaFinal wraps av_sha_final.
/*
  Finish hashing and output digest value.

  @param context hash function context
  @param digest  buffer where output digest value is stored
*/
func AVShaFinal(context *AVSHA, digest unsafe.Pointer) {
	var tmpcontext *C.struct_AVSHA
	if context != nil {
		tmpcontext = context.ptr
	}
	C.av_sha_final(tmpcontext, (*C.uint8_t)(digest))
}

// --- Function av_sha512_alloc ---

// AVSha512Alloc wraps av_sha512_alloc.
//
//	Allocate an AVSHA512 context.
func AVSha512Alloc() *AVSHA512 {
	ret := C.av_sha512_alloc()
	var retMapped *AVSHA512
	if ret != nil {
		retMapped = &AVSHA512{ptr: ret}
	}
	return retMapped
}

// --- Function av_sha512_init ---

// AVSha512Init wraps av_sha512_init.
/*
  Initialize SHA-2 512 hashing.

  @param context pointer to the function context (of size av_sha512_size)
  @param bits    number of bits in digest (224, 256, 384 or 512 bits)
  @return        zero if initialization succeeded, -1 otherwise
*/
func AVSha512Init(context *AVSHA512, bits int) (int, error) {
	var tmpcontext *C.struct_AVSHA512
	if context != nil {
		tmpcontext = context.ptr
	}
	ret := C.av_sha512_init(tmpcontext, C.int(bits))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_sha512_update ---

// AVSha512Update wraps av_sha512_update.
/*
  Update hash value.

  @param context hash function context
  @param data    input data to update hash with
  @param len     input data length
*/
func AVSha512Update(context *AVSHA512, data unsafe.Pointer, len uint64) {
	var tmpcontext *C.struct_AVSHA512
	if context != nil {
		tmpcontext = context.ptr
	}
	C.av_sha512_update(tmpcontext, (*C.uint8_t)(data), C.size_t(len))
}

// --- Function av_sha512_final ---

// AVSha512Final wraps av_sha512_final.
/*
  Finish hashing and output digest value.

  @param context hash function context
  @param digest  buffer where output digest value is stored
*/
func AVSha512Final(context *AVSHA512, digest unsafe.Pointer) {
	var tmpcontext *C.struct_AVSHA512
	if context != nil {
		tmpcontext = context.ptr
	}
	C.av_sha512_final(tmpcontext, (*C.uint8_t)(digest))
}

// --- Function av_spherical_alloc ---

// AVSphericalAlloc wraps av_spherical_alloc.
/*
  Allocate a AVSphericalVideo structure and initialize its fields to default
  values.

  @return the newly allocated struct or NULL on failure
*/
func AVSphericalAlloc(size *uint64) *AVSphericalMapping {
	ret := C.av_spherical_alloc((*C.size_t)(unsafe.Pointer(size)))
	var retMapped *AVSphericalMapping
	if ret != nil {
		retMapped = &AVSphericalMapping{ptr: ret}
	}
	return retMapped
}

// --- Function av_spherical_tile_bounds ---

// av_spherical_tile_bounds skipped due to left (non-output primitive pointer)

// --- Function av_spherical_projection_name ---

// AVSphericalProjectionName wraps av_spherical_projection_name.
/*
  Provide a human-readable name of a given AVSphericalProjection.

  @param projection The input AVSphericalProjection.

  @return The name of the AVSphericalProjection, or "unknown".
*/
func AVSphericalProjectionName(projection AVSphericalProjection) *CStr {
	ret := C.av_spherical_projection_name(C.enum_AVSphericalProjection(projection))
	return wrapCStr(ret)
}

// --- Function av_spherical_from_name ---

// AVSphericalFromName wraps av_spherical_from_name.
/*
  Get the AVSphericalProjection form a human-readable name.

  @param name The input string.

  @return The AVSphericalProjection value, or -1 if not found.
*/
func AVSphericalFromName(name *CStr) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_spherical_from_name(tmpname)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_stereo3d_alloc ---

// AVStereo3DAlloc wraps av_stereo3d_alloc.
/*
  Allocate an AVStereo3D structure and set its fields to default values.
  The resulting struct can be freed using av_freep().

  @return An AVStereo3D filled with default values or NULL on failure.
*/
func AVStereo3DAlloc() *AVStereo3D {
	ret := C.av_stereo3d_alloc()
	var retMapped *AVStereo3D
	if ret != nil {
		retMapped = &AVStereo3D{ptr: ret}
	}
	return retMapped
}

// --- Function av_stereo3d_alloc_size ---

// AVStereo3DAllocSize wraps av_stereo3d_alloc_size.
/*
  Allocate an AVStereo3D structure and set its fields to default values.
  The resulting struct can be freed using av_freep().

  @return An AVStereo3D filled with default values or NULL on failure.
*/
func AVStereo3DAllocSize(size *uint64) *AVStereo3D {
	ret := C.av_stereo3d_alloc_size((*C.size_t)(unsafe.Pointer(size)))
	var retMapped *AVStereo3D
	if ret != nil {
		retMapped = &AVStereo3D{ptr: ret}
	}
	return retMapped
}

// --- Function av_stereo3d_create_side_data ---

// AVStereo3DCreateSideData wraps av_stereo3d_create_side_data.
/*
  Allocate a complete AVFrameSideData and add it to the frame.

  @param frame The frame which side data is added to.

  @return The AVStereo3D structure to be filled by caller.
*/
func AVStereo3DCreateSideData(frame *AVFrame) *AVStereo3D {
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.av_stereo3d_create_side_data(tmpframe)
	var retMapped *AVStereo3D
	if ret != nil {
		retMapped = &AVStereo3D{ptr: ret}
	}
	return retMapped
}

// --- Function av_stereo3d_type_name ---

// AVStereo3DTypeName wraps av_stereo3d_type_name.
/*
  Provide a human-readable name of a given stereo3d type.

  @param type The input stereo3d type value.

  @return The name of the stereo3d value, or "unknown".
*/
func AVStereo3DTypeName(_type uint) *CStr {
	ret := C.av_stereo3d_type_name(C.uint(_type))
	return wrapCStr(ret)
}

// --- Function av_stereo3d_from_name ---

// AVStereo3DFromName wraps av_stereo3d_from_name.
/*
  Get the AVStereo3DType form a human-readable name.

  @param name The input string.

  @return The AVStereo3DType value, or -1 if not found.
*/
func AVStereo3DFromName(name *CStr) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_stereo3d_from_name(tmpname)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_stereo3d_view_name ---

// AVStereo3DViewName wraps av_stereo3d_view_name.
/*
  Provide a human-readable name of a given stereo3d view.

  @param type The input stereo3d view value.

  @return The name of the stereo3d view value, or "unknown".
*/
func AVStereo3DViewName(view uint) *CStr {
	ret := C.av_stereo3d_view_name(C.uint(view))
	return wrapCStr(ret)
}

// --- Function av_stereo3d_view_from_name ---

// AVStereo3DViewFromName wraps av_stereo3d_view_from_name.
/*
  Get the AVStereo3DView form a human-readable name.

  @param name The input string.

  @return The AVStereo3DView value, or -1 if not found.
*/
func AVStereo3DViewFromName(name *CStr) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_stereo3d_view_from_name(tmpname)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_stereo3d_primary_eye_name ---

// AVStereo3DPrimaryEyeName wraps av_stereo3d_primary_eye_name.
/*
  Provide a human-readable name of a given stereo3d primary eye.

  @param type The input stereo3d primary eye value.

  @return The name of the stereo3d primary eye value, or "unknown".
*/
func AVStereo3DPrimaryEyeName(eye uint) *CStr {
	ret := C.av_stereo3d_primary_eye_name(C.uint(eye))
	return wrapCStr(ret)
}

// --- Function av_stereo3d_primary_eye_from_name ---

// AVStereo3DPrimaryEyeFromName wraps av_stereo3d_primary_eye_from_name.
/*
  Get the AVStereo3DPrimaryEye form a human-readable name.

  @param name The input string.

  @return The AVStereo3DPrimaryEye value, or -1 if not found.
*/
func AVStereo3DPrimaryEyeFromName(name *CStr) (int, error) {
	var tmpname *C.char
	if name != nil {
		tmpname = name.ptr
	}
	ret := C.av_stereo3d_primary_eye_from_name(tmpname)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_tdrdi_get_display ---

// AVTdrdiGetDisplay wraps av_tdrdi_get_display.
func AVTdrdiGetDisplay(tdrdi *AV3DReferenceDisplaysInfo, idx uint) *AV3DReferenceDisplay {
	var tmptdrdi *C.AV3DReferenceDisplaysInfo
	if tdrdi != nil {
		tmptdrdi = tdrdi.ptr
	}
	ret := C.av_tdrdi_get_display(tmptdrdi, C.uint(idx))
	var retMapped *AV3DReferenceDisplay
	if ret != nil {
		retMapped = &AV3DReferenceDisplay{ptr: ret}
	}
	return retMapped
}

// --- Function av_tdrdi_alloc ---

// AVTdrdiAlloc wraps av_tdrdi_alloc.
/*
  Allocate a AV3DReferenceDisplaysInfo structure and initialize its fields to default
  values.

  @return the newly allocated struct or NULL on failure
*/
func AVTdrdiAlloc(nbDisplays uint, size *uint64) *AV3DReferenceDisplaysInfo {
	ret := C.av_tdrdi_alloc(C.uint(nbDisplays), (*C.size_t)(unsafe.Pointer(size)))
	var retMapped *AV3DReferenceDisplaysInfo
	if ret != nil {
		retMapped = &AV3DReferenceDisplaysInfo{ptr: ret}
	}
	return retMapped
}

// --- Function av_tea_alloc ---

// AVTeaAlloc wraps av_tea_alloc.
/*
  Allocate an AVTEA context
  To free the struct: av_free(ptr)
*/
func AVTeaAlloc() *AVTEA {
	ret := C.av_tea_alloc()
	var retMapped *AVTEA
	if ret != nil {
		retMapped = &AVTEA{ptr: ret}
	}
	return retMapped
}

// --- Function av_tea_init ---

// av_tea_init skipped due to const array param key

// --- Function av_tea_crypt ---

// AVTeaCrypt wraps av_tea_crypt.
/*
  Encrypt or decrypt a buffer using a previously initialized context.

  @param ctx an AVTEA context
  @param dst destination array, can be equal to src
  @param src source array, can be equal to dst
  @param count number of 8 byte blocks
  @param iv initialization vector for CBC mode, if NULL then ECB will be used
  @param decrypt 0 for encryption, 1 for decryption
*/
func AVTeaCrypt(ctx *AVTEA, dst unsafe.Pointer, src unsafe.Pointer, count int, iv unsafe.Pointer, decrypt int) {
	var tmpctx *C.struct_AVTEA
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	C.av_tea_crypt(tmpctx, (*C.uint8_t)(dst), (*C.uint8_t)(src), C.int(count), (*C.uint8_t)(iv), C.int(decrypt))
}

// --- Function av_thread_message_queue_alloc ---

// AVThreadMessageQueueAlloc wraps av_thread_message_queue_alloc.
/*
  Allocate a new message queue.

  @param mq      pointer to the message queue
  @param nelem   maximum number of elements in the queue
  @param elsize  size of each element in the queue
  @return  >=0 for success; <0 for error, in particular AVERROR(ENOSYS) if
           lavu was built without thread support
*/
func AVThreadMessageQueueAlloc(mq **AVThreadMessageQueue, nelem uint, elsize uint) (int, error) {
	var ptrmq **C.AVThreadMessageQueue
	var tmpmq *C.AVThreadMessageQueue
	var oldTmpmq *C.AVThreadMessageQueue
	if mq != nil {
		innermq := *mq
		if innermq != nil {
			tmpmq = innermq.ptr
			oldTmpmq = tmpmq
		}
		ptrmq = &tmpmq
	}
	ret := C.av_thread_message_queue_alloc(ptrmq, C.uint(nelem), C.uint(elsize))
	if tmpmq != oldTmpmq && mq != nil {
		if tmpmq != nil {
			*mq = &AVThreadMessageQueue{ptr: tmpmq}
		} else {
			*mq = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function av_thread_message_queue_free ---

// AVThreadMessageQueueFree wraps av_thread_message_queue_free.
/*
  Free a message queue.

  The message queue must no longer be in use by another thread.
*/
func AVThreadMessageQueueFree(mq **AVThreadMessageQueue) {
	var ptrmq **C.AVThreadMessageQueue
	var tmpmq *C.AVThreadMessageQueue
	var oldTmpmq *C.AVThreadMessageQueue
	if mq != nil {
		innermq := *mq
		if innermq != nil {
			tmpmq = innermq.ptr
			oldTmpmq = tmpmq
		}
		ptrmq = &tmpmq
	}
	C.av_thread_message_queue_free(ptrmq)
	if tmpmq != oldTmpmq && mq != nil {
		if tmpmq != nil {
			*mq = &AVThreadMessageQueue{ptr: tmpmq}
		} else {
			*mq = nil
		}
	}
}

// --- Function av_thread_message_queue_send ---

// AVThreadMessageQueueSend wraps av_thread_message_queue_send.
//
//	Send a message on the queue.
func AVThreadMessageQueueSend(mq *AVThreadMessageQueue, msg unsafe.Pointer, flags uint) (int, error) {
	var tmpmq *C.AVThreadMessageQueue
	if mq != nil {
		tmpmq = mq.ptr
	}
	ret := C.av_thread_message_queue_send(tmpmq, msg, C.uint(flags))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_thread_message_queue_recv ---

// AVThreadMessageQueueRecv wraps av_thread_message_queue_recv.
//
//	Receive a message from the queue.
func AVThreadMessageQueueRecv(mq *AVThreadMessageQueue, msg unsafe.Pointer, flags uint) (int, error) {
	var tmpmq *C.AVThreadMessageQueue
	if mq != nil {
		tmpmq = mq.ptr
	}
	ret := C.av_thread_message_queue_recv(tmpmq, msg, C.uint(flags))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_thread_message_queue_set_err_send ---

// AVThreadMessageQueueSetErrSend wraps av_thread_message_queue_set_err_send.
/*
  Set the sending error code.

  If the error code is set to non-zero, av_thread_message_queue_send() will
  return it immediately. Conventional values, such as AVERROR_EOF or
  AVERROR(EAGAIN), can be used to cause the sending thread to stop or
  suspend its operation.
*/
func AVThreadMessageQueueSetErrSend(mq *AVThreadMessageQueue, err int) {
	var tmpmq *C.AVThreadMessageQueue
	if mq != nil {
		tmpmq = mq.ptr
	}
	C.av_thread_message_queue_set_err_send(tmpmq, C.int(err))
}

// --- Function av_thread_message_queue_set_err_recv ---

// AVThreadMessageQueueSetErrRecv wraps av_thread_message_queue_set_err_recv.
/*
  Set the receiving error code.

  If the error code is set to non-zero, av_thread_message_queue_recv() will
  return it immediately when there are no longer available messages.
  Conventional values, such as AVERROR_EOF or AVERROR(EAGAIN), can be used
  to cause the receiving thread to stop or suspend its operation.
*/
func AVThreadMessageQueueSetErrRecv(mq *AVThreadMessageQueue, err int) {
	var tmpmq *C.AVThreadMessageQueue
	if mq != nil {
		tmpmq = mq.ptr
	}
	C.av_thread_message_queue_set_err_recv(tmpmq, C.int(err))
}

// --- Function av_thread_message_queue_set_free_func ---

// av_thread_message_queue_set_free_func skipped due to free_func.

// --- Function av_thread_message_queue_nb_elems ---

// AVThreadMessageQueueNbElems wraps av_thread_message_queue_nb_elems.
/*
  Return the current number of messages in the queue.

  @return the current number of messages or AVERROR(ENOSYS) if lavu was built
          without thread support
*/
func AVThreadMessageQueueNbElems(mq *AVThreadMessageQueue) (int, error) {
	var tmpmq *C.AVThreadMessageQueue
	if mq != nil {
		tmpmq = mq.ptr
	}
	ret := C.av_thread_message_queue_nb_elems(tmpmq)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_thread_message_flush ---

// AVThreadMessageFlush wraps av_thread_message_flush.
/*
  Flush the message queue

  This function is mostly equivalent to reading and free-ing every message
  except that it will be done in a single operation (no lock/unlock between
  reads).
*/
func AVThreadMessageFlush(mq *AVThreadMessageQueue) {
	var tmpmq *C.AVThreadMessageQueue
	if mq != nil {
		tmpmq = mq.ptr
	}
	C.av_thread_message_flush(tmpmq)
}

// --- Function av_gettime ---

// AVGettime wraps av_gettime.
//
//	Get the current time in microseconds.
func AVGettime() int64 {
	ret := C.av_gettime()
	return int64(ret)
}

// --- Function av_gettime_relative ---

// AVGettimeRelative wraps av_gettime_relative.
/*
  Get the current time in microseconds since some unspecified starting point.
  On platforms that support it, the time comes from a monotonic clock
  This property makes this time source ideal for measuring relative time.
  The returned values may not be monotonic on platforms where a monotonic
  clock is not available.
*/
func AVGettimeRelative() int64 {
	ret := C.av_gettime_relative()
	return int64(ret)
}

// --- Function av_gettime_relative_is_monotonic ---

// AVGettimeRelativeIsMonotonic wraps av_gettime_relative_is_monotonic.
/*
  Indicates with a boolean result if the av_gettime_relative() time source
  is monotonic.
*/
func AVGettimeRelativeIsMonotonic() (int, error) {
	ret := C.av_gettime_relative_is_monotonic()
	return int(ret), WrapErr(int(ret))
}

// --- Function av_usleep ---

// AVUsleep wraps av_usleep.
/*
  Sleep for a period of time.  Although the duration is expressed in
  microseconds, the actual delay may be rounded to the precision of the
  system timer.

  @param  usec Number of microseconds to sleep.
  @return zero on success or (negative) error code.
*/
func AVUsleep(usec uint) (int, error) {
	ret := C.av_usleep(C.uint(usec))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_timecode_adjust_ntsc_framenum2 ---

// AVTimecodeAdjustNtscFramenum2 wraps av_timecode_adjust_ntsc_framenum2.
/*
  Adjust frame number for NTSC drop frame time code.

  @param framenum frame number to adjust
  @param fps      frame per second, multiples of 30
  @return         adjusted frame number
  @warning        adjustment is only valid for multiples of NTSC 29.97
*/
func AVTimecodeAdjustNtscFramenum2(framenum int, fps int) (int, error) {
	ret := C.av_timecode_adjust_ntsc_framenum2(C.int(framenum), C.int(fps))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_timecode_get_smpte_from_framenum ---

// AVTimecodeGetSmpteFromFramenum wraps av_timecode_get_smpte_from_framenum.
/*
  Convert frame number to SMPTE 12M binary representation.

  @param tc       timecode data correctly initialized
  @param framenum frame number
  @return         the SMPTE binary representation

  See SMPTE ST 314M-2005 Sec 4.4.2.2.1 "Time code pack (TC)"
  the format description as follows:
  bits 0-5:   hours, in BCD(6bits)
  bits 6:     BGF1
  bits 7:     BGF2 (NTSC) or FIELD (PAL)
  bits 8-14:  minutes, in BCD(7bits)
  bits 15:    BGF0 (NTSC) or BGF2 (PAL)
  bits 16-22: seconds, in BCD(7bits)
  bits 23:    FIELD (NTSC) or BGF0 (PAL)
  bits 24-29: frames, in BCD(6bits)
  bits 30:    drop  frame flag (0: non drop,    1: drop)
  bits 31:    color frame flag (0: unsync mode, 1: sync mode)
  @note BCD numbers (6 or 7 bits): 4 or 5 lower bits for units, 2 higher bits for tens.
  @note Frame number adjustment is automatically done in case of drop timecode,
        you do NOT have to call av_timecode_adjust_ntsc_framenum2().
  @note The frame number is relative to tc->start.
  @note Color frame (CF) and binary group flags (BGF) bits are set to zero.
*/
func AVTimecodeGetSmpteFromFramenum(tc *AVTimecode, framenum int) uint32 {
	var tmptc *C.AVTimecode
	if tc != nil {
		tmptc = tc.ptr
	}
	ret := C.av_timecode_get_smpte_from_framenum(tmptc, C.int(framenum))
	return uint32(ret)
}

// --- Function av_timecode_get_smpte ---

// AVTimecodeGetSmpte wraps av_timecode_get_smpte.
/*
  Convert sei info to SMPTE 12M binary representation.

  @param rate     frame rate in rational form
  @param drop     drop flag
  @param hh       hour
  @param mm       minute
  @param ss       second
  @param ff       frame number
  @return         the SMPTE binary representation
*/
func AVTimecodeGetSmpte(rate *AVRational, drop int, hh int, mm int, ss int, ff int) uint32 {
	ret := C.av_timecode_get_smpte(rate.value, C.int(drop), C.int(hh), C.int(mm), C.int(ss), C.int(ff))
	return uint32(ret)
}

// --- Function av_timecode_make_string ---

// AVTimecodeMakeString wraps av_timecode_make_string.
/*
  Load timecode string in buf.

  @param tc       timecode data correctly initialized
  @param buf      destination buffer, must be at least AV_TIMECODE_STR_SIZE long
  @param framenum frame number
  @return         the buf parameter

  @note Timecode representation can be a negative timecode and have more than
        24 hours, but will only be honored if the flags are correctly set.
  @note The frame number is relative to tc->start.
*/
func AVTimecodeMakeString(tc *AVTimecode, buf *CStr, framenum int) *CStr {
	var tmptc *C.AVTimecode
	if tc != nil {
		tmptc = tc.ptr
	}
	var tmpbuf *C.char
	if buf != nil {
		tmpbuf = buf.ptr
	}
	ret := C.av_timecode_make_string(tmptc, tmpbuf, C.int(framenum))
	return wrapCStr(ret)
}

// --- Function av_timecode_make_smpte_tc_string2 ---

// AVTimecodeMakeSmpteTcString2 wraps av_timecode_make_smpte_tc_string2.
/*
  Get the timecode string from the SMPTE timecode format.

  In contrast to av_timecode_make_smpte_tc_string this function supports 50/60
  fps timecodes by using the field bit.

  @param buf        destination buffer, must be at least AV_TIMECODE_STR_SIZE long
  @param rate       frame rate of the timecode
  @param tcsmpte    the 32-bit SMPTE timecode
  @param prevent_df prevent the use of a drop flag when it is known the DF bit
                    is arbitrary
  @param skip_field prevent the use of a field flag when it is known the field
                    bit is arbitrary (e.g. because it is used as PC flag)
  @return           the buf parameter
*/
func AVTimecodeMakeSmpteTcString2(buf *CStr, rate *AVRational, tcsmpte uint32, preventDf int, skipField int) *CStr {
	var tmpbuf *C.char
	if buf != nil {
		tmpbuf = buf.ptr
	}
	ret := C.av_timecode_make_smpte_tc_string2(tmpbuf, rate.value, C.uint32_t(tcsmpte), C.int(preventDf), C.int(skipField))
	return wrapCStr(ret)
}

// --- Function av_timecode_make_smpte_tc_string ---

// AVTimecodeMakeSmpteTcString wraps av_timecode_make_smpte_tc_string.
/*
  Get the timecode string from the SMPTE timecode format.

  @param buf        destination buffer, must be at least AV_TIMECODE_STR_SIZE long
  @param tcsmpte    the 32-bit SMPTE timecode
  @param prevent_df prevent the use of a drop flag when it is known the DF bit
                    is arbitrary
  @return           the buf parameter
*/
func AVTimecodeMakeSmpteTcString(buf *CStr, tcsmpte uint32, preventDf int) *CStr {
	var tmpbuf *C.char
	if buf != nil {
		tmpbuf = buf.ptr
	}
	ret := C.av_timecode_make_smpte_tc_string(tmpbuf, C.uint32_t(tcsmpte), C.int(preventDf))
	return wrapCStr(ret)
}

// --- Function av_timecode_make_mpeg_tc_string ---

// AVTimecodeMakeMpegTcString wraps av_timecode_make_mpeg_tc_string.
/*
  Get the timecode string from the 25-bit timecode format (MPEG GOP format).

  @param buf     destination buffer, must be at least AV_TIMECODE_STR_SIZE long
  @param tc25bit the 25-bits timecode
  @return        the buf parameter
*/
func AVTimecodeMakeMpegTcString(buf *CStr, tc25Bit uint32) *CStr {
	var tmpbuf *C.char
	if buf != nil {
		tmpbuf = buf.ptr
	}
	ret := C.av_timecode_make_mpeg_tc_string(tmpbuf, C.uint32_t(tc25Bit))
	return wrapCStr(ret)
}

// --- Function av_timecode_init ---

// AVTimecodeInit wraps av_timecode_init.
/*
  Init a timecode struct with the passed parameters.

  @param tc          pointer to an allocated AVTimecode
  @param rate        frame rate in rational form
  @param flags       miscellaneous flags such as drop frame, +24 hours, ...
                     (see AVTimecodeFlag)
  @param frame_start the first frame number
  @param log_ctx     a pointer to an arbitrary struct of which the first field
                     is a pointer to an AVClass struct (used for av_log)
  @return            0 on success, AVERROR otherwise
*/
func AVTimecodeInit(tc *AVTimecode, rate *AVRational, flags int, frameStart int, logCtx unsafe.Pointer) (int, error) {
	var tmptc *C.AVTimecode
	if tc != nil {
		tmptc = tc.ptr
	}
	ret := C.av_timecode_init(tmptc, rate.value, C.int(flags), C.int(frameStart), logCtx)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_timecode_init_from_components ---

// AVTimecodeInitFromComponents wraps av_timecode_init_from_components.
/*
  Init a timecode struct from the passed timecode components.

  @param tc          pointer to an allocated AVTimecode
  @param rate        frame rate in rational form
  @param flags       miscellaneous flags such as drop frame, +24 hours, ...
                     (see AVTimecodeFlag)
  @param hh          hours
  @param mm          minutes
  @param ss          seconds
  @param ff          frames
  @param log_ctx     a pointer to an arbitrary struct of which the first field
                     is a pointer to an AVClass struct (used for av_log)
  @return            0 on success, AVERROR otherwise
*/
func AVTimecodeInitFromComponents(tc *AVTimecode, rate *AVRational, flags int, hh int, mm int, ss int, ff int, logCtx unsafe.Pointer) (int, error) {
	var tmptc *C.AVTimecode
	if tc != nil {
		tmptc = tc.ptr
	}
	ret := C.av_timecode_init_from_components(tmptc, rate.value, C.int(flags), C.int(hh), C.int(mm), C.int(ss), C.int(ff), logCtx)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_timecode_init_from_string ---

// AVTimecodeInitFromString wraps av_timecode_init_from_string.
/*
  Parse timecode representation (hh:mm:ss[:;.]ff).

  @param tc      pointer to an allocated AVTimecode
  @param rate    frame rate in rational form
  @param str     timecode string which will determine the frame start
  @param log_ctx a pointer to an arbitrary struct of which the first field is a
                 pointer to an AVClass struct (used for av_log).
  @return        0 on success, AVERROR otherwise
*/
func AVTimecodeInitFromString(tc *AVTimecode, rate *AVRational, str *CStr, logCtx unsafe.Pointer) (int, error) {
	var tmptc *C.AVTimecode
	if tc != nil {
		tmptc = tc.ptr
	}
	var tmpstr *C.char
	if str != nil {
		tmpstr = str.ptr
	}
	ret := C.av_timecode_init_from_string(tmptc, rate.value, tmpstr, logCtx)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_timecode_check_frame_rate ---

// AVTimecodeCheckFrameRate wraps av_timecode_check_frame_rate.
/*
  Check if the timecode feature is available for the given frame rate

  @return 0 if supported, <0 otherwise
*/
func AVTimecodeCheckFrameRate(rate *AVRational) (int, error) {
	ret := C.av_timecode_check_frame_rate(rate.value)
	return int(ret), WrapErr(int(ret))
}

// --- Function av_ts_make_string ---

// AVTsMakeString wraps av_ts_make_string.
/*
  Fill the provided buffer with a string containing a timestamp
  representation.

  @param buf a buffer with size in bytes of at least AV_TS_MAX_STRING_SIZE
  @param ts the timestamp to represent
  @return the buffer in input
*/
func AVTsMakeString(buf *CStr, ts int64) *CStr {
	var tmpbuf *C.char
	if buf != nil {
		tmpbuf = buf.ptr
	}
	ret := C.av_ts_make_string(tmpbuf, C.int64_t(ts))
	return wrapCStr(ret)
}

// --- Function av_ts_make_time_string2 ---

// AVTsMakeTimeString2 wraps av_ts_make_time_string2.
/*
  Fill the provided buffer with a string containing a timestamp time
  representation.

  @param buf a buffer with size in bytes of at least AV_TS_MAX_STRING_SIZE
  @param ts the timestamp to represent
  @param tb the timebase of the timestamp
  @return the buffer in input
*/
func AVTsMakeTimeString2(buf *CStr, ts int64, tb *AVRational) *CStr {
	var tmpbuf *C.char
	if buf != nil {
		tmpbuf = buf.ptr
	}
	ret := C.av_ts_make_time_string2(tmpbuf, C.int64_t(ts), tb.value)
	return wrapCStr(ret)
}

// --- Function av_ts_make_time_string ---

// av_ts_make_time_string skipped due to tb

// --- Function av_tree_node_alloc ---

// AVTreeNodeAlloc wraps av_tree_node_alloc.
//
//	Allocate an AVTreeNode.
func AVTreeNodeAlloc() *AVTreeNode {
	ret := C.av_tree_node_alloc()
	var retMapped *AVTreeNode
	if ret != nil {
		retMapped = &AVTreeNode{ptr: ret}
	}
	return retMapped
}

// --- Function av_tree_find ---

// av_tree_find skipped due to cmp.

// --- Function av_tree_insert ---

// av_tree_insert skipped due to cmp.

// --- Function av_tree_destroy ---

// AVTreeDestroy wraps av_tree_destroy.
func AVTreeDestroy(t *AVTreeNode) {
	var tmpt *C.struct_AVTreeNode
	if t != nil {
		tmpt = t.ptr
	}
	C.av_tree_destroy(tmpt)
}

// --- Function av_tree_enumerate ---

// av_tree_enumerate skipped due to cmp.

// --- Function av_twofish_alloc ---

// AVTwofishAlloc wraps av_twofish_alloc.
/*
  Allocate an AVTWOFISH context
  To free the struct: av_free(ptr)
*/
func AVTwofishAlloc() *AVTWOFISH {
	ret := C.av_twofish_alloc()
	var retMapped *AVTWOFISH
	if ret != nil {
		retMapped = &AVTWOFISH{ptr: ret}
	}
	return retMapped
}

// --- Function av_twofish_init ---

// AVTwofishInit wraps av_twofish_init.
/*
  Initialize an AVTWOFISH context.

  @param ctx an AVTWOFISH context
  @param key a key of size ranging from 1 to 32 bytes used for encryption/decryption
  @param key_bits number of keybits: 128, 192, 256 If less than the required, padded with zeroes to nearest valid value; return value is 0 if key_bits is 128/192/256, -1 if less than 0, 1 otherwise
*/
func AVTwofishInit(ctx *AVTWOFISH, key unsafe.Pointer, keyBits int) (int, error) {
	var tmpctx *C.struct_AVTWOFISH
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	ret := C.av_twofish_init(tmpctx, (*C.uint8_t)(key), C.int(keyBits))
	return int(ret), WrapErr(int(ret))
}

// --- Function av_twofish_crypt ---

// AVTwofishCrypt wraps av_twofish_crypt.
/*
  Encrypt or decrypt a buffer using a previously initialized context

  @param ctx an AVTWOFISH context
  @param dst destination array, can be equal to src
  @param src source array, can be equal to dst
  @param count number of 16 byte blocks
  @param iv initialization vector for CBC mode, NULL for ECB mode
  @param decrypt 0 for encryption, 1 for decryption
*/
func AVTwofishCrypt(ctx *AVTWOFISH, dst unsafe.Pointer, src unsafe.Pointer, count int, iv unsafe.Pointer, decrypt int) {
	var tmpctx *C.struct_AVTWOFISH
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	C.av_twofish_crypt(tmpctx, (*C.uint8_t)(dst), (*C.uint8_t)(src), C.int(count), (*C.uint8_t)(iv), C.int(decrypt))
}

// --- Function av_tx_init ---

// AVTxInit wraps av_tx_init.
/*
  Initialize a transform context with the given configuration
  (i)MDCTs with an odd length are currently not supported.

  @param ctx the context to allocate, will be NULL on error
  @param tx pointer to the transform function pointer to set
  @param type type the type of transform
  @param inv whether to do an inverse or a forward transform
  @param len the size of the transform in samples
  @param scale pointer to the value to scale the output if supported by type
  @param flags a bitmask of AVTXFlags or 0

  @return 0 on success, negative error code on failure
*/
func AVTxInit(ctx **AVTXContext, tx *AVTxFn, _type AVTXType, inv int, len int, scale unsafe.Pointer, flags uint64) (int, error) {
	var ptrctx **C.AVTXContext
	var tmpctx *C.AVTXContext
	var oldTmpctx *C.AVTXContext
	if ctx != nil {
		innerctx := *ctx
		if innerctx != nil {
			tmpctx = innerctx.ptr
			oldTmpctx = tmpctx
		}
		ptrctx = &tmpctx
	}
	var tmptx *C.av_tx_fn
	if tx != nil {
		tmptx = (*C.av_tx_fn)(unsafe.Pointer(tx))
	}
	ret := C.av_tx_init(ptrctx, tmptx, C.enum_AVTXType(_type), C.int(inv), C.int(len), scale, C.uint64_t(flags))
	if tmpctx != oldTmpctx && ctx != nil {
		if tmpctx != nil {
			*ctx = &AVTXContext{ptr: tmpctx}
		} else {
			*ctx = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function av_tx_uninit ---

// AVTxUninit wraps av_tx_uninit.
//
//	Frees a context and sets *ctx to NULL, does nothing when *ctx == NULL.
func AVTxUninit(ctx **AVTXContext) {
	var ptrctx **C.AVTXContext
	var tmpctx *C.AVTXContext
	var oldTmpctx *C.AVTXContext
	if ctx != nil {
		innerctx := *ctx
		if innerctx != nil {
			tmpctx = innerctx.ptr
			oldTmpctx = tmpctx
		}
		ptrctx = &tmpctx
	}
	C.av_tx_uninit(ptrctx)
	if tmpctx != oldTmpctx && ctx != nil {
		if tmpctx != nil {
			*ctx = &AVTXContext{ptr: tmpctx}
		} else {
			*ctx = nil
		}
	}
}

// --- Function av_uuid_parse ---

// av_uuid_parse skipped due to array typedef (manually wrapped in custom.go)

// --- Function av_uuid_urn_parse ---

// av_uuid_urn_parse skipped due to array typedef (manually wrapped in custom.go)

// --- Function av_uuid_parse_range ---

// av_uuid_parse_range skipped due to array typedef (manually wrapped in custom.go)

// --- Function av_uuid_unparse ---

// av_uuid_unparse skipped due to array typedef (manually wrapped in custom.go)

// --- Function av_uuid_equal ---

// av_uuid_equal skipped due to array typedef (manually wrapped in custom.go)

// --- Function av_uuid_copy ---

// av_uuid_copy skipped due to array typedef (manually wrapped in custom.go)

// --- Function av_uuid_nil ---

// av_uuid_nil skipped due to array typedef (manually wrapped in custom.go)

// --- Function av_video_enc_params_block ---

// AVVideoEncParamsBlock wraps av_video_enc_params_block.
//
//	Get the block at the specified {@code idx}. Must be between 0 and nb_blocks - 1.
func AVVideoEncParamsBlock(par *AVVideoEncParams, idx uint) *AVVideoBlockParams {
	var tmppar *C.AVVideoEncParams
	if par != nil {
		tmppar = par.ptr
	}
	ret := C.av_video_enc_params_block(tmppar, C.uint(idx))
	var retMapped *AVVideoBlockParams
	if ret != nil {
		retMapped = &AVVideoBlockParams{ptr: ret}
	}
	return retMapped
}

// --- Function av_video_enc_params_alloc ---

// AVVideoEncParamsAlloc wraps av_video_enc_params_alloc.
/*
  Allocates memory for AVVideoEncParams of the given type, plus an array of
  {@code nb_blocks} AVVideoBlockParams and initializes the variables. Can be
  freed with a normal av_free() call.

  @param out_size if non-NULL, the size in bytes of the resulting data array is
  written here.
*/
func AVVideoEncParamsAlloc(_type AVVideoEncParamsType, nbBlocks uint, outSize *uint64) *AVVideoEncParams {
	ret := C.av_video_enc_params_alloc(C.enum_AVVideoEncParamsType(_type), C.uint(nbBlocks), (*C.size_t)(unsafe.Pointer(outSize)))
	var retMapped *AVVideoEncParams
	if ret != nil {
		retMapped = &AVVideoEncParams{ptr: ret}
	}
	return retMapped
}

// --- Function av_video_enc_params_create_side_data ---

// AVVideoEncParamsCreateSideData wraps av_video_enc_params_create_side_data.
/*
  Allocates memory for AVEncodeInfoFrame plus an array of
  {@code nb_blocks} AVEncodeInfoBlock in the given AVFrame {@code frame}
  as AVFrameSideData of type AV_FRAME_DATA_VIDEO_ENC_PARAMS
  and initializes the variables.
*/
func AVVideoEncParamsCreateSideData(frame *AVFrame, _type AVVideoEncParamsType, nbBlocks uint) *AVVideoEncParams {
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.av_video_enc_params_create_side_data(tmpframe, C.enum_AVVideoEncParamsType(_type), C.uint(nbBlocks))
	var retMapped *AVVideoEncParams
	if ret != nil {
		retMapped = &AVVideoEncParams{ptr: ret}
	}
	return retMapped
}

// --- Function av_video_hint_rects ---

// AVVideoHintRects wraps av_video_hint_rects.
func AVVideoHintRects(hints *AVVideoHint) *AVVideoRect {
	var tmphints *C.AVVideoHint
	if hints != nil {
		tmphints = hints.ptr
	}
	ret := C.av_video_hint_rects(tmphints)
	var retMapped *AVVideoRect
	if ret != nil {
		retMapped = &AVVideoRect{ptr: ret}
	}
	return retMapped
}

// --- Function av_video_hint_get_rect ---

// AVVideoHintGetRect wraps av_video_hint_get_rect.
func AVVideoHintGetRect(hints *AVVideoHint, idx uint64) *AVVideoRect {
	var tmphints *C.AVVideoHint
	if hints != nil {
		tmphints = hints.ptr
	}
	ret := C.av_video_hint_get_rect(tmphints, C.size_t(idx))
	var retMapped *AVVideoRect
	if ret != nil {
		retMapped = &AVVideoRect{ptr: ret}
	}
	return retMapped
}

// --- Function av_video_hint_alloc ---

// AVVideoHintAlloc wraps av_video_hint_alloc.
/*
  Allocate memory for the AVVideoHint struct along with an nb_rects-sized
  arrays of AVVideoRect.

  The side data contains a list of rectangles for the portions of the frame
  which changed from the last encoded one (and the remainder are assumed to be
  changed), or, alternately (depending on the type parameter) the unchanged
  ones (and the remaining ones are those which changed).
  Macroblocks will thus be hinted either to be P_SKIP-ped or go through the
  regular encoding procedure.

  It's responsibility of the caller to fill the AVRects accordingly, and to set
  the proper AVVideoHintType field.

  @param out_size if non-NULL, the size in bytes of the resulting data array is
                  written here

  @return newly allocated AVVideoHint struct (must be freed by the caller using
          av_free()) on success, NULL on memory allocation failure
*/
func AVVideoHintAlloc(nbRects uint64, outSize *uint64) *AVVideoHint {
	ret := C.av_video_hint_alloc(C.size_t(nbRects), (*C.size_t)(unsafe.Pointer(outSize)))
	var retMapped *AVVideoHint
	if ret != nil {
		retMapped = &AVVideoHint{ptr: ret}
	}
	return retMapped
}

// --- Function av_video_hint_create_side_data ---

// AVVideoHintCreateSideData wraps av_video_hint_create_side_data.
/*
  Same as av_video_hint_alloc(), except newly-allocated AVVideoHint is attached
  as side data of type AV_FRAME_DATA_VIDEO_HINT_INFO to frame.
*/
func AVVideoHintCreateSideData(frame *AVFrame, nbRects uint64) *AVVideoHint {
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.av_video_hint_create_side_data(tmpframe, C.size_t(nbRects))
	var retMapped *AVVideoHint
	if ret != nil {
		retMapped = &AVVideoHint{ptr: ret}
	}
	return retMapped
}

// --- Function av_xtea_alloc ---

// AVXteaAlloc wraps av_xtea_alloc.
//
//	Allocate an AVXTEA context.
func AVXteaAlloc() *AVXTEA {
	ret := C.av_xtea_alloc()
	var retMapped *AVXTEA
	if ret != nil {
		retMapped = &AVXTEA{ptr: ret}
	}
	return retMapped
}

// --- Function av_xtea_init ---

// av_xtea_init skipped due to const array param key

// --- Function av_xtea_le_init ---

// av_xtea_le_init skipped due to const array param key

// --- Function av_xtea_crypt ---

// AVXteaCrypt wraps av_xtea_crypt.
/*
  Encrypt or decrypt a buffer using a previously initialized context,
  in big endian format.

  @param ctx an AVXTEA context
  @param dst destination array, can be equal to src
  @param src source array, can be equal to dst
  @param count number of 8 byte blocks
  @param iv initialization vector for CBC mode, if NULL then ECB will be used
  @param decrypt 0 for encryption, 1 for decryption
*/
func AVXteaCrypt(ctx *AVXTEA, dst unsafe.Pointer, src unsafe.Pointer, count int, iv unsafe.Pointer, decrypt int) {
	var tmpctx *C.AVXTEA
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	C.av_xtea_crypt(tmpctx, (*C.uint8_t)(dst), (*C.uint8_t)(src), C.int(count), (*C.uint8_t)(iv), C.int(decrypt))
}

// --- Function av_xtea_le_crypt ---

// AVXteaLeCrypt wraps av_xtea_le_crypt.
/*
  Encrypt or decrypt a buffer using a previously initialized context,
  in little endian format.

  @param ctx an AVXTEA context
  @param dst destination array, can be equal to src
  @param src source array, can be equal to dst
  @param count number of 8 byte blocks
  @param iv initialization vector for CBC mode, if NULL then ECB will be used
  @param decrypt 0 for encryption, 1 for decryption
*/
func AVXteaLeCrypt(ctx *AVXTEA, dst unsafe.Pointer, src unsafe.Pointer, count int, iv unsafe.Pointer, decrypt int) {
	var tmpctx *C.AVXTEA
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	C.av_xtea_le_crypt(tmpctx, (*C.uint8_t)(dst), (*C.uint8_t)(src), C.int(count), (*C.uint8_t)(iv), C.int(decrypt))
}

// --- Function swr_get_class ---

// SwrGetClass wraps swr_get_class.
/*
  Get the AVClass for SwrContext. It can be used in combination with
  AV_OPT_SEARCH_FAKE_OBJ for examining options.

  @see av_opt_find().
  @return the AVClass of SwrContext
*/
func SwrGetClass() *AVClass {
	ret := C.swr_get_class()
	var retMapped *AVClass
	if ret != nil {
		retMapped = &AVClass{ptr: ret}
	}
	return retMapped
}

// --- Function swr_alloc ---

// SwrAlloc wraps swr_alloc.
/*
  Allocate SwrContext.

  If you use this function you will need to set the parameters (manually or
  with swr_alloc_set_opts2()) before calling swr_init().

  @see swr_alloc_set_opts2(), swr_init(), swr_free()
  @return NULL on error, allocated context otherwise
*/
func SwrAlloc() *SwrContext {
	ret := C.swr_alloc()
	var retMapped *SwrContext
	if ret != nil {
		retMapped = &SwrContext{ptr: ret}
	}
	return retMapped
}

// --- Function swr_init ---

// SwrInit wraps swr_init.
/*
  Initialize context after user parameters have been set.
  @note The context must be configured using the AVOption API.

  @see av_opt_set_int()
  @see av_opt_set_dict()

  @param[in,out]   s Swr context to initialize
  @return AVERROR error code in case of failure.
*/
func SwrInit(s *SwrContext) (int, error) {
	var tmps *C.SwrContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.swr_init(tmps)
	return int(ret), WrapErr(int(ret))
}

// --- Function swr_is_initialized ---

// SwrIsInitialized wraps swr_is_initialized.
/*
  Check whether an swr context has been initialized or not.

  @param[in]       s Swr context to check
  @see swr_init()
  @return positive if it has been initialized, 0 if not initialized
*/
func SwrIsInitialized(s *SwrContext) (int, error) {
	var tmps *C.SwrContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.swr_is_initialized(tmps)
	return int(ret), WrapErr(int(ret))
}

// --- Function swr_alloc_set_opts2 ---

// SwrAllocSetOpts2 wraps swr_alloc_set_opts2.
/*
  Allocate SwrContext if needed and set/reset common parameters.

  This function does not require *ps to be allocated with swr_alloc(). On the
  other hand, swr_alloc() can use swr_alloc_set_opts2() to set the parameters
  on the allocated context.

  @param ps              Pointer to an existing Swr context if available, or to NULL if not.
                         On success, *ps will be set to the allocated context.
  @param out_ch_layout   output channel layout (e.g. AV_CHANNEL_LAYOUT_*)
  @param out_sample_fmt  output sample format (AV_SAMPLE_FMT_*).
  @param out_sample_rate output sample rate (frequency in Hz)
  @param in_ch_layout    input channel layout (e.g. AV_CHANNEL_LAYOUT_*)
  @param in_sample_fmt   input sample format (AV_SAMPLE_FMT_*).
  @param in_sample_rate  input sample rate (frequency in Hz)
  @param log_offset      logging level offset
  @param log_ctx         parent logging context, can be NULL

  @see swr_init(), swr_free()
  @return 0 on success, a negative AVERROR code on error.
          On error, the Swr context is freed and *ps set to NULL.
*/
func SwrAllocSetOpts2(ps **SwrContext, outChLayout *AVChannelLayout, outSampleFmt AVSampleFormat, outSampleRate int, inChLayout *AVChannelLayout, inSampleFmt AVSampleFormat, inSampleRate int, logOffset int, logCtx unsafe.Pointer) (int, error) {
	var ptrps **C.SwrContext
	var tmpps *C.SwrContext
	var oldTmpps *C.SwrContext
	if ps != nil {
		innerps := *ps
		if innerps != nil {
			tmpps = innerps.ptr
			oldTmpps = tmpps
		}
		ptrps = &tmpps
	}
	var tmpoutChLayout *C.AVChannelLayout
	if outChLayout != nil {
		tmpoutChLayout = outChLayout.ptr
	}
	var tmpinChLayout *C.AVChannelLayout
	if inChLayout != nil {
		tmpinChLayout = inChLayout.ptr
	}
	ret := C.swr_alloc_set_opts2(ptrps, tmpoutChLayout, C.enum_AVSampleFormat(outSampleFmt), C.int(outSampleRate), tmpinChLayout, C.enum_AVSampleFormat(inSampleFmt), C.int(inSampleRate), C.int(logOffset), logCtx)
	if tmpps != oldTmpps && ps != nil {
		if tmpps != nil {
			*ps = &SwrContext{ptr: tmpps}
		} else {
			*ps = nil
		}
	}
	return int(ret), WrapErr(int(ret))
}

// --- Function swr_free ---

// SwrFree wraps swr_free.
/*
  Free the given SwrContext and set the pointer to NULL.

  @param[in] s a pointer to a pointer to Swr context
*/
func SwrFree(s **SwrContext) {
	var ptrs **C.SwrContext
	var tmps *C.SwrContext
	var oldTmps *C.SwrContext
	if s != nil {
		inners := *s
		if inners != nil {
			tmps = inners.ptr
			oldTmps = tmps
		}
		ptrs = &tmps
	}
	C.swr_free(ptrs)
	if tmps != oldTmps && s != nil {
		if tmps != nil {
			*s = &SwrContext{ptr: tmps}
		} else {
			*s = nil
		}
	}
}

// --- Function swr_close ---

// SwrClose wraps swr_close.
/*
  Closes the context so that swr_is_initialized() returns 0.

  The context can be brought back to life by running swr_init(),
  swr_init() can also be used without swr_close().
  This function is mainly provided for simplifying the usecase
  where one tries to support libavresample and libswresample.

  @param[in,out] s Swr context to be closed
*/
func SwrClose(s *SwrContext) {
	var tmps *C.SwrContext
	if s != nil {
		tmps = s.ptr
	}
	C.swr_close(tmps)
}

// --- Function swr_convert ---

// swr_convert skipped due to out

// --- Function swr_next_pts ---

// SwrNextPts wraps swr_next_pts.
/*
  Convert the next timestamp from input to output
  timestamps are in 1/(in_sample_rate * out_sample_rate) units.

  @note There are 2 slightly differently behaving modes.
        @li When automatic timestamp compensation is not used, (min_compensation >= FLT_MAX)
               in this case timestamps will be passed through with delays compensated
        @li When automatic timestamp compensation is used, (min_compensation < FLT_MAX)
               in this case the output timestamps will match output sample numbers.
               See ffmpeg-resampler(1) for the two modes of compensation.

  @param[in] s     initialized Swr context
  @param[in] pts   timestamp for the next input sample, INT64_MIN if unknown
  @see swr_set_compensation(), swr_drop_output(), and swr_inject_silence() are
       function used internally for timestamp compensation.
  @return the output timestamp for the next output sample
*/
func SwrNextPts(s *SwrContext, pts int64) int64 {
	var tmps *C.SwrContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.swr_next_pts(tmps, C.int64_t(pts))
	return int64(ret)
}

// --- Function swr_set_compensation ---

// SwrSetCompensation wraps swr_set_compensation.
/*
  Activate resampling compensation ("soft" compensation). This function is
  internally called when needed in swr_next_pts().

  @param[in,out] s             allocated Swr context. If it is not initialized,
                               or SWR_FLAG_RESAMPLE is not set, swr_init() is
                               called with the flag set.
  @param[in]     sample_delta  delta in PTS per sample
  @param[in]     compensation_distance number of samples to compensate for
  @return    >= 0 on success, AVERROR error codes if:
             @li @c s is NULL,
             @li @c compensation_distance is less than 0,
             @li @c compensation_distance is 0 but sample_delta is not,
             @li compensation unsupported by resampler, or
             @li swr_init() fails when called.
*/
func SwrSetCompensation(s *SwrContext, sampleDelta int, compensationDistance int) (int, error) {
	var tmps *C.SwrContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.swr_set_compensation(tmps, C.int(sampleDelta), C.int(compensationDistance))
	return int(ret), WrapErr(int(ret))
}

// --- Function swr_set_channel_mapping ---

// SwrSetChannelMapping wraps swr_set_channel_mapping.
/*
  Set a customized input channel mapping.

  @param[in,out] s           allocated Swr context, not yet initialized
  @param[in]     channel_map customized input channel mapping (array of channel
                             indexes, -1 for a muted channel)
  @return >= 0 on success, or AVERROR error code in case of failure.
*/
func SwrSetChannelMapping(s *SwrContext, channelMap *int) (int, error) {
	var tmps *C.SwrContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.swr_set_channel_mapping(tmps, (*C.int)(unsafe.Pointer(channelMap)))
	return int(ret), WrapErr(int(ret))
}

// --- Function swr_build_matrix2 ---

// swr_build_matrix2 skipped due to matrix (non-output primitive pointer)

// --- Function swr_set_matrix ---

// swr_set_matrix skipped due to matrix (non-output primitive pointer)

// --- Function swr_drop_output ---

// SwrDropOutput wraps swr_drop_output.
/*
  Drops the specified number of output samples.

  This function, along with swr_inject_silence(), is called by swr_next_pts()
  if needed for "hard" compensation.

  @param s     allocated Swr context
  @param count number of samples to be dropped

  @return >= 0 on success, or a negative AVERROR code on failure
*/
func SwrDropOutput(s *SwrContext, count int) (int, error) {
	var tmps *C.SwrContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.swr_drop_output(tmps, C.int(count))
	return int(ret), WrapErr(int(ret))
}

// --- Function swr_inject_silence ---

// SwrInjectSilence wraps swr_inject_silence.
/*
  Injects the specified number of silence samples.

  This function, along with swr_drop_output(), is called by swr_next_pts()
  if needed for "hard" compensation.

  @param s     allocated Swr context
  @param count number of samples to be dropped

  @return >= 0 on success, or a negative AVERROR code on failure
*/
func SwrInjectSilence(s *SwrContext, count int) (int, error) {
	var tmps *C.SwrContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.swr_inject_silence(tmps, C.int(count))
	return int(ret), WrapErr(int(ret))
}

// --- Function swr_get_delay ---

// SwrGetDelay wraps swr_get_delay.
/*
  Gets the delay the next input sample will experience relative to the next output sample.

  Swresample can buffer data if more input has been provided than available
  output space, also converting between sample rates needs a delay.
  This function returns the sum of all such delays.
  The exact delay is not necessarily an integer value in either input or
  output sample rate. Especially when downsampling by a large value, the
  output sample rate may be a poor choice to represent the delay, similarly
  for upsampling and the input sample rate.

  @param s     swr context
  @param base  timebase in which the returned delay will be:
               @li if it's set to 1 the returned delay is in seconds
               @li if it's set to 1000 the returned delay is in milliseconds
               @li if it's set to the input sample rate then the returned
                   delay is in input samples
               @li if it's set to the output sample rate then the returned
                   delay is in output samples
               @li if it's the least common multiple of in_sample_rate and
                   out_sample_rate then an exact rounding-free delay will be
                   returned
  @returns     the delay in 1 / @c base units.
*/
func SwrGetDelay(s *SwrContext, base int64) int64 {
	var tmps *C.SwrContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.swr_get_delay(tmps, C.int64_t(base))
	return int64(ret)
}

// --- Function swr_get_out_samples ---

// SwrGetOutSamples wraps swr_get_out_samples.
/*
  Find an upper bound on the number of samples that the next swr_convert
  call will output, if called with in_samples of input samples. This
  depends on the internal state, and anything changing the internal state
  (like further swr_convert() calls) will may change the number of samples
  swr_get_out_samples() returns for the same number of input samples.

  @param in_samples    number of input samples.
  @note any call to swr_inject_silence(), swr_convert(), swr_next_pts()
        or swr_set_compensation() invalidates this limit
  @note it is recommended to pass the correct available buffer size
        to all functions like swr_convert() even if swr_get_out_samples()
        indicates that less would be used.
  @returns an upper bound on the number of samples that the next swr_convert
           will output or a negative value to indicate an error
*/
func SwrGetOutSamples(s *SwrContext, inSamples int) (int, error) {
	var tmps *C.SwrContext
	if s != nil {
		tmps = s.ptr
	}
	ret := C.swr_get_out_samples(tmps, C.int(inSamples))
	return int(ret), WrapErr(int(ret))
}

// --- Function swresample_version ---

// SwresampleVersion wraps swresample_version.
/*
  Return the @ref LIBSWRESAMPLE_VERSION_INT constant.

  This is useful to check if the build-time libswresample has the same version
  as the run-time one.

  @returns     the unsigned int-typed version
*/
func SwresampleVersion() uint {
	ret := C.swresample_version()
	return uint(ret)
}

// --- Function swresample_configuration ---

// SwresampleConfiguration wraps swresample_configuration.
/*
  Return the swr build-time configuration.

  @returns     the build-time @c ./configure flags
*/
func SwresampleConfiguration() *CStr {
	ret := C.swresample_configuration()
	return wrapCStr(ret)
}

// --- Function swresample_license ---

// SwresampleLicense wraps swresample_license.
/*
  Return the swr license.

  @returns     the license of libswresample, determined at build-time
*/
func SwresampleLicense() *CStr {
	ret := C.swresample_license()
	return wrapCStr(ret)
}

// --- Function swr_convert_frame ---

// SwrConvertFrame wraps swr_convert_frame.
/*
  Convert the samples in the input AVFrame and write them to the output AVFrame.

  Input and output AVFrames must have channel_layout, sample_rate and format set.

  If the output AVFrame does not have the data pointers allocated the nb_samples
  field will be set using av_frame_get_buffer()
  is called to allocate the frame.

  The output AVFrame can be NULL or have fewer allocated samples than required.
  In this case, any remaining samples not written to the output will be added
  to an internal FIFO buffer, to be returned at the next call to this function
  or to swr_convert().

  If converting sample rate, there may be data remaining in the internal
  resampling delay buffer. swr_get_delay() tells the number of
  remaining samples. To get this data as output, call this function or
  swr_convert() with NULL input.

  If the SwrContext configuration does not match the output and
  input AVFrame settings the conversion does not take place and depending on
  which AVFrame is not matching AVERROR_OUTPUT_CHANGED, AVERROR_INPUT_CHANGED
  or the result of a bitwise-OR of them is returned.

  @see swr_delay()
  @see swr_convert()
  @see swr_get_delay()

  @param swr             audio resample context
  @param output          output AVFrame
  @param input           input AVFrame
  @return                0 on success, AVERROR on failure or nonmatching
                         configuration.
*/
func SwrConvertFrame(swr *SwrContext, output *AVFrame, input *AVFrame) (int, error) {
	var tmpswr *C.SwrContext
	if swr != nil {
		tmpswr = swr.ptr
	}
	var tmpoutput *C.AVFrame
	if output != nil {
		tmpoutput = output.ptr
	}
	var tmpinput *C.AVFrame
	if input != nil {
		tmpinput = input.ptr
	}
	ret := C.swr_convert_frame(tmpswr, tmpoutput, tmpinput)
	return int(ret), WrapErr(int(ret))
}

// --- Function swr_config_frame ---

// SwrConfigFrame wraps swr_config_frame.
/*
  Configure or reconfigure the SwrContext using the information
  provided by the AVFrames.

  The original resampling context is reset even on failure.
  The function calls swr_close() internally if the context is open.

  @see swr_close();

  @param swr             audio resample context
  @param out             output AVFrame
  @param in              input AVFrame
  @return                0 on success, AVERROR on failure.
*/
func SwrConfigFrame(swr *SwrContext, out *AVFrame, in *AVFrame) (int, error) {
	var tmpswr *C.SwrContext
	if swr != nil {
		tmpswr = swr.ptr
	}
	var tmpout *C.AVFrame
	if out != nil {
		tmpout = out.ptr
	}
	var tmpin *C.AVFrame
	if in != nil {
		tmpin = in.ptr
	}
	ret := C.swr_config_frame(tmpswr, tmpout, tmpin)
	return int(ret), WrapErr(int(ret))
}

// --- Function swscale_version ---

// SwscaleVersion wraps swscale_version.
/*

  Return the LIBSWSCALE_VERSION_INT constant.
*/
func SwscaleVersion() uint {
	ret := C.swscale_version()
	return uint(ret)
}

// --- Function swscale_configuration ---

// SwscaleConfiguration wraps swscale_configuration.
//
//	Return the libswscale build-time configuration.
func SwscaleConfiguration() *CStr {
	ret := C.swscale_configuration()
	return wrapCStr(ret)
}

// --- Function swscale_license ---

// SwscaleLicense wraps swscale_license.
//
//	Return the libswscale license.
func SwscaleLicense() *CStr {
	ret := C.swscale_license()
	return wrapCStr(ret)
}

// --- Function sws_get_class ---

// SwsGetClass wraps sws_get_class.
/*
  Get the AVClass for SwsContext. It can be used in combination with
  AV_OPT_SEARCH_FAKE_OBJ for examining options.

  @see av_opt_find().
*/
func SwsGetClass() *AVClass {
	ret := C.sws_get_class()
	var retMapped *AVClass
	if ret != nil {
		retMapped = &AVClass{ptr: ret}
	}
	return retMapped
}

// --- Function sws_alloc_context ---

// SwsAllocContext wraps sws_alloc_context.
//
//	Allocate an empty SwsContext and set its fields to default values.
func SwsAllocContext() *SwsContext {
	ret := C.sws_alloc_context()
	var retMapped *SwsContext
	if ret != nil {
		retMapped = &SwsContext{ptr: ret}
	}
	return retMapped
}

// --- Function sws_free_context ---

// SwsFreeContext wraps sws_free_context.
/*
  Free the context and everything associated with it, and write NULL
  to the provided pointer.
*/
func SwsFreeContext(ctx **SwsContext) {
	var ptrctx **C.SwsContext
	var tmpctx *C.SwsContext
	var oldTmpctx *C.SwsContext
	if ctx != nil {
		innerctx := *ctx
		if innerctx != nil {
			tmpctx = innerctx.ptr
			oldTmpctx = tmpctx
		}
		ptrctx = &tmpctx
	}
	C.sws_free_context(ptrctx)
	if tmpctx != oldTmpctx && ctx != nil {
		if tmpctx != nil {
			*ctx = &SwsContext{ptr: tmpctx}
		} else {
			*ctx = nil
		}
	}
}

// --- Function sws_test_format ---

// SwsTestFormat wraps sws_test_format.
/*
  Test if a given pixel format is supported.

  @param output  If 0, test if compatible with the source/input frame;
                 otherwise, with the destination/output frame.
  @param format  The format to check.

  @return A positive integer if supported, 0 otherwise.
*/
func SwsTestFormat(format AVPixelFormat, output int) (int, error) {
	ret := C.sws_test_format(C.enum_AVPixelFormat(format), C.int(output))
	return int(ret), WrapErr(int(ret))
}

// --- Function sws_test_colorspace ---

// SwsTestColorspace wraps sws_test_colorspace.
/*
  Test if a given color space is supported.

  @param output  If 0, test if compatible with the source/input frame;
                 otherwise, with the destination/output frame.
  @param colorspace The colorspace to check.

  @return A positive integer if supported, 0 otherwise.
*/
func SwsTestColorspace(colorspace AVColorSpace, output int) (int, error) {
	ret := C.sws_test_colorspace(C.enum_AVColorSpace(colorspace), C.int(output))
	return int(ret), WrapErr(int(ret))
}

// --- Function sws_test_primaries ---

// SwsTestPrimaries wraps sws_test_primaries.
/*
  Test if a given set of color primaries is supported.

  @param output  If 0, test if compatible with the source/input frame;
                 otherwise, with the destination/output frame.
  @param primaries The color primaries to check.

  @return A positive integer if supported, 0 otherwise.
*/
func SwsTestPrimaries(primaries AVColorPrimaries, output int) (int, error) {
	ret := C.sws_test_primaries(C.enum_AVColorPrimaries(primaries), C.int(output))
	return int(ret), WrapErr(int(ret))
}

// --- Function sws_test_transfer ---

// SwsTestTransfer wraps sws_test_transfer.
/*
  Test if a given color transfer function is supported.

  @param output  If 0, test if compatible with the source/input frame;
                 otherwise, with the destination/output frame.
  @param trc     The color transfer function to check.

  @return A positive integer if supported, 0 otherwise.
*/
func SwsTestTransfer(trc AVColorTransferCharacteristic, output int) (int, error) {
	ret := C.sws_test_transfer(C.enum_AVColorTransferCharacteristic(trc), C.int(output))
	return int(ret), WrapErr(int(ret))
}

// --- Function sws_test_frame ---

// SwsTestFrame wraps sws_test_frame.
/*
  Helper function to run all sws_test_* against a frame, as well as testing
  the basic frame properties for sanity. Ignores irrelevant properties - for
  example, AVColorSpace is not checked for RGB frames.
*/
func SwsTestFrame(frame *AVFrame, output int) (int, error) {
	var tmpframe *C.AVFrame
	if frame != nil {
		tmpframe = frame.ptr
	}
	ret := C.sws_test_frame(tmpframe, C.int(output))
	return int(ret), WrapErr(int(ret))
}

// --- Function sws_frame_setup ---

// SwsFrameSetup wraps sws_frame_setup.
/*
  Like `sws_scale_frame`, but without actually scaling. It will instead
  merely initialize internal state that *would* be required to perform the
  operation, as well as returning the correct error code for unsupported
  frame combinations.

  @param ctx   The scaling context.
  @param dst   The destination frame to consider.
  @param src   The source frame to consider.
  @return 0 on success, a negative AVERROR code on failure.
*/
func SwsFrameSetup(ctx *SwsContext, dst *AVFrame, src *AVFrame) (int, error) {
	var tmpctx *C.SwsContext
	if ctx != nil {
		tmpctx = ctx.ptr
	}
	var tmpdst *C.AVFrame
	if dst != nil {
		tmpdst = dst.ptr
	}
	var tmpsrc *C.AVFrame
	if src != nil {
		tmpsrc = src.ptr
	}
	ret := C.sws_frame_setup(tmpctx, tmpdst, tmpsrc)
	return int(ret), WrapErr(int(ret))
}

// --- Function sws_is_noop ---

// SwsIsNoop wraps sws_is_noop.
/*
  Check if a given conversion is a noop. Returns a positive integer if
  no operation needs to be performed, 0 otherwise.
*/
func SwsIsNoop(dst *AVFrame, src *AVFrame) (int, error) {
	var tmpdst *C.AVFrame
	if dst != nil {
		tmpdst = dst.ptr
	}
	var tmpsrc *C.AVFrame
	if src != nil {
		tmpsrc = src.ptr
	}
	ret := C.sws_is_noop(tmpdst, tmpsrc)
	return int(ret), WrapErr(int(ret))
}

// --- Function sws_scale_frame ---

// SwsScaleFrame wraps sws_scale_frame.
/*
  Scale source data from `src` and write the output to `dst`.

  This function can be used directly on an allocated context, without setting
  up any frame properties or calling `sws_init_context()`. Such usage is fully
  dynamic and does not require reallocation if the frame properties change.

  Alternatively, this function can be called on a context that has been
  explicitly initialized. However, this is provided only for backwards
  compatibility. In this usage mode, all frame properties must be correctly
  set at init time, and may no longer change after initialization.

  @param ctx   The scaling context.
  @param dst   The destination frame. The data buffers may either be already
               allocated by the caller or left clear, in which case they will
               be allocated by the scaler. The latter may have performance
               advantages - e.g. in certain cases some (or all) output planes
               may be references to input planes, rather than copies.
  @param src   The source frame. If the data buffers are set to NULL, then
               this function behaves identically to `sws_frame_setup`.
  @return >= 0 on success, a negative AVERROR code on failure.
*/
func SwsScaleFrame(c *SwsContext, dst *AVFrame, src *AVFrame) (int, error) {
	var tmpc *C.SwsContext
	if c != nil {
		tmpc = c.ptr
	}
	var tmpdst *C.AVFrame
	if dst != nil {
		tmpdst = dst.ptr
	}
	var tmpsrc *C.AVFrame
	if src != nil {
		tmpsrc = src.ptr
	}
	ret := C.sws_scale_frame(tmpc, tmpdst, tmpsrc)
	return int(ret), WrapErr(int(ret))
}

// --- Function sws_getCoefficients ---

// SwsGetcoefficients wraps sws_getCoefficients.
/*
  Return a pointer to yuv<->rgb coefficients for the given colorspace
  suitable for sws_setColorspaceDetails().

  @param colorspace One of the SWS_CS_* macros. If invalid,
  SWS_CS_DEFAULT is used.
*/
func SwsGetcoefficients(colorspace int) *int {
	ret := C.sws_getCoefficients(C.int(colorspace))
	return (*int)(unsafe.Pointer(ret))
}

// --- Function sws_isSupportedInput ---

// SwsIssupportedinput wraps sws_isSupportedInput.
/*
  Return a positive value if pix_fmt is a supported input format, 0
  otherwise.
*/
func SwsIssupportedinput(pixFmt AVPixelFormat) (int, error) {
	ret := C.sws_isSupportedInput(C.enum_AVPixelFormat(pixFmt))
	return int(ret), WrapErr(int(ret))
}

// --- Function sws_isSupportedOutput ---

// SwsIssupportedoutput wraps sws_isSupportedOutput.
/*
  Return a positive value if pix_fmt is a supported output format, 0
  otherwise.
*/
func SwsIssupportedoutput(pixFmt AVPixelFormat) (int, error) {
	ret := C.sws_isSupportedOutput(C.enum_AVPixelFormat(pixFmt))
	return int(ret), WrapErr(int(ret))
}

// --- Function sws_isSupportedEndiannessConversion ---

// SwsIssupportedendiannessconversion wraps sws_isSupportedEndiannessConversion.
/*
  @param[in]  pix_fmt the pixel format
  @return a positive value if an endianness conversion for pix_fmt is
  supported, 0 otherwise.
*/
func SwsIssupportedendiannessconversion(pixFmt AVPixelFormat) (int, error) {
	ret := C.sws_isSupportedEndiannessConversion(C.enum_AVPixelFormat(pixFmt))
	return int(ret), WrapErr(int(ret))
}

// --- Function sws_init_context ---

// SwsInitContext wraps sws_init_context.
/*
  Initialize the swscaler context sws_context.

  This function is considered deprecated, and provided only for backwards
  compatibility with sws_scale() and sws_start_frame(). The preferred way to
  use libswscale is to set all frame properties correctly and call
  sws_scale_frame() directly, without explicitly initializing the context.

  @return zero or positive value on success, a negative value on
  error
*/
func SwsInitContext(swsContext *SwsContext, srcFilter *SwsFilter, dstFilter *SwsFilter) (int, error) {
	var tmpswsContext *C.SwsContext
	if swsContext != nil {
		tmpswsContext = swsContext.ptr
	}
	var tmpsrcFilter *C.SwsFilter
	if srcFilter != nil {
		tmpsrcFilter = srcFilter.ptr
	}
	var tmpdstFilter *C.SwsFilter
	if dstFilter != nil {
		tmpdstFilter = dstFilter.ptr
	}
	ret := C.sws_init_context(tmpswsContext, tmpsrcFilter, tmpdstFilter)
	return int(ret), WrapErr(int(ret))
}

// --- Function sws_freeContext ---

// SwsFreecontext wraps sws_freeContext.
/*
  Free the swscaler context swsContext.
  If swsContext is NULL, then does nothing.
*/
func SwsFreecontext(swsContext *SwsContext) {
	var tmpswsContext *C.SwsContext
	if swsContext != nil {
		tmpswsContext = swsContext.ptr
	}
	C.sws_freeContext(tmpswsContext)
}

// --- Function sws_getContext ---

// sws_getContext skipped due to param (non-output primitive pointer)

// --- Function sws_scale ---

// sws_scale skipped due to srcSlice

// --- Function sws_frame_start ---

// SwsFrameStart wraps sws_frame_start.
/*
  Initialize the scaling process for a given pair of source/destination frames.
  Must be called before any calls to sws_send_slice() and sws_receive_slice().
  Requires a context that has been previously been initialized with
  sws_init_context().

  This function will retain references to src and dst, so they must both use
  refcounted buffers (if allocated by the caller, in case of dst).

  @param c   The scaling context
  @param dst The destination frame.

             The data buffers may either be already allocated by the caller or
             left clear, in which case they will be allocated by the scaler.
             The latter may have performance advantages - e.g. in certain cases
             some output planes may be references to input planes, rather than
             copies.

             Output data will be written into this frame in successful
             sws_receive_slice() calls.
  @param src The source frame. The data buffers must be allocated, but the
             frame data does not have to be ready at this point. Data
             availability is then signalled by sws_send_slice().
  @return 0 on success, a negative AVERROR code on failure

  @see sws_frame_end()
*/
func SwsFrameStart(c *SwsContext, dst *AVFrame, src *AVFrame) (int, error) {
	var tmpc *C.SwsContext
	if c != nil {
		tmpc = c.ptr
	}
	var tmpdst *C.AVFrame
	if dst != nil {
		tmpdst = dst.ptr
	}
	var tmpsrc *C.AVFrame
	if src != nil {
		tmpsrc = src.ptr
	}
	ret := C.sws_frame_start(tmpc, tmpdst, tmpsrc)
	return int(ret), WrapErr(int(ret))
}

// --- Function sws_frame_end ---

// SwsFrameEnd wraps sws_frame_end.
/*
  Finish the scaling process for a pair of source/destination frames previously
  submitted with sws_frame_start(). Must be called after all sws_send_slice()
  and sws_receive_slice() calls are done, before any new sws_frame_start()
  calls.

  @param c   The scaling context
*/
func SwsFrameEnd(c *SwsContext) {
	var tmpc *C.SwsContext
	if c != nil {
		tmpc = c.ptr
	}
	C.sws_frame_end(tmpc)
}

// --- Function sws_send_slice ---

// SwsSendSlice wraps sws_send_slice.
/*
  Indicate that a horizontal slice of input data is available in the source
  frame previously provided to sws_frame_start(). The slices may be provided in
  any order, but may not overlap. For vertically subsampled pixel formats, the
  slices must be aligned according to subsampling.

  @param c   The scaling context
  @param slice_start first row of the slice
  @param slice_height number of rows in the slice

  @return a non-negative number on success, a negative AVERROR code on failure.
*/
func SwsSendSlice(c *SwsContext, sliceStart uint, sliceHeight uint) (int, error) {
	var tmpc *C.SwsContext
	if c != nil {
		tmpc = c.ptr
	}
	ret := C.sws_send_slice(tmpc, C.uint(sliceStart), C.uint(sliceHeight))
	return int(ret), WrapErr(int(ret))
}

// --- Function sws_receive_slice ---

// SwsReceiveSlice wraps sws_receive_slice.
/*
  Request a horizontal slice of the output data to be written into the frame
  previously provided to sws_frame_start().

  @param c   The scaling context
  @param slice_start first row of the slice; must be a multiple of
                     sws_receive_slice_alignment()
  @param slice_height number of rows in the slice; must be a multiple of
                      sws_receive_slice_alignment(), except for the last slice
                      (i.e. when slice_start+slice_height is equal to output
                      frame height)

  @return a non-negative number if the data was successfully written into the output
          AVERROR(EAGAIN) if more input data needs to be provided before the
                          output can be produced
          another negative AVERROR code on other kinds of scaling failure
*/
func SwsReceiveSlice(c *SwsContext, sliceStart uint, sliceHeight uint) (int, error) {
	var tmpc *C.SwsContext
	if c != nil {
		tmpc = c.ptr
	}
	ret := C.sws_receive_slice(tmpc, C.uint(sliceStart), C.uint(sliceHeight))
	return int(ret), WrapErr(int(ret))
}

// --- Function sws_receive_slice_alignment ---

// SwsReceiveSliceAlignment wraps sws_receive_slice_alignment.
/*
  Get the alignment required for slices. Requires a context that has been
  previously been initialized with sws_init_context().

  @param c   The scaling context
  @return alignment required for output slices requested with sws_receive_slice().
          Slice offsets and sizes passed to sws_receive_slice() must be
          multiples of the value returned from this function.
*/
func SwsReceiveSliceAlignment(c *SwsContext) uint {
	var tmpc *C.SwsContext
	if c != nil {
		tmpc = c.ptr
	}
	ret := C.sws_receive_slice_alignment(tmpc)
	return uint(ret)
}

// --- Function sws_setColorspaceDetails ---

// sws_setColorspaceDetails skipped due to const array param invTable

// --- Function sws_getColorspaceDetails ---

// sws_getColorspaceDetails skipped due to invTable

// --- Function sws_allocVec ---

// SwsAllocvec wraps sws_allocVec.
//
//	Allocate and return an uninitialized vector with length coefficients.
func SwsAllocvec(length int) *SwsVector {
	ret := C.sws_allocVec(C.int(length))
	var retMapped *SwsVector
	if ret != nil {
		retMapped = &SwsVector{ptr: ret}
	}
	return retMapped
}

// --- Function sws_getGaussianVec ---

// SwsGetgaussianvec wraps sws_getGaussianVec.
/*
  Return a normalized Gaussian curve used to filter stuff
  quality = 3 is high quality, lower is lower quality.
*/
func SwsGetgaussianvec(variance float64, quality float64) *SwsVector {
	ret := C.sws_getGaussianVec(C.double(variance), C.double(quality))
	var retMapped *SwsVector
	if ret != nil {
		retMapped = &SwsVector{ptr: ret}
	}
	return retMapped
}

// --- Function sws_scaleVec ---

// SwsScalevec wraps sws_scaleVec.
//
//	Scale all the coefficients of a by the scalar value.
func SwsScalevec(a *SwsVector, scalar float64) {
	var tmpa *C.SwsVector
	if a != nil {
		tmpa = a.ptr
	}
	C.sws_scaleVec(tmpa, C.double(scalar))
}

// --- Function sws_normalizeVec ---

// SwsNormalizevec wraps sws_normalizeVec.
//
//	Scale all the coefficients of a so that their sum equals height.
func SwsNormalizevec(a *SwsVector, height float64) {
	var tmpa *C.SwsVector
	if a != nil {
		tmpa = a.ptr
	}
	C.sws_normalizeVec(tmpa, C.double(height))
}

// --- Function sws_freeVec ---

// SwsFreevec wraps sws_freeVec.
func SwsFreevec(a *SwsVector) {
	var tmpa *C.SwsVector
	if a != nil {
		tmpa = a.ptr
	}
	C.sws_freeVec(tmpa)
}

// --- Function sws_getDefaultFilter ---

// SwsGetdefaultfilter wraps sws_getDefaultFilter.
func SwsGetdefaultfilter(lumaGblur float32, chromaGblur float32, lumaSharpen float32, chromaSharpen float32, chromaHshift float32, chromaVshift float32, verbose int) *SwsFilter {
	ret := C.sws_getDefaultFilter(C.float(lumaGblur), C.float(chromaGblur), C.float(lumaSharpen), C.float(chromaSharpen), C.float(chromaHshift), C.float(chromaVshift), C.int(verbose))
	var retMapped *SwsFilter
	if ret != nil {
		retMapped = &SwsFilter{ptr: ret}
	}
	return retMapped
}

// --- Function sws_freeFilter ---

// SwsFreefilter wraps sws_freeFilter.
func SwsFreefilter(filter *SwsFilter) {
	var tmpfilter *C.SwsFilter
	if filter != nil {
		tmpfilter = filter.ptr
	}
	C.sws_freeFilter(tmpfilter)
}

// --- Function sws_getCachedContext ---

// sws_getCachedContext skipped due to param (non-output primitive pointer)

// --- Function sws_convertPalette8ToPacked32 ---

// SwsConvertpalette8Topacked32 wraps sws_convertPalette8ToPacked32.
/*
  Convert an 8-bit paletted frame into a frame with a color depth of 32 bits.

  The output frame will have the same packed format as the palette.

  @param src        source frame buffer
  @param dst        destination frame buffer
  @param num_pixels number of pixels to convert
  @param palette    array with [256] entries, which must match color arrangement (RGB or BGR) of src
*/
func SwsConvertpalette8Topacked32(src unsafe.Pointer, dst unsafe.Pointer, numPixels int, palette unsafe.Pointer) {
	C.sws_convertPalette8ToPacked32((*C.uint8_t)(src), (*C.uint8_t)(dst), C.int(numPixels), (*C.uint8_t)(palette))
}

// --- Function sws_convertPalette8ToPacked24 ---

// SwsConvertpalette8Topacked24 wraps sws_convertPalette8ToPacked24.
/*
  Convert an 8-bit paletted frame into a frame with a color depth of 24 bits.

  With the palette format "ABCD", the destination frame ends up with the format "ABC".

  @param src        source frame buffer
  @param dst        destination frame buffer
  @param num_pixels number of pixels to convert
  @param palette    array with [256] entries, which must match color arrangement (RGB or BGR) of src
*/
func SwsConvertpalette8Topacked24(src unsafe.Pointer, dst unsafe.Pointer, numPixels int, palette unsafe.Pointer) {
	C.sws_convertPalette8ToPacked24((*C.uint8_t)(src), (*C.uint8_t)(dst), C.int(numPixels), (*C.uint8_t)(palette))
}
