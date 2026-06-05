package main

import (
	"fmt"
	"log"
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

	shape := g.classifyFieldShape(field)

	switch shape.kind {
	case fieldShapeIdentPrimitive:
		getRetType = []jen.Code{jen.Id(shape.goType)}
		setParams = append(setParams, jen.Id("value").Id(shape.goType))

		getBody = append(getBody, jen.Return(jen.Id(shape.goType).Params(jen.Id("value"))))

		if shape.anonEnum {
			setBody = append(
				setBody,
				tgt.Op("=").Id("value"),
			)
		} else {
			cType := normalizedFieldCType(st, field, shape.typeName, shape.goType)
			setBody = append(
				setBody,
				tgt.Op("=").Params(jen.Qual("C", cType)).Params(jen.Id("value")),
			)
		}

	case fieldShapeIdentByValueStruct:
		getRetType = []jen.Code{
			jen.Op("*").Id(shape.st.Name),
		}
		setParams = append(setParams, jen.Id("value").Op("*").Id(shape.st.Name))

		getBody = append(
			getBody,
			jen.Return(jen.Op("&").Id(shape.st.Name).Values(jen.Dict{
				jen.Id("value"): jen.Id("value"),
			})),
		)

		setBody = append(
			setBody,
			tgt.Op("=").Id("value").Dot("value"),
		)

	case fieldShapeIdentRefStruct:
		refField = true

		getRetType = []jen.Code{
			jen.Op("*").Id(shape.st.Name),
		}

		getBody = append(
			getBody,
			jen.Return(jen.Op("&").Id(shape.st.Name).Values(jen.Dict{
				jen.Id("ptr"): jen.Id("value"),
			})),
		)

	case fieldShapeIdentCallbackSkip:
		o.Commentf("%v skipped due to ident callback", field.Name)
		o.Line()
		g.skips.Record(goName+"."+field.Name, shape.reason)

		return

	case fieldShapeIdentEnum:
		getRetType = []jen.Code{jen.Id(shape.typeName)}
		setParams = append(setParams, jen.Id("value").Id(shape.typeName))

		getBody = append(getBody, jen.Return(jen.Id(shape.typeName).Params(jen.Id("value"))))

		setBody = append(
			setBody,
			tgt.Op("=").Params(jen.Qual("C", shape.enum.CName())).Params(jen.Id("value")),
		)

	case fieldShapeIdentUnknownRef:
		log.Printf("unexpected IdentType in struct field - struct: %v, field: %v, type: %v (treating as struct reference)\n", st.Name, field.Name, shape.typeName)

		refField = true

		getRetType = []jen.Code{
			jen.Op("*").Id(shape.typeName),
		}

		getBody = append(
			getBody,
			jen.Return(jen.Op("&").Id(shape.typeName).Values(jen.Dict{
				jen.Id("ptr"): jen.Id("value"),
			})),
		)

	case fieldShapePointer:
		var pSkip bool
		getRetType, setParams, getBody, setBody, pSkip = g.marshalPointerField(o, st, field, shape.ptr, tgt)
		if pSkip {
			return
		}

	case fieldShapeConstArray:
		refField = true
		var aSkip bool
		getRetType, getBody, aSkip = g.marshalConstArrayField(o, st, field, shape.array)
		if aSkip {
			return
		}

	case fieldShapeUnionSkip, fieldShapeUnhandledSkip:
		o.Commentf("%v skipped due to %v", field.Name, shape.reason)
		o.Line()
		g.skips.Record(goName+"."+field.Name, shape.reason)

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

func normalizedFieldCType(st *Struct, field *Field, typeName, goType string) string {
	cType := field.CTypeName
	cType = strings.ReplaceAll(cType, "unsigned long", "ulong")
	cType = strings.ReplaceAll(cType, "unsigned int", "uint")
	cType = strings.ReplaceAll(cType, "unsigned char", "uchar")
	cType = strings.ReplaceAll(cType, "unsigned short", "ushort")

	if cType == "" || cType == "int" || cType == "uint" {
		if overrides, ok := fieldCTypeOverrides[st.Name]; ok {
			if override, ok := overrides[field.Name]; ok {
				cType = override
			}
		}
	}
	if cType == "" || cType == "int" || cType == "uint" {
		cType = getCType(typeName, goType)
	}

	return cType
}

func (g *Generator) marshalConstArrayField(o *jen.File, st *Struct, field *Field, v *ConstArray) (getRetType, getBody []jen.Code, skip bool) {
	shape := g.classifyConstArrayFieldShape(v)

	switch shape.kind {
	case constArrayFieldShapePrimitive:
		getRetType = []jen.Code{
			jen.Op("*").Id("Array").Types(jen.Id(shape.goType)),
		}

		goName := strcase.ToCamel(shape.goType)

		getBody = append(
			getBody,
			jen.Return(jen.Id(fmt.Sprintf("To%vArray", goName)).Params(
				jen.Qual("unsafe", "Pointer").Params(jen.Id("value")),
			)),
		)

	case constArrayFieldShapeEnum:
		getRetType = []jen.Code{
			jen.Op("*").Id("Array").Types(jen.Id(shape.enum.Name)),
		}

		getBody = append(
			getBody,
			jen.Return(jen.Id(fmt.Sprintf("To%vArray", shape.enum.Name)).Params(
				jen.Qual("unsafe", "Pointer").Params(jen.Id("value")),
			)),
		)

	case constArrayFieldShapeByValueStruct:
		getRetType = []jen.Code{
			jen.Op("*").Id("Array").Types(jen.Op("*").Id(shape.st.Name)),
		}

		getBody = append(
			getBody,
			jen.Return(jen.Id(fmt.Sprintf("To%vArray", shape.st.Name)).Params(
				jen.Qual("unsafe", "Pointer").Params(jen.Id("value")),
			)),
		)

	case constArrayFieldShapeUint8PointerArray:
		getRetType = []jen.Code{
			jen.Op("*").Id("Array").Types(jen.Qual("unsafe", "Pointer")),
		}

		getBody = append(
			getBody,
			jen.Return(jen.Id("ToUint8PtrArray").Params(
				jen.Qual("unsafe", "Pointer").Params(jen.Id("value")),
			)),
		)

	case constArrayFieldShapeStructPointerArray:
		getRetType = []jen.Code{
			jen.Op("*").Id("Array").Types(jen.Op("*").Id(shape.typeName)),
		}

		getBody = append(
			getBody,
			jen.Return(jen.Id(fmt.Sprintf("To%vArray", shape.typeName)).Params(
				jen.Qual("unsafe", "Pointer").Params(jen.Id("value")),
			)),
		)

	case constArrayFieldShapeStructPointerSkip, constArrayFieldShapeValueIdentPointerSkip,
		constArrayFieldShapeIdentPointerSkip, constArrayFieldShapeMultiDimSkip,
		constArrayFieldShapeUnknownTypeSkip, constArrayFieldShapeUnhandledPointerInnerSkip,
		constArrayFieldShapeUnhandledInnerSkip:
		o.Commentf("%v skipped due to %v", field.Name, shape.reason)
		o.Line()
		g.skips.Record(st.Name+"."+field.Name, shape.reason)

		return getRetType, getBody, true
	}

	return getRetType, getBody, false
}

func (g *Generator) marshalPointerField(o *jen.File, st *Struct, field *Field, v *PointerType, tgt *jen.Statement) (getRetType, setParams, getBody, setBody []jen.Code, skip bool) {
	parentName := st.Name
	shape := g.classifyPointerFieldShape(v)

	switch shape.kind {
	case pointerFieldShapeUnsafePointer:
		getRetType = []jen.Code{
			jen.Qual("unsafe", "Pointer"),
		}
		setParams = append(setParams, jen.Id("value").Qual("unsafe", "Pointer"))
		getBody = append(getBody, jen.Return(jen.Id("value")))
		setBody = append(setBody, tgt.Op("=").Id("value"))

	case pointerFieldShapeIgnoredTypeSkip, pointerFieldShapePrimitivePointerSkip,
		pointerFieldShapeCallbackPointerSkip, pointerFieldShapeStructValuePointerSkip,
		pointerFieldShapeFuncPointerSkip, pointerFieldShapeUnknownPointerPointerSkip,
		pointerFieldShapeUnhandledPointerPointerSkip, pointerFieldShapeUnhandledSkip:
		o.Commentf("%v skipped due to %v", field.Name, shape.reason)
		o.Line()
		g.skips.Record(parentName+"."+field.Name, shape.reason)
		return getRetType, setParams, getBody, setBody, true

	case pointerFieldShapeCStr:
		getRetType = []jen.Code{
			jen.Op("*").Id("CStr"),
		}
		setParams = append(setParams, jen.Id("value").Op("*").Id("CStr"))
		getBody = append(getBody, jen.Return(jen.Id("wrapCStr").Params(jen.Id("value"))))
		setBody = append(setBody, tgt.Op("=").Id("value").Dot("ptr"))

	case pointerFieldShapeUint8Pointer:
		getRetType = []jen.Code{
			jen.Qual("unsafe", "Pointer"),
		}
		getBody = append(getBody, jen.Return(jen.Qual("unsafe", "Pointer").Params(jen.Id("value"))))

		setParams = append(setParams, jen.Id("value").Qual("unsafe", "Pointer"))
		setBody = append(setBody, tgt.Op("=").Params(jen.Op("*").Qual("C", shape.typeName)).Params(jen.Id("value")))

	case pointerFieldShapeEnumArray:
		getRetType = []jen.Code{
			jen.Op("*").Id("Array").Types(jen.Id(shape.typeName)),
		}

		getBody = append(
			getBody,
			jen.Return(jen.Id(fmt.Sprintf("To%vArray", shape.typeName)).Params(
				jen.Qual("unsafe", "Pointer").Params(jen.Id("value")),
			)),
		)

		setParams = append(setParams, jen.Id("value").Op("*").Id("Array").Types(jen.Id(shape.typeName)))

		cName := shape.enum.CName()

		setBody = append(
			setBody,
			jen.If(jen.Id("value").Op("!=").Id("nil")).Block(
				tgt.Clone().Op("=").Params(jen.Op("*").Qual("C", cName)).Params(jen.Id("value").Dot("ptr")),
			).Else().Block(
				tgt.Clone().Op("=").Id("nil"),
			),
		)

	case pointerFieldShapeStructPointer:
		getRetType = []jen.Code{
			jen.Op("*").Id(shape.typeName),
		}
		setParams = append(setParams, jen.Id("value").Op("*").Id(shape.typeName))

		getBody = append(
			getBody,
			jen.Var().Id("valueMapped").Op("*").Id(shape.typeName),
			jen.If(jen.Id("value").Op("!=").Id("nil")).Block(
				jen.Id("valueMapped").Op("=").Op("&").Id(shape.typeName).Values(jen.Dict{
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

	case pointerFieldShapeUnknownIdentSkip:
		log.Printf("unexpected IdentType in struct setter - struct: %v, field: %v, type: %v\n", st.Name, field.Name, shape.typeName)
		o.Commentf("%v skipped due to unexpected IdentType %v", field.Name, shape.typeName)
		o.Line()
		g.skips.Record(parentName+"."+field.Name, shape.reason)
		return getRetType, setParams, getBody, setBody, true

	case pointerFieldShapeStructPointerArray:
		getRetType = []jen.Code{
			jen.Op("*").Id("Array").Types(jen.Op("*").Id(shape.typeName)),
		}

		getBody = append(
			getBody,
			jen.Return(jen.Id(fmt.Sprintf("To%vArray", shape.typeName)).Params(
				jen.Qual("unsafe", "Pointer").Params(jen.Id("value")),
			)),
		)

		setParams = append(setParams, jen.Id("value").Op("*").Id("Array").Types(jen.Id(shape.typeName)))

		setBody = append(
			setBody,
			jen.If(jen.Id("value").Op("!=").Id("nil")).Block(
				tgt.Clone().Op("=").Params(jen.Op("**").Qual("C", shape.cName)).Params(jen.Id("value").Dot("ptr")),
			).Else().Block(
				tgt.Clone().Op("=").Id("nil"),
			),
		)
	}

	return getRetType, setParams, getBody, setBody, false
}
