# Standalone CLI downloads

Use the managed Haberdasher package repositories for automatic updates, or download direct `lib-X.Y.Z.N` release archives for pinned and offline installs.

## Managed package repositories

Managed packages install the `default` (GPLv3) profile as `ffmpeg-statigo` and `ffprobe-statigo` so they coexist with a system `ffmpeg`. Homebrew and Scoop also provide `ffprobe-statigo`.

| Platform | Repository |
| --- | --- |
| Debian and Ubuntu | [apt.wimpys.world](https://apt.wimpys.world/#ffmpeg-statigo) |
| Fedora and openSUSE | [rpm.wimpys.world](https://rpm.wimpys.world/#ffmpeg-statigo) |
| Arch Linux | [pkg.wimpys.world](https://pkg.wimpys.world/#ffmpeg-statigo) |
| macOS (Homebrew) | [brew.wimpys.world](https://brew.wimpys.world/#ffmpeg-statigo) |
| Windows (Scoop) | [scoop.wimpys.world](https://scoop.wimpys.world/#ffmpeg-statigo) |
| Nix and NixOS | [nix.wimpys.world](https://nix.wimpys.world/#ffmpeg-statigo) |

APT, RPM, Arch, and Nix repositories are signed. Homebrew and Scoop verify packages over HTTPS with SHA-256 checksums. Each repository page gives the setup steps and keeps `ffmpeg-statigo` up to date.

On any platform, you can also install and update the release binary with [bin](https://github.com/marcosnils/bin):

```sh
bin install github.com/linuxmatters/ffmpeg-statigo
```

## Direct archives

Pre-built static `ffmpeg` and `ffprobe` executables ship on each `lib-X.Y.Z.N` release on [GitHub Releases](https://github.com/linuxmatters/ffmpeg-statigo/releases) alongside the static libraries. Each `lib-X.Y.Z.N` release sets `make_latest: false`, because `vX.Y.Z.N` releases track the Go module as Latest.

Direct downloads use this URL pattern:

```text
https://github.com/linuxmatters/ffmpeg-statigo/releases/download/<lib-tag>/<archive>
```

Each `lib-X.Y.Z.N` release publishes 12 CLI archives across six platforms and two profiles (`.tar.gz` on Linux and macOS, `.zip` on Windows):

| Platform | Default (GPLv3) | Embedded (LGPLv2.1) |
| --- | --- | --- |
| `linux/amd64` | `ffmpeg-bin-linux-amd64.tar.gz` | `ffmpeg-bin-embedded-linux-amd64.tar.gz` |
| `linux/arm64` | `ffmpeg-bin-linux-arm64.tar.gz` | `ffmpeg-bin-embedded-linux-arm64.tar.gz` |
| `darwin/amd64` | `ffmpeg-bin-darwin-amd64.tar.gz` | `ffmpeg-bin-embedded-darwin-amd64.tar.gz` |
| `darwin/arm64` | `ffmpeg-bin-darwin-arm64.tar.gz` | `ffmpeg-bin-embedded-darwin-arm64.tar.gz` |
| `windows/amd64` | `ffmpeg-bin-windows-amd64.zip` | `ffmpeg-bin-embedded-windows-amd64.zip` |
| `windows/386` | `ffmpeg-bin-windows-386.zip` | `ffmpeg-bin-embedded-windows-386.zip` |

## Archive layout

Each archive unpacks into one top-level directory named `<stem>/`, which is the archive filename without its extension (for example `ffmpeg-bin-linux-amd64/`).

Inside `<stem>/`, each archive contains four files:

- `ffmpeg` (`ffmpeg.exe` on Windows)
- `ffprobe` (`ffprobe.exe` on Windows)
- `LICENSE`
- `NOTICE.txt`

## Profiles and licensing

- **`default` (GPLv3)**: Includes the full codec set (`x264`, `x265`, `rav1e`, `dav1d`) and hardware acceleration. Bundles `COPYING.GPLv3` as `LICENSE`.
- **`embedded` (LGPLv2.1)**: Uses software-only OpenH264 with no GPL build flags, omits AV1 and hardware acceleration, and bundles `COPYING.LGPLv2.1` as `LICENSE`. Its `NOTICE.txt` includes LGPLv2.1 Section 4 relinking instructions that reference the matching `ffmpeg-embedded-<os>-<arch>.tar.gz` archive on the same release. This self-built OpenH264 archive does not qualify for the AVC patent royalty coverage that Cisco provides for its separately downloaded binary.

Standalone CLI archives do not use the repository root MIT licence, which applies only to the Go binding code. See [Codec and hardware matrix](Codec-and-Hardware-Matrix) for full codec and hardware details.

## Runtime requirements

| Operating system | Minimum requirement | Dynamic linking behaviour |
| --- | --- | --- |
| Linux | `glibc >= 2.39` (for example Ubuntu 24.04, Debian 13, or Fedora 40) | Static `libstdc++` and `libgcc`. GPU hardware acceleration drivers load through `dlopen` at runtime. |
| macOS | macOS 13.0 or newer | Links system libraries and frameworks only. |
| Windows | Built-in OS DLLs only | `amd64` uses UCRT, and `386` uses MSVCRT. |

## Checksum verification

Each `lib-X.Y.Z.N` release includes one unified `SHA256SUMS` file covering all 24 library and CLI archives. Download your archive and `SHA256SUMS` into the same directory, then verify the checksum for your platform.

Linux:

```sh
sha256sum -c --ignore-missing SHA256SUMS
```

macOS:

```sh
grep ffmpeg-bin-darwin-arm64.tar.gz SHA256SUMS | shasum -a 256 -c -
```

Windows (PowerShell):

```powershell
(Get-FileHash ffmpeg-bin-windows-amd64.zip -Algorithm SHA256).Hash
```

Compare the PowerShell hash output against the matching line in `SHA256SUMS`.

## macOS Gatekeeper

> [!NOTE]
> macOS binaries are ad-hoc signed, not Apple Developer ID notarised. When you download a macOS archive in a browser, remove the Gatekeeper quarantine attribute after extraction:
>
> ```sh
> xattr -dr com.apple.quarantine ffmpeg-bin-darwin-arm64
> ```
>
> Replace `ffmpeg-bin-darwin-arm64` with the extracted directory name for `darwin-amd64` or `embedded` archives.
