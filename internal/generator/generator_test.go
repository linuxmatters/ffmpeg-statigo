package main

import (
	"testing"

	"github.com/dave/jennifer/jen"
)

// TestConvCamelStructCollision locks the rule that a camel-cased name colliding
// with an already-registered struct gets a "_" suffix to stay unique in Go's
// single package-level identifier namespace. Both real collisions in the API
// (avfilter_link, av_buffer_ref) and the non-colliding default are covered.
func TestConvCamelStructCollision(t *testing.T) {
	g := &Generator{
		input: &Module{
			structs: map[string]*Struct{
				"AVFilterLink": {Name: "AVFilterLink"},
				"AVBufferRef":  {Name: "AVBufferRef"},
			},
		},
	}

	tests := []struct {
		name string
		in   string
		want string
	}{
		{"avfilter_link collision", "avfilter_link", "AVFilterLink_"},
		{"av_buffer_ref collision", "av_buffer_ref", "AVBufferRef_"},
		{"no collision", "av_some_function", "AVSomeFunction"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := g.convCamel(tt.in); got != tt.want {
				t.Errorf("convCamel(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// skipGen returns a Generator whose module registers one by-value struct
// (Frame), one by-pointer struct (Buf) and one callback so the skip predicates
// below exercise every registry-dependent branch.
func skipGen() *Generator {
	return &Generator{
		input: &Module{
			structs: map[string]*Struct{
				"AVFrame": {Name: "AVFrame", ByValue: true},
				"AVBuf":   {Name: "AVBuf", ByValue: false},
			},
			callbacks: map[string]*Function{
				"av_log_cb": {Name: "av_log_cb"},
			},
			enums: map[string]*Enum{},
		},
	}
}

func ptr(inner Type) *PointerType  { return &PointerType{Inner: inner} }
func ident(name string) *IdentType { return &IdentType{Name: name} }

// TestMarshalArgSkip locks the per-argument skip decisions extracted into
// marshalArg. These predicates decide whether a function is emittable at all,
// so an inverted condition silently drops (or wrongly emits) bindings. The
// byte-identical regeneration gate would catch that too, but only against the
// live FFmpeg headers; this pins the logic to named type shapes so a regression
// names the exact branch that broke. Emitted code is deliberately not asserted.
func TestMarshalArgSkip(t *testing.T) {
	g := skipGen()

	tests := []struct {
		name     string
		argName  string
		argType  Type
		wantSkip bool
	}{
		{"primitive by value kept", "count", ident("int"), false},
		{"by-value struct pointer skipped", "frame", ptr(ident("AVFrame")), true},
		{"by-pointer struct pointer kept", "buf", ptr(ident("AVBuf")), false},
		{"callback by value skipped", "cb", ident("av_log_cb"), true},
		{"char pointer kept", "name", ptr(ident("char")), false},
		{"non-allowlisted output-shaped primitive pointer skipped", "out_size", ptr(ident("int")), true},
		{"non-allowlisted primitive pointer skipped", "level", ptr(ident("int")), true},
		{"double pointer to struct kept", "frame", ptr(ptr(ident("AVBuf"))), false},
		{"double pointer to char skipped", "argv", ptr(ptr(ident("char"))), true},
		{"double pointer to primitive skipped", "vals", ptr(ptr(ident("int"))), true},
		{"array param skipped", "vals", &Array{Inner: ident("int")}, true},
		{"const array param skipped", "vals", &ConstArray{Inner: ident("int"), Size: 4}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := newFile()
			fn := &Function{Name: "fn"}
			arg := &Param{Name: tt.argName, Type: tt.argType}
			_, _, _, _, skip := g.marshalArg(o, fn, arg)
			if skip != tt.wantSkip {
				t.Errorf("marshalArg(%s %v) skip = %v, want %v", tt.argName, tt.argType, skip, tt.wantSkip)
			}
		})
	}
}

// TestMarshalReturnSkip locks the return-value skip decisions extracted into
// marshalReturn. A by-pointer (non-ByValue) struct returned by value and a
// pointer-to-pointer return are the only two cases that make a function
// unemittable; everything else must stay emittable. Emitted code is not
// asserted, only the skip boolean.
func TestMarshalReturnSkip(t *testing.T) {
	g := skipGen()

	tests := []struct {
		name     string
		result   Type
		wantSkip bool
	}{
		{"void return kept", nil, false},
		{"int return kept", ident("int"), false},
		{"primitive return kept", ident("int64_t"), false},
		{"by-value struct return kept", ident("AVFrame"), false},
		{"by-pointer struct return skipped", ident("AVBuf"), true},
		{"callback return kept", ident("av_log_cb"), false},
		{"char pointer return kept", ptr(ident("char")), false},
		{"struct pointer return kept", ptr(ident("AVBuf")), false},
		{"pointer to pointer return skipped", ptr(ptr(ident("AVBuf"))), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := newFile()
			fn := &Function{Name: "fn"}
			_, _, skip := g.marshalReturn(o, fn, tt.result, jen.Id("C").Dot("fn").Call(), nil, nil)
			if skip != tt.wantSkip {
				t.Errorf("marshalReturn(%v) skip = %v, want %v", tt.result, skip, tt.wantSkip)
			}
		})
	}
}

// TestMarshalArgOutputPointerAllowlist pins the (function, parameter) allowlist
// decision in marshalPointerArg. After Task 2.2 the allowlist is the only
// output-pointer routing signal: a substring match like `width`, `size`, or
// `_ptr` no longer keeps a parameter emittable on its own. This test asserts
// that representative allowlisted pairs spanning the original heuristic
// patterns still resolve to an emitted output pointer, and that near-miss
// pairs whose name would have matched the old substring sweep now skip.
//
// The sample is deliberately small (the byte-identical regeneration gate
// pins the full 64-pair table); the goal here is to name the exact branch
// that breaks when the lookup is inverted, dropped, or rekeyed.
func TestMarshalArgOutputPointerAllowlist(t *testing.T) {
	g := skipGen()

	tests := []struct {
		name     string
		fnName   string
		argName  string
		argType  Type
		wantSkip bool
	}{
		// Allowlisted pairs spanning the original substring patterns.
		{"out_-prefixed allowlisted", "av_detection_bbox_alloc", "out_size", ptr(ident("int")), false},
		{"size allowlisted (size_t fixup site)", "av_cpb_properties_alloc", "size", ptr(ident("int")), false},
		{"w-suffixed allowlisted", "av_opt_get_image_size", "w_out", ptr(ident("int")), false},
		{"h+ptr allowlisted", "av_parse_video_size", "height_ptr", ptr(ident("int")), false},
		{"ptr allowlisted", "av_url_split", "port_ptr", ptr(ident("int")), false},

		// Near-miss pairs: name matches the retired substring sweep
		// (size/w/h) but the (function, parameter) is not in the allowlist.
		{"size param, function not in allowlist", "av_unknown_future_fn", "size", ptr(ident("int")), true},
		{"width param, function in allowlist but param not", "av_cpb_properties_alloc", "width", ptr(ident("int")), true},
		{"h param, function not in allowlist", "av_unknown_future_fn", "height", ptr(ident("int")), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := newFile()
			fn := &Function{Name: tt.fnName}
			arg := &Param{Name: tt.argName, Type: tt.argType}
			_, _, _, _, skip := g.marshalArg(o, fn, arg)
			if skip != tt.wantSkip {
				t.Errorf("marshalArg(%s.%s %v) skip = %v, want %v",
					tt.fnName, tt.argName, tt.argType, skip, tt.wantSkip)
			}
		})
	}
}

// render returns the Go source for a single jen.Code fragment so tests can
// string-match the emitted parameter type and C cast.
func render(c jen.Code) string {
	return (&jen.Statement{}).Add(c).GoString()
}

// TestMarshalArgSizeTOutputParams pins the size_t rewrite to the explicit
// sizeTOutputParams lookup introduced in Task 1.1. marshalArg rewrites a
// pointer-to-int output parameter to size_t only when the (function, parameter)
// pair is present in sizeTOutputParams; the former strings.Contains(fn.Name,
// "_alloc") substring heuristic is gone. The positive case proves the lookup
// drives the rewrite (param becomes *uint64, C cast uses C.size_t). The
// regression case proves an _alloc function name alone no longer triggers the
// rewrite when the pair is absent from the table (param stays *int, cast uses
// C.int), so a reintroduced substring heuristic fails here.
func TestMarshalArgSizeTOutputParams(t *testing.T) {
	g := skipGen()

	// Positive: av_dovi_alloc.size is in both outputPointerAllowlist and
	// sizeTOutputParams, so the int* output pointer is rewritten to size_t.
	t.Run("allowlisted size_t pair rewritten", func(t *testing.T) {
		o := newFile()
		fn := &Function{Name: "av_dovi_alloc"}
		arg := &Param{Name: "size", Type: ptr(ident("int"))}
		params, args, _, _, skip := g.marshalArg(o, fn, arg)
		if skip {
			t.Fatalf("av_dovi_alloc.size unexpectedly skipped")
		}
		if got := render(params[0]); got != "size * uint64" {
			t.Errorf("param = %q, want %q", got, "size * uint64")
		}
		if got := render(args[0]); got != "(*C.size_t)(unsafe.Pointer(size))" {
			t.Errorf("arg = %q, want %q", got, "(*C.size_t)(unsafe.Pointer(size))")
		}
	})

	// Regression: an _alloc-named function that is allowlisted (so its output
	// pointer is emitted, not skipped) but absent from sizeTOutputParams must
	// stay *int. outputPointerAllowlist is a package var, so inject and restore
	// the fake entry with a defer to keep other tests unaffected. The pair is
	// deliberately never added to sizeTOutputParams.
	t.Run("allocnamed pair absent from size_t table stays int", func(t *testing.T) {
		const fnName = "av_fake_alloc"
		saved := outputPointerAllowlist[fnName]
		outputPointerAllowlist[fnName] = map[string]bool{"size": true}
		defer func() {
			if saved == nil {
				delete(outputPointerAllowlist, fnName)
			} else {
				outputPointerAllowlist[fnName] = saved
			}
		}()

		if sizeTOutputParams[fnName] != nil {
			t.Fatalf("test invariant broken: %s present in sizeTOutputParams", fnName)
		}

		o := newFile()
		fn := &Function{Name: fnName}
		arg := &Param{Name: "size", Type: ptr(ident("int"))}
		params, args, _, _, skip := g.marshalArg(o, fn, arg)
		if skip {
			t.Fatalf("%s.size unexpectedly skipped", fnName)
		}
		if got := render(params[0]); got != "size * int" {
			t.Errorf("param = %q, want %q (substring heuristic must be gone)", got, "size * int")
		}
		if got := render(args[0]); got != "(*C.int)(unsafe.Pointer(size))" {
			t.Errorf("arg = %q, want %q (substring heuristic must be gone)", got, "(*C.int)(unsafe.Pointer(size))")
		}
	})
}

// marshalField is intentionally left to the byte-identical regeneration gate.
// It returns nothing (void) and signals its three early-exit decisions
// (skippedFields, bitfield, ident-callback) only by suppressing emitted output.
// Asserting them would require rendering and inspecting the emitted bytes, which
// duplicates the gate with brittle string matching and no added bug-catching
// power, so no independent test is added here.
