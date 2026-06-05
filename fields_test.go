package ffmpeg

import (
	"testing"
)

// TestQuantMatricesNil asserts the AVCodecContext quant-matrix getters return
// nil when their backing pointers are unset, the state of a freshly allocated
// context.
func TestQuantMatricesNil(t *testing.T) {
	ctx := AVCodecAllocContext3(nil)
	if ctx == nil {
		t.Fatal("AVCodecAllocContext3 returned nil")
	}
	defer AVCodecFreeContext(&ctx)

	if got := ctx.IntraMatrix(); got != nil {
		t.Errorf("IntraMatrix on fresh context = %v, want nil", got)
	}
	if got := ctx.InterMatrix(); got != nil {
		t.Errorf("InterMatrix on fresh context = %v, want nil", got)
	}
	if got := ctx.ChromaIntraMatrix(); got != nil {
		t.Errorf("ChromaIntraMatrix on fresh context = %v, want nil", got)
	}
}

// TestExtendedData asserts AVFrame.ExtendedData, after allocating planar-audio
// buffers, aliases the per-channel plane pointers without freeing them.
func TestExtendedData(t *testing.T) {
	frame := AVFrameAlloc()
	if frame == nil {
		t.Fatal("AVFrameAlloc returned nil")
	}
	defer AVFrameFree(&frame)

	const nbChannels = 2
	layout := newStereoLayout(t)
	defer freeChannelLayout(layout)
	if _, err := AVChannelLayoutCopy(frame.ChLayout(), layout); err != nil {
		t.Fatalf("AVChannelLayoutCopy failed: %v", err)
	}
	frame.SetFormat(int(AVSampleFmtFltp))
	frame.SetNbSamples(1024)

	if ret, err := AVFrameGetBuffer(frame, 0); err != nil {
		t.Fatalf("AVFrameGetBuffer failed: %v (ret=%d)", err, ret)
	}

	ext := frame.ExtendedData()
	if ext == nil {
		t.Fatal("ExtendedData after buffer alloc = nil, want non-nil")
	}
	// Planar FLTP holds one plane pointer per channel; assert each is non-nil and
	// that extended_data[0] aliases data[0] (the documented relationship).
	data := frame.Data()
	for i := range nbChannels {
		if ext.Get(uintptr(i)) == nil {
			t.Errorf("extended_data[%d] is nil, want a plane pointer", i)
		}
	}
	if ext.Get(0) != data.Get(0) {
		t.Errorf("extended_data[0]=%p does not alias data[0]=%p", ext.Get(0), data.Get(0))
	}
}

// TestComp asserts AVPixFmtDescriptor.Comp aliases the inline comp[4] array and
// bounds-checks its index.
func TestComp(t *testing.T) {
	desc := AVPixFmtDescGet(AVPixFmtYuv420P)
	if desc == nil {
		t.Fatal("AVPixFmtDescGet(YUV420P) returned nil")
	}

	n := int(desc.NbComponents())
	if n != 3 {
		t.Fatalf("YUV420P NbComponents = %d, want 3", n)
	}

	// Component 0 is the luma plane: plane 0, depth 8.
	c0 := desc.Comp(0)
	if c0 == nil {
		t.Fatal("Comp(0) = nil, want a descriptor")
	}
	if c0.Plane() != 0 {
		t.Errorf("Comp(0).Plane() = %d, want 0", c0.Plane())
	}
	if c0.Depth() != 8 {
		t.Errorf("Comp(0).Depth() = %d, want 8", c0.Depth())
	}

	// Chroma components live on planes 1 and 2.
	if p := desc.Comp(1).Plane(); p != 1 {
		t.Errorf("Comp(1).Plane() = %d, want 1", p)
	}
	if p := desc.Comp(2).Plane(); p != 2 {
		t.Errorf("Comp(2).Plane() = %d, want 2", p)
	}

	if got := desc.Comp(-1); got != nil {
		t.Errorf("Comp(-1) = %v, want nil", got)
	}
	if got := desc.Comp(4); got != nil {
		t.Errorf("Comp(4) = %v, want nil", got)
	}
}

// TestDisplayPrimaries round-trips the AVRational display_primaries[3][2] array
// through the row accessor, asserting it aliases writable struct-owned memory and
// bounds-checks its index.
func TestDisplayPrimaries(t *testing.T) {
	mdm := allocMasteringDisplayMetadata()
	if mdm == nil {
		t.Fatal("allocMasteringDisplayMetadata returned nil")
	}
	defer freeMasteringDisplayMetadata(mdm)

	if got := mdm.DisplayPrimaries(-1); got != nil {
		t.Errorf("DisplayPrimaries(-1) = %v, want nil", got)
	}
	if got := mdm.DisplayPrimaries(3); got != nil {
		t.Errorf("DisplayPrimaries(3) = %v, want nil", got)
	}

	for row := range 3 {
		arr := mdm.DisplayPrimaries(row)
		if arr == nil {
			t.Fatalf("DisplayPrimaries(%d) = nil, want a 2-element array", row)
		}
		wantX := &AVRational{}
		wantX.SetNum(row + 1)
		wantX.SetDen(2)
		wantY := &AVRational{}
		wantY.SetNum(row + 3)
		wantY.SetDen(4)
		arr.Set(0, wantX)
		arr.Set(1, wantY)

		got := mdm.DisplayPrimaries(row)
		if got.Get(0).Num() != row+1 || got.Get(0).Den() != 2 {
			t.Errorf("row %d x = %d/%d, want %d/2", row, got.Get(0).Num(), got.Get(0).Den(), row+1)
		}
		if got.Get(1).Num() != row+3 || got.Get(1).Den() != 4 {
			t.Errorf("row %d y = %d/%d, want %d/4", row, got.Get(1).Num(), got.Get(1).Den(), row+3)
		}
	}
}

// TestPosition round-trips the int16_t position[3][2] array through the row
// accessor, asserting it aliases writable struct-owned memory and bounds-checks
// its index.
func TestPosition(t *testing.T) {
	ps := allocPanScan()
	if ps == nil {
		t.Fatal("allocPanScan returned nil")
	}
	defer freePanScan(ps)

	if got := ps.Position(-1); got != nil {
		t.Errorf("Position(-1) = %v, want nil", got)
	}
	if got := ps.Position(3); got != nil {
		t.Errorf("Position(3) = %v, want nil", got)
	}

	for row := range 3 {
		arr := ps.Position(row)
		if arr == nil {
			t.Fatalf("Position(%d) = nil, want a 2-element array", row)
		}
		wantX := int16(row*10 + 1)
		wantY := int16(row*10 + 2)
		arr.Set(0, wantX)
		arr.Set(1, wantY)

		got := ps.Position(row)
		if got.Get(0) != wantX {
			t.Errorf("row %d x = %d, want %d", row, got.Get(0), wantX)
		}
		if got.Get(1) != wantY {
			t.Errorf("row %d y = %d, want %d", row, got.Get(1), wantY)
		}
	}
}
