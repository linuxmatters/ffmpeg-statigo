package main

import (
	"fmt"
	"log"
	"reflect"
	"slices"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
)

func (g *Generator) generateStructs() {
	i := g.input

	o := newFile()

	for _, stName := range i.structOrder {
		st := i.structs[stName]

		// Skip structs that require manual bindings
		if skippedStructs[st.Name] {
			log.Println("Skipping struct", st.Name, "(manual binding in streamgroup.go)")
			continue
		}

		log.Println("Generating struct", st.Name)
		o.Commentf("--- Struct %v ---", st.Name)
		o.Line()

		goName := st.Name

		o.Commentf("%v wraps %v.", goName, st.Name)

		if st.Comment != "" {
			o.Comment(st.Comment)
		}

		cName := st.CName()

		if st.ByValue {
			o.Type().Id(goName).Struct(
				jen.Id("value").Qual("C", cName),
			)

			o.Line()

			// Generate ToXArray helper for ByValue structs
			o.Func().
				Id(fmt.Sprintf("To%vArray", goName)).
				Params(jen.Id("ptr").Qual("unsafe", "Pointer")).
				Op("*").Id("Array").Types(jen.Op("*").Id(goName)).
				Block(
					jen.If(jen.Id("ptr").Op("==").Id("nil")).Block(
						jen.Return(jen.Id("nil")),
					),
					jen.Line(),
					jen.Return(
						jen.Op("&").Id("Array").Types(jen.Op("*").Id(goName)).Values(jen.Dict{
							jen.Id("ptr"):      jen.Id("ptr"),
							jen.Id("elemSize"): jen.Qual("C", fmt.Sprintf("sizeof_%v", cName)),

							jen.Id("loadPtr"): jen.Func().
								Params(jen.Id("pointer").Qual("unsafe", "Pointer")).
								Op("*").Id(goName).
								Block(
									jen.Id("cValue").Op(":=").Parens(jen.Op("*").Qual("C", cName)).Parens(jen.Id("pointer")),
									jen.Return(
										jen.Op("&").Id(goName).Values(jen.Dict{
											jen.Id("value"): jen.Op("*").Id("cValue"),
										}),
									),
								),

							jen.Id("storePtr"): jen.Func().
								Params(
									jen.Id("pointer").Qual("unsafe", "Pointer"),
									jen.Id("value").Op("*").Id(goName),
								).
								Block(
									jen.Id("dest").Op(":=").Parens(jen.Op("*").Qual("C", cName)).Parens(jen.Id("pointer")),
									jen.Op("*").Id("dest").Op("=").Id("value").Dot("value"),
								),
						}),
					),
				)

			o.Line()

		} else {
			o.Type().Id(goName).Struct(
				jen.Id("ptr").Op("*").Qual("C", cName),
			)

			o.Func().
				Params(jen.Id("s").Op("*").Id(goName)).
				Id("RawPtr").
				Params().
				Qual("unsafe", "Pointer").
				Block(jen.Return(jen.Qual("unsafe", "Pointer").Params(jen.Id("s").Dot("ptr"))))

			o.Line()

			o.Func().
				Id(fmt.Sprintf("To%vArray", goName)).
				Params(jen.Id("ptr").Qual("unsafe", "Pointer")).
				Op("*").Id("Array").Types(jen.Op("*").Id(goName)).
				Block(
					jen.If(jen.Id("ptr").Op("==").Id("nil")).Block(
						jen.Return(jen.Id("nil")),
					),
					jen.Line(),
					jen.Return(
						jen.Op("&").Id("Array").Types(jen.Op("*").Id(goName)).Values(jen.Dict{
							jen.Id("ptr"):      jen.Id("ptr"),
							jen.Id("elemSize"): jen.Id("ptrSize"),

							jen.Id("loadPtr"): jen.Func().
								Params(jen.Id("pointer").Qual("unsafe", "Pointer")).
								Op("*").Id(goName).
								Block(
									jen.Id("ptr").Op(":=").Parens(jen.Op("**").Qual("C", cName)).Parens(jen.Id("pointer")),
									jen.Id("value").Op(":=").Op("*").Id("ptr"),

									jen.Var().Id("valueMapped").Op("*").Id(goName),
									jen.If(jen.Id("value").Op("!=").Id("nil")).Block(
										jen.Id("valueMapped").Op("=").Op("&").Id(goName).Values(jen.Dict{
											jen.Id("ptr"): jen.Id("value"),
										}),
									),
									jen.Return(jen.Id("valueMapped")),
								),

							jen.Id("storePtr"): jen.Func().
								Params(
									jen.Id("pointer").Qual("unsafe", "Pointer"),
									jen.Id("value").Op("*").Id(goName),
								).
								Block(
									jen.Id("ptr").Op(":=").Parens(jen.Op("**").Qual("C", cName)).Parens(jen.Id("pointer")),

									jen.If(jen.Id("value").Op("!=").Id("nil")).Block(
										jen.Op("*").Id("ptr").Op("=").Id("value").Dot("ptr"),
									).Else().Block(
										jen.Op("*").Id("ptr").Op("=").Id("nil"),
									),
								),
						}),
					),
				)

			o.Line()

		}

		for _, field := range st.Fields {
			g.marshalField(o, st, field)
		}
	}

	err := saveFormatted(o, "structs.gen.go")
	if err != nil {
		log.Panicln(err)
	}
}

func (g *Generator) marshalField(o *jen.File, st *Struct, field *Field) {
	goName := st.Name

	// Check if this field should be skipped (manual binding in streamgroup.go)
	if fields, ok := skippedFields[st.Name]; ok {
		if fields[field.Name] {
			o.Commentf("%v skipped (manual binding in streamgroup.go)", field.Name)
			o.Line()
			return
		}
	}

	// Fields the generator cannot emit but that have hand-written accessors in
	// fields.go. Record the skip (keeping the ceiling count) but annotate the
	// reason so the summary stops counting them as coverage gaps.
	if fields, ok := manuallyWrappedFields[st.Name]; ok {
		if fields[field.Name] {
			const reason = "manually wrapped in fields.go"
			o.Commentf("%v skipped due to %v", field.Name, reason)
			o.Line()
			g.skips.Record(goName+"."+field.Name, reason)
			return
		}
	}

	fName := strcase.ToCamel(field.Name)

	cName := field.Name
	if cName == "type" || cName == "range" {
		cName = fmt.Sprintf("_%v", cName)
	}

	var (
		getBody    []jen.Code
		getRetType []jen.Code
		getParams  []jen.Code
		setBody    []jen.Code
		setParams  []jen.Code
		tgt        *jen.Statement

		refField bool
	)

	if field.BitWidth != -1 {
		o.Commentf("%v skipped due to bitfield", field.Name)
		o.Line()
		g.skips.Record(goName+"."+field.Name, "bitfield")

		return
	}

	if st.ByValue {
		tgt = jen.Id("s").Dot("value").Dot(cName)
	} else {
		tgt = jen.Id("s").Dot("ptr").Dot(cName)
	}

	switch v := field.Type.(type) {

	case *IdentType:

		if m, ok := primTypes[v.Name]; ok {
			getRetType = []jen.Code{jen.Id(m)}
			setParams = append(setParams, jen.Id("value").Id(m))

			getBody = append(getBody, jen.Return(jen.Id(m).Params(jen.Id("value"))))

			if v.IsAnonEnum {
				setBody = append(
					setBody,
					tgt.Op("=").Id("value"),
				)
			} else {
				// Use the original C type name if available, otherwise try override, otherwise map from Go type
				cType := field.CTypeName
				// Normalize C type names that have spaces (not valid after "C.")
				cType = strings.ReplaceAll(cType, "unsigned long", "ulong")
				cType = strings.ReplaceAll(cType, "unsigned int", "uint")
				cType = strings.ReplaceAll(cType, "unsigned char", "uchar")
				cType = strings.ReplaceAll(cType, "unsigned short", "ushort")

				if cType == "" || cType == "int" || cType == "uint" {
					// Check for struct-specific override
					if overrides, ok := fieldCTypeOverrides[st.Name]; ok {
						if override, ok := overrides[field.Name]; ok {
							cType = override
						}
					}
				}
				if cType == "" || cType == "int" || cType == "uint" {
					cType = getCType(v.Name, m)
				}
				setBody = append(
					setBody,
					tgt.Op("=").Params(jen.Qual("C", cType)).Params(jen.Id("value")),
				)
			}
		} else if s, ok := g.input.structs[v.Name]; ok {
			if s.ByValue {
				getRetType = []jen.Code{
					jen.Op("*").Id(s.Name),
				}
				setParams = append(setParams, jen.Id("value").Op("*").Id(s.Name))

				getBody = append(
					getBody,
					jen.Return(jen.Op("&").Id(s.Name).Values(jen.Dict{
						jen.Id("value"): jen.Id("value"),
					})),
				)

				setBody = append(
					setBody,
					tgt.Op("=").Id("value").Dot("value"),
				)
			} else {
				refField = true

				getRetType = []jen.Code{
					jen.Op("*").Id(s.Name),
				}

				getBody = append(
					getBody,
					jen.Return(jen.Op("&").Id(s.Name).Values(jen.Dict{
						jen.Id("ptr"): jen.Id("value"),
					})),
				)
			}
		} else if _, ok := g.input.callbacks[v.Name]; ok {
			o.Commentf("%v skipped due to ident callback", field.Name)
			o.Line()
			g.skips.Record(goName+"."+field.Name, "ident callback")

			return
		} else if e, ok := g.input.enums[v.Name]; ok {
			getRetType = []jen.Code{jen.Id(v.Name)}
			setParams = append(setParams, jen.Id("value").Id(v.Name))

			getBody = append(getBody, jen.Return(jen.Id(v.Name).Params(jen.Id("value"))))

			setBody = append(
				setBody,
				tgt.Op("=").Params(jen.Qual("C", e.CName())).Params(jen.Id("value")),
			)
		} else {
			// Unknown IdentType - might be a typedef alias or external type
			// Try to use it as-is, treating it like a struct passed by reference
			structName := st.Name
			fieldName := field.Name
			typeName := v.Name

			log.Printf("unexpected IdentType in struct field - struct: %v, field: %v, type: %v (treating as struct reference)\n", structName, fieldName, typeName)

			// Treat as a struct passed by reference (like typedef aliases)
			refField = true

			getRetType = []jen.Code{
				jen.Op("*").Id(typeName),
			}

			getBody = append(
				getBody,
				jen.Return(jen.Op("&").Id(typeName).Values(jen.Dict{
					jen.Id("ptr"): jen.Id("value"),
				})),
			)
		}

	case *PointerType:
		var pSkip bool
		getRetType, setParams, getBody, setBody, pSkip = g.marshalPointerField(o, st, field, v, tgt)
		if pSkip {
			return
		}

	case *ConstArray:
		refField = true

		switch iv := v.Inner.(type) {
		case *IdentType:
			if pt, ok := primTypes[iv.Name]; ok {
				// Primitive type array (e.g., uint8_t[64])
				getRetType = []jen.Code{
					jen.Op("*").Id("Array").Types(jen.Id(pt)),
				}

				goName := strcase.ToCamel(pt)

				getBody = append(
					getBody,
					jen.Return(jen.Id(fmt.Sprintf("To%vArray", goName)).Params(
						jen.Qual("unsafe", "Pointer").Params(jen.Id("value")),
					)),
				)
			} else if en, ok := g.input.enums[iv.Name]; ok {
				// Enum type array (e.g., AVDOVIMappingMethod[AV_DOVI_MAX_PIECES])
				getRetType = []jen.Code{
					jen.Op("*").Id("Array").Types(jen.Id(en.Name)),
				}

				getBody = append(
					getBody,
					jen.Return(jen.Id(fmt.Sprintf("To%vArray", en.Name)).Params(
						jen.Qual("unsafe", "Pointer").Params(jen.Id("value")),
					)),
				)
			} else if st, ok := g.input.structs[iv.Name]; ok {
				// Struct type array (e.g., AVDOVIReshapingCurve[3])
				if st.ByValue {
					// Array of struct values
					getRetType = []jen.Code{
						jen.Op("*").Id("Array").Types(jen.Op("*").Id(st.Name)),
					}

					getBody = append(
						getBody,
						jen.Return(jen.Id(fmt.Sprintf("To%vArray", st.Name)).Params(
							jen.Qual("unsafe", "Pointer").Params(jen.Id("value")),
						)),
					)
				} else {
					// Array of struct pointers - skip for now (complex)
					o.Commentf("%v skipped due to const array of struct pointers", field.Name)
					o.Line()
					g.skips.Record(goName+"."+field.Name, "const array of struct pointers")

					return
				}
			} else {
				// Unknown type - might be typedef or forward declaration
				o.Commentf("%v skipped due to unknown const array type %v", field.Name, iv.Name)
				o.Line()
				g.skips.Record(goName+"."+field.Name, fmt.Sprintf("unknown const array type %v", iv.Name))

				return
			}

		case *PointerType:

			switch iiv := iv.Inner.(type) {
			case *IdentType:
				if iiv.Name == "uint8_t" {
					getRetType = []jen.Code{
						jen.Op("*").Id("Array").Types(jen.Qual("unsafe", "Pointer")),
					}

					getBody = append(
						getBody,
						jen.Return(jen.Id("ToUint8PtrArray").Params(
							jen.Qual("unsafe", "Pointer").Params(jen.Id("value")),
						)),
					)
				} else if st, ok := g.input.structs[iiv.Name]; ok {
					if st.ByValue {
						o.Commentf("%v skipped due to value ident ptr const array", field.Name)
						o.Line()
						g.skips.Record(goName+"."+field.Name, "value ident ptr const array")

						return
					}

					getRetType = []jen.Code{
						jen.Op("*").Id("Array").Types(jen.Op("*").Id(iiv.Name)),
					}

					getBody = append(
						getBody,
						jen.Return(jen.Id(fmt.Sprintf("To%vArray", iiv.Name)).Params(
							jen.Qual("unsafe", "Pointer").Params(jen.Id("value")),
						)),
					)
				} else {
					o.Commentf("%v skipped due to ident pointer const array", field.Name)
					o.Line()
					g.skips.Record(goName+"."+field.Name, "ident pointer const array")

					return
				}

			default:
				reason := fmt.Sprintf("unhandled const array pointer inner type %v", reflect.TypeOf(iv.Inner))
				o.Commentf("%v skipped due to %v", field.Name, reason)
				o.Line()
				g.skips.Record(goName+"."+field.Name, reason)

				return
			}

		case *ConstArray:
			o.Commentf("%v skipped due to multi dim const array", field.Name)
			o.Line()
			g.skips.Record(goName+"."+field.Name, "multi dim const array")

			return

		default:
			reason := fmt.Sprintf("unhandled const array inner type %v", reflect.TypeOf(v.Inner))
			o.Commentf("%v skipped due to %v", field.Name, reason)
			o.Line()
			g.skips.Record(goName+"."+field.Name, reason)

			return
		}

	case *UnionType:
		o.Commentf("%v skipped due to union type", field.Name)
		o.Line()
		g.skips.Record(goName+"."+field.Name, "union type")

		return

	default:
		reason := fmt.Sprintf("unhandled field type %v", reflect.TypeOf(field.Type))
		o.Commentf("%v skipped due to %v", field.Name, reason)
		o.Line()
		g.skips.Record(goName+"."+field.Name, reason)

		return
	}

	switch {
	case refField && st.ByValue:
		o.Commentf("%v skipped due to ref field of value struct", field.Name)
		o.Line()
		g.skips.Record(goName+"."+field.Name, "ref field of value struct")

		return
	case refField:
		getBody = slices.Insert(
			getBody, 0,
			jen.Code(jen.Id("value").Op(":=").Op("&").Id("s").Dot("ptr").Dot(cName)),
		)
	case st.ByValue:
		getBody = slices.Insert(
			getBody, 0,
			jen.Code(jen.Id("value").Op(":=").Id("s").Dot("value").Dot(cName)),
		)
	default:
		getBody = slices.Insert(
			getBody, 0,
			jen.Code(jen.Id("value").Op(":=").Id("s").Dot("ptr").Dot(cName)),
		)
	}

	o.Commentf("%v gets the %v field.", fName, field.Name)

	if field.Comment != "" {
		o.Comment(field.Comment)
	}

	o.Func().
		Params(jen.Id("s").Op("*").Id(goName)).
		Id(fName).
		Params(getParams...).
		Add(getRetType...).
		Block(getBody...)

	o.Line()

	if len(setBody) > 0 {
		o.Commentf("Set%v sets the %v field.", fName, field.Name)

		if field.Comment != "" {
			o.Comment(field.Comment)
		}

		o.Func().
			Params(jen.Id("s").Op("*").Id(goName)).
			Id(fmt.Sprintf("Set%v", fName)).
			Params(setParams...).
			Block(setBody...)

		o.Line()
	}
}

func (g *Generator) marshalPointerField(o *jen.File, st *Struct, field *Field, v *PointerType, tgt *jen.Statement) (getRetType, setParams, getBody, setBody []jen.Code, skip bool) {
	// Capture before the `st, ok := g.input.structs[...]` branches below
	// shadow the outer struct binding. The skip collector needs the
	// containing struct name, not the array-element struct name.
	parentName := st.Name

	switch iv := v.Inner.(type) {
	case nil:
		getRetType = []jen.Code{
			jen.Qual("unsafe", "Pointer"),
		}
		setParams = append(setParams, jen.Id("value").Qual("unsafe", "Pointer"))
		getBody = append(getBody, jen.Return(jen.Id("value")))
		setBody = append(setBody, tgt.Op("=").Id("value"))

	case *IdentType:

		if iv.Name == "URLContext" || iv.Name == "AVFilterCommand" || iv.Name == "AVCodecInternal" {
			o.Commentf("%v skipped due to ptr to ignored type", field.Name)
			o.Line()
			g.skips.Record(parentName+"."+field.Name, "ptr to ignored type")

			return getRetType, setParams, getBody, setBody, true
		} else if iv.Name == "char" {
			getRetType = []jen.Code{
				jen.Op("*").Id("CStr"),
			}
			setParams = append(setParams, jen.Id("value").Op("*").Id("CStr"))
			getBody = append(getBody, jen.Return(jen.Id("wrapCStr").Params(jen.Id("value"))))
			setBody = append(setBody, tgt.Op("=").Id("value").Dot("ptr"))

		} else if iv.Name == "uint8_t" {
			getRetType = []jen.Code{
				jen.Qual("unsafe", "Pointer"),
			}
			getBody = append(getBody, jen.Return(jen.Qual("unsafe", "Pointer").Params(jen.Id("value"))))

			setParams = append(setParams, jen.Id("value").Qual("unsafe", "Pointer"))
			setBody = append(setBody, tgt.Op("=").Params(jen.Op("*").Qual("C", iv.Name)).Params(jen.Id("value")))
		} else if _, ok := primTypes[iv.Name]; ok {
			o.Commentf("%v skipped due to prim ptr", field.Name)
			o.Line()
			g.skips.Record(parentName+"."+field.Name, "prim ptr")

			return getRetType, setParams, getBody, setBody, true
		} else if enum, ok := g.input.enums[iv.Name]; ok {
			getRetType = []jen.Code{
				jen.Op("*").Id("Array").Types(jen.Id(iv.Name)),
			}

			getBody = append(
				getBody,
				jen.Return(jen.Id(fmt.Sprintf("To%vArray", iv.Name)).Params(
					jen.Qual("unsafe", "Pointer").Params(jen.Id("value")),
				)),
			)

			setParams = append(setParams, jen.Id("value").Op("*").Id("Array").Types(jen.Id(iv.Name)))

			cName := enum.CName()

			setBody = append(
				setBody,
				jen.If(jen.Id("value").Op("!=").Id("nil")).Block(
					tgt.Clone().Op("=").Params(jen.Op("*").Qual("C", cName)).Params(jen.Id("value").Dot("ptr")),
				).Else().Block(
					tgt.Clone().Op("=").Id("nil"),
				),
			)

		} else if _, ok := g.input.callbacks[iv.Name]; ok {
			o.Commentf("%v skipped due to callback ptr", field.Name)
			o.Line()
			g.skips.Record(parentName+"."+field.Name, "callback ptr")

			return getRetType, setParams, getBody, setBody, true
		} else if ist, ok := g.input.structs[iv.Name]; ok {
			if ist.ByValue {
				o.Commentf("%v skipped due to struct value ptr", field.Name)
				o.Line()
				g.skips.Record(parentName+"."+field.Name, "struct value ptr")

				return getRetType, setParams, getBody, setBody, true
			}

			getRetType = []jen.Code{
				jen.Op("*").Id(iv.Name),
			}
			setParams = append(setParams, jen.Id("value").Op("*").Id(iv.Name))

			getBody = append(
				getBody,
				jen.Var().Id("valueMapped").Op("*").Id(iv.Name),
				jen.If(jen.Id("value").Op("!=").Id("nil")).Block(
					jen.Id("valueMapped").Op("=").Op("&").Id(iv.Name).Values(jen.Dict{
						jen.Id("ptr"): jen.Id("value"),
					}),
				),
				jen.Return(jen.Id("valueMapped")),
			)

			setBody = append(
				setBody,
				jen.If(jen.Id("value").Op("!=").Id("nil")).Block(
					tgt.Clone().Op("=").Id("value").Dot("ptr"),
				).Else().Block(
					tgt.Clone().Op("=").Id("nil"),
				),
			)
		} else {
			structName := st.Name
			fieldName := field.Name
			typeName := iv.Name

			log.Printf("unexpected IdentType in struct setter - struct: %v, field: %v, type: %v\n", structName, fieldName, typeName)
			o.Commentf("%v skipped due to unexpected IdentType %v", fieldName, typeName)
			o.Line()
			g.skips.Record(parentName+"."+field.Name, fmt.Sprintf("unexpected IdentType %v", typeName))
			return getRetType, setParams, getBody, setBody, true
		}

	case *FuncType:
		o.Commentf("%v skipped due to func ptr", field.Name)
		o.Line()
		g.skips.Record(parentName+"."+field.Name, "func ptr")

		return getRetType, setParams, getBody, setBody, true

	case *PointerType:

		switch iiv := iv.Inner.(type) {
		case *IdentType:
			if st, ok := g.input.structs[iiv.Name]; ok {

				getRetType = []jen.Code{
					jen.Op("*").Id("Array").Types(jen.Op("*").Id(iiv.Name)),
				}

				getBody = append(
					getBody,
					jen.Return(jen.Id(fmt.Sprintf("To%vArray", iiv.Name)).Params(
						jen.Qual("unsafe", "Pointer").Params(jen.Id("value")),
					)),
				)

				setParams = append(setParams, jen.Id("value").Op("*").Id("Array").Types(jen.Id(iiv.Name)))

				cName := st.CName()

				setBody = append(
					setBody,
					jen.If(jen.Id("value").Op("!=").Id("nil")).Block(
						tgt.Clone().Op("=").Params(jen.Op("**").Qual("C", cName)).Params(jen.Id("value").Dot("ptr")),
					).Else().Block(
						tgt.Clone().Op("=").Id("nil"),
					),
				)
			} else {
				//nolint:dupword // "ptr ptr" describes a pointer-to-pointer field; changing it would alter generated output
				o.Commentf("%v skipped due to unknown ptr ptr", field.Name)
				o.Line()
				//nolint:dupword // "ptr ptr" describes a pointer-to-pointer field; the skip reason must match the generated comment
				g.skips.Record(parentName+"."+field.Name, "unknown ptr ptr")

				return getRetType, setParams, getBody, setBody, true
			}

		default:
			reason := fmt.Sprintf("unhandled pointer-pointer inner type %v", reflect.TypeOf(iv.Inner))
			o.Commentf("%v skipped due to %v", field.Name, reason)
			o.Line()
			g.skips.Record(parentName+"."+field.Name, reason)

			return getRetType, setParams, getBody, setBody, true
		}

	default:
		reason := fmt.Sprintf("unhandled pointer inner type %v", reflect.TypeOf(v.Inner))
		o.Commentf("%v skipped due to %v", field.Name, reason)
		o.Line()
		g.skips.Record(parentName+"."+field.Name, reason)

		return getRetType, setParams, getBody, setBody, true
	}

	return getRetType, setParams, getBody, setBody, false
}
