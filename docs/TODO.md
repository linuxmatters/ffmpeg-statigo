# TODO

- [x] Refactor the internal build tool
- [x] Reorganise project structure
- [x] FFmpeg build argument resolver.to enable/disable codecs
- [x] FFmpeg feature enablement
  --enable-libvmaf         enable vmaf filter via libvmaf [no]
  --enable-whisper         enable whisper filter [no]
- [x] Rename to ffmpeg-statigo
- [x] How to embed/distribute the static FFmpeg library?
- [x] Review default codecs:
  - https://ffmpeg.martin-riedl.de/
  - https://github.com/markus-perl/ffmpeg-build-script/blob/master/build-ffmpeg
- [x] Create some test cases that exercise some of the FFmpeg API surface

## From the original author

- [x] Expose more headers.
- [ ] Expose platform specific headers.
- [ ] Cleanup internal packages.

---

# Static Library Distribution Implementation Plan

**Decision:** GitHub Releases + Auto-Download approach

**Rationale:**
- Git LFS is NOT idiomatic for Go modules - `go get` doesn't fetch LFS files
- Research found zero examples of popular Go CGO libraries using LFS
- mattn/go-sqlite3 (canonical example) bundles C source directly, not via LFS
- GitHub Releases provide unlimited bandwidth vs LFS 10 GiB/month limit
- Platform-specific downloads (100MB) vs downloading all platforms (400MB)

## Phase 1: GitHub Actions - Release Workflow

**Goal:** Create/updated automated workflow to build and publish static libraries as GitHub Release assets. See the existing GitHub: `.github/workflows` and update as appropriated.

**Testing:**
```bash
# Test via workflow_dispatch before tagging
gh workflow run release.yml -f tag=v0.1.0-test

# Once validated, create actual release
git tag v0.1.0
git push origin v0.1.0
```

## Phase 2: Library Download Infrastructure

**Goal:** Implement automatic library download on first import, perhaps via a new file `lib/fetch.go`. See `internal/builder/download.go` for the Go modules already used in the project to download tarballs. The download must be able to resolve the download that matches it's internal tagged release and get the correct architecture version of the library.

**Testing:**
```bash
# Remove existing libraries to test download
rm -rf lib/linux_amd64 lib/darwin_arm64

# Import will trigger download
go build ./examples/metadata
```

### Update .gitignore

**File:** `.gitignore`

Add the following:

```
lib/*/
!lib/fetch.go
!lib/README.md
```

### Create Library Directory README

**File:** `lib/README.md`

```markdown
# FFmpeg Static Libraries

This directory contains platform-specific FFmpeg static libraries.

## Automatic Download

Libraries are automatically downloaded on first import of the `ffmpeg-statigo` package.
No manual intervention required.

## Manual Download

If automatic download fails:

    # Linux amd64
    wget https://github.com/linuxmatters/ffmpeg-statigo/releases/download/v0.1.0/ffmpeg-linux-amd64.tar.gz
    tar -xzf ffmpeg-linux-amd64.tar.gz

    # macOS arm64
    wget https://github.com/linuxmatters/ffmpeg-statigo/releases/download/v0.1.0/ffmpeg-darwin-arm64.tar.gz
    tar -xzf ffmpeg-darwin-arm64.tar.gz

## Verification

    # Verify libraries exist and are valid
    ls -lh lib/*/libffmpeg.a

    # Each library should be ~75-100MB

## Structure

    lib/
    ├── linux_amd64/
    │   └── libffmpeg.a
    ├── linux_arm64/
    │   └── libffmpeg.a
    ├── darwin_amd64/
    │   └── libffmpeg.a
    └── darwin_arm64/
        └── libffmpeg.a
```

## Phase 3: Update Existing Workflows

**Goal:** Modify ffmpeg.yml to NOT commit libraries to git

### Update FFmpeg Build Workflow

**File:** `.github/workflows/ffmpeg.yml`

**Changes:**
1. Remove pull request creation step
2. Libraries only built for validation, not committed
3. Actual release happens via release.yml workflow

## Phase 4: Update Documentation

### Update README.md

**Add Installation Section:**

```markdown
## Installation

    go get github.com/linuxmatters/ffmpeg-statigo

**No additional steps required.** FFmpeg static libraries are automatically downloaded on first import.

### Manual Installation

If you need to download libraries manually:

    # Download for your platform
    wget https://github.com/linuxmatters/ffmpeg-statigo/releases/download/v0.1.0/ffmpeg-linux-amd64.tar.gz

    # Extract to lib directory
    tar -xzf ffmpeg-linux-amd64.tar.gz

    # Verify checksum
    wget https://github.com/linuxmatters/ffmpeg-statigo/releases/download/v0.1.0/ffmpeg-linux-amd64.tar.gz.sha256
    sha256sum -c ffmpeg-linux-amd64.tar.gz.sha256

### Offline / Air-Gapped Environments

For environments without internet access:

1. Download libraries on a connected machine
2. Transfer tarball to air-gapped environment
3. Extract to `lib/` directory before building
4. Set build tag: `go build -tags libffmpeg_bundled`

## Phase 5: Testing & Validation

### Pre-Release Testing Checklist

- [ ] Test automatic download on clean checkout

```bash
cd /tmp
git clone https://github.com/linuxmatters/ffmpeg-statigo.git
cd ffmpeg-statigo
go build ./examples/metadata
# Should auto-download libraries
```

- [ ] Test manual download instructions

```bash
rm -rf lib/*/
# Follow README manual download steps
go build ./examples/metadata
```

- [ ] Test checksum validation

```bash
# Corrupt tarball and verify download fails
```

- [ ] Test all platforms via GitHub Actions

```bash
gh workflow run release.yml -f tag=v0.1.0-rc1
```

- [ ] Test Go module import
```bash
cd /tmp/test-project
go mod init test
go get github.com/linuxmatters/ffmpeg-statigo@v0.1.0-rc1
```

## Success Criteria

- ✅ `go get github.com/linuxmatters/ffmpeg-statigo` works without manual steps
- ✅ Libraries automatically download on first import
- ✅ No static libraries committed to git repository
- ✅ Release workflow creates GitHub Release with all platform libraries
- ✅ Checksums verify download integrity
- ✅ Clear error messages if download fails
- ✅ Manual fallback documented and tested
- ✅ CI/CD workflows updated and tested
- ✅ Documentation complete and accurate

---

# Onboard

This is a hard fork of the [ffmpeg-go](https://github.com/csnewman/ffmpeg-go) project, now called **ffmpeg-statigo**.

ffmpeg-statigo has been updated to Go 1.24, all dependencies uplifted to current versions, GitHub CI is fixed so the ffmpeg static libraries are built for all supported architectures and the generator has been updated to support current `clang` on Linux. Mostly critically this project has been updated from FFmpeg 6.1 to FFmpeg 8.0 including the addition os `x265`, `dav1d`, `rav1e` and hardware acceleration support for NVENC, QuickSync and Vulkan to complement the exist VideoToolbox support.

The git history has been purged all tags and historical commits of the static ffmpeg libraries as we are preparing this project to be launched under it's new name in a different GitHUb organisation. This brought the git repo down in size from 1.9GB to 9MB.

NixOS is the host development workstation. There is a `flake.nix` in the project to enable the required software in a Nix devevelopment shell, which is automatically activated via `direnv`. The `justfile` is used for all build and test commands. `fish` is the default shell.

Build the project using `just build`. Using `just build` is the only valid way to build. Never build the project using adhoc commands. Never rediect builds elsewhere, the build system does the right thing for you. Work with the provide build tool, not against it. The following directory structure is populated with the build assets and log files for debugging.

```
.build
├── build
│   ├── ass
│   │   └── build.log
│   ├── dav1d
│   │   └── build.log
│   ├── ffmpeg
│   │   └── build.log
│   ├── fontconfig
│   │   └── build.log
│   ├── freetype
│   │   └── build.log
│   ├── fribidi
│   │   └── build.log
│   ├── glslang
│   │   └── build.log
│   ├── harfbuzz
│   │   └── build.log
│   ├── iconv
│   │   └── build.log
│   ├── lame
│   │   └── build.log
│   ├── nvcodec
│   │   └── build.log
│   ├── ogg
│   │   └── build.log
│   ├── opus
│   │   └── build.log
│   ├── png
│   │   └── build.log
│   ├── rav1e
│   │   └── build.log
│   ├── theora
│   │   └── build.log
│   ├── unibreak
│   │   └── build.log
│   ├── vorbis
│   │   └── build.log
│   ├── vpl
│   │   └── build.log
│   ├── vpx
│   │   └── build.log
│   ├── vulkan
│   │   └── build.log
│   ├── x264
│   │   └── build.log
│   ├── x265
│   │   └── build.log
│   ├── xml2
│   │   └── build.log
│   └── zlib
│       └── build.log
├── downloads
│   └── <all tarball are here>
├── src
│   ├── ass
│   ├── dav1d
│   ├── ffmpeg
│   ├── fontconfig
│   ├── freetype
│   ├── fribidi
│   ├── glslang
│   ├── harfbuzz
│   ├── lame
│   ├── nvcodec
│   ├── ogg
│   ├── opus
│   ├── png
│   ├── rav1e
│   ├── theora
│   ├── unibreak
│   ├── vorbis
│   ├── vpl
│   ├── vpx
│   ├── vulkan
│   ├── x264
│   ├── x265
│   ├── xml2
│   └── zlib
└── staging
    ├── bin
    ├── etc
    ├── include
    ├── lib
    ├── lib64
    ├── share
    └── var
```

Read the `README.md`, `TODO.md` and analyse the code to get a full understanding of the project. Let me know when you are ready to collaborate.
