package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"modernc.org/cc/v4"
)

type Type interface {
	typeMarker()

	String() string

	Equals(other Type) bool
}

func typeEquals(a Type, b Type) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	return a.Equals(b)
}

type IdentType struct {
	Name       string
	IsAnonEnum bool
}

func (t *IdentType) typeMarker()    {}
func (t *IdentType) String() string { return t.Name }

func (t *IdentType) Equals(other Type) bool {
	ot, ok := other.(*IdentType)
	if !ok {
		return false
	}

	return t.Name == ot.Name
}

type PointerType struct {
	Inner Type
}

func (t *PointerType) typeMarker()    {}
func (t *PointerType) String() string { return fmt.Sprintf("*%v", t.Inner) }

func (t *PointerType) Equals(other Type) bool {
	op, ok := other.(*PointerType)
	if !ok {
		return false
	}

	return typeEquals(t.Inner, op.Inner)
}

type FuncType struct {
	Args   []Type
	Result Type
}

func (t *FuncType) typeMarker()    {}
func (t *FuncType) String() string { return "func" }

func (t *FuncType) Equals(other Type) bool {
	of, ok := other.(*FuncType)
	if !ok {
		return false
	}

	if !typeEquals(t.Result, of.Result) {
		return false
	}

	if len(t.Args) != len(of.Args) {
		return false
	}

	for i := 0; i < len(t.Args); i++ {
		if !typeEquals(t.Args[i], of.Args[i]) {
			return false
		}
	}

	return true
}

type ConstArray struct {
	Inner Type
	Size  int64
}

func (t *ConstArray) typeMarker()    {}
func (t *ConstArray) String() string { return fmt.Sprintf("%v[%v]", t.Inner, t.Size) }

func (t *ConstArray) Equals(other Type) bool {
	ot, ok := other.(*ConstArray)
	if !ok {
		return false
	}

	return t.Size == ot.Size && typeEquals(t.Inner, ot.Inner)
}

type Array struct {
	Inner Type
}

func (t *Array) typeMarker()    {}
func (t *Array) String() string { return fmt.Sprintf("%v[?]", t.Inner) }

func (t *Array) Equals(other Type) bool {
	ot, ok := other.(*Array)
	if !ok {
		return false
	}

	return typeEquals(t.Inner, ot.Inner)
}

type UnionType struct {
	Fields []*Field
}

func (t *UnionType) typeMarker()    {}
func (t *UnionType) String() string { return "union" }

func (t *UnionType) Equals(other Type) bool {
	of, ok := other.(*UnionType)
	if !ok {
		return false
	}

	if len(t.Fields) != len(of.Fields) {
		return false
	}

	for i := 0; i < len(t.Fields); i++ {
		if t.Fields[i].Name != of.Fields[i].Name {
			return false
		}

		if !typeEquals(t.Fields[i].Type, of.Fields[i].Type) {
			return false
		}
	}

	return true
}

// unnamedStructComment is the Comment a synthetic record carries. It is a
// literal, not a parsed comment: the record has no declaration of its own to
// take a comment from, so registerUnnamedStruct writes this exact string and
// the golden pins it.
const unnamedStructComment = "Synthetic type for unnamed struct"

// parseType converts a cc/v4 type into the generator's Type tree. renderType
// writes that tree straight into the IR golden, so every node it builds is
// pinned.
//
// A nil result is C void. renderType prints it as "void".
func (w *ccWalk) parseType(indent string, t cc.Type) Type {
	if t == nil {
		return nil
	}

	// Typedef sugar wins at every level, and is tested before the kind switch.
	// This is the _t-suffix guard: the written typedef spelling is kept, so
	// size_t, ptrdiff_t and the intN_t family reach the IR as themselves rather
	// than as the widths they canonicalise to. It also covers every FFmpeg
	// typedef the golden names directly, from AVRational to the function
	// typedefs such as av_pixelutils_sad_fn.
	//
	// Keep it. WW-141 probe answer 4 measured that this guard, the sizeT flags
	// in outputParams and the actualTypeName == "int" test in
	// correctPointerArgTypeName are one unit. Drop the sugar here and
	// actualTypeName falls back to "int", at which point
	// correctPointerArgTypeName silently emits int for every size_t output
	// parameter. Phase 5.2 owns the table entries; do not change one alone.
	if td := t.Typedef(); td != nil && td.Name() != "" {
		log.Println(indent, "Parsing type as typedef", td.Name())

		return &IdentType{Name: td.Name()}
	}

	switch t.Kind() {
	case cc.Void:
		log.Println(indent, "Parsing type as void")

		return nil

	// The base kinds the IR names. They take the canonical C spelling with
	// "unsigned " rewritten to a "u" prefix, so "unsigned int" is uint and
	// "unsigned char" is uchar. Any other base kind falls through to the default
	// branch and records a skip.
	case cc.Int:
		return &IdentType{Name: "int"}

	case cc.UInt:
		return &IdentType{Name: "uint"}

	case cc.Long:
		return &IdentType{Name: "long"}

	case cc.ULong:
		return &IdentType{Name: "ulong"}

	case cc.Char:
		return &IdentType{Name: "char"}

	case cc.UChar:
		return &IdentType{Name: "uchar"}

	case cc.Float:
		return &IdentType{Name: "float"}

	case cc.Double:
		return &IdentType{Name: "double"}

	case cc.Ptr:
		pt, ok := t.(*cc.PointerType)
		if !ok {
			return w.skipType(indent, t)
		}

		return &PointerType{Inner: w.parseType(indent, pt.Elem())}

	case cc.Array:
		at, ok := t.(*cc.ArrayType)
		if !ok {
			return w.skipType(indent, t)
		}

		inner := w.parseType(indent, at.Elem())

		// Len is negative for an array of unspecified size, which the IR spells
		// as a plain Array rather than a ConstArray.
		if at.Len() < 0 {
			return &Array{Inner: inner}
		}

		return &ConstArray{Inner: inner, Size: at.Len()}

	case cc.Function:
		ft, ok := t.(*cc.FunctionType)
		if !ok {
			return w.skipType(indent, t)
		}

		var args []Type

		for _, p := range ccParameters(ft) {
			args = append(args, w.parseType(indent, ccParamType(p.Type())))
		}

		return &FuncType{Args: args, Result: w.parseType(indent, ft.Result())}

	case cc.Enum:
		return w.parseEnumType(indent, t)

	case cc.Struct:
		st, ok := t.(*cc.StructType)
		if !ok {
			return w.skipType(indent, t)
		}

		tag := st.Tag()
		if name := tag.SrcStr(); name != "" {
			return &IdentType{Name: name}
		}

		return w.unnamedStructType(indent)

	case cc.Union:
		return w.parseUnionType(indent, t)

	default:
		return w.skipType(indent, t)
	}
}

// parseEnumType routes "enum X" and "enum {...}" apart. cc/v4 has no node for
// the difference, so the tag token decides.
//
// A tagged enum is its own IdentType. A tagless one is not a type the emitter
// can name, so its members become loose constants and the field takes the
// uint32_t/anonymous-enum form. The golden holds no anonymous-enum type today,
// so this path adds nothing to constantOrder on the pinned headers.
func (w *ccWalk) parseEnumType(indent string, t cc.Type) Type {
	et, ok := t.(*cc.EnumType)
	if !ok {
		return w.skipType(indent, t)
	}

	tag := et.Tag()
	if tag.Seq() != 0 {
		return &IdentType{Name: tag.SrcStr()}
	}

	log.Println(indent, "Treating unnamed enum as constants")

	for _, en := range et.Enumerators() {
		name := en.Token.SrcStr()
		if name == "" {
			continue
		}

		w.mod.constants[name] = &Constant{Name: name, FromEnum: true}
		w.mod.constantOrder = append(w.mod.constantOrder, name)
	}

	return &IdentType{Name: "uint32_t", IsAnonEnum: true}
}

// parseUnionType expands a union written in place, field by field, tagged or
// not. A union reached through a typedef never gets here: parseType returns its
// typedef name before the kind switch.
func (w *ccWalk) parseUnionType(indent string, t cc.Type) Type {
	ut, ok := t.(*cc.UnionType)
	if !ok {
		return w.skipType(indent, t)
	}

	u := &UnionType{}

	for i := 0; i < ut.NumFields(); i++ {
		f := ut.FieldByIndex(i)
		if f == nil {
			continue
		}

		name := f.Name()
		if name == "" {
			failLog.Fatal("no field name")
		}

		u.Fields = append(u.Fields, &Field{
			Name: name,
			Type: w.parseType(fmt.Sprintf("%v[%v]", indent, name), f.Type()),
		})
	}

	return u
}

// unnamedStructType names and registers the tagless struct written inline as
// the base type of the field currently being walked, and returns a reference to
// it.
//
// The position has to come from the syntax node structFields put on the walk.
// A StructType carries no tag token when the record is tagless, and nothing
// reachable from the type gives the "struct" keyword, so there is no other
// route to a name.
func (w *ccWalk) unnamedStructType(indent string) Type {
	if w.anon == nil {
		// Only a field declaration puts a specifier on the walk, so a tagless
		// struct written directly as a function's return type or parameter
		// reaches here with none. C allows that shape and cc/v4 accepts it, but
		// the record has no compatible type any other declaration can name, so
		// no caller can pass or receive one. FFmpeg's public headers hold none;
		// every tagless record in include/ is a struct field or a typedef. A
		// tagless union is not affected either way, because parseUnionType
		// expands it in place and never asks for a name.
		//
		// Abort rather than record a skip. skipType returns nil, which is the
		// parser's encoding of C void: a skipped return type would emit a Go
		// function that silently drops its result, and the IR would spell it
		// "void". A skipped parameter does reach the emitter as a loud
		// "unhandled arg type" skip, but the return type does not, so stopping
		// is the one answer that is safe in both positions.
		//
		// The indent breadcrumb names the symbol and the slot, so the abort
		// reads as "[av_thing][return] ..." and points straight at the
		// declaration. TestUnnamedStructWithoutSpecifierAborts pins that.
		log.Panicln(indent, "tagless record type with no specifier on the walk")
	}

	name := ccSyntheticStructName(w.anon)
	if name == "" {
		log.Panicln(indent, "No name")
	}

	w.registerUnnamedStruct(indent, w.anon, name)

	return &IdentType{Name: name}
}

// ccSyntheticStructName builds the name a tagless record is registered under,
// from the position of its "struct" or "union" keyword.
//
// The record has no name of its own, so the position supplies one: the file
// stem, the line and the column. cc/v4 carries all three on the keyword token
// (WW-141 probe answer 5), so libavformat/avformat.h:986:5 gives
// UnnamedStruct_avformat_986_5, the string generator.go hardcodes in
// skippedStructs.
func ccSyntheticStructName(sus *cc.StructOrUnionSpecifier) string {
	if sus == nil || sus.StructOrUnion == nil {
		return ""
	}

	p := sus.StructOrUnion.Token.Position()

	base := filepath.Base(p.Filename)
	base = cleanIdentifier(strings.TrimSuffix(base, filepath.Ext(base)))

	return fmt.Sprintf("UnnamedStruct_%s_%d_%d", base, p.Line, p.Column)
}

// registerUnnamedStruct adds a tagless record as a struct of its own, under the
// synthetic name, so the emitter has a type to refer to. The golden pins all
// three of the following:
//
//   - Typedefd is false. The record is not written as a typedef.
//   - Comment is the unnamedStructComment literal. The record has no comment of
//     its own to claim, so the matcher is never asked for one. Its fields do,
//     and claim theirs the way any other field does.
//   - CTypeName is left unset on every field, which is why the golden pins all
//     three of this record's ctype values as empty.
//
// structOrder takes the entry where parseType first meets the type, which is
// while the owning struct's own fields are being walked, so the synthetic record
// lands immediately ahead of its owner.
func (w *ccWalk) registerUnnamedStruct(indent string, sus *cc.StructOrUnionSpecifier, name string) {
	if _, ok := w.mod.structs[name]; ok {
		log.Println(indent, "Unnamed struct already registered:", name)

		return
	}

	log.Println(indent, "Registering unnamed struct:", name)

	s := &Struct{
		Name:     name,
		Typedefd: false,
		Comment:  unnamedStructComment,
	}

	ccEachFieldDeclarator(sus, func(sd *cc.StructDeclarator, anon *cc.StructOrUnionSpecifier) {
		dcl := sd.Declarator

		prev := w.anon
		w.anon = anon

		s.Fields = append(s.Fields, &Field{
			Name:     dcl.Name(),
			Type:     w.parseType(fmt.Sprintf("%v[%v]", indent, dcl.Name()), dcl.Type()),
			BitWidth: ccFieldBitWidth(sd),
			Comment:  w.comment(dcl.NameTok(), 0),
		})

		w.anon = prev
	})

	w.mod.structs[name] = s
	w.mod.structOrder = append(w.mod.structOrder, name)

	log.Println(indent, "Successfully registered unnamed struct:", name, "with", len(s.Fields), "fields")
}

// skipType records an unhandled type kind and returns nil. The pinned headers
// produce zero parse-layer skips, so this must never fire; it exists so a new
// type kind shows up in skips.txt rather than as a silent void.
func (w *ccWalk) skipType(indent string, t cc.Type) Type {
	reason := fmt.Sprintf("unhandled type kind %v", t.Kind())
	log.Printf("%v skipped due to %v", indent, reason)
	w.skips.Record(indent, reason)

	return nil
}

// ccParameters lists a function type's parameters, dropping the single unnamed
// void of an empty parameter list. cc/v4's newFunctionType keeps "f(void)" as
// one Void parameter and sets maxArgs to zero; the IR takes it as no arguments
// at all, and the golden holds no "func()" expansion.
func ccParameters(ft *cc.FunctionType) []*cc.Parameter {
	fp := ft.Parameters()

	if len(fp) == 1 {
		if pt := fp[0].Type(); pt != nil && pt.Kind() == cc.Void {
			return nil
		}
	}

	return fp
}

// ccParamType returns the parameter type to walk.
//
// cc/v4 decays an array or function parameter to a pointer, but the IR carries
// the written type. So "uint8_t *data[4]" has to read as an array of pointers,
// and a parameter of a function typedef as that typedef.
//
// The decayed pointer is tested for typedef sugar first, because the sugar sits
// on different sides for different typedefs: AVUUID, an array typedef, can keep
// its name on the pointer, while AVFifoCB, a function typedef, keeps it only on
// the undecayed function type (WW-141 probe answer 4).
func ccParamType(t cc.Type) cc.Type {
	if t == nil {
		return nil
	}

	if t.Typedef() != nil {
		return t
	}

	return t.Undecay()
}

// cleanIdentifier converts a string into a valid Go identifier
func cleanIdentifier(s string) string {
	s = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return '_'
	}, s)

	for strings.Contains(s, "__") {
		s = strings.ReplaceAll(s, "__", "_")
	}

	s = strings.Trim(s, "_")

	return s
}
