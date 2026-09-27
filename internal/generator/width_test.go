package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dave/jennifer/jen"
)

func TestVariableWidthChecks(t *testing.T) {
	g := skipGen()
	for _, tt := range []struct{ cType, goType string }{
		{"size_t", "uint64"}, {"ulong", "uint64"}, {"ptrdiff_t", "int64"},
	} {
		t.Run(tt.cType, func(t *testing.T) {
			_, _, body, _, skip := g.marshalArg(newFile(), &Function{Name: "consume"}, &Param{Name: "value", Type: ident(tt.cType)})
			if skip || len(body) != 1 {
				t.Fatalf("missing input check: skip=%v body=%v", skip, body)
			}
			want := tt.goType + "(C." + tt.cType + "(value)) != value"
			if got := render(body[0]); !strings.Contains(got, want) {
				t.Fatalf("input check = %s, want %s", got, want)
			}
			o := jen.NewFile("fixture")
			g.marshalField(o, &Struct{Name: "Sample"}, &Field{Name: "count", Type: ident(tt.cType), CTypeName: tt.cType, BitWidth: -1})
			code := o.GoString()
			check := strings.Index(code, want)
			write := strings.Index(code, "s.ptr.count =")
			if check < 0 || write < check {
				t.Fatalf("setter must check before write:\n%s", code)
			}
		})
	}
	if got := primitiveWidthCheck("uint64_t", "uint64", jen.Id("value"), "fixed"); len(got) != 0 {
		t.Fatal("fixed-width input gained a check")
	}
}

// Compile the emitted conversions against both C integer widths without cgo.
// This exercises the generated statements, not a second copy of their logic.
func TestVariableWidthConversions(t *testing.T) {
	g := skipGen()
	fn := &Function{Name: "av_dovi_alloc"}
	params, args, body, post, skip := g.marshalArg(newFile(), fn, &Param{Name: "size", Type: ptr(ident("size_t"))})
	if skip {
		t.Fatal("size output skipped")
	}
	ret, body, skip := g.marshalReturn(newFile(), fn, ident("uint64_t"), jen.Id("writeSize").Params(args...), body, post)
	if skip {
		t.Fatal("return skipped")
	}
	output := render(jen.Func().Id("output").Params(params...).Add(ret...).Block(body...))
	if strings.Contains(output, "unsafe.Pointer") {
		t.Fatalf("output still aliases Go storage: %s", output)
	}

	var checks strings.Builder
	for _, tt := range []struct{ cType, goType string }{
		{"size_t", "uint64"}, {"ulong", "uint64"}, {"ptrdiff_t", "int64"},
	} {
		params, args, body, _, skip := g.marshalArg(newFile(), &Function{Name: "consume"}, &Param{Name: "value", Type: ident(tt.cType)})
		if skip {
			t.Fatal("scalar input skipped")
		}
		body = append(body, jen.Return(jen.Id(tt.goType).Params(args[0])))
		checks.WriteString(render(jen.Func().Id("check_" + tt.cType).Params(params...).Id(tt.goType).Block(body...)))
		checks.WriteByte('\n')
	}

	for _, bits := range []int{32, 64} {
		t.Run(fmt.Sprint(bits), func(t *testing.T) {
			dir := t.TempDir()
			source := "package fixture\n" + output + "\n" + checks.String()
			source = strings.NewReplacer("C.size_t", fmt.Sprintf("uint%d", bits), "C.ulong", "uint32", "C.ptrdiff_t", fmt.Sprintf("int%d", bits)).Replace(source)
			source += fmt.Sprintf(`
var noWrite bool
var nilSeen bool
func writeSize(p *uint%d) uint64 {
	if p == nil { nilSeen = true } else if !noWrite { *p = 7 }
	return 9
}
`, bits)
			tests := fmt.Sprintf(`package fixture
import "testing"
func TestOutput(t *testing.T) {
	v := ^uint64(0)
	if ret := output(&v); ret != 9 || v != 7 { t.Fatalf("ret=%%d value=%%d", ret, v) }
	output(nil)
	if !nilSeen { t.Fatal("nil pointer changed") }
	noWrite = true
	v = ^uint64(0)
	output(&v)
	if v != 0 { t.Fatalf("unwritten output=%%d", v) }
}
func TestBounds(t *testing.T) {
	unsignedMax := ^uint64(0) >> (64 - %d)
	signedMax := int64(unsignedMax >> 1)
	signedMin := -signedMax - 1
	if check_size_t(unsignedMax) != unsignedMax || check_size_t(0) != 0 { t.Fatal("size boundary") }
	if check_ptrdiff_t(signedMin) != signedMin || check_ptrdiff_t(signedMax) != signedMax { t.Fatal("signed boundary") }
	if check_ulong(1<<32-1) != 1<<32-1 { t.Fatal("long boundary") }
	panics(t, func() { check_ulong(1<<32) })
	if %d == 32 {
		panics(t, func() { check_size_t(unsignedMax+1) })
		panics(t, func() { check_ptrdiff_t(signedMin-1) })
		panics(t, func() { check_ptrdiff_t(signedMax+1) })
	}
}
func panics(t *testing.T, f func()) {
	t.Helper()
	defer func() { if recover() == nil { t.Error("overflow did not panic") } }()
	f()
}
`, bits, bits)
			for name, content := range map[string]string{"go.mod": "module fixture\n\ngo 1.25\n", "width.go": source, "width_test.go": tests} {
				if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			cmd := exec.CommandContext(t.Context(), "go", "test", "-count=1", ".")
			cmd.Dir = dir
			cmd.Env = append(os.Environ(), "GOWORK=off", "CGO_ENABLED=0")
			if out, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("%d-bit emitted conversions: %v\n%s\n%s", bits, err, out, source)
			}
		})
	}
}
