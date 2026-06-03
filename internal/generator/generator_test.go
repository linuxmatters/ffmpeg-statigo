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
		{"output primitive pointer kept", "out_size", ptr(ident("int")), false},
		{"non-output primitive pointer skipped", "level", ptr(ident("int")), true},
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

// marshalField is intentionally left to the byte-identical regeneration gate.
// It returns nothing (void) and signals its three early-exit decisions
// (skippedFields, bitfield, ident-callback) only by suppressing emitted output.
// Asserting them would require rendering and inspecting the emitted bytes, which
// duplicates the gate with brittle string matching and no added bug-catching
// power, so no independent test is added here.
