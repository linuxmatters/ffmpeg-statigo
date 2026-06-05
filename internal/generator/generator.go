package main

import "C"
import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dave/jennifer/jen"
)

// saveFormatted writes a jennifer file atomically: it saves to a temporary
// sibling path then renames over the destination, so a crash mid-write never
// leaves a partial `*.gen.go` behind. Formatting is delegated to jennifer's
// Save so the gofmt output is byte-identical to a direct o.Save(path).
func saveFormatted(o *jen.File, path string) error {
	tmp := path + ".tmp"
	if err := o.Save(tmp); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, path)
}

var primTypes = map[string]string{
	"int":       "int",
	"uint":      "uint",
	"char":      "uint8",
	"uchar":     "uint8",
	"ulong":     "uint64",
	"int8_t":    "int8",
	"int16_t":   "int16",
	"int32_t":   "int32",
	"int64_t":   "int64",
	"uint8_t":   "uint8",
	"uint16_t":  "uint16",
	"uint32_t":  "uint32",
	"uint64_t":  "uint64",
	"size_t":    "uint64",
	"ptrdiff_t": "int64",
	"float":     "float32",
	"double":    "float64",
}

// getCType returns the correct C type name to use in C.xxx() conversions.
// This handles cases where libclang might report the wrong type.
// fieldCTypeOverrides provides explicit C type names for fields that libclang misreports.
// This is needed when libclang can't find system headers and resolves typedefs incorrectly.
var fieldCTypeOverrides = map[string]map[string]string{
	"AVPixFmtDescriptor": {
		"nb_components": "uint8_t",
		"log2_chroma_w": "uint8_t",
		"log2_chroma_h": "uint8_t",
		"flags":         "uint64_t",
	},
}

// skippedStructs contains struct names that should not be generated.
// These require manual bindings in streamgroup.go, typically because they are
// anonymous C structs that CGO cannot directly reference by name.
var skippedStructs = map[string]bool{
	// Anonymous struct inside AVStreamGroupTileGrid - CGO generates an internal
	// name like "struct___12" which is unstable across compilations.
	"UnnamedStruct_avformat_986_5": true,
}

// skippedFields contains struct.field combinations that should not be generated.
// These require manual bindings in streamgroup.go, typically because they reference
// types that CGO cannot directly access (like anonymous structs).
var skippedFields = map[string]map[string]bool{
	"AVStreamGroupTileGrid": {
		// The offsets field is a pointer to an anonymous struct that CGO
		// assigns an internal name. Manual binding in streamgroup.go handles this.
		"offsets": true,
	},
}

// manuallyWrappedFields lists struct.field combinations the generator cannot
// emit but that have hand-written accessors in fields.go. They stay recorded as
// skips (so the ceiling count is unchanged) but carry a "manually wrapped in
// fields.go" reason so the skip summary stops counting them as coverage gaps.
// Mirrors the func.go "manually wrapped in uuid.go" convention.
var manuallyWrappedFields = map[string]map[string]bool{
	"AVCodecContext": {
		"intra_matrix":        true,
		"inter_matrix":        true,
		"chroma_intra_matrix": true,
	},
	"AVFrame": {
		"extended_data": true,
	},
	"AVPixFmtDescriptor": {
		"comp": true,
	},
	"AVMasteringDisplayMetadata": {
		"display_primaries": true,
	},
	"AVPanScan": {
		"position": true,
	},
}

// outputParam describes how marshalPointerArg and marshalArg treat an
// output-pointer parameter. Membership in outputParams marks the parameter as
// an output pointer that marshalPointerArg must emit as
// (*C.<type>)(unsafe.Pointer(p)). sizeT marks the size_t-width subset that
// libclang misreports as int but whose FFmpeg headers declare size_t;
// marshalArg rewrites those to size_t.
type outputParam struct {
	sizeT bool
}

// outputParams enumerates the (function, parameter) pairs whose
// primitive-pointer parameter is an output pointer. Mirrors the shape of
// skippedFields. Seeded against the pinned FFmpeg headers. Entries with
// sizeT: true are the size_t-width subset that libclang misreports as int
// (replacing the former substring heuristic).
// Sorted by function name then parameter name for reviewability.
var outputParams = map[string]map[string]outputParam{
	"av_ambient_viewing_environment_alloc": {
		"size": {sizeT: true},
	},
	"av_buffersink_get_side_data": {
		"nb_side_data": {},
	},
	"av_content_light_metadata_alloc": {
		"size": {sizeT: true},
	},
	"av_cpb_properties_alloc": {
		"size": {sizeT: true},
	},
	"av_detection_bbox_alloc": {
		"out_size": {sizeT: true},
	},
	"av_dovi_alloc": {
		"size": {sizeT: true},
	},
	"av_dovi_metadata_alloc": {
		"size": {sizeT: true},
	},
	"av_dynamic_hdr_plus_alloc": {
		"size": {sizeT: true},
	},
	"av_dynamic_hdr_vivid_alloc": {
		"size": {sizeT: true},
	},
	"av_dynarray_add": {
		"nb_ptr": {},
	},
	"av_dynarray_add_nofree": {
		"nb_ptr": {},
	},
	"av_encryption_info_add_side_data": {
		"side_data_size": {},
	},
	"av_encryption_init_info_add_side_data": {
		"side_data_size": {},
	},
	"av_expr_count_func": {
		"counter": {},
	},
	"av_expr_count_vars": {
		"counter": {},
	},
	"av_expr_parse_and_eval": {
		"res": {},
	},
	"av_fast_malloc": {
		"size": {},
	},
	"av_fast_mallocz": {
		"size": {},
	},
	"av_fast_padded_malloc": {
		"size": {},
	},
	"av_fast_padded_mallocz": {
		"size": {},
	},
	"av_fast_realloc": {
		"size": {},
	},
	"av_film_grain_params_alloc": {
		"size": {sizeT: true},
	},
	"av_find_best_pix_fmt_of_2": {
		"loss_ptr": {},
	},
	"av_get_output_timestamp": {
		"dts":  {},
		"wall": {},
	},
	"av_iamf_param_definition_alloc": {
		"size": {sizeT: true},
	},
	"av_lzo1x_decode": {
		"outlen": {},
	},
	"av_mastering_display_metadata_alloc_size": {
		"size": {sizeT: true},
	},
	"av_opt_eval_double": {
		"double_out": {},
	},
	"av_opt_eval_flags": {
		"flags_out": {},
	},
	"av_opt_eval_float": {
		"float_out": {},
	},
	"av_opt_eval_int": {
		"int_out": {},
	},
	"av_opt_eval_int64": {
		"int64_out": {},
	},
	"av_opt_eval_uint": {
		"uint_out": {},
	},
	"av_opt_get_array_size": {
		"out_val": {},
	},
	"av_opt_get_double": {
		"out_val": {},
	},
	"av_opt_get_image_size": {
		"h_out": {},
		"w_out": {},
	},
	"av_opt_get_int": {
		"out_val": {},
	},
	"av_packet_get_side_data": {
		"size": {},
	},
	"av_packet_pack_dictionary": {
		"size": {},
	},
	"av_packet_side_data_add": {
		"nb_sd": {},
	},
	"av_packet_side_data_free": {
		"nb_sd": {},
	},
	"av_packet_side_data_from_frame": {
		"nb_sd": {},
	},
	"av_packet_side_data_new": {
		"pnb_sd": {},
	},
	"av_packet_side_data_remove": {
		"nb_sd": {},
	},
	"av_parse_time": {
		"timeval": {},
	},
	"av_parse_video_size": {
		"height_ptr": {},
		"width_ptr":  {},
	},
	"av_pix_fmt_get_chroma_sub_sample": {
		"h_shift": {},
		"v_shift": {},
	},
	"av_probe_input_format3": {
		"score_ret": {},
	},
	"av_reduce": {
		"dst_den": {},
		"dst_num": {},
	},
	"av_samples_get_buffer_size": {
		"linesize": {},
	},
	"av_spherical_alloc": {
		"size": {sizeT: true},
	},
	"av_stereo3d_alloc_size": {
		"size": {sizeT: true},
	},
	"av_tdrdi_alloc": {
		"size": {sizeT: true},
	},
	"av_url_split": {
		"port_ptr": {},
	},
	"av_video_enc_params_alloc": {
		"out_size": {sizeT: true},
	},
	"av_video_hint_alloc": {
		"out_size": {sizeT: true},
	},
	"avcodec_align_dimensions": {
		"height": {},
		"width":  {},
	},
	"avcodec_align_dimensions2": {
		"height": {},
		"width":  {},
	},
	"avcodec_decode_subtitle2": {
		"got_sub_ptr": {},
	},
	"avcodec_find_best_pix_fmt_of_list": {
		"loss_ptr": {},
	},
	"swr_set_channel_mapping": {
		"channel_map": {},
	},
}

// getCType returns the C type name to emit after "C." in a generated cast.
//
// The special cases below exist because libclang and CGO disagree about type
// identity. When libclang cannot fully resolve a system-header typedef it
// canonicalises or misreports the underlying type, collapsing distinct named
// types onto a same-width primitive. CGO then rejects the result: it treats
// same-width named C types as distinct (C.char vs C.uint8_t, C.ptrdiff_t vs
// C.int64_t, C.ulong/C.long vs C.uint64_t/C.int64_t), so a cast spelled with
// the wrong-but-same-width name fails to compile. Each branch preserves the
// exact C spelling the cast needs. This is essential complexity compensating
// for libclang/CGO type-distinctness, not accidental cruft; the per-branch
// notes record the specific type pair each guards.
func getCType(typeName string, goType string) string {
	// Special case: char should stay as char, not become uint8_t
	// This is important for function parameters where char != uint8_t
	if typeName == "char" {
		return "char"
	}

	// Special case: ptrdiff_t and size_t must be preserved
	// On macOS, ptrdiff_t is a distinct type from int64_t even though they're the same size
	// Using C.int64_t instead of C.ptrdiff_t causes type mismatch errors
	if typeName == "ptrdiff_t" {
		return "ptrdiff_t"
	}

	// Special case: ulong and long must be preserved as their named C types.
	// CGO treats C.ulong/C.long as distinct from C.uint64_t/C.int64_t even at
	// the same width, so map them before the goType switch to avoid mismatches.
	if typeName == "ulong" {
		return "ulong"
	}
	if typeName == "long" {
		return "long"
	}

	// Map Go types to their correct C types
	// Returns the type name to use after "C." in generated code
	switch goType {
	case "int":
		return "int"
	case "uint":
		return "uint" // CGO accepts C.uint
	case "int8":
		return "int8_t"
	case "int16":
		return "int16_t"
	case "int32":
		return "int32_t"
	case "int64":
		return "int64_t"
	case "uint8":
		return "uint8_t"
	case "uint16":
		return "uint16_t"
	case "uint32":
		return "uint32_t"
	case "uint64":
		// Could be uint64_t or size_t - prefer size_t for size-related fields
		if typeName == "size_t" {
			return "size_t"
		}
		return "uint64_t"
	case "float32":
		return "float"
	case "float64":
		return "double"
	}

	// For other types, map Go pseudo-types to CGO types
	// Special case: uchar kept as its named C type; CGO treats C.uchar as
	// distinct from C.uint8_t, so preserve the spelling libclang reported.
	if typeName == "uchar" {
		return "uchar"
	}

	// Default: use the reported type name
	return typeName
}

// SkipEntry pairs a skipped C symbol with the short reason that explains why
// the generator could not emit a binding for it. The reason text mirrors the
// `skipped due to ...` comment emitted into the generated file at the same
// site, so a summary line and the generated comment always agree.
type SkipEntry struct {
	Symbol string
	Reason string
	// Manual names the hand-written Go binding that already covers this skipped
	// symbol, or is empty when none was found. It distinguishes a skip that is
	// genuinely missing a binding from one a topic file in the repo root already
	// wraps. Populated by enrichManualBindings after a run records every skip.
	Manual string
}

// SkipCollector aggregates every `skipped due to ...` decision the generator
// makes during a single run. It observes the existing emit sites; it does not
// change emitted code. enforceSkipCeiling consumes the total to enforce a
// regression ceiling on FFmpeg upgrades.
type SkipCollector struct {
	Entries []SkipEntry
}

// Record appends a skip decision. A nil receiver is a no-op so call sites can
// stay unconditional even when no collector is wired (e.g. focused tests).
func (c *SkipCollector) Record(symbol, reason string) {
	if c == nil {
		return
	}
	c.Entries = append(c.Entries, SkipEntry{Symbol: symbol, Reason: reason})
}

// Total returns the number of recorded skip markers. Matches the count of
// `skipped due to ...` comments in the generated `*.gen.go` output.
func (c *SkipCollector) Total() int {
	if c == nil {
		return 0
	}
	return len(c.Entries)
}

type Generator struct {
	input *Module
	skips *SkipCollector
}

// Gen runs the codegen pass. The skips collector is shared with Parse so a
// single SkipCollector accumulates both parse-layer (Task 3.2) and codegen-layer
// skips. A nil skips is tolerated: a fresh collector is allocated. The returned
// collector is the active one (the argument when non-nil, otherwise the freshly
// allocated one) so callers can read aggregated state without tracking it
// separately.
func Gen(i *Module, skips *SkipCollector) *SkipCollector {
	if skips == nil {
		skips = &SkipCollector{}
	}
	g := &Generator{
		input: i,
		skips: skips,
	}

	g.generateConstants()
	g.generateEnums()
	g.generateCallbacks()
	g.generateStructs()
	g.generateFuncs()

	g.enrichManualBindings(".")

	return g.skips
}

// enrichManualBindings annotates each recorded skip with the hand-written Go
// binding that already covers it, when one exists. dir is the generator's
// working directory, where the hand-written `package ffmpeg` files sit beside
// the freshly written `*.gen.go` output.
//
// A skip symbol takes one of two shapes (see func.go and struct.go):
//   - a function skip records the C name, e.g. "av_display_rotation_get"; the
//     expected binding is convCamel(symbol), matched against top-level funcs.
//   - a struct-field skip records "GoStructName.c_field_name", e.g.
//     "AVCodecContext.intra_matrix"; the expected accessor is
//     convCamel(fieldPart) on that receiver, matched against the type's methods.
//
// convCamel depends on g.input.structs, so enrichment runs as a Generator method
// after Gen has populated the module. A scan failure leaves every Manual field
// empty, degrading to the prior behaviour rather than aborting the run.
func (g *Generator) enrichManualBindings(dir string) {
	mb, err := scanManualBindings(dir)
	if err != nil {
		return
	}

	for i := range g.skips.Entries {
		e := &g.skips.Entries[i]

		if recv, field, ok := strings.Cut(e.Symbol, "."); ok {
			method := g.convCamel(field)
			if mb.hasMethod(recv, method) {
				e.Manual = method
			}
			continue
		}

		fn := g.convCamel(e.Symbol)
		if mb.hasFunc(fn) {
			e.Manual = fn
		}
	}
}

func newFile() *jen.File {
	o := jen.NewFile("ffmpeg")

	o.HeaderComment("Code generated by the ffmpeg-statigo generator. DO NOT EDIT.")

	for _, file := range files {
		o.CgoPreamble(fmt.Sprintf("#include <%v>", file))
	}

	return o
}

func (g *Generator) generateConstants() {
	i := g.input

	o := newFile()

	for _, constName := range i.constantOrder {
		constant := i.constants[constName]

		// Skip constants that conflict with Go's math package
		if constName == "NAN" || constName == "INFINITY" {
			log.Println("Skipping constant", constant.Name, "(conflicts with Go math package)")
			continue
		}

		log.Println("Generating constant", constant.Name)

		goName := g.convCamel(constant.Name)

		if strings.HasPrefix(constName, "AVERROR_") {
			goName = fmt.Sprintf("%vConst", goName)
		}

		o.Commentf("%v wraps %v.", goName, constant.Name)

		o.Const().Id(goName).Op("=").Qual("C", constName)
	}

	err := saveFormatted(o, "constants.gen.go")
	if err != nil {
		log.Panicln(err)
	}
}

func (g *Generator) generateEnums() {
	i := g.input

	o := newFile()

	for _, enumName := range i.enumOrder {
		enum := i.enums[enumName]

		log.Println("Generating enum", enum.Name)
		o.Commentf("--- Enum %v ---", enum.Name)
		o.Line()

		goName := enumName

		o.Commentf("%v wraps %v.", goName, enum.Name)

		if enum.Comment != "" {
			o.Comment(enum.Comment)
		}

		cName := enum.CName()

		o.Type().Id(goName).Qual("C", cName)

		o.Const().Id(fmt.Sprintf("SizeOf%v", goName)).Op("=").Qual("C", fmt.Sprintf("sizeof_%v", cName))

		o.Func().
			Id(fmt.Sprintf("To%vArray", goName)).
			Params(jen.Id("ptr").Qual("unsafe", "Pointer")).
			Op("*").Id("Array").Types(jen.Id(goName)).
			Block(
				jen.If(jen.Id("ptr").Op("==").Id("nil")).Block(
					jen.Return(jen.Id("nil")),
				),
				jen.Line(),
				jen.Return(
					jen.Op("&").Id("Array").Types(jen.Id(goName)).Values(jen.Dict{
						jen.Id("ptr"):      jen.Id("ptr"),
						jen.Id("elemSize"): jen.Id(fmt.Sprintf("SizeOf%v", goName)),

						jen.Id("loadPtr"): jen.Func().
							Params(jen.Id("pointer").Qual("unsafe", "Pointer")).
							Id(goName).
							Block(
								jen.Id("ptr").Op(":=").Parens(jen.Op("*").Id(goName)).Parens(jen.Id("pointer")),
								jen.Return(jen.Op("*").Id("ptr")),
							),

						jen.Id("storePtr"): jen.Func().
							Params(
								jen.Id("pointer").Qual("unsafe", "Pointer"),
								jen.Id("value").Id(goName),
							).
							Block(
								jen.Id("ptr").Op(":=").Parens(jen.Op("*").Id(goName)).Parens(jen.Id("pointer")),
								jen.Op("*").Id("ptr").Op("=").Id("value"),
							),
					}),
				),
			)

		o.Line()

		o.Func().
			Id(fmt.Sprintf("Alloc%vArray", goName)).
			Params(jen.Id("size").Id("uint64")).
			Op("*").Id("Array").Types(jen.Id(goName)).
			Block(
				jen.Return(jen.Id(fmt.Sprintf("To%vArray", goName)).Params(
					jen.Id("AVCalloc").Params(jen.Id("size"), jen.Id(fmt.Sprintf("SizeOf%v", goName))),
				)),
			)

		var valDefs []jen.Code

		for _, constant := range enum.Constants {
			constName := g.convCamel(constant.Name)

			valDefs = append(valDefs, jen.Commentf("%v wraps %v.", constName, constant.Name))

			if constant.Comment != "" {
				valDefs = append(valDefs, jen.Comment(constant.Comment))
			}

			valDefs = append(valDefs, jen.Id(constName).Id(goName).Op("=").Qual("C", constant.Name))
		}

		if len(valDefs) > 0 {
			o.Const().Defs(valDefs...)
		}
	}

	err := saveFormatted(o, "enums.gen.go")
	if err != nil {
		log.Panicln(err)
	}
}

func (g *Generator) generateCallbacks() {
	i := g.input

	o := newFile()

	for _, callbackName := range i.callbackOrder {
		callback := i.callbacks[callbackName]

		log.Println("Generating callback", callback.Name)

		goName := g.convCamel(callback.Name)

		// Generate as type alias to unsafe.Pointer since we can't represent C function pointers directly
		o.Commentf("%v is a function pointer typedef for %v.", goName, callback.Name)
		if callback.Comment != "" {
			o.Comment(callback.Comment)
		}
		o.Type().Id(goName).Qual("unsafe", "Pointer")
		o.Line()
	}

	err := saveFormatted(o, "callbacks.gen.go")
	if err != nil {
		log.Panicln(err)
	}
}
