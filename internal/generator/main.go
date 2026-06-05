package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"maps"
	"os"
	"slices"

	"github.com/Newbluecake/bootstrap/clang"
)

// skipCeiling caps the total skip-marker count a clean generator run is
// allowed to record. It is an upper bound (`count > skipCeiling` trips), not
// an exact equality, so it acts as a regression ceiling rather than a brittle
// target: it tolerates a handful of legitimate per-symbol skips a future header
// bump might introduce while still catching wholesale degradation, such as an
// allowlist regression that drops dozens of bindings.
//
// The current count is the residue of shapes the generator cannot safely emit:
// int32 matrix-pointer functions, the AVExifEntry union field, and the
// fixed-size struct-array fields in the DRM descriptors
// (AVDRMFrameDescriptor.objects/layers, AVDRMLayerDescriptor.planes). Most are
// covered by hand-written bindings; the skip marker records that the generator
// declined them, not that the symbol is unavailable.
//
// Bumping this constant is a curation decision. An FFmpeg upgrade that binds new
// headers carrying unemittable shapes needs an intentional bump alongside the
// header update; do not raise it to silence a regression.
const skipCeiling = 245

func main() {
	skips, err := run(os.Args[1:])
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return // usage already written to stderr by the flag package
		}
		log.Panicln(err)
	}

	printSkipSummary(os.Stderr, skips)

	if err := enforceSkipCeiling(skips.Total(), skipCeiling); err != nil {
		log.Panicln(err)
	}
}

// enforceSkipCeiling returns a non-nil error when total exceeds ceiling.
// Extracted so the ceiling policy is testable without invoking the libclang
// parse path: pass any fabricated total and ceiling. The message names both
// values so a tripped ceiling surfaces the actionable numbers in the run log.
func enforceSkipCeiling(total, ceiling int) error {
	if total > ceiling {
		return fmt.Errorf("skip-count regression: %d skipped symbols exceeds ceiling of %d", total, ceiling)
	}
	return nil
}

// run drives one end-to-end generator pass: parse the FFmpeg headers, apply
// the in-tree fixups, and emit the five `*.gen.go` files. It returns the skip
// collector populated during emission so callers (main, the ceiling check,
// tests) can inspect every `skipped due to ...` decision the run made.
//
// It exists as a separate function so the whole run path is testable without
// invoking the package binary.
func run(args []string) (*SkipCollector, error) {
	fs := flag.NewFlagSet("generator", flag.ContinueOnError)
	verbose := fs.Bool("v", false, "verbose trace logging")
	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if !*verbose {
		log.SetOutput(io.Discard)
	}

	log.Println("Bindings generator")
	log.Printf("libclang: %s", clang.GetClangVersion())
	log.Printf("platform args: %v", getPlatformArgs())
	log.Printf("system includes: %v", getSystemIncludes())

	skips := &SkipCollector{}

	m := Parse(skips)

	if err := applyManualFixups(m); err != nil {
		return nil, err
	}

	return Gen(m, skips), nil
}

// applyManualFixups is the single home for per-type corrections that libclang
// resolves incorrectly on the pinned FFmpeg headers. Each mutation patches a
// parsed type whose shape the generator cannot infer from the headers alone:
//   - AVRational must pass by value (libclang reports it as a pointer-style
//     aggregate, but the C API takes and returns it as a value type).
//   - AVOptionType carries a stray doc comment that libclang attaches; clearing
//     it keeps the generated enum free of unwanted leading text.
//
// Apply after Parse and before Gen. Add new corrections here rather than
// mutating m.structs/m.enums inline elsewhere.
//
// It returns an error when a fixup target is missing from the parsed module:
// the lookups are this layer's single load-bearing fragility, so a header
// reshape that removes or renames AVRational or AVOptionType trips loudly here
// with the offending symbol named, rather than nil-panicking on the assignment
// or emitting a quietly wrong binding.
func applyManualFixups(m *Module) error {
	if _, ok := m.structs["AVRational"]; !ok {
		return fmt.Errorf("manual fixup target %q absent from parsed structs: FFmpeg headers changed shape, update applyManualFixups", "AVRational")
	}
	m.structs["AVRational"].ByValue = true

	if _, ok := m.enums["AVOptionType"]; !ok {
		return fmt.Errorf("manual fixup target %q absent from parsed enums: FFmpeg headers changed shape, update applyManualFixups", "AVOptionType")
	}
	m.enums["AVOptionType"].Comment = ""

	return nil
}

// printSkipSummary writes the end-of-run skip aggregation: a total count line
// followed by the sorted, deduplicated list of skipped symbol names. Symbol
// text mirrors the `Symbol` field on each SkipEntry so a future regression
// surfaces the offending name without having to grep the generated files.
// Reads stderr-friendly: the summary is the only output the generator must
// produce on a non-verbose run.
func printSkipSummary(w io.Writer, c *SkipCollector) {
	total := c.Total()

	fmt.Fprintf(w, "\nSkip summary: %d markers across %d unique symbols\n", total, uniqueSymbolCount(c))

	if total == 0 {
		return
	}

	manual := manualBindingBySymbol(c)

	for _, sym := range sortedUniqueSymbols(c) {
		if name := manual[sym]; name != "" {
			fmt.Fprintf(w, "  %s (manual binding: %s)\n", sym, name)
			continue
		}
		fmt.Fprintf(w, "  %s\n", sym)
	}
}

// manualBindingBySymbol maps each skipped symbol to the hand-written Go binding
// that covers it, when one was detected. Only entries with a non-empty Manual
// field appear, so the map doubles as a presence test in printSkipSummary.
func manualBindingBySymbol(c *SkipCollector) map[string]string {
	if c == nil {
		return nil
	}

	out := make(map[string]string)
	for _, e := range c.Entries {
		if e.Manual != "" {
			out[e.Symbol] = e.Manual
		}
	}
	return out
}

// sortedUniqueSymbols returns the unique skipped symbol names in lexical
// order, feeding both the run summary and the ceiling diagnostic when the cap
// trips.
func sortedUniqueSymbols(c *SkipCollector) []string {
	if c == nil {
		return nil
	}

	seen := make(map[string]bool, len(c.Entries))

	for _, e := range c.Entries {
		seen[e.Symbol] = true
	}

	return slices.Sorted(maps.Keys(seen))
}

// uniqueSymbolCount reports the number of distinct symbol names skipped,
// which is smaller than the total marker count when one symbol (a struct)
// has several fields skipped or one function trips two skip predicates.
func uniqueSymbolCount(c *SkipCollector) int {
	return len(sortedUniqueSymbols(c))
}
