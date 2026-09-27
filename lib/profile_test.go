package lib

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go/build"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

type testTransport func(*http.Request) (*http.Response, error)

func (f testTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func mockHTTP(t *testing.T, handler testTransport) {
	t.Helper()
	previous := http.DefaultClient
	http.DefaultClient = &http.Client{Transport: handler}
	t.Cleanup(func() { http.DefaultClient = previous })
}

func testHTTPResponse(status int, body []byte) *http.Response {
	return &http.Response{
		StatusCode:    status,
		Status:        fmt.Sprintf("%d %s", status, http.StatusText(status)),
		Body:          io.NopCloser(bytes.NewReader(body)),
		ContentLength: int64(len(body)),
	}
}

func TestDownloadLibsProfile(t *testing.T) {
	for _, platform := range []string{"linux", "darwin"} {
		for _, arch := range []string{"amd64", "arm64"} {
			for _, mode := range []string{"digest", "checksum_file", "mismatch", "missing_checksum", "wrong_profile", "download_failure", "cached"} {
				t.Run(platform+"_"+arch+"/"+mode, func(t *testing.T) {
					t.Chdir(t.TempDir())
					t.Setenv("GOOS", platform)
					t.Setenv("GOARCH", arch)
					t.Setenv("GITHUB_TOKEN", "")
					t.Setenv("GH_TOKEN", "")
					platArch := platform + "_" + arch
					payloadPath := filepath.Join(libraryProfile, platArch, "libffmpeg.a")
					asset := "ffmpeg-" + platform + "-" + arch + ".tar.gz"
					otherProfile := "embedded"
					if libraryProfile == "embedded" {
						asset = "ffmpeg-embedded-" + platform + "-" + arch + ".tar.gz"
						otherProfile = ""
					}
					otherPath := filepath.Join("lib", otherProfile, platArch, "libffmpeg.a")
					if err := os.MkdirAll(filepath.Dir(otherPath), 0o755); err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(otherPath, []byte("other profile"), 0o644); err != nil {
						t.Fatal(err)
					}
					libPath := filepath.Join("lib", payloadPath)
					if mode == "cached" {
						if err := os.MkdirAll(filepath.Dir(libPath), 0o755); err != nil {
							t.Fatal(err)
						}
						if err := os.WriteFile(libPath, []byte("selected profile"), 0o644); err != nil {
							t.Fatal(err)
						}
					}
					if mode == "wrong_profile" {
						payloadPath = filepath.Join(otherProfile, platArch, "libffmpeg.a")
					}
					tarball := createValidTarball(t, map[string][]byte{filepath.ToSlash(payloadPath): []byte("selected profile")})
					t.Cleanup(func() { _ = os.Remove(tarball) })
					payload, err := os.ReadFile(tarball)
					if err != nil {
						t.Fatal(err)
					}
					sum := sha256.Sum256(payload)
					digest := hex.EncodeToString(sum[:])
					if mode == "mismatch" {
						digest = strings.Repeat("0", 64)
					}
					const release = "lib-8.1.2.0"
					requests := 0
					mockHTTP(t, func(req *http.Request) (*http.Response, error) {
						requests++
						var data any
						switch req.URL.Path {
						case "/repos/linuxmatters/ffmpeg-statigo/releases":
							data = []GitHubRelease{{TagName: release}}
						case "/repos/linuxmatters/ffmpeg-statigo/releases/tags/" + release:
							assets := []GitHubAsset{{Name: asset, Digest: "sha256:" + digest}}
							switch mode {
							case "checksum_file":
								assets = []GitHubAsset{{Name: asset}, {Name: "SHA256SUMS", BrowserDownloadURL: "https://github.com/checksums"}}
							case "missing_checksum":
								assets = []GitHubAsset{{Name: asset}}
							}
							data = GitHubReleaseDetail{Assets: assets}
						case "/checksums":
							return testHTTPResponse(http.StatusOK, []byte("wrong  "+asset+".bak\n"+digest+" *"+asset+"\n")), nil
						case "/linuxmatters/ffmpeg-statigo/releases/download/" + release + "/" + asset:
							if mode == "download_failure" {
								return testHTTPResponse(http.StatusNotFound, nil), nil
							}
							return testHTTPResponse(http.StatusOK, payload), nil
						default:
							t.Fatalf("unexpected request: %s", req.URL)
						}
						body, err := json.Marshal(data)
						if err != nil {
							t.Fatal(err)
						}
						return testHTTPResponse(http.StatusOK, body), nil
					})
					err = DownloadLibs()
					wantError := map[string]string{
						"mismatch": "checksum mismatch", "missing_checksum": "no checksum",
						"wrong_profile": "missing library", "download_failure": "404",
					}[mode]
					if wantError != "" {
						if err == nil || !strings.Contains(err.Error(), wantError) {
							t.Fatalf("DownloadLibs() = %v, want %q", err, wantError)
						}
						if _, err := os.Stat(libPath); !os.IsNotExist(err) {
							t.Fatalf("failed download left a library: %v", err)
						}
					} else {
						if err != nil {
							t.Fatal(err)
						}
						got, err := os.ReadFile(libPath)
						if err != nil || string(got) != "selected profile" {
							t.Fatalf("installed library = %q, %v", got, err)
						}
					}
					if mode == "cached" && requests != 0 {
						t.Errorf("cached library made %d requests", requests)
					}
					other, err := os.ReadFile(otherPath)
					if err != nil || string(other) != "other profile" {
						t.Fatalf("other profile changed: %q, %v", other, err)
					}
					staging, err := filepath.Glob("lib/.download-*")
					if err != nil || len(staging) != 0 {
						t.Fatalf("download directories remain: %v, %v", staging, err)
					}
				})
			}
		}
	}
}

func TestProfileBuildSelection(t *testing.T) {
	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}
	for _, platform := range []string{"linux", "darwin"} {
		for _, arch := range []string{"amd64", "arm64"} {
			for _, profile := range []string{"default", "embedded"} {
				t.Run(platform+"_"+arch+"/"+profile, func(t *testing.T) {
					ctx := build.Default
					ctx.GOOS, ctx.GOARCH, ctx.CgoEnabled = platform, arch, true
					ctx.BuildTags = nil
					if profile == "embedded" {
						ctx.BuildTags = []string{"embedded"}
					}
					pkg, err := ctx.ImportDir(root, 0)
					if err != nil {
						t.Fatal(err)
					}
					if !slices.Contains(pkg.CgoFiles, "ffmpeg_link_"+profile+".go") {
						t.Fatalf("wrong linker file selection: %v", pkg.CgoFiles)
					}
					var want []string
					if profile == "embedded" {
						want = []string{filepath.Join(root, "lib", "embedded", platform+"_"+arch, "libffmpeg.a"), "-lm"}
						if platform == "linux" {
							want = append(want, "-ldl", "-lstdc++")
						} else {
							want = append(want, "-lc++")
						}
						want = append(want, "-lpthread")
					} else {
						want = []string{"-L" + filepath.Join(root, "lib", platform+"_"+arch), "-lffmpeg"}
						if platform == "linux" {
							want = append(want, "-lm", "-ldl", "-lstdc++", "-lpthread")
						} else {
							want = append(want, "-lstdc++", "-lm", "-framework", "ApplicationServices", "-framework", "CoreVideo", "-framework", "CoreMedia", "-framework", "VideoToolbox", "-framework", "AudioToolbox")
						}
					}
					if !slices.Equal(pkg.CgoLDFLAGS, want) {
						t.Errorf("link flags = %v, want %v", pkg.CgoLDFLAGS, want)
					}
					pkg, err = ctx.ImportDir(filepath.Join(root, "lib"), 0)
					if err != nil {
						t.Fatal(err)
					}
					if !slices.Contains(pkg.GoFiles, "profile_"+profile+".go") {
						t.Fatalf("wrong downloader profile selection: %v", pkg.GoFiles)
					}
				})
			}
		}
	}
}
