// High-level port of doc/examples/transcode.c using the av package wrappers.
package main

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
	"github.com/linuxmatters/ffmpeg-statigo/av"
)

// stream holds the per-input-stream pipeline. For audio and video it transcodes
// through decoder -> filter graph -> encoder; for every other media type it
// remuxes by copying packets.
type stream struct {
	index     int
	transcode bool

	decoder *av.Decoder
	filter  *av.FilterGraph
	encoder *av.Encoder

	inStream  *ffmpeg.AVStream
	outStream *ffmpeg.AVStream
}

func main() {
	slog.Info("Transcode")

	if len(os.Args) < 3 {
		log.Fatal("usage: transcode-hl <input> <output>")
	}

	if err := run(os.Args[1], os.Args[2]); err != nil {
		log.Fatal(err)
	}
}

func run(inPath, outPath string) error {
	input, err := av.Open(inPath)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := av.CreateOutput(outPath)
	if err != nil {
		return err
	}
	defer output.Close()

	globalHeader := output.Raw().Oformat().Flags()&ffmpeg.AVFmtGlobalheader != 0

	inStreams := input.Streams()
	streams := make([]*stream, input.NbStreams())

	for i := range streams {
		s, err := setupStream(input, output, inStreams.Get(uintptr(i)), i, globalHeader)
		if err != nil {
			return err
		}
		defer s.close()

		streams[i] = s
	}

	if err := output.WriteHeader(); err != nil {
		return err
	}

	packet := ffmpeg.AVPacketAlloc()
	defer ffmpeg.AVPacketFree(&packet)

	for {
		if err := input.ReadPacket(packet); err != nil {
			if errors.Is(err, ffmpeg.AVErrorEOF) {
				slog.Info("End of file")
				break
			}
			return err
		}

		s := streams[packet.StreamIndex()]
		err := s.process(packet, output)
		ffmpeg.AVPacketUnref(packet)
		if err != nil {
			return err
		}
	}

	for _, s := range streams {
		if err := s.flush(output); err != nil {
			return err
		}
	}

	return output.WriteTrailer()
}

// setupStream builds the decode/filter/encode pipeline for audio and video, or
// copies codec parameters for a remuxed stream, matching examples/transcode.
func setupStream(input *av.Input, output *av.Output, in *ffmpeg.AVStream, index int, globalHeader bool) (*stream, error) {
	s := &stream{index: index, inStream: in}

	dec, err := av.NewDecoder(in)
	if err != nil {
		return nil, err
	}
	s.decoder = dec
	decCtx := dec.Raw()

	mediaType := decCtx.CodecType()

	switch mediaType {
	case ffmpeg.AVMediaTypeVideo:
		decCtx.SetFramerate(ffmpeg.AVGuessFrameRate(input.Raw(), in, nil))
	case ffmpeg.AVMediaTypeAudio:
	case ffmpeg.AVMediaTypeUnknown:
		s.close()
		return nil, fmt.Errorf("unknown media type for stream %d", index)
	default:
		// Remux: copy codec parameters from the input stream.
		s.decoder.Close()
		s.decoder = nil

		outStream, err := remuxStream(output, in)
		if err != nil {
			s.close()
			return nil, err
		}
		s.outStream = outStream
		return s, nil
	}

	s.transcode = true

	codec := ffmpeg.AVCodecFindEncoder(decCtx.CodecId())
	if codec == nil {
		s.close()
		return nil, fmt.Errorf("failed to find encoder for stream %d", index)
	}

	var layoutErr error
	enc, err := av.NewEncoder(codec, func(encCtx *ffmpeg.AVCodecContext) {
		if mediaType == ffmpeg.AVMediaTypeVideo {
			encCtx.SetHeight(decCtx.Height())
			encCtx.SetWidth(decCtx.Width())
			encCtx.SetSampleAspectRatio(decCtx.SampleAspectRatio())

			if fmts := codec.PixFmts(); fmts != nil {
				encCtx.SetPixFmt(fmts.Get(0))
			} else {
				encCtx.SetPixFmt(decCtx.PixFmt())
			}

			encCtx.SetTimeBase(ffmpeg.AVInvQ(decCtx.Framerate()))
		} else {
			encCtx.SetSampleRate(decCtx.SampleRate())
			if _, cerr := ffmpeg.AVChannelLayoutCopy(encCtx.ChLayout(), decCtx.ChLayout()); cerr != nil {
				layoutErr = cerr
				return
			}
			encCtx.SetSampleFmt(codec.SampleFmts().Get(0))
			encCtx.SetTimeBase(ffmpeg.AVMakeQ(1, encCtx.SampleRate()))
		}

		if globalHeader {
			encCtx.SetFlags(encCtx.Flags() | ffmpeg.AVCodecFlagGlobalHeader)
		}
	})
	if err == nil && layoutErr != nil {
		err = fmt.Errorf("copy channel layout: %w", layoutErr)
	}
	if err != nil {
		s.close()
		return nil, err
	}
	s.encoder = enc

	outStream, err := output.AddStream(enc)
	if err != nil {
		s.close()
		return nil, err
	}
	s.outStream = outStream

	if mediaType == ffmpeg.AVMediaTypeVideo {
		s.filter, err = av.NewVideoFilterGraphFromContext(decCtx, "null", enc.Raw().PixFmt())
	} else {
		s.filter, err = av.NewAudioFilterGraphFromContext(decCtx, "anull", enc.Raw().SampleFmt(), enc.Raw().SampleRate(), enc.Raw().ChLayout())
	}
	if err != nil {
		s.close()
		return nil, err
	}

	return s, nil
}

// remuxStream adds an output stream copying the input stream's codec parameters
// and time base. This is the one forced Raw() drop-through: the av package has
// no stream-copy wrapper.
func remuxStream(output *av.Output, in *ffmpeg.AVStream) (*ffmpeg.AVStream, error) {
	outStream := ffmpeg.AVFormatNewStream(output.Raw(), nil)
	if outStream == nil {
		return nil, errors.New("failed to allocate remux output stream")
	}

	if _, err := ffmpeg.AVCodecParametersCopy(outStream.Codecpar(), in.Codecpar()); err != nil {
		return nil, err
	}

	outStream.SetTimeBase(in.TimeBase())

	return outStream, nil
}

// process handles one input packet: decode -> filter -> encode -> mux for
// transcoded streams, or rescale-and-write for remuxed streams.
func (s *stream) process(packet *ffmpeg.AVPacket, output *av.Output) error {
	if !s.transcode {
		ffmpeg.AVPacketRescaleTs(packet, s.inStream.TimeBase(), s.outStream.TimeBase())
		return output.WritePacket(packet)
	}

	return s.decoder.Decode(packet, func(frame *ffmpeg.AVFrame) error {
		frame.SetPts(frame.BestEffortTimestamp())
		return s.filterEncode(frame, output)
	})
}

// filterEncode pushes a decoded frame through the filter graph then encodes and
// muxes every filtered frame.
func (s *stream) filterEncode(frame *ffmpeg.AVFrame, output *av.Output) error {
	if err := s.filter.Push(frame); err != nil {
		return err
	}

	return s.filter.Pull(func(filtered *ffmpeg.AVFrame) error {
		filtered.SetTimeBase(s.sinkTimeBase())
		filtered.SetPictType(ffmpeg.AVPictureTypeNone)
		return s.encodeWrite(filtered, output)
	})
}

// sinkTimeBase reads the buffersink time base so the filtered frame's PTS can be
// rescaled into the encoder time base. The av FilterGraph wrapper does not
// surface this, so reach the sink (named "out" by the wrapper) through Raw.
func (s *stream) sinkTimeBase() *ffmpeg.AVRational {
	sink := ffmpeg.AVFilterGraphGetFilter(s.filter.Raw(), ffmpeg.GlobalCStr("out"))
	return ffmpeg.AVBuffersinkGetTimeBase(sink)
}

// encodeWrite encodes one frame (or flushes when frame is nil) and muxes each
// resulting packet.
func (s *stream) encodeWrite(frame *ffmpeg.AVFrame, output *av.Output) error {
	encCtx := s.encoder.Raw()

	if frame != nil && frame.Pts() != ffmpeg.AVNoptsValue {
		frame.SetPts(ffmpeg.AVRescaleQ(frame.Pts(), frame.TimeBase(), encCtx.TimeBase()))
	}

	mux := func(pkt *ffmpeg.AVPacket) error {
		pkt.SetStreamIndex(s.index)
		ffmpeg.AVPacketRescaleTs(pkt, encCtx.TimeBase(), s.outStream.TimeBase())
		return output.WritePacket(pkt)
	}

	if frame == nil {
		return s.encoder.Flush(mux)
	}

	return s.encoder.Encode(frame, mux)
}

// flush drains the decoder, filter graph and encoder at end of stream, matching
// the original's tail handling.
func (s *stream) flush(output *av.Output) error {
	if !s.transcode {
		return nil
	}

	if err := s.decoder.Flush(func(frame *ffmpeg.AVFrame) error {
		frame.SetPts(frame.BestEffortTimestamp())
		return s.filterEncode(frame, output)
	}); err != nil {
		return err
	}

	// Flush the filter graph: push a nil frame, then drain and encode.
	if err := s.filter.Push(nil); err != nil {
		return err
	}
	if err := s.filter.Pull(func(filtered *ffmpeg.AVFrame) error {
		filtered.SetTimeBase(s.sinkTimeBase())
		filtered.SetPictType(ffmpeg.AVPictureTypeNone)
		return s.encodeWrite(filtered, output)
	}); err != nil {
		return err
	}

	if s.encoder.Raw().Codec().Capabilities()&ffmpeg.AVCodecCapDelay != 0 {
		if err := s.encodeWrite(nil, output); err != nil {
			return err
		}
	}

	return nil
}

func (s *stream) close() {
	if s.filter != nil {
		s.filter.Close()
	}
	if s.encoder != nil {
		s.encoder.Close()
	}
	if s.decoder != nil {
		s.decoder.Close()
	}
}
