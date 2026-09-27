package main

import (
	"slices"
	"testing"
)

func TestFFmpegFeatureSetAppendArgs(t *testing.T) {
	set := ffmpegFeatureSet{
		Encoders:         []string{"libx264", "libx264rgb"},
		Decoders:         []string{"h264"},
		Parsers:          []string{"h264"},
		Demuxers:         []string{"h264", "mov"},
		Muxers:           []string{"h264", "mov"},
		BitstreamFilters: []string{"h264_metadata", "trace_headers"},
		InputDevices:     []string{"v4l2"},
		OutputDevices:    []string{"v4l2"},
		HWAccels:         []string{"h264_videotoolbox"},
	}

	got := set.appendArgs(nil)
	want := []string{
		"--enable-encoder=libx264,libx264rgb",
		"--enable-decoder=h264",
		"--enable-parser=h264",
		"--enable-demuxer=h264,mov",
		"--enable-muxer=h264,mov",
		"--enable-bsf=h264_metadata,trace_headers",
		"--enable-indev=v4l2",
		"--enable-outdev=v4l2",
		"--enable-hwaccel=h264_videotoolbox",
	}

	if !slices.Equal(got, want) {
		t.Errorf("appendArgs() = %v, want %v", got, want)
	}
}

func TestEmbeddedComponentSelection(t *testing.T) {
	removed := map[string][]string{
		"encoder": {"cfhd", "dnxhd", "exr", "h263", "pbm", "prores", "prores_aw", "prores_ks", "tiff"},
		"decoder": {"cfhd", "dirac", "dnxhd", "exr", "hevc", "pbm", "prores", "prores_raw", "tiff", "vc1", "vvc"},
		"parser":  {"dirac", "dnxhd", "h263", "prores", "prores_raw", "vc1", "vp3", "vvc"},
		"demuxer": {"dirac", "dnxhd", "h263", "hevc", "rm", "vc1", "vvc", "rtp", "rtsp", "sap", "sdp"},
		"muxer":   {"dirac", "dnxhd", "h263", "hevc", "rm", "vc1", "vvc"},
		"bsf":     {"hevc_metadata", "prores_metadata", "vvc_metadata", "dovi_rpu"},
	}
	retained := map[string][]string{
		"encoder": {"libopenh264", "aac", "ppm", "ffv1", "libvpx_vp9", "libopus"},
		"decoder": {"h264", "aac", "mpeg4", "theora", "ppm"},
		"parser":  {"h264", "hevc"},
		"demuxer": {"mov", "mp4", "mpegts", "ogg"},
		"muxer":   {"mov", "mp4", "mpegts", "rtp", "rtsp", "sap", "rtp_mpegts"},
		"bsf":     {"dts2pts", "hevc_mp4toannexb", "vvc_mp4toannexb", "extract_extradata", "filter_units", "trace_headers"},
	}
	for _, targetOS := range []string{"linux", "darwin"} {
		t.Run(targetOS, func(t *testing.T) {
			defaults := FFmpegArgsCommon(targetOS)
			args := ffmpegArgsCommon(targetOS, true)
			disabled := func(kind, name string) bool {
				return featureConfigured(args, "disable", kind, name)
			}
			for kind, names := range removed {
				for _, name := range names {
					if featureEnabled(args, kind, name) || !disabled(kind, name) {
						t.Errorf("embedded %s %s must be explicitly disabled and not enabled", kind, name)
					}
					if (kind != "encoder" || name != "h263") && featureConfigured(defaults, "disable", kind, name) {
						t.Errorf("embedded exclusion affected default %s %s", kind, name)
					}
				}
			}
			for kind, names := range retained {
				for _, name := range names {
					if !featureEnabled(args, kind, name) || disabled(kind, name) {
						t.Errorf("embedded %s %s must remain enabled", kind, name)
					}
				}
			}
			for _, name := range []string{"h263", "vp3"} {
				if disabled("decoder", name) {
					t.Errorf("shared decoder %s must remain available for dependency selection", name)
				}
			}
			for _, set := range commonFFmpegFeatureSets {
				for _, arg := range set.appendArgs(nil) {
					if !slices.Contains(defaults, arg) {
						t.Errorf("default profile lost %s", arg)
					}
				}
			}
			if !slices.Equal(defaults, FFmpegArgsCommon(targetOS)) {
				t.Error("embedded configuration mutated the default profile")
			}
		})
	}
}

func TestEmbeddedFeatureSetMatchesExactComponents(t *testing.T) {
	got := embeddedFeatureSet(ffmpegFeatureSet{
		Encoders:         []string{"pbm", "ppm"},
		Decoders:         []string{"hevc", "h263", "vp3"},
		Parsers:          []string{"hevc", "h263"},
		Demuxers:         []string{"rtp", "rm"},
		Muxers:           []string{"rtp", "rtsp", "sap"},
		BitstreamFilters: []string{"hevc_metadata", "hevc_mp4toannexb", "vvc_mp4toannexb"},
	}).appendArgs(nil)
	want := []string{
		"--enable-encoder=ppm", "--enable-decoder=h263,vp3", "--enable-parser=hevc",
		"--enable-muxer=rtp,rtsp,sap", "--enable-bsf=hevc_mp4toannexb,vvc_mp4toannexb",
	}
	if !slices.Equal(got, want) {
		t.Errorf("embedded components = %v, want %v", got, want)
	}
}

func TestFFmpegArgsCommonPlatformOrdering(t *testing.T) {
	common := FFmpegArgsCommon("freebsd")
	linux := FFmpegArgsCommon("linux")
	darwin := FFmpegArgsCommon("darwin")

	if slices.Contains(common, "--enable-encoder=av1_nvenc,av1_qsv,av1_vaapi") {
		t.Error("generic args include linux-only AV1 hardware encoders")
	}
	if slices.Contains(common, "--enable-hwaccel=h264_videotoolbox") {
		t.Error("generic args include darwin-only VideoToolbox hwaccel")
	}

	wantDisable := slices.Index(common, "--disable-encoder=h263")
	if wantDisable == -1 {
		t.Fatal("generic args missing --disable-encoder=h263")
	}
	if wantDisable != len(common)-1 {
		t.Errorf("--disable-encoder=h263 index = %d, want generic tail index %d", wantDisable, len(common)-1)
	}

	linuxFirst := slices.Index(linux, "--enable-encoder=av1_nvenc,av1_qsv,av1_vaapi")
	if linuxFirst != len(common) {
		t.Errorf("first linux-only arg index = %d, want %d", linuxFirst, len(common))
	}

	darwinFirst := slices.Index(darwin, "--enable-encoder=h264_videotoolbox")
	if darwinFirst != len(common) {
		t.Errorf("first darwin-only arg index = %d, want %d", darwinFirst, len(common))
	}
}
