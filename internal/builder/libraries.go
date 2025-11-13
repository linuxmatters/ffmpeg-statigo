package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// All libraries in build order
var AllLibraries = []*Library{
	// Compression libraries
	zlib,

	// Image libraries
	png,

	// Font libraries & rendering
	expat,
	iconv,
	fribidi,
	unibreak,
	freetype,
	fontconfig,
	harfbuzz,
	libass,

	// Hardware acceleration headers
	nvcodec,
	vulkanHeaders,
	libvpl,

	// Audio codecs
	lame,
	opus,
	ogg,
	vorbis,

	// Video codecs
	theora,
	vpx,
	x264,
	x265,
	dav1d,
	rav1e,

	// FFmpeg itself (must be last)
	ffmpeg,
}

// iconv - character encoding conversion (macOS only)
var iconv = &Library{
	Name:        "iconv",
	URL:         "https://ftp.gnu.org/pub/gnu/libiconv/libiconv-1.18.tar.gz",
	Platform:    []string{"darwin"},
	BuildSystem: &AutoconfBuild{},
	ConfigureArgs: func(os string) []string {
		return []string{
			"--disable-dependency-tracking",
			"--disable-debug",
			"--enable-extra-encodings",
			"--enable-static",
		}
	},
	LinkLibs: []string{"libiconv"},
}

// expat - XML parser (Linux only, needed for fontconfig)
var expat = &Library{
	Name:        "expat",
	URL:         "https://github.com/libexpat/libexpat/releases/download/R_2_7_3/expat-2.7.3.tar.gz",
	Platform:    []string{"linux"},
	BuildSystem: &AutoconfBuild{},
	ConfigureArgs: func(os string) []string {
		return []string{
			"--enable-static",
			"--disable-shared",
			"--without-xmlwf",
		}
	},
	LinkLibs: []string{"libexpat"},
}

// zlib - compression library
var zlib = &Library{
	Name:          "zlib",
	URL:           "https://github.com/madler/zlib/releases/download/v1.3.1/zlib-1.3.1.tar.gz",
	BuildSystem:   &AutoconfBuild{},
	SkipAutoFlags: true, // zlib has a custom configure script that rejects CFLAGS/LDFLAGS
	ConfigureArgs: func(os string) []string {
		return []string{
			"--static",
		}
	},
	LinkLibs: []string{"libz"},
}

// png - PNG image library
var png = &Library{
	Name:        "png",
	URL:         "https://github.com/pnggroup/libpng/archive/refs/tags/v1.6.50.tar.gz",
	BuildSystem: &AutoconfBuild{},
	ConfigureArgs: func(os string) []string {
		return []string{
			"--disable-dependency-tracking",
			"--disable-silent-rules",
			"--disable-shared",
			"--enable-static",
		}
	},
	LinkLibs: []string{"libpng16"},
}

// fribidi - Unicode bidirectional algorithm library (needed by libass)
var fribidi = &Library{
	Name:        "fribidi",
	URL:         "https://github.com/fribidi/fribidi/releases/download/v1.0.16/fribidi-1.0.16.tar.xz",
	BuildSystem: &AutoconfBuild{},
	ConfigureArgs: func(os string) []string {
		return []string{
			"--disable-dependency-tracking",
			"--disable-debug",
			"--disable-silent-rules",
			"--enable-static",
		}
	},
	LinkLibs: []string{"libfribidi"},
}

// unibreak - line breaking library (needed by libass)
var unibreak = &Library{
	Name:        "unibreak",
	URL:         "https://github.com/adah1972/libunibreak/releases/download/libunibreak_6_1/libunibreak-6.1.tar.gz",
	BuildSystem: &AutoconfBuild{},
	ConfigureArgs: func(os string) []string {
		return []string{
			"--enable-static",
			"--disable-shared",
		}
	},
	LinkLibs: []string{"libunibreak"},
}

// fontconfig - font configuration library (Linux only)
var fontconfig = &Library{
	Name:        "fontconfig",
	URL:         "https://www.freedesktop.org/software/fontconfig/release/fontconfig-2.16.0.tar.xz",
	Platform:    []string{"linux"},
	BuildSystem: &AutoconfBuild{},
	ConfigureArgs: func(os string) []string {
		return []string{
			"--enable-static",
			"--disable-shared",
		}
	},
	LinkLibs: []string{"libfontconfig"},
}

// freetype - font rendering library
var freetype = &Library{
	Name:        "freetype",
	URL:         "https://download.savannah.gnu.org/releases/freetype/freetype-2.14.1.tar.xz",
	BuildSystem: &AutoconfBuild{},
	ConfigureArgs: func(os string) []string {
		return []string{
			"--enable-static",
			"--disable-shared",
			"--without-brotli",
			"--without-bzip2",
			"--without-harfbuzz",
		}
	},
	LinkLibs: []string{"libfreetype"},
}

// harfbuzz - text shaping library (needed by libass)
var harfbuzz = &Library{
	Name:        "harfbuzz",
	URL:         "https://github.com/harfbuzz/harfbuzz/releases/download/12.2.0/harfbuzz-12.2.0.tar.xz",
	BuildSystem: &MesonBuild{},
	ConfigureArgs: func(os string) []string {
		return []string{
			"--buildtype=release",
			"--default-library=static",
			"-Dcairo=disabled",
			"-Dcoretext=enabled",
			"-Dfreetype=enabled",
			"-Dintrospection=disabled",
			"-Dtests=disabled",
		}
	},
	LinkLibs: []string{"libharfbuzz"},
}

// libass - subtitle rendering library
var libass = &Library{
	Name:        "ass",
	URL:         "https://github.com/libass/libass/releases/download/0.17.4/libass-0.17.4.tar.gz",
	BuildSystem: &AutoconfBuild{},
	ConfigureArgs: func(os string) []string {
		args := []string{
			"--disable-shared",
		}

		// libass uses coretext on macOS, fontconfig on Linux
		if os == "darwin" {
			args = append(args, "--disable-fontconfig")
		}

		return args
	},
	PostExtract: func(srcPath string) error {
		// Prevent automake regeneration
		return touchAutomakeFiles(srcPath)
	},
	LinkLibs: []string{"libass"},
}

// nvcodec - NVIDIA codec SDK headers (Linux only)
var nvcodec = &Library{
	Name:     "nvcodec",
	URL:      "https://github.com/FFmpeg/nv-codec-headers/releases/download/n12.2.72.0/nv-codec-headers-12.2.72.0.tar.gz",
	Platform: []string{"linux"},
	BuildSystem: &MakefileBuild{
		Targets: nil, // No build targets, just install
		InstallFunc: func(srcPath, installDir string) error {
			return runCommand(srcPath, os.Stdout, installDir, "make", fmt.Sprintf("PREFIX=%s", installDir), "install")
		},
	},
	LinkLibs: nil, // Headers only
}

// vulkanHeaders - Vulkan API headers (cross-platform)
var vulkanHeaders = &Library{
	Name:        "vulkan",
	URL:         "https://github.com/KhronosGroup/Vulkan-Headers/archive/refs/tags/v1.4.332.tar.gz",
	BuildSystem: &CMakeBuild{},
	ConfigureArgs: func(os string) []string {
		return []string{
			"-DVULKAN_HEADERS_ENABLE_MODULE=OFF",
		}
	},
	LinkLibs: nil, // Headers only
}

// libvpl - Intel VPL/oneVPL headers (Linux only, for QuickSync)
var libvpl = &Library{
	Name:        "vpl",
	URL:         "https://github.com/intel/libvpl/archive/refs/tags/v2.15.0.tar.gz",
	Platform:    []string{"linux"},
	BuildSystem: &CMakeBuild{},
	ConfigureArgs: func(os string) []string {
		return []string{
			"-DBUILD_SHARED_LIBS=OFF",
			"-DBUILD_TOOLS=OFF",
			"-DBUILD_TESTS=OFF",
			"-DINSTALL_EXAMPLE_CODE=OFF",
		}
	},
	PostExtract: func(srcPath string) error {
		// Patch vpl.pc.in to add -lstdc++ for C++ static library linking
		vplPcIn := filepath.Join(srcPath, "libvpl", "pkgconfig", "vpl.pc.in")
		content, err := os.ReadFile(vplPcIn)
		if err != nil {
			return fmt.Errorf("failed to read vpl.pc.in: %w", err)
		}

		// Add -lstdc++ after -l@OUTPUT_NAME@ since libvpl is C++ code
		patched := strings.ReplaceAll(string(content), "-l@OUTPUT_NAME@ @VPL_PKGCONFIG_DEPENDENT_LIBS@", "-l@OUTPUT_NAME@ -lstdc++ @VPL_PKGCONFIG_DEPENDENT_LIBS@")

		if err := os.WriteFile(vplPcIn, []byte(patched), 0644); err != nil {
			return fmt.Errorf("failed to write patched vpl.pc.in: %w", err)
		}

		return nil
	},
	LinkLibs: []string{"libvpl"},
}

// lame - MP3 encoder
var lame = &Library{
	Name:        "lame",
	URL:         "https://downloads.sourceforge.net/project/lame/lame/3.100/lame-3.100.tar.gz",
	BuildSystem: &AutoconfBuild{},
	ConfigureArgs: func(os string) []string {
		return []string{
			"--disable-debug",
			"--enable-static",
			"--disable-shared",
		}
	},
	LinkLibs: []string{"libmp3lame"},
}

// opus - Opus audio codec
var opus = &Library{
	Name:        "opus",
	URL:         "https://downloads.xiph.org/releases/opus/opus-1.5.2.tar.gz",
	BuildSystem: &AutoconfBuild{},
	ConfigureArgs: func(os string) []string {
		return []string{
			"--disable-debug",
			"--disable-doc",
			"--enable-static",
			"--disable-shared",
		}
	},
	LinkLibs: []string{"libopus"},
}

// ogg - Ogg container format
var ogg = &Library{
	Name:        "ogg",
	URL:         "https://downloads.xiph.org/releases/ogg/libogg-1.3.6.tar.xz",
	BuildSystem: &AutoconfBuild{},
	ConfigureArgs: func(os string) []string {
		return []string{
			"--enable-static",
			"--disable-shared",
		}
	},
	LinkLibs: []string{"libogg"},
}

// vorbis - Vorbis audio codec
var vorbis = &Library{
	Name:        "vorbis",
	URL:         "https://downloads.xiph.org/releases/vorbis/libvorbis-1.3.7.tar.xz",
	BuildSystem: &AutoconfBuild{},
	ConfigureArgs: func(os string) []string {
		return []string{
			"--enable-static",
			"--disable-shared",
			"--disable-examples",
		}
	},
	LinkLibs: []string{"libvorbis"},
}

// theora - Theora video codec
var theora = &Library{
	Name:        "theora",
	URL:         "https://downloads.xiph.org/releases/theora/libtheora-1.2.0.tar.xz",
	BuildSystem: &AutoconfBuild{},
	ConfigureArgs: func(os string) []string {
		return []string{
			"--enable-static",
			"--disable-shared",
			"--disable-oggtest",
			"--disable-vorbistest",
			"--disable-examples",
		}
	},
	PostExtract: func(srcPath string) error {
		// Prevent automake regeneration
		return touchAutomakeFiles(srcPath)
	},
	LinkLibs: []string{"libtheora"},
}

// vpx - VP8/VP9 video codec
var vpx = &Library{
	Name:          "vpx",
	URL:           "https://github.com/webmproject/libvpx/archive/refs/tags/v1.15.2.tar.gz",
	BuildSystem:   &AutoconfBuild{},
	SkipAutoFlags: true, // vpx has a custom configure script that rejects CFLAGS/LDFLAGS
	ConfigureArgs: func(os string) []string {
		return []string{
			"--enable-static",
			"--disable-shared",
			"--disable-examples",
			"--enable-vp9-highbitdepth",
			"--disable-unit-tests",
			"--as=yasm",
		}
	},
	LinkLibs: []string{"libvpx"},
}

// x264 - H.264/AVC video encoder
var x264 = &Library{
	Name:          "x264",
	URL:           "https://code.videolan.org/videolan/x264/-/archive/master/x264-master.tar.bz2",
	BuildSystem:   &AutoconfBuild{},
	SkipAutoFlags: true, // x264 has a custom configure script that rejects CFLAGS/LDFLAGS
	ConfigureArgs: func(os string) []string {
		return []string{
			"--disable-avs",
			"--disable-cli",
			"--disable-ffms",
			"--disable-gpac",
			"--disable-lavf",
			"--disable-lsmash",
			"--disable-swscale",
			"--enable-static",
			"--enable-strip",
		}
	},
	PostExtract: func(srcPath string) error {
		// x264 needs to find nasm explicitly on x86/x86_64
		// ARM architectures use the C compiler as assembler instead
		if runtime.GOARCH == "amd64" || runtime.GOARCH == "386" {
			os.Setenv("AS", "nasm")
		}
		return nil
	},
	LinkLibs: []string{"libx264"},
}

// x265 - H.265/HEVC video encoder
var x265 = &Library{
	Name: "x265",
	URL:  "https://bitbucket.org/multicoreware/x265_git/get/ffba52bab55dce9b1b3a97dd08d12e70297e2180.tar.bz2",
	BuildSystem: &CMakeBuild{
		SourceSubdir: "source", // x265 source is in source/ subdirectory
	},
	ConfigureArgs: func(os string) []string {
		return []string{
			"-DENABLE_SHARED=OFF",
			"-DENABLE_CLI=OFF",
			"-DENABLE_AGGRESSIVE_CHECKS=OFF",
		}
	},
	LinkLibs: []string{"libx265"},
}

// dav1d - AV1 video decoder
var dav1d = &Library{
	Name:        "dav1d",
	URL:         "https://code.videolan.org/videolan/dav1d/-/archive/1.5.2/dav1d-1.5.2.tar.bz2",
	BuildSystem: &MesonBuild{},
	ConfigureArgs: func(os string) []string {
		return []string{
			"--default-library=static",
			"--buildtype=release",
			"-Denable_tools=false",
			"-Denable_tests=false",
		}
	},
	LinkLibs: []string{"libdav1d"},
}

// rav1e - AV1 video encoder
var rav1e = &Library{
	Name: "rav1e",
	URL:  "https://github.com/xiph/rav1e/archive/refs/tags/v0.8.1.tar.gz",
	BuildSystem: &CargoBuild{
		InstallFunc: func(srcPath, installDir string) error {
			// Set RUSTFLAGS for native CPU optimization
			os.Setenv("RUSTFLAGS", "-C target-cpu=native")

			// cargo cinstall for C library installation
			return runCommand(srcPath, os.Stdout, installDir, "cargo", "cinstall",
				fmt.Sprintf("--prefix=%s", installDir),
				"--libdir=lib",
				"--library-type=staticlib",
				"--crt-static",
				"--release")
		},
	},
	LinkLibs: []string{"librav1e"},
}

// ffmpeg - FFmpeg multimedia framework
var ffmpeg = &Library{
	Name:          "ffmpeg",
	URL:           "https://github.com/FFmpeg/FFmpeg/archive/refs/tags/n8.0.tar.gz",
	BuildSystem:   &AutoconfBuild{},
	SkipAutoFlags: true, // FFmpeg uses --extra-cflags and --extra-ldflags instead
	ConfigureArgs: func(os string) []string {
		// FFmpeg needs explicit paths to headers and libraries
		stagingDir, _ := filepath.Abs(".build/staging")
		incDir := filepath.Join(stagingDir, "include")
		libDir := filepath.Join(stagingDir, "lib")

		args := []string{
			"--pkg-config-flags=--static",
			fmt.Sprintf("--extra-cflags=-I%s", incDir),
			fmt.Sprintf("--extra-ldflags=-L%s", libDir),
			"--disable-autodetect",
			"--disable-debug",
			"--disable-doc",
			"--disable-htmlpages",
			"--disable-manpages",
			"--disable-podpages",
			"--disable-programs",
			"--disable-txtpages",
			// Disable everything, then selectively enable
			"--disable-everything",
			// Build options
			"--enable-pic",
			"--enable-gpl",
			"--enable-version3",
			"--enable-static",
			// Enable all protocols and filters
			"--enable-protocols",
			"--enable-filters",
			// TEMPORARILY Enable all muxers, demuxers, parsers, bsfs, encoders, decoders
			"--enable-muxers",
			"--enable-demuxers",
			"--enable-parsers",
			"--enable-bsfs",
			"--enable-encoders",
			"--enable-decoders",
			// Enable external libraries
			"--enable-libass",
			"--enable-libdav1d",
			"--enable-libfreetype",
			"--enable-libfribidi",
			"--enable-libharfbuzz",
			"--enable-libmp3lame",
			"--enable-libopus",
			"--enable-librav1e",
			"--enable-libtheora",
			"--enable-libvpx",
			"--enable-vulkan",
			"--enable-libx264",
			"--enable-libx265",
			"--enable-zlib",
		}

		// Platform-specific options
		if os == "linux" {
			args = append(args,
				"--enable-libfontconfig",
				"--enable-ffnvcodec",
				"--enable-nvdec",
				"--enable-nvenc",
				"--enable-libvpl",
			)
		} else if os == "darwin" {
			args = append(args,
				"--enable-avfoundation",
				"--enable-audiotoolbox",
				"--enable-videotoolbox",
			)
		}

		return args
	},
	LinkLibs: []string{
		"libavcodec",
		"libavdevice",
		"libavfilter",
		"libavformat",
		"libavutil",
		"libswresample",
		"libswscale",
	},
}

// touchAutomakeFiles touches all automake-related files to prevent regeneration
func touchAutomakeFiles(srcPath string) error {
	now := time.Now()

	// Touch top-level files
	files := []string{
		"Makefile.in",
		"aclocal.m4",
		"config.h.in",
		"configure",
	}

	for _, file := range files {
		fullPath := filepath.Join(srcPath, file)
		if _, err := os.Stat(fullPath); err == nil {
			// File exists, update its timestamp
			if err := os.Chtimes(fullPath, now, now); err != nil {
				// Log warning but continue
				fmt.Fprintf(os.Stderr, "Warning: failed to touch %s: %v\n", file, err)
			}
		}
	}

	// Also touch any Makefile.in files in subdirectories
	filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == "Makefile.in" {
			os.Chtimes(path, now, now)
		}
		return nil
	})

	return nil
}
