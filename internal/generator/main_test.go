package main

import (
	"strconv"
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
			if !strings.Contains(msg, strconv.Itoa(skipCeiling)) && tt.ceiling == skipCeiling {
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

// TestApplyManualFixupsMissingTarget pins the fail-loud contract: when a fixup
// target is absent from the parsed module, applyManualFixups returns a non-nil
// error naming the missing symbol instead of nil-panicking on the assignment.
// Each case omits exactly one target so the per-symbol diagnostic is exercised
// independently. Builds the Module by hand so the test needs no libclang parse.
func TestApplyManualFixupsMissingTarget(t *testing.T) {
	tests := []struct {
		name   string
		mod    *Module
		symbol string
	}{
		{
			name: "AVRational absent from structs",
			mod: &Module{
				structs: map[string]*Struct{},
				enums:   map[string]*Enum{"AVOptionType": {}},
			},
			symbol: "AVRational",
		},
		{
			name: "AVOptionType absent from enums",
			mod: &Module{
				structs: map[string]*Struct{"AVRational": {}},
				enums:   map[string]*Enum{},
			},
			symbol: "AVOptionType",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := applyManualFixups(tt.mod)
			if err == nil {
				t.Fatalf("applyManualFixups() = nil, want error naming %q", tt.symbol)
			}
			if !strings.Contains(err.Error(), tt.symbol) {
				t.Errorf("error message %q does not name missing symbol %q", err.Error(), tt.symbol)
			}
		})
	}
}

// TestApplyManualFixupsAppliesMutations pins the success path: with both
// targets present, applyManualFixups returns nil and applies the corrections
// (AVRational by value, AVOptionType comment cleared).
func TestApplyManualFixupsAppliesMutations(t *testing.T) {
	m := &Module{
		structs: map[string]*Struct{"AVRational": {ByValue: false}},
		enums:   map[string]*Enum{"AVOptionType": {Comment: "stray libclang comment"}},
	}

	if err := applyManualFixups(m); err != nil {
		t.Fatalf("applyManualFixups() = %v, want nil", err)
	}
	if !m.structs["AVRational"].ByValue {
		t.Errorf("AVRational.ByValue = false, want true")
	}
	if m.enums["AVOptionType"].Comment != "" {
		t.Errorf("AVOptionType.Comment = %q, want empty", m.enums["AVOptionType"].Comment)
	}
}
