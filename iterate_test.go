package ffmpeg

import (
	"slices"
	"testing"
	"unsafe"
)

// =============================================================================
// Test 5.1: Codec Iterator Exhaustiveness
// =============================================================================

// TestAVCodecIterate_FindsExpectedCodecs verifies that codec iteration finds
// the critical codecs that should be present in the FFmpeg build.
func TestAVCodecIterate_FindsExpectedCodecs(t *testing.T) {
	// Collect all codec names
	codecNames := make(map[string]bool)
	var opaque unsafe.Pointer

	for {
		codec := AVCodecIterate(&opaque)
		if codec == nil {
			break
		}
		name := codec.Name()
		if name != nil {
			codecNames[name.String()] = true
		}
	}

	t.Logf("Found %d codecs", len(codecNames))

	t.Run("finds_critical_video_decoders", func(t *testing.T) {
		// These decoders should be present in any reasonable FFmpeg build
		criticalDecoders := []string{
			"h264",
			"vp9",
			"mpeg2video",
			"mjpeg",
			"png",
		}

		if !embeddedBuild {
			criticalDecoders = append(criticalDecoders, "av1", "hevc")
		}
		for _, codec := range criticalDecoders {
			if !codecNames[codec] {
				t.Errorf("Critical video decoder %q not found in codec list", codec)
			}
		}
		// These names identify components that the embedded profile excludes.
		for _, name := range []string{"libdav1d", "librav1e", "av1_vulkan", "ffv1_vulkan", "h264_vulkan", "hevc_vulkan"} {
			if got := codecNames[name]; got != !embeddedBuild {
				t.Errorf("codec %q present = %t, want %t (embedded = %t)", name, got, !embeddedBuild, embeddedBuild)
			}
		}
		if embeddedBuild {
			for _, name := range []string{"av1", "cfhd", "dirac", "dnxhd", "exr", "hevc", "pbm", "prores", "prores_aw", "prores_ks", "prores_raw", "tiff", "vc1", "vvc"} {
				if codecNames[name] {
					t.Errorf("embedded build unexpectedly includes codec %s", name)
				}
			}
			// MPEG-4 Part 2 and Theora require the H.263 and VP3 decoders.
			for _, name := range []string{"mpeg4", "h263", "theora", "vp3", "ppm"} {
				if !codecNames[name] {
					t.Errorf("embedded build is missing retained codec %s", name)
				}
			}
			if !codecNames["libopenh264"] {
				t.Error("embedded build is missing OpenH264 encoding")
			}
			if AVCodecFindEncoder(AVCodecIdH263) != nil {
				t.Error("embedded build unexpectedly includes H.263 encoding")
			}
			if AVCodecFindDecoder(AVCodecIdH264) == nil || AVCodecFindDecoder(AVCodecIdAac) == nil || AVCodecFindEncoder(AVCodecIdAac) == nil {
				t.Error("embedded build must retain H.264 decoding and AAC decoding and encoding")
			}
			if codecNames["libx264"] || codecNames["libx265"] {
				t.Error("embedded build unexpectedly includes GPL encoders")
			}
		}
	})

	t.Run("finds_critical_audio_codecs", func(t *testing.T) {
		criticalAudio := []string{
			"aac",
			"mp3",
			"opus",
			"flac",
			"ac3",
			"pcm_s16le",
		}

		for _, codec := range criticalAudio {
			if !codecNames[codec] {
				t.Errorf("Critical audio codec %q not found in codec list", codec)
			}
		}
	})

	t.Run("finds_subtitle_codecs", func(t *testing.T) {
		subtitleCodecs := []string{
			"webvtt",
		}

		for _, codec := range subtitleCodecs {
			if !codecNames[codec] {
				t.Errorf("Subtitle codec %q not found in codec list", codec)
			}
		}
	})

	t.Run("iteration_returns_multiple_codecs", func(t *testing.T) {
		// Sanity check: should have a reasonable number of codecs
		if len(codecNames) < 50 {
			t.Errorf("Expected at least 50 codecs, found %d", len(codecNames))
		}
	})

	t.Run("codec_has_valid_properties", func(t *testing.T) {
		// Reset iteration
		var opaque unsafe.Pointer
		codec := AVCodecIterate(&opaque)

		if codec == nil {
			t.Fatal("First codec iteration returned nil")
		}

		// Check name is valid
		name := codec.Name()
		if name == nil || name.String() == "" {
			t.Error("First codec has empty name")
		}

		// Check type is valid (video, audio, subtitle, etc.)
		codecType := codec.Type()
		validTypes := []AVMediaType{
			AVMediaTypeVideo,
			AVMediaTypeAudio,
			AVMediaTypeSubtitle,
			AVMediaTypeData,
			AVMediaTypeAttachment,
		}

		found := slices.Contains(validTypes, codecType)

		if !found && codecType != AVMediaTypeUnknown {
			t.Errorf("Codec %s has invalid media type %d", name.String(), codecType)
		}
	})
}

// =============================================================================
// Test 5.2: Muxer/Demuxer Iterator Completeness
// =============================================================================

// TestAVMuxerIterate_FindsExpectedFormats verifies that muxer iteration finds
// the critical container formats that should be present.
func TestAVMuxerIterate_FindsExpectedFormats(t *testing.T) {
	// Collect all muxer names
	muxerNames := make(map[string]bool)
	var opaque unsafe.Pointer

	for {
		muxer := AVMuxerIterate(&opaque)
		if muxer == nil {
			break
		}
		name := muxer.Name()
		if name != nil {
			muxerNames[name.String()] = true
		}
	}

	t.Logf("Found %d muxers", len(muxerNames))
	for _, name := range []string{"avif", "obu"} {
		if got := muxerNames[name]; got != !embeddedBuild {
			t.Errorf("muxer %q present = %t, want %t (embedded = %t)", name, got, !embeddedBuild, embeddedBuild)
		}
	}
	if embeddedBuild {
		for _, name := range []string{"dirac", "dnxhd", "h263", "hevc", "rm", "vc1", "vvc"} {
			if muxerNames[name] {
				t.Errorf("embedded build unexpectedly includes muxer %s", name)
			}
		}
		for _, name := range []string{"rtp", "rtsp", "sap", "rtp_mpegts"} {
			if !muxerNames[name] {
				t.Errorf("embedded build is missing retained muxer %s", name)
			}
		}
	}

	t.Run("finds_critical_video_containers", func(t *testing.T) {
		criticalFormats := []string{
			"mp4",
			"webm",
			"matroska",
			"mov",
			"avi",
			"mpegts",
			"hls",
		}

		for _, format := range criticalFormats {
			if !muxerNames[format] {
				t.Errorf("Critical video container %q not found in muxer list", format)
			}
		}
	})

	t.Run("finds_audio_formats", func(t *testing.T) {
		audioFormats := []string{
			"mp3",
			"flac",
			"ogg",
			"wav",
		}

		for _, format := range audioFormats {
			if !muxerNames[format] {
				t.Errorf("Audio format %q not found in muxer list", format)
			}
		}
	})

	t.Run("finds_streaming_formats", func(t *testing.T) {
		streamingFormats := []string{
			"hls",
			"dash",
			"rtp",
		}

		for _, format := range streamingFormats {
			if !muxerNames[format] {
				t.Errorf("Streaming format %q not found in muxer list", format)
			}
		}
	})

	t.Run("iteration_returns_multiple_muxers", func(t *testing.T) {
		if len(muxerNames) < 30 {
			t.Errorf("Expected at least 30 muxers, found %d", len(muxerNames))
		}
	})
}

// TestAVDemuxerIterate_FindsExpectedFormats verifies that demuxer iteration
// finds the critical input formats.
func TestAVDemuxerIterate_FindsExpectedFormats(t *testing.T) {
	// Collect all demuxer names
	demuxerNames := make(map[string]bool)
	var opaque unsafe.Pointer

	for {
		demuxer := AVDemuxerIterate(&opaque)
		if demuxer == nil {
			break
		}
		name := demuxer.Name()
		if name != nil {
			demuxerNames[name.String()] = true
		}
	}

	t.Logf("Found %d demuxers", len(demuxerNames))
	// The default archive uses the MOV demuxer for AVIF and has no separate "avif" entry.
	if embeddedBuild && demuxerNames["avif"] {
		t.Error("embedded build unexpectedly includes AVIF demuxer")
	}
	if got := demuxerNames["obu"]; got != !embeddedBuild {
		t.Errorf("demuxer %q present = %t, want %t (embedded = %t)", "obu", got, !embeddedBuild, embeddedBuild)
	}
	if embeddedBuild {
		for _, name := range []string{"dirac", "dnxhd", "h263", "hevc", "rm", "vc1", "vvc", "rtp", "rtsp", "sap", "sdp"} {
			if demuxerNames[name] {
				t.Errorf("embedded build unexpectedly includes demuxer %s", name)
			}
		}
	}

	t.Run("finds_critical_input_formats", func(t *testing.T) {
		criticalFormats := []string{
			"mov,mp4,m4a,3gp,3g2,mj2", // QuickTime/MP4 demuxer
			"matroska,webm",
			"avi",
			"mpegts",
			"ogg",
			"flac",
			"mp3",
			"wav",
		}

		for _, format := range criticalFormats {
			if !demuxerNames[format] {
				t.Errorf("Critical input format %q not found in demuxer list", format)
			}
		}
	})

	t.Run("iteration_returns_multiple_demuxers", func(t *testing.T) {
		minimum := 30
		if embeddedBuild {
			minimum = 25
		}
		if len(demuxerNames) < minimum {
			t.Errorf("Expected at least %d demuxers, found %d", minimum, len(demuxerNames))
		}
	})
}

// TestAVFilterIterate_FindsExpectedFilters verifies that filter iteration
// finds critical filters used for video/audio processing.
func TestAVFilterIterate_FindsExpectedFilters(t *testing.T) {
	// Collect all filter names
	filterNames := make(map[string]bool)
	var opaque unsafe.Pointer

	for {
		filter := AVFilterIterate(&opaque)
		if filter == nil {
			break
		}
		name := filter.Name()
		if name != nil {
			filterNames[name.String()] = true
		}
	}

	t.Logf("Found %d filters", len(filterNames))

	t.Run("finds_essential_video_filters", func(t *testing.T) {
		essentialFilters := []string{
			"scale",
			"format",
			"null",
			"fps",
		}

		for _, filter := range essentialFilters {
			if !filterNames[filter] {
				t.Errorf("Essential video filter %q not found in filter list", filter)
			}
		}
	})

	t.Run("finds_essential_audio_filters", func(t *testing.T) {
		essentialFilters := []string{
			"aformat",
			"anull",
			"volume",
		}

		for _, filter := range essentialFilters {
			if !filterNames[filter] {
				t.Errorf("Essential audio filter %q not found in filter list", filter)
			}
		}
	})

	t.Run("finds_buffer_filters", func(t *testing.T) {
		// Buffer filters are essential for filter graph construction
		bufferFilters := []string{
			"buffer",
			"buffersink",
			"abuffer",
			"abuffersink",
		}

		for _, filter := range bufferFilters {
			if !filterNames[filter] {
				t.Errorf("Buffer filter %q not found in filter list", filter)
			}
		}
	})

	t.Run("iteration_returns_multiple_filters", func(t *testing.T) {
		if len(filterNames) < 50 {
			t.Errorf("Expected at least 50 filters, found %d", len(filterNames))
		}
	})
}

// TestAVBSFIterate_FindsBitstreamFilters verifies that bitstream filter
// iteration works correctly.
func TestAVBSFIterate_FindsBitstreamFilters(t *testing.T) {
	// Collect all BSF names
	bsfNames := make(map[string]bool)
	var opaque unsafe.Pointer

	for {
		bsf := AVBSFIterate(&opaque)
		if bsf == nil {
			break
		}
		name := bsf.Name()
		if name != nil {
			bsfNames[name.String()] = true
		}
	}

	t.Logf("Found %d bitstream filters", len(bsfNames))
	for _, name := range []string{"av1_frame_merge", "av1_frame_split", "av1_metadata"} {
		if got := bsfNames[name]; got != !embeddedBuild {
			t.Errorf("bitstream filter %q present = %t, want %t (embedded = %t)", name, got, !embeddedBuild, embeddedBuild)
		}
	}
	if embeddedBuild {
		for _, name := range []string{"hevc_metadata", "prores_metadata", "vvc_metadata", "dovi_rpu"} {
			if bsfNames[name] {
				t.Errorf("embedded build unexpectedly includes bitstream filter %s", name)
			}
		}
		for _, name := range []string{"dts2pts", "hevc_mp4toannexb", "vvc_mp4toannexb", "extract_extradata", "filter_units", "trace_headers"} {
			if !bsfNames[name] {
				t.Errorf("embedded build is missing retained bitstream filter %s", name)
			}
		}
	}

	t.Run("finds_common_bitstream_filters", func(t *testing.T) {
		commonBSF := []string{
			"null",
			"h264_mp4toannexb",
		}

		for _, bsf := range commonBSF {
			if !bsfNames[bsf] {
				t.Errorf("Common bitstream filter %q not found", bsf)
			}
		}
	})

	t.Run("iteration_returns_multiple_bsf", func(t *testing.T) {
		if len(bsfNames) < 5 {
			t.Errorf("Expected at least 5 bitstream filters, found %d", len(bsfNames))
		}
	})
}

func TestAVParserProfileSelection(t *testing.T) {
	for _, id := range []AVCodecID{AVCodecIdAV1, AVCodecIdH264, AVCodecIdHevc, AVCodecIdDirac, AVCodecIdDnxhd, AVCodecIdProresRaw, AVCodecIdVc1, AVCodecIdVvc} {
		parser := AVParserInit(int(id))
		present := parser != nil
		if present {
			AVParserClose(parser)
		}
		// dts2pts requires the HEVC parser even without HEVC decoding.
		want := !embeddedBuild || id == AVCodecIdH264 || id == AVCodecIdHevc
		if present != want {
			t.Errorf("parser for codec %d present = %v, want %v", id, present, want)
		}
	}
	if embeddedBuild {
		for _, id := range []AVCodecID{AVCodecIdH263, AVCodecIdProres, AVCodecIdVp3} {
			if parser := AVParserInit(int(id)); parser != nil {
				AVParserClose(parser)
				t.Errorf("embedded build unexpectedly includes parser for codec %d", id)
			}
		}
	}
}

// TestAVChannelLayoutStandard_IteratesStandardLayouts verifies that standard
// channel-layout iteration enumerates the well-known layouts and terminates.
func TestAVChannelLayoutStandard_IteratesStandardLayouts(t *testing.T) {
	count := 0
	channelCounts := make(map[int]bool)
	var opaque unsafe.Pointer

	for {
		layout := AVChannelLayoutStandard(&opaque)
		if layout == nil {
			break
		}
		count++
		channelCounts[layout.NbChannels()] = true

		if count > 1000 {
			t.Fatal("iteration did not terminate")
		}
	}

	t.Logf("Found %d standard channel layouts", count)

	if count < 5 {
		t.Errorf("Expected at least 5 standard layouts, found %d", count)
	}

	t.Run("includes_mono_and_stereo", func(t *testing.T) {
		if !channelCounts[1] {
			t.Error("Expected a 1-channel (mono) standard layout")
		}
		if !channelCounts[2] {
			t.Error("Expected a 2-channel (stereo) standard layout")
		}
	})
}

// TestAVIOEnumProtocols_FindsExpectedProtocols verifies that protocol
// enumeration finds the expected I/O protocols.
func TestAVIOEnumProtocols_FindsExpectedProtocols(t *testing.T) {
	t.Run("finds_input_protocols", func(t *testing.T) {
		protocolNames := make(map[string]bool)
		var opaque unsafe.Pointer

		for {
			name := AVIOEnumProtocols(&opaque, 0) // 0 = input
			if name == "" {
				break
			}
			protocolNames[name] = true
		}

		t.Logf("Found %d input protocols", len(protocolNames))

		expectedProtocols := []string{
			"file",
			"pipe",
		}

		for _, proto := range expectedProtocols {
			if !protocolNames[proto] {
				t.Errorf("Input protocol %q not found", proto)
			}
		}
	})

	t.Run("finds_output_protocols", func(t *testing.T) {
		protocolNames := make(map[string]bool)
		var opaque unsafe.Pointer

		for {
			name := AVIOEnumProtocols(&opaque, 1) // 1 = output
			if name == "" {
				break
			}
			protocolNames[name] = true
		}

		t.Logf("Found %d output protocols", len(protocolNames))

		expectedProtocols := []string{
			"file",
			"pipe",
		}

		for _, proto := range expectedProtocols {
			if !protocolNames[proto] {
				t.Errorf("Output protocol %q not found", proto)
			}
		}
	})
}
