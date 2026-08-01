package main

// Declaration walk for the modernc.org/cc/v4 parser backend (WW-141 phase 2.1).
//
// This file owns the walk: which symbols exist, under which names, in which
// order, and with which comment. The type tree is in type.go, the C type
// spelling in ctypename.go and the comment table and its claim rules in
// comments.go.

import (
	"fmt"
	"log"
	"math"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"slices"
	"strings"

	"modernc.org/cc/v4"
)

// ccModulePath is the parser module, looked up in the build info so Version
// reports the pinned release rather than a constant that drifts from go.mod.
const ccModulePath = "modernc.org/cc/v4"

// ccParser is the cc/v4-backed HeaderParser implementation. It holds no state;
// ccWalk carries the per-run accumulator.
type ccParser struct{}

// Version names the backend and the parser module release in use.
func (ccParser) Version() string {
	return ccModulePath + " " + ccModuleVersion()
}

// ccModuleVersion reads the parser module version out of the build info. A
// binary built without module information reports the fallback rather than
// failing, because the string only feeds the verbose run log.
func ccModuleVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "(version unknown)"
	}

	for _, d := range info.Deps {
		if d.Path == ccModulePath {
			return d.Version
		}
	}

	return "(version unknown)"
}

// Parse translates every pinned FFmpeg header with cc/v4 and walks the
// resulting AST into a Module. It is the HeaderParser implementation for the
// cc/v4 backend.
//
// Each header is a translation unit of its own, so a symbol reached through an
// #include is filtered out by the location test in declaredHere and recorded
// only when its own header is parsed.
func (ccParser) Parse(skips *SkipCollector) *Module {
	cfg, err := newCCConfig(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		failLog.Fatalf("cc/v4 configuration failed: %v", err)
	}

	// The verbose run log carries the discovered include paths. The config is
	// built here rather than in run, so the trace is written here too.
	log.Printf("include paths: %v", hostIncludePaths(cfg))

	w := &ccWalk{
		mod: &Module{
			functions: make(map[string]*Function),
			structs:   make(map[string]*Struct),
			enums:     make(map[string]*Enum),
			callbacks: make(map[string]*Function),
			constants: make(map[string]*Constant),
		},
		skips: skips,
	}

	for _, file := range files {
		w.abs = filepath.Join(AVLibPath, file)
		w.tags = make(map[string]bool)

		log.Println(w.abs)

		w.table, err = BuildCommentTable(w.abs)
		if err != nil {
			failLog.Fatalf("failed to read header %v: %v", w.abs, err)
		}

		ast, err := ccTranslate(cfg, w.abs)
		if err != nil {
			failLog.Fatalf("failed to parse header %v: %v", w.abs, err)
		}

		w.walkAST(ast)
	}

	return w.mod
}

// ccTranslate translates one header into an AST.
//
// cc.Translate panics through its todo() path on branches it does not handle,
// so the panic is recovered and returned as an error: a parse failure must
// surface as a named header the caller can report, never as a stack trace or a
// silently empty AST.
//
// The predefined source has to be named exactly "<predefined>" or cc/v4 panics
// in its preprocessor. cc.Builtin supplies the compiler builtins the module
// ships in place of system headers.
func ccTranslate(cfg *cc.Config, abs string) (ast *cc.AST, err error) {
	defer func() {
		if r := recover(); r != nil {
			ast = nil
			err = fmt.Errorf("panic translating %s: %v", abs, r)
		}
	}()

	return cc.Translate(cfg, []cc.Source{
		{Name: "<predefined>", Value: cfg.Predefined},
		{Name: "<builtin>", Value: cc.Builtin},
		{Name: abs},
	})
}

// ccWalk accumulates one Module across the whole files run. abs is the header
// currently being translated and is the location filter every record passes
// through. tags holds the struct and union tags already declared in the current
// translation unit and is reset for every header.
//
// anon is the tagless struct or union specifier of the field declaration being
// walked, or nil. The type walk needs it: a tagless record carries no tag token,
// so its syntax node is the only route to the position both its synthetic name
// and its C spelling are built from.
//
// table is the comment table of the header being translated, rebuilt per header
// alongside tags. Every doc comment a record carries is claimed from it.
//
// macroTable is the translation unit's macro table. nameSpan needs it to undo
// the position every macro expansion collapses onto its invocation site.
type ccWalk struct {
	mod        *Module
	abs        string
	skips      *SkipCollector // carried for the type walk added in a later phase
	tags       map[string]bool
	anon       *cc.StructOrUnionSpecifier
	table      *CommentTable
	macroTable map[string]*cc.Macro
}

// walkAST records one translation unit. Macros are collected before the
// declarations because constantOrder holds every macro of a header ahead of
// that header's unnamed-enum members.
func (w *ccWalk) walkAST(ast *cc.AST) {
	w.macroTable = ast.Macros

	w.macros(ast)

	for tu := ast.TranslationUnit; tu != nil; tu = tu.TranslationUnit {
		ed := tu.ExternalDeclaration
		if ed == nil {
			continue
		}

		switch ed.Case {
		case cc.ExternalDeclarationDecl:
			w.declaration(ed.Declaration)

		case cc.ExternalDeclarationFuncDef:
			// A static inline body in a header still declares a function, so
			// av_make_q, av_cmp_q and av_q2d are bound like any other. Walking
			// only ExternalDeclarationDecl loses them.
			fd := ed.FunctionDefinition
			if fd == nil {
				continue
			}

			w.declarationSpecifiers(fd.DeclarationSpecifiers, "", false)
			w.recordFunction(fd.Declarator)
		}
	}
}

// declaration handles one top-level declaration. The type specifiers come
// first, because a tagless struct or enum takes its name from the typedef
// declarator that follows them.
func (w *ccWalk) declaration(d *cc.Declaration) {
	if d == nil || d.Case != cc.DeclarationDecl {
		return
	}

	isTypedef := ccIsTypedef(d)

	fallback := ""
	if isTypedef {
		fallback = ccFirstDeclaratorName(d)
	}

	w.declarationSpecifiers(d.DeclarationSpecifiers, fallback, isTypedef)

	if isTypedef {
		// A typedef over a tag alone, "typedef struct AVFrame AVFrame;", still
		// declares the tag, so the tag is recorded from the underlying type.
		w.recordTags(d.DeclarationSpecifiers, true, true)

		// A function typedef belongs in mod.callbacks, never in mod.functions,
		// so the declarators of a typedef are not functions.
		w.recordCallbacks(d)

		return
	}

	// A tag with no body in the declaration specifiers is recorded as an opaque
	// type wherever the tag is first introduced. That covers both the bare
	// "struct AVCodecTag;" form and a tag introduced by the return type of a
	// declaration, which is the only place struct AVMurMur3 ever appears.
	w.recordTags(d.DeclarationSpecifiers, false, d.InitDeclaratorList != nil)

	for idl := d.InitDeclaratorList; idl != nil; idl = idl.InitDeclaratorList {
		if idl.InitDeclarator == nil {
			continue
		}

		w.recordFunction(idl.InitDeclarator.Declarator)
	}
}

// ccIsTypedef reports whether the declaration carries the typedef storage
// class.
func ccIsTypedef(d *cc.Declaration) bool {
	for ds := d.DeclarationSpecifiers; ds != nil; ds = ds.DeclarationSpecifiers {
		if ds.StorageClassSpecifier != nil && ds.StorageClassSpecifier.Case == cc.StorageClassSpecifierTypedef {
			return true
		}
	}

	return false
}

// ccFirstDeclaratorName returns the first declared name, which for
// "typedef struct { ... } AVTimecode;" is the name the struct is registered
// under.
func ccFirstDeclaratorName(d *cc.Declaration) string {
	for idl := d.InitDeclaratorList; idl != nil; idl = idl.InitDeclaratorList {
		if idl.InitDeclarator == nil || idl.InitDeclarator.Declarator == nil {
			continue
		}

		if name := idl.InitDeclarator.Declarator.Name(); name != "" {
			return name
		}
	}

	return ""
}

func (w *ccWalk) declarationSpecifiers(ds *cc.DeclarationSpecifiers, fallback string, isTypedef bool) {
	for ; ds != nil; ds = ds.DeclarationSpecifiers {
		if ds.Case != cc.DeclarationSpecifiersTypeSpec || ds.TypeSpecifier == nil {
			continue
		}

		ts := ds.TypeSpecifier

		switch ts.Case {
		case cc.TypeSpecifierStructOrUnion:
			if sus := ts.StructOrUnionSpecifier; sus != nil && sus.Case == cc.StructOrUnionSpecifierDef {
				w.recordStruct(sus, fallback, isTypedef)
			}

		case cc.TypeSpecifierEnum:
			if es := ts.EnumSpecifier; es != nil && es.Case == cc.EnumSpecifierDef {
				w.recordEnum(es, fallback, isTypedef)
			}
		}
	}
}

// recordTags records the "StructOrUnion IDENTIFIER" form: a tag with no body.
// A bare enum tag is never recorded, because a forward-declared enum has no
// members and an unnamed enum becomes loose constants instead.
//
// embedded says the declaration carries declarators, so the tag is written
// inside another declaration's decl-specifier-seq rather than standing alone.
// Such a tag claims no comment: clang's getDeclLocForCommentSearch returns no
// search location for a TagDecl that is embedded in a declarator and is not a
// complete definition, so "typedef struct AVBSFList AVBSFList;" gives the block
// above it to the typedef alone, and "struct AVMurMur3 *av_murmur3_alloc(void);"
// gives it to the function. A bare "struct AVCodecTag;" is not embedded and does
// claim one.
func (w *ccWalk) recordTags(ds *cc.DeclarationSpecifiers, isTypedef, embedded bool) {
	for ; ds != nil; ds = ds.DeclarationSpecifiers {
		if ds.Case != cc.DeclarationSpecifiersTypeSpec || ds.TypeSpecifier == nil {
			continue
		}

		ts := ds.TypeSpecifier
		if ts.Case != cc.TypeSpecifierStructOrUnion || ts.StructOrUnionSpecifier == nil {
			continue
		}

		sus := ts.StructOrUnionSpecifier
		if sus.Case != cc.StructOrUnionSpecifierTag {
			continue
		}

		name := sus.Token.SrcStr()
		if name == "" {
			continue
		}

		// A tag is declared once, where it enters scope. A later
		// "struct AVCodecTag *" in a return type or a parameter only refers to
		// the tag already declared, so it must not move the entry to the end of
		// structOrder.
		if w.tags[name] {
			continue
		}

		w.tags[name] = true

		if isTypedef {
			// A typedef over a bare tag resolves to the tag's definition,
			// wherever it sits in the file, because a tag type names its whole
			// redeclaration chain. So "typedef struct AVExifEntry AVExifEntry;"
			// registers the fields at the typedef, and the definition that
			// follows hits putStruct's already-has-fields guard instead of
			// re-appending.
			if def := ccTagDefinition(sus); def != nil {
				w.recordTagDefinition(def)

				continue
			}
		} else if ccIsUnion(sus) {
			// A top-level union is not recorded as a struct. Only a union
			// reached through a typedef is.
			continue
		}

		if !w.declaredHere(sus.Token) {
			continue
		}

		comment := ""
		if !embedded {
			comment = w.comment(sus.Token, 0)
		}

		w.putStruct(name, &Struct{Name: name, Typedefd: isTypedef, Comment: comment})
	}
}

// ccTagDefinition returns the specifier that defines the referenced tag in this
// translation unit, or nil when the tag stays opaque.
func ccTagDefinition(sus *cc.StructOrUnionSpecifier) *cc.StructOrUnionSpecifier {
	name := sus.Token.SrcStr()
	if name == "" {
		return nil
	}

	for s := sus.LexicalScope(); s != nil; s = s.Parent {
		for _, n := range s.Nodes[name] {
			def, ok := n.(*cc.StructOrUnionSpecifier)
			if ok && def.Case == cc.StructOrUnionSpecifierDef {
				return def
			}
		}
	}

	return nil
}

// recordTagDefinition records the record a typedef's tag resolves to. It runs no
// location filter: the filter has already been applied to the typedef itself.
//
// The record is reached through a typedef, so it is Typedefd. This is the only
// caller, and it is the typedef branch of recordTags.
func (w *ccWalk) recordTagDefinition(def *cc.StructOrUnionSpecifier) {
	name := def.Token.SrcStr()
	if name == "" {
		return
	}

	w.tags[name] = true

	// The comment comes from the definition, not from the typedef that reached
	// it. AVExifEntry and AVFilterLink are the two records this decides, and the
	// comment golden pins both.
	w.putStruct(name, &Struct{
		Name:     name,
		Typedefd: true,
		Comment:  w.comment(def.Token, w.spanFirstLine(def)),
		Fields:   w.structFields(name, def),
	})
}

// recordStruct records a struct or union definition and its field names.
func (w *ccWalk) recordStruct(sus *cc.StructOrUnionSpecifier, fallback string, isTypedef bool) {
	// A definition declares its tag, so a tag reference later in the same
	// translation unit is no longer an introduction.
	if tag := sus.Token.SrcStr(); tag != "" {
		w.tags[tag] = true
	}

	if !isTypedef && ccIsUnion(sus) {
		return
	}

	name := sus.Token.SrcStr()

	if sus.Token.Seq() == 0 {
		// A tagless record that is the specifier of a typedef takes the typedef
		// name. Otherwise it is a nested anonymous record, which the type walk
		// registers under a synthetic name.
		if fallback == "" {
			return
		}

		name = fallback
	}

	if name == "" || !w.declaredHere(ccStructKeyword(sus)) {
		return
	}

	w.putStruct(name, &Struct{
		Name:     name,
		Typedefd: isTypedef,
		Comment:  w.comment(sus.Token, w.spanFirstLine(sus)),
		Fields:   w.structFields(name, sus),
	})
}

// structFields lists the named fields a record definition declares, with the
// parsed type, the written C spelling and the declared bit width of each.
func (w *ccWalk) structFields(owner string, sus *cc.StructOrUnionSpecifier) []*Field {
	var fields []*Field

	ccEachFieldDeclarator(sus, func(sd *cc.StructDeclarator, anon *cc.StructOrUnionSpecifier) {
		dcl := sd.Declarator

		prev := w.anon
		w.anon = anon

		fields = append(fields, &Field{
			Name:      dcl.Name(),
			CTypeName: ccCTypeName(dcl, anon),
			BitWidth:  ccFieldBitWidth(sd),
			Comment:   w.comment(dcl.NameTok(), 0),
			Type:      w.parseType(fmt.Sprintf("[%v][%v]", owner, dcl.Name()), dcl.Type()),
		})

		w.anon = prev
	})

	return fields
}

// ccEachFieldDeclarator calls fn for every named field declarator of a record
// definition, with the tagless struct or union its declaration writes inline as
// the base type, or nil.
//
// The anonymous specifier is per declaration rather than per declarator, so it
// is resolved once and handed to each name that shares it.
func ccEachFieldDeclarator(sus *cc.StructOrUnionSpecifier, fn func(sd *cc.StructDeclarator, anon *cc.StructOrUnionSpecifier)) {
	for sdl := sus.StructDeclarationList; sdl != nil; sdl = sdl.StructDeclarationList {
		sd := sdl.StructDeclaration
		if sd == nil || sd.Case != cc.StructDeclarationDecl {
			continue
		}

		// A tagless struct or union written inline as the base type is named and
		// spelled from its syntax node, which the field's type cannot reach.
		anon := ccAnonSpecifier(sd)

		for l := sd.StructDeclaratorList; l != nil; l = l.StructDeclaratorList {
			if l.StructDeclarator == nil || l.StructDeclarator.Declarator == nil {
				continue
			}

			if l.StructDeclarator.Declarator.Name() == "" {
				continue
			}

			fn(l.StructDeclarator, anon)
		}
	}
}

// ccFieldBitWidth returns the declared width of a bit field, and -1 for a field
// that is not one. That -1 is the encoding the IR pins, so a plain field must
// never read 0.
//
// The width comes from the syntax node rather than from cc.Field.ValueBits,
// which reports a storage width for every field whether it is a bit field or
// not.
func ccFieldBitWidth(sd *cc.StructDeclarator) int32 {
	if sd.Case != cc.StructDeclaratorBitField || sd.ConstantExpression == nil {
		return -1
	}

	v, ok := sd.ConstantExpression.Value().(cc.Int64Value)
	if !ok {
		return -1
	}

	// C caps a bit field at the width of its declared type, so a value outside
	// int32 cannot come from a header this generator can bind. Truncating it
	// would emit a plausible width for a field the parser did not understand, so
	// abort instead, matching how the walk treats a duplicate registration.
	w := int64(v)
	if w < math.MinInt32 || w > math.MaxInt32 {
		log.Panicln("bit field width", w, "out of int32 range at", sd.Position())

		return -1
	}

	return int32(w)
}

// putStruct stores a struct, keeping the entry that carries fields. An opaque
// forward declaration must never displace a definition.
//
// The kept entry still takes Typedefd from the incoming record, so a struct
// defined first and typedef'd afterwards ends up Typedefd all the same.
func (w *ccWalk) putStruct(name string, s *Struct) {
	if existing, ok := w.mod.structs[name]; ok && len(existing.Fields) > 0 {
		if s.Typedefd {
			existing.Typedefd = true
		}

		return
	}

	w.mod.structOrder = slices.DeleteFunc(w.mod.structOrder, func(n string) bool {
		return n == name
	})

	w.mod.structs[name] = s
	w.mod.structOrder = append(w.mod.structOrder, name)
}

// recordEnum records a named enum and its members. An unnamed enum has no
// members of its own in the Module: it has no name the emitter can refer to, so
// its enumerators become loose constants instead.
func (w *ccWalk) recordEnum(es *cc.EnumSpecifier, fallback string, isTypedef bool) {
	name := es.Token2.SrcStr()

	if es.Token2.Seq() == 0 {
		if fallback == "" {
			if w.declaredHere(es.Token) {
				w.recordAnonEnumConstants(es)
			}

			return
		}

		name = fallback
	}

	if name == "" || !w.declaredHere(es.Token) {
		return
	}

	// Typedefd is set on the kept entry too, so an enum defined first and
	// typedef'd afterwards is still Typedefd.
	if existing, ok := w.mod.enums[name]; ok && len(existing.Constants) > 0 {
		if isTypedef {
			existing.Typedefd = true
		}

		return
	}

	e := &Enum{
		Name:     name,
		Typedefd: isTypedef,
		Comment:  w.comment(es.Token2, w.spanFirstLine(es)),
	}

	for el := es.EnumeratorList; el != nil; el = el.EnumeratorList {
		en := el.Enumerator
		if en == nil || en.Token.SrcStr() == "" {
			continue
		}

		e.Constants = append(e.Constants, &Constant{
			Name:    en.Token.SrcStr(),
			Comment: w.comment(en.Token, 0),
		})
	}

	w.mod.enumOrder = slices.DeleteFunc(w.mod.enumOrder, func(n string) bool {
		return n == name
	})

	w.mod.enums[name] = e
	w.mod.enumOrder = append(w.mod.enumOrder, name)
}

// recordAnonEnumConstants appends the members of an unnamed enum to the
// constant table, marked as enum-sourced.
func (w *ccWalk) recordAnonEnumConstants(es *cc.EnumSpecifier) {
	for el := es.EnumeratorList; el != nil; el = el.EnumeratorList {
		en := el.Enumerator
		if en == nil {
			continue
		}

		name := en.Token.SrcStr()
		if name == "" || !w.declaredHere(en.Token) {
			continue
		}

		w.mod.constants[name] = &Constant{Name: name, FromEnum: true}
		w.mod.constantOrder = append(w.mod.constantOrder, name)
	}
}

// recordFunction records a function declarator.
//
// A duplicate registration is fatal: it means the walk visited one declaration
// twice, so every binding downstream is suspect.
func (w *ccWalk) recordFunction(d *cc.Declarator) {
	if d == nil {
		return
	}

	ft, ok := d.Type().(*cc.FunctionType)
	if !ok {
		return
	}

	name := d.Name()
	if name == "" || !w.declaredHere(d.NameTok()) {
		return
	}

	if _, ok := w.mod.functions[name]; ok {
		log.Panicln("Function", name, "already exists")
	}

	w.mod.functions[name] = &Function{
		Name:     name,
		Result:   w.parseType(fmt.Sprintf("[%v][return]", name), ft.Result()),
		Args:     w.params(name, ft),
		Variadic: ft.IsVariadic(),
		Comment:  w.comment(d.NameTok(), 0),
	}
	w.mod.functionOrder = append(w.mod.functionOrder, name)
}

// params walks a function type's parameter list.
//
// A parameter declared without a name takes "param" followed by its index,
// counted over the parameter list itself, so one rule covers functions and
// callbacks alike.
//
// The name is read off the written declarator, so the decay ccParamType has to
// reverse for the type does not reach it.
func (w *ccWalk) params(owner string, ft *cc.FunctionType) []*Param {
	var args []*Param

	for i, p := range ccParameters(ft) {
		name := p.Name()
		if name == "" {
			name = fmt.Sprintf("param%v", i)
		}

		args = append(args, &Param{
			Name: name,
			Type: w.parseType(fmt.Sprintf("[%v][%v]", owner, name), ccParamType(p.Type())),
		})
	}

	return args
}

// recordCallbacks records the function typedefs of a typedef declaration.
//
// A typedef reaches mod.callbacks only when it declares parameters, so a
// function typedef taking no arguments is not a callback.
func (w *ccWalk) recordCallbacks(d *cc.Declaration) {
	for idl := d.InitDeclaratorList; idl != nil; idl = idl.InitDeclaratorList {
		if idl.InitDeclarator == nil || idl.InitDeclarator.Declarator == nil {
			continue
		}

		dcl := idl.InitDeclarator.Declarator

		ft, ptr := ccFuncType(dcl.Type())
		if ft == nil || len(ft.Parameters()) == 0 {
			continue
		}

		name := dcl.Name()
		if name == "" || !w.declaredHere(dcl.NameTok()) {
			continue
		}

		if _, ok := w.mod.callbacks[name]; ok {
			log.Panicln("callback already exists")
		}

		w.mod.callbacks[name] = &Function{
			Name:   name,
			Result: w.parseType(fmt.Sprintf("[%v][return]", name), ft.Result()),
			Args:   w.params(name, ft),
			Ptr:    ptr,
		}
		w.mod.callbackOrder = append(w.mod.callbackOrder, name)
	}
}

// macros appends the object-like macros a header defines to the constant
// table, in source order.
//
// ast.Macros is a map, so the entries are sorted by position rather than
// ranged over directly. The name filter is ccMacroSkipped.
func (w *ccWalk) macros(ast *cc.AST) {
	type entry struct {
		name string
		line int
		col  int
	}

	var entries []entry

	for name, m := range ast.Macros {
		if m == nil || m.IsFnLike || ccMacroSkipped(name) {
			continue
		}

		// __FILE__ and __LINE__ are dynamic builtins and are not constants the
		// bindings can carry. cc/v4 records them in ast.Macros at the position
		// of their last expansion, which lands inside an av_assert0 in the
		// header being parsed, so without this test they enter the table.
		if name == "__FILE__" || name == "__LINE__" {
			continue
		}

		p := m.Name.Position()
		if p.Filename != w.abs {
			continue
		}

		entries = append(entries, entry{name: name, line: p.Line, col: p.Column})
	}

	slices.SortFunc(entries, func(a, b entry) int {
		if a.line != b.line {
			return a.line - b.line
		}

		return a.col - b.col
	})

	for _, e := range entries {
		w.mod.constants[e.name] = &Constant{Name: e.name}
		w.mod.constantOrder = append(w.mod.constantOrder, e.name)
	}
}

// ccMacroSkipped is the macro-name filter: header guards, version
// and identification strings, the channel-layout compound literals, the
// compiler attribute macros and the math and utility function macros.
func ccMacroSkipped(name string) bool {
	switch {
	case strings.HasSuffix(name, "_H"),
		strings.HasPrefix(name, "av_malloc_"),
		strings.HasPrefix(name, "AV_CHANNEL_LAYOUT_"),
		strings.HasSuffix(name, "_VERSION"),
		strings.HasSuffix(name, "_IDENT"),
		name == "AV_TIME_BASE_Q",
		name == "AV_CH_LAYOUT_NATIVE":
		return true
	}

	if strings.HasPrefix(name, "av_") && (strings.HasSuffix(name, "_inline") ||
		strings.Contains(name, "_unused") ||
		strings.Contains(name, "_deprecated") ||
		strings.Contains(name, "_result") ||
		name == "av_const" || name == "av_pure" ||
		name == "av_cold" || name == "av_flatten" ||
		name == "av_alias" || name == "av_noreturn" ||
		name == "av_used" || name == "av_noinline") {
		return true
	}

	return strings.HasPrefix(name, "av_clip") ||
		strings.HasPrefix(name, "av_sat_") ||
		strings.HasPrefix(name, "av_popcount") ||
		strings.HasPrefix(name, "av_ceil_log2") ||
		strings.HasPrefix(name, "av_mod_") ||
		strings.HasPrefix(name, "av_zero_extend") ||
		strings.HasPrefix(name, "av_parity") ||
		strings.HasPrefix(name, "FF_CEIL_RSHIFT") ||
		strings.HasPrefix(name, "av_builtin_")
}

// declaredHere reports whether a token sits in the header being translated. It
// is the location filter: every header is preprocessed with its includes, so
// without it glibc's __mbstate_t and atomic_wide_counter are recorded ahead of
// the FFmpeg types.
func (w *ccWalk) declaredHere(tok cc.Token) bool {
	if tok.Seq() == 0 {
		return false
	}

	return tok.Position().Filename == w.abs
}

// ccIsUnion reports whether the specifier is a union rather than a struct.
func ccIsUnion(sus *cc.StructOrUnionSpecifier) bool {
	return sus.StructOrUnion != nil && sus.StructOrUnion.Case == cc.StructOrUnionUnion
}

// ccStructKeyword returns the "struct" or "union" keyword token, which is the
// one position a tagless record always has.
func ccStructKeyword(sus *cc.StructOrUnionSpecifier) cc.Token {
	if sus.StructOrUnion == nil {
		return sus.Token
	}

	return sus.StructOrUnion.Token
}

// ccFuncType returns the function type a declarator declares, unwrapping one
// pointer level so both "typedef int f(int)" and "typedef int (*f)(int)"
// resolve. It returns nil for anything else.
//
// The second result is Function.Ptr: whether the pointer level was there. An
// invalid pointee counts as "not a pointer".
func ccFuncType(t cc.Type) (*cc.FunctionType, bool) {
	if t == nil {
		return nil, false
	}

	ptr := false

	if p, ok := t.(*cc.PointerType); ok {
		t = p.Elem()
		ptr = true
	}

	ft, ok := t.(*cc.FunctionType)
	if !ok {
		return nil, false
	}

	return ft, ptr
}

// comment returns the doc comment claimed for a declared name (WW-141 phase
// 4.2). tok is the name token and fallbackLine is the declaration's first line,
// used only when the name token is absent.
//
// A fallbackLine of 0 says the caller guarantees the token: an absent token then
// claims nothing rather than silently matching against line 0, which the leading
// pass would answer with the last comment in the file.
func (w *ccWalk) comment(tok cc.Token, fallbackLine int) string {
	line, col := w.nameSpan(tok, fallbackLine)
	if line == 0 {
		return ""
	}

	return w.table.CommentAt(line, col)
}

// nameSpan returns the line and exclusive end column of a name token, when the
// token exists and sits in the header being translated.
//
// Both claim passes anchor here, because clang searches from
// Decl->getLocation(), which for a C declaration is the identifier rather than
// the start or the end of the declaration.
//
// A tagless struct or enum has no name token, so it falls back to the start of
// the declaration, which is where an unnamed record is reported. Absence is
// Token.Seq() == 0 and not an empty Token.SrcStr(), because a token can carry no
// source text and still be a real position.
// A name written inside a macro invocation reports the invocation's position
// rather than its own, so the written position is recovered by spellingSite
// before either pass runs.
func (w *ccWalk) nameSpan(tok cc.Token, fallbackLine int) (line, col int) {
	p := tok.Position()

	if tok.Seq() == 0 || p.Filename != w.abs {
		return fallbackLine, 1
	}

	name := tok.SrcStr()

	line, col = w.spellingSite(name, p.Line, p.Column)
	if line == 0 {
		return 0, 0
	}

	return line, col + len(name)
}

// spellingSite returns the position at which name is written, given the
// position cc/v4 reports for it. It returns line 0 when the name is written in
// a header other than the one being translated, because the comment table
// cannot see that header and the comment to claim would be in it.
//
// cc/v4 stamps every token a macro expands to with the invocation token's
// offset, so a declaration written inside a macro invocation collapses onto the
// invocation's line and column. A macro argument has to go back to its own
// position, or the FF_PAD_STRUCTURE fields all claim the comment above the
// invocation instead of their own.
//
// The test is that the source at the reported position does not spell the name.
// It is the invocation's own identifier that is written there instead, so the
// name is looked for in three places, in the order clang resolves them:
//
//  1. the invocation's argument text, which is where a macro argument is
//     written, and which gives the AVBPrint fields their own trailing comments;
//  2. the macro's replacement list, which is where a name the macro body
//     declares is written, and which puts AVBPrint.reserved_padding next to the
//     "#define" that a leading comment can never cross;
//  3. the invocation itself, unchanged, for a name no source position spells,
//     which is what "##" pasting produces. clang gives such a token a scratch
//     location, and struct ff_pad_helper_AVBPrint claims the comment above its
//     invocation.
func (w *ccWalk) spellingSite(name string, line, col int) (int, int) {
	if w.table.IdentifierAt(line, col) == name {
		return line, col
	}

	m := w.macroTable[w.table.IdentifierAt(line, col)]
	if m == nil {
		return line, col
	}

	if l, c, ok := w.table.MacroArgumentSite(name, line, col); ok {
		return l, c
	}

	for _, tok := range m.ReplacementList() {
		if tok.SrcStr() != name {
			continue
		}

		p := tok.Position()
		if p.Filename != w.abs {
			return 0, 0
		}

		return p.Line, p.Column
	}

	return line, col
}

// spanFirstLine returns the first line a node occupies, counting only tokens in
// the header being translated. Without that file filter a node reached through
// an #include reports a line in the wrong file, which is the same filter
// declaredHere applies to a single token.
//
// The line is the minimum token line rather than n.Position().Line, because a
// macro-expanded __attribute__ reports the invocation site, which can sit a line
// above the position the node itself reports. 20 of the 1375 top-level
// declarations disagree, every one for that reason.
//
// Only the first line is returned. The claim passes anchor on the name, so the
// end of the declaration's span is never asked for.
func (w *ccWalk) spanFirstLine(n cc.Node) int {
	toks := cc.NodeTokens(n)

	first := 0

	for i := range toks {
		p := toks[i].Position()
		if p.Filename != w.abs {
			continue
		}

		if first == 0 || p.Line < first {
			first = p.Line
		}
	}

	return first
}
