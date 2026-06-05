# asciiplayer

Decodes a video file and renders every frame as a greyscale ASCII image in your terminal using [tcell](https://github.com/gdamore/tcell).

Back: [examples/](../README.md)

## What it demonstrates

- Opening a media file with `AVFormatOpenInput` and finding the best video stream
- Building a filter graph manually: a `buffer` source feeds a `scale → fps → settb` chain into a `buffersink` configured for `AVPixFmtGray8`
- Setting sink pixel format options with `AVOptSetSlice` _before_ `AVFilterGraphConfig`, the correct order for non-runtime options
- Decoding with the send/receive packet-frame loop (`AVCodecSendPacket` / `AVCodecReceiveFrame`)
- Pushing decoded frames into the filter graph and pulling filtered frames out with `AVBuffersrcAddFrameFlags` / `AVBuffersinkGetFrame`
- Copying raw luma bytes from a frame with `unsafe.Slice` and mapping greyscale values to terminal colour styles

The filter graph scales the video to the terminal's current character cell dimensions, caps the frame rate at 10 fps, and converts pixel format to 8-bit greyscale. Each character cell becomes one pixel, coloured using the closest terminal palette entry.

## Run signature

```
asciiplayer <file>
```

`<file>` can be any URL that FFmpeg can open: a local file, an HTTP URL, or any other supported protocol.

## Build

From the repo root, inside `nix develop`:

```bash
just build-examples
```

The binary is written to `examples/asciiplayer/asciiplayer`.

> The static libraries must be present first. Run `go run ./cmd/download-lib` if you have not done so already.

## Running

```bash
./examples/asciiplayer/asciiplayer /path/to/video.mp4
```

Resize the terminal before launching. The player queries `tcell` for the screen dimensions at startup and scales the video to fit.

Press `ESC` or `Ctrl-C` to exit.

## Expected output

The terminal clears and plays the video as a greyscale ASCII image. A status line at the top-left shows the current frame number, sleep duration, and the display width/height in characters:

```
Press ESC to exit - frame=42 sleep=12ms w=220 h=50
```

Frame timing is driven by the PTS difference between consecutive filtered frames, so playback runs at approximately the rate the filter graph produces frames (capped at 10 fps). If the terminal is too slow to redraw, frames are dropped rather than buffered.
