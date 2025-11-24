package main

import (
	"fmt"
	"log"
	"path"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"unicode"

	"github.com/Newbluecake/bootstrap/clang"
)

var (
	AVLibPath, _ = filepath.Abs("include")
)

var files = []string{
	"libavcodec/avcodec.h",
	"libavcodec/bsf.h",
	"libavcodec/codec.h",
	"libavcodec/codec_desc.h",
	"libavcodec/codec_id.h",
	"libavcodec/codec_par.h",
	"libavcodec/defs.h",
	"libavcodec/packet.h",
	"libavcodec/version.h",
	"libavcodec/version_major.h",
	"libavdevice/avdevice.h",
	"libavdevice/version.h",
	"libavdevice/version_major.h",
	"libavfilter/avfilter.h",
	"libavfilter/buffersink.h",
	"libavfilter/buffersrc.h",
	"libavfilter/version.h",
	"libavfilter/version_major.h",
	"libavformat/avformat.h",
	"libavformat/avio.h",
	"libavformat/version.h",
	"libavformat/version_major.h",
	"libavutil/adler32.h",
	"libavutil/aes.h",
	"libavutil/aes_ctr.h",
	"libavutil/ambient_viewing_environment.h",
	////"libavutil/attributes.h", // a compiler attribute macro collection, not an API header
	"libavutil/audio_fifo.h",
	"libavutil/avassert.h",
	"libavutil/avconfig.h",
	"libavutil/avstring.h",
	"libavutil/avutil.h",
	"libavutil/base64.h",
	"libavutil/blowfish.h",
	"libavutil/bprint.h",
	"libavutil/bswap.h",
	"libavutil/buffer.h",
	"libavutil/camellia.h",
	"libavutil/cast5.h",
	"libavutil/channel_layout.h",
	"libavutil/common.h",
	"libavutil/container_fifo.h",
	"libavutil/cpu.h",
	"libavutil/crc.h",
	"libavutil/csp.h",
	"libavutil/des.h",
	"libavutil/detection_bbox.h",
	"libavutil/dict.h",
	"libavutil/display.h",
	"libavutil/dovi_meta.h",
	"libavutil/downmix_info.h",
	"libavutil/encryption_info.h",
	"libavutil/error.h",
	"libavutil/eval.h",
	"libavutil/executor.h",
	"libavutil/ffversion.h",
	"libavutil/fifo.h",
	"libavutil/file.h",
	"libavutil/film_grain_params.h",
	"libavutil/frame.h",
	"libavutil/hash.h",
	"libavutil/hdr_dynamic_metadata.h",
	"libavutil/hdr_dynamic_vivid_metadata.h",
	"libavutil/hmac.h",
	"libavutil/hwcontext.h",
	////"libavutil/hwcontext_cuda.h",
	////"libavutil/hwcontext_qsv.h",
	////"libavutil/hwcontext_videotoolbox.h",
	////"libavutil/hwcontext_vulkan.h",
	"libavutil/iamf.h",
	"libavutil/imgutils.h",
	"libavutil/intfloat.h",
	////"libavutil/intreadwrite.h", //Union types - CGO doesn't expose union fields
	"libavutil/lfg.h",
	"libavutil/log.h",
	"libavutil/lzo.h",
	"libavutil/macros.h",
	"libavutil/mastering_display_metadata.h",
	"libavutil/mathematics.h",
	"libavutil/mem.h",
	"libavutil/motion_vector.h",
	"libavutil/murmur3.h",
	"libavutil/opt.h",
	"libavutil/parseutils.h",
	"libavutil/pixdesc.h",
	"libavutil/pixelutils.h",
	"libavutil/pixfmt.h",
	"libavutil/random_seed.h",
	"libavutil/rational.h",
	"libavutil/rc4.h",
	////"libavutil/refstruct.h", // reference-counted object API introduced as an alternative to AVBuffer for managing complex objects
	"libavutil/replaygain.h",
	"libavutil/ripemd.h",
	"libavutil/samplefmt.h",
	"libavutil/sha.h",
	"libavutil/sha512.h",
	"libavutil/spherical.h",
	"libavutil/stereo3d.h",
	"libavutil/tdrdi.h",
	"libavutil/tea.h",
	"libavutil/threadmessage.h",
	"libavutil/time.h",
	"libavutil/timecode.h",
	"libavutil/timestamp.h",
	"libavutil/tree.h",
	"libavutil/twofish.h",
	"libavutil/tx.h",
	"libavutil/uuid.h",
	"libavutil/version.h",
	"libavutil/video_enc_params.h",
	"libavutil/video_hint.h",
	"libavutil/xtea.h",
	"libswresample/version.h",
	"libswresample/version_major.h",
	"libswresample/swresample.h",
	"libswscale/version.h",
	"libswscale/version_major.h",
	"libswscale/swscale.h",
}

func Parse() *Module {
	p := &Parser{
		mod: &Module{
			functions: make(map[string]*Function),
			structs:   make(map[string]*Struct),
			enums:     make(map[string]*Enum),
			callbacks: make(map[string]*Function),
			constants: make(map[string]*Constant),
		},
	}

	for _, file := range files {
		filePath := path.Join(AVLibPath, file)

		fmt.Println(filePath)

		p.parseFile(fmt.Sprintf("[%v]", file), filePath)
	}

	return p.mod
}

type Parser struct {
	path string
	mod  *Module
	tu   clang.TranslationUnit
}

func getPlatformArgs() []string {
	args := []string{
		"-fparse-all-comments",
		fmt.Sprintf("-I%v", AVLibPath),
		"-D_GNU_SOURCE",
		"-D_DEFAULT_SOURCE",
		"-D__STDC_CONSTANT_MACROS",
	}

	// Add platform-specific includes
	switch runtime.GOOS {
	case "darwin":
		args = append(args, "-I/Library/Developer/CommandLineTools/SDKs/MacOSX.sdk/System/Library/Frameworks/Kernel.framework/Headers/")
	case "linux":
		// For standard Linux distributions, add common include paths
		// These will be ignored if they don't exist, so it's safe to add them
		// NixOS will use its own paths through environment variables
		if runtime.GOARCH == "amd64" {
			args = append(args, "-I/usr/include/x86_64-linux-gnu")
		} else if runtime.GOARCH == "arm64" {
			args = append(args, "-I/usr/include/aarch64-linux-gnu")
		}
		args = append(args, "-I/usr/include")
	}

	return args
}

func (p *Parser) parseFile(indent string, path string) {
	p.path = path

	idx := clang.NewIndex(0, 1)
	defer idx.Dispose()

	tu := idx.ParseTranslationUnit(
		path,
		getPlatformArgs(),
		nil,
		clang.TranslationUnit_IncludeBriefCommentsInCodeCompletion|clang.TranslationUnit_DetailedPreprocessingRecord,
	)
	defer tu.Dispose()

	diagnostics := tu.Diagnostics()
	for _, d := range diagnostics {
		fmt.Println("PROBLEM:", d.Spelling())
	}

	p.tu = tu

	tu.
		TranslationUnitCursor().
		Visit(func(cursor, parent clang.Cursor) (status clang.ChildVisitResult) {
			p.parseTopLevel(indent, cursor)

			return clang.ChildVisit_Continue
		})
}

func (p *Parser) parseTopLevel(indent string, c clang.Cursor) {
	loc, _, _, _ := c.Location().FileLocation()
	if loc.Name() != p.path {
		return
	}

	switch c.Kind() {
	case clang.Cursor_TypedefDecl:
		p.parseTypedef(indent, c)

	case clang.Cursor_VarDecl:
		log.Println("vardecl", "name", c.Spelling())
		log.Println(" ", c.Type().Spelling())

	case clang.Cursor_FunctionDecl:
		p.parseFunction(indent, c)

	case clang.Cursor_MacroDefinition:
		if c.IsMacroFunctionLike() {
			return
		}

		name := c.Spelling()

		log.Println("macro def", "name", name)

		if strings.HasSuffix(name, "_H") {
			return
		}

		if strings.HasPrefix(name, "av_malloc_") {
			return
		}

		if strings.HasPrefix(name, "AV_CHANNEL_LAYOUT_") {
			return
		}

		if strings.HasSuffix(name, "_VERSION") || strings.HasSuffix(name, "_IDENT") {
			return
		}

		if name == "AV_TIME_BASE_Q" || name == "AV_CH_LAYOUT_NATIVE" {
			return
		}

		// Skip compiler attribute macros (from libavutil/attributes.h and transitively included headers)
		if strings.HasPrefix(name, "av_") && (strings.HasSuffix(name, "_inline") ||
			strings.Contains(name, "_unused") ||
			strings.Contains(name, "_deprecated") ||
			strings.Contains(name, "_result") ||
			name == "av_const" || name == "av_pure" ||
			name == "av_cold" || name == "av_flatten" ||
			name == "av_alias" || name == "av_noreturn" ||
			name == "av_used" || name == "av_noinline") {
			return
		}

		// Skip math/utility function macros (from libavutil/common.h and transitively included headers)
		if strings.HasPrefix(name, "av_clip") ||
			strings.HasPrefix(name, "av_sat_") ||
			strings.HasPrefix(name, "av_popcount") ||
			strings.HasPrefix(name, "av_ceil_log2") ||
			strings.HasPrefix(name, "av_mod_") ||
			strings.HasPrefix(name, "av_zero_extend") ||
			strings.HasPrefix(name, "av_parity") ||
			strings.HasPrefix(name, "FF_CEIL_RSHIFT") ||
			strings.HasPrefix(name, "av_builtin_") {
			return
		}

		p.mod.constants[name] = &Constant{
			Name: name,
		}
		p.mod.constantOrder = append(p.mod.constantOrder, name)

	case clang.Cursor_EnumDecl:
		p.parseEnum(indent, c, false)

	case clang.Cursor_StructDecl:
		p.parseStruct(indent, c, false)

	case clang.Cursor_MacroExpansion, clang.Cursor_InclusionDirective, clang.Cursor_UnionDecl:
		// ignore

	default:
		log.Panicln("Unexpected top level", "kind", c.Kind().String())
	}
}

func (p *Parser) parseTypedef(indent string, c clang.Cursor) {
	name := c.Spelling()
	indent = fmt.Sprintf("%v[%v]", indent, name)

	log.Println("typedef", "name", c.Spelling())

	var params []*Param
	paramIndex := 0

	c.Visit(func(cursor, parent clang.Cursor) (status clang.ChildVisitResult) {
		log.Println("  --- ", "kind", cursor.Kind().String(), "name", cursor.Spelling())

		if cursor.Kind() == clang.Cursor_ParmDecl {
			name := cursor.Spelling()
			if name == "" {
				log.Println(indent, "no param name, generating one")
				name = fmt.Sprintf("param%v", paramIndex)
			}
			paramIndex++

			ty := p.parseType(fmt.Sprintf("%v[%v]", indent, name), cursor.Type())

			params = append(params, &Param{
				Name: name,
				Type: ty,
			})
		}

		return clang.ChildVisit_Continue
	})

	log.Println("dk ", c.Definition().Kind())
	log.Println("ds ", c.Definition().Spelling())
	log.Println("s ", c.TypedefDeclUnderlyingType().Spelling())
	log.Println("k ", c.TypedefDeclUnderlyingType().Kind())
	log.Println("cs ", c.TypedefDeclUnderlyingType().CanonicalType().Spelling())
	log.Println("ck ", c.TypedefDeclUnderlyingType().CanonicalType().Kind())

	if len(params) > 0 {
		ut := c.TypedefDeclUnderlyingType()
		pt := ut.PointeeType()

		ptr := true

		if pt.Kind() == clang.Type_Invalid {
			pt = ut
			ptr = false
		}

		if pt.NumArgTypes() != int32(len(params)) {
			log.Panicln("arg mismatch", pt.NumArgTypes(), int32(len(params)))
		}

		result := p.parseType(fmt.Sprintf("%v[%v]", indent, name), pt.ResultType())

		if _, ok := p.mod.callbacks[name]; ok {
			log.Panicln("callback already exists")

			return
		}

		p.mod.callbacks[name] = &Function{
			Name:   name,
			Args:   params,
			Result: result,
			Ptr:    ptr,
		}
		p.mod.callbackOrder = append(p.mod.callbackOrder, name)

		return
	}

	dec := c.TypedefDeclUnderlyingType().Declaration()

	dec.Visit(func(cursor, parent clang.Cursor) (status clang.ChildVisitResult) {
		log.Println("  Inner ", "kind", cursor.Kind().String(), "name", cursor.Spelling())

		return clang.ChildVisit_Continue
	})

	switch dec.Kind() {
	case clang.Cursor_StructDecl:
		p.parseStruct(indent, dec, true)

	case clang.Cursor_UnionDecl:
		// Unions are similar to structs but all fields share the same memory location.
		// NOTE: CGO does not expose union fields directly - unions are treated as opaque types.
		// Field accessor generation will fail at compile time. Headers with only union-based
		// utilities (like intreadwrite.h) should remain excluded from the files list.
		// We parse the structure anyway to avoid panics and for potential future use.
		p.parseStruct(indent, dec, true)

	case clang.Cursor_EnumDecl:
		p.parseEnum(indent, dec, true)

	case clang.Cursor_TypedefDecl:
		// This is a typedef alias to another typedef (e.g., typedef AVCIExy AVWhitepointCoefficients)
		// We don't need to generate anything for this, as the underlying type already exists
		log.Println(indent, "Skipping typedef alias to", dec.Spelling())
		return

	case clang.Cursor_NoDeclFound:
		// This can happen for built-in types or types defined elsewhere
		log.Println(indent, "Skipping typedef with no declaration found")
		return

	default:
		log.Panicln("Unknown typedef", "kind", dec.Kind())
	}

}

func processComment(in string) string {
	txt := in
	txt = strings.ReplaceAll(txt, "\r\n", "\n")
	txt = strings.TrimSpace(txt)
	txt = strings.TrimPrefix(txt, "/**\n")
	txt = strings.TrimSuffix(txt, "*/")

	txt = strings.TrimRightFunc(txt, unicode.IsSpace)

	var rebuilt []string

	for _, s := range strings.Split(txt, "\n") {
		s = strings.TrimSpace(s)
		s = strings.TrimPrefix(s, "* ")
		s = strings.TrimPrefix(s, "/// ")
		s = strings.TrimPrefix(s, "///")
		s = strings.TrimPrefix(s, "// ")
		s = strings.TrimPrefix(s, "//")

		if strings.HasPrefix(s, "/**") {
			rebuilt = nil
			s = strings.TrimPrefix(s, "/** ")
			s = strings.TrimPrefix(s, "/**")

			if strings.TrimSpace(s) == "" {
				continue
			}
		}

		if strings.HasPrefix(s, "/*") {
			rebuilt = nil
			s = strings.TrimPrefix(s, "/* ")
			s = strings.TrimPrefix(s, "/*")

			if strings.TrimSpace(s) == "" {
				continue
			}
		}

		if strings.HasPrefix(s, "@defgroup") || strings.HasPrefix(s, "@ingroup") ||
			strings.HasPrefix(s, "@addtogroup") || strings.HasPrefix(s, "@}") {
			continue
		}

		if strings.HasPrefix(s, "@{") {
			rebuilt = nil
			continue
		}

		if s == "*" {
			s = ""
		}

		// Double space to enter verbatim mode
		rebuilt = append(rebuilt, fmt.Sprintf("  %v", s))
	}

	txt = strings.Join(rebuilt, "\n")
	txt = strings.TrimRightFunc(txt, unicode.IsSpace)

	if strings.Count(txt, "\n") == 0 {
		txt = strings.TrimLeftFunc(txt, unicode.IsSpace)
		txt = strings.TrimPrefix(txt, "<")
		txt = fmt.Sprintf("  %v", txt)
	}

	if strings.TrimSpace(txt) == "" {
		return ""
	}

	return txt
}

func (p *Parser) parseFunction(indent string, c clang.Cursor) {
	name := c.Spelling()

	if _, ok := p.mod.functions[name]; ok {
		log.Panicln("Function", name, "already exists")
	}

	log.Println("Parsing function", name)
	indent = fmt.Sprintf("%v[%v]", indent, name)

	comment := processComment(c.RawCommentText())

	result := p.parseType(fmt.Sprintf("%v[return]", indent), c.ResultType())

	var args []*Param

	for i := 0; i < int(c.NumArguments()); i++ {
		arg := c.Argument(uint32(i))

		if arg.Kind() != clang.Cursor_ParmDecl {
			log.Panicln(indent, "Argument not of parmdecl type", arg.Kind())
		}

		name := arg.Spelling()
		if name == "" {
			log.Println(indent, "no param name")

			name = fmt.Sprintf("param%v", i)
		}

		aIndent := fmt.Sprintf("%v[%v]", indent, name)

		ty := p.parseType(aIndent, arg.Type())

		args = append(args, &Param{
			Name: name,
			Type: ty,
		})
	}

	p.mod.functions[name] = &Function{
		Name:     name,
		Args:     args,
		Result:   result,
		Variadic: c.IsVariadic(),
		Comment:  comment,
	}
	p.mod.functionOrder = append(p.mod.functionOrder, name)
}

func (p *Parser) parseEnum(indent string, c clang.Cursor, typedef bool) {
	log.Println("enum", "name", c.Spelling())

	name := c.Spelling()
	indent = fmt.Sprintf("%v[%v]", indent, name)

	if strings.HasPrefix(name, "enum (unnamed") {
		log.Println(indent, "Treating unnamed enum as constants")

		c.Visit(func(cursor, parent clang.Cursor) (status clang.ChildVisitResult) {
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

		return
	}

	if strings.Contains(name, " ") {
		log.Panicln(indent, "Name contains spaces")
	}

	if val, ok := p.mod.enums[name]; ok && len(val.Constants) > 0 {
		log.Println(indent, "already exists")

		if typedef {
			val.Typedefd = true
		}

		return
	}

	comment := processComment(c.RawCommentText())

	enum := &Enum{
		Name:     name,
		Typedefd: typedef,
		Comment:  comment,
	}

	c.Visit(func(cursor, parent clang.Cursor) (status clang.ChildVisitResult) {
		if cursor.Kind() != clang.Cursor_EnumConstantDecl {
			log.Panicln("Unknown enum type", "kind", cursor.Kind().String())
		}

		comment := processComment(cursor.RawCommentText())

		enum.Constants = append(enum.Constants, &Constant{
			Name:    cursor.Spelling(),
			Comment: comment,
		})

		return clang.ChildVisit_Continue
	})

	p.mod.enumOrder = slices.DeleteFunc(p.mod.enumOrder, func(s string) bool {
		return s == name
	})

	p.mod.enums[name] = enum
	p.mod.enumOrder = append(p.mod.enumOrder, name)
}

func (p *Parser) parseStruct(indent string, c clang.Cursor, typedef bool) {
	log.Println("struct", "name", c.Spelling())

	name := c.Spelling()
	indent = fmt.Sprintf("%v[%v]", indent, name)

	if val, ok := p.mod.structs[name]; ok && len(val.Fields) > 0 {
		log.Println(indent, "already exists")

		if typedef {
			val.Typedefd = true
		}

		return
	}

	comment := processComment(c.RawCommentText())

	s := &Struct{
		Name:     name,
		Typedefd: typedef,
		Comment:  comment,
	}

	c.Visit(func(cursor, parent clang.Cursor) (status clang.ChildVisitResult) {
		if cursor.Kind() == clang.Cursor_FieldDecl {
			name := cursor.Spelling()
			if name == "" {
				log.Fatal("no field name")
			}

			cmt := processComment(cursor.RawCommentText())

			fIndent := fmt.Sprintf("%v[%v]", indent, name)

			// Capture the original C type name before it gets canonicalized
			cursorType := cursor.Type()
			cTypeName := strings.TrimSpace(cursorType.Spelling())
			// Remove "const " prefix if present
			cTypeName = strings.TrimPrefix(cTypeName, "const ")

			ty := p.parseType(fIndent, cursorType)

			s.Fields = append(s.Fields, &Field{
				Name:      name,
				Type:      ty,
				CTypeName: cTypeName,
				BitWidth:  cursor.FieldDeclBitWidth(),
				Comment:   cmt,
			})
		}

		return clang.ChildVisit_Continue
	})

	p.mod.structOrder = slices.DeleteFunc(p.mod.structOrder, func(s string) bool {
		return s == name
	})

	p.mod.structs[name] = s
	p.mod.structOrder = append(p.mod.structOrder, name)
}
