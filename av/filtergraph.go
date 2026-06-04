package av

import (
	"errors"
	"fmt"
	"io"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
)

// chLayoutBufSize bounds the channel-layout description buffer, matching the
// example transcoder's 64-byte allocation.
const chLayoutBufSize = 64

// VideoFilterParams describes the source frames feeding a video FilterGraph and
// the single pixel format the sink must emit. Mirror these from the decoder
// context (see [ffmpeg.AVCodecContext]).
type VideoFilterParams struct {
	Width             int
	Height            int
	PixFmt            ffmpeg.AVPixelFormat
	TimeBase          *ffmpeg.AVRational
	SampleAspectRatio *ffmpeg.AVRational
	OutPixFmt         ffmpeg.AVPixelFormat
}

// AudioFilterParams describes the source frames feeding an audio FilterGraph and
// the format constraints the sink must emit.
type AudioFilterParams struct {
	TimeBase      *ffmpeg.AVRational
	SampleRate    int
	SampleFmt     ffmpeg.AVSampleFormat
	ChLayout      *ffmpeg.AVChannelLayout
	OutSampleFmt  ffmpeg.AVSampleFormat
	OutSampleRate int
	OutChLayout   *ffmpeg.AVChannelLayout
}

// FilterGraph owns an [ffmpeg.AVFilterGraph] plus its single buffersrc and
// buffersink contexts. It handles the common single-input/single-output video
// and audio cases through one cohesive type with an internal media branch.
// FilterGraph implements [io.Closer]; Close is idempotent and frees the graph.
// Reach any unsurfaced capability through Raw.
type FilterGraph struct {
	graph   *ffmpeg.AVFilterGraph
	srcCtx  *ffmpeg.AVFilterContext
	sinkCtx *ffmpeg.AVFilterContext
	scratch *ffmpeg.AVFrame
}

var _ io.Closer = (*FilterGraph)(nil)

// NewVideoFilterGraph builds a single-in/single-out video graph applying spec
// (for example "null"). On any error after allocation it unwinds via Close so
// no handle leaks.
func NewVideoFilterGraph(params VideoFilterParams, spec string) (*FilterGraph, error) {
	g, err := newFilterGraph()
	if err != nil {
		return nil, err
	}

	args := fmt.Sprintf(
		"video_size=%dx%d:pix_fmt=%d:time_base=%d/%d:pixel_aspect=%d/%d",
		params.Width, params.Height,
		params.PixFmt,
		params.TimeBase.Num(), params.TimeBase.Den(),
		params.SampleAspectRatio.Num(), params.SampleAspectRatio.Den(),
	)

	if err := g.createEndpoints("buffer", "buffersink", args); err != nil {
		_ = g.Close()
		return nil, err
	}

	pixFmts := []ffmpeg.AVPixelFormat{params.OutPixFmt}
	if _, err := ffmpeg.AVOptSetSlice(g.sinkCtx.RawPtr(), ffmpeg.GlobalCStr("pix_fmts"), pixFmts, ffmpeg.AVOptSearchChildren); err != nil {
		_ = g.Close()
		return nil, fmt.Errorf("set sink pix_fmts: %w", err)
	}

	if err := g.initSink(); err != nil {
		_ = g.Close()
		return nil, err
	}

	if err := g.parseAndConfigure(spec); err != nil {
		_ = g.Close()
		return nil, err
	}

	return g, nil
}

// NewAudioFilterGraph builds a single-in/single-out audio graph applying spec
// (for example "anull"). On any error after allocation it unwinds via Close.
func NewAudioFilterGraph(params AudioFilterParams, spec string) (*FilterGraph, error) {
	g, err := newFilterGraph()
	if err != nil {
		return nil, err
	}

	srcLayout, err := describeLayout(params.ChLayout)
	if err != nil {
		_ = g.Close()
		return nil, err
	}

	args := fmt.Sprintf(
		"time_base=%d/%d:sample_rate=%d:sample_fmt=%s:channel_layout=%s",
		params.TimeBase.Num(), params.TimeBase.Den(),
		params.SampleRate,
		ffmpeg.AVGetSampleFmtName(params.SampleFmt).String(),
		srcLayout,
	)

	if err := g.createEndpoints("abuffer", "abuffersink", args); err != nil {
		_ = g.Close()
		return nil, err
	}

	sampleFmts := []ffmpeg.AVSampleFormat{params.OutSampleFmt}
	if _, err := ffmpeg.AVOptSetSlice(g.sinkCtx.RawPtr(), ffmpeg.GlobalCStr("sample_fmts"), sampleFmts, ffmpeg.AVOptSearchChildren); err != nil {
		_ = g.Close()
		return nil, fmt.Errorf("set sink sample_fmts: %w", err)
	}

	if params.OutSampleRate != 0 {
		sampleRates := []int{params.OutSampleRate}
		if _, err := ffmpeg.AVOptSetSlice(g.sinkCtx.RawPtr(), ffmpeg.GlobalCStr("sample_rates"), sampleRates, ffmpeg.AVOptSearchChildren); err != nil {
			_ = g.Close()
			return nil, fmt.Errorf("set sink sample_rates: %w", err)
		}
	}

	if params.OutChLayout != nil {
		outLayout, err := describeLayout(params.OutChLayout)
		if err != nil {
			_ = g.Close()
			return nil, err
		}

		layoutC := ffmpeg.ToCStr(outLayout)
		defer layoutC.Free()

		if _, err := ffmpeg.AVOptSet(g.sinkCtx.RawPtr(), ffmpeg.GlobalCStr("ch_layouts"), layoutC, ffmpeg.AVOptSearchChildren); err != nil {
			_ = g.Close()
			return nil, fmt.Errorf("set sink ch_layouts: %w", err)
		}
	}

	if err := g.initSink(); err != nil {
		_ = g.Close()
		return nil, err
	}

	if err := g.parseAndConfigure(spec); err != nil {
		_ = g.Close()
		return nil, err
	}

	return g, nil
}

// NewVideoFilterGraphFromContext builds a video graph reading the source
// parameters from a decoder context, mirroring the example transcoder.
func NewVideoFilterGraphFromContext(decCtx *ffmpeg.AVCodecContext, spec string, outPixFmt ffmpeg.AVPixelFormat) (*FilterGraph, error) {
	if decCtx == nil {
		return nil, errors.New("nil decoder context")
	}

	return NewVideoFilterGraph(VideoFilterParams{
		Width:             decCtx.Width(),
		Height:            decCtx.Height(),
		PixFmt:            decCtx.PixFmt(),
		TimeBase:          decCtx.PktTimebase(),
		SampleAspectRatio: decCtx.SampleAspectRatio(),
		OutPixFmt:         outPixFmt,
	}, spec)
}

// NewAudioFilterGraphFromContext builds an audio graph reading the source
// parameters from a decoder context. outChLayout may be nil to leave the sink
// layout unconstrained.
func NewAudioFilterGraphFromContext(decCtx *ffmpeg.AVCodecContext, spec string, outSampleFmt ffmpeg.AVSampleFormat, outSampleRate int, outChLayout *ffmpeg.AVChannelLayout) (*FilterGraph, error) {
	if decCtx == nil {
		return nil, errors.New("nil decoder context")
	}

	chLayout := decCtx.ChLayout()
	if chLayout.Order() == ffmpeg.AVChannelOrderUnspec {
		ffmpeg.AVChannelLayoutDefault(chLayout, chLayout.NbChannels())
	}

	return NewAudioFilterGraph(AudioFilterParams{
		TimeBase:      decCtx.PktTimebase(),
		SampleRate:    decCtx.SampleRate(),
		SampleFmt:     decCtx.SampleFmt(),
		ChLayout:      chLayout,
		OutSampleFmt:  outSampleFmt,
		OutSampleRate: outSampleRate,
		OutChLayout:   outChLayout,
	}, spec)
}

// Push sends a frame into the graph's buffersrc.
func (g *FilterGraph) Push(frame *ffmpeg.AVFrame) error {
	if _, err := ffmpeg.AVBuffersrcAddFrameFlags(g.srcCtx, frame, ffmpeg.AVBuffersrcFlagKeepRef); err != nil {
		return fmt.Errorf("push frame to filter graph: %w", err)
	}

	return nil
}

// Pull drains all filtered frames currently available from the buffersink,
// calling fn once per frame. The frame passed to fn is valid only for the
// callback's duration; Pull unrefs the scratch frame after each call. EAgain
// and end-of-stream terminate the loop without leaking to the caller.
func (g *FilterGraph) Pull(fn func(*ffmpeg.AVFrame) error) error {
	for {
		if _, err := ffmpeg.AVBuffersinkGetFrame(g.sinkCtx, g.scratch); err != nil {
			if errors.Is(err, ffmpeg.EAgain) || errors.Is(err, ffmpeg.AVErrorEOF) {
				return nil
			}

			return fmt.Errorf("pull frame from filter graph: %w", err)
		}

		cbErr := fn(g.scratch)
		ffmpeg.AVFrameUnref(g.scratch)

		if cbErr != nil {
			return cbErr
		}
	}
}

// Raw returns the underlying filter graph. Do not free it; Close owns it.
func (g *FilterGraph) Raw() *ffmpeg.AVFilterGraph {
	return g.graph
}

// Close frees the filter graph and nils the handle. A second call is a no-op.
func (g *FilterGraph) Close() error {
	if g.graph == nil {
		return nil
	}

	if g.scratch != nil {
		ffmpeg.AVFrameFree(&g.scratch)
		g.scratch = nil
	}

	ffmpeg.AVFilterGraphFree(&g.graph)
	g.graph = nil
	g.srcCtx = nil
	g.sinkCtx = nil

	return nil
}

func newFilterGraph() (*FilterGraph, error) {
	graph := ffmpeg.AVFilterGraphAlloc()
	if graph == nil {
		return nil, fmt.Errorf("allocate filter graph")
	}

	scratch := ffmpeg.AVFrameAlloc()
	if scratch == nil {
		ffmpeg.AVFilterGraphFree(&graph)
		return nil, fmt.Errorf("allocate scratch frame")
	}

	return &FilterGraph{graph: graph, scratch: scratch}, nil
}

// createEndpoints creates the initialised buffersrc and allocates the
// buffersink uninitialised, so format options can be set before initSink runs.
func (g *FilterGraph) createEndpoints(srcName, sinkName, srcArgs string) error {
	src := ffmpeg.AVFilterGetByName(ffmpeg.GlobalCStr(srcName))
	sink := ffmpeg.AVFilterGetByName(ffmpeg.GlobalCStr(sinkName))
	if src == nil || sink == nil {
		return fmt.Errorf("filter %q/%q not found", srcName, sinkName)
	}

	argsC := ffmpeg.ToCStr(srcArgs)
	defer argsC.Free()

	if _, err := ffmpeg.AVFilterGraphCreateFilter(&g.srcCtx, src, ffmpeg.GlobalCStr("in"), argsC, nil, g.graph); err != nil {
		return fmt.Errorf("create buffersrc: %w", err)
	}

	g.sinkCtx = ffmpeg.AVFilterGraphAllocFilter(g.graph, sink, ffmpeg.GlobalCStr("out"))
	if g.sinkCtx == nil {
		return fmt.Errorf("allocate buffersink")
	}

	return nil
}

// initSink initialises the buffersink after its format options are set.
func (g *FilterGraph) initSink() error {
	if _, err := ffmpeg.AVFilterInitStr(g.sinkCtx, nil); err != nil {
		return fmt.Errorf("init buffersink: %w", err)
	}

	return nil
}

// parseAndConfigure wires the buffersrc to the buffersink through spec, then
// validates the graph. It frees the inout descriptors it allocates.
func (g *FilterGraph) parseAndConfigure(spec string) error {
	outputs := ffmpeg.AVFilterInoutAlloc()
	inputs := ffmpeg.AVFilterInoutAlloc()
	if outputs == nil || inputs == nil {
		if outputs != nil {
			ffmpeg.AVFilterInoutFree(&outputs)
		}
		if inputs != nil {
			ffmpeg.AVFilterInoutFree(&inputs)
		}
		return fmt.Errorf("allocate filter inout")
	}

	defer ffmpeg.AVFilterInoutFree(&outputs)
	defer ffmpeg.AVFilterInoutFree(&inputs)

	outputs.SetName(ffmpeg.ToCStr("in"))
	outputs.SetFilterCtx(g.srcCtx)
	outputs.SetPadIdx(0)
	outputs.SetNext(nil)

	inputs.SetName(ffmpeg.ToCStr("out"))
	inputs.SetFilterCtx(g.sinkCtx)
	inputs.SetPadIdx(0)
	inputs.SetNext(nil)

	specC := ffmpeg.ToCStr(spec)
	defer specC.Free()

	if _, err := ffmpeg.AVFilterGraphParsePtr(g.graph, specC, &inputs, &outputs, nil); err != nil {
		return fmt.Errorf("parse filter spec %q: %w", spec, err)
	}

	if _, err := ffmpeg.AVFilterGraphConfig(g.graph, nil); err != nil {
		return fmt.Errorf("configure filter graph: %w", err)
	}

	return nil
}

// describeLayout renders a channel layout to its FFmpeg string form.
func describeLayout(layout *ffmpeg.AVChannelLayout) (string, error) {
	buf := ffmpeg.AllocCStr(chLayoutBufSize)
	defer buf.Free()

	if _, err := ffmpeg.AVChannelLayoutDescribe(layout, buf, chLayoutBufSize); err != nil {
		return "", fmt.Errorf("describe channel layout: %w", err)
	}

	return buf.String(), nil
}
