package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// All libraries in build order
var AllLibraries = []*Library{
	// Compression
	zlib,

	// XML parsing and formatting
	libiconv,
	libxml2,

	// Hardware acceleration
	libvpl,
	nvcodecheaders,
	vulkanheaders,
	glslang,

	// Image processing
	zimg,

	// Audio codecs
	lame,
	opus,

	// Video codecs
	dav1d,
	libvpx,
	rav1e,
	x264,
	x265,

	// TLS/SSL
	openssl,

	// Streaming protocols
	libsrt,

	// FFmpeg (must be last)
	ffmpeg,
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

// libiconv - character encoding conversion (macOS only)
var libiconv = &Library{
	Name:        "libiconv",
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

// libxml2 - XML parsing library
var libxml2 = &Library{
	Name:        "libxml2",
	URL:         "https://download.gnome.org/sources/libxml2/2.15/libxml2-2.15.1.tar.xz",
	BuildSystem: &AutoconfBuild{},
	ConfigureArgs: func(os string) []string {
		return []string{
			"--enable-static",
			"--disable-shared",
			"--with-zlib",
			"--without-python",
			"--without-catalog",  // Don't need XML catalog resolution
			"--without-debug",    // Disable debug module
			"--without-modules",  // Don't need dynamic module loading
			"--without-sax1",     // Don't need legacy SAX1 interface
			"--without-xinclude", // Don't need XInclude processing
			"--without-xptr",     // Don't need XPointer support
		}
	},
	LinkLibs: []string{"libxml2"},
}

// nvcodecheaders - NVIDIA codec SDK headers (Linux only)
var nvcodecheaders = &Library{
	Name:     "nv-codec-headers",
	URL:      "https://github.com/FFmpeg/nv-codec-headers/releases/download/n11.1.5.3/nv-codec-headers-11.1.5.3.tar.gz",
	Platform: []string{"linux"},
	BuildSystem: &MakefileBuild{
		Targets: nil, // No build targets, just install
		InstallFunc: func(srcPath, installDir string) error {
			return runCommand(srcPath, os.Stdout, installDir, "make", fmt.Sprintf("PREFIX=%s", installDir), "install")
		},
	},
	LinkLibs: nil, // Headers only
}

// vulkanheaders - Vulkan API headers (cross-platform)
var vulkanheaders = &Library{
	Name:        "Vulkan-Headers",
	URL:         "https://github.com/KhronosGroup/Vulkan-Headers/archive/refs/tags/v1.4.332.tar.gz",
	BuildSystem: &CMakeBuild{},
	ConfigureArgs: func(os string) []string {
		return []string{
			"-DVULKAN_HEADERS_ENABLE_MODULE=OFF",
		}
	},
	LinkLibs: nil, // Headers only
}

// glslang - Khronos GLSL/SPIR-V shader compiler (required for Vulkan encoders/decoders/filters)
// NOTE: Pinned to 15.4.0 because FFmpeg 8.0 requires libSPVRemapper which was removed in glslang 16.0.0
// (functionality moved to SPIRV-Tools). Upgrade to 16.x requires FFmpeg to update their spirv_compiler detection.
var glslang = &Library{
	Name:        "glslang",
	URL:         "https://github.com/KhronosGroup/glslang/archive/refs/tags/15.4.0.tar.gz",
	BuildSystem: &CMakeBuild{},
	PostExtract: func(srcPath string) error {
		// Run update_glslang_sources.py to fetch external dependencies
		pythonScript := filepath.Join(srcPath, "update_glslang_sources.py")
		if _, err := os.Stat(pythonScript); err == nil {
			fmt.Fprintf(os.Stderr, "Running update_glslang_sources.py to fetch glslang dependencies...\n")
			cmd := exec.Command("python3", pythonScript)
			cmd.Dir = srcPath
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to run update_glslang_sources.py: %w", err)
			}
		}
		return nil
	},
	ConfigureArgs: func(os string) []string {
		return []string{
			"-DBUILD_SHARED_LIBS=OFF",
			"-DENABLE_GLSLANG_BINARIES=OFF", // Don't build CLI tools
			"-DENABLE_HLSL=OFF",             // Don't need DirectX HLSL support for Vulkan
			"-DGLSLANG_TESTS=OFF",           // Don't build tests
			"-DSPIRV_SKIP_EXECUTABLES=ON",
			"-DSPIRV_SKIP_TESTS=ON",
		}
	},
	LinkLibs: []string{
		"libglslang",
		"libGenericCodeGen",
		"libMachineIndependent",
		"libOSDependent",
		"libSPIRV",
		"libSPVRemapper",
		"libSPIRV-Tools",
		"libSPIRV-Tools-opt",
	},
}

// libvpl - Intel VPL/oneVPL headers (Linux only, for QuickSync)
var libvpl = &Library{
	Name:        "libvpl",
	URL:         "https://github.com/intel/libvpl/archive/refs/tags/v2.15.0.tar.gz",
	Platform:    []string{"linux"},
	BuildSystem: &CMakeBuild{},
	ConfigureArgs: func(os string) []string {
		return []string{
			"-DBUILD_SHARED_LIBS=OFF",
			"-DBUILD_TESTS=OFF",
			"-DBUILD_EXPERIMENTAL=OFF",
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

// zimg - High-quality image scaling and colorspace conversion library
var zimg = &Library{
	Name:        "zimg",
	URL:         "https://github.com/sekrit-twc/zimg/archive/refs/tags/release-3.0.6.tar.gz",
	BuildSystem: &AutoconfBuild{},
	PostExtract: func(srcPath string) error {
		// Run autogen.sh to generate configure script
		autogenScript := filepath.Join(srcPath, "autogen.sh")
		if _, err := os.Stat(autogenScript); err == nil {
			fmt.Fprintf(os.Stderr, "Running autogen.sh to generate configure script...\n")
			cmd := exec.Command("sh", autogenScript)
			cmd.Dir = srcPath
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to run autogen.sh: %w", err)
			}
		}

		// Patch zimg.pc.in to add -lm to Libs.private (zimg uses math functions like log10f)
		zimgPcIn := filepath.Join(srcPath, "zimg.pc.in")
		content, err := os.ReadFile(zimgPcIn)
		if err != nil {
			return fmt.Errorf("failed to read zimg.pc.in: %w", err)
		}

		// Add -lm after @STL_LIBS@ in Libs.private
		patched := strings.ReplaceAll(string(content), "Libs.private: @STL_LIBS@", "Libs.private: @STL_LIBS@ -lm")

		if err := os.WriteFile(zimgPcIn, []byte(patched), 0644); err != nil {
			return fmt.Errorf("failed to write patched zimg.pc.in: %w", err)
		}

		fmt.Fprintf(os.Stderr, "Patched zimg.pc.in to include -lm for math library\n")
		return nil
	},
	ConfigureArgs: func(os string) []string {
		return []string{
			"--disable-shared",
			"--enable-static",
			"--disable-testapp",   // Don't build test application
			"--disable-example",   // Don't build example programs
			"--disable-unit-test", // Don't build unit tests
		}
	},
	LinkLibs: []string{"libzimg"},
}

// lame - MP3 encoder
var lame = &Library{
	Name:        "lame",
	URL:         "https://downloads.sourceforge.net/project/lame/lame/3.100/lame-3.100.tar.gz",
	BuildSystem: &AutoconfBuild{},
	ConfigureArgs: func(os string) []string {
		return []string{
			"--disable-debug",
			"--disable-decoder",        // Only need MP3 encoder
			"--disable-frontend",       // Don't build lame CLI tool
			"--disable-analyzer-hooks", // Exclude debugging hooks
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
			"--disable-doc",
			"--disable-extra-programs",
			"--disable-shared",
			"--enable-static",
		}
	},
	LinkLibs: []string{"libopus"},
}

// libvpx - VP8/VP9 video codec
var libvpx = &Library{
	Name:          "libvpx",
	URL:           "https://github.com/webmproject/libvpx/archive/refs/tags/v1.15.2.tar.gz",
	BuildSystem:   &AutoconfBuild{},
	SkipAutoFlags: true, // vpx has a custom configure script that rejects CFLAGS/LDFLAGS
	ConfigureArgs: func(os string) []string {
		return []string{
			"--as=yasm",
			"--disable-docs",
			"--disable-examples",
			"--disable-install-bins",
			"--disable-libyuv",   // FFmpeg handles color conversion
			"--disable-postproc", // Decoder visual enhancement - FFmpeg doesn't use
			"--disable-shared",
			"--disable-tools", // Don't build vpxenc/vpxdec
			"--disable-unit-tests",
			"--disable-vp9-postproc", // VP9 decoder postprocessing - FFmpeg doesn't use
			"--enable-static",
			"--enable-vp9-highbitdepth",
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

// x265 - H.265/HEVC video encoder 7.9M
var x265 = &Library{
	Name: "x265",
	URL:  "https://bitbucket.org/multicoreware/x265_git/get/ffba52bab55dce9b1b3a97dd08d12e70297e2180.tar.bz2",
	BuildSystem: &CMakeBuild{
		SourceSubdir: "source", // x265 source is in source/ subdirectory
	},
	ConfigureArgs: func(os string) []string {
		return []string{
			"-DENABLE_AGGRESSIVE_CHECKS=OFF",
			"-DENABLE_CLI=OFF",
			"-DENABLE_SHARED=OFF",
			"-DENABLE_TESTS=OFF",      // Don't build test suite
			"-DLOGGED_PRIMITIVES=OFF", // Reduce logging overhead
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
			os.Setenv("CARGO_PROFILE_RELEASE_DEBUG", "false")

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

// openssl - TLS/SSL library for HTTPS, RTMPS, SRT, RIST protocols
// Keep these essential features (explicitly enabled):
// - TLS 1.2, TLS 1.3, DTLS 1.2
// - EC (for secp256r1 key generation)
// - DH (for RTMP Diffie-Hellman)
// - AES (for SRTP and TLS ciphers)
// - SHA-256, SHA-384, SHA-512
// - X509 certificates
// - BIO, EVP APIs
var openssl = &Library{
	Name:        "openssl",
	URL:         "https://github.com/openssl/openssl/releases/download/openssl-3.6.0/openssl-3.6.0.tar.gz",
	BuildSystem: &OpenSSLBuild{},
	ConfigureArgs: func(os string) []string {
		return []string{
			"no-shared",
			"no-apps",
			"no-comp",       // Disable compression (CRIME vulnerability, rarely used)
			"no-ct",         // No Certificate Transparency
			"no-deprecated", // Exclude deprecated APIs (FFmpeg patched for OpenSSL 3.6 compatibility)
			"no-docs",
			"no-engine",       // Disable engine support (legacy plugin system)
			"no-err",          // Don't compile error strings (~200-300KB savings)
			"no-nextprotoneg", // Deprecated NPN (use ALPN)
			"no-tests",
			"no-ts",         // No Time Stamping Authority
			"no-ui-console", // No interactive console prompts

			// Disable old/insecure SSL/TLS versions (keep TLS 1.2, 1.3, DTLS 1.2)
			"no-ssl3",
			"no-tls1",
			"no-tls1_1",
			"no-dtls1",

			// Disable weak/obsolete/uncommon algorithms
			"no-bf",        // Blowfish (old)
			"no-blake2",    // BLAKE2 (uncommon hash)
			"no-camellia",  // Uncommon cipher
			"no-cast",      // Ancient cipher
			"no-dsa",       // FFmpeg uses EC, not DSA
			"no-ec2m",      // Binary field ECC (rarely used)
			"no-idea",      // Ancient cipher
			"no-md2",       // MD2 (broken)
			"no-md4",       // MD4 (broken)
			"no-mdc2",      // Ancient hash
			"no-ocb",       // OCB mode (uncommon)
			"no-rc2",       // Ancient cipher
			"no-rc4",       // Broken cipher
			"no-rc5",       // Ancient cipher
			"no-rmd160",    // RIPEMD-160 (uncommon)
			"no-seed",      // Uncommon cipher
			"no-siv",       // SIV mode (uncommon)
			"no-sm2",       // Chinese SM2 (algorithm)
			"no-sm3",       // Chinese SM3 (hash)
			"no-sm4",       // Chinese SM4 (cipher)
			"no-whirlpool", // Uncommon hash
			"no-gost",      // Russian crypto standard (not needed for streaming)

			// Disable advanced PKI features not used by FFmpeg
			"no-cms",  // Cryptographic Message Syntax (email)
			"no-ocsp", // Online Certificate Status Protocol

			// Disable authentication schemes not needed for HTTPS/RTMPS
			"no-srp", // Secure Remote Password
			"no-psk", // Pre-shared keys
		}
	},
	LinkLibs: []string{"libssl", "libcrypto"},
}

// libsrt - Secure Reliable Transport (SRT) protocol library
var libsrt = &Library{
	Name:        "libsrt",
	URL:         "https://github.com/Haivision/srt/archive/refs/tags/v1.5.5-rc.0a.tar.gz",
	BuildSystem: &CMakeBuild{},
	ConfigureArgs: func(os string) []string {
		return []string{
			"-DENABLE_SHARED=OFF",              // Static library only
			"-DENABLE_STATIC=ON",               // Build static library
			"-DENABLE_APPS=OFF",                // Skip srt-live-transmit and CLI tools (saves ~1MB and build time)
			"-DENABLE_BONDING=OFF",             // Advanced bonding feature not used by FFmpeg
			"-DENABLE_TESTING=OFF",             // No test applications
			"-DENABLE_UNITTESTS=OFF",           // No unit tests
			"-DENABLE_LOGGING=OFF",             // Minimal logging (FFmpeg has its own)
			"-DENABLE_HEAVY_LOGGING=OFF",       // Minimal logging (FFmpeg has its own)
			"-DUSE_STATIC_LIBSTDCXX=OFF",       // Static C++ stdlib linking
			"-DSRT_USE_OPENSSL_STATIC_LIBS=ON", // Link OpenSSL statically
			"-DUSE_OPENSSL_PC=ON",              // Use pkg-config to find OpenSSL
		}
	},
	LinkLibs:     []string{"libsrt"},
	Dependencies: []*Library{openssl}, // SRT requires OpenSSL for encryption
}

// ffmpeg - FFmpeg multimedia framework
var ffmpeg = &Library{
	Name:          "ffmpeg",
	URL:           "https://github.com/FFmpeg/FFmpeg/archive/refs/tags/n8.0.tar.gz",
	BuildSystem:   &AutoconfBuild{},
	SkipAutoFlags: true, // FFmpeg uses --extra-cflags and --extra-ldflags instead
	PostExtract: func(srcPath string) error {
		// Apply OpenSSL 3.6 compatibility patch
		patchFile := filepath.Join(srcPath, "libavformat", "tls_openssl.c")

		// Read the file
		content, err := os.ReadFile(patchFile)
		if err != nil {
			return fmt.Errorf("failed to read tls_openssl.c: %w", err)
		}

		// Apply the replacements
		patched := strings.ReplaceAll(string(content),
			"X509_gmtime_adj(X509_get_notBefore(*cert)",
			"X509_gmtime_adj(X509_get0_notBefore(*cert)")
		patched = strings.ReplaceAll(patched,
			"X509_gmtime_adj(X509_get_notAfter(*cert)",
			"X509_gmtime_adj(X509_get0_notAfter(*cert)")

		// Write back
		if err := os.WriteFile(patchFile, []byte(patched), 0644); err != nil {
			return fmt.Errorf("failed to write patched tls_openssl.c: %w", err)
		}

		fmt.Printf("Applied OpenSSL 3.6 compatibility patch to tls_openssl.c\n")
		return nil
	},
	ConfigureArgs: func(os string) []string {
		// FFmpeg needs explicit paths to headers and libraries
		stagingDir, _ := filepath.Abs(".build/staging")
		incDir := filepath.Join(stagingDir, "include")
		libDir := filepath.Join(stagingDir, "lib")
		lib64Dir := filepath.Join(stagingDir, "lib64")

		args := []string{
			"--pkg-config-flags=--static",
			fmt.Sprintf("--extra-cflags=-I%s", incDir),
			fmt.Sprintf("--extra-ldflags=-L%s -L%s", libDir, lib64Dir),
		}

		// Add common FFmpeg arguments
		args = append(args, FFmpegArgsCommon()...)

		// Add platform-specific arguments
		if os == "linux" {
			args = append(args, FFmpegArgsLinux()...)
		} else if os == "darwin" {
			args = append(args, FFmpegArgsDarwin()...)
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
