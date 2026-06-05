package main

import (
	"fmt"
	"reflect"
)

type identShapeKind int

const (
	identShapePrimitive identShapeKind = iota
	identShapeEnum
	identShapeByValueStruct
	identShapeByPointerStruct
	identShapeCallback
	identShapeUnknown
)

type identShape struct {
	kind   identShapeKind
	name   string
	goType string
	enum   *Enum
	st     *Struct
}

func (g *Generator) classifyIdent(name string) identShape {
	if goType, ok := primTypes[name]; ok {
		return identShape{kind: identShapePrimitive, name: name, goType: goType}
	}
	if enum, ok := g.input.enums[name]; ok {
		return identShape{kind: identShapeEnum, name: name, enum: enum}
	}
	if st, ok := g.input.structs[name]; ok {
		if st.ByValue {
			return identShape{kind: identShapeByValueStruct, name: name, st: st}
		}
		return identShape{kind: identShapeByPointerStruct, name: name, st: st}
	}
	if _, ok := g.input.callbacks[name]; ok {
		return identShape{kind: identShapeCallback, name: name}
	}
	return identShape{kind: identShapeUnknown, name: name}
}

func (g *Generator) classifyFunctionArgPreSkip(t Type) bool {
	if typeEquals(t, fileType) || typeEquals(t, fileType2) || typeEquals(t, vaListType) || typeEquals(t, tmType) {
		return true
	}

	ptr, ok := t.(*PointerType)
	if !ok {
		return false
	}

	switch inner := ptr.Inner.(type) {
	case *FuncType:
		return true
	case *IdentType:
		return inner.Name == "tm"
	default:
		return false
	}
}

type returnShapeKind int

const (
	returnShapeVoid returnShapeKind = iota
	returnShapeInt
	returnShapePrimitive
	returnShapeByValueStruct
	returnShapeByPointerStructSkip
	returnShapeCallback
	returnShapeIdent
	returnShapeUnsafePointer
	returnShapeCStrPointer
	returnShapeUint8Pointer
	returnShapePrimitivePointer
	returnShapeStructPointer
	returnShapeUnknownPointer
	returnShapePointerPointerSkip
	returnShapeUnhandledPointerSkip
	returnShapeUnhandledSkip
)

type returnShape struct {
	kind       returnShapeKind
	name       string
	goType     string
	goTypeName string
	reason     string
}

func (g *Generator) classifyReturnShape(result Type) returnShape {
	switch v := result.(type) {
	case nil:
		return returnShape{kind: returnShapeVoid}
	case *IdentType:
		if v.Name == "int" {
			return returnShape{kind: returnShapeInt, name: v.Name}
		}

		shape := g.classifyIdent(v.Name)
		switch shape.kind {
		case identShapePrimitive:
			return returnShape{kind: returnShapePrimitive, name: v.Name, goType: shape.goType}
		case identShapeByValueStruct:
			return returnShape{kind: returnShapeByValueStruct, name: v.Name}
		case identShapeByPointerStruct:
			return returnShape{kind: returnShapeByPointerStructSkip, name: v.Name, reason: "return"}
		case identShapeCallback:
			return returnShape{kind: returnShapeCallback, name: v.Name, goTypeName: g.convCamel(v.Name)}
		default:
			return returnShape{kind: returnShapeIdent, name: v.Name}
		}
	case *PointerType:
		switch inner := v.Inner.(type) {
		case nil:
			return returnShape{kind: returnShapeUnsafePointer}
		case *IdentType:
			if inner.Name == "char" {
				return returnShape{kind: returnShapeCStrPointer, name: inner.Name}
			}
			if inner.Name == "uint8_t" {
				return returnShape{kind: returnShapeUint8Pointer, name: inner.Name}
			}

			shape := g.classifyIdent(inner.Name)
			switch shape.kind {
			case identShapePrimitive:
				return returnShape{kind: returnShapePrimitivePointer, name: inner.Name, goType: shape.goType}
			case identShapeByValueStruct, identShapeByPointerStruct:
				return returnShape{kind: returnShapeStructPointer, name: inner.Name}
			default:
				return returnShape{kind: returnShapeUnknownPointer, name: inner.Name}
			}
		case *PointerType:
			return returnShape{kind: returnShapePointerPointerSkip, reason: "pointer-to-pointer return type"}
		default:
			return returnShape{
				kind:   returnShapeUnhandledPointerSkip,
				reason: fmt.Sprintf("unhandled return pointer inner type %v", reflect.TypeOf(v.Inner)),
			}
		}
	default:
		return returnShape{
			kind:   returnShapeUnhandledSkip,
			reason: fmt.Sprintf("unhandled return type %v", reflect.TypeOf(result)),
		}
	}
}

type argShapeKind int

const (
	argShapeIdentPrimitive argShapeKind = iota
	argShapeIdentEnum
	argShapeIdentByValueStruct
	argShapeIdentByPointerStructSkip
	argShapeIdentCallbackByValueSkip
	argShapeIdentUnknown
	argShapePointer
	argShapeArraySkip
	argShapeConstArraySkip
	argShapeUnhandledSkip
)

type argShape struct {
	kind           argShapeKind
	name           string
	typeName       string
	goType         string
	cType          string
	enum           *Enum
	st             *Struct
	ptr            *PointerType
	actualTypeName string
	reason         string
}

func (g *Generator) classifyArgShape(fn *Function, arg *Param) argShape {
	pName := convParamName(arg.Name)

	switch v := arg.Type.(type) {
	case *IdentType:
		typeName := g.correctIdentArgTypeName(fn, arg, v.Name)
		shape := g.classifyIdent(typeName)
		switch shape.kind {
		case identShapePrimitive:
			return argShape{
				kind:     argShapeIdentPrimitive,
				name:     pName,
				typeName: typeName,
				goType:   shape.goType,
				cType:    getCType(typeName, shape.goType),
			}
		case identShapeEnum:
			return argShape{kind: argShapeIdentEnum, name: pName, typeName: typeName, enum: shape.enum}
		case identShapeByValueStruct:
			return argShape{kind: argShapeIdentByValueStruct, name: pName, typeName: typeName, st: shape.st}
		case identShapeByPointerStruct:
			return argShape{kind: argShapeIdentByPointerStructSkip, name: pName, reason: pName}
		case identShapeCallback:
			return argShape{kind: argShapeIdentCallbackByValueSkip, name: pName, reason: fmt.Sprintf("%v (callback by value)", pName)}
		default:
			return argShape{kind: argShapeIdentUnknown, name: pName, typeName: v.Name}
		}
	case *PointerType:
		return argShape{
			kind:           argShapePointer,
			name:           pName,
			ptr:            v,
			actualTypeName: g.correctPointerArgTypeName(fn, arg, v),
		}
	case *Array:
		return argShape{kind: argShapeArraySkip, name: pName, reason: pName}
	case *ConstArray:
		return argShape{kind: argShapeConstArraySkip, name: pName, reason: fmt.Sprintf("const array param %v", pName)}
	default:
		return argShape{
			kind:   argShapeUnhandledSkip,
			name:   pName,
			reason: fmt.Sprintf("unhandled arg type %v (%v)", reflect.TypeOf(arg.Type), pName),
		}
	}
}

func (g *Generator) correctIdentArgTypeName(fn *Function, arg *Param, typeName string) string {
	if typeName == "int" && arg.Name == "buf_size" {
		if fn.Name == "av_channel_name" || fn.Name == "av_channel_description" ||
			fn.Name == "av_channel_layout_describe" {
			return "size_t"
		}
	} else if typeName == "int" && arg.Name == "max_size" {
		if fn.Name == "avio_read_to_bprint" {
			return "size_t"
		}
	}

	return typeName
}

func (g *Generator) correctPointerArgTypeName(fn *Function, arg *Param, ptr *PointerType) string {
	ident, ok := ptr.Inner.(*IdentType)
	if !ok {
		return ""
	}

	actualTypeName := ident.Name
	if p, ok := outputParams[fn.Name][arg.Name]; ok && p.sizeT && actualTypeName == "int" {
		return "size_t"
	}

	return actualTypeName
}

type pointerArgShapeKind int

const (
	pointerArgShapeUnsafePointer pointerArgShapeKind = iota
	pointerArgShapeCStr
	pointerArgShapeBytePointer
	pointerArgShapeOutputPrimitive
	pointerArgShapeNonOutputPrimitiveSkip
	pointerArgShapeByValueStructSkip
	pointerArgShapeByPointerStruct
	pointerArgShapeEnum
	pointerArgShapeCallback
	pointerArgShapeUnknownIdent
	pointerArgShapePointerToStruct
	pointerArgShapePointerToPointerSkip
	pointerArgShapeUnhandledSkip
)

type pointerArgShape struct {
	kind     pointerArgShapeKind
	name     string
	typeName string
	goType   string
	cType    string
	enum     *Enum
	st       *Struct
	reason   string
}

func (g *Generator) classifyPointerArgShape(fn *Function, arg *Param, ptr *PointerType, pName, actualTypeName string) pointerArgShape {
	switch inner := ptr.Inner.(type) {
	case nil:
		return pointerArgShape{kind: pointerArgShapeUnsafePointer, name: pName}
	case *IdentType:
		switch inner.Name {
		case "char":
			return pointerArgShape{kind: pointerArgShapeCStr, name: pName, typeName: inner.Name}
		case "uint8_t", "uchar":
			return pointerArgShape{kind: pointerArgShapeBytePointer, name: pName, typeName: inner.Name}
		}

		typeName := inner.Name
		if actualTypeName != "" {
			typeName = actualTypeName
		}

		shape := g.classifyIdent(typeName)
		switch shape.kind {
		case identShapePrimitive:
			if _, ok := outputParams[fn.Name][arg.Name]; ok {
				return pointerArgShape{
					kind:     pointerArgShapeOutputPrimitive,
					name:     pName,
					typeName: typeName,
					goType:   shape.goType,
					cType:    getCType(typeName, shape.goType),
				}
			}
			return pointerArgShape{kind: pointerArgShapeNonOutputPrimitiveSkip, name: pName, reason: fmt.Sprintf("%v (non-output primitive pointer)", pName)}
		case identShapeByValueStruct:
			return pointerArgShape{kind: pointerArgShapeByValueStructSkip, name: pName, reason: pName}
		case identShapeByPointerStruct:
			return pointerArgShape{kind: pointerArgShapeByPointerStruct, name: pName, typeName: inner.Name, st: shape.st}
		case identShapeEnum:
			return pointerArgShape{kind: pointerArgShapeEnum, name: pName, typeName: inner.Name, enum: shape.enum}
		case identShapeCallback:
			return pointerArgShape{kind: pointerArgShapeCallback, name: pName, typeName: inner.Name}
		default:
			return pointerArgShape{kind: pointerArgShapeUnknownIdent, name: pName, typeName: inner.Name}
		}
	case *PointerType:
		switch innerInner := inner.Inner.(type) {
		case *IdentType:
			if innerInner.Name == "uint8_t" || innerInner.Name == "char" {
				return pointerArgShape{kind: pointerArgShapePointerToPointerSkip, name: pName, reason: pName}
			}
			if _, ok := primTypes[innerInner.Name]; ok {
				return pointerArgShape{kind: pointerArgShapePointerToPointerSkip, name: pName, reason: pName}
			}
			if st, ok := g.input.structs[innerInner.Name]; ok {
				return pointerArgShape{kind: pointerArgShapePointerToStruct, name: pName, typeName: innerInner.Name, st: st}
			}
			return pointerArgShape{kind: pointerArgShapePointerToPointerSkip, name: pName, reason: pName}
		default:
			return pointerArgShape{kind: pointerArgShapePointerToPointerSkip, name: pName, reason: pName}
		}
	default:
		return pointerArgShape{
			kind:   pointerArgShapeUnhandledSkip,
			name:   pName,
			reason: fmt.Sprintf("unhandled pointer arg inner type %v (%v)", reflect.TypeOf(ptr.Inner), pName),
		}
	}
}

type fieldShapeKind int

const (
	fieldShapeIdentPrimitive fieldShapeKind = iota
	fieldShapeIdentByValueStruct
	fieldShapeIdentRefStruct
	fieldShapeIdentCallbackSkip
	fieldShapeIdentEnum
	fieldShapeIdentUnknownRef
	fieldShapePointer
	fieldShapeConstArray
	fieldShapeUnionSkip
	fieldShapeUnhandledSkip
)

type fieldShape struct {
	kind     fieldShapeKind
	typeName string
	goType   string
	anonEnum bool
	enum     *Enum
	st       *Struct
	ptr      *PointerType
	array    *ConstArray
	reason   string
}

func (g *Generator) classifyFieldShape(field *Field) fieldShape {
	switch v := field.Type.(type) {
	case *IdentType:
		shape := g.classifyIdent(v.Name)
		switch shape.kind {
		case identShapePrimitive:
			return fieldShape{kind: fieldShapeIdentPrimitive, typeName: v.Name, goType: shape.goType, anonEnum: v.IsAnonEnum}
		case identShapeByValueStruct:
			return fieldShape{kind: fieldShapeIdentByValueStruct, typeName: v.Name, st: shape.st}
		case identShapeByPointerStruct:
			return fieldShape{kind: fieldShapeIdentRefStruct, typeName: v.Name, st: shape.st}
		case identShapeCallback:
			return fieldShape{kind: fieldShapeIdentCallbackSkip, reason: "ident callback"}
		case identShapeEnum:
			return fieldShape{kind: fieldShapeIdentEnum, typeName: v.Name, enum: shape.enum}
		default:
			return fieldShape{kind: fieldShapeIdentUnknownRef, typeName: v.Name}
		}
	case *PointerType:
		return fieldShape{kind: fieldShapePointer, ptr: v}
	case *ConstArray:
		return fieldShape{kind: fieldShapeConstArray, array: v}
	case *UnionType:
		return fieldShape{kind: fieldShapeUnionSkip, reason: "union type"}
	default:
		return fieldShape{
			kind:   fieldShapeUnhandledSkip,
			reason: fmt.Sprintf("unhandled field type %v", reflect.TypeOf(field.Type)),
		}
	}
}

type pointerFieldShapeKind int

const (
	pointerFieldShapeUnsafePointer pointerFieldShapeKind = iota
	pointerFieldShapeIgnoredTypeSkip
	pointerFieldShapeCStr
	pointerFieldShapeUint8Pointer
	pointerFieldShapePrimitivePointerSkip
	pointerFieldShapeEnumArray
	pointerFieldShapeCallbackPointerSkip
	pointerFieldShapeStructValuePointerSkip
	pointerFieldShapeStructPointer
	pointerFieldShapeUnknownIdentSkip
	pointerFieldShapeFuncPointerSkip
	pointerFieldShapeStructPointerArray
	pointerFieldShapeUnknownPointerPointerSkip
	pointerFieldShapeUnhandledPointerPointerSkip
	pointerFieldShapeUnhandledSkip
)

type pointerFieldShape struct {
	kind     pointerFieldShapeKind
	typeName string
	cName    string
	enum     *Enum
	st       *Struct
	reason   string
}

func (g *Generator) classifyPointerFieldShape(ptr *PointerType) pointerFieldShape {
	switch inner := ptr.Inner.(type) {
	case nil:
		return pointerFieldShape{kind: pointerFieldShapeUnsafePointer}
	case *IdentType:
		if inner.Name == "URLContext" || inner.Name == "AVFilterCommand" || inner.Name == "AVCodecInternal" {
			return pointerFieldShape{kind: pointerFieldShapeIgnoredTypeSkip, typeName: inner.Name, reason: "ptr to ignored type"}
		}
		if inner.Name == "char" {
			return pointerFieldShape{kind: pointerFieldShapeCStr, typeName: inner.Name}
		}
		if inner.Name == "uint8_t" {
			return pointerFieldShape{kind: pointerFieldShapeUint8Pointer, typeName: inner.Name}
		}

		shape := g.classifyIdent(inner.Name)
		switch shape.kind {
		case identShapePrimitive:
			return pointerFieldShape{kind: pointerFieldShapePrimitivePointerSkip, typeName: inner.Name, reason: "prim ptr"}
		case identShapeEnum:
			return pointerFieldShape{kind: pointerFieldShapeEnumArray, typeName: inner.Name, enum: shape.enum}
		case identShapeCallback:
			return pointerFieldShape{kind: pointerFieldShapeCallbackPointerSkip, typeName: inner.Name, reason: "callback ptr"}
		case identShapeByValueStruct:
			return pointerFieldShape{kind: pointerFieldShapeStructValuePointerSkip, typeName: inner.Name, reason: "struct value ptr"}
		case identShapeByPointerStruct:
			return pointerFieldShape{kind: pointerFieldShapeStructPointer, typeName: inner.Name, st: shape.st}
		default:
			return pointerFieldShape{kind: pointerFieldShapeUnknownIdentSkip, typeName: inner.Name, reason: fmt.Sprintf("unexpected IdentType %v", inner.Name)}
		}
	case *FuncType:
		return pointerFieldShape{kind: pointerFieldShapeFuncPointerSkip, reason: "func ptr"}
	case *PointerType:
		switch innerInner := inner.Inner.(type) {
		case *IdentType:
			if st, ok := g.input.structs[innerInner.Name]; ok {
				return pointerFieldShape{kind: pointerFieldShapeStructPointerArray, typeName: innerInner.Name, st: st, cName: st.CName()}
			}
			return pointerFieldShape{kind: pointerFieldShapeUnknownPointerPointerSkip, reason: "unknown pointer pair"}
		default:
			return pointerFieldShape{
				kind:   pointerFieldShapeUnhandledPointerPointerSkip,
				reason: fmt.Sprintf("unhandled pointer-pointer inner type %v", reflect.TypeOf(inner.Inner)),
			}
		}
	default:
		return pointerFieldShape{
			kind:   pointerFieldShapeUnhandledSkip,
			reason: fmt.Sprintf("unhandled pointer inner type %v", reflect.TypeOf(ptr.Inner)),
		}
	}
}

type constArrayFieldShapeKind int

const (
	constArrayFieldShapePrimitive constArrayFieldShapeKind = iota
	constArrayFieldShapeEnum
	constArrayFieldShapeByValueStruct
	constArrayFieldShapeStructPointerSkip
	constArrayFieldShapeUint8PointerArray
	constArrayFieldShapeValueIdentPointerSkip
	constArrayFieldShapeStructPointerArray
	constArrayFieldShapeIdentPointerSkip
	constArrayFieldShapeMultiDimSkip
	constArrayFieldShapeUnknownTypeSkip
	constArrayFieldShapeUnhandledPointerInnerSkip
	constArrayFieldShapeUnhandledInnerSkip
)

type constArrayFieldShape struct {
	kind     constArrayFieldShapeKind
	typeName string
	goType   string
	enum     *Enum
	st       *Struct
	reason   string
}

func (g *Generator) classifyConstArrayFieldShape(array *ConstArray) constArrayFieldShape {
	switch inner := array.Inner.(type) {
	case *IdentType:
		shape := g.classifyIdent(inner.Name)
		switch shape.kind {
		case identShapePrimitive:
			return constArrayFieldShape{kind: constArrayFieldShapePrimitive, typeName: inner.Name, goType: shape.goType}
		case identShapeEnum:
			return constArrayFieldShape{kind: constArrayFieldShapeEnum, typeName: inner.Name, enum: shape.enum}
		case identShapeByValueStruct:
			return constArrayFieldShape{kind: constArrayFieldShapeByValueStruct, typeName: inner.Name, st: shape.st}
		case identShapeByPointerStruct:
			return constArrayFieldShape{kind: constArrayFieldShapeStructPointerSkip, typeName: inner.Name, reason: "const array of struct pointers"}
		default:
			return constArrayFieldShape{
				kind:     constArrayFieldShapeUnknownTypeSkip,
				typeName: inner.Name,
				reason:   fmt.Sprintf("unknown const array type %v", inner.Name),
			}
		}
	case *PointerType:
		switch innerInner := inner.Inner.(type) {
		case *IdentType:
			if innerInner.Name == "uint8_t" {
				return constArrayFieldShape{kind: constArrayFieldShapeUint8PointerArray, typeName: innerInner.Name}
			}
			shape := g.classifyIdent(innerInner.Name)
			switch shape.kind {
			case identShapeByValueStruct:
				return constArrayFieldShape{kind: constArrayFieldShapeValueIdentPointerSkip, typeName: innerInner.Name, reason: "value ident ptr const array"}
			case identShapeByPointerStruct:
				return constArrayFieldShape{kind: constArrayFieldShapeStructPointerArray, typeName: innerInner.Name, st: shape.st}
			default:
				return constArrayFieldShape{kind: constArrayFieldShapeIdentPointerSkip, typeName: innerInner.Name, reason: "ident pointer const array"}
			}
		default:
			return constArrayFieldShape{
				kind:   constArrayFieldShapeUnhandledPointerInnerSkip,
				reason: fmt.Sprintf("unhandled const array pointer inner type %v", reflect.TypeOf(inner.Inner)),
			}
		}
	case *ConstArray:
		return constArrayFieldShape{kind: constArrayFieldShapeMultiDimSkip, reason: "multi dim const array"}
	default:
		return constArrayFieldShape{
			kind:   constArrayFieldShapeUnhandledInnerSkip,
			reason: fmt.Sprintf("unhandled const array inner type %v", reflect.TypeOf(array.Inner)),
		}
	}
}
