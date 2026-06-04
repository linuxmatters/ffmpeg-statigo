package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"

	"github.com/Newbluecake/bootstrap/clang"
)

// skipCeiling caps the total skip-marker count a clean generator run is
// allowed to record. It is an upper bound (`count > skipCeiling` trips), not
// an exact equality, so the Phase 1 baseline (148 functions + 85 structs =
// 233) acts as a regression ceiling rather than a brittle target.
//
// Rationale: the unemittable-shape panics Task 3.2 converted to skip-with-reason
// add zero markers on the pinned FFmpeg headers (none of those code paths
// fire today), so the live count stays at the 233 baseline. An upper bound
// tolerates a handful of legitimate per-symbol skips a future header bump
// might introduce while still catching wholesale degradation (e.g. an
// allowlist regression that drops dozens of bindings).
//
// Bumping this constant is a curation decision. A legitimate FFmpeg upgrade
// that introduces new unemittable symbols requires an intentional bump
// alongside the header update; do not raise it to silence a regression.
const skipCeiling = 233

func main() {
	skips, err := run(os.Args[1:], os.Stderr)
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
		return fmt.Errorf("skip-count regression: %d skipped symbols exceeds ceiling of %d (see IMPROVE-PLAN.md Task 3.4)", total, ceiling)
	}
	return nil
}

// run drives one end-to-end generator pass: parse the FFmpeg headers, apply
// the in-tree fixups, and emit the five `*.gen.go` files. It returns the skip
// collector populated during emission so callers (main, tests, the Task 3.4
// ceiling check) can inspect every `skipped due to ...` decision the run made.
//
// The function is the single extraction point Task 3.1 added so the run path
// is testable without invoking the package binary. summaryOut is reserved for
// callers that want non-stderr toolchain logging; today only the verbose flag
// touches it via log.SetOutput.
func run(args []string, summaryOut io.Writer) (*SkipCollector, error) {
	_ = summaryOut

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

	m.structs["AVRational"].ByValue = true
	m.enums["AVOptionType"].Comment = ""

	return Gen(m, skips), nil
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

	for _, sym := range sortedUniqueSymbols(c) {
		fmt.Fprintf(w, "  %s\n", sym)
	}
}

// sortedUniqueSymbols returns the unique skipped symbol names in lexical
// order. Used by the run summary and exposed for Task 3.4 to feed the same
// list into the ceiling diagnostic when the cap trips.
func sortedUniqueSymbols(c *SkipCollector) []string {
	if c == nil {
		return nil
	}

	seen := make(map[string]bool, len(c.Entries))
	out := make([]string, 0, len(c.Entries))

	for _, e := range c.Entries {
		if seen[e.Symbol] {
			continue
		}
		seen[e.Symbol] = true
		out = append(out, e.Symbol)
	}

	sort.Strings(out)
	return out
}

// uniqueSymbolCount reports the number of distinct symbol names skipped,
// which is smaller than the total marker count when one symbol (a struct)
// has several fields skipped or one function trips two skip predicates.
func uniqueSymbolCount(c *SkipCollector) int {
	return len(sortedUniqueSymbols(c))
}
