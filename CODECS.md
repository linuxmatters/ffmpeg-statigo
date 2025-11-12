# Codec Inclusion Policy

## Mission Statement

**ffmpeg-go provides a curated FFmpeg static library with Go bindings for contemporary audio and video processing, encoding, conversion and streaming.**

This is a purpose-built media toolkit for Go developers building modern applications—optimised for size, focused on current workflows, and opinionated about what matters in 2025 and beyond.

---

## Target Audience

**Go developers building:**
- Streaming platforms (live and on-demand)
- Content management systems
- Social media applications
- Video conferencing services
- Media transcoding pipelines
- Web-based media players
- Broadcasting tools
- Modern content creation workflows

**NOT for:**
- Video archivists preserving historical formats
- Retro gaming emulator developers
- Digital archaeology projects
- Legacy format conversion utilities
- Film restoration specialists

---

## Codec Inclusion Criteria

A codec is **included** if it meets **ALL** of these:

### 1. Contemporary Relevance
At least **ONE** must be true:
- Standardised by major industry body (MPEG, ITU, W3C, IETF) in last 15 years
- Actively used in modern streaming platforms (YouTube, Netflix, Twitch, etc.)
- Required for current browser/device compatibility (iOS, Android, modern browsers)
- Standard in modern cameras/production equipment (< 5 years old)
- Essential for broadcast/professional workflows

### 2. Ecosystem Health
- Actively maintained (commits/releases in last 2 years)
- Broad adoption across multiple platforms/vendors
- Clear future path (not abandoned technology)

### 3. Go Developer Use Cases
- Solves common problems in web/cloud applications
- Enables modern streaming workflows
- Supports cross-platform delivery requirements

---

## Codec Exclusion Criteria

A codec is **excluded** if it meets **ANY** of these:

### 1. Historical Only
- Primary use case is archival/preservation of old content
- Associated with discontinued platforms (Flash, Silverlight, Windows Media, RealMedia)
- Last significant adoption was before 2010

### 2. Niche/Specialised
- Game engine proprietary formats (Bink, Smacker, RoQ, etc.)
- Console-specific codecs (PlayStation ADPCM, etc.)
- Retro gaming archival (Doom/Quake formats, C&C, etc.)
- Industrial/vertical-specific (surveillance DVR, medical imaging)
- Regional standards with limited global adoption

### 3. Superseded Technology
- Better alternatives exist and are widely adopted
- Format was transitional/bridge technology
- Industry has moved on (DivX, Xvid, VC-1, etc.)

### 4. Maintenance Burden
- Adds significant binary size for limited contemporary use
- No active maintenance or security updates
- Creates compatibility issues

---

## Decision Framework

When evaluating a codec:

```
1. Is it required for modern web/streaming workflows?
   YES → Include
   NO  → Continue to 2

2. Is it a current industry standard (< 15 years old)?
   YES → Continue to 3
   NO  → Exclude

3. Is it actively maintained by its community/vendor?
   YES → Continue to 4
   NO  → Exclude

4. Would a typical Go web/streaming app need it?
   YES → Include
   NO  → Continue to 5

5. Is the primary use case archival/legacy support?
   YES → Exclude
   NO  → Include (edge case - document rationale)
```

---

## Included Codecs

### Video Decoders/Encoders

**Modern Compression (2010s+):**
- **H.264/AVC** (x264) - Universal standard, web/mobile/broadcast
- **H.265/HEVC** (x265) - 4K/8K, modern streaming, broadcast
- **VP8** - WebRTC, open-source, YouTube legacy
- **VP9** - YouTube, open-source, browser support
- **AV1** (libaom) - Next-gen compression, future-proof, Netflix/YouTube

**Production/Intermediate:**
- **ProRes** - Professional editing, broadcast contribution
- **DNxHD/DNxHR** - Professional editing, Avid workflows
- **MJPEG** - IP cameras, simple streaming, editing proxies

**Compatibility/Legacy (Still Active):**
- **H.263** - Basic WebRTC fallback, legacy mobile (considering removal)
- **MPEG-2** - DVB broadcast, some cameras (widespread device support)
- **MPEG-1** - VCD, retro but universal playback

### Audio Decoders/Encoders

**Modern Audio:**
- **Opus** - VoIP, WebRTC, low-latency streaming
- **AAC** - Universal mobile/web audio, podcasts
- **MP3** (LAME) - Universal music delivery, legacy but essential

**Open Standards:**
- **Vorbis** - Ogg container, gaming, open-source stacks
- **FLAC** - Lossless archival, music streaming
- **Speex** - VoIP (legacy Opus precursor, considering removal)

**Professional:**
- **PCM** (various) - Uncompressed audio, production
- **ALAC** - Apple Lossless, iTunes ecosystem

### Containers/Formats
- **MP4/MOV** - Universal delivery
- **WebM** - Web standard (VP8/VP9 + Opus/Vorbis)
- **Matroska (MKV)** - Flexible container
- **HLS/DASH** - Adaptive streaming
- **RTMP/SRT** - Live streaming protocols
- **Ogg** - Open-source, Vorbis/Theora

---

## Excluded Codecs

### Discontinued Platforms

**Windows Media (Microsoft discontinued):**
- WMV1, WMV2, WMV3 - Windows Media Video 7/8/9
- VC-1 - Windows Media Video Advanced Profile (Blu-ray archival only)
- WMA Pro, WMA Lossless, WMA Voice - Windows Media Audio
- MS-MPEG4v1/v2/v3 - Pre-WMV era

**RealMedia (Platform ended 2018):**
- RealVideo 1.0-6.0 (rv10, rv20, rv30, rv40, rv60)
- RealAudio 14.4, 28.8 (ra_144, ra_288)
- RealAudio Cook, Sipro (cook, sipr)

**Flash (Adobe EOL 2020):**
- VP3, VP4, VP5, VP6, VP6A, VP6F - Early On2 codecs
- VP7 - Flash-era bridge codec (superseded by VP8)
- Sorenson Spark (FLV1) - Flash H.263 variant

**Apple Legacy (QuickTime Era):**
- Sorenson Video 1/3 (SVQ1/SVQ3)
- Apple Video (RPZA, SMC)
- Microsoft Video 1 (msvideo1)
- QDesign Music 1/2 (QDM, QDMC)
- MACE 3:1/6:1 (Macintosh Audio Compression)

**Sony Proprietary:**
- ATRAC1/3/3AL/3P/3PAL/9 - MiniDisc, PlayStation platforms

### Retro Gaming / Game Engines

**Id Software (Doom/Quake Era):**
- Id CIN video/audio (idcin, dsicinaudio)
- RoQ video/audio (roq, roq_dpcm) - Quake III modding

**Game Companies:**
- Bink Video/Audio (bink, binkaudio) - RAD Game Tools
- Smacker Video/Audio (smacker, smackaud) - RAD Game Tools
- Westwood VQA - Command & Conquer
- Interplay MVE (interplay_video, interplay_acm, interplay_dpcm)
- 4X Movie (4xm, adpcm_4xm)

**Electronic Arts Proprietary:**
- EA CMV, MAD, TGQ, TGV, TQI (video)
- EA ADPCM variants (Maxis XA, R1, R2, R3, XAS, EACS, SEAD)

**Other Legacy:**
- Bethesda VID (bethsoftvid) - Pre-2002 games
- Delphine Software CIN (dsicinvideo, dsicinaudio)
- Autodesk FLIC (flic) - 1990s animation
- Deluxe Paint Animation (anm) - 1980s

### Legacy Video Codecs

**Pre-2000s:**
- Cinepak - CD-ROM era
- Indeo 2/3/4/5 - Intel video
- TrueMotion 1/2/2RT - Duck Corporation

**Early 2000s (Superseded):**
- MPEG-4 Part 2 (mpeg4) - DivX/Xvid era, replaced by H.264
- H.261 - 1990s videoconferencing
- Chinese AVS (CAVS) - Regional standard, limited global adoption

### Legacy Audio Codecs

**Mobile/Cellular:**
- QCELP - Qualcomm PureVoice (1990s)
- EVRC - Enhanced Variable Rate Codec
- TrueSpeech - DSP Group (1990s)

**Console/Proprietary:**
- PlayStation 1 ADPCM (adpcm_psx, adpcm_xa)
- Ulead DV Audio (dvaudio)

---

## Rationale for Borderline Cases

### Included Despite Age

**MPEG-2:**
- Still used in DVB broadcast (Europe, Asia)
- Legacy camera support
- Universal hardware decode support

**MJPEG:**
- IP cameras still widely use it
- Simple editing proxies
- Low-latency streaming

**Theora:**
- Ogg ecosystem completeness
- Wikipedia uses it
- Open-source requirement for some projects

### Excluded Despite Some Modern Use

**VC-1:**
- Blu-ray discs use it (2000s-2010s)
- But: No contemporary **encoding** use case
- Transcoding Blu-ray rips is archival, not production
- H.264/HEVC replaced it completely

**MPEG-4 Part 2 (DivX/Xvid):**
- Widespread in 2000s video files
- But: H.264 superior and universal
- Contemporary workflows don't create Part 2
- Archival playback only

**Bink/Smacker:**
- Still used in AAA games (2020s)
- But: Game engine integration, not web/streaming
- Not a Go developer use case
- Games bundle their own decoders

---

## Library Size Impact

### Current Optimizations (as of 2025-11-12)

**Starting point:** 69,102,128 bytes (69.1 MB)

**Optimization passes:**
1. **Pass 1 - Ancient video codecs:** 989 KB saved
   - RealVideo, VP3-6, TrueMotion, Cinepak, Indeo, MS-MPEG4v1/v2, H.261/H.263
2. **Pass 2 - Game/multimedia codecs:** 151 KB saved
   - EA formats, Bethesda, FLIC, ANM, 4XM, Interplay, VQA, Delphine
3. **Pass 3 - Legacy audio codecs:** 100 KB saved
   - RealAudio, TrueSpeech, QCELP, game audio, EA audio, PS1 ADPCM

**Current size:** 67,861,966 bytes (67.8 MB)
**Total savings:** 1,240,162 bytes (1.24 MB, 1.79% reduction)

### Additional Planned Removals (Pass 4)

**Windows Media codecs:** ~500-800 KB estimated
**Apple Legacy codecs:** ~200-300 KB estimated
**Remaining game codecs:** ~300-400 KB estimated
**Legacy MPEG-4/VC-1:** ~400-500 KB estimated

**Projected final size:** ~65-66 MB (4-5 MB total savings, ~6-7% reduction)

---

## Future Considerations

### Potential Additions
- **x265 (HEVC encoder)** - Modern 4K/8K workflows
- **SVT-AV1** - Alternative AV1 encoder (faster than libaom)
- **rav1e** - Rust AV1 encoder (quality focus)
- **Native AAC encoder** - Replace libfdk-aac restriction

### Potential Removals
- **Speex** - Fully superseded by Opus, consider removal
- **H.263** - Obsolete, minimal WebRTC fallback use
- **Theora** - Declining use, VP9/AV1 preferred

### Monitoring
- Review codec decisions annually
- Track adoption of new standards (VVC, EVC, LCEVC)
- Monitor browser/device support changes
- Community feedback on edge cases

---

## Contributing

If you believe a codec should be included or excluded:

1. Open an issue with rationale
2. Provide evidence of contemporary use (or lack thereof)
3. Cite industry adoption or standards body backing
4. Describe the Go developer use case
5. Estimate binary size impact if known

We welcome discussion, but maintain our opinionated stance on contemporary workflows over archival completeness.

---

## License Note

All included codecs respect FFmpeg's GPL/LGPL licensing. We explicitly exclude:
- **libfdk-aac** - Redistribution licensing restrictions (despite superior quality)
- **Non-redistributable codecs** - Any format requiring paid licenses

---

**Last updated:** 2025-11-12
**FFmpeg version:** 8.0
**Go version:** 1.24+


Full Journey Summary
Pass	Focus	Size Before	Size After	Savings
Baseline	-	69,102,128	-	-
Pass 1	Ancient video	69,102,128	68,112,556	989 KB
Pass 2	Game/multimedia	68,112,556	67,961,766	151 KB
Pass 3	Legacy audio	67,961,766	67,861,966	100 KB
Pass 4	Contemporary framework ⭐	67,861,966	66,557,910	1,304 KB
Pass 5	Parsers/encoders/containers	66,557,910	65,343,814	1,214 KB
Pass 6	Telecom/legacy codecs	65,343,814	65,143,862	200 KB
TOTAL	6 optimization passes	69,102,128	65,143,862	3,958 KB (3.96 MB, 5.73%)
Size Progression


69.1 MB ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ 100%
        ↓ -989 KB  (Pass 1: Ancient video)
68.1 MB ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  98.6%
        ↓ -151 KB  (Pass 2: Game video)
68.0 MB ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━   98.4%
        ↓ -100 KB  (Pass 3: Legacy audio)
67.9 MB ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━    98.2%
        ↓ -1,304 KB (Pass 4: WMV/Flash/retro ⭐)
66.6 MB ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━     96.4%
        ↓ -1,214 KB (Pass 5: Parsers/encoders)
65.3 MB ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━      94.5%
        ↓ -200 KB  (Pass 6: Telecom/legacy)
65.1 MB ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━       94.3% ✅
