# FFmpeg Statigo Wiki

`ffmpeg-statigo` provides static FFmpeg 8.1.3 C library bindings for Go and standalone static `ffmpeg` and `ffprobe` binaries for Linux, macOS, and Windows.

## Users

| Guide | Topics |
| --- | --- |
| [Standalone CLI downloads](Standalone-CLI-Downloads) | Managed package repositories, direct `lib-X.Y.Z.N` CLI archives, checksum verification, and macOS Gatekeeper steps |
| [Codec and hardware matrix](Codec-and-Hardware-Matrix) | Bundled libraries, enabled codecs, hardware acceleration support, and `embedded` profile differences |
| [High-level `av` package](High-Level-av-Package) | Owned `io.Closer` pipeline types (`Input`, `Decoder`, `FilterGraph`, `Encoder`, `Output`) and hardware decoding with `HWDevice` |
| [Go consumer integration](Go-Consumer-Integration) | Submodule setup, consumer `justfile`, GitHub Actions CI matrix, cross-compilation, and `-tags=embedded` builds |

## Contributors

| Guide | Topics |
| --- | --- |
| [Architecture and bindings](Architecture-and-Bindings) | Three-tier binding layout, CGO memory rules, and FFmpeg API coverage |
| [Building from source](Building-From-Source) | `just` recipes, Tailor integration, pinned build dependency digests, and Nix development shells |
| [Binding generator and IR goldens](Binding-Generator-and-IR-Goldens) | Hermetic header parsing with `modernc.org/cc/v4` and intermediate representation (IR) golden tests |
| [Release process and versioning](Release-Process-and-Versioning) | Dual `lib-X.Y.Z.N` and `vX.Y.Z.N` version schemes and GitHub Actions release workflows |
| [Troubleshooting and FAQ](Troubleshooting-and-FAQ) | Linker errors, NixOS CGO setup, module cache cleanup, and common questions |
