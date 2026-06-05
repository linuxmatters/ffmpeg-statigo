package main

import (
	"fmt"
	"log"
	"regexp"
	"slices"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
)

var (
	fileType = &PointerType{
		Inner: &IdentType{Name: "FILE"},
	}
	fileType2 = &PointerType{
		Inner: &IdentType{Name: "_IO_FILE"},
	}
	vaListType = &IdentType{Name: "va_list"}
	tmType     = &IdentType{Name: "tm"}
)

func (g *Generator) generateFuncs() {
	i := g.input

	o := newFile()

outer:
	for _, fnName := range i.functionOrder {
		fn := i.functions[fnName]

		log.Println("Generating func", fn.Name)
		o.Commentf("--- Function %v ---", fn.Name)
		o.Line()

		// Check if function contains unsupported features
		if fn.Variadic {
			o.Commentf("%v skipped due to variadic arg.", fn.Name)
			o.Line()
			g.skips.Record(fn.Name, "variadic arg")

			continue
		}

		// WORKAROUND: libclang on Linux reports FILE* as int* for av_fopen_utf8
		// Pinned by TestGeneratorSkipPatterns in bindings_test.go (asserts AVFopenUtf8 absent).
		if fn.Name == "av_fopen_utf8" {
			o.Commentf("%v skipped due to return", fn.Name)
			o.Line()
			g.skips.Record(fn.Name, "return")

			continue outer
		}

		// WORKAROUND: UUID functions use AVUUID (array typedef) which requires pointer conversion in CGO
		// These are manually wrapped in uuid.go with proper pointer handling
		// Pinned by TestUUIDBindingsNotDuplicated in bindings_test.go (asserts the seven manual wrappers are not also generated).
		if fn.Name == "av_uuid_parse" || fn.Name == "av_uuid_urn_parse" || fn.Name == "av_uuid_parse_range" ||
			fn.Name == "av_uuid_unparse" || fn.Name == "av_uuid_equal" || fn.Name == "av_uuid_copy" || fn.Name == "av_uuid_nil" {
			o.Commentf("%v skipped due to array typedef (manually wrapped in uuid.go)", fn.Name)
			o.Line()
			g.skips.Record(fn.Name, "array typedef (manually wrapped in uuid.go)")

			continue outer
		}

		if typeEquals(fn.Result, fileType) || typeEquals(fn.Result, fileType2) {
			o.Commentf("%v skipped due to return", fn.Name)
			o.Line()
			g.skips.Record(fn.Name, "return")

			continue outer
		}

		for _, arg := range fn.Args {
			if g.classifyFunctionArgPreSkip(arg.Type) {
				o.Commentf("%v skipped due to %v.", fn.Name, arg.Name)
				o.Line()
				g.skips.Record(fn.Name, fmt.Sprintf("arg %v", arg.Name))

				continue outer
			}
		}

		goName := g.convCamel(fn.Name)

		var (
			params   []jen.Code
			args     []jen.Code
			body     []jen.Code
			postCall []jen.Code
		)

		for _, arg := range fn.Args {
			params2, args2, body2, postCall2, skip := g.marshalArg(o, fn, arg)
			if skip {
				continue outer
			}

			params = append(params, params2...)
			args = append(args, args2...)
			body = append(body, body2...)
			postCall = append(postCall, postCall2...)
		}

		cc := jen.Qual("C", fn.Name).Params(args...)

		retType, body, skip := g.marshalReturn(o, fn, fn.Result, cc, body, postCall)
		if skip {
			continue outer
		}

		o.Commentf("%v wraps %v.", goName, fn.Name)

		if fn.Comment != "" {
			o.Comment(fn.Comment)
		}

		o.Func().
			// Params(jen.Id("s").Op("*").Id(goName)).
			Id(goName).
			Params(params...).
			Add(retType...).
			Block(body...)
	}

	err := saveFormatted(o, "functions.gen.go")
	if err != nil {
		log.Println("ERROR saving functions.gen.go:", err)
		log.Panicln(err)
	}
}

func (g *Generator) marshalReturn(o *jen.File, fn *Function, result Type, cc jen.Code, body, postCall []jen.Code) (retType, outBody []jen.Code, skip bool) {
	shape := g.classifyReturnShape(result)

	switch shape.kind {
	case returnShapeVoid:
		// nothing
		body = append(body, cc)
		body = append(body, postCall...)

	case returnShapeInt:
		body = append(body, jen.Id("ret").Op(":=").Add(cc))
		body = append(body, postCall...)
		retType = []jen.Code{jen.Params(jen.Id("int"), jen.Id("error"))}
		body = append(
			body,
			jen.Return(
				jen.Id("int").Params(jen.Id("ret")).Op(",").
					Id("WrapErr").Params(jen.Id("int").Params(jen.Id("ret"))),
			),
		)

	case returnShapePrimitive:
		body = append(body, jen.Id("ret").Op(":=").Add(cc))
		body = append(body, postCall...)
		retType = []jen.Code{jen.Id(shape.goType)}
		body = append(body, jen.Return(jen.Id(shape.goType).Params(jen.Id("ret"))))

	case returnShapeByValueStruct:
		body = append(body, jen.Id("ret").Op(":=").Add(cc))
		body = append(body, postCall...)
		retType = []jen.Code{jen.Op("*").Id(shape.name)}
		body = append(
			body,
			jen.Return(jen.Op("&").Id(shape.name).Values(jen.Dict{
				jen.Id("value"): jen.Id("ret"),
			})),
		)

	case returnShapeByPointerStructSkip:
		o.Commentf("%v skipped due to return", fn.Name)
		o.Line()
		g.skips.Record(fn.Name, shape.reason)
		return nil, nil, true

	case returnShapeCallback:
		body = append(body, jen.Id("ret").Op(":=").Add(cc))
		body = append(body, postCall...)
		retType = []jen.Code{jen.Id(shape.goTypeName)}
		body = append(body, jen.Return(jen.Id(shape.goTypeName).Params(jen.Id("ret"))))

	case returnShapeIdent:
		body = append(body, jen.Id("ret").Op(":=").Add(cc))
		body = append(body, postCall...)
		retType = []jen.Code{jen.Id(shape.name)}
		body = append(body, jen.Return(jen.Id(shape.name).Params(jen.Id("ret"))))

	case returnShapeUnsafePointer:
		body = append(
			body,
			jen.Id("ret").Op(":=").Add(cc),
		)
		body = append(body, postCall...)
		retType = []jen.Code{
			jen.Qual("unsafe", "Pointer"),
		}
		body = append(body, jen.Return(jen.Id("ret")))

	case returnShapeCStrPointer:
		body = append(
			body,
			jen.Id("ret").Op(":=").Add(cc),
		)
		body = append(body, postCall...)
		retType = []jen.Code{
			jen.Op("*").Id("CStr"),
		}
		body = append(body, jen.Return(jen.Id("wrapCStr").Params(jen.Id("ret"))))

	case returnShapeUint8Pointer:
		body = append(
			body,
			jen.Id("ret").Op(":=").Add(cc),
		)
		body = append(body, postCall...)
		retType = []jen.Code{
			jen.Qual("unsafe", "Pointer"),
		}
		body = append(
			body,
			jen.Return(jen.Qual("unsafe", "Pointer").Params(jen.Id("ret"))),
		)

	case returnShapePrimitivePointer:
		body = append(
			body,
			jen.Id("ret").Op(":=").Add(cc),
		)
		body = append(body, postCall...)
		retType = []jen.Code{
			jen.Op("*").Id(shape.goType),
		}
		body = append(
			body,
			jen.Return(jen.Params(jen.Op("*").Id(shape.goType)).Params(
				jen.Qual("unsafe", "Pointer").Params(jen.Id("ret")),
			)),
		)

	case returnShapeStructPointer:
		body = append(
			body,
			jen.Id("ret").Op(":=").Add(cc),
		)
		body = append(body, postCall...)
		retType = []jen.Code{
			jen.Op("*").Id(shape.name),
		}

		body = append(
			body,
			jen.Var().Id("retMapped").Op("*").Id(shape.name),
			jen.If(jen.Id("ret").Op("!=").Id("nil")).Block(
				jen.Id("retMapped").Op("=").Op("&").Id(shape.name).Values(jen.Dict{
					jen.Id("ptr"): jen.Id("ret"),
				}),
			),
			jen.Return(jen.Id("retMapped")),
		)

	case returnShapeUnknownPointer:
		body = append(
			body,
			jen.Id("ret").Op(":=").Add(cc),
		)
		body = append(body, postCall...)
		retType = []jen.Code{
			jen.Op("*").Id(shape.name),
		}
		body = append(
			body,
			jen.Return(jen.Params(jen.Op("*").Id(shape.name)).Params(
				jen.Qual("unsafe", "Pointer").Params(jen.Id("ret")),
			)),
		)

	case returnShapePointerPointerSkip, returnShapeUnhandledPointerSkip, returnShapeUnhandledSkip:
		o.Commentf("%v skipped due to %v", fn.Name, shape.reason)
		o.Line()
		g.skips.Record(fn.Name, shape.reason)
		return nil, nil, true
	}

	return retType, body, false
}

func (g *Generator) marshalArg(o *jen.File, fn *Function, arg *Param) (params, args, body, postCall []jen.Code, skip bool) {
	shape := g.classifyArgShape(fn, arg)

	switch shape.kind {
	case argShapeIdentPrimitive:
		params = append(params, jen.Id(shape.name).Id(shape.goType))
		args = append(args, jen.Qual("C", shape.cType).Params(jen.Id(shape.name)))

	case argShapeIdentEnum:
		params = append(params, jen.Id(shape.name).Id(shape.typeName))
		args = append(args, jen.Qual("C", shape.enum.CName()).Params(jen.Id(shape.name)))

	case argShapeIdentByValueStruct:
		params = append(params, jen.Id(shape.name).Op("*").Id(shape.st.Name))
		args = append(args, jen.Id(shape.name).Dot("value"))

	case argShapeIdentByPointerStructSkip, argShapeIdentCallbackByValueSkip, argShapeArraySkip:
		o.Commentf("%v skipped due to %v", fn.Name, shape.reason)
		o.Line()
		g.skips.Record(fn.Name, shape.reason)
		return params, args, body, postCall, true

	case argShapeIdentUnknown:
		params = append(params, jen.Id(shape.name).Id(shape.typeName))
		args = append(args, jen.Qual("C", shape.typeName).Params(jen.Id(shape.name)))

	case argShapePointer:
		var pSkip bool
		params, args, body, postCall, pSkip = g.marshalPointerArg(o, fn, arg, shape.ptr, shape.name, shape.actualTypeName)
		if pSkip {
			return params, args, body, postCall, true
		}

	case argShapeConstArraySkip, argShapeUnhandledSkip:
		o.Commentf("%v skipped due to %v", fn.Name, shape.reason)
		o.Line()
		g.skips.Record(fn.Name, shape.reason)
		return params, args, body, postCall, true
	}

	return params, args, body, postCall, false
}

func (g *Generator) marshalPointerArg(o *jen.File, fn *Function, arg *Param, v *PointerType, pName, actualTypeName string) (params, args, body, postCall []jen.Code, skip bool) {
	shape := g.classifyPointerArgShape(fn, arg, v, pName, actualTypeName)

	switch shape.kind {
	case pointerArgShapeUnsafePointer:
		params = append(params, jen.Id(pName).Qual("unsafe", "Pointer"))
		args = append(args, jen.Id(pName))

	case pointerArgShapeCStr:
		params = append(params, jen.Id(pName).Op("*").Id("CStr"))
		convName := fmt.Sprintf("tmp%v", pName)

		body = append(
			body,
			jen.Var().Id(convName).Op("*").Qual("C", "char"),
			jen.If(jen.Id(pName).Op("!=").Id("nil")).Block(
				jen.Id(convName).Op("=").Id(pName).Dot("ptr"),
			),
		)

		args = append(args, jen.Id(convName))

	case pointerArgShapeBytePointer:
		params = append(params, jen.Id(pName).Qual("unsafe", "Pointer"))
		args = append(args, jen.Params(jen.Op("*").Qual("C", shape.typeName)).Params(jen.Id(pName)))

	case pointerArgShapeOutputPrimitive:
		params = append(params, jen.Id(pName).Op("*").Id(shape.goType))
		args = append(args, jen.Params(jen.Op("*").Qual("C", shape.cType)).Params(jen.Qual("unsafe", "Pointer").Params(jen.Id(pName))))

	case pointerArgShapeNonOutputPrimitiveSkip, pointerArgShapeByValueStructSkip, pointerArgShapePointerToPointerSkip, pointerArgShapeUnhandledSkip:
		o.Commentf("%v skipped due to %v", fn.Name, shape.reason)
		o.Line()
		g.skips.Record(fn.Name, shape.reason)
		return params, args, body, postCall, true

	case pointerArgShapeByPointerStruct:
		params = append(params, jen.Id(pName).Op("*").Id(shape.typeName))

		convName := fmt.Sprintf("tmp%v", pName)

		body = append(
			body,
			jen.Var().Id(convName).Op("*").Qual("C", shape.st.CName()),
			jen.If(jen.Id(pName).Op("!=").Id("nil")).Block(
				jen.Id(convName).Op("=").Id(pName).Dot("ptr"),
			),
		)

		args = append(args, jen.Id(convName))

	case pointerArgShapeEnum:
		params = append(params, jen.Id(pName).Op("*").Id(shape.typeName))

		convName := fmt.Sprintf("tmp%v", pName)

		body = append(
			body,
			jen.Var().Id(convName).Op("*").Qual("C", shape.enum.CName()),
			jen.If(jen.Id(pName).Op("!=").Id("nil")).Block(
				jen.Id(convName).Op("=").Params(jen.Op("*").Qual("C", shape.enum.CName())).Params(jen.Qual("unsafe", "Pointer").Params(jen.Id(pName))),
			),
		)

		args = append(args, jen.Id(convName))

	case pointerArgShapeCallback:
		goTypeName := g.convCamel(shape.typeName)
		params = append(params, jen.Id(pName).Op("*").Id(goTypeName))

		convName := fmt.Sprintf("tmp%v", pName)

		body = append(
			body,
			jen.Var().Id(convName).Op("*").Qual("C", shape.typeName),
			jen.If(jen.Id(pName).Op("!=").Id("nil")).Block(
				jen.Id(convName).Op("=").Params(jen.Op("*").Qual("C", shape.typeName)).Params(jen.Qual("unsafe", "Pointer").Params(jen.Id(pName))),
			),
		)

		args = append(args, jen.Id(convName))

	case pointerArgShapeUnknownIdent:
		params = append(params, jen.Id(pName).Op("*").Id(shape.typeName))

		convName := fmt.Sprintf("tmp%v", pName)

		body = append(
			body,
			jen.Var().Id(convName).Op("*").Qual("C", shape.typeName),
			jen.If(jen.Id(pName).Op("!=").Id("nil")).Block(
				jen.Id(convName).Op("=").Params(jen.Op("*").Qual("C", shape.typeName)).Params(jen.Qual("unsafe", "Pointer").Params(jen.Id(pName))),
			),
		)

		args = append(args, jen.Id(convName))

	case pointerArgShapePointerToStruct:
		params = append(params, jen.Id(pName).Op("**").Id(shape.typeName))

		ptrName := fmt.Sprintf("ptr%v", pName)
		tmpName := fmt.Sprintf("tmp%v", pName)
		oldName := fmt.Sprintf("oldTmp%v", pName)
		innerName := fmt.Sprintf("inner%v", pName)

		body = append(
			body,
			jen.Var().Id(ptrName).Op("**").Qual("C", shape.st.CName()),
			jen.Var().Id(tmpName).Op("*").Qual("C", shape.st.CName()),
			jen.Var().Id(oldName).Op("*").Qual("C", shape.st.CName()),
			jen.If(jen.Id(pName).Op("!=").Id("nil")).Block(
				jen.Id(innerName).Op(":=").Op("*").Id(pName),
				jen.If(jen.Id(innerName).Op("!=").Id("nil")).Block(
					jen.Id(tmpName).Op("=").Id(innerName).Dot("ptr"),
					jen.Id(oldName).Op("=").Id(tmpName),
				),
				jen.Id(ptrName).Op("=").Op("&").Id(tmpName),
			),
		)

		postCall = append(
			postCall,
			jen.If(jen.Id(tmpName).Op("!=").Id(oldName).Op("&&").Id(pName).Op("!=").Id("nil")).Block(

				jen.If(jen.Id(tmpName).Op("!=").Id("nil")).Block(
					jen.Op("*").Id(pName).Op("=").Op("&").Id(shape.typeName).Values(jen.Dict{
						jen.Id("ptr"): jen.Id(tmpName),
					}),
				).Else().Block(
					jen.Op("*").Id(pName).Op("=").Id("nil"),
				),
			),
		)

		args = append(args, jen.Id(ptrName))
	}

	return params, args, body, postCall, false
}

func convParamName(val string) string {
	val = strcase.ToLowerCamel(val)

	if val == "type" || val == "range" {
		val = fmt.Sprintf("_%v", val)
	}

	return val
}

var acronyms = []string{
	"av", "hw", "lib", "ff", "io", "api",
}

func (g *Generator) convPart(val string) string {
	val = strings.ToLower(val)

	removed := true

	var prefixes []string

	for removed {
		removed = false

		for _, acronym := range acronyms {
			if after, ok := strings.CutPrefix(val, acronym); ok {
				val = after
				prefixes = append(prefixes, acronym)

				removed = true
			}
		}
	}

	if len(val) > 0 {
		a := val[0:1]
		b := val[1:]
		val = strings.ToUpper(a) + b
	}

	for _, prefix := range slices.Backward(prefixes) {
		val = strings.ToUpper(prefix) + val
	}

	return val
}

var digitRegex = regexp.MustCompile(`(\d)`)

func (g *Generator) convCamel(val string) string {
	divs := digitRegex.ReplaceAllString(val, "${1}_")
	parts := strings.Split(divs, "_")

	var newParts []string

	for _, part := range parts {
		newParts = append(newParts, g.convPart(part))
	}

	res := strings.Join(newParts, "")

	// Go has a single package-level identifier namespace, so a generated name
	// (typically a function) that camel-cases onto an already-registered struct
	// name must be disambiguated. Appending "_" keeps it exported and stable.
	// structs is fully populated by Parse() before any conversion runs, so this
	// lookup is complete and deterministic. Real collisions today:
	// avfilter_link -> AVFilterLink_, av_buffer_ref -> AVBufferRef_.
	if _, ok := g.input.structs[res]; ok {
		res += "_"
	}

	return res
}
