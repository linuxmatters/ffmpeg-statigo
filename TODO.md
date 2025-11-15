# TODO

- [x] Refactor the internal build tool
- [x] Reorganise project structure
- [x] FFmpeg build argument resolver.to enable/disable codecs
- [ ] Enable Dolby Vision/HDR10+ metadata handling via `--enable-libdovi` (if available)
- [ ] Add `libwebp` - web-optimised image format (part of WebM ecosystem) <https://chromium.googlesource.com/webm/libwebp>
- [ ] Add `libavif` - AVIF image support (AV1-based) <https://github.com/AOMediaCodec/libavif>
- [ ] Add `vvcenv` -  H.266/VVC encoder <https://github.com/fraunhoferhhi/vvenc>
- [ ] Rebrand to ffmpeg-statigo
- [ ] How to embed/distribute the static FFmpeg library?
- [ ] Enable Whisper for transcription workflows (New in FFmpgeg in 8.0 ✨)
- [ ] Review default codecs:
  - https://ffmpeg.martin-riedl.de/
  - https://github.com/markus-perl/ffmpeg-build-script/blob/master/build-ffmpeg

## More headers

### Adding Headers to the Generator

1. Modify the Generator Configuration

The header list is defined in a generator configuration file `generator/parser.go`. AI dd entries like:

```go
headers = append(headers,
    "libswscale/swscale.h",
    "libswresample/swresample.h",
    "libavformat/rtmp.h",
    "libavdevice/avdevice.h"
)
```

2. Regenerate Bindings

Run the generator with the updated clang (which you've already fixed for Linux compatibility):

```bash
just generate
```

### Critical Missing Headers for Streaming

#### Network Protocol Headers

- `libavformat/rtmp*.h` - RTMP/RTMPS protocol internals for custom handshaking and stream key management
- `libavformat/srt*.h` - SRT (Secure Reliable Transport) for low-latency streaming, increasingly required by platforms
- `libavformat/hls*.h` & libavformat/dash*.h - For generating adaptive bitrate streams with proper segment control

#### Real-time Processing Headers

- `libavfilter/buffersrc.h` & `libavfilter/buffersink.h` - Essential for building custom filter graphs for overlays, watermarks, and scene transitions
- `libavdevice/avdevice.h` - Capture device integration (webcams, capture cards, screen recording)
- `libavutil/fifo.h` - Thread-safe FIFO buffers for managing multiple output streams

#### Advanced Streaming Control

- `libavformat/avio_internal.h` - Custom I/O contexts for authentication and stream routing
- `libavcodec/bsf.h` - Bitstream filters for H.264/HEVC Annex B conversion (required by some platforms)
- `libavutil/threadmessage.h` - Inter-thread communication for multiple simultaneous outputs

### Missing Headers for Scaling and Resampling

- `libswresample/swresample.h` and `libswscale/swscale.h`
  - Creating adaptive bitrate ladders - scaling from a single 1080p input to 720p, 480p, 360p variants
  - Audio normalisation - resampling between 44.1kHz and 48kHz, converting 5.1 to stereo
  - Pixel format conversion - converting between YUV420P, NV12, and platform-specific requirements

## From the original author

- [ ] Expose more headers.
- [ ] Expose platform specific headers.
- [ ] Cleanup internal packages.

## Testing

- [ ] Create some test cases that exercise some of the FFmpeg API surface

### Basic Validation
# Rebuild static librarie for FFmpeg
```bash
just build-ffmpeg
```

**Regenerate bindings:**
```bash
just generate
```

**Run existing tests:**
```bash
just test
```

### Examples Validation
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

# CODECS

## Enabled

DE  NAME                     DESCRIPTION                                TYPE

DE  aac                      AAC (Advanced Audio Coding)                [AUDIO]
DE  ac3                      ATSC A/52A (AC-3)                          [AUDIO]
DE  alac                     ALAC (Apple Lossless Audio Codec)          [AUDIO]
DE  anull                    Null audio codec                           [AUDIO]
DE  aptx                     aptX (Audio Processing Technology for Blue [AUDIO]
DE  aptx_hd                  aptX HD (Audio Processing Technology for B [AUDIO]
DE  apng                     APNG (Animated Portable Network Graphics)  [VIDEO]
DE  ass                      ASS (Advanced SSA) subtitle                [SUBTITLE]
DE  av1                      Alliance for Open Media AV1                [VIDEO]
DE  bmp                      BMP (Windows and OS/2 bitmap)              [VIDEO]
DE  cfhd                     GoPro CineForm HD                          [VIDEO]
DE  dnxhd                    VC3/DNxHD                                  [VIDEO]
DE  dpx                      DPX (Digital Picture Exchange) image       [VIDEO]
DE  dts                      DCA (DTS Coherent Acoustics)               [AUDIO]
DE  dvb_subtitle             DVB subtitles                              [SUBTITLE]
DE  dvd_subtitle             DVD subtitles                              [SUBTITLE]
DE  eac3                     ATSC A/52B (AC-3, E-AC-3)                  [AUDIO]
DE  exr                      OpenEXR image                              [VIDEO]
DE  ffv1                     FFmpeg video codec #1                      [VIDEO]
DE  flac                     FLAC (Free Lossless Audio Codec)           [AUDIO]
DE  gif                      CompuServe GIF (Graphics Interchange Forma [VIDEO]
DE  h264                     H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10  [VIDEO]
DE  hevc                     H.265 / HEVC (High Efficiency Video Coding [VIDEO]
DE  jpeg2000                 JPEG 2000                                  [VIDEO]
DE  jpegls                   JPEG-LS                                    [VIDEO]
.E  ljpeg                    Lossless JPEG                              [VIDEO]
DE  mjpeg                    Motion JPEG                                [VIDEO]
DE  mp2                      MP2 (MPEG audio layer 2)                   [AUDIO]
DE  mp3                      MP3 (MPEG audio layer 3)                   [AUDIO]
DE  mpeg2video               MPEG-2 video                               [VIDEO]
DE  opus                     Opus (Opus Interactive Audio Codec)        [AUDIO]
DE  pbm                      PBM (Portable BitMap) image                [VIDEO]
DE  pcm_alaw                 PCM A-law / G.711 A-law                    [AUDIO]
DE  pcm_f32be                PCM 32-bit floating point big-endian       [AUDIO]
DE  pcm_f32le                PCM 32-bit floating point little-endian    [AUDIO]
DE  pcm_mulaw                PCM mu-law / G.711 mu-law                  [AUDIO]
DE  pcm_s16be                PCM signed 16-bit big-endian               [AUDIO]
DE  pcm_s16be_planar         PCM signed 16-bit big-endian planar        [AUDIO]
DE  pcm_s16le                PCM signed 16-bit little-endian            [AUDIO]
DE  pcm_s16le_planar         PCM signed 16-bit little-endian planar     [AUDIO]
DE  pcm_s24be                PCM signed 24-bit big-endian               [AUDIO]
DE  pcm_s24le                PCM signed 24-bit little-endian            [AUDIO]
DE  pcm_s24le_planar         PCM signed 24-bit little-endian planar     [AUDIO]
DE  pcm_s32be                PCM signed 32-bit big-endian               [AUDIO]
DE  pcm_s32le                PCM signed 32-bit little-endian            [AUDIO]
DE  pcm_s32le_planar         PCM signed 32-bit little-endian planar     [AUDIO]
DE  png                      PNG (Portable Network Graphics) image      [VIDEO]
DE  ppm                      PPM (Portable PixelMap) image              [VIDEO]
DE  prores                   Apple ProRes (iCodec Pro)                  [VIDEO]
D.  prores_raw               Apple ProRes RAW                           [VIDEO]
DE  qoi                      QOI (Quite OK Image)                       [VIDEO]
DE  sbc                      SBC (low-complexity subband codec)         [AUDIO]
DE  text                     raw UTF-8 text                             [SUBTITLE]
DE  theora                   Theora                                     [VIDEO]
DE  tiff                     TIFF image                                 [VIDEO]
DE  truehd                   TrueHD                                     [AUDIO]
D.  vc1                      SMPTE VC-1                                 [VIDEO]
DE  vnull                    Null video codec                           [VIDEO]
DE  vorbis                   Vorbis                                     [AUDIO]
DE  vp8                      On2 VP8                                    [VIDEO]
DE  vp9                      Google VP9                                 [VIDEO]
D.  vvc                      H.266 / VVC (Versatile Video Coding)       [VIDEO]
D.  webp                     WebP                                       [VIDEO]
DE  webvtt                   WebVTT subtitle                            [SUBTITLE]
DE  yuv4                     Uncompressed packed 4:2:0                  [VIDEO]

## Excluded

These are the codecs I've excluded

 DE  NAME                     DESCRIPTION                                TYPE

 D.  012v                     Uncompressed 4:2:2 10-bit                  [VIDEO]
 D.  4xm                      4X Movie                                   [VIDEO]
 D.  8bps                     QuickTime 8BPS video                       [VIDEO]
 D.  8svx_exp                 8SVX exponential                           [AUDIO]
 D.  8svx_fib                 8SVX fibonacci                             [AUDIO]
 .E  a64_multi                Multicolor charset for Commodore 64        [VIDEO]
 .E  a64_multi5               Multicolor charset for Commodore 64, exten [VIDEO]
 D.  aac_latm                 AAC LATM (Advanced Audio Coding LATM synta [AUDIO]
 D.  aasc                     Autodesk RLE                               [VIDEO]
 D.  acelp.kelvin             Sipro ACELP.KELVIN                         [AUDIO]
 D.  adpcm_4xm                ADPCM 4X Movie                             [AUDIO]
 DE  adpcm_adx                SEGA CRI ADX ADPCM                         [AUDIO]
 D.  adpcm_afc                ADPCM Nintendo Gamecube AFC                [AUDIO]
 D.  adpcm_agm                ADPCM AmuseGraphics Movie AGM              [AUDIO]
 D.  adpcm_aica               ADPCM Yamaha AICA                          [AUDIO]
 DE  adpcm_argo               ADPCM Argonaut Games                       [AUDIO]
 D.  adpcm_ct                 ADPCM Creative Technology                  [AUDIO]
 D.  adpcm_dtk                ADPCM Nintendo Gamecube DTK                [AUDIO]
 D.  adpcm_ea                 ADPCM Electronic Arts                      [AUDIO]
 D.  adpcm_ea_maxis_xa        ADPCM Electronic Arts Maxis CDROM XA       [AUDIO]
 D.  adpcm_ea_r1              ADPCM Electronic Arts R1                   [AUDIO]
 D.  adpcm_ea_r2              ADPCM Electronic Arts R2                   [AUDIO]
 D.  adpcm_ea_r3              ADPCM Electronic Arts R3                   [AUDIO]
 D.  adpcm_ea_xas             ADPCM Electronic Arts XAS                  [AUDIO]
 DE  adpcm_g722               G.722 ADPCM                                [AUDIO]
 DE  adpcm_g726               G.726 ADPCM                                [AUDIO]
 DE  adpcm_g726le             G.726 ADPCM little-endian                  [AUDIO]
 D.  adpcm_ima_acorn          ADPCM IMA Acorn Replay                     [AUDIO]
 DE  adpcm_ima_alp            ADPCM IMA High Voltage Software ALP        [AUDIO]
 DE  adpcm_ima_amv            ADPCM IMA AMV                              [AUDIO]
 D.  adpcm_ima_apc            ADPCM IMA CRYO APC                         [AUDIO]
 DE  adpcm_ima_apm            ADPCM IMA Ubisoft APM                      [AUDIO]
 D.  adpcm_ima_cunning        ADPCM IMA Cunning Developments             [AUDIO]
 D.  adpcm_ima_dat4           ADPCM IMA Eurocom DAT4                     [AUDIO]
 D.  adpcm_ima_dk3            ADPCM IMA Duck DK3                         [AUDIO]
 D.  adpcm_ima_dk4            ADPCM IMA Duck DK4                         [AUDIO]
 D.  adpcm_ima_ea_eacs        ADPCM IMA Electronic Arts EACS             [AUDIO]
 D.  adpcm_ima_ea_sead        ADPCM IMA Electronic Arts SEAD             [AUDIO]
 D.  adpcm_ima_iss            ADPCM IMA Funcom ISS                       [AUDIO]
 D.  adpcm_ima_moflex         ADPCM IMA MobiClip MOFLEX                  [AUDIO]
 D.  adpcm_ima_mtf            ADPCM IMA Capcom's MT Framework            [AUDIO]
 D.  adpcm_ima_oki            ADPCM IMA Dialogic OKI                     [AUDIO]
 DE  adpcm_ima_qt             ADPCM IMA QuickTime                        [AUDIO]
 D.  adpcm_ima_rad            ADPCM IMA Radical                          [AUDIO]
 D.  adpcm_ima_smjpeg         ADPCM IMA Loki SDL MJPEG                   [AUDIO]
 DE  adpcm_ima_ssi            ADPCM IMA Simon & Schuster Interactive     [AUDIO]
 DE  adpcm_ima_wav            ADPCM IMA WAV                              [AUDIO]
 DE  adpcm_ima_ws             ADPCM IMA Westwood                         [AUDIO]
 D.  adpcm_ima_xbox           ADPCM IMA Xbox                             [AUDIO]
 DE  adpcm_ms                 ADPCM Microsoft                            [AUDIO]
 D.  adpcm_mtaf               ADPCM MTAF                                 [AUDIO]
 D.  adpcm_psx                ADPCM Playstation                          [AUDIO]
 D.  adpcm_sanyo              ADPCM Sanyo                                [AUDIO]
 D.  adpcm_sbpro_2            ADPCM Sound Blaster Pro 2-bit              [AUDIO]
 D.  adpcm_sbpro_3            ADPCM Sound Blaster Pro 2.6-bit            [AUDIO]
 D.  adpcm_sbpro_4            ADPCM Sound Blaster Pro 4-bit              [AUDIO]
 DE  adpcm_swf                ADPCM Shockwave Flash                      [AUDIO]
 D.  adpcm_thp                ADPCM Nintendo THP                         [AUDIO]
 D.  adpcm_thp_le             ADPCM Nintendo THP (Little-Endian)         [AUDIO]
 D.  adpcm_vima               LucasArts VIMA audio                       [AUDIO]
 D.  adpcm_xa                 ADPCM CDROM XA                             [AUDIO]
 D.  adpcm_xmd                ADPCM Konami XMD                           [AUDIO]
 DE  adpcm_yamaha             ADPCM Yamaha                               [AUDIO]
 D.  adpcm_zork               ADPCM Zork                                 [AUDIO]
 D.  agm                      Amuse Graphics Movie                       [VIDEO]
 D.  aic                      Apple Intermediate Codec                   [VIDEO]
 DE  alias_pix                Alias/Wavefront PIX image                  [VIDEO]
 D.  amr_nb                   AMR-NB (Adaptive Multi-Rate NarrowBand)    [AUDIO]
 D.  amr_wb                   AMR-WB (Adaptive Multi-Rate WideBand)      [AUDIO]
 DE  amv                      AMV Video                                  [VIDEO]
 D.  anm                      Deluxe Paint Animation                     [VIDEO]
 D.  ansi                     ASCII/ANSI art                             [VIDEO]
 D.  apac                     Marian's A-pac audio                       [AUDIO]
 D.  ape                      Monkey's Audio                             [AUDIO]
 D.  apv                      Advanced Professional Video                [VIDEO]
 D.  arbc                     Gryphon's Anim Compressor                  [VIDEO]
 D.  argo                     Argonaut Games Video                       [VIDEO]
 DE  asv1                     ASUS V1                                    [VIDEO]
 DE  asv2                     ASUS V2                                    [VIDEO]
 D.  atrac1                   ATRAC1 (Adaptive TRansform Acoustic Coding [AUDIO]
 D.  atrac3                   ATRAC3 (Adaptive TRansform Acoustic Coding [AUDIO]
 D.  atrac3al                 ATRAC3 AL (Adaptive TRansform Acoustic Cod [AUDIO]
 D.  atrac3p                  ATRAC3+ (Adaptive TRansform Acoustic Codin [AUDIO]
 D.  atrac3pal                ATRAC3+ AL (Adaptive TRansform Acoustic Co [AUDIO]
 D.  atrac9                   ATRAC9 (Adaptive TRansform Acoustic Coding [AUDIO]
 D.  aura                     Auravision AURA                            [VIDEO]
 D.  aura2                    Auravision Aura 2                          [VIDEO]
 D.  avc                      On2 Audio for Video Codec                  [AUDIO]
 D.  avrn                     Avid AVI Codec                             [VIDEO]
 DE  avrp                     Avid 1:1 10-bit RGB Packer                 [VIDEO]
 D.  avs                      AVS (Audio Video Standard) video           [VIDEO]
 DE  avui                     Avid Meridien Uncompressed                 [VIDEO]
 D.  bethsoftvid              Bethesda VID video                         [VIDEO]
 D.  bfi                      Brute Force & Ignorance                    [VIDEO]
 D.  binkaudio_dct            Bink Audio (DCT)                           [AUDIO]
 D.  binkaudio_rdft           Bink Audio (RDFT)                          [AUDIO]
 D.  binkvideo                Bink video                                 [VIDEO]
 D.  bintext                  Binary text                                [VIDEO]
 DE  bitpacked                Bitpacked                                  [VIDEO]
 D.  bmv_audio                Discworld II BMV audio                     [AUDIO]
 D.  bmv_video                Discworld II BMV video                     [VIDEO]
 D.  bonk                     Bonk audio                                 [AUDIO]
 D.  brender_pix              BRender PIX image                          [VIDEO]
 D.  c93                      Interplay C93                              [VIDEO]
 D.  cavs                     Chinese AVS (Audio Video Standard) (AVS1-P [VIDEO]
 D.  cbd2_dpcm                DPCM Cuberoot-Delta-Exact                  [AUDIO]
 D.  cdgraphics               CD Graphics video                          [VIDEO]
 D.  cdtoons                  CDToons video                              [VIDEO]
 D.  cdxl                     Commodore CDXL video                       [VIDEO]
 DE  cinepak                  Cinepak                                    [VIDEO]
 D.  clearvideo               Iterated Systems ClearVideo                [VIDEO]
 DE  cljr                     Cirrus Logic AccuPak                       [VIDEO]
 D.  cllc                     Canopus Lossless Codec                     [VIDEO]
 D.  cmv                      Electronic Arts CMV video                  [VIDEO]
 DE  comfortnoise             RFC 3389 Comfort Noise                     [AUDIO]
 D.  cook                     Cook / Cooker / Gecko (RealAudio G2)       [AUDIO]
 D.  cpia                     CPiA video format                          [VIDEO]
 D.  cri                      Cintel RAW                                 [VIDEO]
 D.  cscd                     CamStudio                                  [VIDEO]
 D.  cyuv                     Creative YUV (CYUV)                        [VIDEO]
 D.  dds                      DirectDraw Surface image decoder           [VIDEO]
 D.  derf_dpcm                DPCM Xilam DERF                            [AUDIO]
 D.  dfa                      Chronomaster DFA                           [VIDEO]
 DE  dfpwm                    DFPWM (Dynamic Filter Pulse Width Modulati [AUDIO]
 DE  dirac                    Dirac                                      [VIDEO]
 D.  dolby_e                  Dolby E                                    [AUDIO]
 D.  dsd_lsbf                 DSD (Direct Stream Digital), least signifi [AUDIO]
 D.  dsd_lsbf_planar          DSD (Direct Stream Digital), least signifi [AUDIO]
 D.  dsd_msbf                 DSD (Direct Stream Digital), most signific [AUDIO]
 D.  dsd_msbf_planar          DSD (Direct Stream Digital), most signific [AUDIO]
 D.  dsicinaudio              Delphine Software International CIN audio  [AUDIO]
 D.  dsicinvideo              Delphine Software International CIN video  [VIDEO]
 D.  dss_sp                   Digital Speech Standard - Standard Play mo [AUDIO]
 D.  dst                      DST (Direct Stream Transfer)               [AUDIO]
 D.  dvaudio                  DV audio                                   [AUDIO]
 DE  dvvideo                  DV (Digital Video)                         [VIDEO]
 D.  dxa                      Feeble Files/ScummVM DXA                   [VIDEO]
 D.  dxtory                   Dxtory                                     [VIDEO]
 DE  dxv                      Resolume DXV                               [VIDEO]
 D.  eia_608                  EIA-608 closed captions                    [SUBTITLE]
 D.  escape124                Escape 124                                 [VIDEO]
 D.  escape130                Escape 130                                 [VIDEO]
 D.  evrc                     EVRC (Enhanced Variable Rate Codec)        [AUDIO]
 D.  fastaudio                MobiClip FastAudio                         [AUDIO]
 DE  ffvhuff                  Huffyuv FFmpeg variant                     [VIDEO]
 D.  fic                      Mirillis FIC                               [VIDEO]
 DE  fits                     FITS (Flexible Image Transport System)     [VIDEO]
 DE  flashsv                  Flash Screen Video v1                      [VIDEO]
 DE  flashsv2                 Flash Screen Video v2                      [VIDEO]
 D.  flic                     Autodesk Animator Flic video               [VIDEO]
 DE  flv1                     FLV / Sorenson Spark / Sorenson H.263 (Fla [VIDEO]
 D.  fmvc                     FM Screen Capture Codec                    [VIDEO]
 D.  fraps                    Fraps                                      [VIDEO]
 D.  frwu                     Forward Uncompressed                       [VIDEO]
 D.  ftr                      FTR Voice                                  [AUDIO]
 D.  g2m                      Go2Meeting                                 [VIDEO]
 DE  g723_1                   G.723.1                                    [AUDIO]
 D.  g728                     G.728                                      [AUDIO]
 D.  g729                     G.729                                      [AUDIO]
 D.  gdv                      Gremlin Digital Video                      [VIDEO]
 D.  gem                      GEM Raster image                           [VIDEO]
 D.  gremlin_dpcm             DPCM Gremlin                               [AUDIO]
 D.  gsm                      GSM                                        [AUDIO]
 D.  gsm_ms                   GSM Microsoft variant                      [AUDIO]
 DE  h261                     H.261                                      [VIDEO]
 DE  h263                     H.263 / H.263-1996, H.263+ / H.263-1998 /  [VIDEO]
 D.  h263i                    Intel H.263                                [VIDEO]
 DE  h263p                    H.263+ / H.263-1998 / H.263 version 2      [VIDEO]
 D.  hap                      Vidvox Hap                                 [VIDEO]
 D.  hca                      CRI HCA                                    [AUDIO]
 D.  hcom                     HCOM Audio                                 [AUDIO]
 D.  hdmv_pgs_subtitle        HDMV Presentation Graphic Stream subtitles [SUBTITLE]
 DE  hdr                      HDR (Radiance RGBE format) image           [VIDEO]
 D.  hnm4video                HNM 4 video                                [VIDEO]
 D.  hq_hqa                   Canopus HQ/HQA                             [VIDEO]
 D.  hqx                      Canopus HQX                                [VIDEO]
 DE  huffyuv                  HuffYUV                                    [VIDEO]
 D.  hymt                     HuffYUV MT                                 [VIDEO]
 D.  iac                      IAC (Indeo Audio Coder)                    [AUDIO]
 D.  idcin                    id Quake II CIN video                      [VIDEO]
 D.  idf                      iCEDraw text                               [VIDEO]
 D.  iff_ilbm                 IFF ACBM/ANIM/DEEP/ILBM/PBM/RGB8/RGBN      [VIDEO]
 D.  ilbc                     iLBC (Internet Low Bitrate Codec)          [AUDIO]
 D.  imc                      IMC (Intel Music Coder)                    [AUDIO]
 D.  imm4                     Infinity IMM4                              [VIDEO]
 D.  imm5                     Infinity IMM5                              [VIDEO]
 D.  indeo2                   Intel Indeo 2                              [VIDEO]
 D.  indeo3                   Intel Indeo 3                              [VIDEO]
 D.  indeo4                   Intel Indeo Video Interactive 4            [VIDEO]
 D.  indeo5                   Intel Indeo Video Interactive 5            [VIDEO]
 D.  interplay_dpcm           DPCM Interplay                             [AUDIO]
 D.  interplayacm             Interplay ACM                              [AUDIO]
 D.  interplayvideo           Interplay MVE video                        [VIDEO]
 D.  ipu                      IPU Video                                  [VIDEO]
 D.  jacosub                  JACOsub subtitle                           [SUBTITLE]
 D.  jv                       Bitmap Brothers JV video                   [VIDEO]
 D.  kgv1                     Kega Game Video                            [VIDEO]
 D.  kmvc                     Karl Morton's video codec                  [VIDEO]
 D.  lagarith                 Lagarith lossless                          [VIDEO]
 D.  lead                     LEAD MCMP                                  [VIDEO]
 D.  loco                     LOCO                                       [VIDEO]
 D.  lscr                     LEAD Screen Capture                        [VIDEO]
 D.  m101                     Matrox Uncompressed SD                     [VIDEO]
 D.  mace3                    MACE (Macintosh Audio Compression/Expansio [AUDIO]
 D.  mace6                    MACE (Macintosh Audio Compression/Expansio [AUDIO]
 D.  mad                      Electronic Arts Madcow Video               [VIDEO]
 DE  magicyuv                 MagicYUV video                             [VIDEO]
 D.  mdec                     Sony PlayStation MDEC (Motion DECoder)     [VIDEO]
 D.  media100                 Media 100i                                 [VIDEO]
 D.  metasound                Voxware MetaSound                          [AUDIO]
 D.  microdvd                 MicroDVD subtitle                          [SUBTITLE]
 D.  mimic                    Mimic                                      [VIDEO]
 D.  misc4                    Micronas SC-4 Audio                        [AUDIO]
 D.  mjpegb                   Apple MJPEG-B                              [VIDEO]
 DE  mlp                      MLP (Meridian Lossless Packing)            [AUDIO]
 D.  mmvideo                  American Laser Games MM Video              [VIDEO]
 D.  mobiclip                 MobiClip Video                             [VIDEO]
 D.  motionpixels             Motion Pixels video                        [VIDEO]
 DE  mov_text                 MOV text                                   [SUBTITLE]
 D.  mp1                      MP1 (MPEG audio layer 1)                   [AUDIO]
 D.  mp3adu                   ADU (Application Data Unit) MP3 (MPEG audi [AUDIO]
 D.  mp3on4                   MP3onMP4                                   [AUDIO]
 D.  mp4als                   MPEG-4 Audio Lossless Coding (ALS)         [AUDIO]
 DE  mpeg1video               MPEG-1 video                               [VIDEO]
 DE  mpeg4                    MPEG-4 part 2                              [VIDEO]
 D.  mpl2                     MPL2 subtitle                              [SUBTITLE]
 D.  msa1                     MS ATC Screen                              [VIDEO]
 D.  mscc                     Mandsoft Screen Capture Codec              [VIDEO]
 D.  msmpeg4v1                MPEG-4 part 2 Microsoft variant version 1  [VIDEO]
 DE  msmpeg4v2                MPEG-4 part 2 Microsoft variant version 2  [VIDEO]
 DE  msmpeg4v3                MPEG-4 part 2 Microsoft variant version 3  [VIDEO]
 D.  msnsiren                 MSN Siren                                  [AUDIO]
 D.  msp2                     Microsoft Paint (MSP) version 2            [VIDEO]
 DE  msrle                    Microsoft RLE                              [VIDEO]
 D.  mss1                     MS Screen 1                                [VIDEO]
 D.  mss2                     MS Windows Media Video V9 Screen           [VIDEO]
 DE  msvideo1                 Microsoft Video 1                          [VIDEO]
 D.  mszh                     LCL (LossLess Codec Library) MSZH          [VIDEO]
 D.  mts2                     MS Expression Encoder Screen               [VIDEO]
 D.  musepack7                Musepack SV7                               [AUDIO]
 D.  musepack8                Musepack SV8                               [AUDIO]
 D.  mv30                     MidiVid 3.0                                [VIDEO]
 D.  mvc1                     Silicon Graphics Motion Video Compressor 1 [VIDEO]
 D.  mvc2                     Silicon Graphics Motion Video Compressor 2 [VIDEO]
 D.  mvdv                     MidiVid VQ                                 [VIDEO]
 D.  mvha                     MidiVid Archive Codec                      [VIDEO]
 D.  mwsc                     MatchWare Screen Capture Codec             [VIDEO]
 D.  mxpeg                    Mobotix MxPEG video                        [VIDEO]
 DE  nellymoser               Nellymoser Asao                            [AUDIO]
 D.  notchlc                  NotchLC                                    [VIDEO]
 D.  nuv                      NuppelVideo/RTJPEG                         [VIDEO]
 D.  osq                      OSQ (Original Sound Quality)               [AUDIO]
 D.  paf_audio                Amazing Studio Packed Animation File Audio [AUDIO]
 D.  paf_video                Amazing Studio Packed Animation File Video [VIDEO]
 DE  pam                      PAM (Portable AnyMap) image                [VIDEO]
 DE  pcm_bluray               PCM signed 16|20|24-bit big-endian for Blu [AUDIO]
 DE  pcm_dvd                  PCM signed 20|24-bit big-endian            [AUDIO]
 D.  pcm_f16le                PCM 16.8 floating point little-endian      [AUDIO]
 D.  pcm_f24le                PCM 24.0 floating point little-endian      [AUDIO]
 D.  pcm_lxf                  PCM signed 20-bit little-endian planar     [AUDIO]
 DE  pcm_f64be                PCM 64-bit floating point big-endian       [AUDIO]
 DE  pcm_f64le                PCM 64-bit floating point little-endian    [AUDIO]
 DE  pcm_s24daud              PCM D-Cinema audio signed 24-bit           [AUDIO]
 DE  pcm_s64be                PCM signed 64-bit big-endian               [AUDIO]
 DE  pcm_s64le                PCM signed 64-bit little-endian            [AUDIO]
 DE  pcm_s8                   PCM signed 8-bit                           [AUDIO]
 DE  pcm_s8_planar            PCM signed 8-bit planar                    [AUDIO]
 D.  pcm_sga                  PCM SGA                                    [AUDIO]
 DE  pcm_u16be                PCM unsigned 16-bit big-endian             [AUDIO]
 DE  pcm_u16le                PCM unsigned 16-bit little-endian          [AUDIO]
 DE  pcm_u24be                PCM unsigned 24-bit big-endian             [AUDIO]
 DE  pcm_u24le                PCM unsigned 24-bit little-endian          [AUDIO]
 DE  pcm_u32be                PCM unsigned 32-bit big-endian             [AUDIO]
 DE  pcm_u32le                PCM unsigned 32-bit little-endian          [AUDIO]
 DE  pcm_u8                   PCM unsigned 8-bit                         [AUDIO]
 DE  pcm_vidc                 PCM Archimedes VIDC                        [AUDIO]
 DE  pcx                      PC Paintbrush PCX image                    [VIDEO]
 D.  pdv                      PDV (PlayDate Video)                       [VIDEO]
 DE  pfm                      PFM (Portable FloatMap) image              [VIDEO]
 DE  pgm                      PGM (Portable GrayMap) image               [VIDEO]
 DE  pgmyuv                   PGMYUV (Portable GrayMap YUV) image        [VIDEO]
 D.  pgx                      PGX (JPEG2000 Test Format)                 [VIDEO]
 DE  phm                      PHM (Portable HalfFloatMap) image          [VIDEO]
 D.  photocd                  Kodak Photo CD                             [VIDEO]
 D.  pictor                   Pictor/PC Paint                            [VIDEO]
 D.  pixlet                   Apple Pixlet                               [VIDEO]
 D.  pjs                      PJS (Phoenix Japanimation Society) subtitl [SUBTITLE]
 D.  prosumer                 Brooktree ProSumer Video                   [VIDEO]
 D.  psd                      Photoshop PSD file                         [VIDEO]
 D.  ptx                      V.Flash PTX image                          [VIDEO]
 D.  qcelp                    QCELP / PureVoice                          [AUDIO]
 D.  qdm2                     QDesign Music Codec 2                      [AUDIO]
 D.  qdmc                     QDesign Music                              [AUDIO]
 D.  qdraw                    Apple QuickDraw                            [VIDEO]
 D.  qoa                      QOA (Quite OK Audio)                       [AUDIO]
 D.  qpeg                     Q-team QPEG                                [VIDEO]
 DE  qtrle                    QuickTime Animation (RLE) video            [VIDEO]
 DE  r10k                     AJA Kona 10-bit RGB Codec                  [VIDEO]
 DE  r210                     Uncompressed RGB 10-bit                    [VIDEO]
 DE  ra_144                   RealAudio 1.0 (14.4K)                      [AUDIO]
 D.  ra_288                   RealAudio 2.0 (28.8K)                      [AUDIO]
 D.  ralf                     RealAudio Lossless                         [AUDIO]
 D.  rasc                     RemotelyAnywhere Screen Capture            [VIDEO]
 DE  rawvideo                 raw video                                  [VIDEO]
 D.  realtext                 RealText subtitle                          [SUBTITLE]
 D.  rka                      RKA (RK Audio)                             [AUDIO]
 D.  rl2                      RL2 video                                  [VIDEO]
 DE  roq                      id RoQ video                               [VIDEO]
 DE  roq_dpcm                 DPCM id RoQ                                [AUDIO]
 DE  rpza                     QuickTime video (RPZA)                     [VIDEO]
 D.  rscc                     innoHeim/Rsupport Screen Capture Codec     [VIDEO]
 D.  rtv1                     RTV1 (RivaTuner Video)                     [VIDEO]
 DE  rv10                     RealVideo 1.0                              [VIDEO]
 DE  rv20                     RealVideo 2.0                              [VIDEO]
 D.  rv30                     RealVideo 3.0                              [VIDEO]
 D.  rv40                     RealVideo 4.0                              [VIDEO]
 D.  rv60                     RealVideo 6.0                              [VIDEO]
 DE  s302m                    SMPTE 302M                                 [AUDIO]
 D.  sami                     SAMI subtitle                              [SUBTITLE]
 D.  sanm                     LucasArts SANM/SMUSH video                 [VIDEO]
 D.  scpr                     ScreenPressor                              [VIDEO]
 D.  screenpresso             Screenpresso                               [VIDEO]
 D.  sdx2_dpcm                DPCM Squareroot-Delta-Exact                [AUDIO]
 D.  sga                      Digital Pictures SGA Video                 [VIDEO]
 DE  sgi                      SGI image                                  [VIDEO]
 D.  sgirle                   SGI RLE 8-bit                              [VIDEO]
 D.  sheervideo               BitJazz SheerVideo                         [VIDEO]
 D.  shorten                  Shorten                                    [AUDIO]
 D.  simbiosis_imx            Simbiosis Interactive IMX Video            [VIDEO]
 D.  sipr                     RealAudio SIPR / ACELP.NET                 [AUDIO]
 D.  siren                    Siren                                      [AUDIO]
 D.  smackaudio               Smacker audio                              [AUDIO]
 D.  smackvideo               Smacker video                              [VIDEO]
 DE  smc                      QuickTime Graphics (SMC)                   [VIDEO]
 D.  smvjpeg                  Sigmatel Motion Video                      [VIDEO]
 DE  snow                     Snow                                       [VIDEO]
 D.  sol_dpcm                 DPCM Sol                                   [AUDIO]
 DE  sonic                    Sonic                                      [AUDIO]
 .E  sonicls                  Sonic lossless                             [AUDIO]
 D.  sp5x                     Sunplus JPEG (SP5X)                        [VIDEO]
 DE  speedhq                  NewTek SpeedHQ                             [VIDEO]
 D.  speex                    Speex                                      [AUDIO]
 D.  srgc                     Screen Recorder Gold Codec                 [VIDEO]
 D.  stl                      Spruce subtitle format                     [SUBTITLE]
 DE  subrip                   SubRip subtitle                            [SUBTITLE]
 D.  subviewer                SubViewer subtitle                         [SUBTITLE]
 D.  subviewer1               SubViewer v1 subtitle                      [SUBTITLE]
 DE  sunrast                  Sun Rasterfile image                       [VIDEO]
 DE  svq1                     Sorenson Vector Quantizer 1 / Sorenson Vid [VIDEO]
 D.  svq3                     Sorenson Vector Quantizer 3 / Sorenson Vid [VIDEO]
 D.  tak                      TAK (Tom's lossless Audio Kompressor)      [AUDIO]
 DE  targa                    Truevision Targa image                     [VIDEO]
 D.  targa_y216               Pinnacle TARGA CineWave YUV16              [VIDEO]
 D.  tdsc                     TDSC                                       [VIDEO]
 D.  tgq                      Electronic Arts TGQ video                  [VIDEO]
 D.  tgv                      Electronic Arts TGV video                  [VIDEO]
 D.  thp                      Nintendo Gamecube THP video                [VIDEO]
 D.  tiertexseqvideo          Tiertex Limited SEQ video                  [VIDEO]
 D.  tmv                      8088flex TMV                               [VIDEO]
 D.  tqi                      Electronic Arts TQI video                  [VIDEO]
 D.  truemotion1              Duck TrueMotion 1.0                        [VIDEO]
 D.  truemotion2              Duck TrueMotion 2.0                        [VIDEO]
 D.  truemotion2rt            Duck TrueMotion 2.0 Real Time              [VIDEO]
 D.  truespeech               DSP Group TrueSpeech                       [AUDIO]
 D.  tscc                     TechSmith Screen Capture Codec             [VIDEO]
 D.  tscc2                    TechSmith Screen Codec 2                   [VIDEO]
 DE  tta                      TTA (True Audio)                           [AUDIO]
 .E  ttml                     Timed Text Markup Language                 [SUBTITLE]
 D.  twinvq                   VQF TwinVQ                                 [AUDIO]
 D.  txd                      Renderware TXD (TeXture Dictionary) image  [VIDEO]
 D.  ulti                     IBM UltiMotion                             [VIDEO]
 DE  utvideo                  Ut Video                                   [VIDEO]
 DE  v210                     Uncompressed 4:2:2 10-bit                  [VIDEO]
 D.  v210x                    Uncompressed 4:2:2 10-bit                  [VIDEO]
 DE  v308                     Uncompressed packed 4:4:4                  [VIDEO]
 DE  v408                     Uncompressed packed QT 4:4:4:4             [VIDEO]
 DE  v410                     Uncompressed 4:4:4 10-bit                  [VIDEO]
 D.  vb                       Beam Software VB                           [VIDEO]
 D.  vble                     VBLE Lossless Codec                        [VIDEO]
 DE  vbn                      Vizrt Binary Image                         [VIDEO]
 D.  vc1image                 Windows Media Video 9 Image v2             [VIDEO]
 D.  vcr1                     ATI VCR1                                   [VIDEO]
 D.  vixl                     Miro VideoXL                               [VIDEO]
 D.  vmdaudio                 Sierra VMD audio                           [AUDIO]
 D.  vmdvideo                 Sierra VMD video                           [VIDEO]
 D.  vmix                     vMix Video                                 [VIDEO]
 D.  vmnc                     VMware Screen Codec / VMware Video         [VIDEO]
 D.  vp3                      On2 VP3                                    [VIDEO]
 D.  vp4                      On2 VP4                                    [VIDEO]
 D.  vp5                      On2 VP5                                    [VIDEO]
 D.  vp6                      On2 VP6                                    [VIDEO]
 D.  vp6a                     On2 VP6 (Flash version, with alpha channel [VIDEO]
 D.  vp6f                     On2 VP6 (Flash version)                    [VIDEO]
 D.  vp7                      On2 VP7                                    [VIDEO]
 D.  vplayer                  VPlayer subtitle                           [SUBTITLE]
 D.  vqc                      ViewQuest VQC                              [VIDEO]
 D.  wady_dpcm                DPCM Marble WADY                           [AUDIO]
 D.  wavarc                   Waveform Archiver                          [AUDIO]
 D.  wavesynth                Wave synthesis pseudo-codec                [AUDIO]
 DE  wavpack                  WavPack                                    [AUDIO]
 DE  wbmp                     WBMP (Wireless Application Protocol Bitmap [VIDEO]
 D.  wcmv                     WinCAM Motion Video                        [VIDEO]
 D.  westwood_snd1            Westwood Audio (SND1)                      [AUDIO]
 D.  wmalossless              Windows Media Audio Lossless               [AUDIO]
 D.  wmapro                   Windows Media Audio 9 Professional         [AUDIO]
 DE  wmav1                    Windows Media Audio 1                      [AUDIO]
 DE  wmav2                    Windows Media Audio 2                      [AUDIO]
 D.  wmavoice                 Windows Media Audio Voice                  [AUDIO]
 DE  wmv1                     Windows Media Video 7                      [VIDEO]
 DE  wmv2                     Windows Media Video 8                      [VIDEO]
 D.  wmv3                     Windows Media Video 9                      [VIDEO]
 D.  wmv3image                Windows Media Video 9 Image                [VIDEO]
 D.  wnv1                     Winnov WNV1                                [VIDEO]
 DE  wrapped_avframe          AVFrame to AVPacket passthrough            [VIDEO]
 D.  ws_vqa                   Westwood Studios VQA (Vector Quantized Ani [VIDEO]
 D.  xan_dpcm                 DPCM Xan                                   [AUDIO]
 D.  xan_wc3                  Wing Commander III / Xan                   [VIDEO]
 D.  xan_wc4                  Wing Commander IV / Xxan                   [VIDEO]
 D.  xbin                     eXtended BINary text                       [VIDEO]
 DE  xbm                      XBM (X BitMap) image                       [VIDEO]
 DE  xface                    X-face image                               [VIDEO]
 D.  xma1                     Xbox Media Audio 1                         [AUDIO]
 D.  xma2                     Xbox Media Audio 2                         [AUDIO]
 D.  xpm                      XPM (X PixMap) image                       [VIDEO]
 DE  xsub                     XSUB                                       [SUBTITLE]
 DE  xwd                      XWD (X Window Dump) image                  [VIDEO]
 DE  y41p                     Uncompressed YUV 4:1:1 12-bit              [VIDEO]
 D.  ylc                      YUY2 Lossless Codec                        [VIDEO]
 D.  yop                      Psygnosis YOP Video                        [VIDEO]
 D.  zerocodec                ZeroCodec Lossless Video                   [VIDEO]
 DE  zlib                     LCL (LossLess Codec Library) ZLIB          [VIDEO]
 DE  zmbv                     Zip Motion Blocks Video                    [VIDEO]

# Onboard

The is my hard fork of the [ffmpeg-go](https://github.com/csnewman/ffmpeg-go) project called ffmpeg-statigo.

ffmpeg-statigo has been updated to Go 1.24, all dependencies uplifted to current versions, GitHub CI is fixed so the ffmpeg static libraries are built for all supported architectures and the generator has been updated to support current `clang` on Linux. Mostly critically this project has been updated from FFmpeg 6.1 to FFmpeg 8.0 including the addition os `x265`, `dav1d`, `rav1e` and hardware acceleration support for NVENC, QuickSync and Vulkan to complement the exist VideoToolbox support.

The git history has been purged all tags and historical commits of the static ffmpeg libraries as we are preparing this project to be launched under it's new name in a different GitHUb organisation. This brought the git repo down in size from 1.9GB to 9MB.

NixOS is the host development workstation. There is a `flake.nix` in the project to enable the required software in a Nix devevelopment shell, which is automatically activated via `direnv`. The `justfile` is used for all build and test commands. `fish` is the default shell.

Read the `README.md`, `TODO.md` and analyse the code to get a full understanding of the project. Let me know when you are ready to collaborate.
