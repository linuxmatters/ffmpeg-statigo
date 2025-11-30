{
  description = "ffmpeg-statigo development shell";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:

    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        devShells.default = pkgs.mkShell {
          packages =
            with pkgs;
            [
              # Build tools
              autoconf
              automake
              libtool
              cmake
              curl
              ffmpeg
              file
              gcc
              git-filter-repo
              gnumake
              go
              gperf
              just
              meson
              nasm
              ninja
              # LLVM/Clang for code generator (go-clang requires libclang)
              # llvmPackages.libclang provides: clang compiler + libclang library
              # llvmPackages.llvm provides: llvm-config command
              llvmPackages_18.libclang
              llvmPackages_18.llvm
              pkg-config
              python3
              yasm
              # Rust toolchain for rav1e AV1 encoder
              cargo
              cargo-c
              rustc
            ]
            ++ pkgs.lib.optionals pkgs.stdenv.isLinux [
              # Hardware acceleration runtime (Linux only)
              vulkan-loader # Required for Vulkan accelerated encoders
              intel-media-driver # VA-API driver for Intel GPUs (iHD_drv_video.so)
              vpl-gpu-rt # oneVPL runtime for Intel GPUs (QSV backend)
            ];

          # Environment for go-clang CGO compilation and hardware acceleration
          shellHook = ''
            export CGO_LDFLAGS="-L${pkgs.llvmPackages.libclang.lib}/lib"
            export CPATH="${pkgs.llvmPackages.libclang.dev}/include"
            # Ensure vpx build finds yasm
            export PATH="${pkgs.yasm}/bin:${pkgs.nasm}/bin:$PATH"
          ''
          + pkgs.lib.optionalString pkgs.stdenv.isLinux ''
            # Hardware acceleration: Make GPU drivers visible
            # NixOS mounts GPU drivers under /run/opengl-driver/lib
            if [ -d "/run/opengl-driver/lib" ]; then
              if [ -z "$LD_LIBRARY_PATH" ]; then
                export LD_LIBRARY_PATH="/run/opengl-driver/lib"
              else
                export LD_LIBRARY_PATH="/run/opengl-driver/lib:$LD_LIBRARY_PATH"
              fi
            fi
            # Vulkan loader for h264_vulkan, hevc_vulkan, av1_vulkan encoders
            export LD_LIBRARY_PATH="${pkgs.vulkan-loader}/lib:$LD_LIBRARY_PATH"
            # Intel media driver and oneVPL runtime for QuickSync
            export LD_LIBRARY_PATH="${pkgs.intel-media-driver}/lib:$LD_LIBRARY_PATH"
            export LD_LIBRARY_PATH="${pkgs.vpl-gpu-rt}/lib:$LD_LIBRARY_PATH"
            export ONEVPL_SEARCH_PATH="${pkgs.vpl-gpu-rt}/lib"

            # Vulkan ICD discovery: tell vulkan-loader where to find GPU drivers
            # NixOS installs ICDs under /run/opengl-driver/share/vulkan/icd.d/
            if [ -d "/run/opengl-driver/share/vulkan/icd.d" ]; then
              export VK_DRIVER_FILES=$(find /run/opengl-driver/share/vulkan/icd.d -name '*.json' 2>/dev/null | tr '\n' ':' | sed 's/:$//')
            fi
          '';
        };
      }
    );
}
