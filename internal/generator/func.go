package main

import (
	"fmt"
	"log"
	"reflect"
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
		// These are manually wrapped in custom.go with proper pointer handling
		// Pinned by TestUUIDBindingsNotDuplicated in bindings_test.go (asserts the seven manual wrappers are not also generated).
		if fn.Name == "av_uuid_parse" || fn.Name == "av_uuid_urn_parse" || fn.Name == "av_uuid_parse_range" ||
			fn.Name == "av_uuid_unparse" || fn.Name == "av_uuid_equal" || fn.Name == "av_uuid_copy" || fn.Name == "av_uuid_nil" {
			o.Commentf("%v skipped due to array typedef (manually wrapped in custom.go)", fn.Name)
			o.Line()
			g.skips.Record(fn.Name, "array typedef (manually wrapped in custom.go)")

			continue outer
		}

		if typeEquals(fn.Result, fileType) || typeEquals(fn.Result, fileType2) {
			o.Commentf("%v skipped due to return", fn.Name)
			o.Line()
			g.skips.Record(fn.Name, "return")

			continue outer
		}

		for _, arg := range fn.Args {
			skip := false

			if typeEquals(arg.Type, fileType) || typeEquals(arg.Type, fileType2) || typeEquals(arg.Type, vaListType) || typeEquals(arg.Type, tmType) {
				skip = true
			}

			if v, ok := arg.Type.(*PointerType); ok {
				switch iv := v.Inner.(type) {
				case *FuncType:
					skip = true
				case *IdentType:
					// Skip pointer to tm (C standard library type)
					if iv.Name == "tm" {
						skip = true
					}
				}
			}

			if skip {
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
	switch v := result.(type) {
	case nil:
		// nothing
		body = append(body, cc)
		body = append(body, postCall...)

	case *IdentType:

		body = append(body, jen.Id("ret").Op(":=").Add(cc))
		body = append(body, postCall...)

		if v.Name == "int" {
			retType = []jen.Code{jen.Params(jen.Id("int"), jen.Id("error"))}
			body = append(
				body,
				jen.Return(
					jen.Id("int").Params(jen.Id("ret")).Op(",").
						Id("WrapErr").Params(jen.Id("int").Params(jen.Id("ret"))),
				),
			)
		} else if m, ok := primTypes[v.Name]; ok {
			retType = []jen.Code{jen.Id(m)}
			body = append(body, jen.Return(jen.Id(m).Params(jen.Id("ret"))))
		} else if s, ok := g.input.structs[v.Name]; ok {
			if s.ByValue {
				retType = []jen.Code{jen.Op("*").Id(v.Name)}
				body = append(
					body,
					jen.Return(jen.Op("&").Id(v.Name).Values(jen.Dict{
						jen.Id("value"): jen.Id("ret"),
					})),
				)
			} else {
				o.Commentf("%v skipped due to return", fn.Name)
				o.Line()
				g.skips.Record(fn.Name, "return")
				return nil, nil, true
			}
		} else if _, ok := g.input.callbacks[v.Name]; ok {
			// Callback type - convert C name to Go name
			goTypeName := g.convCamel(v.Name)
			retType = []jen.Code{jen.Id(goTypeName)}
			body = append(body, jen.Return(jen.Id(goTypeName).Params(jen.Id("ret"))))
		} else {
			retType = []jen.Code{jen.Id(v.Name)}
			body = append(body, jen.Return(jen.Id(v.Name).Params(jen.Id("ret"))))
		}

	case *PointerType:
		body = append(
			body,
			jen.Id("ret").Op(":=").Add(cc),
		)
		body = append(body, postCall...)

		switch iv := v.Inner.(type) {
		case nil:
			retType = []jen.Code{
				jen.Qual("unsafe", "Pointer"),
			}
			body = append(body, jen.Return(jen.Id("ret")))

		case *IdentType:

			if iv.Name == "char" {
				retType = []jen.Code{
					jen.Op("*").Id("CStr"),
				}
				body = append(body, jen.Return(jen.Id("wrapCStr").Params(jen.Id("ret"))))
			} else if iv.Name == "uint8_t" {
				retType = []jen.Code{
					jen.Qual("unsafe", "Pointer"),
				}
				body = append(
					body,
					jen.Return(jen.Qual("unsafe", "Pointer").Params(jen.Id("ret"))),
				)
			} else if m, ok := primTypes[iv.Name]; ok {
				// Handle pointer to primitive type (e.g., *int, *float64)
				retType = []jen.Code{
					jen.Op("*").Id(m),
				}
				body = append(
					body,
					jen.Return(jen.Params(jen.Op("*").Id(m)).Params(
						jen.Qual("unsafe", "Pointer").Params(jen.Id("ret")),
					)),
				)
			} else if _, ok := g.input.structs[iv.Name]; ok {
				retType = []jen.Code{
					jen.Op("*").Id(iv.Name),
				}

				body = append(
					body,
					jen.Var().Id("retMapped").Op("*").Id(iv.Name),
					jen.If(jen.Id("ret").Op("!=").Id("nil")).Block(
						jen.Id("retMapped").Op("=").Op("&").Id(iv.Name).Values(jen.Dict{
							jen.Id("ptr"): jen.Id("ret"),
						}),
					),
					jen.Return(jen.Id("retMapped")),
				)
			} else {
				// Unknown type - could be a typedef alias or enum
				// Cast through unsafe.Pointer to handle typedef aliases correctly
				retType = []jen.Code{
					jen.Op("*").Id(iv.Name),
				}
				body = append(
					body,
					jen.Return(jen.Params(jen.Op("*").Id(iv.Name)).Params(
						jen.Qual("unsafe", "Pointer").Params(jen.Id("ret")),
					)),
				)
			}

		case *PointerType:
			// Pointer to pointer return type (e.g., AVFrameSideData *const *)
			// Skip for now - these are complex array returns that need special handling
			o.Commentf("%v skipped due to pointer-to-pointer return type", fn.Name)
			o.Line()
			g.skips.Record(fn.Name, "pointer-to-pointer return type")
			return nil, nil, true

		default:
			reason := fmt.Sprintf("unhandled return pointer inner type %v", reflect.TypeOf(v.Inner))
			o.Commentf("%v skipped due to %v", fn.Name, reason)
			o.Line()
			g.skips.Record(fn.Name, reason)
			return nil, nil, true
		}

	default:
		reason := fmt.Sprintf("unhandled return type %v", reflect.TypeOf(fn.Result))
		o.Commentf("%v skipped due to %v", fn.Name, reason)
		o.Line()
		g.skips.Record(fn.Name, reason)
		return nil, nil, true
	}

	return retType, body, false
}

func (g *Generator) marshalArg(o *jen.File, fn *Function, arg *Param) (params, args, body, postCall []jen.Code, skip bool) {
	pName := convParamName(arg.Name)

	// WORKAROUND: libclang on Linux incorrectly reports size_t as int for some parameters
	// Track actual type names for pointer-to-primitive cases
	actualTypeName := ""
	if ptrType, ok := arg.Type.(*PointerType); ok {
		if identType, ok := ptrType.Inner.(*IdentType); ok {
			actualTypeName = identType.Name
			// libclang misreports some size_t* output parameters as int*. The
			// affected (function, parameter) pairs are the sizeT entries in
			// outputParams rather than matched by substring, so the rewrite is
			// exact. Pinned by TestGeneratorOutputParameters in bindings_test.go
			// and by the byte-identical regen gate.
			if p, ok := outputParams[fn.Name][arg.Name]; ok && p.sizeT && actualTypeName == "int" {
				actualTypeName = "size_t"
			}
		}
	}

	switch v := arg.Type.(type) {
	case *IdentType:
		// WORKAROUND: libclang on Linux incorrectly reports size_t as int
		// Only fix specific known cases where FFmpeg headers use size_t
		typeName := v.Name
		if typeName == "int" && arg.Name == "buf_size" {
			// These functions use size_t buf_size per FFmpeg headers.
			// Pinned by TestGeneratorOutputParameters in bindings_test.go (size_t buf_size branch).
			if fn.Name == "av_channel_name" || fn.Name == "av_channel_description" ||
				fn.Name == "av_channel_layout_describe" {
				typeName = "size_t"
			}
		} else if typeName == "int" && arg.Name == "max_size" {
			// avio_read_to_bprint uses size_t max_size per FFmpeg headers.
			// Pinned by TestGeneratorOutputParameters in bindings_test.go (size_t max_size branch).
			if fn.Name == "avio_read_to_bprint" {
				typeName = "size_t"
			}
		}

		if m, ok := primTypes[typeName]; ok {
			params = append(params, jen.Id(pName).Id(m))
			cType := getCType(typeName, m)
			args = append(args, jen.Qual("C", cType).Params(jen.Id(pName)))
		} else if e, ok := g.input.enums[typeName]; ok {
			params = append(params, jen.Id(pName).Id(typeName))
			args = append(args, jen.Qual("C", e.CName()).Params(jen.Id(pName)))
		} else if s, ok := g.input.structs[typeName]; ok {
			if s.ByValue {
				params = append(params, jen.Id(pName).Op("*").Id(s.Name))
				args = append(args, jen.Id(pName).Dot("value"))
			} else {
				o.Commentf("%v skipped due to %v", fn.Name, pName)
				o.Line()
				g.skips.Record(fn.Name, fmt.Sprintf("%v", pName))

				return params, args, body, postCall, true
			}
		} else if _, ok := g.input.callbacks[typeName]; ok {
			// Callback type passed by value - CGO doesn't allow conversion from unsafe.Pointer to function pointer
			// Skip these for now
			o.Commentf("%v skipped due to %v (callback by value)", fn.Name, pName)
			o.Line()
			g.skips.Record(fn.Name, fmt.Sprintf("%v (callback by value)", pName))
			return params, args, body, postCall, true
		} else {
			params = append(params, jen.Id(pName).Id(typeName))
			args = append(args, jen.Qual("C", v.Name).Params(jen.Id(pName)))
		}

	case *PointerType:
		var pSkip bool
		params, args, body, postCall, pSkip = g.marshalPointerArg(o, fn, arg, v, pName, actualTypeName)
		if pSkip {
			return params, args, body, postCall, true
		}

	case *Array:
		o.Commentf("%v skipped due to %v", fn.Name, pName)
		o.Line()
		g.skips.Record(fn.Name, fmt.Sprintf("%v", pName))
		return params, args, body, postCall, true

	case *ConstArray:
		o.Commentf("%v skipped due to const array param %v", fn.Name, pName)
		o.Line()
		g.skips.Record(fn.Name, fmt.Sprintf("const array param %v", pName))
		return params, args, body, postCall, true

	default:
		reason := fmt.Sprintf("unhandled arg type %v (%v)", reflect.TypeOf(arg.Type), pName)
		o.Commentf("%v skipped due to %v", fn.Name, reason)
		o.Line()
		g.skips.Record(fn.Name, reason)
		return params, args, body, postCall, true
	}

	return params, args, body, postCall, false
}

func (g *Generator) marshalPointerArg(o *jen.File, fn *Function, arg *Param, v *PointerType, pName, actualTypeName string) (params, args, body, postCall []jen.Code, skip bool) {
	switch iv := v.Inner.(type) {
	case nil:
		params = append(params, jen.Id(pName).Qual("unsafe", "Pointer"))
		args = append(args, jen.Id(pName))

	case *IdentType:
		switch iv.Name {
		case "char":
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
		case "uint8_t", "uchar":
			params = append(params, jen.Id(pName).Qual("unsafe", "Pointer"))
			args = append(args, jen.Params(jen.Op("*").Qual("C", iv.Name)).Params(jen.Id(pName)))
		default:

			// Check if we have a corrected type name from the workaround
			typeNameForParam := iv.Name
			if actualTypeName != "" {
				typeNameForParam = actualTypeName
			}

			if m, ok := primTypes[typeNameForParam]; ok {
				// Pointer to primitive type - check if it's an output parameter
				if _, ok := outputParams[fn.Name][arg.Name]; ok {
					// This is likely an output parameter
					// We'll generate a wrapper function that handles the output
					params = append(params, jen.Id(pName).Op("*").Id(m))

					// For output parameters, we pass the Go pointer directly
					// and let CGO handle the conversion
					cType := getCType(typeNameForParam, m)
					args = append(args, jen.Params(jen.Op("*").Qual("C", cType)).Params(jen.Qual("unsafe", "Pointer").Params(jen.Id(pName))))
				} else {
					// Not an output parameter - skip as it's ambiguous
					o.Commentf("%v skipped due to %v (non-output primitive pointer)", fn.Name, pName)
					o.Line()
					g.skips.Record(fn.Name, fmt.Sprintf("%v (non-output primitive pointer)", pName))
					return params, args, body, postCall, true
				}
			} else if s, ok := g.input.structs[iv.Name]; ok {
				if s.ByValue {
					o.Commentf("%v skipped due to %v", fn.Name, pName)
					o.Line()
					g.skips.Record(fn.Name, fmt.Sprintf("%v", pName))
					return params, args, body, postCall, true
				}

				params = append(params, jen.Id(pName).Op("*").Id(iv.Name))

				convName := fmt.Sprintf("tmp%v", pName)

				body = append(
					body,
					jen.Var().Id(convName).Op("*").Qual("C", s.CName()),
					jen.If(jen.Id(pName).Op("!=").Id("nil")).Block(
						jen.Id(convName).Op("=").Id(pName).Dot("ptr"),
					),
				)

				args = append(args, jen.Id(convName))
			} else if e, ok := g.input.enums[iv.Name]; ok {
				// Pointer to enum type
				params = append(params, jen.Id(pName).Op("*").Id(iv.Name))

				convName := fmt.Sprintf("tmp%v", pName)

				body = append(
					body,
					jen.Var().Id(convName).Op("*").Qual("C", e.CName()),
					jen.If(jen.Id(pName).Op("!=").Id("nil")).Block(
						jen.Id(convName).Op("=").Params(jen.Op("*").Qual("C", e.CName())).Params(jen.Qual("unsafe", "Pointer").Params(jen.Id(pName))),
					),
				)

				args = append(args, jen.Id(convName))
			} else if _, ok := g.input.callbacks[iv.Name]; ok {
				// Pointer to callback type - use Go callback type name
				goTypeName := g.convCamel(iv.Name)
				params = append(params, jen.Id(pName).Op("*").Id(goTypeName))

				convName := fmt.Sprintf("tmp%v", pName)

				body = append(
					body,
					jen.Var().Id(convName).Op("*").Qual("C", iv.Name),
					jen.If(jen.Id(pName).Op("!=").Id("nil")).Block(
						jen.Id(convName).Op("=").Params(jen.Op("*").Qual("C", iv.Name)).Params(jen.Qual("unsafe", "Pointer").Params(jen.Id(pName))),
					),
				)

				args = append(args, jen.Id(convName))
			} else {
				// Unknown IdentType - could be a typedef alias defined in custom.go
				// Try to use it directly (e.g., AVCRC, AVAdler)
				// Cast through C type for the call
				params = append(params, jen.Id(pName).Op("*").Id(iv.Name))

				convName := fmt.Sprintf("tmp%v", pName)

				body = append(
					body,
					jen.Var().Id(convName).Op("*").Qual("C", iv.Name),
					jen.If(jen.Id(pName).Op("!=").Id("nil")).Block(
						jen.Id(convName).Op("=").Params(jen.Op("*").Qual("C", iv.Name)).Params(jen.Qual("unsafe", "Pointer").Params(jen.Id(pName))),
					),
				)

				args = append(args, jen.Id(convName))
			}
		}

	case *PointerType:

		switch iiv := iv.Inner.(type) {
		case *IdentType:

			if iiv.Name == "uint8_t" || iiv.Name == "char" {
				o.Commentf("%v skipped due to %v", fn.Name, pName)
				o.Line()
				g.skips.Record(fn.Name, fmt.Sprintf("%v", pName))
				return params, args, body, postCall, true
			}

			if _, ok := primTypes[iiv.Name]; ok {
				o.Commentf("%v skipped due to %v", fn.Name, pName)
				o.Line()
				g.skips.Record(fn.Name, fmt.Sprintf("%v", pName))
				return params, args, body, postCall, true
			} else if s, ok := g.input.structs[iiv.Name]; ok {
				params = append(params, jen.Id(pName).Op("**").Id(iiv.Name))

				ptrName := fmt.Sprintf("ptr%v", pName)
				tmpName := fmt.Sprintf("tmp%v", pName)
				oldName := fmt.Sprintf("oldTmp%v", pName)
				innerName := fmt.Sprintf("inner%v", pName)

				body = append(
					body,
					jen.Var().Id(ptrName).Op("**").Qual("C", s.CName()),
					jen.Var().Id(tmpName).Op("*").Qual("C", s.CName()),
					jen.Var().Id(oldName).Op("*").Qual("C", s.CName()),
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
							jen.Op("*").Id(pName).Op("=").Op("&").Id(iiv.Name).Values(jen.Dict{
								jen.Id("ptr"): jen.Id(tmpName),
							}),
						).Else().Block(
							jen.Op("*").Id(pName).Op("=").Id("nil"),
						),
					),
				)

				args = append(args, jen.Id(ptrName))
			} else {
				o.Commentf("%v skipped due to %v", fn.Name, pName)
				o.Line()
				g.skips.Record(fn.Name, fmt.Sprintf("%v", pName))
				return params, args, body, postCall, true
			}

		default:
			o.Commentf("%v skipped due to %v", fn.Name, pName)
			o.Line()
			g.skips.Record(fn.Name, fmt.Sprintf("%v", pName))
			return params, args, body, postCall, true

		}

	default:
		reason := fmt.Sprintf("unhandled pointer arg inner type %v (%v)", reflect.TypeOf(v.Inner), pName)
		o.Commentf("%v skipped due to %v", fn.Name, reason)
		o.Line()
		g.skips.Record(fn.Name, reason)
		return params, args, body, postCall, true
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
