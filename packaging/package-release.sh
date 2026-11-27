#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
nfpm_config="${script_dir}/nfpm.yaml"
version=""
in_dir=""
out_dir=""
lib_release=""
nix_sources=""

usage() {
  cat <<'EOF'
Usage: package-release.sh --version <ver> --in-dir <dir> --out-dir <dir> [options]

Package the 6 prebuilt default (GPLv3) ffmpeg-bin-* CLI archives into 18
Haberdasher release assets, a unified SHA256SUMS file, and an optional
nix/sources.json manifest.

Required arguments:
  --version <ver>        Release version (e.g. 8.1.2.1 or v8.1.2.1)
  --in-dir <dir>         Directory containing the 6 ffmpeg-bin-* archives
  --out-dir <dir>        Output directory for the 18 assets and SHA256SUMS

Options:
  --lib-release <tag>    Library release tag for nix/sources.json
                         (default: lib-<version>)
  --nix-sources <path>   Write Nix sources JSON manifest to <path>
  --nfpm-config <path>   Path to nfpm.yaml (default: packaging/nfpm.yaml)
  -h, --help             Show this help text
EOF
}

die() {
  printf 'error: %s\n' "$*" >&2
  exit 1
}

while (( $# > 0 )); do
  case "$1" in
    --version)
      (( $# >= 2 )) || die "--version requires an argument"
      version="$2"
      shift 2
      ;;
    --in-dir)
      (( $# >= 2 )) || die "--in-dir requires an argument"
      in_dir="$2"
      shift 2
      ;;
    --out-dir)
      (( $# >= 2 )) || die "--out-dir requires an argument"
      out_dir="$2"
      shift 2
      ;;
    --lib-release)
      (( $# >= 2 )) || die "--lib-release requires an argument"
      lib_release="$2"
      shift 2
      ;;
    --nix-sources)
      (( $# >= 2 )) || die "--nix-sources requires an argument"
      nix_sources="$2"
      shift 2
      ;;
    --nfpm-config)
      (( $# >= 2 )) || die "--nfpm-config requires an argument"
      nfpm_config="$2"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      die "unknown argument: $1"
      ;;
  esac
done

[[ -n "$version" ]] || die "--version is required"
[[ -n "$in_dir" ]] || die "--in-dir is required"
[[ -n "$out_dir" ]] || die "--out-dir is required"

version="${version#v}"
if [[ ! "$version" =~ ^[0-9]+(\.[0-9]+)+$ ]]; then
  die "version must be dot-separated digits (e.g. 8.1.2.1), got: $version"
fi

if [[ -z "$lib_release" ]]; then
  lib_release="lib-${version}"
fi

[[ -d "$in_dir" ]] || die "input directory does not exist: $in_dir"
[[ -f "$nfpm_config" ]] || die "nfpm config does not exist: $nfpm_config"

for tool in nfpm tar python3 sha256sum install; do
  command -v "$tool" >/dev/null 2>&1 || die "required tool not found in PATH: $tool"
done

required_inputs=(
  "ffmpeg-bin-linux-amd64.tar.gz"
  "ffmpeg-bin-linux-arm64.tar.gz"
  "ffmpeg-bin-darwin-amd64.tar.gz"
  "ffmpeg-bin-darwin-arm64.tar.gz"
  "ffmpeg-bin-windows-amd64.zip"
  "ffmpeg-bin-windows-386.zip"
)

for archive in "${required_inputs[@]}"; do
  archive_path="${in_dir}/${archive}"
  [[ -f "$archive_path" && -s "$archive_path" ]] || die "missing or empty input archive: $archive_path"
done

mkdir -p "$out_dir"
work_dir="$(mktemp -d)"
trap 'rm -rf "$work_dir"' EXIT

extract_unix_archive() {
  local os="$1"
  local arch="$2"
  local stem="ffmpeg-bin-${os}-${arch}"
  local archive="${in_dir}/${stem}.tar.gz"
  local dest="${work_dir}/extract/${os}-${arch}"
  local stage_dir="${dest}/${stem}"
  local entry

  mkdir -p "$dest"
  tar -xzf "$archive" -C "$dest"

  for entry in ffmpeg ffprobe LICENSE NOTICE.txt; do
    [[ -f "${stage_dir}/${entry}" && -s "${stage_dir}/${entry}" ]] || \
      die "missing or empty ${stem}/${entry} in ${archive}"
  done
}

package_linux_arch() {
  local arch="$1"
  local pkg_arch="$2"
  local stage_dir="${work_dir}/extract/linux-${arch}/ffmpeg-bin-linux-${arch}"

  extract_unix_archive "linux" "$arch"

  STAGE_DIR="$stage_dir" ARCH="$arch" VERSION="$version" \
    nfpm package -f "$nfpm_config" -p deb \
      -t "${out_dir}/ffmpeg-statigo_${version}_${arch}.deb"

  STAGE_DIR="$stage_dir" ARCH="$arch" VERSION="$version" \
    nfpm package -f "$nfpm_config" -p rpm \
      -t "${out_dir}/ffmpeg-statigo-${version}-1.${pkg_arch}.rpm"

  STAGE_DIR="$stage_dir" ARCH="$arch" VERSION="$version" \
    nfpm package -f "$nfpm_config" -p archlinux \
      -t "${out_dir}/ffmpeg-statigo-${version}-1-${pkg_arch}.pkg.tar.zst"

  install -m 0755 "${stage_dir}/ffmpeg" "${out_dir}/ffmpeg-statigo-linux-${arch}"
  install -m 0755 "${stage_dir}/ffprobe" "${out_dir}/ffprobe-statigo-linux-${arch}"
}

package_darwin_arch() {
  local arch="$1"
  local stage_dir="${work_dir}/extract/darwin-${arch}/ffmpeg-bin-darwin-${arch}"

  extract_unix_archive "darwin" "$arch"

  install -m 0755 "${stage_dir}/ffmpeg" "${out_dir}/ffmpeg-statigo-darwin-${arch}"
  install -m 0755 "${stage_dir}/ffprobe" "${out_dir}/ffprobe-statigo-darwin-${arch}"
}

package_windows_arch() {
  local arch="$1"
  local stem="ffmpeg-bin-windows-${arch}"
  local archive="${in_dir}/${stem}.zip"

  python3 - "$archive" "$stem" "$arch" "$out_dir" <<'PY'
import pathlib
import sys
import zipfile

archive_path = pathlib.Path(sys.argv[1])
stem = sys.argv[2]
arch = sys.argv[3]
out_dir = pathlib.Path(sys.argv[4])

required = ("ffmpeg.exe", "ffprobe.exe", "LICENSE", "NOTICE.txt")
payloads = {}

with zipfile.ZipFile(archive_path, "r") as zin:
    for name in required:
        member = f"{stem}/{name}"
        try:
            data = zin.read(member)
        except KeyError as exc:
            raise SystemExit(f"error: missing {member} in {archive_path}") from exc
        if not data:
            raise SystemExit(f"error: empty {member} in {archive_path}")
        payloads[name] = data

for src_exe, tool in (("ffmpeg.exe", "ffmpeg-statigo"), ("ffprobe.exe", "ffprobe-statigo")):
    out_zip = out_dir / f"{tool}-windows-{arch}.zip"
    entries = (
        ("LICENSE", payloads["LICENSE"]),
        ("NOTICE.txt", payloads["NOTICE.txt"]),
        (f"{tool}.exe", payloads[src_exe]),
    )
    with zipfile.ZipFile(out_zip, "w", compression=zipfile.ZIP_DEFLATED) as zout:
        for arcname, data in entries:
            info = zipfile.ZipInfo(filename=arcname, date_time=(1980, 1, 1, 0, 0, 0))
            info.compress_type = zipfile.ZIP_DEFLATED
            info.external_attr = 0o644 << 16
            zout.writestr(info, data)
PY
}

package_linux_arch "amd64" "x86_64"
package_linux_arch "arm64" "aarch64"

package_darwin_arch "amd64"
package_darwin_arch "arm64"

package_windows_arch "amd64"
package_windows_arch "386"

expected_assets=(
  "ffmpeg-statigo_${version}_amd64.deb"
  "ffmpeg-statigo_${version}_arm64.deb"
  "ffmpeg-statigo-${version}-1.x86_64.rpm"
  "ffmpeg-statigo-${version}-1.aarch64.rpm"
  "ffmpeg-statigo-${version}-1-x86_64.pkg.tar.zst"
  "ffmpeg-statigo-${version}-1-aarch64.pkg.tar.zst"
  "ffmpeg-statigo-linux-amd64"
  "ffmpeg-statigo-linux-arm64"
  "ffprobe-statigo-linux-amd64"
  "ffprobe-statigo-linux-arm64"
  "ffmpeg-statigo-darwin-amd64"
  "ffmpeg-statigo-darwin-arm64"
  "ffprobe-statigo-darwin-amd64"
  "ffprobe-statigo-darwin-arm64"
  "ffmpeg-statigo-windows-amd64.zip"
  "ffmpeg-statigo-windows-386.zip"
  "ffprobe-statigo-windows-amd64.zip"
  "ffprobe-statigo-windows-386.zip"
)

mapfile -t sorted_assets < <(printf '%s\n' "${expected_assets[@]}" | LC_ALL=C sort)

for asset in "${sorted_assets[@]}"; do
  asset_path="${out_dir}/${asset}"
  [[ -f "$asset_path" && -s "$asset_path" ]] || die "missing or empty output asset: $asset_path"
done

(
  cd -- "$out_dir"
  LC_ALL=C sha256sum "${sorted_assets[@]}" > SHA256SUMS
  sha256sum --check --strict SHA256SUMS >/dev/null
)

if [[ -n "$nix_sources" ]]; then
  nix_sources_dir="$(dirname -- "$nix_sources")"
  mkdir -p "$nix_sources_dir"

  python3 - "$version" "$lib_release" "$in_dir" "$nix_sources" <<'PY'
import base64
import hashlib
import json
import pathlib
import sys

version = sys.argv[1]
lib_release = sys.argv[2]
in_dir = pathlib.Path(sys.argv[3])
out_path = pathlib.Path(sys.argv[4])

base_url = f"https://github.com/linuxmatters/ffmpeg-statigo/releases/download/{lib_release}"
system_map = (
    ("x86_64-linux", "linux", "amd64"),
    ("aarch64-linux", "linux", "arm64"),
    ("x86_64-darwin", "darwin", "amd64"),
    ("aarch64-darwin", "darwin", "arm64"),
)

systems = {}
for nix_system, os_name, arch in system_map:
    filename = f"ffmpeg-bin-{os_name}-{arch}.tar.gz"
    digest = hashlib.sha256((in_dir / filename).read_bytes()).digest()
    sri_hash = "sha256-" + base64.b64encode(digest).decode("ascii")
    systems[nix_system] = {
        "url": f"{base_url}/{filename}",
        "hash": sri_hash,
    }

manifest = {
    "version": version,
    "libRelease": lib_release,
    "systems": systems,
}

out_path.write_text(json.dumps(manifest, indent=2) + "\n", encoding="utf-8")
PY
fi
