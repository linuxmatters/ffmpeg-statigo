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
	argShapeConstArray
	argShapeUnhandledSkip
)

type argShape struct {
	kind     argShapeKind
	name     string
	typeName string
	goType   string
	cType    string
	enum     *Enum
	st       *Struct
	ptr      *PointerType
	array    *ConstArray
	reason   string
}

func (g *Generator) classifyArgShape(arg *Param) argShape {
	pName := convParamName(arg.Name)

	switch v := arg.Type.(type) {
	case *IdentType:
		shape := g.classifyIdent(v.Name)
		switch shape.kind {
		case identShapePrimitive:
			return argShape{
				kind:     argShapeIdentPrimitive,
				name:     pName,
				typeName: v.Name,
				goType:   shape.goType,
				cType:    getCType(v.Name, shape.goType),
			}
		case identShapeEnum:
			return argShape{kind: argShapeIdentEnum, name: pName, typeName: v.Name, enum: shape.enum}
		case identShapeByValueStruct:
			return argShape{kind: argShapeIdentByValueStruct, name: pName, typeName: v.Name, st: shape.st}
		case identShapeByPointerStruct:
			return argShape{kind: argShapeIdentByPointerStructSkip, name: pName, reason: pName}
		case identShapeCallback:
			return argShape{kind: argShapeIdentCallbackByValueSkip, name: pName, reason: fmt.Sprintf("%v (callback by value)", pName)}
		default:
			return argShape{kind: argShapeIdentUnknown, name: pName, typeName: v.Name}
		}
	case *PointerType:
		return argShape{
			kind: argShapePointer,
			name: pName,
			ptr:  v,
		}
	case *Array:
		return argShape{kind: argShapeArraySkip, name: pName, reason: pName}
	case *ConstArray:
		return argShape{kind: argShapeConstArray, name: pName, array: v}
	default:
		return argShape{
			kind:   argShapeUnhandledSkip,
			name:   pName,
			reason: fmt.Sprintf("unhandled arg type %v (%v)", reflect.TypeOf(arg.Type), pName),
		}
	}
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

func (g *Generator) classifyPointerArgShape(fn *Function, arg *Param, ptr *PointerType, pName string) pointerArgShape {
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

		shape := g.classifyIdent(inner.Name)
		switch shape.kind {
		case identShapePrimitive:
			if outputParams[fn.Name][arg.Name] {
				return pointerArgShape{
					kind:     pointerArgShapeOutputPrimitive,
					name:     pName,
					typeName: inner.Name,
					goType:   shape.goType,
					cType:    getCType(inner.Name, shape.goType),
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

type constArrayArgShapeKind int

const (
	constArrayArgShapeFixedPrimitive constArrayArgShapeKind = iota
	constArrayArgShapeBouncePrimitive
	constArrayArgShapeUint8Pointer
	constArrayArgShapeSkip
)

type constArrayArgShape struct {
	kind      constArrayArgShapeKind
	size      int64
	elemConst bool
	goType    string
	cType     string
	reason    string
}

func (g *Generator) classifyConstArrayArgShape(array *ConstArray, pName string) constArrayArgShape {
	skip := constArrayArgShape{
		kind:   constArrayArgShapeSkip,
		reason: fmt.Sprintf("const array param %v", pName),
	}
	if array == nil || array.Size <= 0 {
		return skip
	}

	switch inner := array.Inner.(type) {
	case *IdentType:
		shape := g.classifyIdent(inner.Name)
		if shape.kind != identShapePrimitive {
			return skip
		}
		cType := getCType(inner.Name, shape.goType)
		kind := constArrayArgShapeFixedPrimitive
		if needsArrayBounceBuffer(cType) {
			kind = constArrayArgShapeBouncePrimitive
		}
		return constArrayArgShape{
			kind:      kind,
			size:      array.Size,
			elemConst: array.ElemConst,
			goType:    shape.goType,
			cType:     cType,
		}
	case *PointerType:
		if innerIdent, ok := inner.Inner.(*IdentType); ok && innerIdent != nil && innerIdent.Name == "uint8_t" {
			return constArrayArgShape{
				kind:      constArrayArgShapeUint8Pointer,
				size:      array.Size,
				elemConst: array.ElemConst,
			}
		}
		return skip
	default:
		return skip
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
	constArrayFieldShapeByPointerStruct
	constArrayFieldShapeUint8PointerArray
	constArrayFieldShapeValueIdentPointerSkip
	constArrayFieldShapeStructPointerArray
	constArrayFieldShapeIdentPointerSkip
	constArrayFieldShapeMultiDimPrimitive
	constArrayFieldShapeMultiDimByValueStruct
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
	dims     []int64
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
			return constArrayFieldShape{kind: constArrayFieldShapeByPointerStruct, typeName: inner.Name, st: shape.st}
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
		return g.classifyMultiDimConstArrayFieldShape(array.Size, inner)
	default:
		return constArrayFieldShape{
			kind:   constArrayFieldShapeUnhandledInnerSkip,
			reason: fmt.Sprintf("unhandled const array inner type %v", reflect.TypeOf(array.Inner)),
		}
	}
}

func (g *Generator) classifyMultiDimConstArrayFieldShape(outerSize int64, inner *ConstArray) constArrayFieldShape {
	dims := []int64{outerSize}
	cur := Type(inner)
	for {
		nested, ok := cur.(*ConstArray)
		if !ok {
			break
		}
		dims = append(dims, nested.Size)
		cur = nested.Inner
	}
	if len(dims) < 2 || len(dims) > 3 {
		return constArrayFieldShape{kind: constArrayFieldShapeMultiDimSkip, reason: "multi dim const array"}
	}
	for _, d := range dims {
		if d <= 0 {
			return constArrayFieldShape{kind: constArrayFieldShapeMultiDimSkip, reason: "multi dim const array"}
		}
	}
	leaf, ok := cur.(*IdentType)
	if !ok {
		return constArrayFieldShape{kind: constArrayFieldShapeMultiDimSkip, reason: "multi dim const array"}
	}
	shape := g.classifyIdent(leaf.Name)
	switch shape.kind {
	case identShapePrimitive:
		return constArrayFieldShape{
			kind:     constArrayFieldShapeMultiDimPrimitive,
			typeName: leaf.Name,
			goType:   shape.goType,
			dims:     dims,
		}
	case identShapeByValueStruct:
		return constArrayFieldShape{
			kind:     constArrayFieldShapeMultiDimByValueStruct,
			typeName: leaf.Name,
			st:       shape.st,
			dims:     dims,
		}
	default:
		return constArrayFieldShape{kind: constArrayFieldShapeMultiDimSkip, reason: "multi dim const array"}
	}
}
