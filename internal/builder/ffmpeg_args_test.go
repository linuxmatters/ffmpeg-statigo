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
