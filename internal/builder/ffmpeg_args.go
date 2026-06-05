package main

import "strings"

type ffmpegFeatureSet struct {
	Encoders         []string
	Decoders         []string
	Parsers          []string
	Demuxers         []string
	Muxers           []string
	BitstreamFilters []string
	InputDevices     []string
	OutputDevices    []string
	HWAccels         []string
}

func (set ffmpegFeatureSet) appendArgs(args []string) []string {
	args = appendEnableArg(args, "encoder", set.Encoders)
	args = appendEnableArg(args, "decoder", set.Decoders)
	args = appendEnableArg(args, "parser", set.Parsers)
	args = appendEnableArg(args, "demuxer", set.Demuxers)
	args = appendEnableArg(args, "muxer", set.Muxers)
	args = appendEnableArg(args, "bsf", set.BitstreamFilters)
	args = appendEnableArg(args, "indev", set.InputDevices)
	args = appendEnableArg(args, "outdev", set.OutputDevices)
	args = appendEnableArg(args, "hwaccel", set.HWAccels)
	return args
}

func appendEnableArg(args []string, kind string, values []string) []string {
	if len(values) == 0 {
		return args
	}
	return append(args, "--enable-"+kind+"="+strings.Join(values, ","))
}

func appendFeatureSets(args []string, sets []ffmpegFeatureSet) []string {
	for _, set := range sets {
		args = set.appendArgs(args)
	}
	return args
}

var commonFFmpegFeatureSets = []ffmpegFeatureSet{
	{BitstreamFilters: []string{"chomp", "dump_extra", "mov2textsub", "noise", "null", "remove_extra", "setts", "showinfo", "text2movsub"}},
	{
		Encoders:         []string{"aac"},
		Decoders:         []string{"aac", "aac_fixed", "aac_latm"},
		Parsers:          []string{"aac", "aac_latm"},
		Demuxers:         []string{"adts", "dash", "hds", "hls", "ipod", "ismv", "latm", "mov", "mp4", "rtp_mpegts", "rtsp", "sap", "smoothstreaming"},
		Muxers:           []string{"adts", "dash", "hds", "hls", "ipod", "ismv", "latm", "mov", "mp4", "rtp_mpegts", "rtsp", "sap", "smoothstreaming"},
		BitstreamFilters: []string{"aac_adtstoasc"},
	},
	{
		Encoders:         []string{"ac3", "ac3_fixed", "eac3"},
		Decoders:         []string{"ac3", "ac3_fixed", "eac3"},
		Parsers:          []string{"ac3"},
		Demuxers:         []string{"ac3", "eac3", "matroska", "spdif"},
		Muxers:           []string{"ac3", "eac3", "matroska", "spdif"},
		BitstreamFilters: []string{"eac3_core"},
	},
	{
		Encoders: []string{"alac"},
		Decoders: []string{"alac"},
	},
	{
		Encoders:         []string{"av1_vulkan", "librav1e"},
		Decoders:         []string{"av1", "libdav1d"},
		Parsers:          []string{"av1"},
		Demuxers:         []string{"avif", "obu"},
		Muxers:           []string{"avif", "obu"},
		BitstreamFilters: []string{"av1_frame_merge", "av1_frame_split", "av1_metadata", "dovi_rpu", "extract_extradata", "filter_units", "trace_headers"},
	},
	{
		Decoders: []string{"dirac"},
		Parsers:  []string{"dirac"},
		Demuxers: []string{"dirac"},
	},
	{
		Encoders: []string{"cfhd"},
		Decoders: []string{"cfhd"},
	},
	{
		Encoders: []string{"dnxhd"},
		Decoders: []string{"dnxhd"},
		Parsers:  []string{"dnxhd"},
		Demuxers: []string{"dnxhd"},
		Muxers:   []string{"dnxhd"},
	},
	{
		Encoders:         []string{"dca", "dts"},
		Decoders:         []string{"dca", "dts"},
		Demuxers:         []string{"dts", "dtshd", "matroska", "mov", "mp4", "rtp_mpegts"},
		Muxers:           []string{"dts", "matroska", "mov", "mp4", "rtp_mpegts"},
		BitstreamFilters: []string{"dca_core", "dts2pts"},
	},
	{
		Encoders: []string{"dvbsub", "dvdsub"},
		Decoders: []string{"dvbsub", "dvdsub"},
	},
	{
		Encoders:         []string{"eac3"},
		Decoders:         []string{"eac3"},
		Demuxers:         []string{"eac3"},
		Muxers:           []string{"eac3"},
		BitstreamFilters: []string{"eac3_core"},
	},
	{
		Encoders: []string{"exr"},
		Decoders: []string{"exr"},
	},
	{
		Encoders: []string{"ffv1", "ffv1_vulkan"},
		Decoders: []string{"ffv1"},
		Parsers:  []string{"ffv1"},
	},
	{
		Encoders: []string{"flac"},
		Decoders: []string{"flac"},
		Parsers:  []string{"flac"},
		Demuxers: []string{"flac", "oga", "ogg", "ogv"},
		Muxers:   []string{"flac", "oga", "ogg", "ogv"},
	},
	{
		Encoders: []string{"gif"},
		Decoders: []string{"gif"},
		Parsers:  []string{"gif"},
		Demuxers: []string{"gif"},
		Muxers:   []string{"gif"},
	},
	{
		Encoders:         []string{"h264_vulkan", "libx264", "libx264rgb"},
		Decoders:         []string{"h264"},
		Parsers:          []string{"h264"},
		Demuxers:         []string{"dash", "h264", "hds", "hls", "ismv", "matroska", "mov", "mp4", "smoothstreaming", "whip"},
		Muxers:           []string{"dash", "h264", "hds", "hls", "ismv", "matroska", "mov", "mp4", "smoothstreaming", "whip"},
		BitstreamFilters: []string{"dts2pts", "extract_extradata", "filter_units", "h264_metadata", "h264_mp4toannexb", "h264_redundant_pps", "trace_headers"},
	},
	{
		Encoders:         []string{"hevc_vulkan", "libx265"},
		Decoders:         []string{"hevc"},
		Parsers:          []string{"hevc"},
		Demuxers:         []string{"hevc"},
		Muxers:           []string{"hevc"},
		BitstreamFilters: []string{"dovi_rpu", "extract_extradata", "filter_units", "hevc_metadata", "hevc_mp4toannexb", "trace_headers"},
	},
	{
		Encoders:         []string{"mjpeg"},
		Decoders:         []string{"mjpeg", "mjpegb"},
		Parsers:          []string{"mjpeg"},
		Demuxers:         []string{"image2", "image2pipe", "mjpeg", "mpjpeg"},
		Muxers:           []string{"image2", "image2pipe", "mjpeg", "mpjpeg"},
		BitstreamFilters: []string{"mjpeg2jpeg", "mjpegadump"},
	},
	{
		Encoders: []string{"mp2", "mp2fixed"},
		Decoders: []string{"mp2", "mp2float"},
		Demuxers: []string{"dvd", "mp2", "mpegts", "vob"},
		Muxers:   []string{"dvd", "mp2", "mpegts", "vob"},
	},
	{
		Encoders: []string{"libmp3lame"},
		Decoders: []string{"mp3", "mp3float"},
		Demuxers: []string{"avi", "mp3"},
		Muxers:   []string{"avi", "mp3"},
	},
	{
		Encoders:         []string{"mpeg2video"},
		Decoders:         []string{"mpeg2video", "mpegvideo"},
		Demuxers:         []string{"dvd", "mpeg2video", "mpegts", "vob"},
		Muxers:           []string{"dvd", "mpeg2video", "mpegts", "vob"},
		BitstreamFilters: []string{"filter_units", "imxdump", "mov2textsub", "mpeg2_metadata", "trace_headers"},
	},
	{
		Decoders:         []string{"mpeg4"},
		Demuxers:         []string{"avi", "m4v", "rtp", "rtp_mpegts", "rtsp"},
		Muxers:           []string{"avi", "m4v", "rtp", "rtp_mpegts", "rtsp"},
		BitstreamFilters: []string{"extract_extradata", "mpeg4_unpack_bframes"},
	},
	{
		Encoders: []string{"anull", "vnull"},
		Decoders: []string{"anull", "vnull"},
	},
	{
		Encoders:         []string{"libopus", "opus"},
		Decoders:         []string{"libopus", "opus"},
		Parsers:          []string{"opus"},
		Demuxers:         []string{"iamf", "opus", "webm", "whip"},
		Muxers:           []string{"iamf", "opus", "webm", "whip"},
		BitstreamFilters: []string{"opus_metadata"},
	},
	{
		Encoders: []string{"pbm", "ppm"},
		Decoders: []string{"pbm", "ppm"},
	},
	{
		Encoders:         []string{"pcm_alaw", "pcm_mulaw"},
		Decoders:         []string{"pcm_alaw", "pcm_mulaw"},
		Demuxers:         []string{"alaw", "mulaw", "rtp"},
		Muxers:           []string{"alaw", "mulaw", "rtp"},
		BitstreamFilters: []string{"pcm_rechunk"},
	},
	{
		Encoders:         []string{"pcm_f32be", "pcm_f32le"},
		Decoders:         []string{"pcm_f32be", "pcm_f32le"},
		BitstreamFilters: []string{"pcm_rechunk"},
	},
	{
		Encoders:         []string{"pcm_s16be", "pcm_s16be_planar"},
		Decoders:         []string{"pcm_s16be", "pcm_s16be_planar"},
		Demuxers:         []string{"aiff", "caf", "s16be"},
		Muxers:           []string{"aiff", "caf", "s16be"},
		BitstreamFilters: []string{"pcm_rechunk"},
	},
	{
		Encoders:         []string{"pcm_s16le", "pcm_s16le_planar"},
		Decoders:         []string{"pcm_s16le", "pcm_s16le_planar"},
		Demuxers:         []string{"null", "s16le", "wav"},
		Muxers:           []string{"null", "s16le", "wav"},
		BitstreamFilters: []string{"pcm_rechunk"},
	},
	{
		Encoders:         []string{"pcm_s24be", "pcm_s24le", "pcm_s24le_planar"},
		Decoders:         []string{"pcm_s24be", "pcm_s24le", "pcm_s24le_planar"},
		BitstreamFilters: []string{"pcm_rechunk"},
	},
	{
		Encoders: []string{"pcm_s32be", "pcm_s32le", "pcm_s32le_planar"},
		Decoders: []string{"pcm_s32be", "pcm_s32le", "pcm_s32le_planar"},
	},
	{
		Encoders: []string{"apng", "png"},
		Decoders: []string{"apng", "png"},
		Parsers:  []string{"png"},
		Demuxers: []string{"apng"},
		Muxers:   []string{"apng"},
	},
	{
		Encoders:         []string{"prores", "prores_aw", "prores_ks"},
		Decoders:         []string{"prores", "prores_raw"},
		Parsers:          []string{"prores_raw"},
		BitstreamFilters: []string{"prores_metadata"},
	},
	{
		Encoders: []string{"vorbis"},
		Decoders: []string{"theora", "vorbis"},
		Parsers:  []string{"theora", "vorbis"},
		Demuxers: []string{"ogg"},
		Muxers:   []string{"ogg"},
	},
	{
		Encoders: []string{"tiff"},
		Decoders: []string{"tiff"},
	},
	{
		Encoders:         []string{"truehd"},
		Decoders:         []string{"truehd"},
		Demuxers:         []string{"truehd"},
		Muxers:           []string{"truehd"},
		BitstreamFilters: []string{"truehd_core"},
	},
	{
		Decoders:         []string{"vc1"},
		Parsers:          []string{"vc1"},
		Demuxers:         []string{"vc1"},
		Muxers:           []string{"vc1"},
		BitstreamFilters: []string{"extract_extradata"},
	},
	{
		Decoders:         []string{"vp8"},
		Parsers:          []string{"vp8"},
		Muxers:           []string{"ogv"},
		BitstreamFilters: []string{"filter_units", "trace_headers"},
	},
	{
		Encoders:         []string{"libvpx_vp9"},
		Decoders:         []string{"libvpx_vp9", "vp9"},
		Parsers:          []string{"vp9"},
		Muxers:           []string{"webm"},
		BitstreamFilters: []string{"filter_units", "trace_headers", "vp9_metadata", "vp9_raw_reorder", "vp9_superframe", "vp9_superframe_split"},
	},
	{
		Decoders:         []string{"vvc"},
		Parsers:          []string{"vvc"},
		Demuxers:         []string{"vvc"},
		Muxers:           []string{"vvc"},
		BitstreamFilters: []string{"extract_extradata", "filter_units", "trace_headers", "vvc_metadata", "vvc_mp4toannexb"},
	},
	{
		Encoders: []string{"libwebp", "libwebp_anim"},
		Decoders: []string{"webp"},
		Parsers:  []string{"webp"},
		Muxers:   []string{"webp"},
	},
	{
		Encoders: []string{"webvtt"},
		Decoders: []string{"webvtt"},
		Demuxers: []string{"hls", "webvtt"},
		Muxers:   []string{"hls", "webm", "webvtt"},
	},
	{
		InputDevices:  []string{"v4l2"},
		OutputDevices: []string{"v4l2"},
	},
}

var linuxFFmpegFeatureSets = []ffmpegFeatureSet{
	{
		Encoders: []string{"av1_nvenc", "av1_qsv", "av1_vaapi"},
		Decoders: []string{"av1_cuvid", "av1_qsv"},
	},
	{
		Encoders: []string{"h264_nvenc", "h264_qsv", "h264_vaapi"},
		Decoders: []string{"h264_cuvid", "h264_qsv"},
	},
	{
		Encoders: []string{"hevc_nvenc", "hevc_qsv", "hevc_vaapi"},
		Decoders: []string{"hevc_cuvid", "hevc_qsv"},
	},
	{
		Encoders: []string{"mjpeg_qsv", "mjpeg_vaapi"},
		Decoders: []string{"mjpeg_cuvid", "mjpeg_qsv"},
	},
	{
		Encoders: []string{"mpeg2_qsv", "mpeg2_vaapi"},
		Decoders: []string{"mpeg2_cuvid", "mpeg2_qsv"},
	},
	{Decoders: []string{"mpeg4_cuvid"}},
	{Decoders: []string{"vc1_cuvid", "vc1_qsv"}},
	{
		Encoders: []string{"vp8_vaapi"},
		Decoders: []string{"vp8_cuvid", "vp8_qsv"},
	},
	{
		Encoders: []string{"vp9_qsv", "vp9_vaapi"},
		Decoders: []string{"vp9_cuvid", "vp9_qsv"},
	},
	{Decoders: []string{"vvc_qsv"}},
}

var darwinFFmpegFeatureSets = []ffmpegFeatureSet{
	{
		Encoders: []string{"h264_videotoolbox"},
		HWAccels: []string{"h264_videotoolbox"},
	},
	{
		Encoders: []string{"hevc_videotoolbox"},
		HWAccels: []string{"hevc_videotoolbox"},
	},
	{
		Encoders: []string{"prores_videotoolbox"},
		HWAccels: []string{"prores_videotoolbox"},
	},
	{HWAccels: []string{"av1_videotoolbox"}},
	{HWAccels: []string{"vp9_videotoolbox"}},
	{HWAccels: []string{"mpeg2_videotoolbox"}},
	{HWAccels: []string{"mpeg4_videotoolbox"}},
	{HWAccels: []string{"h263_videotoolbox"}},
	{HWAccels: []string{"mpeg1_videotoolbox"}},
}

// FFmpegArgsCommon returns common FFmpeg configure arguments for all platforms
// os parameter enables platform-specific hardware acceleration (linux, darwin)
func FFmpegArgsCommon(os string) []string {
	args := []string{
		"--enable-pic",
		"--enable-gpl",
		"--enable-version3",
		"--enable-static",
		"--disable-autodetect",
		"--disable-debug",
		"--disable-doc",
		"--disable-htmlpages",
		"--disable-manpages",
		"--disable-podpages",
		"--disable-programs",
		"--disable-txtpages",
		"--disable-everything",
		"--enable-filters",
		"--enable-protocols",
	}

	args = appendFeatureSets(args, commonFFmpegFeatureSets)
	args = append(args, "--disable-encoder=h263")

	if os == "linux" {
		args = appendFeatureSets(args, linuxFFmpegFeatureSets)
	}
	if os == "darwin" {
		args = appendFeatureSets(args, darwinFFmpegFeatureSets)
	}

	return args
}
