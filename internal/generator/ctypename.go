package main

// C type spelling for the modernc.org/cc/v4 parser backend (WW-141 phase 3.2).
//
// Field.CTypeName holds what clang's QualType::getAsString() printed for the
// written type of a struct field. struct.go reads that spelling to decide how a
// field is wrapped, and the IR golden pins all 1335 of them, so the string has
// to come back byte for byte off a parser that has no such printer.
//
// The rules, all measured against the golden over the 120 pinned headers
// (WW-141 probe answer 3, 1335 matched, 0 mismatched):
//
//   - The written type wins over the canonical one. Typedef sugar is taken from
//     Type.Typedef() and stops the walk, which is what keeps uint8_t, size_t and
//     AVRational out of the widths they canonicalise to. Type.String() is
//     documented as unspecified and is not used.
//   - Pointer-level const comes from the declarator syntax, never from
//     Type.Attributes().IsConst(). See ccPointerQualifiers.
//   - A tagless struct or union needs its syntax node, because the type carries
//     no tag token and nothing on it reaches the "struct" keyword.
//   - A nested type, meaning a function-pointer parameter, keeps its const from
//     the type tree and is not trimmed.

import (
	"fmt"
	"strings"

	"modernc.org/cc/v4"
)

// ccPredefinedNames maps a cc/v4 predefined kind to the spelling clang's
// QualType::getAsString produces for it.
var ccPredefinedNames = map[cc.Kind]string{
	cc.Void:       "void",
	cc.Bool:       "_Bool",
	cc.Char:       "char",
	cc.SChar:      "signed char",
	cc.UChar:      "unsigned char",
	cc.Short:      "short",
	cc.UShort:     "unsigned short",
	cc.Int:        "int",
	cc.UInt:       "unsigned int",
	cc.Long:       "long",
	cc.ULong:      "unsigned long",
	cc.LongLong:   "long long",          //nolint:dupword // "long long" is the C type name, not a repeated word
	cc.ULongLong:  "unsigned long long", //nolint:dupword // "long long" is the C type name, not a repeated word
	cc.Float:      "float",
	cc.Double:     "double",
	cc.LongDouble: "long double",
	cc.Int128:     "__int128",
	cc.UInt128:    "unsigned __int128",
}

// ccCTypeName spells the written type of one struct field. anon is the tagless
// struct or union the field declaration writes inline as its base type, or nil.
//
// The result is trimmed and loses one leading "const ", because the IR spells a
// const-qualified base type without it.
func ccCTypeName(d *cc.Declarator, anon *cc.StructOrUnionSpecifier) string {
	if d == nil {
		return ""
	}

	s := ccTypeSpeller{anon: anon, ptrQuals: ccPointerQualifiers(d), fromDeclarator: true}

	return strings.TrimPrefix(strings.TrimSpace(s.name(d.Type())), "const ")
}

// ccPointerQualifiers lists the qualifier keywords written on each '*' of a
// declarator, left to right in source order.
//
// The type tree cannot supply this. cc/v4 drops the const of
// "const char *const *" outright and hoists the base const of
// "const uint8_t (*)[9]" onto the pointer, so reading
// Type.Attributes().IsConst() per level spells both wrongly. The declarator
// carries a TypeQualifiers list on every '*' and carries both correctly.
//
// A parenthesised declarator holds its '*' on a nested Declarator rather than
// on this one, which is where "int (*const table)[4]" keeps its const, so the
// walk takes the outer Pointer chain first and then recurses into that nested
// declarator, appending what it carries after. Both steps run left to right, so
// the whole list stays in source order.
//
// The leftmost '*' is the innermost derivation, so index i names the pointer
// level len-1-i counted from the outermost. Base-level const is not collected at
// all, because ccCTypeName trims a leading "const " and it could never survive
// into the golden.
func ccPointerQualifiers(d *cc.Declarator) []string {
	var quals []string

	for ; d != nil; d = ccParenthesisedDeclarator(d.DirectDeclarator) {
		for p := d.Pointer; p != nil; p = p.Pointer {
			quals = append(quals, ccQualifierNames(p.TypeQualifiers))
		}
	}

	return quals
}

// ccParenthesisedDeclarator returns the declarator a '(' Declarator ')' form
// wraps, following the DirectDeclarator chain down to it past any array or
// function derivation, or nil when the declarator writes no parentheses.
func ccParenthesisedDeclarator(dd *cc.DirectDeclarator) *cc.Declarator {
	for ; dd != nil; dd = dd.DirectDeclarator {
		if dd.Case == cc.DirectDeclaratorDecl {
			return dd.Declarator
		}
	}

	return nil
}

// ccQualifierNames flattens a TypeQualifiers list to its keywords.
func ccQualifierNames(q *cc.TypeQualifiers) string {
	var out []string

	for ; q != nil; q = q.TypeQualifiers {
		if q.TypeQualifier != nil {
			out = append(out, q.TypeQualifier.Token.SrcStr())
		}
	}

	return strings.Join(out, " ")
}

// ccTypeSpeller renders one written type.
//
// fromDeclarator selects where const comes from: the ptrQuals collected off the
// declarator when set, the type tree when not. A nested type has no declarator
// of its own, so it takes the second route.
type ccTypeSpeller struct {
	anon           *cc.StructOrUnionSpecifier
	ptrQuals       []string
	fromDeclarator bool
}

// name renders a type with an empty declarator name, which is the form
// QualType::getAsString produces.
func (s ccTypeSpeller) name(t cc.Type) string {
	return ccJoinDeclarator(s.split(t))
}

// ccJoinDeclarator puts a base spelling and a declarator back together. An array
// declarator abuts its base ("int[3]"); anything else takes a space
// ("char *", "void (*)(int)").
func ccJoinDeclarator(base, decl string) string {
	switch {
	case decl == "":
		return base
	case decl[0] == '[':
		return base + decl
	default:
		return base + " " + decl
	}
}

// split walks the type outward-in, accumulating the C declarator, and returns
// the base spelling plus that declarator. A pointer wrapped by an array or a
// function is parenthesised, which is C declarator syntax and is what produces
// "uint8_t (*)[9]" and "void (*)(int16_t *)".
func (s ccTypeSpeller) split(t cc.Type) (base, decl string) {
	level := 0

	for t != nil {
		if td := t.Typedef(); td != nil && td.Name() != "" {
			return s.baseQual(t) + td.Name(), decl
		}

		switch x := t.(type) {
		case *cc.PointerType:
			decl = "*" + s.pointerQual(level, x) + decl
			t = x.Elem()
			level++

		case *cc.ArrayType:
			decl = ccParenPointer(decl)

			if x.IsIncomplete() {
				decl += "[]"
			} else {
				decl += fmt.Sprintf("[%d]", x.Len())
			}

			t = x.Elem()

		case *cc.FunctionType:
			decl = ccParenPointer(decl)

			params := make([]string, 0, len(x.Parameters()))
			for _, p := range x.Parameters() {
				params = append(params, ccNestedTypeName(p.Type()))
			}

			if x.IsVariadic() {
				params = append(params, "...")
			}

			decl += "(" + strings.Join(params, ", ") + ")"

			t = x.Result()

		default:
			return s.baseQual(t) + ccBaseTypeName(t, s.anon), decl
		}
	}

	return "void", decl
}

// ccParenPointer wraps a pointer declarator in the parentheses C needs when an
// array or a function derives from it, and leaves anything else alone. A
// pointer qualifier ends in a space because a further '*' would follow it, so
// the wrap drops it again: clang spells "int (*const)[4]", not "int (*const )[4]".
func ccParenPointer(decl string) string {
	if !strings.HasPrefix(decl, "*") {
		return decl
	}

	return "(" + strings.TrimRight(decl, " ") + ")"
}

// ccNestedTypeName renders a function-pointer parameter type. It reads const off
// the type tree, because a parameter of a function type reached through a field
// has no declarator to read, and it does not trim: the leading const of
// "const enum AVPixelFormat *" is part of the spelling here.
func ccNestedTypeName(t cc.Type) string {
	return ccTypeSpeller{}.name(t)
}

// pointerQual spells the const of one pointer level.
func (s ccTypeSpeller) pointerQual(level int, t cc.Type) string {
	if !s.fromDeclarator {
		return ccConstQual(t)
	}

	i := len(s.ptrQuals) - 1 - level
	if i >= 0 && i < len(s.ptrQuals) && strings.Contains(s.ptrQuals[i], "const") {
		return "const "
	}

	return ""
}

// baseQual spells the const of a typedef name or a base type. In declarator mode
// it is always empty: ccCTypeName trims a leading "const " away, so the base
// qualifier can never reach the result.
func (s ccTypeSpeller) baseQual(t cc.Type) string {
	if s.fromDeclarator {
		return ""
	}

	return ccConstQual(t)
}

// ccConstQual reads const off the type tree.
func ccConstQual(t cc.Type) string {
	if t != nil && t.Attributes().IsConst() {
		return "const "
	}

	return ""
}

// ccBaseTypeName spells a type that derives from nothing: a record, an enum or a
// predefined type. anon supplies the syntax node of a tagless record, because
// StructType.Tag() is the zero Token there and the type reaches neither the
// keyword nor a position.
func ccBaseTypeName(t cc.Type, anon *cc.StructOrUnionSpecifier) string {
	switch x := t.(type) {
	case *cc.StructType:
		tok := x.Tag()
		if tag := tok.SrcStr(); tag != "" {
			return "struct " + tag
		}

		return ccUnnamedRecordName("struct", anon)

	case *cc.UnionType:
		tok := x.Tag()
		if tag := tok.SrcStr(); tag != "" {
			return "union " + tag
		}

		return ccUnnamedRecordName("union", anon)

	case *cc.EnumType:
		tok := x.Tag()
		if tag := tok.SrcStr(); tag != "" {
			return "enum " + tag
		}

		return "enum (unnamed enum)"

	case *cc.PredefinedType:
		if name, ok := ccPredefinedNames[x.Kind()]; ok {
			return name
		}

		return "PREDEFINED?" + x.Kind().String()
	}

	if t == nil {
		return "void"
	}

	// An unhandled implementation must be loud in the golden diff rather than
	// spelling as something plausible, so name the Go type.
	return "UNHANDLED?" + fmt.Sprintf("%T", t)
}

// ccUnnamedRecordName reproduces clang's
// "union (unnamed union at <absolute header path>:<line>:<col>)" form from the
// position of the struct or union keyword.
//
// The path is absolute on purpose. normalizeCTypeName (dump.go) cuts it back to
// "include/..." for the golden, so emitting a relative path here would leave it
// nothing to cut and the golden would move.
func ccUnnamedRecordName(kw string, anon *cc.StructOrUnionSpecifier) string {
	if anon == nil || anon.StructOrUnion == nil {
		return fmt.Sprintf("%s (unnamed %s at ?)", kw, kw)
	}

	p := anon.StructOrUnion.Token.Position()

	return fmt.Sprintf("%s (unnamed %s at %s:%d:%d)", kw, kw, p.Filename, p.Line, p.Column)
}

// ccAnonSpecifier returns the tagless struct or union specifier a struct
// declaration writes inline as its base type, if any. It is the anon argument
// ccCTypeName needs.
//
// Taglessness is tested with Token.Seq(), not Token.SrcStr(): the tag token of a
// tagless record is the zero Token, whose SrcStr is empty but so is that of a
// token the lexer never produced.
func ccAnonSpecifier(sd *cc.StructDeclaration) *cc.StructOrUnionSpecifier {
	for sql := sd.SpecifierQualifierList; sql != nil; sql = sql.SpecifierQualifierList {
		if sql.Case != cc.SpecifierQualifierListTypeSpec || sql.TypeSpecifier == nil {
			continue
		}

		sus := sql.TypeSpecifier.StructOrUnionSpecifier
		if sus != nil && sus.Case == cc.StructOrUnionSpecifierDef && sus.Token.Seq() == 0 {
			return sus
		}
	}

	return nil
}
