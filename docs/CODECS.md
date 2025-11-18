# ffmpeg-statigo

## CODECS

```
 DE  NAME                     DESCRIPTION                                TYPE

 DE  aac                      AAC (Advanced Audio Coding)                [AUDIO]
 D.  aac_latm                 AAC LATM (Advanced Audio Coding LATM synta [AUDIO]
 DE  ac3                      ATSC A/52A (AC-3)                          [AUDIO]
 DE  alac                     ALAC (Apple Lossless Audio Codec)          [AUDIO]
 DE  anull                    Null audio codec                           [AUDIO]
 DE  apng                     APNG (Animated Portable Network Graphics)  [VIDEO]
 DE  av1                      Alliance for Open Media AV1                [VIDEO]
 DE  cfhd                     GoPro CineForm HD                          [VIDEO]
 D.  dirac                    Dirac                                      [VIDEO]
 DE  dnxhd                    VC3/DNxHD                                  [VIDEO]
 DE  dts                      DCA (DTS Coherent Acoustics)               [AUDIO]
 DE  dvb_subtitle             DVB subtitles                              [SUBTITLE]
 DE  dvd_subtitle             DVD subtitles                              [SUBTITLE]
 DE  eac3                     ATSC A/52B (AC-3, E-AC-3)                  [AUDIO]
 DE  exr                      OpenEXR image                              [VIDEO]
 DE  ffv1                     FFmpeg video codec #1                      [VIDEO]
 DE  flac                     FLAC (Free Lossless Audio Codec)           [AUDIO]
 DE  gif                      CompuServe GIF (Graphics Interchange Forma [VIDEO]
 D.  h263                     H.263 / H.263-1996, H.263+ / H.263-1998 /  [VIDEO]
 DE  h264                     H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10  [VIDEO]
 DE  hevc                     H.265 / HEVC (High Efficiency Video Coding [VIDEO]
 DE  mjpeg                    Motion JPEG                                [VIDEO]
 D.  mjpegb                   Apple MJPEG-B                              [VIDEO]
 DE  mp2                      MP2 (MPEG audio layer 2)                   [AUDIO]
 DE  mp3                      MP3 (MPEG audio layer 3)                   [AUDIO]
 DE  mpeg2video               MPEG-2 video                               [VIDEO]
 D.  mpeg4                    MPEG-4 part 2                              [VIDEO]
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
 D.  theora                   Theora                                     [VIDEO]
 DE  tiff                     TIFF image                                 [VIDEO]
 DE  truehd                   TrueHD                                     [AUDIO]
 D.  vc1                      SMPTE VC-1                                 [VIDEO]
 DE  vnull                    Null video codec                           [VIDEO]
 DE  vorbis                   Vorbis                                     [AUDIO]
 D.  vp3                      On2 VP3                                    [VIDEO]
 D.  vp8                      On2 VP8                                    [VIDEO]
 DE  vp9                      Google VP9                                 [VIDEO]
 D.  vvc                      H.266 / VVC (Versatile Video Coding)       [VIDEO]
 DE  webp                     WebP                                       [VIDEO]
 DE  webvtt                   WebVTT subtitle                            [SUBTITLE]
```

**Summary:**
- Total codecs: 59
- Decoders: 59 (Video: 29, Audio: 27, Subtitle: 3, Other: 0)
- Encoders: 48 (Video: 19, Audio: 26, Subtitle: 3, Other: 0)

**Flags:**
- D - Decoder available
- E - Encoder available

## HARDWARE ACCELERATORS

```
    NAME

    cuda
    qsv
    vulkan
```

### HARDWARE CODECS

```
 DEH  NAME                     DESCRIPTION                                TYPE PRESENT

 D..  av1_cuvid                Nvidia CUVID AV1 decoder                   [VIDEO]  N
 .E.  av1_nvenc                NVIDIA NVENC av1 encoder                   [VIDEO]  N
 D..  av1_qsv                  AV1 video (Intel Quick Sync Video accelera [VIDEO]  N
 .E.  av1_qsv                  AV1 (Intel Quick Sync Video acceleration)  [VIDEO]  N
 .E.  av1_vulkan               AV1 (Vulkan)                               [VIDEO]  N
 .E.  ffv1_vulkan              FFmpeg video codec #1 (Vulkan)             [VIDEO]  N
 D..  h264_cuvid               Nvidia CUVID H264 decoder                  [VIDEO]  N
 .E.  h264_nvenc               NVIDIA NVENC H.264 encoder                 [VIDEO]  N
 D..  h264_qsv                 H264 video (Intel Quick Sync Video acceler [VIDEO]  N
 .E.  h264_qsv                 H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10  [VIDEO]  N
 .E.  h264_vulkan              H.264/AVC (Vulkan)                         [VIDEO]  N
 D..  hevc_cuvid               Nvidia CUVID HEVC decoder                  [VIDEO]  N
 .E.  hevc_nvenc               NVIDIA NVENC hevc encoder                  [VIDEO]  N
 D..  hevc_qsv                 HEVC video (Intel Quick Sync Video acceler [VIDEO]  N
 .E.  hevc_qsv                 HEVC (Intel Quick Sync Video acceleration) [VIDEO]  N
 .E.  hevc_vulkan              H.265/HEVC (Vulkan)                        [VIDEO]  N
 D..  mjpeg_cuvid              Nvidia CUVID MJPEG decoder                 [VIDEO]  N
 D..  mjpeg_qsv                MJPEG video (Intel Quick Sync Video accele [VIDEO]  N
 .E.  mjpeg_qsv                MJPEG (Intel Quick Sync Video acceleration [VIDEO]  N
 D..  mpeg2_cuvid              Nvidia CUVID MPEG2VIDEO decoder            [VIDEO]  N
 .E.  mpeg2_qsv                MPEG-2 video (Intel Quick Sync Video accel [VIDEO]  N
 D..  mpeg2_qsv                MPEG2VIDEO video (Intel Quick Sync Video a [VIDEO]  N
 D..  mpeg4_cuvid              Nvidia CUVID MPEG4 decoder                 [VIDEO]  N
 D..  vc1_cuvid                Nvidia CUVID VC1 decoder                   [VIDEO]  N
 D..  vc1_qsv                  VC1 video (Intel Quick Sync Video accelera [VIDEO]  N
 D..  vp8_cuvid                Nvidia CUVID VP8 decoder                   [VIDEO]  N
 D..  vp8_qsv                  VP8 video (Intel Quick Sync Video accelera [VIDEO]  N
 D..  vp9_cuvid                Nvidia CUVID VP9 decoder                   [VIDEO]  N
 .E.  vp9_qsv                  VP9 video (Intel Quick Sync Video accelera [VIDEO]  N
 D..  vp9_qsv                  VP9 video (Intel Quick Sync Video accelera [VIDEO]  N
 D..  vvc_qsv                  VVC video (Intel Quick Sync Video accelera [VIDEO]  N

Summary:
  Total hardware device types: 3
  Total hwaccels: 0
  Total hardware codecs: 31
  Hardware decoders: 18
  Hardware encoders: 13
```

## FORMATS

```
DE  NAME                     DESCRIPTION                                CODECS                              MIME TYPE

D.  aac                      raw ADTS AAC (Advanced Audio Coding)
DE  ac3                      raw AC-3                                   audio:ac3                           audio/x-ac3
.E  adts                     ADTS AAC (Advanced Audio Coding)           audio:aac                           audio/aac
DE  aiff                     Audio IFF                                  video:png,audio:pcm_s16be           audio/aiff
DE  apng                     Animated Portable Network Graphics         video:apng                          image/png
D.  asf                      ASF (Advanced / Active Streaming Format)
DE  avi                      AVI (Audio Video Interleaved)              video:mpeg4,audio:mp3               video/x-msvideo
.E  avif                     AVIF                                       video:av1                           image/avif
DE  caf                      Apple CAF (Core Audio Format)              audio:pcm_s16be                     audio/x-caf
DE  dash                     DASH Muxer                                 video:h264,audio:aac
D.  dirac                    raw Dirac
DE  dnxhd                    raw DNxHD (SMPTE VC-3)                     video:dnxhd
DE  dts                      raw DTS                                    audio:dts                           audio/x-dca
D.  dtshd                    raw DTS-HD
DE  eac3                     raw E-AC-3                                 audio:eac3                          audio/x-eac3
DE  flac                     raw FLAC                                   video:png,audio:flac                audio/x-flac
.E  flv                      FLV (Flash Video)                          video:flv1,audio:mp3                video/x-flv
DE  gif                      CompuServe Graphics Interchange Format (GI video:gif                           image/gif
DE  h264                     raw H.264 video                            video:h264
.E  hds                      HDS Muxer                                  video:h264,audio:aac
DE  hevc                     raw HEVC video                             video:hevc
DE  hls                      Apple HTTP Live Streaming                  video:h264,audio:aac,subtitle:webvt
DE  iamf                     Raw Immersive Audio Model and Formats      audio:opus
DE  image2                   image2 sequence                            video:mjpeg
DE  image2pipe               piped image2 sequence                      video:mjpeg
.E  ipod                     iPod H.264 MP4 (MPEG-4 Part 14)            video:h264,audio:aac                video/mp4
.E  ismv                     ISMV/ISMA (Smooth Streaming)               video:h264,audio:aac                video/mp4
.E  latm                     LOAS/LATM                                  audio:aac                           audio/MP4A-LATM
DE  m4v                      raw MPEG-4 video                           video:mpeg4
.E  matroska                 Matroska                                   video:h264,audio:ac3,subtitle:ass   video/x-matroska
D.  matroska,webm            Matroska / WebM
DE  mjpeg                    raw MJPEG video                            video:mjpeg                         video/x-mjpeg
.E  mov                      QuickTime / MOV                            video:h264,audio:aac
D.  mov,mp4,m4a,3gp,3g2,mj2  QuickTime / MOV
.E  mp2                      MP2 (MPEG audio layer 2)                   audio:mp2                           audio/mpeg
DE  mp3                      MP3 (MPEG audio layer 3)                   video:png,audio:mp3                 audio/mpeg
.E  mp4                      MP4 (MPEG-4 Part 14)                       video:h264,audio:aac                video/mp4
.E  mpeg2video               raw MPEG-2 video                           video:mpeg2video
DE  mpegts                   MPEG-TS (MPEG-2 Transport Stream)          video:mpeg2video,audio:mp2          video/MP2T
DE  mpjpeg                   MIME multipart JPEG                        video:mjpeg                         multipart/x-mixed-re
.E  null                     raw null video                             video:wrapped_avframe,audio:pcm_s16
DE  obu                      AV1 low overhead OBU                       video:av1
.E  oga                      Ogg Audio                                  audio:flac                          audio/ogg
DE  ogg                      Ogg                                        video:theora,audio:flac             application/ogg
.E  ogv                      Ogg Video                                  video:vp8,audio:flac                video/ogg
.E  opus                     Ogg Opus                                   audio:opus                          audio/ogg
D.  rm                       RealMedia
DE  rtp                      RTP output                                 video:mpeg4,audio:pcm_mulaw
.E  rtp_mpegts               RTP/mpegts output format                   video:mpeg4,audio:aac
DE  rtsp                     RTSP output                                video:mpeg4,audio:aac
DE  sap                      SAP output                                 video:mpeg4,audio:aac
D.  sdp                      SDP
.E  smoothstreaming          Smooth Streaming Muxer                     video:h264,audio:aac
DE  spdif                    IEC 61937 (used on S/PDIF - IEC958)        audio:ac3
DE  truehd                   raw TrueHD                                 audio:truehd
DE  vc1                      raw VC-1 video                             video:vc1
DE  vvc                      raw H.266/VVC video                        video:vvc
DE  wav                      WAV / WAVE (Waveform Audio)                audio:pcm_s16le                     audio/x-wav
.E  webm                     WebM                                       video:vp8,audio:opus,subtitle:webvt video/webm
.E  webp                     WebP                                       video:webp
DE  webvtt                   WebVTT subtitle                            subtitle:webvtt                     text/vtt
.E  whip                     WHIP(WebRTC-HTTP ingestion protocol) muxer video:h264,audio:opus
```

**Summary:**
- Total formats: 62
- Total demuxers: 41
- Total muxers: 54

## FILTERS

```
 TSHM  NAME                     DESCRIPTION                                TYPE

 .... a3dscope                 Convert input audio to 3d scope video outp AUDIO
 .S.. aap                      Apply Affine Projection algorithm to first AUDIO
 ...M abench                   Benchmark part of a filtergraph.           AUDIO
 .... abitscope                Convert input audio to audio bit scope vid AUDIO
 .... abuffer                  Buffer audio frames, and make them accessi AUDIO
 .... abuffersink              Buffer audio frames, and make them availab AUDIO
 .... acompressor              Audio compressor.                          AUDIO
 .... acontrast                Simple audio dynamic range compression/exp AUDIO
 ...M acopy                    Copy the input audio unchanged to the outp AUDIO
 .... acrossfade               Cross fade two input audio streams.        AUDIO
 .S.. acrossover               Split audio into per-bands streams.        AUDIO
 .... acrusher                 Reduce audio bit resolution.               AUDIO
 ...M acue                     Delay filtering to match a cue.            AUDIO
 ...M addroi                   Add region of interest to frame.           VIDEO
 .S.. adeclick                 Remove impulsive noise from input audio.   AUDIO
 .S.. adeclip                  Remove clipping from input audio.          AUDIO
 TS.. adecorrelate             Apply decorrelation to input audio.        AUDIO
 .... adelay                   Delay one or more audio channels.          AUDIO
 TS.. adenorm                  Remedy denormals by adding extremely low-l AUDIO
 .... aderivative              Compute derivative of input audio.         AUDIO
 .... adrawgraph               Draw a graph using input audio metadata.   AUDIO
 .S.. adrc                     Audio Spectral Dynamic Range Controller.   AUDIO
 .S.. adynamicequalizer        Apply Dynamic Equalization of input audio. AUDIO
 .... adynamicsmooth           Apply Dynamic Smoothing of input audio.    AUDIO
 .... aecho                    Add echoing to the audio.                  AUDIO
 TS.. aemphasis                Audio emphasis.                            AUDIO
 T... aeval                    Filter audio signal according to a specifi AUDIO
 .... aevalsrc                 Generate an audio signal generated by an e AUDIO
 .... aexciter                 Enhance high frequency part of audio.      AUDIO
 T... afade                    Fade in/out input audio.                   AUDIO
 .... afdelaysrc               Generate a Fractional delay FIR coefficien AUDIO
 .S.. afftdn                   Denoise audio samples using FFT.           AUDIO
 .S.. afftfilt                 Apply arbitrary expressions to samples in  AUDIO
 .S.. afir                     Apply Finite Impulse Response filter with  AUDIO
 .... afireqsrc                Generate a FIR equalizer coefficients audi AUDIO
 .... afirsrc                  Generate a FIR coefficients audio stream.  AUDIO
 ...M aformat                  Convert the input audio to one of the spec AUDIO
 TS.. afreqshift               Apply frequency shifting to input audio.   AUDIO
 .S.. afwtdn                   Denoise audio stream using Wavelets.       AUDIO
 T... agate                    Audio gate.                                AUDIO
 .... agraphmonitor            Show various filtergraph stats.            AUDIO
 .... ahistogram               Convert input audio to histogram video out AUDIO
 .S.. aiir                     Apply Infinite Impulse Response filter wit AUDIO
 .... aintegral                Compute integral of input audio.           AUDIO
 .... ainterleave              Temporally interleave audio inputs.        AUDIO
 .... alatency                 Report audio filtering latency.            AUDIO
 T... alimiter                 Audio lookahead limiter.                   AUDIO
 .S.. allpass                  Apply a two-pole all-pass filter.          AUDIO
 .... allrgb                   Generate all RGB colors.                   VIDEO
 .... allyuv                   Generate all yuv colors.                   VIDEO
 .... aloop                    Loop audio samples.                        AUDIO
 .... alphaextract             Extract an alpha channel as a grayscale im VIDEO
 .... alphamerge               Copy the luma value of the second input in VIDEO
 .... amerge                   Merge two or more audio streams into a sin AUDIO
 T..M ametadata                Manipulate audio frame metadata.           AUDIO
 .... amix                     Audio mixing.                              AUDIO
 .... amovie                   Read audio from a movie source.            UNKNOWN
 .S.. amplify                  Amplify changes between successive video f VIDEO
 .... amultiply                Multiply two audio streams.                AUDIO
 .S.. anequalizer              Apply high-order audio parametric multi ba AUDIO
 .S.. anlmdn                   Reduce broadband noise from stream using N AUDIO
 .S.. anlmf                    Apply Normalized Least-Mean-Fourth algorit AUDIO
 .S.. anlms                    Apply Normalized Least-Mean-Squares algori AUDIO
 .... anoisesrc                Generate a noise audio signal.             AUDIO
 ...M anull                    Pass the source unchanged to the output.   AUDIO
 .... anullsink                Do absolutely nothing with the input audio AUDIO
 .... anullsrc                 Null audio source, return empty audio fram AUDIO
 .... apad                     Pad audio with silence.                    AUDIO
 T..M aperms                   Set permissions for the output audio frame AUDIO
 .... aphasemeter              Convert input audio to phase meter video o AUDIO
 .... aphaser                  Add a phasing effect to the audio.         AUDIO
 TS.. aphaseshift              Apply phase shifting to input audio.       AUDIO
 .S.M apsnr                    Measure Audio Peak Signal-to-Noise Ratio.  AUDIO
 .S.. apsyclip                 Audio Psychoacoustic Clipper.              AUDIO
 .... apulsator                Audio pulsator.                            AUDIO
 ...M arealtime                Slow down filtering to match realtime.     AUDIO
 .... aresample                Resample audio data.                       AUDIO
 .... areverse                 Reverse an audio clip.                     AUDIO
 .S.. arls                     Apply Recursive Least Squares algorithm to AUDIO
 .S.. arnndn                   Reduce noise from speech using Recurrent N AUDIO
 .S.M asdr                     Measure Audio Signal-to-Distortion Ratio.  AUDIO
 ...M asegment                 Segment audio stream.                      AUDIO
 .... aselect                  Select audio frames to pass in output.     AUDIO
 ...M asendcmd                 Send commands to filters.                  AUDIO
 .... asetnsamples             Set the number of samples for each output  AUDIO
 ...M asetpts                  Set PTS for the output audio frame.        AUDIO
 ...M asetrate                 Change the sample rate without altering th AUDIO
 ...M asettb                   Set timebase for the audio output link.    AUDIO
 ...M ashowinfo                Show textual information for each audio fr AUDIO
 T..M asidedata                Manipulate audio frame side data.          AUDIO
 .S.M asisdr                   Measure Audio Scale-Invariant Signal-to-Di AUDIO
 TS.. asoftclip                Audio Soft Clipper.                        AUDIO
 .S.. aspectralstats           Show frequency domain statistics about aud AUDIO
 ...M asplit                   Pass on the audio input to N audio outputs AUDIO
 .S.M astats                   Show time domain statistics about audio fr AUDIO
 .... astreamselect            Select audio streams                       UNKNOWN
 .S.. asubboost                Boost subwoofer frequencies.               AUDIO
 TS.. asubcut                  Cut subwoofer frequencies.                 AUDIO
 TS.. asupercut                Cut super frequencies.                     AUDIO
 TS.. asuperpass               Apply high order Butterworth band-pass fil AUDIO
 TS.. asuperstop               Apply high order Butterworth band-stop fil AUDIO
 .S.. atadenoise               Apply an Adaptive Temporal Averaging Denoi VIDEO
 .... atempo                   Adjust audio tempo.                        AUDIO
 TS.. atilt                    Apply spectral tilt to audio.              AUDIO
 ...M atrim                    Pick one continuous section from the input AUDIO
 .S.. avectorscope             Convert input audio to vectorscope video o AUDIO
 T... avgblur                  Apply Average Blur filter.                 VIDEO
 ..H. avgblur_vulkan           Apply avgblur mask to input video          VIDEO
 .... avsynctest               Generate an Audio Video Sync Test.         AUDIO
 .... axcorrelate              Cross-correlate two audio streams.         AUDIO
 TS.. backgroundkey            Turns a static background into transparenc VIDEO
 .S.. bandpass                 Apply a two-pole Butterworth band-pass fil AUDIO
 .S.. bandreject               Apply a two-pole Butterworth band-reject f AUDIO
 .S.. bass                     Boost or cut lower frequencies.            AUDIO
 T..M bbox                     Compute bounding box for each frame.       VIDEO
 ...M bench                    Benchmark part of a filtergraph.           VIDEO
 TS.. bilateral                Apply Bilateral filter.                    VIDEO
 .S.. biquad                   Apply a biquad IIR filter with the given c AUDIO
 T... bitplanenoise            Measure bit plane noise.                   VIDEO
 .S.M blackdetect              Detect video intervals that are (almost) b VIDEO
 ..H. blackdetect_vulkan       Detect video intervals that are (almost) b VIDEO
 ...M blackframe               Detect frames that are (almost) black.     VIDEO
 .S.. blend                    Blend two video frames into each other.    VIDEO
 ..H. blend_vulkan             Blend two video frames in Vulkan           VIDEO
 ...M blockdetect              Blockdetect filter.                        VIDEO
 ...M blurdetect               Blurdetect filter.                         VIDEO
 .S.. bm3d                     Block-Matching 3D denoiser.                VIDEO
 T... boxblur                  Blur the input.                            VIDEO
 .... buffer                   Buffer video frames, and make them accessi VIDEO
 .... buffersink               Buffer video frames, and make them availab VIDEO
 .S.. bwdif                    Deinterlace the input image.               VIDEO
 ..H. bwdif_vulkan             Deinterlace Vulkan frames via bwdif        VIDEO
 TS.. cas                      Contrast Adaptive Sharpen.                 VIDEO
 .... ccrepack                 Repack CEA-708 closed caption metadata     VIDEO
 .... cellauto                 Create pattern generated by an elementary  VIDEO
 .... channelmap               Remap audio channels.                      AUDIO
 .... channelsplit             Split audio into per-channel streams.      AUDIO
 .... chorus                   Add a chorus effect to the audio.          AUDIO
 ..H. chromaber_vulkan         Offset chroma of input video (chromatic ab VIDEO
 TS.. chromahold               Turns a certain color range into gray.     VIDEO
 TS.. chromakey                Turns a certain color into transparency. O VIDEO
 TS.. chromanr                 Reduce chrominance noise.                  VIDEO
 TS.. chromashift              Shift chroma.                              VIDEO
 .... ciescope                 Video CIE scope.                           VIDEO
 T... codecview                Visualize information about some codecs.   VIDEO
 .... color                    Provide an uniformly colored input.        VIDEO
 ..H. color_vulkan             Generate a constant color (Vulkan)         VIDEO
 TS.. colorbalance             Adjust the color balance.                  VIDEO
 TS.. colorchannelmixer        Adjust colors by mixing color channels.    VIDEO
 .... colorchart               Generate color checker chart.              VIDEO
 TS.. colorcontrast            Adjust color contrast between RGB componen VIDEO
 TS.. colorcorrect             Adjust color white balance selectively for VIDEO
 .S.M colordetect              Detect video color properties.             VIDEO
 TS.. colorhold                Turns a certain color range into gray. Ope VIDEO
 TS.. colorize                 Overlay a solid color on the video stream. VIDEO
 TS.. colorkey                 Turns a certain color into transparency. O VIDEO
 TS.. colorlevels              Adjust the color levels.                   VIDEO
 .S.. colormap                 Apply custom Color Maps to video stream.   VIDEO
 TS.. colormatrix              Convert color matrix.                      VIDEO
 TS.. colorspace               Convert between colorspaces.               VIDEO
 .... colorspectrum            Generate colors spectrum.                  VIDEO
 TS.. colortemperature         Adjust color temperature of video.         VIDEO
 .... compand                  Compress or expand audio dynamic range.    AUDIO
 .... compensationdelay        Audio Compensation Delay Line.             AUDIO
 .... concat                   Concatenate audio and video streams.       UNKNOWN
 TS.. convolution              Apply convolution filter.                  VIDEO
 .S.. convolve                 Convolve first video stream with second vi VIDEO
 ...M copy                     Copy the input video unchanged to the outp VIDEO
 .S.M corr                     Calculate the correlation between two vide VIDEO
 .... cover_rect               Find and cover a user specified object.    VIDEO
 .... crop                     Crop the input video.                      VIDEO
 T..M cropdetect               Auto-detect crop size.                     VIDEO
 .... crossfeed                Apply headphone crossfeed filter.          AUDIO
 .S.. crystalizer              Simple audio noise sharpening filter.      AUDIO
 .... cue                      Delay filtering to match a cue.            VIDEO
 TS.. curves                   Adjust components curves.                  VIDEO
 .S.. datascope                Video data analysis.                       VIDEO
 T... dblur                    Apply Directional Blur filter.             VIDEO
 T... dcshift                  Apply a DC shift to the audio.             AUDIO
 TS.. dctdnoiz                 Denoise frames using 2D DCT.               VIDEO
 TS.. deband                   Debands video.                             VIDEO
 T... deblock                  Deblock video.                             VIDEO
 .... decimate                 Decimate frames (post field matching filte VIDEO
 .S.. deconvolve               Deconvolve first video stream with second  VIDEO
 .S.. dedot                    Reduce cross-luminance and cross-color.    VIDEO
 .... deesser                  Apply de-essing to the audio.              AUDIO
 TS.. deflate                  Apply deflate effect.                      VIDEO
 .... deflicker                Remove temporal frame luminance variations VIDEO
 ..H. deinterlace_qsv          Quick Sync Video "deinterlacing"           VIDEO
 .... dejudder                 Remove judder produced by pullup.          VIDEO
 T... delogo                   Remove logo from input video.              VIDEO
 .... deshake                  Stabilize shaky video.                     VIDEO
 TS.. despill                  Despill video.                             VIDEO
 .... detelecine               Apply an inverse telecine pattern.         VIDEO
 .... dialoguenhance           Audio Dialogue Enhancement.                AUDIO
 TS.. dilation                 Apply dilation effect.                     VIDEO
 .S.. displace                 Displace pixels.                           VIDEO
 .S.. doubleweave              Weave input video fields into double numbe VIDEO
 T... drawbox                  Draw a colored box on the input video.     VIDEO
 .... drawgraph                Draw a graph using input video metadata.   VIDEO
 T... drawgrid                 Draw a colored grid on the input video.    VIDEO
 ...M drmeter                  Measure audio dynamic range.               AUDIO
 .S.. dynaudnorm               Dynamic Audio Normalizer.                  AUDIO
 .... earwax                   Widen the stereo image.                    AUDIO
 .... ebur128                  EBU R128 scanner.                          AUDIO
 T... edgedetect               Detect and draw edge.                      VIDEO
 .... elbg                     Apply posterize effect, using the ELBG alg VIDEO
 T..M entropy                  Measure video frames entropy.              VIDEO
 .S.. epx                      Scale the input using EPX algorithm.       VIDEO
 T... eq                       Adjust brightness, contrast, gamma, and sa VIDEO
 .S.. equalizer                Apply two-pole peaking equalization (EQ) f AUDIO
 TS.. erosion                  Apply erosion effect.                      VIDEO
 .S.. estdif                   Apply Edge Slope Tracing deinterlace.      VIDEO
 TS.. exposure                 Adjust exposure of the video stream.       VIDEO
 .... extractplanes            Extract planes as grayscale frames.        VIDEO
 T... extrastereo              Increase difference between stereo audio c AUDIO
 TS.. fade                     Fade in/out input video.                   VIDEO
 .... feedback                 Apply feedback video filter.               VIDEO
 .S.. fftdnoiz                 Denoise frames using 3D FFT.               VIDEO
 TS.. fftfilt                  Apply arbitrary expressions to pixels in f VIDEO
 .... field                    Extract a field from the input video.      VIDEO
 .... fieldhint                Field matching using hints.                VIDEO
 .... fieldmatch               Field matching for inverse telecine.       VIDEO
 T... fieldorder               Set the field order.                       VIDEO
 T... fillborders              Fill borders of the input video.           VIDEO
 ...M find_rect                Find a user specified object.              VIDEO
 .... firequalizer             Finite Impulse Response Equalizer.         AUDIO
 .... flanger                  Apply a flanging effect to the audio.      AUDIO
 ..H. flip_vulkan              Flip both horizontally and vertically      VIDEO
 T... floodfill                Fill area with same color with another col VIDEO
 ...M format                   Convert the input video to one of the spec VIDEO
 ...M fps                      Force constant framerate.                  VIDEO
 .... framepack                Generate a frame packed stereoscopic video VIDEO
 .S.. framerate                Upsamples or downsamples progressive sourc VIDEO
 T..M framestep                Select one frame every N frames.           VIDEO
 ...M freezedetect             Detects frozen video input.                VIDEO
 .... freezeframes             Freeze video frames.                       VIDEO
 .... fspp                     Apply Fast Simple Post-processing filter.  VIDEO
 ...M fsync                    Synchronize video frames from external sou VIDEO
 TS.. gblur                    Apply Gaussian Blur filter.                VIDEO
 ..H. gblur_vulkan             Gaussian Blur in Vulkan                    VIDEO
 TS.. geq                      Apply generic equation to each pixel.      VIDEO
 T... gradfun                  Debands video quickly using gradients.     VIDEO
 .S.. gradients                Draw a gradients.                          VIDEO
 .... graphmonitor             Show various filtergraph stats.            VIDEO
 TS.. grayworld                Adjust white balance using LAB gray world  VIDEO
 TS.. greyedge                 Estimates scene illumination by grey edge  VIDEO
 .S.. guided                   Apply Guided filter.                       VIDEO
 .... haas                     Apply Haas Stereo Enhancer.                AUDIO
 .S.. haldclut                 Adjust colors using a Hald CLUT.           VIDEO
 .... haldclutsrc              Provide an identity Hald CLUT.             VIDEO
 .... hdcd                     Apply High Definition Compatible Digital ( AUDIO
 .S.. headphone                Apply headphone binaural spatialization wi AUDIO
 TS.. hflip                    Horizontally flip the input video.         VIDEO
 .... hflip_vulkan             Horizontally flip the input video in Vulka VIDEO
 .S.. highpass                 Apply a high-pass filter with 3dB point fr AUDIO
 .S.. highshelf                Apply a high shelf filter.                 AUDIO
 .... hilbert                  Generate a Hilbert transform FIR coefficie AUDIO
 T... histeq                   Apply global color histogram equalization. VIDEO
 .... histogram                Compute and draw a histogram.              VIDEO
 .S.. hqdn3d                   Apply a High Quality 3D Denoiser.          VIDEO
 .S.. hqx                      Scale the input by 2, 3 or 4 using the hq* VIDEO
 .S.. hstack                   Stack video inputs horizontally.           VIDEO
 ..H. hstack_qsv               "Quick Sync Video" hstack                  VIDEO
 TS.. hsvhold                  Turns a certain HSV range into gray.       VIDEO
 TS.. hsvkey                   Turns a certain HSV range into transparenc VIDEO
 T... hue                      Adjust the hue and saturation of the input VIDEO
 TS.. huesaturation            Apply hue-saturation-intensity adjustments VIDEO
 .... hwdownload               Download a hardware frame to a normal fram VIDEO
 ..H. hwmap                    Map hardware frames                        VIDEO
 ..H. hwupload                 Upload a normal frame to a hardware frame  VIDEO
 .... hwupload_cuda            Upload a system memory frame to a CUDA dev VIDEO
 .... hysteresis               Grow first stream into second stream by co VIDEO
 .S.M identity                 Calculate the Identity between two video s VIDEO
 ...M idet                     Interlace detect Filter.                   VIDEO
 T... il                       Deinterleave or interleave fields.         VIDEO
 TS.. inflate                  Apply inflate effect.                      VIDEO
 .... interlace                Convert progressive video into interlaced. VIDEO
 ..H. interlace_vulkan         Convert progressive video into interlaced. VIDEO
 .... interleave               Temporally interleave video inputs.        VIDEO
 .... join                     Join multiple audio streams into multi-cha AUDIO
 .... kerndeint                Apply kernel deinterlacing to the input.   VIDEO
 TS.. kirsch                   Apply kirsch operator.                     VIDEO
 .S.. lagfun                   Slowly update darker pixels.               VIDEO
 ...M latency                  Report video filtering latency.            VIDEO
 TS.. lenscorrection           Rectify the image by correcting for lens d VIDEO
 .... life                     Create life.                               VIDEO
 .S.. limitdiff                Apply filtering with limiting difference.  VIDEO
 TS.. limiter                  Limit pixels components to the specified r VIDEO
 .... loop                     Loop video frames.                         VIDEO
 .... loudnorm                 EBU R128 loudness normalization            AUDIO
 .S.. lowpass                  Apply a low-pass filter with 3dB point fre AUDIO
 .S.. lowshelf                 Apply a low shelf filter.                  AUDIO
 TS.. lumakey                  Turns a certain luma into transparency.    VIDEO
 TS.. lut                      Compute and apply a lookup table to the RG VIDEO
 TS.. lut1d                    Adjust colors using a 1D LUT.              VIDEO
 .S.. lut2                     Compute and apply a lookup table from two  VIDEO
 TS.. lut3d                    Adjust colors using a 3D LUT.              VIDEO
 TS.. lutrgb                   Compute and apply a lookup table to the RG VIDEO
 TS.. lutyuv                   Compute and apply a lookup table to the YU VIDEO
 .... mandelbrot               Render a Mandelbrot fractal.               VIDEO
 .S.. maskedclamp              Clamp first stream with second stream and  VIDEO
 .S.. maskedmax                Apply filtering with maximum difference of VIDEO
 .S.. maskedmerge              Merge first stream with second stream usin VIDEO
 .S.. maskedmin                Apply filtering with minimum difference of VIDEO
 .S.. maskedthreshold          Pick pixels comparing absolute difference  VIDEO
 TS.. maskfun                  Create Mask.                               VIDEO
 .... mcdeint                  Apply motion compensating deinterlacing.   VIDEO
 .... mcompand                 Multiband Compress or expand audio dynamic AUDIO
 TS.. median                   Apply Median filter.                       VIDEO
 .... mergeplanes              Merge planes.                              VIDEO
 ...M mestimate                Generate motion vectors.                   VIDEO
 T..M metadata                 Manipulate video frame metadata.           VIDEO
 .... midequalizer             Apply Midway Equalization.                 VIDEO
 .... minterpolate             Frame rate conversion using Motion Interpo VIDEO
 .S.. mix                      Mix video inputs.                          VIDEO
 TS.. monochrome               Convert video to gray using custom color f VIDEO
 .S.. morpho                   Apply Morphological filter.                VIDEO
 .... movie                    Read from a movie source.                  UNKNOWN
 .... mpdecimate               Remove near-duplicate frames.              VIDEO
 .... mptestsrc                Generate various test pattern.             VIDEO
 .S.M msad                     Calculate the MSAD between two video strea VIDEO
 .S.. multiply                 Multiply first video stream with second vi VIDEO
 TS.. negate                   Negate input video.                        VIDEO
 TS.. nlmeans                  Non-local means denoiser.                  VIDEO
 ..H. nlmeans_vulkan           Non-local means denoiser (Vulkan)          VIDEO
 .S.. nnedi                    Apply neural network edge directed interpo VIDEO
 ...M noformat                 Force libavfilter not to use any of the sp VIDEO
 TS.. noise                    Add noise.                                 VIDEO
 .... normalize                Normalize RGB video.                       VIDEO
 ...M null                     Pass the source unchanged to the output.   VIDEO
 .... nullsink                 Do absolutely nothing with the input video VIDEO
 .... nullsrc                  Null video source, return unprocessed vide VIDEO
 T... oscilloscope             2D Video Oscilloscope.                     VIDEO
 .S.. overlay                  Overlay a video source on top of the input VIDEO
 ..H. overlay_qsv              Quick Sync Video overlay.                  VIDEO
 ..H. overlay_vulkan           Overlay a source on top of another         VIDEO
 T... owdenoise                Denoise using wavelets.                    VIDEO
 .... pad                      Pad the input video.                       VIDEO
 .... pal100bars               Generate PAL 100% color bars.              VIDEO
 .... pal75bars                Generate PAL 75% color bars.               VIDEO
 .... palettegen               Find the optimal palette for a given strea VIDEO
 .... paletteuse               Use a palette to downsample an input video VIDEO
 .... pan                      Remix channels with coefficients (panning) AUDIO
 .... perlin                   Generate Perlin noise                      VIDEO
 T..M perms                    Set permissions for the output video frame VIDEO
 TS.. perspective              Correct the perspective of video.          VIDEO
 .... phase                    Phase shift fields.                        VIDEO
 .... photosensitivity         Filter out photosensitive epilepsy seizure VIDEO
 .... pixdesctest              Test pixel format definitions.             VIDEO
 TS.. pixelize                 Pixelize video.                            VIDEO
 T... pixscope                 Pixel data analysis.                       VIDEO
 .... pp7                      Apply Postprocessing 7 filter.             VIDEO
 .S.. premultiply              PreMultiply first stream with first plane  VIDEO
 TS.. prewitt                  Apply prewitt operator.                    VIDEO
 TS.. pseudocolor              Make pseudocolored video frames.           VIDEO
 .S.M psnr                     Calculate the PSNR between two video strea VIDEO
 .... pullup                   Pullup from field sequence to frames.      VIDEO
 ...M qp                       Change video quantization parameters.      VIDEO
 .... random                   Return random frames.                      VIDEO
 TS.M readeia608               Read EIA-608 Closed Caption codes from inp VIDEO
 ...M readvitc                 Read vertical interval timecode and write  VIDEO
 ...M realtime                 Slow down filtering to match realtime.     VIDEO
 .S.. remap                    Remap pixels.                              VIDEO
 TS.. removegrain              Remove grain.                              VIDEO
 T... removelogo               Remove a TV logo based on a mask image.    VIDEO
 .... repeatfields             Hard repeat fields based on MPEG repeat fi VIDEO
 ...M replaygain               ReplayGain scanner.                        AUDIO
 .... reverse                  Reverse a clip.                            VIDEO
 TS.. rgbashift                Shift RGBA.                                VIDEO
 .... rgbtestsrc               Generate RGB test pattern.                 VIDEO
 TS.. roberts                  Apply roberts cross operator.              VIDEO
 TS.. rotate                   Rotate the input image.                    VIDEO
 T... sab                      Apply shape adaptive blur.                 VIDEO
 .... scale                    Scale the input video size and/or convert  VIDEO
 .... scale2ref                Scale the input video size and/or convert  VIDEO
 ..H. scale_qsv                Quick Sync Video "scaling and format conve VIDEO
 ..H. scale_vulkan             Scale Vulkan frames                        VIDEO
 ...M scdet                    Detect video scene change                  VIDEO
 ..H. scdet_vulkan             Detect video scene change                  VIDEO
 TS.. scharr                   Apply scharr operator.                     VIDEO
 TS.. scroll                   Scroll input video.                        VIDEO
 ...M segment                  Segment video stream.                      VIDEO
 ...M select                   Select video frames to pass in output.     VIDEO
 TS.. selectivecolor           Apply CMYK adjustments to specific color r VIDEO
 ...M sendcmd                  Send commands to filters.                  VIDEO
 .... separatefields           Split input video frames into fields.      VIDEO
 ...M setdar                   Set the frame display aspect ratio.        VIDEO
 ...M setfield                 Force field for the output video frame.    VIDEO
 ...M setparams                Force field, or color property for the out VIDEO
 ...M setpts                   Set PTS for the output video frame.        VIDEO
 ...M setrange                 Force color range for the output video fra VIDEO
 ...M setsar                   Set the pixel sample aspect ratio.         VIDEO
 ...M settb                    Set timebase for the video output link.    VIDEO
 TS.. shear                    Shear transform the input image.           VIDEO
 .... showcqt                  Convert input audio to a CQT (Constant/Cla AUDIO
 .S.. showcwt                  Convert input audio to a CWT (Continuous W AUDIO
 .... showfreqs                Convert input audio to a frequencies video AUDIO
 ...M showinfo                 Show textual information for each video fr VIDEO
 .... showpalette              Display frame palette.                     VIDEO
 .S.. showspatial              Convert input audio to a spatial video out AUDIO
 .S.. showspectrum             Convert input audio to a spectrum video ou AUDIO
 .S.. showspectrumpic          Convert input audio to a spectrum video ou AUDIO
 .... showvolume               Convert input audio volume to video output AUDIO
 .... showwaves                Convert input audio to a video output.     AUDIO
 .... showwavespic             Convert input audio to a video output sing AUDIO
 T... shuffleframes            Shuffle video frames.                      VIDEO
 TS.. shufflepixels            Shuffle video pixels.                      VIDEO
 T... shuffleplanes            Shuffle video planes.                      VIDEO
 .... sidechaincompress        Sidechain compressor.                      AUDIO
 .... sidechaingate            Audio sidechain gate.                      AUDIO
 T..M sidedata                 Manipulate video frame side data.          VIDEO
 .S.. sierpinski               Render a Sierpinski fractal.               VIDEO
 .S.. signalstats              Generate statistics from video analysis.   VIDEO
 .... signature                Calculate the MPEG-7 video signature       VIDEO
 ...M silencedetect            Detect silence.                            AUDIO
 .... silenceremove            Remove silence.                            AUDIO
 .... sinc                     Generate a sinc kaiser-windowed low-pass,  AUDIO
 .... sine                     Generate sine wave audio signal.           AUDIO
 ...M siti                     Calculate spatial information (SI) and tem VIDEO
 T... smartblur                Blur the input video without impacting the VIDEO
 .... smptebars                Generate SMPTE color bars.                 VIDEO
 .... smptehdbars              Generate SMPTE HD color bars.              VIDEO
 TS.. sobel                    Apply sobel operator.                      VIDEO
 .... spectrumsynth            Convert input spectrum videos to audio out VIDEO
 .... speechnorm               Speech Normalizer.                         AUDIO
 ...M split                    Pass on the input to N video outputs.      VIDEO
 .... spp                      Apply a simple post processing filter.     VIDEO
 .S.M ssim                     Calculate the SSIM between two video strea VIDEO
 .... ssim360                  Calculate the SSIM between two 360 video s VIDEO
 .S.. stereo3d                 Convert video stereoscopic 3D view.        VIDEO
 .... stereotools              Apply various stereo tools.                AUDIO
 .... stereowiden              Apply stereo widening effect.              AUDIO
 .... streamselect             Select video streams                       UNKNOWN
 .S.. super2xsai               Scale the input by 2x using the Super2xSaI VIDEO
 .... superequalizer           Apply 18 band equalization filter.         AUDIO
 .S.. surround                 Apply audio surround upmix filter.         AUDIO
 T... swaprect                 Swap 2 rectangular objects in video.       VIDEO
 T... swapuv                   Swap U and V components.                   VIDEO
 .S.. tblend                   Blend successive frames.                   VIDEO
 .... telecine                 Apply a telecine pattern.                  VIDEO
 .... testsrc                  Generate test pattern.                     VIDEO
 .... testsrc2                 Generate another test pattern.             VIDEO
 .... thistogram               Compute and draw a temporal histogram.     VIDEO
 .S.. threshold                Threshold first video stream using other v VIDEO
 TS.. thumbnail                Select the most representative frame in a  VIDEO
 .... tile                     Tile several successive frames together.   VIDEO
 .... tiltandshift             Generate a tilt-and-shift'd video.         VIDEO
 .S.. tiltshelf                Apply a tilt shelf filter.                 AUDIO
 .... tinterlace               Perform temporal field interlacing.        VIDEO
 .S.. tlut2                    Compute and apply a lookup table from two  VIDEO
 .S.. tmedian                  Pick median pixels from successive frames. VIDEO
 .... tmidequalizer            Apply Temporal Midway Equalization.        VIDEO
 .S.. tmix                     Mix successive video frames.               VIDEO
 .S.. tonemap                  Conversion to/from different dynamic range VIDEO
 .... tpad                     Temporarily pad video frames.              VIDEO
 .S.. transpose                Transpose input video.                     VIDEO
 ..H. transpose_vulkan         Transpose Vulkan Filter                    VIDEO
 .S.. treble                   Boost or cut upper frequencies.            AUDIO
 T... tremolo                  Apply tremolo effect.                      AUDIO
 ...M trim                     Pick one continuous section from the input VIDEO
 .S.. unpremultiply            UnPreMultiply first stream with first plan VIDEO
 TS.. unsharp                  Sharpen or blur the input video.           VIDEO
 .... untile                   Untile a frame into a sequence of frames.  VIDEO
 .S.. uspp                     Apply Ultra Simple / Slow Post-processing  VIDEO
 .S.. v360                     Convert 360 projection of video.           VIDEO
 T... vaguedenoiser            Apply a Wavelet based Denoiser.            VIDEO
 .S.. varblur                  Apply Variable Blur filter.                VIDEO
 .... vectorscope              Video vectorscope.                         VIDEO
 T... vflip                    Flip the input video vertically.           VIDEO
 .... vflip_vulkan             Vertically flip the input video in Vulkan  VIDEO
 ...M vfrdet                   Variable frame rate detect filter.         VIDEO
 TS.. vibrance                 Boost or alter saturation.                 VIDEO
 T... vibrato                  Apply vibrato effect.                      AUDIO
 .S.M vif                      Calculate the VIF between two video stream VIDEO
 T... vignette                 Make or reverse a vignette effect.         VIDEO
 .... virtualbass              Audio Virtual Bass.                        AUDIO
 ...M vmafmotion               Calculate the VMAF Motion score.           VIDEO
 T... volume                   Change input volume.                       AUDIO
 ...M volumedetect             Detect audio volume.                       AUDIO
 ..H. vpp_qsv                  Quick Sync Video "VPP"                     VIDEO
 .S.. vstack                   Stack video inputs vertically.             VIDEO
 ..H. vstack_qsv               "Quick Sync Video" vstack                  VIDEO
 .S.. w3fdif                   Apply Martin Weston three field deinterlac VIDEO
 .S.. waveform                 Video waveform monitor.                    VIDEO
 .S.. weave                    Weave input video fields into frames.      VIDEO
 .S.. xbr                      Scale the input using xBR algorithm.       VIDEO
 .S.. xcorrelate               Cross-correlate first video stream with se VIDEO
 .S.. xfade                    Cross fade one video with another video.   VIDEO
 ..H. xfade_vulkan             Cross fade one video with another video.   VIDEO
 .S.. xmedian                  Pick median pixels from several video inpu VIDEO
 ...M xpsnr                    Calculate the extended perceptually weight VIDEO
 .S.. xstack                   Stack video inputs into custom layout.     VIDEO
 ..H. xstack_qsv               "Quick Sync Video" xstack                  VIDEO
 .S.. yadif                    Deinterlace the input image.               VIDEO
 TS.. yaepblur                 Yet another edge preserving blur filter.   VIDEO
 .... yuvtestsrc               Generate YUV test pattern.                 VIDEO
 .S.. zoneplate                Generate zone-plate.                       VIDEO
 .... zoompan                  Apply Zoom & Pan effect.                   VIDEO
 .S.. zscale                   Apply resizing, colorspace and bit depth c VIDEO
```

**Summary:**
- Total filters: 500
- Timeline support: 133
- Slice threading: 200
- Hardware filters: 24
- Metadata-only filters: 78

**By media type:**
- Video filters: 334
- Audio filters: 161
- Subtitle filters: 0
- Data filters: 0

**Flags:**
- T - Timeline support
- S - Slice threading
- H - Hardware device required
- M - Metadata only (does not modify frame data)

# BITSTREAM FILTERS

```
    NAME                     SUPPORTED CODECS

    aac_adtstoasc            aac
    av1_frame_merge          av1
    av1_frame_split          av1
    av1_metadata             av1
    chomp                    all
    dca_core                 dts
    dovi_rpu                 hevc, av1
    dts2pts                  h264
    eac3_core                eac3
    extract_extradata        av1, avs2, avs3, cavs, h264, hevc, mpeg1video, mpeg2video, mp...
    filter_units             av1, h264, hevc, vvc, mpeg2video, vp8, vp9
    h264_metadata            h264
    h264_mp4toannexb         h264
    h264_redundant_pps       h264
    hevc_metadata            hevc
    hevc_mp4toannexb         hevc
    mjpeg2jpeg               mjpeg
    mpeg2_metadata           mpeg2video
    mpeg4_unpack_bframes     mpeg4
    mov2textsub              all
    noise                    all
    null                     all
    opus_metadata            opus
    pcm_rechunk              pcm_alaw, pcm_f16le, pcm_f24le, pcm_f32be, pcm_f32le, pcm_f64...
    pgs_frame_merge          hdmv_pgs_subtitle
    prores_metadata          prores
    setts                    all
    showinfo                 all
    text2movsub              all
    trace_headers            av1, h264, hevc, vvc, mpeg2video, vp8, vp9
    truehd_core              truehd
    vp9_metadata             vp9
    vp9_raw_reorder          vp9
    vp9_superframe           vp9
    vp9_superframe_split     vp9
    vvc_metadata             vvc
    vvc_mp4toannexb          vvc
```

**Summary:**
- Total bitstream filters: 37
- Codec-specific filters: 30
- Generic filters: 7

## PARSERS

```
    NAME                     SUPPORTED CODECS

    aac                      aac
    aac_latm                 aac_latm
    ac3                      ac3, eac3
    av1                      av1
    dirac                    dirac
    dnxhd                    dnxhd
    dts                      dts
    ffv1                     ffv1
    flac                     flac
    gif                      gif
    h263                     h263
    h264                     h264
    hevc                     hevc
    mjpeg                    mjpeg, jpegls
    mlp                      mlp, truehd
    mp1                      mp1, mp2, mp3, mp3adu
    opus                     opus
    png                      png
    prores_raw               prores_raw
    vc1                      vc1
    vorbis                   vorbis
    vp8                      vp8
    vp9                      vp9
    vvc                      vvc
    webp                     webp
```

**Summary:**
- Total parsers: 25


## PROTOCOLS

```
IO  NAME

I.  async
I.  cache
I.  concat
I.  concatf
IO  crypto
I.  data
IO  dtls
IO  fd
IO  ffrtmpcrypt
IO  ffrtmphttp
IO  file
IO  ftp
IO  gopher
IO  gophers
I.  hls
IO  http
IO  httpproxy
IO  https
.O  icecast
I.  ipfs
I.  ipns
.O  md5
I.  mmsh
I.  mmst
IO  pipe
.O  prompeg
IO  rtmp
IO  rtmpe
IO  rtmps
IO  rtmpt
IO  rtmpte
IO  rtmpts
IO  rtp
IO  srt
IO  srtp
I.  subfile
IO  tcp
.O  tee
IO  tls
IO  udp
IO  udplite
IO  unix
```

**Summary:**
- Total protocols: 42
- Total input protocols: 38
- Total output protocols: 31
