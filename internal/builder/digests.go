package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sort"
)

// archiveDigests pins the expected SHA-256 of every third-party source archive
// the builder downloads, keyed by download URL. Each compiled-in dependency must
// have an entry here or the build aborts (see verifyDigest); a changed URL has no
// pin and so fails closed rather than silently trusting an unverified download.
//
// SECURITY: these values are the supply-chain trust anchor for the shipped
// libffmpeg.a. Entries seeded by `--update-digests` are trust-on-first-use: their
// integrity is only as good as the single fetch that produced them. A maintainer
// MUST independently verify each digest against the upstream-published checksum
// before relying on it. Where upstream publishes official checksums (e.g. openssl,
// gnu.org), those values supersede any TOFU-seeded entry.
//
// To regenerate after a version or URL bump, run in a trusted environment:
//
//	go run ./internal/builder --update-digests
//
// then review the printed entries against upstream checksums and commit this file.
//
// Verification status (2026-06-05):
//   - openssl, opus, lame: independently confirmed against the upstream-published
//     SHA-256 (openssl .sha256 asset, xiph SHA256SUMS.txt, sourceforge published
//     value). Treat these as authoritative.
//   - libiconv 1.19: verified against the GNU-published GPG signature (good
//     signature from maintainer Bruno Haible) and cross-checked against the
//     Homebrew formula digest and a second GNU mirror. Treat as authoritative.
//   - All other entries: TOFU-seeded by `--update-digests` in an automated
//     environment. A maintainer should independently re-verify each before
//     relying on it. Several upstreams (gnome, videolan/gitlab on-demand
//     tarballs, github auto-generated archives) do not publish per-archive
//     checksums, so cross-checking may require a clean second fetch.
var archiveDigests = map[string]string{
	"https://bitbucket.org/multicoreware/x265_git/get/4.2.tar.bz2":                                                                                     "04978f795943e49fcea76eb5ede9c1bd0fe9b6c073518897be6fc43b44f60850",
	"https://code.videolan.org/videolan/dav1d/-/archive/1.5.3/dav1d-1.5.3.tar.bz2":                                                                     "e099f53253f6c247580c554d53a13f1040638f2066edc3c740e4c2f15174ce22",
	"https://code.videolan.org/videolan/x264/-/archive/0480cb05fa188d37ae87e8f4fd8f1aea3711f7ee/x264-0480cb05fa188d37ae87e8f4fd8f1aea3711f7ee.tar.bz2": "f05c59f2e83d494c36307025dca2d3afc6b4d185f3a3453d06cc4fecd7094057",
	"https://download.gnome.org/sources/libxml2/2.15/libxml2-2.15.3.tar.xz":                                                                            "78262a6e7ac170d6528ebfe2efccdf220191a5af6a6cd61ea4a9a9a5042c7a07",
	"https://downloads.sourceforge.net/project/lame/lame/3.100/lame-3.100.tar.gz":                                                                      "ddfe36cab873794038ae2c1210557ad34857a4b6bdc515785d1da9e175b1da1e",
	"https://downloads.xiph.org/releases/opus/opus-1.6.1.tar.gz":                                                                                       "6ffcb593207be92584df15b32466ed64bbec99109f007c82205f0194572411a1",
	"https://ftp.gnu.org/pub/gnu/libiconv/libiconv-1.19.tar.gz":                                                                                        "88dd96a8c0464eca144fc791ae60cd31cd8ee78321e67397e25fc095c4a19aa6",
	"https://github.com/FFmpeg/FFmpeg/archive/refs/tags/n8.1.2.tar.gz":                                                                                 "9fd092511605bbebafe095ea6d38d9e40f34d12f7386e1258372df8be0576eb7",
	"https://github.com/FFmpeg/nv-codec-headers/releases/download/n13.0.19.0/nv-codec-headers-13.0.19.0.tar.gz":                                        "13da39edb3a40ed9713ae390ca89faa2f1202c9dda869ef306a8d4383e242bee",
	"https://github.com/Haivision/srt/archive/refs/tags/v1.5.5.tar.gz":                                                                                 "c3518bc43a71b5289032395b2db4c3e09e73d78b54247d56c14553a503b491cf",
	"https://github.com/KhronosGroup/Vulkan-Headers/archive/refs/tags/v1.4.356.tar.gz":                                                                 "f0aac83f32b2895a15fb0686defc16755810e2705a3fd917cb9535ca79c71d4f",
	"https://github.com/KhronosGroup/glslang/archive/refs/tags/16.3.0.tar.gz":                                                                          "efff5a15258dce1ca2d323bf64c974f5fca03778174615dbc30c8d36db645bf5",
	"https://github.com/intel/libva/releases/download/2.24.0/libva-2.24.0.tar.bz2":                                                                     "56fab4e482dca2c9e8280d5057294b9faa789d637f97cc394a0c6ec08159060c",
	"https://github.com/intel/libvpl/archive/refs/tags/v2.17.0.tar.gz":                                                                                 "4de3e2faf1e8307fb282e4a43f443191810f6a6b0a484fffa7995ba1c814c6ec",
	"https://github.com/madler/zlib/releases/download/v1.3.2/zlib-1.3.2.tar.gz":                                                                        "bb329a0a2cd0274d05519d61c667c062e06990d72e125ee2dfa8de64f0119d16",
	"https://github.com/openssl/openssl/releases/download/openssl-3.6.3/openssl-3.6.3.tar.gz":                                                          "243a86649cf6f23eeb6a2ff2456e09e5d77dd9018a54d3d96b0c6bdd6ba6c7f1",
	"https://github.com/sekrit-twc/zimg/archive/refs/tags/release-3.0.6.tar.gz":                                                                        "be89390f13a5c9b2388ce0f44a5e89364a20c1c57ce46d382b1fcc3967057577",
	"https://github.com/webmproject/libvpx/archive/refs/tags/v1.16.0.tar.gz":                                                                           "7a479a3c66b9f5d5542a4c6a1b7d3768a983b1e5c14c60a9396edc9b649e015c",
	"https://github.com/xiph/rav1e/archive/refs/tags/v0.8.1.tar.gz":                                                                                    "06d1523955fb6ed9cf9992eace772121067cca7e8926988a1ee16492febbe01e",
	"https://gitlab.freedesktop.org/mesa/drm/-/archive/libdrm-2.4.134/drm-libdrm-2.4.134.tar.gz":                                                       "6b18e4834b0c061232cb5c11e98a6ecdc72ebc6bc282d124406b7a9d4e089ce2",
	"https://storage.googleapis.com/downloads.webmproject.org/releases/webp/libwebp-1.6.0.tar.gz":                                                      "e4ab7009bf0629fd11982d4c2aa83964cf244cffba7347ecd39019a9e38c4564",
}

// expectedDigest returns the pinned SHA-256 for url and whether a pin exists.
func expectedDigest(url string) (string, bool) {
	d, ok := archiveDigests[url]
	return d, ok
}

// hashFile streams path through SHA-256 and returns the lowercase hex digest.
// It never slurps the whole archive into memory.
func hashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// verifyDigest checks the archive at path against the pinned SHA-256 for url.
//
// Fail-closed contract:
//   - No pin for url: returns an actionable error; the file is left in place so a
//     maintainer can inspect it, but the build does not proceed.
//   - Digest mismatch: deletes the file (so a re-run re-downloads rather than
//     trusting poison) and returns a descriptive error.
//   - Match: returns nil.
func verifyDigest(url, path string) error {
	expected, ok := expectedDigest(url)
	if !ok {
		return fmt.Errorf("no pinned SHA-256 for %s; run 'go run ./internal/builder --update-digests' in a trusted environment and commit internal/builder/digests.go", url)
	}

	actual, err := hashFile(path)
	if err != nil {
		return fmt.Errorf("hashing %s: %w", path, err)
	}

	if actual != expected {
		os.Remove(path)
		return fmt.Errorf("SHA-256 mismatch for %s: expected %s, got %s (deleted %s; re-run to re-download)", url, expected, actual, path)
	}

	return nil
}

// formatDigestEntries renders a sorted Go map-literal body for the given
// url→digest pairs, suitable for pasting into archiveDigests.
func formatDigestEntries(digests map[string]string) string {
	urls := make([]string, 0, len(digests))
	for url := range digests {
		urls = append(urls, url)
	}
	sort.Strings(urls)

	var b []byte
	for _, url := range urls {
		b = append(b, fmt.Sprintf("\t%q: %q,\n", url, digests[url])...)
	}
	return string(b)
}
