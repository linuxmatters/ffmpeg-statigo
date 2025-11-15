## Codec Dependency Analysis for FFmpeg Configure Flags

### Overview
The introspect tool generates consolidated FFmpeg configure flags for selective codec enablement/disablement with all dependencies. It uses intelligent matching and reverse lookups to discover all related codecs and components.

### Command Line Interface

**Enable Mode (Additive Pattern):**
```bash
./introspect --enable h264
./introspect --enable opus
./introspect --enable av1
```

**Disable Mode (Reductive Pattern):**
```bash
./introspect --disable vp7
./introspect --disable rv10
```

### Intelligent Codec Matching

The tool uses multiple strategies to find relevant codecs:

1. **Exact Match:** `av1` matches codec named exactly `av1`
2. **Variant Pattern:** `av1_qsv`, `av1_nvenc`, `av1_vulkan` (underscore-separated variants)
3. **Library Pattern:** `libdav1d`, `librav1e`, `libx264` (lib* prefix with search term)
4. **Prefix Match:** Searching `h26` finds `h264`, `h265`

**Excludes false positives:**
- Searching `av1` no longer matches `wmav1` (Windows Media Audio)
- Uses word boundary detection to avoid substring pollution

### Reverse Lookup Discovery

The tool performs comprehensive reverse lookups to find all related codecs:

1. **Parser Lookup:** Finds codecs supported by matching parsers
   - Search "av1" → finds `av1` parser → discovers all AV1 codec variants

2. **Format Lookup:** Finds codecs used by matching demuxers/muxers
   - Search "av1" → finds `avif`, `obu` formats → discovers AV1 codecs

3. **BSF Lookup:** Finds codecs used by matching bitstream filters
   - Search "av1" → finds `av1_metadata`, `av1_frame_split` → discovers AV1 codecs

This ensures comprehensive codec discovery even if you don't know all variant names.

### Algorithm

1. **Codec Discovery:**
   - Direct name matching using improved patterns (exact, variant, library, prefix)
   - Reverse lookup from parsers with matching names
   - Reverse lookup from formats with matching names
   - Reverse lookup from BSFs with matching names
   - Deduplicate all discovered codecs

2. **Smart Ordering:**
   - **Exact match first:** `av1` before `av1_qsv`
   - **Software codecs:** `libdav1d`, `librav1e`
   - **Hardware codecs:** `av1_qsv`, `av1_nvenc`, `av1_vulkan`

3. **Dependency Resolution:**

   For each discovered codec (by codec ID):

   - **Encoders:** Collect all encoders for this codec ID
   - **Decoders:** Collect all decoders for this codec ID
   - **Parsers:** Include parsers where codec ID appears in `parser.CodecIds()` array
   - **Demuxers:** Include formats where codec ID matches `format.VideoCodec`, `format.AudioCodec`, or `format.SubtitleCodec`
   - **Muxers:** Include formats where codec ID matches `format.VideoCodec`, `format.AudioCodec`, or `format.SubtitleCodec`
   - **Bitstream Filters:**
     - Include BSFs where codec ID appears in `bsf.CodecIds()` array
     - Include ALL BSFs where `bsf.CodecIds() == nil` (generic filters applicable to all codecs)

4. **Consolidation:**
   - Merge all codec dependencies into single output block
   - Deduplicate all component lists
   - Sort all items alphanumerically

### Output Format

Single consolidated block containing:

1. **Codec descriptions:** All matching codec descriptions as Go-style comments
2. **Encoders:** `--enable-encoder=` with comma-delimited, alphanumerically sorted list
3. **Decoders:** `--enable-decoder=` with comma-delimited, alphanumerically sorted list
4. **Parsers:** `--enable-parser=` with comma-delimited, alphanumerically sorted list
5. **Demuxers:** `--enable-demuxer=` with comma-delimited, alphanumerically sorted list
6. **Muxers:** `--enable-muxer=` with comma-delimited, alphanumerically sorted list
7. **BSFs:** `--enable-bsf=` with comma-delimited, alphanumerically sorted list

**Component types only appear if non-empty** (e.g., no decoder line if codec has no decoders)

**FFmpeg Configure Flag Names:**
- `--enable-encoder=` / `--disable-encoder=`
- `--enable-decoder=` / `--disable-decoder=`
- `--enable-parser=` / `--disable-parser=`
- `--enable-demuxer=` / `--disable-demuxer=`
- `--enable-muxer=` / `--disable-muxer=`
- `--enable-bsf=` / `--disable-bsf=`

### Example Output

**Enable Mode (AV1 - Multiple Variants):**
```bash
$ ./introspect --enable av1
```

```go
// Alliance for Open Media AV1
// librav1e AV1
// dav1d AV1 decoder by VideoLAN
// AV1 (Intel Quick Sync Video acceleration)
// AV1 (Vulkan)
// NVIDIA NVENC av1 encoder
--enable-encoder=av1_nvenc,av1_qsv,av1_vulkan,librav1e
--enable-decoder=av1,av1_qsv,libdav1d
--enable-parser=av1
--enable-demuxer=avif,obu
--enable-muxer=avif,obu
--enable-bsf=av1_frame_merge,av1_frame_split,av1_metadata,chomp,dovi_rpu,dump_extra,extract_extradata,filter_units,mov2textsub,noise,null,remove_extra,setts,showinfo,text2movsub,trace_headers
```

**Enable Mode (H.264 - Software + Hardware):**
```bash
$ ./introspect --enable h264
```

```go
// H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10
// libx264 H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10
// libx264 H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10 RGB
// H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10 (Intel Quick Sync Video acceleration)
// NVIDIA NVENC H.264 encoder
// H.264/AVC (Vulkan)
--enable-encoder=h264_nvenc,h264_qsv,h264_vulkan,libx264,libx264rgb
--enable-decoder=h264,h264_qsv
--enable-parser=h264
--enable-demuxer=dash,f4v,h264,hds,hls,ipod,ismv,matroska,mov,mp4,psp,smoothstreaming
--enable-muxer=dash,f4v,h264,hds,hls,ipod,ismv,matroska,mov,mp4,psp,smoothstreaming
--enable-bsf=chomp,dts2pts,dump_extra,extract_extradata,filter_units,h264_metadata,h264_mp4toannexb,h264_redundant_pps,mov2textsub,noise,null,remove_extra,setts,showinfo,text2movsub,trace_headers
```

**Enable Mode (AAC - Audio Codec):**
```bash
$ ./introspect --enable aac
```

```go
// AAC (Advanced Audio Coding)
// AAC LATM (Advanced Audio Coding LATM syntax)
--enable-encoder=aac
--enable-decoder=aac,aac_fixed,aac_latm
--enable-parser=aac,aac_latm
--enable-demuxer=adts,dash,f4v,hds,hls,ipod,ismv,latm,mov,mp4,psp,rtp_mpegts,rtsp,sap,smoothstreaming
--enable-muxer=adts,dash,f4v,hds,hls,ipod,ismv,latm,mov,mp4,psp,rtp_mpegts,rtsp,sap,smoothstreaming
--enable-bsf=aac_adtstoasc,chomp,dump_extra,mov2textsub,noise,null,remove_extra,setts,showinfo,text2movsub
```

**Disable Mode:**
```bash
$ ./introspect --disable vp7
```

```go
// On2 VP7
--disable-decoder=vp7
--disable-demuxer=avi,ivf,matroska
--disable-muxer=avi,ivf,matroska
```

### Edge Cases

1. **No matches found:** Print error message: `Error: No codec found matching '<codec>'`
2. **Generic BSFs:** Always included for all codecs (e.g., `chomp`, `dump_extra`, `null`, `trace_headers`)
3. **Codec with no encoder:** Only shows `--enable-decoder` line
4. **Codec with no decoder:** Only shows `--enable-encoder` line
5. **Duplicate descriptions:** Automatically deduplicated in comment section
6. **Hardware-only codecs:** Properly identified and included (e.g., `av1_nvenc`, `h264_qsv`)

### Key Improvements

**Version 2.0 Changes:**

1. **Consolidated Output:** All codec variants merged into single output block instead of multiple separate blocks
2. **Improved Matching:** Word boundary detection prevents false positives like `wmav1` matching `av1` search
3. **Reverse Lookups:** Discovers codecs through parser/format/BSF name matching
4. **Smart Ordering:** Exact match → software → hardware variants
5. **Deduplication:** Descriptions, components, and flags automatically deduplicated
6. **Alphanumeric Sorting:** All component lists sorted for consistency

### Implementation Notes

- Uses existing introspection functions: `listCodecs()`, `listParsers()`, `listFormats()`, `listBSFs()`
- New function: `analyzeCodecDependencies(codecName string, disable bool)`
- Helper functions:
  - `findMatchingCodecs()` - Intelligent matching with reverse lookups
  - `sortCodecsByPriority()` - Orders exact → software → hardware
  - `consolidateCodecDependencies()` - Merges all dependencies with deduplication
  - `outputConsolidatedDependencies()` - Formats clean output
- Command-line flag parsing intercepts `--enable`/`--disable` before normal introspection
