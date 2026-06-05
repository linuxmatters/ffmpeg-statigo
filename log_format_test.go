package ffmpeg

import (
	"strings"
	"sync"
	"testing"
)

// payloadWithDirectives carries printf conversion directives that must survive
// the "%s" indirection untouched. If any reached C's format parser, %n would
// crash or corrupt memory and %s/%x would consume a (non-existent) argument.
const payloadWithDirectives = "value=100% done %n %s %x literal"

// TestAVLogPassesFormatDirectivesVerbatim registers a log callback and asserts
// AVLog delivers caller text intact, including embedded printf directives. This
// proves the "%s" indirection blocks format-string injection: the directives
// arrive as plain bytes rather than being interpreted by C's printf parser.
func TestAVLogPassesFormatDirectivesVerbatim(t *testing.T) {
	var (
		mu      sync.Mutex
		gotMsg  string
		gotCall bool
	)

	prevLevel, _ := AVLogGetLevel()
	AVLogSetLevel(AVLogTrace)
	defer AVLogSetLevel(prevLevel)

	AVLogSetCallback(func(_ *LogCtx, _ int, msg string) {
		mu.Lock()
		defer mu.Unlock()
		gotMsg += msg
		gotCall = true
	})
	defer AVLogSetCallback(nil)

	AVLog(nil, AVLogInfo, payloadWithDirectives+"\n")

	mu.Lock()
	defer mu.Unlock()

	if !gotCall {
		t.Fatal("log callback never fired")
	}

	if !strings.Contains(gotMsg, payloadWithDirectives) {
		t.Fatalf("log payload mangled: want substring %q, got %q", payloadWithDirectives, gotMsg)
	}
}

// TestAVAsprintfPayloadVerbatim asserts AVAsprintf returns the caller text
// unchanged, including printf directives, and that the result is freeable.
func TestAVAsprintfPayloadVerbatim(t *testing.T) {
	got := AVAsprintf(payloadWithDirectives)
	if got == nil {
		t.Fatal("AVAsprintf returned nil")
	}
	defer got.Free()

	if s := got.String(); s != payloadWithDirectives {
		t.Fatalf("AVAsprintf mangled payload: want %q, got %q", payloadWithDirectives, s)
	}
}

// TestAVStrlcatfAppendsVerbatim asserts AVStrlcatf appends caller text intact
// and reports the untruncated length.
func TestAVStrlcatfAppendsVerbatim(t *testing.T) {
	const prefix = "prefix:"
	const bufSize = 256

	dst := AllocCStr(bufSize)
	defer dst.Free()

	// Seed the buffer with a known prefix via a first append onto the empty
	// NUL-initialised buffer.
	first := AVStrlcatf(dst, bufSize, prefix)
	if int(first) != len(prefix) {
		t.Fatalf("first AVStrlcatf length: want %d, got %d", len(prefix), first)
	}

	total := AVStrlcatf(dst, bufSize, payloadWithDirectives)

	want := prefix + payloadWithDirectives
	if int(total) != len(want) {
		t.Fatalf("AVStrlcatf length: want %d, got %d", len(want), total)
	}

	if s := dst.String(); s != want {
		t.Fatalf("AVStrlcatf mangled buffer: want %q, got %q", want, s)
	}
}

// TestAVStrlcatfReportsTruncation asserts AVStrlcatf returns the would-be length
// when the buffer is too small, mirroring the C contract.
func TestAVStrlcatfReportsTruncation(t *testing.T) {
	const bufSize = 8
	const msg = "0123456789abcdef"

	dst := AllocCStr(bufSize)
	defer dst.Free()

	got := AVStrlcatf(dst, bufSize, msg)
	if int(got) != len(msg) {
		t.Fatalf("truncation length: want %d, got %d", len(msg), got)
	}

	// Stored content is truncated to bufSize-1 bytes plus a NUL terminator.
	if s := dst.String(); len(s) != bufSize-1 {
		t.Fatalf("truncated stored length: want %d, got %d (%q)", bufSize-1, len(s), s)
	}
}

// TestAVBprintfPayloadVerbatim asserts AVBprintf appends caller text intact to a
// growable buffer, including printf directives.
func TestAVBprintfPayloadVerbatim(t *testing.T) {
	buf := newTestAVBPrint()
	defer freeTestAVBPrint(buf)

	AVBprintInit(buf, 0, 1)

	AVBprintf(buf, payloadWithDirectives)

	if s := buf.Str().String(); s != payloadWithDirectives {
		t.Fatalf("AVBprintf mangled payload: want %q, got %q", payloadWithDirectives, s)
	}
}
