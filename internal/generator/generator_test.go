package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
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
		{"primitive const array param kept", "vals", &ConstArray{Inner: ident("int"), Size: 4}, false},
		{"struct const array param skipped", "bufs", &ConstArray{Inner: ident("AVBuf"), Size: 4}, true},
		{"void pointer const array param skipped", "ptrs", &ConstArray{Inner: ptr(nil), Size: 2}, true},
		{"nested 2D const array param skipped", "matrix", &ConstArray{Inner: &ConstArray{Inner: ident("int"), Size: 4}, Size: 2}, true},
		{"non-positive const array param skipped", "empty", &ConstArray{Inner: ident("int"), Size: 0}, true},
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

func TestMarshalReturnCStrAllocator(t *testing.T) {
	g := skipGen()
	for _, tt := range []struct {
		name string
		want string
	}{
		{"av_strdup", "wrapAVCStr"},
		{"av_strndup", "wrapAVCStr"},
		{"av_version_info", "wrapCStr"},
		{"av_stristr", "wrapCStr"},
		{"av_strdup_other", "wrapCStr"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			fn := &Function{Name: tt.name}
			_, body, skip := g.marshalReturn(newFile(), fn, ptr(ident("char")), jen.Qual("C", tt.name).Call(), nil, nil)
			if skip {
				t.Fatal("string return unexpectedly skipped")
			}
			if got, want := render(body[len(body)-1]), "return "+tt.want+"(ret)"; got != want {
				t.Errorf("return = %q, want %q", got, want)
			}
		})
	}
}

// TestMarshalArgOutputPointerAllowlist pins the (function, parameter) allowlist
// decision in marshalPointerArg. The allowlist is the only output-pointer
// routing signal: a substring match like `width`, `size`, or
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
		{"size allowlisted", "av_cpb_properties_alloc", "size", ptr(ident("size_t")), false},
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

// TestMarshalArgOutputPointers pins the output-pointer table entries:
// av_parse_time.timeval and av_get_output_timestamp.{dts,wall} are pure-output
// int64_t* parameters present in outputParams, so marshalArg must emit them
// (skip == false) rather than record a "non-output primitive pointer" skip.
// Dropping a table entry flips skip to true here and trips skipCeiling in the
// generator run. av_expr_parse_and_eval.res is deliberately excluded: res is a
// genuine output pointer, but the function stays skipped on its later
// const char *const * params, so it never emits a callable binding.
func TestMarshalArgOutputPointers(t *testing.T) {
	g := skipGen()

	tests := []struct {
		name    string
		fnName  string
		argName string
		argType Type
	}{
		{"av_parse_time.timeval emits", "av_parse_time", "timeval", ptr(ident("int64_t"))},
		{"av_get_output_timestamp.dts emits", "av_get_output_timestamp", "dts", ptr(ident("int64_t"))},
		{"av_get_output_timestamp.wall emits", "av_get_output_timestamp", "wall", ptr(ident("int64_t"))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !outputParams[tt.fnName][tt.argName] {
				t.Fatalf("outputParams[%q][%q] absent: output-pointer table entry missing", tt.fnName, tt.argName)
			}

			o := newFile()
			fn := &Function{Name: tt.fnName}
			arg := &Param{Name: tt.argName, Type: tt.argType}
			_, _, _, _, skip := g.marshalArg(o, fn, arg)
			if skip {
				t.Errorf("marshalArg(%s.%s) skip = true, want false (output pointer must emit)",
					tt.fnName, tt.argName)
			}
		})
	}
}

// render returns the Go source for a single jen.Code fragment so tests can
// string-match the emitted parameter type and C cast.
func render(c jen.Code) string {
	return (&jen.Statement{}).Add(c).GoString()
}

// TestMarshalArgPreservesParsedSizeT verifies that marshalArg emits the exact
// primitive type reported by the parser for both output pointers and by-value
// parameters, without name-based int-to-size_t rewrites.
func TestMarshalArgPreservesParsedSizeT(t *testing.T) {
	g := skipGen()

	tests := []struct {
		name      string
		fnName    string
		argName   string
		argType   Type
		wantParam string
		wantArg   string
	}{
		{
			name:      "allowlisted size_t output pointer",
			fnName:    "av_dovi_alloc",
			argName:   "size",
			argType:   ptr(ident("size_t")),
			wantParam: "var size *uint64",
			wantArg:   "ptrsize",
		},
		{
			name:      "allowlisted int output pointer stays int",
			fnName:    "av_dovi_alloc",
			argName:   "size",
			argType:   ptr(ident("int")),
			wantParam: "var size *int",
			wantArg:   "(*C.int)(unsafe.Pointer(size))",
		},
		{
			name:      "by-value size_t",
			fnName:    "av_channel_name",
			argName:   "buf_size",
			argType:   ident("size_t"),
			wantParam: "var bufSize uint64",
			wantArg:   "C.size_t(bufSize)",
		},
		{
			name:      "by-value int on buf_size stays int",
			fnName:    "av_channel_name",
			argName:   "buf_size",
			argType:   ident("int"),
			wantParam: "var bufSize int",
			wantArg:   "C.int(bufSize)",
		},
		{
			name:      "by-value int on max_size stays int",
			fnName:    "avio_read_to_bprint",
			argName:   "max_size",
			argType:   ident("int"),
			wantParam: "var maxSize int",
			wantArg:   "C.int(maxSize)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := newFile()
			fn := &Function{Name: tt.fnName}
			arg := &Param{Name: tt.argName, Type: tt.argType}
			params, args, _, _, skip := g.marshalArg(o, fn, arg)
			if skip {
				t.Fatalf("%s.%s unexpectedly skipped", tt.fnName, tt.argName)
			}
			if got := render(jen.Var().Add(params[0])); got != tt.wantParam {
				t.Errorf("param = %q, want %q", got, tt.wantParam)
			}
			if got := render(args[0]); got != tt.wantArg {
				t.Errorf("arg = %q, want %q", got, tt.wantArg)
			}
		})
	}
}

// TestClassifyConstArrayFieldShape verifies that inline const-array fields
// classify into the expected shape kind, target struct, Go type, dimensions,
// and skip reason.
func TestClassifyConstArrayFieldShape(t *testing.T) {
	g := skipGen()

	tests := []struct {
		name           string
		array          *ConstArray
		wantKind       constArrayFieldShapeKind
		wantStructName string
		wantGoType     string
		wantDims       []int64
		wantReason     string
	}{
		{
			name:           "by-pointer struct inline array",
			array:          &ConstArray{Inner: ident("AVBuf"), Size: 4},
			wantKind:       constArrayFieldShapeByPointerStruct,
			wantStructName: "AVBuf",
		},
		{
			name:           "by-value struct inline array",
			array:          &ConstArray{Inner: ident("AVFrame"), Size: 3},
			wantKind:       constArrayFieldShapeByValueStruct,
			wantStructName: "AVFrame",
		},
		{
			name:           "struct pointer array",
			array:          &ConstArray{Inner: ptr(ident("AVBuf")), Size: 4},
			wantKind:       constArrayFieldShapeStructPointerArray,
			wantStructName: "AVBuf",
		},
		{
			name:       "2D primitive const array",
			array:      &ConstArray{Inner: &ConstArray{Inner: ident("int16_t"), Size: 2}, Size: 3},
			wantKind:   constArrayFieldShapeMultiDimPrimitive,
			wantGoType: "int16",
			wantDims:   []int64{3, 2},
		},
		{
			name:       "2D char const array",
			array:      &ConstArray{Inner: &ConstArray{Inner: ident("char"), Size: 64}, Size: 4},
			wantKind:   constArrayFieldShapeMultiDimPrimitive,
			wantGoType: "uint8",
			wantDims:   []int64{4, 64},
		},
		{
			name:           "2D by-value struct const array",
			array:          &ConstArray{Inner: &ConstArray{Inner: ident("AVFrame"), Size: 2}, Size: 3},
			wantKind:       constArrayFieldShapeMultiDimByValueStruct,
			wantStructName: "AVFrame",
			wantDims:       []int64{3, 2},
		},
		{
			name:       "3D primitive const array",
			array:      &ConstArray{Inner: &ConstArray{Inner: &ConstArray{Inner: ident("int64_t"), Size: 7}, Size: 3}, Size: 8},
			wantKind:   constArrayFieldShapeMultiDimPrimitive,
			wantGoType: "int64",
			wantDims:   []int64{8, 3, 7},
		},
		{
			name:       "4D primitive const array skipped",
			array:      &ConstArray{Inner: &ConstArray{Inner: &ConstArray{Inner: &ConstArray{Inner: ident("int16_t"), Size: 5}, Size: 4}, Size: 3}, Size: 2},
			wantKind:   constArrayFieldShapeMultiDimSkip,
			wantReason: "multi dim const array",
		},
		{
			name:       "2D by-pointer struct const array skipped",
			array:      &ConstArray{Inner: &ConstArray{Inner: ident("AVBuf"), Size: 2}, Size: 3},
			wantKind:   constArrayFieldShapeMultiDimSkip,
			wantReason: "multi dim const array",
		},
		{
			name:       "2D pointer leaf const array skipped",
			array:      &ConstArray{Inner: &ConstArray{Inner: ptr(ident("uint8_t")), Size: 2}, Size: 3},
			wantKind:   constArrayFieldShapeMultiDimSkip,
			wantReason: "multi dim const array",
		},
		{
			name:       "2D non-positive inner extent skipped",
			array:      &ConstArray{Inner: &ConstArray{Inner: ident("int16_t"), Size: 0}, Size: 3},
			wantKind:   constArrayFieldShapeMultiDimSkip,
			wantReason: "multi dim const array",
		},
		{
			name:       "2D non-positive outer extent skipped",
			array:      &ConstArray{Inner: &ConstArray{Inner: ident("int16_t"), Size: 2}, Size: 0},
			wantKind:   constArrayFieldShapeMultiDimSkip,
			wantReason: "multi dim const array",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := g.classifyConstArrayFieldShape(tt.array)
			if got.kind != tt.wantKind {
				t.Errorf("kind = %v, want %v", got.kind, tt.wantKind)
			}
			if tt.wantStructName != "" {
				if got.st == nil {
					t.Fatalf("st = nil, want %q", tt.wantStructName)
				}
				if got.st.Name != tt.wantStructName {
					t.Errorf("st.Name = %q, want %q", got.st.Name, tt.wantStructName)
				}
			} else if got.st != nil {
				t.Errorf("st = %+v, want nil", got.st)
			}
			if got.goType != tt.wantGoType {
				t.Errorf("goType = %q, want %q", got.goType, tt.wantGoType)
			}
			if !slices.Equal(got.dims, tt.wantDims) {
				t.Errorf("dims = %v, want %v", got.dims, tt.wantDims)
			}
			if got.reason != tt.wantReason {
				t.Errorf("reason = %q, want %q", got.reason, tt.wantReason)
			}
		})
	}
}

// TestMarshalFieldInlineStructArray verifies indexed accessor emission,
// Go-keyword C field escaping, and by-value parent skips for inline T[N]
// struct-array fields.
func TestMarshalFieldInlineStructArray(t *testing.T) {
	t.Run("by-pointer parent emits guarded indexed getter without setter", func(t *testing.T) {
		g := skipGen()
		g.skips = &SkipCollector{}
		o := jen.NewFile("ffmpeg")
		parent := &Struct{Name: "AVParent", ByValue: false}
		field := &Field{
			Name:     "items",
			Type:     &ConstArray{Inner: ident("AVBuf"), Size: 4},
			BitWidth: -1,
		}

		g.marshalField(o, parent, field)

		if got := g.skips.Total(); got != 0 {
			t.Fatalf("skips.Total() = %d, want 0 (entries: %+v)", got, g.skips.Entries)
		}

		code := o.GoString()
		for _, want := range []string{
			"// The returned pointer aliases struct-owned memory; do not free it.",
			"func (s *AVParent) Items(i int) *AVBuf",
			"if s == nil || s.ptr == nil || i < 0 || i >= 4 {",
			"return nil",
			"return &AVBuf{ptr: &s.ptr.items[i]}",
		} {
			if !strings.Contains(code, want) {
				t.Errorf("emitted code missing %q:\n%s", want, code)
			}
		}
		for _, unwanted := range []string{"SetItems", "value :="} {
			if strings.Contains(code, unwanted) {
				t.Errorf("emitted code unexpectedly contains %q:\n%s", unwanted, code)
			}
		}
	})

	t.Run("go-keyword field name uses escaped cgo identifier", func(t *testing.T) {
		g := skipGen()
		g.skips = &SkipCollector{}
		o := jen.NewFile("ffmpeg")
		parent := &Struct{Name: "AVParent", ByValue: false}
		field := &Field{
			Name:     "type",
			Type:     &ConstArray{Inner: ident("AVBuf"), Size: 2},
			BitWidth: -1,
		}

		g.marshalField(o, parent, field)

		if got := g.skips.Total(); got != 0 {
			t.Fatalf("skips.Total() = %d, want 0 (entries: %+v)", got, g.skips.Entries)
		}

		code := o.GoString()
		for _, want := range []string{
			"func (s *AVParent) Type(i int) *AVBuf",
			"if s == nil || s.ptr == nil || i < 0 || i >= 2 {",
			"return &AVBuf{ptr: &s.ptr._type[i]}",
		} {
			if !strings.Contains(code, want) {
				t.Errorf("emitted code missing %q:\n%s", want, code)
			}
		}
	})

	t.Run("by-value parent records skip and emits no method", func(t *testing.T) {
		g := skipGen()
		g.skips = &SkipCollector{}
		o := jen.NewFile("ffmpeg")
		parent := &Struct{Name: "AVParent", ByValue: true}
		field := &Field{
			Name:     "items",
			Type:     &ConstArray{Inner: ident("AVBuf"), Size: 4},
			BitWidth: -1,
		}

		g.marshalField(o, parent, field)

		if got := g.skips.Total(); got != 1 {
			t.Fatalf("skips.Total() = %d, want 1 (entries: %+v)", got, g.skips.Entries)
		}
		if got, want := g.skips.Entries[0].Reason, "ref field of value struct"; got != want {
			t.Errorf("skip reason = %q, want %q", got, want)
		}

		code := o.GoString()
		if strings.Contains(code, "func (s *AVParent)") {
			t.Errorf("emitted code unexpectedly contains a method:\n%s", code)
		}
	})
}

// TestMarshalFieldMultiDimConstArray verifies indexed accessor emission,
// Go-keyword C field escaping, and by-value parent skips for 2D and 3D
// fixed-size array fields.
func TestMarshalFieldMultiDimConstArray(t *testing.T) {
	emittedTests := []struct {
		name         string
		fieldName    string
		fieldType    *ConstArray
		wantContains []string
		wantAbsent   []string
	}{
		{
			name:      "2D primitive emits guarded indexed getter without setter",
			fieldName: "position",
			fieldType: &ConstArray{Inner: &ConstArray{Inner: ident("int16_t"), Size: 2}, Size: 3},
			wantContains: []string{
				"func (s *AVParent) Position(i int) *Array[int16]",
				"if s == nil || s.ptr == nil || i < 0 || i >= 3 {",
				"return nil",
				"return ToInt16Array(unsafe.Pointer(&s.ptr.position[i][0]))",
			},
			wantAbsent: []string{"SetPosition", "value :="},
		},
		{
			name:      "2D char emits uint8 array getter",
			fieldName: "classify_labels",
			fieldType: &ConstArray{Inner: &ConstArray{Inner: ident("char"), Size: 64}, Size: 4},
			wantContains: []string{
				"func (s *AVParent) ClassifyLabels(i int) *Array[uint8]",
				"if s == nil || s.ptr == nil || i < 0 || i >= 4 {",
				"return nil",
				"return ToUint8Array(unsafe.Pointer(&s.ptr.classify_labels[i][0]))",
			},
			wantAbsent: []string{"SetClassifyLabels", "value :="},
		},
		{
			name:      "2D by-value struct emits struct array getter",
			fieldName: "display_primaries",
			fieldType: &ConstArray{Inner: &ConstArray{Inner: ident("AVFrame"), Size: 2}, Size: 3},
			wantContains: []string{
				"func (s *AVParent) DisplayPrimaries(i int) *Array[*AVFrame]",
				"if s == nil || s.ptr == nil || i < 0 || i >= 3 {",
				"return nil",
				"return ToAVFrameArray(unsafe.Pointer(&s.ptr.display_primaries[i][0]))",
			},
			wantAbsent: []string{"SetDisplayPrimaries", "value :="},
		},
		{
			name:      "3D primitive emits two-index guarded getter without setter",
			fieldName: "mmr_coef",
			fieldType: &ConstArray{Inner: &ConstArray{Inner: &ConstArray{Inner: ident("int64_t"), Size: 7}, Size: 3}, Size: 8},
			wantContains: []string{
				"func (s *AVParent) MmrCoef(i, j int) *Array[int64]",
				"if s == nil || s.ptr == nil || i < 0 || i >= 8 || j < 0 || j >= 3 {",
				"return nil",
				"return ToInt64Array(unsafe.Pointer(&s.ptr.mmr_coef[i][j][0]))",
			},
			wantAbsent: []string{"SetMmrCoef", "value :="},
		},
		{
			name:      "go-keyword field name uses escaped cgo identifier",
			fieldName: "type",
			fieldType: &ConstArray{Inner: &ConstArray{Inner: ident("uint8_t"), Size: 4}, Size: 2},
			wantContains: []string{
				"func (s *AVParent) Type(i int) *Array[uint8]",
				"if s == nil || s.ptr == nil || i < 0 || i >= 2 {",
				"return ToUint8Array(unsafe.Pointer(&s.ptr._type[i][0]))",
			},
		},
	}

	for _, tt := range emittedTests {
		t.Run(tt.name, func(t *testing.T) {
			g := skipGen()
			g.skips = &SkipCollector{}
			o := jen.NewFile("ffmpeg")
			parent := &Struct{Name: "AVParent", ByValue: false}
			field := &Field{
				Name:     tt.fieldName,
				Type:     tt.fieldType,
				BitWidth: -1,
			}

			g.marshalField(o, parent, field)

			if got := g.skips.Total(); got != 0 {
				t.Fatalf("skips.Total() = %d, want 0 (entries: %+v)", got, g.skips.Entries)
			}

			code := o.GoString()
			const wantDoc = "// The returned Array aliases struct-owned memory; do not call Free on it."
			if !strings.Contains(code, wantDoc) {
				t.Errorf("emitted code missing %q:\n%s", wantDoc, code)
			}
			for _, want := range tt.wantContains {
				if !strings.Contains(code, want) {
					t.Errorf("emitted code missing %q:\n%s", want, code)
				}
			}
			for _, unwanted := range tt.wantAbsent {
				if strings.Contains(code, unwanted) {
					t.Errorf("emitted code unexpectedly contains %q:\n%s", unwanted, code)
				}
			}
		})
	}

	t.Run("by-value parent records skip and emits no method", func(t *testing.T) {
		g := skipGen()
		g.skips = &SkipCollector{}
		o := jen.NewFile("ffmpeg")
		parent := &Struct{Name: "AVParent", ByValue: true}
		field := &Field{
			Name:     "position",
			Type:     &ConstArray{Inner: &ConstArray{Inner: ident("int16_t"), Size: 2}, Size: 3},
			BitWidth: -1,
		}

		g.marshalField(o, parent, field)

		if got := g.skips.Total(); got != 1 {
			t.Fatalf("skips.Total() = %d, want 1 (entries: %+v)", got, g.skips.Entries)
		}
		if got, want := g.skips.Entries[0].Reason, "ref field of value struct"; got != want {
			t.Errorf("skip reason = %q, want %q", got, want)
		}

		code := o.GoString()
		if strings.Contains(code, "func (s *AVParent)") {
			t.Errorf("emitted code unexpectedly contains a method:\n%s", code)
		}
	})
}

// TestClassifyConstArrayArgShape verifies that fixed-size array parameters
// classify into fixed-width primitive, bounce-buffer primitive, uint8_t
// pointer array, or skipped shapes.
func TestClassifyConstArrayArgShape(t *testing.T) {
	g := skipGen()

	tests := []struct {
		name          string
		array         *ConstArray
		paramName     string
		wantKind      constArrayArgShapeKind
		wantSize      int64
		wantElemConst bool
		wantGoType    string
		wantCType     string
		wantReason    string
	}{
		{
			name:          "fixed-width const uint8_t[16]",
			array:         &ConstArray{Inner: ident("uint8_t"), Size: 16, ElemConst: true},
			paramName:     "key",
			wantKind:      constArrayArgShapeFixedPrimitive,
			wantSize:      16,
			wantElemConst: true,
			wantGoType:    "uint8",
			wantCType:     "uint8_t",
		},
		{
			name:          "fixed-width const uint32_t[4]",
			array:         &ConstArray{Inner: ident("uint32_t"), Size: 4, ElemConst: true},
			paramName:     "color",
			wantKind:      constArrayArgShapeFixedPrimitive,
			wantSize:      4,
			wantElemConst: true,
			wantGoType:    "uint32",
			wantCType:     "uint32_t",
		},
		{
			name:          "fixed-width double[2]",
			array:         &ConstArray{Inner: ident("double"), Size: 2},
			paramName:     "out",
			wantKind:      constArrayArgShapeFixedPrimitive,
			wantSize:      2,
			wantElemConst: false,
			wantGoType:    "float64",
			wantCType:     "double",
		},
		{
			name:          "bounce-buffer int[4]",
			array:         &ConstArray{Inner: ident("int"), Size: 4},
			paramName:     "linesizes",
			wantKind:      constArrayArgShapeBouncePrimitive,
			wantSize:      4,
			wantElemConst: false,
			wantGoType:    "int",
			wantCType:     "int",
		},
		{
			name:          "bounce-buffer const int[4]",
			array:         &ConstArray{Inner: ident("int"), Size: 4, ElemConst: true},
			paramName:     "linesizes",
			wantKind:      constArrayArgShapeBouncePrimitive,
			wantSize:      4,
			wantElemConst: true,
			wantGoType:    "int",
			wantCType:     "int",
		},
		{
			name:          "bounce-buffer uint[4]",
			array:         &ConstArray{Inner: ident("uint"), Size: 4},
			paramName:     "counts",
			wantKind:      constArrayArgShapeBouncePrimitive,
			wantSize:      4,
			wantElemConst: false,
			wantGoType:    "uint",
			wantCType:     "uint",
		},
		{
			name:          "bounce-buffer const ptrdiff_t[4]",
			array:         &ConstArray{Inner: ident("ptrdiff_t"), Size: 4, ElemConst: true},
			paramName:     "strides",
			wantKind:      constArrayArgShapeBouncePrimitive,
			wantSize:      4,
			wantElemConst: true,
			wantGoType:    "int64",
			wantCType:     "ptrdiff_t",
		},
		{
			name:          "bounce-buffer size_t[4]",
			array:         &ConstArray{Inner: ident("size_t"), Size: 4},
			paramName:     "size",
			wantKind:      constArrayArgShapeBouncePrimitive,
			wantSize:      4,
			wantElemConst: false,
			wantGoType:    "uint64",
			wantCType:     "size_t",
		},
		{
			name:          "bounce-buffer long[4]",
			array:         &ConstArray{Inner: ident("long"), Size: 4},
			paramName:     "offsets",
			wantKind:      constArrayArgShapeBouncePrimitive,
			wantSize:      4,
			wantElemConst: false,
			wantGoType:    "int64",
			wantCType:     "long",
		},
		{
			name:          "bounce-buffer ulong[4]",
			array:         &ConstArray{Inner: ident("ulong"), Size: 4},
			paramName:     "masks",
			wantKind:      constArrayArgShapeBouncePrimitive,
			wantSize:      4,
			wantElemConst: false,
			wantGoType:    "uint64",
			wantCType:     "ulong",
		},
		{
			name:          "uint8_t pointer array",
			array:         &ConstArray{Inner: ptr(ident("uint8_t")), Size: 4},
			paramName:     "dstData",
			wantKind:      constArrayArgShapeUint8Pointer,
			wantSize:      4,
			wantElemConst: false,
		},
		{
			name:          "const uint8_t pointer array",
			array:         &ConstArray{Inner: ptr(ident("uint8_t")), Size: 4, ElemConst: true},
			paramName:     "srcData",
			wantKind:      constArrayArgShapeUint8Pointer,
			wantSize:      4,
			wantElemConst: true,
		},
		{
			name:       "nil array skipped",
			array:      nil,
			paramName:  "missing",
			wantKind:   constArrayArgShapeSkip,
			wantReason: "const array param missing",
		},
		{
			name:       "zero-length array skipped",
			array:      &ConstArray{Inner: ident("int"), Size: 0},
			paramName:  "empty",
			wantKind:   constArrayArgShapeSkip,
			wantReason: "const array param empty",
		},
		{
			name:       "negative-length array skipped",
			array:      &ConstArray{Inner: ident("int"), Size: -1},
			paramName:  "bad",
			wantKind:   constArrayArgShapeSkip,
			wantReason: "const array param bad",
		},
		{
			name:       "struct element array skipped",
			array:      &ConstArray{Inner: ident("AVBuf"), Size: 4},
			paramName:  "bufs",
			wantKind:   constArrayArgShapeSkip,
			wantReason: "const array param bufs",
		},
		{
			name:       "void pointer array skipped",
			array:      &ConstArray{Inner: ptr(nil), Size: 2},
			paramName:  "ptrs",
			wantKind:   constArrayArgShapeSkip,
			wantReason: "const array param ptrs",
		},
		{
			name:       "int pointer array skipped",
			array:      &ConstArray{Inner: ptr(ident("int")), Size: 4},
			paramName:  "iptrs",
			wantKind:   constArrayArgShapeSkip,
			wantReason: "const array param iptrs",
		},
		{
			name:       "nested 2D array skipped",
			array:      &ConstArray{Inner: &ConstArray{Inner: ident("int"), Size: 4}, Size: 2},
			paramName:  "matrix",
			wantKind:   constArrayArgShapeSkip,
			wantReason: "const array param matrix",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := g.classifyConstArrayArgShape(tt.array, tt.paramName)
			if got.kind != tt.wantKind {
				t.Errorf("kind = %v, want %v", got.kind, tt.wantKind)
			}
			if got.size != tt.wantSize {
				t.Errorf("size = %d, want %d", got.size, tt.wantSize)
			}
			if got.elemConst != tt.wantElemConst {
				t.Errorf("elemConst = %v, want %v", got.elemConst, tt.wantElemConst)
			}
			if got.goType != tt.wantGoType {
				t.Errorf("goType = %q, want %q", got.goType, tt.wantGoType)
			}
			if got.cType != tt.wantCType {
				t.Errorf("cType = %q, want %q", got.cType, tt.wantCType)
			}
			if got.reason != tt.wantReason {
				t.Errorf("reason = %q, want %q", got.reason, tt.wantReason)
			}
		})
	}
}

// TestMarshalConstArrayArg verifies emitted Go parameter types, nil guards,
// zero-copy casts, stack bounce buffers, width checks, post-call copy-back
// gating on ElemConst, and [4]*C.uint8_t marshalling.
func TestMarshalConstArrayArg(t *testing.T) {
	tests := []struct {
		name             string
		fnName           string
		argName          string
		array            *ConstArray
		wantParam        string
		wantArg          string
		wantBodyContains []string
		wantBodyAbsent   []string
		wantPostContains []string
		wantNoPost       bool
	}{
		{
			name:      "fixed-width primitive emits zero-copy cast with nil guard",
			fnName:    "av_tea_init",
			argName:   "key",
			array:     &ConstArray{Inner: ident("uint8_t"), Size: 16, ElemConst: true},
			wantParam: "var key *[16]uint8",
			wantArg:   "tmpkey",
			wantBodyContains: []string{
				"var tmpkey *C.uint8_t",
				"if key != nil {",
				"tmpkey = (*C.uint8_t)(unsafe.Pointer(&key[0]))",
			},
			wantBodyAbsent: []string{"[16]C.uint8_t", "overflows"},
			wantNoPost:     true,
		},
		{
			name:      "mutable int[4] emits stack bounce buffer and post-call copy-back",
			fnName:    "av_image_fill_linesizes",
			argName:   "linesizes",
			array:     &ConstArray{Inner: ident("int"), Size: 4, ElemConst: false},
			wantParam: "var linesizes *[4]int",
			wantArg:   "ptrlinesizes",
			wantBodyContains: []string{
				"var tmplinesizes [4]C.int",
				"var ptrlinesizes *C.int",
				"if linesizes != nil {",
				"for idx := range tmplinesizes {",
				"tmplinesizes[idx] = C.int(linesizes[idx])",
				"ptrlinesizes = &tmplinesizes[0]",
			},
			wantBodyAbsent: []string{"unsafe.Pointer", "overflows"},
			wantPostContains: []string{
				"if linesizes != nil {",
				"for idx := range tmplinesizes {",
				"linesizes[idx] = int(tmplinesizes[idx])",
			},
		},
		{
			name:      "const int[4] emits stack bounce buffer without post-call copy-back",
			fnName:    "av_const_ints",
			argName:   "vals",
			array:     &ConstArray{Inner: ident("int"), Size: 4, ElemConst: true},
			wantParam: "var vals *[4]int",
			wantArg:   "ptrvals",
			wantBodyContains: []string{
				"var tmpvals [4]C.int",
				"var ptrvals *C.int",
				"if vals != nil {",
				"tmpvals[idx] = C.int(vals[idx])",
				"ptrvals = &tmpvals[0]",
			},
			wantBodyAbsent: []string{"overflows"},
			wantNoPost:     true,
		},
		{
			name:      "const ptrdiff_t[4] emits width check without post-call copy-back",
			fnName:    "av_image_fill_plane_sizes",
			argName:   "linesizes",
			array:     &ConstArray{Inner: ident("ptrdiff_t"), Size: 4, ElemConst: true},
			wantParam: "var linesizes *[4]int64",
			wantArg:   "ptrlinesizes",
			wantBodyContains: []string{
				"var tmplinesizes [4]C.ptrdiff_t",
				"var ptrlinesizes *C.ptrdiff_t",
				"if linesizes != nil {",
				"if int64(C.ptrdiff_t(linesizes[idx])) != linesizes[idx] {",
				`panic("av_image_fill_plane_sizes.linesizes overflows C.ptrdiff_t")`,
				"tmplinesizes[idx] = C.ptrdiff_t(linesizes[idx])",
				"ptrlinesizes = &tmplinesizes[0]",
			},
			wantNoPost: true,
		},
		{
			name:      "mutable size_t[4] emits width check and post-call copy-back",
			fnName:    "av_image_fill_plane_sizes",
			argName:   "size",
			array:     &ConstArray{Inner: ident("size_t"), Size: 4, ElemConst: false},
			wantParam: "var size *[4]uint64",
			wantArg:   "ptrsize",
			wantBodyContains: []string{
				"var tmpsize [4]C.size_t",
				"var ptrsize *C.size_t",
				"if size != nil {",
				"if uint64(C.size_t(size[idx])) != size[idx] {",
				`panic("av_image_fill_plane_sizes.size overflows C.size_t")`,
				"tmpsize[idx] = C.size_t(size[idx])",
				"ptrsize = &tmpsize[0]",
			},
			wantPostContains: []string{
				"if size != nil {",
				"for idx := range tmpsize {",
				"size[idx] = uint64(tmpsize[idx])",
			},
		},
		{
			name:      "mutable uint8_t*[4] emits [4]*C.uint8_t bounce buffer and copy-back",
			fnName:    "av_image_fill_black",
			argName:   "dst_data",
			array:     &ConstArray{Inner: ptr(ident("uint8_t")), Size: 4, ElemConst: false},
			wantParam: "var dstData *[4]unsafe.Pointer",
			wantArg:   "ptrdstData",
			wantBodyContains: []string{
				"var tmpdstData [4]*C.uint8_t",
				"var ptrdstData **C.uint8_t",
				"if dstData != nil {",
				"for idx := range tmpdstData {",
				"tmpdstData[idx] = (*C.uint8_t)(dstData[idx])",
				"ptrdstData = &tmpdstData[0]",
			},
			wantPostContains: []string{
				"if dstData != nil {",
				"for idx := range tmpdstData {",
				"dstData[idx] = unsafe.Pointer(tmpdstData[idx])",
			},
		},
		{
			name:      "const uint8_t*[4] emits [4]*C.uint8_t bounce buffer without copy-back",
			fnName:    "av_image_read_only",
			argName:   "src_data",
			array:     &ConstArray{Inner: ptr(ident("uint8_t")), Size: 4, ElemConst: true},
			wantParam: "var srcData *[4]unsafe.Pointer",
			wantArg:   "ptrsrcData",
			wantBodyContains: []string{
				"var tmpsrcData [4]*C.uint8_t",
				"var ptrsrcData **C.uint8_t",
				"if srcData != nil {",
				"tmpsrcData[idx] = (*C.uint8_t)(srcData[idx])",
				"ptrsrcData = &tmpsrcData[0]",
			},
			wantNoPost: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := skipGen()
			g.skips = &SkipCollector{}
			fn := &Function{Name: tt.fnName}
			arg := &Param{Name: tt.argName, Type: tt.array}

			params, args, body, postCall, skip := g.marshalArg(newFile(), fn, arg)
			if skip {
				t.Fatalf("%s.%s unexpectedly skipped (entries: %+v)", tt.fnName, tt.argName, g.skips.Entries)
			}
			if got := g.skips.Total(); got != 0 {
				t.Fatalf("skips.Total() = %d, want 0", got)
			}
			if got := render(jen.Var().Add(params[0])); got != tt.wantParam {
				t.Errorf("param = %q, want %q", got, tt.wantParam)
			}
			if got := render(args[0]); got != tt.wantArg {
				t.Errorf("arg = %q, want %q", got, tt.wantArg)
			}

			bodyCode := render(jen.Block(body...))
			for _, want := range tt.wantBodyContains {
				if !strings.Contains(bodyCode, want) {
					t.Errorf("body missing %q:\n%s", want, bodyCode)
				}
			}
			for _, unwanted := range tt.wantBodyAbsent {
				if strings.Contains(bodyCode, unwanted) {
					t.Errorf("body unexpectedly contains %q:\n%s", unwanted, bodyCode)
				}
			}

			if tt.wantNoPost {
				if len(postCall) != 0 {
					t.Errorf("postCall len = %d, want 0:\n%s", len(postCall), render(jen.Block(postCall...)))
				}
				return
			}

			postCode := render(jen.Block(postCall...))
			for _, want := range tt.wantPostContains {
				if !strings.Contains(postCode, want) {
					t.Errorf("postCall missing %q:\n%s", want, postCode)
				}
			}
		})
	}
}

// TestManuallyWrappedFuncsSkipped pins the 17-entry manuallyWrappedFuncs table
// and verifies that generateFuncs skips every entry with its expected reason
// without emitting a colliding Go wrapper.
func TestManuallyWrappedFuncsSkipped(t *testing.T) {
	want := map[string]string{
		"av_display_matrix_flip":   "manually wrapped in display.go",
		"av_display_rotation_get":  "manually wrapped in display.go",
		"av_display_rotation_set":  "manually wrapped in display.go",
		"av_image_alloc":           "manually wrapped in image.go",
		"av_image_copy":            "manually wrapped in image.go",
		"av_image_copy2":           "manually wrapped in image.go",
		"av_image_copy_to_buffer":  "manually wrapped in image.go",
		"av_image_fill_arrays":     "manually wrapped in image.go",
		"av_image_fill_pointers":   "manually wrapped in image.go",
		"av_uuid_copy":             "array typedef (manually wrapped in uuid.go)",
		"av_uuid_equal":            "array typedef (manually wrapped in uuid.go)",
		"av_uuid_nil":              "array typedef (manually wrapped in uuid.go)",
		"av_uuid_parse":            "array typedef (manually wrapped in uuid.go)",
		"av_uuid_parse_range":      "array typedef (manually wrapped in uuid.go)",
		"av_uuid_unparse":          "array typedef (manually wrapped in uuid.go)",
		"av_uuid_urn_parse":        "array typedef (manually wrapped in uuid.go)",
		"sws_setColorspaceDetails": "manually wrapped in swscale.go",
	}

	if len(manuallyWrappedFuncs) != len(want) {
		t.Fatalf("len(manuallyWrappedFuncs) = %d, want %d", len(manuallyWrappedFuncs), len(want))
	}
	for name, wantReason := range want {
		gotReason, ok := manuallyWrappedFuncs[name]
		if !ok {
			t.Errorf("manuallyWrappedFuncs missing %q", name)
			continue
		}
		if gotReason != wantReason {
			t.Errorf("manuallyWrappedFuncs[%q] = %q, want %q", name, gotReason, wantReason)
		}
	}

	names := make([]string, 0, len(want))
	funcs := make(map[string]*Function, len(want))
	for name := range want {
		names = append(names, name)
		funcs[name] = &Function{Name: name}
	}
	slices.Sort(names)

	dir := t.TempDir()
	t.Chdir(dir)

	origOut := log.Writer()
	log.SetOutput(io.Discard)
	defer log.SetOutput(origOut)

	g := &Generator{
		input: &Module{
			functions:     funcs,
			functionOrder: names,
			structs:       map[string]*Struct{},
			enums:         map[string]*Enum{},
			callbacks:     map[string]*Function{},
		},
		skips: &SkipCollector{},
	}

	g.generateFuncs()

	if got := g.skips.Total(); got != len(want) {
		t.Fatalf("skips.Total() = %d, want %d", got, len(want))
	}

	gotReasons := make(map[string]string, len(g.skips.Entries))
	for _, entry := range g.skips.Entries {
		gotReasons[entry.Symbol] = entry.Reason
	}
	for name, wantReason := range want {
		if gotReason, ok := gotReasons[name]; !ok {
			t.Errorf("generateFuncs did not record skip for %q", name)
		} else if gotReason != wantReason {
			t.Errorf("skip reason for %q = %q, want %q", name, gotReason, wantReason)
		}
	}

	raw, err := os.ReadFile(filepath.Join(dir, "functions.gen.go"))
	if err != nil {
		t.Fatalf("read generated functions.gen.go: %v", err)
	}
	code := string(raw)
	if strings.Contains(code, "\nfunc ") {
		t.Errorf("generateFuncs unexpectedly emitted a wrapper for a manually wrapped function:\n%s", code)
	}
	for name, wantReason := range want {
		wantComment := "// " + name + " skipped due to " + wantReason
		if !strings.Contains(code, wantComment) {
			t.Errorf("functions.gen.go missing skip comment %q", wantComment)
		}
	}
}

// Other marshalField early-exit decisions (skippedFields, bitfield,
// ident-callback) stay covered by the byte-identical regeneration gate rather
// than separate unit tests.
