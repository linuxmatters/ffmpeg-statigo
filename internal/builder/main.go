package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	tempDir, _   = filepath.Abs("temp")
	downloadsDir = path.Join(tempDir, "downloads")
	buildDir     = path.Join(tempDir, "build")
	tgtDir       = path.Join(tempDir, "tgt")
	incDir       = path.Join(tgtDir, "include")
	libDir       = path.Join(tgtDir, "lib")

	libs = []string{
		"librav1e",
		"libass",
		"libavcodec",
		"libavdevice",
		"libavfilter",
		"libavformat",
		"libavutil",
		"libbrotlicommon",
		"libbrotlidec",
		"libbrotlienc",
		"libbz2",
		"libfreetype",
		"libfribidi",
		"libharfbuzz",
		"libmp3lame",
		"libogg",
		"libopus",
		"libpng16",
		"libspeex",
		"libswresample",
		"libswscale",
		"libtheora",
		"libunibreak",
		"libvpx",
		"libx264",
		"libx265",
		"libdav1d",
		"libz",
	}
)

type OS string

var (
	MacOS OS = "macos"
	Linux OS = "linux"
)

type Builder struct {
	os OS
}

func main() {
	targetOutput, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Panicln(err)
	}

	//os.RemoveAll(tempDir)
	os.RemoveAll(buildDir)
	os.RemoveAll(tgtDir)

	if err := os.MkdirAll(downloadsDir, 0755); err != nil {
		log.Panicln(err)
	}

	if err := os.MkdirAll(buildDir, 0755); err != nil {
		log.Panicln(err)
	}

	if err := os.MkdirAll(tgtDir, 0755); err != nil {
		log.Panicln(err)
	}

	if err := os.MkdirAll(libDir, 0755); err != nil {
		log.Panicln(err)
	}

	os := Linux

	if runtime.GOOS == "darwin" {
		os = MacOS
	}

	b := Builder{
		os: os,
	}

	if os == MacOS {
		buildIconv()
		libs = append(libs, "libiconv")
	}

	buildZlib()
	buildBZ2()
	buildBrotli()
	buildPng()
	buildFribidi()
	buildUnibreak()
	buildFreetype()

	if os == Linux {
		buildExpat()
		buildFontconfig()

		libs = append(libs, "libexpat", "libfontconfig")
	}

	buildHarfbuzz()
	b.buildASS()
	buildRav1e()

	// NVENC/NVDEC support (Linux only, requires NVIDIA GPU drivers)
	if os == Linux {
		b.buildNVCodec()
	}

	buildLame()
	buildOpus()
	buildOgg()
	buildVorbis()
	buildSpeex()
	buildTheora()
	buildVpx()
	buildX264()
	buildX265()
	buildDav1d()

	b.buildFFmpeg()

	if b.os == Linux {
		combineLinux(targetOutput)
	} else if b.os == MacOS {
		combineMac(targetOutput)
	}
}

func buildX264() {
	zipPath := path.Join(downloadsDir, "x264.tar.bz2")
	srcPath := path.Join(buildDir, "x264")

	if !exists(zipPath) {
		download("https://code.videolan.org/videolan/x264/-/archive/master/x264-master.tar.bz2", zipPath)
	}

	untar(zipPath, srcPath, "x264-master/")

	{
		log.Println("Running configure")

		cmd := cmd(
			path.Join(srcPath, "configure"),
			srcPath,
			fmt.Sprintf("--prefix=%v", tgtDir),
			"--disable-avs",
			"--disable-cli",
			"--disable-ffms",
			"--disable-gpac",
			"--disable-lavf",
			"--disable-lsmash",
			"--disable-swscale",
			"--enable-static",
			"--enable-strip",
		)
		// x264 needs to find nasm explicitly on x86/x86_64
		// ARM architectures use the C compiler as assembler instead
		if runtime.GOARCH == "amd64" || runtime.GOARCH == "386" {
			cmd.Env = append(cmd.Env, "AS=nasm")
		}

		run("[x264 configure]", cmd)
	}

	{
		log.Println("Running make")

		cmd := cmd(
			"make",
			srcPath,
			"-j8",
			"install",
		)

		run("[x264 make]", cmd)
	}
}

func buildX265() {
	zipPath := path.Join(downloadsDir, "x265.tar.bz2")
	srcPath := path.Join(buildDir, "x265")
	buildPath := path.Join(buildDir, "x265-build")

	if !exists(zipPath) {
		// x265 master branch HEAD (ffba52bab55d) - newer with CMake fixes
		download("https://bitbucket.org/multicoreware/x265_git/get/ffba52bab55dce9b1b3a97dd08d12e70297e2180.tar.bz2", zipPath)
	}

	untar(zipPath, srcPath, "multicoreware-x265_git-ffba52bab55d/")

	if err := os.MkdirAll(buildPath, 0755); err != nil {
		log.Panicln(err)
	}

	{
		log.Println("Running cmake")

		cmakeArgs := []string{
			"-G", "Unix Makefiles",
			fmt.Sprintf("-DCMAKE_INSTALL_PREFIX=%v", tgtDir),
			"-DCMAKE_BUILD_TYPE=Release", // Optimized build
			"-DENABLE_SHARED=OFF",        // Static library
			"-DENABLE_CLI=OFF",           // No CLI tools needed
			"-DENABLE_AGGRESSIVE_CHECKS=OFF",
			path.Join(srcPath, "source"), // x265 source is in source/ subdirectory
		}

		cmdArgs := append([]string{buildPath}, cmakeArgs...)
		cmd := cmd("cmake", cmdArgs[0], cmdArgs[1:]...)

		run("[x265 cmake]", cmd)
	}

	{
		log.Println("Running make")

		cmd := cmd(
			"make",
			buildPath,
			"-j8",
			"install",
		)

		run("[x265 make]", cmd)
	}
}

func buildDav1d() {
	zipPath := path.Join(downloadsDir, "dav1d.tar.bz2")
	srcPath := path.Join(buildDir, "dav1d")
	buildPath := path.Join(buildDir, "dav1d-build")

	if !exists(zipPath) {
		// dav1d 1.5.2 - fast AV1 decoder from VideoLAN
		download("https://code.videolan.org/videolan/dav1d/-/archive/1.5.2/dav1d-1.5.2.tar.bz2", zipPath)
	}

	untar(zipPath, srcPath, "dav1d-1.5.2/")

	{
		log.Println("Running meson setup")

		cmd := cmd(
			"meson",
			"",
			"setup",
			buildPath,
			srcPath,
			fmt.Sprintf("--prefix=%v", tgtDir),
			"--default-library=static",
			"--buildtype=release",
			"-Denable_tools=false",
			"-Denable_tests=false",
		)

		run("[dav1d meson]", cmd)
	}

	{
		log.Println("Running ninja install")

		cmd := cmd(
			"ninja",
			buildPath,
			"-C", buildPath,
			"install",
		)

		run("[dav1d ninja]", cmd)
	}
}

func buildVpx() {
	zipPath := path.Join(downloadsDir, "vpx.zip")
	srcPath := path.Join(buildDir, "vpx")
	buildDir := path.Join(srcPath, "ffbuild")

	if !exists(zipPath) {
		download("https://github.com/webmproject/libvpx/archive/refs/tags/v1.15.2.zip", zipPath)
	}

	unzip(zipPath, srcPath)

	if err := os.MkdirAll(buildDir, 0755); err != nil {
		log.Panicln(err)
	}

	{
		log.Println("Running configure")

		cmd := cmd(
			path.Join(srcPath, "configure"),
			buildDir,
			fmt.Sprintf("--prefix=%v", tgtDir),
			"--disable-dependency-tracking",
			"--enable-static",
			"--disable-shared",
			"--disable-examples",
			"--enable-vp9-highbitdepth",
			"--disable-unit-tests",
			"--as=yasm",
		)

		// TODO: target/cpudetect

		run("[vpx configure]", cmd)
	}

	{
		log.Println("Running make")

		cmd := cmd(
			"make",
			buildDir,
			"-j8",
			"install",
		)

		run("[vps make]", cmd)
	}
}

func buildTheora() {
	zipPath := path.Join(downloadsDir, "theora.tar.xz")
	srcPath := path.Join(buildDir, "theora")

	if !exists(zipPath) {
		download("https://ftp.osuosl.org/pub/xiph/releases/theora/libtheora-1.2.0.tar.xz", zipPath)
	}

	untar(zipPath, srcPath, "libtheora-1.2.0/")

	{
		log.Println("Running configure")

		cmd := cmd(
			path.Join(srcPath, "configure"),
			srcPath,
			fmt.Sprintf("--prefix=%v", tgtDir),
			"--disable-dependency-tracking",
			"--enable-static",
			"--disable-shared",
			"--disable-oggtest",
			"--disable-vorbistest",
			"--disable-examples",
			fmt.Sprintf("CFLAGS=-I%v", incDir),
			fmt.Sprintf("CPPFLAGS=-I%v", incDir),
			fmt.Sprintf("LDFLAGS=-L%v", libDir),
		)

		run("[theora configure]", cmd)
	}

	// Prevent automake regeneration
	touchAutomakeFiles(srcPath)

	{
		log.Println("Running make")

		cmd := cmd(
			"make",
			srcPath,
			"-j8",
			"install",
		)

		run("[theora make]", cmd)
	}
}

func buildOgg() {
	zipPath := path.Join(downloadsDir, "ogg.tar.xz")
	srcPath := path.Join(buildDir, "ogg")

	if !exists(zipPath) {
		download("https://ftp.osuosl.org/pub/xiph/releases/ogg/libogg-1.3.6.tar.xz", zipPath)
	}

	untar(zipPath, srcPath, "libogg-1.3.6/")

	{
		log.Println("Running configure")

		cmd := cmd(
			path.Join(srcPath, "configure"),
			srcPath,
			fmt.Sprintf("--prefix=%v", tgtDir),
			"--disable-dependency-tracking",
			"--enable-static",
			"--disable-shared",
			fmt.Sprintf("CFLAGS=-I%v", incDir),
			fmt.Sprintf("CPPFLAGS=-I%v", incDir),
			fmt.Sprintf("LDFLAGS=-L%v", libDir),
		)

		run("[ogg configure]", cmd)
	}

	// Prevent automake regeneration
	touchAutomakeFiles(srcPath)

	{
		log.Println("Running make")

		cmd := cmd(
			"make",
			srcPath,
			"-j8",
			"install",
		)

		run("[ogg make]", cmd)
	}
}

func buildVorbis() {
	zipPath := path.Join(downloadsDir, "vorbis.tar.gz")
	srcPath := path.Join(buildDir, "vorbis")

	if !exists(zipPath) {
		download("https://ftp.osuosl.org/pub/xiph/releases/vorbis/libvorbis-1.3.7.tar.gz", zipPath)
	}

	untar(zipPath, srcPath, "libvorbis-1.3.7/")

	modify(
		path.Join(srcPath, "configure"),
		func(bytes []byte) []byte {
			return []byte(strings.ReplaceAll(string(bytes), "-force_cpusubtype_ALL", ""))
		},
	)

	{
		log.Println("Running configure")

		cmd := cmd(
			path.Join(srcPath, "configure"),
			srcPath,
			fmt.Sprintf("--prefix=%v", tgtDir),
			"--disable-dependency-tracking",
			"--enable-static",
			"--disable-shared",
			fmt.Sprintf("CFLAGS=-I%v", incDir),
			fmt.Sprintf("CPPFLAGS=-I%v", incDir),
			fmt.Sprintf("LDFLAGS=-L%v", libDir),
		)

		run("[vorbis configure]", cmd)
	}

	{
		log.Println("Running make")

		cmd := cmd(
			"make",
			srcPath,
			"-j8",
			"install",
		)

		run("[vorbis make]", cmd)
	}
}

func buildSpeex() {
	zipPath := path.Join(downloadsDir, "speex.tar.gz")
	srcPath := path.Join(buildDir, "speex")

	if !exists(zipPath) {
		download("http://downloads.xiph.org/releases/speex/speex-1.2.1.tar.gz", zipPath)
	}

	untar(zipPath, srcPath, "speex-1.2.1/")

	{
		log.Println("Running configure")

		cmd := cmd(
			path.Join(srcPath, "configure"),
			srcPath,
			fmt.Sprintf("--prefix=%v", tgtDir),
			"--disable-dependency-tracking",
			"--enable-static",
			"--disable-shared",
			fmt.Sprintf("CFLAGS=-I%v", incDir),
			fmt.Sprintf("CPPFLAGS=-I%v", incDir),
			fmt.Sprintf("LDFLAGS=-L%v", libDir),
		)

		run("[speex configure]", cmd)
	}

	// Prevent automake regeneration
	touchAutomakeFiles(srcPath)

	{
		log.Println("Running make")

		cmd := cmd(
			"make",
			srcPath,
			"-j8",
			"install",
		)

		run("[speex make]", cmd)
	}
}

func buildOpus() {
	zipPath := path.Join(downloadsDir, "opus.tar.gz")
	srcPath := path.Join(buildDir, "opus")

	if !exists(zipPath) {
		download("https://ftp.osuosl.org/pub/xiph/releases/opus/opus-1.5.2.tar.gz", zipPath)
	}

	untar(zipPath, srcPath, "opus-1.5.2/")

	{
		log.Println("Running configure")

		cmd := cmd(
			path.Join(srcPath, "configure"),
			srcPath,
			fmt.Sprintf("--prefix=%v", tgtDir),
			"--disable-dependency-tracking",
			"--disable-debug",
			"--disable-doc",
			"--enable-static",
			"--disable-shared",
			fmt.Sprintf("CFLAGS=-I%v", incDir),
			fmt.Sprintf("CPPFLAGS=-I%v", incDir),
			fmt.Sprintf("LDFLAGS=-L%v", libDir),
		)

		run("[opus configure]", cmd)
	}

	// Prevent automake regeneration
	touchAutomakeFiles(srcPath)

	{
		log.Println("Running make")

		cmd := cmd(
			"make",
			srcPath,
			"-j8",
			"install",
		)

		run("[opus make]", cmd)
	}
}

func buildLame() {
	zipPath := path.Join(downloadsDir, "lame.tar.gz")
	srcPath := path.Join(buildDir, "lame")

	if !exists(zipPath) {
		download("https://downloads.sourceforge.net/project/lame/lame/3.100/lame-3.100.tar.gz", zipPath)
	}

	untar(zipPath, srcPath, "lame-3.100/")

	{
		log.Println("Running configure")

		cmd := cmd(
			path.Join(srcPath, "configure"),
			srcPath,
			fmt.Sprintf("--prefix=%v", tgtDir),
			"--disable-dependency-tracking",
			"--disable-debug",
			"--enable-static",
			"--disable-shared",
			fmt.Sprintf("CFLAGS=-I%v", incDir),
			fmt.Sprintf("CPPFLAGS=-I%v", incDir),
			fmt.Sprintf("LDFLAGS=-L%v", libDir),
		)

		run("[lame configure]", cmd)
	}

	{
		log.Println("Running make")

		cmd := cmd(
			"make",
			srcPath,
			"-j8",
			"install",
		)

		run("[lame make]", cmd)
	}
}

func buildExpat() {
	zipPath := path.Join(downloadsDir, "expat.tar.gz")
	srcPath := path.Join(buildDir, "expat")

	if !exists(zipPath) {
		download("https://github.com/libexpat/libexpat/releases/download/R_2_7_3/expat-2.7.3.tar.gz", zipPath)
	}

	untar(zipPath, srcPath, "expat-2.7.3/")

	{
		log.Println("Running configure")

		cmd := cmd(
			path.Join(srcPath, "configure"),
			srcPath,
			fmt.Sprintf("--prefix=%v", tgtDir),
			"--enable-static",
			"--disable-shared",
			"--without-xmlwf",
			fmt.Sprintf("CFLAGS=-I%v", incDir),
			fmt.Sprintf("CPPFLAGS=-I%v", incDir),
			fmt.Sprintf("LDFLAGS=-L%v", libDir),
		)

		run("[expat configure]", cmd)
	}

	// Prevent automake regeneration
	touchAutomakeFiles(srcPath)

	{
		log.Println("Running make")

		cmd := cmd(
			"make",
			srcPath,
			"-j8",
			"install",
		)

		run("[expat make]", cmd)
	}
}

func buildFontconfig() {
	zipPath := path.Join(downloadsDir, "fontconfig.tar.xz")
	srcPath := path.Join(buildDir, "fontconfig")

	if !exists(zipPath) {
		download("https://www.freedesktop.org/software/fontconfig/release/fontconfig-2.16.0.tar.xz", zipPath)
	}

	untar(zipPath, srcPath, "fontconfig-2.16.0/")

	{
		log.Println("Running configure")

		cmd := cmd(
			path.Join(srcPath, "configure"),
			srcPath,
			fmt.Sprintf("--prefix=%v", tgtDir),
			"--disable-debug",
			"--enable-static",
			"--disable-shared",
			fmt.Sprintf("CFLAGS=-I%v", incDir),
			fmt.Sprintf("CPPFLAGS=-I%v", incDir),
			fmt.Sprintf("LDFLAGS=-L%v", libDir),
		)

		run("[fontconfig configure]", cmd)
	}

	// Prevent automake regeneration
	touchAutomakeFiles(srcPath)

	{
		log.Println("Running make")

		cmd := cmd(
			"make",
			srcPath,
			"-j8",
			"install",
		)

		run("[fontconfig make]", cmd)
	}
}

func buildFribidi() {
	zipPath := path.Join(downloadsDir, "fribidi.tar.xz")
	srcPath := path.Join(buildDir, "fribidi")

	if !exists(zipPath) {
		download("https://github.com/fribidi/fribidi/releases/download/v1.0.16/fribidi-1.0.16.tar.xz", zipPath)
	}

	untar(zipPath, srcPath, "fribidi-1.0.16/")

	{
		log.Println("Running configure")

		cmd := cmd(
			path.Join(srcPath, "configure"),
			srcPath,
			fmt.Sprintf("--prefix=%v", tgtDir),
			"--disable-dependency-tracking",
			"--disable-debug",
			"--disable-silent-rules",
			"--enable-static",
			fmt.Sprintf("CFLAGS=-I%v", incDir),
			fmt.Sprintf("CPPFLAGS=-I%v", incDir),
			fmt.Sprintf("LDFLAGS=-L%v", libDir),
		)

		run("[fribidi configure]", cmd)
	}

	// Prevent automake regeneration
	touchAutomakeFiles(srcPath)

	{
		log.Println("Running make")

		cmd := cmd(
			"make",
			srcPath,
			"-j8",
			"install",
		)

		run("[fribidi make]", cmd)
	}
}

func buildIconv() {
	zipPath := path.Join(downloadsDir, "iconv.tar.gz")
	srcPath := path.Join(buildDir, "iconv")

	if !exists(zipPath) {
		download("https://ftp.mirrorservice.org/pub/gnu/libiconv/libiconv-1.18.tar.gz", zipPath)
	}

	untar(zipPath, srcPath, "libiconv-1.18/")

	{
		log.Println("Running configure")

		cmd := cmd(
			path.Join(srcPath, "configure"),
			srcPath,
			fmt.Sprintf("--prefix=%v", tgtDir),
			"--disable-dependency-tracking",
			"--disable-debug",
			"--enable-extra-encodings",
			"--enable-static",
			fmt.Sprintf("CFLAGS=-I%v", incDir),
			fmt.Sprintf("CPPFLAGS=-I%v", incDir),
			fmt.Sprintf("LDFLAGS=-L%v", libDir),
		)

		run("[iconv configure]", cmd)
	}

	//{
	//	log.Println("Running make")
	//
	//	cmd := cmd(
	//		"make",
	//		srcPath,
	//		"-j8",
	//		"-f", "Makefile.devel",
	//	)
	//	run("[iconv make]", cmd)
	//}

	{
		log.Println("Running install")

		cmd := cmd(
			"make",
			srcPath,
			"install",
		)

		run("[iconv install]", cmd)
	}
}

func buildZlib() {
	zipPath := path.Join(downloadsDir, "zlib.zip")
	srcPath := path.Join(buildDir, "zlib")

	if !exists(zipPath) {
		download("https://github.com/madler/zlib/releases/download/v1.3.1/zlib131.zip", zipPath)
	}

	unzip(zipPath, srcPath)

	{
		log.Println("Running configure ??")

		cmd := cmd(
			path.Join(srcPath, "configure"),
			srcPath,
			fmt.Sprintf("--prefix=%v", tgtDir),
			"--static",
		)

		run("[zlib configure]", cmd)
	}

	{
		log.Println("Running make")

		cmd := cmd(
			"make",
			srcPath,
			"-j8",
			"install",
		)

		run("[zlib make]", cmd)
	}
}

func buildBZ2() {
	zipPath := path.Join(downloadsDir, "bz2.zip")
	srcPath := path.Join(buildDir, "bz2")

	if !exists(zipPath) {
		download("https://gitlab.com/bzip2/bzip2/-/archive/bzip2-1.0.8/bzip2-bzip2-1.0.8.zip", zipPath)
	}

	unzip(zipPath, srcPath)

	{
		log.Println("Running make")

		cmd := cmd(
			"make",
			srcPath,
			"-j8",
			"install",
			fmt.Sprintf("PREFIX=%v", tgtDir),
			fmt.Sprintf("CFLAGS=-I%v", incDir),
			fmt.Sprintf("CPPFLAGS=-I%v", incDir),
			fmt.Sprintf("LDFLAGS=-L%v", libDir),
		)

		run("[bz2 make]", cmd)
	}
}

func buildBrotli() {
	zipPath := path.Join(downloadsDir, "brotli.zip")
	srcPath := path.Join(buildDir, "brotli")

	if !exists(zipPath) {
		download("https://github.com/google/brotli/archive/refs/tags/v1.2.0.zip", zipPath)
	}

	unzip(zipPath, srcPath)

	{
		log.Println("Running cmake")

		cmd := cmd(
			"cmake",
			srcPath,
			"-G",
			"Unix Makefiles",
			fmt.Sprintf("-DCMAKE_INSTALL_PREFIX=%v", tgtDir),
			"-DCMAKE_BUILD_TYPE=Release",
			"-DCMAKE_FIND_FRAMEWORK=last",
			"-DCMAKE_VERBOSE_MAKEFILE=ON",
			"-Wno-dev",
			"-DBUILD_SHARED_LIBS=OFF",
			".",
		)

		run("[brotli cmake]", cmd)
	}

	{
		log.Println("Running make")

		cmd := cmd(
			"make",
			srcPath,
			"-j8",
			"install",
		)

		run("[brotli make]", cmd)
	}
}

func buildPng() {
	zipPath := path.Join(downloadsDir, "png.zip")
	srcPath := path.Join(buildDir, "png")

	if !exists(zipPath) {
		download("https://github.com/pnggroup/libpng/archive/refs/tags/v1.6.50.zip", zipPath)
	}

	unzip(zipPath, srcPath)

	{
		log.Println("Running configure")

		cmd := cmd(
			path.Join(srcPath, "configure"),
			srcPath,
			fmt.Sprintf("--prefix=%v", tgtDir),
			"--disable-dependency-tracking",
			"--disable-silent-rules",
			"--disable-shared",
			"--enable-static",
			fmt.Sprintf("CFLAGS=-I%v", incDir),
			fmt.Sprintf("CPPFLAGS=-I%v", incDir),
			fmt.Sprintf("LDFLAGS=-L%v", libDir),
		)

		run("[png configure]", cmd)
	}

	{
		log.Println("Running make")

		cmd := cmd(
			"make",
			srcPath,
			"-j8",
			"install",
		)

		run("[png make]", cmd)
	}
}

func buildUnibreak() {
	zipPath := path.Join(downloadsDir, "unibreak.tar.gz")
	srcPath := path.Join(buildDir, "unibreak")

	if !exists(zipPath) {
		download("https://github.com/adah1972/libunibreak/releases/download/libunibreak_6_1/libunibreak-6.1.tar.gz", zipPath)
	}

	untar(zipPath, srcPath, "libunibreak-6.1/")

	{
		log.Println("Running configure")

		cmd := cmd(
			path.Join(srcPath, "configure"),
			srcPath,
			fmt.Sprintf("--prefix=%v", tgtDir),
			"--disable-shared",
			"--enable-static",
			fmt.Sprintf("CFLAGS=-I%v", incDir),
			fmt.Sprintf("CPPFLAGS=-I%v", incDir),
			fmt.Sprintf("LDFLAGS=-L%v", libDir),
		)

		run("[unibreak configure]", cmd)
	}

	// Prevent automake regeneration
	touchAutomakeFiles(srcPath)

	{
		log.Println("Running make")

		cmd := cmd(
			"make",
			srcPath,
			"-j8",
			"install",
		)

		run("[unibreak make]", cmd)
	}
}

func buildFreetype() {
	zipPath := path.Join(downloadsDir, "freetype.tar.xz")
	srcPath := path.Join(buildDir, "freetype")

	if !exists(zipPath) {
		download("https://download-mirror.savannah.gnu.org/releases/freetype/freetype-2.14.1.tar.xz", zipPath)
	}

	untar(zipPath, srcPath, "freetype-2.14.1/")

	{
		log.Println("Running configure")

		cmd := cmd(
			path.Join(srcPath, "configure"),
			srcPath,
			fmt.Sprintf("--prefix=%v", tgtDir),
			"--without-harfbuzz",
			// brotlii breaks harfbuzz building
			"--without-brotli",
			"--disable-shared",
			"--enable-static",
			fmt.Sprintf("CFLAGS=-I%v", incDir),
			fmt.Sprintf("CPPFLAGS=-I%v", incDir),
			fmt.Sprintf("LDFLAGS=-L%v", libDir),
		)

		run("[freetype configure]", cmd)
	}

	{
		log.Println("Running make")

		cmd := cmd(
			"make",
			srcPath,
			"-j8",
			"install",
		)

		run("[freetype make]", cmd)
	}
}

func buildHarfbuzz() {
	zipPath := path.Join(downloadsDir, "harfbuzz.tar.xz")
	srcPath := path.Join(buildDir, "harfbuzz")

	if !exists(zipPath) {
		download("https://github.com/harfbuzz/harfbuzz/releases/download/12.2.0/harfbuzz-12.2.0.tar.xz", zipPath)
	}

	untar(zipPath, srcPath, "harfbuzz-12.2.0/")

	{
		log.Println("Running setup")

		cmd := cmd(
			"meson",
			srcPath,
			"setup",
			"build",
			fmt.Sprintf("--prefix=%v", tgtDir),
			fmt.Sprintf("--libdir=%v", libDir),
			"--buildtype=release",
			"--default-library=static",
			"-Dcairo=disabled",
			"-Dcoretext=enabled",
			"-Dfreetype=enabled",
			//"-Dglib=enabled",
			//"-Dgobject=enabled",
			//"-Dgraphite=enabled",
			//"-Dicu=enabled",
			//"-Dintrospection=enabled",
			"-Dtests=disabled",
		)

		run("[harfbuzz setup]", cmd)
	}

	{
		log.Println("Running compile")

		cmd := cmd(
			"meson",
			srcPath,
			"compile",
			"-C",
			"build",
			"--verbose",
		)

		run("[harfbuzz compile]", cmd)
	}

	{
		log.Println("Running install")

		cmd := cmd(
			"meson",
			srcPath,
			"install",
			"-C",
			"build",
		)

		run("[harfbuzz install]", cmd)
	}
}

func (b *Builder) buildASS() {
	zipPath := path.Join(downloadsDir, "ass.tar.gz")
	srcPath := path.Join(buildDir, "ass")

	if !exists(zipPath) {
		download("https://github.com/libass/libass/releases/download/0.17.4/libass-0.17.4.tar.gz", zipPath)
	}

	untar(zipPath, srcPath, "libass-0.17.4/")

	{
		log.Println("Running configure")

		cmd := cmd(
			path.Join(srcPath, "configure"),
			srcPath,
			fmt.Sprintf("--prefix=%v", tgtDir),
			"--disable-shared",
		)

		// libass uses coretext on macOS, fontconfig on Linux
		if b.os == MacOS {
			cmd.Args = append(
				cmd.Args,
				"--disable-fontconfig",
			)
		}

		cmd.Args = append(cmd.Args, fmt.Sprintf("CFLAGS=-I%v", incDir))

		run("[ass configure]", cmd)
	}

	// Prevent automake regeneration
	touchAutomakeFiles(srcPath)

	{
		log.Println("Running make")

		cmd := cmd(
			"make",
			srcPath,
			"-j8",
			"install",
		)

		run("[ass make]", cmd)
	}
}

func buildRav1e() {
	zipPath := path.Join(downloadsDir, "rav1e.tar.gz")
	srcPath := path.Join(buildDir, "rav1e")

	if !exists(zipPath) {
		// rav1e v0.8.1 - Fastest and safest AV1 encoder from Xiph
		download("https://github.com/xiph/rav1e/archive/refs/tags/v0.8.1.tar.gz", zipPath)
	}

	untar(zipPath, srcPath, "rav1e-0.8.1/")

	{
		log.Println("Running cargo cinstall for rav1e")

		// Set RUSTFLAGS for native CPU optimization
		cmd := cmd(
			"cargo",
			srcPath,
			"cinstall",
			fmt.Sprintf("--prefix=%v", tgtDir),
			"--libdir=lib",
			"--library-type=staticlib",
			"--crt-static",
			"--release",
		)

		// Add RUSTFLAGS environment variable for native CPU optimization
		cmd.Env = append(os.Environ(), "RUSTFLAGS=-C target-cpu=native")

		run("[rav1e cargo]", cmd)
	}
}

func (b *Builder) buildNVCodec() {
	// nv-codec-headers provides NVENC/NVDEC headers without requiring CUDA runtime
	// This is a header-only library that enables hardware accelerated encoding/decoding
	// on NVIDIA GPUs via the Video Codec SDK
	zipPath := path.Join(downloadsDir, "nv-codec-headers.tar.gz")
	srcPath := path.Join(buildDir, "nv-codec-headers")

	if !exists(zipPath) {
		// nv-codec-headers 12.2.72.0 (latest stable)
		download("https://github.com/FFmpeg/nv-codec-headers/releases/download/n12.2.72.0/nv-codec-headers-12.2.72.0.tar.gz", zipPath)
	}

	untar(zipPath, srcPath, "nv-codec-headers-12.2.72.0/")

	{
		log.Println("Installing nv-codec-headers")

		// nv-codec-headers just installs header files, no compilation needed
		cmd := cmd(
			"make",
			srcPath,
			fmt.Sprintf("PREFIX=%v", tgtDir),
			"install",
		)

		run("[nvcodec install]", cmd)
	}
}

func (b *Builder) buildFFmpeg() {
	zipPath := path.Join(downloadsDir, "ffmpeg.zip")
	buildPath := path.Join(buildDir, "ffmpeg")

	if !exists(zipPath) {
		download("https://codeload.github.com/FFmpeg/FFmpeg/zip/refs/heads/release/8.0", zipPath)
	}

	unzip(zipPath, buildPath)

	{
		log.Println("Running configure")

		cmd := cmd(
			path.Join(buildPath, "configure"),
			buildPath,
			fmt.Sprintf("--prefix=%v", tgtDir),
			"--pkg-config-flags=--static",
			fmt.Sprintf("--extra-cflags=-I%v", incDir),
			fmt.Sprintf("--extra-ldflags=-L%v/lib", tgtDir),
			"--disable-autodetect",
			"--disable-debug",
			"--disable-doc",
			"--disable-htmlpages",
			"--disable-manpages",
			"--disable-podpages",
			"--disable-programs",
			"--disable-txtpages",
			"--disable-decoder=rv10,rv20,rv30,rv40,rv60",                            // RealVideo
			"--disable-decoder=ra_144,ra_288",                                       // RealAudio 14.4/28.8 kbps (1995-1997)
			"--disable-demuxer=rm",                                                  // RealMedia demuxer
			"--disable-muxer=rm",                                                    // RealMedia muxer
			"--disable-decoder=cook,sipr",                                           // RealAudio variants
			"--disable-parser=cook,sipr",                                            // RealAudio variants
			"--disable-decoder=vp3,vp4,vp5,vp6,vp6a,vp6f,vp7",                       // Early VP
			"--disable-demuxer=ivf",                                                 // IVF container (VP8/VP9 test format)
			"--disable-parser=vp3",                                                  // Early VP
			"--disable-decoder=truemotion1,truemotion2,truemotion2rt",               // Duck corp
			"--disable-decoder=cinepak,indeo2,indeo3,indeo4,indeo5,msrle",           // 1990s legacy
			"--disable-encoder=cinepak,msrle",                                       // 1990s legacy
			"--disable-decoder=msvideo1,wmv1,wmv2,wmv3,wmv3image",                   // Windows Media
			"--disable-encoder=msvideo1,wmv1,wmv2,wmav1,wmav2",                      // Windows Media
			"--disable-decoder=wmalossless,wmapro,wmav1,wmav2,wmavoice",             // Windows Media
			"--disable-decoder=vc1,vc1_cuvid,vc1_mmal,vc1_qsv,vc1_v4l2m2m,vc1image", // VC-1
			"--disable-demuxer=vc1,vc1t",                                            // VC-1
			"--disable-muxer=vc1,vc1t",                                              // VC-1
			"--disable-parser=vc1",                                                  // VC-1
			"--disable-decoder=asf",                                                 // ASF
			"--disable-demuxer=asf,asf_o",                                           // ASF
			"--disable-muxer=asf,asf_stream",                                        // ASF
			"--disable-encoder=mpeg4,msmpeg4v2,msmpeg4v3",                           // Old MS-MPEG-4 encoders
			"--disable-decoder=mpeg4,msmpeg4v1,msmpeg4v2,msmpeg4v3",                 // Old MS-MPEG-4 decoders
			"--disable-encoder=h261,h263,h263_v4l2m2m,h263p",                        // H.261/H.263
			"--disable-decoder=h261,h263,h263_v4l2m2m,h263i,h263p",                  // H.261/H.263
			"--disable-demuxer=h261,h263",                                           // H.261/H.263
			"--disable-muxer=h261,h263",                                             // H.261/H.263
			"--disable-parser=h261,h263",                                            // H.261/H.263
			"--disable-encoder=flv",                                                 // Flash Sorenson H.263
			"--disable-decoder=flv",                                                 // Flash Sorenson H.263
			"--disable-demuxer=flv,live_flv",                                        // Flash Sorenson H.263
			"--disable-muxer=flv",                                                   // Flash Sorenson H.263
			"--disable-decoder=eacmv,eamad,eatgq,eatgv,eatqi",                       // Electronic Arts
			"--disable-decoder=adpcm_ea_maxis_xa,adpcm_ea_r1,adpcm_ea_r2,adpcm_ea_r3,adpcm_ea_xas,adpcm_ima_ea_eacs,adpcm_ima_ea_sead", // Electronic Arts
			"--disable-decoder=flic",                                                      // Autodesk FLIC 1990s
			"--disable-demuxer=flic",                                                      // Autodesk FLIC 1990s
			"--disable-decoder=anm",                                                       // Deluxe Paint 1980s
			"--disable-decoder=adpcm_4xm,4xm",                                             // 4X Movie
			"--disable-decoder=interplay_acm,interplay_dpcm,interplay_video",              // Interplay 1990s
			"--disable-decoder=bethsoftvid",                                               // Bethesda pre-2002
			"--disable-demuxer=bethsoftvid",                                               // Bethesda pre-2002
			"--disable-decoder=vqa",                                                       // Westwood
			"--disable-demuxer=wsvqa",                                                     // Westwood
			"--disable-decoder=dsicinaudio,dsicinvideo",                                   // Delphine Software
			"--disable-decoder=dsicin",                                                    // Delphine Software
			"--disable-decoder=idcin,roq,roq_dpcm",                                        // Id Software
			"--disable-encoder=roq,roq_dpcm",                                              // Id Software
			"--disable-demuxer=idcin,roq",                                                 // Id Software
			"--disable-muxer=roq",                                                         // Id Software
			"--disable-decoder=bink,binkaudio_dct,binkaudio_rdft,smackaud,smacker,",       // RAD Game Tools
			"--disable-demuxer=bink,binka,smacker",                                        // RAD Game Tools
			"--disable-decoder=adpcm_argo,argo",                                           // Argonaut Games
			"--disable-encoder=adpcm_argo",                                                // Argonaut Games
			"--disable-demuxer=argo_asf,argo_brp,argo_cvg",                                // Argonaut Games
			"--disable-muxer=argo_asf,argo_cvg",                                           // Argonaut Games
			"--disable-encoder=dvvideo",                                                   // DV tape
			"--disable-decoder=dvaudio,dvvideo",                                           // DV tape
			"--disable-demuxer=dv",                                                        // DV tape
			"--disable-muxer=dv",                                                          // DV tape
			"--disable-parser=dvaudio",                                                    // DV tape
			"--disable-decoder=adpcm_psx,adpcm_xa",                                        // PlayStation 1 / CD-ROM XA
			"--disable-decoder=atrac1,atrac3,atrac3al,atrac3p,atrac3pal,atrac9",           // Sony ATRAC (MiniDisc/PS)
			"--disable-encoder=svq1",                                                      // Sorenson Video (QuickTime)
			"--disable-decoder=svq1,svq3",                                                 // Sorenson Video (QuickTime)
			"--disable-decoder=qdm2,qdm2_at,qdmc,qdmc_at",                                 // QDesign Music (QuickTime)
			"--disable-encoder=rpza,smc",                                                  // Apple legacy
			"--disable-decoder=mace3,mace6,rpza,smc",                                      // Apple legacy
			"--disable-decoder=cavs",                                                      // Chinese AVS
			"--disable-demuxer=cavsvideo",                                                 // Chinese AVS
			"--disable-muxer=cavsvideo",                                                   // Chinese AVS
			"--disable-parser=cavsvideo",                                                  // Chinese AVS
			"--disable-encoder=asv1,asv2",                                                 // ASUS V1/V2 encoders
			"--disable-decoder=asv1,asv2",                                                 // ASUS V1/V2 decoders
			"--disable-encoder=amv",                                                       // AMV video
			"--disable-decoder=amv",                                                       // AMV video
			"--disable-muxer=amv",                                                         // AMV video
			"--disable-decoder=amr_nb_at,amrnb,amrnb_mediacodec,amrwb,amrwb_mediacodec",   // AMR mobile codecs
			"--disable-demuxer=amr,amrnb,amrwb",                                           // AMR mobile codecs
			"--disable-muxer=amr",                                                         // AMR mobile codecs
			"--disable-parser=amr",                                                        // AMR mobile codecs
			"--disable-decoder=libopencore_amrnb,libopencore_amrwb",                       // AMR mobile codecs
			"--disable-encoder=adpcm_g722,adpcm_g726,adpcm_g726le,g723_1",                 // Telecom encoders
			"--disable-decoder=adpcm_g722,adpcm_g726,adpcm_g726le,g723_1,g728,g729,qcelp", // Telecom decoders
			"--disable-parser=g723_1,g729",                                                // Telecom parsers
			"--disable-muxer=g722,g723_1,g726,g726le",                                     // Telecom muxers
			"--disable-demuxer=g722,g723_1,g726,g726le,g728,g729",                         // Telecom demuxers
			"--disable-decoder=truespeech",                                                // DSP Group TrueSpeech (1990s)
			"--disable-encoder=a64multi,a64multi5",                                        // Commodore 64
			"--disable-muxer=a64",                                                         // Commodore 64
			"--enable-pic",
			"--enable-gpl",
			"--enable-version3",
			"--enable-static",
			"--enable-librav1e",
			"--enable-libass",
			"--enable-libfreetype",
			"--enable-libfribidi",
			"--enable-libharfbuzz",
			"--enable-libmp3lame",
			"--enable-libopus",
			"--enable-libspeex",
			"--enable-libtheora",
			"--enable-libvpx",
			"--enable-libx264",
			"--enable-libx265",
			"--enable-libdav1d",
		)

		if b.os == Linux {
			cmd.Args = append(
				cmd.Args,
				"--enable-libfontconfig",
				// NVENC/NVDEC: Hardware accelerated encode/decode for NVIDIA GPUs
				// These flags enable NVENC/NVDEC support using nv-codec-headers
				// No CUDA runtime required - works with NVIDIA drivers only
				"--enable-ffnvcodec", // Enable NVENC/NVDEC support
				"--enable-nvdec",     // Enable NVDEC hardware decoder
				"--enable-nvenc",     // Enable NVENC hardware encoder
			)
		} else if b.os == MacOS {
			cmd.Args = append(
				cmd.Args,
				"--enable-avfoundation",
				"--enable-audiotoolbox",
				"--enable-videotoolbox",
			)
		}

		// Set PKG_CONFIG_PATH to find all pkg-config files (both lib and lib64)
		pkgConfigPath := fmt.Sprintf("%s/lib/pkgconfig:%s/lib64/pkgconfig", tgtDir, tgtDir)
		cmd.Env = append(cmd.Env, fmt.Sprintf("PKG_CONFIG_PATH=%s", pkgConfigPath))

		run("[ffmpeg configure]", cmd)
	}

	{
		log.Println("Running make")

		cmd := cmd(
			"make",
			buildPath,
			"-j8",
			"install",
		)

		run("[ffmpeg make]", cmd)
	}
}
