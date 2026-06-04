package main

import (
	"testing"
)

// buildLibraryOrder reads the package-global allLibraryDefinitions rather than
// an injected graph, and it always holds the package-global ffmpeg var aside to
// append last. To exercise custom graphs, these tests swap allLibraryDefinitions
// (restored via defer) and include the real ffmpeg var as the terminal node so
// the hold-aside path is honoured.

// indexOf returns the position of lib in order, or -1 if absent.
func indexOf(order []*Library, lib *Library) int {
	for i, l := range order {
		if l == lib {
			return i
		}
	}
	return -1
}

// withDefinitions swaps allLibraryDefinitions for the duration of fn.
func withDefinitions(defs []*Library, fn func()) {
	saved := allLibraryDefinitions
	allLibraryDefinitions = defs
	defer func() { allLibraryDefinitions = saved }()
	fn()
}

// TestBuildLibraryOrder covers dependency ordering, the ffmpeg-last invariant,
// the real definitions, and cycle detection.
func TestBuildLibraryOrder(t *testing.T) {
	t.Run("dependencies precede dependents", func(t *testing.T) {
		base := &Library{Name: "base"}
		mid := &Library{Name: "mid", Dependencies: []*Library{base}}
		top := &Library{Name: "top", Dependencies: []*Library{mid, base}}

		withDefinitions([]*Library{top, mid, base, ffmpeg}, func() {
			order, err := buildLibraryOrder()
			if err != nil {
				t.Fatalf("buildLibraryOrder() error = %v", err)
			}

			for _, pair := range [][2]*Library{{base, mid}, {mid, top}, {base, top}} {
				dep, dependent := pair[0], pair[1]
				di, ddi := indexOf(order, dep), indexOf(order, dependent)
				if di == -1 || ddi == -1 {
					t.Fatalf("missing %s (%d) or %s (%d) in order", dep.Name, di, dependent.Name, ddi)
				}
				if di >= ddi {
					t.Errorf("%s (index %d) should precede %s (index %d)", dep.Name, di, dependent.Name, ddi)
				}
			}
		})
	})

	t.Run("ffmpeg is last", func(t *testing.T) {
		a := &Library{Name: "a"}
		b := &Library{Name: "b", Dependencies: []*Library{a}}

		withDefinitions([]*Library{ffmpeg, b, a}, func() {
			order, err := buildLibraryOrder()
			if err != nil {
				t.Fatalf("buildLibraryOrder() error = %v", err)
			}
			if len(order) == 0 || order[len(order)-1] != ffmpeg {
				t.Errorf("last element = %v, want ffmpeg", order[len(order)-1])
			}
		})
	})

	t.Run("real definitions sort cleanly with ffmpeg last", func(t *testing.T) {
		order, err := buildLibraryOrder()
		if err != nil {
			t.Fatalf("buildLibraryOrder() over real definitions error = %v", err)
		}
		if len(order) != len(allLibraryDefinitions) {
			t.Errorf("ordered %d libraries, want %d", len(order), len(allLibraryDefinitions))
		}
		if order[len(order)-1] != ffmpeg {
			t.Errorf("last element = %v, want ffmpeg", order[len(order)-1])
		}
		// Every dependency precedes its dependents.
		for _, lib := range order {
			for _, dep := range lib.Dependencies {
				if indexOf(order, dep) >= indexOf(order, lib) {
					t.Errorf("dependency %s does not precede %s", dep.Name, lib.Name)
				}
			}
		}
	})

	t.Run("cyclic graph returns error", func(t *testing.T) {
		x := &Library{Name: "x"}
		y := &Library{Name: "y"}
		x.Dependencies = []*Library{y}
		y.Dependencies = []*Library{x}

		withDefinitions([]*Library{x, y, ffmpeg}, func() {
			order, err := buildLibraryOrder()
			if err == nil {
				t.Fatalf("buildLibraryOrder() over cyclic graph returned nil error, order = %v", order)
			}
		})
	})
}
