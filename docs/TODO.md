# TODO

- [x] Refactor the internal build tool
- [x] Reorganise project structure
- [x] FFmpeg build argument resolver.to enable/disable codecs
- [x] FFmpeg feature enablement
  --enable-libvmaf         enable vmaf filter via libvmaf [no]
  --enable-whisper         enable whisper filter [no]
- [x] Rename to ffmpeg-statigo
- [x] How to embed/distribute the static FFmpeg library?
- [x] Review default codecs:
  - https://ffmpeg.martin-riedl.de/
  - https://github.com/markus-perl/ffmpeg-build-script/blob/master/build-ffmpeg
- [x] Create some test cases that exercise some of the FFmpeg API surface

## From the original author

- [x] Expose more headers.
- [ ] Expose platform specific headers.
- [ ] Cleanup internal packages.

---

# Onboard

This is a hard fork of the [ffmpeg-go](https://github.com/csnewman/ffmpeg-go) project, now called **ffmpeg-statigo**.

ffmpeg-statigo has been updated to Go 1.24, all dependencies uplifted to current versions, GitHub CI is fixed so the ffmpeg static libraries are built for all supported architectures and the generator has been updated to support current `clang` on Linux. Mostly critically this project has been updated from FFmpeg 6.1 to FFmpeg 8.0 including the addition os `x265`, `dav1d`, `rav1e` and hardware acceleration support for NVENC, QuickSync and Vulkan to complement the exist VideoToolbox support.

The git history has been purged all tags and historical commits of the static ffmpeg libraries as we are preparing this project to be launched under it's new name in a different GitHUb organisation. This brought the git repo down in size from 1.9GB to 9MB.

NixOS is the host development workstation. There is a `flake.nix` in the project to enable the required software in a Nix devevelopment shell, which is automatically activated via `direnv`. The `justfile` is used for all build and test commands. `fish` is the default shell.

Build the project using `just build`. Using `just build` is the only valid way to build. Never build the project using adhoc commands. Never rediect builds elsewhere, the build system does the right thing for you. Work with the provide build tool, not against it. The following directory structure is populated with the build assets and log files for debugging.

```
.build
├── build
│   ├── ass
│   │   └── build.log
│   ├── dav1d
│   │   └── build.log
│   ├── ffmpeg
│   │   └── build.log
│   ├── fontconfig
│   │   └── build.log
│   ├── freetype
│   │   └── build.log
│   ├── fribidi
│   │   └── build.log
│   ├── glslang
│   │   └── build.log
│   ├── harfbuzz
│   │   └── build.log
│   ├── iconv
│   │   └── build.log
│   ├── lame
│   │   └── build.log
│   ├── nvcodec
│   │   └── build.log
│   ├── ogg
│   │   └── build.log
│   ├── opus
│   │   └── build.log
│   ├── png
│   │   └── build.log
│   ├── rav1e
│   │   └── build.log
│   ├── theora
│   │   └── build.log
│   ├── unibreak
│   │   └── build.log
│   ├── vorbis
│   │   └── build.log
│   ├── vpl
│   │   └── build.log
│   ├── vpx
│   │   └── build.log
│   ├── vulkan
│   │   └── build.log
│   ├── x264
│   │   └── build.log
│   ├── x265
│   │   └── build.log
│   ├── xml2
│   │   └── build.log
│   └── zlib
│       └── build.log
├── downloads
│   └── <all tarball are here>
├── src
│   ├── ass
│   ├── dav1d
│   ├── ffmpeg
│   ├── fontconfig
│   ├── freetype
│   ├── fribidi
│   ├── glslang
│   ├── harfbuzz
│   ├── lame
│   ├── nvcodec
│   ├── ogg
│   ├── opus
│   ├── png
│   ├── rav1e
│   ├── theora
│   ├── unibreak
│   ├── vorbis
│   ├── vpl
│   ├── vpx
│   ├── vulkan
│   ├── x264
│   ├── x265
│   ├── xml2
│   └── zlib
└── staging
    ├── bin
    ├── etc
    ├── include
    ├── lib
    ├── share
    └── var
```

Read the `README.md`, `TODO.md` and analyse the code to get a full understanding of the project. Let me know when you are ready to collaborate.
