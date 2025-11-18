package main

// FFmpegArgsCommon returns common FFmpeg configure arguments for all platforms
// os parameter enables platform-specific hardware acceleration (linux, darwin)
func FFmpegArgsCommon(os string) []string {
	args := []string{
		// Build options
		"--enable-pic",
		"--enable-gpl",
		"--enable-version3",
		"--enable-static",

		// Disable everything first
		"--disable-autodetect",
		"--disable-debug",
		"--disable-doc",
		"--disable-htmlpages",
		"--disable-manpages",
		"--disable-podpages",
		"--disable-programs",
		"--disable-txtpages",
		"--disable-everything",

		// Enable all filters, hwaccels and protocols
		"--enable-filters",
		//"--enable-hwaccels",
		"--enable-protocols",

		// TEMPORARILY Enable all muxers, demuxers, parsers, bsfs, encoders, decoders, hwaccels
		//"--enable-muxers",
		//"--enable-demuxers",
		//"--enable-parsers",
		//"--enable-bsfs",
		//"--enable-encoders",
		//"--enable-decoders",

		// Common bitstream filters
		"--enable-bsf=chomp,dump_extra,mov2textsub,noise,null,remove_extra,setts,showinfo,text2movsub",

		// AAC (Advanced Audio Coding)
		// AAC LATM (Advanced Audio Coding LATM syntax)
		"--enable-encoder=aac",
		"--enable-decoder=aac,aac_fixed,aac_latm",
		"--enable-parser=aac,aac_latm",
		"--enable-demuxer=adts,dash,hds,hls,ipod,ismv,latm,mov,mp4,rtp_mpegts,rtsp,sap,smoothstreaming",
		"--enable-muxer=adts,dash,hds,hls,ipod,ismv,latm,mov,mp4,rtp_mpegts,rtsp,sap,smoothstreaming",
		"--enable-bsf=aac_adtstoasc",

		// ATSC A/52A (AC-3)
		// ATSC A/52 E-AC-3
		// ATRAC3 (Adaptive TRansform Acoustic Coding 3)
		"--enable-encoder=ac3,ac3_fixed,eac3",
		"--enable-decoder=ac3,ac3_fixed,eac3",
		"--enable-parser=ac3",
		"--enable-demuxer=ac3,eac3,matroska,spdif",
		"--enable-muxer=ac3,eac3,matroska,spdif",
		"--enable-bsf=eac3_core",

		// ALAC (Apple Lossless Audio Codec)
		"--enable-encoder=alac",
		"--enable-decoder=alac",

		// Alliance for Open Media AV1
		// librav1e AV1
		// dav1d AV1 decoder by VideoLAN
		// Nvidia CUVID AV1 decoder
		// AV1 (Vulkan)
		// AV1 (Intel Quick Sync Video acceleration)
		"--enable-encoder=av1_vulkan,librav1e",
		"--enable-decoder=av1,libdav1d",
		"--enable-parser=av1",
		"--enable-demuxer=avif,obu",
		"--enable-muxer=avif,obu",
		"--enable-bsf=av1_frame_merge,av1_frame_split,av1_metadata,dovi_rpu,extract_extradata,filter_units,trace_headers",

		// BBC Dirac VC-2
		"--enable-decoder=dirac",
		"--enable-parser=dirac",
		"--enable-demuxer=dirac",

		// GoPro CineForm HD
		"--enable-encoder=cfhd",
		"--enable-decoder=cfhd",

		// VC3/DNxHD
		"--enable-encoder=dnxhd",
		"--enable-decoder=dnxhd",
		"--enable-parser=dnxhd",
		"--enable-demuxer=dnxhd",
		"--enable-muxer=dnxhd",

		// DCA (DTS Coherent Acoustics)
		"--enable-encoder=dca,dts",
		"--enable-decoder=dca,dts",
		"--enable-demuxer=dts,dtshd,matroska,mov,mp4,rtp_mpegts",
		"--enable-muxer=dts,matroska,mov,mp4,rtp_mpegts",
		"--enable-bsf=dca_core,dts2pts",

		// DVB subtitles
		// DVD subtitles
		"--enable-encoder=dvbsub,dvdsub",
		"--enable-decoder=dvbsub,dvdsub",

		// ATSC A/52 E-AC-3
		"--enable-encoder=eac3",
		"--enable-decoder=eac3",
		"--enable-demuxer=eac3",
		"--enable-muxer=eac3",
		"--enable-bsf=eac3_core",

		// OpenEXR image
		"--enable-encoder=exr",
		"--enable-decoder=exr",

		// FFmpeg video codec #1
		// FFmpeg video codec #1 (Vulkan)
		"--enable-encoder=ffv1,ffv1_vulkan",
		"--enable-decoder=ffv1",
		"--enable-parser=ffv1",

		// FLAC (Free Lossless Audio Codec)
		"--enable-encoder=flac",
		"--enable-decoder=flac",
		"--enable-parser=flac",
		"--enable-demuxer=flac,oga,ogg,ogv",
		"--enable-muxer=flac,oga,ogg,ogv",

		// GIF (Graphics Interchange Format)
		"--enable-encoder=gif",
		"--enable-decoder=gif",
		"--enable-parser=gif",
		"--enable-demuxer=gif",
		"--enable-muxer=gif",

		// H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10
		// libx264 H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10
		// libx264 H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10 RGB
		// Nvidia CUVID H264 decoder
		// NVIDIA NVENC H.264 encoder
		// H.264/AVC (Vulkan)
		// H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10 (Intel Quick Sync Video acceleration)
		"--enable-encoder=h264_vulkan,libx264,libx264rgb",
		"--enable-decoder=h264",
		"--enable-parser=h264",
		"--enable-demuxer=dash,h264,hds,hls,ismv,matroska,mov,mp4,smoothstreaming,whip",
		"--enable-muxer=dash,h264,hds,hls,ismv,matroska,mov,mp4,smoothstreaming,whip",
		"--enable-bsf=dts2pts,extract_extradata,filter_units,h264_metadata,h264_mp4toannexb,h264_redundant_pps,trace_headers",

		// HEVC (High Efficiency Video Coding)
		// Nvidia CUVID HEVC decoder
		// libx265 H.265 / HEVC
		// H.265/HEVC (Vulkan)
		// NVIDIA NVENC hevc encoder
		// HEVC (Intel Quick Sync Video acceleration)
		"--enable-encoder=hevc_vulkan,libx265",
		"--enable-decoder=hevc",
		"--enable-parser=hevc",
		"--enable-demuxer=hevc",
		"--enable-muxer=hevc",
		"--enable-bsf=dovi_rpu,extract_extradata,filter_units,hevc_metadata,hevc_mp4toannexb,trace_headers",

		// MJPEG (Motion JPEG)
		// Apple MJPEG-B
		// Nvidia CUVID MJPEG decoder
		// Media 100
		// MJPEG (Intel Quick Sync Video acceleration)
		"--enable-encoder=mjpeg",
		"--enable-decoder=mjpeg,mjpegb",
		"--enable-parser=mjpeg",
		"--enable-demuxer=image2,image2pipe,mjpeg,mpjpeg",
		"--enable-muxer=image2,image2pipe,mjpeg,mpjpeg",
		"--enable-bsf=mjpeg2jpeg,mjpegadump",

		// MP2 (MPEG audio layer 2)
		// MP2 fixed point (MPEG audio layer 2)
		"--enable-encoder=mp2,mp2fixed",
		"--enable-decoder=mp2,mp2float",
		"--enable-demuxer=dvd,mp2,mpegts,vob",
		"--enable-muxer=dvd,mp2,mpegts,vob",

		// MP3 (MPEG audio layer 3)
		// libmp3lame MP3 (MPEG audio layer 3)
		// ADU (Application Data Unit) MP3 (MPEG audio layer 3)
		"--enable-encoder=libmp3lame",
		"--enable-decoder=mp3,mp3float",
		"--enable-demuxer=avi,mp3",
		"--enable-muxer=avi,mp3",

		// MPEG-2 video
		// Nvidia CUVID MPEG2VIDEO decoder
		// MPEG-2 video (Intel Quick Sync Video acceleration)
		"--enable-encoder=mpeg2video",
		"--enable-decoder=mpeg2video,mpegvideo",
		"--enable-demuxer=dvd,mpeg2video,mpegts,vob",
		"--enable-muxer=dvd,mpeg2video,mpegts,vob",
		"--enable-bsf=filter_units,imxdump,mov2textsub,mpeg2_metadata,trace_headers",

		// MPEG-4 part 2
		// Nvidia CUVID MPEG4 decoder
		// MPEG-4 part 2 Microsoft variant version 3
		"--enable-decoder=mpeg4",
		"--enable-demuxer=avi,m4v,rtp,rtp_mpegts,rtsp",
		"--enable-muxer=avi,m4v,rtp,rtp_mpegts,rtsp",
		"--enable-bsf=extract_extradata,mpeg4_unpack_bframes",

		// null audio
		// null video
		"--enable-encoder=anull,vnull",
		"--enable-decoder=anull,vnull",

		// Opus
		// libopus Opus
		"--enable-encoder=libopus,opus",
		"--enable-decoder=libopus,opus",
		"--enable-parser=opus",
		"--enable-demuxer=iamf,opus,webm,whip",
		"--enable-muxer=iamf,opus,webm,whip",
		"--enable-bsf=opus_metadata",

		// PBM (Portable BitMap) image
		// PPM (Portable PixelMap) image
		"--enable-encoder=pbm,ppm",
		"--enable-decoder=pbm,ppm",

		// PCM A-law / G.711 A-law
		// PCM mu-law / G.711 mu-law
		"--enable-encoder=pcm_alaw,pcm_mulaw",
		"--enable-decoder=pcm_alaw,pcm_mulaw",
		"--enable-demuxer=alaw,mulaw,rtp",
		"--enable-muxer=alaw,mulaw,rtp",
		"--enable-bsf=pcm_rechunk",

		// PCM 32-bit floating point big-endian
		// PCM 32-bit floating point little-endian
		"--enable-encoder=pcm_f32be,pcm_f32le",
		"--enable-decoder=pcm_f32be,pcm_f32le",
		"--enable-bsf=pcm_rechunk",

		// PCM signed 16-bit big-endian
		// PCM signed 16-bit big-endian planar
		"--enable-encoder=pcm_s16be,pcm_s16be_planar",
		"--enable-decoder=pcm_s16be,pcm_s16be_planar",
		"--enable-demuxer=aiff,caf,s16be",
		"--enable-muxer=aiff,caf,s16be",
		"--enable-bsf=pcm_rechunk",

		// PCM signed 16-bit little-endian
		// PCM signed 16-bit little-endian planar
		"--enable-encoder=pcm_s16le,pcm_s16le_planar",
		"--enable-decoder=pcm_s16le,pcm_s16le_planar",
		"--enable-demuxer=null,s16le,wav",
		"--enable-muxer=null,s16le,wav",
		"--enable-bsf=pcm_rechunk",

		// PCM signed 24-bit big-endian
		// PCM signed 24-bit little-endian
		// PCM signed 24-bit little-endian planar
		"--enable-encoder=pcm_s24be,pcm_s24le,pcm_s24le_planar",
		"--enable-decoder=pcm_s24be,pcm_s24le,pcm_s24le_planar",
		"--enable-bsf=pcm_rechunk",

		// PCM signed 32-bit big-endian
		// PCM signed 32-bit little-endian
		// PCM signed 32-bit little-endian planar
		"--enable-encoder=pcm_s32be,pcm_s32le,pcm_s32le_planar",
		"--enable-decoder=pcm_s32be,pcm_s32le,pcm_s32le_planar",

		// PNG (Portable Network Graphics) image
		// APNG (Animated Portable Network Graphics) image
		"--enable-encoder=apng,png",
		"--enable-decoder=apng,png",
		"--enable-parser=png",
		"--enable-demuxer=apng",
		"--enable-muxer=apng",

		// Apple ProRes
		// Apple ProRes (iCodec Pro)
		// Apple ProRes RAW
		"--enable-encoder=prores,prores_aw,prores_ks",
		"--enable-decoder=prores,prores_raw",
		"--enable-parser=prores_raw",
		"--enable-bsf=prores_metadata",

		// Theora
		// Vorbis
		"--enable-encoder=vorbis",
		"--enable-decoder=theora,vorbis",
		"--enable-parser=theora,vorbis",
		"--enable-demuxer=ogg",
		"--enable-muxer=ogg",

		// TIFF image
		"--enable-encoder=tiff",
		"--enable-decoder=tiff",

		// TrueHD
		"--enable-encoder=truehd",
		"--enable-decoder=truehd",
		"--enable-demuxer=truehd",
		"--enable-muxer=truehd",
		"--enable-bsf=truehd_core",

		// SMPTE VC-1
		// Nvidia CUVID VC1 decoder
		// Windows Media Video 9 Image v2
		// VC1 video (Intel Quick Sync Video acceleration)
		"--enable-decoder=vc1",
		"--enable-parser=vc1",
		"--enable-demuxer=vc1",
		"--enable-muxer=vc1",
		"--enable-bsf=extract_extradata",

		// On2 VP8
		// libvpx VP8
		// Nvidia CUVID VP8 decoder
		// VP8 video (Intel Quick Sync Video acceleration)
		"--enable-decoder=vp8",
		"--enable-parser=vp8",
		"--enable-muxer=ogv",
		"--enable-bsf=filter_units,trace_headers",

		// Google VP9
		// Nvidia CUVID VP9 decoder
		// libvpx VP9
		// VP9 video (Intel Quick Sync Video acceleration)
		"--enable-encoder=libvpx-vp9",
		"--enable-decoder=libvpx-vp9,vp9",
		"--enable-parser=vp9",
		"--enable-muxer=webm",
		"--enable-bsf=filter_units,trace_headers,vp9_metadata,vp9_raw_reorder,vp9_superframe,vp9_superframe_split",

		// VVC (Versatile Video Coding)
		// libvvenc H.266 / VVC
		// VVC video (Intel Quick Sync Video acceleration)
		//"--enable-encoder=libvvenc",
		"--enable-decoder=vvc",
		"--enable-parser=vvc",
		"--enable-demuxer=vvc",
		"--enable-muxer=vvc",
		"--enable-bsf=extract_extradata,filter_units,trace_headers,vvc_metadata,vvc_mp4toannexb",

		// WebP image
		// libwebp WebP image
		"--enable-encoder=libwebp,libwebp_anim",
		"--enable-decoder=webp",
		"--enable-parser=webp",
		"--enable-muxer=webp",

		// WebVTT subtitle
		"--enable-encoder=webvtt",
		"--enable-decoder=webvtt",
		"--enable-demuxer=hls,webvtt",
		"--enable-muxer=hls,webm,webvtt",

		// Devices
		"--enable-indev=v4l2",
		"--enable-outdev=v4l2",

		// Remove remaining unwanted cruft
		"--disable-encoder=h263",
	}

	// Linux-specific hardware acceleration (NVENC, CUVID, QuickSync)
	if os == "linux" {
		args = append(args,
			// AV1 hardware acceleration
			"--enable-encoder=av1_nvenc,av1_qsv",
			"--enable-decoder=av1_cuvid,av1_qsv",
			// H.264 hardware acceleration
			"--enable-encoder=h264_nvenc,h264_qsv",
			"--enable-decoder=h264_cuvid,h264_qsv",
			// H.265 hardware acceleration
			"--enable-encoder=hevc_nvenc,hevc_qsv",
			"--enable-decoder=hevc_cuvid,hevc_qsv",
			// MJPEG hardware acceleration
			"--enable-encoder=mjpeg_qsv",
			"--enable-decoder=mjpeg_cuvid,mjpeg_qsv",
			// MPEG-2 hardware acceleration
			"--enable-encoder=mpeg2_qsv",
			"--enable-decoder=mpeg2_cuvid,mpeg2_qsv",
			// MPEG-4 Part 2 hardware acceleration
			"--enable-decoder=mpeg4_cuvid",
			// VC-1 hardware acceleration
			"--enable-decoder=vc1_cuvid,vc1_qsv",
			// VP8 hardware acceleration
			"--enable-decoder=vp8_cuvid,vp8_qsv",
			// VP9 hardware acceleration
			"--enable-encoder=vp9_qsv",
			"--enable-decoder=vp9_cuvid,vp9_qsv",
			// VVC hardware acceleration
			"--enable-decoder=vvc_qsv",
		)
	}

	return args
}
