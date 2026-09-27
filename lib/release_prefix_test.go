package lib

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestReleasePrefixPatchBoundary(t *testing.T) {
	for _, tc := range []struct {
		name     string
		tags     []string
		want     string
		wantMiss bool
	}{
		{
			name: "reject_later_patch_with_compatible_release",
			tags: []string{"lib-8.1.20.0", "lib-8.1.2.0"},
			want: "lib-8.1.2.0",
		},
		{
			name:     "reject_later_patch_without_compatible_release",
			tags:     []string{"lib-8.1.20.0", "lib-8.1.21.3", "lib-8.1.200.0"},
			want:     "lib-8.1.2.0",
			wantMiss: true,
		},
		{
			name: "preserve_lexicographic_build_sort",
			tags: []string{"lib-8.1.2.2", "lib-8.1.20.0", "lib-8.1.2.10", "lib-8.1.2.0"},
			want: "lib-8.1.2.2",
		},
		{
			name:     "reject_tags_without_build_separator",
			tags:     []string{"lib-8.1.2", "lib-8.1.2-rc1", "v8.1.2.0"},
			want:     "lib-8.1.2.0",
			wantMiss: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			releases := make([]GitHubRelease, len(tc.tags))
			for i, tag := range tc.tags {
				releases[i] = GitHubRelease{TagName: tag}
			}
			body, err := json.Marshal(releases)
			if err != nil {
				t.Fatal(err)
			}
			mockHTTP(t, func(req *http.Request) (*http.Response, error) {
				if req.URL.Host != "api.github.com" || req.URL.Path != "/repos/linuxmatters/ffmpeg-statigo/releases" {
					t.Fatalf("unexpected request: %s", req.URL)
				}
				return testHTTPResponse(http.StatusOK, body), nil
			})

			got, err := findViaAPI("lib-8.1.2")
			if tc.wantMiss {
				if err == nil || !strings.Contains(err.Error(), "no releases found matching") || got != "" {
					t.Fatalf("findViaAPI() = %q, %v, want no matching release", got, err)
				}
			} else if err != nil || got != tc.want {
				t.Fatalf("findViaAPI() = %q, %v, want %q, nil", got, err, tc.want)
			}

			got, err = findCompatibleRelease("8.1.2")
			if err != nil || got != tc.want {
				t.Fatalf("findCompatibleRelease() = %q, %v, want %q, nil", got, err, tc.want)
			}
		})
	}
}
