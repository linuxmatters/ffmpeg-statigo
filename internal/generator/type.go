package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/go-clang/bootstrap/clang"
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

func (p *Parser) parseType(indent string, t clang.Type) Type {
	switch t.Kind() {
	case clang.Type_Void:
		log.Println(indent, "Parsing type", t.Spelling(), "as void")
		return nil

	case clang.Type_Int, clang.Type_UInt, clang.Type_Long, clang.Type_ULong, clang.Type_UChar, clang.Type_Char_S,
		clang.Type_Float, clang.Type_Double, clang.Type_Enum, clang.Type_Record:
		log.Println(indent, "Parsing type", t.Spelling(), "as ident")
		name := t.CanonicalType().Spelling()
		name = strings.TrimPrefix(name, "const ")
		name = strings.TrimPrefix(name, "struct ")
		name = strings.TrimPrefix(name, "enum ")

		if strings.HasPrefix(name, "unsigned ") {
			name = strings.TrimPrefix(name, "unsigned ")
			name = fmt.Sprintf("u%v", name)
		}

		if name == "" {
			log.Panicln(indent, "No name")
		}

		// Handle unnamed structs/unions by generating a synthetic name
		if strings.Contains(name, " ") {
			if strings.HasPrefix(name, "(unnamed") || strings.HasPrefix(name, "(anonymous") {
				syntheticName := p.generateUnnamedStructName(t)
				log.Println(indent, "Generated synthetic name for unnamed struct:", syntheticName)

				// Register the unnamed struct so it gets generated
				p.registerUnnamedStruct(indent, t, syntheticName)

				return &IdentType{
					Name: syntheticName,
				}
			}

			log.Panicln(indent, "name", name, "contains space")
		}

		return &IdentType{
			Name: name,
		}

	case clang.Type_Typedef:
		log.Println(indent, "Parsing type", t.Spelling(), "as typedef")

		name := t.DefName()
		if name == "" {
			log.Panicln(indent, "Unexpected def name", name)
		}

		return &IdentType{
			Name: t.DefName(),
		}

	case clang.Type_Pointer:
		log.Println(indent, "Parsing type", t.Spelling(), "as pointer")

		inner := p.parseType(indent, t.PointeeType())

		return &PointerType{
			Inner: inner,
		}

	case clang.Type_Elaborated:
		log.Println(indent, "Parsing type", t.Spelling(), "as elaborated")

		dec := t.Declaration()

		switch dec.Kind() {
		case clang.Cursor_EnumDecl:
			name := dec.Spelling()

			if name == "" {
				log.Panicln(indent, "No name")
			}

			if strings.HasPrefix(name, "enum (unnamed") {
				dec.Visit(func(cursor, parent clang.Cursor) (status clang.ChildVisitResult) {
					if cursor.Kind() != clang.Cursor_EnumConstantDecl {
						log.Panicln("Unknown enum type", "kind", cursor.Kind().String())
					}

					name := cursor.Spelling()

					p.mod.constants[name] = &Constant{
						Name:     name,
						FromEnum: true,
					}
					p.mod.constantOrder = append(p.mod.constantOrder, name)

					return clang.ChildVisit_Continue
				})

				return &IdentType{Name: "uint32_t", IsAnonEnum: true}
			}

			if strings.Contains(name, " ") {
				log.Panicln(indent, "name", name, "contains space")
			}

			return &IdentType{Name: name}

		case clang.Cursor_StructDecl, clang.Cursor_TypedefDecl:
			name := dec.Spelling()

			if name == "" {
				log.Panicln(indent, "No name")
			}

			// Handle unnamed structs
			if strings.Contains(name, " ") {
				if strings.HasPrefix(name, "struct (unnamed") || strings.HasPrefix(name, "union (unnamed") {
					syntheticName := p.generateUnnamedStructName(t)
					log.Println(indent, "Generated synthetic name for unnamed struct:", syntheticName)

					// Register the unnamed struct so it gets generated
					p.registerUnnamedStruct(indent, t, syntheticName)

					return &IdentType{
						Name: syntheticName,
					}
				}

				log.Panicln(indent, "name", name, "contains space")
			}

			return &IdentType{Name: name}

		case clang.Cursor_UnionDecl:
			u := &UnionType{}

			dec.Visit(func(cursor, parent clang.Cursor) (status clang.ChildVisitResult) {
				if cursor.Kind() == clang.Cursor_FieldDecl {
					name := cursor.Spelling()
					if name == "" {
						log.Fatal("no field name")
					}

					ty := p.parseType(fmt.Sprintf("%v[%v]", indent, name), cursor.Type())

					u.Fields = append(u.Fields, &Field{
						Name: name,
						Type: ty,
					})
				}

				return clang.ChildVisit_Continue
			})

			return u

		default:
			log.Panicln("Unknown dec", "kind", dec.Kind())

			return nil
		}

	case clang.Type_FunctionProto:
		log.Println(indent, "Parsing type", t.Spelling(), "as funcproto")

		var args []Type

		for i := int32(0); i < t.NumArgTypes(); i++ {
			arg := t.ArgType(uint32(i))
			parsed := p.parseType(indent, arg)
			args = append(args, parsed)
		}

		result := p.parseType(indent, t.ResultType())

		return &FuncType{
			Args:   args,
			Result: result,
		}

	case clang.Type_ConstantArray:
		log.Println(indent, "Parsing type", t.Spelling(), "as const array")
		inner := p.parseType(indent, t.ArrayElementType())

		return &ConstArray{
			Inner: inner,
			Size:  t.ArraySize(),
		}

	case clang.Type_IncompleteArray:
		log.Println(indent, "Parsing type", t.Spelling(), "as array")
		inner := p.parseType(indent, t.ArrayElementType())

		return &Array{
			Inner: inner,
		}

	default:
		log.Println(indent, "Parsing type", t.Spelling(), "as ???")
		log.Panicln(indent, "Unknown type", t.Kind().Spelling())
		return nil
	}
}

// generateUnnamedStructName creates a synthetic name for an unnamed struct/union
func (p *Parser) generateUnnamedStructName(t clang.Type) string {
	// Get the canonical spelling which contains location info
	spelling := t.CanonicalType().Spelling()

	// Fallback: try to parse location from the spelling string
	// Format: "struct (unnamed at .../avformat.h:1022:5)"
	if strings.Contains(spelling, "unnamed at ") {
		parts := strings.Split(spelling, "unnamed at ")
		if len(parts) == 2 {
			locStr := strings.TrimSuffix(parts[1], ")")
			locParts := strings.Split(locStr, ":")
			if len(locParts) >= 3 {
				filename := filepath.Base(locParts[0])
				filename = strings.TrimSuffix(filename, filepath.Ext(filename))
				filename = cleanIdentifier(filename)
				line := locParts[len(locParts)-2]
				col := locParts[len(locParts)-1]

				return fmt.Sprintf("UnnamedStruct_%s_%s_%s", filename, line, col)
			}
		}
	}

	// Last resort fallback
	return fmt.Sprintf("UnnamedStruct_%d", len(p.mod.structs))
}

// cleanIdentifier converts a string into a valid Go identifier
func cleanIdentifier(s string) string {
	// Replace common special characters with underscores
	s = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return '_'
	}, s)

	// Remove consecutive underscores
	for strings.Contains(s, "__") {
		s = strings.ReplaceAll(s, "__", "_")
	}

	// Trim underscores from start/end
	s = strings.Trim(s, "_")

	return s
}

// registerUnnamedStruct creates a struct definition for an unnamed struct
func (p *Parser) registerUnnamedStruct(indent string, t clang.Type, syntheticName string) {
	// Check if we've already registered this struct
	if _, exists := p.mod.structs[syntheticName]; exists {
		log.Println(indent, "Unnamed struct already registered:", syntheticName)
		return
	}

	log.Println(indent, "Registering unnamed struct:", syntheticName)

	// Get the declaration to extract fields
	decl := t.Declaration()
	if decl.Spelling() == "" {
		log.Println(indent, "Warning: unnamed struct has no valid declaration")
		return
	}

	s := &Struct{
		Name:     syntheticName,
		Typedefd: false,
		Comment:  "Synthetic type for unnamed struct",
	}

	// Visit fields of the struct
	decl.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		if cursor.Kind() == clang.Cursor_FieldDecl {
			fieldName := cursor.Spelling()
			if fieldName == "" {
				log.Println(indent, "Warning: unnamed struct has field with no name")
				return clang.ChildVisit_Continue
			}

			comment := processComment(cursor.RawCommentText())
			fIndent := fmt.Sprintf("%v[%v]", indent, fieldName)
			fieldType := p.parseType(fIndent, cursor.Type())

			s.Fields = append(s.Fields, &Field{
				Name:     fieldName,
				Type:     fieldType,
				BitWidth: cursor.FieldDeclBitWidth(),
				Comment:  comment,
			})
		}

		return clang.ChildVisit_Continue
	})

	// Register the struct
	p.mod.structs[syntheticName] = s
	p.mod.structOrder = append(p.mod.structOrder, syntheticName)

	log.Println(indent, "Successfully registered unnamed struct:", syntheticName, "with", len(s.Fields), "fields")
}
