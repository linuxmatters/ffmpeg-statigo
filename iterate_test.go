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
			"hevc",
			"vp9",
			"av1",
			"mpeg2video",
			"mjpeg",
			"png",
		}

		for _, codec := range criticalDecoders {
			if !codecNames[codec] {
				t.Errorf("Critical video decoder %q not found in codec list", codec)
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
		if len(demuxerNames) < 30 {
			t.Errorf("Expected at least 30 demuxers, found %d", len(demuxerNames))
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
