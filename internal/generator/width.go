package main

import "github.com/dave/jennifer/jen"

func variableWidthPrimitive(cType string) bool {
	switch cType {
	case "size_t", "ptrdiff_t", "ulong", "long":
		return true
	default:
		return false
	}
}

// A round trip through the build target's C type detects narrowing without
// assuming that pointer width and long width are equal (Windows uses LLP64).
func primitiveWidthCheck(cType, goType string, value jen.Code, name string) []jen.Code {
	if !variableWidthPrimitive(cType) {
		return nil
	}
	return []jen.Code{
		jen.If(jen.Id(goType).Params(jen.Qual("C", cType).Params(value)).Op("!=").Add(value)).Block(
			jen.Panic(jen.Lit(name + " overflows C." + cType)),
		),
	}
}
