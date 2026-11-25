package main

import (
	"bytes"
	"context"
	"debug/pe"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	if os.Getenv("FFMPEG_BUILDER_ARCHIVE_FIXTURE") != "1" {
		os.Exit(m.Run())
	}
	if slices.Equal(os.Args[1:], []string{"-M"}) {
		script, err := io.ReadAll(os.Stdin)
		if err != nil || string(script) != "create combined.a\naddlib input-000000.a\nsave\nend\n" {
			os.Exit(2)
		}
		if err := copyMergeArchive("input-000000.a", "combined.a"); err != nil {
			os.Exit(2)
		}
		os.Exit(0)
	}
	if slices.Equal(os.Args[1:], []string{"--strip-unneeded", "--remove-section=.llvmbc", "--remove-section=.llvmcmd", "combined.a"}) {
		f, err := os.Open("combined.a")
		if err != nil {
			os.Exit(3)
		}
		_, pruned, scanErr := scanArchiveMembers(f)
		_ = f.Close()
		if scanErr != nil || pruned != 0 {
			os.Exit(4)
		}
		if err := os.WriteFile(filepath.Join("..", "strip-called"), []byte("called"), 0o644); err != nil {
			os.Exit(2)
		}
		if os.Getenv("FFMPEG_BUILDER_STRIP_FAIL") == "1" {
			os.Exit(1)
		}
		os.Exit(0)
	}
	if slices.Equal(os.Args[1:], []string{"s", "combined.a"}) {
		if !fileExists(filepath.Join("..", "strip-called")) {
			os.Exit(5)
		}
		if err := os.WriteFile(filepath.Join("..", "ar-index-called"), []byte("called"), 0o644); err != nil {
			os.Exit(2)
		}
		if os.Getenv("FFMPEG_BUILDER_AR_INDEX_FAIL") == "1" {
			os.Exit(1)
		}
		os.Exit(0)
	}
	os.Exit(2)
}

type testARMember struct {
	name    string
	payload []byte
}

func buildTestArchive(members ...testARMember) []byte {
	var buf bytes.Buffer
	buf.WriteString(arMagic)
	for _, m := range members {
		if len(m.name) > 16 {
			panic("archive member name too long: " + m.name)
		}
		fmt.Fprintf(&buf, "%-16s%-12s%-6s%-6s%-8s%-10d`\n", m.name, "0", "0", "0", "100644", len(m.payload))
		buf.Write(m.payload)
		if len(m.payload)&1 != 0 {
			buf.WriteByte('\n')
		}
	}
	return buf.Bytes()
}

func parseTestArchive(t *testing.T, data []byte) []testARMember {
	t.Helper()
	if len(data) < len(arMagic) || string(data[:len(arMagic)]) != arMagic {
		t.Fatalf("invalid archive magic: %q", data)
	}
	var members []testARMember
	off := len(arMagic)
	for off < len(data) {
		if off+arHeaderSize > len(data) {
			t.Fatalf("truncated header at offset %d", off)
		}
		hdr := data[off : off+arHeaderSize]
		if hdr[58] != '`' || hdr[59] != '\n' {
			t.Fatalf("malformed header trailer at offset %d: %q", off, hdr[58:60])
		}
		name := strings.TrimSpace(string(hdr[0:16]))
		size, err := strconv.Atoi(strings.TrimSpace(string(hdr[48:58])))
		if err != nil || size < 0 || off+arHeaderSize+size > len(data) {
			t.Fatalf("invalid member size at offset %d: %v", off, err)
		}
		payload := bytes.Clone(data[off+arHeaderSize : off+arHeaderSize+size])
		members = append(members, testARMember{name: name, payload: payload})
		off += arHeaderSize + size
		if size&1 != 0 {
			if off >= len(data) || data[off] != '\n' {
				t.Fatalf("missing odd-size pad at offset %d", off)
			}
			off++
		}
	}
	return members
}

func makeCOFFObject(machine, numSections uint16, totalLen int) []byte {
	buf := make([]byte, totalLen)
	binary.LittleEndian.PutUint16(buf[0:2], machine)
	binary.LittleEndian.PutUint16(buf[2:4], numSections)
	for i := coffHeaderSize; i < len(buf); i++ {
		buf[i] = byte(i)
	}
	return buf
}

func makeCOFFBigObj(machine, version uint16, numSections uint32, classID [16]byte) []byte {
	buf := make([]byte, coffBigObjHeaderSize)
	binary.LittleEndian.PutUint16(buf[0:2], 0x0000)
	binary.LittleEndian.PutUint16(buf[2:4], 0xffff)
	binary.LittleEndian.PutUint16(buf[4:6], version)
	binary.LittleEndian.PutUint16(buf[6:8], machine)
	copy(buf[12:28], classID[:])
	binary.LittleEndian.PutUint32(buf[44:48], numSections)
	return buf
}

func makePEShortImportObject(machine uint16) []byte {
	const stringsData = "av_version_info\x00avutil.dll\x00"
	buf := make([]byte, 20+len(stringsData))
	binary.LittleEndian.PutUint16(buf[0:2], 0x0000)
	binary.LittleEndian.PutUint16(buf[2:4], 0xffff)
	binary.LittleEndian.PutUint16(buf[4:6], 0)
	binary.LittleEndian.PutUint16(buf[6:8], machine)
	binary.LittleEndian.PutUint32(buf[12:16], uint32(len(stringsData)))
	copy(buf[20:], stringsData)
	return buf
}

func TestPortabilityWindowsMergeStrip(t *testing.T) {
	tool, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name            string
		machine         uint16
		stripFails      bool
		indexFails      bool
		wantErrContains string
		wantStrip       bool
		wantIndex       bool
	}{
		{name: "amd64", machine: pe.IMAGE_FILE_MACHINE_AMD64, wantStrip: true, wantIndex: true},
		{name: "386", machine: pe.IMAGE_FILE_MACHINE_I386, wantStrip: true, wantIndex: true},
		{name: "strip-failure", machine: pe.IMAGE_FILE_MACHINE_AMD64, stripFails: true, wantErrContains: "strip failed", wantStrip: true, wantIndex: false},
		{name: "ar-index-failure", machine: pe.IMAGE_FILE_MACHINE_I386, indexFails: true, wantErrContains: "archive index failed", wantStrip: true, wantIndex: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			input, output := filepath.Join(root, "input.a"), filepath.Join(root, "output.a")
			stripMarker := filepath.Join(root, "strip-called")
			indexMarker := filepath.Join(root, "ar-index-called")

			longNames := []byte("very_long_zero_section_member_name.o/\n")
			validObj := makeCOFFObject(tc.machine, 1, 61)
			zeroObj := makeCOFFObject(tc.machine, 0, coffHeaderSize)
			importObj := makePEShortImportObject(tc.machine)

			archiveBytes := buildTestArchive(
				testARMember{name: "/", payload: []byte{0, 0, 0, 0}},
				testARMember{name: "//", payload: longNames},
				testARMember{name: "valid.o/", payload: validObj},
				testARMember{name: "/0", payload: zeroObj},
				testARMember{name: "imp.o/", payload: importObj},
			)

			if err := os.WriteFile(input, archiveBytes, 0o644); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(output, []byte("original archive"), 0o644); err != nil {
				t.Fatal(err)
			}

			t.Setenv("AR", tool)
			t.Setenv("STRIP", tool)
			t.Setenv("FFMPEG_BUILDER_ARCHIVE_FIXTURE", "1")
			t.Setenv("FFMPEG_BUILDER_STRIP_FAIL", "0")
			t.Setenv("FFMPEG_BUILDER_AR_INDEX_FAIL", "0")
			if tc.stripFails {
				t.Setenv("FFMPEG_BUILDER_STRIP_FAIL", "1")
			}
			if tc.indexFails {
				t.Setenv("FFMPEG_BUILDER_AR_INDEX_FAIL", "1")
			}

			err := combineWindows(context.Background(), []string{input}, output)
			if tc.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErrContains) {
					t.Fatalf("merge error = %v, want substring %q", err, tc.wantErrContains)
				}
			} else if err != nil {
				t.Fatal(err)
			}

			if got := fileExists(stripMarker); got != tc.wantStrip {
				t.Fatalf("strip invoked = %v, want %v", got, tc.wantStrip)
			}
			if got := fileExists(indexMarker); got != tc.wantIndex {
				t.Fatalf("ar index invoked = %v, want %v", got, tc.wantIndex)
			}

			got, err := os.ReadFile(output)
			if err != nil {
				t.Fatal(err)
			}
			if tc.wantErrContains != "" {
				if string(got) != "original archive" {
					t.Fatalf("output changed on error: %q", got)
				}
			} else {
				members := parseTestArchive(t, got)
				wantMembers := []testARMember{
					{name: "//", payload: longNames},
					{name: "valid.o/", payload: validObj},
					{name: "imp.o/", payload: importObj},
				}
				if len(members) != len(wantMembers) {
					t.Fatalf("member count = %d, want %d", len(members), len(wantMembers))
				}
				for i, want := range wantMembers {
					if members[i].name != want.name || !bytes.Equal(members[i].payload, want.payload) {
						t.Fatalf("member[%d] = (%q, %d bytes), want (%q, %d bytes)", i, members[i].name, len(members[i].payload), want.name, len(want.payload))
					}
				}
			}

			entries, err := filepath.Glob(filepath.Join(root, ".ffmpeg-merge-*"))
			if err != nil || len(entries) != 0 {
				t.Fatalf("temporary merge directory remains: %v, %v", entries, err)
			}
		})
	}
}

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

func TestPruneZeroSectionCOFFMembers(t *testing.T) {
	t.Run("empty-archive", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "empty.a")
		if err := os.WriteFile(path, []byte(arMagic), 0o644); err != nil {
			t.Fatal(err)
		}
		pruned, err := pruneZeroSectionCOFFMembers(path)
		if err != nil || pruned != 0 {
			t.Fatalf("pruneZeroSectionCOFFMembers = (%d, %v), want (0, nil)", pruned, err)
		}
		got, err := os.ReadFile(path)
		if err != nil || string(got) != arMagic {
			t.Fatalf("empty archive mutated: %q, %v", got, err)
		}
	})

	t.Run("noop-keeps-symbol-table", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "noop.a")
		original := buildTestArchive(
			testARMember{name: "/", payload: []byte{0, 0, 0, 0}},
			testARMember{name: "valid.o/", payload: makeCOFFObject(pe.IMAGE_FILE_MACHINE_AMD64, 1, 40)},
		)
		if err := os.WriteFile(path, original, 0o644); err != nil {
			t.Fatal(err)
		}
		pruned, err := pruneZeroSectionCOFFMembers(path)
		if err != nil || pruned != 0 {
			t.Fatalf("pruneZeroSectionCOFFMembers = (%d, %v), want (0, nil)", pruned, err)
		}
		got, err := os.ReadFile(path)
		if err != nil || !bytes.Equal(got, original) {
			t.Fatalf("no-op archive changed: got %d bytes, want %d bytes, err=%v", len(got), len(original), err)
		}
	})

	t.Run("bigobj-zero-section", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "bigobj.a")
		validObj := makeCOFFObject(pe.IMAGE_FILE_MACHINE_AMD64, 2, 41)
		zeroBigObj := makeCOFFBigObj(pe.IMAGE_FILE_MACHINE_AMD64, 2, 0, coffBigObjClassID)
		validBigObj := makeCOFFBigObj(pe.IMAGE_FILE_MACHINE_AMD64, 2, 3, coffBigObjClassID)
		archive := buildTestArchive(
			testARMember{name: "/", payload: []byte{0, 0, 0, 0}},
			testARMember{name: "valid.o/", payload: validObj},
			testARMember{name: "zero_big.o/", payload: zeroBigObj},
			testARMember{name: "valid_big.o/", payload: validBigObj},
		)
		if err := os.WriteFile(path, archive, 0o644); err != nil {
			t.Fatal(err)
		}
		pruned, err := pruneZeroSectionCOFFMembers(path)
		if err != nil || pruned != 1 {
			t.Fatalf("pruneZeroSectionCOFFMembers = (%d, %v), want (1, nil)", pruned, err)
		}
		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		members := parseTestArchive(t, got)
		if len(members) != 2 ||
			members[0].name != "valid.o/" || !bytes.Equal(members[0].payload, validObj) ||
			members[1].name != "valid_big.o/" || !bytes.Equal(members[1].payload, validBigObj) {
			t.Fatalf("unexpected remaining members: %+v", members)
		}
	})

	t.Run("bsd-extended-names", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "bsd.a")
		symdefPayload := append([]byte("__.SYMDEF SORTED\x00\x00\x00\x00"), 0, 0, 0, 0)
		validName := []byte("valid_member.o\x00\x00")
		zeroName := []byte("empty_member.o\x00\x00")
		validPayload := append(bytes.Clone(validName), makeCOFFObject(pe.IMAGE_FILE_MACHINE_I386, 1, 32)...)
		zeroPayload := append(bytes.Clone(zeroName), makeCOFFObject(pe.IMAGE_FILE_MACHINE_I386, 0, coffHeaderSize)...)
		archive := buildTestArchive(
			testARMember{name: "#1/20", payload: symdefPayload},
			testARMember{name: "#1/16", payload: validPayload},
			testARMember{name: "#1/16", payload: zeroPayload},
		)
		if err := os.WriteFile(path, archive, 0o644); err != nil {
			t.Fatal(err)
		}
		pruned, err := pruneZeroSectionCOFFMembers(path)
		if err != nil || pruned != 1 {
			t.Fatalf("pruneZeroSectionCOFFMembers = (%d, %v), want (1, nil)", pruned, err)
		}
		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		members := parseTestArchive(t, got)
		if len(members) != 1 || members[0].name != "#1/16" || !bytes.Equal(members[0].payload, validPayload) {
			t.Fatalf("unexpected remaining BSD members: %+v", members)
		}
	})

	t.Run("malformed-ar-headers", func(t *testing.T) {
		badTrailer := buildTestArchive(testARMember{name: "a.o/", payload: makeCOFFObject(pe.IMAGE_FILE_MACHINE_I386, 1, 20)})
		badTrailer[len(arMagic)+58] = 'X'

		badSize := buildTestArchive(testARMember{name: "a.o/", payload: makeCOFFObject(pe.IMAGE_FILE_MACHINE_I386, 1, 20)})
		copy(badSize[len(arMagic)+48:len(arMagic)+58], "notanumber")

		truncatedPayload := buildTestArchive(testARMember{name: "a.o/", payload: makeCOFFObject(pe.IMAGE_FILE_MACHINE_I386, 1, 20)})
		truncatedPayload = truncatedPayload[:len(truncatedPayload)-5]

		truncatedOddPad := buildTestArchive(testARMember{name: "a.o/", payload: makeCOFFObject(pe.IMAGE_FILE_MACHINE_I386, 1, 21)})
		truncatedOddPad = truncatedOddPad[:len(truncatedOddPad)-1]

		for _, tc := range []struct {
			name string
			data []byte
		}{
			{name: "bad-magic", data: []byte("!<thin>\n")},
			{name: "truncated-header", data: []byte(arMagic + "short")},
			{name: "bad-trailer", data: badTrailer},
			{name: "bad-size", data: badSize},
			{name: "truncated-payload", data: truncatedPayload},
			{name: "truncated-odd-pad", data: truncatedOddPad},
		} {
			t.Run(tc.name, func(t *testing.T) {
				dir := t.TempDir()
				path := filepath.Join(dir, "bad.a")
				if err := os.WriteFile(path, tc.data, 0o644); err != nil {
					t.Fatal(err)
				}
				if _, err := pruneZeroSectionCOFFMembers(path); err == nil {
					t.Fatalf("expected error for %s archive", tc.name)
				}
				temps, err := filepath.Glob(filepath.Join(dir, ".filtered-*.a"))
				if err != nil || len(temps) != 0 {
					t.Fatalf("temporary filtered archive remains: %v, %v", temps, err)
				}
			})
		}
	})
}

func TestIsZeroSectionCOFF(t *testing.T) {
	elfObj := append([]byte{0x7f, 'E', 'L', 'F'}, make([]byte, 60)...)
	mzImage := append([]byte{'M', 'Z', 0x90, 0x00}, make([]byte, 60)...)
	wrongClassID := coffBigObjClassID
	wrongClassID[0] ^= 0xff

	cases := []struct {
		name        string
		prefix      []byte
		payloadSize int64
		want        bool
	}{
		{name: "std-i386-zero", prefix: makeCOFFObject(pe.IMAGE_FILE_MACHINE_I386, 0, 20), payloadSize: 20, want: true},
		{name: "std-amd64-zero", prefix: makeCOFFObject(pe.IMAGE_FILE_MACHINE_AMD64, 0, 56), payloadSize: 120, want: true},
		{name: "std-arm64-zero", prefix: makeCOFFObject(pe.IMAGE_FILE_MACHINE_ARM64, 0, 20), payloadSize: 20, want: true},
		{name: "std-riscv64-zero", prefix: makeCOFFObject(pe.IMAGE_FILE_MACHINE_RISCV64, 0, 20), payloadSize: 20, want: true},
		{name: "std-loongarch64-zero", prefix: makeCOFFObject(pe.IMAGE_FILE_MACHINE_LOONGARCH64, 0, 20), payloadSize: 20, want: true},
		{name: "std-amd64-nonzero", prefix: makeCOFFObject(pe.IMAGE_FILE_MACHINE_AMD64, 1, 40), payloadSize: 40, want: false},
		{name: "std-truncated", prefix: makeCOFFObject(pe.IMAGE_FILE_MACHINE_I386, 0, 20)[:19], payloadSize: 19, want: false},
		{name: "bigobj-zero", prefix: makeCOFFBigObj(pe.IMAGE_FILE_MACHINE_AMD64, 2, 0, coffBigObjClassID), payloadSize: 56, want: true},
		{name: "bigobj-i386-zero", prefix: makeCOFFBigObj(pe.IMAGE_FILE_MACHINE_I386, 2, 0, coffBigObjClassID), payloadSize: 56, want: true},
		{name: "bigobj-nonzero", prefix: makeCOFFBigObj(pe.IMAGE_FILE_MACHINE_AMD64, 2, 4, coffBigObjClassID), payloadSize: 100, want: false},
		{name: "bigobj-v1", prefix: makeCOFFBigObj(pe.IMAGE_FILE_MACHINE_AMD64, 1, 0, coffBigObjClassID), payloadSize: 56, want: false},
		{name: "bigobj-wrong-clsid", prefix: makeCOFFBigObj(pe.IMAGE_FILE_MACHINE_AMD64, 2, 0, wrongClassID), payloadSize: 56, want: false},
		{name: "bigobj-truncated", prefix: makeCOFFBigObj(pe.IMAGE_FILE_MACHINE_AMD64, 2, 0, coffBigObjClassID)[:55], payloadSize: 55, want: false},
		{name: "pe-short-import", prefix: makePEShortImportObject(pe.IMAGE_FILE_MACHINE_I386), payloadSize: 47, want: false},
		{name: "elf", prefix: elfObj, payloadSize: int64(len(elfObj)), want: false},
		{name: "mz", prefix: mzImage, payloadSize: int64(len(mzImage)), want: false},
		{name: "unknown-machine", prefix: makeCOFFObject(pe.IMAGE_FILE_MACHINE_UNKNOWN, 0, 20), payloadSize: 20, want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isZeroSectionCOFF(tc.prefix, tc.payloadSize); got != tc.want {
				t.Fatalf("isZeroSectionCOFF = %v, want %v", got, tc.want)
			}
		})
	}
}
