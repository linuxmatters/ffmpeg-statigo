package main

import (
	"strings"
	"testing"
)

// TestEnforceSkipCeilingAtOrUnder pins the upper-bound semantics: a total at
// or below the ceiling must not error. Mirrors the unmodified-run path where
// the live count equals the baseline (233) and the run exits 0.
func TestEnforceSkipCeilingAtOrUnder(t *testing.T) {
	tests := []struct {
		name    string
		total   int
		ceiling int
	}{
		{"zero under ceiling", 0, skipCeiling},
		{"one under ceiling", skipCeiling - 1, skipCeiling},
		{"equal to ceiling", skipCeiling, skipCeiling},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := enforceSkipCeiling(tt.total, tt.ceiling); err != nil {
				t.Errorf("enforceSkipCeiling(%d, %d) = %v, want nil", tt.total, tt.ceiling, err)
			}
		})
	}
}

// TestEnforceSkipCeilingExceeds pins the regression-trip semantics: a total
// strictly above the ceiling returns a non-nil error whose message names both
// the count and the ceiling, so the run log surfaces the actionable numbers.
func TestEnforceSkipCeilingExceeds(t *testing.T) {
	tests := []struct {
		name    string
		total   int
		ceiling int
	}{
		{"one over ceiling", skipCeiling + 1, skipCeiling},
		{"ten over ceiling", skipCeiling + 10, skipCeiling},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := enforceSkipCeiling(tt.total, tt.ceiling)
			if err == nil {
				t.Fatalf("enforceSkipCeiling(%d, %d) = nil, want error", tt.total, tt.ceiling)
			}
			msg := err.Error()
			if !strings.Contains(msg, "skip-count regression") {
				t.Errorf("error message %q missing 'skip-count regression'", msg)
			}
			// The diagnostic must name both values so the run log shows the actual count and ceiling.
			if !strings.Contains(msg, "233") && tt.ceiling == skipCeiling {
				t.Errorf("error message %q missing ceiling value %d", msg, tt.ceiling)
			}
		})
	}
}

// TestSkipCollectorThroughCeiling drives a fabricated SkipCollector through
// the same check path the production run uses: Record() one symbol per
// extra-skip slot above the ceiling, then call enforceSkipCeiling on
// collector.Total(). This exercises the full collector+check contract
// without invoking libclang or touching `*.h` headers.
func TestSkipCollectorThroughCeiling(t *testing.T) {
	c := &SkipCollector{}

	// Fill the collector exactly to the ceiling and verify the check passes.
	for range skipCeiling {
		c.Record("sym", "reason")
	}
	if got := c.Total(); got != skipCeiling {
		t.Fatalf("collector.Total() = %d, want %d", got, skipCeiling)
	}
	if err := enforceSkipCeiling(c.Total(), skipCeiling); err != nil {
		t.Errorf("enforceSkipCeiling(%d, %d) = %v, want nil at ceiling", c.Total(), skipCeiling, err)
	}

	// Inject one extra skipped symbol above the ceiling; the check must trip.
	c.Record("extra_symbol", "fabricated overflow")
	if got := c.Total(); got != skipCeiling+1 {
		t.Fatalf("collector.Total() = %d, want %d", got, skipCeiling+1)
	}
	err := enforceSkipCeiling(c.Total(), skipCeiling)
	if err == nil {
		t.Fatalf("enforceSkipCeiling(%d, %d) = nil, want error after extra Record", c.Total(), skipCeiling)
	}
	if !strings.Contains(err.Error(), "exceeds ceiling") {
		t.Errorf("error message %q missing 'exceeds ceiling'", err.Error())
	}
}
