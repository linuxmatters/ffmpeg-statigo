# introspect

Queries the compiled FFmpeg library and prints every component it contains: codecs, hardware accelerators, container formats, filters, bitstream filters, parsers, and I/O protocols. It also generates FFmpeg `./configure` flags for enabling or disabling a codec and all its dependencies.

Back: [examples/](../README.md)

## What it demonstrates

- Iterating codec descriptors with `AVCodecDescriptorNext` and matching encoders/decoders with `AVCodecFindEncoder` / `AVCodecFindDecoder`
- Walking muxers and demuxers with `AVMuxerIterate` / `AVDemuxerIterate`
- Listing filters with `AVFilterIterate` and reading pad types with `AVFilterPadGetType`
- Listing bitstream filters with `AVBSFIterate` and reading per-BSF codec ID arrays
- Listing parsers with `AVParserIterate`
- Enumerating I/O protocols with `AVIOEnumProtocols`
- Probing hardware device availability at runtime with `AVHWDeviceCtxCreate`

## Run signatures

### Full introspection

```
introspect
```

Prints: build configuration string, license, then each of the following sections with counts and flag legends:

- **CODECS**: all codecs, with `D` (decoder) and `E` (encoder) flags and media type
- **HARDWARE ACCELERATORS**: compiled-in device type names, plus HARDWARE CODECS with presence probed at runtime (`Y`/`N`)
- **FORMATS**: muxers and demuxers, with `D` and `E` flags, default codecs, and MIME type
- **FILTERS**: all filters, with `T` (timeline), `S` (slice threads), `H` (hardware), `M` (metadata-only) flags
- **BITSTREAM FILTERS**: each BSF and its supported codec list
- **PARSERS**: each parser and its supported codec list
- **PROTOCOLS**: all I/O protocols, with `I` (input) and `O` (output) flags

### Codec dependency analysis

```
introspect --enable <codec>
introspect --disable <codec>
```

Finds all codecs, encoders, decoders, parsers, demuxers, muxers, and BSFs that relate to `<codec>` and prints the corresponding `./configure` flags as a single consolidated block.

Matching is multi-strategy: exact name, variant (`h264_nvenc`), library prefix (`libx264`), prefix (`h26` finds `h264`/`h265`), and reverse lookup from parser, format, and BSF names.

## Build

From the repo root, inside `nix develop`:

```bash
just build-examples
```

The binary is written to `examples/introspect/introspect`.

> The static libraries must be present first. Run `go run ./cmd/download-lib` if you have not done so already.

## Running

```bash
# Full listing
./examples/introspect/introspect

# Generate configure flags for AV1
./examples/introspect/introspect --enable av1

# Generate disable flags for VP7
./examples/introspect/introspect --disable vp7
```

## Expected output

**Full mode** prints each section in order. Example (abbreviated):

```
ffmpeg-statigo
==============

Configuration:
--enable-gpl --enable-nonfree ...

CODECS
==================================================
 DE  aac                      AAC (Advanced Audio Coding)                [AUDIO]
 DE  av1                      Alliance for Open Media AV1                [VIDEO]
 ...
Summary:
  Total codecs: 89
  Decoders: 87 (Video: 42, Audio: 39, Subtitle: 4, Other: 2)
  Encoders: 61 (Video: 29, Audio: 26, Subtitle: 4, Other: 2)
```

**Dependency mode** prints codec descriptions as comments, then one `--enable-*` (or `--disable-*`) line per component type:

```
// Alliance for Open Media AV1
// librav1e AV1
// dav1d AV1 decoder by VideoLAN
--enable-encoder=av1_nvenc,av1_qsv,av1_vulkan,librav1e
--enable-decoder=av1,av1_qsv,libdav1d
--enable-parser=av1
--enable-demuxer=avif,obu
--enable-muxer=avif,obu
--enable-bsf=av1_frame_merge,av1_frame_split,...
```

Component lines only appear when non-empty. If no codec matches the search term, the tool exits with status 1 and prints an error to stderr.
