# TODO

## FFmpeg 8.0 Upgrade

To update from FFmpeg 7.1 to 8.0 in this project, you'll need to address several areas.

**Status:** FFmpeg 7.1.2 successfully completed (commits e5615cf, 34609cd). Ready for 8.0 upgrade.

**Why 8.0?** FFmpeg 8.0 "Huffman" (released August 22, 2025) is a major release with significant new features and codec support. It removes APIs deprecated before 6.0, following the new 3-major-release deprecation policy.

## 1. Update Build System

**In `internal/builder/main.go`:**

```go
// Line ~1090
func (b *Builder) buildFFmpeg() {
    zipPath := path.Join(downloadsDir, "ffmpeg.zip")
    buildPath := path.Join(buildDir, "ffmpeg")

    if !exists(zipPath) {
        // Change from 7.1 to 8.0
        download("https://codeload.github.com/FFmpeg/FFmpeg/zip/refs/heads/release/8.0", zipPath)
    }
    // ...existing code...
}
```

Expected library versions after upgrade:
```
FFmpeg 7.1.2 → FFmpeg 8.0:
libavutil      59.39.100 → 60. 8.100  (MAJOR BUMP)
libavcodec     61.19.100 → 62.11.100  (MAJOR BUMP)
libavformat    61. 7.100 → 62. 5.100  (MAJOR BUMP)
libavdevice    61. 3.100 → 62. 0.100  (MAJOR BUMP)
libavfilter    10. 4.100 → 11. 0.100  (MAJOR BUMP)
libswscale      8. 3.100 →  9. 0.100  (MAJOR BUMP)
libswresample   5. 3.100 →  6. 0.100  (MAJOR BUMP)
libpostproc    58. 3.100 → 59. 0.100  (MAJOR BUMP)
```

**Key Version Jump:** All libraries bump major version. This is a **breaking release**.

## 2. Critical API Removals & Breaking Changes

FFmpeg 8.0 removes all APIs deprecated before version 60 (3-major-release policy). **This is more aggressive than 7.0 which removed APIs before 60.**

### 2.1 Deprecation Policy Change (CRITICAL)

**New in FFmpeg 8.0:**
- Deprecated APIs now removed after **3 major releases** instead of indefinite deprecation
- APIs deprecated in FFmpeg 5.x are now **REMOVED** in 8.0
- Future deprecations in 8.x will be removed in 11.0

**Impact:** More frequent API churn. Must stay current with deprecation warnings.

### 2.2 Channel Layout API (Already Fixed in 7.1)

✅ **No action needed** - FFmpeg 7.1 upgrade already addressed:
- Removed `AVFrame.channel_layout`, `AVFrame.channels` → Using `AVFrame.ch_layout`
- Generator fix (commit 34609cd) handles unnamed structs
- jivedrop already uses modern `AVChannelLayout` API

### 2.3 Removed Deprecated AVFrame Fields

**Already removed in 7.1.2 (our current version):**
- `AVFrame.pkt_pts` → Using `AVFrame.pts`
- `AVFrame.top_field_first` → Using `AVFrame.flags`
- `AVFrame.interlaced_frame` → Using `AV_FRAME_FLAG_INTERLACED`

✅ **No action needed** - These were addressed in 7.1 upgrade.

### 2.4 New Deprecations in 8.0 (Monitor for Future)

**Newly deprecated in 8.0 (will be removed in 11.0):**
- Check FFmpeg's `APIchanges` and deprecation warnings after upgrade
- Document any new deprecation warnings for future 9.0/10.0/11.0 planning

**Impact:** Start tracking new deprecations now. They'll be removed in 3 major releases.

### 2.5 Codec/Format Changes

**Potential removals:**
- Some legacy codec implementations may be removed
- Check for encoder/decoder availability after upgrade
- Verify LAME MP3, AAC, H.264 encoders still work (likely fine)

**Impact:** Low - Core codecs are stable. More concerned with exotic/legacy formats.

## 3. Major New Features in FFmpeg 8.0

Understanding new features helps prioritize testing and documentation.

### 3.1 New Native Decoders

**APV (Acorn Replay Video)**
- Retro computing video format
- Low priority for modern use

**ProRes RAW**
- Professional video production codec
- **HIGH VALUE** - Apple ecosystem workflows
- Complements existing ProRes support

**RealVideo 6.0**
- Legacy streaming video codec
- Low priority unless supporting old content

**Sanyo LD-ADPCM & G.728**
- Audio codecs for specific hardware/telephony
- Low priority

### 3.2 Hardware Acceleration (MAJOR)

**Vulkan Compute Implementations:**
- **Vulkan VP9 decoder** - Cross-platform GPU decoding
- **Vulkan AV1 encoder** - Modern codec, GPU accelerated
- **Vulkan ProRes RAW** - Hardware accelerated RAW decoding

**VA-API Additions:**
- **VVC (H.266) VA-API decoder** - Next-gen codec support
- Requires modern Intel/AMD GPUs

**OpenHarmony Support:**
- H.264/H.265 decoding on HarmonyOS devices
- Low priority unless targeting Huawei ecosystem

**Impact:**
- Consider enabling Vulkan support in build (cross-platform HW accel)
- VA-API useful for Linux server encoding workloads
- ProRes RAW valuable for video production workflows

### 3.3 VVC (H.266) Improvements

**Enhancements:**
- IBC (Intra Block Copy) support
- ACT (Adaptive Color Transform) support
- Decoder stability improvements

**Impact:** VVC is the successor to HEVC. Early adoption opportunity for next-gen codec.

### 3.4 Encoder Enhancements

**libx265 Alpha Layer Encoding:**
- HEVC with transparency support
- Useful for compositing workflows

**CENC AV1 Encryption:**
- Common Encryption for AV1
- DRM/content protection workflows

**Impact:** Niche features but show ecosystem maturity.

### 3.5 New Filters

**Whisper Filter:**
- OpenAI Whisper speech-to-text integration
- **HIGH VALUE** - Automatic transcription in FFmpeg pipeline
- Could be useful for podcast workflows (jivedrop/jivetalking)

**Impact:** Consider testing Whisper filter for podcast metadata generation.
## 4. Update Generated Bindings

After updating FFmpeg headers, regenerate bindings using the **fixed generator** (v0.6.1):

```bash
# Use justfile for generation
just generate

# Or run generator directly
go run internal/generator/*.go
```

This regenerates:
- `constants.gen.go` - Will add new FFmpeg 8.0 constants
- `enums.gen.go` - May have new enum values
- `functions.gen.go` - Updated function signatures, removed deprecated functions
- `structs.gen.go` - Updated struct definitions (API changes from removed pre-6.0 APIs)

**Expected changes from 7.1 → 8.0:**
- New codec IDs for native decoders (APV, ProRes RAW, RealVideo 6.0, etc.)
- VVC codec constants updated for IBC/ACT support
- Hardware acceleration constants for Vulkan/VA-API/OpenHarmony
- New filter constants (Whisper filter)
- Potential struct changes if any pre-6.0 deprecated fields were retained
- Note: The 7.1 fixes (unnamed structs, channel layout) carry forward ✓

## 5. Update Version Constants

**In `ffmpeg_test.go`:**

```go
func TestVersions(t *testing.T) {
    // Update expected version constants for FFmpeg 8.0
    // All libraries bump major versions for FFmpeg 8.0
    // Exact values from FFmpeg 8.0 "Huffman":
    // - libavutil:    62.11.100
    // - libavcodec:   62. 8.100
    // - libavformat:  62. 5.100
    // - libavdevice:  62. 0.100
    // - libavfilter:  11. 0.100
    // - libswscale:    9. 0.100
    // - libswresample: 6. 0.100
    // - libpostproc:   7. 0.100

    assert.Equal(t, expectedCodecVersion, int(ffmpeg.AVCodecVersion()),
        "AVCodec version should match FFmpeg 8.0")
    assert.Equal(t, ffmpeg.LIBAVCodecVersionInt, int(ffmpeg.AVCodecVersion()),
        "AVCodec version func and const should match")
}
```

## 6. Update Documentation

**In `README.md`:**

Update library versions table:

```md
### Library versions

| Library       | Version   |
|---------------|-----------|
| FFmpeg        | 8.0       |
| aom           | 3.8.1     |
| opus          | 1.4       |
| vpx           | 1.14.0    |
| x264          | (latest)  |
| freetype      | 2.13.2    |
```

Update any documentation referring to FFmpeg version support or features.

## 7. Verify jivedrop Compatibility

Your jivedrop encoder at `/home/martin/Development/linuxmatters/jivedrop/internal/encoder/encoder.go` **should work without changes** because:

✅ **Already using modern APIs from FFmpeg 7.1:**
- `AVChannelLayoutDescribe()` - Modern channel layout API (added in 7.0)
- `AVChannelLayoutCopy()` - Modern channel layout copying
- `av_packet_rescale_ts()` - Standard time base handling
- Modern AVCodecContext usage
- All deprecated channel layout APIs already removed in 7.1 migration

🔍 **Things to verify after upgrade:**
1. MP3 encoder options still work (LAME encoder)
2. AAC encoder options still work (native AAC)
3. Metadata/ID3 tag writing still works
4. Frame/packet time base handling unchanged

**Note:** FFmpeg 8.0 primarily removes APIs deprecated *before 6.0*. Since jivedrop was updated for 7.1 (which already removed those APIs), you're in good shape.

## 8. Testing Strategy

**Phase 1: Basic Validation**
```bash
# 1. Rebuild static libraries for FFmpeg 8.0
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
# Test jivedrop encoder with new FFmpeg 8.0
cd ~/Development/linuxmatters/jivedrop
just build
just test

# Manual test: encode a real podcast episode
./jivedrop --input test.flac --output test.mp3 --episode "LMP-XX" --title "Test"
```

**Phase 4: New Features Exploration (Optional)**
```bash
# Test ProRes RAW if you have sample files
# Test Vulkan acceleration if available
# Test Whisper filter integration
```
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

## 9. Potential Issues & Mitigation

### Issue 1: Generator Parsing Errors
**Symptom:** Generator fails to parse new headers or headers from new features
**Fix:** Update go-clang/bootstrap if needed, check for new libclang compatibility issues
**Likelihood:** Low (generator fixes in v0.6.1 handle unnamed structs and platform differences)

### Issue 2: Struct Field Access Breaks
**Symptom:** Compilation errors accessing removed fields from pre-6.0 deprecated APIs
**Fix:** Should not occur - FFmpeg 7.1 already removed channel_layout/channels/pkt_pts fields
**Likelihood:** Very low (your code already uses modern APIs, 7.1 migration complete)

### Issue 3: Encoder Option Changes
**Symptom:** Encoder initialization fails or produces different output
**Fix:** Check encoder documentation, update AVOptions as needed
**Likelihood:** Low (LAME/AAC encoders are stable, no major changes in 8.0)

### Issue 4: Example Code Breaks
**Symptom:** Examples fail to compile or run
**Fix:** Update examples to use new APIs (if they use deprecated APIs for demonstration)
**Likelihood:** Low (examples likely updated during 7.1 migration)

### Issue 5: Hardware Acceleration Compatibility
**Symptom:** Vulkan/VA-API/VideoToolbox failures
**Fix:** Check hardware support, update driver/runtime versions if needed
**Likelihood:** Medium (new hardware features may require newer drivers)
**Note:** This is environmental, not code-related

### Issue 6: New Deprecations in 8.0
**Symptom:** Compiler warnings about deprecated APIs (newly deprecated in 8.0)
**Fix:** Track for future 9.0/10.0 migration, not urgent (3-major-release policy)
**Likelihood:** High (8.0 introduces new deprecations)
**Impact:** Low (warnings only, functionality works)

## 10. Recommended Upgrade Path

**Timeline: ~2-4 hours for careful upgrade**

1. **Preparation (30 min)**
   - Ensure current code is committed (clean working directory)
   - Create `ffmpeg-8.0` feature branch
   - Document current test results as baseline
   - Note: Already on clean main branch with 9MB git history ✓

2. **Build FFmpeg 8.0 (1 hour)**
   - Update builder to download FFmpeg 8.0 (release/8.0 branch)
   - Rebuild static libraries for your platform(s):
     - linux-amd64
     - linux-arm64
     - darwin-amd64
     - darwin-arm64
   - Verify libraries are correct version (check version constants)

3. **Regenerate Bindings (30 min)**
   - Run generator on new FFmpeg 8.0 headers
   - Review generated file diffs (expect new constants, possible struct changes)
   - Commit generated files separately for clean diff

4. **Fix & Test (1-2 hours)**
   - Update version constants in tests (62.x major versions)
   - Fix any compilation errors (unlikely)
   - Run all tests and examples
   - Test jivedrop integration
   - Check for new deprecation warnings (track for future)

5. **Static Library Integration (15 min)**
   - Add all 4 platform static libraries in single commit
   - Clean commit message documenting FFmpeg 8.0
   - Verify build and tests pass with libraries

6. **Fork to LinuxMatters (30 min)**
   - Create linuxmatters/static-ffmpeg-go repository
   - Push clean main branch
   - Update README with new project name
   - Update module path if needed
   - Tag as v0.8.0 or v1.0.0 (major version for FFmpeg major version)
   - Update GitHub release notes

## 11. Success Criteria

✅ All tests pass with FFmpeg 8.0
✅ All examples compile and run correctly
✅ jivedrop builds and encodes correctly
✅ Static libraries are correct FFmpeg 8.0 version
✅ No compilation errors (minor deprecation warnings acceptable)
✅ Generated bindings are clean (no parser errors)
✅ Version constants match FFmpeg 8.0 (62.x major versions)
✅ Clean git history maintained (single commit for libraries)
✅ Successfully forked to linuxmatters/static-ffmpeg-go
✅ Documentation updated (README, version table)

## 12. Post-Upgrade: Consider Additional Codecs for FFmpeg 8.0

Based on FFmpeg 8.0 features and current multimedia landscape, here are codecs/features to consider:

## High-Priority Additions

**AV1 Encoding (SVT-AV1 or rav1e)**
- Modern, royalty-free codec with excellent compression
- Growing adoption (YouTube, Netflix, streaming platforms)
- Superior to VP9, competitive with HEVC without licensing issues
- **Already have libaom decoder**, adding SVT-AV1 encoder would complete the story
- **FFmpeg 8.0 adds Vulkan AV1 encoder** - could leverage GPU acceleration

**HEVC/H.265 (x265)**
- Industry standard for 4K/HDR content
- Better compression than H.264
- Widely used in streaming, broadcasting, mobile
- There's even a commented-out reference in your FFmpeg configure: `--enable-libx265`
- **FFmpeg 8.0 adds x265 alpha layer encoding** - transparency support

**VVC/H.266 Support**
- **FFmpeg 8.0 adds VA-API VVC decoder** - next-gen codec
- Superior compression to HEVC
- Early adoption opportunity
- Consider native VVC support (vvdec/vvenc)

**Vulkan Hardware Acceleration**
- **FFmpeg 8.0 adds Vulkan VP9, AV1, ProRes RAW support**
- Cross-platform GPU acceleration (Linux/Windows/macOS)
- Significantly faster than CPU-only
- Consider `--enable-vulkan` in build

## Medium-Priority Additions

**WebP/WebM image support**
- `libwebp` - web-optimised image format
- Useful for thumbnail generation, web workflows
- Already have VP9 (part of WebM ecosystem)

**AVIF Support**
- `libavif` - AVIF image support (AV1-based)
- Better compression than WebP/JPEG
- Growing browser support
- Aligns with AV1 ecosystem

**ProRes RAW Decoding**
- **FFmpeg 8.0 adds native ProRes RAW decoder** - already available!
- Professional video production workflows
- Apple ecosystem integration
- No additional codec needed, but verify it's enabled

**Dolby Vision/HDR10+ metadata handling**
- Increasingly important for video processing
- Requires no additional codecs, just metadata support
- Enable with: `--enable-libdovi` (if available)
- **FFmpeg 8.0 improves HDR metadata handling**

**Hardware acceleration**
- **VA-API** (Linux Intel/AMD) - already common in CI
  - **FFmpeg 8.0 adds VVC VA-API decoder**
- **VideoToolbox** (macOS) - you have this enabled already
- **QSV** (Intel Quick Sync) - cross-platform
- **OpenHarmony** - **FFmpeg 8.0 adds H.264/H.265 support**
- Significantly faster encoding/decoding

**FDK-AAC**
- Higher quality AAC encoder than FFmpeg's native
- Important for audio workflows requiring AAC
- License-compatible with GPLv3
- Useful for podcast encoding (jivedrop)

## Lower-Priority But Useful

**Whisper Filter**
- **FFmpeg 8.0 adds OpenAI Whisper integration**
- Automatic speech-to-text transcription
- **HIGH VALUE for podcast workflows** (jivedrop/jivetalking)
- Consider enabling for metadata generation

**Faster AV1 Decoding**
- `libdav1d` - faster AV1 decoder (alternative to libaom)
- Better performance for AV1 playback
- Already widely adopted

**Codecs for specific workflows**
- **libkvazaar** - open HEVC encoder
- **RealVideo 6.0** - **FFmpeg 8.0 adds decoder** (legacy support)
- **APV** - **FFmpeg 8.0 adds decoder** (retro computing)

**Subtitle/accessibility**
- Already have `libass` ✓
- Consider `libfribidi`, `harfbuzz` for complex text (you have these ✓)

## Probably Skip

**Patent/License concerns:**
- AC-3/E-AC-3 (Dolby) - licensing restrictions
- DTS - licensing issues
- ProRes encoding - Apple proprietary (decoder is fine, encoder problematic)

**Niche/declining:**
- WMV/WMA (unless targeting legacy Windows content)
- Indeo, Cinepak, etc.
- **Sanyo LD-ADPCM, G.728** - **FFmpeg 8.0 adds these** but very niche

## Recommended Additions for FFmpeg 8.0 Build

Given you're building a **Go FFmpeg binding for general multimedia use** and FFmpeg 8.0's new features:

1. **Vulkan support** - Enable for cross-platform GPU acceleration (VP9, AV1, ProRes RAW)
2. **x265** - Essential for modern 4K/HEVC video
3. **SVT-AV1** - Future-proofing for royalty-free encoding, complements Vulkan AV1
4. **libfdk-aac** - Professional audio quality (podcast workflows)
5. **libwebp/libavif** - Web workflows and modern image formats
6. **Whisper filter** - Podcast transcription (jivedrop integration potential)
7. **VA-API** (Linux) - Hardware acceleration including new VVC decoder

These would cover:
- Modern streaming (AV1, HEVC, VVC)
- Hardware acceleration (Vulkan, VA-API)
- Web optimisation (WebP, AVIF, AV1)
- Professional audio (FDK-AAC)
- AI features (Whisper transcription)
- All royalty-free or GPL-compatible

Your current codec selection is already solid for **open source/royalty-free workflows**. The main gaps for FFmpeg 8.0 are:
- **HEVC/x265** for 4K content
- **SVT-AV1** for next-gen encoding
- **Vulkan** for cross-platform GPU acceleration (NEW in 8.0!)
- **Whisper** for transcription workflows (NEW in 8.0!)
