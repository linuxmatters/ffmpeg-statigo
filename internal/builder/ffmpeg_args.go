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
	}
}
