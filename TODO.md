# TODO

## FFmpeg 7.1 Upgrade

To update from FFmpeg 6.1 to 7.1 in this project, you'll need to address several areas.

**Why 7.1 instead of 7.0?** Going directly to 7.1 is the same effort as 7.0 (both share the 60→61 major version bump), but 7.1 includes 6+ months of additional bug fixes and the VVC decoder promoted to stable. Latest point release is 7.1.2 (September 2025).

## 1. Update Build System

**In `internal/builder/main.go`:**

```go
// Line ~1090
func (b *Builder) buildFFmpeg() {
    zipPath := path.Join(downloadsDir, "ffmpeg.zip")
    buildPath := path.Join(buildDir, "ffmpeg")

    if !exists(zipPath) {
        // Change from 6.1 to 7.1.2
        download("https://codeload.github.com/FFmpeg/FFmpeg/zip/refs/heads/release/7.1", zipPath)
    }
    // ...existing code...
}
```

Expected library versions after upgrade:
```
FFmpeg 6.1.3 → FFmpeg 7.1.2:
libavutil      58.29.100 → 59.39.100
libavcodec     60.31.102 → 61.19.100  (MAJOR BUMP)
libavformat    60.16.100 → 61. 7.100  (MAJOR BUMP)
libavdevice    60. 3.100 → 61. 3.100  (MAJOR BUMP)
libavfilter     9.12.100 → 10. 4.100  (MAJOR BUMP)
libswscale      7. 5.100 →  8. 3.100  (MAJOR BUMP)
libswresample   4.12.100 →  5. 3.100  (MAJOR BUMP)
libpostproc    57. 3.100 → 58. 3.100  (MAJOR BUMP)
```

## 2. Critical API Removals & Breaking Changes

FFmpeg 7.0/7.1 removes all APIs deprecated before version 60. These are the **confirmed** breaking changes:

### 2.1 Channel Layout API (CRITICAL)

**Removed in FFmpeg 7:**
- `AVFrame.channel_layout` field (int64 channel mask) → **REMOVED**
- `AVFrame.channels` field (int channel count) → **REMOVED**
- Old `av_get_channel_layout()` functions → **REMOVED**

**Replacement (already in FFmpeg 6.1):**
- Use `AVFrame.ch_layout` (AVChannelLayout struct)
- Use `av_channel_layout_*()` family of functions
- ✅ **jivedrop already uses new API** (`AVChannelLayoutDescribe`, `AVChannelLayoutCopy`)

**Impact:** Generator should not expose removed fields. The new `AVChannelLayout` struct bindings should work correctly.

### 2.2 AVFrame Field Removals

**Removed deprecated fields:**
- `AVFrame.pkt_pts` → Replaced by `AVFrame.pts` (since FFmpeg 4.0)
- `AVFrame.top_field_first` → Deprecated field for interlaced video
- `AVFrame.interlaced_frame` → Use `AVFrame.flags` with `AV_FRAME_FLAG_INTERLACED`

**Impact:** Generator will skip these fields automatically. Verify no custom code references them.

### 2.3 AVPacket Field Changes

**Potential removals:**
- Check for deprecated packet fields that may reference old channel layout

**Impact:** Verify AVPacket struct generation after upgrade.

### 2.4 Codec API Changes

**CLI Options Removed (less relevant for library use):**
- `-map_channel` CLI option removed
- `-psnr` CLI option removed (use `-flags +psnr` instead)

**Library API:**
- Some encoder/decoder options may have changed
- QSV encoder default changed from VBR to CQP

**Impact:** Minimal for Go bindings - these are mostly CLI concerns.

### 2.5 Filter API Changes

**Deprecated filter options:**
- Some audio filter channel layout options changed
- May need to update filter graph construction

**Impact:** Low - most filters use modern APIs already.

## 3. Update Generated Bindings

After updating FFmpeg headers, regenerate bindings using the **fixed generator** (v0.6.1):

```bash
# Use justfile for generation
just generate

# Or run generator directly
go run internal/generator/*.go
```

This regenerates:
- `constants.gen.go` - Will add new FFmpeg 7.1 constants
- `enums.gen.go` - May have new enum values
- `functions.gen.go` - Updated function signatures, removed deprecated functions
- `structs.gen.go` - **Critical:** Updated AVFrame, AVPacket structs (removed fields)

**Expected changes:**
- AVFrame struct will lose `channel_layout`, `channels`, `pkt_pts` fields
- AVPacket struct may have field changes
- Some function signatures may update (especially channel layout related)
- New constants for FFmpeg 7.1 features (VVC codec IDs, etc.)

## 4. Update Version Constants

**In `ffmpeg_test.go`:**

```go
func TestVersions(t *testing.T) {
    // Update expected version constants for FFmpeg 7.1
    // LIBAVCODEC_VERSION_MAJOR will be 61 for FFmpeg 7.1
    // Exact values:
    // - libavutil:    59.39.100
    // - libavcodec:   61.19.100
    // - libavformat:  61. 7.100
    // - libavfilter:  10. 4.100

    assert.Equal(t, expectedCodecVersion, int(ffmpeg.AVCodecVersion()),
        "AVCodec version should match FFmpeg 7.1.2")
    assert.Equal(t, ffmpeg.LIBAVCodecVersionInt, int(ffmpeg.AVCodecVersion()),
        "AVCodec version func and const should match")
}
```

## 5. Update Documentation

**In `README.md`:**

Update library versions table:

```md
### Library versions

| Library       | Version   |
|---------------|-----------|
| FFmpeg        | 7.1.2     |
| aom           | 3.8.1     |
| opus          | 1.4       |
| vpx           | 1.14.0    |
| x264          | (latest)  |
| freetype      | 2.13.2    |
```

## 6. Verify jivedrop Compatibility

Your jivedrop encoder at `/home/martin/Development/linuxmatters/jivedrop/internal/encoder/encoder.go` **should work without changes** because:

✅ **Already using new APIs:**
- `AVChannelLayoutDescribe()` - Modern channel layout API
- `AVChannelLayoutCopy()` - Modern channel layout copying
- `av_packet_rescale_ts()` - Standard time base handling
- Modern AVCodecContext usage

🔍 **Things to verify after upgrade:**
1. MP3 encoder options still work (LAME encoder)
2. AAC encoder options still work (native AAC)
3. Metadata/ID3 tag writing still works
4. Frame/packet time base handling unchanged

## 7. Testing Strategy

**Phase 1: Basic Validation**
```bash
# 1. Rebuild static libraries
cd internal/builder
go run . --platforms=linux-amd64  # or your platform

# 2. Regenerate bindings
just generate

# 3. Run existing tests
just test
```

**Phase 2: Example Validation**
```bash
# Run transcode example
cd examples/transcode
go run main.go

# Run asciiplayer example
cd examples/asciiplayer
go run main.go

# Run metadata example
cd examples/metadata
go run main.go
```

**Phase 3: jivedrop Integration Test**
```bash
# Test jivedrop encoder with new FFmpeg
cd ~/Development/linuxmatters/jivedrop
just build
just test

# Manual test: encode a real podcast episode
./jivedrop --input test.flac --output test.mp3 --episode "LMP-XX" --title "Test"
```

## 8. Potential Issues & Mitigation

### Issue 1: Generator Parsing Errors
**Symptom:** Generator fails to parse new headers
**Fix:** Update go-clang/bootstrap if needed, check for new libclang compatibility issues
**Likelihood:** Low (generator fixes in v0.6.1 handle platform differences)

### Issue 2: Struct Field Access Breaks
**Symptom:** Compilation errors accessing removed AVFrame fields
**Fix:** Update code to use replacement fields (pts instead of pkt_pts, ch_layout instead of channel_layout)
**Likelihood:** Very low (your code already uses modern APIs)

### Issue 3: Encoder Option Changes
**Symptom:** Encoder initialization fails or produces different output
**Fix:** Check encoder documentation, update AVOptions as needed
**Likelihood:** Low (LAME/AAC encoders are stable)

### Issue 4: Example Code Breaks
**Symptom:** Examples fail to compile or run
**Fix:** Update examples to use new APIs
**Likelihood:** Medium (examples may use deprecated APIs for demonstration)

## 9. Recommended Upgrade Path

**Timeline: ~2-4 hours for careful upgrade**

1. **Preparation (30 min)**
   - Tag current v0.6.1 release (generator fixes)
   - Create `ffmpeg-7.1` feature branch
   - Document current test results as baseline

2. **Build FFmpeg 7.1.2 (1 hour)**
   - Update builder to download FFmpeg 7.1.2
   - Rebuild static libraries for your platform(s)
   - Verify libraries are correct version

3. **Regenerate Bindings (30 min)**
   - Run generator on new headers
   - Review generated file diffs (expect struct changes)
   - Commit generated files separately for clean diff

4. **Fix & Test (1-2 hours)**
   - Update version constants in tests
   - Fix any compilation errors
   - Run all tests and examples
   - Test jivedrop integration

5. **Release (30 min)**
   - Update README and CHANGELOG
   - Merge to master
   - Tag as v0.7.0 (major version for FFmpeg major version)
   - Update GitHub release notes

## 10. Success Criteria

✅ All tests pass
✅ All examples compile and run
✅ jivedrop builds and encodes correctly
✅ Static libraries are correct FFmpeg 7.1.2 version
✅ No deprecated API warnings during compilation
✅ Generated bindings are clean (no parser errors)

## 11. Post-Upgrade: Consider Additional Codecs

Based on current multimedia landscape and practical use cases, here are codecs/features that would add genuine value:

## High-Priority Additions

**AV1 Encoding (SVT-AV1 or rav1e)**
- Modern, royalty-free codec with excellent compression
- Growing adoption (YouTube, Netflix, streaming platforms)
- Superior to VP9, competitive with HEVC without licensing issues
- **Already have libaom decoder**, adding SVT-AV1 encoder would complete the story

**HEVC/H.265 (x265)**
- Industry standard for 4K/HDR content
- Better compression than H.264
- Widely used in streaming, broadcasting, mobile
- There's even a commented-out reference in your FFmpeg configure: `--enable-libx265`

**Zstandard (zstd)**
- Fast, modern compression for muxing/demuxing
- Better than zlib for many use cases
- Low overhead, excellent for real-time scenarios

## Medium-Priority Additions

**WebP/WebM image support**
- `libwebp` - web-optimised image format
- Useful for thumbnail generation, web workflows
- Already have VP9 (part of WebM ecosystem)

**Dolby Vision/HDR10+ metadata handling**
- Increasingly important for video processing
- Requires no additional codecs, just metadata support
- Enable with: `--enable-libdovi` (if available)

**Hardware acceleration**
- **VAAPI** (Linux Intel/AMD) - already common in CI
- **VideoToolbox** (macOS) - you have this enabled already
- **QSV** (Intel Quick Sync) - cross-platform
- Significantly faster encoding/decoding for supported codecs

**FDK-AAC**
- Higher quality AAC encoder than FFmpeg's native
- Important for audio workflows requiring AAC
- License-compatible with GPLv3

## Lower-Priority But Useful

**AV1 image formats**
- `libavif` - AVIF image support (AV1-based)
- Better compression than WebP/JPEG
- Growing browser support

**Codecs for specific workflows**
- **Xvid** - legacy MPEG-4 support
- **libkvazaar** - open HEVC encoder
- **libdav1d** - faster AV1 decoder (alternative to libaom)

**Subtitle/accessibility**
- Already have `libass` ✓
- Consider `libfribidi`, `harfbuzz` for complex text (you have these ✓)

## Probably Skip

**Patent/License concerns:**
- AC-3/E-AC-3 (Dolby) - licensing restrictions
- DTS - licensing issues
- ProRes encoding - Apple proprietary (decoder is fine)

**Niche/declining:**
- Real Media codecs
- WMV/WMA (unless targeting legacy Windows content)
- Indeo, Cinepak, etc.

## My Recommendations for Your Build

Given you're building a **Go FFmpeg binding for general multimedia use**, I'd prioritise:

1. **x265** - Essential for modern video
2. **SVT-AV1** - Future-proofing for royalty-free encoding
3. **libfdk-aac** - Professional audio quality
4. **libwebp** - Web workflows
5. **Hardware acceleration** (VAAPI on Linux, keep VideoToolbox on macOS)

These would cover:
- Modern streaming (AV1, HEVC)
- Web optimisation (WebP, AV1)
- Professional audio (FDK-AAC)
- Performance (HW accel)
- All royalty-free or GPL-compatible

Your current codec selection is already solid for **open source/royalty-free workflows**. The main gap is **HEVC/x265** for 4K content and **SVT-AV1** for next-gen encoding.
