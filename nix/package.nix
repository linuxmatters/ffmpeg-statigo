{
  lib,
  stdenvNoCC,
  fetchurl,
  autoPatchelfHook,
  sources ? builtins.fromJSON (builtins.readFile ./sources.json),
}:

let
  system = stdenvNoCC.hostPlatform.system;
  srcInfo =
    sources.systems.${system}
      or (throw "ffmpeg-statigo: unsupported system ${system}");
in
stdenvNoCC.mkDerivation {
  pname = "ffmpeg-statigo";
  version = sources.version;

  src = fetchurl {
    inherit (srcInfo) url hash;
  };

  nativeBuildInputs = lib.optionals stdenvNoCC.hostPlatform.isLinux [
    autoPatchelfHook
  ];

  dontConfigure = true;
  dontBuild = true;

  installPhase = ''
    runHook preInstall

    install -Dm755 ffmpeg "$out/bin/ffmpeg-statigo"
    install -Dm755 ffprobe "$out/bin/ffprobe-statigo"
    install -Dm644 LICENSE "$out/share/licenses/ffmpeg-statigo/LICENSE"
    install -Dm644 NOTICE.txt "$out/share/licenses/ffmpeg-statigo/NOTICE.txt"

    runHook postInstall
  '';

  meta = {
    description = "Static FFmpeg and FFprobe binaries with modern codecs and hardware acceleration";
    homepage = "https://github.com/linuxmatters/ffmpeg-statigo";
    license = lib.licenses.gpl3Plus;
    platforms = builtins.attrNames sources.systems;
    mainProgram = "ffmpeg-statigo";
  };
}
