# FFmpeg Builder Refactoring Plan

## Overview

Refactoring the internal builder from 2,129 lines of duplicated code to ~850 lines of clean, maintainable code with incremental builds, per-library logging, and selective build support.

## Approach

**Strategy:** Rename existing `internal/builder` → `internal/old-builder`, create fresh implementation in `internal/builder`.

**Benefits:**
- No risk to working code
- Easy side-by-side comparison
- Clean slate without legacy constraints
- Simple rollback if needed

## Architecture

### New Structure

```
internal/builder/
├── main.go           # CLI entry point (~100 lines)
├── library.go        # Library struct and core build logic (~200 lines)
├── buildsystems.go   # Autoconf, CMake, Meson, Cargo implementations (~150 lines)
├── libraries.go      # All 22 library definitions (~300 lines)
├── state.go          # Incremental build state management (~50 lines)
└── download.go       # grab-based downloader (~50 lines)
```

### Key Design Decisions

1. **Library as Data:** Each library is a struct with configuration, not a function
2. **Platform Support:** Libraries specify platforms (empty = all platforms)
3. **Incremental Builds:** Hash-based state tracking (URL + config hash)
4. **Simple Dependencies:** Linear dependency order, platform-aware skipping
5. **One External Dep:** `github.com/cavaliergopher/grab/v3` for downloads

## Implementation Phases

### Phase 1: Core Framework

**Files:** `library.go`, `buildsystems.go`, `state.go`, `download.go`

**Key types:**
```go
type Library struct {
    Name         string
    URL          string
    ArchiveType  string      // "tar.gz", "tar.bz2", "tar.xz", "zip"
    StripPrefix  string
    Platform     []string    // empty = all platforms, ["linux"], ["darwin"], etc.
    BuildSystem  BuildSystem
    ConfigureArgs func(os string) []string
    PostExtract  func(srcPath string) error // optional patches
    Dependencies []*Library
}

type BuildSystem interface {
    Configure(lib *Library, srcPath string) error
    Build(lib *Library, srcPath string) error
}
```

**Incremental build logic:**
- Store `.build-state` JSON in each library's build directory
- Track: URL, config hash, build timestamp
- Skip if: state exists + URL matches + config hash matches + output exists
- Rebuild if: any of above fail

**Download:**
- Use `grab/v3` for robust downloads with resume support
- Simple progress logging (no animations)
- Per-library log file via `io.MultiWriter`

### Phase 2: Library Definitions

**File:** `libraries.go`

Convert all 22 library build functions to data structures:

```go
var x264 = &Library{
    Name: "x264",
    URL: "https://code.videolan.org/videolan/x264/-/archive/master/x264-master.tar.bz2",
    ArchiveType: "tar.bz2",
    StripPrefix: "x264-master/",
    BuildSystem: &AutoconfBuild{},
    ConfigureArgs: func(os string) []string {
        return []string{
            "--disable-cli",
            "--enable-static",
            "--enable-strip",
        }
    },
}

var iconv = &Library{
    Name: "iconv",
    URL: "https://ftp.gnu.org/pub/gnu/libiconv/libiconv-1.18.tar.gz",
    Platform: []string{"darwin"}, // macOS only
    // ... rest of config
}
```

**Platform handling:**
- Libraries with `Platform: []string{"linux"}` auto-skip on macOS
- Libraries with `Platform: []string{"darwin"}` auto-skip on Linux
- Libraries with `Platform: nil` or empty build on all platforms

**Dependency order:**
- Define in `libraries.go` as ordered slice
- Builder respects order but skips platform-incompatible libs

### Phase 3: Main + Logging

**File:** `main.go`

**CLI interface:**
```bash
# Build everything
go run ./internal/builder lib/linux_amd64/libffmpeg.a

# Build specific libraries only (with dependencies)
go run ./internal/builder lib/linux_amd64/libffmpeg.a x264 vpx

# Force rebuild (ignore state)
go run ./internal/builder --force lib/linux_amd64/libffmpeg.a
```

**Per-library logging:**
- Create `.build/build/{lib-name}/build.log`
- `io.MultiWriter` to both stdout and log file
- Clean separation for debugging

**Progress reporting:**
```
[1/22] x264: Checking build state... SKIP (already built)
[2/22] x265: Downloading...
[2/22] x265: Extracting...
[2/22] x265: Configuring...
[2/22] x265: Building...
[2/22] x265: Complete (45.2s)
```

### Phase 4: Testing & Validation

**Validation approach:**

1. **Symbol comparison** (primary):
   - Extract well-known symbols: `avcodec_version`, `avformat_version`, etc.
   - Compare symbol counts and key exports
   - Warning if different, not failure

2. **Binary hash** (secondary):
   - SHA256 hash comparison
   - **WARNING ONLY** (due to HEAD snapshots for x264/x265)
   - Documents that hashes may differ for HEAD builds

3. **Build success** (required):
   - All libraries build without errors
   - Final static library created
   - All expected `.a` files in staging

**Test script:**
```bash
#!/usr/bin/env bash
# Compare old vs new builder outputs

# Build with old builder
mv internal/builder internal/builder-new
mv internal/old-builder internal/builder
go run ./internal/builder lib/test-old.a
mv internal/builder internal/old-builder
mv internal/builder-new internal/builder

# Build with new builder
rm -rf .build
go run ./internal/builder lib/test-new.a

# Compare symbols (primary validation)
nm lib/test-old.a | grep -E "avcodec_version|avformat_version" > old-symbols.txt
nm lib/test-new.a | grep -E "avcodec_version|avformat_version" > new-symbols.txt
diff old-symbols.txt new-symbols.txt

# Compare hashes (warning only)
sha256sum lib/test-old.a lib/test-new.a
echo "Note: Hash differences expected for x264/x265 HEAD snapshots"
```

## Migration Plan

### Step 1: Backup (5 minutes)
```bash
git checkout -b builder-refactor
mv internal/builder internal/old-builder
git add internal/old-builder
git commit -m "refactor: backup old builder before refactor"
```

### Step 2: Implement New Builder (3 days)
- Create `internal/builder/` with new implementation
- Test incrementally as each phase completes
- Keep old builder functional for comparison

### Step 3: Update Justfile (5 minutes)
```just
# Old builder (temporary)
build-ffmpeg-old:
    go run ./internal/old-builder "lib/${GOOS}_${GOARCH}/libffmpeg.a"

# New builder
build-ffmpeg:
    go run ./internal/builder "lib/${GOOS}_${GOARCH}/libffmpeg.a"
```

### Step 4: Validation (1 day)
- Build with both builders
- Compare outputs (symbols + hashes)
- Fix any discrepancies
- Document expected differences (HEAD snapshots)

### Step 5: Cleanup (5 minutes)
```bash
rm -rf internal/old-builder
git add .
git commit -m "refactor: remove old builder after validation"
```

## Success Criteria

- ✅ All 22 libraries build successfully
- ✅ Final static library contains expected symbols
- ✅ Incremental builds work (skip already-built libs)
- ✅ Selective builds work (build only specified libs)
- ✅ Per-library log files created
- ✅ Platform-specific libs auto-skip on wrong platform
- ✅ Code reduced by ~63% (2,129 → ~850 lines)
- ✅ Maintainable: adding new library = add struct, not function

## Known Differences from Old Builder

1. **Binary hashes may differ** for x264/x265 (using HEAD snapshots) - WARNING ONLY
2. **Build order optimized** based on dependencies
3. **Better error messages** with per-library context
4. **Progress reporting** shows current library being built

## Future Enhancements (Not in Scope)

- Parallel library builds (would require more complexity)
- Build system plugins (YAGNI - 4 build systems is enough)
- Remote caching (not needed for this use case)
- Fancy progress bars (conflicts with CI simplicity goal)

---

**Start Date:** 13 November 2025
**Estimated Completion:** 16 November 2025
**Actual Completion:** TBD
