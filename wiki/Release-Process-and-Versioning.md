# Release process and versioning

`ffmpeg-statigo` uses two release schemes so heavy native archives and lightweight Go module tags evolve cleanly on separate tracks.

## Dual version schemes

| Scheme | Tag format | Trigger recipe | GitHub Latest | Published assets |
| --- | --- | --- | --- | --- |
| Library and CLI release | `lib-X.Y.Z.N` | `just ffmpeg-release lib-X.Y.Z.N` | `false` (`make_latest: false`) | 24 archives (12 static library tarballs + 12 standalone CLI archives) + 1 unified `SHA256SUMS` |
| Go module and package release | `vX.Y.Z.N` | `just go-release X.Y.Z.N` | `true` | 18 Haberdasher CLI packages and binaries + 1 `SHA256SUMS` (18 lines), plus detached `latest` Nix tag |

In both schemes, `X.Y.Z` is the upstream FFmpeg version (for example `8.1.2`) and `N` is the build revision (for example `1`). Use `just ffmpeg-release` and `just go-release` rather than the generic three-part `just release` recipe from Tailor.

`lib/fetch.go` pins the active pair:

- `const Version = "X.Y.Z"` holds the FFmpeg version.
- `const LibRelease = "lib-X.Y.Z.N"` pins the exact GitHub Release tag that `cmd/download-lib` fetches.

### `lib-X.Y.Z.N` release assets (24 archives + `SHA256SUMS`)

Each library release ships across six targets (`linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`, `windows/amd64`, `windows/386`) and two profiles (`default` GPLv3 and `embedded` LGPLv2.1):

| Target | Static library (default / embedded) | Standalone CLI (default / embedded) |
| --- | --- | --- |
| `linux/amd64` | `ffmpeg-linux-amd64.tar.gz`, `ffmpeg-embedded-linux-amd64.tar.gz` | `ffmpeg-bin-linux-amd64.tar.gz`, `ffmpeg-bin-embedded-linux-amd64.tar.gz` |
| `linux/arm64` | `ffmpeg-linux-arm64.tar.gz`, `ffmpeg-embedded-linux-arm64.tar.gz` | `ffmpeg-bin-linux-arm64.tar.gz`, `ffmpeg-bin-embedded-linux-arm64.tar.gz` |
| `darwin/amd64` | `ffmpeg-darwin-amd64.tar.gz`, `ffmpeg-embedded-darwin-amd64.tar.gz` | `ffmpeg-bin-darwin-amd64.tar.gz`, `ffmpeg-bin-embedded-darwin-amd64.tar.gz` |
| `darwin/arm64` | `ffmpeg-darwin-arm64.tar.gz`, `ffmpeg-embedded-darwin-arm64.tar.gz` | `ffmpeg-bin-darwin-arm64.tar.gz`, `ffmpeg-bin-embedded-darwin-arm64.tar.gz` |
| `windows/amd64` | `ffmpeg-windows-amd64.tar.gz`, `ffmpeg-embedded-windows-amd64.tar.gz` | `ffmpeg-bin-windows-amd64.zip`, `ffmpeg-bin-embedded-windows-amd64.zip` |
| `windows/386` | `ffmpeg-windows-386.tar.gz`, `ffmpeg-embedded-windows-386.tar.gz` | `ffmpeg-bin-windows-386.zip`, `ffmpeg-bin-embedded-windows-386.zip` |

Every CLI archive unpacks into `<stem>/` containing `ffmpeg` (`ffmpeg.exe` on Windows), `ffprobe` (`ffprobe.exe` on Windows), `LICENSE` (`COPYING.GPLv3` or `COPYING.LGPLv2.1`), and `NOTICE.txt`. One unified `SHA256SUMS` file covers all 24 archives.

### `vX.Y.Z.N` release assets (18 Haberdasher assets + `SHA256SUMS`)

When the pinned `LibRelease` contains `SHA256SUMS` and all six default (GPLv3) `ffmpeg-bin-*` archives, `go-release.yml` verifies those six archives against `SHA256SUMS` and runs `packaging/package-release.sh` with pinned `nfpm` v2.49.0 and `packaging/nfpm.yaml`.

The script produces 18 assets and an 18-line `SHA256SUMS`:

- **Linux packages (6)**: `ffmpeg-statigo_<ver>_{amd64,arm64}.deb`, `ffmpeg-statigo-<ver>-1.{x86_64,aarch64}.rpm`, and `ffmpeg-statigo-<ver>-1-{x86_64,aarch64}.pkg.tar.zst`. Each package installs `/usr/bin/ffmpeg-statigo`, `/usr/bin/ffprobe-statigo`, `LICENSE`, and `NOTICE.txt`, and declares a `glibc >= 2.39` dependency (`libc6 (>= 2.39)`, `libc.so.6(GLIBC_2.39)(64bit)`, or `glibc>=2.39`).
- **Linux standalone binaries (4)**: `ffmpeg-statigo-linux-{amd64,arm64}` and `ffprobe-statigo-linux-{amd64,arm64}`.
- **macOS standalone binaries (4)**: `ffmpeg-statigo-darwin-{amd64,arm64}` and `ffprobe-statigo-darwin-{amd64,arm64}`.
- **Windows zips (4)**: `ffmpeg-statigo-windows-{amd64,386}.zip` and `ffprobe-statigo-windows-{amd64,386}.zip`, each bundling the renamed `.exe`, `LICENSE`, and `NOTICE.txt`.

Haberdasher publishes these assets to `apt.wimpys.world`, `rpm.wimpys.world`, `pkg.wimpys.world`, `brew.wimpys.world`, `scoop.wimpys.world`, and `nix.wimpys.world`.

### Detached `latest` Nix tag

For every non-prerelease `vX.Y.Z.N` release with CLI assets, `go-release.yml` updates the detached `latest` git tag:

1. Fetches `refs/tags/latest` if it exists and compares its `nix/sources.json` version with the new release using `sort -V`. If `latest` is already newer, the step skips the update so maintenance releases on older branches do not downgrade `latest`.
2. Checks out a detached HEAD at `vX.Y.Z.N`, writes the generated `nix/sources.json` (holding `version`, `libRelease`, and `url` plus SRI `sha256-...` hashes for `x86_64-linux`, `aarch64-linux`, `x86_64-darwin`, and `aarch64-darwin` pointing at `ffmpeg-bin-*` on `lib-X.Y.Z.N`), and commits it on the detached HEAD.
3. Force-pushes `refs/tags/latest`.

`flake.nix` and `nix/package.nix` expose `packages.default` and `packages.ffmpeg-statigo` whenever `nix/sources.json` is present on that tag.

## GitHub Actions workflows

| Workflow | Trigger | Purpose |
| --- | --- | --- |
| `.github/workflows/_build-libs.yml` | `workflow_call` | Reusable native builder for all 12 platform and profile combinations across Linux, macOS, and Windows |
| `.github/workflows/ffmpeg-release.yml` | `lib-X.Y.Z.N` tag push or `workflow_dispatch` | Validates 3-way version parity, calls `_build-libs.yml`, checks all 24 archives, writes unified `SHA256SUMS`, and publishes the `lib-` release (`make_latest: false`) |
| `.github/workflows/go-test.yml` | Push or pull request on `main`, or `workflow_dispatch` | Tests Go packages against downloaded archives or freshly built artefacts, runs race tests, lints, and checks byte-identical binding regeneration |
| `.github/workflows/go-release.yml` | `vX.Y.Z.N` tag push or `workflow_dispatch` | Tests Unix and Windows across both profiles, packages the 18 Haberdasher CLI assets, creates the `v` tag and GitHub Release, and updates the detached `latest` Nix tag |
| `.github/workflows/tailor-wiki.yml` | Push to `wiki/**` or `.github/workflows/tailor-wiki.yml` on `main`, or `workflow_dispatch` | Publishes `wiki/` to `<repo>.wiki.git` after verifying `.tailor-wiki-base` adoption and checking for external wiki edits |

### `_build-libs.yml` details

`_build-libs.yml` accepts four optional inputs: `clean_ffmpeg` (default `false`), `test` (default `false`), `build_programs` (default `true`), and `release_tag` (default empty).

The matrix runs natively on `ubuntu-24.04` (`linux/amd64`), `ubuntu-24.04-arm` (`linux/arm64`), `macos-15-intel` (`darwin/amd64`), `macos-15` (`darwin/arm64`), and `windows-2025` (MSYS2 `UCRT64` for `windows/amd64`, MSYS2 `MINGW32` with Rust `1.93.1` targeting `i686-pc-windows-gnu` for `windows/386`).

Before uploading artefacts, `_build-libs.yml` enforces binary portability gates:

- **Linux**: Checks `LICENSE` text, verifies `DT_NEEDED` allows only core glibc libraries (`libc`, `libm`, `libmvec`, `libpthread`, `libdl`, `librt`, dynamic loader), rejects any ELF symbol version above `GLIBC_2.39`, and runs an `ffmpeg` plus `ffprobe` H.264 encode/probe smoke test (`libx264` or `libopenh264`).
- **macOS**: Checks `LICENSE` text, runs `codesign --verify --strict`, verifies `otool -L` lists only `/usr/lib/*` and `/System/Library/Frameworks/*`, checks `vtool -show-build` reports `minos <= 13.0`, and runs the H.264 smoke test.
- **Windows**: Verifies archive COFF object architectures, checks `objdump -p` DLL imports against the built-in OS allowlist (`ucrtbase.dll` / `api-ms-win-crt-*.dll` on `amd64` only, `msvcrt.dll` on `386` only), and runs the H.264 smoke test in PowerShell with `PATH` restricted to `C:\Windows\System32;C:\Windows`.

### `go-test.yml` build modes

`go-test.yml` uses fast download mode by default (fetching pre-built archives from `LibRelease`). It switches to from-source mode (`_build-libs.yml`) automatically when a commit or pull request modifies any of these paths or declarations:

- `internal/builder/*`
- `include/*`
- `ffmpeg_link_*.go`
- `.github/workflows/_build-libs.yml`
- `const LibRelease` in `lib/fetch.go`

Dispatch `go-test.yml` with `build_libs: true` to force a from-source build manually. On `linux/amd64`, `go-test.yml` also runs `go test -race -count=1 ./...` and `golangci-lint`. On all four Unix default-profile jobs, it runs `go run ./internal/generator` and fails if `git diff --quiet -- ./*.gen.go` shows any difference.

### `ffmpeg-release.yml` parity check

Before building, `ffmpeg-release.yml` verifies that three values agree:

1. The `lib-X.Y.Z.N` release tag matches `const LibRelease` in `lib/fetch.go`.
2. The `X.Y.Z` portion matches `const Version` in `lib/fetch.go`.
3. The `X.Y.Z` portion matches the FFmpeg tarball URL (`nX.Y.Z.tar.*`) in `internal/builder/libraries.go`.

### `tailor-wiki.yml` publication guard

Tailor manages `.github/workflows/tailor-wiki.yml`. On pushes to `main` touching `wiki/**`, the job checks `wiki/.tailor-wiki-base` against the remote wiki HEAD and refuses to push if the remote wiki changed outside the workflow. Edit wiki pages in `wiki/` rather than directly in the GitHub web editor.

## FFmpeg upgrade checklist

Follow these steps in order when bumping FFmpeg or a static dependency:

1. Update the library URL or version in `internal/builder/libraries.go` (or `internal/builder/openh264.go`).
2. Refresh SHA-256 pins in `internal/builder/digests.go` with `go run ./internal/builder --update-digests` (and `--embedded --update-digests` when OpenH264 changes), then verify each digest against upstream checksums.
3. Rebuild locally with `just build-static ffmpeg --clean && just build-static` (and `just build-embedded` if embedded changed) to verify compilation.
4. Copy updated FFmpeg headers into `include/` when upgrading the FFmpeg version.
5. Regenerate bindings and IR goldens with `go run ./internal/generator -dump-ir`. Inspect `git diff -- internal/generator/testdata/ir/` to confirm every change is expected. Adjust `skipCeiling` in `internal/generator/main.go` only if new FFmpeg headers add legitimate unemittable symbols.
6. Update `const Version` and `const LibRelease` in `lib/fetch.go` and the library version table in `README.md`.
7. Run `just build`, `just test`, and `just lint`.
8. Open a pull request. Because `internal/builder/*`, `include/*`, or `LibRelease` changed, `go-test.yml` builds all 12 platform and profile combinations from source and runs the full test and generator-determinism matrix.
9. When CI passes, merge the pull request to `main`.
10. Trigger the library and CLI release: `just ffmpeg-release lib-X.Y.Z.N`. Wait until all 24 archives and `SHA256SUMS` publish.
11. Trigger the Go module and Haberdasher release: `just go-release X.Y.Z.N`. Confirm the 18 CLI assets publish and the detached `latest` Nix tag updates.

Next: read [Troubleshooting and FAQ](Troubleshooting-and-FAQ) if a local build, runtime link, or platform gate fails.
