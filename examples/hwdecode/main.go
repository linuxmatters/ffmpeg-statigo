// Command hwdecode demonstrates hardware-accelerated video decoding through the
// av/ pipeline layer. It opens an input, picks a hardware device, decodes the
// best video stream on it, transfers each frame back to software memory and
// prints frame metadata. With no usable hardware it falls back to software
// decoding so the example still runs.
package main

import (
	"errors"
	"fmt"
	"os"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
	"github.com/linuxmatters/ffmpeg-statigo/av"
)

const maxFrames = 10

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: hwdecode <input> [device-type]")
		os.Exit(2)
	}

	input := os.Args[1]

	in, err := av.Open(input)
	if err != nil {
		return err
	}
	defer in.Close()

	stream, err := in.BestStream(ffmpeg.AVMediaTypeVideo)
	if err != nil {
		return err
	}

	var (
		dev      *av.HWDevice
		dec      *av.Decoder
		hwPixFmt = ffmpeg.AVPixFmtNone
	)

	if len(os.Args) >= 3 {
		dev, err = av.NewHWDevice(os.Args[2])
		if err != nil {
			return err
		}
	} else {
		for _, name := range av.HWDeviceTypes() {
			if d, derr := av.NewHWDevice(name); derr == nil {
				dev = d
				break
			}
		}
	}

	if dev != nil {
		defer dev.Close()

		dec, err = av.NewHWDecoder(stream, dev)
		if err != nil {
			fmt.Printf("hw decoder unavailable for device %v (%v); falling back to software\n", ffmpeg.AVHWDeviceGetTypeName(dev.Type()).String(), err)
		} else {
			hwPixFmt = decoderHWPixFmt(stream, dev)
			fmt.Printf("decoding on hardware device: %v\n", ffmpeg.AVHWDeviceGetTypeName(dev.Type()).String())
		}
	}

	if dec == nil {
		fmt.Println("no usable hardware device; decoding in software")

		dec, err = av.NewDecoder(stream)
		if err != nil {
			return err
		}
	}
	defer dec.Close()

	pkt := ffmpeg.AVPacketAlloc()
	if pkt == nil {
		return errors.New("alloc packet failed")
	}
	defer ffmpeg.AVPacketFree(&pkt)

	frames := 0
	emit := func(frame *ffmpeg.AVFrame) error {
		if frames >= maxFrames {
			return nil
		}

		if ffmpeg.AVPixelFormat(frame.Format()) == hwPixFmt && hwPixFmt != ffmpeg.AVPixFmtNone { //nolint:gosec // G115: AVFrame.Format holds a small AVPixelFormat enum value
			sw, terr := av.TransferToSoftware(frame)
			if terr != nil {
				return terr
			}

			printFrame(frames, sw)
			ffmpeg.AVFrameFree(&sw)
		} else {
			printFrame(frames, frame)
		}

		frames++

		return nil
	}

	for frames < maxFrames {
		if err := in.ReadPacket(pkt); err != nil {
			if errors.Is(err, ffmpeg.AVErrorEOF) {
				break
			}

			return err
		}

		if pkt.StreamIndex() == stream.Index() {
			if err := dec.Decode(pkt, emit); err != nil {
				ffmpeg.AVPacketUnref(pkt)
				return err
			}
		}

		ffmpeg.AVPacketUnref(pkt)
	}

	if err := dec.Flush(emit); err != nil {
		return err
	}

	return nil
}

// decoderHWPixFmt returns the hardware pixel format the decoder for stream uses
// on dev, or AVPixFmtNone if none matches.
func decoderHWPixFmt(stream *ffmpeg.AVStream, dev *av.HWDevice) ffmpeg.AVPixelFormat {
	codec := ffmpeg.AVCodecFindDecoder(stream.Codecpar().CodecId())
	if codec == nil {
		return ffmpeg.AVPixFmtNone
	}

	for i := 0; ; i++ {
		config := ffmpeg.AVCodecGetHWConfig(codec, i)
		if config == nil {
			return ffmpeg.AVPixFmtNone
		}

		if config.Methods()&ffmpeg.AVCodecHWConfigMethodHWDeviceCtx != 0 && config.DeviceType() == dev.Type() {
			return config.PixFmt()
		}
	}
}

// printFrame writes a one-line summary of frame at index i.
func printFrame(i int, frame *ffmpeg.AVFrame) {
	fmt.Printf("frame %d: pts=%d %dx%d fmt=%s\n",
		i,
		frame.Pts(),
		frame.Width(),
		frame.Height(),
		ffmpeg.AVGetPixFmtName(ffmpeg.AVPixelFormat(frame.Format())).String(), //nolint:gosec // G115: AVFrame.Format holds a small AVPixelFormat enum value
	)
}
