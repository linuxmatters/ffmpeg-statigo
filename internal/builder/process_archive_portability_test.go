package main

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestPortabilityConfigureCommand(t *testing.T) {
	for _, targetOS := range []string{"linux", "darwin", "windows"} {
		for _, script := range []struct{ interpreter, path string }{{"sh", "./configure"}, {"perl", "./Configure"}} {
			args := []string{"--prefix=C:/build directory/staging", "--enable-static"}
			name, got := configureCommand(targetOS, script.interpreter, script.path, args)
			wantName, wantArgs := script.path, args
			if targetOS == "windows" {
				wantName = script.interpreter
				wantArgs = append([]string{script.path}, args...)
			}
			if name != wantName || !slices.Equal(got, wantArgs) {
				t.Fatalf("%s: configureCommand = %q %q, want %q %q", targetOS, name, got, wantName, wantArgs)
			}
			if args[0] != "--prefix=C:/build directory/staging" {
				t.Fatal("configureCommand changed the input arguments")
			}
		}
	}
}

func TestPortabilityBuildToolPath(t *testing.T) {
	for _, tc := range []struct{ targetOS, input, want string }{
		{"windows", `C:\build directory\staging`, "C:/build directory/staging"},
		{"windows", `\\server\share\staging`, "//server/share/staging"},
		{"windows", "/ucrt64/bin", "/ucrt64/bin"},
		{"linux", `/build/back\slash`, `/build/back\slash`},
		{"darwin", "/build directory/staging", "/build directory/staging"},
	} {
		if got := buildToolPath(tc.input, tc.targetOS); got != tc.want {
			t.Errorf("buildToolPath(%q, %s) = %q, want %q", tc.input, tc.targetOS, got, tc.want)
		}
	}
}

func TestPortabilityBuildPathArgs(t *testing.T) {
	args := []string{`--extra-cflags=-IC:\build directory\staging\include`, `PREFIX=C:\build directory\staging`, `--pattern=\d+`}
	for _, targetOS := range []string{"linux", "darwin", "windows"} {
		for _, stage := range []string{`C:\build directory\staging`, "C:/build directory/staging"} {
			got := buildPathArgs(args, stage, targetOS)
			want := slices.Clone(args)
			if targetOS == "windows" {
				want[0] = "--extra-cflags=-IC:/build directory/staging/include"
				want[1] = "PREFIX=C:/build directory/staging"
			}
			if !slices.Equal(got, want) {
				t.Errorf("%s buildPathArgs = %q, want %q", targetOS, got, want)
			}
		}
	}
	if args[1] != `PREFIX=C:\build directory\staging` {
		t.Fatal("buildPathArgs changed the input arguments")
	}
}

func TestPortabilityBuildPathEnv(t *testing.T) {
	for _, targetOS := range []string{"linux", "darwin", "windows"} {
		stage, pathKey, separator := "/build directory/staging", "PATH", ":"
		if targetOS == "windows" {
			stage, pathKey, separator = `C:\build directory\staging`, "Path", ";"
		}
		for _, existing := range []bool{false, true} {
			env := []string{"UNCHANGED=value"}
			if existing {
				env = append(env, pathKey+"=old-bin", "PKG_CONFIG_PATH=old-pkg")
			}
			got := buildPathEnv(env, stage, targetOS)
			want := []string{
				"PATH=" + buildToolPath(filepath.Join(stage, "bin"), targetOS),
				"PKG_CONFIG_PATH=" + buildToolPath(filepath.Join(stage, "lib", "pkgconfig"), targetOS),
			}
			if existing {
				want[0] += separator + "old-bin"
				want[1] += separator + "old-pkg"
			}
			for _, entry := range append(want, "UNCHANGED=value") {
				if !slices.Contains(got, entry) {
					t.Errorf("%s buildPathEnv = %q, missing %q", targetOS, got, entry)
				}
			}
			if len(got) != 3 {
				t.Errorf("%s buildPathEnv has duplicate entries: %q", targetOS, got)
			}
		}
	}
}

func TestPortabilityConfiguredArchiveTool(t *testing.T) {
	for _, key := range []string{"AR", "STRIP"} {
		t.Setenv(key, "")
		if got := configuredArchiveTool(key, "fallback"); got != "fallback" {
			t.Fatalf("unset %s = %q", key, got)
		}
		tool := `C:\tool directory\x86_64-w64-mingw32-` + strings.ToLower(key) + ".exe"
		t.Setenv(key, tool)
		if got := configuredArchiveTool(key, "fallback"); got != tool {
			t.Fatalf("configured %s = %q, want %q", key, got, tool)
		}
	}
}

func TestPortabilityPrepareArchiveMerge(t *testing.T) {
	root := t.TempDir()
	var sources []string
	for _, directory := range []string{"first directory", "second directory"} {
		dir := filepath.Join(root, directory)
		if err := os.Mkdir(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(dir, "same name.a")
		if err := os.WriteFile(path, []byte("!<arch>\n"+directory), 0o644); err != nil {
			t.Fatal(err)
		}
		sources = append(sources, path)
	}
	work := filepath.Join(root, "merge directory")
	if err := os.Mkdir(work, 0o755); err != nil {
		t.Fatal(err)
	}
	script, err := prepareArchiveMerge(sources, work)
	if err != nil {
		t.Fatal(err)
	}
	want := "create combined.a\naddlib input-000000.a\naddlib input-000001.a\nsave\nend\n"
	if script != want {
		t.Fatalf("MRI script = %q, want %q", script, want)
	}
	for i, alias := range []string{"input-000000.a", "input-000001.a"} {
		got, err := os.ReadFile(filepath.Join(work, alias))
		if err != nil {
			t.Fatal(err)
		}
		original, err := os.ReadFile(sources[i])
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != string(original) {
			t.Fatalf("copied archive %s differs", alias)
		}
	}
}

func TestPortabilityMergeRejectsInvalidInputWithoutReplacingOutput(t *testing.T) {
	for _, contents := range []string{"!<thin>\n", "not an archive", ""} {
		t.Run(contents, func(t *testing.T) {
			root := t.TempDir()
			input, output := filepath.Join(root, "input.a"), filepath.Join(root, "output.a")
			for path, data := range map[string]string{input: contents, output: "original archive"} {
				if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			if err := combineWindows(context.Background(), []string{input}, output); err == nil {
				t.Fatal("invalid archive was accepted")
			}
			got, err := os.ReadFile(output)
			if err != nil || string(got) != "original archive" {
				t.Fatalf("existing output changed: %q, %v", got, err)
			}
			entries, err := os.ReadDir(root)
			if err != nil || len(entries) != 2 {
				t.Fatalf("temporary merge directory remains: %v, %v", entries, err)
			}
		})
	}
}
