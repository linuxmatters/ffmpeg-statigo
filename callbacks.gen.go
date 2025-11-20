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

// AVFilterActionFunc is a function pointer typedef for avfilter_action_func.
type AVFilterActionFunc unsafe.Pointer

// AVFilterExecuteFunc is a function pointer typedef for avfilter_execute_func.
type AVFilterExecuteFunc unsafe.Pointer

// AVFormatControlMessage is a function pointer typedef for av_format_control_message.
type AVFormatControlMessage unsafe.Pointer

// AVOpencallback is a function pointer typedef for AVOpenCallback.
type AVOpencallback unsafe.Pointer

// AVCspTrcFunction is a function pointer typedef for av_csp_trc_function.
type AVCspTrcFunction unsafe.Pointer

// AVCspEotfFunction is a function pointer typedef for av_csp_eotf_function.
type AVCspEotfFunction unsafe.Pointer

// AVFifocb is a function pointer typedef for AVFifoCB.
type AVFifocb unsafe.Pointer

// AVPixelutilsSadFn is a function pointer typedef for av_pixelutils_sad_fn.
type AVPixelutilsSadFn unsafe.Pointer

// AVTxFn is a function pointer typedef for av_tx_fn.
type AVTxFn unsafe.Pointer
