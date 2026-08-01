package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

var failLog = log.New(os.Stderr, "", 0)

var AVLibPath, _ = filepath.Abs("include")

var files = []string{
	"libavcodec/ac3_parser.h",
	"libavcodec/adts_parser.h",
	"libavcodec/avcodec.h",
	"libavcodec/avdct.h",
	"libavcodec/bsf.h",
	"libavcodec/codec.h",
	"libavcodec/codec_desc.h",
	"libavcodec/codec_id.h",
	"libavcodec/codec_par.h",
	////"libavcodec/d3d11va.h", // needs <d3d11.h> (Windows SDK)
	"libavcodec/defs.h",
	"libavcodec/dirac.h",
	"libavcodec/dv_profile.h",
	////"libavcodec/dxva2.h", // needs <d3d9.h>/<dxva2api.h> (Windows SDK)
	"libavcodec/exif.h",
	////"libavcodec/jni.h", // Android-only symbols, link risk on Linux/macOS static lib
	////"libavcodec/mediacodec.h", // Android-only symbols, link risk
	"libavcodec/packet.h",
	////"libavcodec/qsv.h", // needs <mfxvideo.h> (Intel Media SDK)
	////"libavcodec/smpte_436m.h", // header vendored but symbols absent from FFmpeg 8.1.2 static lib (link error); needs a build that ships the SMPTE 436M/291M code
	////"libavcodec/vdpau.h", // needs <vdpau/vdpau.h>
	"libavcodec/version.h",
	"libavcodec/version_major.h",
	////"libavcodec/videotoolbox.h", // needs <VideoToolbox/VideoToolbox.h> (Apple)
	"libavcodec/vorbis_parser.h",
	"libavdevice/avdevice.h",
	"libavdevice/version.h",
	"libavdevice/version_major.h",
	"libavfilter/avfilter.h",
	"libavfilter/buffersink.h",
	"libavfilter/buffersrc.h",
	"libavfilter/version.h",
	"libavfilter/version_major.h",
	"libavformat/avformat.h",
	"libavformat/avio.h",
	"libavformat/version.h",
	"libavformat/version_major.h",
	"libavutil/adler32.h",
	"libavutil/aes.h",
	"libavutil/aes_ctr.h",
	"libavutil/ambient_viewing_environment.h",
	////"libavutil/attributes.h", // a compiler attribute macro collection, not an API header
	"libavutil/audio_fifo.h",
	"libavutil/avassert.h",
	"libavutil/avconfig.h",
	"libavutil/avstring.h",
	"libavutil/avutil.h",
	"libavutil/base64.h",
	"libavutil/blowfish.h",
	"libavutil/bprint.h",
	"libavutil/bswap.h",
	"libavutil/buffer.h",
	"libavutil/camellia.h",
	"libavutil/cast5.h",
	"libavutil/channel_layout.h",
	"libavutil/common.h",
	"libavutil/container_fifo.h",
	"libavutil/cpu.h",
	"libavutil/crc.h",
	"libavutil/csp.h",
	"libavutil/des.h",
	"libavutil/detection_bbox.h",
	"libavutil/dict.h",
	"libavutil/display.h",
	"libavutil/dovi_meta.h",
	"libavutil/downmix_info.h",
	"libavutil/encryption_info.h",
	"libavutil/error.h",
	"libavutil/eval.h",
	"libavutil/executor.h",
	"libavutil/ffversion.h",
	"libavutil/fifo.h",
	"libavutil/file.h",
	"libavutil/film_grain_params.h",
	"libavutil/frame.h",
	"libavutil/hash.h",
	"libavutil/hdr_dynamic_metadata.h",
	"libavutil/hdr_dynamic_vivid_metadata.h",
	"libavutil/hmac.h",
	"libavutil/hwcontext.h",
	////"libavutil/hwcontext_amf.h", // needs <AMF/...> SDK
	////"libavutil/hwcontext_cuda.h",
	////"libavutil/hwcontext_d3d11va.h", // needs <d3d11.h> (Windows SDK)
	////"libavutil/hwcontext_d3d12va.h", // needs <d3d12.h> (Windows SDK)
	"libavutil/hwcontext_drm.h",
	////"libavutil/hwcontext_dxva2.h", // needs <d3d9.h>/<dxva2api.h> (Windows SDK)
	////"libavutil/hwcontext_mediacodec.h", // Android-only, link risk
	////"libavutil/hwcontext_oh.h", // OpenHarmony-only, link risk
	////"libavutil/hwcontext_opencl.h", // needs <CL/cl.h>
	////"libavutil/hwcontext_qsv.h",
	////"libavutil/hwcontext_vaapi.h", // needs <va/va.h>
	////"libavutil/hwcontext_vdpau.h", // needs <vdpau/vdpau.h>
	////"libavutil/hwcontext_videotoolbox.h",
	////"libavutil/hwcontext_vulkan.h",
	"libavutil/iamf.h",
	"libavutil/imgutils.h",
	"libavutil/intfloat.h",
	////"libavutil/intreadwrite.h", //Union types - CGO doesn't expose union fields
	"libavutil/lfg.h",
	"libavutil/log.h",
	"libavutil/lzo.h",
	"libavutil/macros.h",
	"libavutil/mastering_display_metadata.h",
	"libavutil/mathematics.h",
	"libavutil/md5.h",
	"libavutil/mem.h",
	"libavutil/motion_vector.h",
	"libavutil/murmur3.h",
	"libavutil/opt.h",
	"libavutil/parseutils.h",
	"libavutil/pixdesc.h",
	"libavutil/pixelutils.h",
	"libavutil/pixfmt.h",
	"libavutil/random_seed.h",
	"libavutil/rational.h",
	"libavutil/rc4.h",
	////"libavutil/refstruct.h", // reference-counted object API introduced as an alternative to AVBuffer for managing complex objects
	"libavutil/replaygain.h",
	"libavutil/ripemd.h",
	"libavutil/samplefmt.h",
	"libavutil/sha.h",
	"libavutil/sha512.h",
	"libavutil/spherical.h",
	"libavutil/stereo3d.h",
	"libavutil/tdrdi.h",
	"libavutil/tea.h",
	"libavutil/threadmessage.h",
	"libavutil/time.h",
	"libavutil/timecode.h",
	"libavutil/timestamp.h",
	"libavutil/tree.h",
	"libavutil/twofish.h",
	"libavutil/tx.h",
	"libavutil/uuid.h",
	"libavutil/version.h",
	"libavutil/video_enc_params.h",
	"libavutil/video_hint.h",
	"libavutil/xtea.h",
	"libswresample/version.h",
	"libswresample/version_major.h",
	"libswresample/swresample.h",
	"libswscale/version.h",
	"libswscale/version_major.h",
	"libswscale/swscale.h",
}

func processComment(in string) string {
	txt := in
	txt = strings.ReplaceAll(txt, "\r\n", "\n")
	txt = strings.TrimSpace(txt)
	txt = strings.TrimPrefix(txt, "/**\n")
	txt = strings.TrimSuffix(txt, "*/")

	txt = strings.TrimRightFunc(txt, unicode.IsSpace)

	var rebuilt []string

	for s := range strings.SplitSeq(txt, "\n") {
		s = strings.TrimSpace(s)
		s = strings.TrimPrefix(s, "* ")
		s = strings.TrimPrefix(s, "/// ")
		s = strings.TrimPrefix(s, "///")
		s = strings.TrimPrefix(s, "// ")
		s = strings.TrimPrefix(s, "//")

		if strings.HasPrefix(s, "/**") {
			rebuilt = nil
			s = strings.TrimPrefix(s, "/** ")
			s = strings.TrimPrefix(s, "/**")

			if strings.TrimSpace(s) == "" {
				continue
			}
		}

		if strings.HasPrefix(s, "/*") {
			rebuilt = nil
			s = strings.TrimPrefix(s, "/* ")
			s = strings.TrimPrefix(s, "/*")

			if strings.TrimSpace(s) == "" {
				continue
			}
		}

		if strings.HasPrefix(s, "@defgroup") || strings.HasPrefix(s, "@ingroup") ||
			strings.HasPrefix(s, "@addtogroup") || strings.HasPrefix(s, "@}") {
			continue
		}

		if strings.HasPrefix(s, "@{") {
			rebuilt = nil
			continue
		}

		if s == "*" {
			s = ""
		}

		// Double space to enter verbatim mode
		rebuilt = append(rebuilt, fmt.Sprintf("  %v", s))
	}

	txt = strings.Join(rebuilt, "\n")
	txt = strings.TrimRightFunc(txt, unicode.IsSpace)

	if strings.Count(txt, "\n") == 0 {
		txt = strings.TrimLeftFunc(txt, unicode.IsSpace)
		txt = strings.TrimPrefix(txt, "<")
		txt = fmt.Sprintf("  %v", txt)
	}

	if strings.TrimSpace(txt) == "" {
		return ""
	}

	return txt
}
