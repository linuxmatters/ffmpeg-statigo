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
            ++ pkgs.lib.optionals pkgs.stdenv.isDarwin [
              # C++ standard library headers for building C++ dependencies (zimg, etc.)
              # The .dev output contains include/c++/v1/ with <algorithm>, <iostream>, etc.
              llvmPackages_18.libcxx
            ]
            ++ pkgs.lib.optionals pkgs.stdenv.isLinux [
              # Hardware acceleration runtime (Linux only)
              vulkan-loader # Required for Vulkan accelerated encoders
              intel-media-driver # VA-API driver for Intel GPUs (iHD_drv_video.so)
              vpl-gpu-rt # oneVPL runtime for 11th gen+ Intel (Tiger Lake+) QSV
            ];

          # Environment for go-clang CGO compilation and hardware acceleration
          shellHook = ''
            export CGO_LDFLAGS="-L${pkgs.llvmPackages_18.libclang.lib}/lib"
            export CPATH="${pkgs.llvmPackages_18.libclang.dev}/include"
            # Ensure vpx build finds yasm
            export PATH="${pkgs.yasm}/bin:${pkgs.nasm}/bin:$PATH"
          ''
          + pkgs.lib.optionalString pkgs.stdenv.isDarwin ''
            # macOS: Export SDK paths for C/C++ compilation
            # The SDKROOT is set by Nix's stdenv wrapper, use it or fall back to xcrun
            # CGO needs both the SDK path and clang's builtin headers (stdarg.h, stddef.h, etc.)
            export CGO_CFLAGS="-isysroot ''${SDKROOT:-$(xcrun --show-sdk-path)} -I${pkgs.llvmPackages_18.libclang.lib}/lib/clang/18/include"
            export CPATH="${pkgs.llvmPackages_18.libclang.dev}/include:$CPATH"
            # C++ standard library headers path for building C++ dependencies (zimg, etc.)
            # The builder uses this to add -I flag for <algorithm>, <iostream>, etc.
            export LIBCXX_INCLUDE="${pkgs.llvmPackages_18.libcxx.dev}/include/c++/v1"
            # Set deployment target to match ffmpeg-statigo build (macOS 13.0+)
            export MACOSX_DEPLOYMENT_TARGET="13.0"
            # macOS uses DYLD_LIBRARY_PATH instead of LD_LIBRARY_PATH
            export DYLD_LIBRARY_PATH="${pkgs.llvmPackages_18.libclang.lib}/lib:''${DYLD_LIBRARY_PATH:-}"
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
            # oneVPL search path for QSV (11th gen+ Intel only)
            export ONEVPL_SEARCH_PATH="${pkgs.vpl-gpu-rt}/lib"
            # VA-API driver discovery for libva
            # Use system drivers if available, fall back to nix package for Intel
            if [ -d "/run/opengl-driver/lib/dri" ]; then
              export LIBVA_DRIVERS_PATH="/run/opengl-driver/lib/dri"
            fi
            # Auto-detect VA-API driver based on GPU vendor (prefer Intel for VA-API)
            if lspci -d ::0300 2>/dev/null | grep -qi intel; then
              export LIBVA_DRIVER_NAME="iHD"
              # Ensure Intel driver path is set even without system drivers
              export LIBVA_DRIVERS_PATH="${pkgs.intel-media-driver}/lib/dri:''${LIBVA_DRIVERS_PATH:-}"
            elif lspci -d ::0300 2>/dev/null | grep -qi amd; then
              export LIBVA_DRIVER_NAME="radeonsi"
            elif lspci -d ::0300 2>/dev/null | grep -qi nvidia; then
              export LIBVA_DRIVER_NAME="nvidia"
            fi

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
