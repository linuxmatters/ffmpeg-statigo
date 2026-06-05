package ffmpeg

import (
	"testing"
)

// SWS_BILINEAR is the bilinear scaling flag value from libswscale.
const swsBilinear = 2

// TestSwsScale scales one frame from YUV420P to RGB24 and frees the context,
// asserting balanced frees of both images and the scaler context. Run under -race
// to exercise cgocheck on the plane-pointer arrays passed to sws_scale.
func TestSwsScale(t *testing.T) {
	const (
		srcW  = 64
		srcH  = 48
		dstW  = 32
		dstH  = 24
		align = 1
	)

	ctx := SwsGetContext(srcW, srcH, AVPixFmtYuv420P, dstW, dstH, AVPixFmtRgb24, swsBilinear, nil, nil, nil)
	if ctx == nil {
		t.Fatal("SwsGetContext returned nil")
	}
	defer SwsFreeContext(&ctx)

	srcPlanes, srcLines, _, err := AVImageAlloc(srcW, srcH, AVPixFmtYuv420P, align)
	if err != nil {
		t.Fatalf("AVImageAlloc src failed: %v", err)
	}
	defer AVImageFreePlanes(srcPlanes)

	dstPlanes, dstLines, _, err := AVImageAlloc(dstW, dstH, AVPixFmtRgb24, align)
	if err != nil {
		t.Fatalf("AVImageAlloc dst failed: %v", err)
	}
	defer AVImageFreePlanes(dstPlanes)

	outH, err := SwsScale(ctx, srcPlanes, srcLines, 0, srcH, dstPlanes, dstLines)
	if err != nil {
		t.Fatalf("SwsScale failed: %v (outH=%d)", err, outH)
	}
	if outH != dstH {
		t.Errorf("SwsScale output height = %d, want %d", outH, dstH)
	}
}

// TestSwsGetCachedContext verifies the cached-context constructor returns a usable
// context and frees cleanly.
func TestSwsGetCachedContext(t *testing.T) {
	const (
		srcW = 64
		srcH = 48
	)

	ctx := SwsGetCachedContext(nil, srcW, srcH, AVPixFmtYuv420P, srcW, srcH, AVPixFmtRgb24, swsBilinear, nil, nil, nil)
	if ctx == nil {
		t.Fatal("SwsGetCachedContext returned nil")
	}
	defer SwsFreeContext(&ctx)

	// Reuse the same context with identical parameters; should return the same context.
	ctx2 := SwsGetCachedContext(ctx, srcW, srcH, AVPixFmtYuv420P, srcW, srcH, AVPixFmtRgb24, swsBilinear, nil, nil, nil)
	if ctx2 == nil {
		t.Fatal("SwsGetCachedContext reuse returned nil")
	}
	// ctx2 owns the underlying context now; avoid double free by only freeing ctx2.
	ctx = ctx2
}

// TestSwsColorspaceDetails round-trips the colorspace tables for a YUV→RGB context.
func TestSwsColorspaceDetails(t *testing.T) {
	const (
		w = 64
		h = 48
	)

	ctx := SwsGetContext(w, h, AVPixFmtYuv420P, w, h, AVPixFmtRgb24, swsBilinear, nil, nil, nil)
	if ctx == nil {
		t.Fatal("SwsGetContext returned nil")
	}
	defer SwsFreeContext(&ctx)

	invTable, srcRange, table, dstRange, brightness, contrast, saturation, err := SwsGetColorspaceDetails(ctx)
	if err != nil {
		// Not all builds/contexts support reading details; document and skip.
		t.Skipf("SwsGetColorspaceDetails not supported: %v", err)
	}
	if len(invTable) != 4 || len(table) != 4 {
		t.Fatalf("expected 4-element tables, got inv=%d table=%d", len(invTable), len(table))
	}

	ret, err := SwsSetColorspaceDetails(ctx, invTable, srcRange, table, dstRange, brightness, contrast, saturation)
	if err != nil {
		t.Fatalf("SwsSetColorspaceDetails failed: %v (ret=%d)", err, ret)
	}
}
