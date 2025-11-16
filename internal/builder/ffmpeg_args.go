package main

// FFmpegArgsCommon returns common FFmpeg configure arguments for all platforms
func FFmpegArgsCommon() []string {
	return []string{
		// Build options
		"--enable-pic",
		"--enable-gpl",
		"--enable-version3",
		"--enable-static",

		// Disable everything first
		"--disable-autodetect",
		"--disable-debug",
		"--disable-doc",
		"--disable-htmlpages",
		"--disable-manpages",
		"--disable-podpages",
		"--disable-programs",
		"--disable-txtpages",
		"--disable-everything",

		// Enable all filters, hwaccels and protocols
		"--enable-filters",
		"--enable-hwaccels",
		"--enable-protocols",

		// TEMPORARILY Enable all muxers, demuxers, parsers, bsfs, encoders, decoders, hwaccels
		"--enable-muxers",
		"--enable-demuxers",
		"--enable-parsers",
		"--enable-bsfs",
		"--enable-encoders",
		"--enable-decoders",

		// External libraries
		"--enable-libdav1d",
		"--enable-libglslang",
		"--enable-libmp3lame",
		"--enable-libopus",
		"--enable-librav1e",
		"--enable-libsrt",
		"--enable-libvpx",
		"--enable-libwebp",
		"--enable-libx264",
		"--enable-libx265",
		"--enable-libzimg",
		"--enable-openssl",
		"--enable-vulkan",
		"--enable-zlib",
	}
}

// FFmpegArgsLinux returns Linux-specific FFmpeg configure arguments
func FFmpegArgsLinux() []string {
	return []string{
		"--enable-cuvid",
		"--enable-ffnvcodec",
		"--enable-nvdec",
		"--enable-nvenc",
		"--enable-libvpl",
	}
}

// FFmpegArgsDarwin returns macOS-specific FFmpeg configure arguments
func FFmpegArgsDarwin() []string {
	return []string{
		"--enable-avfoundation",
		"--enable-audiotoolbox",
		"--enable-videotoolbox",
	}
}
