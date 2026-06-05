package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// withPinnedDigest temporarily adds url→digest to archiveDigests and restores the
// prior state when the test finishes.
func withPinnedDigest(t *testing.T, url, digest string) {
	t.Helper()
	prev, had := archiveDigests[url]
	archiveDigests[url] = digest
	t.Cleanup(func() {
		if had {
			archiveDigests[url] = prev
		} else {
			delete(archiveDigests, url)
		}
	})
}

// writeFixture writes content to a file in a fresh temp dir and returns its path.
func writeFixture(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "archive.bin")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return path
}

// SHA-256 of the byte string "hello".
const helloDigest = "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"

func TestVerifyDigestMatchPasses(t *testing.T) {
	const url = "https://example.test/hello.tar.gz"
	withPinnedDigest(t, url, helloDigest)

	path := writeFixture(t, "hello")
	if err := verifyDigest(url, path); err != nil {
		t.Fatalf("verifyDigest = %v, want nil", err)
	}
	if !fileExists(path) {
		t.Fatal("matching archive was removed; want preserved")
	}
}

func TestVerifyDigestMismatchRemovesFile(t *testing.T) {
	const url = "https://example.test/hello.tar.gz"
	withPinnedDigest(t, url, helloDigest)

	path := writeFixture(t, "poisoned")
	err := verifyDigest(url, path)
	if err == nil {
		t.Fatal("expected mismatch error, got nil")
	}
	if !strings.Contains(err.Error(), "SHA-256 mismatch") {
		t.Fatalf("error = %v, want SHA-256 mismatch", err)
	}
	if fileExists(path) {
		t.Fatal("mismatched archive was not removed; a re-run would trust poison")
	}
}

func TestVerifyDigestMissingPinFailsClosed(t *testing.T) {
	const url = "https://example.test/unpinned.tar.gz"
	if _, ok := expectedDigest(url); ok {
		t.Fatalf("test url unexpectedly pinned: %s", url)
	}

	path := writeFixture(t, "anything")
	err := verifyDigest(url, path)
	if err == nil {
		t.Fatal("expected fail-closed error for missing pin, got nil")
	}
	if !strings.Contains(err.Error(), "no pinned SHA-256") {
		t.Fatalf("error = %v, want no pinned SHA-256", err)
	}
	if !strings.Contains(err.Error(), "--update-digests") {
		t.Fatalf("error = %v, want actionable --update-digests hint", err)
	}
	// A missing pin leaves the file in place for inspection.
	if !fileExists(path) {
		t.Fatal("file removed on missing pin; want preserved for inspection")
	}
}

func TestDownloadFileCacheHitVerifies(t *testing.T) {
	const url = "https://example.test/cached.tar.gz"
	withPinnedDigest(t, url, helloDigest)

	dir := t.TempDir()
	dest := filepath.Join(dir, "downloads", "cached.tar.gz")
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	// Seed the download cache with mismatched (poisoned) contents.
	if err := os.WriteFile(dest, []byte("poisoned"), 0o644); err != nil {
		t.Fatalf("seed cache: %v", err)
	}

	err := DownloadFile(url, dest, io.Discard)
	if err == nil {
		t.Fatal("expected cached poisoned archive to be rejected, got nil")
	}
	if !strings.Contains(err.Error(), "SHA-256 mismatch") {
		t.Fatalf("error = %v, want SHA-256 mismatch", err)
	}
	if fileExists(dest) {
		t.Fatal("poisoned cached archive was not removed")
	}
}

func TestFormatDigestEntriesSorted(t *testing.T) {
	got := formatDigestEntries(map[string]string{
		"https://b.test/x.tgz": "bb",
		"https://a.test/x.tgz": "aa",
	})
	want := "\t\"https://a.test/x.tgz\": \"aa\",\n\t\"https://b.test/x.tgz\": \"bb\",\n"
	if got != want {
		t.Fatalf("formatDigestEntries =\n%q\nwant\n%q", got, want)
	}
}
