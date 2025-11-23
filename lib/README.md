# FFmpeg Static Libraries

This directory contains pre-built FFmpeg static libraries for different platforms.

## Automatic Download

Libraries are automatically downloaded when you build ffmpeg-statigo for the first time.

To manually trigger the download:
```bash
just download-lib
```

## Manual Download

If automatic download fails, download libraries from [GitHub Releases](https://github.com/linuxmatters/ffmpeg-statigo/releases).

Extract the appropriate tarball to the lib directory:
- `ffmpeg-linux-amd64.tar.gz` → `lib/linux_amd64/`
- `ffmpeg-darwin-amd64.tar.gz` → `lib/darwin_amd64/`
- `ffmpeg-darwin-arm64.tar.gz` → `lib/darwin_arm64/`
- `ffmpeg-linux-arm64.tar.gz` → `lib/linux_arm64/`

## Structure

```
lib/
├── <platform>_<arch>/
│   └── libffmpeg.a
├── fetch.go         # Auto-download implementation
└── README.md        # This file
```
